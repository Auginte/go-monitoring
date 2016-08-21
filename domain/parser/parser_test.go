package parser

import (
	"testing"
)

func TestGetUniqueValues(t *testing.T) {
	input := map[string]string{
		"A": "A",
		"B": "Same",
		"C": "A",
		"D": "Same",
		"E": "Different",
	}
	expected := []string{
		"A",
		"Same",
		"Different",
	}

	assert(len(GetUniqueValues(input)) == len(expected), t)
}

func assert(expression bool, t *testing.T) {
	if !expression {
		t.Failed()
	}
}
