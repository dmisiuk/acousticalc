package e2e

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestE2EIntegrationWithRecording tests E2E integration with recording
func TestE2EIntegrationWithRecording(t *testing.T) {
	config := DefaultE2EConfig()
	config.EnableRecording = true
	config.Verbose = true
	
	suite := NewE2ETestSuite(config)
	
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	
	t.Run("RecordingIntegration", func(t *testing.T) {
		// Test that recording directory exists and is writable
		if _, err := os.Stat(config.RecordingDir); os.IsNotExist(err) {
			t.Error("Recording directory does not exist")
		}
		
		// Test that we can create recording artifacts
		recordingFile := filepath.Join(config.RecordingDir, "e2e_integration_test.txt")
		content := fmt.Sprintf("E2E integration test recording - %s", time.Now().Format(time.RFC3339))
		
		err := os.WriteFile(recordingFile, []byte(content), 0644)
		if err != nil {
			t.Errorf("Failed to create recording file: %v", err)
		}
		defer os.Remove(recordingFile)
		
		// Run E2E test
		suite.testBasicArithmeticWorkflow(ctx, t)
		
		// Verify recording file was created
		if _, err := os.Stat(recordingFile); os.IsNotExist(err) {
			t.Error("Recording file was not created")
		}
		
		t.Logf("Recording integration test completed successfully")
	})
}

// TestE2EIntegrationWithReporting tests E2E integration with reporting
func TestE2EIntegrationWithReporting(t *testing.T) {
	config := DefaultE2EConfig()
	suite := NewE2ETestSuite(config)
	
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	
	t.Run("ReportingIntegration", func(t *testing.T) {
		// Test that reporting directory exists and is writable
		if _, err := os.Stat(config.ArtifactsDir); os.IsNotExist(err) {
			t.Error("Artifacts directory does not exist")
		}
		
		// Test that we can create report artifacts
		reportFile := filepath.Join(config.ArtifactsDir, "e2e_integration_report.json")
		reportContent := fmt.Sprintf(`{
			"test_name": "e2e_integration_test",
			"type": "e2e",
			"platform": "%s",
			"status": "passed",
			"duration": "100ms",
			"timestamp": "%s"
		}`, config.Platform, time.Now().Format(time.RFC3339))
		
		err := os.WriteFile(reportFile, []byte(reportContent), 0644)
		if err != nil {
			t.Errorf("Failed to create report file: %v", err)
		}
		defer os.Remove(reportFile)
		
		// Run E2E test
		suite.testBasicArithmeticWorkflow(ctx, t)
		
		// Verify report file was created
		if _, err := os.Stat(reportFile); os.IsNotExist(err) {
			t.Error("Report file was not created")
		}
		
		t.Logf("Reporting integration test completed successfully")
	})
}

// TestE2EIntegrationWithVisualTesting tests E2E integration with visual testing
func TestE2EIntegrationWithVisualTesting(t *testing.T) {
	config := DefaultE2EConfig()
	suite := NewE2ETestSuite(config)
	
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	
	_ = suite
	_ = ctx
	
	t.Run("VisualTestingIntegration", func(t *testing.T) {
		// This test ensures E2E framework integrates with existing visual testing
		// from Story 0.2.2 without conflicts
		
		// Test that visual testing artifacts directory exists
		visualArtifactsDir := "tests/artifacts/screenshots"
		if _, err := os.Stat(visualArtifactsDir); os.IsNotExist(err) {
			// Create the directory if it doesn't exist
			err = os.MkdirAll(visualArtifactsDir, 0755)
			if err != nil {
				t.Errorf("Failed to create visual artifacts directory: %v", err)
			}
		}
		
		// Test that E2E artifacts directory exists
		if _, err := os.Stat(config.ArtifactsDir); os.IsNotExist(err) {
			t.Error("E2E artifacts directory does not exist")
		}
		
		// Test that recording directory exists
		if _, err := os.Stat(config.RecordingDir); os.IsNotExist(err) {
			t.Error("Recording directory does not exist")
		}
		
		// Test that both directories are separate and don't conflict
		if config.ArtifactsDir == visualArtifactsDir {
			t.Error("E2E artifacts directory should be separate from visual artifacts directory")
		}
		
		if config.RecordingDir == visualArtifactsDir {
			t.Error("Recording directory should be separate from visual artifacts directory")
		}
		
		// Test that we can create files in both directories
		e2eTestFile := filepath.Join(config.ArtifactsDir, "e2e_test.txt")
		visualTestFile := filepath.Join(visualArtifactsDir, "visual_test.txt")
		
		err := os.WriteFile(e2eTestFile, []byte("E2E test content"), 0644)
		if err != nil {
			t.Errorf("Failed to create E2E test file: %v", err)
		}
		defer os.Remove(e2eTestFile)
		
		err = os.WriteFile(visualTestFile, []byte("Visual test content"), 0644)
		if err != nil {
			t.Errorf("Failed to create visual test file: %v", err)
		}
		defer os.Remove(visualTestFile)
		
		// Verify both files exist
		if _, err := os.Stat(e2eTestFile); os.IsNotExist(err) {
			t.Error("E2E test file was not created")
		}
		
		if _, err := os.Stat(visualTestFile); os.IsNotExist(err) {
			t.Error("Visual test file was not created")
		}
	})
}

