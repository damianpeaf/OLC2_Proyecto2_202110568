package repl

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/value"

	"github.com/antlr4-go/antlr/v4"
)

type ObjectBuiltInFunction struct {
	*Function
	Object     *ObjectValue
	CustomExec func(builtinRef *ObjectBuiltInFunction, visitor *ReplVisitor, args map[string]*Argument, token antlr.Token)
}

// implementing ivor
func (b ObjectBuiltInFunction) Type() string {
	return value.IVOR_FUNCTION
}

func (b ObjectBuiltInFunction) Value() interface{} {
	return b
}

func (b ObjectBuiltInFunction) Copy() value.IVOR {
	return b
}

func (f *ObjectBuiltInFunction) Exec(visitor *ReplVisitor, args []*Argument, token antlr.Token) {

	context := visitor.GetReplContext()

	// validate args
	argsOk, argsMap := f.ValidateArgs(context, args, token)

	if !argsOk {
		f.ReturnValue = value.DefaultNilValue
		return
	}

	f.CustomExec(f, visitor, argsMap, token)

}

// * Vector Built In Functions

// 1. Append
// vector.append(value)

// Parameters:

var appendParams = []*Param{
	// Just one positional argument
	{
		ExternName:      "_",
		InnerName:       "_",
		Type:            value.IVOR_ANY,
		PassByReference: false,
		Token:           nil,
	},
}

func appendCustomExec(builtinRef *ObjectBuiltInFunction, visitor *ReplVisitor, args map[string]*Argument, token antlr.Token) {

	builtinRef.ReturnValue = value.DefaultNilValue

	// get the vector
	vector := builtinRef.Object.AuxObject.(*VectorValue)

	// get the value
	arg := args["_"]

	if vector.ItemType != arg.Value.Type() {
		visitor.ErrorTable.NewSemanticError(arg.Token, "No se puede agregar un valor de tipo "+arg.Value.Type()+" a un vector de tipo "+vector.ItemType)
		return
	}
	vector.InternalValue = append(vector.InternalValue, arg.Value)
	vector.updateProps()
}

// 2. vector.remove(at: Int) -> nil

// parameters:

var removeParams = []*Param{
	// at: Int
	{
		ExternName:      "at",
		InnerName:       "at",
		Type:            value.IVOR_INT,
		PassByReference: false,
		Token:           nil,
	},
}

func removeCustomExec(builtinRef *ObjectBuiltInFunction, visitor *ReplVisitor, args map[string]*Argument, token antlr.Token) {

	builtinRef.ReturnValue = value.DefaultNilValue

	// get the vector
	vector := builtinRef.Object.AuxObject.(*VectorValue)

	// get the value
	arg := args["at"]

	if arg.Value.Type() != value.IVOR_INT {
		visitor.ErrorTable.NewSemanticError(arg.Token, "El argumento 'at' debe ser de tipo Int")
		return
	}

	// out of bounds
	if arg.Value.Value().(int) >= vector.Size() || arg.Value.Value().(int) < 0 {
		visitor.ErrorTable.NewSemanticError(arg.Token, "El indice esta fuera de rango")
		return
	}

	// remove the element
	vector.InternalValue = append(vector.InternalValue[:arg.Value.Value().(int)], vector.InternalValue[arg.Value.Value().(int)+1:]...)
	vector.updateProps()
}

// 3. removeLast
// vector.removeLast() -> nil

// parameters:

var removeLastParams = []*Param{}

func removeLastCustomExec(builtinRef *ObjectBuiltInFunction, visitor *ReplVisitor, args map[string]*Argument, token antlr.Token) {

	builtinRef.ReturnValue = value.DefaultNilValue

	// get the vector
	vector := builtinRef.Object.AuxObject.(*VectorValue)

	if vector.Size() == 0 {
		visitor.ErrorTable.NewSemanticError(token, "El vector esta vacio y no se puede remover el ultimo elemento")
		return
	}

	// remove the last element
	vector.InternalValue = vector.InternalValue[:vector.Size()-1]
	vector.updateProps()
}

func AddVectorBuiltins(vectorRef *VectorValue) {

	vectorScope := NewVectorScope()

	vectorInternalObject := &ObjectValue{
		InternalScope: vectorScope,
		AuxObject:     vectorRef,
	}

	// Register built in functions
	vectorScope.AddFunction("append", &ObjectBuiltInFunction{
		Function: &Function{
			Param: appendParams,
		},
		Object:     vectorInternalObject,
		CustomExec: appendCustomExec,
	})

	vectorScope.AddFunction("remove", &ObjectBuiltInFunction{
		Function: &Function{
			Param: removeParams,
		},
		Object:     vectorInternalObject,
		CustomExec: removeCustomExec,
	})

	vectorScope.AddFunction("removeLast", &ObjectBuiltInFunction{
		Function: &Function{
			Param: removeLastParams,
		},
		Object:     vectorInternalObject,
		CustomExec: removeLastCustomExec,
	})

	// make isEmpty a property
	vectorScope.AddVariable("isEmpty", value.IVOR_BOOL, vectorRef.IsEmpty, true, false, nil)

	// make count a property
	vectorScope.AddVariable("count", value.IVOR_INT, vectorRef.SizeValue, true, false, nil)

	vectorRef.ObjectValue = vectorInternalObject
}
