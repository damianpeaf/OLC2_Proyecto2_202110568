package visitor

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/abstract"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/tac"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/value"
)

type evalFunc func(*IrVisitor, *value.ValueWrapper, *value.ValueWrapper) (bool, *value.ValueWrapper) // takes 2 values and returns a value
type conversionFunc func(*value.ValueWrapper)                                                        // takes a value and returns a value (different type)

type BinaryValidation struct {
	LeftType        string // allowed left type
	RightType       string // allowed right type
	LeftConversion  conversionFunc
	RightConversion conversionFunc
	Eval            evalFunc
	ReturnType      string // if the return type is different from the left/right type
}

type BinaryStrategy struct {
	Name        string
	Validations []BinaryValidation
	Viceversa   bool // if true, the validation is also performed in the opposite order
	DefaultEval evalFunc
	v           *IrVisitor
}

func (s *BinaryStrategy) Validate(left, right *value.ValueWrapper) (bool, *value.ValueWrapper) {

	for _, valid := range s.Validations {

		if valid.LeftType == left.Metadata && valid.RightType == right.Metadata {

			if valid.LeftConversion != nil {
				valid.LeftConversion(left)
			}

			if valid.RightConversion != nil {
				valid.RightConversion(right)
			}

			if valid.Eval != nil {
				return valid.Eval(s.v, left, right)
			}

			ok, dvw := s.DefaultEval(s.v, left, right)
			dvw.Metadata = valid.ReturnType
			return ok, dvw
		}

		if s.Viceversa && valid.LeftType == right.Metadata && valid.RightType == left.Metadata {

			if valid.LeftConversion != nil {
				valid.LeftConversion(right)
			}

			if valid.RightConversion != nil {
				valid.RightConversion(left)
			}

			if valid.Eval != nil {
				return valid.Eval(s.v, left, right)
			}

			ok, dvw := s.DefaultEval(s.v, left, right)
			dvw.Metadata = valid.ReturnType
			return ok, dvw
		}

	}

	return false, &value.ValueWrapper{
		Val:      s.v.Utility.NilValue(),
		Metadata: abstract.IVOR_NIL,
	}
}

// * arithmetic operators

// int + int; float + float; float + int (viceversa); string + string
var addStrategy = BinaryStrategy{
	Name:      "+",
	Viceversa: true,
	DefaultEval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
		return true, &value.ValueWrapper{
			Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.PLUS),
			Metadata: "",
		}
	},
	Validations: []BinaryValidation{
		{
			LeftType:        abstract.IVOR_INT,
			RightType:       abstract.IVOR_INT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval:            nil,
			ReturnType:      abstract.IVOR_INT,
		},
		{
			LeftType:        abstract.IVOR_FLOAT,
			RightType:       abstract.IVOR_FLOAT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval:            nil,
			ReturnType:      abstract.IVOR_FLOAT,
		},
		{
			LeftType:       abstract.IVOR_FLOAT,
			RightType:      abstract.IVOR_INT,
			LeftConversion: nil,
			RightConversion: func(vw *value.ValueWrapper) {
				tac.AddCastToSimpleValue(vw.Val, tac.CAST_FLOAT)
			},
			Eval:       nil,
			ReturnType: abstract.IVOR_FLOAT,
		},
		{
			LeftType:        abstract.IVOR_STRING,
			RightType:       abstract.IVOR_STRING,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
				return true, &value.ValueWrapper{
					Val:      iv.Utility.ConcatStrings(vw1.Val, vw2.Val),
					Metadata: abstract.IVOR_STRING,
				}
			},
		},
		{
			LeftType:        abstract.IVOR_CHARACTER,
			RightType:       abstract.IVOR_STRING,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
				charFirst := true
				if vw1.Metadata == abstract.IVOR_STRING {
					charFirst = false
				}
				return true, &value.ValueWrapper{
					Val:      iv.Utility.ConcatCharStrings(vw1.Val, vw2.Val, charFirst),
					Metadata: abstract.IVOR_STRING,
				}
			},
		},
	},
}

