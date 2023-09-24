package repl

import "github.com/damianpeaf/OLC2_Proyecto2_202110568/value"

type PointerValue struct {
	AssocVariable *Variable
}

func (p PointerValue) Value() interface{} {
	return p
}

func (p PointerValue) Type() string {
	return value.IVOR_POINTER
}

func (p PointerValue) Copy() value.IVOR {
	return p
}
