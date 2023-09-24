package visitor

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"
)

type IrVisitor struct {
	compiler.BaseTSwiftLanguageVisitor
}

func NewIrVisitor() *IrVisitor {
	return &IrVisitor{}
}
