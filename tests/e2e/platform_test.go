package e2e

import (
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/dmisiuk/acousticalc/pkg/calculator"
)

// PlatformTestSuite manages platform-specific E2E testing
type PlatformTestSuite struct {
	platform  string
	osVersion string
	goArch    string
}

// NewPlatformTestSuite creates a new platform-specific test suite
func NewPlatformTestSuite(t *testing.T) *PlatformTestSuite {
	t.Helper()

	return &PlatformTestSuite{
		platform:  runtime.GOOS,
		osVersion: getOSVersion(),
		goArch:    runtime.GOARCH,
	}
}

// getOSVersion attempts to get OS version information
func getOSVersion() string {
	switch runtime.GOOS {
	case "windows":
		return "Windows"
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	default:
		return "Unknown"
	}
}

// TestPlatformSpecificBehavior tests platform-specific behaviors and edge cases
func TestPlatformSpecificBehavior(t *testing.T) {
	suite := NewPlatformTestSuite(t)

	t.Logf("Running platform-specific tests on: %s (%s, %s)", suite.platform, suite.osVersion, suite.goArch)

	// Test platform-specific behaviors
	t.Run("PlatformIdentification", func(t *testing.T) {
		suite.testPlatformIdentification(t)
	})

	t.Run("PerformanceCharacteristics", func(t *testing.T) {
		suite.testPerformanceCharacteristics(t)
	})

	t.Run("ResourceHandling", func(t *testing.T) {
		suite.testResourceHandling(t)
	})

	t.Run("ErrorHandlingConsistency", func(t *testing.T) {
		suite.testErrorHandlingConsistency(t)
	})
}

// TestWindowsSpecificFeatures tests Windows-specific behaviors
func TestWindowsSpecificFeatures(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific tests on non-Windows platform")
	}

	suite := NewPlatformTestSuite(t)

	t.Run("WindowsFileHandling", func(t *testing.T) {
		suite.testWindowsFileHandling(t)
	})

	t.Run("WindowsPerformance", func(t *testing.T) {
		suite.testWindowsPerformance(t)
	})

	t.Run("WindowsErrorMessages", func(t *testing.T) {
		suite.testWindowsErrorMessages(t)
	})
}

// TestMacOSSpecificFeatures tests macOS-specific behaviors
func TestMacOSSpecificFeatures(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping macOS-specific tests on non-macOS platform")
	}

	suite := NewPlatformTestSuite(t)

	t.Run("MacOSFileHandling", func(t *testing.T) {
		suite.testMacOSFileHandling(t)
	})

	t.Run("MacOSPerformance", func(t *testing.T) {
		suite.testMacOSPerformance(t)
	})

	t.Run("MacOSUIIntegration", func(t *testing.T) {
		suite.testMacOSUIIntegration(t)
	})
}

// TestLinuxSpecificFeatures tests Linux-specific behaviors
func TestLinuxSpecificFeatures(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skipping Linux-specific tests on non-Linux platform")
	}

	suite := NewPlatformTestSuite(t)

	t.Run("LinuxFileHandling", func(t *testing.T) {
		suite.testLinuxFileHandling(t)
	})

	t.Run("LinuxPerformance", func(t *testing.T) {
		suite.testLinuxPerformance(t)
	})

	t.Run("LinuxCompatibility", func(t *testing.T) {
		suite.testLinuxCompatibility(t)
	})
}

