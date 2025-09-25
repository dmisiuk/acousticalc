package recording

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Recorder handles terminal recording operations.
type Recorder struct {
	outputDir string
	testName  string
	command   []string
	cmd       *exec.Cmd
}

// NewRecorder creates a new Recorder instance.
// If no command is provided, it will start an interactive recording.
func NewRecorder(testName, outputDir string, command ...string) *Recorder {
	return &Recorder{
		testName:  testName,
		outputDir: outputDir,
		command:   command,
	}
}

// Start begins a new terminal recording.
func (r *Recorder) Start() (string, error) {
	if err := os.MkdirAll(r.outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	filename := fmt.Sprintf("%s_recording_%s.cast",
		r.testName,
		time.Now().Format("20060102_150405"))
	filePath := filepath.Join(r.outputDir, filename)

	var cmdArgs []string
	cmdArgs = append(cmdArgs, "rec", "--overwrite", "-q")
	if len(r.command) > 0 {
		cmdArgs = append(cmdArgs, "-c", strings.Join(r.command, " "))
	}
	cmdArgs = append(cmdArgs, filePath)

	r.cmd = exec.Command("asciinema", cmdArgs...)

	if err := r.cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start asciinema recording: %w", err)
	}

	return filePath, nil
}

// Stop terminates the current recording.
func (r *Recorder) Stop() error {
	if r.cmd == nil || r.cmd.Process == nil {
		return fmt.Errorf("recording not started")
	}

	// If no command was specified, it's an interactive recording that needs to be interrupted.
	if len(r.command) == 0 {
		// A short sleep to ensure the process is ready to be interrupted.
		time.Sleep(250 * time.Millisecond)
		if err := r.cmd.Process.Signal(os.Interrupt); err != nil {
			// It might have already exited, so we wait and check the error.
			if waitErr := r.cmd.Wait(); waitErr != nil {
				return fmt.Errorf("failed to interrupt and wait for asciinema: %v, initial error: %w", waitErr, err)
			}
			return nil // Exited cleanly after all
		}
	}

	// For both command-based and interrupted recordings, we must Wait().
	if err := r.cmd.Wait(); err != nil {
		// Ignore "exit status 1" for interactive recordings, as this is expected on interrupt.
		if len(r.command) == 0 && strings.Contains(err.Error(), "exit status 1") {
			return nil
		}
		return fmt.Errorf("failed to wait for asciinema command: %w", err)
	}
	return nil
}
