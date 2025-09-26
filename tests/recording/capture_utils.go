package recording

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// RecordingConfig defines configuration for terminal recording
type RecordingConfig struct {
	OutputDir   string
	TestName    string
	Format      string // "cast" (asciinema), "gif", "mp4"
	Quality     string
	ShowInput   bool
	ShowTiming  bool
	MaxDuration time.Duration
	Platform    string
}

// TerminalRecorder handles cross-platform terminal recording
type TerminalRecorder struct {
	config     *RecordingConfig
	process    *os.Process
	outputFile string
	startTime  time.Time
	active     bool
	mu         sync.RWMutex
	eventLog   []RecordingEvent
}

// RecordingEvent represents an event during recording
type RecordingEvent struct {
	Timestamp   time.Time              `json:"timestamp"`
	Type        string                 `json:"type"` // "start", "stop", "input", "output", "error"
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// InputVisualizer handles input visualization for recordings
type InputVisualizer struct {
	enabled     bool
	inputEvents []InputEvent
	sync.RWMutex
}

// InputEvent represents a user input event
type InputEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // "keyboard", "mouse", "clipboard"
	Value     string    `json:"value"`
	X         int       `json:"x,omitempty"`
	Y         int       `json:"y,omitempty"`
}

// NewTerminalRecorder creates a new terminal recorder
func NewTerminalRecorder(config *RecordingConfig) *TerminalRecorder {
	if config.Platform == "" {
		config.Platform = runtime.GOOS
	}

	if config.OutputDir == "" {
		config.OutputDir = filepath.Join("tests", "artifacts", "recordings")
	}

	if config.Format == "" {
		config.Format = "cast"
	}

	if config.Quality == "" {
		config.Quality = "medium"
	}

	return &TerminalRecorder{
		config:   config,
		eventLog: make([]RecordingEvent, 0),
	}
}

// StartRecording starts the terminal recording
func (tr *TerminalRecorder) StartRecording() error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	if tr.active {
		return fmt.Errorf("recording is already active")
	}

	// Ensure output directory exists
	if err := os.MkdirAll(tr.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate output filename
	timestamp := time.Now().Format("20060102_150405")
	tr.outputFile = filepath.Join(tr.config.OutputDir,
		fmt.Sprintf("%s_terminal_%s.%s", tr.config.TestName, timestamp, tr.config.Format))

	tr.startTime = time.Now()

	// Start recording based on platform and format
	switch tr.config.Platform {
	case "linux", "darwin":
		if tr.config.Format == "cast" {
			if err := tr.startAsciinemaRecording(); err != nil {
				return fmt.Errorf("failed to start asciinema recording: %w", err)
			}
		} else {
			if err := tr.startFFmpegRecording(); err != nil {
				return fmt.Errorf("failed to start ffmpeg recording: %w", err)
			}
		}
	case "windows":
		if err := tr.startWindowsRecording(); err != nil {
			return fmt.Errorf("failed to start windows recording: %w", err)
		}
	default:
		return fmt.Errorf("unsupported platform for recording: %s", tr.config.Platform)
	}

	tr.active = true
	tr.logEvent("start", "Terminal recording started", map[string]interface{}{
		"output_file": tr.outputFile,
		"format":      tr.config.Format,
		"platform":    tr.config.Platform,
	})

	return nil
}

// startAsciinemaRecording starts recording using asciinema
func (tr *TerminalRecorder) startAsciinemaRecording() error {
	// Check if asciinema is available
	if _, err := exec.LookPath("asciinema"); err != nil {
		return fmt.Errorf("asciinema not found: %w", err)
	}

	// Prepare asciinema command
	args := []string{"rec"}

	if tr.config.ShowInput {
		args = append(args, "--stdin")
	}

	if tr.config.Quality == "high" {
		args = append(args, "--cols", "120", "--rows", "40")
	}

	args = append(args, tr.outputFile)

	cmd := exec.Command("asciinema", args...)

	// Start the process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start asciinema: %w", err)
	}

	tr.process = cmd.Process
	return nil
}