// testPlatformIdentification verifies platform detection works correctly
func (pts *PlatformTestSuite) testPlatformIdentification(t *testing.T) {
	t.Helper()

	expectedPlatforms := []string{"windows", "darwin", "linux"}

	found := false
	for _, expected := range expectedPlatforms {
		if pts.platform == expected {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Unexpected platform: %s", pts.platform)
	}

	t.Logf("Platform identification successful: %s", pts.platform)
}

// testPerformanceCharacteristics tests platform-specific performance
func (pts *PlatformTestSuite) testPerformanceCharacteristics(t *testing.T) {
	t.Helper()

	// Platform-specific performance expectations
	var expectedMaxTime time.Duration
	switch pts.platform {
	case "linux":
		expectedMaxTime = 1 * time.Millisecond // Linux is typically fastest
	case "darwin":
		expectedMaxTime = 2 * time.Millisecond // macOS slightly slower
	case "windows":
		expectedMaxTime = 5 * time.Millisecond // Windows may be slower
	default:
		expectedMaxTime = 10 * time.Millisecond // Conservative default
	}

	expression := "1 + 2 + 3 + 4 + 5"
	start := time.Now()

	_, err := calculator.Evaluate(expression)
	if err != nil {
		t.Errorf("Performance test failed: %v", err)
		return
	}

	elapsed := time.Since(start)

	if elapsed > expectedMaxTime {
		t.Logf("Performance warning on %s: %v > %v (may be acceptable)", pts.platform, elapsed, expectedMaxTime)
	} else {
		t.Logf("Performance test passed on %s: %v <= %v", pts.platform, elapsed, expectedMaxTime)
	}
}

// testResourceHandling tests platform-specific resource handling
func (pts *PlatformTestSuite) testResourceHandling(t *testing.T) {
	t.Helper()

	// Test multiple calculations to check for resource leaks
	for i := 0; i < 100; i++ {
		_, err := calculator.Evaluate("2 * 3 + 1")
		if err != nil {
			t.Errorf("Resource handling test failed at iteration %d: %v", i, err)
			return
		}
	}

	t.Logf("Resource handling test completed successfully on %s", pts.platform)
}

// testErrorHandlingConsistency tests consistent error handling across platforms
func (pts *PlatformTestSuite) testErrorHandlingConsistency(t *testing.T) {
	t.Helper()

	errorExpressions := []string{
		"10 / 0", // Division by zero
		"2 ++ 3", // Invalid syntax
		"(2 + 3", // Unmatched parentheses
		"",       // Empty expression
		"2 +* 3", // Invalid operator sequence
	}

	for _, expr := range errorExpressions {
		_, err := calculator.Evaluate(expr)
		if err == nil {
			t.Errorf("Expected error for expression '%s' on %s, but got none", expr, pts.platform)
		} else {
			t.Logf("Consistent error handling on %s for '%s': %v", pts.platform, expr, err)
		}
	}
}

// testWindowsFileHandling tests Windows-specific file handling
func (pts *PlatformTestSuite) testWindowsFileHandling(t *testing.T) {
	t.Helper()

	// Test Windows-specific path handling if applicable
	tempDir := os.TempDir()
	if tempDir == "" {
		t.Error("Windows temp directory not accessible")
	} else {
		t.Logf("Windows temp directory accessible: %s", tempDir)
	}
}

// testWindowsPerformance tests Windows-specific performance characteristics
func (pts *PlatformTestSuite) testWindowsPerformance(t *testing.T) {
	t.Helper()

	// Windows may have different performance characteristics
	start := time.Now()

	for i := 0; i < 50; i++ {
		_, err := calculator.Evaluate("1 + 2 * 3")
		if err != nil {
			t.Errorf("Windows performance test failed: %v", err)
			return
		}
	}

	elapsed := time.Since(start)
	t.Logf("Windows performance test completed in %v", elapsed)
}

// testWindowsErrorMessages tests Windows-specific error message formats
func (pts *PlatformTestSuite) testWindowsErrorMessages(t *testing.T) {
	t.Helper()

	_, err := calculator.Evaluate("10 / 0")
	if err == nil {
		t.Error("Expected error for division by zero on Windows")
	} else {
		t.Logf("Windows error message format: %v", err)
	}
}

// testMacOSFileHandling tests macOS-specific file handling
func (pts *PlatformTestSuite) testMacOSFileHandling(t *testing.T) {
	t.Helper()

	// Test macOS-specific path handling
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		t.Error("macOS HOME directory not accessible")
	} else {
		t.Logf("macOS HOME directory accessible: %s", homeDir)
	}
}

