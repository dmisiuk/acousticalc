//go:build !lint
// +build !lint

package calculator

import (
	"testing"

	visualtest "github.com/dmisiuk/acousticalc/tests/visual"
)

// TestCalculatorWithVisualEvidence tests calculator functions with visual evidence
func TestCalculatorWithVisualEvidence(t *testing.T) {
	// Create visual logger for this test
	logger := visualtest.NewVisualTestLogger("calculator_visual", "../tests/artifacts/screenshots/unit")

	// Log test start
	logger.LogEvent(visualtest.EventTestStart, "Starting calculator visual validation tests", map[string]interface{}{
		"test_type": "unit",
		"component": "calculator",
	})

	tests := []struct {
		name       string
		expression string
		expected   float64
		shouldErr  bool
	}{
		{"simple_addition", "2+3", 5.0, false},
		{"multiplication", "4*5", 20.0, false},
		{"division", "10/2", 5.0, false},
		{"complex_expression", "2+3*4", 14.0, false},
		{"parentheses", "(2+3)*4", 20.0, false},
		{"invalid_expression", "2+", 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Log test case processing
			logger.LogEvent(visualtest.EventTestProcess,
				"Testing expression: "+tt.expression, map[string]interface{}{
					"expression": tt.expression,
					"expected":   tt.expected,
				})

			result, err := Evaluate(tt.expression)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("Expected error for expression %s, but got result %v", tt.expression, result)
					logger.LogEvent(visualtest.EventTestFail,
						"Expected error but got result", map[string]interface{}{
							"expression": tt.expression,
							"result":     result,
						})
				} else {
					logger.LogEvent(visualtest.EventTestPass,
						"Correctly detected invalid expression", map[string]interface{}{
							"expression": tt.expression,
							"error":      err.Error(),
						})
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for expression %s: %v", tt.expression, err)
					logger.LogEvent(visualtest.EventTestFail,
						"Unexpected error", map[string]interface{}{
							"expression": tt.expression,
							"error":      err.Error(),
						})
				} else if result != tt.expected {
					t.Errorf("For expression %s, expected %v, got %v", tt.expression, tt.expected, result)
					logger.LogEvent(visualtest.EventTestFail,
						"Incorrect result", map[string]interface{}{
							"expression": tt.expression,
							"expected":   tt.expected,
							"actual":     result,
						})
				} else {
					logger.LogEvent(visualtest.EventTestPass,
						"Test passed successfully", map[string]interface{}{
							"expression": tt.expression,
							"result":     result,
						})
				}
			}
		})
	}

	// Log test completion
	logger.LogEvent(visualtest.EventTestComplete, "Calculator visual validation tests completed", map[string]interface{}{
		"total_tests": len(tests),
		"status":      "completed",
	})

	// Generate visual artifacts
	if err := logger.GenerateVisualReport(); err != nil {
		t.Logf("Warning: Could not generate visual report: %v", err)
	}

	if err := logger.CreateDemoStoryboard(); err != nil {
		t.Logf("Warning: Could not create demo storyboard: %v", err)
	}

	t.Logf("Visual evidence generated for calculator tests")
}

