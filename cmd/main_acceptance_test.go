package main

import (
	th "actioneer/internal/testing_helper"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"

	"testing"
)

func Test_Fails_NoConfig(t *testing.T) {
	binName := BuildApp(t)

	cmd := []string{"./" + binName}
	err := sh(cmd, false)
	th.AssertNotNil(t, err)
	th.AssertStringContains(t, "config.yaml", err.Error())
	th.AssertStringContains(t, "no such file or directory", err.Error())

	CleanUp(t, binName)
}

func Test_Fails_NoVersionInConfig(t *testing.T) {
	binName := BuildApp(t)

	MakeTestConfig(t, "")
	cmd := []string{"./" + binName, "-config-path", "config.yaml"}
	err := sh(cmd, false)
	th.AssertStringContains(t, "\"wrong config version: \"", err.Error())
	th.AssertNotNil(t, err)

	CleanUp(t, binName)
}

func Test_Fails_InvalidVersionInConfig(t *testing.T) {
	binName := BuildApp(t)

	config := `version: 9999999
`
	MakeTestConfig(t, config)
	cmd := []string{"./" + binName, "-config-path", "config.yaml"}
	err := sh(cmd, false)

	th.AssertStringContains(t, "\"wrong config version: 9999999\"", err.Error())
	th.AssertNotNil(t, err)

	CleanUp(t, binName)
}

func Test_Fails_NoActionsDefined(t *testing.T) {
	binName := BuildApp(t)

	config := `version: 1
`
	MakeTestConfig(t, config)
	cmd := []string{"./" + binName, "-config-path", "config.yaml"}
	err := sh(cmd, false)

	th.AssertStringContains(t, "\"no actions defined\"", err.Error())
	th.AssertNotNil(t, err)

	CleanUp(t, binName)
}

func Test_Fails_IncorrectActionInConfig(t *testing.T) {
	binName := BuildApp(t)

	config := `version: 1
actions:
  - test
`
	MakeTestConfig(t, config)
	cmd := []string{"./" + binName, "-config-path", "config.yaml"}
	err := sh(cmd, false)

	th.AssertStringContains(t, "cannot unmarshal", err.Error())
	th.AssertNotNil(t, err)

	CleanUp(t, binName)
}

func Test_Fails_NoAlertnameInAction(t *testing.T) {
	binName := BuildApp(t)

	config := `version: 1
actions:
  - name: test
`
	MakeTestConfig(t, config)
	cmd := []string{"./" + binName, "-config-path", "config.yaml"}
	err := sh(cmd, false)

	th.AssertStringContains(t, "empty alertname in action", err.Error())
	th.AssertNotNil(t, err)

	CleanUp(t, binName)
}

func Test_Fails_NoNameInAction(t *testing.T) {
	binName := BuildApp(t)

	config := `version: 1
actions:
  - alertname: test
`
	MakeTestConfig(t, config)
	cmd := []string{"./" + binName, "-config-path", "config.yaml"}
	err := sh(cmd, false)

	th.AssertStringContains(t, "empty name in action", err.Error())
	th.AssertNotNil(t, err)

	CleanUp(t, binName)
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
	err := sh(cmd, false)

	th.AssertStringContains(t, "empty command in action", err.Error())
	th.AssertNotNil(t, err)

	CleanUp(t, binName)
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
	// cmd := []string{"./" + binName, "-config-path", "config.yaml"}
	// err := sh(cmd, false)

	// th.AssertNil(t, err)
	th.AssertNil(t, errors.New("there is no way to test health endpoint currently"))

	CleanUp(t, binName)
}

func BuildApp(t *testing.T) string {
	binName := "actioneer-" + randomString(10)
	cmd := []string{"go", "build", "-o", binName, "main.go"}

	errBuild := sh(cmd, false)
	th.AssertNil(t, errBuild)

	return binName
}

func MakeTestConfig(t *testing.T, config string) {
	err := os.WriteFile("config.yaml", []byte(config), 0644)
	th.AssertNil(t, err)
}

func CleanUp(t *testing.T, binName string) {
	cmdRmBin := []string{"rm", binName}
	sh(cmdRmBin, false)

	cmdRmCfg := []string{"rm", "config.yaml"}
	sh(cmdRmCfg, false)
}

func randomString(n int) string {
	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func sh(args []string, verbose bool) error {
	cmd := exec.Command(args[0], args[1:]...)

	var stdout strings.Builder
	var stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if verbose {
		fmt.Println("running command:", cmd.String())
	}

	errReturn := error(nil)
	if errInternal := cmd.Run(); errInternal != nil {
		errReturn = fmt.Errorf("%s: %s:%s", errInternal, cmd.Stderr, cmd.Stdout)
	}

	if verbose {
		fmt.Println("stdout:", stdout.String())
		fmt.Println("stderr:", stderr.String())
	}
	return errReturn
}
