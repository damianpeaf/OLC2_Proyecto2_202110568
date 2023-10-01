package abstract

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/tac"
)

const (
	ExternNameParam = iota
	ExternEqualInnerParam
	PositionalParam
)

type Param struct {
	ExternName      string
	InnerName       string
	Type            string
	PassByReference bool
	Temp            *tac.Temp
}

// 3 types of paramantlr
// 1. extern and inner name declared 			[args are specified with the extern name]
// 2. extern name = "_", inner name declared 	[args are specified with the order]
// 3. extern name = "", inner name declared	[args are specified with the inner name]

func (p *Param) ParamType() int {
	if p.ExternName != "" {
		if p.ExternName == "_" {
			return PositionalParam
		} else {
			return ExternNameParam
		}
	} else {
		return ExternEqualInnerParam
	}
}
