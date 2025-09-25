package e2e

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/dmisiuk/acousticalc/pkg/calculator"
)

// E2ETestConfig holds configuration for E2E tests
type E2ETestConfig struct {
	Timeout        time.Duration
	RecordingDir   string
	ArtifactsDir   string
	Platform       string
	Verbose        bool
	EnableRecording bool
}

// DefaultE2EConfig returns default configuration for E2E tests
func DefaultE2EConfig() *E2ETestConfig {
	return &E2ETestConfig{
		Timeout:        30 * time.Second,
		RecordingDir:   "tests/artifacts/recordings",
		ArtifactsDir:   "tests/artifacts/e2e",
		Platform:       runtime.GOOS,
		Verbose:        false,
		EnableRecording: true,
	}
}

// E2ETestSuite manages E2E test execution
type E2ETestSuite struct {
	config *E2ETestConfig
	calc   *calculator.Calculator
}

// NewE2ETestSuite creates a new E2E test suite
func NewE2ETestSuite(config *E2ETestConfig) *E2ETestSuite {
	if config == nil {
		config = DefaultE2EConfig()
	}
	
	// Ensure directories exist
	os.MkdirAll(config.RecordingDir, 0755)
	os.MkdirAll(config.ArtifactsDir, 0755)
	
	return &E2ETestSuite{
		config: config,
		calc:   calculator.NewCalculator(),
	}
}

// TestCompleteCalculatorWorkflow tests the complete user journey through the calculator
func TestCompleteCalculatorWorkflow(t *testing.T) {
	config := DefaultE2EConfig()
	config.Verbose = true
	
	suite := NewE2ETestSuite(config)
	
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	
	t.Run("BasicArithmeticWorkflow", func(t *testing.T) {
		suite.testBasicArithmeticWorkflow(ctx, t)
	})
	
	t.Run("ComplexExpressionWorkflow", func(t *testing.T) {
		suite.testComplexExpressionWorkflow(ctx, t)
	})
	
	t.Run("ErrorHandlingWorkflow", func(t *testing.T) {
		suite.testErrorHandlingWorkflow(ctx, t)
	})
	
	t.Run("PerformanceWorkflow", func(t *testing.T) {
		suite.testPerformanceWorkflow(ctx, t)
	})
}

// testBasicArithmeticWorkflow tests basic arithmetic operations
func (suite *E2ETestSuite) testBasicArithmeticWorkflow(ctx context.Context, t *testing.T) {
	operations := []struct {
		name     string
		expr     string
		expected float64
	}{
		{"Addition", "2 + 3", 5.0},
		{"Subtraction", "10 - 4", 6.0},
		{"Multiplication", "3 * 4", 12.0},
		{"Division", "15 / 3", 5.0},
		{"Mixed Operations", "2 + 3 * 4", 14.0},
		{"Parentheses", "(2 + 3) * 4", 20.0},
	}
	
	for _, op := range operations {
		t.Run(op.name, func(t *testing.T) {
			start := time.Now()
			
			result, err := suite.calc.Evaluate(op.expr)
			duration := time.Since(start)
			
			if err != nil {
				t.Errorf("Operation %s failed: %v", op.name, err)
				return
			}
			
			if result != op.expected {
				t.Errorf("Operation %s: expected %v, got %v", op.name, op.expected, result)
			}
			
			// Record performance metrics
			if suite.config.Verbose {
				t.Logf("Operation %s completed in %v", op.name, duration)
			}
			
			// Validate performance (should complete within 100ms)
			if duration > 100*time.Millisecond {
				t.Errorf("Operation %s took too long: %v", op.name, duration)
			}
		})
	}
}

