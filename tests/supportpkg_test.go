package tests

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestKubectlCommands(t *testing.T) {
	// Execute `kubectl` command
	cmd := exec.Command("kubectl", "nginx-supportpkg")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	_ = cmd.Run()

	// Compare the output with expected result
	expectedOutput := "required flag(s)"
	actualOutput := stdout.String()
	if !strings.Contains(actualOutput, expectedOutput) {
		t.Errorf("Expected output should contain \"%s\" but got \"%s\"", expectedOutput, actualOutput)
	}
}
