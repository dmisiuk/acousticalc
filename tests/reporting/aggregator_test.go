package reporting

import (
	"os"
	"path/filepath"
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

	// --- Setup fake artifacts for Linux and Windows ---
	e2eDir := filepath.Join(artifactsDir, "e2e")

	// Linux artifact
	linuxTestDir := filepath.Join(e2eDir, "linux", "TestE2EWorkflow_TestSimpleAddition", "recordings")
	if err := os.MkdirAll(linuxTestDir, 0755); err != nil {
		t.Fatalf("Failed to create linux test dir: %v", err)
	}
	linuxRecPath := filepath.Join(linuxTestDir, "rec1.cast")
	if _, err := os.Create(linuxRecPath); err != nil {
		t.Fatalf("Failed to create linux rec: %v", err)
	}

	// Windows artifact
	windowsTestDir := filepath.Join(e2eDir, "windows", "TestE2EWorkflow_TestComplexExpression", "recordings")
	if err := os.MkdirAll(windowsTestDir, 0755); err != nil {
		t.Fatalf("Failed to create windows test dir: %v", err)
	}
	windowsRecPath := filepath.Join(windowsTestDir, "rec2.txt")
	if _, err := os.Create(windowsRecPath); err != nil {
		t.Fatalf("Failed to create windows rec: %v", err)
	}

	t.Run("TestDiscoverCrossPlatformArtifacts", func(t *testing.T) {
		// Run the discovery function
		reportData, err := DiscoverE2EArtifacts(artifactsDir)
		if err != nil {
			t.Fatalf("DiscoverE2EArtifacts failed: %v", err)
		}

		// Assert the results
		if len(reportData.TestRuns) != 2 {
			t.Fatalf("Expected 2 test runs, but got %d", len(reportData.TestRuns))
		}

		// Check for the Linux run
		foundLinuxRun := false
		for _, run := range reportData.TestRuns {
			if run.Platform == "linux" {
				foundLinuxRun = true
				if run.Name != "TestE2EWorkflow_TestSimpleAddition" {
					t.Errorf("Incorrect test name for linux run: got %s", run.Name)
				}
				if !strings.HasSuffix(run.Recording, ".cast") {
					t.Errorf("Incorrect recording file for linux run: got %s", run.Recording)
				}
			}
		}
		if !foundLinuxRun {
			t.Error("Did not find the linux test run")
		}

		// Check for the Windows run
		foundWindowsRun := false
		for _, run := range reportData.TestRuns {
			if run.Platform == "windows" {
				foundWindowsRun = true
				if run.Name != "TestE2EWorkflow_TestComplexExpression" {
					t.Errorf("Incorrect test name for windows run: got %s", run.Name)
				}
				if !strings.HasSuffix(run.Recording, ".txt") {
					t.Errorf("Incorrect recording file for windows run: got %s", run.Recording)
				}
			}
		}
		if !foundWindowsRun {
			t.Error("Did not find the windows test run")
		}
	})

	t.Run("TestGenerateE2EReportWithCrossPlatformData", func(t *testing.T) {
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

		// Verify the report content
		content, err := os.ReadFile(reportPath)
		if err != nil {
			t.Fatalf("Failed to read report file: %v", err)
		}
		htmlContent := string(content)

		if !strings.Contains(htmlContent, "<td>linux</td>") {
			t.Error("Report does not contain the linux platform")
		}
		if !strings.Contains(htmlContent, "<td>windows</td>") {
			t.Error("Report does not contain the windows platform")
		}
		if !strings.Contains(htmlContent, ".cast") {
			t.Error("Report does not contain a link to the .cast file")
		}
		if !strings.Contains(htmlContent, ".txt") {
			t.Error("Report does not contain a link to the .txt file")
		}
	})
}