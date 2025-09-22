package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Styles for the calculator UI
var (
	// Display styles
	displayStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1).
			Width(30).
			Align(lipgloss.Right).
			Bold(true)

	errorDisplayStyle = displayStyle.Copy().
				Foreground(lipgloss.Color("196")).
				BorderForeground(lipgloss.Color("196"))

	// Button styles
	buttonStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1).
			Width(6).
			Height(2).
			Align(lipgloss.Center).
			MarginRight(1).
			MarginBottom(1)

	selectedButtonStyle = buttonStyle.Copy().
				BorderForeground(lipgloss.Color("86")).
				Foreground(lipgloss.Color("86")).
				Bold(true)

	pressedButtonStyle = buttonStyle.Copy().
				BorderForeground(lipgloss.Color("196")).
				Foreground(lipgloss.Color("196")).
				Bold(true).
				Background(lipgloss.Color("52"))

	operatorButtonStyle = buttonStyle.Copy().
				Background(lipgloss.Color("94")).
				Foreground(lipgloss.Color("15")).
				Bold(true)

	selectedOperatorButtonStyle = operatorButtonStyle.Copy().
					BorderForeground(lipgloss.Color("86")).
					Background(lipgloss.Color("100"))

	pressedOperatorButtonStyle = operatorButtonStyle.Copy().
					BorderForeground(lipgloss.Color("196")).
					Background(lipgloss.Color("160"))

	numberButtonStyle = buttonStyle.Copy().
				Background(lipgloss.Color("236")).
				Foreground(lipgloss.Color("15"))

	selectedNumberButtonStyle = numberButtonStyle.Copy().
					BorderForeground(lipgloss.Color("86")).
					Background(lipgloss.Color("240"))

	pressedNumberButtonStyle = numberButtonStyle.Copy().
					BorderForeground(lipgloss.Color("196")).
					Background(lipgloss.Color("52"))

	// Special button styles
	clearButtonStyle = buttonStyle.Copy().
				Background(lipgloss.Color("160")).
				Foreground(lipgloss.Color("15")).
				Bold(true)

	selectedClearButtonStyle = clearButtonStyle.Copy().
					BorderForeground(lipgloss.Color("86")).
					Background(lipgloss.Color("196"))

	pressedClearButtonStyle = clearButtonStyle.Copy().
					BorderForeground(lipgloss.Color("196")).
					Background(lipgloss.Color("124"))

	equalsButtonStyle = buttonStyle.Copy().
				Background(lipgloss.Color("33")).
				Foreground(lipgloss.Color("15")).
				Bold(true)

	selectedEqualsButtonStyle = equalsButtonStyle.Copy().
					BorderForeground(lipgloss.Color("86")).
					Background(lipgloss.Color("39"))

	pressedEqualsButtonStyle = equalsButtonStyle.Copy().
					BorderForeground(lipgloss.Color("196")).
					Background(lipgloss.Color("27"))

	// Wide button style for "0"
	wideButtonStyle = buttonStyle.Copy().
			Width(14) // Double width plus margin

	selectedWideButtonStyle = selectedButtonStyle.Copy().
				Width(14)

	pressedWideButtonStyle = pressedButtonStyle.Copy().
				Width(14)

	// Container styles
	containerStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true).
			Align(lipgloss.Center).
			Width(30)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			Align(lipgloss.Center).
			Width(30).
			MarginTop(1)
)

// View renders the calculator UI
func (m Model) View() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("🔊 Acoustic Calculator"))
	b.WriteString("\n\n")

	// Display
	displayText := m.display
	if m.hasError {
		b.WriteString(errorDisplayStyle.Render(displayText))
	} else {
		b.WriteString(displayStyle.Render(displayText))
	}
	b.WriteString("\n\n")

	// Button grid
	b.WriteString(m.renderButtonGrid())

	// Help text
	b.WriteString("\n")
	helpText := "Arrow keys: Navigate • Enter/Space: Press • Mouse: Click • q/Ctrl+C: Quit"
	b.WriteString(helpStyle.Render(helpText))

	// Error message if any
	if m.hasError && m.errorMsg != "" {
		b.WriteString("\n\n")
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Italic(true).
			Align(lipgloss.Center).
			Width(30)
		b.WriteString(errorStyle.Render("Error: "+m.errorMsg))
	}

	return containerStyle.Render(b.String())
}

