package basic

import "math/bits"

// Bitboard is a set representing the 8x8 chess board squares
type Bitboard uint64

// useful bitboards
const (
	BbEmpty              Bitboard = 0x0000000000000000
	BbFull               Bitboard = 0xffffffffffffffff
	BbBorder             Bitboard = 0xff818181818181ff
	BbPawnStartRank      Bitboard = 0x00ff00000000ff00
	BbPawnStartRankBlack Bitboard = 0x00ff000000000000
	BbPawnStartRankWhite Bitboard = 0x000000000000ff00
	BbPawnDoubleRank     Bitboard = 0x000000ffff000000
	BbBlackSquares       Bitboard = 0xaa55aa55aa55aa55
	BbWhiteSquares       Bitboard = 0x55aa55aa55aa55aa
)

const (
	BbFileA Bitboard = 0x101010101010101 << iota
	BbFileB
	BbFileC
	BbFileD
	BbFileE
	BbFileF
	BbFileG
	BbFileH
)

const (
	BbRank1 Bitboard = 0x0000000000000FF << (8 * iota)
	BbRank2
	BbRank3
	BbRank4
	BbRank5
	BbRank6
	BbRank7
	BbRank8
)

// RankBb returns a bitboard with all bits on rank set
func RankBb(rank int) Bitboard {
	return BbRank1 << uint(8*rank)
}

// FileBb returns a bitboard with all bits on file set
func FileBb(file int) Bitboard {
	return BbFileA << uint(file)
}

// North shifts all squares one rank up
func North(bb Bitboard) Bitboard {
	return bb << 8
}

// South shifts all squares one rank down
func South(bb Bitboard) Bitboard {
	return bb >> 8
}

// meaning of &^ operator
// https://stackoverflow.com/questions/34459450/what-is-the-operator-in-golang
// turn mask bits off

// East shifts all squares one file right
// delete h-file, then shift left
func East(bb Bitboard) Bitboard {
	return bb &^ BbFileH << 1
}

// West shifts all squares one file left
// delete a-file, then shift right
func West(bb Bitboard) Bitboard {
	return bb &^ BbFileA >> 1
}

// Fill returns a bitboard with all files with squares filled.
func Fill(bb Bitboard) Bitboard {
	return NorthFill(bb) | SouthFill(bb)
}

// ForwardSpan computes forward span wrt color.
func ForwardSpan(col Color, bb Bitboard) Bitboard {
	if col == White {
		return NorthSpan(bb)
	}
	if col == Black {
		return SouthSpan(bb)
	}
	return bb
}

// ForwardFill computes forward fill wrt color.
func ForwardFill(col Color, bb Bitboard) Bitboard {
	if col == White {
		return NorthFill(bb)
	}
	if col == Black {
		return SouthFill(bb)
	}
	return bb
}

// BackwardSpan computes backward span wrt color
func BackwardSpan(col Color, bb Bitboard) Bitboard {
	if col == White {
		return SouthSpan(bb)
	}
	if col == Black {
		return NorthSpan(bb)
	}
	return bb
}

// BackwardFill computes forward fill wrt color
func BackwardFill(col Color, bb Bitboard) Bitboard {
	if col == White {
		return SouthFill(bb)
	}
	if col == Black {
		return NorthFill(bb)
	}
	return bb
}

// NorthFill returns a bitboard with all north bits set
func NorthFill(bb Bitboard) Bitboard {
	bb |= (bb << 8)
	bb |= (bb << 16)
	bb |= (bb << 32)
	return bb
}

// NorthSpan is like NorthFill shifted on up
func NorthSpan(bb Bitboard) Bitboard {
	return NorthFill(North(bb))
}

// SouthFill returns a bitboard with all south bits set
func SouthFill(bb Bitboard) Bitboard {
	bb |= (bb >> 8)
	bb |= (bb >> 16)
	bb |= (bb >> 32)
	return bb
}

// SouthSpan is like SouthFill shifted on up
func SouthSpan(bb Bitboard) Bitboard {
	return SouthFill(South(bb))
}

// Has returns bb if sq is occupied in bitboard
func (bb Bitboard) Has(sq Square) bool {
	return bb>>sq&1 != 0
}

// AsSquare returns the occupied square if the bitboard has a single piece
// if the board has more then one piece the result is undefined
// https://golang.org/pkg/math/bits/#TrailingZeros64
func (bb Bitboard) AsSquare() Square {
	return Square(bits.TrailingZeros64(uint64(bb)) & 0x3f)
}

// LSB picks a square in the board
// returns empty board for empty board
func (bb Bitboard) LSB() Bitboard {
	return bb & (-bb)
}

// count returns the number of squares set in bb
// https://golang.org/pkg/math/bits/#OnesCount64
func (bb Bitboard) Count() int32 {
	return int32(bits.OnesCount64(uint64(bb)))
}

// Pop pops a set square from the bitboard
func (bb *Bitboard) Pop() Square {
	sq := *bb & (-*bb)
	*bb -= sq
	return Square(bits.TrailingZeros64(uint64(sq)) & 0x3f)
}

// String return the string representation of a bitboard
func (bb Bitboard) String() string {
	buff := "**********\n"

	for rank := LAST_RANK; rank >= 0; rank-- {
		buff += "*"
		for file := 0; file < NUM_FILES; file++ {
			sq := RankFile[rank][file]
			mask := sq.Bitboard()

			if bb&mask != 0 {
				buff += "1"
			} else {
				buff += "0"
			}
		}
		buff += "*"
		buff += "\n"
	}

	return buff + "**********\n"
}
