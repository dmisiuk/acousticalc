package e2e

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestE2EWorkflow(t *testing.T) {
	t.Run("TestFullWorkflow", func(t *testing.T) {
		// Step 1: Run a simple calculation
		cmd1 := exec.Command("go", "run", "main.go", "-e", "10 / 2")
		var out1 bytes.Buffer
		cmd1.Stdout = &out1
		if err := cmd1.Run(); err != nil {
			t.Fatalf("Workflow step 1 failed: %v", err)
		}
		if !strings.Contains(out1.String(), "Result: 5") {
			t.Errorf("Workflow step 1 expected 'Result: 5', got: %s", out1.String())
		}

		// Step 2: Run another calculation with different input
		cmd2 := exec.Command("go", "run", "main.go", "-e", "3 * (4 + 1)")
		var out2 bytes.Buffer
		cmd2.Stdout = &out2
		if err := cmd2.Run(); err != nil {
			t.Fatalf("Workflow step 2 failed: %v", err)
		}
		if !strings.Contains(out2.String(), "Result: 15") {
			t.Errorf("Workflow step 2 expected 'Result: 15', got: %s", out2.String())
		}

		// Step 3: Verify error handling
		cmd3 := exec.Command("go", "run", "main.go", "-e", "invalid")
		var errOut3 bytes.Buffer
		cmd3.Stderr = &errOut3
		if err := cmd3.Run(); err == nil {
			t.Fatal("Workflow step 3 expected an error, but got none")
		}
		if !strings.Contains(errOut3.String(), "Error: Invalid expression") {
			t.Errorf("Workflow step 3 expected error message not found, got: %s", errOut3.String())
		}
	})
}
