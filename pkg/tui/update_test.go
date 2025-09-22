package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestHandleKeyPress(t *testing.T) {
	model := NewModel()

	// Test quit keys
	quitKeys := []string{"ctrl+c", "q"}
	for _, key := range quitKeys {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
		if key == "ctrl+c" {
			msg = tea.KeyMsg{Type: tea.KeyCtrlC}
		}
		
		newModel, cmd := model.handleKeyPress(msg)
		if cmd == nil {
			t.Errorf("Expected quit command for key %s", key)
		}
		_ = newModel // Use the returned model
	}

	// Test navigation keys
	navigationTests := []struct {
		key         tea.KeyType
		expectedCmd bool
	}{
		{tea.KeyUp, false},
		{tea.KeyDown, false},
		{tea.KeyLeft, false},
		{tea.KeyRight, false},
	}

	for _, test := range navigationTests {
		msg := tea.KeyMsg{Type: test.key}
		newModel, cmd := model.handleKeyPress(msg)
		
		if (cmd != nil) != test.expectedCmd {
			t.Errorf("Navigation key %v: expected cmd=%v, got cmd=%v", test.key, test.expectedCmd, cmd != nil)
		}
		
		// Verify model is returned (even if unchanged)
		if newModel == nil {
			t.Errorf("Expected model to be returned for navigation key %v", test.key)
		}
	}

	// Test enter and space
	actionKeys := []tea.KeyType{tea.KeyEnter, tea.KeySpace}
	for _, keyType := range actionKeys {
		msg := tea.KeyMsg{Type: keyType}
		newModel, cmd := model.handleKeyPress(msg)
		
		// Should return model and potentially a command
		if newModel == nil {
			t.Errorf("Expected model to be returned for action key %v", keyType)
		}
		_ = cmd // Command may or may not be present depending on button selection
	}

	// Test number input
	numberKeys := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	for _, num := range numberKeys {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(num)}
		newModel, cmd := model.handleKeyPress(msg)
		
		if newModel == nil {
			t.Errorf("Expected model to be returned for number key %s", num)
		}
		if cmd != nil {
			t.Errorf("Expected no command for number key %s", num)
		}
		
		// Check that the number was added to expression
		typedModel := newModel.(Model)
		if typedModel.expression != num {
			t.Errorf("Expected expression to be %s, got %s", num, typedModel.expression)
		}
		
		// Reset for next test
		model = NewModel()
	}

	// Test operator input
	operatorTests := []struct {
		key      string
		operator string
	}{
		{"+", "+"},
		{"-", "-"},
		{"*", "*"},
		{"/", "/"},
	}

	for _, test := range operatorTests {
		model = NewModel()
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(test.key)}
		newModel, cmd := model.handleKeyPress(msg)
		
		if newModel == nil {
			t.Errorf("Expected model to be returned for operator key %s", test.key)
		}
		if cmd != nil {
			t.Errorf("Expected no command for operator key %s", test.key)
		}
		
		typedModel := newModel.(Model)
		if typedModel.expression != test.operator {
			t.Errorf("Expected expression to be %s, got %s", test.operator, typedModel.expression)
		}
	}

	// Test special keys
	specialKeyTests := []struct {
		key         string
		keyType     tea.KeyType
		expectCmd   bool
		description string
	}{
		{".", tea.KeyRunes, false, "decimal point"},
		{"=", tea.KeyRunes, false, "equals"},
		{"c", tea.KeyRunes, false, "clear lowercase"},
		{"C", tea.KeyRunes, false, "clear uppercase"},
		{"", tea.KeyBackspace, false, "backspace"},
		{"", tea.KeyEsc, false, "escape"},
	}

	for _, test := range specialKeyTests {
		model = NewModel()
		var msg tea.KeyMsg
		if test.keyType == tea.KeyRunes {
			msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(test.key)}
		} else {
			msg = tea.KeyMsg{Type: test.keyType}
		}
		
		newModel, cmd := model.handleKeyPress(msg)
		
		if newModel == nil {
			t.Errorf("Expected model to be returned for %s", test.description)
		}
		if (cmd != nil) != test.expectCmd {
			t.Errorf("%s: expected cmd=%v, got cmd=%v", test.description, test.expectCmd, cmd != nil)
		}
	}
}