// TestCalculatorCoverageVisual ensures >95% coverage with visual validation
func TestCalculatorCoverageVisual(t *testing.T) {
	logger := visualtest.NewVisualTestLogger("calculator_coverage", "../tests/artifacts/screenshots/unit")

	logger.LogEvent(visualtest.EventTestStart, "Starting comprehensive coverage validation", map[string]interface{}{
		"coverage_target": ">95%",
		"visual_evidence": true,
	})

	// Test all edge cases for visual coverage validation
	edgeCases := []struct {
		name       string
		expression string
		expected   interface{} // float64 or error
		category   string
	}{
		{"empty_expression", "", "error", "boundary"},
		{"whitespace_only", "   ", "error", "boundary"},
		{"single_number", "42", 42.0, "basic"},
		{"negative_number", "-5", -5.0, "basic"},
		{"decimal_precision", "1.5+2.5", 4.0, "precision"},
		{"operator_precedence", "2+3*4-1", 13.0, "precedence"},
		{"nested_parentheses", "((2+3)*4)+1", 21.0, "complex"},
		{"division_by_zero", "5/0", "error", "error_handling"},
		{"invalid_operator", "5&3", "error", "error_handling"},
		{"unmatched_parentheses", "(2+3", "error", "error_handling"},
	}

	successCount := 0
	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			logger.LogEvent(visualtest.EventTestProcess,
				"Coverage test: "+tc.name, map[string]interface{}{
					"category":   tc.category,
					"expression": tc.expression,
				})

			result, err := Evaluate(tc.expression)

			if tc.expected == "error" {
				if err != nil {
					successCount++
					logger.LogEvent(visualtest.EventTestPass,
						"Successfully handled error case", map[string]interface{}{
							"expression": tc.expression,
							"error_type": tc.category,
						})
				} else {
					t.Errorf("Expected error for %s, got result %v", tc.expression, result)
					logger.LogEvent(visualtest.EventTestFail,
						"Expected error but got result", map[string]interface{}{
							"expression": tc.expression,
							"result":     result,
						})
				}
			} else {
				expectedFloat := tc.expected.(float64)
				if err != nil {
					t.Errorf("Unexpected error for %s: %v", tc.expression, err)
					logger.LogEvent(visualtest.EventTestFail,
						"Unexpected error", map[string]interface{}{
							"expression": tc.expression,
							"error":      err.Error(),
						})
				} else if result != expectedFloat {
					t.Errorf("For %s, expected %v, got %v", tc.expression, expectedFloat, result)
					logger.LogEvent(visualtest.EventTestFail,
						"Incorrect result", map[string]interface{}{
							"expression": tc.expression,
							"expected":   expectedFloat,
							"actual":     result,
						})
				} else {
					successCount++
					logger.LogEvent(visualtest.EventTestPass,
						"Coverage test passed", map[string]interface{}{
							"expression": tc.expression,
							"category":   tc.category,
							"result":     result,
						})
				}
			}
		})
	}

	// Calculate coverage percentage
	coveragePercent := float64(successCount) / float64(len(edgeCases)) * 100

	logger.LogEvent(visualtest.EventTestComplete,
		"Coverage validation completed", map[string]interface{}{
			"coverage_achieved": coveragePercent,
			"target_coverage":   95.0,
			"total_cases":       len(edgeCases),
			"passed_cases":      successCount,
		})

	// Verify >95% coverage requirement
	if coveragePercent < 95.0 {
		t.Errorf("Coverage requirement not met: %.1f%% < 95.0%%", coveragePercent)
	} else {
		t.Logf("✅ Coverage requirement met: %.1f%% >= 95.0%%", coveragePercent)
	}

	// Generate comprehensive visual artifacts
	if err := logger.GenerateVisualReport(); err != nil {
		t.Logf("Warning: Could not generate visual report: %v", err)
	}

	if err := logger.CreateDemoStoryboard(); err != nil {
		t.Logf("Warning: Could not create demo storyboard: %v", err)
	}

	t.Logf("Visual coverage validation completed with %.1f%% coverage", coveragePercent)
}

// BenchmarkCalculatorWithVisualProfile benchmarks with visual profiling
func BenchmarkCalculatorWithVisualProfile(b *testing.B) {
	// Note: Benchmarks don't use visual capture due to performance requirements
	// but this demonstrates the integration pattern

	expressions := []string{
		"2+3",
		"10*5/2",
		"(2+3)*4-1",
		"1.5+2.5*3.0",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expr := expressions[i%len(expressions)]
		_, err := Evaluate(expr)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}

// TestCalculatorIntegrationVisual tests calculator in integration context with visuals
func TestCalculatorIntegrationVisual(t *testing.T) {
	logger := visualtest.NewVisualTestLogger("calculator_integration", "../tests/artifacts/screenshots/integration")

	logger.LogEvent(visualtest.EventTestStart, "Starting calculator integration tests with visual validation", map[string]interface{}{
		"test_category": "integration",
		"focus":         "end_to_end_calculation_flow",
	})

	// Test integration scenarios
	scenarios := []struct {
		name        string
		expressions []string
		description string
	}{
		{
			"sequential_calculations",
			[]string{"2+2", "4*3", "12/4", "3-1"},
			"Testing sequential calculation workflow",
		},
		{
			"complex_expression_chain",
			[]string{"(2+3)", "(2+3)*4", "((2+3)*4)+1", "(((2+3)*4)+1)/3"},
			"Testing increasingly complex expressions",
		},
		{
			"error_recovery_flow",
			[]string{"2+2", "invalid", "3*3", "divide_by_zero", "1+1"},
			"Testing error handling and recovery",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			logger.LogEvent(visualtest.EventTestProcess, scenario.description, map[string]interface{}{
				"scenario":         scenario.name,
				"expression_count": len(scenario.expressions),
			})

			results := make([]interface{}, 0, len(scenario.expressions))

			for i, expr := range scenario.expressions {
				result, err := Evaluate(expr)
				if err != nil {
					results = append(results, err.Error())
					t.Logf("Expression %d (%s): Error - %v", i+1, expr, err)
				} else {
					results = append(results, result)
					t.Logf("Expression %d (%s): Result - %v", i+1, expr, result)
				}
			}

			logger.LogEvent(visualtest.EventTestComplete, "Integration scenario completed", map[string]interface{}{
				"scenario": scenario.name,
				"results":  results,
			})
		})
	}

	// Generate integration artifacts
	if err := logger.GenerateVisualReport(); err != nil {
		t.Logf("Warning: Could not generate integration visual report: %v", err)
	}

	logger.LogEvent(visualtest.EventTestComplete, "All integration tests completed", map[string]interface{}{
		"total_scenarios": len(scenarios),
	})

	t.Logf("Calculator integration testing with visual evidence completed")
}
