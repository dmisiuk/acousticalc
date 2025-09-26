package recording

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// RecordingConfig holds configuration for terminal recording
type RecordingConfig struct {
	OutputDir      string
	Format         string
	Quality        string
	EnableOverlay  bool
	OverlayStyle   string
	Compression    bool
	MaxDuration    time.Duration
	Platform       string
}

// DefaultRecordingConfig returns default configuration for terminal recording
func DefaultRecordingConfig() *RecordingConfig {
	return &RecordingConfig{
		OutputDir:     "tests/artifacts/recordings",
		Format:        "asciinema",
		Quality:       "high",
		EnableOverlay: true,
		OverlayStyle:  "modern",
		Compression:   true,
		MaxDuration:   5 * time.Minute,
		Platform:      runtime.GOOS,
	}
}

// RecordingSession manages a terminal recording session
type RecordingSession struct {
	config     *RecordingConfig
	sessionID  string
	startTime  time.Time
	outputFile string
	process    *exec.Cmd
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewRecordingSession creates a new recording session
func NewRecordingSession(config *RecordingConfig) (*RecordingSession, error) {
	if config == nil {
		config = DefaultRecordingConfig()
	}
	
	// Ensure output directory exists
	err := os.MkdirAll(config.OutputDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create recording directory: %v", err)
	}
	
	// Generate unique session ID
	sessionID := fmt.Sprintf("e2e_session_%d", time.Now().UnixNano())
	outputFile := filepath.Join(config.OutputDir, fmt.Sprintf("%s.%s", sessionID, config.Format))
	
	ctx, cancel := context.WithTimeout(context.Background(), config.MaxDuration)
	
	return &RecordingSession{
		config:     config,
		sessionID:  sessionID,
		startTime:  time.Now(),
		outputFile: outputFile,
		ctx:        ctx,
		cancel:     cancel,
	}, nil
}

// Start begins the recording session
func (rs *RecordingSession) Start() error {
	var cmd *exec.Cmd
	
	switch rs.config.Platform {
	case "windows":
		cmd = rs.createWindowsRecordingCommand()
	case "darwin":
		cmd = rs.createMacOSRecordingCommand()
	case "linux":
		cmd = rs.createLinuxRecordingCommand()
	default:
		return fmt.Errorf("unsupported platform for recording: %s", rs.config.Platform)
	}
	
	if cmd == nil {
		return fmt.Errorf("failed to create recording command for platform: %s", rs.config.Platform)
	}
	
	rs.process = cmd
	
	// Start the recording process
	err := rs.process.Start()
	if err != nil {
		return fmt.Errorf("failed to start recording: %v", err)
	}
	
	return nil
}

// Stop ends the recording session
func (rs *RecordingSession) Stop() error {
	if rs.process == nil {
		return fmt.Errorf("recording session not started")
	}
	
	// Cancel the context to stop the recording
	rs.cancel()
	
	// Wait for the process to finish
	err := rs.process.Wait()
	if err != nil {
		// Don't return error if process was terminated by context
		if rs.ctx.Err() == context.Canceled {
			return nil
		}
		return fmt.Errorf("recording process failed: %v", err)
	}
	
	return nil
}

// GetOutputFile returns the path to the recording output file
func (rs *RecordingSession) GetOutputFile() string {
	return rs.outputFile
}

// GetSessionID returns the session ID
func (rs *RecordingSession) GetSessionID() string {
	return rs.sessionID
}

// GetDuration returns the recording duration
func (rs *RecordingSession) GetDuration() time.Duration {
	return time.Since(rs.startTime)
}

// createWindowsRecordingCommand creates a Windows recording command
func (rs *RecordingSession) createWindowsRecordingCommand() *exec.Cmd {
	// For Windows, we'll use PowerShell to capture terminal output
	// This is a simplified approach - in a real implementation, you might use
	// more sophisticated tools like Windows Terminal's recording features
	
	script := fmt.Sprintf(`
		$outputFile = "%s"
		$startTime = Get-Date
		$sessionID = "%s"
		
		# Create a simple recording by capturing command output
		Write-Host "Recording session $sessionID started at $startTime" | Out-File -FilePath $outputFile -Append
		
		# Wait for the recording to complete
		Start-Sleep -Seconds 1
		
		Write-Host "Recording session $sessionID completed" | Out-File -FilePath $outputFile -Append
	`, rs.outputFile, rs.sessionID)
	
	return exec.CommandContext(rs.ctx, "powershell", "-Command", script)
}

// createMacOSRecordingCommand creates a macOS recording command
func (rs *RecordingSession) createMacOSRecordingCommand() *exec.Cmd {
	// For macOS, we'll use script command for terminal recording
	// This captures terminal input/output to a file
	
	args := []string{
		"-q",                    // Quiet mode
		"-t", "0",              // No timing file
		rs.outputFile,          // Output file
		"echo", "Recording started", // Command to record
	}
	
	return exec.CommandContext(rs.ctx, "script", args...)
}

// createLinuxRecordingCommand creates a Linux recording command
func (rs *RecordingSession) createLinuxRecordingCommand() *exec.Cmd {
	// For Linux, we'll use script command for terminal recording
	// This is similar to macOS but with Linux-specific options
	
	args := []string{
		"-q",                    // Quiet mode
		"-t", "0",              // No timing file
		rs.outputFile,          // Output file
		"echo", "Recording started", // Command to record
	}
	
	return exec.CommandContext(rs.ctx, "script", args...)
}

// RecordingManager manages multiple recording sessions
type RecordingManager struct {
	config    *RecordingConfig
	sessions  map[string]*RecordingSession
	active    bool
}

// NewRecordingManager creates a new recording manager
func NewRecordingManager(config *RecordingConfig) *RecordingManager {
	if config == nil {
		config = DefaultRecordingConfig()
	}
	
	return &RecordingManager{
		config:   config,
		sessions: make(map[string]*RecordingSession),
		active:   false,
	}
}

// StartSession starts a new recording session
func (rm *RecordingManager) StartSession(sessionName string) (*RecordingSession, error) {
	if rm.active {
		return nil, fmt.Errorf("recording manager is already active")
	}
	
	session, err := NewRecordingSession(rm.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create recording session: %v", err)
	}
	
	err = session.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start recording session: %v", err)
	}
	
	rm.sessions[sessionName] = session
	rm.active = true
	
	return session, nil
}

