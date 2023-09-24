package repl

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/value"

	"github.com/antlr4-go/antlr/v4"
)

type Argument struct {
	Name            string
	Value           value.IVOR
	PassByReference bool
	Token           antlr.Token
	VariableRef     *Variable
}
