package visitor

import (
	"fmt"

	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/abstract"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/tac"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/value"
)

func (v *IrVisitor) BuiltinHandler(name string, args []*abstract.Argument) *value.ValueWrapper {

	if name == "print" {
		return v.PrintBuiltin(args)
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
			fmt.Println("bool VAAAAAAAAL")
			fmt.Println(arg.Wrapper.Metadata)
			fmt.Println(arg.Wrapper.Val)
		}
		v.Utility.PrintSpace()
	}
	// print new line
	v.Utility.PrintNewLine()

	return v.GetNilVW()
}