// StopSession stops a recording session
func (rm *RecordingManager) StopSession(sessionName string) error {
	session, exists := rm.sessions[sessionName]
	if !exists {
		return fmt.Errorf("recording session %s not found", sessionName)
	}
	
	err := session.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop recording session: %v", err)
	}
	
	delete(rm.sessions, sessionName)
	
	if len(rm.sessions) == 0 {
		rm.active = false
	}
	
	return nil
}

// StopAllSessions stops all active recording sessions
func (rm *RecordingManager) StopAllSessions() error {
	var errors []string
	
	for sessionName, session := range rm.sessions {
		err := session.Stop()
		if err != nil {
			errors = append(errors, fmt.Sprintf("session %s: %v", sessionName, err))
		}
	}
	
	rm.sessions = make(map[string]*RecordingSession)
	rm.active = false
	
	if len(errors) > 0 {
		return fmt.Errorf("failed to stop some sessions: %s", strings.Join(errors, "; "))
	}
	
	return nil
}

// GetSession returns a recording session by name
func (rm *RecordingManager) GetSession(sessionName string) (*RecordingSession, bool) {
	session, exists := rm.sessions[sessionName]
	return session, exists
}

// IsActive returns whether the recording manager is active
func (rm *RecordingManager) IsActive() bool {
	return rm.active
}

// GetActiveSessions returns the names of active sessions
func (rm *RecordingManager) GetActiveSessions() []string {
	var sessions []string
	for name := range rm.sessions {
		sessions = append(sessions, name)
	}
	return sessions
}

// InputVisualizationOverlay creates input visualization overlays for recordings
type InputVisualizationOverlay struct {
	config     *RecordingConfig
	overlayData []OverlayEvent
}

// OverlayEvent represents a user input event for visualization
type OverlayEvent struct {
	Timestamp time.Time
	Type      string // "keypress", "mouse", "command"
	Data      string
	Position  struct {
		X int
		Y int
	}
}

// NewInputVisualizationOverlay creates a new input visualization overlay
func NewInputVisualizationOverlay(config *RecordingConfig) *InputVisualizationOverlay {
	if config == nil {
		config = DefaultRecordingConfig()
	}
	
	return &InputVisualizationOverlay{
		config:     config,
		overlayData: make([]OverlayEvent, 0),
	}
}

// AddEvent adds an input event to the overlay
func (ivo *InputVisualizationOverlay) AddEvent(eventType, data string, x, y int) {
	event := OverlayEvent{
		Timestamp: time.Now(),
		Type:      eventType,
		Data:      data,
		Position: struct {
			X int
			Y int
		}{X: x, Y: y},
	}
	
	ivo.overlayData = append(ivo.overlayData, event)
}

