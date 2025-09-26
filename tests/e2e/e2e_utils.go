//go:build !lint
// +build !lint

package e2e

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
)

// E2ETestConfig defines configuration for E2E tests
type E2ETestConfig struct {
	AppPath          string
	TestName         string
	OutputDir        string
	Timeout          time.Duration
	ScreenshotOnPass bool
	ScreenshotOnFail bool
	RecordTerminal   bool
	Platform         string
}

// E2ETestRunner handles end-to-end test execution
type E2ETestRunner struct {
	config      *E2ETestConfig
	appProcess  *os.Process
	terminalPID int
	recorder    *TerminalRecorder
	mu          sync.RWMutex
	startTime   time.Time
	events      []E2ETestEvent
	observers   []E2ETestObserver
}

// E2ETestEvent represents an event during E2E test execution
type E2ETestEvent struct {
	Type        string      `json:"type"`
	Timestamp   time.Time   `json:"timestamp"`
	Description string      `json:"description"`
	Screenshot  string      `json:"screenshot,omitempty"`
	Recording   string      `json:"recording,omitempty"`
	Metadata    interface{} `json:"metadata,omitempty"`
	Error       error       `json:"error,omitempty"`
}

// E2ETestObserver interface for observing E2E test events
type E2ETestObserver interface {
	OnEvent(event E2ETestEvent)
	OnScreenshot(path string, eventType string)
	OnRecording(path string, eventType string)
	OnTestComplete(runner *E2ETestRunner)
}

// TerminalRecorder handles cross-platform terminal recording
type TerminalRecorder struct {
	Active     bool
	OutputPath string
	Process    *os.Process
	Platform   string
}

// NewE2ETestRunner creates a new E2E test runner
func NewE2ETestRunner(config *E2ETestConfig) *E2ETestRunner {
	if config.Platform == "" {
		config.Platform = runtime.GOOS
	}

	if config.OutputDir == "" {
		config.OutputDir = filepath.Join("tests", "artifacts", "e2e")
	}

	if config.Timeout == 0 {
		config.Timeout = 5 * time.Minute
	}

	return &E2ETestRunner{
		config:    config,
		startTime: time.Now(),
		events:    make([]E2ETestEvent, 0),
		observers: make([]E2ETestObserver, 0),
	}
}

// AddObserver adds an observer to the E2E test runner
func (r *E2ETestRunner) AddObserver(observer E2ETestObserver) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.observers = append(r.observers, observer)
}

// StartApplication starts the application under test
func (r *E2ETestRunner) StartApplication() error {
	r.logEvent("app_start", "Starting application", nil)

	// Ensure output directory exists
	if err := os.MkdirAll(r.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Start terminal recording if enabled
	if r.config.RecordTerminal {
		if err := r.startTerminalRecording(); err != nil {
			r.logEvent("recording_error", "Failed to start terminal recording", map[string]interface{}{"error": err.Error()})
		}
	}

	// Build the application first
	buildCmd := exec.Command("go", "build", "-o", "acousticalc-test", "./cmd/acousticalc")
	buildCmd.Dir = "."
	if output, err := buildCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to build application: %w, output: %s", err, string(output))
	}

	// Start the application
	appPath := "./acousticalc-test"
	if r.config.AppPath != "" {
		appPath = r.config.AppPath
	}

	cmd := exec.Command(appPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}

	r.appProcess = cmd.Process
	r.logEvent("app_started", "Application started successfully", map[string]interface{}{"pid": cmd.Process.Pid})

	// Wait a moment for app to initialize
	time.Sleep(2 * time.Second)

	return nil
}

// StopApplication stops the application under test
func (r *E2ETestRunner) StopApplication() error {
	r.logEvent("app_stop", "Stopping application", nil)

	if r.appProcess != nil {
		if err := r.appProcess.Kill(); err != nil {
			return fmt.Errorf("failed to kill application process: %w", err)
		}
		r.appProcess = nil
	}

	// Stop terminal recording
	if r.recorder != nil && r.recorder.Active {
		r.stopTerminalRecording()
	}

	// Clean up test binary
	os.Remove("./acousticalc-test")

	r.logEvent("app_stopped", "Application stopped", nil)
	return nil
}

