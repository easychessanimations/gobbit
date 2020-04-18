package basic

// State records the state of a position
type State struct {
	Pieces [NUM_RANKS][NUM_FILES]Piece
}

// PrettyPlacementString returns the pretty string representation of the board
func (st State) PrettyPlacementString() string {
	buff := ""

	for rank := LAST_RANK; rank >= 0; rank-- {
		for file := 0; file < NUM_FILES; file++ {
			buff += st.Pieces[rank][file].PrettySymbol()
		}
		buff += "\n"
	}

	return buff
}
