package visual

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// PlatformCompatibility tests visual testing across different platforms
type PlatformCompatibility struct {
	Platform     string
	Architecture string
	OutputDir    string
}

// NewPlatformCompatibility creates a new platform compatibility tester
func NewPlatformCompatibility(outputDir string) *PlatformCompatibility {
	return &PlatformCompatibility{
		Platform:     runtime.GOOS,
		Architecture: runtime.GOARCH,
		OutputDir:    outputDir,
	}
}

// TestPlatformScreenshotCompatibility tests screenshot capture across platforms
func TestPlatformScreenshotCompatibility(t *testing.T) {
	outputDir := "../artifacts/screenshots/unit"
	compat := NewPlatformCompatibility(outputDir)

	t.Logf("Testing on platform: %s/%s", compat.Platform, compat.Architecture)

	// Test basic cross-platform functionality
	capture := NewScreenshotCapture("platform_compat", outputDir)

	// Test PNG format consistency
	t.Run("png_format_consistency", func(t *testing.T) {
		filepath := capture.CaptureTestEvent(t, "png_test")
		if filepath == "" {
			t.Skip("Screenshot capture not available")
			return
		}

		// Verify PNG format
		if ext := filepath[len(filepath)-4:]; ext != ".png" {
			t.Errorf("Expected .png extension, got: %s", ext)
		}

		// Verify file exists and has content
		if info, err := os.Stat(filepath); err != nil {
			t.Errorf("Screenshot file error: %v", err)
		} else if info.Size() < 1024 {
			t.Errorf("Screenshot file too small (%d bytes), may be corrupted", info.Size())
		}
	})

	// Test sRGB color space (metadata check)
	t.Run("srgb_color_space", func(t *testing.T) {
		filepath := capture.CaptureTestEvent(t, "srgb_test")
		if filepath == "" {
			t.Skip("Screenshot capture not available")
			return
		}

		// For now, just verify file was created properly
		// Future enhancement: check color profile metadata
		if _, err := os.Stat(filepath); err != nil {
			t.Errorf("sRGB test screenshot failed: %v", err)
		} else {
			t.Logf("sRGB format test completed: %s", filepath)
		}
	})

	// Test filename compatibility across platforms
	t.Run("filename_compatibility", func(t *testing.T) {
		testName := "cross_platform_filename_test"
		timestamp := time.Now()

		// Create filename following our convention
		filename := fmt.Sprintf("%s_compatibility_%s.png",
			testName,
			timestamp.Format("20060102_150405"))

		// Verify filename is valid for current platform
		testPath := filepath.Join(outputDir, filename)

		// Test path creation
		if err := os.MkdirAll(filepath.Dir(testPath), 0755); err != nil {
			t.Errorf("Failed to create directory path: %v", err)
		}

		// Test file creation (mock)
		if file, err := os.Create(testPath); err != nil {
			t.Errorf("Failed to create test file: %v", err)
		} else {
			file.Close()
			os.Remove(testPath) // Cleanup
			t.Logf("Filename compatibility verified: %s", filename)
		}
	})
}

// TestVisualArtifactCompatibility tests artifact generation compatibility
func TestVisualArtifactCompatibility(t *testing.T) {
	outputDir := "../artifacts/screenshots/unit"
	logger := NewVisualTestLogger("artifact_compat_test", outputDir)

	// Log test events
	logger.LogEvent(EventTestStart, "Starting cross-platform artifact test", map[string]interface{}{
		"platform":     runtime.GOOS,
		"architecture": runtime.GOARCH,
	})

	logger.LogEvent(EventTestProcess, "Processing visual artifacts", map[string]interface{}{
		"test_type": "compatibility",
	})

	logger.LogEvent(EventTestComplete, "Completed cross-platform test", map[string]interface{}{
		"status": "success",
	})

	// Test HTML report generation
	t.Run("html_report_generation", func(t *testing.T) {
		err := logger.GenerateVisualReport()
		if err != nil {
			t.Errorf("Failed to generate visual report: %v", err)
		} else {
			t.Logf("Visual report generated successfully")
		}
	})

	// Test demo storyboard generation
	t.Run("demo_storyboard_generation", func(t *testing.T) {
		err := logger.CreateDemoStoryboard()
		if err != nil {
			t.Errorf("Failed to create demo storyboard: %v", err)
		} else {
			t.Logf("Demo storyboard generated successfully")
		}
	})
}

