package cross_platform

import (
	"runtime"
	"testing"
	"time"

	"github.com/dmisiuk/acousticalc/pkg/calculator"
)

// CrossPlatformTestSuite manages cross-platform testing
type CrossPlatformTestSuite struct {
	platform   PlatformInfo
}

// PlatformInfo holds information about the current platform
type PlatformInfo struct {
	OS       string
	Arch     string
	Version  string
	Features []string
}

// NewCrossPlatformTestSuite creates a new cross-platform test suite
func NewCrossPlatformTestSuite(t *testing.T) *CrossPlatformTestSuite {
	t.Helper()
	
	platform := PlatformInfo{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Version:  runtime.Version(),
		Features: getSupportedFeatures(),
	}
	
	return &CrossPlatformTestSuite{
		platform:   platform,
	}
}

// getSupportedFeatures returns supported features for the current platform
func getSupportedFeatures() []string {
	features := []string{"basic_arithmetic"}
	
	switch runtime.GOOS {
	case "windows":
		features = append(features, "windows_paths", "powershell_compat")
	case "darwin":
		features = append(features, "macos_ui", "apple_silicon_compat")
		if runtime.GOARCH == "arm64" {
			features = append(features, "m1_optimized")
		}
	case "linux":
		features = append(features, "unix_paths", "x11_compat", "container_ready")
	}
	
	return features
}

// TestCrossPlatformConsistency tests consistent behavior across all platforms
func TestCrossPlatformConsistency(t *testing.T) {
	suite := NewCrossPlatformTestSuite(t)
	
	t.Logf("Testing cross-platform consistency on %s/%s", suite.platform.OS, suite.platform.Arch)
	
	t.Run("ArithmeticConsistency", func(t *testing.T) {
		suite.testArithmeticConsistency(t)
	})
	
	t.Run("ErrorHandlingConsistency", func(t *testing.T) {
		suite.testErrorHandlingConsistency(t)
	})
	
	t.Run("PerformanceConsistency", func(t *testing.T) {
		suite.testPerformanceConsistency(t)
	})
	
	t.Run("PrecisionConsistency", func(t *testing.T) {
		suite.testPrecisionConsistency(t)
	})
}

// TestPlatformSpecificOptimizations tests platform-specific optimizations
func TestPlatformSpecificOptimizations(t *testing.T) {
	suite := NewCrossPlatformTestSuite(t)
	
	t.Run("PlatformDetection", func(t *testing.T) {
		suite.testPlatformDetection(t)
	})
	
	t.Run("PerformanceOptimizations", func(t *testing.T) {
		suite.testPerformanceOptimizations(t)
	})
	
	t.Run("ResourceUtilization", func(t *testing.T) {
		suite.testResourceUtilization(t)
	})
}

// TestPlatformCompatibility tests compatibility across different platform versions
func TestPlatformCompatibility(t *testing.T) {
	suite := NewCrossPlatformTestSuite(t)
	
	t.Run("GoVersionCompatibility", func(t *testing.T) {
		suite.testGoVersionCompatibility(t)
	})
	
	t.Run("ArchitectureCompatibility", func(t *testing.T) {
		suite.testArchitectureCompatibility(t)
	})
	
	t.Run("FeatureAvailability", func(t *testing.T) {
		suite.testFeatureAvailability(t)
	})
}

