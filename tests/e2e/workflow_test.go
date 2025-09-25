package e2e

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/dmisiuk/acousticalc/pkg/calculator"
)

// WorkflowTestSuite represents a complete E2E test suite for application workflows
type WorkflowTestSuite struct {
	testEnvironment *TestEnvironment
	recordingActive bool
}

// TestEnvironment manages the test environment setup and cleanup
type TestEnvironment struct {
	TempDir     string
	StartTime   time.Time
	TestContext context.Context
	Cancel      context.CancelFunc
}

// NewWorkflowTestSuite creates a new E2E workflow test suite
func NewWorkflowTestSuite(t *testing.T) *WorkflowTestSuite {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

	tempDir, err := os.MkdirTemp("", "e2e_workflow_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	testEnv := &TestEnvironment{
		TempDir:     tempDir,
		StartTime:   time.Now(),
		TestContext: ctx,
		Cancel:      cancel,
	}

	suite := &WorkflowTestSuite{
		testEnvironment: testEnv,
		recordingActive: os.Getenv("E2E_RECORDING") == "true",
	}

	t.Cleanup(func() {
		suite.cleanup(t)
	})

	return suite
}

// cleanup performs cleanup of test resources
func (wts *WorkflowTestSuite) cleanup(t *testing.T) {
	t.Helper()

	if wts.testEnvironment != nil {
		wts.testEnvironment.Cancel()
	}

	if wts.testEnvironment.TempDir != "" {
		os.RemoveAll(wts.testEnvironment.TempDir)
	}
}

// TestCompleteCalculatorWorkflow tests the complete user journey for calculator operations
func TestCompleteCalculatorWorkflow(t *testing.T) {
	suite := NewWorkflowTestSuite(t)

	// Test scenarios representing complete user workflows
	workflows := []struct {
		name        string
		expression  string
		expected    float64
		expectError bool
		description string
	}{
		{
			name:        "BasicArithmetic",
			expression:  "2 + 3 * 4",
			expected:    14,
			expectError: false,
			description: "User performs basic arithmetic with precedence",
		},
		{
			name:        "ComplexExpression",
			expression:  "(10 + 5) * 2 - 8 / 4",
			expected:    28,
			expectError: false,
			description: "User performs complex expression with parentheses",
		},
		{
			name:        "DecimalCalculation",
			expression:  "3.14 * 2.5 + 1.86",
			expected:    9.71,
			expectError: false,
			description: "User performs decimal calculations",
		},
		{
			name:        "ErrorHandling",
			expression:  "10 / 0",
			expected:    0,
			expectError: true,
			description: "User encounters division by zero error",
		},
		{
			name:        "InvalidSyntax",
			expression:  "2 +* 3",
			expected:    0,
			expectError: true,
			description: "User encounters syntax error",
		},
	}

	for _, workflow := range workflows {
		t.Run(workflow.name, func(t *testing.T) {
			suite.runWorkflowScenario(t, workflow.expression, workflow.expected, workflow.expectError, workflow.description)
		})
	}
}

// runWorkflowScenario executes a complete workflow scenario
func (wts *WorkflowTestSuite) runWorkflowScenario(t *testing.T, expression string, expected float64, expectError bool, description string) {
	t.Helper()

	// Start recording if enabled
	if wts.recordingActive {
		t.Logf("Starting recording for scenario: %s", description)
	}

	// Execute the calculation workflow
	result, err := calculator.Evaluate(expression)

	// Validate the workflow outcome
	if expectError {
		if err == nil {
			t.Errorf("Expected error for expression '%s', but got result: %f", expression, result)
		} else {
			t.Logf("Expected error occurred for '%s': %v", expression, err)
		}
	} else {
		if err != nil {
			t.Errorf("Unexpected error for expression '%s': %v", expression, err)
		} else {
			// Use a small epsilon for floating point comparison
			const epsilon = 1e-9
			if abs(result-expected) > epsilon {
				t.Errorf("For expression '%s': expected %f, got %f", expression, expected, result)
			} else {
				t.Logf("Workflow successful for '%s': %f", expression, result)
			}
		}
	}

	// Stop recording if enabled
	if wts.recordingActive {
		t.Logf("Stopping recording for scenario: %s", description)
	}
}

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// TestUserJourneyIntegration tests multiple operations in sequence (user session simulation)
func TestUserJourneyIntegration(t *testing.T) {
	_ = NewWorkflowTestSuite(t) // Initialize but we don't need to store it

	// Simulate a user session with multiple calculations
	userSession := []struct {
		step       int
		expression string
		expected   float64
		note       string
	}{
		{1, "5 + 3", 8, "User starts with simple addition"},
		{2, "8 * 2", 16, "User multiplies previous result"},
		{3, "16 - 4", 12, "User subtracts from result"},
		{4, "12 / 3", 4, "User completes with division"},
		{5, "4 + 6 * 2", 16, "User tests operator precedence"},
		{6, "(4 + 6) * 2", 20, "User uses parentheses"},
	}

	for _, step := range userSession {
		t.Run(step.note, func(t *testing.T) {
			t.Logf("Step %d: %s - Expression: '%s'", step.step, step.note, step.expression)

			result, err := calculator.Evaluate(step.expression)
			if err != nil {
				t.Errorf("Step %d failed with error: %v", step.step, err)
				return
			}

			const epsilon = 1e-9
			if abs(result-step.expected) > epsilon {
				t.Errorf("Step %d: expected %f, got %f", step.step, step.expected, result)
			} else {
				t.Logf("Step %d completed successfully: %f", step.step, result)
			}
		})
	}
}

