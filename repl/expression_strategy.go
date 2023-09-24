package repl

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/value"
)

type evalFunc func(value.IVOR, value.IVOR) (bool, string, value.IVOR) // takes 2 values and returns a value
type conversionFunc func(value.IVOR) value.IVOR                       // takes a value and returns a value (different type)

type BinaryValidation struct {
	LeftType        string // allowed left type
	RightType       string // allowed right type
	LeftConversion  conversionFunc
	RightConversion conversionFunc
	Eval            evalFunc
}

type BinaryStrategy struct {
	Name        string
	Validations []BinaryValidation
	Viceversa   bool // if true, the validation is also performed in the opposite order
	DefaultEval evalFunc
}

func (s *BinaryStrategy) Validate(left, right value.IVOR) (bool, string, value.IVOR) {

	// nil in any side is, by default return nil

	if left.Type() == value.IVOR_NIL || right.Type() == value.IVOR_NIL {
		return false, "No es posible realizar operaciones con valores nulos", value.DefaultNilValue
	}

	for _, valid := range s.Validations {

		if valid.LeftType == left.Type() && valid.RightType == right.Type() {

			if valid.LeftConversion != nil {
				left = valid.LeftConversion(left)
			}

			if valid.RightConversion != nil {
				right = valid.RightConversion(right)
			}

			if valid.Eval != nil {
				return valid.Eval(left, right)
			}

			return s.DefaultEval(left, right)
		}

		if s.Viceversa && valid.LeftType == right.Type() && valid.RightType == left.Type() {

			if valid.LeftConversion != nil {
				right = valid.LeftConversion(right)
			}

			if valid.RightConversion != nil {
				left = valid.RightConversion(left)
			}

			if valid.Eval != nil {
				return valid.Eval(left, right)
			}

			return s.DefaultEval(left, right)
		}

	}

	msg := "No es posible realizar la operación '" + s.Name + "' con los tipos '" + left.Type() + "' y '" + right.Type() + "'"

	return false, msg, value.DefaultNilValue
}

// * arithmetic operators

// int + int; float + float; float + int (viceversa); string + string
var addStrategy = BinaryStrategy{
	Name:        "+",
	Viceversa:   true,
	DefaultEval: nil,
	Validations: []BinaryValidation{
		{
			LeftType:        value.IVOR_INT,
			RightType:       value.IVOR_INT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.IntValue{
					InternalValue: left.(*value.IntValue).InternalValue + right.(*value.IntValue).InternalValue,
				}
			},
		},
		{
			LeftType:        value.IVOR_FLOAT,
			RightType:       value.IVOR_FLOAT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.FloatValue{
					InternalValue: left.(*value.FloatValue).InternalValue + right.(*value.FloatValue).InternalValue,
				}
			},
		},
		{
			LeftType:       value.IVOR_FLOAT,
			RightType:      value.IVOR_INT,
			LeftConversion: nil,
			RightConversion: func(v value.IVOR) value.IVOR {
				return &value.FloatValue{
					InternalValue: float64(v.(*value.IntValue).InternalValue),
				}
			},
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.FloatValue{
					InternalValue: left.(*value.FloatValue).InternalValue + right.(*value.FloatValue).InternalValue,
				}
			},
		},
		{
			LeftType:        value.IVOR_STRING,
			RightType:       value.IVOR_STRING,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.StringValue{
					InternalValue: left.(*value.StringValue).InternalValue + right.(*value.StringValue).InternalValue,
				}
			},
		},
		{
			LeftType:  value.IVOR_CHARACTER,
			RightType: value.IVOR_STRING,
			LeftConversion: func(v value.IVOR) value.IVOR {
				return &value.StringValue{
					InternalValue: string(v.(*value.CharacterValue).InternalValue),
				}
			},
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.StringValue{
					InternalValue: left.(*value.StringValue).InternalValue + right.(*value.StringValue).InternalValue,
				}
			},
		},
	},
}

