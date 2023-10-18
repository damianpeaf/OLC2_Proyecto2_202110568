package abstract

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/tac"
)

const (
	USER_DEFINED_FUNCTION = "user_defined_function"
	BUILTIN_FUNCTION      = "builtin_function"
)

type Function struct {
	Name       string
	Params     []*Param
	Type       string
	ScopeTrace *ScopeTrace
	ReturnType string
	ReturnTemp *tac.Temp
}
