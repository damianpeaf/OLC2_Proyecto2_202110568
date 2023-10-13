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
	previousBlock := v.Factory.MainBlock
	conditionBlock := make(tac.TACBlock, 0)
	v.Factory.MainBlock = &conditionBlock

	trueLabel := v.Factory.NewLabel()
	wrapper := v.Visit(ctx.Expr()).(*value.ValueWrapper)
	condition := v.Factory.NewBoolExpression().SetLeft(wrapper.Val).SetRight(v.Factory.NewLiteral().SetValue("1")).SetLeftCast("int").SetOp(tac.EQ)
	conditionalJmp := v.Factory.NewConditionalJump().SetCondition(condition).SetTarget(trueLabel)

	v.ScopeTrace.NextScope()
	innerBlock := make(tac.TACBlock, 0)
	v.Factory.MainBlock = &innerBlock

	for _, stmt := range ctx.AllStmt() {
		v.Visit(stmt)
	}

	v.ScopeTrace.PrevScope()
	v.Factory.MainBlock = previousBlock
	return &ifStmt{
		Condition:      conditionalJmp,
		TrueLabel:      trueLabel,
		InnerBlock:     innerBlock,
		ConditionBlock: conditionBlock,
	}
}

func (v *IrVisitor) VisitElseStmt(ctx *compiler.ElseStmtContext) interface{} {
	prevBlock := v.Factory.MainBlock
	auxBlock := make(tac.TACBlock, 0)
	v.Factory.MainBlock = &auxBlock

	v.ScopeTrace.NextScope()
	for _, stmt := range ctx.AllStmt() {
		v.Visit(stmt)
	}

	v.Factory.MainBlock = prevBlock
	v.ScopeTrace.PrevScope()
	return auxBlock
}

func (v *IrVisitor) VisitSwitchStmt(ctx *compiler.SwitchStmtContext) interface{} {

	v.Factory.AppendToBlock(v.Factory.NewComment().SetComment("switch"))
	v.ScopeTrace.NextScope()
	wrapper := v.Visit(ctx.Expr()).(*value.ValueWrapper)

	endLabel := v.Factory.NewLabel()
	prevLabel := v.Transfer.BreakLabel
	v.Transfer.BreakLabel = endLabel

	for _, switchCase := range ctx.AllSwitch_case() {
		v.TraverseCase(wrapper, switchCase)
	}

	if ctx.Default_case() != nil {
		v.Visit(ctx.Default_case())
	}

	v.Factory.AppendToBlock(endLabel)
	v.Transfer.BreakLabel = prevLabel
	return nil
}

func (v *IrVisitor) TraverseCase(wrapper *value.ValueWrapper, tree antlr.ParseTree) interface{} {

	switch tree.(type) {
	case *compiler.SwitchCaseContext:
		// if expr != case; goto next case
		// ... case
		// goto end
		// next case
		// ...
		// end:
		v.Factory.AppendToBlock(v.Factory.NewComment().SetComment("case"))

		caseCtx := tree.(*compiler.SwitchCaseContext)
		caseWrapper := v.Visit(caseCtx.Expr()).(*value.ValueWrapper)

		nextCaseLabel := v.Factory.NewLabel()
		strat := v.Strats["=="]

		ok, result := strat.Validate(wrapper, caseWrapper)

		if !ok {
			fmt.Println("Error: Invalid operation between", wrapper.Metadata, "and", caseWrapper.Metadata)
			return v.GetNilVW()
		}

		condition := v.Factory.NewBoolExpression().SetLeft(result.Val).SetLeftCast("int").SetRight(v.Factory.NewLiteral().SetValue("1")).SetOp(tac.NEQ)
		conditionalJmp := v.Factory.NewConditionalJump().SetCondition(condition).SetTarget(nextCaseLabel)
		v.Factory.AppendToBlock(conditionalJmp)

		for _, stmt := range caseCtx.AllStmt() {
			v.Visit(stmt)
		}

		v.Factory.AppendToBlock(v.Factory.NewUnconditionalJump().SetTarget(v.Transfer.BreakLabel))
		v.Factory.AppendToBlock(nextCaseLabel)
		return nil
	default:
		return nil
	}
}

func (v *IrVisitor) VisitDefaultCase(ctx *compiler.DefaultCaseContext) interface{} {
	v.Factory.AppendToBlock(v.Factory.NewComment().SetComment("default case"))
	for _, stmt := range ctx.AllStmt() {
		v.Visit(stmt)
	}
	return nil
}

func (v *IrVisitor) VisitGuardStmt(ctx *compiler.GuardStmtContext) interface{} {
	/*
		if(expr == 1) goto end
		... block
		end:
	*/

	wrapper := v.Visit(ctx.Expr()).(*value.ValueWrapper)
	endLabel := v.Factory.NewLabel()

	condition := v.Factory.NewBoolExpression().SetLeft(wrapper.Val).SetLeftCast("int").SetRight(v.Factory.NewLiteral().SetValue("1")).SetOp(tac.EQ)
	conditionalJmp := v.Factory.NewConditionalJump().SetCondition(condition).SetTarget(endLabel)
	v.Factory.AppendToBlock(conditionalJmp)

	v.ScopeTrace.NextScope()
	for _, stmt := range ctx.AllStmt() {
		v.Visit(stmt)
	}

	v.ScopeTrace.PrevScope()
	v.Factory.AppendToBlock(endLabel)
	return nil
}