// startTerminalRecording starts terminal recording based on platform
func (r *E2ETestRunner) startTerminalRecording() error {
	timestamp := time.Now().Format("20060102_150405")
	recordingPath := filepath.Join(r.config.OutputDir, fmt.Sprintf("%s_terminal_%s.cast", r.config.TestName, timestamp))

	r.recorder = &TerminalRecorder{
		OutputPath: recordingPath,
		Platform:   r.config.Platform,
	}

	switch r.config.Platform {
	case "linux", "darwin":
		// Use asciinema for recording
		if _, err := exec.LookPath("asciinema"); err == nil {
			cmd := exec.Command("asciinema", "rec", "-c", "bash", recordingPath)
			if err := cmd.Start(); err != nil {
				return fmt.Errorf("failed to start asciinema: %w", err)
			}
			r.recorder.Process = cmd.Process
			r.recorder.Active = true
			r.logEvent("recording_started", "Terminal recording started", map[string]interface{}{"path": recordingPath})
		} else {
			return fmt.Errorf("asciinema not found for terminal recording")
		}
	case "windows":
		// Windows terminal recording would be implemented here
		return fmt.Errorf("Windows terminal recording not yet implemented")
	default:
		return fmt.Errorf("unsupported platform for terminal recording: %s", r.config.Platform)
	}

	return nil
}

// stopTerminalRecording stops the terminal recording
func (r *E2ETestRunner) stopTerminalRecording() {
	if r.recorder == nil || !r.recorder.Active {
		return
	}

	if r.recorder.Process != nil {
		r.recorder.Process.Signal(os.Interrupt)
		r.recorder.Process.Wait()
		r.recorder.Active = false
		r.logEvent("recording_stopped", "Terminal recording stopped", map[string]interface{}{"path": r.recorder.OutputPath})
	}
}

// SimulateInput simulates keyboard input
func (r *E2ETestRunner) SimulateInput(input string) error {
	r.logEvent("input_simulate", "Simulating keyboard input", map[string]interface{}{"input": input})

	// Use robotgo for cross-platform input simulation
	for _, char := range input {
		key := string(char)
		switch key {
		case " ":
			key = "space"
		case "\n":
			key = "enter"
		case "\t":
			key = "tab"
		}

		robotgo.KeyTap(key)
		time.Sleep(50 * time.Millisecond) // Small delay between keystrokes
	}

	return nil
}

// SimulateMouseClick simulates a mouse click at specified coordinates
func (r *E2ETestRunner) SimulateMouseClick(x, y int) error {
	r.logEvent("mouse_click", "Simulating mouse click", map[string]interface{}{"x": x, "y": y})

	robotgo.MoveMouse(x, y)

	time.Sleep(100 * time.Millisecond)
	robotgo.MouseClick("left")
	return nil
}

// CaptureScreenshot captures a screenshot for the current state
func (r *E2ETestRunner) CaptureScreenshot(eventType string) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s_%s.png", r.config.TestName, eventType, timestamp)
	screenshotPath := filepath.Join(r.config.OutputDir, filename)

	// Use robotgo for cross-platform screenshot capture
	bitmap := robotgo.CaptureScreen()
	if bitmap == nil {
		return "", fmt.Errorf("failed to capture screen")
	}
	defer robotgo.FreeBitmap(bitmap)

	// Convert CBitmap to image.Image before saving
	img := robotgo.ToImage(bitmap)
	if err := robotgo.SavePng(img, screenshotPath); err != nil {
		return "", fmt.Errorf("failed to save screenshot: %w", err)
	}

	r.logEvent("screenshot_captured", "Screenshot captured", map[string]interface{}{
		"path":       screenshotPath,
		"event_type": eventType,
	})

	r.notifyScreenshotObservers(screenshotPath, eventType)
	return screenshotPath, nil
}

// WaitFor waits for a condition to be true within timeout
func (r *E2ETestRunner) WaitFor(condition func() bool, description string) error {
	r.logEvent("wait_start", "Waiting for condition", map[string]interface{}{"condition": description})

	ctx, cancel := context.WithTimeout(context.Background(), r.config.Timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			err := fmt.Errorf("timeout waiting for condition: %s", description)
			r.logEvent("wait_timeout", "Wait timeout", map[string]interface{}{"condition": description, "error": err.Error()})
			return err
		case <-ticker.C:
			if condition() {
				r.logEvent("wait_complete", "Condition met", map[string]interface{}{"condition": description})
				return nil
			}
		}
	}
}

// logEvent logs an E2E test event
func (r *E2ETestRunner) logEvent(eventType, description string, metadata interface{}) {
	event := E2ETestEvent{
		Type:        eventType,
		Timestamp:   time.Now(),
		Description: description,
		Metadata:    metadata,
	}

	r.mu.Lock()
	r.events = append(r.events, event)
	r.mu.Unlock()

	r.notifyObservers(event)
}

// notifyObservers notifies all observers of an event
func (r *E2ETestRunner) notifyObservers(event E2ETestEvent) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, observer := range r.observers {
		observer.OnEvent(event)
	}
}

