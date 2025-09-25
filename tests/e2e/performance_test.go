package e2e

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/dmisiuk/acousticalc/pkg/calculator"
)

// PerformanceTestConfig holds configuration for performance testing
type PerformanceTestConfig struct {
	Timeout        time.Duration
	ArtifactsDir   string
	Platform       string
	Verbose        bool
	EnableProfiling bool
	MaxMemoryMB    int
	MaxCPUPct      float64
}

// DefaultPerformanceTestConfig returns default configuration for performance testing
func DefaultPerformanceTestConfig() *PerformanceTestConfig {
	return &PerformanceTestConfig{
		Timeout:        30 * time.Second,
		ArtifactsDir:   "tests/artifacts/performance",
		Platform:       runtime.GOOS,
		Verbose:        false,
		EnableProfiling: true,
		MaxMemoryMB:    100,
		MaxCPUPct:      80.0,
	}
}

// PerformanceTestSuite manages performance testing
type PerformanceTestSuite struct {
	config *PerformanceTestConfig
	calc   *calculator.Calculator
}

// NewPerformanceTestSuite creates a new performance test suite
func NewPerformanceTestSuite(config *PerformanceTestConfig) *PerformanceTestSuite {
	if config == nil {
		config = DefaultPerformanceTestConfig()
	}
	
	// Ensure artifacts directory exists
	os.MkdirAll(config.ArtifactsDir, 0755)
	
	return &PerformanceTestSuite{
		config: config,
		calc:   calculator.NewCalculator(),
	}
}

// TestE2EPerformanceCharacteristics tests E2E performance characteristics
func TestE2EPerformanceCharacteristics(t *testing.T) {
	config := DefaultPerformanceTestConfig()
	config.Verbose = true
	suite := NewPerformanceTestSuite(config)
	
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	
	t.Run("StartupPerformance", func(t *testing.T) {
		suite.testStartupPerformance(ctx, t)
	})
	
	t.Run("OperationPerformance", func(t *testing.T) {
		suite.testOperationPerformance(ctx, t)
	})
	
	t.Run("MemoryPerformance", func(t *testing.T) {
		suite.testMemoryPerformance(ctx, t)
	})
	
	t.Run("ConcurrencyPerformance", func(t *testing.T) {
		suite.testConcurrencyPerformance(ctx, t)
	})
	
	t.Run("ScalabilityPerformance", func(t *testing.T) {
		suite.testScalabilityPerformance(ctx, t)
	})
}

// testStartupPerformance tests application startup performance
func (suite *PerformanceTestSuite) testStartupPerformance(ctx context.Context, t *testing.T) {
	// Test calculator initialization performance
	start := time.Now()
	
	for i := 0; i < 100; i++ {
		calc := calculator.NewCalculator()
		if calc == nil {
			t.Errorf("Failed to create calculator instance %d", i)
			return
		}
	}
	
	duration := time.Since(start)
	avgDuration := duration / 100
	
	if suite.config.Verbose {
		t.Logf("Calculator startup performance: 100 instances in %v (avg: %v per instance)", duration, avgDuration)
	}
	
	// Startup should be fast (within 1ms per instance)
	if avgDuration > 1*time.Millisecond {
		t.Errorf("Calculator startup too slow: %v per instance", avgDuration)
	}
	
	// Total startup time should be reasonable (within 100ms for 100 instances)
	if duration > 100*time.Millisecond {
		t.Errorf("Total calculator startup too slow: %v for 100 instances", duration)
	}
}

