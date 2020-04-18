package basic

import "fmt"

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

// ParsePlacementString parses a placement string and sets pieces accordingly
// returns an error if there are not enough pieces to fill the board
// liberal otherwise
func (st *State) ParsePlacementString(ps string) error {
	t := Tokenizer{}
	t.Init(ps)

	rank := LAST_RANK
	file := 0

	for ps := t.GetFenPiece(); len(ps) > 0; {
		if len(ps) > 0 {
			for _, p := range ps {
				st.Pieces[rank][file] = p
				file++
				if file > LAST_FILE {
					file = 0
					rank--
				}
				if rank < 0 {
					return nil
				}
			}

			ps = t.GetFenPiece()
		}
	}

	return fmt.Errorf("too few pieces in placement string")
}
