package e2e

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

var (
	repoRootOnce sync.Once
	repoRootPath string
	repoRootErr  error
)

func repositoryRoot() (string, error) {
	repoRootOnce.Do(func() {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			repoRootErr = fmt.Errorf("unable to determine caller information for E2E harness")
			return
		}

		repoRootPath = filepath.Clean(filepath.Join(filepath.Dir(filename), "..", ".."))
	})

	return repoRootPath, repoRootErr
}

// CLIRunner orchestrates CLI invocations using `go run` so the end-to-end tests
// remain hermetic without depending on pre-built binaries. The runner supports
// context based cancellation to avoid hanging CI jobs when something goes wrong.
type CLIRunner struct {
	repo    string
	env     []string
	timeout time.Duration
}

// NewCLIRunner constructs a runner with repository-relative execution.
func NewCLIRunner(timeout time.Duration) (*CLIRunner, error) {
	root, err := repositoryRoot()
	if err != nil {
		return nil, err
	}

	return &CLIRunner{
		repo:    root,
		env:     os.Environ(),
		timeout: timeout,
	}, nil
}

// Run executes the acousticalc CLI with the provided arguments.
type RunResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Started  time.Time
	Duration time.Duration
}

func (r *CLIRunner) Run(args ...string) (*RunResult, error) {
	if r.repo == "" {
		return nil, fmt.Errorf("runner is not initialised with repository root")
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	cmdArgs := append([]string{"run", "./cmd/acousticalc"}, args...)
	cmd := exec.CommandContext(ctx, "go", cmdArgs...)
	cmd.Dir = r.repo
	cmd.Env = r.env

	stdout, stderr := &safeBuffer{}, &safeBuffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	started := time.Now()
	err := cmd.Run()
	duration := time.Since(started)

	exitCode := 0
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("command timeout after %s", r.timeout)
		}

		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return nil, err
		}
	}

	return &RunResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
		Started:  started,
		Duration: duration,
	}, nil
}

// safeBuffer is a threadsafe bytes.Buffer alternative that protects against
// concurrent writes from stdout/stderr pipes when commands emit interleaving
// output.
type safeBuffer struct {
	mu  sync.Mutex
	buf []byte
}

func (b *safeBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.buf = append(b.buf, p...)
	return len(p), nil
}

func (b *safeBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return string(b.buf)
}
