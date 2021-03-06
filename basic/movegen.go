package basic

import "fmt"

func FmtDummy(){
	fmt.Println("")
}

type Move uint32

type MoveType uint8

const (
	Normal = MoveType(iota)
	Promotion
	SentryPush
	Castling
	Null
)

// stack phases
const (
	PopNull = iota	
	GenAll
	PopAll	
	GenDone
)

func (st *State) InitStack(nmp bool, pvTable *PvHash, ignoreMoves []Move){	
	_, entry, ok := pvTable.Get(st.Zobrist)

	st.StackPvMoves = [MAX_PV_MOVES]Move{}

	if ok{
		st.StackPvMoves = entry.Moves
	}

	st.StackIgnoreMoves = ignoreMoves

	st.StackReduceDepth = 0

	st.StackPhase = GenAll

	if nmp{
		st.StackPhase = PopNull
	}
}

func (st *State) PopStackBuff() (StackBuffEntry, bool){
	l := len(st.StackBuff)

	if l <= 0{
		return StackBuffEntry{}, false
	}

	sbe := st.StackBuff[l-1]
	st.StackBuff = st.StackBuff[0:l-1]

	return sbe, true
}

const NullMove = Move(Null) << MOVE_TYPE_SHIFT

const MIN_REDUCE_LIMIT = 8

func (st *State) PopStack(pos *Position) Move{
	if st.StackPhase == PopNull{
		st.StackPhase = GenAll

		_, _, hasPvMove := pos.PvTable.Get(st.Zobrist)

		if !hasPvMove{			
			return NullMove
		}
	}

	if st.StackPhase == GenAll{
		st.SetStackBuff(pos, st.GenerateMoves())
		numAll := len(st.StackBuff)

		rF := 1
		for rF * rF < numAll{
			rF++
		}

		reduceLimit := rF / pos.PruningAgressivity

		if reduceLimit < MIN_REDUCE_LIMIT / pos.PruningAgressivity{
			reduceLimit = MIN_REDUCE_LIMIT / pos.PruningAgressivity
		}

		if reduceLimit < 1{
			reduceLimit = 1
		}

		st.StackReduceFrom = numAll - reduceLimit
		
		st.StackReduceFactor = rF

		st.StackPhase = PopAll
	}

	if st.StackPhase == PopAll{
		sbe, ok := st.PopStackBuff()
		if ok{
			st.StackReduceDepth = 0
			if sbe.SubTree > 0 && len(st.StackBuff) <= st.StackReduceFrom{
				st.StackReduceDepth = pos.PruningReduction
			}
			return sbe.Move
		}else{
			st.StackPhase = GenDone
		}
	}

	if st.StackPhase == GenDone{
		return Move(0)
	}

	panic("invalid stack phase")
}

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
	} else if move.MoveType() == SentryPush {
		buff += SymbolOf[FigureOf[move.PromotionPiece()]] + "@" + move.PromotionSquare().UCI()
	}

	return buff
}

func MakeMoveFT(fromSq, toSq Square) Move {
	return Move(fromSq + toSq<<TO_SQUARE_SHIFT)
}

func MakeMoveFTP(fromSq, toSq Square, pp Piece) Move {
	return Move(fromSq) + Move(toSq)<<TO_SQUARE_SHIFT + Move(pp)<<PROMOTION_PIECE_SHIFT + Move(Promotion)<<MOVE_TYPE_SHIFT
}

func MakeMoveFTPS(fromSq, toSq Square, pp Piece, ps Square) Move {
	return Move(fromSq) + Move(toSq)<<TO_SQUARE_SHIFT + Move(pp)<<PROMOTION_PIECE_SHIFT + Move(ps)<<PROMOTION_SQUARE_SHIFT + Move(SentryPush)<<MOVE_TYPE_SHIFT
}

