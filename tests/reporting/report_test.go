package reporting

import (
	"os"
	"strings"
	"testing"
)

func TestGenerateE2EReport(t *testing.T) {
	outputDir := t.TempDir()
	testName := "test_e2e_report"

	results := []E2ETestResult{
		{TestName: "TestSimpleWorkflow", Passed: true, Recording: "../recording/simple_workflow_recording.cast"},
		{TestName: "TestComplexWorkflow", Passed: false, Recording: "../recording/complex_workflow_recording.cast"},
	}

	reportPath, err := GenerateE2EReport(testName, outputDir, results)
	if err != nil {
		t.Fatalf("failed to generate E2E report: %v", err)
	}

	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		t.Errorf("report file was not created: %s", reportPath)
	}

	content, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("failed to read report file: %v", err)
	}

	if !strings.Contains(string(content), "TestSimpleWorkflow") {
		t.Errorf("report does not contain the first test name")
	}

	if !strings.Contains(string(content), "TestComplexWorkflow") {
		t.Errorf("report does not contain the second test name")
	}

	if !strings.Contains(string(content), "PASSED") {
		t.Errorf("report does not contain the PASSED status")
	}

	if !strings.Contains(string(content), "FAILED") {
		t.Errorf("report does not contain the FAILED status")
	}
}
