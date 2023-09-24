package value

type StringValue struct {
	InternalValue string
}

func (s StringValue) Value() interface{} {
	return s.InternalValue
}

func (s StringValue) Type() string {
	return IVOR_STRING
}

func (s StringValue) Copy() IVOR {
	return &StringValue{InternalValue: s.InternalValue}
}
