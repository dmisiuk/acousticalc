package reporting

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Result captures the outcome of a single E2E workflow execution.
type Result struct {
	Name          string    `json:"name"`
	Platform      string    `json:"platform"`
	ExitCode      int       `json:"exit_code"`
	Status        string    `json:"status"`
	Stdout        string    `json:"stdout"`
	Stderr        string    `json:"stderr"`
	RecordingPath string    `json:"recording_path"`
	DurationMS    int64     `json:"duration_ms"`
	Timestamp     time.Time `json:"timestamp"`
}

// Reporter aggregates workflow results into a single JSON artifact that can be
// published from CI and consumed by downstream reporting pipelines.
type Reporter struct {
	mu         sync.Mutex
	outputPath string
	results    []Result
}

// NewReporter constructs a Reporter that will emit JSON to the provided path.
func NewReporter(outputPath string) *Reporter {
	return &Reporter{outputPath: outputPath}
}

// Record adds a new result to the aggregation and immediately flushes it to
// disk. Immediate flushing keeps the artifact in sync even if an individual test
// fails after the recording is created.
func (r *Reporter) Record(result Result) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.results = append(r.results, result)

	if err := os.MkdirAll(filepath.Dir(r.outputPath), 0o755); err != nil {
		return fmt.Errorf("create report directory: %w", err)
	}

	f, err := os.Create(r.outputPath)
	if err != nil {
		return fmt.Errorf("create report file: %w", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(r.results); err != nil {
		return fmt.Errorf("encode report: %w", err)
	}

	return nil
}
