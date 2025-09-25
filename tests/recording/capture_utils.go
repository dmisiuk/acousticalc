package recording

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// RecordingManager manages terminal recording operations
type RecordingManager struct {
	outputDir   string
	sessions    map[string]*RecordingSession
	mutex       sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	isRecording bool
}

// RecordingSession represents a single recording session
type RecordingSession struct {
	ID        string            `json:"id"`
	StartTime time.Time         `json:"start_time"`
	EndTime   time.Time         `json:"end_time"`
	Platform  string            `json:"platform"`
	TestName  string            `json:"test_name"`
	Metadata  map[string]string `json:"metadata"`
	FilePath  string            `json:"file_path"`
	Status    SessionStatus     `json:"status"`
	mutex     sync.RWMutex
}

// SessionStatus represents the status of a recording session
type SessionStatus string

const (
	StatusPending   SessionStatus = "pending"
	StatusRecording SessionStatus = "recording"
	StatusCompleted SessionStatus = "completed"
	StatusFailed    SessionStatus = "failed"
	StatusCancelled SessionStatus = "cancelled"
)

// RecordingConfig holds configuration for recording operations
type RecordingConfig struct {
	OutputDir        string
	Platform         string
	MaxDuration      time.Duration
	CompressionLevel int
	IncludeInput     bool
	IncludeOutput    bool
	MetadataEnabled  bool
}

// NewRecordingManager creates a new recording manager
func NewRecordingManager(config RecordingConfig) *RecordingManager {
	ctx, cancel := context.WithCancel(context.Background())

	// Ensure output directory exists
	os.MkdirAll(filepath.Join(config.OutputDir, "sessions"), 0755)
	os.MkdirAll(filepath.Join(config.OutputDir, "metadata"), 0755)
	os.MkdirAll(filepath.Join(config.OutputDir, "demos"), 0755)

	return &RecordingManager{
		outputDir:   config.OutputDir,
		sessions:    make(map[string]*RecordingSession),
		ctx:         ctx,
		cancel:      cancel,
		isRecording: false,
	}
}

// StartRecording starts a new recording session
func (rm *RecordingManager) StartRecording(testName string, metadata map[string]string) (*RecordingSession, error) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	sessionID := fmt.Sprintf("%s_%s", testName, time.Now().Format("20060102_150405"))

	session := &RecordingSession{
		ID:        sessionID,
		StartTime: time.Now(),
		Platform:  runtime.GOOS,
		TestName:  testName,
		Metadata:  metadata,
		Status:    StatusPending,
	}

	// Create recording file path
	filename := fmt.Sprintf("%s.cast", sessionID)
	session.FilePath = filepath.Join(rm.outputDir, "sessions", filename)

	// Add session to manager
	rm.sessions[sessionID] = session

	// Start recording process
	if err := rm.startRecordingProcess(session); err != nil {
		session.setStatus(StatusFailed)
		return nil, fmt.Errorf("failed to start recording: %w", err)
	}

	session.setStatus(StatusRecording)
	rm.isRecording = true

	return session, nil
}

// StopRecording stops a recording session
func (rm *RecordingManager) StopRecording(sessionID string) error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	session, exists := rm.sessions[sessionID]
	if !exists {
		return fmt.Errorf("recording session not found: %s", sessionID)
	}

	session.mutex.Lock()
	defer session.mutex.Unlock()

	if session.Status != StatusRecording {
		return fmt.Errorf("session is not recording: %s", session.Status)
	}

	// Stop recording process
	if err := rm.stopRecordingProcess(session); err != nil {
		session.Status = StatusFailed
		return fmt.Errorf("failed to stop recording: %w", err)
	}

	session.EndTime = time.Now()
	session.Status = StatusCompleted
	rm.isRecording = false

	// Generate metadata file
	if err := rm.generateMetadata(session); err != nil {
		// Don't fail the recording stop, but log the error
		fmt.Printf("Warning: failed to generate metadata for session %s: %v\n", sessionID, err)
	}

	return nil
}

// GetSession retrieves a recording session by ID
func (rm *RecordingManager) GetSession(sessionID string) (*RecordingSession, error) {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	session, exists := rm.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	return session, nil
}

// ListSessions returns all recording sessions
func (rm *RecordingManager) ListSessions() []*RecordingSession {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	sessions := make([]*RecordingSession, 0, len(rm.sessions))
	for _, session := range rm.sessions {
		sessions = append(sessions, session)
	}

	return sessions
}

// IsRecording returns whether any recording is active
func (rm *RecordingManager) IsRecording() bool {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	return rm.isRecording
}

// Close gracefully shuts down the recording manager
func (rm *RecordingManager) Close() error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Stop all active recordings
	for _, session := range rm.sessions {
		if session.Status == StatusRecording {
			session.setStatus(StatusCancelled)
			rm.stopRecordingProcess(session)
		}
	}

	if rm.cancel != nil {
		rm.cancel()
	}

	return nil
}

