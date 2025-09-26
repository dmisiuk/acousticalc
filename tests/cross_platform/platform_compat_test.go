package cross_platform

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PlatformInfo contains information about the current platform
type PlatformInfo struct {
	OS           string          `json:"os"`
	Arch         string          `json:"arch"`
	GoVersion    string          `json:"go_version"`
	HasGUI       bool            `json:"has_gui"`
	ScreenWidth  int             `json:"screen_width"`
	ScreenHeight int             `json:"screen_height"`
	HasFFmpeg    bool            `json:"has_ffmpeg"`
	HasAsciinema bool            `json:"has_asciinema"`
	Dependencies map[string]bool `json:"dependencies"`
}

// PlatformTestConfig defines configuration for cross-platform tests
type PlatformTestConfig struct {
	TestName        string
	Timeout         time.Duration
	SkipPlatforms   []string
	RequiredTools   []string
	ValidateGUI     bool
	ValidateAudio   bool
	ValidateFileOps bool
}

// PlatformTestResult contains results of platform compatibility tests
type PlatformTestResult struct {
	Platform      string                   `json:"platform"`
	Passed        bool                     `json:"passed"`
	Errors        []string                 `json:"errors"`
	Warnings      []string                 `json:"warnings"`
	TestDuration  time.Duration            `json:"test_duration"`
	Compatibility map[string]bool          `json:"compatibility"`
	Performance   map[string]time.Duration `json:"performance"`
}

// PlatformValidator handles cross-platform validation
type PlatformValidator struct {
	info      *PlatformInfo
	results   map[string]*PlatformTestResult
	startTime time.Time
}

// NewPlatformValidator creates a new platform validator
func NewPlatformValidator() (*PlatformValidator, error) {
	info, err := gatherPlatformInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to gather platform info: %w", err)
	}

	return &PlatformValidator{
		info:    info,
		results: make(map[string]*PlatformTestResult),
	}, nil
}

// gatherPlatformInfo collects information about the current platform
func gatherPlatformInfo() (*PlatformInfo, error) {
	info := &PlatformInfo{
		OS:           runtime.GOOS,
		Arch:         runtime.GOARCH,
		Dependencies: make(map[string]bool),
	}

	// Get Go version
	if goVersion, err := exec.Command("go", "version").Output(); err == nil {
		info.GoVersion = strings.TrimSpace(string(goVersion))
	}

	// Check for GUI support
	info.HasGUI = checkGUISupport()

	// Get screen dimensions (simplified)
	info.ScreenWidth, info.ScreenHeight = getScreenDimensions()

	// Check for required tools
	tools := []string{"ffmpeg", "asciinema", "xdotool", "xclip", "screencapture", "powershell"}
	for _, tool := range tools {
		info.Dependencies[tool] = checkToolAvailable(tool)
	}

	info.HasFFmpeg = info.Dependencies["ffmpeg"]
	info.HasAsciinema = info.Dependencies["asciinema"]

	return info, nil
}

// checkGUISupport checks if the platform has GUI support
func checkGUISupport() bool {
	switch runtime.GOOS {
	case "linux":
		// Check if DISPLAY is set
		return os.Getenv("DISPLAY") != ""
	case "darwin":
		// macOS typically has GUI
		return true
	case "windows":
		// Windows typically has GUI
		return true
	default:
		return false
	}
}

// getScreenDimensions gets screen dimensions (simplified implementation)
func getScreenDimensions() (width, height int) {
	// This is a simplified implementation
	// In a real implementation, you'd use platform-specific methods
	switch runtime.GOOS {
	case "linux":
		if os.Getenv("DISPLAY") != "" {
			// Try to get dimensions from xrandr
			if cmd := exec.Command("xrandr"); cmd.Run() == nil {
				// Parse output to get dimensions (simplified)
				return 1920, 1080 // Default fallback
			}
		}
	case "darwin":
		// macOS dimensions
		return 1920, 1080 // Default fallback
	case "windows":
		// Windows dimensions
		return 1920, 1080 // Default fallback
	}
	return 1024, 768 // Default fallback
}

// checkToolAvailable checks if a tool is available in PATH
func checkToolAvailable(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}

