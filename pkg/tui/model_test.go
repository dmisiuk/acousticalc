package tui

import (
	"testing"
)

func TestNewModel(t *testing.T) {
	model := NewModel()

	// Test initial state
	if model.display != "0" {
		t.Errorf("Expected initial display to be '0', got '%s'", model.display)
	}

	if model.expression != "" {
		t.Errorf("Expected initial expression to be empty, got '%s'", model.expression)
	}

	if model.hasResult {
		t.Error("Expected hasResult to be false initially")
	}

	if model.hasError {
		t.Error("Expected hasError to be false initially")
	}

	if model.selectedRow != 0 || model.selectedCol != 0 {
		t.Errorf("Expected initial selection to be (0,0), got (%d,%d)", model.selectedRow, model.selectedCol)
	}

	if !model.mouseEnabled {
		t.Error("Expected mouse to be enabled by default")
	}

	// Test buttons initialization
	if len(model.buttons) == 0 {
		t.Error("Expected buttons to be initialized")
	}

	// Check for essential buttons
	hasNumbers := false
	hasOperators := false
	hasClear := false
	hasEquals := false

	for _, button := range model.buttons {
		switch button.Value {
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			hasNumbers = true
		case "+", "-", "*", "/":
			hasOperators = true
		case "C":
			hasClear = true
		case "=":
			hasEquals = true
		}
	}

	if !hasNumbers {
		t.Error("Expected number buttons to be present")
	}
	if !hasOperators {
		t.Error("Expected operator buttons to be present")
	}
	if !hasClear {
		t.Error("Expected clear button to be present")
	}
	if !hasEquals {
		t.Error("Expected equals button to be present")
	}
}

func TestGetButtonAt(t *testing.T) {
	model := NewModel()

	// Test getting a button that exists
	button := model.getButtonAt(0, 0)
	if button == nil {
		t.Error("Expected to find button at (0,0)")
	}

	// Test getting a button that doesn't exist
	button = model.getButtonAt(10, 10)
	if button != nil {
		t.Error("Expected no button at (10,10)")
	}
}

func TestMoveSelection(t *testing.T) {
	model := NewModel()
	
	// Test moving right
	initialCol := model.selectedCol
	model.moveSelection(0, 1)
	if model.selectedCol <= initialCol && model.getButtonAt(model.selectedRow, model.selectedCol) != nil {
		// Only check if there's actually a button to move to
		t.Error("Expected selection to move right")
	}

	// Test moving down
	initialRow := model.selectedRow
	model.moveSelection(1, 0)
	if model.selectedRow <= initialRow && model.getButtonAt(model.selectedRow, model.selectedCol) != nil {
		t.Error("Expected selection to move down")
	}
}

func TestClear(t *testing.T) {
	model := NewModel()
	
	// Set some state
	model.display = "123"
	model.expression = "1+2"
	model.result = 3
	model.hasResult = true
	model.hasError = true
	model.errorMsg = "test error"

	// Clear
	model.clear()

	// Check state is reset
	if model.display != "0" {
		t.Errorf("Expected display to be '0' after clear, got '%s'", model.display)
	}
	if model.expression != "" {
		t.Errorf("Expected expression to be empty after clear, got '%s'", model.expression)
	}
	if model.result != 0 {
		t.Errorf("Expected result to be 0 after clear, got %f", model.result)
	}
	if model.hasResult {
		t.Error("Expected hasResult to be false after clear")
	}
	if model.hasError {
		t.Error("Expected hasError to be false after clear")
	}
	if model.errorMsg != "" {
		t.Errorf("Expected errorMsg to be empty after clear, got '%s'", model.errorMsg)
	}
}

func TestAppendToExpression(t *testing.T) {
	model := NewModel()

	// Test appending to empty expression
	model.appendToExpression("1")
	if model.expression != "1" {
		t.Errorf("Expected expression to be '1', got '%s'", model.expression)
	}
	if model.display != "1" {
		t.Errorf("Expected display to be '1', got '%s'", model.display)
	}

	// Test appending more
	model.appendToExpression("+")
	if model.expression != "1+" {
		t.Errorf("Expected expression to be '1+', got '%s'", model.expression)
	}

	model.appendToExpression("2")
	if model.expression != "1+2" {
		t.Errorf("Expected expression to be '1+2', got '%s'", model.expression)
	}
}