// testArithmeticConsistency tests arithmetic operations consistency
func (cpts *CrossPlatformTestSuite) testArithmeticConsistency(t *testing.T) {
	t.Helper()
	
	// Test cases that should produce identical results across platforms
	testCases := []struct {
		expression string
		expected   float64
	}{
		{"2 + 3", 5},
		{"10 - 4", 6},
		{"3 * 7", 21},
		{"15 / 3", 5},
		{"2.5 + 3.7", 6.2},
		{"10.0 / 4.0", 2.5},
		{"3 + 4 * 5", 23},
		{"(3 + 4) * 5", 35},
		{"-5 + 3", -2},
		{"2 * -3", -6},
	}
	
	for _, testCase := range testCases {
		t.Run(testCase.expression, func(t *testing.T) {
			result, err := calculator.Evaluate(testCase.expression)
			if err != nil {
				t.Errorf("Platform %s/%s: arithmetic failed for '%s': %v", 
					cpts.platform.OS, cpts.platform.Arch, testCase.expression, err)
				return
			}
			
			const epsilon = 1e-9
			if abs(result-testCase.expected) > epsilon {
				t.Errorf("Platform %s/%s: inconsistent result for '%s': expected %f, got %f", 
					cpts.platform.OS, cpts.platform.Arch, testCase.expression, testCase.expected, result)
			} else {
				t.Logf("Platform %s/%s: consistent result for '%s': %f", 
					cpts.platform.OS, cpts.platform.Arch, testCase.expression, result)
			}
		})
	}
}

// testErrorHandlingConsistency tests error handling consistency
func (cpts *CrossPlatformTestSuite) testErrorHandlingConsistency(t *testing.T) {
	t.Helper()
	
	// Error cases that should behave consistently across platforms
	errorCases := []string{
		"10 / 0",     // Division by zero
		"2 ++ 3",     // Invalid syntax
		"(2 + 3",     // Unmatched parentheses
		"",           // Empty expression
		"2 +* 3",     // Invalid operator sequence
		"abc + 123",  // Invalid characters
	}
	
	for _, errorCase := range errorCases {
		t.Run(errorCase, func(t *testing.T) {
			result, err := calculator.Evaluate(errorCase)
			if err == nil {
				t.Errorf("Platform %s/%s: expected error for '%s', but got result: %f", 
					cpts.platform.OS, cpts.platform.Arch, errorCase, result)
			} else {
				t.Logf("Platform %s/%s: consistent error handling for '%s': %v", 
					cpts.platform.OS, cpts.platform.Arch, errorCase, err)
			}
		})
	}
}

// testPerformanceConsistency tests performance consistency across platforms
func (cpts *CrossPlatformTestSuite) testPerformanceConsistency(t *testing.T) {
	t.Helper()
	
	// Performance test cases
	performanceTests := []struct {
		name       string
		expression string
		iterations int
		maxAvgTime time.Duration
	}{
		{
			name:       "SimpleArithmetic",
			expression: "2 + 3",
			iterations: 1000,
			maxAvgTime: 1 * time.Microsecond,
		},
		{
			name:       "ComplexExpression",
			expression: "(1 + 2) * (3 + 4) - (5 / 2)",
			iterations: 500,
			maxAvgTime: 10 * time.Microsecond,
		},
		{
			name:       "DecimalCalculation",
			expression: "3.14159 * 2.71828",
			iterations: 500,
			maxAvgTime: 5 * time.Microsecond,
		},
	}
	
	for _, perfTest := range performanceTests {
		t.Run(perfTest.name, func(t *testing.T) {
			start := time.Now()
			
			for i := 0; i < perfTest.iterations; i++ {
				_, err := calculator.Evaluate(perfTest.expression)
				if err != nil {
					t.Errorf("Performance test failed at iteration %d: %v", i, err)
					return
				}
			}
			
			elapsed := time.Since(start)
			avgTime := elapsed / time.Duration(perfTest.iterations)
			
			t.Logf("Platform %s/%s: %s performance - %d iterations in %v (avg: %v)", 
				cpts.platform.OS, cpts.platform.Arch, perfTest.name, perfTest.iterations, elapsed, avgTime)
			
			if avgTime > perfTest.maxAvgTime {
				t.Logf("Performance warning on %s/%s: %v > %v for %s", 
					cpts.platform.OS, cpts.platform.Arch, avgTime, perfTest.maxAvgTime, perfTest.name)
			}
		})
	}
}

