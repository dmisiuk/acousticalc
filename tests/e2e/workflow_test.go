package e2e

import (
	"os"
	"strings"
	"testing"

	"github.com/dmisiuk/acousticalc/tests/recording"
)

func TestSimpleWorkflow(t *testing.T) {
	outputDir := t.TempDir()
	testName := "simple_workflow"

	// The command to run in the E2E test.
	// In a real scenario, this would be the command to run the application.
	command := []string{"echo", "starting application..."}

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

	if !strings.Contains(string(content), "starting application...") {
		t.Errorf("recording file does not contain the expected output")
	}
}