package visual

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
)

// TestScreenshotCapture tests the screenshot capture functionality
func TestScreenshotCapture(t *testing.T) {
	// Set up virtual display for CI environments
	if err := setupVirtualDisplay(); err != nil {
		t.Logf("Warning: Could not setup virtual display: %v", err)
		// Continue with test - robotgo might still work in some environments
	}

	outputDir := "../artifacts/screenshots/unit"

	capture := NewScreenshotCapture("screenshot_test", outputDir)

	tests := []struct {
		name      string
		eventType string
	}{
		{"test_start", "start"},
		{"test_process", "process"},
		{"test_end", "end"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Small delay to ensure different timestamps
			time.Sleep(10 * time.Millisecond)

			filepath := capture.CaptureTestEvent(t, tt.eventType)
			if filepath == "" {
				t.Skip("Screenshot capture not available in test environment")
				return
			}

			// Verify file was created
			if _, err := os.Stat(filepath); os.IsNotExist(err) {
				t.Errorf("Screenshot file was not created: %s", filepath)
			}

			// Verify file has reasonable size (>1KB indicates actual content)
			if info, err := os.Stat(filepath); err == nil {
				if info.Size() < 1024 {
					t.Errorf("Screenshot file appears to be empty or corrupted: %d bytes", info.Size())
				} else {
					t.Logf("Screenshot successfully created: %s (%d bytes)", filepath, info.Size())
				}
			}
		})
	}
}

// TestCrossPlatformScreenshot tests screenshot capture across platforms
func TestCrossPlatformScreenshot(t *testing.T) {
	// Set up virtual display for CI environments
	if err := setupVirtualDisplay(); err != nil {
		t.Logf("Warning: Could not setup virtual display: %v", err)
	}

	outputDir := "../artifacts/screenshots/unit"
	capture := NewScreenshotCapture("cross_platform_test", outputDir)

	// Test basic screenshot capture
	filepath := capture.CaptureTestEvent(t, "cross_platform")
	if filepath == "" {
		t.Skip("Screenshot capture not available in test environment")
		return
	}

	// Verify PNG format
	if filepath := filepath; filepath[len(filepath)-4:] != ".png" {
		t.Errorf("Expected PNG format, got: %s", filepath)
	}

	t.Logf("Cross-platform screenshot test completed: %s", filepath)
}

// BenchmarkScreenshotCapture benchmarks screenshot performance
func BenchmarkScreenshotCapture(b *testing.B) {
	// Set up virtual display for CI environments
	if err := setupVirtualDisplay(); err != nil {
		b.Logf("Warning: Could not setup virtual display: %v", err)
	}

	outputDir := "../artifacts/screenshots/unit"
	capture := NewScreenshotCapture("benchmark_test", outputDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := capture.CaptureScreen("benchmark")
		if err != nil {
			b.Skipf("Screenshot capture not available: %v", err)
		}
	}
}

// setupVirtualDisplay configures a virtual display for headless CI environments
// Note: In GitHub Actions, xvfb-action handles the Xvfb setup automatically
func setupVirtualDisplay() error {
	// Check if we already have a display (either real or virtual)
	if os.Getenv("DISPLAY") != "" {
		return nil // Display is available (real or virtual)
	}

	// Only set up on Linux (Ubuntu CI runners)
	if runtime.GOOS != "linux" {
		return fmt.Errorf("virtual display only supported on Linux, current OS: %s", runtime.GOOS)
	}

	// In GitHub Actions with xvfb-action, DISPLAY will be set automatically
	// For local testing, try to detect and use existing Xvfb or start one
	if os.Getenv("CI") != "" {
		// Running in CI - assume xvfb-action will provide display
		return fmt.Errorf("no display available in CI environment")
	}

	// Local development - try to start Xvfb
	if _, err := exec.LookPath("Xvfb"); err != nil {
		return fmt.Errorf("Xvfb not available for local testing: %w", err)
	}

	displayNum := ":99"
	cmd := exec.Command("Xvfb", displayNum, "-screen", "0", "1920x1080x24", "-ac")
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start local Xvfb: %w", err)
	}

	os.Setenv("DISPLAY", displayNum)
	time.Sleep(2 * time.Second)

	return nil
}
