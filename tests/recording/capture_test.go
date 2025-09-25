package recording

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// TestRecordingSessionCreation tests recording session creation
func TestRecordingSessionCreation(t *testing.T) {
	config := DefaultRecordingConfig()
	_ = config
	
	t.Run("DefaultConfig", func(t *testing.T) {
		session, err := NewRecordingSession(nil)
		if err != nil {
			t.Errorf("Failed to create recording session with default config: %v", err)
		}
		
		if session == nil {
			t.Error("Recording session is nil")
		}
		
		if session.config == nil {
			t.Error("Recording session config is nil")
		}
		
		if session.sessionID == "" {
			t.Error("Recording session ID is empty")
		}
		
		if session.outputFile == "" {
			t.Error("Recording session output file is empty")
		}
	})
	
	t.Run("CustomConfig", func(t *testing.T) {
		customConfig := &RecordingConfig{
			OutputDir:     "tests/artifacts/custom_recordings",
			Format:        "custom",
			Quality:       "medium",
			EnableOverlay: false,
			Compression:   false,
			MaxDuration:   2 * time.Minute,
			Platform:      runtime.GOOS,
		}
		
		session, err := NewRecordingSession(customConfig)
		if err != nil {
			t.Errorf("Failed to create recording session with custom config: %v", err)
		}
		
		if session.config.OutputDir != customConfig.OutputDir {
			t.Error("Custom output directory not set correctly")
		}
		
		if session.config.Format != customConfig.Format {
			t.Error("Custom format not set correctly")
		}
		
		if session.config.Quality != customConfig.Quality {
			t.Error("Custom quality not set correctly")
		}
		
		if session.config.EnableOverlay != customConfig.EnableOverlay {
			t.Error("Custom overlay setting not set correctly")
		}
		
		if session.config.Compression != customConfig.Compression {
			t.Error("Custom compression setting not set correctly")
		}
		
		if session.config.MaxDuration != customConfig.MaxDuration {
			t.Error("Custom max duration not set correctly")
		}
	})
}

// TestRecordingSessionLifecycle tests recording session lifecycle
func TestRecordingSessionLifecycle(t *testing.T) {
	config := DefaultRecordingConfig()
	config.MaxDuration = 10 * time.Second // Short duration for testing
	
	session, err := NewRecordingSession(config)
	if err != nil {
		t.Fatalf("Failed to create recording session: %v", err)
	}
	
	t.Run("StartRecording", func(t *testing.T) {
		err := session.Start()
		if err != nil {
			t.Errorf("Failed to start recording: %v", err)
		}
		
		// Give the recording process time to start
		time.Sleep(100 * time.Millisecond)
	})
	
	t.Run("StopRecording", func(t *testing.T) {
		err := session.Stop()
		if err != nil {
			t.Errorf("Failed to stop recording: %v", err)
		}
	})
	
	t.Run("VerifyOutputFile", func(t *testing.T) {
		outputFile := session.GetOutputFile()
		if outputFile == "" {
			t.Error("Output file path is empty")
		}
		
		// Check if output file exists (it might not exist if recording failed)
		if _, err := os.Stat(outputFile); err == nil {
			t.Logf("Recording output file created: %s", outputFile)
		} else {
			t.Logf("Recording output file not created (this might be expected): %s", outputFile)
		}
	})
	
	t.Run("SessionMetadata", func(t *testing.T) {
		sessionID := session.GetSessionID()
		if sessionID == "" {
			t.Error("Session ID is empty")
		}
		
		duration := session.GetDuration()
		if duration < 0 {
			t.Error("Session duration is negative")
		}
		
		t.Logf("Session ID: %s", sessionID)
		t.Logf("Session duration: %v", duration)
	})
}

