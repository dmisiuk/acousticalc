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
	// Calculate exact button positions based on the rendered layout
	// Layout structure:
	// - Container border: 1 line at top
	// - Title: 2 lines (title + blank line)  
	// - Display: 3 lines (blank + display + blank)
	// - Button grid starts at row 8 (1 + 2 + 3 + 2 for spacing)
	
	startY := 8 // First button row starts at line 8
	
	// Check if click is within button grid area (5 rows * 3 lines each = 15 lines)
	if y < startY || y >= startY+15 {
		return nil
	}
	
	// Calculate button row based on exact y positions
	// Buttons are at specific y positions: 8-10, 11-13, 14-16, 17-19, 20-22
	// Each button is 3 rows tall (including borders)
	var buttonRow int = -1
	
	if y >= 8 && y <= 10 {
		buttonRow = 0
	} else if y >= 11 && y <= 13 {
		buttonRow = 1
	} else if y >= 14 && y <= 16 {
		buttonRow = 2
	} else if y >= 17 && y <= 19 {
		buttonRow = 3
	} else if y >= 20 && y <= 22 {
		buttonRow = 4
	}
	
	if buttonRow == -1 {
		return nil // Click was not on a button row
	}
	
	// Calculate button column based on exact positions
	// Container has left padding of 2, then buttons start
	// Button positions: 2-9, 11-18, 20-27, 29-36 (each button is 8 chars wide)
	// There's a 1-character gap between buttons at positions 10, 19, 28
	
	var buttonCol int = -1
	
	// Check which button column was clicked (using absolute coordinates)
	if x >= 2 && x <= 9 {
		buttonCol = 0
	} else if x >= 11 && x <= 18 {
		buttonCol = 1
	} else if x >= 20 && x <= 27 {
		buttonCol = 2
	} else if x >= 29 && x <= 36 {
		buttonCol = 3
	}
	
	if buttonCol == -1 {
		return nil // Click was not on a button
	}
	
	// Handle special case for wide "0" button in row 4
	if buttonRow == 4 && buttonCol <= 1 {
		// The "0" button spans columns 0-1, so return it for both
		if button := m.getButtonAt(4, 0); button != nil && button.Label == "0" {
			return button
		}
	}
	
	return m.getButtonAt(buttonRow, buttonCol)
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