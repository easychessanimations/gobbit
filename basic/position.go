package basic

import (
	"fmt"
	"strings"
	"time"
)

const SEARCH_MAX_DEPTH = 100

const MAX_STATES = SEARCH_MAX_DEPTH + 1

type Position struct {
	States          [MAX_STATES]State
	StatePtr        int
	MaxStatePtr     int
	Nodes           int
	SearchStopped   bool
	NullMovePruning bool
}

func (pos *Position) Current() *State {
	return &pos.States[pos.StatePtr]
}

func (pos *Position) Init(variant Variant) {
	pos.StatePtr = 0
	pos.Current().Init(variant)
	pos.Current().Ply = 0
}

func (pos *Position) Reset(){
	pos.Init(pos.Current().Variant)
}

func (pos *Position) ParseFen(fen string){
	pos.Reset()
	pos.Current().ParseFen(fen)
}

func (pos Position) Line() string {
	sans := []string{}

	for ptr := 1; ptr <= pos.StatePtr; ptr++ {
		prevSt := pos.States[ptr-1]
		san := prevSt.MoveToSan(pos.States[ptr].Move)
		if prevSt.Turn == White {
			san = fmt.Sprintf("%d. %s", prevSt.FullmoveNumber, san)
		}
		sans = append(sans, san)
	}

	return strings.Join(sans, " ")
}

func (pos Position) PrettyPrintString() string {
	buff := pos.Current().PrettyPrintString()

	buff += fmt.Sprintf("\n\nLine ( %d ) : %s\n", pos.StatePtr, pos.Line())

	return buff
}

func (st *State) UpdateMaterialBalance() {
	st.Material[NoColor] = st.Material[White].Sub(st.Material[Black])
}

func (st *State) Put(p Piece, sq Square) {
	if p == NoPiece {
		return
	}

	st.Pieces[RankOf[sq]][FileOf[sq]] = p

	sqbb := sq.Bitboard()

	color := ColorOf[p]

	st.ByColor[color] |= sqbb

	st.ByFigure[FigureOf[p]] |= sqbb

	if p.IsLancer(){
		st.ByLancer |= sqbb
	}

	mat := PieceMaterialTables[p][sq]

	st.Material[color].Merge(mat)

	st.UpdateMaterialBalance()

	st.Zobrist ^= zobristPiece[p][sq]

	if FigureOf[p] == King {
		st.KingInfos[color] = KingInfo{
			IsCaptured: false,
			Square:     sq,
		}
	}
}

func (st *State) Remove(sq Square) {
	if st.Pieces[RankOf[sq]][FileOf[sq]] == NoPiece {
		return
	}

	p := st.PieceAtSquare(sq)

	st.Pieces[RankOf[sq]][FileOf[sq]] = NoPiece

	color := ColorOf[p]

	sqbb := sq.Bitboard()

	st.ByColor[color] &^= sqbb

	st.ByFigure[FigureOf[p]] &^= sqbb

	if p.IsLancer(){
		st.ByLancer &^= sqbb
	}

	mat := PieceMaterialTables[p][sq]

	st.Material[color].UnMerge(mat)

	st.UpdateMaterialBalance()

	st.Zobrist ^= zobristPiece[p][sq]

	if FigureOf[p] == King {
		st.KingInfos[ColorOf[p]] = KingInfo{
			IsCaptured: true,
			Square:     SquareA1,
		}
	}
}

const INFINITE_SCORE = Score(20000)
const MATE_SCORE = Score(10000)

func (pos Position) GameEnd(ply int) (bool, Score) {
	st := pos.Current()

	if st.KingInfos[st.Turn].IsCaptured {
		return true, -MATE_SCORE + Score(ply)
	}

	return false, 0
}

