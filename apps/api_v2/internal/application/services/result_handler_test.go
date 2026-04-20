package services

import (
	"testing"

	exam_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
)

func TestCompareOutputValues_ArrayIgnoresWhitespace(t *testing.T) {
	actual := "[1, 2, 3]"
	expected := "[1,2,3]"

	if !compareOutputValues(actual, expected, exam_constants.VariableFormatArray) {
		t.Fatalf("expected array outputs to be equivalent")
	}
}

func TestCompareOutputValues_BooleanNormalization(t *testing.T) {
	if !compareOutputValues("True", "true", exam_constants.VariableFormatBoolean) {
		t.Fatalf("expected boolean outputs to be equivalent")
	}
	if !compareOutputValues("1", "true", exam_constants.VariableFormatBoolean) {
		t.Fatalf("expected 1 and true to be equivalent")
	}
	if compareOutputValues("false", "true", exam_constants.VariableFormatBoolean) {
		t.Fatalf("expected false and true to be different")
	}
}
