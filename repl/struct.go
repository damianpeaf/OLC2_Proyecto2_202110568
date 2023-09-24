package repl

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"

	"github.com/antlr4-go/antlr/v4"
)

type Struct struct {
	Name   string
	Fields []compiler.IStruct_propContext
	Token  antlr.Token
}
