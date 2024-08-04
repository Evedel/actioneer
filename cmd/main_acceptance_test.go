package main

import (
	th "actioneer/internal/testing_helper"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"testing"
)

const verbose = false

func Test_Fails_NoConfig(t *testing.T) {
	binName := BuildApp(t)

	cmd := []string{"./" + binName}
	cmdChanel := make(chan int)
	cmdHandle, errCMDCreate := sh(cmd, verbose, cmdChanel)
	th.AssertNil(t, errCMDCreate)

	combinedStd, exited := wait_sh(cmdHandle, cmdChanel, verbose)
	th.AssertEqual(t, true, exited)
	th.AssertStringContains(t, "config.yaml", combinedStd)
	th.AssertStringContains(t, "no such file or directory", combinedStd)

	CleanUp(t, binName, cmdHandle)
}

func Test_Fails_NoVersionInConfig(t *testing.T) {
	binName := BuildApp(t)

	MakeTestConfig(t, "")
	cmd := []string{"./" + binName, "-config-path", "config.yaml"}
	cmdChanel := make(chan int)
	cmdHandle, errCMDCreate := sh(cmd, verbose, cmdChanel)
	th.AssertNil(t, errCMDCreate)

	combinedStd, exited := wait_sh(cmdHandle, cmdChanel, verbose)
	th.AssertEqual(t, true, exited)
	th.AssertStringContains(t, "\"wrong config version: \"", combinedStd)

	CleanUp(t, binName, cmdHandle)
}

func Test_Fails_InvalidVersionInConfig(t *testing.T) {
	binName := BuildApp(t)

	config := `version: 9999999
`
	MakeTestConfig(t, config)
	cmd := []string{"./" + binName, "-config-path", "config.yaml"}
	cmdChanel := make(chan int)
	cmdHandle, errCMDCreate := sh(cmd, verbose, cmdChanel)
	th.AssertNil(t, errCMDCreate)

	combinedStd, exited := wait_sh(cmdHandle, cmdChanel, verbose)
	th.AssertEqual(t, true, exited)
	th.AssertStringContains(t, "\"wrong config version: 9999999\"", combinedStd)

	CleanUp(t, binName, cmdHandle)
}

func Test_Fails_NoActionsDefined(t *testing.T) {
	binName := BuildApp(t)

	config := `version: 1
`
	MakeTestConfig(t, config)
	cmd := []string{"./" + binName, "-config-path", "config.yaml"}
	cmdChanel := make(chan int)
	cmdHandle, errCMDCreate := sh(cmd, verbose, cmdChanel)
	th.AssertNil(t, errCMDCreate)

	combinedStd, exited := wait_sh(cmdHandle, cmdChanel, verbose)
	th.AssertEqual(t, true, exited)
	th.AssertStringContains(t, "\"no actions defined\"", combinedStd)

	CleanUp(t, binName, cmdHandle)
}

func Test_Fails_IncorrectActionInConfig(t *testing.T) {
	binName := BuildApp(t)

	config := `version: 1
actions:
  - test
`
	MakeTestConfig(t, config)
	cmd := []string{"./" + binName, "-config-path", "config.yaml"}
	cmdChanel := make(chan int)
	cmdHandle, errCMDCreate := sh(cmd, verbose, cmdChanel)
	th.AssertNil(t, errCMDCreate)

	combinedStd, exited := wait_sh(cmdHandle, cmdChanel, verbose)
	th.AssertEqual(t, true, exited)
	th.AssertStringContains(t, "cannot unmarshal", combinedStd)

	CleanUp(t, binName, cmdHandle)
}

func Test_Fails_NoAlertnameInAction(t *testing.T) {
	binName := BuildApp(t)

	config := `version: 1
actions:
  - name: test
`
	MakeTestConfig(t, config)
	cmd := []string{"./" + binName, "-config-path", "config.yaml"}
	cmdChanel := make(chan int)
	cmdHandle, errCMDCreate := sh(cmd, verbose, cmdChanel)
	th.AssertNil(t, errCMDCreate)

	combinedStd, exited := wait_sh(cmdHandle, cmdChanel, verbose)
	th.AssertEqual(t, true, exited)
	th.AssertStringContains(t, "empty alertname in action", combinedStd)

	CleanUp(t, binName, cmdHandle)
}

func Test_Fails_NoNameInAction(t *testing.T) {
	binName := BuildApp(t)

	config := `version: 1
actions:
  - alertname: test
`
	MakeTestConfig(t, config)
	cmd := []string{"./" + binName, "-config-path", "config.yaml"}
	cmdChanel := make(chan int)
	cmdHandle, errCMDCreate := sh(cmd, verbose, cmdChanel)
	th.AssertNil(t, errCMDCreate)

	combinedStd, exited := wait_sh(cmdHandle, cmdChanel, verbose)
	th.AssertEqual(t, true, exited)
	th.AssertStringContains(t, "empty name in action", combinedStd)

	CleanUp(t, binName, cmdHandle)
}