// ValidatePlatform validates the current platform for E2E testing
func (pv *PlatformValidator) ValidatePlatform(config *PlatformTestConfig) *PlatformTestResult {
	result := &PlatformTestResult{
		Platform:      pv.info.OS,
		Compatibility: make(map[string]bool),
		Performance:   make(map[string]time.Duration),
	}

	start := time.Now()
	defer func() {
		result.TestDuration = time.Since(start)
	}()

	// Check if platform should be skipped
	for _, skipPlatform := range config.SkipPlatforms {
		if pv.info.OS == skipPlatform {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("Platform %s is in skip list", pv.info.OS))
			return result
		}
	}

	// Validate required tools
	for _, tool := range config.RequiredTools {
		if !pv.info.Dependencies[tool] {
			result.Errors = append(result.Errors,
				fmt.Sprintf("Required tool '%s' not available", tool))
		}
		result.Compatibility[tool] = pv.info.Dependencies[tool]
	}

	// Validate GUI if required
	if config.ValidateGUI {
		guiStart := time.Now()
		if !pv.info.HasGUI {
			result.Errors = append(result.Errors, "GUI support not available")
		}
		result.Performance["gui_check"] = time.Since(guiStart)
		result.Compatibility["gui"] = pv.info.HasGUI
	}

	// Validate audio if required
	if config.ValidateAudio {
		audioStart := time.Now()
		audioSupported := checkAudioSupport()
		if !audioSupported {
			result.Warnings = append(result.Warnings, "Audio support not available")
		}
		result.Performance["audio_check"] = time.Since(audioStart)
		result.Compatibility["audio"] = audioSupported
	}

	// Validate file operations if required
	if config.ValidateFileOps {
		fileStart := time.Now()
		fileOpsSupported := checkFileOperations()
		if !fileOpsSupported {
			result.Errors = append(result.Errors, "File operations not supported")
		}
		result.Performance["file_ops_check"] = time.Since(fileStart)
		result.Compatibility["file_ops"] = fileOpsSupported
	}

	// Check basic Go functionality
	goStart := time.Now()
	goSupported := checkGoFunctionality()
	result.Performance["go_check"] = time.Since(goStart)
	result.Compatibility["go"] = goSupported

	// Set overall pass/fail status
	result.Passed = len(result.Errors) == 0

	return result
}

// checkAudioSupport checks if audio is supported
func checkAudioSupport() bool {
	switch runtime.GOOS {
	case "linux":
		// Check for ALSA or PulseAudio
		tools := []string{"aplay", "paplay", "pactl"}
		for _, tool := range tools {
			if checkToolAvailable(tool) {
				return true
			}
		}
	case "darwin":
		// macOS audio support
		return checkToolAvailable("afplay")
	case "windows":
		// Windows audio support
		return true
	}
	return false
}

// checkFileOperations checks if file operations work properly
func checkFileOperations() bool {
	// Create a temporary file to test file operations
	tempDir := os.TempDir()
	testFile := filepath.Join(tempDir, "acousticalc_test_"+fmt.Sprintf("%d", time.Now().Unix()))

	// Test file creation
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return false
	}

	// Test file reading
	if _, err := os.ReadFile(testFile); err != nil {
		os.Remove(testFile)
		return false
	}

	// Test file deletion
	if err := os.Remove(testFile); err != nil {
		return false
	}

	return true
}

