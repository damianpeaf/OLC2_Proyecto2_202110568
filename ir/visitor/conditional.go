package visitor

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/tac"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/value"
)

type ifStmt struct {
	Condition      *tac.ConditionalJump
	TrueLabel      *tac.Label
	ConditionBlock tac.TACBlock
	InnerBlock     tac.TACBlock
}

func (v *IrVisitor) VisitIfStmt(ctx *compiler.IfStmtContext) interface{} {

	v.Factory.AppendToBlock(v.Factory.NewComment().SetComment("if chain"))

	ifChain := make([]*ifStmt, 0)

	for _, child := range ctx.AllIf_chain() {
		ifChain = append(ifChain, v.Visit(child).(*ifStmt))
	}

	elseBlock := make([]tac.TACStmtI, 0)

	if ctx.Else_stmt() != nil {
		elseBlock = v.Visit(ctx.Else_stmt()).([]tac.TACStmtI)
	}

	finalLabel := v.Factory.NewLabel()
	// add chain to block
	for _, stmt := range ifChain {

		nextLabel := v.Factory.NewLabel()

		// add condition block
		v.Factory.AppendBlock(stmt.ConditionBlock)

		// add condition
		v.Factory.AppendToBlock(v.Factory.NewComment().SetComment("if condition"))
		v.Factory.AppendToBlock(stmt.Condition)
		// jmp to next condition
		v.Factory.AppendToBlock(v.Factory.NewUnconditionalJump().SetTarget(nextLabel))

		// add true label
		v.Factory.AppendToBlock(stmt.TrueLabel)

		// add block
		fmt.Println("adding block", stmt.InnerBlock)
		v.Factory.AppendBlock(stmt.InnerBlock)

		// add jmp to final label
		v.Factory.AppendToBlock(v.Factory.NewUnconditionalJump().SetTarget(finalLabel))

		// add next label
		v.Factory.AppendToBlock(nextLabel)

	}

	// add else block
	if len(elseBlock) > 0 {
		v.Factory.AppendToBlock(v.Factory.NewComment().SetComment("else block"))
		v.Factory.AppendBlock(elseBlock)
	}

	// add final label
	v.Factory.AppendToBlock(finalLabel)

	return nil
}

func (v *IrVisitor) VisitIfChain(ctx *compiler.IfChainContext) interface{} {
	// TODO: new scope
	previousBlock := v.Factory.MainBlock
	conditionBlock := make(tac.TACBlock, 0)
	v.Factory.MainBlock = &conditionBlock

	trueLabel := v.Factory.NewLabel()
	wrapper := v.Visit(ctx.Expr()).(*value.ValueWrapper)
	condition := v.Factory.NewBoolExpression().SetLeft(wrapper.Val).SetRight(v.Factory.NewLiteral().SetValue("1")).SetLeftCast("int").SetOp(tac.EQ)
	conditionalJmp := v.Factory.NewConditionalJump().SetCondition(condition).SetTarget(trueLabel)

	innerBlock := make(tac.TACBlock, 0)
	v.Factory.MainBlock = &innerBlock

	for _, stmt := range ctx.AllStmt() {
		v.Visit(stmt)
	}

	v.Factory.MainBlock = previousBlock
	return &ifStmt{
		Condition:      conditionalJmp,
		TrueLabel:      trueLabel,
		InnerBlock:     innerBlock,
		ConditionBlock: conditionBlock,
	}
}

func (v *IrVisitor) VisitElseStmt(ctx *compiler.ElseStmtContext) interface{} {
	// TODO: new scope
	prevBlock := v.Factory.MainBlock
	auxBlock := make(tac.TACBlock, 0)
	v.Factory.MainBlock = &auxBlock

	for _, stmt := range ctx.AllStmt() {
		v.Visit(stmt)
	}

	v.Factory.MainBlock = prevBlock
	return auxBlock
}

func (v *IrVisitor) VisitSwitchStmt(ctx *compiler.SwitchStmtContext) interface{} {

	return nil
}

func (v *IrVisitor) GetCaseValue(tree antlr.ParseTree) *value.ValueWrapper {
	return nil
}

func (v *IrVisitor) VisitSwitchCase(ctx *compiler.SwitchCaseContext) interface{} {
	return nil
}

func (v *IrVisitor) VisitDefaultCase(ctx *compiler.DefaultCaseContext) interface{} {

	return nil
}

func (v *IrVisitor) VisitWhileStmt(ctx *compiler.WhileStmtContext) interface{} {
	return nil
}
