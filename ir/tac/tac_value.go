package tac

import (
	"strconv"
)

type CastType string

const (
	CAST_INT   CastType = "int"
	CAST_FLOAT CastType = "float"
)

type ValueType string

const (
	TempType     ValueType = "Temp"
	HeapPtrType  ValueType = "HeapPtr"
	StackPtrType ValueType = "StackPtr"
	LiteralType  ValueType = "Literal"
	HeapIndex    ValueType = "HeapIndex"
	StackIndex   ValueType = "StackIndex"
)

// * Simple values
type SimpleValue interface {
	TACStmtI
	Cast() CastType
	Type() ValueType
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

func (t *Temp) Type() ValueType {
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
type HeapPtr struct {
	CastType CastType
}

func (h *HeapPtr) String() string {
	if h.CastType != "" {
		return "(" + string(h.CastType) + ") " + "H"
	}
	return "H"
}

func (h *HeapPtr) Type() ValueType {
	return HeapPtrType
}

func (h *HeapPtr) Cast() CastType {
	return h.CastType
}

// builder utils
func (h *HeapPtr) SetCast(castType CastType) *HeapPtr {
	h.CastType = castType
	return h
}

// ** StackPtr
type StackPtr struct {
	CastType CastType
}

func (s *StackPtr) String() string {
	if s.CastType != "" {
		return "(" + string(s.CastType) + ") " + "P"
	}
	return "P"
}

func (s *StackPtr) Type() ValueType {
	return StackPtrType
}

func (s *StackPtr) Cast() CastType {
	return s.CastType
}

// builder utils
func (s *StackPtr) SetCast(castType CastType) *StackPtr {
	s.CastType = castType
	return s
}

// ** Literal
type Literal struct {
	Value string
}

func (l *Literal) String() string {
	return l.Value
}

func (l *Literal) Type() ValueType {
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

func (t *Literal) SetCast(castType CastType) *Literal {
	// ?
	return t
}

// * Indexed values

type IndexedValue interface {
	SimpleValue
	Index() SimpleValue
	Cast() CastType
}

// ** HeapPtrIndexed

type HeapIndexedValue struct {
	IndexValue SimpleValue
	CastType   CastType
}

func (h *HeapIndexedValue) Cast() CastType {
	return h.CastType
}

func (h *HeapIndexedValue) Type() ValueType {
	return HeapIndex
}

func (h *HeapIndexedValue) String() string {
	if h.CastType != "" {
		return "(" + string(h.CastType) + ") " + h.getName()
	}
	return h.getName()
}

func (h *HeapIndexedValue) getName() string {
	return "heap[ (int) " + h.IndexValue.String() + "]"
}

func (h *HeapIndexedValue) Index() SimpleValue {
	return h.IndexValue
}

// builder utils
func (h *HeapIndexedValue) SetIndex(indexValue SimpleValue) *HeapIndexedValue {
	h.IndexValue = indexValue
	return h
}

func (t *HeapIndexedValue) SetCast(castType CastType) *HeapIndexedValue {
	// ?
	return t
}

// ** StackPtrIndexed

type StackIndexedValue struct {
	IndexValue SimpleValue
	CastType   CastType
}

func (s *StackIndexedValue) String() string {
	if s.CastType != "" {
		return "(" + string(s.CastType) + ") " + s.getName()
	}
	return s.getName()
}

func (s *StackIndexedValue) getName() string {
	return "stack[ (int) " + s.IndexValue.String() + "]"
}
func (s *StackIndexedValue) Index() SimpleValue {
	return s.IndexValue
}

func (s *StackIndexedValue) Cast() CastType {
	return ""
}

func (s *StackIndexedValue) Type() ValueType {
	return StackIndex
}

// builder utils
func (s *StackIndexedValue) SetIndex(indexValue SimpleValue) *StackIndexedValue {
	s.IndexValue = indexValue
	return s
}

func (t *StackIndexedValue) SetCast(castType CastType) *StackIndexedValue {
	// ?
	return t
}

// --- utils ---

func AddCastToSimpleValue(sv SimpleValue, cast CastType) SimpleValue {

	switch val := sv.(type) {
	case *Temp:
		return val.SetCast(cast)
	case *HeapPtr:
		return val.SetCast(cast)
	case *StackPtr:
		return val.SetCast(cast)
	case *Literal:
		return val.SetCast(cast)
	case *HeapIndexedValue:
		return val.SetCast(cast)
	case *StackIndexedValue:
		return val.SetCast(cast)
	default:
		panic("Unknown simple value type")
	}
}