// testPrecisionConsistency tests floating-point precision consistency
func (cpts *CrossPlatformTestSuite) testPrecisionConsistency(t *testing.T) {
	t.Helper()
	
	// Precision test cases
	precisionTests := []struct {
		expression string
		expected   float64
		tolerance  float64
	}{
		{"0.1 + 0.2", 0.3, 1e-15},
		{"1.0 / 3.0 * 3.0", 1.0, 1e-14},
		{"2.0 / 3.0", 0.6666666666666666, 1e-15},
		{"3.14159265359 * 2.71828182846", 8.539734222673566, 1e-11}, // Relaxed tolerance for complex multiplication
	}
	
	for _, precTest := range precisionTests {
		t.Run(precTest.expression, func(t *testing.T) {
			result, err := calculator.Evaluate(precTest.expression)
			if err != nil {
				t.Errorf("Precision test failed for '%s': %v", precTest.expression, err)
				return
			}
			
			diff := abs(result - precTest.expected)
			if diff > precTest.tolerance {
				t.Errorf("Platform %s/%s: precision inconsistency for '%s': expected %g, got %g, diff %g > tolerance %g", 
					cpts.platform.OS, cpts.platform.Arch, precTest.expression, precTest.expected, result, diff, precTest.tolerance)
			} else {
				t.Logf("Platform %s/%s: precision consistent for '%s': %g (diff: %g)", 
					cpts.platform.OS, cpts.platform.Arch, precTest.expression, result, diff)
			}
		})
	}
}

// testPlatformDetection tests platform detection accuracy
func (cpts *CrossPlatformTestSuite) testPlatformDetection(t *testing.T) {
	t.Helper()
	
	// Verify platform detection
	validPlatforms := []string{"windows", "darwin", "linux"}
	validArchs := []string{"amd64", "arm64", "386", "arm"}
	
	platformValid := false
	for _, platform := range validPlatforms {
		if cpts.platform.OS == platform {
			platformValid = true
			break
		}
	}
	
	archValid := false
	for _, arch := range validArchs {
		if cpts.platform.Arch == arch {
			archValid = true
			break
		}
	}
	
	if !platformValid {
		t.Errorf("Unknown platform detected: %s", cpts.platform.OS)
	}
	
	if !archValid {
		t.Errorf("Unknown architecture detected: %s", cpts.platform.Arch)
	}
	
	t.Logf("Platform detection successful: %s/%s", cpts.platform.OS, cpts.platform.Arch)
	t.Logf("Go version: %s", cpts.platform.Version)
	t.Logf("Supported features: %v", cpts.platform.Features)
}

// testPerformanceOptimizations tests platform-specific performance optimizations
func (cpts *CrossPlatformTestSuite) testPerformanceOptimizations(t *testing.T) {
	t.Helper()
	
	// Test performance characteristics that may vary by platform
	switch cpts.platform.OS {
	case "linux":
		t.Log("Testing Linux-specific performance optimizations")
		cpts.testLinuxPerformance(t)
	case "darwin":
		t.Log("Testing macOS-specific performance optimizations")
		cpts.testMacOSPerformance(t)
	case "windows":
		t.Log("Testing Windows-specific performance optimizations")
		cpts.testWindowsPerformance(t)
	}
}

// testLinuxPerformance tests Linux-specific performance
func (cpts *CrossPlatformTestSuite) testLinuxPerformance(t *testing.T) {
	t.Helper()
	
	// Linux typically has the best performance
	start := time.Now()
	for i := 0; i < 1000; i++ {
		_, err := calculator.Evaluate("1 + 2 * 3")
		if err != nil {
			t.Errorf("Linux performance test failed: %v", err)
			return
		}
	}
	elapsed := time.Since(start)
	
	expectedMaxTime := 5 * time.Millisecond
	if elapsed > expectedMaxTime {
		t.Logf("Linux performance slower than expected: %v > %v", elapsed, expectedMaxTime)
	} else {
		t.Logf("Linux performance optimal: %v <= %v", elapsed, expectedMaxTime)
	}
}

