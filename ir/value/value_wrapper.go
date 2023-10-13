package value

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/tac"
)

type ValueWrapper struct {
	Val      tac.SimpleValue
	Metadata string
	Aux      interface{}
}
