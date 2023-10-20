package utils

import (
	"regexp"
	"strings"
)

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
