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

	fmt.Println(Wizards)
}
