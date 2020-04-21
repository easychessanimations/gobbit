package main

import (
	"fmt"

	. "github.com/easychessanimations/gobbit/basic"
)

func main() {
	st := State{}
	st.Init(VariantEightPiece)

	//fmt.Println(st)

	//fmt.Println(st.PrettyPrintString())

	/*fmt.Println(BISHOP_MAGICS)
	fmt.Println(ROOK_MAGICS)*/

	//fmt.Println(TotalMagicEntries)

	//fmt.Println(RookMobility(SquareD4, SquareD5.Bitboard()|SquareF4.Bitboard()|SquareD2.Bitboard()|SquareB4.Bitboard()))

	/*fmt.Println("rook attack b2")
	fmt.Println(RookAttack[SquareB2])
	fmt.Println("bishop attack e4")
	fmt.Println(BishopAttack[SquareE4])
	fmt.Println("queen attack f7")
	fmt.Println(QueenAttack[SquareF7])*/

	/*m := MakeMoveFT(SquareG1, SquareF3)

	fmt.Println(st.MoveLAN(m))*/

	fmt.Println(st.PslmsForPieceAtSquare(Violent|Quiet, WhiteKnight, SquareE5, st.OccupUs(), st.OccupThem()))
}
