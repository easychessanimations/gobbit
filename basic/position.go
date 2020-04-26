package basic

import (
	"fmt"
	"strings"
	"time"
)

const MAX_STATES = 100

type Position struct {
	States        [MAX_STATES]State
	StatePtr      int
	MaxStatePtr   int
	Nodes         int
	SearchStopped bool
}

func (pos *Position) Current() *State {
	return &pos.States[pos.StatePtr]
}

func (pos *Position) Init(variant Variant) {
	pos.StatePtr = 0
	pos.Current().Init(variant)
	pos.Current().Ply = 0
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

	st.Remove(move.FromSq())

	st.Remove(move.ToSq())

	if move.MoveType() == Promotion {
		st.Put(move.PromotionPiece(), move.ToSq())
	} else {
		st.Put(p, move.ToSq())
	}

	st.Move = move

	st.Ply++

	st.SetSideToMove(st.Turn.Inverse())

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

func (pos *Position) Pop() {
	pos.StatePtr--
}

func (pos *Position) PerfRec(remDepth int) {
	pos.Nodes++

	if remDepth == 0 {
		return
	}

	for _, move := range pos.Current().GenerateMoves() {
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
