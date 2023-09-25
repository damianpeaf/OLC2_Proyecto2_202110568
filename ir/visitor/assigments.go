package visitor

import (
	"fmt"

	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"
)

func (v *IrVisitor) VisitType(ctx *compiler.TypeContext) interface{} {
	return ctx.GetText()
}

func (v *IrVisitor) VisitTypeValueDecl(ctx *compiler.TypeValueDeclContext) interface{} {

	varName := ctx.ID().GetText()
	valueWrapper := v.Visit(ctx.Expr()).(*ValueWrapper)

	v.ScopeTrace.CurrentScope.NewVariable(varName, valueWrapper.Val, valueWrapper.Metadata)

	return nil
}

func (v *IrVisitor) VisitValueDecl(ctx *compiler.ValueDeclContext) interface{} {

	varName := ctx.ID().GetText()
	varValue := v.Visit(ctx.Expr()).(*ValueWrapper)

	v.ScopeTrace.CurrentScope.NewVariable(varName, varValue.Val, varValue.Metadata)

	return nil
}

func (v *IrVisitor) VisitTypeDecl(ctx *compiler.TypeDeclContext) interface{} {

	varName := ctx.ID().GetText()
	varType := v.Visit(ctx.Type_()).(string)

	v.ScopeTrace.CurrentScope.NewVariable(varName, v.Utility.NilValue(), varType)
	fmt.Println("Type decl: ", varName, varType)
	return nil
}
