package e2e

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dmisiuk/acousticalc/tests/cross_platform"
	"github.com/dmisiuk/acousticalc/tests/recording"
	"github.com/dmisiuk/acousticalc/tests/reporting"
)

type workflow struct {
	name                string
	args                []string
	expectedExitCode    int
	expectStdoutSubstrs []string
	expectStderrSubstrs []string
	platformOverrides   map[string]platformExpectations
}

type platformExpectations struct {
	expectedExitCode    *int
	expectStdoutSubstrs []string
	expectStderrSubstrs []string
}

func TestAcoustiCalcWorkflows(t *testing.T) {
	runner, err := NewCLIRunner(15 * time.Second)
	if err != nil {
		t.Fatalf("failed to create CLI runner: %v", err)
	}

	root, err := repositoryRoot()
	if err != nil {
		t.Fatalf("failed to resolve repository root: %v", err)
	}

	recorder, err := recording.NewSessionRecorder(filepath.Join(root, "tests", "recording", "artifacts"))
	if err != nil {
		t.Fatalf("failed to create session recorder: %v", err)
	}

	reporter := reporting.NewReporter(filepath.Join(root, "tests", "reporting", "artifacts", "e2e_report.json"))

	workflows := []workflow{
		{
			name:             "basic-evaluation",
			args:             []string{"2+3*4"},
			expectedExitCode: 0,
			expectStdoutSubstrs: []string{
				"Result: 14",
			},
		},
		{
			name:                "whitespace-evaluation",
			args:                []string{"2 + 3 * 4"},
			expectedExitCode:    0,
			expectStdoutSubstrs: []string{"Result: 14"},
		},
		{
			name:                "invalid-expression",
			args:                []string{"2 + (3"},
			expectedExitCode:    1,
			expectStdoutSubstrs: []string{"Error:"},
		},
		{
			name:                "no-arguments",
			args:                []string{},
			expectedExitCode:    1,
			expectStdoutSubstrs: []string{"Usage: acousticalc"},
			platformOverrides: map[string]platformExpectations{
				"windows": {
					expectStdoutSubstrs: []string{"Usage: acousticalc"},
				},
			},
		},
	}

	platform := cross_platform.Platform()

	for _, wf := range workflows {
		wf := wf
		t.Run(wf.name, func(t *testing.T) {
			t.Parallel()

			status := "failed"
			var runResult *RunResult
			var recordingPath string

			t.Cleanup(func() {
				if runResult != nil {
					rel := recordingPath
					if rel == "" {
						path, err := recorder.RecordSession(wf.name, append([]string{"acousticalc"}, wf.args...), runResult.Stdout, runResult.Stderr, runResult.ExitCode, runResult.Started, runResult.Duration)
						if err != nil {
							t.Logf("failed to create recording: %v", err)
						} else {
							recordingPath = path
						}
						rel = recordingPath
					}

					if rel != "" {
						if relative, err := filepath.Rel(root, rel); err == nil {
							rel = relative
						}
					}

					if err := reporter.Record(reporting.Result{
						Name:          wf.name,
						Platform:      platform,
						ExitCode:      runResult.ExitCode,
						Status:        status,
						Stdout:        runResult.Stdout,
						Stderr:        runResult.Stderr,
						RecordingPath: rel,
						DurationMS:    runResult.Duration.Milliseconds(),
						Timestamp:     runResult.Started,
					}); err != nil {
						t.Logf("failed to write report: %v", err)
					}
				}
			})

			result, err := runner.Run(wf.args...)
			if err != nil {
				t.Fatalf("failed to run CLI: %v", err)
			}

			runResult = result

			expectations := wf
			if override, ok := wf.platformOverrides[platform]; ok {
				if override.expectedExitCode != nil {
					expectations.expectedExitCode = *override.expectedExitCode
				}
				if len(override.expectStdoutSubstrs) > 0 {
					expectations.expectStdoutSubstrs = override.expectStdoutSubstrs
				}
				if len(override.expectStderrSubstrs) > 0 {
					expectations.expectStderrSubstrs = override.expectStderrSubstrs
				}
			}

			if expectations.expectedExitCode != result.ExitCode {
				t.Fatalf("expected exit code %d, got %d", expectations.expectedExitCode, result.ExitCode)
			}

			normalizedStdout := cross_platform.NormalizeNewlines(strings.TrimSpace(result.Stdout))
			for _, substr := range expectations.expectStdoutSubstrs {
				if !strings.Contains(normalizedStdout, substr) {
					t.Fatalf("expected stdout to contain %q, got %q", substr, normalizedStdout)
				}
			}

			normalizedStderr := cross_platform.NormalizeNewlines(strings.TrimSpace(result.Stderr))
			for _, substr := range expectations.expectStderrSubstrs {
				if !strings.Contains(normalizedStderr, substr) {
					t.Fatalf("expected stderr to contain %q, got %q", substr, normalizedStderr)
				}
			}

			status = "passed"
		})
	}
}
