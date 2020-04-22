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

func (st *State) Put(p Piece, sq Square) {
	if p == NoPiece {
		return
	}

	st.Pieces[RankOf[sq]][FileOf[sq]] = p

	sqbb := sq.Bitboard()

	st.ByColor[ColorOf[p]] |= sqbb

	st.ByFigure[FigureOf[p]] |= sqbb
}

func (st *State) Remove(sq Square) {
	st.Pieces[RankOf[sq]][FileOf[sq]] = NoPiece

	p := st.PieceAtSquare(sq)

	sqbb := sq.Bitboard()

	st.ByColor[ColorOf[p]] &^= sqbb

	st.ByFigure[FigureOf[p]] &^= sqbb
}

func (st *State) MakeMove(move Move) {
	p := st.PieceAtSquare(move.FromSq())

	st.Remove(move.FromSq())

	st.Put(p, move.ToSq())

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

	fmt.Println(elapsed, pos.Nodes)
}
