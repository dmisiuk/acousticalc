package integration

import (
	"fmt"
	"math"
	"os"
	"testing"
)

// MockCalculator simulates a calculator for integration testing
type MockCalculator struct {
	results map[string]float64
	errors  map[string]error
}

// NewMockCalculator creates a new mock calculator for testing
func NewMockCalculator() *MockCalculator {
	return &MockCalculator{
		results: make(map[string]float64),
		errors:  make(map[string]error),
	}
}

// SetResult sets a predefined result for a given expression
func (m *MockCalculator) SetResult(expression string, result float64) {
	m.results[expression] = result
}

// SetError sets a predefined error for a given expression
func (m *MockCalculator) SetError(expression string, err error) {
	m.errors[expression] = err
}

// Evaluate simulates calculator evaluation with mock data
func (m *MockCalculator) Evaluate(expression string) (float64, error) {
	if result, exists := m.results[expression]; exists {
		return result, nil
	}
	if err, exists := m.errors[expression]; exists {
		return 0, err
	}
	// Default behavior: return a simple calculation
	return simpleMockEvaluate(expression)
}

// simpleMockEvaluate provides basic mock evaluation
func simpleMockEvaluate(expression string) (float64, error) {
	// Simple mock implementation for common cases
	switch expression {
	case "2 + 3":
		return 5, nil
	case "10 - 4":
		return 6, nil
	case "3 * 4":
		return 12, nil
	case "15 / 3":
		return 5, nil
	case "10 / 0":
		return 0, fmt.Errorf("division by zero")
	default:
		return 0, fmt.Errorf("mock: unknown expression: %s", expression)
	}
}

// TestDataProvider provides test data for integration scenarios
type TestDataProvider struct {
	validExpressions   []string
	invalidExpressions []string
	complexExpressions []string
}

// NewTestDataProvider creates a new test data provider
func NewTestDataProvider() *TestDataProvider {
	return &TestDataProvider{
		validExpressions: []string{
			"2 + 3",
			"10 - 4",
			"3 * 4",
			"15 / 3",
			"3.5 + 2.1",
			"-5 + 3",
			"2 * (3 + 4)",
			"(2 + 3) * 4",
		},
		invalidExpressions: []string{
			"10 / 0",
			"2 +",
			"2 + a",
			"",
			"(2 + 3",
		},
		complexExpressions: []string{
			"2 * (3 + 4) - 5 / 2",
			"((2 + 3) * 4) - 5",
			"10 + 5 - 3 * 2 / 4",
			"-3 * (2 + -4) - -5",
		},
	}
}

// GetValidExpressions returns valid test expressions
func (p *TestDataProvider) GetValidExpressions() []string {
	return p.validExpressions
}

// GetInvalidExpressions returns invalid test expressions
func (p *TestDataProvider) GetInvalidExpressions() []string {
	return p.invalidExpressions
}

// GetComplexExpressions returns complex test expressions
func (p *TestDataProvider) GetComplexExpressions() []string {
	return p.complexExpressions
}

// PerformanceBenchmark provides performance testing utilities
type PerformanceBenchmark struct {
	thresholds map[string]timeThreshold
}

type timeThreshold struct {
	maxDurationMs int64
	warningMs     int64
}

// NewPerformanceBenchmark creates a new performance benchmark utility
func NewPerformanceBenchmark() *PerformanceBenchmark {
	return &PerformanceBenchmark{
		thresholds: map[string]timeThreshold{
			"simple_operation": {maxDurationMs: 1, warningMs: 0},
			"complex_operation": {maxDurationMs: 5, warningMs: 2},
			"error_handling":    {maxDurationMs: 1, warningMs: 0},
			"data_access":       {maxDurationMs: 10, warningMs: 5},
		},
	}
}

// ValidatePerformance checks if an operation meets performance requirements
func (p *PerformanceBenchmark) ValidatePerformance(operation string, durationMs int64) error {
	threshold, exists := p.thresholds[operation]
	if !exists {
		return fmt.Errorf("unknown operation type: %s", operation)
	}

	if durationMs > threshold.maxDurationMs {
		return fmt.Errorf("performance threshold exceeded: %s took %dms (max: %dms)",
			operation, durationMs, threshold.maxDurationMs)
	}

	if durationMs > threshold.warningMs {
		fmt.Printf("Warning: %s performance is slow: %dms (warning threshold: %dms)\n",
			operation, durationMs, threshold.warningMs)
	}

	return nil
}

