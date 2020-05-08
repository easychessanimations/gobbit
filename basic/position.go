package basic

import (
	"fmt"
	"strings"
	"time"
	"os"
)

const SEARCH_MAX_DEPTH = 100

const INFINITE_DEPTH = int8(SEARCH_MAX_DEPTH + 1)

const MAX_STATES = SEARCH_MAX_DEPTH + 1

const MAX_MULTIPV = 20

type MultiPvInfo struct{
	Index                    int
	Depth                    int
	Time                     int
	Nodes                    int	
	Nps                      float32
	Score                    Score
	Pv                       []Move
	PvUCI                    string
}

func (mpi MultiPvInfo) String() string{
	return fmt.Sprintf("info multipv %d depth %d time %d nodes %d nps %.0f score cp %d pv %v", mpi.Index, mpi.Depth, mpi.Time, mpi.Nodes, mpi.Nps, mpi.Score, mpi.PvUCI)
}

type MultiPvInfos [MAX_MULTIPV]MultiPvInfo

func (mpis MultiPvInfos) Len() int{
	return len(mpis)
}

func (mpis MultiPvInfos) Swap(i, j int){
	mpis[i], mpis[j] = mpis[j], mpis[i]
}

func (mpis MultiPvInfos) Less(i, j int) bool{
	if mpis[j].Depth != mpis[i].Depth{
		return mpis[j].Depth > mpis[i].Depth
	}

	return mpis[j].Score > mpis[i].Score
}

type Position struct {
	States                   [MAX_STATES]State
	StatePtr                 int
	SearchRootPtr            int
	MaxStatePtr              int
	Nodes                    int
	SearchStopped            bool
	NullMovePruning          bool
	NullMovePruningMinDepth  int
	NullMoveDepthReduction   int
	StackReduction           bool
	AspirationWindow         bool
	PvTable                  *PvHash
	PosMoveTable             *PosMoveHash
	LastRootPvScore          Score
	LastGoodPv               []Move
	Start                    time.Time
	CheckPoint               time.Time
	Depth                    int
	Verbose                  bool
	IgnoreRootMoves          []Move
	PruningAgressivity       int
	PruningReduction         int
	MultiPV                  int
	MultiPvInfos             MultiPvInfos
	OldMultiPvInfos          MultiPvInfos
	MultiPvIndex             int
	LogFilePath              string
}

func (pos Position) Log(content string){
	if content == ""{
		return
	}

	fmt.Println(content)

	if pos.LogFilePath != ""{
		f, err := os.OpenFile(pos.LogFilePath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString(content + "\n"); err != nil {
			fmt.Println(err)
		}
	}
}

func (pos *Position) ClearPvTable(){
	for i := 0; i < len(pos.PvTable.Entries); i++{
		pos.PvTable.Entries[i].Depth = INFINITE_DEPTH
	}
}

func (pos *Position) ClearPosMoveTable(){
	for i := 0; i < len(pos.PosMoveTable.Entries); i++{
		pos.PosMoveTable.Entries[i].Used = false
	}
}

func (pos Position) SearchRoot() *State{
	return &pos.States[pos.SearchRootPtr]
}

func (sc Score) IsMateInN() bool{
	return sc < -MAX_SCORE || sc > MAX_SCORE
}

func (pos Position) IsMateInN() bool{
	return pos.LastRootPvScore.IsMateInN()
}

func (pos Position) PvUCI() string{
	buff := []string{}

	for _, testMove := range pos.LastGoodPv {
		buff = append(buff, testMove.UCI())
	}

	pv := strings.Join(buff, " ")

	return pv
}

func (pos Position) Time() float32{
	return float32(time.Now().Sub(pos.Start)) / 1e9
}

func (pos Position) Nps() float32{
	return float32(pos.Nodes) / pos.Time()
}

func (pos Position) TimeMs() int{
	return int(float32(time.Now().Sub(pos.Start)) / 1e6)
}

func (pos Position) CheckTime() float32{
	return float32(time.Now().Sub(pos.CheckPoint)) / 1e9
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

	for ptr := pos.StatePtr; ptr > 0; ptr-- {
		prevSt := pos.States[ptr-1]
		san := prevSt.MoveToSan(pos.States[ptr].Move)		
		sans = append(sans, san)
	}

	return strings.Join(sans, " ")
}

func (pos Position) PrettyPrintString() string {
	buff := pos.Current().PrettyPrintString()

	buff += fmt.Sprintf("\n\n%s\n", pos.Line())

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

	currentZobrist := pos.Zobrist()

	for testPtr := 0; testPtr < ( pos.StatePtr - 1 ); testPtr++{
		if pos.States[testPtr].Zobrist == currentZobrist{
			// repetition
			return true, 0
		}
	}

	return false, 0
}

func (st *State) MakeMove(move Move) {
	if move == NullMove{		
		st.SetSideToMove(st.Turn.Inverse())

		if st.Turn == White {
			st.FullmoveNumber++
		}

		return
	}

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
		if st.CastlingRights[pCol].CanCastle() && move.MoveType() != Castling{
			// if king lost castling rights without castling, note that
			st.LostCastlingForColor[pCol] = true
		}

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
