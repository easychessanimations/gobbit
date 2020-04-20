package basic

// Delta is a square vector
type Delta struct {
	dRank int
	dFile int
}

// NormalizedBishopDirection tells the normalized bishop direction from fromSq to toSq
// returns true in the second return parameter if such a direction exists, false otherwise
func NormalizedBishopDirection(fromSq, toSq Square) (Delta, bool) {
	dRank := RankOf[toSq] - RankOf[fromSq]
	dFile := FileOf[toSq] - FileOf[fromSq]

	if dRank == 0 && dFile == 0 {
		// original and destination square are the same
		return Delta{}, false
	}

	if dRank*dRank != dFile*dFile {
		// not a bishop direction
		return Delta{}, false
	}

	if dRank > 0 && dFile > 0 {
		return Delta{1, 1}, true
	} else if dRank > 0 && dFile < 0 {
		return Delta{1, -1}, true
	} else if dFile > 0 {
		return Delta{-1, 1}, true
	} else {
		return Delta{-1, -1}, true
	}
}

// NormalizedRookDirection tells the normalized rook direction from fromSq to toSq
// returns true in the second return parameter if such a direction exists, false otherwise
func NormalizedRookDirection(fromSq, toSq Square) (Delta, bool) {
	dRank := RankOf[toSq] - RankOf[fromSq]
	dFile := FileOf[toSq] - FileOf[fromSq]

	if dRank == 0 && dFile == 0 {
		// original and destination square are the same
		return Delta{}, false
	}

	if dRank != 0 && dFile != 0 {
		// not a rook direction
		return Delta{}, false
	}

	if dRank > 0 {
		return Delta{1, 0}, true
	} else if dRank < 0 {
		return Delta{-1, 0}, true
	} else if dFile > 0 {
		return Delta{0, 1}, true
	} else {
		return Delta{0, -1}, true
	}
}

// Bitboard returns a bitboard that has sq set
func (sq Square) Bitboard() Bitboard {
	return 1 << sq
}

// RankFile constructs a square from rank and file
var RankFile [NUM_RANKS][NUM_FILES]Square

// FileOf tells the file of a square
var FileOf [BOARD_AREA]int

// RankOf tells the rank of a square
var RankOf [BOARD_AREA]int

// Rank tells the rank of a square
func (sq Square) Rank() int {
	return int((sq & Square(RANK_MASK)) >> RANK_SHIFT_IN_BITS)
}

// File tells the file of a square
func (sq Square) File() int {
	return int(sq & Square(FILE_MASK))
}

// UCI tells the UCI representation of a square
func (sq Square) UCI() string {
	/*buff := []byte{}

	buff = append(buff, byte(FileOf[sq]+'a'))
	buff = append(buff, byte(RankOf[sq]+'1'))

	return string(buff)*/

	return UCIOf[sq]
}

// String tells the string representation of a square, defaults to UCI
func (sq Square) String() string {
	return sq.UCI()
}

func init() {
	UCIToSquare = make(map[string]Square)
	for rank := 0; rank < NUM_RANKS; rank++ {
		for file := 0; file < NUM_FILES; file++ {
			sq := rank*NUM_FILES + file
			RankFile[rank][file] = Square(sq)
			FileOf[sq] = file
			RankOf[sq] = rank
			uci := FileLetterOf[file] + RankLetterOf[rank]
			UCIOf[sq] = uci
			UCIToSquare[uci] = Square(sq)
		}
	}
}
