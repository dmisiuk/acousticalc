package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dmisiuk/acousticalc/pkg/calculator"
)

// Button represents a calculator button
type Button struct {
	Label    string
	Value    string
	Row      int
	Col      int
	Width    int
	Height   int
	IsActive bool
}

// Model represents the TUI state
type Model struct {
	// Calculator state
	display     string
	expression  string
	result      float64
	hasResult   bool
	hasError    bool
	errorMsg    string

	// UI state
	buttons       []Button
	selectedRow   int
	selectedCol   int
	width         int
	height        int
	mouseEnabled  bool

	// Visual feedback
	pressedButton *Button
	feedbackTimer int
}

// NewModel creates a new TUI model
func NewModel() Model {
	m := Model{
		display:      "0",
		expression:   "",
		selectedRow:  0,
		selectedCol:  0,
		mouseEnabled: true,
	}

	// Initialize calculator buttons
	m.initButtons()

	return m
}

// initButtons initializes the calculator button layout
func (m *Model) initButtons() {
	// Calculator button layout (4x5 grid)
	buttonLayout := [][]string{
		{"C", "±", "%", "÷"},
		{"7", "8", "9", "×"},
		{"4", "5", "6", "-"},
		{"1", "2", "3", "+"},
		{"0", "", ".", "="},
	}

	// Map display symbols to actual operators
	operatorMap := map[string]string{
		"÷": "/",
		"×": "*",
		"±": "±", // Special case for sign toggle
		"%": "%", // Special case for percentage
		"C": "C", // Special case for clear
		"=": "=", // Special case for equals
	}

	m.buttons = make([]Button, 0)

	for row, buttonRow := range buttonLayout {
		for col, label := range buttonRow {
			if label == "" {
				continue // Skip empty buttons
			}

			value := label
			if mappedValue, exists := operatorMap[label]; exists {
				value = mappedValue
			}

			// Special width for "0" button
			width := 1
			if label == "0" {
				width = 2
			}

			button := Button{
				Label:  label,
				Value:  value,
				Row:    row,
				Col:    col,
				Width:  width,
				Height: 1,
			}

			m.buttons = append(m.buttons, button)
		}
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.EnableMouseCellMotion
}

// getButtonAt returns the button at the given row and column
func (m *Model) getButtonAt(row, col int) *Button {
	for i := range m.buttons {
		button := &m.buttons[i]
		if button.Row == row && button.Col == col {
			return button
		}
		// Handle wide buttons (like "0")
		if button.Row == row && button.Col <= col && col < button.Col+button.Width {
			return button
		}
	}
	return nil
}

// getSelectedButton returns the currently selected button
func (m *Model) getSelectedButton() *Button {
	return m.getButtonAt(m.selectedRow, m.selectedCol)
}

// moveSelection moves the selection in the given direction
func (m *Model) moveSelection(deltaRow, deltaCol int) {
	newRow := m.selectedRow + deltaRow
	newCol := m.selectedCol + deltaCol

	// Clamp to valid ranges
	if newRow < 0 {
		newRow = 4 // Wrap to bottom
	} else if newRow > 4 {
		newRow = 0 // Wrap to top
	}

	if newCol < 0 {
		newCol = 3 // Wrap to right
	} else if newCol > 3 {
		newCol = 0 // Wrap to left
	}

	// Check if there's a button at the new position
	if m.getButtonAt(newRow, newCol) != nil {
		m.selectedRow = newRow
		m.selectedCol = newCol
	}
}

// pressButton handles button press logic
func (m *Model) pressButton(button *Button) {
	if button == nil {
		return
	}

	// Set visual feedback
	m.pressedButton = button
	m.feedbackTimer = 3 // Show feedback for 3 update cycles

	// Handle button logic
	switch button.Value {
	case "C":
		m.clear()
	case "=":
		m.calculate()
	case "±":
		m.toggleSign()
	case "%":
		m.percentage()
	default:
		m.appendToExpression(button.Value)
	}
}

// clear resets the calculator
func (m *Model) clear() {
	m.display = "0"
	m.expression = ""
	m.result = 0
	m.hasResult = false
	m.hasError = false
	m.errorMsg = ""
}

// calculate evaluates the current expression
func (m *Model) calculate() {
	if m.expression == "" {
		return
	}

	result, err := calculator.Evaluate(m.expression)
	if err != nil {
		m.hasError = true
		m.errorMsg = err.Error()
		m.display = "Error"
	} else {
		m.result = result
		m.hasResult = true
		m.hasError = false
		m.display = formatNumber(result)
		m.expression = m.display // Allow continuing with result
	}
}

// toggleSign toggles the sign of the current number
func (m *Model) toggleSign() {
	if m.hasResult {
		m.result = -m.result
		m.display = formatNumber(m.result)
		m.expression = m.display
	} else if m.expression != "" {
		// Simple implementation: prepend minus or remove it
		if strings.HasPrefix(m.expression, "-") {
			m.expression = m.expression[1:]
		} else {
			m.expression = "-" + m.expression
		}
		m.display = m.expression
	}
}

// percentage converts the current number to percentage
func (m *Model) percentage() {
	if m.hasResult {
		m.result = m.result / 100
		m.display = formatNumber(m.result)
		m.expression = m.display
	} else if m.expression != "" {
		// Try to evaluate current expression and convert to percentage
		if result, err := calculator.Evaluate(m.expression); err == nil {
			m.result = result / 100
			m.hasResult = true
			m.display = formatNumber(m.result)
			m.expression = m.display
		}
	}
}

// appendToExpression adds a character to the current expression
func (m *Model) appendToExpression(value string) {
	if m.hasResult && isOperator(value) {
		// Continue with result if adding an operator
		m.expression = m.display + value
		m.hasResult = false
	} else if m.hasResult {
		// Start new expression if adding a number
		m.expression = value
		m.hasResult = false
	} else {
		m.expression += value
	}

	m.display = m.expression
	m.hasError = false
}

// isOperator checks if a string is an operator
func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/"
}

// formatNumber formats a number for display
func formatNumber(n float64) string {
	// Remove trailing zeros and decimal point if not needed
	str := fmt.Sprintf("%.10f", n)
	str = strings.TrimRight(str, "0")
	str = strings.TrimRight(str, ".")
	return str
}