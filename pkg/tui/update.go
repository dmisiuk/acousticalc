package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// tickMsg is sent on every timer tick for visual feedback
type tickMsg time.Time

// Update handles all incoming messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	
	case tea.MouseMsg:
		return m.handleMouseEvent(msg)
	
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	
	case tickMsg:
		return m.handleTick()
	}

	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "up":
		m.moveSelection(-1, 0)
		return m, nil

	case "down":
		m.moveSelection(1, 0)
		return m, nil

	case "left":
		m.moveSelection(0, -1)
		return m, nil

	case "right":
		m.moveSelection(0, 1)
		return m, nil

	case "enter", " ":
		// Press the currently selected button
		if button := m.getSelectedButton(); button != nil {
			m.pressButton(button)
			return m, m.startFeedbackTimer()
		}
		return m, nil

	// Direct number input
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		m.appendToExpression(msg.String())
		return m, nil

	// Direct operator input
	case "+":
		m.appendToExpression("+")
		return m, nil
	case "-":
		m.appendToExpression("-")
		return m, nil
	case "*":
		m.appendToExpression("*")
		return m, nil
	case "/":
		m.appendToExpression("/")
		return m, nil

	// Special keys
	case ".":
		m.appendToExpression(".")
		return m, nil
	case "=":
		m.calculate()
		return m, nil
	case "c", "C":
		m.clear()
		return m, nil
	case "backspace":
		m.backspace()
		return m, nil
	case "escape":
		m.clear()
		return m, nil
	}

	return m, nil
}

// handleMouseEvent processes mouse input
func (m Model) handleMouseEvent(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if !m.mouseEnabled {
		return m, nil
	}

	switch msg.Type {
	case tea.MouseLeft:
		// Calculate which button was clicked based on mouse position
		if button := m.getButtonFromMousePos(msg.X, msg.Y); button != nil {
			// Update selection to clicked button
			m.selectedRow = button.Row
			m.selectedCol = button.Col
			
			// Press the button
			m.pressButton(button)
			return m, m.startFeedbackTimer()
		}
	}

	return m, nil
}

// handleTick processes timer ticks for visual feedback
func (m Model) handleTick() (tea.Model, tea.Cmd) {
	if m.feedbackTimer > 0 {
		m.feedbackTimer--
		if m.feedbackTimer == 0 {
			m.pressedButton = nil
		}
		return m, m.startFeedbackTimer()
	}
	return m, nil
}

// getButtonFromMousePos determines which button was clicked based on mouse coordinates
func (m Model) getButtonFromMousePos(x, y int) *Button {
	// This is a simplified implementation
	// In a real implementation, you'd need to calculate exact positions
	// based on the rendered layout and terminal coordinates
	
	// Calculate approximate button positions
	// The layout starts after title (2 lines) + display (3 lines) + padding
	startY := 6 // Approximate start of button grid
	
	if y < startY || y >= startY+15 { // 5 rows * 3 height each
		return nil
	}
	
	// Calculate row (each button is 3 lines tall including margin)
	row := (y - startY) / 3
	if row < 0 || row > 4 {
		return nil
	}
	
	// Calculate column (each button is about 7 characters wide including margin)
	// Account for container padding (2 characters from left)
	adjustedX := x - 2
	col := adjustedX / 7
	if col < 0 || col > 3 {
		return nil
	}
	
	// Handle special case for wide "0" button
	if row == 4 && col >= 0 && col <= 1 {
		// Check if clicking on the "0" button area
		if button := m.getButtonAt(4, 0); button != nil && button.Label == "0" {
			return button
		}
	}
	
	return m.getButtonAt(row, col)
}

// backspace removes the last character from the expression
func (m *Model) backspace() {
	if m.hasResult {
		// If showing result, clear it and start fresh
		m.clear()
		return
	}
	
	if len(m.expression) > 0 {
		m.expression = m.expression[:len(m.expression)-1]
		if m.expression == "" {
			m.display = "0"
		} else {
			m.display = m.expression
		}
		m.hasError = false
	}
}

// startFeedbackTimer starts the visual feedback timer
func (m Model) startFeedbackTimer() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// SetMouseEnabled enables or disables mouse support
func (m *Model) SetMouseEnabled(enabled bool) {
	m.mouseEnabled = enabled
}

// GetDisplayValue returns the current display value
func (m Model) GetDisplayValue() string {
	return m.display
}

// GetExpression returns the current expression
func (m Model) GetExpression() string {
	return m.expression
}

// HasError returns whether there's currently an error
func (m Model) HasError() bool {
	return m.hasError
}

// GetErrorMessage returns the current error message
func (m Model) GetErrorMessage() string {
	return m.errorMsg
}

// HasResult returns whether there's a calculated result
func (m Model) HasResult() bool {
	return m.hasResult
}

// GetResult returns the calculated result
func (m Model) GetResult() float64 {
	return m.result
}