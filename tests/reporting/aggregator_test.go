package reporting

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestE2EReporting(t *testing.T) {
	// Create a temporary directory for our fake artifacts
	artifactsDir, err := os.MkdirTemp("", "test-artifacts-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(artifactsDir)

	// --- Setup fake artifacts ---
	platform := runtime.GOOS
	e2eDir := filepath.Join(artifactsDir, "e2e")
	test1Dir := filepath.Join(e2eDir, platform, "TestE2EWorkflow_TestSimpleAddition", "recordings")
	test2Dir := filepath.Join(e2eDir, platform, "TestE2EWorkflow_TestComplexExpression", "recordings")
	if err := os.MkdirAll(test1Dir, 0755); err != nil {
		t.Fatalf("Failed to create test1 dir: %v", err)
	}
	if err := os.MkdirAll(test2Dir, 0755); err != nil {
		t.Fatalf("Failed to create test2 dir: %v", err)
	}
	rec1Path := filepath.Join(test1Dir, "rec1.cast")
	rec2Path := filepath.Join(test2Dir, "rec2.cast")
	if _, err := os.Create(rec1Path); err != nil {
		t.Fatalf("Failed to create rec1: %v", err)
	}
	if _, err := os.Create(rec2Path); err != nil {
		t.Fatalf("Failed to create rec2: %v", err)
	}

	t.Run("TestDiscoverE2EArtifacts", func(t *testing.T) {
		// Run the discovery function
		reportData, err := DiscoverE2EArtifacts(artifactsDir)
		if err != nil {
			t.Fatalf("DiscoverE2EArtifacts failed: %v", err)
		}

		// Assert the results
		if len(reportData.TestRuns) != 2 {
			t.Errorf("Expected 2 test runs, but got %d", len(reportData.TestRuns))
		}

		// Check the details of the first test run
		foundTest1 := false
		for _, run := range reportData.TestRuns {
			if run.Name == "TestE2EWorkflow_TestSimpleAddition" {
				foundTest1 = true
				if run.Platform != platform {
					t.Errorf("Expected platform '%s', but got '%s'", platform, run.Platform)
				}
				expectedRelPath := filepath.Join("e2e", platform, "TestE2EWorkflow_TestSimpleAddition", "recordings", "rec1.cast")
				if run.Recording != expectedRelPath {
					t.Errorf("Expected recording path '%s', but got '%s'", expectedRelPath, run.Recording)
				}
			}
		}
		if !foundTest1 {
			t.Error("Test run 'TestE2EWorkflow_TestSimpleAddition' not found")
		}
	})

	t.Run("TestGenerateE2EReport", func(t *testing.T) {
		// First, get the data
		reportData, err := DiscoverE2EArtifacts(artifactsDir)
		if err != nil {
			t.Fatalf("DiscoverE2EArtifacts failed: %v", err)
		}

		// Generate the report
		reportOutputDir := filepath.Join(artifactsDir, "reports")
		reportPath, err := GenerateE2EReport(reportData, reportOutputDir)
		if err != nil {
			t.Fatalf("GenerateE2EReport failed: %v", err)
		}

		// Verify the report was created
		if _, err := os.Stat(reportPath); os.IsNotExist(err) {
			t.Fatalf("Report file was not created at %s", reportPath)
		}

		// Verify the report content
		content, err := os.ReadFile(reportPath)
		if err != nil {
			t.Fatalf("Failed to read report file: %v", err)
		}
		htmlContent := string(content)

		if !strings.Contains(htmlContent, "TestE2EWorkflow_TestSimpleAddition") {
			t.Error("Report does not contain 'TestE2EWorkflow_TestSimpleAddition'")
		}
		if !strings.Contains(htmlContent, platform) {
			t.Errorf("Report does not contain platform '%s'", platform)
		}
	})
}