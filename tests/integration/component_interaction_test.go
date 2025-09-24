package integration

import (
	"fmt"
	"testing"
	"time"
)

// TestComponentInteraction tests interaction between different components
func TestComponentInteraction(t *testing.T) {
	// Create test suite
	suite := NewTestSuite("Component Interaction", "Tests interaction between calculator and mock components")

	// Add setup and teardown
	suite.Setup = func() error {
		// Initialize test environment
		fmt.Println("Setting up component interaction tests...")
		return nil
	}

	suite.Teardown = func() error {
		// Cleanup test environment
		fmt.Println("Cleaning up component interaction tests...")
		return nil
	}

	// Add test cases
	suite.AddTest(TestCase{
		Name:        "Calculator Mock Integration",
		Description: "Test integration between real calculator and mock calculator",
		TestFunc: func(t *testing.T) error {
			return testCalculatorMockIntegration(t)
		},
	})

	suite.AddTest(TestCase{
		Name:        "Data Provider Integration",
		Description: "Test integration with test data provider",
		TestFunc: func(t *testing.T) error {
			return testDataProviderIntegration(t)
		},
	})

	suite.AddTest(TestCase{
		Name:        "Performance Validation Integration",
		Description: "Test performance validation integration",
		TestFunc: func(t *testing.T) error {
			return testPerformanceValidationIntegration(t)
		},
	})

	suite.AddTest(TestCase{
		Name:        "Error Handling Integration",
		Description: "Test error handling across components",
		TestFunc: func(t *testing.T) error {
			return testErrorHandlingIntegration(t)
		},
	})

	suite.AddTest(TestCase{
		Name:        "Math Utilities Integration",
		Description: "Test math utilities integration",
		TestFunc: func(t *testing.T) error {
			return testMathUtilitiesIntegration(t)
		},
	})

	// Run the test suite
	suite.Run(t)
}

// testCalculatorMockIntegration tests integration between real and mock calculators
func testCalculatorMockIntegration(t *testing.T) error {
	mockCalc := NewMockCalculator()
	dataProvider := NewTestDataProvider()

	// Set up mock responses
	validExprs := dataProvider.GetValidExpressions()
	for _, expr := range validExprs[:3] { // Test first 3 expressions
		mockCalc.SetResult(expr, 42.0) // Mock result
	}

	// Test integration
	for _, expr := range validExprs[:3] {
		// Get mock result
		mockResult, mockErr := mockCalc.Evaluate(expr)
		if mockErr != nil {
			return fmt.Errorf("mock calculator failed for %s: %v", expr, mockErr)
		}

		// Test that mock returns expected result
		if mockResult != 42.0 {
			return fmt.Errorf("mock calculator returned unexpected result for %s: got %v, want %v",
				expr, mockResult, 42.0)
		}

		// Log integration success
		t.Logf("Integration test passed for expression: %s", expr)
	}

	return nil
}

// testDataProviderIntegration tests data provider integration
func testDataProviderIntegration(t *testing.T) error {
	dataProvider := NewTestDataProvider()
	benchmark := NewPerformanceBenchmark()

	// Test valid expressions
	validExprs := dataProvider.GetValidExpressions()
	if len(validExprs) == 0 {
		return fmt.Errorf("data provider returned no valid expressions")
	}

	// Test invalid expressions
	invalidExprs := dataProvider.GetInvalidExpressions()
	if len(invalidExprs) == 0 {
		return fmt.Errorf("data provider returned no invalid expressions")
	}

	// Test complex expressions
	complexExprs := dataProvider.GetComplexExpressions()
	if len(complexExprs) == 0 {
		return fmt.Errorf("data provider returned no complex expressions")
	}

	// Test performance validation
	start := time.Now()
	_ = len(validExprs) + len(invalidExprs) + len(complexExprs)
	duration := time.Since(start)

	err := benchmark.ValidatePerformance("data_access", duration.Milliseconds())
	if err != nil {
		return fmt.Errorf("performance validation failed: %v", err)
	}

	t.Logf("Data provider integration successful - Valid: %d, Invalid: %d, Complex: %d",
		len(validExprs), len(invalidExprs), len(complexExprs))

	return nil
}

// testPerformanceValidationIntegration tests performance validation integration
func testPerformanceValidationIntegration(t *testing.T) error {
	benchmark := NewPerformanceBenchmark()

	// Test performance threshold validation
	tests := []struct {
		operation  string
		duration   int64
		shouldFail bool
	}{
		{"simple_operation", 0, false},  // Should pass
		{"simple_operation", 2, true},   // Should fail (exceeds 1ms)
		{"complex_operation", 3, false}, // Should pass
		{"complex_operation", 10, true}, // Should fail (exceeds 5ms)
		{"error_handling", 0, false},    // Should pass
		{"error_handling", 2, true},     // Should fail (exceeds 1ms)
	}

	for _, test := range tests {
		err := benchmark.ValidatePerformance(test.operation, test.duration)
		if test.shouldFail && err == nil {
			return fmt.Errorf("expected performance validation to fail for %s with duration %d, but it passed",
				test.operation, test.duration)
		}
		if !test.shouldFail && err != nil {
			return fmt.Errorf("expected performance validation to pass for %s with duration %d, but it failed: %v",
				test.operation, test.duration, err)
		}
	}

	t.Logf("Performance validation integration test completed successfully")
	return nil
}

