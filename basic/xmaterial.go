package basic

import (
	"fmt"
	"math/rand"
	"strings"
)

const INITIAL_MATERIAL = 2 * 4220

const CastleArraySize = 4

// the zobrist* arrays contain magic numbers used for Zobrist hashing
// more information on Zobrist hashing can be found in the paper:
// http://research.cs.wisc.edu/techreports/1970/TR88.pdf
var (
	zobristPiece     [PieceArraySize + 2][BOARD_AREA]uint64
	zobristEnpassant [BOARD_AREA]uint64
	zobristCastle    [CastleArraySize]uint64
	zobristColor     [ColorArraySize]uint64
)

// Zobrist returns the zobrist key of the position, never returns 0
func (st *State) GetZobrist() uint64 {
	if st.Zobrist != 0 {
		return st.Zobrist
	}
	return 0x4204fa763da3abeb
}

const (
	CastleWhiteKingIndex = iota
	CastleWhiteQueenIndex
	CastleBlackKingIndex
	CastleBlackQueenIndex
)

// SetCastlingAbility sets the side to move, correctly updating the Zobrist key
func (st *State) SetCastlingAbility(newCastlingRighs CastlingRights) {
	// unmark old castling rights if any
	if st.CastlingRights[White][CastlingSideKing].CanCastle{
		st.Zobrist ^= zobristCastle[CastleWhiteKingIndex]
	}	
	if st.CastlingRights[White][CastlingSideQueen].CanCastle{
		st.Zobrist ^= zobristCastle[CastleWhiteQueenIndex]
	}	
	if st.CastlingRights[Black][CastlingSideKing].CanCastle{
		st.Zobrist ^= zobristCastle[CastleBlackKingIndex]
	}	
	if st.CastlingRights[Black][CastlingSideQueen].CanCastle{
		st.Zobrist ^= zobristCastle[CastleBlackQueenIndex]
	}	

	// mark new castling rights
	if newCastlingRighs[White][CastlingSideKing].CanCastle{
		st.Zobrist ^= zobristCastle[CastleWhiteKingIndex]
	}	
	if newCastlingRighs[White][CastlingSideQueen].CanCastle{
		st.Zobrist ^= zobristCastle[CastleWhiteQueenIndex]
	}	
	if newCastlingRighs[Black][CastlingSideKing].CanCastle{
		st.Zobrist ^= zobristCastle[CastleBlackKingIndex]
	}	
	if newCastlingRighs[Black][CastlingSideQueen].CanCastle{
		st.Zobrist ^= zobristCastle[CastleBlackQueenIndex]
	}	

	st.CastlingRights = newCastlingRighs
}

// SetSideToMove sets the side to move, correctly updating the Zobrist key
func (st *State) SetSideToMove(color Color) {
	st.Zobrist ^= zobristColor[st.Turn]
	st.Turn = color
	st.Zobrist ^= zobristColor[st.Turn]
}

// SetEnpassantSquare sets the en passant square correctly updating the Zobrist key
func (st *State) SetEpSquare(epsq Square) {

	// in polyglot the hash key for en passant is updated only if
	// an en passant capture is possible next move; in other words
	// if there is an enemy pawn next to the end square of the move
	// TODO

	st.Zobrist ^= zobristEnpassant[st.EpSquare]
	st.EpSquare = epsq
	st.Zobrist ^= zobristEnpassant[st.EpSquare]
}

func init() {
	r := rand.New(rand.NewSource(5))
	f := func() uint64 { return uint64(r.Int63())<<32 ^ uint64(r.Int63()) }
	initZobristPiece(f)
	initZobristEnpassant(f)
	initZobristCastle(f)
	initZobristColor(f)
}

func initZobristPiece(f func() uint64) {
	for p := PieceMinValue; p <= PieceMaxValue; p++ {
		for sq := SquareMinValue; sq <= SquareMaxValue; sq++ {
			zobristPiece[p][sq] = f()
		}
	}
}

func initZobristEnpassant(f func() uint64) {
	for i := 0; i < 8; i++ {
		zobristEnpassant[SquareA3+Square(i)] = f()
		zobristEnpassant[SquareA6+Square(i)] = f()
	}
}