// TestEnvironment represents the test environment configuration
type TestEnvironment struct {
	IsUnix      bool
	IsCI        bool
	TestMode    string // "unit", "integration", "e2e"
	Coverage    bool
	Benchmark   bool
	Parallelism int
}

// NewTestEnvironment creates a new test environment
func NewTestEnvironment() *TestEnvironment {
	env := &TestEnvironment{
		IsUnix:      isUnixEnvironment(),
		IsCI:        isCIEnvironment(),
		TestMode:    "integration",
		Coverage:    false,
		Benchmark:   false,
		Parallelism: 1,
	}
	return env
}

// isUnixEnvironment checks if running on Unix-like system
func isUnixEnvironment() bool {
	// This would use runtime.GOOS in real implementation
	// For testing purposes, we'll assume Unix
	return true
}

// isCIEnvironment checks if running in CI environment
func isCIEnvironment() bool {
	// Check common CI environment variables
	ciVars := []string{"CI", "GITHUB_ACTIONS", "JENKINS_URL", "TRAVIS"}
	for _, v := range ciVars {
		if os.Getenv(v) != "" {
			return true
		}
	}
	return false
}

// ValidationResult represents the result of a validation test
type ValidationResult struct {
	IsValid      bool
	ErrorMessage string
	Details      map[string]interface{}
}

// NewValidationResult creates a new validation result
func NewValidationResult(isValid bool, message string) *ValidationResult {
	return &ValidationResult{
		IsValid:      isValid,
		ErrorMessage: message,
		Details:      make(map[string]interface{}),
	}
}

// AddDetail adds a detail to the validation result
func (v *ValidationResult) AddDetail(key string, value interface{}) {
	v.Details[key] = value
}

// TestSuite represents a collection of integration tests
type TestSuite struct {
	Name        string
	Description string
	Tests       []TestCase
	Setup       func() error
	Teardown    func() error
}

// TestCase represents a single integration test case
type TestCase struct {
	Name        string
	Description string
	TestFunc    func(t *testing.T) error
	Dependencies []string
	Tags        []string
}

// NewTestSuite creates a new test suite
func NewTestSuite(name, description string) *TestSuite {
	return &TestSuite{
		Name:        name,
		Description: description,
		Tests:       make([]TestCase, 0),
	}
}

// AddTest adds a test case to the suite
func (s *TestSuite) AddTest(test TestCase) {
	s.Tests = append(s.Tests, test)
}

// Run executes all tests in the suite
func (s *TestSuite) Run(t *testing.T) {
	if s.Setup != nil {
		if err := s.Setup(); err != nil {
			t.Fatalf("Test suite setup failed: %v", err)
		}
		defer func() {
			if s.Teardown != nil {
				if err := s.Teardown(); err != nil {
					t.Errorf("Test suite teardown failed: %v", err)
				}
			}
		}()
	}

	for _, test := range s.Tests {
		t.Run(test.Name, func(t *testing.T) {
			if err := test.TestFunc(t); err != nil {
				t.Errorf("Test %s failed: %v", test.Name, err)
			}
		})
	}
}

// MathUtilities provides mathematical utilities for testing
type MathUtilities struct{}

// NewMathUtilities creates a new math utilities instance
func NewMathUtilities() *MathUtilities {
	return &MathUtilities{}
}

// AlmostEqual checks if two floating point numbers are approximately equal
func (m *MathUtilities) AlmostEqual(a, b, epsilon float64) bool {
	// Handle the exact equality case (including both zero)
	if a == b {
		return true
	}

	// Calculate absolute difference
	diff := math.Abs(a - b)

	// For zero values, only use absolute comparison
	if a == 0 || b == 0 {
		return diff <= epsilon
	}

	// For non-zero values, use both absolute and relative comparison
	// The relative difference is often more meaningful for larger numbers
	relativeDiff := diff / math.Max(math.Abs(a), math.Abs(b))

	// Return true if either absolute or relative difference is within epsilon
	return diff <= epsilon || relativeDiff <= epsilon
}

// DefaultEpsilon returns the default epsilon for floating point comparison
func (m *MathUtilities) DefaultEpsilon() float64 {
	return 1e-9
}