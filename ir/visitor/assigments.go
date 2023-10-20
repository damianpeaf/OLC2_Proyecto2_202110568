package visitor

import (
	"log"

	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/value"
)

func (v *IrVisitor) VisitType(ctx *compiler.TypeContext) interface{} {
	return ctx.GetText()
}

// now, every decl will ony assign the value (and/or type), because it was "statically" analyzed

func (v *IrVisitor) VisitTypeValueDecl(ctx *compiler.TypeValueDeclContext) interface{} {
	varName := ctx.ID().GetText()
	vw := v.Visit(ctx.Expr()).(*value.ValueWrapper)
	variable := v.ScopeTrace.GetVariable(varName)
	variable.Type = vw.Metadata
	stackAddress := variable.GetStackStmt(v.Factory)
	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(stackAddress).SetVal(vw.Val))
	return nil
}

func (v *IrVisitor) VisitValueDecl(ctx *compiler.ValueDeclContext) interface{} {
	varName := ctx.ID().GetText()
	vw := v.Visit(ctx.Expr()).(*value.ValueWrapper)
	variable := v.ScopeTrace.GetVariable(varName)
	variable.Type = vw.Metadata
	stackAddress := variable.GetStackStmt(v.Factory)
	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(stackAddress).SetVal(vw.Val))

	return nil
}

func (v *IrVisitor) VisitTypeDecl(ctx *compiler.TypeDeclContext) interface{} {
	// Just set the type
	varType := ctx.Type_().GetText()
	variable := v.ScopeTrace.GetVariable(ctx.ID().GetText())
	variable.Type = varType
	return nil
}

func (v *IrVisitor) VisitDirectAssign(ctx *compiler.DirectAssignContext) interface{} {

	pattern := ctx.Id_pattern().GetText()
	varValue := v.Visit(ctx.Expr()).(*value.ValueWrapper)

	variable := v.ScopeTrace.GetVariable(pattern)

	if variable == nil {
		return v.GetNilVW()
	} else {

		// TODO: copy object ifs is a struct or vector

		// assign
		stackAddress := variable.GetStackStmt(v.Factory)
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(stackAddress).SetVal(varValue.Val))
	}

	return nil

}

func (v *IrVisitor) VisitArithmeticAssign(ctx *compiler.ArithmeticAssignContext) interface{} {
	pattern := ctx.Id_pattern().GetText()

	variable := v.ScopeTrace.GetVariable(pattern)

	if variable == nil {
		return v.GetNilVW()
	} else {

		stackAddress := variable.GetStackStmt(v.Factory)
		varTemp := v.Factory.NewTemp()
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(varTemp).SetVal(stackAddress))

		leftValue := &value.ValueWrapper{
			Val:      varTemp,
			Metadata: variable.Type,
		}
		rightValue := v.Visit(ctx.Expr()).(*value.ValueWrapper)

		op := string(ctx.GetOp().GetText()[0])

		strat, ok := v.Strats[op]

		if !ok {
			log.Fatal("Binary operator not found")
		}

		ok, varValue := strat.Validate(leftValue, rightValue)

		if !ok {
			return v.GetNilVW()
		}

		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(stackAddress).SetVal(varValue.Val))
	}

	return v.GetNilVW()
}
