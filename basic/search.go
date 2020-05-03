package basic

import (
	"fmt"
	"strings"
	"time"
)

type AlphaBetaInfo struct {
	Alpha         Score
	Beta          Score
	CurrentDepth  int
	MaxDepth      int
	NullMoveMade  bool
	NullMoveDepth int
}

var PvTable map[uint64][]Move

type PosMove struct{
	Zobrist uint64
	Move    Move
}

var PosMoveTable = make(map[PosMove]int)

type TranspositionTableEntry struct{
	Score    Score
	RemDepth int
	Zobrist  uint64
}

func (st State) Phase() float32{
	mat := st.Material[White]
	mat.Merge(st.Material[Black])

	return float32(mat.M) / float32(INITIAL_MATERIAL)
}

const MAX_SCORE = 9000

func (st State) Score() Score {
	phase := st.Phase()

	mat := st.MaterialPOV()

	mf := float32(mat.M)
	ef := float32(mat.E)

	scoref := mf * phase + ( 1 - phase ) * ef

	score := Score(scoref)

	mob := st.MobilityPOV()

	score += mob.M + mob.E

	if score > MAX_SCORE{
		score = MAX_SCORE
	}

	if score < -MAX_SCORE{
		score = -MAX_SCORE
	}

	return score
}

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
	allowNMP := pos.NullMovePruning && (!abi.NullMoveMade) && abi.CurrentDepth >= pos.NullMovePruningMinDepth

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
			if move != NullMove{
				hasMove = true
			}

			maxDepth := abi.MaxDepth
			nullMoveMade := abi.NullMoveMade
			nullMoveDepth := abi.NullMoveDepth

			if pos.NullMovePruning{
				if nullMoveMade{
					if move != NullMove && abi.CurrentDepth <= abi.NullMoveDepth{
						nullMoveMade = false
					}
				}

				if move == NullMove{
					nullMoveDepth = abi.CurrentDepth
					abi.NullMoveMade = true

					maxDepth -= pos.NullMoveDepthReduction
				}
			}			

			nodesStart := pos.Nodes

			stackReduceDepth := st.StackReduceDepth

			if !pos.StackReduction{
				stackReduceDepth = 0
			}

			score = -pos.AlphaBetaRec(AlphaBetaInfo{
				Alpha:         -abi.Beta,
				Beta:          -abi.Alpha,
				CurrentDepth:  abi.CurrentDepth + 1,
				MaxDepth:      maxDepth - stackReduceDepth,
				NullMoveMade:  nullMoveMade,
				NullMoveDepth: nullMoveDepth,
			})

			pos.Pop()

			subTree := pos.Nodes - nodesStart

			if stackReduceDepth > 0{
				for i := 0; i < stackReduceDepth; i++{
					subTree *= st.StackReduceFactor
				}				
			}

			if abi.CurrentDepth < 7{
				PosMoveTable[PosMove{st.Zobrist, move}] = subTree
			}

			if score > abi.Alpha {
				// alpha improvement
				abi.Alpha = score

				if move != NullMove{
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

	start := time.Now()

	for depth := 1; depth <= maxDepth; depth++ {

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

		totalPvTableMoves := 0
		maxPvItemLength := 0
		for _, item := range PvTable{
			l := len(item)
			totalPvTableMoves += l
			if l > maxPvItemLength{
				maxPvItemLength = l
			}
		}

		fmt.Printf("info pvtablesize %d pvtablemoves %d maxpvitemlength %d\n", len(PvTable), totalPvTableMoves, maxPvItemLength)
		fmt.Printf("info depth %d time %.0f nodes %d nps %.0f score cp %d pv %v\n", depth, elapsed*1000, pos.Nodes, nps, score, pv)
	}

	pos.PrintBestMove(lastGoodPv)
}


func init(){	
}