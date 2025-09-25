package recording

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

// RecordCommand executes a command within a recording session.
// It uses asciinema on Linux/macOS and PowerShell Start-Transcript on Windows.
func (r *Recorder) RecordCommand(command string, args ...string) error {
	if err := os.MkdirAll(r.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create recording output directory: %w", err)
	}

	if runtime.GOOS == "windows" {
		return r.recordWithPowerShell(command, args...)
	}
	return r.recordWithAsciinema(command, args...)
}

func (r *Recorder) recordWithAsciinema(command string, args ...string) error {
	filename := fmt.Sprintf("%s_%s.cast", r.testName, time.Now().Format("20060102_150405"))
	filePath := filepath.Join(r.outputDir, filename)

	fullCommand := fmt.Sprintf("%s %s", command, strings.Join(args, " "))

	// Use --quiet to suppress asciinema's own output.
	cmd := exec.Command("asciinema", "rec", "--overwrite", "--quiet", "-c", fullCommand, filePath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("asciinema recording failed: %w", err)
	}
	return nil
}

func (r *Recorder) recordWithPowerShell(command string, args ...string) error {
	filename := fmt.Sprintf("%s_%s.txt", r.testName, time.Now().Format("20060102_150405"))
	filePath := filepath.Join(r.outputDir, filename)

	// Escape arguments for PowerShell
	var escapedArgs []string
	for _, arg := range args {
		escapedArgs = append(escapedArgs, fmt.Sprintf("'%s'", strings.ReplaceAll(arg, "'", "''")))
	}

	// Construct the PowerShell command string
	// Start-Transcript logs the session.
	// & executes the command.
	// Stop-Transcript stops logging.
	// We use -join to handle expressions with spaces correctly.
	fullCommand := fmt.Sprintf("& '%s' %s", command, strings.Join(escapedArgs, " "))
	psCommand := fmt.Sprintf("Start-Transcript -Path '%s' -NoClobber; try { %s } finally { Stop-Transcript }", filePath, fullCommand)

	cmd := exec.Command("powershell.exe", "-NoProfile", "-Command", psCommand)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("powershell recording failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}
