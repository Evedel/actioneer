package testing_helper

import (
	"actioneer/internal/state"
	"strings"
)

type Dict map[string]string

type Itesting interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

func GetState(options map[string]string) state.State {
	alertNameKey := "alertname"
	if options["alertname"] != "" {
		alertNameKey = options["alertname"]
	}
	return state.State{
		AlertNameKey: alertNameKey,
	}
}

func AssertNil(t Itesting, err error) {
	if err != nil {
		t.Error("expected no error, got: " + err.Error())
	}
}

func AssertNotNil(t Itesting, err error) {
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func AssertEqual(t Itesting, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func AssertStringContains(t Itesting, expected, actual string) {
	if !strings.Contains(actual, expected) {
		t.Errorf("expected string to contain '%s', got '%s'", expected, actual)
	}
}
