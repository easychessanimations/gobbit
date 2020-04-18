package basic

import "strings"

// Figure tells the Figure of a Piece
func (p Piece) Figure() Figure {
	return Figure(p >> 1)
}

// Color tells the color of a Piece
func (p Piece) Color() Color {
	return Color(p & Piece(COLOR_MASK))
}

// SanSymbol tells the SAN symbol of a Piece ( letter always upper case )
func (p Piece) SanSymbol() string {
	sym := SymbolOf[FigureOf[p]]

	return strings.ToUpper(sym[0:1]) + sym[1:]
}

// IsLancer tells whether a Piece is a lancer
func (p Piece) IsLancer() bool {
	return (FigureOf[p] &^ Figure(LANCER_DIRECTION_MASK)) == LancerMinValue
}

// PrettySymbol tells the pretty print symbol of a Piece
func (p Piece) PrettySymbol() string {
	sym := p.FenSymbol()

	if len(sym) == 1 {
		return " " + sym + " "
	}

	if len(sym) == 2 {
		return sym + " "
	}

	return sym
}

// SanLetter tells the SAN letter of a Piece
func (p Piece) SanLetter() string {
	return p.SanSymbol()[0:1]
}

// FenSymbol tells the FEN symbol of a Piece
func (p Piece) FenSymbol() string {
	if ColorOf[p] == White {
		return p.SanSymbol()
	}
	return SymbolOf[FigureOf[p]]
}

// UCI tells the UCI symbol of a Piece ( letter always lower case )
func (p Piece) UCI() string {
	return SymbolOf[FigureOf[p]]
}

func init() {
	for i := 0; i <= int(FigureMaxValue); i++ {
		FigureOf[i*2] = Figure(i)
		FigureOf[i*2+1] = Figure(i)

		ColorOf[i*2] = Black
		ColorOf[i*2+1] = White

		ColorFigure[Black][i] = Piece(i * 2)
		ColorFigure[White][i] = Piece(i*2 + 1)
	}
}
