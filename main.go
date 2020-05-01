package main

import (
	. "github.com/easychessanimations/gobbit/basic"
	. "github.com/easychessanimations/gobbit/uci"
)

const ENGINE_NAME = "gobbit"
const ENGINE_AUTHOR = "easychessanimations"

var UCI_COMMAND_ALIASES = map[string]string{
	"vs" : "setoption name UCI_Variant value Standard",
	"ve" : "setoption name UCI_Variant value Eightpiece",
}

func main() {
	uci := Uci{}

	uci.Init(ENGINE_NAME, ENGINE_AUTHOR, VariantEightPiece, UCI_COMMAND_ALIASES)

	uci.Welcome()

	uci.UciLoop()
}
