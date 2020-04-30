package basic

import (
	"fmt"
	"strings"
	"time"
)

type AlphaBetaInfo struct {
	Alpha        Score
	Beta         Score
	CurrentDepth int
	MaxDepth     int
}

var PvTable map[uint64]Move

type TranspositionTableEntry struct{
	Score    Score
	RemDepth int
	Zobrist  uint64
}

const TRANSPOSITION_TABLE_KEY_SIZE_IN_BITS = 26
const TRANSPOSITION_TABLE_SIZE = 1 << TRANSPOSITION_TABLE_KEY_SIZE_IN_BITS
const TRANSPOSITION_TABLE_KEY_MASK = TRANSPOSITION_TABLE_SIZE - 1

var TranspTable [TRANSPOSITION_TABLE_SIZE]TranspositionTableEntry

func TranspositionTableKey(zobrist uint64) uint64{
	return zobrist & TRANSPOSITION_TABLE_KEY_MASK
}

func TranspGet(zobrist uint64) (TranspositionTableEntry, bool){
	entry := TranspTable[TranspositionTableKey(zobrist)]
	if entry.Zobrist == zobrist{
		return entry, true
	}
	return entry, false
}

func TranspSet(zobrist uint64, entry TranspositionTableEntry){
	oldEntry, ok := TranspGet(zobrist)
	if ok{
		if entry.RemDepth > oldEntry.RemDepth{
			entry.Zobrist = zobrist
			TranspTable[TranspositionTableKey(zobrist)] = entry
		}
	}else{
		entry.Zobrist = zobrist
		TranspTable[TranspositionTableKey(zobrist)] = entry
	}
}

func (st State) Score() Score {
	mat := st.Material[White]
	mat.Merge(st.Material[Black])

	mf := float32(mat.M)

	phase := mf / float32(INITIAL_MATERIAL)

	mat = st.MaterialPOV()

	mf = float32(mat.M)
	ef := float32(mat.E)

	scoref := mf*phase + (1-phase)*ef

	score := Score(scoref)

	mob := st.MobilityPOV()

	score += mob.M + mob.E

	return score
}

const ALLOW_FAIL_SOFT = false

const ALLOW_TRANSPOSITION_TABLE = true

const TRANSPOSITION_TABLE_MAX_DEPTH = 5

func (pos *Position) AlphaBetaRec(abi AlphaBetaInfo) Score {
	pos.Nodes++

	st := pos.Current()

	if ALLOW_TRANSPOSITION_TABLE{
		entry, ok := TranspGet(st.Zobrist)

		if ok{
			if entry.RemDepth >= abi.MaxDepth - abi.CurrentDepth{
				return entry.Score
			}
		}
	}

	end, score := pos.GameEnd(abi.CurrentDepth)

	if end {
		// if game ended, return final score
		return score
	}

	if abi.CurrentDepth >= abi.MaxDepth || pos.SearchStopped {
		// if reached max depth or search stopped, return material score
		return st.Score()
	}

	ms := st.GenerateMoves()

	pvMove, ok := PvTable[pos.Zobrist()]

	if ok {
		newMs := []Move{pvMove}

		for _, testMove := range ms {
			if testMove != pvMove {
				newMs = append(newMs, testMove)
			}
		}

		ms = newMs
	}
	
	bestScore := -INFINITE_SCORE

	hasMove := false

	st.InitStack()

	for st.StackPhase != GenDone {
		move := st.PopStack()

		if st.StackPhase != GenDone{
			pos.Push(move)
		}

		if st.StackPhase == GenDone || pos.Current().IsCheckedThem(){
			if st.StackPhase != GenDone{
				pos.Pop()
			}			
		}else{
			hasMove = true

			score = -pos.AlphaBetaRec(AlphaBetaInfo{
				Alpha:        -abi.Beta,
				Beta:         -abi.Alpha,
				CurrentDepth: abi.CurrentDepth + 1,
				MaxDepth:     abi.MaxDepth,
			})

			if ALLOW_TRANSPOSITION_TABLE && abi.CurrentDepth < TRANSPOSITION_TABLE_MAX_DEPTH  {	
				TranspSet(pos.Zobrist(), TranspositionTableEntry{
					Score: -score,
					RemDepth: abi.MaxDepth - abi.CurrentDepth - 1,
				})
			}

			pos.Pop()

			if score >= abi.Beta {
				// beta cut
				return abi.Beta
			}
			
			if score > bestScore{
				bestScore = score
			}

			if score > abi.Alpha {
				// alpha improvement
				abi.Alpha = score
				PvTable[st.Zobrist] = move
			}
		}
	}

	if !hasMove{
		if pos.Current().IsCheckedUs(){
			return -MATE_SCORE + Score(abi.CurrentDepth)
		}else{
			return 0
		}
	}

	if ALLOW_FAIL_SOFT{
		return bestScore
	}else{
		return abi.Alpha
	}	
}

func (pos *Position) AlphaBeta(maxDepth int) Score {
	delete(PvTable, pos.Zobrist())

	pos.Nodes = 0

	return pos.AlphaBetaRec(AlphaBetaInfo{
		Alpha:        -INFINITE_SCORE,
		Beta:         INFINITE_SCORE,
		CurrentDepth: 0,
		MaxDepth:     maxDepth,
	})
}

func (pos Position) Zobrist() uint64 {
	return pos.Current().Zobrist
}

func (pos Position) GetPvRec(depthRemaining int, pvSoFar []Move) []Move {
	if depthRemaining <= 0 {
		return pvSoFar
	}

	move, ok := PvTable[pos.Zobrist()]

	if ok {
		pos.Push(move)
		return pos.GetPvRec(depthRemaining-1, append(pvSoFar, move))
	}

	return pvSoFar
}

func (pos Position) GetPv(maxDepth int) []Move {
	oldStatePtr := pos.StatePtr
	pv := pos.GetPvRec(maxDepth, []Move{})
	pos.StatePtr = oldStatePtr
	return pv
}

func (pos *Position) PrintBestMove(pv []Move) {
	if len(pv) == 0 {
		fmt.Println("bestmove (none)")
		return
	}

	if len(pv) == 1 {
		fmt.Println("bestmove", pv[0].UCI())
		return
	}

	fmt.Println("bestmove", pv[0].UCI(), "ponder", pv[1].UCI())
}

var lastGoodPv []Move

func PrintPvTable() {
	for zobrist, move := range PvTable {
		fmt.Printf("%016X %s\n", zobrist, move.UCI())
	}
}

func (pos *Position) Search(maxDepth int) {
	PvTable = make(map[uint64]Move)

	if ALLOW_TRANSPOSITION_TABLE{
		TranspTable = [TRANSPOSITION_TABLE_SIZE]TranspositionTableEntry{}
	}

	lastGoodPv = []Move{}

	pos.SearchStopped = false

	for depth := 1; depth <= maxDepth; depth++ {

		start := time.Now()

		score := pos.AlphaBeta(depth)

		if pos.SearchStopped {
			pos.PrintBestMove(lastGoodPv)
			return
		}

		lastGoodPv = pos.GetPv(depth)

		elapsed := float32(time.Now().Sub(start)) / 1e9

		nps := float32(pos.Nodes) / elapsed

		buff := []string{}

		for _, testMove := range lastGoodPv {
			buff = append(buff, testMove.UCI())
		}

		pv := strings.Join(buff, " ")

		fmt.Printf("info depth %d time %.0f nodes %d nps %.0f score cp %d pv %v\n", depth, elapsed*1000, pos.Nodes, nps, score, pv)
	}

	pos.PrintBestMove(lastGoodPv)
}


func init(){	
}