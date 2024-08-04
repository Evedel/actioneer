package config

import (
	"actioneer/internal/logging"
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

type FakeConfigReader struct {
	openError      error
	readAllError   error
	unmarshalError error
}
type FakeReadCloser struct{}

func (frc FakeReadCloser) Read(p []byte) (n int, err error) {
	return 0, nil
}
func (frc FakeReadCloser) Close() error {
	return nil
}

func (fcr FakeConfigReader) Open(path string) (io.ReadCloser, error) {
	return FakeReadCloser{}, fcr.openError
}
func (fcr FakeConfigReader) ReadAll(r io.Reader) ([]byte, error) {
	return []byte{}, fcr.readAllError
}
func (fcr FakeConfigReader) Unmarshal(bytes []byte, v interface{}) error {
	return fcr.unmarshalError
}

func TestRead_Fake_NoErrors(t *testing.T) {
	fcr := FakeConfigReader{
		openError:      nil,
		readAllError:   nil,
		unmarshalError: nil,
	}

	_, err := Read(fcr, "some/path")
	if err != nil {
		t.Error("expected no error, got: " + err.Error())
	}
}

func TestRead_Fake_OpenErorr(t *testing.T) {
	fcr := FakeConfigReader{
		openError:      errors.New("fake open error"),
		readAllError:   nil,
		unmarshalError: nil,
	}

	var buf bytes.Buffer
	logging.Init("error", &buf)

	_, err := Read(fcr, "some/path")
	if err != fcr.openError {
		t.Error("expected " + fcr.openError.Error() + " error, got nil")
	}
	if !strings.Contains(buf.String(), "fake open error") {
		t.Error("expected " + fcr.openError.Error() + " error, got " + buf.String())
	}
}

func TestRead_Fake_ReadAllErorr(t *testing.T) {
	fcr := FakeConfigReader{
		openError:      nil,
		readAllError:   errors.New("fake read all error"),
		unmarshalError: nil,
	}

	var buf bytes.Buffer
	logging.Init("error", &buf)

	_, err := Read(fcr, "some/path")
	if err != fcr.readAllError {
		t.Error("expected " + fcr.readAllError.Error() + " error, got nil")
	}
	if !strings.Contains(buf.String(), "fake read all error") {
		t.Error("expected " + fcr.readAllError.Error() + " error, got " + buf.String())
	}
}

func TestRead_Fake_UnmarshalErorr(t *testing.T) {
	fcr := FakeConfigReader{
		openError:      nil,
		readAllError:   nil,
		unmarshalError: errors.New("fake unmarshal error"),
	}

	var buf bytes.Buffer
	logging.Init("error", &buf)

	_, err := Read(fcr, "some/path")
	if err != fcr.unmarshalError {
		t.Error("expected " + fcr.unmarshalError.Error() + " error, got nil")
	}
	if !strings.Contains(buf.String(), "fake unmarshal error") {
		t.Error("expected " + fcr.unmarshalError.Error() + " error, got " + buf.String())
	}
}

func TestIsValid_Ok(t *testing.T) {
	cfg := Config{
		Version: "1",
		Actions: []Action{
			{
				Name:      "Test Action",
				Alertname: "Test Alert",
				Command:   "echo \"test\"",
			},
		},
	}

	var buf bytes.Buffer
	logging.Init("error", &buf)

	if !IsValid(cfg) {
		t.Error("expected valid config, got invalid")
	}
	if buf.String() != "" {
		t.Error("expected empty log, got: " + buf.String())
	}
}

func TestIsValid_WrongVersion(t *testing.T) {
	cfg := Config{
		Version: "v2",
		Actions: []Action{
			{
				Alertname: "Test Alert",
				Command:   "echo \"test\"",
			},
		},
	}

	var buf bytes.Buffer
	logging.Init("error", &buf)

	if IsValid(cfg) {
		t.Error("expected invalid config, got valid")
	}
	if !strings.Contains(buf.String(), "wrong config version: v2") {
		t.Error("expected wrong config version error, got: " + buf.String())
	}
}

func TestIsValid_NoActions(t *testing.T) {
	cfg := Config{
		Version: "1",
		Actions: []Action{},
	}

	var buf bytes.Buffer
	logging.Init("error", &buf)

	if IsValid(cfg) {
		t.Error("expected invalid config, got valid")
	}
	if !strings.Contains(buf.String(), "no actions defined") {
		t.Error("expected no actions defined error, got: " + buf.String())
	}
}

func TestIsValid_EmptyAlertname(t *testing.T) {
	cfg := Config{
		Version: "1",
		Actions: []Action{
			{
				Name:    "Test Action",
				Command: "echo \"test\"",
			},
		},
	}

	var buf bytes.Buffer
	logging.Init("error", &buf)

	if IsValid(cfg) {
		t.Error("expected invalid config, got valid")
	}
	if !strings.Contains(buf.String(), "empty alertname in action") {
		t.Error("expected empty alertname error, got: " + buf.String())
	}
}

func TestIsValid_EmptyCommand(t *testing.T) {
	cfg := Config{
		Version: "1",
		Actions: []Action{
			{
				Name:      "Test Action",
				Alertname: "Test Alert",
			},
		},
	}

	var buf bytes.Buffer
	logging.Init("error", &buf)

	if IsValid(cfg) {
		t.Error("expected invalid config, got valid")
	}
	if !strings.Contains(buf.String(), "empty command in action") {
		t.Error("expected empty command error, got: " + buf.String())
	}
}

func TestIsValid_EmptyName(t *testing.T) {
	cfg := Config{
		Version: "1",
		Actions: []Action{
			{
				Alertname: "Test Alert",
				Command:   "echo \"test\"",
			},
		},
	}

	var buf bytes.Buffer
	logging.Init("error", &buf)

	if IsValid(cfg) {
		t.Error("expected invalid config, got valid")
	}
	if !strings.Contains(buf.String(), "empty name in action") {
		t.Error("expected empty command error, got: " + buf.String())
	}
}

func TestIsValid_DuplicateAlertname(t *testing.T) {
	cfg := Config{
		Version: "1",
		Actions: []Action{
			{
				Name:      "Test Action 1",
				Alertname: "Test Alert",
				Command:   "echo \"test\"",
			},
			{
				Name:      "Test Action 2",
				Alertname: "Test Alert",
				Command:   "echo \"test\"",
			},
		},
	}

	var buf bytes.Buffer
	logging.Init("error", &buf)

	if IsValid(cfg) {
		t.Error("expected invalid config, got valid")
	}
	if !strings.Contains(buf.String(), "multiple actions are not allowed for the same alertname") {
		t.Error("expected duplicate alertname error, got: " + buf.String())
	}
}
