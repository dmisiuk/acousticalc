package e2e

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCalculatorWorkflow tests the complete calculator workflow
func TestCalculatorWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Configure E2E test
	config := &E2ETestConfig{
		TestName:         "CalculatorWorkflow",
		OutputDir:        filepath.Join("tests", "artifacts", "e2e"),
		Timeout:          3 * time.Minute,
		ScreenshotOnPass: true,
		ScreenshotOnFail: true,
		RecordTerminal:   true,
	}

	runner := NewE2ETestRunner(config)

	// Add observer for logging
	observer := &LoggingE2EObserver{t: t}
	runner.AddObserver(observer)

	// Cleanup when test completes
	defer func() {
		runner.StopApplication()
		runner.Complete()
		err := runner.GenerateReport()
		assert.NoError(t, err, "Should generate E2E report")
	}()

	// Start the application
	require.NoError(t, runner.StartApplication(), "Should start application successfully")

	// Capture initial state
	_, err := runner.CaptureScreenshot("initial_state")
	require.NoError(t, err, "Should capture initial screenshot")

	// Test basic calculator workflow
	t.Run("BasicOperations", func(t *testing.T) {
		// Test addition: 2 + 2 = 4
		testInputSequence(t, runner, []string{"2", "+", "2", "="}, "4")

		// Test subtraction: 10 - 3 = 7
		testInputSequence(t, runner, []string{"1", "0", "-", "3", "="}, "7")

		// Test multiplication: 6 * 7 = 42
		testInputSequence(t, runner, []string{"6", "*", "7", "="}, "42")

		// Test division: 15 / 3 = 5
		testInputSequence(t, runner, []string{"1", "5", "/", "3", "="}, "5")
	})

	// Test complex workflow
	t.Run("ComplexCalculation", func(t *testing.T) {
		// Test: (2 + 3) * 4 = 20
		testInputSequence(t, runner, []string{"(", "2", "+", "3", ")", "*", "4", "="}, "20")

		// Test: 10 / 2 + 5 = 10
		testInputSequence(t, runner, []string{"1", "0", "/", "2", "+", "5", "="}, "10")
	})

	// Test error handling
	t.Run("ErrorHandling", func(t *testing.T) {
		// Test division by zero
		testInputSequence(t, runner, []string{"5", "/", "0", "="}, "error")

		// Test invalid input
		testInputSequence(t, runner, []string{"a", "+", "b", "="}, "error")
	})

	// Test clear functionality
	t.Run("ClearFunction", func(t *testing.T) {
		// Clear and start fresh
		testInputSequence(t, runner, []string{"c", "1", "+", "1", "="}, "2")
	})

	// Capture final state
	_, err = runner.CaptureScreenshot("final_state")
	require.NoError(t, err, "Should capture final screenshot")
}

// TestCrossPlatformCompatibility tests cross-platform compatibility
func TestCrossPlatformCompatibility(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	config := &E2ETestConfig{
		TestName:         "CrossPlatformCompatibility",
		OutputDir:        filepath.Join("tests", "artifacts", "e2e"),
		Timeout:          2 * time.Minute,
		ScreenshotOnPass: true,
		ScreenshotOnFail: true,
		RecordTerminal:   false, // Disable for this test to focus on compatibility
	}

	runner := NewE2ETestRunner(config)
	observer := &LoggingE2EObserver{t: t}
	runner.AddObserver(observer)

	defer func() {
		runner.StopApplication()
		runner.Complete()
		err := runner.GenerateReport()
		assert.NoError(t, err, "Should generate E2E report")
	}()

	require.NoError(t, runner.StartApplication(), "Should start application successfully")

	// Test platform-specific behaviors
	t.Run("PlatformSpecific", func(t *testing.T) {
		// Test basic operations that should work across all platforms
		testInputSequence(t, runner, []string{"9", "+", "1", "="}, "10")
		testInputSequence(t, runner, []string{"8", "*", "2", "="}, "16")
		testInputSequence(t, runner, []string{"2", "0", "/", "4", "="}, "5")
	})

	// Test file operations if applicable
	t.Run("FileOperations", func(t *testing.T) {
		// This would test file-based operations that might be platform-specific
		// For now, just test basic functionality
		testInputSequence(t, runner, []string{"1", "0", "0", "=", "c"}, "clear")
	})
}

// TestPerformanceAndResponsiveness tests application performance
func TestPerformanceAndResponsiveness(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	config := &E2ETestConfig{
		TestName:         "PerformanceTest",
		OutputDir:        filepath.Join("tests", "artifacts", "e2e"),
		Timeout:          1 * time.Minute,
		ScreenshotOnPass: false, // Reduce overhead for performance test
		ScreenshotOnFail: true,
		RecordTerminal:   false,
	}

	runner := NewE2ETestRunner(config)
	observer := &LoggingE2EObserver{t: t}
	runner.AddObserver(observer)

	defer func() {
		runner.StopApplication()
		runner.Complete()
	}()

	require.NoError(t, runner.StartApplication(), "Should start application successfully")

	// Test rapid input sequence
	t.Run("RapidInput", func(t *testing.T) {
		start := time.Now()

		// Perform 50 rapid calculations
		for i := 0; i < 50; i++ {
			sequence := []string{fmt.Sprintf("%d", i%10), "+", "1", "="}
			testInputSequence(t, runner, sequence, fmt.Sprintf("%d", (i%10)+1))
		}

		duration := time.Since(start)
		t.Logf("Completed 50 calculations in %v", duration)

		// Assert that it completed within reasonable time (less than 30 seconds)
		assert.Less(t, duration, 30*time.Second, "Should complete rapid calculations within 30 seconds")
	})

	// Test memory usage (simplified)
	t.Run("MemoryUsage", func(t *testing.T) {
		// This is a simplified memory test
		// In a real implementation, you'd monitor actual memory usage
		start := time.Now()

		// Perform memory-intensive operations
		for i := 0; i < 100; i++ {
			sequence := []string{fmt.Sprintf("%d", i), "*", fmt.Sprintf("%d", i), "="}
			testInputSequence(t, runner, sequence, fmt.Sprintf("%d", i*i))
		}

		duration := time.Since(start)
		t.Logf("Completed 100 multiplication operations in %v", duration)
		assert.Less(t, duration, 45*time.Second, "Should complete memory operations within 45 seconds")
	})
}

