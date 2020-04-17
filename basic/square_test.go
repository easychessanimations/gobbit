package basic

import "testing"

const ITERATION_COUNT = 100000

func TestSquare(t *testing.T) {
	flc := FileLetterOf[FileC]
	rl6 := RankLetterOf[Rank6]
	if rl6 != "6" {
		t.Errorf("wrong rank letter of Rank6, expected 6, got %s", rl6)
	}
	if flc != "c" {
		t.Errorf("wrong file letter of FileC, expected c, got %s", flc)
	}
	sqUci := UCIOf[SquareE4]
	if sqUci != "e4" {
		t.Errorf("wrong UCI of square e4, expected e4, got %s", sqUci)
	}
	sq, ok := UCIToSquare["e4"]
	if !ok {
		t.Errorf("UCI e4 could not be converted to Square")
	}
	if sq != SquareE4 {
		t.Errorf("wrong square for UCI e4, expected %v, got %v", SquareE4, sq)
	}
	sq = RankFile[Rank7][FileG]
	if sq != SquareG7 {
		t.Errorf("wrong square for rank 7, file g, expected %v, got %v", SquareG7, sq)
	}
	rof := RankOf[SquareD4]
	sqr := SquareD4.Rank()
	if rof != sqr {
		t.Errorf("array ( %v ) and member ( %v ) access to rank of square differ", rof, sqr)
	}
	fof := FileOf[SquareD4]
	sqf := SquareD4.File()
	if fof != sqf {
		t.Errorf("array ( %v ) and member ( %v ) access to file of square differ", fof, sqf)
	}
}

func BenchmarkSquareArray(b *testing.B) {
	for cnt := 0; cnt < ITERATION_COUNT; cnt++ {
		for i := 0; i < BOARD_AREA; i++ {
			rank, file := RankOf[Square(i)], FileOf[Square(i)]

			file, rank = rank, file
		}
	}
}

func BenchmarkSquareMember(b *testing.B) {
	for cnt := 0; cnt < ITERATION_COUNT; cnt++ {
		for i := 0; i < BOARD_AREA; i++ {
			rank, file := Square(i).Rank(), Square(i).File()

			file, rank = rank, file
		}
	}
}