// int - int; float - float; float - int (viceversa)
var subStrategy = BinaryStrategy{
	Name:        "-",
	Viceversa:   true,
	DefaultEval: nil,
	Validations: []BinaryValidation{
		{
			LeftType:        value.IVOR_INT,
			RightType:       value.IVOR_INT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.IntValue{
					InternalValue: left.(*value.IntValue).InternalValue * right.(*value.IntValue).InternalValue,
				}
			},
		},
		{
			LeftType:        value.IVOR_FLOAT,
			RightType:       value.IVOR_FLOAT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.FloatValue{
					InternalValue: left.(*value.FloatValue).InternalValue - right.(*value.FloatValue).InternalValue,
				}
			},
		},
		{
			LeftType:       value.IVOR_FLOAT,
			RightType:      value.IVOR_INT,
			LeftConversion: nil,
			RightConversion: func(v value.IVOR) value.IVOR {
				return &value.FloatValue{
					InternalValue: float64(v.(*value.IntValue).InternalValue),
				}
			},
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.FloatValue{
					InternalValue: left.(*value.FloatValue).InternalValue - right.(*value.FloatValue).InternalValue,
				}
			},
		},
	},
}

// int * int; float * float; float * int (viceversa)
var mulStrategy = BinaryStrategy{
	Name:        "*",
	Viceversa:   true,
	DefaultEval: nil,
	Validations: []BinaryValidation{
		{
			LeftType:        value.IVOR_INT,
			RightType:       value.IVOR_INT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.IntValue{
					InternalValue: left.(*value.IntValue).InternalValue * right.(*value.IntValue).InternalValue,
				}
			},
		},
		{
			LeftType:        value.IVOR_FLOAT,
			RightType:       value.IVOR_FLOAT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.FloatValue{
					InternalValue: left.(*value.FloatValue).InternalValue * right.(*value.FloatValue).InternalValue,
				}
			},
		},
		{
			LeftType:       value.IVOR_FLOAT,
			RightType:      value.IVOR_INT,
			LeftConversion: nil,
			RightConversion: func(v value.IVOR) value.IVOR {
				return &value.FloatValue{
					InternalValue: float64(v.(*value.IntValue).InternalValue),
				}
			},
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.FloatValue{
					InternalValue: left.(*value.FloatValue).InternalValue * right.(*value.FloatValue).InternalValue,
				}
			},
		},
	},
}

// int / int; float / float; float / int (viceversa) !division by zero
var divStrategy = BinaryStrategy{
	Name:        "/",
	Viceversa:   true,
	DefaultEval: nil,
	Validations: []BinaryValidation{
		{
			LeftType:        value.IVOR_INT,
			RightType:       value.IVOR_INT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {

				if right.(*value.IntValue).InternalValue == 0 {
					return false, "No se puede dividir entre cero", value.DefaultNilValue
				}

				return true, "", &value.IntValue{
					InternalValue: left.(*value.IntValue).InternalValue / right.(*value.IntValue).InternalValue,
				}
			},
		},
		{
			LeftType:        value.IVOR_FLOAT,
			RightType:       value.IVOR_FLOAT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {

				if right.(*value.FloatValue).InternalValue == 0 {
					return false, "No se puede dividir entre cero", value.DefaultNilValue
				}

				return true, "", &value.FloatValue{
					InternalValue: left.(*value.FloatValue).InternalValue / right.(*value.FloatValue).InternalValue,
				}
			},
		},
		{
			LeftType:       value.IVOR_FLOAT,
			RightType:      value.IVOR_INT,
			LeftConversion: nil,
			RightConversion: func(v value.IVOR) value.IVOR {
				return &value.FloatValue{
					InternalValue: float64(v.(*value.IntValue).InternalValue),
				}
			},
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {

				if right.(*value.FloatValue).InternalValue == 0 {
					return false, "No se puede dividir entre cero", value.DefaultNilValue
				}

				return true, "", &value.FloatValue{
					InternalValue: left.(*value.FloatValue).InternalValue / right.(*value.FloatValue).InternalValue,
				}
			},
		},
	},
}

