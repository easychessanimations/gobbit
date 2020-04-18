package basic

// Figure tells the Figure of a Piece
func (p Piece) Figure() Figure {
	return Figure(p >> 1)
}
