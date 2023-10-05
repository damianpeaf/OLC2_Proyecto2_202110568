package abstract

import (
	"strconv"

	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/tac"
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
)

// IVOR stands for Internal Value Object Representation
type IVOR struct {
	Name    string
	Type    string
	Address int
	// Copy() IVOR // ? it would be interesting
}

func (i *IVOR) GetStackStmt(f *tac.TACFactory) *tac.StackIndexedValue {
	return f.NewStackIndexed().SetIndex(f.NewLiteral().SetValue(strconv.Itoa(i.Address)))
}
