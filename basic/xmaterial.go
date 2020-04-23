package basic

import (
	"fmt"
	"strings"
)

type Score int16

type Accum struct {
	M Score
	E Score
}

func (acc Accum) String() string {
	return fmt.Sprintf("M %3d E %3d", acc.M, acc.E)
}

var PAWN_VALUE = Accum{100, 120}
var KNIGHT_VALUE = Accum{300, 300}
var BISHOP_VALUE = Accum{300, 320}
var ROOK_VALUE = Accum{500, 520}
var QUEEN_VALUE = Accum{900, 920}
var LANCER_VALUE = Accum{700, 720}
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

func init() {
	for color := Black; color <= White; color++ {
		for fig := FigureMinValue; fig <= FigureMaxValue; fig++ {
			p := ColorFigure[color][fig]
			mt := PieceMaterialTable{}
			switch fig {
			case Pawn:
				mt.Fill(PAWN_VALUE)
				PieceMaterialTables[p] = mt.POV(color)
				break
			case Knight:
				mt.Fill(KNIGHT_VALUE)
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
					PieceMaterialTables[p] = mt
				}
			}
		}
	}
}
