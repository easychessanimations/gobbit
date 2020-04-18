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

	fof := FigureOf[BlackSentry]

	pfm := BlackSentry.Figure()

	if fof != pfm {
		t.Errorf("figure of array %v and figure member %v are not equal", fof, pfm)
	}

	cf := ColorFigure[White][Knight]

	if cf != WhiteKnight {
		t.Errorf("color figure of white knight %v does not equal white knight %v", cf, WhiteKnight)
	}

	colWk := WhiteKing.Color()

	if colWk != White {
		t.Errorf("white king color %v is not white %v", colWk, White)
	}

	colBs := BlackSentry.Color()

	if colBs != Black {
		t.Errorf("black sentry color %v is not black %v", colBs, Black)
	}

	cof := ColorOf[WhiteRook]

	pcm := WhiteRook.Color()

	if cof != pcm {
		t.Errorf("color of white rook %v is not equal to white rook color %v", cof, pcm)
	}

	fsBp := BlackPawn.FenSymbol()

	if fsBp != "p" {
		t.Errorf("expected fen symbol of black pawn p, got %v", fsBp)
	}

	fsWp := WhitePawn.FenSymbol()

	if fsWp != "P" {
		t.Errorf("expected fen symbol of white pawn P, got %v", fsWp)
	}

	us := WhiteBishop.UCI()

	if us != "b" {
		t.Errorf("expected b as UCI symbol of white bishop, got %v", us)
	}

	sanBlackLancerNW := BlackLancerNW.SanSymbol()

	if sanBlackLancerNW != "Lnw" {
		t.Errorf("expected Lnw as san symbol of black lancer north west, got %v", sanBlackLancerNW)
	}

	sanLetterBlackLancerSW := BlackLancerSW.SanLetter()

	if sanLetterBlackLancerSW != "L" {
		t.Errorf("expected L as san letter of black lancer south west, got %v", sanLetterBlackLancerSW)
	}

	if !WhiteLancerE.IsLancer() {
		t.Errorf("white lancer east is not a lancer")
	}

	if BlackQueen.IsLancer() {
		t.Errorf("black queen is a lancer")
	}

	p, ok := SymbolToPiece["Lnw"]

	if !ok {
		t.Errorf("figure for symbol lnw not found")
	}

	if p != WhiteLancerNW {
		t.Errorf("expected %v as piece for Lnw, got %v", LancerNW, p)
	}
}
