package tui

import (
	"strings"
	"testing"
)

func TestView(t *testing.T) {
	model := NewModel()

	// Test basic view rendering
	view := model.View()
	
	if view == "" {
		t.Error("Expected view to return non-empty string")
	}

	// Check for essential components
	if !strings.Contains(view, "Acoustic Calculator") {
		t.Error("Expected view to contain title 'Acoustic Calculator'")
	}

	if !strings.Contains(view, "0") {
		t.Error("Expected view to contain initial display value '0'")
	}

	// Test view with error
	model.hasError = true
	model.errorMsg = "Test error"
	model.display = "Error"

	errorView := model.View()
	
	if !strings.Contains(errorView, "Error") {
		t.Error("Expected error view to contain 'Error'")
	}
	if !strings.Contains(errorView, "Test error") {
		t.Error("Expected error view to contain error message")
	}

	// Test view with result
	model = NewModel()
	model.hasResult = true
	model.result = 42.5
	model.display = "42.5"

	resultView := model.View()
	
	if !strings.Contains(resultView, "42.5") {
		t.Error("Expected result view to contain result value")
	}
}

func TestRenderButtonGrid(t *testing.T) {
	model := NewModel()

	grid := model.renderButtonGrid()
	
	if grid == "" {
		t.Error("Expected button grid to return non-empty string")
	}

	// Check that essential buttons are rendered
	expectedButtons := []string{"C", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "+", "-", "×", "÷", "="}
	
	for _, button := range expectedButtons {
		if !strings.Contains(grid, button) {
			t.Errorf("Expected button grid to contain button '%s'", button)
		}
	}
}

func TestRenderButtonRow(t *testing.T) {
	model := NewModel()

	// Create test buttons for a row
	testButtons := []Button{
		{Label: "1", Value: "1", Row: 0, Col: 0, Width: 1, Height: 1},
		{Label: "2", Value: "2", Row: 0, Col: 1, Width: 1, Height: 1},
		{Label: "3", Value: "3", Row: 0, Col: 2, Width: 1, Height: 1},
		{Label: "+", Value: "+", Row: 0, Col: 3, Width: 1, Height: 1},
	}

	row := model.renderButtonRow(testButtons, 0)
	
	if row == "" {
		t.Error("Expected button row to return non-empty string")
	}

	// Check that all buttons are present
	for _, button := range testButtons {
		if !strings.Contains(row, button.Label) {
			t.Errorf("Expected button row to contain button '%s'", button.Label)
		}
	}
}

func TestRenderButton(t *testing.T) {
	model := NewModel()

	// Test regular number button
	numberButton := Button{Label: "5", Value: "5", Row: 1, Col: 1, Width: 1, Height: 1}
	rendered := model.renderButton(numberButton)
	
	if rendered == "" {
		t.Error("Expected rendered button to return non-empty string")
	}
	if !strings.Contains(rendered, "5") {
		t.Error("Expected rendered button to contain label '5'")
	}

	// Test operator button
	operatorButton := Button{Label: "+", Value: "+", Row: 1, Col: 3, Width: 1, Height: 1}
	rendered = model.renderButton(operatorButton)
	
	if !strings.Contains(rendered, "+") {
		t.Error("Expected rendered operator button to contain label '+'")
	}

	// Test clear button
	clearButton := Button{Label: "C", Value: "C", Row: 0, Col: 0, Width: 1, Height: 1}
	rendered = model.renderButton(clearButton)
	
	if !strings.Contains(rendered, "C") {
		t.Error("Expected rendered clear button to contain label 'C'")
	}

	// Test equals button
	equalsButton := Button{Label: "=", Value: "=", Row: 4, Col: 3, Width: 1, Height: 1}
	rendered = model.renderButton(equalsButton)
	
	if !strings.Contains(rendered, "=") {
		t.Error("Expected rendered equals button to contain label '='")
	}

	// Test wide button (like "0")
	wideButton := Button{Label: "0", Value: "0", Row: 4, Col: 0, Width: 2, Height: 1}
	rendered = model.renderButton(wideButton)
	
	if !strings.Contains(rendered, "0") {
		t.Error("Expected rendered wide button to contain label '0'")
	}
}