// testComplexExpressionWorkflow tests complex mathematical expressions
func (suite *E2ETestSuite) testComplexExpressionWorkflow(ctx context.Context, t *testing.T) {
	complexExpressions := []struct {
		name     string
		expr     string
		expected float64
	}{
		{"Nested Parentheses", "((2 + 3) * 4) - (5 / 2)", 17.5},
		{"Multiple Operations", "2 * 3 + 4 * 5 - 6 / 2", 23.0},
		{"Decimal Operations", "3.14 * 2.0 + 1.86", 8.14},
		{"Large Numbers", "1000 * 1000 + 500 * 500", 1250000.0},
	}
	
	for _, expr := range complexExpressions {
		t.Run(expr.name, func(t *testing.T) {
			start := time.Now()
			
			result, err := suite.calc.Evaluate(expr.expr)
			duration := time.Since(start)
			
			if err != nil {
				t.Errorf("Complex expression %s failed: %v", expr.name, err)
				return
			}
			
			if result != expr.expected {
				t.Errorf("Complex expression %s: expected %v, got %v", expr.name, expr.expected, result)
			}
			
			if suite.config.Verbose {
				t.Logf("Complex expression %s completed in %v", expr.name, duration)
			}
			
			// Validate performance (should complete within 200ms)
			if duration > 200*time.Millisecond {
				t.Errorf("Complex expression %s took too long: %v", expr.name, duration)
			}
		})
	}
}

// testErrorHandlingWorkflow tests error handling scenarios
func (suite *E2ETestSuite) testErrorHandlingWorkflow(ctx context.Context, t *testing.T) {
	errorCases := []struct {
		name        string
		expr        string
		expectError bool
		errorType   string
	}{
		{"Division by Zero", "10 / 0", true, "division by zero"},
		{"Invalid Syntax", "2 + + 3", true, "syntax error"},
		{"Unmatched Parentheses", "(2 + 3", true, "unmatched parentheses"},
		{"Empty Expression", "", true, "empty expression"},
		{"Invalid Characters", "2 + abc", true, "invalid character"},
	}
	
	for _, testCase := range errorCases {
		t.Run(testCase.name, func(t *testing.T) {
			start := time.Now()
			
			result, err := suite.calc.Evaluate(testCase.expr)
			duration := time.Since(start)
			
			if testCase.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, but got result: %v", testCase.name, result)
				} else {
					if suite.config.Verbose {
						t.Logf("Correctly caught error for %s: %v", testCase.name, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for %s: %v", testCase.name, err)
				}
			}
			
			if suite.config.Verbose {
				t.Logf("Error handling test %s completed in %v", testCase.name, duration)
			}
			
			// Error handling should be fast (within 50ms)
			if duration > 50*time.Millisecond {
				t.Errorf("Error handling test %s took too long: %v", testCase.name, duration)
			}
		})
	}
}

// testPerformanceWorkflow tests performance characteristics
func (suite *E2ETestSuite) testPerformanceWorkflow(ctx context.Context, t *testing.T) {
	// Test rapid successive operations
	start := time.Now()
	
	for i := 0; i < 100; i++ {
		expr := fmt.Sprintf("%d + %d", i, i+1)
		result, err := suite.calc.Evaluate(expr)
		if err != nil {
			t.Errorf("Performance test failed at iteration %d: %v", i, err)
			return
		}
		expected := float64(i + i + 1)
		if result != expected {
			t.Errorf("Performance test failed at iteration %d: expected %v, got %v", i, expected, result)
			return
		}
	}
	
	totalDuration := time.Since(start)
	avgDuration := totalDuration / 100
	
	if suite.config.Verbose {
		t.Logf("Performance test: 100 operations completed in %v (avg: %v per operation)", totalDuration, avgDuration)
	}
	
	// Average operation should complete within 10ms
	if avgDuration > 10*time.Millisecond {
		t.Errorf("Performance test: average operation took too long: %v", avgDuration)
	}
	
	// Total test should complete within 1 second
	if totalDuration > 1*time.Second {
		t.Errorf("Performance test: total duration too long: %v", totalDuration)
	}
}

// TestE2EFrameworkIntegration tests integration with existing test infrastructure
func TestE2EFrameworkIntegration(t *testing.T) {
	config := DefaultE2EConfig()
	suite := NewE2ETestSuite(config)
	
	// Test that E2E framework can work with existing visual testing
	t.Run("VisualIntegration", func(t *testing.T) {
		// This test ensures E2E framework doesn't conflict with visual testing
		// from Story 0.2.2
		if suite.config.EnableRecording {
			// Verify recording directory exists and is writable
			if _, err := os.Stat(suite.config.RecordingDir); os.IsNotExist(err) {
				t.Errorf("Recording directory does not exist: %s", suite.config.RecordingDir)
			}
		}
	})
	
	// Test artifact management integration
	t.Run("ArtifactIntegration", func(t *testing.T) {
		// Verify artifact directory exists and is writable
		if _, err := os.Stat(suite.config.ArtifactsDir); os.IsNotExist(err) {
			t.Errorf("Artifacts directory does not exist: %s", suite.config.ArtifactsDir)
		}
		
		// Test artifact creation
		testFile := filepath.Join(suite.config.ArtifactsDir, "e2e_test_artifact.txt")
		content := fmt.Sprintf("E2E test artifact created at %s", time.Now().Format(time.RFC3339))
		
		err := os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			t.Errorf("Failed to create test artifact: %v", err)
		}
		
		// Clean up
		defer os.Remove(testFile)
	})
}

