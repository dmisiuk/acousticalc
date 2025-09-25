package e2e

import (
	"testing"

	"github.com/dmisiuk/acousticalc/tests/reporting"
)

func TestGenerateE2EReport(t *testing.T) {
	outputDir := "../../tests/artifacts/reports"
	testName := "e2e_workflow"

	// In a real scenario, these results would be dynamically generated.
	results := []reporting.E2ETestResult{
		{TestName: "TestSimpleWorkflow", Passed: true, Recording: "../recordings/simple_workflow_recording.cast"},
	}

	_, err := reporting.GenerateE2EReport(testName, outputDir, results)
	if err != nil {
		t.Fatalf("failed to generate E2E report: %v", err)
	}
}