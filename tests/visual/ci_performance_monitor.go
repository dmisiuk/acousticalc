package visual

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// CIPerformanceMonitor tracks visual testing performance in CI environments
type CIPerformanceMonitor struct {
	StartTime       time.Time            `json:"start_time"`
	EndTime         time.Time            `json:"end_time"`
	TotalDuration   time.Duration        `json:"total_duration"`
	Platform        string               `json:"platform"`
	ScreenshotCount int                  `json:"screenshot_count"`
	ArtifactCount   int                  `json:"artifact_count"`
	Metrics         map[string]float64   `json:"metrics"`
	Thresholds      PerformanceThreshold `json:"thresholds"`
	IsCI            bool                 `json:"is_ci"`
}

// PerformanceThreshold defines acceptable performance limits
type PerformanceThreshold struct {
	MaxCIOverhead     time.Duration `json:"max_ci_overhead"`     // 30 seconds max
	MaxScreenshotTime time.Duration `json:"max_screenshot_time"` // 5 seconds max per screenshot
	MaxArtifactTime   time.Duration `json:"max_artifact_time"`   // 10 seconds max for artifact generation
}

// NewCIPerformanceMonitor creates a new performance monitor
func NewCIPerformanceMonitor() *CIPerformanceMonitor {
	return &CIPerformanceMonitor{
		Platform: fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Metrics:  make(map[string]float64),
		Thresholds: PerformanceThreshold{
			MaxCIOverhead:     30 * time.Second,
			MaxScreenshotTime: 5 * time.Second,
			MaxArtifactTime:   10 * time.Second,
		},
		IsCI: isCI(),
	}
}

// Start begins performance monitoring
func (m *CIPerformanceMonitor) Start() {
	m.StartTime = time.Now()
}

// RecordScreenshot records screenshot capture performance
func (m *CIPerformanceMonitor) RecordScreenshot(duration time.Duration) {
	m.ScreenshotCount++
	m.Metrics[fmt.Sprintf("screenshot_%d_ms", m.ScreenshotCount)] = float64(duration.Milliseconds())
}

// RecordArtifact records artifact generation performance
func (m *CIPerformanceMonitor) RecordArtifact(artifactType string, duration time.Duration) {
	m.ArtifactCount++
	m.Metrics[fmt.Sprintf("artifact_%s_ms", artifactType)] = float64(duration.Milliseconds())
}

// Finish completes monitoring and validates thresholds
func (m *CIPerformanceMonitor) Finish() error {
	m.EndTime = time.Now()
	m.TotalDuration = m.EndTime.Sub(m.StartTime)

	// Validate against thresholds
	if m.TotalDuration > m.Thresholds.MaxCIOverhead {
		return fmt.Errorf("CI overhead %v exceeds threshold %v",
			m.TotalDuration, m.Thresholds.MaxCIOverhead)
	}

	// Check individual screenshot times
	for i := 1; i <= m.ScreenshotCount; i++ {
		key := fmt.Sprintf("screenshot_%d_ms", i)
		if ms, exists := m.Metrics[key]; exists {
			duration := time.Duration(ms) * time.Millisecond
			if duration > m.Thresholds.MaxScreenshotTime {
				return fmt.Errorf("screenshot %d took %v, exceeds threshold %v",
					i, duration, m.Thresholds.MaxScreenshotTime)
			}
		}
	}

	return nil
}

// SaveReport saves performance report to artifacts
func (m *CIPerformanceMonitor) SaveReport(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	// Replace slashes in platform name to avoid directory issues
	safePlatform := strings.ReplaceAll(m.Platform, "/", "_")
	filename := fmt.Sprintf("ci_performance_%s_%s.json", safePlatform, timestamp)
	filePath := filepath.Join(outputDir, filename)

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal performance data: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write performance report: %w", err)
	}

	// Also create a summary file for quick CI checking
	summaryFile := filepath.Join(outputDir, "ci_performance_summary.txt")
	summary := fmt.Sprintf(`Visual Testing Performance Report
Platform: %s
Total Duration: %v (Threshold: %v)
Screenshot Count: %d
Artifact Count: %d
Status: %s
Timestamp: %s
`,
		m.Platform,
		m.TotalDuration,
		m.Thresholds.MaxCIOverhead,
		m.ScreenshotCount,
		m.ArtifactCount,
		m.getStatus(),
		m.EndTime.Format(time.RFC3339))

	if err := os.WriteFile(summaryFile, []byte(summary), 0644); err != nil {
		return fmt.Errorf("failed to write performance summary: %w", err)
	}

	return nil
}

// getStatus returns PASS/FAIL based on threshold validation
func (m *CIPerformanceMonitor) getStatus() string {
	if m.TotalDuration > m.Thresholds.MaxCIOverhead {
		return "FAIL"
	}
	return "PASS"
}

// isCI detects if running in CI environment
func isCI() bool {
	ciEnvVars := []string{
		"CI",
		"CONTINUOUS_INTEGRATION",
		"GITHUB_ACTIONS",
		"TRAVIS",
		"CIRCLECI",
		"JENKINS_URL",
		"GITLAB_CI",
	}

	for _, envVar := range ciEnvVars {
		if os.Getenv(envVar) != "" {
			return true
		}
	}
	return false
}

// MonitorWithContext runs a function with performance monitoring
func MonitorWithContext(ctx context.Context, fn func() error, outputDir string) error {
	monitor := NewCIPerformanceMonitor()
	monitor.Start()

	// Execute the monitored function
	err := fn()

	// Complete monitoring
	if finishErr := monitor.Finish(); finishErr != nil {
		if err == nil {
			err = finishErr
		} else {
			err = fmt.Errorf("original error: %v; monitoring error: %v", err, finishErr)
		}
	}

	// Save report regardless of errors
	if saveErr := monitor.SaveReport(outputDir); saveErr != nil {
		fmt.Printf("Warning: failed to save performance report: %v\n", saveErr)
	}

	return err
}
