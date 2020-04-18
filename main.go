package main

import (
	"fmt"

	. "github.com/easychessanimations/gobbit/basic"
)

func main() {
	st := State{}
	st.Init(VariantEightPiece)

	fmt.Println(st.PrettyPrintString())
}
