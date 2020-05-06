package basic

import (
	"fmt"
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

const LOST_CASTLING_DEDUCTION = 150

func (st State) LostCastlingDeductionForColor(color Color, phase float32) Score{
	if st.LostCastlingForColor[color]{
		return Score(phase * float32(LOST_CASTLING_DEDUCTION))
	}	

	return 0
}

func (st State) LostCastlingDeductionBalance(phase float32) Score{	
	return st.LostCastlingDeductionForColor(White, phase) - st.LostCastlingDeductionForColor(Black, phase)
}

func (st State) LostCastlingDeductionPOV(phase float32) Score{	
	bal := st.LostCastlingDeductionBalance(phase)
	
	if st.Turn == White{
		return bal
	}

	return -bal
}

func (st State) Score() Score {
	phase := st.Phase()

	mat := st.MaterialPOV()

	mf := float32(mat.M)
	ef := float32(mat.E)

	scoref := mf * phase + ( 1 - phase ) * ef

	score := Score(scoref)

	mob := st.MobilityPOV()

	score += mob.M + mob.E

	score -= st.LostCastlingDeductionPOV(phase)

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

	currPvMove := NullMove

	for st.StackPhase != GenDone {
		move := st.PopStack(pos)

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
			
			pos.PosMoveTable.Set(st.Zobrist, move, PosMoveEntry{
				Depth: int8(abi.CurrentDepth),
				SubTree: subTree,
			})

			if score > abi.Alpha {
				// alpha improvement
				abi.Alpha = score

				if abi.CurrentDepth == 0 && score.IsMateInN(){
					// stop at forced mate
					if pos.Verbose{
						fmt.Printf("info skip %d\n", len(st.StackBuff))
					}					
					st.StackPhase = GenDone
				}

				if abi.CurrentDepth == 0{					
					pos.LastRootPvScore = score
					currPvMove = move
				}

				if move != NullMove{
					_, entry, ok := pos.PvTable.Get(st.Zobrist)
					pvMoves := entry.Moves
					if ok{
						newPvMoves := [MAX_PV_MOVES]Move{move}
						ptr := 1
						for _, testMove := range pvMoves{
							if testMove != move && ptr < MAX_PV_MOVES{
								newPvMoves[ptr] = testMove
								ptr++
							}
						}
						pos.PvTable.Set(st.Zobrist, PvEntry{
							Depth: int8(abi.CurrentDepth),
							Moves: newPvMoves,
						})
					}else{						
						pos.PvTable.Set(st.Zobrist, PvEntry{
							Depth: int8(abi.CurrentDepth),
							Moves: [MAX_PV_MOVES]Move{move},
						})
					}					
				}				
			}

			if pos.CheckTime() > 30{
				currPvUci := "none"
				if currPvMove != NullMove{
					currPvUci = currPvMove.UCI()
				}

				if pos.Verbose{
					fmt.Printf("info currstack %d currdepth %d time %d currpvmove %s latestrootscore cp %d oldpv %v\n", len(pos.SearchRoot().StackBuff), pos.Depth, pos.TimeMs(), currPvUci, pos.LastRootPvScore, pos.PvUCI())
				}				

				pos.CheckPoint = time.Now()
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
	pos.PvTable.Set(pos.Zobrist(), PvEntry{		
		Depth: INFINITE_DEPTH,
	})	

	pos.Nodes = 0

	pos.SearchRootPtr = pos.StatePtr

	// for low depth perform normal search
	if (!pos.AspirationWindow) || maxDepth < 4{
		return pos.AlphaBetaRec(AlphaBetaInfo{
			Alpha:        -INFINITE_SCORE,
			Beta:         INFINITE_SCORE,
			CurrentDepth: 0,
			MaxDepth:     maxDepth,
		})
	}

	windowLow := Score(50)
	windowHigh := Score(50)

	// for higher depths try aspiration window
	// https://www.chessprogramming.org/Aspiration_Windows
	for asp := 1; asp <= 3; asp++{
		alpha := pos.LastRootPvScore - windowLow
		beta := pos.LastRootPvScore + windowHigh

		if pos.Verbose{
			fmt.Printf("info asp %d windowlow %d windowhigh %d est %d alpha %d beta %d\n", asp, windowLow, windowHigh, pos.LastRootPvScore, alpha, beta)
		}		

		pos.CheckPoint = time.Now()

		score := pos.AlphaBetaRec(AlphaBetaInfo{
			Alpha:        alpha,
			Beta:         beta,
			CurrentDepth: 0,
			MaxDepth:     maxDepth,
		})

		if score > alpha && score < beta{
			return score
		}

		if score == alpha{
			if pos.Verbose{
				fmt.Println("info failed low")	
			}			

			windowLow *= 3
		}

		if score == beta{
			if pos.Verbose{
				fmt.Println("info failed high")
			}			

			windowHigh *= 3
		}
	}

	if pos.Verbose{
		fmt.Println("info asp full")
	}	

	pos.CheckPoint = time.Now()

	// aspiration window failed to return a pv, fall back to normal search
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

	_, entry, ok := pos.PvTable.Get(pos.Zobrist())

	moves := entry.Moves

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

var PvTable PvHash
var PosMoveTable PosMoveHash

func (pos *Position) Search(maxDepth int) {
	pos.PvTable = &PvTable
	pos.ClearPvTable()
	pos.PosMoveTable = &PosMoveTable
	pos.ClearPosMoveTable()

	pos.LastGoodPv = []Move{}

	pos.SearchStopped = false

	pos.Start = time.Now()
	pos.CheckPoint = pos.Start

	for pos.Depth = 1; pos.Depth <= maxDepth; pos.Depth++ {

		pos.AlphaBeta(pos.Depth)

		if pos.SearchStopped {
			pos.PrintBestMove(pos.LastGoodPv)
			return
		}

		pos.LastGoodPv = pos.GetPv(pos.Depth)

		pvTableSize := 0
		for _, item := range pos.PvTable.Entries{
			if item.Depth != INFINITE_DEPTH{
				pvTableSize++
			}			
		}

		if pos.Verbose {
			fmt.Printf("info pvtablesize %d\n", pvTableSize)
		}		
		fmt.Printf("info depth %d time %d nodes %d nps %.0f score cp %d pv %v\n", pos.Depth, pos.TimeMs(), pos.Nodes, pos.Nps(), pos.LastRootPvScore, pos.PvUCI())

		pos.CheckPoint = time.Now()

		if pos.IsMateInN(){
			break
		}
	}

	pos.PrintBestMove(pos.LastGoodPv)
}


func init(){	
}
