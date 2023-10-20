package abstract

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/tac"
)

const (
	USER_DEFINED_FUNCTION = "user_defined_function"
	BUILTIN_FUNCTION      = "builtin_function"
)

type Function struct {
	Name          string
	Params        []*Param
	Type          string
	ScopeTrace    *ScopeTrace
	ReturnType    string
	ReturnTemp    *tac.Temp
	StructPointer *tac.Temp
	StructRef     *IVOR
}

func (f *Function) ValidateArgs(args []*Argument) (bool, map[string]*Argument) {

	// validate arg count
	if len(args) != len(f.Params) {
		return false, nil
	}

	argsMap := make(map[string]*Argument)
	finalArgsMap := make(map[string]*Argument)

	for _, arg := range args {
		argsMap[arg.Name] = arg
	}

	errorFound := false

	for i, param := range f.Params {

		// determine param type
		var argToValidate *Argument = nil

		if param.ExternName == "" {
			// inner = arg name
			argToValidate = argsMap[param.InnerName]

		} else if param.ExternName == "_" {
			// positional arg
			argToValidate = args[i]
		} else {
			// extern = arg name
			argToValidate = argsMap[param.ExternName]
		}

		// validate arg exists
		if argToValidate == nil {
			errorFound = true
			continue
		}

		// validate type
		if argToValidate.Wrapper.Metadata != param.Type && param.Type != IVOR_ANY {
			errorFound = true
			continue
		}

		// validate pass by reference
		if argToValidate.PassByReference != param.PassByReference {
			errorFound = true
			continue
		}

		// add to final args map
		finalArgsMap[param.InnerName] = argToValidate
	}

	if errorFound {
		return false, nil
	}

	return true, finalArgsMap
}
