package visual

import (
	"context"
	"os"
	"testing"
	"time"
)

// TestRefactoredArchitecture validates the architectural improvements
func TestRefactoredArchitecture(t *testing.T) {
	testDir := t.TempDir()
	ctx := context.Background()

	t.Run("interface_compliance", func(t *testing.T) {
		// Test interface compliance
		var _ ScreenshotCapturer = &ScreenshotCapture{}
		var _ ArtifactGeneratorInterface = &ArtifactGenerator{}

		// Test screenshot capturer interface
		capture := NewScreenshotCapture("arch_test", testDir)
		capture.SetOutputDir(testDir)
		capture.SetTestName("interface_test")

		if capture.TestName != "interface_test" {
			t.Errorf("SetTestName didn't work: got %s, want interface_test", capture.TestName)
		}
	})

	t.Run("performance_monitoring", func(t *testing.T) {
		monitor := NewPerformanceMonitor(ctx)
		defer monitor.Stop()

		// Track a mock operation
		err := monitor.TrackOperation("test_operation", func() error {
			time.Sleep(10 * time.Millisecond)
			return nil
		})

		if err != nil {
			t.Fatalf("TrackOperation failed: %v", err)
		}

		// Verify metrics were recorded
		metrics := monitor.GetMetrics()
		if len(metrics) != 1 {
			t.Errorf("Expected 1 metric, got %d", len(metrics))
		}

		testMetric, exists := metrics["test_operation"]
		if !exists {
			t.Error("test_operation metric not found")
		}

		if testMetric.Count != 1 {
			t.Errorf("Expected count 1, got %d", testMetric.Count)
		}

		if testMetric.AverageTime < 5*time.Millisecond {
			t.Errorf("Average time seems too low: %v", testMetric.AverageTime)
		}

		// Test threshold validation
		violations := monitor.CheckThresholds()
		// Should be no violations for our test operation
		if len(violations) > 0 {
			t.Logf("Threshold violations (expected for test): %v", violations)
		}
	})

	t.Run("artifact_generator_enhancements", func(t *testing.T) {
		generator := NewArtifactGeneratorWithContext(ctx, "arch_test", testDir)
		defer generator.Close()

		// Test directory creation
		err := generator.CreateDirectoryStructure()
		if err != nil {
			t.Fatalf("CreateDirectoryStructure failed: %v", err)
		}

		// Verify directories were created
		expectedDirs := []string{
			"reports", "demo_content/storyboards", "demo_content/assets",
			"demo_content/metadata", "charts",
		}

		for _, dir := range expectedDirs {
			fullPath := testDir + "/" + dir
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				t.Errorf("Expected directory %s was not created", fullPath)
			}
		}

		// Test artifact summary
		summary := generator.GetArtifactSummary()
		if summary.TotalArtifacts < 0 {
			t.Error("Invalid total artifacts count")
		}

		// Test metadata generation
		err = generator.GenerateMetadataFiles()
		if err != nil {
			t.Fatalf("GenerateMetadataFiles failed: %v", err)
		}

		// Test thread-safe report addition
		generator.addReport(ReportInfo{
			Filename:    "test_report.html",
			Type:        "test_report",
			Timestamp:   time.Now(),
			Size:        1024,
			Description: "Test report for architecture validation",
		})

		// Verify report was added
		summaryAfter := generator.GetArtifactSummary()
		if summaryAfter.TotalArtifacts != summary.TotalArtifacts+1 {
			t.Errorf("Report was not added correctly: before=%d, after=%d",
				summary.TotalArtifacts, summaryAfter.TotalArtifacts)
		}
	})

	t.Run("context_cancellation", func(t *testing.T) {
		testCtx, cancel := context.WithCancel(ctx)

		generator := NewArtifactGeneratorWithContext(testCtx, "cancel_test", testDir)
		monitor := NewPerformanceMonitor(testCtx)

		// Cancel the context
		cancel()

		// Test that components handle cancellation gracefully
		err := generator.Close()
		if err != nil {
			t.Errorf("Generator.Close() returned error: %v", err)
		}

		monitor.Stop()

		// Verify context was cancelled
		select {
		case <-testCtx.Done():
			// Good, context was cancelled
		default:
			t.Error("Context was not cancelled")
		}
	})

	t.Run("concurrent_safety", func(t *testing.T) {
		generator := NewArtifactGeneratorWithContext(ctx, "concurrent_test", testDir)
		defer generator.Close()

		// Test concurrent artifact addition
		done := make(chan bool, 10)

		for i := 0; i < 10; i++ {
			go func(id int) {
				generator.addReport(ReportInfo{
					Filename:    "concurrent_report.html",
					Type:        "concurrent_test",
					Timestamp:   time.Now(),
					Size:        int64(100 * id),
					Description: "Concurrent test report",
				})
				done <- true
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < 10; i++ {
			<-done
		}

		// Verify all reports were added
		summary := generator.GetArtifactSummary()
		if summary.TotalArtifacts != 10 {
			t.Errorf("Expected 10 artifacts, got %d", summary.TotalArtifacts)
		}
	})
}

// TestArchitecturalPatterns validates the design patterns implemented
func TestArchitecturalPatterns(t *testing.T) {
	testDir := t.TempDir()
	ctx := context.Background()

	t.Run("factory_pattern", func(t *testing.T) {
		// Test different factory methods produce compatible instances
		generator1 := NewArtifactGenerator("test1", testDir)
		generator2 := NewArtifactGeneratorWithContext(ctx, "test2", testDir)

		defer generator1.Close()
		defer generator2.Close()

		// Both should implement the same interface
		var iface1 ArtifactGeneratorInterface = generator1
		var iface2 ArtifactGeneratorInterface = generator2

		// Both should have working functionality
		if err := iface1.CreateDirectoryStructure(); err != nil {
			t.Errorf("Factory method 1 failed: %v", err)
		}

		if err := iface2.CreateDirectoryStructure(); err != nil {
			t.Errorf("Factory method 2 failed: %v", err)
		}
	})

	t.Run("dependency_injection", func(t *testing.T) {
		// Test that components can be injected with different implementations
		capture := NewScreenshotCapture("di_test", testDir)

		// The capturer engine can be swapped out
		originalEngine := capture.capturer
		if originalEngine == nil {
			t.Error("No default engine was injected")
		}

		// Test engine interface
		var _ ScreenshotEngine = &RobotGoEngine{}
	})

	t.Run("template_method_pattern", func(t *testing.T) {
		// The GenerateComprehensiveArtifacts method follows template method pattern
		generator := NewArtifactGenerator("template_test", testDir)
		defer generator.Close()

		logger := NewVisualTestLogger("template_test", testDir)

		// This method orchestrates multiple smaller operations
		err := generator.GenerateComprehensiveArtifacts(logger)
		if err != nil {
			t.Errorf("Template method failed: %v", err)
		}

		// Verify that the template method coordinated multiple operations
		summary := generator.GetArtifactSummary()
		if len(summary.ArtifactTypes) == 0 {
			t.Error("Template method did not generate any artifact types")
		}
	})
}