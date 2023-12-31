package visitor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/abstract"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/value"
)

func (v *IrVisitor) VisitIntLiteral(ctx *compiler.IntLiteralContext) interface{} {
	return &value.ValueWrapper{
		Val:      v.Factory.NewLiteral().SetValue(ctx.GetText()),
		Metadata: abstract.IVOR_INT,
	}
}

func (v *IrVisitor) VisitFloatLiteral(ctx *compiler.FloatLiteralContext) interface{} {
	return &value.ValueWrapper{
		Val:      v.Factory.NewLiteral().SetValue(ctx.GetText()),
		Metadata: abstract.IVOR_FLOAT,
	}
}

func (v *IrVisitor) VisitStringLiteral(ctx *compiler.StringLiteralContext) interface{} {
	stringVal := ctx.GetText()[1 : len(ctx.GetText())-1]

	// \" \\ \n \r \
	stringVal = strings.ReplaceAll(stringVal, "\\\"", "\"")
	stringVal = strings.ReplaceAll(stringVal, "\\\\", "\\")
	stringVal = strings.ReplaceAll(stringVal, "\\n", "\n")
	stringVal = strings.ReplaceAll(stringVal, "\\r", "\r")

	// Character literal
	if len(stringVal) == 1 {
		asciiVal := int(stringVal[0])
		return &value.ValueWrapper{
			Val:      v.Factory.NewLiteral().SetValue(strconv.Itoa(asciiVal)),
			Metadata: abstract.IVOR_CHARACTER,
		}
	}
	return &value.ValueWrapper{
		Val:      v.Utility.SaveString(stringVal),
		Metadata: abstract.IVOR_STRING,
	}
}

func (v *IrVisitor) VisitBoolLiteral(ctx *compiler.BoolLiteralContext) interface{} {
	if ctx.GetText() == "true" {
		return &value.ValueWrapper{
			Val:      v.Factory.NewLiteral().SetValue("1"),
			Metadata: abstract.IVOR_BOOL,
		}
	}
	return &value.ValueWrapper{
		Val:      v.Factory.NewLiteral().SetValue("0"),
		Metadata: abstract.IVOR_BOOL,
	}
}

func (v *IrVisitor) VisitNilLiteral(ctx *compiler.NilLiteralContext) interface{} {
	return &value.ValueWrapper{
		Val:      v.Utility.NilValue(),
		Metadata: abstract.IVOR_NIL,
	}
}

func (v *IrVisitor) VisitLiteralExp(ctx *compiler.LiteralExpContext) interface{} {
	return v.Visit(ctx.Literal())
}

func (v *IrVisitor) VisitBinaryExp(ctx *compiler.BinaryExpContext) interface{} {

	op := ctx.GetOp().GetText()

	var left *value.ValueWrapper
	var right *value.ValueWrapper

	switch ctx.GetRight().(type) {
	case *compiler.FuncCallExpContext:
		right = v.Visit(ctx.GetRight()).(*value.ValueWrapper) // just in case its a recursive call
		left = v.Visit(ctx.GetLeft()).(*value.ValueWrapper)
	default:
		left = v.Visit(ctx.GetLeft()).(*value.ValueWrapper)
		right = v.Visit(ctx.GetRight()).(*value.ValueWrapper)
	}

	strat, ok := v.Strats[op]

	if !ok {
		panic("No strategy for " + op)
	}

	ok, result := strat.Validate(left, right)

	if !ok {
		fmt.Println("Error: Invalid operation between", left.Metadata, "and", right.Metadata, "binexp")
		return &value.ValueWrapper{
			Val:      v.Utility.NilValue(),
			Metadata: abstract.IVOR_NIL,
		}
	}

	return result
}

func (v *IrVisitor) VisitUnaryExp(ctx *compiler.UnaryExpContext) interface{} {

	exp := v.Visit(ctx.Expr()).(*value.ValueWrapper)

	strat, ok := v.UnStrats[ctx.GetOp().GetText()]

	if !ok {
		panic("Unary operator not found")
	}

	ok, result := strat.Validate(exp)

	if !ok {
		fmt.Println("Error: Invalid operation between", exp.Metadata)
		return &value.ValueWrapper{
			Val:      v.Utility.NilValue(),
			Metadata: abstract.IVOR_NIL,
		}
	}

	return result

}

func (v *IrVisitor) VisitParenExp(ctx *compiler.ParenExpContext) interface{} {
	return v.Visit(ctx.Expr())
}

func (v *IrVisitor) VisitIdExp(ctx *compiler.IdExpContext) interface{} {
	varName := ctx.Id_pattern().GetText()

	variable := v.ScopeTrace.GetVariable(varName)

	if variable == nil {
		return v.GetNilVW()
	}

	fmt.Println("IDEXP: ", variable.Name, variable.Type, variable.Address, variable.FrameRelative, variable.Offset)

	temp := v.Factory.NewTemp()
	stackAddress := variable.GetStackStmt(v.Factory)
	assign := v.Factory.NewSimpleAssignment().SetAssignee(temp).SetVal(stackAddress)
	v.Factory.AppendToBlock(assign)

	// ? pointer
	return &value.ValueWrapper{
		Val:      temp,
		Metadata: variable.Type,
	}
}

func (v *IrVisitor) VisitFuncCallExp(ctx *compiler.FuncCallExpContext) interface{} {
	return v.Visit(ctx.Func_call())
}
