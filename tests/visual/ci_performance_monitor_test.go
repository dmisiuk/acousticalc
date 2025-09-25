package visual

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCIPerformanceMonitor(t *testing.T) {
	t.Run("basic_monitoring", func(t *testing.T) {
		monitor := NewCIPerformanceMonitor()

		// Verify initial state
		if monitor.Platform == "" {
			t.Error("Platform should be set")
		}
		if monitor.Thresholds.MaxCIOverhead != 30*time.Second {
			t.Errorf("Expected 30s threshold, got %v", monitor.Thresholds.MaxCIOverhead)
		}

		// Test monitoring lifecycle
		monitor.Start()
		if monitor.StartTime.IsZero() {
			t.Error("Start time should be set")
		}

		// Simulate some work
		time.Sleep(10 * time.Millisecond)
		monitor.RecordScreenshot(500 * time.Millisecond)
		monitor.RecordArtifact("html", 100*time.Millisecond)

		err := monitor.Finish()
		if err != nil {
			t.Fatalf("Finish should not error for fast operations: %v", err)
		}

		if monitor.EndTime.IsZero() {
			t.Error("End time should be set")
		}
		if monitor.TotalDuration == 0 {
			t.Error("Total duration should be > 0")
		}
		if monitor.ScreenshotCount != 1 {
			t.Errorf("Expected 1 screenshot, got %d", monitor.ScreenshotCount)
		}
		if monitor.ArtifactCount != 1 {
			t.Errorf("Expected 1 artifact, got %d", monitor.ArtifactCount)
		}
	})

	t.Run("threshold_validation", func(t *testing.T) {
		monitor := NewCIPerformanceMonitor()
		monitor.Thresholds.MaxCIOverhead = 50 * time.Millisecond // Very short for testing

		monitor.Start()
		time.Sleep(100 * time.Millisecond) // Exceed threshold
		err := monitor.Finish()

		if err == nil {
			t.Error("Should have failed threshold validation")
		}
		if monitor.GetStatus() != "FAIL" {
			t.Error("Status should be FAIL")
		}
	})

	t.Run("screenshot_performance_tracking", func(t *testing.T) {
		monitor := NewCIPerformanceMonitor()

		// Record multiple screenshots with different durations
		durations := []time.Duration{
			400 * time.Millisecond,
			600 * time.Millisecond,
			300 * time.Millisecond,
		}

		for _, duration := range durations {
			monitor.RecordScreenshot(duration)
		}

		if monitor.ScreenshotCount != 3 {
			t.Errorf("Expected 3 screenshots, got %d", monitor.ScreenshotCount)
		}

		// Check metrics
		for i, expected := range durations {
			key := func(i int) string { return fmt.Sprintf("screenshot_%d_ms", i+1) }(i)
			if actual, exists := monitor.Metrics[key]; !exists {
				t.Errorf("Missing metric %s", key)
			} else if actual != float64(expected.Milliseconds()) {
				t.Errorf("Expected %v ms, got %v ms", expected.Milliseconds(), actual)
			}
		}
	})

	t.Run("save_report", func(t *testing.T) {
		// Create temporary directory
		tempDir := filepath.Join(os.TempDir(), "test_ci_performance")
		defer os.RemoveAll(tempDir)

		monitor := NewCIPerformanceMonitor()
		monitor.Start()
		monitor.RecordScreenshot(500 * time.Millisecond)
		monitor.RecordArtifact("demo", 200*time.Millisecond)
		if err := monitor.Finish(); err != nil {
			t.Logf("Warning: failed to finish monitor: %v", err)
		}

		err := monitor.SaveReport(tempDir)
		if err != nil {
			t.Fatalf("Failed to save report: %v", err)
		}

		// Check that files were created
		summaryFile := filepath.Join(tempDir, "ci_performance_summary.txt")
		if _, err := os.Stat(summaryFile); os.IsNotExist(err) {
			t.Error("Summary file was not created")
		}

		// Check JSON report exists
		files, err := filepath.Glob(filepath.Join(tempDir, "ci_performance_*.json"))
		if err != nil || len(files) == 0 {
			t.Error("JSON report file was not created")
		}

		// Verify JSON content
		if len(files) > 0 {
			data, err := os.ReadFile(files[0])
			if err != nil {
				t.Fatalf("Failed to read report file: %v", err)
			}

			var report CIPerformanceMonitor
			if err := json.Unmarshal(data, &report); err != nil {
				t.Fatalf("Failed to parse JSON report: %v", err)
			}

			if report.ScreenshotCount != 1 {
				t.Errorf("Expected 1 screenshot in report, got %d", report.ScreenshotCount)
			}
			if report.ArtifactCount != 1 {
				t.Errorf("Expected 1 artifact in report, got %d", report.ArtifactCount)
			}
		}
	})

	t.Run("monitor_with_context", func(t *testing.T) {
		tempDir := filepath.Join(os.TempDir(), "test_monitor_context")
		defer os.RemoveAll(tempDir)

		ctx := context.Background()
		executed := false

		err := MonitorWithContext(ctx, func() error {
			executed = true
			time.Sleep(10 * time.Millisecond)
			return nil
		}, tempDir)

		if err != nil {
			t.Fatalf("MonitorWithContext failed: %v", err)
		}
		if !executed {
			t.Error("Function was not executed")
		}

		// Check that report was saved
		summaryFile := filepath.Join(tempDir, "ci_performance_summary.txt")
		if _, err := os.Stat(summaryFile); os.IsNotExist(err) {
			t.Error("Performance report was not saved")
		}
	})
}

