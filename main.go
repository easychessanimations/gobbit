package main

import (
	"fmt"

	. "github.com/easychessanimations/gobbit/basic"
)

func main() {
	st := State{}
	st.ParsePlacementString("jlsesqkbnr/pppppppp/8/8/8/8/PPPPPPPP/JLneSQKBNR w KQkq - 0 1 -")

	fmt.Println(st.PrettyPlacementString())
}
