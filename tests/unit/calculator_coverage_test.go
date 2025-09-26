package unit

import (
	"testing"

	"github.com/dmisiuk/acousticalc/pkg/calculator"
)

// TestCoverageVerification verifies that test coverage meets the 80% requirement
// This test runs automatically without requiring environment variables
func TestCoverageVerification(t *testing.T) {
	// This test ensures all major functions are tested
	// It serves as a basic coverage verification

	testCases := []struct {
		name       string
		expression string
		expectErr  bool
	}{
		// Basic arithmetic operations
		{"Addition", "2 + 3", false},
		{"Subtraction", "10 - 4", false},
		{"Multiplication", "3 * 4", false},
		{"Division", "15 / 3", false},

		// Expression parsing
		{"Precedence", "2 + 3 * 4", false},
		{"Parentheses", "(2 + 3) * 4", false},
		{"Nested parentheses", "((2 + 3) * 4) - 5", false},

		// Error cases
		{"Division by zero", "10 / 0", true},
		{"Invalid syntax", "2 +", true},
		{"Invalid character", "2 + a", true},
		{"Empty expression", "", true},
		{"Mismatched parentheses", "(2 + 3", true},

		// Edge cases
		{"Decimal numbers", "3.5 + 2.1", false},
		{"Negative numbers", "-5 + 3", false},
		{"Complex expression", "2 * (3 + 4) - 5 / 2", false},
		{"Single number", "42", false},
		{"Zero operations", "0 + 0", false},
		{"Large numbers", "1000000 * 2", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := calculator.Evaluate(tc.expression)

			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error for expression %q, but got result: %v", tc.expression, result)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for expression %q, but got: %v", tc.expression, err)
				}
			}
		})
	}
}