// GenerateOverlay generates the overlay data for the recording
func (ivo *InputVisualizationOverlay) GenerateOverlay() ([]byte, error) {
	// This is a simplified implementation
	// In a real implementation, you would generate proper overlay data
	// that can be rendered on top of the terminal recording
	
	overlay := fmt.Sprintf(`
# Input Visualization Overlay
# Generated at: %s
# Platform: %s
# Style: %s

`, time.Now().Format(time.RFC3339), ivo.config.Platform, ivo.config.OverlayStyle)
	
	for _, event := range ivo.overlayData {
		overlay += fmt.Sprintf("%s: %s at (%d, %d) - %s\n",
			event.Timestamp.Format("15:04:05.000"),
			event.Type,
			event.Position.X,
			event.Position.Y,
			event.Data)
	}
	
	return []byte(overlay), nil
}

// SaveOverlay saves the overlay data to a file
func (ivo *InputVisualizationOverlay) SaveOverlay(filename string) error {
	overlayData, err := ivo.GenerateOverlay()
	if err != nil {
		return fmt.Errorf("failed to generate overlay: %v", err)
	}
	
	err = os.WriteFile(filename, overlayData, 0644)
	if err != nil {
		return fmt.Errorf("failed to save overlay: %v", err)
	}
	
	return nil
}

// RecordingCompressor handles compression of recording files
type RecordingCompressor struct {
	config *RecordingConfig
}

// NewRecordingCompressor creates a new recording compressor
func NewRecordingCompressor(config *RecordingConfig) *RecordingCompressor {
	if config == nil {
		config = DefaultRecordingConfig()
	}
	
	return &RecordingCompressor{
		config: config,
	}
}

// CompressRecording compresses a recording file
func (rc *RecordingCompressor) CompressRecording(inputFile, outputFile string) error {
	if !rc.config.Compression {
		// If compression is disabled, just copy the file
		return rc.copyFile(inputFile, outputFile)
	}
	
	// Use platform-specific compression tools
	switch rc.config.Platform {
	case "windows":
		return rc.compressWindows(inputFile, outputFile)
	case "darwin", "linux":
		return rc.compressUnix(inputFile, outputFile)
	default:
		return fmt.Errorf("unsupported platform for compression: %s", rc.config.Platform)
	}
}

// copyFile copies a file from source to destination
func (rc *RecordingCompressor) copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source file: %v", err)
	}
	
	err = os.WriteFile(dst, input, 0644)
	if err != nil {
		return fmt.Errorf("failed to write destination file: %v", err)
	}
	
	return nil
}

// compressWindows compresses a file on Windows
func (rc *RecordingCompressor) compressWindows(inputFile, outputFile string) error {
	// For Windows, we'll use PowerShell compression
	script := fmt.Sprintf(`
		$inputFile = "%s"
		$outputFile = "%s"
		
		# Use .NET compression
		$bytes = [System.IO.File]::ReadAllBytes($inputFile)
		$compressed = [System.IO.Compression.GzipStream]::new(
			[System.IO.File]::Create($outputFile),
			[System.IO.Compression.CompressionMode]::Compress
		)
		$compressed.Write($bytes, 0, $bytes.Length)
		$compressed.Close()
	`, inputFile, outputFile)
	
	cmd := exec.Command("powershell", "-Command", script)
	return cmd.Run()
}

// compressUnix compresses a file on Unix-like systems
func (rc *RecordingCompressor) compressUnix(inputFile, outputFile string) error {
	// Use gzip for compression
	cmd := exec.Command("gzip", "-c", inputFile)
	
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer output.Close()
	
	cmd.Stdout = output
	
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to compress file: %v", err)
	}
	
	return nil
}

// GetCompressionRatio returns the compression ratio for a file
func (rc *RecordingCompressor) GetCompressionRatio(originalFile, compressedFile string) (float64, error) {
	originalInfo, err := os.Stat(originalFile)
	if err != nil {
		return 0, fmt.Errorf("failed to stat original file: %v", err)
	}
	
	compressedInfo, err := os.Stat(compressedFile)
	if err != nil {
		return 0, fmt.Errorf("failed to stat compressed file: %v", err)
	}
	
	originalSize := float64(originalInfo.Size())
	compressedSize := float64(compressedInfo.Size())
	
	if originalSize == 0 {
		return 0, nil
	}
	
	return compressedSize / originalSize, nil
}