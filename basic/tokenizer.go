package basic

import "fmt"

type Tokenizer struct {
	Content string
}

func (t *Tokenizer) Init(content string) {
	t.Content = content
}

func (t *Tokenizer) GetToken() (string, bool){	
	buff := ""

	var i int

	hasToken := false

	for i = 0; i < len(t.Content); i++{
		c := t.Content[i:i+1]		
		if c != " "{
			hasToken = true
			buff += c	
		}else{
			i++
			break
		}		
	}

	if hasToken{		
		t.Content = t.Content[i:]
		return buff, true
	}

	return "", false
}

func (t *Tokenizer) GetTokensUpTo(upTo string) []string{	
	tokens := []string{}

	for true{
		token, ok := t.GetToken()

		if (!ok) || token == upTo{
			return tokens
		}

		tokens = append(tokens, token)
	}

	return tokens
}

func (t *Tokenizer) GetSquare() Square {
	file := 0
	rank := 0

	parseFile := true

	for {
		if len(t.Content) == 0 {
			return RankFile[rank][file]
		}

		c := t.Content[0]

		if parseFile {
			if c >= 'a' && c <= 'z' {
				file = int(c - 'a')
				t.Content = t.Content[1:]
				parseFile = false
			} else {
				return RankFile[rank][file]
			}
		} else {
			if c >= '1' && c <= '9' {
				rank = int(c - '1')
				t.Content = t.Content[1:]
			}
			return RankFile[rank][file]
		}
	}
}

func (t *Tokenizer) GetInt() int {
	num := 0

	for {
		if len(t.Content) == 0 {
			return num
		}

		if t.Content[0] >= '0' && t.Content[0] <= '9' {
			num *= 10
			num += int(t.Content[0] - '0')
			t.Content = t.Content[1:]
		} else {
			return num
		}
	}
}

func (t *Tokenizer) GetCastlingRights() CastlingRights {
	ccrs := [2]ColorCastlingRights{}

	for {
		if len(t.Content) == 0 {
			return ccrs
		}

		c := t.Content[0]

		if c == 'K' {
			ccrs[White][CastlingSideKing].CanCastle = true
			t.Content = t.Content[1:]
		} else if c == 'Q' {
			ccrs[White][CastlingSideQueen].CanCastle = true
			t.Content = t.Content[1:]
		} else if c == 'k' {
			ccrs[Black][CastlingSideKing].CanCastle = true
			t.Content = t.Content[1:]
		} else if c == 'q' {
			ccrs[Black][CastlingSideQueen].CanCastle = true
			t.Content = t.Content[1:]
		} else {
			return ccrs
		}
	}
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
				}else{
					sym = t.Content[:2]
					t.Content = t.Content[2:]
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
