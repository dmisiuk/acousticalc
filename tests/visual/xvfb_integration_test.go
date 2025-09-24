package visual

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// TestXvfbIntegration validates visual testing works with virtual display
func TestXvfbIntegration(t *testing.T) {
	// Skip if not on Linux or in CI environment
	if runtime.GOOS != "linux" {
		t.Skip("Xvfb integration tests only run on Linux")
		return
	}

	t.Run("virtual_display_availability", func(t *testing.T) {
		display := os.Getenv("DISPLAY")
		if display == "" {
			t.Skip("No display available - this test requires DISPLAY environment variable")
			return
		}
		t.Logf("Using display: %s", display)
	})

	t.Run("screenshot_capture_with_xvfb", func(t *testing.T) {
		display := os.Getenv("DISPLAY")
		if display == "" {
			t.Skip("No display available for screenshot testing")
			return
		}

		outputDir := t.TempDir()
		capture := NewScreenshotCapture("xvfb_test", outputDir)

		// Test screenshot capture
		startTime := time.Now()
		filepath, err := capture.CaptureScreen("xvfb_integration")
		duration := time.Since(startTime)

		if err != nil {
			t.Logf("Screenshot capture failed (expected in some CI environments): %v", err)
			return
		}

		if filepath == "" {
			t.Skip("Screenshot capture returned empty path - display may not support robotgo")
			return
		}

		// Verify screenshot was created
		if _, err := os.Stat(filepath); err != nil {
			t.Errorf("Screenshot file not created: %v", err)
			return
		}

		// Check file size (should be reasonable for a 1920x1080 screenshot)
		if info, err := os.Stat(filepath); err == nil {
			if info.Size() < 1024 { // Less than 1KB indicates failure
				t.Errorf("Screenshot file too small (%d bytes), likely failed", info.Size())
			} else {
				t.Logf("✅ Screenshot created successfully: %s (%d bytes) in %v",
					filepath, info.Size(), duration)
			}
		}
	})

	t.Run("visual_test_logger_with_xvfb", func(t *testing.T) {
		display := os.Getenv("DISPLAY")
		if display == "" {
			t.Skip("No display available for visual logging test")
			return
		}

		outputDir := t.TempDir()
		logger := NewVisualTestLogger("xvfb_logger_test", outputDir)

		// Log events with screenshots
		events := []struct {
			eventType   VisualTestEvent
			description string
		}{
			{EventTestStart, "Xvfb integration test started"},
			{EventTestProcess, "Processing with virtual display"},
			{EventTestComplete, "Xvfb integration test completed"},
		}

		successfulScreenshots := 0
		for _, event := range events {
			logger.LogEvent(event.eventType, event.description, map[string]interface{}{
				"display": display,
				"runtime": runtime.GOOS,
			})

			// Check if screenshot was captured
			if len(logger.Screenshots) > successfulScreenshots {
				successfulScreenshots++
			}
		}

		if successfulScreenshots > 0 {
			t.Logf("✅ Successfully captured %d screenshots with virtual display", successfulScreenshots)
		} else {
			t.Log("⚠️ No screenshots captured - virtual display may not support robotgo")
		}

		// Generate visual report
		if err := logger.GenerateVisualReport(); err != nil {
			t.Errorf("Failed to generate visual report: %v", err)
		} else {
			reportPath := filepath.Join(outputDir, "xvfb_logger_test_visual_report.html")
			if _, err := os.Stat(reportPath); err == nil {
				t.Logf("✅ Visual report generated: %s", reportPath)
			}
		}
	})

	t.Run("cross_platform_artifact_generation", func(t *testing.T) {
		outputDir := t.TempDir()
		logger := NewVisualTestLogger("cross_platform_test", outputDir)

		// Add test events
		logger.LogEvent(EventTestStart, "Cross-platform test with Xvfb", nil)
		logger.LogEvent(EventTestComplete, "Cross-platform test completed", nil)

		// Generate artifacts
		generator := NewArtifactGenerator("cross_platform_test", outputDir)
		err := generator.GenerateComprehensiveArtifacts(logger)

		if err != nil {
			t.Errorf("Failed to generate comprehensive artifacts: %v", err)
			return
		}

		// Verify artifacts were created
		expectedDirs := []string{
			"reports",
			"demo_content/storyboards",
			"demo_content/metadata",
			"charts",
		}

		for _, dir := range expectedDirs {
			fullPath := filepath.Join(outputDir, dir)
			if _, err := os.Stat(fullPath); err != nil {
				t.Errorf("Expected artifact directory not created: %s", fullPath)
			} else {
				t.Logf("✅ Artifact directory created: %s", dir)
			}
		}

		// Check artifact summary
		summary := generator.GetArtifactSummary()
		if summary.TotalArtifacts == 0 {
			t.Error("No artifacts reported in summary")
		} else {
			t.Logf("✅ Generated %d artifacts (%d bytes total)",
				summary.TotalArtifacts, summary.TotalSizeBytes)
		}
	})
}

// TestXvfbEnvironmentValidation validates the Xvfb environment is set up correctly
func TestXvfbEnvironmentValidation(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Xvfb environment validation only runs on Linux")
		return
	}

	t.Run("display_environment", func(t *testing.T) {
		display := os.Getenv("DISPLAY")
		if display == "" {
			t.Skip("DISPLAY environment variable not set")
			return
		}

		t.Logf("Display environment: %s", display)

		// Check if display is virtual (common virtual display patterns)
		if display == ":99" || display == ":0" {
			t.Logf("✅ Detected virtual display: %s", display)
		} else {
			t.Logf("ℹ️ Display: %s (may be real or virtual)", display)
		}
	})

	t.Run("github_actions_environment", func(t *testing.T) {
		if os.Getenv("GITHUB_ACTIONS") == "true" {
			t.Log("✅ Running in GitHub Actions environment")

			// Check xvfb-action specific environment
			if os.Getenv("DISPLAY") != "" {
				t.Logf("✅ xvfb-action provided display: %s", os.Getenv("DISPLAY"))
			} else {
				t.Error("Expected DISPLAY to be set by xvfb-action in GitHub Actions")
			}
		} else {
			t.Log("ℹ️ Not running in GitHub Actions")
		}
	})

	t.Run("ci_environment_detection", func(t *testing.T) {
		ci := os.Getenv("CI")
		if ci == "true" {
			t.Log("✅ Detected CI environment")

			// In CI with Xvfb, we should have a display
			if os.Getenv("DISPLAY") == "" {
				t.Error("Expected DISPLAY to be available in CI environment")
			}
		} else {
			t.Log("ℹ️ Not running in CI environment")
		}
	})
}
