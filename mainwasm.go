// +build wasm

package main

import (		
	"syscall/js"
	"time"	
	. "github.com/easychessanimations/gobbit/uci"
)

var uci Uci

var UciLoopStopped bool

func ExecUciCommandLineWasm(this js.Value, commandLineArray []js.Value) interface{} {
	err := uci.ExecUciCommandLine(commandLineArray[0].String())
	if err != nil{
		UciLoopStopped = true
	}
	return js.ValueOf("")
}

func main() {
	uci = Uci{}

	uci.Init(ENGINE_NAME, ENGINE_AUTHOR, UCI_COMMAND_ALIASES)

	uci.Welcome(" [ wasm build ]")

	js.Global().Set("ExecUciCommandLineWasm", js.FuncOf(ExecUciCommandLineWasm))

	UciLoopStopped = false

	for !UciLoopStopped{
		time.Sleep(1 * time.Second)
	}
}
