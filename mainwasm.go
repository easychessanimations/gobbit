// +build wasm

package main

import (		
	"syscall/js"
	"time"	
	. "github.com/easychessanimations/gobbit/uci"
)

var uci Uci

func ExecUciCommandLineWasm(this js.Value, commandLineArray []js.Value) interface{} {
	uci.ExecUciCommandLine(commandLineArray[0].String())

	return js.ValueOf("")
}

func main() {
	uci = Uci{}

	uci.Init(ENGINE_NAME, ENGINE_AUTHOR, UCI_COMMAND_ALIASES)

	uci.Welcome(" [ wasm build ]")

	js.Global().Set("ExecUciCommandLineWasm", js.FuncOf(ExecUciCommandLineWasm))

	for{
		time.Sleep(time.Second)
	}
}