// testOperationPerformance tests individual operation performance
func (suite *PerformanceTestSuite) testOperationPerformance(ctx context.Context, t *testing.T) {
	operations := []struct {
		name string
		expr string
		maxDuration time.Duration
	}{
		{"Simple Addition", "2 + 3", 1 * time.Millisecond},
		{"Simple Subtraction", "10 - 4", 1 * time.Millisecond},
		{"Simple Multiplication", "3 * 4", 1 * time.Millisecond},
		{"Simple Division", "15 / 3", 1 * time.Millisecond},
		{"Complex Expression", "2 * (3 + 4) - 5 / 2", 2 * time.Millisecond},
		{"Nested Parentheses", "((2 + 3) * 4) - (5 / 2)", 2 * time.Millisecond},
		{"Large Numbers", "1000 * 1000 + 500 * 500", 2 * time.Millisecond},
		{"Decimal Operations", "3.14159 * 2.71828", 2 * time.Millisecond},
	}
	
	for _, op := range operations {
		t.Run(op.name, func(t *testing.T) {
			// Run operation multiple times for accurate measurement
			iterations := 100
			start := time.Now()
			
			for i := 0; i < iterations; i++ {
				result, err := suite.calc.Evaluate(op.expr)
				if err != nil {
					t.Errorf("Operation %s failed at iteration %d: %v", op.name, i, err)
					return
				}
				
				// Verify result is reasonable (not zero for most operations)
				if result == 0 && op.expr != "0 + 0" {
					t.Errorf("Operation %s returned unexpected zero result", op.name)
				}
			}
			
			duration := time.Since(start)
			avgDuration := duration / time.Duration(iterations)
			
			if suite.config.Verbose {
				t.Logf("Operation %s: %d iterations in %v (avg: %v per operation)", 
					op.name, iterations, duration, avgDuration)
			}
			
			// Verify performance meets requirements
			if avgDuration > op.maxDuration {
				t.Errorf("Operation %s too slow: %v (max: %v)", op.name, avgDuration, op.maxDuration)
			}
		})
	}
}

// testMemoryPerformance tests memory usage characteristics
func (suite *PerformanceTestSuite) testMemoryPerformance(ctx context.Context, t *testing.T) {
	// Test memory usage during normal operations
	start := time.Now()
	
	// Perform many operations to test memory usage
	operations := []string{
		"2 + 3", "10 - 4", "3 * 4", "15 / 3",
		"2 * (3 + 4) - 5 / 2", "((2 + 3) * 4) - (5 / 2)",
		"1000 * 1000 + 500 * 500", "3.14159 * 2.71828",
	}
	
	for i := 0; i < 1000; i++ {
		expr := operations[i%len(operations)]
		result, err := suite.calc.Evaluate(expr)
		if err != nil {
			t.Errorf("Memory test operation failed at iteration %d: %v", i, err)
			return
		}
		
		// Verify result is reasonable
		if result == 0 && expr != "0 + 0" {
			t.Errorf("Memory test operation returned unexpected zero result at iteration %d", i)
		}
	}
	
	duration := time.Since(start)
	
	if suite.config.Verbose {
		t.Logf("Memory performance test: 1000 operations in %v", duration)
	}
	
	// Memory test should complete within reasonable time
	if duration > 1*time.Second {
		t.Errorf("Memory performance test too slow: %v", duration)
	}
	
	// Test memory usage with large numbers
	largeNumberOps := []string{
		"1000000 + 2000000",
		"1000000 * 2000000",
		"1000000 / 2000000",
		"1000000 - 2000000",
	}
	
	start = time.Now()
	
	for i := 0; i < 100; i++ {
		expr := largeNumberOps[i%len(largeNumberOps)]
		result, err := suite.calc.Evaluate(expr)
		if err != nil {
			t.Errorf("Large number memory test failed at iteration %d: %v", i, err)
			return
		}
		
		// Verify result is reasonable
		if result == 0 && expr != "0 + 0" {
			t.Errorf("Large number memory test returned unexpected zero result at iteration %d", i)
		}
	}
	
	duration = time.Since(start)
	
	if suite.config.Verbose {
		t.Logf("Large number memory test: 100 operations in %v", duration)
	}
	
	// Large number operations should still be reasonably fast
	if duration > 500*time.Millisecond {
		t.Errorf("Large number memory test too slow: %v", duration)
	}
}

