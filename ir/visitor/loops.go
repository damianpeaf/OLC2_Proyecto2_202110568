package visitor

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"
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

	v.ScopeTrace.PushScope("while")
	for _, stmt := range ctx.AllStmt() {
		v.Visit(stmt)
	}

	v.ScopeTrace.PopScope()
	v.Factory.AppendToBlock(v.Factory.NewUnconditionalJump().SetTarget(startLabel))
	v.Factory.AppendToBlock(endLabel)

	v.Transfer.ContinueLabel = prevContinue
	v.Transfer.BreakLabel = prevBreak
	return nil
}
