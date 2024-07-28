package testing_helper

import (
	"errors"
	"testing"
)

func Test_GetState(t *testing.T) {
	options := Dict{
		"alertname": "test-alert",
	}
	state := GetState(options)
	AssertEqual(t, "test-alert", state.AlertNameKey)
}

type MockTesting struct {
	errorMessage string
}

func (m *MockTesting) Error(args ...interface{}) {
	m.errorMessage = args[0].(string)
}

func (m *MockTesting) Errorf(format string, args ...interface{}) {
	m.errorMessage = format
}

func Test_AssertNil_Pass(t *testing.T) {
	mt := &MockTesting{}
	AssertNil(mt, nil)
	if mt.errorMessage != "" {
		t.Error("expected no error, got: " + mt.errorMessage)
	}
}

func Test_AssertNil_Fail(t *testing.T) {
	mt := &MockTesting{}
	AssertNil(mt, errors.New("test"))
	if mt.errorMessage == "" {
		t.Error("expected [test] error , got nil")
	}
}

func Test_AssertNotNil_Pass(t *testing.T) {
	mt := &MockTesting{}
	AssertNotNil(mt, errors.New("test"))
	if mt.errorMessage != "" {
		t.Error("expected no error, got " + mt.errorMessage)
	}
}

func Test_AssertNotNil_Fail(t *testing.T) {
	mt := &MockTesting{}
	AssertNotNil(mt, nil)
	if mt.errorMessage == "" {
		t.Error("expected error, got nil")
	}
}

func Test_AssertEqual_Pass(t *testing.T) {
	mt := &MockTesting{}
	AssertEqual(mt, 1, 1)
	if mt.errorMessage != "" {
		t.Error("expected no error, got " + mt.errorMessage)
	}
}

func Test_AssertEqual_Fail(t *testing.T) {
	mt := &MockTesting{}
	AssertEqual(mt, 1, 2)
	if mt.errorMessage == "" {
		t.Error("expected 1 not equal 2 error, got nil")
	}
}
func Test_AssertStringContains_Pass(t *testing.T) {
	mt := &MockTesting{}
	AssertStringContains(mt, "test", "test string that shold be found")
	if mt.errorMessage != "" {
		t.Error("expected [test] to be found, got " + mt.errorMessage)
	}
}

func Test_AssertStringContains_Fail(t *testing.T) {
	mt := &MockTesting{}
	AssertStringContains(mt, "test", "string that should not be found")
	if mt.errorMessage == "" {
		t.Error("expected [test] to be not found, got nil (i.e. found)")
	}
}
