package recording

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Recorder handles the terminal recording process.
type Recorder struct {
	outputDir string
	testName  string
}

// NewRecorder creates a new recorder instance.
func NewRecorder(outputDir, testName string) *Recorder {
	return &Recorder{
		outputDir: outputDir,
		testName:  testName,
	}
}

// RecordCommand executes a command within an asciinema recording session.
func (r *Recorder) RecordCommand(command string, args ...string) error {
	if err := os.MkdirAll(r.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create recording output directory: %w", err)
	}

	filename := fmt.Sprintf("%s_%s.cast", r.testName, time.Now().Format("20060102_150405"))
	filePath := filepath.Join(r.outputDir, filename)

	fullCommand := fmt.Sprintf("%s %s", command, strings.Join(args, " "))

	// Use --quiet to suppress asciinema's own output.
	cmd := exec.Command("asciinema", "rec", "--overwrite", "--quiet", "-c", fullCommand, filePath)

	// We run this command for the side effect of creating the recording.
	// We check for an error, but don't need the output here.
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("asciinema recording failed: %w", err)
	}

	return nil
}