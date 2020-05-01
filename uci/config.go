package uci


import (
	"fmt"
	"strings"

	. "github.com/easychessanimations/gobbit/basic"
)

const DEFAULT_VARIANT = VariantEightPiece

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
