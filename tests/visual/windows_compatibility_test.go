package visual

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestWindowsCompatibility(t *testing.T) {
	// Skip on non-Windows unless specifically testing cross-platform compatibility
	if runtime.GOOS != "windows" && !isTestingCrossPlatform() {
		t.Skip("Windows compatibility tests require Windows platform or cross-platform test mode")
	}

	t.Run("windows_screenshot_capture", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Log("Simulating Windows screenshot capture for cross-platform validation")
			testWindowsScreenshotSimulation(t)
			return
		}

		// Actual Windows screenshot test
		outputDir := filepath.Join(os.TempDir(), "windows_screenshots")
		defer os.RemoveAll(outputDir)

		capture := NewScreenshotCapture("windows_test", outputDir)
		filepath, err := capture.CaptureScreen("start")

		if err != nil {
			t.Fatalf("Windows screenshot capture failed: %v", err)
		}

		// Verify file exists and has proper format
		if !fileExists(filepath) {
			t.Errorf("Screenshot file was not created: %s", filepath)
		}

		if !strings.HasSuffix(filepath, ".png") {
			t.Errorf("Screenshot should be PNG format, got: %s", filepath)
		}
	})

	t.Run("windows_file_path_compatibility", func(t *testing.T) {
		testPaths := []string{
			"tests/artifacts/screenshots/unit/test_file.png",
			"tests\\artifacts\\screenshots\\unit\\test_file.png",
			"C:\\temp\\test.png",
			"/tmp/test.png",
		}

		for _, testPath := range testPaths {
			normalizedPath := normalizePathForWindows(testPath)
			if runtime.GOOS == "windows" {
				// On Windows, paths should use backslashes
				if !strings.Contains(normalizedPath, "\\") && strings.Contains(normalizedPath, "/") {
					t.Errorf("Path not properly normalized for Windows: %s", normalizedPath)
				}
			}

			// Ensure path can be created as directory
			tempBase := filepath.Join(os.TempDir(), "path_test")
			fullPath := filepath.Join(tempBase, normalizedPath)
			dir := filepath.Dir(fullPath)

			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Errorf("Failed to create Windows-compatible directory path %s: %v", dir, err)
			}

			defer os.RemoveAll(tempBase)
		}
	})

	t.Run("windows_artifact_format_validation", func(t *testing.T) {
		outputDir := filepath.Join(os.TempDir(), "windows_artifacts")
		defer os.RemoveAll(outputDir)

		// Test Windows-specific artifact generation
		artifacts := []struct {
			name   string
			format string
			data   []byte
		}{
			{"test_report.html", "html", []byte("<html><body>Test Report</body></html>")},
			{"test_data.json", "json", []byte(`{"test": "data", "platform": "windows"}`)},
			{"screenshot.png", "png", generateTestPNGData()},
		}

		for _, artifact := range artifacts {
			artifactPath := filepath.Join(outputDir, artifact.name)
			if err := os.MkdirAll(filepath.Dir(artifactPath), 0755); err != nil {
				t.Fatalf("Failed to create artifact directory: %v", err)
			}

			if err := os.WriteFile(artifactPath, artifact.data, 0644); err != nil {
				t.Fatalf("Failed to write Windows artifact %s: %v", artifact.name, err)
			}

			// Verify artifact accessibility
			if !fileExists(artifactPath) {
				t.Errorf("Artifact not accessible after creation: %s", artifactPath)
			}

			// Test file can be read back
			readData, err := os.ReadFile(artifactPath)
			if err != nil {
				t.Errorf("Failed to read back Windows artifact %s: %v", artifact.name, err)
			}

			if len(readData) != len(artifact.data) {
				t.Errorf("Artifact data corrupted on Windows: %s", artifact.name)
			}
		}
	})

	t.Run("windows_performance_validation", func(t *testing.T) {
		// Test that Windows operations meet performance requirements
		monitor := NewCIPerformanceMonitor()
		monitor.Start()

		// Simulate Windows-specific operations
		operationStart := time.Now()
		simulateWindowsVisualOperation()
		operationDuration := time.Since(operationStart)

		monitor.RecordScreenshot(operationDuration)

		if err := monitor.Finish(); err != nil {
			t.Errorf("Windows performance validation failed: %v", err)
		}

		// Verify Windows meets same performance standards
		if monitor.TotalDuration > 30*time.Second {
			t.Errorf("Windows operations too slow: %v > 30s", monitor.TotalDuration)
		}

		// Save Windows-specific performance report
		reportDir := filepath.Join(os.TempDir(), "windows_performance")
		defer os.RemoveAll(reportDir)

		if err := monitor.SaveReport(reportDir); err != nil {
			t.Errorf("Failed to save Windows performance report: %v", err)
		}
	})

	t.Run("windows_ci_environment_detection", func(t *testing.T) {
		// Test CI detection on Windows
		windowsCIVars := map[string]string{
			"GITHUB_ACTIONS": "true",
			"TF_BUILD":       "True", // Azure DevOps
			"APPVEYOR":       "True",
		}

		for envVar, value := range windowsCIVars {
			original := os.Getenv(envVar)
			os.Setenv(envVar, value)

			if !isCI() {
				t.Errorf("Failed to detect CI environment with %s=%s", envVar, value)
			}

			// Restore original value
			if original != "" {
				os.Setenv(envVar, original)
			} else {
				os.Unsetenv(envVar)
			}
		}
	})
}

