package basic

import (
	"fmt"
	"strings"
)

type Variant int

func (v Variant) String() string {
	return VariantInfos[v].DisplayName
}

const (
	VariantStandard = Variant(iota)
	VariantEightPiece
)

type VariantInfo struct {
	StartFen    string
	DisplayName string
}

var VariantInfos = []VariantInfo{
	{ // standard
		StartFen:    "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		DisplayName: "Standard",
	},
	{ // eightpiece
		StartFen:    "jlsesqkbnr/pppppppp/8/8/8/8/PPPPPPPP/JLneSQKBNR w KQkq - 0 1 -",
		DisplayName: "Eightpiece",
	},
}

const (
	CastlingSideKing = iota
	CastlingSideQueen
)

type ColorCastlingRights [2]bool
type CastlingRights [2]ColorCastlingRights

type MoveBuffItem struct {
	Move Move
	Uci  string
	Lan  string
	San  string
}

type MoveBuff []MoveBuffItem

type KingInfo struct {
	IsCaptured bool
	Square     Square
}

// State records the state of a position
type State struct {
	Variant           Variant
	Pieces            [NUM_RANKS][NUM_FILES]Piece
	Turn              Color
	CastlingRights    CastlingRights
	EpSquare          Square
	HalfmoveClock     int
	FullmoveNumber    int
	HasDisabledMove   bool
	DisableFromSquare Square
	DisableToSquare   Square
	ByFigure          [FigureArraySize]Bitboard
	ByColor           [ColorArraySize]Bitboard
	Ply               int
	Move              Move
	MoveBuff          MoveBuff
	Material          [ColorArraySize]Accum
	Zobrist           uint64
	KingInfos         [ColorArraySize]KingInfo
}

// CastlingRights.String() reports castling rights in fen format
func (crs CastlingRights) String() string {
	buff := ""

	if crs[White][CastlingSideKing] {
		buff += "K"
	}

	if crs[White][CastlingSideQueen] {
		buff += "Q"
	}

	if crs[Black][CastlingSideKing] {
		buff += "k"
	}

	if crs[Black][CastlingSideQueen] {
		buff += "q"
	}

	if buff == "" {
		return "-"
	}

	return buff
}

// Color.String() converts a color to "b" for black, "w" for white and "-" for no color
func (color Color) String() string {
	if color == Black {
		return "b"
	}

	if color == White {
		return "w"
	}

	return "-"
}

// Init initializes state
// sets itself up from variant start fen
func (st *State) Init(variant Variant) {
	st.Variant = variant
	st.ParseFen(VariantInfos[st.Variant].StartFen)
}

// ParseFen sets up state from a fen
func (st *State) ParseFen(fen string) error {
	fenParts := strings.Split(fen, " ")

	if len(fenParts) > 0 {
		err := st.ParsePlacementString(fenParts[0])
		if err != nil {
			return err
		}
	}

	if len(fenParts) > 1 {
		st.ParseTurnString(fenParts[1])
	}

	if len(fenParts) > 2 {
		st.ParseCastlingRights(fenParts[2])
	}

	if len(fenParts) > 3 {
		st.ParseEpSquare(fenParts[3])
	}

	if len(fenParts) > 4 {
		st.ParseHalfmoveClock(fenParts[4])
	}

	if len(fenParts) > 5 {
		st.ParseFullmoveNumber(fenParts[5])
	}

	if len(fenParts) > 6 {
		st.ParseDisabledMove(fenParts[6])
	}

	return nil
}

func (st *State) ParseDisabledMove(dms string) {
	st.HasDisabledMove = false

	if len(dms) == 0 {
		return
	}

	if dms[:1] == "-" {
		return
	}

	t := Tokenizer{}
	t.Init(dms)

	st.DisableFromSquare = t.GetSquare()
	st.DisableToSquare = t.GetSquare()

	if st.DisableFromSquare != st.DisableToSquare {
		st.HasDisabledMove = true
	}
}

func (st *State) ParseEpSquare(epsqs string) {
	t := Tokenizer{}
	t.Init(epsqs)
	st.SetEpSquare(t.GetSquare())
}

func (st *State) ParseHalfmoveClock(hmcs string) {
	t := Tokenizer{}
	t.Init(hmcs)
	st.HalfmoveClock = t.GetInt()
}

func (st *State) ParseFullmoveNumber(fmns string) {
	t := Tokenizer{}
	t.Init(fmns)
	st.FullmoveNumber = t.GetInt()
}

func (st *State) ParseCastlingRights(crs string) {
	t := Tokenizer{}
	t.Init(crs)

	st.CastlingRights = t.GetCastlingRights()
	// TODO: Zobrist
}

func (st *State) ParseTurnString(ts string) {
	t := Tokenizer{}
	t.Init(ts)

	color := t.GetColor()

	st.SetSideToMove(White)

	if color != NoColor {
		st.SetSideToMove(color)
	}
}

