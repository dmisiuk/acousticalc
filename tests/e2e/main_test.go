package e2e

import (
	"os"
	"testing"

	"github.com/dmisiuk/acousticalc/tests/reporting"
)

var testResults []reporting.E2ETestResult

func TestMain(m *testing.M) {
	// Run all tests
	exitCode := m.Run()

	// Generate the report after all tests have run
	_, err := reporting.GenerateE2EReport("e2e_workflow", "../../tests/artifacts/reports", testResults)
	if err != nil {
		// We can't fail the test here, but we can log the error.
		println("failed to generate E2E report:", err.Error())
	}

	os.Exit(exitCode)
}