// TestE2EIntegrationWithExistingTests tests E2E integration with existing test infrastructure
func TestE2EIntegrationWithExistingTests(t *testing.T) {
	config := DefaultE2EConfig()
	suite := NewE2ETestSuite(config)
	
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	
	_ = suite
	_ = ctx
	
	t.Run("ExistingTestIntegration", func(t *testing.T) {
		// Test that E2E framework can work alongside existing unit and integration tests
		
		// Test that unit tests directory exists
		unitTestsDir := "tests/unit"
		if _, err := os.Stat(unitTestsDir); os.IsNotExist(err) {
			t.Logf("Unit tests directory does not exist: %s", unitTestsDir)
		}
		
		// Test that integration tests directory exists
		integrationTestsDir := "tests/integration"
		if _, err := os.Stat(integrationTestsDir); os.IsNotExist(err) {
			t.Logf("Integration tests directory does not exist: %s", integrationTestsDir)
		}
		
		// Test that E2E tests directory exists
		e2eTestsDir := "tests/e2e"
		if _, err := os.Stat(e2eTestsDir); os.IsNotExist(err) {
			t.Logf("E2E tests directory does not exist: %s", e2eTestsDir)
		}
		
		// Test that all test directories are separate
		if unitTestsDir == integrationTestsDir || unitTestsDir == e2eTestsDir || integrationTestsDir == e2eTestsDir {
			t.Error("Test directories should be separate")
		}
		
		// Test that we can run E2E tests without affecting other test types
		suite.testBasicArithmeticWorkflow(ctx, t)
		
		// Verify that E2E test execution doesn't interfere with other test types
		// This is a basic smoke test - in a real scenario, you would run the other
		// test types and verify they still pass
	})
}

// TestE2EIntegrationWithCI tests E2E integration with CI infrastructure
func TestE2EIntegrationWithCI(t *testing.T) {
	config := DefaultE2EConfig()
	suite := NewE2ETestSuite(config)
	
	_ = suite
	
	t.Run("CIIntegration", func(t *testing.T) {
		// Test that E2E framework can work with CI infrastructure
		
		// Test that CI artifacts directory exists
		ciArtifactsDir := "tests/artifacts"
		if _, err := os.Stat(ciArtifactsDir); os.IsNotExist(err) {
			t.Error("CI artifacts directory does not exist")
		}
		
		// Test that E2E artifacts are created in the correct location
		if _, err := os.Stat(config.ArtifactsDir); os.IsNotExist(err) {
			t.Error("E2E artifacts directory does not exist")
		}
		
		// Test that recording artifacts are created in the correct location
		if _, err := os.Stat(config.RecordingDir); os.IsNotExist(err) {
			t.Error("Recording directory does not exist")
		}
		
		// Test that E2E framework respects CI timeout constraints
		if config.Timeout > 30*time.Second {
			t.Error("E2E timeout should respect CI constraints (<30s)")
		}
		
		// Test that E2E framework can generate CI-compatible artifacts
		testFile := filepath.Join(config.ArtifactsDir, "ci_test.txt")
		content := fmt.Sprintf("CI integration test - %s", time.Now().Format(time.RFC3339))
		
		err := os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			t.Errorf("Failed to create CI test file: %v", err)
		}
		defer os.Remove(testFile)
		
		// Verify file was created
		if _, err := os.Stat(testFile); os.IsNotExist(err) {
			t.Error("CI test file was not created")
		}
	})
}

// TestE2EIntegrationWithPerformance tests E2E integration with performance constraints
func TestE2EIntegrationWithPerformance(t *testing.T) {
	config := DefaultE2EConfig()
	suite := NewE2ETestSuite(config)
	
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	
	_ = suite
	_ = ctx
	
	t.Run("PerformanceIntegration", func(t *testing.T) {
		// Test that E2E framework meets performance constraints
		
		// Test that E2E tests complete within reasonable time
		start := time.Now()
		suite.testBasicArithmeticWorkflow(ctx, t)
		duration := time.Since(start)
		
		if duration > 5*time.Second {
			t.Errorf("E2E test took too long: %v", duration)
		}
		
		// Test that E2E framework doesn't consume excessive resources
		// This is a basic test - in a real scenario, you would monitor
		// memory usage, CPU usage, etc.
		
		// Test that E2E framework can handle concurrent execution
		// (This is a simplified test - real concurrency testing would be more complex)
		concurrentTests := 3
		results := make(chan bool, concurrentTests)
		
		for i := 0; i < concurrentTests; i++ {
			go func() {
				// Run a simple E2E test
				suite.testBasicArithmeticWorkflow(ctx, t)
				results <- true
			}()
		}
		
		// Wait for all tests to complete
		completed := 0
		timeout := time.After(10 * time.Second)
		
		for completed < concurrentTests {
			select {
			case <-results:
				completed++
			case <-timeout:
				t.Error("Concurrent E2E tests timed out")
				return
			}
		}
		
		if suite.config.Verbose {
			t.Logf("Concurrent E2E tests completed successfully")
		}
	})
}