func initZobristCastle(f func() uint64) {
	r := [...]uint64{f(), f(), f(), f()}
	for i := 0; i < 4; i++ {
		zobristCastle[i] ^= r[i]
	}
}

func initZobristColor(f func() uint64) {
	zobristColor[White] = f()
}

type Score int16

type Accum struct {
	M Score
	E Score
}

func (acc *Accum) Merge(otherAcc Accum) {
	acc.M += otherAcc.M
	acc.E += otherAcc.E
}

func (acc *Accum) UnMerge(otherAcc Accum) {
	acc.M -= otherAcc.M
	acc.E -= otherAcc.E
}

func (acc Accum) Sub(otherAcc Accum) Accum {
	return Accum{acc.M - otherAcc.M, acc.E - otherAcc.E}
}

func (acc Accum) Mult(s Score) Accum {
	return Accum{acc.M * s, acc.E * s}
}

func (acc Accum) String() string {
	return fmt.Sprintf("M %3d E %3d", acc.M, acc.E)
}

var PAWN_VALUE = Accum{100, 120}
var CENTER_PAWN_VALUE = Accum{150, 120}
var SEMI_CENTER_PAWN_VALUE = Accum{125, 120}
var KNIGHT_VALUE = Accum{300, 300}
var KNIGHT_ON_EDGE_DEDUCTION = Accum{50, 50}
var KNIGHT_CLOSE_TO_EDGE_DEDUCTION = Accum{25, 25}
var BISHOP_VALUE = Accum{300, 320}
var ROOK_VALUE = Accum{500, 520}
var QUEEN_VALUE = Accum{900, 920}
var LANCER_VALUE = Accum{700, 720}
var LANCER_HOME_BONUS = Accum{100, 0}
var LANCER_FACING_OUT_VALUE = Accum{0, 0}
var SENTRY_VALUE = Accum{320, 320}
var JAILER_VALUE = Accum{400, 420}

type PieceMaterialTable [BOARD_AREA]Accum

func (pmt PieceMaterialTable) String() string {
	var rank Rank
	var file File
	rankBuff := []string{}
	for rank = LAST_RANK; rank >= 0; rank-- {
		lineBuff := []string{}
		for file = 0; file < LAST_FILE; file++ {
			sq := RankFile[rank][file]
			lineBuff = append(lineBuff, fmt.Sprintf("%s %v", sq.UCI(), pmt[sq].String()))
		}
		rankBuff = append(rankBuff, strings.Join(lineBuff, " , "))
	}
	return strings.Join(rankBuff, "\n") + "\n"
}

var PieceMaterialTables [PieceArraySize + 2]PieceMaterialTable

// POV returns material table from point of view of color
func (mt PieceMaterialTable) POV(color Color) PieceMaterialTable {
	if color == White {
		return mt
	}

	ret := PieceMaterialTable{}

	for rank := 0; rank < NUM_RANKS; rank++ {
		for file := 0; file < NUM_FILES; file++ {
			ret[rank*NUM_FILES+file] = mt[(LAST_RANK-rank)*NUM_FILES+file]
		}
	}

	return ret
}

func (mt *PieceMaterialTable) Fill(accum Accum) {
	for sq := SquareMinValue; sq <= SquareMaxValue; sq++ {
		mt[sq] = accum
	}
}

func PieceMaterialTablesString() string {
	items := []string{}
	for p := PieceMinValue; p <= PieceMaxValue; p++ {
		items = append(items, fmt.Sprintf("%s\n%s", Piece(p).FenSymbol(), PieceMaterialTables[p]))
	}
	return strings.Join(items, "\n")
}

const NUM_LANCER_DIRECTIONS = 8

func GetMaterialForPieceAtSquare(p Piece, sq Square) Accum {
	return PieceMaterialTables[p][sq]
}

const (
	LANCER_DIRECTION_N = iota
	LANCER_DIRECTION_NE
	LANCER_DIRECTION_E
	LANCER_DIRECTION_SE
	LANCER_DIRECTION_S
	LANCER_DIRECTION_SW
	LANCER_DIRECTION_W
	LANCER_DIRECTION_NW
)

