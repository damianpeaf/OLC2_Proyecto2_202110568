package value

func IsPrimitiveType(t string) bool {
	switch t {
	case IVOR_BOOL, IVOR_INT, IVOR_FLOAT, IVOR_STRING, IVOR_NIL, IVOR_CHARACTER:
		return true
	default:
		return false
	}
}

func ImplicitCast(targetType string, value IVOR) (IVOR, bool) {

	if targetType == value.Type() {
		return value, true
	}

	// implicit conversion

	// 1. int can be converted to float
	if targetType == IVOR_FLOAT && value.Type() == IVOR_INT {
		return &FloatValue{
			InternalValue: float64(value.(*IntValue).InternalValue),
		}, true
	}

	// 2. Character can be converted to string
	if targetType == IVOR_STRING && value.Type() == IVOR_CHARACTER {
		return &StringValue{
			InternalValue: string(value.(*CharacterValue).InternalValue),
		}, true
	}

	return nil, false

}
