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
	Strats     map[string]BinaryStrategy
	UnStrats   map[string]UnaryStrategy
	Transfer   *TransferLabels
}

type TransferLabels struct {
	ReturnLabel   *tac.Label
	BreakLabel    *tac.Label
	ContinueLabel *tac.Label
}

func NewIrVisitor() *IrVisitor {
	factory := tac.NewTACFactory()
	util := tac.NewUtility(factory)
	factory.Utility = util

	transfer := &TransferLabels{
		ReturnLabel:   nil,
		BreakLabel:    nil,
		ContinueLabel: nil,
	}

	visitor := &IrVisitor{
		Factory:    factory,
		ScopeTrace: nil,
		Utility:    util,
		Transfer:   transfer,
	}

	strats := NewBinaryStrats(visitor)
	visitor.Strats = strats

	unStrats := NewUnaryStrats(visitor)
	visitor.UnStrats = unStrats

	return visitor
}
