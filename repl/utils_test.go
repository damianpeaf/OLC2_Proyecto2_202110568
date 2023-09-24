package repl

import (
	"testing"
)

func TestIsVector(t *testing.T) {

	matrix := map[string]bool{
		"[int]":   true,
		"int":     false,
		"[[int]]": false,
		"[]":      true,
	}

	for k, v := range matrix {
		if IsVectorType(k) != v {
			t.Errorf("isVector(%s) != %t", k, v)
		}
	}

}

func TestIsMatrix(t *testing.T) {
	matrix := map[string]bool{
		"[int]":     false,
		"int":       false,
		"[[int]]":   true,
		"[]":        false,
		"[[[int]]]": true,
		"[":         false,
		"[[]":       false,
	}

	for k, v := range matrix {
		if IsMatrixType(k) != v {
			t.Errorf("isMatrix(%s) != %t", k, v)
		}
	}
}
