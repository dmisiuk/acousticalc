package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestTestIsolation ensures that tests are properly isolated and don't interfere with each other
func TestTestIsolation(t *testing.T) {
	suite := NewTestSuite("Test Isolation", "Ensures test isolation and proper cleanup")

	suite.Setup = func() error {
		fmt.Println("Setting up isolation tests...")
		return setupTestEnvironment()
	}

	suite.Teardown = func() error {
		fmt.Println("Cleaning up isolation tests...")
		return cleanupTestEnvironment()
	}

	suite.AddTest(TestCase{
		Name:        "File System Isolation",
		Description: "Tests that file operations don't interfere between tests",
		TestFunc: func(t *testing.T) error {
			return testFileSystemIsolation(t)
		},
	})

	suite.AddTest(TestCase{
		Name:        "Mock State Isolation",
		Description: "Tests that mock objects maintain proper state isolation",
		TestFunc: func(t *testing.T) error {
			return testMockStateIsolation(t)
		},
	})

	suite.AddTest(TestCase{
		Name:        "Environment Variable Isolation",
		Description: "Tests that environment variables are properly isolated",
		TestFunc: func(t *testing.T) error {
			return testEnvironmentVariableIsolation(t)
		},
	})

	suite.AddTest(TestCase{
		Name:        "Resource Cleanup",
		Description: "Tests that resources are properly cleaned up after each test",
		TestFunc: func(t *testing.T) error {
			return testResourceCleanup(t)
		},
	})

	suite.AddTest(TestCase{
		Name:        "Parallel Test Execution",
		Description: "Tests that tests can run in parallel without interference",
		TestFunc: func(t *testing.T) error {
			return testParallelExecution(t)
		},
	})

	suite.Run(t)
}

// setupTestEnvironment creates a controlled test environment
func setupTestEnvironment() error {
	// Create temporary directory for test artifacts
	tempDir := filepath.Join(os.TempDir(), "acousticalc-test-isolation")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create test directory: %v", err)
	}

	// Set up environment variables for testing
	os.Setenv("ACOUSTICALC_TEST_MODE", "isolation")
	os.Setenv("ACOUSTICALC_TEST_TEMP_DIR", tempDir)

	return nil
}

// cleanupTestEnvironment cleans up the test environment
func cleanupTestEnvironment() error {
	// Clean up temporary files
	tempDir := os.Getenv("ACOUSTICALC_TEST_TEMP_DIR")
	if tempDir != "" {
		if err := os.RemoveAll(tempDir); err != nil {
			return fmt.Errorf("failed to cleanup test directory: %v", err)
		}
	}

	// Clean up environment variables
	os.Unsetenv("ACOUSTICALC_TEST_MODE")
	os.Unsetenv("ACOUSTICALC_TEST_TEMP_DIR")

	return nil
}

// testFileSystemIsolation tests file system isolation between tests
func testFileSystemIsolation(t *testing.T) error {
	tempDir := os.Getenv("ACOUSTICALC_TEST_TEMP_DIR")
	if tempDir == "" {
		return fmt.Errorf("test environment not properly set up")
	}

	// Each test should have its own isolated directory
	testDir := filepath.Join(tempDir, fmt.Sprintf("test-%d", time.Now().UnixNano()))
	if err := os.MkdirAll(testDir, 0755); err != nil {
		return fmt.Errorf("failed to create test directory: %v", err)
	}

	// Create test files
	testFile := filepath.Join(testDir, "test.txt")
	content := fmt.Sprintf("Test data for isolation test at %d", time.Now().UnixNano())
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write test file: %v", err)
	}

	// Verify file exists and content matches
	readContent, err := os.ReadFile(testFile)
	if err != nil {
		return fmt.Errorf("failed to read test file: %v", err)
	}

	if string(readContent) != content {
		return fmt.Errorf("file content mismatch: got %s, want %s", string(readContent), content)
	}

	// Clean up
	if err := os.RemoveAll(testDir); err != nil {
		return fmt.Errorf("failed to cleanup test directory: %v", err)
	}

	t.Logf("File system isolation test completed successfully")
	return nil
}

