package abstract

import (
	"fmt"
	"strconv"

	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/tac"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/value"
)

const (
	IVOR_INT              = "Int"
	IVOR_FLOAT            = "Float"
	IVOR_STRING           = "String"
	IVOR_BOOL             = "Bool"
	IVOR_CHARACTER        = "Character"
	IVOR_NIL              = "nil"
	IVOR_BUILTIN_FUNCTION = "builtin_function"
	IVOR_FUNCTION         = "function"
	IVOR_VECTOR           = "vector"
	IVOR_OBJECT           = "object"
	IVOR_ANY              = "any"
	IVOR_POINTER          = "pointer"
	IVOR_MATRIX           = "matrix"
	IVOR_SELF             = "self"
	IVOR_UNINITIALIZED    = "uninitialized"
	IVOR_RANGE            = "range"
)

// IVOR stands for Internal Value Object Representation
type IVOR struct {
	Name          string
	Type          string
	Address       int
	FrameRelative bool
	Offset        int
	// Copy() IVOR // ? it would be interesting
}

func (i *IVOR) GetStackStmt(f *tac.TACFactory) *tac.StackIndexedValue {
	var index tac.SimpleValue = f.NewLiteral().SetValue(strconv.Itoa(i.Address))

	fmt.Println("IVOR: ", i.Name, i.Type, i.Address, i.FrameRelative, i.Offset)

	if i.FrameRelative {
		// ? This could be reserved in some way
		// framePointer + (address + offset) -> stack address
		index = f.NewTemp()
		addressComputation := f.NewCompoundAssignment().SetAssignee(index).SetLeft(f.GetFramePointer()).SetRight(f.NewLiteral().SetValue(strconv.Itoa(i.Address + i.Offset))).SetLeftCast("int").SetRightCast("int").SetOperator("+")
		f.AppendToBlock(addressComputation)
	}

	return f.NewStackIndexed().SetIndex(index)
}

type Range struct {
	Init *value.ValueWrapper
	End  *value.ValueWrapper
}