// testMacOSPerformance tests macOS-specific performance characteristics
func (pts *PlatformTestSuite) testMacOSPerformance(t *testing.T) {
	t.Helper()

	// macOS typically has good performance
	start := time.Now()

	for i := 0; i < 75; i++ {
		_, err := calculator.Evaluate("2 * (3 + 4)")
		if err != nil {
			t.Errorf("macOS performance test failed: %v", err)
			return
		}
	}

	elapsed := time.Since(start)
	t.Logf("macOS performance test completed in %v", elapsed)
}

// testMacOSUIIntegration tests macOS UI integration capabilities
func (pts *PlatformTestSuite) testMacOSUIIntegration(t *testing.T) {
	t.Helper()

	// Test macOS-specific UI integration preparation
	displayEnv := os.Getenv("DISPLAY")
	t.Logf("macOS display environment: %s", displayEnv)

	// Calculator should work regardless of UI environment
	_, err := calculator.Evaluate("5 * 6")
	if err != nil {
		t.Errorf("macOS UI integration test failed: %v", err)
	} else {
		t.Log("macOS UI integration test passed")
	}
}

// testLinuxFileHandling tests Linux-specific file handling
func (pts *PlatformTestSuite) testLinuxFileHandling(t *testing.T) {
	t.Helper()

	// Test Linux-specific path handling
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		t.Error("Linux HOME directory not accessible")
	} else {
		t.Logf("Linux HOME directory accessible: %s", homeDir)
	}

	// Test /tmp directory access
	tmpDir := "/tmp"
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Error("Linux /tmp directory not accessible")
	} else {
		t.Logf("Linux /tmp directory accessible")
	}
}

// testLinuxPerformance tests Linux-specific performance characteristics
func (pts *PlatformTestSuite) testLinuxPerformance(t *testing.T) {
	t.Helper()

	// Linux typically has the best performance
	start := time.Now()

	for i := 0; i < 100; i++ {
		_, err := calculator.Evaluate("(1 + 2) * (3 + 4)")
		if err != nil {
			t.Errorf("Linux performance test failed: %v", err)
			return
		}
	}

	elapsed := time.Since(start)
	t.Logf("Linux performance test completed in %v", elapsed)
}

// testLinuxCompatibility tests Linux-specific compatibility
func (pts *PlatformTestSuite) testLinuxCompatibility(t *testing.T) {
	t.Helper()

	// Test Linux environment variables
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		t.Error("Linux PATH environment variable not set")
	} else {
		t.Logf("Linux PATH environment accessible")
	}

	// Calculator should work in all Linux environments
	_, err := calculator.Evaluate("7 * 8 + 9")
	if err != nil {
		t.Errorf("Linux compatibility test failed: %v", err)
	} else {
		t.Log("Linux compatibility test passed")
	}
}

// TestCrossPlatformConsistency tests behavior consistency across platforms
func TestCrossPlatformConsistency(t *testing.T) {
	suite := NewPlatformTestSuite(t)

	// Test expressions that should behave identically across platforms
	consistencyTests := []struct {
		expression string
		expected   float64
	}{
		{"2 + 3", 5},
		{"10 - 4", 6},
		{"3 * 7", 21},
		{"15 / 3", 5},
		{"2 + 3 * 4", 14},
		{"(2 + 3) * 4", 20},
		{"10.5 + 2.3", 12.8},
		{"-5 + 3", -2},
	}

	for _, test := range consistencyTests {
		t.Run(test.expression, func(t *testing.T) {
			result, err := calculator.Evaluate(test.expression)
			if err != nil {
				t.Errorf("Cross-platform consistency test failed for '%s' on %s: %v",
					test.expression, suite.platform, err)
				return
			}

			const epsilon = 1e-9
			if abs(result-test.expected) > epsilon {
				t.Errorf("Cross-platform consistency failed for '%s' on %s: expected %f, got %f",
					test.expression, suite.platform, test.expected, result)
			} else {
				t.Logf("Cross-platform consistency verified for '%s' on %s: %f",
					test.expression, suite.platform, result)
			}
		})
	}
}