// testMockStateIsolation tests that mock objects maintain proper state isolation
func testMockStateIsolation(t *testing.T) error {
	// Create multiple mock calculators
	mock1 := NewMockCalculator()
	mock2 := NewMockCalculator()

	// Set different states for each mock
	mock1.SetResult("2 + 3", 5.0)
	mock1.SetError("10 / 0", fmt.Errorf("division by zero"))

	mock2.SetResult("2 + 3", 99.0) // Different result
	mock2.SetError("10 / 0", fmt.Errorf("custom error")) // Different error

	// Test that mocks maintain separate states
	result1, err1 := mock1.Evaluate("2 + 3")
	result2, err2 := mock2.Evaluate("2 + 3")

	if result1 != 5.0 {
		return fmt.Errorf("mock1 returned unexpected result: got %v, want %v", result1, 5.0)
	}

	if result2 != 99.0 {
		return fmt.Errorf("mock2 returned unexpected result: got %v, want %v", result2, 99.0)
	}

	if err1 != nil {
		return fmt.Errorf("mock1 returned unexpected error: %v", err1)
	}

	if err2 != nil {
		return fmt.Errorf("mock2 returned unexpected error: %v", err2)
	}

	// Test error states
	_, err1 = mock1.Evaluate("10 / 0")
	_, err2 = mock2.Evaluate("10 / 0")

	if err1 == nil || err1.Error() != "division by zero" {
		return fmt.Errorf("mock1 returned unexpected error: got %v, want division by zero", err1)
	}

	if err2 == nil || err2.Error() != "custom error" {
		return fmt.Errorf("mock2 returned unexpected error: got %v, want custom error", err2)
	}

	t.Logf("Mock state isolation test completed successfully")
	return nil
}

// testEnvironmentVariableIsolation tests environment variable isolation
func testEnvironmentVariableIsolation(t *testing.T) error {
	// Save original environment variables
	originalMode := os.Getenv("ACOUSTICALC_TEST_MODE")
	originalDir := os.Getenv("ACOUSTICALC_TEST_TEMP_DIR")

	// Set new values
	os.Setenv("ACOUSTICALC_TEST_MODE", "isolation-test")
	os.Setenv("ACOUSTICALC_TEST_TEMP_DIR", "/tmp/test-isolation")

	// Verify new values are set
	if os.Getenv("ACOUSTICALC_TEST_MODE") != "isolation-test" {
		return fmt.Errorf("environment variable not set correctly")
	}

	if os.Getenv("ACOUSTICALC_TEST_TEMP_DIR") != "/tmp/test-isolation" {
		return fmt.Errorf("environment variable not set correctly")
	}

	// Restore original values
	if originalMode != "" {
		os.Setenv("ACOUSTICALC_TEST_MODE", originalMode)
	} else {
		os.Unsetenv("ACOUSTICALC_TEST_MODE")
	}

	if originalDir != "" {
		os.Setenv("ACOUSTICALC_TEST_TEMP_DIR", originalDir)
	} else {
		os.Unsetenv("ACOUSTICALC_TEST_TEMP_DIR")
	}

	// Verify original values are restored
	if os.Getenv("ACOUSTICALC_TEST_MODE") != originalMode {
		return fmt.Errorf("environment variable not restored correctly")
	}

	if os.Getenv("ACOUSTICALC_TEST_TEMP_DIR") != originalDir {
		return fmt.Errorf("environment variable not restored correctly")
	}

	t.Logf("Environment variable isolation test completed successfully")
	return nil
}