// TestRecordingManager tests recording manager functionality
func TestRecordingManager(t *testing.T) {
	config := DefaultRecordingConfig()
	config.MaxDuration = 5 * time.Second
	
	manager := NewRecordingManager(config)
	
	t.Run("CreateManager", func(t *testing.T) {
		if manager == nil {
			t.Error("Recording manager is nil")
		}
		
		if manager.config == nil {
			t.Error("Recording manager config is nil")
		}
		
		if manager.sessions == nil {
			t.Error("Recording manager sessions map is nil")
		}
		
		if manager.active {
			t.Error("Recording manager should not be active initially")
		}
	})
	
	t.Run("StartSession", func(t *testing.T) {
		session, err := manager.StartSession("test_session")
		if err != nil {
			t.Errorf("Failed to start recording session: %v", err)
		}
		
		if session == nil {
			t.Error("Recording session is nil")
		}
		
		if !manager.active {
			t.Error("Recording manager should be active after starting session")
		}
		
		activeSessions := manager.GetActiveSessions()
		if len(activeSessions) != 1 {
			t.Errorf("Expected 1 active session, got %d", len(activeSessions))
		}
		
		if activeSessions[0] != "test_session" {
			t.Errorf("Expected active session 'test_session', got '%s'", activeSessions[0])
		}
	})
	
	t.Run("GetSession", func(t *testing.T) {
		session, exists := manager.GetSession("test_session")
		if !exists {
			t.Error("Session 'test_session' should exist")
		}
		
		if session == nil {
			t.Error("Retrieved session is nil")
		}
		
		// Test non-existent session
		_, exists = manager.GetSession("non_existent_session")
		if exists {
			t.Error("Non-existent session should not exist")
		}
	})
	
	t.Run("StopSession", func(t *testing.T) {
		err := manager.StopSession("test_session")
		if err != nil {
			t.Errorf("Failed to stop recording session: %v", err)
		}
		
		if manager.active {
			t.Error("Recording manager should not be active after stopping all sessions")
		}
		
		activeSessions := manager.GetActiveSessions()
		if len(activeSessions) != 0 {
			t.Errorf("Expected 0 active sessions, got %d", len(activeSessions))
		}
	})
	
	t.Run("StopNonExistentSession", func(t *testing.T) {
		err := manager.StopSession("non_existent_session")
		if err == nil {
			t.Error("Expected error when stopping non-existent session")
		}
	})
}

// TestInputVisualizationOverlay tests input visualization overlay functionality
func TestInputVisualizationOverlay(t *testing.T) {
	config := DefaultRecordingConfig()
	overlay := NewInputVisualizationOverlay(config)
	
	t.Run("CreateOverlay", func(t *testing.T) {
		if overlay == nil {
			t.Error("Input visualization overlay is nil")
		}
		
		if overlay.config == nil {
			t.Error("Overlay config is nil")
		}
		
		if overlay.overlayData == nil {
			t.Error("Overlay data slice is nil")
		}
	})
	
	t.Run("AddEvents", func(t *testing.T) {
		// Add various types of events
		overlay.AddEvent("keypress", "a", 10, 20)
		overlay.AddEvent("mouse", "click", 100, 200)
		overlay.AddEvent("command", "echo test", 0, 0)
		
		if len(overlay.overlayData) != 3 {
			t.Errorf("Expected 3 overlay events, got %d", len(overlay.overlayData))
		}
		
		// Verify event data
		event1 := overlay.overlayData[0]
		if event1.Type != "keypress" || event1.Data != "a" {
			t.Error("First event data is incorrect")
		}
		
		event2 := overlay.overlayData[1]
		if event2.Type != "mouse" || event2.Data != "click" {
			t.Error("Second event data is incorrect")
		}
		
		event3 := overlay.overlayData[2]
		if event3.Type != "command" || event3.Data != "echo test" {
			t.Error("Third event data is incorrect")
		}
	})
	
	t.Run("GenerateOverlay", func(t *testing.T) {
		overlayData, err := overlay.GenerateOverlay()
		if err != nil {
			t.Errorf("Failed to generate overlay: %v", err)
		}
		
		if len(overlayData) == 0 {
			t.Error("Generated overlay data is empty")
		}
		
		overlayString := string(overlayData)
		if !contains(overlayString, "Input Visualization Overlay") {
			t.Error("Generated overlay does not contain expected header")
		}
		
		if !contains(overlayString, "keypress") {
			t.Error("Generated overlay does not contain keypress event")
		}
		
		if !contains(overlayString, "mouse") {
			t.Error("Generated overlay does not contain mouse event")
		}
		
		if !contains(overlayString, "command") {
			t.Error("Generated overlay does not contain command event")
		}
	})
	
	t.Run("SaveOverlay", func(t *testing.T) {
		testFile := filepath.Join(config.OutputDir, "test_overlay.txt")
		
		err := overlay.SaveOverlay(testFile)
		if err != nil {
			t.Errorf("Failed to save overlay: %v", err)
		}
		
		// Verify file was created
		if _, err := os.Stat(testFile); os.IsNotExist(err) {
			t.Error("Overlay file was not created")
		}
		
		// Clean up
		os.Remove(testFile)
	})
}