func TestHandleMouseEvent(t *testing.T) {
	model := NewModel()

	// Test mouse click
	mouseMsg := tea.MouseMsg{
		X:    10,
		Y:    10,
		Type: tea.MouseLeft,
	}

	newModel, cmd := model.handleMouseEvent(mouseMsg)
	
	if newModel == nil {
		t.Error("Expected model to be returned for mouse event")
	}
	
	// Command may or may not be present depending on whether a button was clicked
	_ = cmd

	// Test with mouse disabled
	model.SetMouseEnabled(false)
	newModel, cmd = model.handleMouseEvent(mouseMsg)
	
	if newModel == nil {
		t.Error("Expected model to be returned even with mouse disabled")
	}
	if cmd != nil {
		t.Error("Expected no command when mouse is disabled")
	}
}

func TestHandleTick(t *testing.T) {
	model := NewModel()

	// Test with feedback timer > 0
	model.feedbackTimer = 2
	model.pressedButton = &Button{Label: "1", Value: "1"}

	newModel, cmd := model.handleTick()
	
	if newModel == nil {
		t.Error("Expected model to be returned for tick")
	}
	
	typedModel := newModel.(Model)
	if typedModel.feedbackTimer != 1 {
		t.Errorf("Expected feedbackTimer to be decremented to 1, got %d", typedModel.feedbackTimer)
	}
	if typedModel.pressedButton == nil {
		t.Error("Expected pressedButton to still be set")
	}
	if cmd == nil {
		t.Error("Expected command to continue timer")
	}

	// Test with feedback timer = 1 (should clear pressed button)
	model.feedbackTimer = 1
	model.pressedButton = &Button{Label: "1", Value: "1"}

	newModel, cmd = model.handleTick()
	typedModel = newModel.(Model)
	
	if typedModel.feedbackTimer != 0 {
		t.Errorf("Expected feedbackTimer to be 0, got %d", typedModel.feedbackTimer)
	}
	if typedModel.pressedButton != nil {
		t.Error("Expected pressedButton to be cleared")
	}

	// Test with feedback timer = 0
	model.feedbackTimer = 0
	model.pressedButton = nil

	newModel, cmd = model.handleTick()
	typedModel = newModel.(Model)
	
	if typedModel.feedbackTimer != 0 {
		t.Errorf("Expected feedbackTimer to remain 0, got %d", typedModel.feedbackTimer)
	}
	if cmd != nil {
		t.Error("Expected no command when timer is 0")
	}
}

func TestUpdate(t *testing.T) {
	model := NewModel()

	// Test KeyMsg
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("1")}
	newModel, cmd := model.Update(keyMsg)
	
	if newModel == nil {
		t.Error("Expected model to be returned for KeyMsg")
	}
	_ = cmd

	// Test MouseMsg
	mouseMsg := tea.MouseMsg{X: 10, Y: 10, Type: tea.MouseLeft}
	newModel, cmd = model.Update(mouseMsg)
	
	if newModel == nil {
		t.Error("Expected model to be returned for MouseMsg")
	}
	_ = cmd

	// Test WindowSizeMsg
	windowMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	newModel, cmd = model.Update(windowMsg)
	
	if newModel == nil {
		t.Error("Expected model to be returned for WindowSizeMsg")
	}
	if cmd != nil {
		t.Error("Expected no command for WindowSizeMsg")
	}
	
	typedModel := newModel.(Model)
	if typedModel.width != 80 || typedModel.height != 24 {
		t.Errorf("Expected dimensions to be set to 80x24, got %dx%d", typedModel.width, typedModel.height)
	}

	// Test tickMsg
	tickMsg := tickMsg(time.Now())
	newModel, cmd = model.Update(tickMsg)
	
	if newModel == nil {
		t.Error("Expected model to be returned for tickMsg")
	}
	_ = cmd

	// Test unknown message type
	unknownMsg := "unknown"
	newModel, cmd = model.Update(unknownMsg)
	
	if newModel == nil {
		t.Error("Expected model to be returned for unknown message")
	}
	if cmd != nil {
		t.Error("Expected no command for unknown message")
	}
}

