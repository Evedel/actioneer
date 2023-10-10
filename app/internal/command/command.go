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
