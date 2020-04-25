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

	pos.Init(VariantStandard)

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
		} else if command == "g" {
			go pos.Search(5)
		} else if command == "s" {
			pos.SearchStopped = true
		} else {
			pos.ExecCommand(command)
		}
	}
}
