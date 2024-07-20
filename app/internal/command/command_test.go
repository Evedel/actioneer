package command

import (
	"actioneer/internal/logging"
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestExecute_SetCommandAndArgs(t *testing.T) {
	fcr := FakeCommandRunner{}
	Execute(&fcr, "fake command", false)

	if fcr.cmd == nil {
		t.Errorf("expected cmd to be non-nil")
	}
	if fcr.cmd.Args[0] != "bash" {
		t.Errorf("expected cmd.Args[0] to be bash, got %+v", fcr.cmd.Args[0])
	}
	if fcr.cmd.Args[1] != "-c" {
		t.Errorf("expected cmd.Args[1] to be -c, got %+v", fcr.cmd.Args[1])
	}
	if fcr.cmd.Args[2] != "fake command" {
		t.Errorf("expected cmd.Args[2] to be fake command, got %+v", fcr.cmd.Args[2])
	}
}

func TestExecute_ExpectedSetStdout(t *testing.T) {
	fcr := FakeCommandRunner{}
	fcr.stdout = "fake stdout"

	var buf bytes.Buffer

	ed := logging.Init("debug", &buf)
	if ed != nil {
		t.Error("unexpected error during logging.Init")
	}

	Execute(&fcr, "fake command", false)

	expected_log := "stdout: fake stdout"
	if !strings.Contains(buf.String(), expected_log) {
		t.Error("\nexpected log not found \n\n all logs: " + buf.String() + "\n expected log: " + expected_log + "\n")
	}

	unexpected_log := "stderr:"
	if strings.Contains(buf.String(), unexpected_log) {
		t.Error("\nunexpected_log log found \n\n all logs: " + buf.String() + "\n logs that should not be there: " + unexpected_log + "\n")
	}
}

func TestExecute_ExpectedSetStderr(t *testing.T) {
	fcr := FakeCommandRunner{}
	fcr.stderr = "fake stderr"

	var buf bytes.Buffer

	ed := logging.Init("debug", &buf)
	if ed != nil {
		t.Error("unexpected error during logging.Init")
	}

	Execute(&fcr, "fake command", false)

	expected_stdout := `"ERROR","msg":"stderr: fake stderr"}`
	if !strings.Contains(buf.String(), expected_stdout) {
		t.Error("\nexpected log not found \n\n all logs: " + buf.String() + "\n expected log: " + expected_stdout + "\n")
	}

	expected_stderr := `"INFO","msg":"stdout: "}`
	if !strings.Contains(buf.String(), expected_stderr) {
		t.Error("\nexpected log not found \n\n all logs: " + buf.String() + "\n expected log: " + expected_stderr + "\n")
	}
}

func TestExecute_ExpectedDryRun(t *testing.T) {
	fcr := FakeCommandRunner{}

	var buf bytes.Buffer

	ed := logging.Init("debug", &buf)
	if ed != nil {
		t.Error("unexpected error during logging.Init")
	}

	Execute(&fcr, "fake command", true)

	expected_log := `"INFO","msg":"dry run: fake command"}`
	if !strings.Contains(buf.String(), expected_log) {
		t.Error("\nexpected log not found \n\n all logs: " + buf.String() + "\n expected log: " + expected_log + "\n")
	}

	if fcr.cmd != nil {
		t.Error("expected cmd to be nil")
	}
}

func TestExecute_ExpectedCommandRunFailure(t *testing.T) {
	fcr := FakeCommandRunner{}
	fcr.runErr = exec.ErrNotFound

	var buf bytes.Buffer

	ed := logging.Init("debug", &buf)
	if ed != nil {
		t.Error("unexpected error during logging.Init")
	}

	Execute(&fcr, "fake command", false)

	unexpected_stdout := `stdout: `
	if strings.Contains(buf.String(), unexpected_stdout) {
		t.Log(buf.String())
		t.Error("should not have found stdout log")
	}
	unexpected_stderr := `stderr: `
	if strings.Contains(buf.String(), unexpected_stderr) {
		t.Log(buf.String())
		t.Error("should not have found stderr log")
	}

	expected_err := `command execution failed: `
	if !strings.Contains(buf.String(), expected_err) {
		t.Error("\nexpected log not found \n\n all logs: " + buf.String() + "\n expected log: " + expected_err + "\n")
	}
}

func TestExecute_ExpectedRealCommandRunSuccess(t *testing.T) {
	cr := CommandRunner{}

	var buf bytes.Buffer
	ed := logging.Init("debug", &buf)
	if ed != nil {
		t.Error("unexpected error during logging.Init")
	}

	Execute(&cr, `echo "hello"`, false)

	expected_stdout := `"INFO","msg":"stdout: hello\n"}`
	if !strings.Contains(buf.String(), expected_stdout) {
		t.Log(buf.String())
		t.Error("should not have found stdout log")
	}
	unexpected_stderr := `stderr: `
	if strings.Contains(buf.String(), unexpected_stderr) {
		t.Log(buf.String())
		t.Error("should not have found stderr log")
	}

	unexpected_err := `command execution failed: `
	if strings.Contains(buf.String(), unexpected_err) {
		t.Error("\nunexpected log found \n\n all logs: " + buf.String() + "\n unexpected log: " + unexpected_err + "\n")
	}
}

func TestExecute_ExpectedRealCommandRunFailure(t *testing.T) {
	cr := CommandRunner{}

	var buf bytes.Buffer
	ed := logging.Init("debug", &buf)
	if ed != nil {
		t.Error("unexpected error during logging.Init")
	}

	Execute(&cr, `echo "hello" && exit 1`, false)

	expected_stderr1 := `"ERROR","msg":"command execution failed: echo \"hello\" && exit 1"}`
	if !strings.Contains(buf.String(), expected_stderr1) {
		t.Log(buf.String())
		t.Error("expected log not found")
	}

	expected_stderr2 := `"ERROR","msg":"exit status 1"}`
	if !strings.Contains(buf.String(), expected_stderr2) {
		t.Log(buf.String())
		t.Error("expected log not found")
	}
}