// int - int; float - float; float - int (viceversa)
var subStrategy = BinaryStrategy{
	Name:      "-",
	Viceversa: true,
	DefaultEval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
		return true, &value.ValueWrapper{
			Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.MINUS),
			Metadata: abstract.IVOR_INT,
		}
	},
	Validations: []BinaryValidation{
		{
			LeftType:        abstract.IVOR_INT,
			RightType:       abstract.IVOR_INT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
				return true, &value.ValueWrapper{
					Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.MINUS),
					Metadata: abstract.IVOR_INT,
				}
			},
		},
		{
			LeftType:        abstract.IVOR_FLOAT,
			RightType:       abstract.IVOR_FLOAT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
				return true, &value.ValueWrapper{
					Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.MINUS),
					Metadata: abstract.IVOR_FLOAT,
				}
			},
		},
		{
			LeftType:       abstract.IVOR_FLOAT,
			RightType:      abstract.IVOR_INT,
			LeftConversion: nil,
			RightConversion: func(vw *value.ValueWrapper) {
				tac.AddCastToSimpleValue(vw.Val, tac.CAST_FLOAT)
			},
			Eval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
				return true, &value.ValueWrapper{
					Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.MINUS),
					Metadata: abstract.IVOR_FLOAT,
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
			LeftType:        abstract.IVOR_INT,
			RightType:       abstract.IVOR_INT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
				return true, &value.ValueWrapper{
					Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.MULTIPLY),
					Metadata: abstract.IVOR_INT,
				}
			},
		},
		{
			LeftType:        abstract.IVOR_FLOAT,
			RightType:       abstract.IVOR_FLOAT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
				return true, &value.ValueWrapper{
					Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.MULTIPLY),
					Metadata: abstract.IVOR_FLOAT,
				}
			},
		},
		{
			LeftType:       abstract.IVOR_FLOAT,
			RightType:      abstract.IVOR_INT,
			LeftConversion: nil,
			RightConversion: func(vw *value.ValueWrapper) {
				tac.AddCastToSimpleValue(vw.Val, tac.CAST_FLOAT)
			},
			Eval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
				return true, &value.ValueWrapper{
					Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.MULTIPLY),
					Metadata: abstract.IVOR_FLOAT,
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
			LeftType:        abstract.IVOR_INT,
			RightType:       abstract.IVOR_INT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
				rvw := &value.ValueWrapper{
					Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.DIVIDE),
					Metadata: abstract.IVOR_INT,
				}
				return true, rvw
			},
		},
		{
			LeftType:        abstract.IVOR_FLOAT,
			RightType:       abstract.IVOR_FLOAT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
				rvw := &value.ValueWrapper{
					Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.DIVIDE),
					Metadata: abstract.IVOR_FLOAT,
				}
				return true, rvw
			},
		},
		{
			LeftType:       abstract.IVOR_FLOAT,
			RightType:      abstract.IVOR_INT,
			LeftConversion: nil,
			RightConversion: func(vw *value.ValueWrapper) {
				tac.AddCastToSimpleValue(vw.Val, tac.CAST_FLOAT)
			},
			Eval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
				rvw := &value.ValueWrapper{
					Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.DIVIDE),
					Metadata: abstract.IVOR_FLOAT,
				}
				return true, rvw
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
			LeftType:        abstract.IVOR_INT,
			RightType:       abstract.IVOR_INT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
				rvw := &value.ValueWrapper{
					Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.MOD),
					Metadata: abstract.IVOR_INT,
				}
				return true, rvw
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
				LeftType:        abstract.IVOR_INT,
				RightType:       abstract.IVOR_INT,
				LeftConversion:  nil,
				RightConversion: nil,
				Eval:            nil,
			},
			{
				LeftType:        abstract.IVOR_FLOAT,
				RightType:       abstract.IVOR_FLOAT,
				LeftConversion:  nil,
				RightConversion: nil,
				Eval:            nil,
			},
			{
				LeftType:        abstract.IVOR_BOOL,
				RightType:       abstract.IVOR_BOOL,
				LeftConversion:  nil,
				RightConversion: nil,
				Eval:            nil,
			},
			{
				LeftType:        abstract.IVOR_STRING,
				RightType:       abstract.IVOR_STRING,
				LeftConversion:  nil,
				RightConversion: nil,
				Eval:            nil,
			},
			{
				LeftType:        abstract.IVOR_CHARACTER,
				RightType:       abstract.IVOR_CHARACTER,
				LeftConversion:  nil,
				RightConversion: nil,
				Eval:            nil,
			},
		},
	}
}

// int == int; float == float; bool == bool; string == string; char == char

var defaultEqValidation = func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
	return true, &value.ValueWrapper{
		Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.EQ),
		Metadata: abstract.IVOR_BOOL,
	}
}

var eqStrategy = BinaryStrategy{
	Name:        "==",
	Viceversa:   true,
	DefaultEval: nil,
	Validations: []BinaryValidation{
		{
			LeftType:        abstract.IVOR_INT,
			RightType:       abstract.IVOR_INT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval:            defaultEqValidation,
		},
		{
			LeftType:        abstract.IVOR_FLOAT,
			RightType:       abstract.IVOR_FLOAT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval:            defaultEqValidation,
		},
		{
			LeftType:        abstract.IVOR_BOOL,
			RightType:       abstract.IVOR_BOOL,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval:            defaultEqValidation,
		},
		{
			LeftType:        abstract.IVOR_STRING,
			RightType:       abstract.IVOR_STRING,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
				return true, &value.ValueWrapper{
					Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.EQ), // TODO: char by char comparison
					Metadata: abstract.IVOR_BOOL,
				}
			},
		},
	},
}