// int % int; !division by zero
var modStrategy = BinaryStrategy{
	Name:        "%",
	Viceversa:   true,
	DefaultEval: nil,
	Validations: []BinaryValidation{
		{
			LeftType:        value.IVOR_INT,
			RightType:       value.IVOR_INT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {

				if right.(*value.IntValue).InternalValue == 0 {
					return false, "No se puede dividir entre cero", value.DefaultNilValue
				}

				return true, "", &value.IntValue{
					InternalValue: left.(*value.IntValue).InternalValue % right.(*value.IntValue).InternalValue,
				}
			},
		},
	},
}

// * comparison operators

// int == int; float == float; bool == bool; string == string; char == char
func sameTypeStrat(name string, eval evalFunc) BinaryStrategy {
	return BinaryStrategy{
		Name:        name,
		Viceversa:   true,
		DefaultEval: eval,
		Validations: []BinaryValidation{
			{
				LeftType:        value.IVOR_INT,
				RightType:       value.IVOR_INT,
				LeftConversion:  nil,
				RightConversion: nil,
				Eval:            nil,
			},
			{
				LeftType:        value.IVOR_FLOAT,
				RightType:       value.IVOR_FLOAT,
				LeftConversion:  nil,
				RightConversion: nil,
				Eval:            nil,
			},
			{
				LeftType:        value.IVOR_BOOL,
				RightType:       value.IVOR_BOOL,
				LeftConversion:  nil,
				RightConversion: nil,
				Eval:            nil,
			},
			{
				LeftType:        value.IVOR_STRING,
				RightType:       value.IVOR_STRING,
				LeftConversion:  nil,
				RightConversion: nil,
				Eval:            nil,
			},
			{
				LeftType:        value.IVOR_CHARACTER,
				RightType:       value.IVOR_CHARACTER,
				LeftConversion:  nil,
				RightConversion: nil,
				Eval:            nil,
			},
		},
	}
}

var eqStrategy = sameTypeStrat("==", func(left, right value.IVOR) (bool, string, value.IVOR) {
	return true, "", &value.BoolValue{
		InternalValue: left.Value() == right.Value(),
	}
})

var notEqStrategy = sameTypeStrat("!=", func(left, right value.IVOR) (bool, string, value.IVOR) {
	return true, "", &value.BoolValue{
		InternalValue: left.Value() != right.Value(),
	}
})

var lessThanStrategy = BinaryStrategy{
	Name:        "<",
	Viceversa:   true,
	DefaultEval: nil,
	Validations: []BinaryValidation{
		{
			LeftType:        value.IVOR_INT,
			RightType:       value.IVOR_INT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.BoolValue{
					InternalValue: left.(*value.IntValue).InternalValue < right.(*value.IntValue).InternalValue,
				}
			},
		},
		{
			LeftType:        value.IVOR_FLOAT,
			RightType:       value.IVOR_FLOAT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.BoolValue{
					InternalValue: left.(*value.FloatValue).InternalValue < right.(*value.FloatValue).InternalValue,
				}
			},
		},
		{
			LeftType:        value.IVOR_STRING,
			RightType:       value.IVOR_STRING,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.BoolValue{
					InternalValue: left.(*value.StringValue).InternalValue < right.(*value.StringValue).InternalValue,
				}
			},
		},
		{
			LeftType:        value.IVOR_CHARACTER,
			RightType:       value.IVOR_CHARACTER,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.BoolValue{
					InternalValue: left.(*value.CharacterValue).InternalValue < right.(*value.CharacterValue).InternalValue,
				}
			},
		},
	},
}