func init() {
	var rank Rank
	var file File
	for color := Black; color <= White; color++ {
		for fig := FigureMinValue; fig <= FigureMaxValue; fig++ {
			p := ColorFigure[color][fig]
			mt := PieceMaterialTable{}
			switch fig {
			case Pawn:
				mt.Fill(PAWN_VALUE)
				mt[SquareE4] = CENTER_PAWN_VALUE
				mt[SquareE5] = CENTER_PAWN_VALUE
				mt[SquareD4] = CENTER_PAWN_VALUE
				mt[SquareD5] = CENTER_PAWN_VALUE
				mt[SquareC3] = SEMI_CENTER_PAWN_VALUE
				mt[SquareE3] = SEMI_CENTER_PAWN_VALUE
				PieceMaterialTables[p] = mt.POV(color)
				break
			case Knight:
				mt.Fill(KNIGHT_VALUE)				
				for rank = 0; rank < NUM_RANKS; rank++{
					mt[RankFile[rank][0]].UnMerge(KNIGHT_ON_EDGE_DEDUCTION)
					mt[RankFile[rank][LAST_FILE]].UnMerge(KNIGHT_ON_EDGE_DEDUCTION)
					if rank > 0 && rank < LAST_RANK{
						mt[RankFile[rank][1]].UnMerge(KNIGHT_CLOSE_TO_EDGE_DEDUCTION)
						mt[RankFile[rank][LAST_FILE-1]].UnMerge(KNIGHT_CLOSE_TO_EDGE_DEDUCTION)
					}
				}
				for file = 0; file < NUM_FILES; file++{
					mt[RankFile[0][file]].UnMerge(KNIGHT_ON_EDGE_DEDUCTION)
					mt[RankFile[LAST_RANK][file]].UnMerge(KNIGHT_ON_EDGE_DEDUCTION)
					if file > 0 && file < LAST_FILE{
						mt[RankFile[1][file]].UnMerge(KNIGHT_CLOSE_TO_EDGE_DEDUCTION)
						mt[RankFile[LAST_RANK-1][file]].UnMerge(KNIGHT_CLOSE_TO_EDGE_DEDUCTION)	
					}
				}
				PieceMaterialTables[p] = mt.POV(color)
				break
			case Bishop:
				mt.Fill(BISHOP_VALUE)
				PieceMaterialTables[p] = mt.POV(color)
				break
			case Rook:
				mt.Fill(ROOK_VALUE)
				PieceMaterialTables[p] = mt.POV(color)
				break
			case Queen:
				mt.Fill(QUEEN_VALUE)
				PieceMaterialTables[p] = mt.POV(color)
				break
			case Sentry:
				mt.Fill(SENTRY_VALUE)
				PieceMaterialTables[p] = mt.POV(color)
				break
			case Jailer:
				mt.Fill(JAILER_VALUE)
				PieceMaterialTables[p] = mt.POV(color)
				break
			default:
				// lancer
				for ld := 0; ld < NUM_LANCER_DIRECTIONS; ld++ {
					p = ColorFigure[color][int(LancerMinValue)+ld]
					mt.Fill(LANCER_VALUE)
					pstr := Rank2
					if color == Black{
						pstr = Rank7
					}		
					delta := LANCER_DELTAS[ld]			
					for file = 0; file < NUM_FILES; file++{
						if ( file < 6 && ld == LANCER_DIRECTION_E ) || ( file > 2 && ld == LANCER_DIRECTION_W ){
							mt[RankFile[pstr][file]].Merge(LANCER_HOME_BONUS)
						}						
						for rank = 0; rank < NUM_RANKS; rank++{
							if (file == 0 && delta.dFile < 0) || (file == LAST_FILE && delta.dFile > 0) || (rank == 0 && delta.dRank < 0) || (rank == LAST_RANK && delta.dRank > 0){
								mt[RankFile[rank][file]] = LANCER_FACING_OUT_VALUE
							}
						}
					}					
					PieceMaterialTables[p] = mt
				}
			}
		}
	}
}
