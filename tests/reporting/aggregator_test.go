package reporting

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResultAggregator(t *testing.T) {
	t.Run("TestAggregateResults", func(t *testing.T) {
		artifactsDir := "tests/artifacts/temp_aggregator"
		reportsDir := filepath.Join(artifactsDir, "reports")
		finalReport := filepath.Join(reportsDir, "final_report.txt")

		// Create dummy artifact files
		os.MkdirAll(filepath.Join(artifactsDir, "e2e"), 0755)
		os.MkdirAll(filepath.Join(artifactsDir, "performance"), 0755)
		os.WriteFile(filepath.Join(artifactsDir, "e2e", "results.xml"), []byte("<testsuite tests=\"1\"></testsuite>"), 0644)
		os.WriteFile(filepath.Join(artifactsDir, "performance", "summary.txt"), []byte("avg_latency: 50ms"), 0644)

		err := AggregateResults(artifactsDir, finalReport)
		if err != nil {
			t.Fatalf("Failed to aggregate results: %v", err)
		}

		if _, err := os.Stat(finalReport); os.IsNotExist(err) {
			t.Errorf("Final report was not created at %s", finalReport)
		}

		// Clean up dummy files
		os.RemoveAll(artifactsDir)
	})
}