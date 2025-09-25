package e2e

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dmisiuk/acousticalc/tests/reporting"
)

func TestGenerateE2EReport(t *testing.T) {
	recordingsDir := "../../tests/artifacts/recordings"
	reportsDir := "../../tests/artifacts/reports"
	testName := "e2e_workflow"

	// Create a dummy recording file
	dummyRecordingPath := filepath.Join(recordingsDir, "simple_workflow_recording_20250924_174752.cast")
	if err := os.WriteFile(dummyRecordingPath, []byte("dummy content"), 0644); err != nil {
		t.Fatalf("failed to create dummy recording file: %v", err)
	}

	// In a real scenario, these results would be dynamically generated.
	results := []reporting.E2ETestResult{
		{TestName: "TestSimpleWorkflow", Passed: true, Recording: findLatestRecording(t, recordingsDir)},
	}

	_, err := reporting.GenerateE2EReport(testName, reportsDir, results)
	if err != nil {
		t.Fatalf("failed to generate E2E report: %v", err)
	}
}

func findLatestRecording(t *testing.T, dir string) string {
	t.Helper()
	var latestFile string
	var latestModTime time.Time

	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read recordings directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".cast" {
			info, err := file.Info()
			if err != nil {
				t.Fatalf("failed to get file info: %v", err)
			}
			if info.ModTime().After(latestModTime) {
				latestModTime = info.ModTime()
				latestFile = filepath.Join("../recordings", file.Name())
			}
		}
	}

	if latestFile == "" {
		t.Fatal("no recording files found")
	}

	return latestFile
}