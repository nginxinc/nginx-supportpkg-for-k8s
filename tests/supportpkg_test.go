package tests

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestKubectlNginxSupportpkgVersion(t *testing.T) {
	cmd := exec.Command("kubectl", "nginx-supportpkg", "-v")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		t.Errorf("Failed to run 'kubectl nginx-supportpkg -v': %v", err)
	}

	expectedOutput := "nginx-supportpkg - version"
	actualOutput := stdout.String()
	if !strings.Contains(actualOutput, expectedOutput) {
		t.Errorf("Expected output should contain '%s' but got '%s'", expectedOutput, actualOutput)
	}
}