func TestRenderButtonWithSelection(t *testing.T) {
	model := NewModel()

	// Test selected button
	model.selectedRow = 1
	model.selectedCol = 1
	
	selectedButton := Button{Label: "5", Value: "5", Row: 1, Col: 1, Width: 1, Height: 1}
	rendered := model.renderButton(selectedButton)
	
	if rendered == "" {
		t.Error("Expected rendered selected button to return non-empty string")
	}
	// We can't easily test the styling without parsing ANSI codes, but we can verify it renders

	// Test non-selected button
	nonSelectedButton := Button{Label: "6", Value: "6", Row: 1, Col: 2, Width: 1, Height: 1}
	rendered = model.renderButton(nonSelectedButton)
	
	if rendered == "" {
		t.Error("Expected rendered non-selected button to return non-empty string")
	}
}

func TestRenderButtonWithPressedState(t *testing.T) {
	model := NewModel()

	// Test pressed button
	pressedButton := Button{Label: "5", Value: "5", Row: 1, Col: 1, Width: 1, Height: 1}
	model.pressedButton = &pressedButton
	model.feedbackTimer = 2

	rendered := model.renderButton(pressedButton)
	
	if rendered == "" {
		t.Error("Expected rendered pressed button to return non-empty string")
	}
	if !strings.Contains(rendered, "5") {
		t.Error("Expected rendered pressed button to contain label '5'")
	}

	// Test button that's not pressed
	notPressedButton := Button{Label: "6", Value: "6", Row: 1, Col: 2, Width: 1, Height: 1}
	rendered = model.renderButton(notPressedButton)
	
	if rendered == "" {
		t.Error("Expected rendered non-pressed button to return non-empty string")
	}
}

func TestGetButtonDimensions(t *testing.T) {
	model := NewModel()

	width, height := model.GetButtonDimensions()
	
	if width <= 0 {
		t.Errorf("Expected positive width, got %d", width)
	}
	if height <= 0 {
		t.Errorf("Expected positive height, got %d", height)
	}

	// Check that dimensions are reasonable for a calculator
	if width < 20 {
		t.Errorf("Expected width to be at least 20 for calculator layout, got %d", width)
	}
	if height < 10 {
		t.Errorf("Expected height to be at least 10 for calculator layout, got %d", height)
	}
}

func TestViewWithDifferentStates(t *testing.T) {
	// Test view with different calculator states
	testCases := []struct {
		name        string
		setupModel  func(*Model)
		expectInView string
	}{
		{
			name: "initial state",
			setupModel: func(m *Model) {
				// Default state
			},
			expectInView: "0",
		},
		{
			name: "with expression",
			setupModel: func(m *Model) {
				m.expression = "2+3"
				m.display = "2+3"
			},
			expectInView: "2+3",
		},
		{
			name: "with result",
			setupModel: func(m *Model) {
				m.hasResult = true
				m.result = 5
				m.display = "5"
			},
			expectInView: "5",
		},
		{
			name: "with error",
			setupModel: func(m *Model) {
				m.hasError = true
				m.errorMsg = "Division by zero"
				m.display = "Error"
			},
			expectInView: "Error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model := NewModel()
			tc.setupModel(&model)
			
			view := model.View()
			
			if !strings.Contains(view, tc.expectInView) {
				t.Errorf("Expected view to contain '%s' for %s", tc.expectInView, tc.name)
			}
		})
	}
}

func TestViewContainsHelpText(t *testing.T) {
	model := NewModel()
	view := model.View()

	// Check for help text elements
	helpElements := []string{
		"Arrow keys",
		"Navigate",
		"Enter",
		"Space",
		"Press",
		"Mouse",
		"Click",
		"Quit",
	}

	for _, element := range helpElements {
		if !strings.Contains(view, element) {
			t.Errorf("Expected view to contain help text element '%s'", element)
		}
	}
}

func TestViewStructure(t *testing.T) {
	model := NewModel()
	view := model.View()

	// Basic structure checks
	if view == "" {
		t.Fatal("View should not be empty")
	}

	// Should contain multiple lines
	lines := strings.Split(view, "\n")
	if len(lines) < 5 {
		t.Errorf("Expected view to have at least 5 lines, got %d", len(lines))
	}

	// Should not contain obvious rendering errors
	if strings.Contains(view, "ERROR") || strings.Contains(view, "PANIC") {
		t.Error("View contains rendering errors")
	}
}