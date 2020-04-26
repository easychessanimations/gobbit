package basic

//import "fmt"

type Move uint32

type MoveType uint8

const (
	Normal = MoveType(iota)
	Promotion
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

func (move Move) UCI() string {
	buff := move.FromSq().UCI() + move.ToSq().UCI()

	if move.MoveType() == Promotion {
		buff += SymbolOf[FigureOf[move.PromotionPiece()]]
	}

	return buff
}

func MakeMoveFT(fromSq, toSq Square) Move {
	return Move(fromSq + toSq<<TO_SQUARE_SHIFT)
}

func MakeMoveFTP(fromSq, toSq Square, pp Piece) Move {
	return Move(fromSq) + Move(toSq)<<TO_SQUARE_SHIFT + Move(pp)<<PROMOTION_PIECE_SHIFT + Move(Promotion)<<MOVE_TYPE_SHIFT
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

func (st State) IsSquareJailedForColor(sq Square, color Color) bool{
	ja := JailerAdjacent[sq]

	if (st.ByFigure[Jailer] & st.ByColor[color.Inverse()] & ja) != 0{
		return true
	}

	return false
}

func (st State) AppendMove(moves []Move, move Move, jailColor Color) []Move{
	if jailColor == NoColor || !st.IsSquareJailedForColor(move.FromSq(), jailColor){
		return append(moves, move)
	}

	return moves
}

func (st State) GenBitboardMoves(sq Square, mobility Bitboard, jailColor Color) []Move {	
	moves := []Move{}

	for _, toSq := range mobility.PopAll() {
		moves = st.AppendMove(moves, MakeMoveFT(sq, toSq), jailColor)
	}

	return moves
}

func MakeLancer(color Color, ld int) Piece {
	return ColorFigure[color][LancerMinValue+Figure(ld)]
}

func (st State) GenLancerMoves(color Color, sq Square, mobility Bitboard, keepDir bool, jailColor Color) []Move {
	moves := []Move{}

	for _, toSq := range mobility.PopAll() {
		if keepDir {
			moves = st.AppendMove(moves, MakeMoveFTP(sq, toSq, st.PieceAtSquare(sq)), jailColor)
		} else {
			for ld := 0; ld < NUM_LANCER_DIRECTIONS; ld++ {
				moves = st.AppendMove(moves, MakeMoveFTP(sq, toSq, MakeLancer(color, ld)), jailColor)
			}
		}
	}

	return moves
}

func (st State) GenPawnMoves(kind MoveKind, color Color, sq Square, occupUs, occupThem Bitboard, jailColor Color) []Move {
	pi := PawnInfos[sq][color]

	moves := []Move{}

	if kind&Violent != 0 {
		for _, captInfo := range pi.Captures {
			if (captInfo.CheckSq.Bitboard() & occupThem) != 0 {
				moves = st.AppendMove(moves, captInfo.Move, jailColor)
			}
		}
	}

	if kind&Quiet != 0 {
		for _, pushInfo := range pi.Pushes {
			if (pushInfo.CheckSq.Bitboard() & (occupUs | occupThem)) == 0 {
				moves = st.AppendMove(moves, pushInfo.Move, jailColor)
			} else {
				break
			}
		}
	}

	return moves
}

func (l Piece) LancerDirection() int {
	return int(FigureOf[l]) & LANCER_DIRECTION_MASK
}

func (st State) PslmsForPieceAtSquare(kind MoveKind, p Piece, sq Square, occupUs, occupThem Bitboard, jailColor Color) []Move {
	switch FigureOf[p] {
	case Bishop:
		return st.GenBitboardMoves(sq, BishopMobility(kind, sq, occupUs, occupThem), jailColor)
	case Rook:
		return st.GenBitboardMoves(sq, RookMobility(kind, sq, occupUs, occupThem), jailColor)
	case Queen:
		return st.GenBitboardMoves(sq, QueenMobility(kind, sq, occupUs, occupThem), jailColor)
	case Knight:
		return st.GenBitboardMoves(sq, KnightMobility(kind, sq, occupUs, occupThem), jailColor)
	case King:
		return st.GenBitboardMoves(sq, KingMobility(kind, sq, occupUs, occupThem), jailColor)
	case Pawn:
		return st.GenPawnMoves(kind, ColorOf[p], sq, occupUs, occupThem, jailColor)
	case LancerN, LancerNE, LancerE, LancerSE, LancerS, LancerSW, LancerW, LancerNW:
		return st.GenLancerMoves(ColorOf[p], sq, LancerMobility(kind, p.LancerDirection(), sq, occupUs, occupThem), false, jailColor)
	case Sentry:
		return st.GenBitboardMoves(sq, BishopMobility(kind, sq, occupUs, occupThem), jailColor)
	case Jailer:
		return st.GenBitboardMoves(sq, RookMobility(kind&^Violent, sq, occupUs, occupThem), jailColor)
	}

	return []Move{}
}

func (st State) PslmsForColor(kind MoveKind, color Color) []Move {
	us := st.ByColor[color]
	them := st.ByColor[color.Inverse()]

	moves := []Move{}

	usbb := us

	for _, sq := range usbb.PopAll() {
		p := st.PieceAtSquare(sq)

		moves = append(moves, st.PslmsForPieceAtSquare(kind, p, sq, us, them, color)...)
	}

	return moves
}

func (st State) Pslms(kind MoveKind) []Move {
	return st.PslmsForColor(kind, st.Turn)
}

func (st State) GenerateMoves() []Move {
	return st.Pslms(Violent | Quiet)
}

func (st State) LegalMoves(stopAtFirst bool) []Move {
	lms := []Move{}

	for _, move := range st.GenerateMoves() {
		newSt := st
		newSt.MakeMove(move)
		if !newSt.IsCheckedThem() {
			lms = append(lms, move)
		}
		if stopAtFirst {
			if len(lms) > 0 {
				return lms
			}
		}
	}

	return lms
}

func (st State) HasLegalMove() bool {
	return len(st.LegalMoves(true)) > 0
}
