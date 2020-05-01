package uci

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	. "github.com/easychessanimations/gobbit/basic"
)

func Uci(){
	fmt.Printf("id name %s\n", ENGINE_NAME)
	fmt.Printf("id author %s\n\n", ENGINE_AUTHOR)

	for _, uo := range UCI_OPTIONS{
		fmt.Println(uo.UciCommandOutputString())
	}
}

func UciLoop(){
	fmt.Println(EngineId())

	pos := Position{}

	pos.Init(DEFAULT_VARIANT)

	//pos.Print()

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
		}else if command == "uci"{
			Uci()
		} else if command == "pmt" {
			fmt.Println(PieceMaterialTablesString())
		} else if command == "g" {
			go pos.Search(10)
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
