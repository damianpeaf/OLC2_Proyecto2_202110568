package repl

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/value"
	"strconv"
)

type BuiltInFunction struct {
	Name string
	Exec func(context *ReplContext, args []*Argument) (value.IVOR, bool, string)
}

// implementing ivor

func (b BuiltInFunction) Type() string {
	return value.IVOR_BUILTIN_FUNCTION
}

func (b BuiltInFunction) Value() interface{} {
	return b
}

func (b BuiltInFunction) Copy() value.IVOR {
	return b
}

// * Print Function
func Print(context *ReplContext, args []*Argument) (value.IVOR, bool, string) {

	var output string

	for i, arg := range args {

		if !value.IsPrimitiveType(arg.Value.Type()) {
			return value.DefaultNilValue, false, "La función print solo acepta tipos primitivos"
		}

		switch arg.Value.Type() {

		case value.IVOR_BOOL:
			output += strconv.FormatBool(arg.Value.Value().(bool))
		case value.IVOR_INT:
			output += strconv.Itoa(arg.Value.Value().(int))
		case value.IVOR_FLOAT:
			output += strconv.FormatFloat(arg.Value.Value().(float64), 'f', 4, 64) // 4 digits of precision
		case value.IVOR_STRING:
			output += arg.Value.Value().(string)
		case value.IVOR_CHARACTER:
			output += arg.Value.Value().(string)
		case value.IVOR_NIL:
			output += "nil"
		}

		// Add a space between each argument
		if i < len(args)-1 {
			output += " "
		}
	}

	context.Console.Print(output)

	return value.DefaultNilValue, true, ""
}

// * Int Function

func Int(context *ReplContext, args []*Argument) (value.IVOR, bool, string) {

	if len(args) != 1 {
		return value.DefaultNilValue, false, "La función int solo acepta un argumento"
	}

	argValue := args[0].Value

	if !(argValue.Type() == value.IVOR_STRING || argValue.Type() == value.IVOR_FLOAT) {
		return value.DefaultNilValue, false, "La función Int solo acepta un argumento de tipo string o float"
	}

	if argValue.Type() == value.IVOR_STRING {
		floatValue, err := strconv.ParseFloat(argValue.Value().(string), 64)

		if err != nil {
			return value.DefaultNilValue, false, "No se pudo convertir el valor a int"
		}

		return &value.IntValue{
			InternalValue: int(floatValue),
		}, true, ""
	}

	if argValue.Type() == value.IVOR_FLOAT {
		// truncate the float

		floatValue := argValue.Value().(float64)

		return &value.IntValue{
			InternalValue: int(floatValue),
		}, true, ""
	}

	return value.DefaultNilValue, false, "No se pudo convertir el valor a int"
}

// * Float Function

func Float(context *ReplContext, args []*Argument) (value.IVOR, bool, string) {

	if len(args) != 1 {
		return value.DefaultNilValue, false, "La función float solo acepta un argumento"
	}

	argValue := args[0].Value

	if !(argValue.Type() == value.IVOR_STRING) {
		return value.DefaultNilValue, false, "La función float solo acepta un argumento de tipo string"
	}

	floatValue, err := strconv.ParseFloat(argValue.Value().(string), 64)

	if err != nil {
		return value.DefaultNilValue, false, "No se pudo convertir el valor a float"
	}

	return &value.FloatValue{
		InternalValue: floatValue,
	}, true, ""
}

// * String Function

func String(context *ReplContext, args []*Argument) (value.IVOR, bool, string) {

	if len(args) != 1 {
		return value.DefaultNilValue, false, "La función string solo acepta un argumento"
	}

	argValue := args[0].Value

	if !(argValue.Type() == value.IVOR_INT || argValue.Type() == value.IVOR_FLOAT || argValue.Type() == value.IVOR_BOOL) {
		return value.DefaultNilValue, false, "La función string solo acepta un argumento de tipo int, float o bool"
	}

	if argValue.Type() == value.IVOR_INT {
		stringValue := strconv.Itoa(argValue.Value().(int))

		return &value.StringValue{
			InternalValue: stringValue,
		}, true, ""
	}

	if argValue.Type() == value.IVOR_FLOAT {
		stringValue := strconv.FormatFloat(argValue.Value().(float64), 'f', 4, 64)

		return &value.StringValue{
			InternalValue: stringValue,
		}, true, ""
	}

	if argValue.Type() == value.IVOR_BOOL {
		stringValue := strconv.FormatBool(argValue.Value().(bool))

		return &value.StringValue{
			InternalValue: stringValue,
		}, true, ""
	}

	return value.DefaultNilValue, false, "No se pudo convertir el valor a string"
}

var DefaultBuiltInFunctions = map[string]*BuiltInFunction{
	"print": {
		Name: "print",
		Exec: Print,
	},
	"Int": {
		Name: "Int",
		Exec: Int,
	},
	"Float": {
		Name: "Float",
		Exec: Float,
	},
	"String": {
		Name: "String",
		Exec: String,
	},
}
