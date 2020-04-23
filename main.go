package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	. "github.com/easychessanimations/gobbit/basic"
)

func main() {
	st := State{}
	st.Init(VariantEightPiece)

	//fmt.Println(st)

	//fmt.Println(st.PrettyPrintString())

	/*fmt.Println(BISHOP_MAGICS)
	fmt.Println(ROOK_MAGICS)*/

	//fmt.Println(TotalMagicEntries)

	//fmt.Println(RookMobility(SquareD4, SquareD5.Bitboard()|SquareF4.Bitboard()|SquareD2.Bitboard()|SquareB4.Bitboard()))

	/*fmt.Println("rook attack b2")
	fmt.Println(RookAttack[SquareB2])
	fmt.Println("bishop attack e4")
	fmt.Println(BishopAttack[SquareE4])
	fmt.Println("queen attack f7")
	fmt.Println(QueenAttack[SquareF7])*/

	/*m := MakeMoveFT(SquareG1, SquareF3)

	fmt.Println(st.MoveLAN(m))*/

	//fmt.Println(st.Pslms(Violent | Quiet))

	pos := Position{}

	pos.Init(VariantEightPiece)

	//fmt.Println(pos.PrettyPrintString())

	//fmt.Println(pos.Current().GenerateMoves())

	//pos.Perf(6)

	//fmt.Println(PieceMaterialTablesString())

	pos.Print()

	scan := bufio.NewScanner(os.Stdin)

	for scan.Scan() {
		line := scan.Text()

		command := strings.TrimSpace(line)

		if command == "x" || command == "q" || command == "quit" {
			break
		} else if command == "h" || command == "help" {
			fmt.Println("h, help = help")
			fmt.Println("x, q, quit = quit")
			fmt.Println("pmt = print material table")
		} else if command == "pmt" {
			fmt.Println(PieceMaterialTablesString())
		} else {
			pos.ExecCommand(command)
		}
	}
}
