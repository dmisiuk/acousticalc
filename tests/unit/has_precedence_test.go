package unit

import (
	"testing"

	"github.com/dmisiuk/acousticalc/pkg/calculator"
)

// TestHasPrecedenceLowerPrecedenceOp1 tests the case where op1 has lower precedence than op2
// This should trigger the final return false in hasPrecedence function
func TestHasPrecedenceLowerPrecedenceOp1(t *testing.T) {
	// Test expressions where addition/subtraction comes before multiplication/division
	// In these cases, op1 ("+" or "-") has lower precedence than op2 ("*" or "/")
	// So hasPrecedence should return false, meaning op2 has higher precedence
	tests := []struct {
		name       string
		expression string
		expected   float64
	}{
		{"Addition before multiplication", "2 + 3 * 4", 14},    // + has lower precedence than *
		{"Subtraction before multiplication", "10 - 3 * 2", 4}, // - has lower precedence than *
		{"Addition before division", "5 + 6 / 2", 8},           // + has lower precedence than /
		{"Subtraction before division", "10 - 8 / 4", 8},       // - has lower precedence than /
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.Evaluate(tt.expression)
			if err != nil {
				t.Fatalf("Unexpected error for expression '%s': %v", tt.expression, err)
			}
			if result != tt.expected {
				t.Errorf("For expression '%s': expected %v, got %v", tt.expression, tt.expected, result)
			}
		})
	}
}

// TestHasPrecedenceHigherPrecedenceOp1 tests the case where op1 has higher precedence than op2
func TestHasPrecedenceHigherPrecedenceOp1(t *testing.T) {
	// Test expressions where multiplication/division comes before addition/subtraction
	// In these cases, op1 ("*" or "/") has higher precedence than op2 ("+" or "-")
	// So hasPrecedence should return true
	tests := []struct {
		name       string
		expression string
		expected   float64
	}{
		{"Multiplication before addition", "2 * 3 + 4", 10},   // * has higher precedence than +
		{"Division before addition", "10 / 2 + 3", 8},         // / has higher precedence than +
		{"Multiplication before subtraction", "3 * 4 - 5", 7}, // * has higher precedence than -
		{"Division before subtraction", "10 / 2 - 3", 2},      // / has higher precedence than -
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.Evaluate(tt.expression)
			if err != nil {
				t.Fatalf("Unexpected error for expression '%s': %v", tt.expression, err)
			}
			if result != tt.expected {
				t.Errorf("For expression '%s': expected %v, got %v", tt.expression, tt.expected, result)
			}
		})
	}
}
