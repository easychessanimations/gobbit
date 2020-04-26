package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	. "github.com/easychessanimations/gobbit/basic"
)

func main() {
	pos := Position{}

	//pos.Init(VariantEightPiece)

	//pos.Init(VariantStandard)
	
	pos.Init(VariantEightPiece)

	//pos.Current().ParseFen("k7/P7/8/K7/8/8/8/8 w KQkq - 0 1")

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
			fmt.Println("g = go depth 6")
			fmt.Println("s = stop")
			fmt.Println("d = del")
			fmt.Println("f = forward")
			fmt.Println("r = reset")
		} else if command == "pmt" {
			fmt.Println(PieceMaterialTablesString())
		} else if command == "g" {
			go pos.Search(6)
		} else if command == "s" {
			pos.SearchStopped = true
		} else if command == "r" {
			pos.StatePtr = 0
			pos.Print()
		} else {
			pos.ExecCommand(command)
		}
	}
}
