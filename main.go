package main

import (
	"fmt"

	. "github.com/easychessanimations/gobbit/basic"
)

func main() {
	st := State{}
	st.Init(VariantEightPiece)

	//fmt.Println(st.PrettyPrintString())

	/*fmt.Println(BISHOP_MAGICS)
	fmt.Println(ROOK_MAGICS)*/

	fmt.Println(TotalMagicEntries)

	fmt.Println(RookMobility(SquareD4, SquareD5.Bitboard()|SquareF4.Bitboard()|SquareD2.Bitboard()|SquareB4.Bitboard()))
}