// TestE2ETestDiscovery tests E2E test discovery mechanism
func TestE2ETestDiscovery(t *testing.T) {
	// Test that E2E tests can be discovered and executed
	t.Run("TestDiscovery", func(t *testing.T) {
		// This test validates that the E2E test framework can be discovered
		// by the Go testing framework
		config := DefaultE2EConfig()
		suite := NewE2ETestSuite(config)
		
		if suite == nil {
			t.Error("Failed to create E2E test suite")
		}
		
		if suite.config == nil {
			t.Error("E2E test suite configuration is nil")
		}
		
		if suite.calc == nil {
			t.Error("E2E test suite calculator is nil")
		}
	})
	
	// Test configuration validation
	t.Run("ConfigurationValidation", func(t *testing.T) {
		config := DefaultE2EConfig()
		
		// Validate timeout
		if config.Timeout <= 0 {
			t.Error("Invalid timeout configuration")
		}
		
		// Validate platform detection
		if config.Platform == "" {
			t.Error("Platform not detected")
		}
		
		// Validate directories
		if config.RecordingDir == "" {
			t.Error("Recording directory not configured")
		}
		
		if config.ArtifactsDir == "" {
			t.Error("Artifacts directory not configured")
		}
	})
}

// TestE2ETestIsolation tests E2E test isolation and cleanup
func TestE2ETestIsolation(t *testing.T) {
	config := DefaultE2EConfig()
	suite := NewE2ETestSuite(config)
	
	// Test that each test run is isolated
	t.Run("TestIsolation", func(t *testing.T) {
		// Create a test artifact
		testFile := filepath.Join(suite.config.ArtifactsDir, "isolation_test.txt")
		content := "isolation test content"
		
		err := os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			t.Errorf("Failed to create isolation test file: %v", err)
		}
		
		// Verify file exists
		if _, err := os.Stat(testFile); os.IsNotExist(err) {
			t.Error("Isolation test file was not created")
		}
		
		// Clean up
		err = os.Remove(testFile)
		if err != nil {
			t.Errorf("Failed to clean up isolation test file: %v", err)
		}
		
		// Verify file is gone
		if _, err := os.Stat(testFile); err == nil {
			t.Error("Isolation test file was not cleaned up")
		}
	})
	
	// Test calculator state isolation
	t.Run("CalculatorStateIsolation", func(t *testing.T) {
		// Each test should start with a clean calculator state
		calc1 := calculator.NewCalculator()
		calc2 := calculator.NewCalculator()
		
		// Both calculators should be independent
		result1, err1 := calc1.Evaluate("2 + 3")
		result2, err2 := calc2.Evaluate("4 + 5")
		
		if err1 != nil || err2 != nil {
			t.Error("Calculator state isolation test failed")
		}
		
		if result1 == result2 {
			t.Error("Calculators are not properly isolated")
		}
	})
}

// BenchmarkE2EOperations benchmarks E2E operations
func BenchmarkE2EOperations(b *testing.B) {
	config := DefaultE2EConfig()
	suite := NewE2ETestSuite(config)
	
	operations := []string{
		"2 + 3",
		"10 - 4",
		"3 * 4",
		"15 / 3",
		"(2 + 3) * 4",
		"2 * 3 + 4 * 5",
	}
	
	for _, op := range operations {
		b.Run(fmt.Sprintf("Operation_%s", strings.ReplaceAll(op, " ", "_")), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := suite.calc.Evaluate(op)
				if err != nil {
					b.Fatalf("Benchmark operation failed: %v", err)
				}
			}
		})
	}
}