var defaultNeqValidation = func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
	return true, &value.ValueWrapper{
		Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.NEQ),
		Metadata: abstract.IVOR_BOOL,
	}
}

var notEqStrategy = BinaryStrategy{
	Name:        "!=",
	Viceversa:   true,
	DefaultEval: nil,
	Validations: []BinaryValidation{
		{
			LeftType:        abstract.IVOR_INT,
			RightType:       abstract.IVOR_INT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval:            defaultNeqValidation,
		},
		{
			LeftType:        abstract.IVOR_FLOAT,
			RightType:       abstract.IVOR_FLOAT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval:            defaultNeqValidation,
		},
		{
			LeftType:        abstract.IVOR_BOOL,
			RightType:       abstract.IVOR_BOOL,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval:            defaultNeqValidation,
		},
		{
			LeftType:        abstract.IVOR_STRING,
			RightType:       abstract.IVOR_STRING,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
				return true, &value.ValueWrapper{
					Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.NEQ), // TODO: char by char comparison
					Metadata: abstract.IVOR_BOOL,
				}
			},
		},
	},
}

var defaultLessThanValidation = func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
	return true, &value.ValueWrapper{
		Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.LT),
		Metadata: abstract.IVOR_BOOL,
	}
}
var lessThanStrategy = BinaryStrategy{
	Name:        "<",
	Viceversa:   true,
	DefaultEval: nil,
	Validations: []BinaryValidation{
		{
			LeftType:        abstract.IVOR_INT,
			RightType:       abstract.IVOR_INT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval:            defaultLessThanValidation,
		},
		{
			LeftType:        abstract.IVOR_FLOAT,
			RightType:       abstract.IVOR_FLOAT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval:            defaultLessThanValidation,
		},
		{
			LeftType:        abstract.IVOR_STRING,
			RightType:       abstract.IVOR_STRING,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
				return true, &value.ValueWrapper{
					Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.LT), // TODO: char by char comparison
					Metadata: abstract.IVOR_BOOL,
				}
			},
		},
		{
			LeftType:        abstract.IVOR_CHARACTER,
			RightType:       abstract.IVOR_CHARACTER,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval:            defaultLessThanValidation,
		},
	},
}

var defaultLessOrEqValidation = func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
	return true, &value.ValueWrapper{
		Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.LTE),
		Metadata: abstract.IVOR_BOOL,
	}
}

var lessOrEqStrategy = BinaryStrategy{
	Name:        "<=",
	Viceversa:   true,
	DefaultEval: nil,
	Validations: []BinaryValidation{
		{
			LeftType:        abstract.IVOR_INT,
			RightType:       abstract.IVOR_INT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval:            defaultLessOrEqValidation,
		},
		{
			LeftType:        abstract.IVOR_FLOAT,
			RightType:       abstract.IVOR_FLOAT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval:            defaultLessOrEqValidation,
		},
		{
			LeftType:        abstract.IVOR_STRING,
			RightType:       abstract.IVOR_STRING,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
				return true, &value.ValueWrapper{
					Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.LTE), // TODO: char by char comparison
					Metadata: abstract.IVOR_BOOL,
				}
			},
		},
		{
			LeftType:        abstract.IVOR_CHARACTER,
			RightType:       abstract.IVOR_CHARACTER,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval:            defaultLessOrEqValidation,
		},
	},
}

var defaultGreaterThanValidation = func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
	return true, &value.ValueWrapper{
		Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.GT),
		Metadata: abstract.IVOR_BOOL,
	}
}

var greaterThanStrategy = BinaryStrategy{
	Name:        ">",
	Viceversa:   true,
	DefaultEval: nil,
	Validations: []BinaryValidation{
		{
			LeftType:        abstract.IVOR_INT,
			RightType:       abstract.IVOR_INT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval:            defaultGreaterThanValidation,
		},
		{
			LeftType:        abstract.IVOR_FLOAT,
			RightType:       abstract.IVOR_FLOAT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval:            defaultGreaterThanValidation,
		},
		{
			LeftType:        abstract.IVOR_STRING,
			RightType:       abstract.IVOR_STRING,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
				return true, &value.ValueWrapper{
					Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.GT), // TODO: char by char comparison
					Metadata: abstract.IVOR_BOOL,
				}
			},
		},
		{
			LeftType:        abstract.IVOR_CHARACTER,
			RightType:       abstract.IVOR_CHARACTER,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval:            defaultGreaterThanValidation,
		},
	},
}

var defaultGreaterOrEqValidation = func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
	return true, &value.ValueWrapper{
		Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.GTE),
		Metadata: abstract.IVOR_BOOL,
	}
}

