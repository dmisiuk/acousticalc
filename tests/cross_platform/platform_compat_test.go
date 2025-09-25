package cross_platform

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/dmisiuk/acousticalc/tests/recording"
)

func TestPlatformIdentification(t *testing.T) {
	t.Run("TestIdentifyOperatingSystem", func(t *testing.T) {
		switch runtime.GOOS {
		case "linux":
			t.Log("Running on Linux")
		case "darwin":
			t.Log("Running on macOS")
		case "windows":
			t.Log("Running on Windows")
		default:
			t.Logf("Running on an unrecognized operating system: %s", runtime.GOOS)
		}
	})

	t.Run("TestPathSeparator", func(t *testing.T) {
		expectedSeparator := "/"
		if runtime.GOOS == "windows" {
			expectedSeparator = "\\"
		}

		// In Go, os.PathSeparator is a rune (int32). We need to convert it to a string for comparison.
		actualSeparator := string(os.PathSeparator)

		if actualSeparator != expectedSeparator {
			t.Errorf("Expected path separator to be '%s' on %s, but got '%s'", expectedSeparator, runtime.GOOS, actualSeparator)
		} else {
			t.Logf("Correct path separator ('%s') found for %s", actualSeparator, runtime.GOOS)
		}
	})

	t.Run("TestRecordingMechanism", func(t *testing.T) {
		tempDir := t.TempDir()
		recorder := recording.NewRecorder(tempDir, "TestRecording")

		var expectedExtension string
		if runtime.GOOS == "windows" {
			expectedExtension = ".txt"
		} else {
			expectedExtension = ".cast"
		}

		// Use 'go version' as a simple, cross-platform command that produces output.
		err := recorder.RecordCommand("go", "version")
		if err != nil {
			t.Fatalf("RecordCommand failed: %v", err)
		}

		files, err := os.ReadDir(tempDir)
		if err != nil {
			t.Fatalf("Failed to read temp dir: %v", err)
		}

		if len(files) != 1 {
			var createdFiles []string
			for _, f := range files {
				createdFiles = append(createdFiles, f.Name())
			}
			t.Fatalf("Expected 1 file in temp dir, but got %d. Files: %v", len(files), createdFiles)
		}

		createdFile := files[0].Name()
		if !strings.HasSuffix(createdFile, expectedExtension) {
			t.Errorf("Expected file with extension '%s', but got '%s'", expectedExtension, createdFile)
		} else {
			t.Logf("Correct recording artifact found: %s", createdFile)
		}
	})
}
