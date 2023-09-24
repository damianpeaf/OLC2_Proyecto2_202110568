package repl

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/value"

	"github.com/antlr4-go/antlr/v4"
)

type ObjectValue struct {
	InternalScope *BaseScope
	AuxObject     interface{}
	ConcretType   string
	v             *ReplVisitor
	t             antlr.Token
}

func (o ObjectValue) Value() interface{} {
	return o
}

func (o ObjectValue) Type() string {
	if o.ConcretType != "" {
		return o.ConcretType
	}

	return value.IVOR_OBJECT
}

func (o *ObjectValue) Copy() value.IVOR {
	args := make([]*Argument, 0)

	for _, prop := range o.InternalScope.variables {
		args = append(args, &Argument{
			Name:  prop.Name,
			Value: prop.Value,
		})
	}

	return NewObjectValue(o.v, o.ConcretType, o.t, args, true)
}

func NewObjectValue(v *ReplVisitor, targetStruct string, targetToken antlr.Token, args []*Argument, allowReinitialize bool) value.IVOR {

	// Check if struct exists

	structTemplate, msg := v.ScopeTrace.GlobalScope.GetStruct(targetStruct)

	if structTemplate == nil {
		v.ErrorTable.NewSemanticError(targetToken, msg)
	}

	internalScope := NewStructScope()

	prevScope := v.ScopeTrace.CurrentScope
	v.ScopeTrace.CurrentScope = internalScope

	defer func() {
		// restore scope
		v.ScopeTrace.CurrentScope = prevScope
	}()

	// Add fields to internal scope
	for _, field := range structTemplate.Fields {
		v.Visit(field)
	}

	// create args map
	argMap := make(map[string]*Argument)

	for _, arg := range args {

		// repeat arg
		if _, ok := argMap[arg.Name]; ok {
			v.ErrorTable.NewSemanticError(arg.Token, "El argumento "+arg.Name+" ya fue definido")
			return value.DefaultNilValue
		}

		argMap[arg.Name] = arg
	}

	// validate constructor args
	wasConst := false
	usedArgs := make(map[string]bool)

	for _, prop := range internalScope.variables {
		arg, found := argMap[prop.Name]

		if !found {
			if prop.Value == value.DefaultUnInitializedValue {
				v.ErrorTable.NewSemanticError(targetToken, "El campo "+prop.Name+" no fue inicializado en el constructor")
				return value.DefaultNilValue
			}

			continue
		}

		// then the arg exists
		if prop.IsConst {
			if (prop.Value != value.DefaultUnInitializedValue) && !allowReinitialize {
				v.ErrorTable.NewSemanticError(targetToken, "El campo "+prop.Name+" es inmutable y ya fue inicializado")
				return value.DefaultNilValue
			}

			wasConst = true
			prop.IsConst = false
		}

		var throwError bool = false
		var msg string = ""
		var assignValue value.IVOR = arg.Value

		// pointer support
		if arg.PassByReference {

			if arg.VariableRef == nil {
				msg = "No es posible pasar por referencia un valor que no este asociado a una variable"
				throwError = true
			}

			// create the pointer
			assignValue = &PointerValue{
				AssocVariable: arg.VariableRef,
			}
		}

		if !throwError {
			throwError, msg = prop.Assign(assignValue, false)
		}

		if wasConst {
			prop.IsConst = true
			wasConst = false
		}

		if !throwError {
			v.ErrorTable.NewSemanticError(targetToken, msg)
			return value.DefaultNilValue
		}

		usedArgs[prop.Name] = true
	}

	// validate unused args
	for _, arg := range args {
		if _, ok := usedArgs[arg.Name]; !ok {
			v.ErrorTable.NewSemanticError(arg.Token, "El argumento "+arg.Name+" no es utilizado en el constructor")
		}
	}

	// mutable related flags
	for _, prop := range internalScope.variables {
		prop.isProp = true
	}

	// self object

	selfObject := &ObjectValue{
		InternalScope: internalScope,
		ConcretType:   value.IVOR_SELF,
	}

	instanceInternalScope := NewStructScope()

	instanceInternalScope.AddVariable("self", value.IVOR_SELF, selfObject, true, false, nil)

	// make functions use the instance scope

	for _, function := range internalScope.functions {
		f, ok := function.(*Function)

		if !ok {
			continue
		}

		f.DefaultScope = instanceInternalScope
	}

	// create instance
	return &ObjectValue{
		InternalScope: internalScope,
		ConcretType:   targetStruct,
		v:             v,
		t:             targetToken,
	}
}

func IsArgValidForStruct(arg []*Argument) bool {
	for _, a := range arg {

		if a.Name == "" {
			return false
		}
	}
	return true
}
