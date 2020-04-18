package basic

import "testing"

func TestPiece(t *testing.T) {
	lnw := LancerNW

	sym := SymbolOf[lnw]

	if sym != "lnw" {
		t.Errorf("wrong symbol for LancerNW, expected lnw, got %v", sym)
	}

	wkf := WhiteKing.Figure()

	if wkf != King {
		t.Errorf("wrong figure for white king, expected %v, got %v", King, wkf)
	}
}