var lessOrEqStrategy = BinaryStrategy{
	Name:        "<=",
	Viceversa:   true,
	DefaultEval: nil,
	Validations: []BinaryValidation{
		{
			LeftType:        value.IVOR_INT,
			RightType:       value.IVOR_INT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.BoolValue{
					InternalValue: left.(*value.IntValue).InternalValue <= right.(*value.IntValue).InternalValue,
				}
			},
		},
		{
			LeftType:        value.IVOR_FLOAT,
			RightType:       value.IVOR_FLOAT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.BoolValue{
					InternalValue: left.(*value.FloatValue).InternalValue <= right.(*value.FloatValue).InternalValue,
				}
			},
		},
		{
			LeftType:        value.IVOR_STRING,
			RightType:       value.IVOR_STRING,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.BoolValue{
					InternalValue: left.(*value.StringValue).InternalValue <= right.(*value.StringValue).InternalValue,
				}
			},
		},
		{
			LeftType:        value.IVOR_CHARACTER,
			RightType:       value.IVOR_CHARACTER,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.BoolValue{
					InternalValue: left.(*value.CharacterValue).InternalValue <= right.(*value.CharacterValue).InternalValue,
				}
			},
		},
	},
}

var greaterThanStrategy = BinaryStrategy{
	Name:        ">",
	Viceversa:   true,
	DefaultEval: nil,
	Validations: []BinaryValidation{
		{
			LeftType:        value.IVOR_INT,
			RightType:       value.IVOR_INT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.BoolValue{
					InternalValue: left.(*value.IntValue).InternalValue > right.(*value.IntValue).InternalValue,
				}
			},
		},
		{
			LeftType:        value.IVOR_FLOAT,
			RightType:       value.IVOR_FLOAT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.BoolValue{
					InternalValue: left.(*value.FloatValue).InternalValue > right.(*value.FloatValue).InternalValue,
				}
			},
		},
		{
			LeftType:        value.IVOR_STRING,
			RightType:       value.IVOR_STRING,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.BoolValue{
					InternalValue: left.(*value.StringValue).InternalValue > right.(*value.StringValue).InternalValue,
				}
			},
		},
		{
			LeftType:        value.IVOR_CHARACTER,
			RightType:       value.IVOR_CHARACTER,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.BoolValue{
					InternalValue: left.(*value.CharacterValue).InternalValue > right.(*value.CharacterValue).InternalValue,
				}
			},
		},
	},
}

var greaterOrEqStrategy = BinaryStrategy{
	Name:        ">=",
	Viceversa:   true,
	DefaultEval: nil,
	Validations: []BinaryValidation{
		{
			LeftType:        value.IVOR_INT,
			RightType:       value.IVOR_INT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.BoolValue{
					InternalValue: left.(*value.IntValue).InternalValue >= right.(*value.IntValue).InternalValue,
				}
			},
		},
		{
			LeftType:        value.IVOR_FLOAT,
			RightType:       value.IVOR_FLOAT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.BoolValue{
					InternalValue: left.(*value.FloatValue).InternalValue >= right.(*value.FloatValue).InternalValue,
				}
			},
		},
		{
			LeftType:        value.IVOR_STRING,
			RightType:       value.IVOR_STRING,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.BoolValue{
					InternalValue: left.(*value.StringValue).InternalValue >= right.(*value.StringValue).InternalValue,
				}
			},
		},
		{
			LeftType:        value.IVOR_CHARACTER,
			RightType:       value.IVOR_CHARACTER,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(left, right value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.BoolValue{
					InternalValue: left.(*value.CharacterValue).InternalValue >= right.(*value.CharacterValue).InternalValue,
				}
			},
		},
	},
}

// * logical operators

func genericBinaryLogicalStrat(name string, eval evalFunc) BinaryStrategy {

	return BinaryStrategy{
		Name:        name,
		Viceversa:   true,
		DefaultEval: eval,
		Validations: []BinaryValidation{
			{
				LeftType:        value.IVOR_BOOL,
				RightType:       value.IVOR_BOOL,
				LeftConversion:  nil,
				RightConversion: nil,
				Eval:            nil,
			},
		},
	}
}

var andStrategy = genericBinaryLogicalStrat("&&", func(left, right value.IVOR) (bool, string, value.IVOR) {
	return true, "", &value.BoolValue{
		InternalValue: left.(*value.BoolValue).InternalValue && right.(*value.BoolValue).InternalValue,
	}
})

