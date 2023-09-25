package visitor

import (
	"strconv"
	"strings"

	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/abstract"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/tac"
)

type ValueWrapper struct {
	Val      tac.SimpleValue
	Metadata string
}

func (v *IrVisitor) VisitIntLiteral(ctx *compiler.IntLiteralContext) interface{} {
	return &ValueWrapper{
		Val:      v.Factory.NewLiteral().SetValue(ctx.GetText()),
		Metadata: abstract.IVOR_INT,
	}
}

func (v *IrVisitor) VisitFloatLiteral(ctx *compiler.FloatLiteralContext) interface{} {
	return &ValueWrapper{
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
		return &ValueWrapper{
			Val:      v.Factory.NewLiteral().SetValue(strconv.Itoa(asciiVal)),
			Metadata: abstract.IVOR_CHARACTER,
		}
	}
	return &ValueWrapper{
		Val:      v.Utility.SaveString(stringVal),
		Metadata: abstract.IVOR_STRING,
	}
}

func (v *IrVisitor) VisitBoolLiteral(ctx *compiler.BoolLiteralContext) interface{} {
	if ctx.GetText() == "true" {
		return &ValueWrapper{
			Val:      v.Factory.NewLiteral().SetValue("1"),
			Metadata: abstract.IVOR_BOOL,
		}
	}
	return &ValueWrapper{
		Val:      v.Factory.NewLiteral().SetValue("0"),
		Metadata: abstract.IVOR_BOOL,
	}
}

func (v *IrVisitor) VisitNilLiteral(ctx *compiler.NilLiteralContext) interface{} {
	return &ValueWrapper{
		Val:      v.Utility.NilValue(),
		Metadata: abstract.IVOR_NIL,
	}
}

func (v *IrVisitor) VisitLiteralExp(ctx *compiler.LiteralExpContext) interface{} {
	return v.Visit(ctx.Literal())
}

func (v *IrVisitor) VisitBinaryExp(ctx *compiler.BinaryExpContext) interface{} {

	op := ctx.GetOp().GetText()
	left := v.Visit(ctx.GetLeft()).(*ValueWrapper)
	right := v.Visit(ctx.GetRight()).(*ValueWrapper)

	return &ValueWrapper{
		Val:      v.Utility.ArithmeticOperation(left.Val, right.Val, op),
		Metadata: left.Metadata,
	}
}