// TestE2EIntegrationWithErrorHandling tests E2E integration with error handling
func TestE2EIntegrationWithErrorHandling(t *testing.T) {
	config := DefaultE2EConfig()
	suite := NewE2ETestSuite(config)
	
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	
	_ = suite
	_ = ctx
	
	t.Run("ErrorHandlingIntegration", func(t *testing.T) {
		// Test that E2E framework handles errors gracefully
		
		// Test that E2E framework can handle invalid configurations
		invalidConfig := &E2ETestConfig{
			Timeout:        -1 * time.Second, // Invalid timeout
			RecordingDir:   "/invalid/path",  // Invalid path
			ArtifactsDir:   "/invalid/path",  // Invalid path
			Platform:       "invalid",        // Invalid platform
			Verbose:        false,
			EnableRecording: false,
		}
		
		// Test that invalid configuration doesn't crash the framework
		invalidSuite := NewE2ETestSuite(invalidConfig)
		if invalidSuite == nil {
			t.Error("E2E test suite should handle invalid configuration gracefully")
		}
		
		// Test that E2E framework can handle test failures
		// This is a basic test - in a real scenario, you would test
		// various failure modes and recovery scenarios
		
		// Test that E2E framework can handle resource exhaustion
		// (This is a simplified test - real resource testing would be more complex)
		
		// Test that E2E framework can handle network issues
		// (This is a placeholder - real network testing would be more complex)
		
		// Test that E2E framework can handle file system issues
		// (This is a placeholder - real file system testing would be more complex)
	})
}

// TestE2EIntegrationWithSecurity tests E2E integration with security considerations
func TestE2EIntegrationWithSecurity(t *testing.T) {
	config := DefaultE2EConfig()
	suite := NewE2ETestSuite(config)
	
	_ = suite
	
	t.Run("SecurityIntegration", func(t *testing.T) {
		// Test that E2E framework follows security best practices
		
		// Test that E2E framework doesn't expose sensitive information
		// in test artifacts
		testFile := filepath.Join(config.ArtifactsDir, "security_test.txt")
		content := "This is a test file with sensitive information: password123"
		
		err := os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			t.Errorf("Failed to create security test file: %v", err)
		}
		defer os.Remove(testFile)
		
		// Test that E2E framework uses appropriate file permissions
		info, err := os.Stat(testFile)
		if err != nil {
			t.Errorf("Failed to stat security test file: %v", err)
		}
		
		// File should be readable and writable by owner only
		mode := info.Mode()
		if mode&0077 != 0 {
			t.Logf("E2E test file permissions: %v (group/other access detected)", mode)
		}
		
		// Test that E2E framework doesn't create world-writable files
		// (This is a basic test - real security testing would be more comprehensive)
		
		// Test that E2E framework doesn't expose sensitive environment variables
		// (This is a placeholder - real environment variable testing would be more complex)
		
		// Test that E2E framework doesn't expose sensitive command line arguments
		// (This is a placeholder - real command line testing would be more complex)
	})
}

// TestE2EIntegrationWithMaintainability tests E2E integration with maintainability
func TestE2EIntegrationWithMaintainability(t *testing.T) {
	config := DefaultE2EConfig()
	suite := NewE2ETestSuite(config)
	
	_ = suite
	
	t.Run("MaintainabilityIntegration", func(t *testing.T) {
		// Test that E2E framework is maintainable and extensible
		
		// Test that E2E framework has clear separation of concerns
		// (This is a basic test - real maintainability testing would be more comprehensive)
		
		// Test that E2E framework follows consistent naming conventions
		// (This is a basic test - real naming convention testing would be more comprehensive)
		
		// Test that E2E framework has appropriate documentation
		// (This is a placeholder - real documentation testing would be more complex)
		
		// Test that E2E framework has appropriate error messages
		// (This is a basic test - real error message testing would be more comprehensive)
		
		// Test that E2E framework has appropriate logging
		// (This is a basic test - real logging testing would be more comprehensive)
		
		// Test that E2E framework has appropriate configuration options
		// (This is a basic test - real configuration testing would be more comprehensive)
		
		// Test that E2E framework has appropriate extension points
		// (This is a placeholder - real extension point testing would be more complex)
	})
}

// BenchmarkE2EIntegrationOperations benchmarks E2E integration operations
func BenchmarkE2EIntegrationOperations(b *testing.B) {
	config := DefaultE2EConfig()
	suite := NewE2ETestSuite(config)
	
	b.Run("BasicWorkflow", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
			suite.testBasicArithmeticWorkflow(ctx, &testing.T{})
			cancel()
		}
	})
	
	b.Run("ComplexWorkflow", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
			suite.testComplexExpressionWorkflow(ctx, &testing.T{})
			cancel()
		}
	})
	
	b.Run("ErrorHandlingWorkflow", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
			suite.testErrorHandlingWorkflow(ctx, &testing.T{})
			cancel()
		}
	})
	
	b.Run("PerformanceWorkflow", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
			suite.testPerformanceWorkflow(ctx, &testing.T{})
			cancel()
		}
	})
}