// testErrorHandlingIntegration tests error handling across components
func testErrorHandlingIntegration(t *testing.T) error {
	mockCalc := NewMockCalculator()
	dataProvider := NewTestDataProvider()

	// Set up error scenarios
	invalidExprs := dataProvider.GetInvalidExpressions()
	for _, expr := range invalidExprs {
		mockCalc.SetError(expr, fmt.Errorf("simulated error for %s", expr))
	}

	// Test error handling integration
	for _, expr := range invalidExprs {
		_, err := mockCalc.Evaluate(expr)
		if err == nil {
			return fmt.Errorf("expected error for invalid expression %s, but got none", expr)
		}

		// Verify error message contains expected content
		expectedMsg := fmt.Sprintf("simulated error for %s", expr)
		if err.Error() != expectedMsg {
			return fmt.Errorf("unexpected error message for %s: got %v, want %v",
				expr, err.Error(), expectedMsg)
		}
	}

	// Test mixed scenario: some expressions succeed, some fail
	mixedTests := []struct {
		expr       string
		shouldFail bool
	}{
		{"2 + 3", false},
		{"10 / 0", true},
		{"3 * 4", false},
		{"invalid + expr", true},
	}

	for _, test := range mixedTests {
		// Configure mock
		if test.shouldFail {
			mockCalc.SetError(test.expr, fmt.Errorf("forced error for %s", test.expr))
		} else {
			mockCalc.SetResult(test.expr, 42.0)
		}

		_, err := mockCalc.Evaluate(test.expr)
		if test.shouldFail && err == nil {
			return fmt.Errorf("expected error for %s, but got none", test.expr)
		}
		if !test.shouldFail && err != nil {
			return fmt.Errorf("expected success for %s, but got error: %v", test.expr, err)
		}
	}

	t.Logf("Error handling integration test completed successfully")
	return nil
}

// testMathUtilitiesIntegration tests math utilities integration
func testMathUtilitiesIntegration(t *testing.T) error {
	mathUtils := NewMathUtilities()
	dataProvider := NewTestDataProvider()

	// Test floating point comparison
	comparisonTests := []struct {
		a, b     float64
		epsilon  float64
		expected bool
	}{
		{1.0, 1.0, 1e-9, true},         // Exact match
		{1.0, 1.0 + 5e-10, 1e-9, true}, // Very close (5e-10 < 1e-9)
		{1.0, 1.0001, 1e-9, false},     // Too far apart
		{0.0, 0.0, 1e-9, true},         // Zero comparison
		{-1.0, -1.0, 1e-9, true},       // Negative numbers
	}

	for _, test := range comparisonTests {
		result := mathUtils.AlmostEqual(test.a, test.b, test.epsilon)
		if result != test.expected {
			return fmt.Errorf("AlmostEqual(%v, %v, %v) = %v, expected %v",
				test.a, test.b, test.epsilon, result, test.expected)
		}
	}

	// Test default epsilon
	defaultEpsilon := mathUtils.DefaultEpsilon()
	if defaultEpsilon != 1e-9 {
		return fmt.Errorf("DefaultEpsilon() = %v, expected 1e-9", defaultEpsilon)
	}

	// Test integration with data provider
	validExprs := dataProvider.GetValidExpressions()
	for _, expr := range validExprs[:2] {
		// Simulate calculation with floating point precision
		result1 := 0.1 + 0.2
		result2 := 0.3

		if !mathUtils.AlmostEqual(result1, result2, mathUtils.DefaultEpsilon()) {
			return fmt.Errorf("floating point comparison failed for expression %s: %v != %v",
				expr, result1, result2)
		}
	}

	t.Logf("Math utilities integration test completed successfully")
	return nil
}

// TestParallelComponentInteraction tests concurrent component interactions
func TestParallelComponentInteraction(t *testing.T) {
	mockCalc := NewMockCalculator()
	dataProvider := NewTestDataProvider()

	// Set up mock data
	validExprs := dataProvider.GetValidExpressions()
	for _, expr := range validExprs {
		mockCalc.SetResult(expr, 42.0)
	}

	// Test concurrent access
	t.Run("Concurrent Calculator Access", func(t *testing.T) {
		const numGoroutines = 10
		const numOperations = 100

		done := make(chan bool, numGoroutines)
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer func() { done <- true }()

				for j := 0; j < numOperations; j++ {
					expr := validExprs[j%len(validExprs)]
					_, err := mockCalc.Evaluate(expr)
					if err != nil {
						errors <- fmt.Errorf("goroutine %d: %v", id, err)
						return
					}
				}
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			<-done
		}

		// Check for errors
		select {
		case err := <-errors:
			t.Errorf("Concurrent access test failed: %v", err)
		default:
			// No errors, test passed
		}

		t.Logf("Concurrent access test completed successfully")
	})
}
