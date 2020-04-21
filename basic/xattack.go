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

func SearchMagic(sq Square, sqs []Square, deltas []Delta, tries int) (int, uint64, bool, int) {
	Rand = rand.New(rand.NewSource(1))

	var enum uint64
	var lastGoodMagic uint64
	var lastGoodShift int
	foundMagic := false
	nodes := 0
	for shift := 22; shift > 6; shift-- {
		found := false
		for loop := 0; loop < tries; loop++ {
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
var KNIGHT_DELTAS = []Delta{{1, 2}, {1, -2}, {-1, 2}, {-1, -2}, {2, 1}, {2, -1}, {-2, 1}, {-2, -1}}
var KING_DELTAS = append(BISHOP_DELTAS, ROOK_DELTAS...)

type Wizard struct {
	Name   string
	Deltas []Delta
	Tries  int
	Magics Magics
}

var Wizards = []Wizard{
	{
		Name:   "bishop",
		Deltas: BISHOP_DELTAS,
		Tries:  50000,
		Magics: BISHOP_MAGICS,
	},
	{
		Name:   "rook",
		Deltas: ROOK_DELTAS,
		Tries:  50000,
		Magics: ROOK_MAGICS,
	},
}

const BISHOP_WIZARD_INDEX = 0
const ROOK_WIZARD_INDEX = 1

func (wiz *Wizard) GenAttacks() {
	fmt.Println("generating attacks for", wiz.Name)
	maxShift := 0
	for sq := SquareMinValue; sq <= SquareMaxValue; sq++ {
		_, sqs := SlidingAttack(sq, wiz.Deltas, BbEmpty)
		//fmt.Println(bb)
		shift, magic, ok, nodes := SearchMagic(sq, sqs, wiz.Deltas, wiz.Tries)
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

var TotalMagicEntries int = 0

func (msq MagicSquare) Key(occup Bitboard) uint64 {
	return (msq.Magic * uint64(occup)) >> (64 - msq.Shift)
}

func (ms *Magics) Attack(sq Square, occup Bitboard) Bitboard {
	msq := ms[sq]
	return msq.Entries[msq.Key(occup)]
}

func BishopMobility(kind MoveKind, sq Square, occupUs, occupThem Bitboard) Bitboard {
	attack := Wizards[BISHOP_WIZARD_INDEX].Magics.Attack(sq, (occupUs|occupThem)&BishopAttack[sq])
	if !kind.IsViolent() {
		attack &^= occupThem
	}
	return attack &^ occupUs
}

func RookMobility(kind MoveKind, sq Square, occupUs, occupThem Bitboard) Bitboard {
	attack := Wizards[ROOK_WIZARD_INDEX].Magics.Attack(sq, (occupUs|occupThem)&RookAttack[sq])
	if !kind.IsViolent() {
		attack &^= occupThem
	}
	return attack &^ occupUs
}

func QueenMobility(kind MoveKind, sq Square, occupUs, occupThem Bitboard) Bitboard {
	return BishopMobility(kind, sq, occupUs, occupThem) | RookMobility(kind, sq, occupUs, occupThem)
}

func KnightMobility(kind MoveKind, sq Square, occupUs, occupThem Bitboard) Bitboard {
	attack := KnightAttack[sq]
	if !kind.IsViolent() {
		attack &^= occupThem
	}
	return attack &^ occupUs
}

func KingMobility(kind MoveKind, sq Square, occupUs, occupThem Bitboard) Bitboard {
	attack := KingAttack[sq]
	if !kind.IsViolent() {
		attack &^= occupThem
	}
	return attack &^ occupUs
}

var BishopAttack [BOARD_AREA]Bitboard
var RookAttack [BOARD_AREA]Bitboard
var QueenAttack [BOARD_AREA]Bitboard
var KnightAttack [BOARD_AREA]Bitboard
var KingAttack [BOARD_AREA]Bitboard

func JumpAttack(sq Square, deltas []Delta) Bitboard {
	bb := BbEmpty

	for _, delta := range deltas {
		rank := RankOf[sq] + delta.dRank
		file := FileOf[sq] + delta.dFile

		if rank >= 0 && rank < NUM_FILES && file >= 0 && file < NUM_FILES {
			bb |= RankFile[rank][file].Bitboard()
		}
	}

	return bb
}

func init() {
	/*for _, wiz := range Wizards {
		wiz.GenAttacks()
	}*/

	for wi := 0; wi < len(Wizards); wi++ {
		for i, msq := range Wizards[wi].Magics {
			size := 1 << msq.Shift
			Wizards[wi].Magics[i].Entries = make([]Bitboard, size)
			TotalMagicEntries += size
			_, sqs := SlidingAttack(msq.Square, Wizards[wi].Deltas, BbEmpty)
			var enum uint64
			for enum = 0; enum < 1<<len(sqs); enum++ {
				trb := Translate(sqs, enum)

				mobility, _ := SlidingAttack(msq.Square, Wizards[wi].Deltas, trb)
				key := msq.Key(trb)

				Wizards[wi].Magics[i].Entries[key] = mobility
			}

			if wi == BISHOP_WIZARD_INDEX {
				BishopAttack[msq.Square] = BishopMobility(Violent|Quiet, msq.Square, BbEmpty, BbEmpty)

				KnightAttack[msq.Square] = JumpAttack(msq.Square, KNIGHT_DELTAS)
				KingAttack[msq.Square] = JumpAttack(msq.Square, KING_DELTAS)
			} else if wi == ROOK_WIZARD_INDEX {
				RookAttack[msq.Square] = RookMobility(Violent|Quiet, msq.Square, BbEmpty, BbEmpty)
				QueenAttack[msq.Square] = BishopAttack[msq.Square] | RookAttack[msq.Square]
			}
		}
	}
}
