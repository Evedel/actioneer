package testing_helper

import (
	"actioneer/internal/state"
	"strings"
	"testing"
)

type Dict map[string]string

func GetState(options map[string]string) state.State {
	alertNameKey := "alertname"
	if options["alertname"] != "" {
		alertNameKey = options["alertname"]
	}
	return state.State{
		AlertNameKey: alertNameKey,
	}
}

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
