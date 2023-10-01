package abstract

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/value"
)

type Argument struct {
	Name               string
	Wrapper            *value.ValueWrapper
	PassByReference    bool
	VariableRefAddress int
}