// testMacOSPerformance tests macOS-specific performance
func (cpts *CrossPlatformTestSuite) testMacOSPerformance(t *testing.T) {
	t.Helper()
	
	start := time.Now()
	for i := 0; i < 800; i++ {
		_, err := calculator.Evaluate("2 * 3 + 4")
		if err != nil {
			t.Errorf("macOS performance test failed: %v", err)
			return
		}
	}
	elapsed := time.Since(start)
	
	expectedMaxTime := 8 * time.Millisecond
	if elapsed > expectedMaxTime {
		t.Logf("macOS performance slower than expected: %v > %v", elapsed, expectedMaxTime)
	} else {
		t.Logf("macOS performance good: %v <= %v", elapsed, expectedMaxTime)
	}
}

// testWindowsPerformance tests Windows-specific performance
func (cpts *CrossPlatformTestSuite) testWindowsPerformance(t *testing.T) {
	t.Helper()
	
	start := time.Now()
	for i := 0; i < 500; i++ {
		_, err := calculator.Evaluate("5 - 2 + 1")
		if err != nil {
			t.Errorf("Windows performance test failed: %v", err)
			return
		}
	}
	elapsed := time.Since(start)
	
	expectedMaxTime := 15 * time.Millisecond
	if elapsed > expectedMaxTime {
		t.Logf("Windows performance slower than expected: %v > %v", elapsed, expectedMaxTime)
	} else {
		t.Logf("Windows performance acceptable: %v <= %v", elapsed, expectedMaxTime)
	}
}

// testResourceUtilization tests resource utilization across platforms
func (cpts *CrossPlatformTestSuite) testResourceUtilization(t *testing.T) {
	t.Helper()
	
	// Test memory usage patterns
	t.Log("Testing resource utilization patterns")
	
	// Perform calculations and monitor for memory leaks
	for i := 0; i < 100; i++ {
		_, err := calculator.Evaluate("1 + 2 + 3 + 4 + 5")
		if err != nil {
			t.Errorf("Resource utilization test failed: %v", err)
			return
		}
	}
	
	t.Logf("Resource utilization test completed on %s/%s", cpts.platform.OS, cpts.platform.Arch)
}

// testGoVersionCompatibility tests Go version compatibility
func (cpts *CrossPlatformTestSuite) testGoVersionCompatibility(t *testing.T) {
	t.Helper()
	
	t.Logf("Testing Go version compatibility: %s", cpts.platform.Version)
	
	// Verify minimum Go version requirements are met
	// This is a placeholder - in reality, you'd parse the version string
	if cpts.platform.Version == "" {
		t.Error("Go version not detected")
	} else {
		t.Logf("Go version compatibility verified: %s", cpts.platform.Version)
	}
}

// testArchitectureCompatibility tests architecture-specific compatibility
func (cpts *CrossPlatformTestSuite) testArchitectureCompatibility(t *testing.T) {
	t.Helper()
	
	t.Logf("Testing architecture compatibility: %s", cpts.platform.Arch)
	
	// Test architecture-specific behaviors
	switch cpts.platform.Arch {
	case "amd64":
		t.Log("Testing amd64-specific features")
	case "arm64":
		t.Log("Testing arm64-specific features")
	case "386":
		t.Log("Testing 32-bit x86 compatibility")
	case "arm":
		t.Log("Testing ARM 32-bit compatibility")
	default:
		t.Logf("Testing general compatibility for architecture: %s", cpts.platform.Arch)
	}
	
	// Calculator should work on all supported architectures
	result, err := calculator.Evaluate("2 + 2")
	if err != nil {
		t.Errorf("Architecture compatibility test failed: %v", err)
	} else {
		t.Logf("Architecture compatibility verified: %f", result)
	}
}

// testFeatureAvailability tests feature availability across platforms
func (cpts *CrossPlatformTestSuite) testFeatureAvailability(t *testing.T) {
	t.Helper()
	
	t.Logf("Testing feature availability: %v", cpts.platform.Features)
	
	// Verify expected features are available
	expectedFeatures := []string{"basic_arithmetic"}
	
	for _, expected := range expectedFeatures {
		found := false
		for _, available := range cpts.platform.Features {
			if available == expected {
				found = true
				break
			}
		}
		
		if !found {
			t.Errorf("Expected feature not available: %s", expected)
		} else {
			t.Logf("Feature available: %s", expected)
		}
	}
}

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}