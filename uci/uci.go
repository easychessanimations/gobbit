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
			}

			if name == "Null Move Pruning"{
				uci.Pos.NullMovePruning = uo.BooleanValue()
			}

			if name == "Null Move Pruning Min Depth"{
				uci.Pos.NullMovePruningMinDepth = uo.IntValue()
			}

			if name == "Null Move Depth Reduction"{
				uci.Pos.NullMoveDepthReduction = uo.IntValue()
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

func (uci *Uci) ExecPositionCommand(t *Tokenizer){
	token, ok := t.GetToken()

	if !ok{
		fmt.Println("missing position specifier")
		return
	}
	if token == "startpos" || token == "s"{
		uci.Pos.Reset()				
	}else if token == "fen" || token == "f"{
		fenParts := t.GetTokensUpTo("moves")
		if len(fenParts) < 4{
			fmt.Println("too few fen fields")
			return
		}
		uci.Pos.ParseFen(strings.Join(fenParts, " "))
	}else{
		fmt.Println("unknown position specifier")		
	}

	moves := t.GetTokensUpTo("")

	for _, move := range moves{
		uci.Pos.PushUci(move)
	}

	uci.Pos.Print()
}

func (uci *Uci) ExecGoCommand(t *Tokenizer){
	depth := DEFAULT_DEPTH

	for true{
		token, ok := t.GetToken()

		if !ok{
			break
		}

		if token == "depth"{
			parsedDepth := t.GetInt()

			if parsedDepth >= 1{
				depth = parsedDepth
			}
		}

		if token == "infinite"{
			depth = SEARCH_MAX_DEPTH
		}
	}

	go uci.Pos.Search(depth)
}

func (uci Uci) ListUciOptionValues(){
	for _, uo := range uci.UciOptions{
		fmt.Printf("%-30s = %s\n", uo.Name, uo.StringValue())
	}
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
		fmt.Println("l = list uci option values")
		fmt.Println("pmt = print material table")
		fmt.Println("g = go depth 10")
		fmt.Println("s = stop")
		fmt.Println("d = del")
		fmt.Println("f = forward")
		fmt.Println("b = to begin")
	}else if command == "uci"{
		uci.ExecUciCommand()
	}else if command == "position" || command == "p"{
		uci.ExecPositionCommand(&t)
	}else if command == "go"{
		uci.ExecGoCommand(&t)
	}else if command == "setoption"{
		uci.ExecSetOptionCommand(&t)
	}else if command == "l"{
		uci.ListUciOptionValues()
	} else if command == "pmt" {
		fmt.Println(PieceMaterialTablesString())
	} else if command == "g" {
		go uci.Pos.Search(10)
	} else if command == "s" || command == "stop" {
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

	for _, uo := range uci.UciOptions{
		uci.SetOption(uo.Name, uo.StringValue())
	}
}

func (uci Uci) Welcome(){
	fmt.Println(uci.Id())
}

func (uci *Uci) ProcessConfigLine(line string){
	if line[0:2] == "//"{
		fmt.Println("--", line[2:])
		return
	}else{
		fmt.Println("++", line)
	}
	uci.ExecUciCommandLine(line)	
}

func (uci *Uci) ProcessConfig(){
	IterateTextFile("engineconfig.txt", uci.ProcessConfigLine)
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
