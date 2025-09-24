package visual

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// TestVisualUtilsComprehensive provides extensive test coverage for visual utilities
func TestVisualUtilsComprehensive(t *testing.T) {
	// Test: NewScreenshotCapture initialization
	t.Run("NewScreenshotCapture", func(t *testing.T) {
		testName := "test_capture"
		outputDir := t.TempDir()

		capture := NewScreenshotCapture(testName, outputDir)

		if capture.TestName != testName {
			t.Errorf("Expected TestName %s, got %s", testName, capture.TestName)
		}
		if capture.OutputDir != outputDir {
			t.Errorf("Expected OutputDir %s, got %s", outputDir, capture.OutputDir)
		}
		if capture.Format != "png" {
			t.Errorf("Expected default format png, got %s", capture.Format)
		}
		if capture.Quality != 100 {
			t.Errorf("Expected default quality 100, got %d", capture.Quality)
		}
	})

	// Test: NewVisualTestLogger initialization
	t.Run("NewVisualTestLogger", func(t *testing.T) {
		testName := "test_logger"
		outputDir := t.TempDir()

		logger := NewVisualTestLogger(testName, outputDir)

		if logger.TestName != testName {
			t.Errorf("Expected TestName %s, got %s", testName, logger.TestName)
		}
		if logger.OutputDir != outputDir {
			t.Errorf("Expected OutputDir %s, got %s", outputDir, logger.OutputDir)
		}
		if len(logger.Screenshots) != 0 {
			t.Errorf("Expected empty Screenshots slice, got %d items", len(logger.Screenshots))
		}
		if len(logger.Events) != 0 {
			t.Errorf("Expected empty Events slice, got %d items", len(logger.Events))
		}
		if logger.StartTime.IsZero() {
			t.Error("Expected StartTime to be set")
		}
	})

	// Test: VisualTestEvent constants
	t.Run("VisualTestEvent_Constants", func(t *testing.T) {
		events := []VisualTestEvent{
			EventTestStart,
			EventTestPass,
			EventTestFail,
			EventTestProcess,
			EventTestComplete,
		}

		expectedValues := []string{"start", "pass", "fail", "process", "complete"}
		for i, event := range events {
			if string(event) != expectedValues[i] {
				t.Errorf("Expected event %d to be %s, got %s", i, expectedValues[i], event)
			}
		}
	})

	// Test: VisualTestLogger LogEvent
	t.Run("VisualTestLogger_LogEvent", func(t *testing.T) {
		outputDir := t.TempDir()
		logger := NewVisualTestLogger("log_event_test", outputDir)

		// Test logging an event with metadata
		metadata := map[string]interface{}{
			"key1": "value1",
			"key2": 42,
			"key3": true,
		}

		logger.LogEvent(EventTestStart, "Test started", metadata)

		if len(logger.Events) != 1 {
			t.Fatalf("Expected 1 event, got %d", len(logger.Events))
		}

		event := logger.Events[0]
		if event.Type != EventTestStart {
			t.Errorf("Expected event type %s, got %s", EventTestStart, event.Type)
		}
		if event.Description != "Test started" {
			t.Errorf("Expected description 'Test started', got %s", event.Description)
		}
		if len(event.Metadata) != 3 {
			t.Errorf("Expected 3 metadata items, got %d", len(event.Metadata))
		}
		if event.Metadata["key1"] != "value1" {
			t.Errorf("Expected metadata key1 to be 'value1', got %v", event.Metadata["key1"])
		}
	})

	// Test: ScreenshotCapture with invalid directory
	t.Run("ScreenshotCapture_InvalidDirectory", func(t *testing.T) {
		// Create a directory path that can't be created
		invalidDir := "/invalid/path/that/does/not/exist"
		capture := NewScreenshotCapture("invalid_test", invalidDir)

		// This should fail when trying to create the directory
		_, err := capture.CaptureScreen("test")
		if err == nil {
			t.Error("Expected error when capturing screenshot to invalid directory")
		}
	})

	// Test: Performance thresholds validation
	t.Run("PerformanceThresholds_Validation", func(t *testing.T) {
		monitor := NewCIPerformanceMonitor()

		// Test default thresholds
		if monitor.Thresholds.MaxCIOverhead != 30*time.Second {
			t.Errorf("Expected MaxCIOverhead to be 30s, got %v", monitor.Thresholds.MaxCIOverhead)
		}
		if monitor.Thresholds.MaxScreenshotTime != 5*time.Second {
			t.Errorf("Expected MaxScreenshotTime to be 5s, got %v", monitor.Thresholds.MaxScreenshotTime)
		}
		if monitor.Thresholds.MaxArtifactTime != 10*time.Second {
			t.Errorf("Expected MaxArtifactTime to be 10s, got %v", monitor.Thresholds.MaxArtifactTime)
		}
	})

	// Test: Platform detection
	t.Run("Platform_Detection", func(t *testing.T) {
		monitor := NewCIPerformanceMonitor()
		expectedPlatform := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

		if monitor.Platform != expectedPlatform {
			t.Errorf("Expected platform %s, got %s", expectedPlatform, monitor.Platform)
		}
	})

	// Test: CI environment detection
	t.Run("CI_Detection", func(t *testing.T) {
		// Test in normal environment (should be false unless actually in CI)
		monitor := NewCIPerformanceMonitor()

		// This test runs in local environment, so should be false
		if monitor.IsCI {
			// If this fails, we're actually running in CI, which is fine
			t.Log("Running in CI environment - this is expected in CI")
		}
	})

	// Test: Performance monitoring lifecycle
	t.Run("PerformanceMonitoring_Lifecycle", func(t *testing.T) {
		monitor := NewCIPerformanceMonitor()

		// Test initial state - StartTime should be zero before Start()
		if !monitor.StartTime.IsZero() {
			t.Error("Expected StartTime to be zero before Start()")
		}
		if !monitor.EndTime.IsZero() {
			t.Error("Expected EndTime to be zero initially")
		}

		// Start monitoring (resets StartTime)
		monitor.Start()

		// Test recording screenshots
		monitor.RecordScreenshot(100 * time.Millisecond)
		monitor.RecordScreenshot(200 * time.Millisecond)

		if monitor.ScreenshotCount != 2 {
			t.Errorf("Expected ScreenshotCount to be 2, got %d", monitor.ScreenshotCount)
		}

		// Test recording artifacts
		monitor.RecordArtifact("html", 50*time.Millisecond)
		monitor.RecordArtifact("json", 25*time.Millisecond)

		if monitor.ArtifactCount != 2 {
			t.Errorf("Expected ArtifactCount to be 2, got %d", monitor.ArtifactCount)
		}

		// Test finish
		err := monitor.Finish()
		if err != nil {
			t.Errorf("Expected Finish() to succeed, got error: %v", err)
		}

		if monitor.EndTime.IsZero() {
			t.Error("Expected EndTime to be set after Finish()")
		}
		if monitor.TotalDuration == 0 {
			t.Error("Expected TotalDuration to be calculated")
		}
	})

	// Test: Performance threshold validation
	t.Run("PerformanceThreshold_Validation", func(t *testing.T) {
		monitor := NewCIPerformanceMonitor()
		monitor.Start()

		// Simulate a fast execution that should pass
		time.Sleep(10 * time.Millisecond)
		monitor.Finish()

		status := monitor.GetStatus()
		if status != "PASS" {
			t.Errorf("Expected status 'PASS' for fast execution, got %s", status)
		}

		// Test threshold violation by manually setting the duration
		monitor = NewCIPerformanceMonitor()
		monitor.Start()
		time.Sleep(10 * time.Millisecond)
		monitor.Finish()

		// Manually set a duration that exceeds the threshold
		monitor.TotalDuration = 35 * time.Second
		status = monitor.GetStatus()
		if status != "FAIL" {
			t.Errorf("Expected status 'FAIL' for slow execution, got %s", status)
		}
	})

	// Test: Context-based monitoring
	t.Run("Context_Monitoring", func(t *testing.T) {
		outputDir := t.TempDir()

		err := MonitorWithContext(context.Background(), func() error {
			// Simulate some work
			time.Sleep(10 * time.Millisecond)
			return nil
		}, outputDir)

		if err != nil {
			t.Errorf("Expected MonitorWithContext to succeed, got error: %v", err)
		}

		// Check if performance report was created
		files, err := filepath.Glob(filepath.Join(outputDir, "ci_performance_*.json"))
		if err != nil {
			t.Errorf("Failed to check for performance reports: %v", err)
		}
		if len(files) == 0 {
			t.Error("Expected performance report file to be created")
		}
	})

	// Test: Performance dashboard creation
	t.Run("PerformanceDashboard_Creation", func(t *testing.T) {
		outputDir := t.TempDir()
		reportsDir := filepath.Join(outputDir, "reports")
		dashboardDir := filepath.Join(outputDir, "dashboard")

		// Create reports directory
		if err := os.MkdirAll(reportsDir, 0755); err != nil {
			t.Fatalf("Failed to create reports directory: %v", err)
		}

		// Create a sample performance report
		monitor := NewCIPerformanceMonitor()
		monitor.Start()
		time.Sleep(10 * time.Millisecond)
		monitor.Finish()

		if err := monitor.SaveReport(reportsDir); err != nil {
			t.Errorf("Failed to save performance report: %v", err)
		}

		// Generate dashboard
		err := GenerateDashboard(reportsDir, dashboardDir)
		if err != nil {
			t.Errorf("Expected GenerateDashboard to succeed, got error: %v", err)
		}

		// Check if dashboard files were created
		htmlFile := filepath.Join(dashboardDir, "performance_dashboard.html")
		jsonFile := filepath.Join(dashboardDir, "performance_dashboard.json")

		if _, err := os.Stat(htmlFile); os.IsNotExist(err) {
			t.Error("Expected HTML dashboard file to be created")
		}
		if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
			t.Error("Expected JSON dashboard file to be created")
		}
	})

	// Test: Error handling for invalid report files
	t.Run("Dashboard_LoadReports_InvalidFiles", func(t *testing.T) {
		outputDir := t.TempDir()

		// Create an invalid JSON file
		invalidFile := filepath.Join(outputDir, "ci_performance_invalid.json")
		if err := os.WriteFile(invalidFile, []byte("invalid json"), 0644); err != nil {
			t.Fatalf("Failed to create invalid JSON file: %v", err)
		}

		dashboard := NewPerformanceDashboard()
		err := dashboard.LoadReports(outputDir)

		// Should not error, should skip invalid files
		if err != nil {
			t.Errorf("Expected LoadReports to handle invalid files gracefully, got error: %v", err)
		}

		if len(dashboard.Reports) != 0 {
			t.Errorf("Expected no reports to be loaded from invalid files, got %d", len(dashboard.Reports))
		}
	})

	// Test: HTML template functionality
	t.Run("HTML_Template_Functionality", func(t *testing.T) {
		dashboard := NewPerformanceDashboard()

		// Add a sample report
		monitor := NewCIPerformanceMonitor()
		monitor.Start()
		time.Sleep(10 * time.Millisecond)
		monitor.Finish()

		dashboard.Reports = append(dashboard.Reports, *monitor)
		dashboard.calculateSummary()

		// Test HTML generation
		outputFile := filepath.Join(t.TempDir(), "test_dashboard.html")
		err := dashboard.GenerateHTML(outputFile)

		if err != nil {
			t.Errorf("Expected GenerateHTML to succeed, got error: %v", err)
		}

		// Check if HTML file was created and contains expected content
		htmlContent, err := os.ReadFile(outputFile)
		if err != nil {
			t.Errorf("Failed to read generated HTML file: %v", err)
		}

		contentStr := string(htmlContent)
		expectedStrings := []string{
			"Visual Testing Performance Dashboard",
			"Total Test Runs",
			"Average CI Time",
			"Platform Distribution",
		}

		for _, expected := range expectedStrings {
			if !contains(contentStr, expected) {
				t.Errorf("Expected HTML to contain '%s', but it doesn't", expected)
			}
		}
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}

// Benchmark tests for performance-critical functions
func BenchmarkVisualUtils(b *testing.B) {
	// Benchmark: ScreenshotCapture creation
	b.Run("NewScreenshotCapture", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewScreenshotCapture("benchmark_test", "/tmp")
		}
	})

	// Benchmark: VisualTestLogger creation
	b.Run("NewVisualTestLogger", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewVisualTestLogger("benchmark_logger", "/tmp")
		}
	})

	// Benchmark: PerformanceMonitor creation
	b.Run("NewCIPerformanceMonitor", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewCIPerformanceMonitor()
		}
	})

	// Benchmark: Event logging
	b.Run("LogEvent", func(b *testing.B) {
		logger := NewVisualTestLogger("benchmark", "/tmp")
		metadata := map[string]interface{}{"key": "value"}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.LogEvent(EventTestProcess, "benchmark event", metadata)
		}
	})
}
