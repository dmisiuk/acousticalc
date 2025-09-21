package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// getExecutablePath returns the path to the acousticalc executable
func getExecutablePath() (string, error) {
	// Get the current file's directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", os.ErrNotExist
	}

	// Build the path to the executable
	// The executable is in the same directory as the test file
	dir := filepath.Dir(filename)
	executableName := "acousticalc"
	if runtime.GOOS == "windows" {
		executableName += ".exe"
	}
	return filepath.Join(dir, executableName), nil
}

// TestCLIValidExpressions tests the CLI with valid mathematical expressions
func TestCLIValidExpressions(t *testing.T) {
	executable, err := getExecutablePath()
	if err != nil {
		t.Skipf("Could not find executable: %v", err)
	}

	testCases := []struct {
		name       string
		expression string
		expected   string
	}{
		{
			name:       "Basic addition",
			expression: "2 + 3",
			expected:   "Result: 5",
		},
		{
			name:       "Operator precedence",
			expression: "2 + 3 * 4",
			expected:   "Result: 14",
		},
		{
			name:       "Parentheses",
			expression: "(2 + 3) * 4",
			expected:   "Result: 20",
		},
		{
			name:       "Decimal numbers",
			expression: "3.5 + 2.1",
			expected:   "Result: 5.6",
		},
		{
			name:       "Complex expression",
			expression: "2 * (3 + 4) - 5 / 2",
			expected:   "Result: 11.5",
		},
		{
			name:       "Negative numbers",
			expression: "-5 + 3",
			expected:   "Result: -2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(executable, tc.expression)
			output, err := cmd.CombinedOutput()

			if err != nil {
				t.Fatalf("CLI command failed: %v", err)
			}

			actual := strings.TrimSpace(string(output))
			if actual != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

// TestCLIInvalidExpressions tests the CLI with invalid mathematical expressions
func TestCLIInvalidExpressions(t *testing.T) {
	executable, err := getExecutablePath()
	if err != nil {
		t.Skipf("Could not find executable: %v", err)
	}

	testCases := []struct {
		name       string
		expression string
		expected   string
	}{
		{
			name:       "Division by zero",
			expression: "10 / 0",
			expected:   "Error: division by zero",
		},
		{
			name:       "Invalid syntax",
			expression: "2 +",
			expected:   "Error:",
		},
		{
			name:       "Invalid character",
			expression: "2 + a",
			expected:   "Error:",
		},
		{
			name:       "Empty expression",
			expression: "",
			expected:   "Error: empty expression",
		},
		{
			name:       "Mismatched parentheses",
			expression: "(2 + 3",
			expected:   "Error:",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(executable, tc.expression)
			output, err := cmd.CombinedOutput()

			// We expect the command to fail for invalid expressions
			if err == nil {
				t.Errorf("Expected command to fail for invalid expression %q", tc.expression)
			}

			actual := strings.TrimSpace(string(output))
			if !strings.Contains(actual, tc.expected) {
				t.Errorf("Expected output to contain %q, got %q", tc.expected, actual)
			}
		})
	}
}

// TestCLINoArguments tests the CLI when no arguments are provided
func TestCLINoArguments(t *testing.T) {
	executable, err := getExecutablePath()
	if err != nil {
		t.Skipf("Could not find executable: %v", err)
	}

	cmd := exec.Command(executable)
	output, err := cmd.CombinedOutput()

	// We expect the command to fail when no arguments are provided
	if err == nil {
		t.Error("Expected command to fail when no arguments are provided")
	}

	actual := strings.TrimSpace(string(output))
	expectedUsage := "Usage: acousticalc <expression>"
	if !strings.Contains(actual, expectedUsage) {
		t.Errorf("Expected usage message, got: %s", actual)
	}
}