// testResourceCleanup tests that resources are properly cleaned up
func testResourceCleanup(t *testing.T) error {
	tempDir := os.Getenv("ACOUSTICALC_TEST_TEMP_DIR")
	if tempDir == "" {
		return fmt.Errorf("test environment not properly set up")
	}

	// Create multiple test resources
	resources := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		resourcePath := filepath.Join(tempDir, fmt.Sprintf("resource-%d", i))
		if err := os.WriteFile(resourcePath, []byte(fmt.Sprintf("resource %d", i)), 0644); err != nil {
			return fmt.Errorf("failed to create resource %d: %v", i, err)
		}
		resources = append(resources, resourcePath)
	}

	// Verify resources exist
	for _, resource := range resources {
		if _, err := os.Stat(resource); os.IsNotExist(err) {
			return fmt.Errorf("resource does not exist: %s", resource)
		}
	}

	// Cleanup resources
	for _, resource := range resources {
		if err := os.Remove(resource); err != nil {
			return fmt.Errorf("failed to cleanup resource %s: %v", resource, err)
		}
	}

	// Verify resources are cleaned up
	for _, resource := range resources {
		if _, err := os.Stat(resource); !os.IsNotExist(err) {
			return fmt.Errorf("resource still exists after cleanup: %s", resource)
		}
	}

	t.Logf("Resource cleanup test completed successfully")
	return nil
}

// testParallelExecution tests that tests can run in parallel without interference
func testParallelExecution(t *testing.T) error {
	tempDir := os.Getenv("ACOUSTICALC_TEST_TEMP_DIR")
	if tempDir == "" {
		return fmt.Errorf("test environment not properly set up")
	}

	const numWorkers = 5
	const numOperations = 20

	done := make(chan bool, numWorkers)
	errors := make(chan error, numWorkers)

	// Launch parallel workers
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			defer func() { done <- true }()

			// Each worker creates its own isolated directory
			workerDir := filepath.Join(tempDir, fmt.Sprintf("worker-%d", workerID))
			if err := os.MkdirAll(workerDir, 0755); err != nil {
				errors <- fmt.Errorf("worker %d failed to create directory: %v", workerID, err)
				return
			}

			// Perform operations in isolation
			for j := 0; j < numOperations; j++ {
				filePath := filepath.Join(workerDir, fmt.Sprintf("operation-%d.txt", j))
				content := fmt.Sprintf("Worker %d, operation %d, time %d", workerID, j, time.Now().UnixNano())

				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					errors <- fmt.Errorf("worker %d failed to write file %d: %v", workerID, j, err)
					return
				}

				// Verify file content
				readContent, err := os.ReadFile(filePath)
				if err != nil {
					errors <- fmt.Errorf("worker %d failed to read file %d: %v", workerID, j, err)
					return
				}

				if string(readContent) != content {
					errors <- fmt.Errorf("worker %d file content mismatch for operation %d", workerID, j)
					return
				}
			}

			// Cleanup worker directory
			if err := os.RemoveAll(workerDir); err != nil {
				errors <- fmt.Errorf("worker %d failed to cleanup directory: %v", workerID, err)
				return
			}
		}(i)
	}

	// Wait for all workers to complete
	for i := 0; i < numWorkers; i++ {
		<-done
	}

	// Check for errors
	select {
	case err := <-errors:
		return fmt.Errorf("parallel execution test failed: %v", err)
	default:
		// No errors, test passed
	}

	t.Logf("Parallel execution test completed successfully")
	return nil
}

// TestLeakDetection tests for resource leaks
func TestLeakDetection(t *testing.T) {
	t.Run("File Handle Leak Detection", func(t *testing.T) {
		tempDir := os.Getenv("ACOUSTICALC_TEST_TEMP_DIR")
		if tempDir == "" {
			t.Skip("Test environment not set up")
		}

		// Create and close many files to check for leaks
		for i := 0; i < 1000; i++ {
			filePath := filepath.Join(tempDir, fmt.Sprintf("leak-test-%d.txt", i))
			file, err := os.Create(filePath)
			if err != nil {
				t.Fatalf("Failed to create file %d: %v", i, err)
			}

			if _, err := file.WriteString(fmt.Sprintf("test content %d", i)); err != nil {
				file.Close()
				t.Fatalf("Failed to write to file %d: %v", i, err)
			}

			if err := file.Close(); err != nil {
				t.Fatalf("Failed to close file %d: %v", i, err)
			}

			if err := os.Remove(filePath); err != nil {
				t.Fatalf("Failed to remove file %d: %v", i, err)
			}
		}

		t.Logf("File handle leak detection test completed successfully")
	})
}