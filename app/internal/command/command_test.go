package command

import (
	"fmt"
	"os/exec"
	"testing"
)

type FakeCommandRunner struct {
	cmd    *exec.Cmd
	stdout string
	stderr string
}

func (c FakeCommandRunner) Command(name string, arg ...string) *exec.Cmd {
	return CommandRunner{}.Command(name, arg...)
}
func (c *FakeCommandRunner) Run(cmd *exec.Cmd) error {
	c.cmd = cmd
	cmd.Stderr.Write([]byte(c.stderr))
	cmd.Stdout.Write([]byte(c.stdout))
	return nil
}

func TestExecute_ExpectedSetCommandAndArgs(t *testing.T) {
	fcr := FakeCommandRunner{}
	Execute(&fcr, "fake command", false)

	fmt.Println(fcr.cmd)

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
