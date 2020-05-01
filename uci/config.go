package uci


import (
	"fmt"
	"strings"

	. "github.com/easychessanimations/gobbit/basic"
)

const DEFAULT_VARIANT = VariantEightPiece

const ENGINE_NAME = "gobbit"
const ENGINE_AUTHOR = "easychessanimations"

func EngineId() string{
	return fmt.Sprintf("%s multi variant uci engine by %s", ENGINE_NAME, ENGINE_AUTHOR)
}

var VARIANT_NAMES = []string{
	"Standard",
	"Eightpiece",
}

type UciOption struct{
	Name string
	Type string
	Default string
	Vars []string
	Value string
}

var UCI_OPTIONS = []UciOption{
	{
		Name: "UCI_Variant",
		Type: "combo",
		Vars: VARIANT_NAMES,
		Default: VARIANT_NAMES[DEFAULT_VARIANT],
	},
}

func (uo UciOption) UciCommandOutputString() string{
	buff := fmt.Sprintf("option name %s type %s default %s", uo.Name, uo.Type, uo.Default)

	vbuff := []string{}

	if uo.Type == "combo"{
		for _, v := range uo.Vars{
			vbuff = append(vbuff, fmt.Sprintf("var %s", v))
		}

		buff += " " + strings.Join(vbuff, " ")
	}

	return buff
}
