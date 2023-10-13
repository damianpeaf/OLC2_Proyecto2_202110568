package visitor

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/abstract"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/tac"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/value"
)

func (v *IrVisitor) VisitReturnStmt(ctx *compiler.ReturnStmtContext) interface{} {
	v.Factory.AppendToBlock(v.Factory.NewComment().SetComment("return stmt"))
	v.Factory.AppendToBlock(v.Factory.NewUnconditionalJump().SetTarget(v.Transfer.ReturnLabel))
	return nil
}

func (v *IrVisitor) VisitBreakStmt(ctx *compiler.BreakStmtContext) interface{} {
	v.Factory.AppendToBlock(v.Factory.NewComment().SetComment("break stmt"))
	v.Factory.AppendToBlock(v.Factory.NewUnconditionalJump().SetTarget(v.Transfer.BreakLabel))
	return nil
}

func (v *IrVisitor) VisitContinueStmt(ctx *compiler.ContinueStmtContext) interface{} {
	v.Factory.AppendToBlock(v.Factory.NewComment().SetComment("continue stmt"))
	v.Factory.AppendToBlock(v.Factory.NewUnconditionalJump().SetTarget(v.Transfer.ContinueLabel))
	return nil
}

func (v *IrVisitor) VisitWhileStmt(ctx *compiler.WhileStmtContext) interface{} {

	// init:
	// if (!condition) goto end
	// body
	// goto init
	// end:

	v.Factory.AppendToBlock(v.Factory.NewComment().SetComment("while loop"))
	startLabel := v.Factory.NewLabel()
	endLabel := v.Factory.NewLabel()
	v.Factory.AppendToBlock(startLabel)

	prevContinue := v.Transfer.ContinueLabel
	v.Transfer.ContinueLabel = startLabel

	prevBreak := v.Transfer.BreakLabel
	v.Transfer.BreakLabel = endLabel

	wrapper := v.Visit(ctx.Expr()).(*value.ValueWrapper)
	condition := v.Factory.NewBoolExpression().SetLeft(wrapper.Val).SetLeftCast("int").SetOp(tac.EQ).SetRight(v.Factory.NewLiteral().SetValue("0"))
	conditional := v.Factory.NewConditionalJump().SetCondition(condition).SetTarget(endLabel)
	v.Factory.AppendToBlock(conditional)

	v.ScopeTrace.NextScope()
	for _, stmt := range ctx.AllStmt() {
		v.Visit(stmt)
	}

	v.ScopeTrace.PrevScope()
	v.Factory.AppendToBlock(v.Factory.NewUnconditionalJump().SetTarget(startLabel))
	v.Factory.AppendToBlock(endLabel)

	v.Transfer.ContinueLabel = prevContinue
	v.Transfer.BreakLabel = prevBreak
	return nil
}

func (v *IrVisitor) VisitNumericRange(ctx *compiler.NumericRangeContext) interface{} {

	leftExpr := v.Visit(ctx.Expr(0)).(*value.ValueWrapper)
	rightExpr := v.Visit(ctx.Expr(1)).(*value.ValueWrapper)

	return &value.ValueWrapper{
		Val:      nil,
		Metadata: abstract.IVOR_RANGE,
		Aux: &abstract.Range{
			Init: leftExpr,
			End:  rightExpr,
		},
	}
}

func (v *IrVisitor) VisitForStmt(ctx *compiler.ForStmtContext) interface{} {
	varName := ctx.ID().GetText()
	v.ScopeTrace.NextScope()

	if ctx.Range_() != nil {
		/*
			iteratorVar = ran.Init
			condition:
			if (iteratorVar > ran.End) goto break
			// body
			continue:
			iteratorVar = iteratorVar + 1
			goto condition
			break:
		*/
		ran := v.Visit(ctx.Range_()).(*value.ValueWrapper).Aux.(*abstract.Range)

		// labels
		conditionLabel := v.Factory.NewLabel()
		continueLabel := v.Factory.NewLabel()
		breakLabel := v.Factory.NewLabel()

		// transfer
		prevContinue := v.Transfer.ContinueLabel
		v.Transfer.ContinueLabel = continueLabel

		prevBreak := v.Transfer.BreakLabel
		v.Transfer.BreakLabel = breakLabel

		v.Factory.AppendToBlock(v.Factory.NewComment().SetComment("for loop"))
		// iteratorVar = ran.Init
		iteratorVar := v.ScopeTrace.GetVariable(varName)
		iteratorVar.Type = abstract.IVOR_INT
		stackAddress := iteratorVar.GetStackStmt(v.Factory)
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(stackAddress).SetVal(ran.Init.Val))

		// condition:
		v.Factory.AppendToBlock(conditionLabel)

		// if (iteratorVar > ran.End) goto break
		tmpVarVal := v.Factory.NewTemp()
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(tmpVarVal).SetVal(stackAddress))
		condition := v.Factory.NewBoolExpression().SetLeft(tmpVarVal).SetLeftCast("int").SetOp(tac.GT).SetRight(ran.End.Val).SetRightCast("int")
		conditional := v.Factory.NewConditionalJump().SetCondition(condition).SetTarget(breakLabel)
		v.Factory.AppendToBlock(conditional)

		// body
		v.Factory.AppendToBlock(v.Factory.NewComment().SetComment("inner loop"))
		v.ScopeTrace.NextScope()
		for _, stmt := range ctx.AllStmt() {
			v.Visit(stmt)
		}
		v.ScopeTrace.PrevScope() // inner for

		// continue:
		v.Factory.AppendToBlock(continueLabel)

		// iteratorVar = iteratorVar + 1
		v.Factory.AppendToBlock(v.Factory.NewComment().SetComment("continue:"))
		resultTemp := v.Factory.NewTemp()
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(tmpVarVal).SetVal(stackAddress))
		one := v.Factory.NewLiteral().SetValue("1")
		increase := v.Factory.NewCompoundAssignment().SetAssignee(resultTemp).SetLeft(tmpVarVal).SetRight(one).SetOperator("+")
		v.Factory.AppendToBlock(increase)
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(stackAddress).SetVal(resultTemp))

		// goto condition
		v.Factory.AppendToBlock(v.Factory.NewUnconditionalJump().SetTarget(conditionLabel))

		// break:
		v.Factory.AppendToBlock(breakLabel)

		// transfer
		v.Transfer.ContinueLabel = prevContinue
		v.Transfer.BreakLabel = prevBreak
		v.ScopeTrace.PrevScope() // outer for
	}

	return nil
}