var greaterOrEqStrategy = BinaryStrategy{
	Name:        ">=",
	Viceversa:   true,
	DefaultEval: nil,
	Validations: []BinaryValidation{
		{
			LeftType:        abstract.IVOR_INT,
			RightType:       abstract.IVOR_INT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval:            defaultGreaterOrEqValidation,
		},
		{
			LeftType:        abstract.IVOR_FLOAT,
			RightType:       abstract.IVOR_FLOAT,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval:            defaultGreaterOrEqValidation,
		},
		{
			LeftType:        abstract.IVOR_STRING,
			RightType:       abstract.IVOR_STRING,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
				return true, &value.ValueWrapper{
					Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.GTE), // TODO: char by char comparison
					Metadata: abstract.IVOR_BOOL,
				}
			},
		},
		{
			LeftType:        abstract.IVOR_CHARACTER,
			RightType:       abstract.IVOR_CHARACTER,
			LeftConversion:  nil,
			RightConversion: nil,
			Eval:            defaultGreaterOrEqValidation,
		},
	},
}

// // * logical operators

func genericBinaryLogicalStrat(name string, eval evalFunc) BinaryStrategy {

	return BinaryStrategy{
		Name:        name,
		Viceversa:   true,
		DefaultEval: eval,
		Validations: []BinaryValidation{
			{
				LeftType:        abstract.IVOR_BOOL,
				RightType:       abstract.IVOR_BOOL,
				LeftConversion:  nil,
				RightConversion: nil,
				Eval:            nil,
			},
		},
	}
}

var andStrategy = genericBinaryLogicalStrat("&&", func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
	return true, &value.ValueWrapper{
		Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.PLUS), // TODO: and built-in function
		Metadata: abstract.IVOR_BOOL,
	}
})

var orStrategy = genericBinaryLogicalStrat("||", func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
	return true, &value.ValueWrapper{
		Val:      iv.Utility.BasicOperation(vw1.Val, vw2.Val, tac.PLUS), // TODO: or built-in function
		Metadata: abstract.IVOR_BOOL,
	}
})

func NewBinaryStrats(v *IrVisitor) map[string]BinaryStrategy {

	addStrategy.v = v
	subStrategy.v = v
	mulStrategy.v = v
	divStrategy.v = v
	modStrategy.v = v
	eqStrategy.v = v
	notEqStrategy.v = v
	lessThanStrategy.v = v
	lessOrEqStrategy.v = v
	greaterThanStrategy.v = v
	greaterOrEqStrategy.v = v
	andStrategy.v = v
	orStrategy.v = v

	return map[string]BinaryStrategy{
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
	v           *IrVisitor
}

func (s *UnaryStrategy) Validate(val *value.ValueWrapper) (bool, *value.ValueWrapper) {

	for _, valid := range s.Validations {

		if valid.Type == val.Metadata {

			if valid.Conversion != nil {
				valid.Conversion(val)
			}

			if valid.Eval != nil {
				return valid.Eval(s.v, val, nil)
			}

			return s.DefaultEval(s.v, val, nil)
		}

	}

	return false, &value.ValueWrapper{
		Val:      s.v.Utility.NilValue(),
		Metadata: abstract.IVOR_NIL,
	}
}

// * Not
var notStrategy = UnaryStrategy{
	Name:        "!",
	DefaultEval: nil,
	Validations: []UnaryValidation{
		{
			Type:       abstract.IVOR_BOOL,
			Conversion: nil,
			Eval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
				return true, &value.ValueWrapper{
					Val:      iv.Utility.BasicOperation(vw1.Val, nil, tac.PLUS), // TODO: not built-in function
					Metadata: abstract.IVOR_BOOL,
				}
			},
		},
	},
}

// * Minus
var minusStrategy = UnaryStrategy{
	Name: "-",
	DefaultEval: func(iv *IrVisitor, vw1, vw2 *value.ValueWrapper) (bool, *value.ValueWrapper) {
		return true, &value.ValueWrapper{
			Val:      iv.Utility.BasicOperation(iv.Factory.NewLiteral().SetValue("0"), vw1.Val, tac.MINUS),
			Metadata: abstract.IVOR_INT,
		}
	},
	Validations: []UnaryValidation{
		{
			Type:       abstract.IVOR_INT,
			Conversion: nil,
			Eval:       nil,
		},
		{
			Type:       abstract.IVOR_FLOAT,
			Conversion: nil,
			Eval:       nil,
		},
	},
}

func NewUnaryStrats(v *IrVisitor) map[string]UnaryStrategy {

	notStrategy.v = v
	minusStrategy.v = v

	return map[string]UnaryStrategy{
		"!": notStrategy,
		"-": minusStrategy,
	}
}