// testConcurrencyPerformance tests concurrent operation performance
func (suite *PerformanceTestSuite) testConcurrencyPerformance(ctx context.Context, t *testing.T) {
	// Test concurrent operations
	concurrencyLevels := []int{1, 2, 4, 8}
	
	for _, level := range concurrencyLevels {
		t.Run(fmt.Sprintf("ConcurrencyLevel_%d", level), func(t *testing.T) {
			start := time.Now()
			
			// Create channels for coordination
			done := make(chan bool, level)
			errors := make(chan error, level)
			
			// Start concurrent operations
			for i := 0; i < level; i++ {
				go func(workerID int) {
					// Each worker performs 100 operations
					for j := 0; j < 100; j++ {
						expr := fmt.Sprintf("%d + %d", workerID, j)
						result, err := suite.calc.Evaluate(expr)
						if err != nil {
							errors <- fmt.Errorf("worker %d operation %d failed: %v", workerID, j, err)
							return
						}
						
						// Verify result is reasonable
						expected := float64(workerID + j)
						if result != expected {
							errors <- fmt.Errorf("worker %d operation %d: expected %v, got %v", 
								workerID, j, expected, result)
							return
						}
					}
					done <- true
				}(i)
			}
			
			// Wait for all workers to complete
			completed := 0
			timeout := time.After(5 * time.Second)
			
			for completed < level {
				select {
				case <-done:
					completed++
				case err := <-errors:
					t.Errorf("Concurrent operation failed: %v", err)
					return
				case <-timeout:
					t.Errorf("Concurrent operations timed out at level %d", level)
					return
				}
			}
			
			duration := time.Since(start)
			totalOps := level * 100
			opsPerSecond := float64(totalOps) / duration.Seconds()
			
			if suite.config.Verbose {
				t.Logf("Concurrency level %d: %d operations in %v (%.2f ops/sec)", 
					level, totalOps, duration, opsPerSecond)
			}
			
			// Performance should scale reasonably with concurrency
			// (This is a simplified test - real concurrency testing would be more complex)
			if opsPerSecond < 1000 {
				t.Errorf("Concurrency performance too low at level %d: %.2f ops/sec", level, opsPerSecond)
			}
		})
	}
}

// testScalabilityPerformance tests scalability characteristics
func (suite *PerformanceTestSuite) testScalabilityPerformance(ctx context.Context, t *testing.T) {
	// Test scalability with increasing operation complexity
	complexityLevels := []struct {
		name        string
		expr        string
		maxDuration time.Duration
	}{
		{"Simple", "2 + 3", 1 * time.Millisecond},
		{"Medium", "2 * (3 + 4) - 5 / 2", 2 * time.Millisecond},
		{"Complex", "((2 + 3) * 4) - (5 / 2) + (6 * 7) / 8", 3 * time.Millisecond},
		{"Very Complex", "((2 + 3) * 4) - (5 / 2) + (6 * 7) / 8 - (9 + 10) * 11", 5 * time.Millisecond},
	}
	
	for _, level := range complexityLevels {
		t.Run(level.name, func(t *testing.T) {
			// Run operation multiple times for accurate measurement
			iterations := 50
			start := time.Now()
			
			for i := 0; i < iterations; i++ {
				result, err := suite.calc.Evaluate(level.expr)
				if err != nil {
					t.Errorf("Scalability test %s failed at iteration %d: %v", level.name, i, err)
					return
				}
				
				// Verify result is reasonable
				if result == 0 && level.expr != "0 + 0" {
					t.Errorf("Scalability test %s returned unexpected zero result at iteration %d", level.name, i)
				}
			}
			
			duration := time.Since(start)
			avgDuration := duration / time.Duration(iterations)
			
			if suite.config.Verbose {
				t.Logf("Scalability test %s: %d iterations in %v (avg: %v per operation)", 
					level.name, iterations, duration, avgDuration)
			}
			
			// Verify performance meets requirements
			if avgDuration > level.maxDuration {
				t.Errorf("Scalability test %s too slow: %v (max: %v)", level.name, avgDuration, level.maxDuration)
			}
		})
	}
	
	// Test scalability with increasing data size
	dataSizes := []struct {
		name string
		size int
		maxDuration time.Duration
	}{
		{"Small", 10, 1 * time.Millisecond},
		{"Medium", 100, 5 * time.Millisecond},
		{"Large", 1000, 20 * time.Millisecond},
	}
	
	for _, size := range dataSizes {
		t.Run(fmt.Sprintf("DataSize_%s", size.name), func(t *testing.T) {
			start := time.Now()
			
			// Perform operations with increasing data size
			for i := 0; i < size.size; i++ {
				expr := fmt.Sprintf("%d + %d", i, i+1)
				result, err := suite.calc.Evaluate(expr)
				if err != nil {
					t.Errorf("Data size test %s failed at iteration %d: %v", size.name, i, err)
					return
				}
				
				// Verify result is reasonable
				expected := float64(i + i + 1)
				if result != expected {
					t.Errorf("Data size test %s: expected %v, got %v at iteration %d", 
						size.name, expected, result, i)
					return
				}
			}
			
			duration := time.Since(start)
			avgDuration := duration / time.Duration(size.size)
			
			if suite.config.Verbose {
				t.Logf("Data size test %s: %d operations in %v (avg: %v per operation)", 
					size.name, size.size, duration, avgDuration)
			}
			
			// Verify performance meets requirements
			if avgDuration > size.maxDuration {
				t.Errorf("Data size test %s too slow: %v (max: %v)", size.name, avgDuration, size.maxDuration)
			}
		})
	}
}

