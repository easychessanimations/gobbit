package basic

import (
	"fmt"
	"strings"
	"sort"
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

var VARIANT_NAMES = make([]string, len(VariantInfos))

func init(){
	for i, vinfo := range VariantInfos{
		VARIANT_NAMES[i] = vinfo.DisplayName
	}
}

func VariantNameToVariant(name string) Variant{
	for i, vinfo := range VariantInfos{
		if vinfo.DisplayName == name{
			return Variant(i)
		}
	}

	return VariantStandard
}

const (
	CastlingSideKing = iota
	CastlingSideQueen
)

type CastlingRight struct{
	CanCastle          bool
	RookOrigSq         Square
	RookOrigPiece      Piece
	BetweenOrigSquares []Square
}

type ColorCastlingRights [2]CastlingRight
type CastlingRights [2]ColorCastlingRights

func (ccr ColorCastlingRights) CanCastle() bool{
	return ccr[0].CanCastle || ccr[1].CanCastle
}

type MoveBuffItem struct {
	Move Move
	Uci  string
	Lan  string
	San  string
}

type MoveBuff []MoveBuffItem

func (mb MoveBuff) Len() int{
	return len(mb)
}

func (mb MoveBuff) Swap(i, j int){
	mb[i], mb[j] = mb[j], mb[i]
}

func (mb MoveBuff) Less(i, j int) bool{
	return strings.ToUpper(mb[i].San) < strings.ToUpper(mb[j].San)
}

type KingInfo struct {
	IsCaptured bool
	Square     Square
}

type StackBuffEntry struct{	
	Move      Move
	IsPv      bool
	PvIndex   int
	IsCapture bool
	Mobility  Accum
	SubTree   int		
}

type StackBuff []StackBuffEntry

func (sb StackBuff) Len() int{
	return len(sb)
}

func (sb StackBuff) Swap(i, j int){
	sb[i], sb[j] = sb[j], sb[i]
}

func (sb StackBuff) Less(i, j int) bool{	
	if sb[j].IsPv && (!sb[i].IsPv){
		return true
	}

	if (!sb[j].IsPv) && sb[i].IsPv{
		return false
	}

	if sb[j].IsPv && sb[i].IsPv{
		return sb[j].PvIndex < sb[i].PvIndex
	}

	if sb[j].IsCapture && (!sb[i].IsCapture){
		return true
	}

	if (!sb[j].IsCapture) && sb[i].IsCapture{
		return false
	}	

	if sb[j].Mobility.E > 0 && sb[i].Mobility.E == 0{
		return true
	}

	if sb[j].Mobility.E == 0 && sb[i].Mobility.E >= 0{
		return false
	}

	return sb[j].SubTree > sb[i].SubTree
}

func (st *State) SetStackBuff(pos *Position, moves []Move){
	st.StackBuff = []StackBuffEntry{}

	for _, move := range moves{
		posMove := PosMove{
			Zobrist: st.Zobrist,
			Move: move,			
		}

		subTree, _ := pos.PosMoveTable[posMove]		

		isPv := false
		pvIndex := 0

		for i, testMove := range st.StackPvMoves{
			if move == testMove{
				isPv = true
				pvIndex = i
				break
			}
		}

		st.StackBuff = append(st.StackBuff, StackBuffEntry{			
			Move: move, 
			IsPv: isPv,
			PvIndex: pvIndex,
			IsCapture: st.PieceAtSquare(move.ToSq()) != NoPiece,
			Mobility: st.MobilityForPieceAtSquare(st.PieceAtSquare(move.FromSq()), move.ToSq()),
			SubTree: subTree,
		})
	}	
	
	sort.Sort(st.StackBuff)
}

// State records the state of a position
type State struct {
	Variant               Variant
	Pieces                [NUM_RANKS][NUM_FILES]Piece
	Turn                  Color
	CastlingRights        CastlingRights
	EpSquare              Square
	HalfmoveClock         int
	FullmoveNumber        int
	HasDisabledMove       bool
	DisableFromSquare     Square
	DisableToSquare       Square
	ByFigure              [FigureArraySize]Bitboard
	ByLancer              Bitboard
	ByColor               [ColorArraySize]Bitboard
	Ply                   int
	Move                  Move
	MoveBuff              MoveBuff
	Material              [ColorArraySize]Accum
	Zobrist               uint64
	KingInfos             [ColorArraySize]KingInfo
	StackPhase            int
	StackBuff             StackBuff
	StackPvMoves          []Move
	StackReduceFrom       int
	StackReduceDepth      int
	StackReduceFactor     int
	LostCastlingForColor  [ColorArraySize]bool
}

func (st State) AddDeltaToSquare(sq Square, delta Delta) (Square, bool){
	rank := RankOf[sq]
	file := FileOf[sq]
	
	newRank := rank + delta.dRank
	newFile := file + delta.dFile

	if newRank >= 0 && newRank < NUM_RANKS && newFile >= 0 && newFile < NUM_FILES{
		return RankFile[newRank][newFile], true
	}

	return sq, false
}

// CastlingRights.String() reports castling rights in fen format
func (crs CastlingRights) String() string {
	buff := ""

	if crs[White][CastlingSideKing].CanCastle {
		buff += "K"
	}

	if crs[White][CastlingSideQueen].CanCastle {
		buff += "Q"
	}

	if crs[Black][CastlingSideKing].CanCastle {
		buff += "k"
	}

	if crs[Black][CastlingSideQueen].CanCastle {
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

// Reset resets the starting position
func (st *State) Reset(){
	st.Init(st.Variant)
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

func (st State) CastlingRank(color Color) Rank{
	if color == White{
		return Rank1
	}

	return Rank8
}

func (st State) IsCastlingPartner(fig Figure) bool{
	return fig == Rook || fig == Jailer
}

// PopulateCastlingRights determines castling rights information based on piece placement
// also does sanity check
func (st *State) PopulateCastlingRights(crs CastlingRights) CastlingRights{
	for color := Black; color <= White; color++{
		wk := st.KingInfos[color].Square		
		cRank := st.CastlingRank(color)
		if RankOf[wk] != cRank{
			// king is on illegal rank, delete castling rights
			crs[color][CastlingSideKing].CanCastle = false
			crs[color][CastlingSideQueen].CanCastle = false
		}else{
			for side := CastlingSideKing; side <= CastlingSideQueen; side++{
				dir := File(1 - (2 * side))
				foundCastlingPartner := false
				betweenSquares := []Square{}
				for testFile := FileOf[wk]; testFile >= 0 && testFile < NUM_FILES; testFile += dir{					
					testSq := RankFile[cRank][testFile]

					betweenSquares = append(betweenSquares, testSq)

					testP := st.PieceAtSquare(testSq)
					testCol := ColorOf[testP]
					testFig := FigureOf[testP]
					if st.IsCastlingPartner(testFig) && testCol == color{
						foundCastlingPartner = true
						crs[color][side].RookOrigSq = testSq
						crs[color][side].RookOrigPiece = testP						
						crs[color][side].BetweenOrigSquares = betweenSquares
					}					
				}
				if !foundCastlingPartner{
					crs[color][side].CanCastle = false
				}
			}
		}		
	}

	return crs
}

func (st *State) ParseCastlingRights(crs string) {
	t := Tokenizer{}
	t.Init(crs)

	st.LostCastlingForColor[White] = false
	st.LostCastlingForColor[Black] = false

	newCastlingRights := t.GetCastlingRights()

	populatedCastlingRights := st.PopulateCastlingRights(newCastlingRights)

	st.SetCastlingAbility(populatedCastlingRights)
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
		mbi.Uci = mbi.Move.UCI()
		st.MoveBuff[i] = mbi
	}

	sort.Sort(st.MoveBuff)
}

func (st *State) UciToMove(uci string) (Move, bool){
	st.GenMoveBuff()
	for _, mbi := range st.MoveBuff{
		if mbi.Uci == uci{
			return mbi.Move, true
		}
	}

	return Move(0), false
}

func (mb MoveBuff) PrettyPrintString() string {
	buff := ""

	cumul := 0

	for i, mbi := range mb {		
		buff += fmt.Sprintf("%-2d. %-12s", i+1, mbi.San)
		cumul++
		if cumul > 6 && i != len(mb)-1{
			cumul = 0
			buff += "\n"
		}		
	}

	return buff
}

const MOBILITY_MULTIPLIER = 5
const ATTACK_MULTIPLIER = 25

func (st State) MobilityBalance() Accum{
	return st.MobilityForColor(White).Sub(st.MobilityForColor(Black))
}

func (st State) MobilityPOV() Accum{
	mobBal := st.MobilityBalance()

	if st.Turn == White{
		return mobBal
	}

	return mobBal.Mult(-1)
}

func (st State) MobilityForPieceAtSquare(p Piece, sq Square) Accum{
	color := ColorOf[p]
	mobility := Accum{}

	occupUs := st.ByColor[color]
	occupThem := st.ByColor[color.Inverse()]

	if !st.IsSquareJailedForColor(sq, color){		
		fig := FigureOf[p]
		var mob Bitboard
		switch fig{
		case Knight:
			mob = KnightMobility(Violent|Quiet, sq, occupUs, occupThem)
			break
		// approximate sentry as bishop
		case Bishop, Sentry:
			mob = BishopMobility(Violent|Quiet, sq, occupUs, occupThem)
			break
		// approximate jailer as rook
		case Rook, Jailer:
			mob = RookMobility(Violent|Quiet, sq, occupUs, occupThem)
			break
		case Queen:
			mob = QueenMobility(Violent|Quiet, sq, occupUs, occupThem)
			break
		case LancerN, LancerNE, LancerE, LancerSE, LancerS, LancerSW, LancerW, LancerNW:
			mob = LancerMobility(Violent|Quiet, p.LancerDirection(), sq, occupUs, occupThem)
			break
		}			
		attack := mob & KingArea[sq]
		mobility.Merge(Accum{Score(mob.Count() * MOBILITY_MULTIPLIER), Score(attack.Count() * ATTACK_MULTIPLIER)})
	}		

	return mobility
}

func (st State) MobilityForColor(color Color) Accum{	
	mobility := Accum{}
	for _, sq := range st.ByColor[color].PopAll(){
		mobility.Merge(st.MobilityForPieceAtSquare(st.PieceAtSquare(sq), sq))		
	}
	return mobility
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

	buff += fmt.Sprintf("\n%s %s\n", VariantInfos[st.Variant].DisplayName, st.ReportFen())

	buff += fmt.Sprintf("\nMat White %v Black %v Balance %v POV %v Score %d\n", st.Material[White], st.Material[Black], st.Material[NoColor], st.MaterialPOV(), st.Score())

	mobW := st.MobilityForColor(White)
	mobB := st.MobilityForColor(Black)

	buff += fmt.Sprintf("Mob White %v Black %v Balance %v POV %v Phase %.2f\n", mobW, mobB, st.MobilityBalance(), st.MobilityPOV(), st.Phase())

	st.GenMoveBuff()

	buff += fmt.Sprintf("\n%s", st.MoveBuff.PrettyPrintString())

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
	return st.PieceAtSquare(move.ToSq()) != NoPiece || (FigureOf[st.PieceAtSquare(move.FromSq())] == Pawn && move.ToSq() == st.EpSquare)
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
	}else if move.MoveType() == SentryPush {
		buff += "=" + move.PromotionPiece().SanSymbol() + "@" + move.PromotionSquare().UCI()
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

			if p.IsLancer(){
				st.ByLancer |= bb
			}

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
		fromSq := mbi.Move.FromSq()						
		toSq := mbi.Move.ToSq()

		if st.PieceAtSquare(fromSq) == p && toSq == move.ToSq() && fromSq != move.FromSq() {
			ambig = true
			if FileOf[fromSq] == FileOf[move.FromSq()] {
				sameFile = true
			}
			if RankOf[fromSq] == RankOf[move.FromSq()] {
				sameRank = true
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

	if move.MoveType() == SentryPush {
		prom = "=" + move.PromotionPiece().SanSymbol() + "@" + move.PromotionSquare().UCI()
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

	if move.MoveType() == Castling{
		if FileOf[move.FromSq()] < FileOf[move.ToSq()]{
			return "O-O" + check
		}else{
			return "O-O-O" + check
		}
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

		if col == color.Inverse() && !st.IsSquareJailedForColor(sq, color.Inverse()) {
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
		if st.PieceAtSquare(captInfo.CheckSq) == ColorFigure[color.Inverse()][Pawn] && !st.IsSquareJailedForColor(captInfo.CheckSq, color.Inverse()) {
			return true
		}
	}

	na := KnightAttack[wk]

	// knight check
	themKnights := st.ByColor[color.Inverse()] & st.ByFigure[Knight] & na
	for _, sq := range themKnights.PopAll() {
		if !st.IsSquareJailedForColor(sq, color.Inverse()){
			return true
		}		
	}

	ka := KingAttack[wk]

	// king check
	themKings := st.ByColor[color.Inverse()] & st.ByFigure[King] & ka
	for _, sq := range themKings.PopAll() {
		if !st.IsSquareJailedForColor(sq, color.Inverse()){
			return true
		}		
	}

	// lancer check
	themLancers := st.ByColor[color.Inverse()] & st.ByLancer
	for _, sq := range themLancers.PopAll(){
		ms := st.PslmsForPieceAtSquare(Violent, st.PieceAtSquare(sq), sq, st.ByColor[color.Inverse()], st.ByColor[color], color.Inverse())
		for _, move := range ms{
			if move.ToSq() == wk{
				return true
			}
		}
	}

	// sentry check
	themSentries := st.ByColor[color.Inverse()] & st.ByFigure[Sentry]
	for _, sq := range themSentries.PopAll() {		
		sentryMoves := st.GenSentryMoves(Violent, color.Inverse(), sq, st.ByColor[color.Inverse()], st.ByColor[color], color)
		for _, sm := range sentryMoves{
			if sm.PromotionSquare() == wk{
				return true
			}
		}
	}

	return false
}

func (st State) IsCheckedUs() bool {
	return st.IsChecked(st.Turn)
}

func (st State) IsCheckedThem() bool {
	return st.IsChecked(st.Turn.Inverse())
}
