package e2e

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/dmisiuk/acousticalc/tests/recording"
	"github.com/dmisiuk/acousticalc/tests/reporting"
)

func TestSimpleWorkflow(t *testing.T) {
	outputDir := "../../tests/artifacts/recordings"
	testName := "simple_workflow"

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	executablePath := filepath.Join(wd, "..", "..", "cmd", "acousticalc", "acousticalc")
	if runtime.GOOS == "windows" {
		executablePath += ".exe"
	}

	// The command to run in the E2E test.
	command := []string{executablePath, "2+2"}

	recorder := recording.NewRecorder(testName, outputDir, command...)

	filePath, err := recorder.Start()
	if err != nil {
		t.Fatalf("failed to start recording: %v", err)
	}

	if err := recorder.Stop(); err != nil {
		t.Logf("got an error while stopping the recorder, but this may not be a failure: %v", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("recording file was not created: %s", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read recording file: %v", err)
	}

	passed := strings.Contains(string(content), "Result: 4")
	if !passed {
		t.Errorf("recording file does not contain the expected output")
	}

	testResults = append(testResults, reporting.E2ETestResult{
		TestName:  "TestSimpleWorkflow",
		Passed:    passed,
		Recording: filePath,
	})
}