// checkGoFunctionality checks if basic Go functionality works
func checkGoFunctionality() bool {
	// Test if we can run a simple Go program
	cmd := exec.Command("go", "run", "-c", "package main; import \"fmt\"; func main() { fmt.Println(\"test\") }")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// GetPlatformInfo returns the gathered platform information
func (pv *PlatformValidator) GetPlatformInfo() *PlatformInfo {
	return pv.info
}

// GetResult returns the validation result for a specific test
func (pv *PlatformValidator) GetResult(testName string) *PlatformTestResult {
	return pv.results[testName]
}

// GetAllResults returns all validation results
func (pv *PlatformValidator) GetAllResults() map[string]*PlatformTestResult {
	return pv.results
}

// TestCrossPlatformCompatibility tests cross-platform compatibility
func TestCrossPlatformCompatibility(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping cross-platform test in short mode")
	}

	validator, err := NewPlatformValidator()
	require.NoError(t, err, "Should create platform validator")

	// Test basic platform compatibility
	t.Run("BasicCompatibility", func(t *testing.T) {
		config := &PlatformTestConfig{
			TestName:        "BasicCompatibility",
			Timeout:         30 * time.Second,
			ValidateGUI:     true,
			ValidateAudio:   true,
			ValidateFileOps: true,
		}

		result := validator.ValidatePlatform(config)

		// Log platform information
		t.Logf("Platform Info:")
		t.Logf("  OS: %s", validator.info.OS)
		t.Logf("  Arch: %s", validator.info.Arch)
		t.Logf("  Go Version: %s", validator.info.GoVersion)
		t.Logf("  Has GUI: %t", validator.info.HasGUI)
		t.Logf("  Screen: %dx%d", validator.info.ScreenWidth, validator.info.ScreenHeight)
		t.Logf("  Has FFmpeg: %t", validator.info.HasFFmpeg)
		t.Logf("  Has Asciinema: %t", validator.info.HasAsciinema)

		// Log compatibility results
		t.Logf("Compatibility Results:")
		for feature, compatible := range result.Compatibility {
			status := "✓"
			if !compatible {
				status = "✗"
			}
			t.Logf("  %s %s: %t", status, feature, compatible)
		}

		// Log performance metrics
		t.Logf("Performance Metrics:")
		for check, duration := range result.Performance {
			t.Logf("  %s: %v", check, duration)
		}

		// Assert that basic functionality works
		assert.True(t, result.Compatibility["go"], "Go functionality should work")
		assert.True(t, result.Compatibility["file_ops"], "File operations should work")

		// Log warnings and errors
		for _, warning := range result.Warnings {
			t.Logf("Warning: %s", warning)
		}

		for _, error := range result.Errors {
			t.Logf("Error: %s", error)
		}

		// Store result
		validator.results[config.TestName] = result
	})

	// Test E2E tool compatibility
	t.Run("E2EToolCompatibility", func(t *testing.T) {
		config := &PlatformTestConfig{
			TestName:      "E2EToolCompatibility",
			Timeout:       30 * time.Second,
			RequiredTools: []string{"ffmpeg", "asciinema"},
			ValidateGUI:   true,
		}

		result := validator.ValidatePlatform(config)

		// Check if we have the minimum required tools for E2E testing
		hasRecordingTools := result.Compatibility["ffmpeg"] || result.Compatibility["asciinema"]

		t.Logf("Has recording tools: %t", hasRecordingTools)
		t.Logf("FFmpeg available: %t", result.Compatibility["ffmpeg"])
		t.Logf("Asciinema available: %t", result.Compatibility["asciinema"])

		// Store result
		validator.results[config.TestName] = result

		// Don't fail the test if recording tools aren't available, just warn
		if !hasRecordingTools {
			t.Logf("Warning: No recording tools available - terminal recording will be disabled")
		}
	})

	// Test platform-specific features
	t.Run("PlatformSpecificFeatures", func(t *testing.T) {
		config := &PlatformTestConfig{
			TestName:      "PlatformSpecificFeatures",
			Timeout:       30 * time.Second,
			ValidateGUI:   true,
			ValidateAudio: true,
		}

		result := validator.ValidatePlatform(config)

		// Test platform-specific assertions
		switch validator.info.OS {
		case "linux":
			t.Logf("Running on Linux - checking for X11 tools")
			if validator.info.HasGUI {
				// Check for Linux-specific tools
				xdotool := validator.info.Dependencies["xdotool"]
				xclip := validator.info.Dependencies["xclip"]
				t.Logf("Xdotool available: %t", xdotool)
				t.Logf("Xclip available: %t", xclip)
			}
		case "darwin":
			t.Logf("Running on macOS - checking for macOS tools")
			screencapture := validator.info.Dependencies["screencapture"]
			t.Logf("Screencapture available: %t", screencapture)
		case "windows":
			t.Logf("Running on Windows - checking for Windows tools")
			powershell := validator.info.Dependencies["powershell"]
			t.Logf("PowerShell available: %t", powershell)
		default:
			t.Logf("Running on unsupported platform: %s", validator.info.OS)
		}

		// Store result
		validator.results[config.TestName] = result
	})
}

// TestPlatformPerformance tests platform-specific performance
func TestPlatformPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping platform performance test in short mode")
	}

	validator, err := NewPlatformValidator()
	require.NoError(t, err, "Should create platform validator")
	_ = validator // Use the validator variable to avoid "declared and not used" error

	t.Run("FileOperationPerformance", func(t *testing.T) {
		// Test file operation performance
		start := time.Now()

		// Perform multiple file operations
		for i := 0; i < 100; i++ {
			testFile := filepath.Join(os.TempDir(), fmt.Sprintf("perf_test_%d", i))
			os.WriteFile(testFile, []byte("performance test"), 0644)
			os.ReadFile(testFile)
			os.Remove(testFile)
		}

		duration := time.Since(start)
		t.Logf("100 file operations completed in %v", duration)

		// Assert reasonable performance (should be much less than 1 second)
		assert.Less(t, duration, time.Second, "File operations should be fast")
	})

	t.Run("ProcessExecutionPerformance", func(t *testing.T) {
		// Test process execution performance
		start := time.Now()

		// Execute multiple processes
		for i := 0; i < 10; i++ {
			cmd := exec.Command("go", "version")
			if err := cmd.Run(); err != nil {
				t.Logf("Warning: Go version command failed: %v", err)
			}
		}

		duration := time.Since(start)
		t.Logf("10 process executions completed in %v", duration)

		// Assert reasonable performance
		assert.Less(t, duration, 5*time.Second, "Process execution should be reasonable")
	})
}

// BenchmarkPlatformOperations benchmarks platform-specific operations
func BenchmarkPlatformOperations(b *testing.B) {
	_, err := NewPlatformValidator()
	if err != nil {
		b.Fatalf("Failed to create platform validator: %v", err)
	}

	b.Run("FileOperations", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			testFile := filepath.Join(os.TempDir(), fmt.Sprintf("bench_%d", i))
			os.WriteFile(testFile, []byte("benchmark"), 0644)
			os.ReadFile(testFile)
			os.Remove(testFile)
		}
	})

	b.Run("ProcessExecution", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cmd := exec.Command("go", "version")
			cmd.Run()
		}
	})
}
