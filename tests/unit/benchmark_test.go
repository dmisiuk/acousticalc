package unit

import (
	"testing"

	"github.com/dmisiuk/acousticalc/pkg/calculator"
)

// BenchmarkBasicOperations benchmarks basic arithmetic operations
func BenchmarkBasicOperations(b *testing.B) {
	operations := []struct {
		name string
		expr string
	}{
		{"Addition", "2 + 3"},
		{"Subtraction", "10 - 4"},
		{"Multiplication", "3 * 4"},
		{"Division", "15 / 3"},
		{"Negative addition", "-2 + 3"},
		{"Decimal addition", "3.5 + 2.1"},
	}

	for _, op := range operations {
		b.Run(op.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := calculator.Evaluate(op.expr)
				if err != nil {
					b.Fatalf("Benchmark failed for %s: %v", op.name, err)
				}
			}
		})
	}
}

// BenchmarkComplexExpressions benchmarks complex mathematical expressions
func BenchmarkComplexExpressions(b *testing.B) {
	expressions := []struct {
		name string
		expr string
	}{
		{"Operator precedence", "2 + 3 * 4"},
		{"Parentheses", "(2 + 3) * 4"},
		{"Nested parentheses", "((2 + 3) * 4) - 5"},
		{"Complex mixed", "2 * (3 + 4) - 5 / 2"},
		{"Multiple operators", "10 + 5 - 3 * 2 / 4"},
		{"Complex negatives", "-3 * (2 + -4) - -5"},
	}

	for _, expr := range expressions {
		b.Run(expr.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := calculator.Evaluate(expr.expr)
				if err != nil {
					b.Fatalf("Benchmark failed for %s: %v", expr.name, err)
				}
			}
		})
	}
}

// BenchmarkEdgeCases benchmarks edge case scenarios
func BenchmarkEdgeCases(b *testing.B) {
	edgeCases := []struct {
		name string
		expr string
	}{
		{"Single number", "42"},
		{"Zero operations", "0 + 0"},
		{"Large numbers", "1000000 * 2"},
		{"Very small decimal", "0.000001 * 1000"},
		{"Negative single", "-42"},
	}

	for _, edge := range edgeCases {
		b.Run(edge.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := calculator.Evaluate(edge.expr)
				if err != nil {
					b.Fatalf("Benchmark failed for %s: %v", edge.name, err)
				}
			}
		})
	}
}

// BenchmarkErrorHandling benchmarks error handling performance
func BenchmarkErrorHandling(b *testing.B) {
	errorCases := []struct {
		name string
		expr string
	}{
		{"Division by zero", "10 / 0"},
		{"Invalid syntax", "2 +"},
		{"Invalid character", "2 + a"},
		{"Empty expression", ""},
		{"Mismatched parentheses", "(2 + 3"},
	}

	for _, errCase := range errorCases {
		b.Run(errCase.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = calculator.Evaluate(errCase.expr) // Expect errors, ignore them
			}
		})
	}
}

// BenchmarkParallel benchmarks concurrent calculator operations
func BenchmarkParallel(b *testing.B) {
	expr := "2 * (3 + 4) - 5 / 2"

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := calculator.Evaluate(expr)
			if err != nil {
				b.Fatalf("Parallel benchmark failed: %v", err)
			}
		}
	})
}

// BenchmarkMemoryUsage measures memory allocation patterns
func BenchmarkMemoryUsage(b *testing.B) {
	var result float64
	var err error

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Test a variety of expressions to exercise different memory patterns
		exprs := []string{
			"2 + 3",
			"2 * (3 + 4) - 5 / 2",
			"3.14159 * 2.71828",
			"((2 + 3) * 4) - 5",
		}

		for _, expr := range exprs {
			result, err = calculator.Evaluate(expr)
			if err != nil {
				b.Fatalf("Memory benchmark failed: %v", err)
			}
		}
	}

	// Use result to prevent compiler optimization
	_ = result
}

// BenchmarkCalculatorPerformance is a comprehensive performance test
// that measures overall calculator performance across different scenarios
func BenchmarkCalculatorPerformance(b *testing.B) {
	// Warm up
	for i := 0; i < 100; i++ {
		_, _ = calculator.Evaluate("2 + 3")
	}

	b.Run("Overall performance", func(b *testing.B) {
		scenarios := []string{
			"2 + 3",                                 // Simple
			"2 * (3 + 4) - 5 / 2",                   // Complex
			"3.14159 * 2.71828 + 1.41421 - 0.57721", // Decimals
			"((2 + 3) * 4) - (6 / 2) + 1",           // Nested
		}

		for i := 0; i < b.N; i++ {
			for _, scenario := range scenarios {
				_, err := calculator.Evaluate(scenario)
				if err != nil {
					b.Fatalf("Performance benchmark failed: %v", err)
				}
			}
		}
	})
}