func MakeMoveFTC(fromSq, toSq Square) Move {
	return Move(fromSq) + Move(toSq)<<TO_SQUARE_SHIFT + Move(Castling)<<MOVE_TYPE_SHIFT
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

func (st State) PromotionFigures() []Figure{
	promFigures := []Figure{Queen, Rook, Bishop, Knight}
	if st.Variant == VariantEightPiece{
		return append(promFigures, []Figure{LancerN, LancerNE, LancerE, LancerSE, LancerS, LancerSW, LancerW, LancerNW, Sentry, Jailer}...)
	}
	return promFigures
}

func (st State) AppendMove(moves *[]Move, move Move, jailColor Color){	
	p := st.PieceAtSquare(move.FromSq())

	fromFig := FigureOf[p]
	fromCol := ColorOf[p]

	if move.MoveType() == SentryPush{
		// for sentry push target square and promotion square have to differ
		if move.ToSq() == move.PromotionSquare(){
			return
		}
	}

	if st.HasDisabledMove && fromFig != Pawn{
		if move.FromSq() == st.DisableFromSquare{
			// quick check for exact match, that'll do it for jumping pieces
			if move.ToSq() == st.DisableToSquare{
				return
			}
			// check sliding pieces
			disabledBishopDir, hasDisabledBishopDir := NormalizedBishopDirection(st.DisableFromSquare, st.DisableToSquare)
			disabledRookDir, hasDisabledRookDir := NormalizedBishopDirection(st.DisableFromSquare, st.DisableToSquare)
			moveBishopDir, moveHasBishopDir := NormalizedBishopDirection(move.FromSq(), move.ToSq())
			moveRookDir, moveHasRookDir := NormalizedRookDirection(move.FromSq(), move.ToSq())
			if hasDisabledBishopDir && moveHasBishopDir{
				if disabledBishopDir == moveBishopDir{
					return
				}
			}
			if hasDisabledRookDir && moveHasRookDir{
				if disabledRookDir == moveRookDir{
					return
				}
			}
		}
	}

	if jailColor == NoColor || !st.IsSquareJailedForColor(move.FromSq(), jailColor){		
		if fromFig == Pawn && RankOf[move.ToSq()] == PromotionRank[fromCol]{
			for _, fig := range st.PromotionFigures(){
				*moves = append(*moves, MakeMoveFTP(move.FromSq(), move.ToSq(), ColorFigure[fromCol][fig]))
			}
		}else{
			*moves = append(*moves, move)
			return
		}
	}

	return
}

func (st State) GenBitboardMoves(sq Square, mobility Bitboard, jailColor Color) []Move {	
	moves := []Move{}

	for _, toSq := range mobility.PopAll() {
		st.AppendMove(&moves, MakeMoveFT(sq, toSq), jailColor)
	}

	return moves
}

func MakeLancer(color Color, ld int) Piece {
	return ColorFigure[color][LancerMinValue+Figure(ld)]
}

func (st State) GenLancerMoves(color Color, sq Square, mobility Bitboard, keepDir bool, lancerDir int, jailColor Color) []Move {
	moves := []Move{}

	for _, toSq := range mobility.PopAll() {
		if keepDir {
			st.AppendMove(&moves, MakeMoveFTP(sq, toSq, MakeLancer(color, lancerDir)), jailColor)
		} else {
			for ld := 0; ld < NUM_LANCER_DIRECTIONS; ld++ {
				st.AppendMove(&moves, MakeMoveFTP(sq, toSq, MakeLancer(color, ld)), jailColor)
			}
		}
	}

	return moves
}

func (st State) GenPawnMoves(kind MoveKind, color Color, sq Square, occupUs, occupThem Bitboard, jailColor Color, disablePushByTwo bool) []Move {
	pi := PawnInfos[sq][color]

	moves := []Move{}

	if kind&Violent != 0 {
		for _, captInfo := range pi.Captures {			
			if (captInfo.CheckSq.Bitboard() & occupThem) != 0 || captInfo.CheckSq == st.EpSquare {				
				st.AppendMove(&moves, captInfo.Move, jailColor)
			}
		}
	}

	if kind&Quiet != 0 {
		for _, pushInfo := range pi.Pushes {
			if (pushInfo.CheckSq.Bitboard() & (occupUs | occupThem)) == 0 {
				st.AppendMove(&moves, pushInfo.Move, jailColor)
				if disablePushByTwo{
					break
				}
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

func (st State) GenSentryMoves(kind MoveKind, color Color, sq Square, occupUs, occupThem Bitboard, jailColor Color) []Move{
	moves := []Move{}

	if jailColor != NoColor{
		if st.IsSquareJailedForColor(sq, color){
			return moves
		}
	}

	if kind & Quiet != 0{
		moves = append(moves, st.GenBitboardMoves(sq, BishopMobility(Quiet, sq, occupUs, occupThem), jailColor)...)
	}

	if kind & Violent != 0{
		// remove sentry for move generation
		sentry := st.PieceAtSquare(sq)
		st.Remove(sq)

		// remove sentry from occupation
		occupUs &^= sq.Bitboard()

		mob := BishopMobility(Violent, sq, occupUs, occupThem)		

		for _, pushSq := range mob.PopAll(){
			pushPiece := st.PieceAtSquare(pushSq)

			pushFig := FigureOf[pushPiece]
			pushCol := ColorOf[pushPiece]

			pushMoves := []Move{}

			if pushFig == Pawn{				
				// pawn handled separately, not to allow push by two
				pushMoves = append(pushMoves, st.GenPawnMoves(Violent|Quiet, color, pushSq, occupUs, occupThem, NoColor, true)...)
			}else if pushFig == Sentry{
				// pushed sentry should not push
				pushMoves = append(pushMoves, st.GenSentryMoves(Quiet, color, pushSq, occupUs, occupThem, NoColor)...)
			}else if pushPiece.IsLancer(){
				// lancer has special moves
				// normal moves keeping own direction
				lancerDir := pushPiece.LancerDirection()
				pushMoves = append(pushMoves, st.GenLancerMoves(pushCol, pushSq, LancerMobility(Violent|Quiet, lancerDir, pushSq, occupUs, occupThem), true, lancerDir, NoColor)...)
				// nudge to adjacent squares				
				for ld := 0; ld < NUM_LANCER_DIRECTIONS; ld++{
					if ld != lancerDir{
						targetSq, ok := st.AddDeltaToSquare(pushSq, LANCER_DELTAS[ld])												
						if ok && st.PieceAtSquare(targetSq) == NoPiece{
							moves = append(moves, MakeMoveFTPS(sq, pushSq, ColorFigure[pushCol][LancerN + Figure(ld)], targetSq))
						}
					}
				}
			}else{
				// all the rest
				pushMoves = append(pushMoves, st.PslmsForPieceAtSquare(Violent|Quiet, ColorFigure[color][pushFig], pushSq, occupUs, occupThem, NoColor)...)				
			}

			for _, pushMove := range pushMoves{
				moves = append(moves, MakeMoveFTPS(sq, pushSq, pushPiece, pushMove.ToSq()))
			}
		}

		// put back sentry
		st.Put(sentry, sq)
	}

	return moves
}

func (st State) CastlingTargetSquares(color Color, side int) [2]Square{
	cRank := st.CastlingRank(color)

	if side == CastlingSideKing{
		return [2]Square{RankFile[cRank][FileG], RankFile[cRank][FileF]}
	}else{
		return [2]Square{RankFile[cRank][FileC], RankFile[cRank][FileD]}	
	}
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
		moves := st.GenBitboardMoves(sq, KingMobility(kind, sq, occupUs, occupThem), jailColor)
		kCol := ColorOf[p]
		ccr := st.CastlingRights[kCol]
		wk := st.KingInfos[kCol].Square
		for side := CastlingSideKing; side <= CastlingSideQueen; side++{			
			cr := ccr[side]
			if cr.CanCastle && !st.IsSquareJailedForColor(cr.RookOrigSq, kCol){
				betweenOrigEmpty := true
				for _, testSq := range cr.BetweenOrigSquares{
					if st.PieceAtSquare(testSq) != NoPiece && testSq != wk && testSq != cr.RookOrigSq {
						betweenOrigEmpty = false
						break
					}
				}
				if betweenOrigEmpty{
					// remove king for tests
					st.Remove(wk)

					checksOk := true					

					cts := st.CastlingTargetSquares(kCol, side)

					kt := cts[0]

					fileDir := File(1)
					if FileOf[kt] < FileOf[wk]{
						fileDir = File(-1)
					}

					testRank := RankOf[wk]
					testFile := FileOf[wk]

					ok := true

					for ok {
						testSq := RankFile[testRank][testFile]
						
						st.Remove(testSq)
						st.Put(p, testSq)

						// now king is on test square
						if st.IsChecked(kCol){
							checksOk = false
						}

						st.Remove(testSq)

						if !checksOk{
							break
						}

						if testSq == kt{
							ok = false
						}else{
							testFile += fileDir
						}
					}

					if checksOk{
						moves = append(moves, MakeMoveFTC(wk, cr.RookOrigSq))
					}		

					// put back king
					st.Put(p, wk)		
				}
			}			
		}
		if st.IsSquareJailedForColor(wk, kCol){
			// jailed king has pass move			
			moves = append(moves, MakeMoveFT(wk, wk))
		}
		return moves
	case Pawn:
		return st.GenPawnMoves(kind, ColorOf[p], sq, occupUs, occupThem, jailColor, false)
	case LancerN, LancerNE, LancerE, LancerSE, LancerS, LancerSW, LancerW, LancerNW:
		moves := st.GenLancerMoves(ColorOf[p], sq, LancerMobility(kind, p.LancerDirection(), sq, occupUs, occupThem), false, 0, jailColor)
		// nudged lancer has special moves
		if st.HasDisabledMove{
			if st.DisableFromSquare == sq{
				for ld := 0; ld < NUM_LANCER_DIRECTIONS; ld++{
					if ld != p.LancerDirection(){
						moves = append(moves, st.GenLancerMoves(ColorOf[p], sq, LancerMobility(kind, ld, sq, occupUs, occupThem), true, ld, jailColor)...)
					}
				}				
			}
		}
		return moves
	case Sentry:
		return st.GenSentryMoves(kind, ColorOf[p], sq, occupUs, occupThem, jailColor)
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
	//return st.Pslms(Violent | Quiet)
	return append(st.Pslms(Violent), st.Pslms(Quiet)...)
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