// TestDirectoryStructureCompatibility tests directory creation across platforms
func TestDirectoryStructureCompatibility(t *testing.T) {
	baseDir := "../artifacts"

	// Test required directory structure
	requiredDirs := []string{
		"screenshots/unit",
		"screenshots/integration",
		"screenshots/e2e",
		"demo_content/storyboards",
		"demo_content/assets",
		"demo_content/metadata",
		"baselines/ui",
		"baselines/terminal",
		"baselines/charts",
	}

	for _, dir := range requiredDirs {
		t.Run(fmt.Sprintf("directory_%s", dir), func(t *testing.T) {
			fullPath := filepath.Join(baseDir, dir)

			// Test directory creation
			err := os.MkdirAll(fullPath, 0755)
			if err != nil {
				t.Errorf("Failed to create directory %s: %v", fullPath, err)
				return
			}

			// Test directory access
			if info, err := os.Stat(fullPath); err != nil {
				t.Errorf("Failed to access directory %s: %v", fullPath, err)
			} else if !info.IsDir() {
				t.Errorf("Path %s is not a directory", fullPath)
			} else {
				t.Logf("Directory structure verified: %s", fullPath)
			}
		})
	}
}

// TestPerformanceConstraints tests that visual testing meets performance requirements
func TestPerformanceConstraints(t *testing.T) {
	outputDir := "../artifacts/screenshots/unit"

	// Performance requirement: <30 seconds additional CI time
	// This translates to very fast screenshot capture
	maxDuration := 5 * time.Second // Per operation

	t.Run("screenshot_performance", func(t *testing.T) {
		capture := NewScreenshotCapture("performance_test", outputDir)

		start := time.Now()
		filepath := capture.CaptureTestEvent(t, "performance")
		duration := time.Since(start)

		if filepath == "" {
			t.Skip("Screenshot capture not available")
			return
		}

		if duration > maxDuration {
			t.Errorf("Screenshot capture took too long: %v (max: %v)", duration, maxDuration)
		} else {
			t.Logf("Screenshot performance acceptable: %v", duration)
		}
	})

	t.Run("visual_report_performance", func(t *testing.T) {
		logger := NewVisualTestLogger("performance_report_test", outputDir)

		// Add some test events
		for i := 0; i < 5; i++ {
			logger.LogEvent(EventTestProcess, fmt.Sprintf("Performance test event %d", i), nil)
		}

		start := time.Now()
		err := logger.GenerateVisualReport()
		duration := time.Since(start)

		if err != nil {
			t.Errorf("Visual report generation failed: %v", err)
			return
		}

		if duration > maxDuration {
			t.Errorf("Visual report generation took too long: %v (max: %v)", duration, maxDuration)
		} else {
			t.Logf("Visual report performance acceptable: %v", duration)
		}
	})
}

// BenchmarkCrossPlatformOperations benchmarks key operations
func BenchmarkCrossPlatformOperations(b *testing.B) {
	outputDir := "../artifacts/screenshots/unit"

	b.Run("ScreenshotCapture", func(b *testing.B) {
		capture := NewScreenshotCapture("benchmark", outputDir)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := capture.CaptureScreen("benchmark")
			if err != nil {
				b.Skipf("Screenshot capture not available: %v", err)
			}
		}
	})

	b.Run("VisualReportGeneration", func(b *testing.B) {
		logger := NewVisualTestLogger("benchmark_report", outputDir)
		logger.LogEvent(EventTestStart, "Benchmark test", nil)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			err := logger.GenerateVisualReport()
			if err != nil {
				b.Errorf("Visual report generation failed: %v", err)
			}
		}
	})
}
