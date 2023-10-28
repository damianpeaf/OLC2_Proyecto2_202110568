package visitor

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/abstract"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/tac"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/value"
)

func (v *IrVisitor) BuiltinHandler(funcObj *abstract.Function, args []*abstract.Argument) *value.ValueWrapper {

	switch funcObj.Name {
	case "print": // thiw should be a pointer
		return v.PrintBuiltin(args)
	case "Int":
		return v.IntBuiltIn(args)
	case "Float":
		return v.FloatBuiltIn(args)
	case "String":
		return v.StringBuiltIn(args)
	}

	// vector builtins
	switch funcObj {
	case abstract.VectorAppendFunc:
		return v.VectorAppend(funcObj, args)
	case abstract.VectorRemoveLastFunc:
		return v.VectorRemoveLast(funcObj, args)
	case abstract.VectorRemoveFunc:
		return v.VectorRemove(funcObj, args)

	}

	return v.GetNilVW()
}

func (v *IrVisitor) PrintBuiltin(args []*abstract.Argument) *value.ValueWrapper {

	for _, arg := range args {

		// check if its nil
		nilLiteral := v.Utility.NilValue()

		// if val is nil, print nil
		printNilLabel := v.Factory.NewLabel()
		endLabel := v.Factory.NewLabel()

		condition := v.Factory.NewBoolExpression().SetLeft(arg.Wrapper.Val).SetRight(nilLiteral).SetOp(tac.EQ)
		v.Factory.AppendToBlock(v.Factory.NewConditionalJump().SetCondition(condition).SetTarget(printNilLabel))

		switch arg.Wrapper.Metadata {
		case abstract.IVOR_INT:
			v.Factory.AppendToBlock(v.Factory.NewPrint().SetMode(tac.PRINT_DIGIT).SetVal(arg.Wrapper.Val).SetCast("int"))
		case abstract.IVOR_FLOAT:
			v.Factory.AppendToBlock(v.Factory.NewPrint().SetMode(tac.PRINT_FLOAT).SetVal(arg.Wrapper.Val))
		case abstract.IVOR_CHARACTER:
			v.Factory.AppendToBlock(v.Factory.NewPrint().SetMode(tac.PRINT_CHAR).SetVal(arg.Wrapper.Val).SetCast("int"))
		case abstract.IVOR_STRING:
			v.Utility.PrintString(arg.Wrapper.Val)
		case abstract.IVOR_BOOL:
			v.Factory.Utility.PrintBool(arg.Wrapper.Val)
		case abstract.IVOR_NIL:
			v.Utility.PrintNil()
		default:
		}
		v.Factory.AppendToBlock(v.Factory.NewUnconditionalJump().SetTarget(endLabel))

		v.Factory.AppendToBlock(printNilLabel)
		v.Utility.PrintNil()

		v.Factory.AppendToBlock(endLabel)
		v.Utility.PrintSpace()
	}
	// print new line
	v.Utility.PrintNewLine()

	return v.GetNilVW()
}

func (v *IrVisitor) IntBuiltIn(args []*abstract.Argument) *value.ValueWrapper {
	// String, float -> Int
	if len(args) != 1 {
		return v.GetNilVW()
	}

	switch args[0].Wrapper.Metadata {
	case abstract.IVOR_STRING:
		params := v.Factory.GetBuiltinParams("__string_to_int")
		stringParam := params[0]
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(stringParam).SetVal(args[0].Wrapper.Val))
		v.Factory.AppendToBlock(v.Factory.NewMethodCall("__string_to_int"))
		return &value.ValueWrapper{
			Val:      stringParam,
			Metadata: abstract.IVOR_INT,
		}
	case abstract.IVOR_FLOAT:
		temp := v.Factory.NewTemp()
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(temp).SetVal(args[0].Wrapper.Val).SetCast("int"))
		return &value.ValueWrapper{
			Val:      temp,
			Metadata: abstract.IVOR_INT,
		}
	default:
		return v.GetNilVW()
	}

}

