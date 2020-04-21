package basic

type Move uint32

type MoveType uint8

const (
	Normal = MoveType(iota)
)

const SQUARE_MASK = (1 << SQUARE_STORAGE_SIZE_IN_BITS) - 1

const PIECE_STORAGE_SIZE_IN_BITS = 6
const PIECE_MASK = (1 << PIECE_STORAGE_SIZE_IN_BITS) - 1

const MOVE_TYPE_STORAGE_SIZE_IN_BITS = 8
const MOVE_TYPE_MASK = (1 << MOVE_TYPE_STORAGE_SIZE_IN_BITS) - 1

// Move : [Move Type 8 bits][Promotion Piece 6 bits][Promotion Square 6 bits][To Square 6 bits][From Square 6 bits]

const FROM_SQUARE_SHIFT = 0                                                        // 0
const TO_SQUARE_SHIFT = SQUARE_STORAGE_SIZE_IN_BITS                                // 6
const PROMOTION_SQUARE_SHIFT = TO_SQUARE_SHIFT + SQUARE_STORAGE_SIZE_IN_BITS       // 12
const PROMOTION_PIECE_SHIFT = PROMOTION_SQUARE_SHIFT + SQUARE_STORAGE_SIZE_IN_BITS // 18
const MOVE_TYPE_SHIFT = PROMOTION_PIECE_SHIFT + PIECE_STORAGE_SIZE_IN_BITS         // 24

func (move Move) FromSq() Square {
	return Square(move & SQUARE_MASK)
}

func (move Move) ToSq() Square {
	return Square((move >> TO_SQUARE_SHIFT) & SQUARE_MASK)
}

func (move Move) PromotionSquare() Square {
	return Square((move >> PROMOTION_SQUARE_SHIFT) & SQUARE_MASK)
}

func (move Move) PromotionPiece() Piece {
	return Piece((move >> PROMOTION_PIECE_SHIFT) & PIECE_MASK)
}

func (move Move) MoveType() MoveType {
	return MoveType((move >> MOVE_TYPE_SHIFT) & MOVE_TYPE_MASK)
}

func (mt MoveType) String() string {
	switch mt {
	case Normal:
		return "Normal Move"
	}

	return "Unknown Type Move"
}

func (move Move) String() string {
	buff := move.FromSq().UCI() + "-" + move.ToSq().UCI()

	if move.MoveType() != Normal {
		buff += "@" + move.PromotionSquare().UCI() + "=" + move.PromotionPiece().FenSymbol()
	}

	return buff
}

func MakeMoveFT(fromSq, toSq Square) Move {
	return Move(fromSq + toSq<<TO_SQUARE_SHIFT)
}

type MoveKind int

const (
	NoMove = MoveKind(iota)
	Quiet
	Violent
)

func (mk MoveKind) IsQuiet() bool {
	return mk&Quiet != 0
}

func (mk MoveKind) IsViolent() bool {
	return mk&Violent != 0
}

func (st State) GenBitboardMoves(sq Square, mobility Bitboard) []Move {
	moves := []Move{}

	for toSq := mobility.Pop(); toSq != 0; toSq = mobility.Pop() {
		moves = append(moves, MakeMoveFT(sq, toSq))
	}

	return moves
}

func (st State) PslmsForPieceAtSquare(kind MoveKind, p Piece, sq Square, occupUs, occupThem Bitboard) []Move {
	switch FigureOf[p] {
	case Bishop:
		return st.GenBitboardMoves(sq, BishopMobility(kind, sq, occupUs, occupThem))
	case Rook:
		return st.GenBitboardMoves(sq, RookMobility(kind, sq, occupUs, occupThem))
	case Queen:
		return st.GenBitboardMoves(sq, QueenMobility(kind, sq, occupUs, occupThem))
	case Knight:
		return st.GenBitboardMoves(sq, KnightMobility(kind, sq, occupUs, occupThem))
	case King:
		return st.GenBitboardMoves(sq, KingMobility(kind, sq, occupUs, occupThem))
	}

	return []Move{}
}

func (st State) PslmsForColor(kind MoveKind, color Color) []Move {
	us := st.ByColor[color]
	them := st.ByColor[color.Inverse()]

	moves := []Move{}

	usbb := us

	for sq := usbb.Pop(); usbb != 0; sq = usbb.Pop() {
		p := st.PieceAtSquare(sq)

		moves = append(moves, st.PslmsForPieceAtSquare(kind, p, sq, us, them)...)
	}

	return moves
}

func (st State) Pslms(kind MoveKind) []Move {
	return st.PslmsForColor(kind, st.Turn)
}
