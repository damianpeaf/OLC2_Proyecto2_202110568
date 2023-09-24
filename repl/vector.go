package repl

import "github.com/damianpeaf/OLC2_Proyecto2_202110568/value"

// implements ivor interface

type VectorValue struct {
	*ObjectValue
	InternalValue []value.IVOR
	CurrentIndex  int
	ItemType      string
	FullType      string
	SizeValue     *value.IntValue
	IsEmpty       *value.BoolValue
}

func (v VectorValue) Value() interface{} {
	return v
}

func (v VectorValue) Type() string {
	return v.FullType
}

func (v VectorValue) Size() int {
	return len(v.InternalValue)
}

func (v VectorValue) ValidIndex(index int) bool {

	if index < 0 || index >= len(v.InternalValue) {
		return false
	}

	return true

}

func (v VectorValue) Get(index int) value.IVOR {
	return v.InternalValue[index]
}

func (v *VectorValue) Next() bool {
	if v.CurrentIndex < len(v.InternalValue) {
		v.CurrentIndex++
		return true
	}
	return false
}

func (v *VectorValue) Current() value.IVOR {
	return v.InternalValue[v.CurrentIndex]
}

func (v *VectorValue) Reset() {
	v.CurrentIndex = 0
}

func (v *VectorValue) Copy() value.IVOR {

	internalCopy := make([]value.IVOR, len(v.InternalValue))

	for i, item := range v.InternalValue {
		internalCopy[i] = item.Copy()
	}

	return NewVectorValue(internalCopy, v.FullType, v.ItemType)

}

func (v *VectorValue) updateProps() {

	v.SizeValue.InternalValue = len(v.InternalValue)
	v.IsEmpty.InternalValue = len(v.InternalValue) == 0

}

func NewVectorValue(vectorItems []value.IVOR, fullType, itemType string) *VectorValue {
	vector := &VectorValue{
		InternalValue: vectorItems,
		CurrentIndex:  0,
		ItemType:      itemType,
		FullType:      fullType,
		SizeValue:     &value.IntValue{InternalValue: len(vectorItems)},
		IsEmpty:       &value.BoolValue{InternalValue: len(vectorItems) == 0},
	}

	AddVectorBuiltins(vector)

	return vector
}

var DefaultEmptyVectorValue = NewVectorValue([]value.IVOR{}, "["+value.IVOR_ANY+"]", value.IVOR_ANY)

type VectorItemReference struct {
	Vector *VectorValue
	Index  int
	Value  value.IVOR
}