func TestCIPerformanceThresholds(t *testing.T) {
	t.Run("default_thresholds", func(t *testing.T) {
		monitor := NewCIPerformanceMonitor()

		expected := PerformanceThreshold{
			MaxCIOverhead:     30 * time.Second,
			MaxScreenshotTime: 5 * time.Second,
			MaxArtifactTime:   10 * time.Second,
		}

		if monitor.Thresholds != expected {
			t.Errorf("Default thresholds mismatch: got %+v, want %+v",
				monitor.Thresholds, expected)
		}
	})

	t.Run("screenshot_threshold_validation", func(t *testing.T) {
		monitor := NewCIPerformanceMonitor()
		monitor.Thresholds.MaxScreenshotTime = 100 * time.Millisecond

		monitor.Start()
		monitor.RecordScreenshot(200 * time.Millisecond) // Exceeds threshold
		err := monitor.Finish()

		if err == nil {
			t.Error("Should have failed screenshot threshold validation")
		}
	})
}

func TestCIDetection(t *testing.T) {
	t.Run("detect_github_actions", func(t *testing.T) {
		// Save original value
		original := os.Getenv("GITHUB_ACTIONS")
		defer os.Setenv("GITHUB_ACTIONS", original)

		// Test with GITHUB_ACTIONS set
		os.Setenv("GITHUB_ACTIONS", "true")
		if !isCI() {
			t.Error("Should detect GitHub Actions as CI")
		}

		// Test without CI env vars
		os.Unsetenv("GITHUB_ACTIONS")
		// Note: This might still return true if other CI vars are set
		// In a real environment, but for testing we assume clean state
	})

	t.Run("monitor_sets_ci_flag", func(t *testing.T) {
		monitor := NewCIPerformanceMonitor()
		// The IsCI flag should be set based on environment
		// We can't easily test this without manipulating all CI env vars
		_ = monitor.IsCI // Just verify the field exists
	})
}

func BenchmarkCIPerformanceMonitor(b *testing.B) {
	b.Run("monitor_overhead", func(b *testing.B) {
		tempDir := filepath.Join(os.TempDir(), "bench_ci_performance")
		defer os.RemoveAll(tempDir)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			monitor := NewCIPerformanceMonitor()
			monitor.Start()
			monitor.RecordScreenshot(100 * time.Millisecond)
			monitor.RecordArtifact("test", 50*time.Millisecond)
			if err := monitor.Finish(); err != nil {
				b.Logf("Warning: failed to finish monitor: %v", err)
			}
			if err := monitor.SaveReport(tempDir); err != nil {
				b.Logf("Warning: failed to save report: %v", err)
			}
		}
	})
}
