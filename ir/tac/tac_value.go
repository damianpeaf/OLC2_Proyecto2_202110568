package tac

import (
	"strconv"
)

type CastType string

const (
	CAST_INT   CastType = "Int"
	CAST_FLOAT CastType = "Float"
)

type SimpleValueType string

const (
	TempType     SimpleValueType = "Temp"
	HeapPtrType  SimpleValueType = "HeapPtr"
	StackPtrType SimpleValueType = "StackPtr"
	LiteralType  SimpleValueType = "Literal"
)

// * Simple values
type SimpleValue interface {
	TACStmtI
	Cast() CastType
	Type() SimpleValueType
}

// ** Temp
type Temp struct {
	ID       int
	CastType CastType
}

func (t *Temp) String() string {
	if t.CastType != "" {
		return "(" + string(t.CastType) + ") " + t.TempName()
	}
	return t.TempName()
}

func (t *Temp) TempName() string {
	return "t" + strconv.Itoa(t.ID)
}

func (t *Temp) Type() SimpleValueType {
	return TempType
}

func (t *Temp) Cast() CastType {
	return t.CastType
}

// builder utils
func (t *Temp) SetCast(castType CastType) *Temp {
	t.CastType = castType
	return t
}

// ** HeapPtr
type HeapPtr struct{}

func (h *HeapPtr) String() string {
	return "H"
}

func (h *HeapPtr) Type() SimpleValueType {
	return HeapPtrType
}

func (h *HeapPtr) Cast() CastType {
	return ""
}

// ** StackPtr
type StackPtr struct{}

func (s *StackPtr) String() string {
	return "P"
}

func (s *StackPtr) Type() SimpleValueType {
	return StackPtrType
}

func (s *StackPtr) Cast() CastType {
	return ""
}

// ** Literal
type Literal struct {
	Value string
}

func (l *Literal) String() string {
	return l.Value
}

func (l *Literal) Type() SimpleValueType {
	return LiteralType
}

func (l *Literal) Cast() CastType {
	return ""
}

// builder utils
func (l *Literal) SetValue(value string) *Literal {
	l.Value = value
	return l
}

// * Indexed values

type IndexedValue interface {
	Index() SimpleValue
}

// ** HeapPtrIndexed

type HeapIndexedValue struct {
	IndexValue SimpleValue
}

func (h *HeapIndexedValue) String() string {
	return "heap[" + h.IndexValue.String() + "]"
}

func (h *HeapIndexedValue) Index() SimpleValue {
	return h.IndexValue
}

// builder utils
func (h *HeapIndexedValue) SetIndex(indexValue SimpleValue) *HeapIndexedValue {
	h.IndexValue = indexValue
	return h
}

// ** StackPtrIndexed

type StackIndexedValue struct {
	IndexValue SimpleValue
}

func (s *StackIndexedValue) String() string {
	return "stack[" + s.IndexValue.String() + "]"
}

func (s *StackIndexedValue) Index() SimpleValue {
	return s.IndexValue
}

// builder utils
func (s *StackIndexedValue) SetIndex(indexValue SimpleValue) *StackIndexedValue {
	s.IndexValue = indexValue
	return s
}
