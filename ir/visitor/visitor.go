package visitor

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/abstract"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/tac"
)

type IrVisitor struct {
	compiler.BaseTSwiftLanguageVisitor
	Factory    *tac.TACFactory
	ScopeTrace *abstract.ScopeTrace
	Utility    *tac.Utility
}

func NewIrVisitor() *IrVisitor {
	factory := tac.NewTACFactory()
	scopeTrace := abstract.NewScopeTrace(factory)
	util := tac.NewUtility(factory)
	factory.Utility = util
	return &IrVisitor{
		Factory:    factory,
		ScopeTrace: scopeTrace,
		Utility:    util,
	}
}