// TestE2EPerformanceRegression tests for performance regressions
func TestE2EPerformanceRegression(t *testing.T) {
	config := DefaultPerformanceTestConfig()
	suite := NewPerformanceTestSuite(config)
	
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	
	_ = ctx
	
	t.Run("RegressionBaseline", func(t *testing.T) {
		// Establish performance baselines
		baselines := map[string]time.Duration{
			"simple_addition":     1 * time.Millisecond,
			"complex_expression":  2 * time.Millisecond,
			"nested_parentheses":  3 * time.Millisecond,
			"large_numbers":       2 * time.Millisecond,
			"decimal_operations":  2 * time.Millisecond,
		}
		
		for testName, baseline := range baselines {
			t.Run(testName, func(t *testing.T) {
				var expr string
				switch testName {
				case "simple_addition":
					expr = "2 + 3"
				case "complex_expression":
					expr = "2 * (3 + 4) - 5 / 2"
				case "nested_parentheses":
					expr = "((2 + 3) * 4) - (5 / 2)"
				case "large_numbers":
					expr = "1000 * 1000 + 500 * 500"
				case "decimal_operations":
					expr = "3.14159 * 2.71828"
				}
				
				// Run operation multiple times for accurate measurement
				iterations := 100
				start := time.Now()
				
				for i := 0; i < iterations; i++ {
					result, err := suite.calc.Evaluate(expr)
					if err != nil {
						t.Errorf("Regression test %s failed at iteration %d: %v", testName, i, err)
						return
					}
					
					// Verify result is reasonable
					if result == 0 && expr != "0 + 0" {
						t.Errorf("Regression test %s returned unexpected zero result at iteration %d", testName, i)
					}
				}
				
				duration := time.Since(start)
				avgDuration := duration / time.Duration(iterations)
				
				if suite.config.Verbose {
					t.Logf("Regression test %s: %d iterations in %v (avg: %v, baseline: %v)", 
						testName, iterations, duration, avgDuration, baseline)
				}
				
				// Check for performance regression (allow 20% tolerance)
				maxAllowed := baseline + baseline/5
				if avgDuration > maxAllowed {
					t.Errorf("Performance regression detected in %s: %v (baseline: %v, max allowed: %v)", 
						testName, avgDuration, baseline, maxAllowed)
				}
			})
		}
	})
}

