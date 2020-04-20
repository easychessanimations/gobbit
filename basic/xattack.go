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

func SearchMagic(sq Square, sqs []Square) (int, uint64, bool, int) {
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
			hash := make(map[uint64]uint64)
			coll := 0
			for enum = 0; enum < 1<<len(sqs); enum++ {
				tr := uint64(Translate(sqs, enum))
				key := (magic * tr) >> (64 - shift)
				mask, foundKey := hash[key]
				if foundKey {
					if mask != tr {
						coll++
						break
					}
				} else {
					hash[key] = tr
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

func init() {
	maxShift := 0
	for sq := SquareMinValue; sq <= SquareMaxValue; sq++ {
		sqs, _ := GenAttackSquares(sq, NormalizedBishopDirection)
		//fmt.Println(bb)
		shift, magic, ok, nodes := SearchMagic(sq, sqs)
		if shift > maxShift {
			maxShift = shift
		}
		if ok {
			fmt.Printf("Found bishop magic for %v shift %2d magic %016x nodes %d sqs  %d\n", sq, shift, magic, nodes, len(sqs))
		} else {
			fmt.Println("Failed", sq)
			break
		}
		sqs, _ = GenAttackSquares(sq, NormalizedRookDirection)
		//fmt.Println(bb)
		shift, magic, ok, nodes = SearchMagic(sq, sqs)
		if shift > maxShift {
			maxShift = shift
		}
		if ok {
			fmt.Printf("Found rook   magic for %v shift %2d magic %016x nodes %d sqs %d\n", sq, shift, magic, nodes, len(sqs))
		} else {
			fmt.Println("Failed", sq)
			break
		}
	}
	fmt.Println("max shift", maxShift)
}