func TestCalculate(t *testing.T) {
	model := NewModel()

	// Test valid calculation
	model.expression = "2+3"
	model.calculate()

	if !model.hasResult {
		t.Error("Expected hasResult to be true after calculation")
	}
	if model.result != 5 {
		t.Errorf("Expected result to be 5, got %f", model.result)
	}
	if model.hasError {
		t.Error("Expected no error for valid calculation")
	}

	// Test invalid calculation
	model.expression = "2/0"
	model.calculate()

	if !model.hasError {
		t.Error("Expected error for division by zero")
	}
	if model.display != "Error" {
		t.Errorf("Expected display to show 'Error', got '%s'", model.display)
	}
}

func TestToggleSign(t *testing.T) {
	model := NewModel()

	// Test with result
	model.result = 5
	model.hasResult = true
	model.display = "5"
	model.expression = "5"

	model.toggleSign()

	if model.result != -5 {
		t.Errorf("Expected result to be -5, got %f", model.result)
	}

	// Test with expression
	model.hasResult = false
	model.expression = "123"
	model.display = "123"

	model.toggleSign()

	if model.expression != "-123" {
		t.Errorf("Expected expression to be '-123', got '%s'", model.expression)
	}

	// Test toggling back
	model.toggleSign()

	if model.expression != "123" {
		t.Errorf("Expected expression to be '123', got '%s'", model.expression)
	}
}

func TestPercentage(t *testing.T) {
	model := NewModel()

	// Test with result
	model.result = 50
	model.hasResult = true
	model.display = "50"
	model.expression = "50"

	model.percentage()

	if model.result != 0.5 {
		t.Errorf("Expected result to be 0.5, got %f", model.result)
	}

	// Test with expression
	model.hasResult = false
	model.expression = "25"

	model.percentage()

	if model.result != 0.25 {
		t.Errorf("Expected result to be 0.25, got %f", model.result)
	}
}

func TestPressButton(t *testing.T) {
	model := NewModel()

	// Test pressing number button
	button := &Button{Label: "5", Value: "5", Row: 1, Col: 1}
	model.pressButton(button)

	if model.expression != "5" {
		t.Errorf("Expected expression to be '5', got '%s'", model.expression)
	}
	if model.pressedButton != button {
		t.Error("Expected pressedButton to be set")
	}
	if model.feedbackTimer == 0 {
		t.Error("Expected feedbackTimer to be set")
	}

	// Test pressing clear button
	clearButton := &Button{Label: "C", Value: "C", Row: 0, Col: 0}
	model.pressButton(clearButton)

	if model.expression != "" {
		t.Errorf("Expected expression to be empty after clear, got '%s'", model.expression)
	}
	if model.display != "0" {
		t.Errorf("Expected display to be '0' after clear, got '%s'", model.display)
	}
}

func TestIsOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"+", true},
		{"-", true},
		{"*", true},
		{"/", true},
		{"1", false},
		{"=", false},
		{"C", false},
		{"", false},
	}

	for _, test := range tests {
		result := isOperator(test.input)
		if result != test.expected {
			t.Errorf("isOperator(%s) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{1.5, "1.5"},
		{1.0, "1"},
		{123.456, "123.456"},
		{123.000, "123"},
		{0.1, "0.1"},
		{-5, "-5"},
		{-5.5, "-5.5"},
	}

	for _, test := range tests {
		result := formatNumber(test.input)
		if result != test.expected {
			t.Errorf("formatNumber(%f) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestGetters(t *testing.T) {
	model := NewModel()
	
	// Set some test state
	model.display = "test display"
	model.expression = "test expression"
	model.hasError = true
	model.errorMsg = "test error"
	model.hasResult = true
	model.result = 42.5

	// Test getters
	if model.GetDisplayValue() != "test display" {
		t.Errorf("GetDisplayValue() = %s, expected 'test display'", model.GetDisplayValue())
	}
	if model.GetExpression() != "test expression" {
		t.Errorf("GetExpression() = %s, expected 'test expression'", model.GetExpression())
	}
	if !model.HasError() {
		t.Error("HasError() = false, expected true")
	}
	if model.GetErrorMessage() != "test error" {
		t.Errorf("GetErrorMessage() = %s, expected 'test error'", model.GetErrorMessage())
	}
	if !model.HasResult() {
		t.Error("HasResult() = false, expected true")
	}
	if model.GetResult() != 42.5 {
		t.Errorf("GetResult() = %f, expected 42.5", model.GetResult())
	}
}

func TestSetMouseEnabled(t *testing.T) {
	model := NewModel()
	
	// Test disabling mouse
	model.SetMouseEnabled(false)
	if model.mouseEnabled {
		t.Error("Expected mouse to be disabled")
	}
	
	// Test enabling mouse
	model.SetMouseEnabled(true)
	if !model.mouseEnabled {
		t.Error("Expected mouse to be enabled")
	}
}