package abstract

const (
	USER_DEFINED_FUNCTION = "user_defined_function"
	BUILTIN_FUNCTION      = "builtin_function"
)

type Function struct {
	Name  string
	Param []*Param
	Type  string
	// DeclScope       *BaseScope // always reference the global scope
	// DefaultScope    *BaseScope
}