func (st *State) GenMoveBuff() {
	st.MoveBuff = MoveBuff{}

	ms := st.LegalMoves(false)

	for _, move := range ms {
		st.MoveBuff = append(st.MoveBuff, MoveBuffItem{
			Move: move,
			Lan:  st.MoveLAN(move),
		})
	}

	for i, mbi := range st.MoveBuff {
		mbi.San = st.MoveToSanBatch(mbi.Move)
		st.MoveBuff[i] = mbi
	}
}

func (mb MoveBuff) PrettyPrintString() string {
	buff := []string{}

	for i, mbi := range mb {
		buff = append(buff, fmt.Sprintf("%d. %s", i+1, mbi.San))
	}

	return strings.Join(buff, " ")
}

func (st State) MaterialPOV() Accum {
	if st.Turn == White {
		return st.Material[NoColor]
	}
	return st.Material[NoColor].Mult(-1)
}

// PrettyPrintString returns the state pretty print string
func (st State) PrettyPrintString() string {
	buff := st.PrettyPlacementString()

	buff += fmt.Sprintf("\n%s : %s : %16X\n", VariantInfos[st.Variant].DisplayName, st.ReportFen(), st.Zobrist)

	buff += fmt.Sprintf("\nWhite %v , Black %v , Balance %v , POV %v , Score %d\n", st.Material[White], st.Material[Black], st.Material[NoColor], st.MaterialPOV(), st.Score())

	st.GenMoveBuff()

	buff += fmt.Sprintf("\nLegal moves ( %d ) : %s", len(st.MoveBuff), st.MoveBuff.PrettyPrintString())

	return buff
}

// ReportFen reports the state as a fen string
func (st State) ReportFen() string {
	buff := ""

	cum := 0

	for rank := LAST_RANK; rank >= 0; rank-- {
		for file := 0; file < NUM_FILES; file++ {
			p := st.Pieces[rank][file]

			if p == NoPiece {
				cum++
			} else {
				if cum > 0 {
					buff += fmt.Sprintf("%d", cum)
				}
				cum = 0
				buff += p.FenSymbol()
			}
		}
		if cum > 0 {
			buff += fmt.Sprintf("%d", cum)
		}
		cum = 0
		if rank > 0 {
			buff += "/"
		}
	}

	buff += " " + st.Turn.String()

	buff += " " + st.CastlingRights.String()

	if st.EpSquare == SquareA1 {
		buff += " -"
	} else {
		buff += " " + st.EpSquare.UCI()
	}

	buff += " " + fmt.Sprintf("%d", st.HalfmoveClock)

	buff += " " + fmt.Sprintf("%d", st.FullmoveNumber)

	if st.Variant == VariantEightPiece {
		if st.HasDisabledMove {
			buff += " " + st.DisableFromSquare.UCI() + st.DisableToSquare.UCI()
		} else {
			buff += " -"
		}

	}

	return buff
}

// PrettyPlacementString returns the pretty string representation of the board
func (st State) PrettyPlacementString() string {
	buff := ""

	for rank := LAST_RANK; rank >= 0; rank-- {
		for file := 0; file < NUM_FILES; file++ {
			buff += st.Pieces[rank][file].PrettySymbol()
		}
		if (rank == 0 && st.Turn == White) || (rank == LAST_RANK && st.Turn == Black) {
			buff += " *"
		}
		buff += "\n"
	}

	return buff
}

// ParsePlacementString parses a placement string and sets pieces accordingly
// returns an error if there are not enough pieces to fill the board
// liberal otherwise
func (st *State) ParsePlacementString(ps string) error {
	t := Tokenizer{}
	t.Init(ps)

	rank := LAST_RANK
	file := 0

	st.Zobrist = 0

	for color := Black; color <= White; color++ {
		st.KingInfos[color] = KingInfo{
			IsCaptured: true,
			Square:     SquareA1,
		}
	}

	for ps := t.GetFenPiece(); len(ps) > 0; {
		if len(ps) > 0 {
			for _, p := range ps {
				sq := RankFile[rank][file]
				st.Pieces[rank][file] = NoPiece
				st.Put(p, sq)
				file++
				if file > LAST_FILE {
					file = 0
					rank--
				}
				if rank < 0 {
					st.CalculateOccupancyAndMaterial()

					return nil
				}
			}

			ps = t.GetFenPiece()
		}
	}

	return fmt.Errorf("too few pieces in placement string")
}

func (st State) PieceAtSquare(sq Square) Piece {
	return st.Pieces[RankOf[sq]][FileOf[sq]]
}

func (st State) IsCapture(move Move) bool {
	return st.PieceAtSquare(move.ToSq()) != NoPiece
}