func TestGetButtonFromMousePos(t *testing.T) {
	model := NewModel()

	// Test coordinates outside button area
	button := model.getButtonFromMousePos(0, 0)
	if button != nil {
		t.Error("Expected no button at coordinates (0,0)")
	}

	button = model.getButtonFromMousePos(100, 100)
	if button != nil {
		t.Error("Expected no button at coordinates (100,100)")
	}

	// Test specific button coordinates based on the corrected layout
	testCases := []struct {
		x, y          int
		expectedLabel string
		description   string
	}{
		// Row 0: C, ±, %, ÷
		{5, 8, "C", "Button C (row 0, col 0)"},
		{14, 8, "±", "Button ± (row 0, col 1)"},
		{23, 8, "%", "Button % (row 0, col 2)"},
		{32, 8, "÷", "Button ÷ (row 0, col 3)"},
		
		// Row 1: 7, 8, 9, ×
		{5, 11, "7", "Button 7 (row 1, col 0)"},
		{14, 11, "8", "Button 8 (row 1, col 1)"},
		{23, 11, "9", "Button 9 (row 1, col 2)"},
		{32, 11, "×", "Button × (row 1, col 3)"},
		
		// Row 2: 4, 5, 6, -
		{5, 14, "4", "Button 4 (row 2, col 0)"},
		{14, 14, "5", "Button 5 (row 2, col 1)"},
		{23, 14, "6", "Button 6 (row 2, col 2)"},
		{32, 14, "-", "Button - (row 2, col 3)"},
		
		// Row 3: 1, 2, 3, +
		{5, 17, "1", "Button 1 (row 3, col 0)"},
		{14, 17, "2", "Button 2 (row 3, col 1)"},
		{23, 17, "3", "Button 3 (row 3, col 2)"},
		{32, 17, "+", "Button + (row 3, col 3)"},
		
		// Row 4: 0 (wide), ., =
		{5, 20, "0", "Button 0 (row 4, col 0) - left side"},
		{14, 20, "0", "Button 0 (row 4, col 0) - right side"},
		{23, 20, ".", "Button . (row 4, col 2)"},
		{32, 20, "=", "Button = (row 4, col 3)"},
	}

	for _, tc := range testCases {
		button := model.getButtonFromMousePos(tc.x, tc.y)
		if button == nil {
			t.Errorf("%s: Expected button at (%d, %d), got nil", tc.description, tc.x, tc.y)
			continue
		}
		if button.Label != tc.expectedLabel {
			t.Errorf("%s: Expected button label '%s', got '%s'", tc.description, tc.expectedLabel, button.Label)
		}
	}

	// Test coordinates between buttons (should return nil)
	button = model.getButtonFromMousePos(10, 8) // Between columns 0 and 1 (gap at x=10)
	if button != nil {
		t.Errorf("Expected no button between columns at (10,8), got button '%s'", button.Label)
	}
	
	button = model.getButtonFromMousePos(5, 7) // Above button area (y=7, buttons start at y=8)
	if button != nil {
		t.Errorf("Expected no button above button area at (5,7), got button '%s'", button.Label)
	}
	
	button = model.getButtonFromMousePos(1, 8) // Left of button area (x=1, buttons start at x=2)
	if button != nil {
		t.Errorf("Expected no button left of button area at (1,8), got button '%s'", button.Label)
	}
}

func TestBackspace(t *testing.T) {
	model := NewModel()

	// Test backspace with result
	model.hasResult = true
	model.display = "42"
	model.expression = "42"
	model.result = 42

	model.backspace()

	if model.display != "0" {
		t.Errorf("Expected display to be '0' after backspace on result, got '%s'", model.display)
	}
	if model.expression != "" {
		t.Errorf("Expected expression to be empty after backspace on result, got '%s'", model.expression)
	}
	if model.hasResult {
		t.Error("Expected hasResult to be false after backspace on result")
	}

	// Test backspace with expression
	model.expression = "123"
	model.display = "123"
	model.hasResult = false

	model.backspace()

	if model.expression != "12" {
		t.Errorf("Expected expression to be '12' after backspace, got '%s'", model.expression)
	}
	if model.display != "12" {
		t.Errorf("Expected display to be '12' after backspace, got '%s'", model.display)
	}

	// Test backspace on single character
	model.expression = "1"
	model.display = "1"

	model.backspace()

	if model.expression != "" {
		t.Errorf("Expected expression to be empty after backspace on single char, got '%s'", model.expression)
	}
	if model.display != "0" {
		t.Errorf("Expected display to be '0' after backspace on single char, got '%s'", model.display)
	}

	// Test backspace on empty expression
	model.expression = ""
	model.display = "0"

	model.backspace()

	if model.expression != "" {
		t.Errorf("Expected expression to remain empty after backspace on empty, got '%s'", model.expression)
	}
	if model.display != "0" {
		t.Errorf("Expected display to remain '0' after backspace on empty, got '%s'", model.display)
	}
}

func TestStartFeedbackTimer(t *testing.T) {
	model := NewModel()

	cmd := model.startFeedbackTimer()
	
	if cmd == nil {
		t.Error("Expected command to be returned from startFeedbackTimer")
	}

	// We can't easily test the actual timer functionality without running the command,
	// but we can verify that a command is returned
}