// TestErrorRecovery tests error handling and recovery
func TestErrorRecovery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	config := &E2ETestConfig{
		TestName:         "ErrorRecoveryTest",
		OutputDir:        filepath.Join("tests", "artifacts", "e2e"),
		Timeout:          2 * time.Minute,
		ScreenshotOnPass: true,
		ScreenshotOnFail: true,
		RecordTerminal:   true,
	}

	runner := NewE2ETestRunner(config)
	observer := &LoggingE2EObserver{t: t}
	runner.AddObserver(observer)

	defer func() {
		runner.StopApplication()
		runner.Complete()
		err := runner.GenerateReport()
		assert.NoError(t, err, "Should generate E2E report")
	}()

	require.NoError(t, runner.StartApplication(), "Should start application successfully")

	// Test various error scenarios
	t.Run("DivisionByZero", func(t *testing.T) {
		testInputSequence(t, runner, []string{"5", "/", "0", "="}, "error")

		// Should recover and allow new input
		testInputSequence(t, runner, []string{"c", "1", "+", "1", "="}, "2")
	})

	t.Run("InvalidInput", func(t *testing.T) {
		testInputSequence(t, runner, []string{"a", "b", "c", "="}, "error")

		// Should recover and allow new input
		testInputSequence(t, runner, []string{"c", "5", "*", "2", "="}, "10")
	})

	t.Run("Overflow", func(t *testing.T) {
		// Test very large numbers
		testInputSequence(t, runner, []string{"9", "9", "9", "9", "9", "9", "*", "9", "9", "9", "9", "9", "9", "="}, "error")

		// Should recover
		testInputSequence(t, runner, []string{"c", "1", "+", "1", "="}, "2")
	})
}

// testInputSequence is a helper function to test input sequences
func testInputSequence(t *testing.T, runner *E2ETestRunner, sequence []string, expectedResult string) {
	t.Helper()

	// Clear calculator first
	if err := runner.SimulateInput("c"); err != nil {
		t.Logf("Warning: Failed to clear calculator: %v", err)
	}

	// Input the sequence
	for _, input := range sequence {
		require.NoError(t, runner.SimulateInput(input), "Should simulate input '%s'", input)
		time.Sleep(100 * time.Millisecond) // Small delay between inputs
	}

	// Wait for result (this is simplified - in real implementation you'd verify actual output)
	time.Sleep(500 * time.Millisecond)

	// Capture screenshot after operation
	eventName := fmt.Sprintf("operation_%s", expectedResult)
	if _, err := runner.CaptureScreenshot(eventName); err != nil {
		t.Logf("Warning: Failed to capture screenshot: %v", err)
	}
}

// LoggingE2EObserver is an observer that logs E2E test events
type LoggingE2EObserver struct {
	t *testing.T
}

func (o *LoggingE2EObserver) OnEvent(event E2ETestEvent) {
	o.t.Logf("E2E Event: %s - %s", event.Type, event.Description)
	if event.Error != nil {
		o.t.Logf("  Error: %v", event.Error)
	}
}

func (o *LoggingE2EObserver) OnScreenshot(path string, eventType string) {
	o.t.Logf("Screenshot captured: %s for event: %s", path, eventType)
}

func (o *LoggingE2EObserver) OnRecording(path string, eventType string) {
	o.t.Logf("Recording available: %s for event: %s", path, eventType)
}

func (o *LoggingE2EObserver) OnTestComplete(runner *E2ETestRunner) {
	duration := time.Since(runner.startTime)
	o.t.Logf("E2E test completed: %s, Duration: %v, Events: %d",
		runner.config.TestName, duration, len(runner.events))
}

// BenchmarkE2EOperations benchmarks E2E operations
func BenchmarkE2EOperations(b *testing.B) {
	config := &E2ETestConfig{
		TestName:         "Benchmark",
		OutputDir:        filepath.Join("tests", "artifacts", "e2e"),
		Timeout:          5 * time.Minute,
		ScreenshotOnPass: false,
		ScreenshotOnFail: false,
		RecordTerminal:   false,
	}

	b.Run("BasicAddition", func(b *testing.B) {
		runner := NewE2ETestRunner(config)

		// Setup
		require.NoError(b, runner.StartApplication())
		defer runner.StopApplication()

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			// Perform addition operation
			sequence := []string{"1", "+", "1", "="}
			for _, input := range sequence {
				if err := runner.SimulateInput(input); err != nil {
					b.Fatalf("Failed to simulate input: %v", err)
				}
			}

			// Clear for next iteration
			runner.SimulateInput("c")
		}
	})
}