// startFFmpegRecording starts recording using ffmpeg
func (tr *TerminalRecorder) startFFmpegRecording() error {
	// Check if ffmpeg is available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found: %w", err)
	}

	// Get terminal dimensions (this is a simplified approach)
	// In a real implementation, you'd get the actual terminal dimensions
	width, height := "800", "600"

	// Prepare ffmpeg command
	args := []string{
		"-f", "x11grab",
		"-r", "30",
		"-s", width + "x" + height,
		"-i", ":0.0",
		"-c:v", "libx264",
		"-preset", tr.config.Quality,
		"-pix_fmt", "yuv420p",
		tr.outputFile,
	}

	cmd := exec.Command("ffmpeg", args...)

	// Start the process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	tr.process = cmd.Process
	return nil
}

// startWindowsRecording starts recording on Windows
func (tr *TerminalRecorder) startWindowsRecording() error {
	// Windows recording implementation
	// This could use PowerShell or other Windows-specific tools
	return fmt.Errorf("Windows recording not yet implemented")
}

// StopRecording stops the terminal recording
func (tr *TerminalRecorder) StopRecording() error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	if !tr.active {
		return fmt.Errorf("recording is not active")
	}

	if tr.process != nil {
		// Try to stop gracefully first
		if err := tr.process.Signal(os.Interrupt); err != nil {
			// If graceful stop fails, kill the process
			if err := tr.process.Kill(); err != nil {
				return fmt.Errorf("failed to kill recording process: %w", err)
			}
		}

		// Wait for process to finish
		if _, err := tr.process.Wait(); err != nil {
			// This is expected if we killed the process
		}

		tr.process = nil
	}

	tr.active = false
	duration := time.Since(tr.startTime)

	tr.logEvent("stop", "Terminal recording stopped", map[string]interface{}{
		"output_file": tr.outputFile,
		"duration":    duration.String(),
	})

	return nil
}

// IsActive returns whether recording is currently active
func (tr *TerminalRecorder) IsActive() bool {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	return tr.active
}

// GetOutputFile returns the path to the output file
func (tr *TerminalRecorder) GetOutputFile() string {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	return tr.outputFile
}

// GetDuration returns the duration of the recording
func (tr *TerminalRecorder) GetDuration() time.Duration {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	if tr.active {
		return time.Since(tr.startTime)
	}
	return time.Duration(0)
}

// LogInput logs an input event for visualization
func (tr *TerminalRecorder) LogInput(inputType, value string, x, y int) {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	_ = InputEvent{
		Timestamp: time.Now(),
		Type:      inputType,
		Value:     value,
		X:         x,
		Y:         y,
	}

	// Note: inputEvents field removed from TerminalRecorder struct
	// If needed, this should be handled by a separate InputVisualizer

	tr.logEvent("input", "User input logged", map[string]interface{}{
		"type":  inputType,
		"value": value,
		"x":     x,
		"y":     y,
	})
}

// logEvent logs a recording event
func (tr *TerminalRecorder) logEvent(eventType, description string, metadata map[string]interface{}) {
	event := RecordingEvent{
		Timestamp:   time.Now(),
		Type:        eventType,
		Description: description,
		Metadata:    metadata,
	}

	tr.eventLog = append(tr.eventLog, event)
}

// GetEventLog returns the event log
func (tr *TerminalRecorder) GetEventLog() []RecordingEvent {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	// Return a copy to avoid race conditions
	events := make([]RecordingEvent, len(tr.eventLog))
	copy(events, tr.eventLog)
	return events
}

// GetInputEvents returns the input event log
func (tr *TerminalRecorder) GetInputEvents() []InputEvent {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	// Note: inputEvents field removed from TerminalRecorder struct
	// Return empty slice for now
	return []InputEvent{}
}

