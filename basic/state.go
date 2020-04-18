package basic

import (
	"fmt"
	"strings"
)

const (
	VariantStandard = iota
	VariantEightPiece
)

type VariantInfo struct {
	StartFen    string
	DisplayName string
}

var VariantInfos = []VariantInfo{
	{ // standard
		StartFen:    "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		DisplayName: "Standard",
	},
	{ // eightpiece
		StartFen:    "jlsesqkbnr/pppppppp/8/8/8/8/PPPPPPPP/JLneSQKBNR w KQkq - 0 1",
		DisplayName: "Eightpiece",
	},
}

// State records the state of a position
type State struct {
	Variant int
	Pieces  [NUM_RANKS][NUM_FILES]Piece
	Turn    Color
}

// Color.String() converts a color to "b" for black, "w" for white and "-" for no color
func (color Color) String() string {
	if color == Black {
		return "b"
	}

	if color == White {
		return "w"
	}

	return "-"
}

// Init initializes state
// sets itself up from variant start fen
func (st *State) Init(variant int) {
	st.Variant = variant
	st.ParseFen(VariantInfos[st.Variant].StartFen)
}

// ParseFen sets up state from a fen
func (st *State) ParseFen(fen string) error {
	fenParts := strings.Split(fen, " ")

	if len(fenParts) > 0 {
		err := st.ParsePlacementString(fenParts[0])
		if err != nil {
			return err
		}
	}

	if len(fenParts) > 1 {
		st.ParseTurnString(fenParts[0])
	}

	return nil
}

func (st *State) ParseTurnString(ts string) {
	t := Tokenizer{}
	t.Init(ts)

	color := t.GetColor()

	st.Turn = White

	if color != NoColor {
		st.Turn = color
	}
}

// PrettyPrintString returns the state pretty print string
func (st State) PrettyPrintString() string {
	buff := st.PrettyPlacementString()

	buff += "\n" + VariantInfos[st.Variant].DisplayName + " : " + st.ReportFen()

	return buff
}

// ReportFen reports the state as a fen string
func (st State) ReportFen() string {
	buff := ""

	cum := 0

	for rank := LAST_RANK; rank >= 0; rank-- {
		for file := 0; file < NUM_FILES; file++ {
			p := st.Pieces[rank][file]

			if p == NoPiece {
				cum++
			} else {
				if cum > 0 {
					buff += fmt.Sprintf("%d", cum)
				}
				cum = 0
				buff += p.FenSymbol()
			}
		}
		if cum > 0 {
			buff += fmt.Sprintf("%d", cum)
		}
		cum = 0
		if rank > 0 {
			buff += "/"
		}
	}

	buff += " " + st.Turn.String()

	return buff
}

// PrettyPlacementString returns the pretty string representation of the board
func (st State) PrettyPlacementString() string {
	buff := ""

	for rank := LAST_RANK; rank >= 0; rank-- {
		for file := 0; file < NUM_FILES; file++ {
			buff += st.Pieces[rank][file].PrettySymbol()
		}
		if (rank == 0 && st.Turn == White) || (rank == LAST_RANK && st.Turn == Black) {
			buff += " *"
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
