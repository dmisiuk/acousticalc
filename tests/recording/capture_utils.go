package recording

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// CaptureScreen captures the entire screen and saves it to a file.
func CaptureScreen(filename string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		// Ensure scrot is installed: sudo apt-get install scrot
		cmd = exec.Command("scrot", "-q", "100", filename)
	case "darwin":
		cmd = exec.Command("screencapture", "-x", filename)
	default:
		return fmt.Errorf("screen capture not supported on %s", runtime.GOOS)
	}

	// Create directory if it doesn't exist
	// This is a simplified approach; consider a more robust directory creation strategy
	if err := os.MkdirAll("tests/artifacts/recordings", 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return cmd.Run()
}

// RecordScreen records the screen for a given duration and saves it to a file.
// This is a placeholder and requires a screen recording tool like ffmpeg.
func RecordScreen(filename string, duration time.Duration) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("screen recording is only supported on Linux with ffmpeg")
	}

	// Ensure ffmpeg is installed: sudo apt-get install ffmpeg
	// Also requires a running X server.
	display := os.Getenv("DISPLAY")
	if display == "" {
		return fmt.Errorf("DISPLAY environment variable not set, cannot record screen")
	}

	cmd := exec.Command(
		"ffmpeg",
		"-y", // Overwrite output file if it exists
		"-f", "x11grab",
		"-s", "1920x1080", // Screen size
		"-i", display,
		"-t", fmt.Sprintf("%.0f", duration.Seconds()),
		"-r", "25", // Frame rate
		filename,
	)

	return cmd.Run()
}