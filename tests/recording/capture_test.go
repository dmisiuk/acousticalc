package recording

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestRecorder_Command(t *testing.T) {
	outputDir := t.TempDir()
	testName := "test_command_recording"
	command := []string{"echo", "hello from test"}

	recorder := NewRecorder(testName, outputDir, command...)

	filePath, err := recorder.Start()
	if err != nil {
		t.Fatalf("failed to start recording: %v", err)
	}

	if err := recorder.Stop(); err != nil {
		// It is possible for the process to exit so quickly that an error is not returned.
		// We will check for the file content to be sure.
		t.Logf("got an error while stopping the recorder, but this may not be a failure: %v", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("recording file was not created: %s", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read recording file: %v", err)
	}

	if !strings.Contains(string(content), "hello from test") {
		t.Errorf("recording file does not contain the expected output")
	}
}

func TestRecorder_Interactive(t *testing.T) {
	outputDir := t.TempDir()
	testName := "test_interactive_recording"

	recorder := NewRecorder(testName, outputDir)

	filePath, err := recorder.Start()
	if err != nil {
		t.Fatalf("failed to start recording: %v", err)
	}

	// Give it a moment to run.
	time.Sleep(500 * time.Millisecond)

	if err := recorder.Stop(); err != nil {
		t.Fatalf("failed to stop recording: %v", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("recording file was not created: %s", filePath)
	}
}

func TestNewRecorder(t *testing.T) {
	outputDir := t.TempDir()
	testName := "test_new_recorder"

	recorder := NewRecorder(testName, outputDir)

	if recorder.testName != testName {
		t.Errorf("expected testName to be %s, but got %s", testName, recorder.testName)
	}

	if recorder.outputDir != outputDir {
		t.Errorf("expected outputDir to be %s, but got %s", outputDir, recorder.outputDir)
	}
}

func TestRecorder_Stop_NotStarted(t *testing.T) {
	outputDir := t.TempDir()
	testName := "test_stop_not_started"

	recorder := NewRecorder(testName, outputDir)

	if err := recorder.Stop(); err == nil {
		t.Errorf("expected an error when stopping a recorder that was not started, but got nil")
	}
}
