package recording

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// SessionRecorder produces asciinema compatible `.cast` files so GitHub
// Actions jobs can publish demo-ready recordings without requiring the
// asciinema binary. The format is JSON based and easy to generate from Go.
type SessionRecorder struct {
	mu      sync.Mutex
	outDir  string
	width   int
	height  int
	started time.Time
}

// Metadata represents the header of an asciinema v2 recording file.
type Metadata struct {
	Version   int    `json:"version"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Timestamp int64  `json:"timestamp"`
	Command   string `json:"command"`
	Title     string `json:"title"`
}

// Event represents a single frame in the asciinema format.
type Event struct {
	Time float64 `json:"time"`
	Type string  `json:"type"`
	Data string  `json:"data"`
}

// NewSessionRecorder initialises a recorder in the provided output directory.
// The directory is created when it does not already exist.
func NewSessionRecorder(outDir string) (*SessionRecorder, error) {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, fmt.Errorf("create recorder directory: %w", err)
	}

	return &SessionRecorder{
		outDir:  outDir,
		width:   80,
		height:  24,
		started: time.Now(),
	}, nil
}

// RecordSession persists the provided command invocation as an asciinema
// recording. The caller supplies stdout and stderr so the recorder remains pure
// and test friendly. The returned path is relative to the repository root so it
// can be published directly as an artifact.
func (r *SessionRecorder) RecordSession(name string, command []string, stdout, stderr string, exitCode int, started time.Time, duration time.Duration) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	safeName := strings.ToLower(name)
	safeName = strings.ReplaceAll(safeName, " ", "-")
	safeName = strings.ReplaceAll(safeName, "_", "-")
	safeName = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return r
		case r >= '0' && r <= '9':
			return r
		case r == '-':
			return r
		default:
			return -1
		}
	}, safeName)
	if safeName == "" {
		safeName = "session"
	}

	filename := fmt.Sprintf("%s-%d.cast", safeName, started.Unix())
	outputPath := filepath.Join(r.outDir, filename)

	f, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("create recording: %w", err)
	}
	defer f.Close()

	metadata := Metadata{
		Version:   2,
		Width:     r.width,
		Height:    r.height,
		Timestamp: started.Unix(),
		Command:   strings.Join(command, " "),
		Title:     name,
	}

	header, err := json.Marshal(metadata)
	if err != nil {
		return "", fmt.Errorf("encode metadata: %w", err)
	}

	if _, err := fmt.Fprintf(f, "%s\n", header); err != nil {
		return "", fmt.Errorf("write metadata: %w", err)
	}

	events := []Event{
		{
			Time: 0,
			Type: "o",
			Data: stdout,
		},
	}

	if stderr != "" {
		events = append(events, Event{
			Time: duration.Seconds(),
			Type: "o",
			Data: fmt.Sprintf("[stderr]%s", stderr),
		})
	}

	events = append(events, Event{
		Time: duration.Seconds(),
		Type: "o",
		Data: fmt.Sprintf("[exit %d]", exitCode),
	})

	encoder := json.NewEncoder(f)
	for _, event := range events {
		if err := encoder.Encode([]interface{}{event.Time, event.Type, event.Data}); err != nil {
			return "", fmt.Errorf("write event: %w", err)
		}
	}

	return outputPath, nil
}
