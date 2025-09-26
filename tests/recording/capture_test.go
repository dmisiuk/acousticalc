package recording

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestScreenCapture(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Screen capture tests are not yet supported on Windows in this test suite")
	}

	t.Run("TestCaptureScreen", func(t *testing.T) {
		dir := "tests/artifacts/recordings"
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create artifacts directory: %v", err)
		}
		filename := fmt.Sprintf("%s/capture_%s.png", dir, time.Now().Format("20060102_150405"))

		err := CaptureScreen(filename)
		if err != nil {
			t.Fatalf("Failed to capture screen: %v", err)
		}

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			t.Errorf("Capture file was not created at %s", filename)
		} else {
			t.Logf("Screen captured successfully to %s", filename)
		}
	})
}

func TestRecordScreen(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Screen recording tests are currently supported only on Linux")
	}

	t.Run("TestRecordScreenForDuration", func(t *testing.T) {
		dir := "tests/artifacts/recordings"
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create artifacts directory: %v", err)
		}
		filename := fmt.Sprintf("%s/recording_%s.mp4", dir, time.Now().Format("20060102_150405"))
		duration := 3 * time.Second

		err := RecordScreen(filename, duration)
		if err != nil {
			t.Fatalf("Failed to record screen: %v", err)
		}

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			t.Errorf("Recording file was not created at %s", filename)
		} else {
			t.Logf("Screen recorded successfully to %s", filename)
		}
	})
}
