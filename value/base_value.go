package value

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
type IVOR interface {
	Value() interface{}
	Type() string
	Copy() IVOR
}