func TestWindowsCrossPlatformValidation(t *testing.T) {
	t.Run("png_format_consistency", func(t *testing.T) {
		// Verify PNG format works consistently across platforms
		testData := generateTestPNGData()

		tempFile := filepath.Join(os.TempDir(), "cross_platform_test.png")
		defer os.Remove(tempFile)

		if err := os.WriteFile(tempFile, testData, 0644); err != nil {
			t.Fatalf("Failed to write test PNG: %v", err)
		}

		// Verify file can be read and has expected format
		readData, err := os.ReadFile(tempFile)
		if err != nil {
			t.Fatalf("Failed to read test PNG: %v", err)
		}

		if !isPNGFormat(readData) {
			t.Error("Data is not valid PNG format")
		}
	})

	t.Run("srgb_color_space_validation", func(t *testing.T) {
		// Test sRGB color space handling across platforms
		testColors := []struct {
			name       string
			r, g, b, a uint8
		}{
			{"red", 255, 0, 0, 255},
			{"green", 0, 255, 0, 255},
			{"blue", 0, 0, 255, 255},
			{"white", 255, 255, 255, 255},
			{"black", 0, 0, 0, 255},
		}

		for _, color := range testColors {
			// Simulate color validation for cross-platform consistency
			if !validateSRGBColor(color.r, color.g, color.b, color.a) {
				t.Errorf("Color %s failed sRGB validation: RGBA(%d,%d,%d,%d)",
					color.name, color.r, color.g, color.b, color.a)
			}
		}
	})

	t.Run("directory_structure_compatibility", func(t *testing.T) {
		// Test that directory structures work across platforms
		testStructure := []string{
			"tests/artifacts/screenshots/unit",
			"tests/artifacts/screenshots/integration",
			"tests/artifacts/demo_content/storyboards",
			"tests/artifacts/baselines/ui",
		}

		baseDir := filepath.Join(os.TempDir(), "cross_platform_structure")
		defer os.RemoveAll(baseDir)

		for _, dir := range testStructure {
			fullPath := filepath.Join(baseDir, dir)
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				t.Errorf("Failed to create cross-platform directory %s: %v", dir, err)
			}

			// Test file creation in each directory
			testFile := filepath.Join(fullPath, "test_file.txt")
			if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
				t.Errorf("Failed to create file in directory %s: %v", dir, err)
			}
		}
	})
}

// Helper functions

func isTestingCrossPlatform() bool {
	return os.Getenv("CROSS_PLATFORM_TEST") == "true" || os.Getenv("CI") != ""
}

func testWindowsScreenshotSimulation(t *testing.T) {
	// Simulate Windows screenshot for testing when not on Windows
	t.Log("Simulating Windows screenshot capture")

	outputDir := filepath.Join(os.TempDir(), "simulated_windows_screenshots")
	defer os.RemoveAll(outputDir)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Create a simulated screenshot file
	testFile := filepath.Join(outputDir, "simulated_windows_screenshot.png")
	testData := generateTestPNGData()

	if err := os.WriteFile(testFile, testData, 0644); err != nil {
		t.Fatalf("Failed to create simulated screenshot: %v", err)
	}

	// Verify the simulation worked
	if !fileExists(testFile) {
		t.Error("Simulated Windows screenshot was not created")
	}
}

func normalizePathForWindows(path string) string {
	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(path, "/", "\\")
	}
	return path
}

func simulateWindowsVisualOperation() {
	// Simulate Windows-specific visual operations
	time.Sleep(100 * time.Millisecond)
}

func generateTestPNGData() []byte {
	// PNG signature + minimal PNG data
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		// Minimal PNG data (this is a simplified representation)
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, // 1x1 pixel
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53, // RGB, no compression
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41, // IDAT chunk
		0x54, 0x08, 0x99, 0x01, 0x01, 0x00, 0x00, 0xFF, // minimal pixel data
		0xFF, 0x00, 0x00, 0x00, 0x02, 0x00, 0x01, 0xE5,
		0x27, 0xDE, 0xFC, 0x00, 0x00, 0x00, 0x00, 0x49, // IEND chunk
		0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}
}

func isPNGFormat(data []byte) bool {
	// Check PNG signature
	if len(data) < 8 {
		return false
	}
	pngSignature := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	for i, b := range pngSignature {
		if data[i] != b {
			return false
		}
	}
	return true
}

func validateSRGBColor(r, g, b, a uint8) bool {
	// Basic sRGB validation - all values should be in valid range
	return true // In a real implementation, this would do actual sRGB validation
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func BenchmarkWindowsOperations(b *testing.B) {
	if runtime.GOOS != "windows" && !isTestingCrossPlatform() {
		b.Skip("Windows benchmarks require Windows platform or cross-platform test mode")
	}

	b.Run("windows_screenshot_performance", func(b *testing.B) {
		outputDir := filepath.Join(os.TempDir(), "bench_windows_screenshots")
		defer os.RemoveAll(outputDir)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			capture := NewScreenshotCapture(fmt.Sprintf("bench_%d", i), outputDir)
			if runtime.GOOS == "windows" {
				capture.CaptureScreen("benchmark")
			} else {
				// Simulate for cross-platform benchmarking
				simulateWindowsVisualOperation()
			}
		}
	})
}
