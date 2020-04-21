package basic

import (
	"fmt"
	"math/rand"
)

func GenAttackSquares(sq Square, normFunc func(Square, Square) (Delta, bool)) ([]Square, Bitboard) {
	sqs := []Square{}

	middle := BbFull

	if RankOf[sq] == 0 {
		middle &^= BbRank8
	}

	if RankOf[sq] == LAST_RANK {
		middle &^= BbRank1
	}

	if FileOf[sq] == 0 {
		middle &^= BbFileH
	}

	if FileOf[sq] == LAST_FILE {
		middle &^= BbFileA
	}

	if RankOf[sq] >= 1 && RankOf[sq] <= (LAST_RANK-1) {
		middle &^= BbRank1
		middle &^= BbRank8
	}

	if FileOf[sq] >= 1 && FileOf[sq] <= (LAST_FILE-1) {
		middle &^= BbFileA
		middle &^= BbFileH
	}

	if (sq.Bitboard() &^ BbBorder) != 0 {
		middle = ^BbBorder
	}

	middle = BbFull

	bb := BbEmpty

	var testSq Square

	for middle != 0 {
		testSq = middle.Pop()
		_, ok := normFunc(sq, testSq)
		if ok {
			sqs = append(sqs, testSq)
			bb |= testSq.Bitboard()
		}
	}

	return sqs, bb
}

func Translate(sqs []Square, occup uint64) Bitboard {
	bb := BbEmpty
	for i := 0; i < len(sqs); i++ {
		if occup&1 != 0 {
			bb |= sqs[i].Bitboard()
		}
		occup >>= 1
	}
	return bb
}

var Rand = rand.New(rand.NewSource(1))

func randMagic() uint64 {
	r := uint64(Rand.Int63())
	r &= uint64(Rand.Int63())
	r &= uint64(Rand.Int63())
	return r << 1
}

func SlidingAttack(sq Square, deltas []Delta, occup Bitboard) (Bitboard, []Square) {
	bb := BbEmpty

	sqs := []Square{}

	for _, delta := range deltas {
		testSq := sq

		rank := RankOf[testSq] + delta.dRank
		file := FileOf[testSq] + delta.dFile

		ok := true

		for rank >= 0 && rank < NUM_RANKS && file >= 0 && file < NUM_FILES && ok {
			testSq = RankFile[rank][file]

			bb |= testSq.Bitboard()

			sqs = append(sqs, testSq)

			rank += delta.dRank
			file += delta.dFile

			if (testSq.Bitboard() & occup) != 0 {
				ok = false
			}
		}
	}

	return bb, sqs
}

func SearchMagic(sq Square, sqs []Square, deltas []Delta) (int, uint64, bool, int) {
	Rand = rand.New(rand.NewSource(1))

	var enum uint64
	var lastGoodMagic uint64
	var lastGoodShift int
	foundMagic := false
	nodes := 0
	for shift := 22; shift > 6; shift-- {
		found := false
		for loop := 0; loop < 5000; loop++ {
			nodes++
			magic := randMagic() >> 6 //+ uint64(64-shift)<<58
			hash := make(map[uint64]Bitboard)
			coll := 0
			for enum = 0; enum < 1<<len(sqs); enum++ {
				trb := Translate(sqs, enum)
				mobility, _ := SlidingAttack(sq, deltas, trb)
				key := (magic * uint64(trb)) >> (64 - shift)
				storedMobility, foundKey := hash[key]
				if foundKey {
					if storedMobility != mobility {
						coll++
						break
					}
				} else {
					hash[key] = mobility
				}
				//fmt.Println(Translate(sqs, enum))
				//fmt.Println()
			}
			if coll == 0 {
				//fmt.Println(shift, loop, magic, coll)
				foundMagic = true
				lastGoodMagic = magic
				lastGoodShift = shift
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	if foundMagic {
		return lastGoodShift, lastGoodMagic, true, nodes
	}

	return 0, 0, false, nodes
}

var BISHOP_DELTAS = []Delta{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
var ROOK_DELTAS = []Delta{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

type Wizard struct {
	Name   string
	Deltas []Delta
	Tries  int
}

var wizards = []Wizard{
	{
		Name:   "bishop",
		Deltas: BISHOP_DELTAS,
		Tries:  5000,
	},
	{
		Name:   "rook",
		Deltas: ROOK_DELTAS,
		Tries:  5000,
	},
}

func (wiz *Wizard) GenAttacks() {
	fmt.Println("generating attacks for", wiz.Name)
	maxShift := 0
	for sq := SquareMinValue; sq <= SquareMaxValue; sq++ {
		_, sqs := SlidingAttack(sq, wiz.Deltas, BbEmpty)
		//fmt.Println(bb)
		shift, magic, ok, nodes := SearchMagic(sq, sqs, wiz.Deltas)
		if shift > maxShift {
			maxShift = shift
		}
		if ok {
			fmt.Printf("found %-6s magic for %v %2d shift %2d max shift %2d magic %016x nodes %6d sqs %2d\n", wiz.Name, sq, sq, shift, maxShift, magic, nodes, len(sqs))
		} else {
			fmt.Println("failed", wiz.Name, "at", sq)
			break
		}
	}
	fmt.Println("max shift for", wiz.Name, "=", maxShift)
}

func init() {
	for _, wiz := range wizards {
		wiz.GenAttacks()
	}
}
