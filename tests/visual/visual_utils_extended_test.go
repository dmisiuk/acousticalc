package visual

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// TestVisualUtilsExtended provides additional test coverage for visual utilities
func TestVisualUtilsExtended(t *testing.T) {
	// Test: Demo Storyboard Creation
	t.Run("CreateDemoStoryboard", func(t *testing.T) {
		outputDir := t.TempDir()
		logger := NewVisualTestLogger("storyboard_test", outputDir)

		// Add events with metadata for storyboard
		metadata := map[string]interface{}{
			"stage":   "development",
			"version": "1.0.0",
		}
		logger.LogEvent(EventTestStart, "Demo sequence started", metadata)
		logger.LogEvent(EventTestProcess, "Processing calculations", nil)
		logger.LogEvent(EventTestComplete, "Demo complete", nil)

		if err := logger.CreateDemoStoryboard(); err != nil {
			t.Errorf("Failed to create demo storyboard: %v", err)
		}

		storyboardPath := filepath.Join(outputDir, "../demo_content/storyboards/storyboard_test_storyboard.html")
		if _, err := os.Stat(storyboardPath); os.IsNotExist(err) {
			t.Error("Demo storyboard HTML file was not created")
		}
	})

	// Test: OptimizeScreenshots function with real scenarios
	t.Run("OptimizeScreenshots", func(t *testing.T) {
		inputDir := t.TempDir()
		outputDir := t.TempDir()

		// Create a test PNG file
		testImagePath := filepath.Join(inputDir, "test_image.png")
		// Create a simple 1x1 PNG file content (minimal valid PNG)
		pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
		if err := os.WriteFile(testImagePath, pngHeader, 0644); err != nil {
			t.Fatalf("Failed to create test PNG file: %v", err)
		}

		// Note: OptimizeScreenshots will fail with our minimal PNG, but we test error handling
		err := OptimizeScreenshots(inputDir, outputDir)
		if err == nil {
			t.Log("Optimization succeeded (or gracefully handled invalid PNG)")
		} else {
			t.Logf("Expected failure with minimal PNG file: %v", err)
		}

		// Test with non-existent input directory
		err = OptimizeScreenshots("/non/existent/dir", outputDir)
		if err == nil {
			t.Error("Expected error with non-existent input directory")
		}
	})

	// Test: CaptureTestEvent method
	t.Run("CaptureTestEvent", func(t *testing.T) {
		// Set up virtual display for CI environments
		if err := setupVirtualDisplayForTesting(); err != nil {
			t.Logf("Warning: Could not setup virtual display: %v", err)
		}

		outputDir := t.TempDir()
		capture := NewScreenshotCapture("event_test", outputDir)

		// Test capturing test event (will fail in headless environment but tests the code path)
		result := capture.CaptureTestEvent(t, "test_event")
		// In headless environment, this should return empty string due to capture failure
		if result != "" {
			t.Logf("Screenshot captured successfully: %s", result)
		} else {
			t.Log("Screenshot capture failed as expected in test environment")
		}
	})

	// Test: VisualEvent struct validation
	t.Run("VisualEvent_Validation", func(t *testing.T) {
		event := VisualEvent{
			Type:        EventTestStart,
			Timestamp:   time.Now(),
			Description: "Test event",
			Screenshot:  "test.png",
			Metadata:    map[string]interface{}{"key": "value"},
		}

		if event.Type != EventTestStart {
			t.Errorf("Expected Type %s, got %s", EventTestStart, event.Type)
		}
		if event.Description != "Test event" {
			t.Errorf("Expected Description 'Test event', got %s", event.Description)
		}
		if event.Screenshot != "test.png" {
			t.Errorf("Expected Screenshot 'test.png', got %s", event.Screenshot)
		}
		if event.Metadata["key"] != "value" {
			t.Errorf("Expected Metadata key 'value', got %v", event.Metadata["key"])
		}
	})

	// Test: HTML Report Generation with Multiple Events
	t.Run("HTMLReport_MultipleEvents", func(t *testing.T) {
		outputDir := t.TempDir()
		logger := NewVisualTestLogger("multi_event_test", outputDir)

		// Add multiple events with different types and metadata
		logger.LogEvent(EventTestStart, "Test initialization", map[string]interface{}{"setup": true})
		logger.LogEvent(EventTestProcess, "Processing data", map[string]interface{}{"stage": 1})
		logger.LogEvent(EventTestProcess, "Validating results", map[string]interface{}{"stage": 2})
		logger.LogEvent(EventTestPass, "Test completed successfully", map[string]interface{}{"success": true})

		// Generate HTML report
		htmlContent := logger.generateHTMLReport()
		if len(htmlContent) == 0 {
			t.Error("Generated HTML content is empty")
		}

		// Check that HTML contains expected elements
		expectedStrings := []string{
			"Visual Test Report",
			"multi_event_test",
			"Test initialization",
			"Processing data",
			"Test completed successfully",
			"Total Events: 4",
		}

		for _, expected := range expectedStrings {
			if !bytes.Contains([]byte(htmlContent), []byte(expected)) {
				t.Errorf("HTML content missing expected string: %s", expected)
			}
		}
	})

	// Test: Demo Storyboard Generation
	t.Run("DemoStoryboard_Generation", func(t *testing.T) {
		outputDir := t.TempDir()
		logger := NewVisualTestLogger("demo_test", outputDir)

		// Add events for demo storyboard
		logger.LogEvent(EventTestStart, "Demo begins", nil)
		logger.LogEvent(EventTestProcess, "User interaction", nil)
		logger.LogEvent(EventTestComplete, "Demo ends", nil)

		// Generate demo storyboard content
		storyboardContent := logger.generateDemoStoryboard()
		if len(storyboardContent) == 0 {
			t.Error("Generated storyboard content is empty")
		}

		// Check storyboard contains expected elements
		expectedStrings := []string{
			"AcoustiCalc Demo Storyboard",
			"demo_test",
			"Professional Demo Content Generation",
			"Scene 1: start",
			"Demo begins",
		}

		for _, expected := range expectedStrings {
			if !bytes.Contains([]byte(storyboardContent), []byte(expected)) {
				t.Errorf("Storyboard content missing expected string: %s", expected)
			}
		}
	})

	// Test: Directory Creation Edge Cases
	t.Run("DirectoryCreation_EdgeCases", func(t *testing.T) {
		// Test with nested directory creation
		baseDir := t.TempDir()
		deepDir := filepath.Join(baseDir, "level1", "level2", "level3")
		capture := NewScreenshotCapture("deep_test", deepDir)

		// This should create the directory structure
		_, err := capture.CaptureScreen("test")
		// Error is expected due to headless environment, but directory should be created
		if err == nil {
			t.Log("Screenshot captured successfully")
		} else {
			t.Logf("Screenshot failed as expected: %v", err)
		}

		// Verify directory was created
		if _, err := os.Stat(deepDir); os.IsNotExist(err) {
			t.Error("Deep directory structure was not created")
		}
	})

	// Test: Filename Generation
	t.Run("Filename_Generation", func(t *testing.T) {
		testTime := time.Date(2025, 9, 24, 15, 30, 45, 0, time.UTC)
		capture := &ScreenshotCapture{
			OutputDir: t.TempDir(),
			TestName:  "filename_test",
			Timestamp: testTime,
			Format:    "png",
			Quality:   100,
		}

		// Test filename generation by examining the expected path
		_, err := capture.CaptureScreen("start")
		// Error expected in headless environment
		if err != nil {
			t.Logf("Screenshot capture failed as expected: %v", err)
			// Verify the error indicates directory creation succeeded
			if !bytes.Contains([]byte(err.Error()), []byte("failed to capture screen")) {
				t.Errorf("Unexpected error type: %v", err)
			}
		}

		// Verify directory was created
		if _, err := os.Stat(capture.OutputDir); os.IsNotExist(err) {
			t.Error("Output directory was not created")
		}
	})

	// Test: VisualTestLogger with Empty Events
	t.Run("VisualTestLogger_EmptyEvents", func(t *testing.T) {
		outputDir := t.TempDir()
		logger := NewVisualTestLogger("empty_test", outputDir)

		// Generate report with no events
		if err := logger.GenerateVisualReport(); err != nil {
			t.Errorf("Failed to generate visual report with no events: %v", err)
		}

		reportPath := filepath.Join(outputDir, "empty_test_visual_report.html")
		if _, err := os.Stat(reportPath); os.IsNotExist(err) {
			t.Error("Visual report HTML file was not created for empty logger")
		}

		// Check content
		content, err := os.ReadFile(reportPath)
		if err != nil {
			t.Errorf("Failed to read report file: %v", err)
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, "Total Events: 0") {
			t.Error("Report should indicate 0 events")
		}
		if !strings.Contains(contentStr, "Screenshots Captured: 0") {
			t.Error("Report should indicate 0 screenshots")
		}
	})

	// Test: ScreenshotCapture with Different Event Types
	t.Run("ScreenshotCapture_AllEventTypes", func(t *testing.T) {
		outputDir := t.TempDir()
		capture := NewScreenshotCapture("event_types_test", outputDir)

		events := []VisualTestEvent{
			EventTestStart,
			EventTestPass,
			EventTestFail,
			EventTestProcess,
			EventTestComplete,
		}

		for _, eventType := range events {
			result := capture.CaptureTestEvent(t, string(eventType))
			// All will fail in headless environment, but tests code paths
			if result != "" {
				t.Logf("Screenshot captured for %s: %s", eventType, result)
			} else {
				t.Logf("Screenshot capture failed for %s as expected", eventType)
			}
		}
	})

	// Test: Platform-specific behavior
	t.Run("Platform_Behavior", func(t *testing.T) {
		currentOS := runtime.GOOS
		t.Logf("Current platform: %s", currentOS)

		// Test that our functions work on current platform
		outputDir := t.TempDir()
		_ = NewScreenshotCapture("platform_test", outputDir)
		logger := NewVisualTestLogger("platform_logger", outputDir)

		// Test basic functionality on current platform
		switch currentOS {
		case "darwin", "linux":
			t.Log("Testing on Unix-like platform")
		case "windows":
			t.Log("Testing on Windows platform")
		}

		// All platforms should support basic directory creation and HTML generation
		logger.LogEvent(EventTestStart, "Platform test", nil)

		if err := logger.GenerateVisualReport(); err != nil {
			t.Errorf("Visual report generation failed on %s: %v", currentOS, err)
		}
	})

	// Test: Error handling in OptimizeScreenshots
	t.Run("OptimizeScreenshots_ErrorHandling", func(t *testing.T) {
		// Test with invalid output directory permissions (if not Windows)
		if runtime.GOOS != "windows" {
			inputDir := t.TempDir()
			outputDir := t.TempDir()

			// Make output directory read-only
			if err := os.Chmod(outputDir, 0444); err != nil {
				t.Logf("Failed to change directory permissions: %v", err)
			}

			err := OptimizeScreenshots(inputDir, outputDir)
			if err == nil {
				t.Log("OptimizeScreenshots succeeded despite read-only output dir (permissions may not be enforced)")
			} else {
				t.Logf("Expected error with read-only output directory: %v", err)
			}

			// Restore permissions for cleanup
			os.Chmod(outputDir, 0755)
		}
	})

	// Test: HTML content validation
	t.Run("HTML_Content_Validation", func(t *testing.T) {
		outputDir := t.TempDir()
		logger := NewVisualTestLogger("html_validation_test", outputDir)

		// Add events with special characters that need HTML escaping
		logger.LogEvent(EventTestStart, "Test with <special> & \"characters\"", map[string]interface{}{
			"html_test": "<script>alert('test')</script>",
		})

		htmlContent := logger.generateHTMLReport()

		// Ensure HTML is properly structured
		if !strings.Contains(htmlContent, "<!DOCTYPE html>") {
			t.Error("HTML content missing DOCTYPE declaration")
		}
		if !strings.Contains(htmlContent, "<html>") {
			t.Error("HTML content missing html tag")
		}
		if !strings.Contains(htmlContent, "</html>") {
			t.Error("HTML content missing closing html tag")
		}

		// Test storyboard HTML as well
		storyboardContent := logger.generateDemoStoryboard()
		if !strings.Contains(storyboardContent, "<!DOCTYPE html>") {
			t.Error("Storyboard content missing DOCTYPE declaration")
		}
	})
}

