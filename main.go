// +build !wasm

package main

import (	
	"fmt"
	. "github.com/easychessanimations/gobbit/uci"
)

func main() {
	fmt.Println()

	uci := Uci{}

	uci.Init(ENGINE_NAME, ENGINE_AUTHOR, UCI_COMMAND_ALIASES)

	uci.Welcome(" [ native build ]")

	uci.ProcessConfig()

	uci.ProcessMatePuzzles()

	uci.ProcessCommandLine()

	uci.UciLoop()
}
