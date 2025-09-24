package visual

import (
	"os"
	"testing"
	"time"
)



// TestScreenshotCapture tests the screenshot capture functionality
func TestScreenshotCapture(t *testing.T) {
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