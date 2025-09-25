package e2e

import (
	"os"
	"strings"
	"testing"
)

// TestMain sets up the E2E test environment.
func TestMain(m *testing.M) {
	teardown, err := setupE2ETests()
	if err != nil {
		os.Stderr.WriteString("Failed to setup E2E tests: " + err.Error() + "\n")
		os.Exit(1)
	}

	exitCode := m.Run()

	teardown()

	os.Exit(exitCode)
}

func TestE2EWorkflow(t *testing.T) {
	t.Run("TestSimpleAddition", func(t *testing.T) {
		run := NewE2ETestRun(t)

		expression := "2 + 2"
		expected := "4"

		// Record the command
		run.RecordCommand(expression)

		// Run the command for assertion
		output := run.RunCommand(expression)

		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s', but got '%s'", expected, output)
		}
	})

	t.Run("TestComplexExpression", func(t *testing.T) {
		run := NewE2ETestRun(t)

		expression := "(5 + 3) * 2 - 4"
		expected := "12"

		// Record the command
		run.RecordCommand(expression)

		// Run the command for assertion
		output := run.RunCommand(expression)

		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s', but got '%s'", expected, output)
		}
	})
}
