package basic

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
	buff := []byte{}

	buff = append(buff, byte(FileOf[sq]+'a'))
	buff = append(buff, byte(RankOf[sq]+'1'))

	return string(buff)
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