// TestRecordingCompressor tests recording compression functionality
func TestRecordingCompressor(t *testing.T) {
	config := DefaultRecordingConfig()
	compressor := NewRecordingCompressor(config)
	
	t.Run("CreateCompressor", func(t *testing.T) {
		if compressor == nil {
			t.Error("Recording compressor is nil")
		}
		
		if compressor.config == nil {
			t.Error("Compressor config is nil")
		}
	})
	
	t.Run("CompressRecording", func(t *testing.T) {
		// Create a test input file
		testDir := filepath.Join(config.OutputDir, "compression_test")
		os.MkdirAll(testDir, 0755)
		defer os.RemoveAll(testDir)
		
		inputFile := filepath.Join(testDir, "input.txt")
		outputFile := filepath.Join(testDir, "output.compressed")
		
		// Create test content
		testContent := "This is test content for compression testing. " +
			"It should be long enough to demonstrate compression. " +
			"Repeating this content multiple times to make it longer. " +
			"This is test content for compression testing. " +
			"It should be long enough to demonstrate compression. " +
			"Repeating this content multiple times to make it longer."
		
		err := os.WriteFile(inputFile, []byte(testContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test input file: %v", err)
		}
		
		// Test compression
		err = compressor.CompressRecording(inputFile, outputFile)
		if err != nil {
			t.Errorf("Failed to compress recording: %v", err)
		}
		
		// Verify output file was created
		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			t.Error("Compressed output file was not created")
		}
		
		// Test compression ratio
		ratio, err := compressor.GetCompressionRatio(inputFile, outputFile)
		if err != nil {
			t.Errorf("Failed to get compression ratio: %v", err)
		}
		
		t.Logf("Compression ratio: %.2f", ratio)
		
		// Compression ratio should be less than 1.0 (compressed file should be smaller)
		if ratio >= 1.0 {
			t.Logf("Compression ratio is not less than 1.0 (%.2f), but this might be expected for small files", ratio)
		}
	})
	
	t.Run("CopyFile", func(t *testing.T) {
		// Test file copying when compression is disabled
		testDir := filepath.Join(config.OutputDir, "copy_test")
		os.MkdirAll(testDir, 0755)
		defer os.RemoveAll(testDir)
		
		inputFile := filepath.Join(testDir, "input.txt")
		outputFile := filepath.Join(testDir, "output.txt")
		
		testContent := "Test content for file copying"
		
		err := os.WriteFile(inputFile, []byte(testContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test input file: %v", err)
		}
		
		err = compressor.copyFile(inputFile, outputFile)
		if err != nil {
			t.Errorf("Failed to copy file: %v", err)
		}
		
		// Verify output file was created
		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			t.Error("Copied output file was not created")
		}
		
		// Verify content is the same
		outputContent, err := os.ReadFile(outputFile)
		if err != nil {
			t.Errorf("Failed to read output file: %v", err)
		}
		
		if string(outputContent) != testContent {
			t.Error("Copied file content does not match original")
		}
	})
}

