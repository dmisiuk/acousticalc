package e2e

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestE2EIntegration(t *testing.T) {
	t.Run("TestHelpCommand", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "--help")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			t.Fatalf("Failed to execute --help command: %v", err)
		}
		if !strings.Contains(out.String(), "Usage:") {
			t.Errorf("Expected 'Usage:' in help output, got: %s", out.String())
		}
	})

	t.Run("TestBasicCalculation", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "-e", "2+2")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			t.Fatalf("Failed to execute basic calculation: %v", err)
		}
		expected := "Result: 4"
		if !strings.Contains(out.String(), expected) {
			t.Errorf("Expected '%s' in output, got: %s", expected, out.String())
		}
	})

	t.Run("TestInvalidExpression", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "-e", "2+")
		var out bytes.Buffer
		cmd.Stderr = &out
		err := cmd.Run()
		if err == nil {
			t.Fatal("Expected error for invalid expression, but got none")
		}
		expected := "Error: Invalid expression"
		if !strings.Contains(out.String(), expected) {
			t.Errorf("Expected '%s' in error output, got: %s", expected, out.String())
		}
	})
}