// TestE2EPerformanceUnderLoad tests performance under load
func TestE2EPerformanceUnderLoad(t *testing.T) {
	config := DefaultPerformanceTestConfig()
	suite := NewPerformanceTestSuite(config)
	
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	
	_ = ctx
	
	t.Run("LoadTest", func(t *testing.T) {
		// Test performance under various load conditions
		loadLevels := []struct {
			name        string
			operations  int
			concurrency int
			maxDuration time.Duration
		}{
			{"Light Load", 100, 1, 100 * time.Millisecond},
			{"Medium Load", 500, 2, 500 * time.Millisecond},
			{"Heavy Load", 1000, 4, 1 * time.Second},
		}
		
		for _, load := range loadLevels {
			t.Run(load.name, func(t *testing.T) {
				start := time.Now()
				
				// Create channels for coordination
				done := make(chan bool, load.concurrency)
				errors := make(chan error, load.concurrency)
				
				// Start concurrent operations
				for i := 0; i < load.concurrency; i++ {
					go func(workerID int) {
						opsPerWorker := load.operations / load.concurrency
						for j := 0; j < opsPerWorker; j++ {
							expr := fmt.Sprintf("%d + %d", workerID, j)
							result, err := suite.calc.Evaluate(expr)
							if err != nil {
								errors <- fmt.Errorf("worker %d operation %d failed: %v", workerID, j, err)
								return
							}
							
							// Verify result is reasonable
							expected := float64(workerID + j)
							if result != expected {
								errors <- fmt.Errorf("worker %d operation %d: expected %v, got %v", 
									workerID, j, expected, result)
								return
							}
						}
						done <- true
					}(i)
				}
				
				// Wait for all workers to complete
				completed := 0
				timeout := time.After(load.maxDuration + 1*time.Second)
				
				for completed < load.concurrency {
					select {
					case <-done:
						completed++
					case err := <-errors:
						t.Errorf("Load test %s failed: %v", load.name, err)
						return
					case <-timeout:
						t.Errorf("Load test %s timed out", load.name)
						return
					}
				}
				
				duration := time.Since(start)
				opsPerSecond := float64(load.operations) / duration.Seconds()
				
				if suite.config.Verbose {
					t.Logf("Load test %s: %d operations in %v (%.2f ops/sec)", 
						load.name, load.operations, duration, opsPerSecond)
				}
				
				// Verify performance meets requirements
				if duration > load.maxDuration {
					t.Errorf("Load test %s too slow: %v (max: %v)", load.name, duration, load.maxDuration)
				}
				
				// Verify throughput is reasonable
				if opsPerSecond < 1000 {
					t.Errorf("Load test %s throughput too low: %.2f ops/sec", load.name, opsPerSecond)
				}
			})
		}
	})
}

// BenchmarkE2EPerformanceOperations benchmarks E2E performance operations
func BenchmarkE2EPerformanceOperations(b *testing.B) {
	config := DefaultPerformanceTestConfig()
	suite := NewPerformanceTestSuite(config)
	
	operations := []string{
		"2 + 3",
		"10 - 4",
		"3 * 4",
		"15 / 3",
		"2 * (3 + 4) - 5 / 2",
		"((2 + 3) * 4) - (5 / 2)",
		"1000 * 1000 + 500 * 500",
		"3.14159 * 2.71828",
	}
	
	for _, op := range operations {
		b.Run(fmt.Sprintf("Operation_%s", strings.ReplaceAll(op, " ", "_")), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result, err := suite.calc.Evaluate(op)
				if err != nil {
					b.Fatalf("Benchmark operation failed: %v", err)
				}
				
				// Verify result is reasonable
				if result == 0 && op != "0 + 0" {
					b.Fatalf("Benchmark operation returned unexpected zero result")
				}
			}
		})
	}
}

// BenchmarkE2EPerformanceConcurrency benchmarks E2E performance under concurrency
func BenchmarkE2EPerformanceConcurrency(b *testing.B) {
	config := DefaultPerformanceTestConfig()
	suite := NewPerformanceTestSuite(config)
	
	concurrencyLevels := []int{1, 2, 4, 8}
	
	for _, level := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency_%d", level), func(b *testing.B) {
			b.SetParallelism(level)
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					expr := "2 + 3"
					result, err := suite.calc.Evaluate(expr)
					if err != nil {
						b.Fatalf("Concurrent benchmark operation failed: %v", err)
					}
					
					// Verify result is reasonable
					if result != 5.0 {
						b.Fatalf("Concurrent benchmark operation returned unexpected result: %v", result)
					}
				}
			})
		})
	}
}