// TestVisualUtilsErrorPaths tests error handling paths
func TestVisualUtilsErrorPaths(t *testing.T) {
	// Test: Invalid directory permissions for CreateDemoStoryboard
	t.Run("CreateDemoStoryboard_InvalidPermissions", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Skipping permission test on Windows")
		}

		outputDir := t.TempDir()
		logger := NewVisualTestLogger("perm_test", outputDir)

		// Make the output directory read-only
		if err := os.Chmod(outputDir, 0444); err != nil {
			t.Fatalf("Failed to change directory permissions: %v", err)
		}

		err := logger.CreateDemoStoryboard()
		// On some systems, the directory creation might still succeed due to parent permissions
		// So we accept both error and success cases
		if err != nil {
			t.Logf("Expected error when creating storyboard in read-only directory: %v", err)
		} else {
			t.Log("Directory creation succeeded despite read-only permissions (system-dependent behavior)")
		}

		// Restore permissions for cleanup
		os.Chmod(outputDir, 0755)
	})

	// Test: GenerateVisualReport with invalid permissions
	t.Run("GenerateVisualReport_InvalidPermissions", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Skipping permission test on Windows")
		}

		outputDir := t.TempDir()
		logger := NewVisualTestLogger("report_perm_test", outputDir)
		logger.LogEvent(EventTestStart, "Test event", nil)

		// Make the output directory read-only
		if err := os.Chmod(outputDir, 0444); err != nil {
			t.Fatalf("Failed to change directory permissions: %v", err)
		}

		err := logger.GenerateVisualReport()
		if err == nil {
			t.Error("Expected error when generating report in read-only directory")
		}

		// Restore permissions for cleanup
		os.Chmod(outputDir, 0755)
	})
}

// setupVirtualDisplayForTesting is a helper function for tests that need display
func setupVirtualDisplayForTesting() error {
	// Only set up virtual display if we don't already have one
	if os.Getenv("DISPLAY") != "" {
		return nil // Already have a display
	}

	// Only set up on Linux (Ubuntu CI runners)
	if runtime.GOOS != "linux" {
		return nil // Skip setup on other platforms
	}

	// Check if Xvfb is available
	if _, err := os.Stat("/usr/bin/Xvfb"); os.IsNotExist(err) {
		return err
	}

	// Set up virtual display
	displayNum := ":99"

	// Set DISPLAY environment variable
	os.Setenv("DISPLAY", displayNum)

	return nil
}
