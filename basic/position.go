package basic

import (
	"fmt"
	"time"
)

const MAX_STATES = 100

type Position struct {
	States   [MAX_STATES]State
	StatePtr int
	Nodes    int
}

func (pos *Position) Current() *State {
	return &pos.States[pos.StatePtr]
}

func (pos *Position) Init(variant Variant) {
	pos.StatePtr = 0
	pos.Current().Init(variant)
	pos.Current().Ply = 0
}

func (pos Position) PrettyPrintString() string {
	buff := pos.Current().PrettyPrintString()

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

	mat := PieceMaterialTables[p][sq]

	st.Material[color].Merge(mat)

	st.UpdateMaterialBalance()
}

func (st *State) Remove(sq Square) {
	st.Pieces[RankOf[sq]][FileOf[sq]] = NoPiece

	p := st.PieceAtSquare(sq)

	color := ColorOf[p]

	sqbb := sq.Bitboard()

	st.ByColor[color] &^= sqbb

	st.ByFigure[FigureOf[p]] &^= sqbb

	mat := PieceMaterialTables[p][sq]

	st.Material[color].UnMerge(mat)

	st.UpdateMaterialBalance()
}

func (st *State) MakeMove(move Move) {
	p := st.PieceAtSquare(move.FromSq())

	st.Remove(move.FromSq())

	if move.MoveType() == Promotion {
		st.Put(move.PromotionPiece(), move.ToSq())
	} else {
		st.Put(p, move.ToSq())
	}

	st.Move = move

	st.Ply++

	st.Turn = st.Turn.Inverse()

	if st.Turn == White {
		st.FullmoveNumber++
	}
}

func (pos *Position) Push(move Move) {
	oldState := pos.States[pos.StatePtr]

	pos.StatePtr++

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
		}
	} else if command == "d" {
		if pos.StatePtr > 0 {
			pos.Pop()
			pos.Print()
		} else {
			fmt.Println("warning : no move to delete")
		}
	} else if command == "perf" {
		pos.Perf(4)
	}
}
