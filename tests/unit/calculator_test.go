package unit

import (
	"testing"

	"github.com/dmisiuk/acousticalc/pkg/calculator"
)

// Test basic arithmetic operations (AC1)
func TestBasicArithmeticOperations(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		expected   float64
	}{
		{"Addition", "2 + 3", 5},
		{"Subtraction", "10 - 4", 6},
		{"Multiplication", "3 * 4", 12},
		{"Division", "15 / 3", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.Evaluate(tt.expression)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test expression parsing with operator precedence (AC2)
func TestExpressionParsing(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		expected   float64
	}{
		{"Operator precedence", "2 + 3 * 4", 14},
		{"Operator precedence with division", "10 - 6 / 2", 7},
		{"Expression with parentheses", "(2 + 3) * 4", 20},
		{"Nested parentheses", "((2 + 3) * 4) - 5", 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.Evaluate(tt.expression)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test error handling (AC4)
func TestDivisionByZero(t *testing.T) {
	_, err := calculator.Evaluate("10 / 0")
	if err == nil {
		t.Error("Expected division by zero error")
		return
	}
	if err.Error() != "division by zero" {
		t.Errorf("Expected 'division by zero' error, got %v", err)
	}
}

func TestInvalidSyntax(t *testing.T) {
	_, err := calculator.Evaluate("2 +")
	if err == nil {
		t.Error("Expected syntax error")
		return
	}
}

func TestInvalidCharacter(t *testing.T) {
	_, err := calculator.Evaluate("2 + a")
	if err == nil {
		t.Error("Expected invalid character error")
		return
	}
}

func TestEmptyExpression(t *testing.T) {
	_, err := calculator.Evaluate("")
	if err == nil {
		t.Error("Expected empty expression error")
		return
	}
	if err.Error() != "empty expression" {
		t.Errorf("Expected 'empty expression' error, got %v", err)
	}
}

func TestMismatchedParentheses(t *testing.T) {
	_, err := calculator.Evaluate("(2 + 3")
	if err == nil {
		t.Error("Expected mismatched parentheses error")
		return
	}
}

func TestMismatchedParenthesesClosing(t *testing.T) {
	_, err := calculator.Evaluate("2 + 3)")
	if err == nil {
		t.Error("Expected mismatched parentheses error")
		return
	}
}

// Test decimal numbers
func TestDecimalNumbers(t *testing.T) {
	result, err := calculator.Evaluate("3.5 + 2.1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 5.6 {
		t.Errorf("Expected 5.6, got %v", result)
	}
}

func TestDecimalDivision(t *testing.T) {
	result, err := calculator.Evaluate("7.5 / 2.5")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 3 {
		t.Errorf("Expected 3, got %v", result)
	}
}

// Test complex expressions
func TestComplexExpression(t *testing.T) {
	result, err := calculator.Evaluate("2 * (3 + 4) - 5 / 2")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 11.5 {
		t.Errorf("Expected 11.5, got %v", result)
	}
}

func TestNegativeNumbers(t *testing.T) {
	result, err := calculator.Evaluate("-5 + 3")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != -2 {
		t.Errorf("Expected -2, got %v", result)
	}
}

func TestMultipleOperators(t *testing.T) {
	result, err := calculator.Evaluate("10 + 5 - 3 * 2 / 4")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	// Following order of operations: 10 + 5 - (3 * 2 / 4) = 10 + 5 - 1.5 = 13.5
	if result != 13.5 {
		t.Errorf("Expected 13.5, got %v", result)
	}
}

// Test edge cases
func TestSingleNumber(t *testing.T) {
	result, err := calculator.Evaluate("42")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 42 {
		t.Errorf("Expected 42, got %v", result)
	}
}

func TestZeroOperations(t *testing.T) {
	result, err := calculator.Evaluate("0 + 0")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 0 {
		t.Errorf("Expected 0, got %v", result)
	}
}

func TestLargeNumbers(t *testing.T) {
	result, err := calculator.Evaluate("1000000 * 2")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 2000000 {
		t.Errorf("Expected 2000000, got %v", result)
	}
}

// Additional tests for negative numbers and complex expressions
func TestNegativeNumberAtStart(t *testing.T) {
	result, err := calculator.Evaluate("-10 + 5")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != -5 {
		t.Errorf("Expected -5, got %v", result)
	}
}

func TestNegativeNumberAfterOperator(t *testing.T) {
	result, err := calculator.Evaluate("10 - -5")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 15 {
		t.Errorf("Expected 15, got %v", result)
	}
}

func TestComplexExpressionWithNegatives(t *testing.T) {
	result, err := calculator.Evaluate("-3 * (2 + -4) - -5")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 11 {
		t.Errorf("Expected 11, got %v", result)
	}
}

func TestExpressionStartingWithParenthesesAndNegative(t *testing.T) {
	result, err := calculator.Evaluate("(-5 + 3) * 2")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != -4 {
		t.Errorf("Expected -4, got %v", result)
	}
}

// Tests for floating point precision issues (TECH-002)
func TestFloatingPointPrecision(t *testing.T) {
	result, err := calculator.Evaluate("0.1 + 0.2")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	// Check that the result is close to 0.3, allowing for floating point precision issues
	if result < 0.299999 || result > 0.300000001 {
		t.Errorf("Expected approximately 0.3, got %v", result)
	}
}

func TestFloatingPointPrecisionSubtraction(t *testing.T) {
	result, err := calculator.Evaluate("0.3 - 0.2")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	// Check that the result is close to 0.1, allowing for floating point precision issues
	if result < 0.09999999 || result > 0.1000001 {
		t.Errorf("Expected approximately 0.1, got %v", result)
	}
}

func TestFloatingPointPrecisionComplex(t *testing.T) {
	result, err := calculator.Evaluate("0.1 + 0.2 - 0.3")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	// Check that the result is close to 0, allowing for floating point precision issues
	if result < -0.0000000001 || result > 0.00000001 {
		t.Errorf("Expected approximately 0, got %v", result)
	}
}
