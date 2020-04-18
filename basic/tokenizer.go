package basic

import "fmt"

type Tokenizer struct {
	Content string
}

func (t *Tokenizer) Init(content string) {
	t.Content = content
}

func (t *Tokenizer) GetColor() Color {
	if len(t.Content) == 0 {
		return NoColor
	}

	if t.Content[0] == 'b' {
		t.Content = t.Content[1:]
		return Black
	}

	if t.Content[0] == 'w' {
		t.Content = t.Content[1:]
		return White
	}

	return NoColor
}

func (t *Tokenizer) GetFenPiece() []Piece {
	if len(t.Content) <= 0 {
		return []Piece{}
	}

	c := t.Content[0]

	if c == '/' {
		t.Content = t.Content[1:]
	}

	if len(t.Content) <= 0 {
		return []Piece{}
	}

	c = t.Content[0]

	if c >= '0' && c <= '9' {
		t.Content = t.Content[1:]

		return make([]Piece, c-'0')
	}

	if c == 'l' || c == 'L' {
		if len(t.Content) <= 1 {
			// lancer without direction
			return []Piece{}
		}

		sym := ""

		if t.Content[1] == 'n' || t.Content[1] == 's' {
			if len(t.Content) <= 2 {
				sym = t.Content[:2]
				t.Content = t.Content[2:]
			} else {
				if t.Content[2] == 'e' || t.Content[2] == 'w' {
					sym = t.Content[:3]
					t.Content = t.Content[3:]
				}
			}
		} else {
			if t.Content[1] == 'e' || t.Content[1] == 'w' {
				sym = t.Content[:2]
				t.Content = t.Content[2:]
			}
		}

		p, ok := SymbolToPiece[sym]

		if !ok {
			panic(fmt.Sprintf("got no piece for %v", sym))
		}

		return []Piece{p}
	}

	p, ok := SymbolToPiece[t.Content[0:1]]

	if ok {
		t.Content = t.Content[1:]
		return []Piece{p}
	}

	return []Piece{}
}
