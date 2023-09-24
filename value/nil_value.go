package value

type NilValue struct {
}

func (s NilValue) Value() interface{} {
	return nil
}

func (s NilValue) Type() string {
	return IVOR_NIL
}

func (s NilValue) Copy() IVOR {
	return DefaultNilValue
}

var DefaultNilValue = &NilValue{}

type UnInitializedValue struct {
}

func (s UnInitializedValue) Value() interface{} {
	return nil
}

func (s UnInitializedValue) Type() string {
	return IVOR_UNINITIALIZED
}

func (s UnInitializedValue) Copy() IVOR {
	return DefaultUnInitializedValue
}

var DefaultUnInitializedValue = &UnInitializedValue{}