func (st State) MoveLAN(move Move) string {
	fromPiece := st.PieceAtSquare(move.FromSq())

	buff := fromPiece.SanLetter() + move.FromSq().UCI()

	if st.IsCapture(move) {
		buff += "x"
	} else {
		buff += "-"
	}

	buff += move.ToSq().UCI()

	if move.MoveType() == Promotion {
		buff += "=" + move.PromotionPiece().SanSymbol()
	}

	return buff
}

func (st *State) CalculateOccupancyAndMaterial() {
	st.ByFigure = [FigureArraySize]Bitboard{}
	st.ByColor = [ColorArraySize]Bitboard{}

	st.Material[White] = Accum{}
	st.Material[Black] = Accum{}

	for sq := SquareMinValue; sq <= SquareMaxValue; sq++ {
		bb := sq.Bitboard()

		p := st.PieceAtSquare(sq)

		if p != NoPiece {
			fig := FigureOf[p]

			col := ColorOf[p]

			st.ByFigure[fig] |= bb
			st.ByColor[col] |= bb

			p := ColorFigure[col][fig]

			mat := GetMaterialForPieceAtSquare(p, sq)

			st.Material[col].Merge(mat)
		}
	}

	st.Material[NoColor] = st.Material[White].Sub(st.Material[Black])
}

func (st State) OccupUs() Bitboard {
	return st.ByColor[st.Turn]
}

func (st State) OccupThem() Bitboard {
	return st.ByColor[st.Turn.Inverse()]
}

func (st *State) MoveToSan(move Move) string {
	st.GenMoveBuff()

	return st.MoveToSanBatch(move)
}

func (st State) MoveToSanBatch(move Move) string {
	p := st.PieceAtSquare(move.FromSq())

	sanLetter := p.SanLetter()

	orig := move.FromSq().UCI()

	sameRank := false
	sameFile := false
	ambig := false

	for _, mbi := range st.MoveBuff {
		if mbi.Move != move {
			fromSq := mbi.Move.FromSq()
			if st.PieceAtSquare(fromSq) == p && mbi.Move.ToSq() == move.ToSq() {
				ambig = true
				if FileOf[fromSq] == FileOf[move.FromSq()] {
					sameFile = true
				}
				if RankOf[fromSq] == RankOf[move.FromSq()] {
					sameRank = true
				}
			}
		}
	}

	if FigureOf[p] == Pawn {
		orig = orig[0:1]
		sanLetter = ""
	} else if ambig {
		if sameRank && sameFile {
			// do nothing, orig already has both rank and file
		} else if sameFile {
			// differentiate by rank
			orig = orig[1:2]
		} else {
			// default is differentiate by file
			orig = orig[0:1]
		}
	} else {
		orig = ""
	}

	dest := move.ToSq().UCI()

	takes := ""

	if st.IsCapture(move) {
		takes = "x"
	} else {
		if FigureOf[p] == Pawn {
			orig = ""
		}
	}

	prom := ""

	if move.MoveType() == Promotion {
		prom = "=" + move.PromotionPiece().SanSymbol()
	}

	check := ""

	newSt := st

	newSt.MakeMove(move)

	if newSt.IsCheckedUs() {
		check = "+"

		if !newSt.HasLegalMove() {
			check = "#"
		}
	} else if !newSt.HasLegalMove() {
		check = "="
	}

	return sanLetter + orig + takes + dest + prom + check
}

func (st State) IsChecked(color Color) bool {
	if st.KingInfos[color].IsCaptured {
		return true
	}

	wk := st.KingInfos[color].Square

	qm := QueenMobility(Violent, wk, st.ByColor[color], st.ByColor[color.Inverse()])

	// bishop, rook, queen
	for _, sq := range qm.PopAll() {
		p := st.PieceAtSquare(sq)

		col := ColorOf[p]
		fig := FigureOf[p]

		if col == color.Inverse() {
			if fig == Queen {
				return true
			}
			_, isBishopDirection := NormalizedBishopDirection(wk, sq)
			if fig == Bishop && isBishopDirection{
				return true
			}
			_, isRookDirection := NormalizedRookDirection(wk, sq)
			if fig == Rook && isRookDirection{
				return true
			}
		}
	}

	pi := PawnInfos[wk][color]

	// pawn check
	for _, captInfo := range pi.Captures {
		if st.PieceAtSquare(captInfo.CheckSq) == ColorFigure[color.Inverse()][Pawn] {
			return true
		}
	}

	na := KnightAttack[wk]

	// knight check
	if (st.ByColor[color.Inverse()] & st.ByFigure[Knight] & na) != 0 {
		return true
	}

	ka := KingAttack[wk]

	// king check
	if (st.ByColor[color.Inverse()] & st.ByFigure[King] & ka) != 0 {
		return true
	}

	return false
}

func (st State) IsCheckedUs() bool {
	return st.IsChecked(st.Turn)
}

func (st State) IsCheckedThem() bool {
	return st.IsChecked(st.Turn.Inverse())
}