// TestPerformanceWorkflow tests application responsiveness under various conditions
func TestPerformanceWorkflow(t *testing.T) {
	_ = NewWorkflowTestSuite(t) // Initialize but we don't need to store it

	// Test performance with increasingly complex expressions
	performanceTests := []struct {
		name       string
		expression string
		maxTime    time.Duration
	}{
		{
			name:       "SimplePerformance",
			expression: "1 + 2 + 3 + 4 + 5",
			maxTime:    1 * time.Millisecond,
		},
		{
			name:       "MediumComplexity",
			expression: "(1 + 2) * (3 + 4) + (5 - 2) * (6 / 3)",
			maxTime:    5 * time.Millisecond,
		},
		{
			name:       "HighComplexity",
			expression: "((1 + 2) * (3 + 4) + (5 - 2)) * ((6 / 3) + (7 - 1)) - ((8 * 2) / (4 + 4))",
			maxTime:    10 * time.Millisecond,
		},
	}

	for _, test := range performanceTests {
		t.Run(test.name, func(t *testing.T) {
			start := time.Now()

			_, err := calculator.Evaluate(test.expression)
			if err != nil {
				t.Errorf("Performance test failed with error: %v", err)
				return
			}

			elapsed := time.Since(start)
			if elapsed > test.maxTime {
				t.Errorf("Performance test '%s' took %v, expected less than %v", test.name, elapsed, test.maxTime)
			} else {
				t.Logf("Performance test '%s' completed in %v (limit: %v)", test.name, elapsed, test.maxTime)
			}
		})
	}
}

// TestErrorRecoveryWorkflow tests graceful error handling and recovery
func TestErrorRecoveryWorkflow(t *testing.T) {
	_ = NewWorkflowTestSuite(t) // Initialize but we don't need to store it

	// Test various error scenarios and recovery
	errorTests := []struct {
		name        string
		expression  string
		expectError bool
		errorType   string
	}{
		{
			name:        "DivisionByZero",
			expression:  "10 / 0",
			expectError: true,
			errorType:   "division by zero",
		},
		{
			name:        "InvalidSyntax",
			expression:  "2 ++ 3",
			expectError: true,
			errorType:   "syntax error",
		},
		{
			name:        "UnmatchedParentheses",
			expression:  "(2 + 3",
			expectError: true,
			errorType:   "unmatched parentheses",
		},
		{
			name:        "EmptyExpression",
			expression:  "",
			expectError: true,
			errorType:   "empty expression",
		},
		{
			name:        "RecoveryAfterError",
			expression:  "2 + 3",
			expectError: false,
			errorType:   "none",
		},
	}

	for _, test := range errorTests {
		t.Run(test.name, func(t *testing.T) {
			result, err := calculator.Evaluate(test.expression)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error for '%s', but got result: %f", test.expression, result)
				} else {
					t.Logf("Expected error occurred for '%s': %v", test.expression, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for '%s': %v", test.expression, err)
				} else {
					t.Logf("Recovery successful for '%s': %f", test.expression, result)
				}
			}
		})
	}
}
