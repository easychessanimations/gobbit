package uci

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	. "github.com/easychessanimations/gobbit/basic"
)

type Uci struct{
	Name string
	Author string
	Id string
	UciOptions []UciOption
	Pos Position
}

func (uci Uci) Uci(){
	fmt.Printf("id name %s\n", uci.Name)
	fmt.Printf("id author %s\n\n", uci.Author)

	for _, uo := range uci.UciOptions{
		fmt.Println(uo.UciCommandOutputString())
	}
}

func (uci *Uci) ExecUciCommandLine(commandLine string) error{
	command := commandLine

	if command == "x" || command == "q" || command == "quit" {
		return fmt.Errorf("exit")
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
		uci.Uci()
	} else if command == "pmt" {
		fmt.Println(PieceMaterialTablesString())
	} else if command == "g" {
		go uci.Pos.Search(10)
	} else if command == "s" {
		uci.Pos.SearchStopped = true
	} else if command == "r" {
		uci.Pos.StatePtr = 0
		uci.Pos.Print()
	} else {
		uci.Pos.ExecCommand(command)
	}

	return nil
}

func (uci *Uci) Init(name string, author string, uciOptions []UciOption, id string, variant Variant){
	uci.Name = name
	uci.Author = author
	uci.UciOptions = uciOptions
	uci.Id = id

	uci.Pos = Position{}

	uci.Pos.Init(variant)
}

func (uci Uci) Welcome(){
	fmt.Println(uci.Id)
}

func (uci *Uci) UciLoop(){	
	scan := bufio.NewScanner(os.Stdin)

	for scan.Scan() {
		line := scan.Text()

		commandLine := strings.TrimSpace(line)

		err := uci.ExecUciCommandLine(commandLine)

		if err != nil{
			break
		}
	}
}
