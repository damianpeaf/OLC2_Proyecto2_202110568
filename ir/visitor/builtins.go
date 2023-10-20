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
		// ? nil check func
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
		v.Utility.PrintSpace()
	}
	// print new line
	v.Utility.PrintNewLine()

	return v.GetNilVW()
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

	if args[0].Name != "at" {
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
	params := v.Factory.GetBuiltinParams("__vector_remove")
	vectorParam := params[0]
	indexParam := params[1]

	vectorStackIndex := funcObj.StructRef.GetStackIndex(v.Factory)
	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(vectorParam).SetVal(vectorStackIndex))

	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(indexParam).SetVal(args[0].Wrapper.Val))

	v.Factory.AppendToBlock(v.Factory.NewMethodCall("__vector_remove"))

	return v.GetNilVW()
}