// GenerateEnhancedRecording generates an enhanced recording with input visualization
func (tr *TerminalRecorder) GenerateEnhancedRecording() (string, error) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	if !tr.active && tr.outputFile == "" {
		return "", fmt.Errorf("no recording available")
	}

	// Generate enhanced version with input visualization
	enhancedFile := strings.TrimSuffix(tr.outputFile, filepath.Ext(tr.outputFile)) + "_enhanced" + filepath.Ext(tr.outputFile)

	// This is a placeholder for enhanced recording generation
	// In a real implementation, this would:
	// 1. Parse the original recording
	// 2. Add input visualization overlays
	// 3. Generate timing information
	// 4. Add metadata and annotations

	// For now, just copy the original file
	if err := tr.copyFile(tr.outputFile, enhancedFile); err != nil {
		return "", fmt.Errorf("failed to create enhanced recording: %w", err)
	}

	tr.logEvent("enhanced", "Enhanced recording generated", map[string]interface{}{
		"original_file": tr.outputFile,
		"enhanced_file": enhancedFile,
		"input_events":  0,
	})

	return enhancedFile, nil
}

// copyFile copies a file from src to dst
func (tr *TerminalRecorder) copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, 0644)
}

// RecordingManager manages multiple recording sessions
type RecordingManager struct {
	recordings map[string]*TerminalRecorder
	mu         sync.RWMutex
	config     *RecordingConfig
}

// NewRecordingManager creates a new recording manager
func NewRecordingManager(config *RecordingConfig) *RecordingManager {
	return &RecordingManager{
		recordings: make(map[string]*TerminalRecorder),
		config:     config,
	}
}

// StartRecording starts a new recording session
func (rm *RecordingManager) StartRecording(sessionID string) (*TerminalRecorder, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if _, exists := rm.recordings[sessionID]; exists {
		return nil, fmt.Errorf("recording session %s already exists", sessionID)
	}

	config := &RecordingConfig{
		OutputDir:   rm.config.OutputDir,
		TestName:    rm.config.TestName + "_" + sessionID,
		Format:      rm.config.Format,
		Quality:     rm.config.Quality,
		ShowInput:   rm.config.ShowInput,
		ShowTiming:  rm.config.ShowTiming,
		MaxDuration: rm.config.MaxDuration,
		Platform:    rm.config.Platform,
	}

	recorder := NewTerminalRecorder(config)
	if err := recorder.StartRecording(); err != nil {
		return nil, fmt.Errorf("failed to start recording: %w", err)
	}

	rm.recordings[sessionID] = recorder
	return recorder, nil
}

// StopRecording stops a recording session
func (rm *RecordingManager) StopRecording(sessionID string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	recorder, exists := rm.recordings[sessionID]
	if !exists {
		return fmt.Errorf("recording session %s not found", sessionID)
	}

	if err := recorder.StopRecording(); err != nil {
		return fmt.Errorf("failed to stop recording: %w", err)
	}

	delete(rm.recordings, sessionID)
	return nil
}

// GetRecording returns a recording session
func (rm *RecordingManager) GetRecording(sessionID string) (*TerminalRecorder, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	recorder, exists := rm.recordings[sessionID]
	if !exists {
		return nil, fmt.Errorf("recording session %s not found", sessionID)
	}

	return recorder, nil
}

// ListRecordings returns all active recording sessions
func (rm *RecordingManager) ListRecordings() map[string]*TerminalRecorder {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	// Return a copy to avoid race conditions
	copy := make(map[string]*TerminalRecorder)
	for id, recorder := range rm.recordings {
		copy[id] = recorder
	}
	return copy
}

// StopAllRecordings stops all active recording sessions
func (rm *RecordingManager) StopAllRecordings() error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	var errors []error
	for sessionID, recorder := range rm.recordings {
		if err := recorder.StopRecording(); err != nil {
			errors = append(errors, fmt.Errorf("failed to stop session %s: %w", sessionID, err))
		}
	}

	rm.recordings = make(map[string]*TerminalRecorder)

	if len(errors) > 0 {
		return fmt.Errorf("encountered %d errors while stopping recordings: %v", len(errors), errors)
	}

	return nil
}
