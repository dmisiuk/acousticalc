package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/dmisiuk/acousticalc/tests/recording"
)

const (
	binaryName = "acousticalc"
)

var (
	binaryPath string
)

// E2ETestRun manages the context for a single E2E test run.
type E2ETestRun struct {
	t         *testing.T
	testName  string
	outputDir string
	recorder  *recording.Recorder
}

// NewE2ETestRun creates a new E2E test run.
func NewE2ETestRun(t *testing.T) *E2ETestRun {
	testName := strings.ReplaceAll(t.Name(), "/", "_")
	platform := runtime.GOOS
	outputDir := filepath.Join("../../tests/artifacts/e2e", platform, testName)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create artifact directory '%s': %v", outputDir, err)
	}

	recDir := filepath.Join(outputDir, "recordings")

	return &E2ETestRun{
		t:         t,
		testName:  testName,
		outputDir: outputDir,
		recorder:  recording.NewRecorder(recDir, testName),
	}
}

// RecordCommand runs the acousticalc binary within a recording session for its side effect.
func (r *E2ETestRun) RecordCommand(args ...string) {
	r.t.Helper()

	if binaryPath == "" {
		r.t.Fatal("E2E test setup failed: binaryPath is not set. Was TestMain run correctly?")
	}

	err := r.recorder.RecordCommand(binaryPath, args...)
	if err != nil {
		r.t.Fatalf("Failed to record command: %v", err)
	}
	r.t.Logf("Command recorded successfully.")
}

// RunCommand runs the acousticalc binary and returns its output for assertions.
func (r *E2ETestRun) RunCommand(args ...string) string {
	r.t.Helper()

	if binaryPath == "" {
		r.t.Fatal("E2E test setup failed: binaryPath is not set. Was TestMain run correctly?")
	}

	cmd := exec.Command(binaryPath, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		r.t.Fatalf("Failed to run command for assertion: %v, output: %s", err, output)
	}

	return strings.TrimSpace(string(output))
}

// setupE2ETests builds the binary for the E2E tests and returns a teardown function.
func setupE2ETests() (func(), error) {
	tempDir, err := os.MkdirTemp("", "acousticalc-e2e-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir for binary: %w", err)
	}

	var binName string
	if runtime.GOOS == "windows" {
		binName = binaryName + ".exe"
	} else {
		binName = binaryName
	}

	path := filepath.Join(tempDir, binName)

	cmd := exec.Command("go", "build", "-o", path, "../../cmd/acousticalc")
	buildOutput, err := cmd.CombinedOutput()
	if err != nil {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("failed to build binary: %v\nOutput: %s", err, buildOutput)
	}

	binaryPath = path

	teardown := func() {
		os.RemoveAll(tempDir)
	}

	return teardown, nil
}
