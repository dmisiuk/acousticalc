package cross_platform

import (
	"os"
	"runtime"
	"testing"
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
}