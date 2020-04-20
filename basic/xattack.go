package basic

import (
	"fmt"
	"math/rand"
)

func BishopAttackSquares(sq Square) []Square {
	sqs := []Square{}

	middle := ^BbBorder

	var testSq Square = middle.Pop()

	for ; testSq != 0; testSq = middle.Pop() {
		_, ok := NormalizedBishopDirection(sq, testSq)
		if ok {
			sqs = append(sqs, testSq)
		}
	}

	return sqs
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

func SearchBishopMagic(sq Square) (int, uint64, bool, int) {
	sqs := BishopAttackSquares(sq)
	var enum uint64
	var lastGoodMagic uint64
	var lastGoodShift int
	foundMagic := false
	nodes := 0
	for shift := 14; shift > 6; shift-- {
		found := false
		for loop := 0; loop < 100; loop++ {
			nodes++
			magic := randMagic() >> 6 //+ uint64(64-shift)<<58
			hash := make(map[uint64]int)
			coll := 0
			for enum = 0; enum < 1<<len(sqs); enum++ {
				key := (magic * enum) >> (64 - shift)
				cnt, found := hash[key]
				if found {
					hash[key] = cnt + 1
					coll++
				} else {
					hash[key] = 1
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
	for sq := SquareMinValue; sq <= SquareMaxValue; sq++ {
		shift, magic, ok, nodes := SearchBishopMagic(sq)
		if ok {
			fmt.Println("Found bishop magic for", sq, shift, magic, nodes)
		} else {
			fmt.Println("Failed", sq)
			break
		}
	}
}
