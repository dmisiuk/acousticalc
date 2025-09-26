package e2e

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// PlatformTestSuite manages platform-specific E2E tests
type PlatformTestSuite struct {
	config   *E2ETestConfig
	platform string
}

// NewPlatformTestSuite creates a new platform test suite
func NewPlatformTestSuite(config *E2ETestConfig) *PlatformTestSuite {
	if config == nil {
		config = DefaultE2EConfig()
	}
	
	return &PlatformTestSuite{
		config:   config,
		platform: runtime.GOOS,
	}
}

// TestPlatformSpecificBehavior tests platform-specific behaviors
func TestPlatformSpecificBehavior(t *testing.T) {
	config := DefaultE2EConfig()
	suite := NewPlatformTestSuite(config)
	
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	
	// Run platform-specific tests based on current platform
	switch runtime.GOOS {
	case "windows":
		suite.testWindowsSpecificBehavior(ctx, t)
	case "darwin":
		suite.testMacOSSpecificBehavior(ctx, t)
	case "linux":
		suite.testLinuxSpecificBehavior(ctx, t)
	default:
		t.Skipf("Platform %s not supported for platform-specific testing", runtime.GOOS)
	}
}

// testWindowsSpecificBehavior tests Windows-specific behaviors
func (suite *PlatformTestSuite) testWindowsSpecificBehavior(ctx context.Context, t *testing.T) {
	t.Run("WindowsFileSystem", func(t *testing.T) {
		// Test Windows file system behavior
		testDir := filepath.Join(suite.config.ArtifactsDir, "windows_test")
		err := os.MkdirAll(testDir, 0755)
		if err != nil {
			t.Errorf("Failed to create Windows test directory: %v", err)
		}
		defer os.RemoveAll(testDir)
		
		// Test Windows path handling
		testFile := filepath.Join(testDir, "test_file.txt")
		content := "Windows E2E test content"
		
		err = os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			t.Errorf("Failed to create Windows test file: %v", err)
		}
		
		// Verify file exists
		if _, err := os.Stat(testFile); os.IsNotExist(err) {
			t.Error("Windows test file was not created")
		}
	})
	
	t.Run("WindowsPowerShellIntegration", func(t *testing.T) {
		// Test PowerShell integration for Windows automation
		cmd := exec.CommandContext(ctx, "powershell", "-Command", "Get-Host | Select-Object Name")
		output, err := cmd.Output()
		if err != nil {
			t.Logf("PowerShell not available or failed: %v", err)
			return
		}
		
		if !strings.Contains(string(output), "PowerShell") {
			t.Error("PowerShell integration test failed")
		}
		
		if suite.config.Verbose {
			t.Logf("PowerShell integration test passed: %s", string(output))
		}
	})
	
	t.Run("WindowsEnvironmentVariables", func(t *testing.T) {
		// Test Windows environment variable handling
		envVars := []string{"PATH", "TEMP", "USERPROFILE"}
		
		for _, envVar := range envVars {
			value := os.Getenv(envVar)
			if value == "" {
				t.Errorf("Windows environment variable %s is not set", envVar)
			}
			
			if suite.config.Verbose {
				t.Logf("Windows environment variable %s: %s", envVar, value)
			}
		}
	})
}

