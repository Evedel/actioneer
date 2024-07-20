package command

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

type ICommandRunner interface {
	Command(name string, arg ...string) *exec.Cmd
	Run(cmd *exec.Cmd) error
}

type CommandRunner struct {
}
func (c CommandRunner) Command(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}
func (c CommandRunner) Run(cmd *exec.Cmd) error {
	return cmd.Run()
}

type FakeCommandRunner struct {
	Calls  []string
	cmd    *exec.Cmd
	runErr error
	stdout string
	stderr string
}
func (c *FakeCommandRunner) Command(name string, arg ...string) *exec.Cmd {
	c.Calls = append(c.Calls, fmt.Sprintf("%s %s", name, strings.Join(arg, " ")))
	return exec.Command(name, arg...)
}
func (c *FakeCommandRunner) Run(cmd *exec.Cmd) error {
	c.cmd = cmd
	cmd.Stderr.Write([]byte(c.stderr))
	cmd.Stdout.Write([]byte(c.stdout))
	return c.runErr
}

func Execute(icr ICommandRunner, command string, isDryRun bool) {
	slog.Debug("processing command: " + fmt.Sprint(command))
	if isDryRun {
		slog.Info("dry run: " + command)
	} else {
		cmd := icr.Command("bash", "-c", command)
		var stdout strings.Builder
		var stderr strings.Builder
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := icr.Run(cmd); err != nil {
			slog.Error("command execution failed: " + command)
			slog.Error(err.Error())
			return
		}
		slog.Info("stdout: " + stdout.String())
		if stderr.String() != "" {
			slog.Error("stderr: " + stderr.String())
		}
	}
}
