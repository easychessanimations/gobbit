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

func (st State) PslmsForPieceAtSquare() {

}
