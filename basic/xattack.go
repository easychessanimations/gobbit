package basic

import "fmt"

func BishopAttackSquares(sq Square) []Square {
	sqs := []Square{}

	for testSq := SquareMinValue; testSq <= SquareMaxValue; testSq++ {
		_, ok := NormalizedBishopDirection(sq, testSq)
		if ok {
			sqs = append(sqs, testSq)
		}
	}

	return sqs
}

func init() {
	fmt.Println(BishopAttackSquares(SquareE4))
}
