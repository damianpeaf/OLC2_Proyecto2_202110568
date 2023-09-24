package value

type FloatValue struct {
	InternalValue float64
}

func (d FloatValue) Value() interface{} {
	return d.InternalValue
}

func (d FloatValue) Type() string {
	return IVOR_FLOAT
}

func (d FloatValue) Copy() IVOR {
	return &FloatValue{d.InternalValue}
}