// testMacOSSpecificBehavior tests macOS-specific behaviors
func (suite *PlatformTestSuite) testMacOSSpecificBehavior(ctx context.Context, t *testing.T) {
	t.Run("MacOSFileSystem", func(t *testing.T) {
		// Test macOS file system behavior
		testDir := filepath.Join(suite.config.ArtifactsDir, "macos_test")
		err := os.MkdirAll(testDir, 0755)
		if err != nil {
			t.Errorf("Failed to create macOS test directory: %v", err)
		}
		defer os.RemoveAll(testDir)
		
		// Test macOS path handling
		testFile := filepath.Join(testDir, "test_file.txt")
		content := "macOS E2E test content"
		
		err = os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			t.Errorf("Failed to create macOS test file: %v", err)
		}
		
		// Verify file exists
		if _, err := os.Stat(testFile); os.IsNotExist(err) {
			t.Error("macOS test file was not created")
		}
	})
	
	t.Run("MacOSScreenshotIntegration", func(t *testing.T) {
		// Test macOS screenshot integration (from Story 0.2.2)
		cmd := exec.CommandContext(ctx, "screencapture", "-h")
		err := cmd.Run()
		if err != nil {
			t.Logf("macOS screencapture not available: %v", err)
			return
		}
		
		// Test screenshot capture
		screenshotFile := filepath.Join(suite.config.ArtifactsDir, "macos_test_screenshot.png")
		cmd = exec.CommandContext(ctx, "screencapture", "-x", screenshotFile)
		err = cmd.Run()
		if err != nil {
			t.Logf("macOS screenshot capture failed: %v", err)
			return
		}
		
		// Verify screenshot was created
		if _, err := os.Stat(screenshotFile); os.IsNotExist(err) {
			t.Error("macOS screenshot was not created")
		} else {
			// Clean up
			os.Remove(screenshotFile)
		}
	})
	
	t.Run("MacOSEnvironmentVariables", func(t *testing.T) {
		// Test macOS environment variable handling
		envVars := []string{"PATH", "HOME", "USER"}
		
		for _, envVar := range envVars {
			value := os.Getenv(envVar)
			if value == "" {
				t.Errorf("macOS environment variable %s is not set", envVar)
			}
			
			if suite.config.Verbose {
				t.Logf("macOS environment variable %s: %s", envVar, value)
			}
		}
	})
}

// testLinuxSpecificBehavior tests Linux-specific behaviors
func (suite *PlatformTestSuite) testLinuxSpecificBehavior(ctx context.Context, t *testing.T) {
	t.Run("LinuxFileSystem", func(t *testing.T) {
		// Test Linux file system behavior
		testDir := filepath.Join(suite.config.ArtifactsDir, "linux_test")
		err := os.MkdirAll(testDir, 0755)
		if err != nil {
			t.Errorf("Failed to create Linux test directory: %v", err)
		}
		defer os.RemoveAll(testDir)
		
		// Test Linux path handling
		testFile := filepath.Join(testDir, "test_file.txt")
		content := "Linux E2E test content"
		
		err = os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			t.Errorf("Failed to create Linux test file: %v", err)
		}
		
		// Verify file exists
		if _, err := os.Stat(testFile); os.IsNotExist(err) {
			t.Error("Linux test file was not created")
		}
	})
	
	t.Run("LinuxXvfbIntegration", func(t *testing.T) {
		// Test Xvfb integration for headless testing (from Story 0.2.2)
		cmd := exec.CommandContext(ctx, "which", "xvfb-run")
		err := cmd.Run()
		if err != nil {
			t.Logf("Xvfb not available: %v", err)
			return
		}
		
		// Test Xvfb execution
		cmd = exec.CommandContext(ctx, "xvfb-run", "-a", "echo", "Xvfb test")
		output, err := cmd.Output()
		if err != nil {
			t.Logf("Xvfb execution failed: %v", err)
			return
		}
		
		if !strings.Contains(string(output), "Xvfb test") {
			t.Error("Linux Xvfb integration test failed")
		}
		
		if suite.config.Verbose {
			t.Logf("Linux Xvfb integration test passed: %s", string(output))
		}
	})
	
	t.Run("LinuxEnvironmentVariables", func(t *testing.T) {
		// Test Linux environment variable handling
		envVars := []string{"PATH", "HOME", "USER", "DISPLAY"}
		
		for _, envVar := range envVars {
			value := os.Getenv(envVar)
			if envVar == "DISPLAY" && value == "" {
				// DISPLAY might not be set in headless environments
				continue
			}
			
			if value == "" {
				t.Errorf("Linux environment variable %s is not set", envVar)
			}
			
			if suite.config.Verbose {
				t.Logf("Linux environment variable %s: %s", envVar, value)
			}
		}
	})
}

