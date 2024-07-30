package main

import (
	th "actioneer/internal/testing_helper"

	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"testing"
)

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

func BuildApp(t *testing.T) string {
	binName := "actioneer-" + randomString(10)
	cmd := []string{"go", "build", "-o", binName, "main.go"}

	errBuild := sh(cmd, false)
	th.AssertNil(t, errBuild)

	return binName
}

func CleanUp(t *testing.T, binName string) {
	cmd := []string{"rm", binName}
	errCleanup := sh(cmd, false)
	th.AssertNil(t, errCleanup)
}

func Test_Fails_WithoutConfig(t *testing.T) {
	binName := BuildApp(t)

	cmd := []string{"./" + binName}
	err := sh(cmd, false)
	th.AssertNotNil(t, err)
	th.AssertStringContains(t, "config.yaml", err.Error())
	th.AssertStringContains(t, "no such file or directory", err.Error())

	CleanUp(t, binName)
}