func (st *State) MakeMove(move Move) {
	p := st.PieceAtSquare(move.FromSq())

	pCol := ColorOf[p]

	top := st.PieceAtSquare(move.ToSq())

	st.Remove(move.FromSq())

	st.Remove(move.ToSq())

	st.HasDisabledMove = false

	if move.MoveType() == Promotion {
		st.Put(move.PromotionPiece(), move.ToSq())
	} else if move.MoveType() == SentryPush{
		st.Put(p, move.ToSq())

		st.Remove(move.PromotionSquare())
		st.Put(move.PromotionPiece(), move.PromotionSquare())

		st.DisableFromSquare = move.PromotionSquare()
		st.DisableToSquare = move.ToSq()
		st.HasDisabledMove = true
	} else if move.MoveType() == Castling{		
		side := CastlingSideKing		
		if FileOf[move.ToSq()] < FileOf[move.FromSq()]{
			side = CastlingSideQueen
		}
		cts := st.CastlingTargetSquares(pCol, side)
		st.Put(p, cts[0])
		st.Put(st.CastlingRights[pCol][side].RookOrigPiece, cts[1])
	} else {
		st.Put(p, move.ToSq())
	}

	st.Move = move

	st.Ply++

	st.SetSideToMove(st.Turn.Inverse())

	st.HalfmoveClock++

	oldEpSq := st.EpSquare

	st.SetEpSquare(SquareA1)

	newCastlingRights := st.CastlingRights

	if FigureOf[p] == King{
		// if king was moved, delete all castling rights
		newCastlingRights[pCol][CastlingSideKing].CanCastle = false
		newCastlingRights[pCol][CastlingSideQueen].CanCastle = false

		st.SetCastlingAbility(newCastlingRights)
	}else{
		// check if castling partner was moved or captured, if so, delete castling right on that side
		for side := CastlingSideKing; side <= CastlingSideQueen; side++{
			testSq := newCastlingRights[pCol][side].RookOrigSq

			if move.FromSq() == testSq || move.ToSq() == testSq{
				newCastlingRights[pCol][side].CanCastle = false

				st.SetCastlingAbility(newCastlingRights)
			}
		}
	}

	if FigureOf[p] == Pawn{
		st.HalfmoveClock = 0

		rankDiff := RankOf[move.ToSq()] - RankOf[move.FromSq()]

		if rankDiff > 1 || rankDiff < -1{
			// push by two
			pi := PawnInfos[move.FromSq()][ColorOf[p]]

			for _, checkEp := range pi.CheckEps{
				if st.PieceAtSquare(checkEp) == p.ColorInverse(){
					st.SetEpSquare(pi.PushOneSq)
					break
				}
			}
		}

		if move.ToSq() == oldEpSq{
			var dir Rank = 1
			if ColorOf[p] == White{
				dir = -1
			}
			epClSq := RankFile[RankOf[move.ToSq()]+dir][FileOf[move.ToSq()]]			
			st.Remove(epClSq)
		}
	}

	if top != NoPiece{
		st.HalfmoveClock = 0
	}

	if st.Turn == White {
		st.FullmoveNumber++
	}
}

func (pos *Position) Push(move Move) {
	oldState := pos.States[pos.StatePtr]

	pos.StatePtr++

	if pos.StatePtr > pos.MaxStatePtr || pos.Current().Move != move {
		pos.MaxStatePtr = pos.StatePtr
	}

	pos.States[pos.StatePtr] = oldState

	pos.Current().MakeMove(move)
}

func (pos *Position) PushUci(uci string){
	move, ok := pos.Current().UciToMove(uci)
	if ok{
		pos.Push(move)
	}
}

func (pos *Position) Pop() {
	pos.StatePtr--
}

func (pos *Position) PerfRec(remDepth int) {
	pos.Nodes++

	if remDepth == 0 {
		return
	}

	for _, move := range pos.Current().LegalMoves(false) {
		pos.Push(move)
		//fmt.Println(pos.PrettyPrintString())
		pos.PerfRec(remDepth - 1)
		pos.Pop()
	}
}

func (pos *Position) Perf(depth int) {
	pos.Nodes = 0

	start := time.Now()

	pos.PerfRec(depth)

	elapsed := time.Now().Sub(start)

	fmt.Printf("elapsed %v , nodes %v , nps %.3f Mn/s\n", elapsed, pos.Nodes, float32(pos.Nodes)/(float32(elapsed)/1e9)/1e6)
}

func (pos *Position) Print() {
	fmt.Println(pos.PrettyPrintString())
}

func (pos *Position) ExecCommand(command string) {
	t := Tokenizer{command}
	i := t.GetInt()

	pos.Current().GenMoveBuff()

	mb := pos.Current().MoveBuff

	if i != 0 {
		i--
		if i < len(mb) {
			pos.Push(mb[i].Move)
			pos.Print()
		} else {
			fmt.Println("warning : move index out of range")
		}
	} else if command == "d" {
		if pos.StatePtr > 0 {
			pos.Pop()
			pos.Print()
		} else {
			fmt.Println("warning : no move to delete")
		}
	} else if command == "f" {
		if pos.StatePtr < pos.MaxStatePtr {
			pos.StatePtr++
			pos.Print()
		} else {
			fmt.Println("warning : no move forward")
		}
	} else if command == "perf" {
		pos.Perf(4)
	} else {
		found := false
		for _, mbi := range pos.Current().MoveBuff {
			if mbi.San == command {
				found = true
				pos.Push(mbi.Move)
				pos.Print()
				break
			}
		}
		if !found {
			fmt.Println("warning : unknown command or illegal move")
		}
	}
}