func Test_Fails_NoCommandInAction(t *testing.T) {
	binName := BuildApp(t)

	config := `version: 1
actions:
  - name: test
    alertname: test
`
	MakeTestConfig(t, config)
	cmd := []string{"./" + binName, "-config-path", "config.yaml"}
	cmdChanel := make(chan int)
	cmdHandle, errCMDCreate := sh(cmd, verbose, cmdChanel)
	th.AssertNil(t, errCMDCreate)

	combinedStd, exited := wait_sh(cmdHandle, cmdChanel, verbose)
	th.AssertEqual(t, true, exited)
	th.AssertStringContains(t, "empty command in action", combinedStd)

	CleanUp(t, binName, cmdHandle)
}

func Test_CanBeStarted(t *testing.T) {
	binName := BuildApp(t)

	config := `version: 1
actions:
  - name: test
    alertname: test
    command: echo "test"
`
	MakeTestConfig(t, config)
	cmd := []string{"./" + binName, "-config-path", "config.yaml"}
	cmdChanel := make(chan int)
	cmdHandle, errCMDCreate := sh(cmd, verbose, cmdChanel)
	th.AssertNil(t, errCMDCreate)

	combinedStd, exited := check_sh(cmdHandle, cmdChanel, verbose)
	th.AssertEqual(t, false, exited)
	th.AssertEqual(t, "", combinedStd)

	CleanUp(t, binName, cmdHandle)
}

func Test_HealthzReturns(t *testing.T) {
	binName := BuildApp(t)

	config := `version: 1
actions:
  - name: test
    alertname: test
    command: echo "test"
`
	MakeTestConfig(t, config)
	cmd := []string{"./" + binName, "-config-path", "config.yaml"}
	cmdChanel := make(chan int)
	cmdHandle, errCMDCreate := sh(cmd, verbose, cmdChanel)
	th.AssertNil(t, errCMDCreate)

	combinedStd, exited := check_sh(cmdHandle, cmdChanel, verbose)
	th.AssertEqual(t, false, exited)
	th.AssertEqual(t, "", combinedStd)

	cmdHealthz := []string{"curl", "-s", "http://localhost:8080/healthz"}
	cmdHealthzChanel := make(chan int)
	cmdHealthzHandle, errCMDHealthzCreate := sh(cmdHealthz, verbose, cmdHealthzChanel)
	th.AssertNil(t, errCMDHealthzCreate)

	combinedStdHealthz, exitedHealthz := wait_sh(cmdHealthzHandle, cmdHealthzChanel, verbose)
	th.AssertEqual(t, true, exitedHealthz)
	th.AssertEqual(t, "0: : ok", combinedStdHealthz)

	CleanUp(t, binName, cmdHandle)
}

func BuildApp(t *testing.T) string {
	binName := "actioneer-" + randomString(10)
	cmd := []string{"go", "build", "-o", binName, "main.go"}

	combinedStd := sh_run(cmd, false)
	th.AssertNil(t, combinedStd)

	return binName
}

func MakeTestConfig(t *testing.T, config string) {
	err := os.WriteFile("config.yaml", []byte(config), 0644)
	th.AssertNil(t, err)
}

func CleanUp(t *testing.T, binName string, cmdHandle *exec.Cmd) {
	cmdHandle.Process.Kill()

	cmdRmBin := []string{"rm", binName}
	sh_run(cmdRmBin, false)

	cmdRmCfg := []string{"rm", "config.yaml"}
	sh_run(cmdRmCfg, false)
}

func randomString(n int) string {
	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func sh_run(args []string, verbose bool) error {
	cmd := exec.Command(args[0], args[1:]...)

	var stdout strings.Builder
	var stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if verbose {
		fmt.Println("running command:", cmd.String())
	}

	combinedStd := error(nil)
	if errCmd := cmd.Run(); errCmd != nil {
		combinedStd = fmt.Errorf("%s: %s: %s", errCmd, cmd.Stderr, cmd.Stdout)
	}

	if verbose {
		fmt.Println("stdout:", cmd.Stdout)
		fmt.Println("stderr:", cmd.Stderr)
	}
	return combinedStd
}

func sh(args []string, verbose bool, output chan int) (*exec.Cmd, error) {
	cmd := exec.Command(args[0], args[1:]...)

	var stdout strings.Builder
	var stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if verbose {
		fmt.Println("running command:", cmd.String())
	}

	errStart := cmd.Start()
	if errStart != nil {
		return nil, errStart
	}

	go func() {
		cmd.Wait()
		output <- 1
	}()

	return cmd, errStart
}

func check_sh(cmd *exec.Cmd, output chan int, verbose bool) (string, bool) {
	exited := false
	conbinedStd := ""

	select {
	default:
		_ = 1
	case <-output:
		conbinedStd = fmt.Sprintf("%d: %s: %s", cmd.ProcessState.ExitCode(), cmd.Stderr, cmd.Stdout)
		exited = true
		if verbose {
			fmt.Println("exit code:", cmd.ProcessState.ExitCode())
			fmt.Println("stdout:", cmd.Stdout)
			fmt.Println("stderr:", cmd.Stderr)
		}
	}

	return conbinedStd, exited
}

func wait_sh(cmd *exec.Cmd, output chan int, verbose bool) (string, bool) {
	conbinedStd := ""
	exited := false

	maxWait := 10
	for i := 0; i < maxWait; i++ {
		time.Sleep(1 * time.Second)
		conbinedStd, exited = check_sh(cmd, output, verbose)
		if exited {
			break
		}
	}

	return conbinedStd, exited
}