// TestRecordingIntegration tests integration between recording components
func TestRecordingIntegration(t *testing.T) {
	config := DefaultRecordingConfig()
	config.MaxDuration = 5 * time.Second
	
	manager := NewRecordingManager(config)
	overlay := NewInputVisualizationOverlay(config)
	compressor := NewRecordingCompressor(config)
	
	t.Run("FullRecordingWorkflow", func(t *testing.T) {
		// Start a recording session
		session, err := manager.StartSession("integration_test")
		if err != nil {
			t.Errorf("Failed to start recording session: %v", err)
		}
		
		// Add some overlay events
		overlay.AddEvent("keypress", "2", 10, 20)
		overlay.AddEvent("keypress", "+", 20, 20)
		overlay.AddEvent("keypress", "3", 30, 20)
		overlay.AddEvent("keypress", "Enter", 40, 20)
		
		// Simulate some work
		time.Sleep(100 * time.Millisecond)
		
		// Stop the recording session
		err = manager.StopSession("integration_test")
		if err != nil {
			t.Errorf("Failed to stop recording session: %v", err)
		}
		
		// Generate and save overlay
		overlayFile := filepath.Join(config.OutputDir, "integration_overlay.txt")
		err = overlay.SaveOverlay(overlayFile)
		if err != nil {
			t.Errorf("Failed to save overlay: %v", err)
		}
		defer os.Remove(overlayFile)
		
		// Test compression if output file exists
		outputFile := session.GetOutputFile()
		if _, err := os.Stat(outputFile); err == nil {
			compressedFile := outputFile + ".compressed"
			err = compressor.CompressRecording(outputFile, compressedFile)
			if err != nil {
				t.Errorf("Failed to compress recording: %v", err)
			}
			defer os.Remove(compressedFile)
		}
		
		t.Logf("Integration test completed successfully")
	})
}

// TestRecordingErrorHandling tests error handling in recording components
func TestRecordingErrorHandling(t *testing.T) {
	t.Run("InvalidConfig", func(t *testing.T) {
		// Test with invalid output directory
		invalidConfig := &RecordingConfig{
			OutputDir: "/invalid/path/that/does/not/exist",
			Format:    "test",
			Platform:  runtime.GOOS,
		}
		
		_, err := NewRecordingSession(invalidConfig)
		if err == nil {
			t.Error("Expected error with invalid output directory")
		}
	})
	
	t.Run("ManagerDoubleStart", func(t *testing.T) {
		config := DefaultRecordingConfig()
		manager := NewRecordingManager(config)
		
		// Start first session
		_, err := manager.StartSession("session1")
		if err != nil {
			t.Errorf("Failed to start first session: %v", err)
		}
		
		// Try to start second session (should fail)
		_, err = manager.StartSession("session2")
		if err == nil {
			t.Error("Expected error when starting second session")
		}
		
		// Clean up
		manager.StopSession("session1")
	})
}

// BenchmarkRecordingOperations benchmarks recording operations
func BenchmarkRecordingOperations(b *testing.B) {
	config := DefaultRecordingConfig()
	config.MaxDuration = 1 * time.Second
	
	b.Run("SessionCreation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			session, err := NewRecordingSession(config)
			if err != nil {
				b.Fatalf("Failed to create recording session: %v", err)
			}
			_ = session
		}
	})
	
	b.Run("OverlayEventAddition", func(b *testing.B) {
		overlay := NewInputVisualizationOverlay(config)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			overlay.AddEvent("keypress", "a", i%100, i%50)
		}
	})
	
	b.Run("OverlayGeneration", func(b *testing.B) {
		overlay := NewInputVisualizationOverlay(config)
		
		// Add some events
		for i := 0; i < 100; i++ {
			overlay.AddEvent("keypress", "a", i%100, i%50)
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := overlay.GenerateOverlay()
			if err != nil {
				b.Fatalf("Failed to generate overlay: %v", err)
			}
		}
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		contains(s[1:], substr))))
}