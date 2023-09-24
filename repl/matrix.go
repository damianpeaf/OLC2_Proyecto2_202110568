package repl

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/value"
)

type MatrixValue struct {
	*VectorValue
}

func (v *MatrixValue) Copy() value.IVOR {

	internalCopy := make([]value.IVOR, len(v.InternalValue))

	for i, item := range v.InternalValue {
		internalCopy[i] = item.Copy()
	}

	return NewMatrixValue(internalCopy, v.FullType, v.ItemType)

}

func (v *MatrixValue) ValidIndexes(indexes []int) bool {

	// check if indexes are valid
	pivot := v.VectorValue

	for i, index := range indexes {
		if index < 0 || index >= pivot.Size() {
			return false
		}

		item := pivot.Get(index)

		// vector, matrix or value
		switch s := item.(type) {
		case *VectorValue:
			pivot = s

			if i == len(indexes)-1 {
				return true
			}

		case *MatrixValue:
			pivot = s.VectorValue

			if i == len(indexes)-1 {
				return true
			}

		default:
			if i != len(indexes)-1 {
				return false
			} else {
				return true
			}
		}
	}

	return false
}

func (v *MatrixValue) Get(index []int) value.IVOR {

	// check if indexes are valid
	if !v.ValidIndexes(index) {
		return nil
	}

	pivot := v.VectorValue

	for i := 0; i < len(index); i++ {
		item := pivot.Get(index[i])

		// vector, matrix or value
		switch s := item.(type) {
		case *VectorValue:
			pivot = s

			if i == len(index)-1 {
				return pivot
			}

		case *MatrixValue:
			pivot = s.VectorValue

			if i == len(index)-1 {
				return pivot
			}
		default:
			return item
		}
	}

	return nil
}

func (v *MatrixValue) Set(index []int, value value.IVOR) bool {

	// check if indexes are valid
	if !v.ValidIndexes(index) {
		return false
	}

	pivot := v.VectorValue

	for i := 0; i < len(index); i++ {
		item := pivot.Get(index[i])

		// vector, matrix or value
		switch s := item.(type) {
		case *VectorValue:
			pivot = s
		case *MatrixValue:
			pivot = s.VectorValue
		default:
			if i == len(index)-1 {
				pivot.InternalValue[index[i]] = value
				return true
			}
		}
	}

	return false
}

func NewMatrixValue(vectorItems []value.IVOR, fullType, itemType string) *MatrixValue {
	vector := &VectorValue{
		InternalValue: vectorItems,
		CurrentIndex:  0,
		ItemType:      itemType,
		FullType:      fullType,
		SizeValue:     &value.IntValue{InternalValue: len(vectorItems)},
		IsEmpty:       &value.BoolValue{InternalValue: len(vectorItems) == 0},
	}

	// remove builtins from vector value
	removeBuiltinsFromVector(vectorItems)

	return &MatrixValue{
		VectorValue: vector,
	}
}

func removeBuiltinsFromVector(vectorItems []value.IVOR) {
	for i := 0; i < len(vectorItems); i++ {
		if item, ok := vectorItems[i].(*VectorValue); ok {
			item.ObjectValue.InternalScope.Reset()
			// removeBuiltinsFromVector(item.InternalValue)
		} else {
			break
		}
	}
}

type MatrixItemReference struct {
	Matrix *MatrixValue
	Index  []int
	Value  value.IVOR
}
