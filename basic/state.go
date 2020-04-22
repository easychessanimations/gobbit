package basic

import (
	"fmt"
	"strings"
)

type Variant int

func (v Variant) String() string {
	return VariantInfos[v].DisplayName
}

const (
	VariantStandard = Variant(iota)
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
		StartFen:    "jlsesqkbnr/pppppppp/8/8/8/8/PPPPPPPP/JLneSQKBNR w KQkq - 0 1 -",
		DisplayName: "Eightpiece",
	},
}

const (
	CastlingSideKing = iota
	CastlingSideQueen
)

type ColorCastlingRights [2]bool
type CastlingRights [2]ColorCastlingRights

// State records the state of a position
type State struct {
	Variant           Variant
	Pieces            [NUM_RANKS][NUM_FILES]Piece
	Turn              Color
	CastlingRights    CastlingRights
	EnPassantSquare   Square
	EpSquare          Square
	HalfmoveClock     int
	FullmoveNumber    int
	HasDisabledMove   bool
	DisableFromSquare Square
	DisableToSquare   Square
	ByFigure          [FigureArraySize]Bitboard
	ByColor           [ColorArraySize]Bitboard
	Ply               int
	Move              Move
}

// CastlingRights.String() reports castling rights in fen format
func (crs CastlingRights) String() string {
	buff := ""

	if crs[White][CastlingSideKing] {
		buff += "K"
	}

	if crs[White][CastlingSideQueen] {
		buff += "Q"
	}

	if crs[Black][CastlingSideKing] {
		buff += "k"
	}

	if crs[Black][CastlingSideQueen] {
		buff += "q"
	}

	if buff == "" {
		return "-"
	}

	return buff
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
func (st *State) Init(variant Variant) {
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
		st.ParseTurnString(fenParts[1])
	}

	if len(fenParts) > 2 {
		st.ParseCastlingRights(fenParts[2])
	}

	if len(fenParts) > 3 {
		st.ParseEnPassantSquare(fenParts[3])
	}

	if len(fenParts) > 4 {
		st.ParseHalfmoveClock(fenParts[4])
	}

	if len(fenParts) > 5 {
		st.ParseFullmoveNumber(fenParts[5])
	}

	if len(fenParts) > 6 {
		st.ParseDisabledMove(fenParts[6])
	}

	return nil
}

func (st *State) ParseDisabledMove(dms string) {
	st.HasDisabledMove = false

	if len(dms) == 0 {
		return
	}

	if dms[:1] == "-" {
		return
	}

	t := Tokenizer{}
	t.Init(dms)

	st.DisableFromSquare = t.GetSquare()
	st.DisableToSquare = t.GetSquare()

	if st.DisableFromSquare != st.DisableToSquare {
		st.HasDisabledMove = true
	}
}

func (st *State) ParseEnPassantSquare(epsqs string) {
	t := Tokenizer{}
	t.Init(epsqs)
	st.EnPassantSquare = t.GetSquare()
}

func (st *State) ParseHalfmoveClock(hmcs string) {
	t := Tokenizer{}
	t.Init(hmcs)
	st.HalfmoveClock = t.GetInt()
}

func (st *State) ParseFullmoveNumber(fmns string) {
	t := Tokenizer{}
	t.Init(fmns)
	st.FullmoveNumber = t.GetInt()
}

func (st *State) ParseCastlingRights(crs string) {
	t := Tokenizer{}
	t.Init(crs)

	st.CastlingRights = t.GetCastlingRights()
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

	buff += " " + st.CastlingRights.String()

	if st.EnPassantSquare == SquareA1 {
		buff += " -"
	} else {
		buff += " " + st.EnPassantSquare.UCI()
	}

	buff += " " + fmt.Sprintf("%d", st.HalfmoveClock)

	buff += " " + fmt.Sprintf("%d", st.FullmoveNumber)

	if st.Variant == VariantEightPiece {
		if st.HasDisabledMove {
			buff += " " + st.DisableFromSquare.UCI() + st.DisableToSquare.UCI()
		} else {
			buff += " -"
		}

	}

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
					st.CalculateOccupancy()

					return nil
				}
			}

			ps = t.GetFenPiece()
		}
	}

	return fmt.Errorf("too few pieces in placement string")
}

func (st State) PieceAtSquare(sq Square) Piece {
	return st.Pieces[RankOf[sq]][FileOf[sq]]
}

func (st State) MoveLAN(move Move) string {
	fromPiece := st.PieceAtSquare(move.FromSq())

	return fromPiece.SanSymbol() + move.String()
}

func (st *State) CalculateOccupancy() {
	st.ByFigure = [FigureArraySize]Bitboard{}
	st.ByColor = [ColorArraySize]Bitboard{}

	for sq := SquareMinValue; sq <= SquareMaxValue; sq++ {
		bb := sq.Bitboard()

		p := st.PieceAtSquare(sq)

		if p != NoPiece {
			fig := FigureOf[p]

			col := ColorOf[p]

			st.ByFigure[fig] |= bb
			st.ByColor[col] |= bb
		}
	}
}

func (st State) OccupUs() Bitboard {
	return st.ByColor[st.Turn]
}

func (st State) OccupThem() Bitboard {
	return st.ByColor[st.Turn.Inverse()]
}
