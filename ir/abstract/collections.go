package abstract

var _countProp = &IVOR{
	Name:          "count",
	Type:          IVOR_INT,
	Address:       0,
	FrameRelative: false,
	Offset:        0,
}

var DefaultVectorScope = &BaseScope{
	Name:     "vector",
	Parent:   nil,
	Children: []*BaseScope{},
	Variables: map[string]*IVOR{
		"count": _countProp,
	},
}

var VectorAppendFunc = &Function{
	Name: "append",
	Type: BUILTIN_FUNCTION,
}

var VectorRemoveLastFunc = &Function{
	Name: "removeLast",
	Type: BUILTIN_FUNCTION,
}

var VectorRemoveFunc = &Function{
	Name: "remove",
	Type: BUILTIN_FUNCTION,
}
