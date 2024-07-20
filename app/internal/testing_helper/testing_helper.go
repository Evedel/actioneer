package testing_helper

import (
	"strings"
	"testing"
)

func AssertNil(t *testing.T, err error) {
	if err != nil {
		t.Error("expected no error, got: " + err.Error())
	}
}

func AssertNotNil(t *testing.T, err error) {
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func AssertEqual(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func AssertStringContains(t *testing.T, expected, actual string) {
	if !strings.Contains(actual, expected) {
		t.Errorf("expected string to contain '%s', got '%s'", expected, actual)
	}
}
