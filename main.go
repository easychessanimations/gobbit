package main

import (
	. "github.com/easychessanimations/gobbit/basic"
	. "github.com/easychessanimations/gobbit/uci"
)

func main() {
	uci := Uci{}

	uci.Init(ENGINE_NAME, ENGINE_AUTHOR, UCI_OPTIONS, EngineId(), VariantEightPiece)

	uci.Welcome()

	uci.UciLoop()
}
