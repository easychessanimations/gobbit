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

var PvTable map[uint64][]Move

type TranspositionTableEntry struct{
	Score    Score
	RemDepth int
	Zobrist  uint64
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

const NULL_MOVE_PRUNING_MIN_DEPTH = 3

func (pos *Position) AlphaBetaRec(abi AlphaBetaInfo) Score {
	pos.Nodes++

	st := pos.Current()

	end, score := pos.GameEnd(abi.CurrentDepth)

	if end {
		// if game ended, return final score
		return score
	}

	if abi.CurrentDepth >= abi.MaxDepth || pos.SearchStopped {
		// if reached max depth or search stopped, return material score
		return st.Score()
	}

	hasMove := false

	// https://www.chessprogramming.org/Null_Move_Pruning
	allowNMP := pos.NullMovePruning && abi.CurrentDepth > NULL_MOVE_PRUNING_MIN_DEPTH

	st.InitStack(allowNMP)

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

			pos.Pop()

			if score > abi.Alpha {
				// alpha improvement
				abi.Alpha = score

				if move.MoveType() != Null{
					pvMoves, ok := PvTable[st.Zobrist]
					if ok{
						newPvMoves := []Move{move}
						for _, testMove := range pvMoves{
							if testMove != move{
								newPvMoves = append(newPvMoves, testMove)
							}
						}
						PvTable[st.Zobrist] = newPvMoves
					}else{
						PvTable[st.Zobrist] = []Move{move}
					}					
				}				
			}

			if score >= abi.Beta {
				// beta cut				
				return abi.Beta
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

	return abi.Alpha
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

	moves, ok := PvTable[pos.Zobrist()]

	if ok {
		pos.Push(moves[0])
		return pos.GetPvRec(depthRemaining-1, append(pvSoFar, moves[0]))
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
	for zobrist, moves := range PvTable {
		fmt.Printf("%016X %s\n", zobrist, moves[0].UCI())
	}
}

func (pos *Position) Search(maxDepth int) {
	PvTable = make(map[uint64][]Move)

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