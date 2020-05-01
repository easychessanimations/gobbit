package uci


import (
	"fmt"
	"strings"
	"strconv"

	. "github.com/easychessanimations/gobbit/basic"
)

const DEFAULT_VARIANT = VariantEightPiece

const DEFAULT_DEPTH = 4

type UciOption struct{
	Name string
	Type string
	Default string
	Vars []string
	Value string
	Min int
	Max int
}

var UCI_OPTIONS = []UciOption{
	{
		Name: "UCI_Variant",
		Type: "combo",
		Vars: VARIANT_NAMES,
		Default: VARIANT_NAMES[DEFAULT_VARIANT],
	},
	{
		Name: "Null Move Pruning",
		Type: "check",		
		Default: "false",
	},
	{
		Name: "Null Move Pruning Min Depth",
		Type: "spin",		
		Min: 2,
		Max: 5,
		Default: "2",
	},
	{
		Name: "Null Move Depth Reduction",
		Type: "spin",		
		Min: 1,
		Max: 3,
		Default: "1",
	},
}

func (uo UciOption) StringValue() string{
	if uo.Value == ""{
		return uo.Default
	}

	return uo.Value
}

func (uo UciOption) BooleanValue() bool{
	return uo.StringValue() == "true"
}

func (uo UciOption) IntValue() int{
	value, err := strconv.ParseInt(uo.StringValue(), 10, 32)

	if err == nil{
		return int(value)
	}

	return 0
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

	if uo.Type == "spin"{
		buff += fmt.Sprintf(" min %d max %d", uo.Min, uo.Max)
	}

	return buff
}