func (v *IrVisitor) FloatBuiltIn(args []*abstract.Argument) *value.ValueWrapper {

	// String, int -> Float
	if len(args) != 1 {
		return v.GetNilVW()
	}

	switch args[0].Wrapper.Metadata {
	case abstract.IVOR_STRING:
		v.Factory.GetBuiltinParams("__string_to_int") // just to reserve the params
		params := v.Factory.GetBuiltinParams("__string_to_float")
		stringParam := params[0]
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(stringParam).SetVal(args[0].Wrapper.Val))
		v.Factory.AppendToBlock(v.Factory.NewMethodCall("__string_to_float"))
		return &value.ValueWrapper{
			Val:      stringParam,
			Metadata: abstract.IVOR_FLOAT,
		}
	case abstract.IVOR_INT:
		temp := v.Factory.NewTemp()
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(temp).SetVal(args[0].Wrapper.Val).SetCast("float"))
		return &value.ValueWrapper{
			Val:      temp,
			Metadata: abstract.IVOR_FLOAT,
		}
	default:
		return v.GetNilVW()
	}
}

func (v *IrVisitor) StringBuiltIn(args []*abstract.Argument) *value.ValueWrapper {

	// Int, float, bool -> String
	if len(args) != 1 {
		return v.GetNilVW()
	}

	switch args[0].Wrapper.Metadata {
	case abstract.IVOR_INT:
		params := v.Factory.GetBuiltinParams("__int_to_string")
		intParam := params[0]
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(intParam).SetVal(args[0].Wrapper.Val))
		v.Factory.AppendToBlock(v.Factory.NewMethodCall("__int_to_string"))
		return &value.ValueWrapper{
			Val:      intParam,
			Metadata: abstract.IVOR_STRING,
		}
	case abstract.IVOR_FLOAT:
		v.Factory.GetBuiltinParams("__int_to_string")
		params := v.Factory.GetBuiltinParams("__float_to_string")
		floatParam := params[0]
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(floatParam).SetVal(args[0].Wrapper.Val))
		v.Factory.AppendToBlock(v.Factory.NewMethodCall("__float_to_string"))
		return &value.ValueWrapper{
			Val:      floatParam,
			Metadata: abstract.IVOR_STRING,
		}
	case abstract.IVOR_BOOL:
		params := v.Factory.GetBuiltinParams("__bool_to_string")
		boolParam := params[0]
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(boolParam).SetVal(args[0].Wrapper.Val))
		v.Factory.AppendToBlock(v.Factory.NewMethodCall("__bool_to_string"))
		return &value.ValueWrapper{
			Val:      boolParam,
			Metadata: abstract.IVOR_STRING,
		}
	default:
		return v.GetNilVW()
	}
}

func (v *IrVisitor) VectorAppend(funcObj *abstract.Function, args []*abstract.Argument) *value.ValueWrapper {

	// just one value to append
	if len(args) != 1 {
		return v.GetNilVW()
	}
	params := v.Factory.GetBuiltinParams("__vector_append")

	vectorParam := params[0]
	itemParam := params[1]

	vectorStackIndex := funcObj.StructRef.GetStackIndex(v.Factory)
	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(vectorParam).SetVal(vectorStackIndex))

	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(itemParam).SetVal(args[0].Wrapper.Val))

	v.Factory.AppendToBlock(v.Factory.NewMethodCall("__vector_append"))

	return v.GetNilVW()
}

func (v *IrVisitor) VectorRemoveLast(funcObj *abstract.Function, args []*abstract.Argument) *value.ValueWrapper {

	// validate arg

	if len(args) != 0 {
		return v.GetNilVW()
	}

	params := v.Factory.GetBuiltinParams("__vector_remove_last")
	vectorParam := params[0]

	vectorStackIndex := funcObj.StructRef.GetStackIndex(v.Factory)
	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(vectorParam).SetVal(vectorStackIndex))

	v.Factory.AppendToBlock(v.Factory.NewMethodCall("__vector_remove_last"))

	return v.GetNilVW()
}

func (v *IrVisitor) VectorRemove(funcObj *abstract.Function, args []*abstract.Argument) *value.ValueWrapper {
	if len(args) != 1 {
		return v.GetNilVW()
	}
	if args[0].Name != "at" {
		return v.GetNilVW()
	}
	params := v.Factory.GetBuiltinParams("__vector_remove")
	vectorParam := params[0]
	indexParam := params[1]

	vectorStackIndex := funcObj.StructRef.GetStackIndex(v.Factory)
	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(vectorParam).SetVal(vectorStackIndex))

	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(indexParam).SetVal(args[0].Wrapper.Val))

	v.Factory.AppendToBlock(v.Factory.NewMethodCall("__vector_remove"))

	return v.GetNilVW()
}
