package main

import (
	"fmt"

	. "github.com/easychessanimations/gobbit/basic"
)

func main() {
	st := State{}
	st.Pieces[Rank8][FileB] = BlackLancerSE
	st.Pieces[Rank1][FileB] = WhiteLancerNE
	fmt.Println(st.PrettyPlacementString())
}