var orStrategy = genericBinaryLogicalStrat("||", func(left, right value.IVOR) (bool, string, value.IVOR) {
	return true, "", &value.BoolValue{
		InternalValue: left.(*value.BoolValue).InternalValue || right.(*value.BoolValue).InternalValue,
	}
})

var BinaryStrats = map[string]BinaryStrategy{
	"+":  addStrategy,
	"-":  subStrategy,
	"*":  mulStrategy,
	"/":  divStrategy,
	"%":  modStrategy,
	"==": eqStrategy,
	"!=": notEqStrategy,
	"<":  lessThanStrategy,
	"<=": lessOrEqStrategy,
	">":  greaterThanStrategy,
	">=": greaterOrEqStrategy,
	"&&": andStrategy,
	"||": orStrategy,
}

// UnaryStrats

type UnaryValidation struct {
	Type       string // allowed type
	Conversion conversionFunc
	Eval       evalFunc
}

type UnaryStrategy struct {
	Name        string
	Validations []UnaryValidation
	DefaultEval evalFunc
}

func (s *UnaryStrategy) Validate(val value.IVOR) (bool, string, value.IVOR) {

	if val.Type() == value.IVOR_NIL {
		return false, "No es posible realizar operaciones con valores nulos", value.DefaultNilValue
	}

	for _, valid := range s.Validations {

		if valid.Type == val.Type() {

			if valid.Conversion != nil {
				val = valid.Conversion(val)
			}

			if valid.Eval != nil {
				return valid.Eval(val, nil)
			}

			return s.DefaultEval(val, nil)
		}

	}

	msg := "No es posible realizar la operación '" + s.Name + "' con el tipo '" + val.Type() + "'"

	return false, msg, value.DefaultNilValue
}

// * Not

var notStrategy = UnaryStrategy{
	Name:        "!",
	DefaultEval: nil,
	Validations: []UnaryValidation{
		{
			Type:       value.IVOR_BOOL,
			Conversion: nil,
			Eval: func(i1, i2 value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.BoolValue{
					InternalValue: !i1.(*value.BoolValue).InternalValue,
				}
			},
		},
	},
}

// * Minus

var minusStrategy = UnaryStrategy{
	Name:        "-",
	DefaultEval: nil,
	Validations: []UnaryValidation{
		{
			Type:       value.IVOR_INT,
			Conversion: nil,
			Eval: func(i1, i2 value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.IntValue{
					InternalValue: -i1.(*value.IntValue).InternalValue,
				}
			},
		},
		{
			Type:       value.IVOR_FLOAT,
			Conversion: nil,
			Eval: func(i1, i2 value.IVOR) (bool, string, value.IVOR) {
				return true, "", &value.FloatValue{
					InternalValue: -i1.(*value.FloatValue).InternalValue,
				}
			},
		},
	},
}

var UnaryStrats = map[string]UnaryStrategy{
	"!": notStrategy,
	"-": minusStrategy,
}

// Early return strats

// * And

var andEarlyReturnStrategy = UnaryStrategy{
	Name: "&&",
	Validations: []UnaryValidation{
		{
			Type:       value.IVOR_BOOL,
			Conversion: nil,
			Eval: func(i1, i2 value.IVOR) (bool, string, value.IVOR) {

				if !i1.(*value.BoolValue).InternalValue {
					return true, "", &value.BoolValue{
						InternalValue: false,
					}
				}

				return false, "", nil
			},
		},
	},
}

// * Or

var orEarlyReturnStrategy = UnaryStrategy{
	Name: "||",
	Validations: []UnaryValidation{
		{
			Type:       value.IVOR_BOOL,
			Conversion: nil,
			Eval: func(i1, i2 value.IVOR) (bool, string, value.IVOR) {

				if i1.(*value.BoolValue).InternalValue {
					return true, "", &value.BoolValue{
						InternalValue: true,
					}
				}

				return false, "", nil
			},
		},
	},
}

var EarlyReturnStrats = map[string]UnaryStrategy{
	"&&": andEarlyReturnStrategy,
	"||": orEarlyReturnStrategy,
}