// TestCrossPlatformBehaviorConsistency tests behavior consistency across platforms
func TestCrossPlatformBehaviorConsistency(t *testing.T) {
	config := DefaultE2EConfig()
	suite := NewPlatformTestSuite(config)
	
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	
	t.Run("FileSystemConsistency", func(t *testing.T) {
		suite.testFileSystemConsistency(ctx, t)
	})
	
	t.Run("EnvironmentConsistency", func(t *testing.T) {
		suite.testEnvironmentConsistency(ctx, t)
	})
	
	t.Run("PathHandlingConsistency", func(t *testing.T) {
		suite.testPathHandlingConsistency(ctx, t)
	})
}

// testFileSystemConsistency tests file system behavior consistency
func (suite *PlatformTestSuite) testFileSystemConsistency(ctx context.Context, t *testing.T) {
	// Test that file operations work consistently across platforms
	testDir := filepath.Join(suite.config.ArtifactsDir, "consistency_test")
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Errorf("Failed to create consistency test directory: %v", err)
	}
	defer os.RemoveAll(testDir)
	
	// Test file creation
	testFile := filepath.Join(testDir, "consistency_test.txt")
	content := fmt.Sprintf("Cross-platform consistency test - %s", time.Now().Format(time.RFC3339))
	
	err = os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Errorf("Failed to create consistency test file: %v", err)
	}
	
	// Test file reading
	readContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Errorf("Failed to read consistency test file: %v", err)
	}
	
	if string(readContent) != content {
		t.Error("File content consistency test failed")
	}
	
	// Test file permissions
	info, err := os.Stat(testFile)
	if err != nil {
		t.Errorf("Failed to stat consistency test file: %v", err)
	}
	
	// Verify file is readable and writable
	if info.Mode()&0400 == 0 {
		t.Error("File is not readable")
	}
	
	if info.Mode()&0200 == 0 {
		t.Error("File is not writable")
	}
}

// testEnvironmentConsistency tests environment variable consistency
func (suite *PlatformTestSuite) testEnvironmentConsistency(ctx context.Context, t *testing.T) {
	// Test that common environment variables are available
	commonEnvVars := []string{"PATH"}
	
	for _, envVar := range commonEnvVars {
		value := os.Getenv(envVar)
		if value == "" {
			t.Errorf("Common environment variable %s is not set", envVar)
		}
		
		if suite.config.Verbose {
			t.Logf("Environment variable %s: %s", envVar, value)
		}
	}
	
	// Test that we can set and read environment variables
	testEnvVar := "E2E_TEST_VAR"
	testValue := "test_value_123"
	
	err := os.Setenv(testEnvVar, testValue)
	if err != nil {
		t.Errorf("Failed to set test environment variable: %v", err)
	}
	
	readValue := os.Getenv(testEnvVar)
	if readValue != testValue {
		t.Error("Environment variable consistency test failed")
	}
	
	// Clean up
	os.Unsetenv(testEnvVar)
}

// testPathHandlingConsistency tests path handling consistency
func (suite *PlatformTestSuite) testPathHandlingConsistency(ctx context.Context, t *testing.T) {
	// Test that path operations work consistently
	testDir := filepath.Join(suite.config.ArtifactsDir, "path_test")
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Errorf("Failed to create path test directory: %v", err)
	}
	defer os.RemoveAll(testDir)
	
	// Test path joining
	subDir := filepath.Join(testDir, "subdir")
	err = os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Errorf("Failed to create subdirectory: %v", err)
	}
	
	// Test path existence
	if _, err := os.Stat(subDir); os.IsNotExist(err) {
		t.Error("Subdirectory was not created")
	}
	
	// Test path cleaning
	cleanPath := filepath.Clean(testDir + "/../path_test")
	if !strings.Contains(cleanPath, "path_test") {
		t.Error("Path cleaning consistency test failed")
	}
}

