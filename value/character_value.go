package value

type CharacterValue struct {
	InternalValue string
}

func (s CharacterValue) Value() interface{} {
	return s.InternalValue
}

func (s CharacterValue) Type() string {
	return IVOR_CHARACTER
}

func (s CharacterValue) Copy() IVOR {
	return &CharacterValue{s.InternalValue}
}