// renderButtonGrid renders the calculator button grid
func (m Model) renderButtonGrid() string {
	var rows []string

	// Group buttons by row
	buttonRows := make(map[int][]Button)
	for _, button := range m.buttons {
		buttonRows[button.Row] = append(buttonRows[button.Row], button)
	}

	// Render each row
	for row := 0; row < 5; row++ {
		if buttons, exists := buttonRows[row]; exists {
			rowStr := m.renderButtonRow(buttons, row)
			rows = append(rows, rowStr)
		}
	}

	return strings.Join(rows, "\n")
}

// renderButtonRow renders a single row of buttons
func (m Model) renderButtonRow(buttons []Button, row int) string {
	var renderedButtons []string

	// Sort buttons by column
	for col := 0; col < 4; col++ {
		var button *Button
		for i := range buttons {
			if buttons[i].Col == col {
				button = &buttons[i]
				break
			}
		}

		if button != nil {
			renderedButton := m.renderButton(*button)
			renderedButtons = append(renderedButtons, renderedButton)

			// Skip next column if this is a wide button
			if button.Width > 1 {
				col++
			}
		} else {
			// Empty space for missing buttons
			emptyStyle := lipgloss.NewStyle().
				Width(6).
				Height(2).
				MarginRight(1).
				MarginBottom(1)
			renderedButtons = append(renderedButtons, emptyStyle.Render(""))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedButtons...)
}

// renderButton renders a single button with appropriate styling
func (m Model) renderButton(button Button) string {
	isSelected := m.selectedRow == button.Row && m.selectedCol == button.Col
	isPressed := m.pressedButton != nil && 
		m.pressedButton.Row == button.Row && 
		m.pressedButton.Col == button.Col && 
		m.feedbackTimer > 0

	// Choose base style based on button type
	var style lipgloss.Style
	switch button.Value {
	case "C":
		if isPressed {
			style = pressedClearButtonStyle
		} else if isSelected {
			style = selectedClearButtonStyle
		} else {
			style = clearButtonStyle
		}
	case "=":
		if isPressed {
			style = pressedEqualsButtonStyle
		} else if isSelected {
			style = selectedEqualsButtonStyle
		} else {
			style = equalsButtonStyle
		}
	case "+", "-", "*", "/", "±", "%":
		if isPressed {
			style = pressedOperatorButtonStyle
		} else if isSelected {
			style = selectedOperatorButtonStyle
		} else {
			style = operatorButtonStyle
		}
	default:
		// Number buttons
		if button.Width > 1 {
			// Wide button (like "0")
			if isPressed {
				style = pressedWideButtonStyle.Copy().Background(lipgloss.Color("52"))
			} else if isSelected {
				style = selectedWideButtonStyle.Copy().Background(lipgloss.Color("240"))
			} else {
				style = wideButtonStyle.Copy().Background(lipgloss.Color("236"))
			}
		} else {
			if isPressed {
				style = pressedNumberButtonStyle
			} else if isSelected {
				style = selectedNumberButtonStyle
			} else {
				style = numberButtonStyle
			}
		}
	}

	return style.Render(button.Label)
}

// GetButtonDimensions returns the dimensions needed for button layout
func (m Model) GetButtonDimensions() (width, height int) {
	// Calculate required width: 4 buttons * 6 width + 3 margins + container padding
	width = 4*6 + 3*1 + 4 // 4 padding (2 on each side)
	
	// Calculate required height: title + display + 5 button rows + help + container padding
	height = 1 + 2 + 3 + 2 + 5*3 + 1 + 2 + 4 // Approximate height
	
	return width, height
}