// startRecordingProcess starts the actual recording process for a session
func (rm *RecordingManager) startRecordingProcess(session *RecordingSession) error {
	// Platform-specific recording implementation would go here
	// For now, we'll simulate the recording process

	switch session.Platform {
	case "linux":
		return rm.startLinuxRecording(session)
	case "darwin":
		return rm.startMacOSRecording(session)
	case "windows":
		return rm.startWindowsRecording(session)
	default:
		return fmt.Errorf("unsupported platform for recording: %s", session.Platform)
	}
}

// stopRecordingProcess stops the actual recording process for a session
func (rm *RecordingManager) stopRecordingProcess(session *RecordingSession) error {
	// Platform-specific recording termination would go here
	// For now, we'll simulate the stop process

	switch session.Platform {
	case "linux":
		return rm.stopLinuxRecording(session)
	case "darwin":
		return rm.stopMacOSRecording(session)
	case "windows":
		return rm.stopWindowsRecording(session)
	default:
		return fmt.Errorf("unsupported platform for recording: %s", session.Platform)
	}
}

// startLinuxRecording starts recording on Linux
func (rm *RecordingManager) startLinuxRecording(session *RecordingSession) error {
	// In a real implementation, this would use tools like:
	// - asciinema for terminal recording
	// - script command for session recording
	// - custom Go implementation for input/output capture

	// For testing purposes, create a placeholder file
	file, err := os.Create(session.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write asciinema-compatible header
	header := map[string]interface{}{
		"version":   2,
		"width":     80,
		"height":    24,
		"timestamp": session.StartTime.Unix(),
		"env": map[string]string{
			"SHELL": "/bin/bash",
			"TERM":  "xterm-256color",
		},
	}

	headerBytes, _ := json.Marshal(header)
	file.Write(headerBytes)
	file.Write([]byte("\n"))

	return nil
}

// stopLinuxRecording stops recording on Linux
func (rm *RecordingManager) stopLinuxRecording(session *RecordingSession) error {
	// In a real implementation, this would terminate the recording process
	// For testing purposes, we'll just verify the file exists

	if _, err := os.Stat(session.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("recording file not found: %s", session.FilePath)
	}

	return nil
}

// startMacOSRecording starts recording on macOS
func (rm *RecordingManager) startMacOSRecording(session *RecordingSession) error {
	// Similar to Linux but with macOS-specific considerations
	return rm.startLinuxRecording(session) // Use same implementation for now
}

// stopMacOSRecording stops recording on macOS
func (rm *RecordingManager) stopMacOSRecording(session *RecordingSession) error {
	return rm.stopLinuxRecording(session) // Use same implementation for now
}

// startWindowsRecording starts recording on Windows
func (rm *RecordingManager) startWindowsRecording(session *RecordingSession) error {
	// Windows-specific recording implementation
	// For testing purposes, create a placeholder file similar to Unix
	return rm.startLinuxRecording(session) // Use same implementation for now
}

// stopWindowsRecording stops recording on Windows
func (rm *RecordingManager) stopWindowsRecording(session *RecordingSession) error {
	return rm.stopLinuxRecording(session) // Use same implementation for now
}

// generateMetadata generates metadata file for a completed session
func (rm *RecordingManager) generateMetadata(session *RecordingSession) error {
	metadataPath := filepath.Join(rm.outputDir, "metadata", fmt.Sprintf("%s_metadata.json", session.ID))

	// Create comprehensive metadata
	metadata := map[string]interface{}{
		"session_id":      session.ID,
		"test_name":       session.TestName,
		"platform":        session.Platform,
		"start_time":      session.StartTime.Format(time.RFC3339),
		"end_time":        session.EndTime.Format(time.RFC3339),
		"duration":        session.EndTime.Sub(session.StartTime).Seconds(),
		"file_path":       session.FilePath,
		"status":          session.Status,
		"go_version":      runtime.Version(),
		"go_arch":         runtime.GOARCH,
		"custom_metadata": session.Metadata,
	}

	// Add file information if recording file exists
	if info, err := os.Stat(session.FilePath); err == nil {
		metadata["file_size"] = info.Size()
		metadata["file_mode"] = info.Mode().String()
	}

	// Write metadata to file
	metadataBytes, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metadataPath, metadataBytes, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

// setStatus safely sets the session status
func (rs *RecordingSession) setStatus(status SessionStatus) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	rs.Status = status
}

// GetStatus safely gets the session status
func (rs *RecordingSession) GetStatus() SessionStatus {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	return rs.Status
}

// GetDuration returns the duration of the recording session
func (rs *RecordingSession) GetDuration() time.Duration {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	if rs.Status == StatusRecording {
		return time.Since(rs.StartTime)
	}

	if !rs.EndTime.IsZero() {
		return rs.EndTime.Sub(rs.StartTime)
	}

	return 0
}

// DefaultRecordingConfig returns a default recording configuration
func DefaultRecordingConfig() RecordingConfig {
	return RecordingConfig{
		OutputDir:        filepath.Join("tests", "artifacts", "recordings"),
		Platform:         runtime.GOOS,
		MaxDuration:      5 * time.Minute,
		CompressionLevel: 5,
		IncludeInput:     true,
		IncludeOutput:    true,
		MetadataEnabled:  true,
	}
}
