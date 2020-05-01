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
	UciOptions []UciOption
	Pos Position
	Aliases map[string]string
}

func (uci Uci) Id() string{
	return fmt.Sprintf("%s multi variant uci engine by %s", uci.Name, uci.Author)
}

func (uci Uci) ExecUciCommand(){
	fmt.Printf("id name %s\n", uci.Name)
	fmt.Printf("id author %s\n\n", uci.Author)

	for _, uo := range uci.UciOptions{
		fmt.Println(uo.UciCommandOutputString())
	}

	fmt.Println("uciok")
}

func (uci *Uci) SetVariant(variant Variant){
	uci.Pos = Position{}

	uci.Pos.Init(variant)
}

func (uci *Uci) SetOption(name, value string){
	for i, uo := range uci.UciOptions{
		if uo.Name == name{
			uo.Value = value
			uci.UciOptions[i] = uo

			if name == "UCI_Variant"{
				uci.SetVariant(VariantNameToVariant(value))

				uci.Pos.Print()
			}

			return
		}
	}

	fmt.Println("unknown option")
}

func (uci *Uci) ExecSetOptionCommand(t *Tokenizer){
	nameToken, ok := t.GetToken()

	if (!ok) || nameToken != "name"{
		fmt.Println("expected name")
		return
	}

	nameParts := t.GetTokensUpTo("value")

	if len(nameParts) == 0{
		fmt.Println("option name missing")
		return
	}

	name := strings.Join(nameParts, " ")

	value := t.Content

	uci.SetOption(name, value)
}

func (uci *Uci) ExecUciCommandLine(commandLine string) error{
	alias, ok := uci.Aliases[commandLine]

	if ok{
		commandLine = alias
		fmt.Println(commandLine)
	}

	t := Tokenizer{commandLine}

	command, ok := t.GetToken()

	if !ok{
		fmt.Println("no command")
		return nil
	}

	if command == "x" || command == "q" || command == "quit" {
		return fmt.Errorf("exit")
	} else if command == "h" || command == "help" {
		fmt.Println("h, help = help")
		fmt.Println("x, q, quit = quit")
		fmt.Println("pmt = print material table")
		fmt.Println("g = go depth 10")
		fmt.Println("s = stop")
		fmt.Println("d = del")
		fmt.Println("f = forward")
		fmt.Println("b = to begin")
	}else if command == "uci"{
		uci.ExecUciCommand()
	}else if command == "setoption"{
		uci.ExecSetOptionCommand(&t)
	} else if command == "pmt" {
		fmt.Println(PieceMaterialTablesString())
	} else if command == "g" {
		go uci.Pos.Search(10)
	} else if command == "s" {
		uci.Pos.SearchStopped = true
	} else if command == "b" {
		uci.Pos.StatePtr = 0
		uci.Pos.Print()
	} else {
		uci.Pos.ExecCommand(command)
	}

	return nil
}

func (uci *Uci) Init(name string, author string, aliases map[string]string){
	uci.Name = name
	uci.Author = author
	uci.UciOptions = UCI_OPTIONS
	uci.Aliases = aliases

	uci.SetVariant(DEFAULT_VARIANT)
}

func (uci Uci) Welcome(){
	fmt.Println(uci.Id())
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