// notifyScreenshotObservers notifies all observers of a screenshot capture
func (r *E2ETestRunner) notifyScreenshotObservers(path, eventType string) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, observer := range r.observers {
		observer.OnScreenshot(path, eventType)
	}
}

// notifyRecordingObservers notifies all observers of a recording
func (r *E2ETestRunner) notifyRecordingObservers(path, eventType string) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, observer := range r.observers {
		observer.OnRecording(path, eventType)
	}
}

// GetEvents returns all logged events
func (r *E2ETestRunner) GetEvents() []E2ETestEvent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy to avoid race conditions
	events := make([]E2ETestEvent, len(r.events))
	copy(events, r.events)
	return events
}

// GenerateReport generates an E2E test report
func (r *E2ETestRunner) GenerateReport() error {
	reportPath := filepath.Join(r.config.OutputDir, fmt.Sprintf("%s_e2e_report.html", r.config.TestName))

	html := r.generateHTMLReport()

	if err := os.WriteFile(reportPath, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write E2E report: %w", err)
	}

	r.logEvent("report_generated", "E2E test report generated", map[string]interface{}{"path": reportPath})
	return nil
}

// generateHTMLReport generates HTML report content
func (r *E2ETestRunner) generateHTMLReport() string {
	duration := time.Since(r.startTime)

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>E2E Test Report: %s</title>
    <style>
        body { font-family: 'Monaco', 'Menlo', monospace; margin: 40px; background: #1e1e1e; color: #d4d4d4; }
        .header { border-bottom: 2px solid #32cd32; padding: 20px 0; margin-bottom: 30px; }
        .event { margin: 20px 0; padding: 15px; border-left: 3px solid #32cd32; background: #2d2d2d; }
        .screenshot { max-width: 400px; border: 1px solid #444; margin: 10px 0; }
        .timestamp { color: #808080; font-size: 12px; }
        .event-type { color: #32cd32; font-weight: bold; }
        .metadata { background: #3c3c3c; padding: 10px; margin: 10px 0; font-size: 12px; }
        .error { color: #ff6b6b; }
        .success { color: #51cf66; }
    </style>
</head>
<body>
    <div class="header">
        <h1>E2E Test Report</h1>
        <h2>Test: %s</h2>
        <p>Platform: %s</p>
        <p>Duration: %s</p>
        <p>Total Events: %d</p>
        <p>Terminal Recording: %s</p>
    </div>`,
		r.config.TestName, r.config.TestName, r.config.Platform, duration, len(r.events),
		func() string {
			if r.recorder != nil && r.recorder.OutputPath != "" {
				return fmt.Sprintf(`<a href="%s">View Recording</a>`, filepath.Base(r.recorder.OutputPath))
			}
			return "Not available"
		}())

	for _, event := range r.events {
		eventClass := "event"
		if event.Error != nil {
			eventClass += " error"
		} else if strings.Contains(event.Type, "complete") || strings.Contains(event.Type, "success") {
			eventClass += " success"
		}

		html += fmt.Sprintf(`
    <div class="%s">
        <div class="timestamp">%s</div>
        <div class="event-type">%s</div>
        <div>%s</div>`,
			eventClass, event.Timestamp.Format("15:04:05.000"), event.Type, event.Description)

		if event.Screenshot != "" {
			html += fmt.Sprintf(`<img src="%s" class="screenshot" alt="Screenshot for %s">`, filepath.Base(event.Screenshot), event.Type)
		}

		if event.Recording != "" {
			html += fmt.Sprintf(`<p><a href="%s">View Recording</a></p>`, filepath.Base(event.Recording))
		}

		if event.Metadata != nil {
			html += `<div class="metadata"><strong>Metadata:</strong><br>`
			if meta, ok := event.Metadata.(map[string]interface{}); ok {
				for key, value := range meta {
					html += fmt.Sprintf("%s: %v<br>", key, value)
				}
			}
			html += `</div>`
		}

		if event.Error != nil {
			html += fmt.Sprintf(`<div class="error">Error: %v</div>`, event.Error)
		}

		html += `</div>`
	}

	html += `
</body>
</html>`

	return html
}

// Complete completes the E2E test and notifies observers
func (r *E2ETestRunner) Complete() {
	r.logEvent("test_complete", "E2E test completed", map[string]interface{}{
		"duration": time.Since(r.startTime),
		"events":   len(r.events),
	})

	// Notify completion observers
	r.mu.RLock()
	for _, observer := range r.observers {
		observer.OnTestComplete(r)
	}
	r.mu.RUnlock()
}
