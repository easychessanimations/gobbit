package main

const ENGINE_NAME = "gobbit"
const ENGINE_AUTHOR = "easychessanimations"

var UCI_COMMAND_ALIASES = map[string]string{
	"vs" : "setoption name UCI_Variant value Standard",
	"ve" : "setoption name UCI_Variant value Eightpiece",
	"va" : "setoption name UCI_Variant value Atomic",
}