// TestPlatformDetection tests platform detection accuracy
func TestPlatformDetection(t *testing.T) {
	config := DefaultE2EConfig()
	suite := NewPlatformTestSuite(config)
	
	t.Run("PlatformDetection", func(t *testing.T) {
		// Test that platform detection works correctly
		detectedPlatform := runtime.GOOS
		expectedPlatform := suite.platform
		
		if detectedPlatform != expectedPlatform {
			t.Errorf("Platform detection mismatch: expected %s, got %s", expectedPlatform, detectedPlatform)
		}
		
		// Test that platform is one of the supported platforms
		supportedPlatforms := []string{"windows", "darwin", "linux"}
		isSupported := false
		
		for _, platform := range supportedPlatforms {
			if detectedPlatform == platform {
				isSupported = true
				break
			}
		}
		
		if !isSupported {
			t.Errorf("Platform %s is not supported", detectedPlatform)
		}
		
		if suite.config.Verbose {
			t.Logf("Platform detection test passed: %s", detectedPlatform)
		}
	})
	
	t.Run("ArchitectureDetection", func(t *testing.T) {
		// Test that architecture detection works correctly
		arch := runtime.GOARCH
		
		if arch == "" {
			t.Error("Architecture detection failed")
		}
		
		// Test that architecture is one of the common architectures
		commonArchs := []string{"amd64", "386", "arm64", "arm"}
		isCommon := false
		
		for _, commonArch := range commonArchs {
			if arch == commonArch {
				isCommon = true
				break
			}
		}
		
		if !isCommon {
			t.Logf("Architecture %s is not in the common list, but that's okay", arch)
		}
		
		if suite.config.Verbose {
			t.Logf("Architecture detection test passed: %s", arch)
		}
	})
}

// TestPlatformSpecificPerformance tests platform-specific performance characteristics
func TestPlatformSpecificPerformance(t *testing.T) {
	config := DefaultE2EConfig()
	suite := NewPlatformTestSuite(config)
	
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	
	t.Run("FileOperationPerformance", func(t *testing.T) {
		suite.testFileOperationPerformance(ctx, t)
	})
	
	t.Run("ProcessCreationPerformance", func(t *testing.T) {
		suite.testProcessCreationPerformance(ctx, t)
	})
}

// testFileOperationPerformance tests file operation performance
func (suite *PlatformTestSuite) testFileOperationPerformance(ctx context.Context, t *testing.T) {
	testDir := filepath.Join(suite.config.ArtifactsDir, "performance_test")
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Errorf("Failed to create performance test directory: %v", err)
	}
	defer os.RemoveAll(testDir)
	
	// Test file creation performance
	start := time.Now()
	
	for i := 0; i < 100; i++ {
		testFile := filepath.Join(testDir, fmt.Sprintf("test_file_%d.txt", i))
		content := fmt.Sprintf("Performance test file %d", i)
		
		err := os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			t.Errorf("Failed to create performance test file %d: %v", i, err)
			return
		}
	}
	
	duration := time.Since(start)
	avgDuration := duration / 100
	
	if suite.config.Verbose {
		t.Logf("File creation performance: 100 files in %v (avg: %v per file)", duration, avgDuration)
	}
	
	// File creation should be reasonably fast (within 10ms per file)
	if avgDuration > 10*time.Millisecond {
		t.Errorf("File creation performance too slow: %v per file", avgDuration)
	}
}

// testProcessCreationPerformance tests process creation performance
func (suite *PlatformTestSuite) testProcessCreationPerformance(ctx context.Context, t *testing.T) {
	// Test process creation performance
	start := time.Now()
	
	for i := 0; i < 10; i++ {
		var cmd *exec.Cmd
		
		switch runtime.GOOS {
		case "windows":
			cmd = exec.CommandContext(ctx, "cmd", "/c", "echo", "test")
		default:
			cmd = exec.CommandContext(ctx, "echo", "test")
		}
		
		err := cmd.Run()
		if err != nil {
			t.Logf("Process creation test %d failed: %v", i, err)
			continue
		}
	}
	
	duration := time.Since(start)
	avgDuration := duration / 10
	
	if suite.config.Verbose {
		t.Logf("Process creation performance: 10 processes in %v (avg: %v per process)", duration, avgDuration)
	}
	
	// Process creation should be reasonably fast (within 100ms per process)
	if avgDuration > 100*time.Millisecond {
		t.Errorf("Process creation performance too slow: %v per process", avgDuration)
	}
}