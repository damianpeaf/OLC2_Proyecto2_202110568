package repl

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/value"
	"regexp"
	"strings"
)

func StringToVector(s *value.StringValue) *VectorValue {

	items := make([]value.IVOR, 0)

	for _, c := range s.InternalValue {
		items = append(items, &value.CharacterValue{InternalValue: string(c)})
	}

	return NewVectorValue(items, "["+value.IVOR_CHARACTER+"]", value.IVOR_CHARACTER)

}

func IsVectorType(_type string) bool {

	// Vector starts with only one [ and ends with only one ]
	vectorPattern := "^\\[.*\\]$"

	// Matrix starts with AT LEAST two [[ and ends with at least two ]]
	matrixPattern := "^\\[\\[.*\\]\\](\\[.*\\]\\])*$"

	// match vector pattern but not matrix pattern

	match, _ := regexp.MatchString(vectorPattern, _type)
	match2, _ := regexp.MatchString(matrixPattern, _type)

	return match && !match2
}

func RemoveBrackets(s string) string {
	return strings.Trim(s, "[]")
}

func IsMatrixType(_type string) bool {

	// Matrix starts with AT LEAST two [[ and ends with at least two ]]
	matrixPattern := "^\\[\\[.*\\]\\](\\[.*\\]\\])*$"

	match, _ := regexp.MatchString(matrixPattern, _type)

	return match
}
