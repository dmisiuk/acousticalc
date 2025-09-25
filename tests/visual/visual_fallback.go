//go:build !gui || (!linux && !darwin)

package visual

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ScreenshotCapturer defines the interface for screenshot capture implementations
type ScreenshotCapturer interface {
	CaptureScreen(eventType string) (string, error)
	SetOutputDir(dir string)
	SetTestName(name string)
}

// ScreenshotCapture provides fallback implementation for non-GUI environments
type ScreenshotCapture struct {
	OutputDir string
	TestName  string
	Timestamp time.Time
	Format    string
	Quality   int
	capturer  ScreenshotEngine
}

// ScreenshotEngine abstracts the underlying screenshot mechanism
type ScreenshotEngine interface {
	Capture() ([]byte, error)
	GetImageData() (interface{}, error)
	GetPlatform() string
	IsAvailable() bool
}

// MockScreenshotEngine provides fallback for unsupported platforms
type MockScreenshotEngine struct{}

func (m *MockScreenshotEngine) Capture() ([]byte, error) {
	return []byte("mock-screenshot-data"), nil
}

func (m *MockScreenshotEngine) GetImageData() (interface{}, error) {
	return "mock-image-data", nil
}

func (m *MockScreenshotEngine) GetPlatform() string {
	return "mock"
}

func (m *MockScreenshotEngine) IsAvailable() bool {
	return true
}

// RobotGoEngine fallback for non-GUI environments
type RobotGoEngine struct {
	platform string
}

func NewRobotGoEngine() *RobotGoEngine {
	return &RobotGoEngine{
		platform: "fallback",
	}
}

func (rg *RobotGoEngine) Capture() ([]byte, error) {
	return []byte("fallback-screenshot-data"), nil
}

func (rg *RobotGoEngine) GetImageData() (interface{}, error) {
	return "fallback-image-data", nil
}

func (rg *RobotGoEngine) GetPlatform() string {
	return rg.platform
}

func (rg *RobotGoEngine) IsAvailable() bool {
	return false // Indicate GUI not available
}

// NewScreenshotCapture creates a new screenshot capture instance with fallback behavior
func NewScreenshotCapture(testName, outputDir string) *ScreenshotCapture {
	return &ScreenshotCapture{
		OutputDir: outputDir,
		TestName:  testName,
		Timestamp: time.Now(),
		Format:    "png",
		Quality:   100,
		capturer:  &MockScreenshotEngine{},
	}
}

// SetOutputDir implements ScreenshotCapturer interface
func (sc *ScreenshotCapture) SetOutputDir(dir string) {
	sc.OutputDir = dir
}

// SetTestName implements ScreenshotCapturer interface  
func (sc *ScreenshotCapture) SetTestName(name string) {
	sc.TestName = name
}

// CaptureScreen provides fallback implementation that creates a placeholder file
func (sc *ScreenshotCapture) CaptureScreen(eventType string) (string, error) {
	// Ensure output directory exists
	if err := os.MkdirAll(sc.OutputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate filename with timestamp and event type
	filename := fmt.Sprintf("%s_%s_%s.txt",
		sc.TestName,
		eventType,
		sc.Timestamp.Format("20060102_150405"))

	filePath := filepath.Join(sc.OutputDir, filename)

	// Create a placeholder file indicating fallback mode
	content := fmt.Sprintf("Screenshot placeholder for %s event\nTest: %s\nTime: %s\nNote: GUI not available, screenshot capture disabled\n",
		eventType, sc.TestName, time.Now().Format("2006-01-02 15:04:05"))

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to create placeholder file: %w", err)
	}

	return filePath, nil
}

// CaptureTestEvent captures screenshot for specific test events (fallback)
func (sc *ScreenshotCapture) CaptureTestEvent(t interface{}, eventType string) string {
	// For tests that expect actual screenshot files, return empty string to trigger skip
	// This allows tests to handle the "no GUI" case gracefully
	return ""
}

// VisualTestEvent represents different test events for screenshot capture
type VisualTestEvent string

const (
	EventTestStart    VisualTestEvent = "start"
	EventTestPass     VisualTestEvent = "pass"
	EventTestFail     VisualTestEvent = "fail"
	EventTestProcess  VisualTestEvent = "process"
	EventTestComplete VisualTestEvent = "complete"
)

// VisualEvent represents a test event with visual context
type VisualEvent struct {
	Type        VisualTestEvent
	Timestamp   time.Time
	Description string
	Screenshot  string
	Metadata    map[string]interface{}
}

// VisualTestObserver interface for event notifications
type VisualTestObserver interface {
	OnEvent(event VisualEvent)
	OnScreenshot(capturedPath string, eventType string)
	OnTestComplete(logger *VisualTestLogger)
}

// VisualTestLogger provides fallback implementation for non-GUI environments
type VisualTestLogger struct {
	TestName    string
	OutputDir   string
	Screenshots []string
	Events      []VisualEvent
	StartTime   time.Time
	observers   []VisualTestObserver
	mu          sync.RWMutex
}

// NewVisualTestLogger creates a new visual test logger with fallback behavior
func NewVisualTestLogger(testName, outputDir string) *VisualTestLogger {
	return &VisualTestLogger{
		TestName:    testName,
		OutputDir:   outputDir,
		Screenshots: make([]string, 0),
		Events:      make([]VisualEvent, 0),
		StartTime:   time.Now(),
		observers:   make([]VisualTestObserver, 0),
	}
}

// AddObserver adds an observer to the visual test logger
func (vtl *VisualTestLogger) AddObserver(observer VisualTestObserver) {
	vtl.mu.Lock()
	defer vtl.mu.Unlock()
	vtl.observers = append(vtl.observers, observer)
}

// RemoveObserver removes an observer from the visual test logger
func (vtl *VisualTestLogger) RemoveObserver(observer VisualTestObserver) {
	vtl.mu.Lock()
	defer vtl.mu.Unlock()
	for i, obs := range vtl.observers {
		if obs == observer {
			vtl.observers = append(vtl.observers[:i], vtl.observers[i+1:]...)
			break
		}
	}
}

// LogEvent logs a visual test event without screenshot capture (fallback)
func (vtl *VisualTestLogger) LogEvent(eventType VisualTestEvent, description string, metadata map[string]interface{}) {
	event := VisualEvent{
		Type:        eventType,
		Timestamp:   time.Now(),
		Description: description,
		Metadata:    metadata,
		Screenshot:  "", // No screenshot in fallback mode
	}

	vtl.Events = append(vtl.Events, event)
	vtl.notifyObservers(event)
	
	// Log to console instead of capturing screenshot
	fmt.Printf("[VISUAL-FALLBACK] %s: %s\n", eventType, description)
}

// notifyObservers notifies all observers of an event
func (vtl *VisualTestLogger) notifyObservers(event VisualEvent) {
	vtl.mu.RLock()
	defer vtl.mu.RUnlock()

	for _, observer := range vtl.observers {
		observer.OnEvent(event)
	}
}

// Complete marks the test as complete and notifies observers
func (vtl *VisualTestLogger) Complete() {
	vtl.mu.RLock()
	defer vtl.mu.RUnlock()

	for _, observer := range vtl.observers {
		observer.OnTestComplete(vtl)
	}
}

// GenerateVisualReport generates a simplified text report (fallback)
func (vtl *VisualTestLogger) GenerateVisualReport() error {
	if err := os.MkdirAll(vtl.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	reportPath := filepath.Join(vtl.OutputDir, fmt.Sprintf("%s_visual_report.txt", vtl.TestName))
	
	content := fmt.Sprintf("Visual Test Report (Fallback Mode)\n")
	content += fmt.Sprintf("Test: %s\n", vtl.TestName)
	content += fmt.Sprintf("Start Time: %s\n", vtl.StartTime.Format("2006-01-02 15:04:05"))
	content += fmt.Sprintf("Total Events: %d\n", len(vtl.Events))
	content += fmt.Sprintf("Note: Screenshot capture not available in this environment\n\n")

	for _, event := range vtl.Events {
		content += fmt.Sprintf("%s - %s: %s\n", 
			event.Timestamp.Format("15:04:05.000"), 
			event.Type, 
			event.Description)
	}

	if err := os.WriteFile(reportPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write visual report: %w", err)
	}

	return nil
}

// CreateDemoStoryboard creates a simplified demo content (fallback)
func (vtl *VisualTestLogger) CreateDemoStoryboard() error {
	storyboardDir := filepath.Join(vtl.OutputDir, "../demo_content/storyboards")
	if err := os.MkdirAll(storyboardDir, 0755); err != nil {
		return fmt.Errorf("failed to create storyboard directory: %w", err)
	}

	storyboardPath := filepath.Join(storyboardDir, fmt.Sprintf("%s_storyboard.txt", vtl.TestName))
	
	content := fmt.Sprintf("Demo Storyboard (Fallback Mode)\n")
	content += fmt.Sprintf("Test: %s\n", vtl.TestName)
	content += fmt.Sprintf("Note: Visual storyboard not available without GUI support\n\n")

	for i, event := range vtl.Events {
		content += fmt.Sprintf("Scene %d: %s - %s\n", i+1, event.Type, event.Description)
	}

	if err := os.WriteFile(storyboardPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write demo storyboard: %w", err)
	}

	return nil
}

// PerformanceMonitor provides fallback performance monitoring
type PerformanceMonitor struct {
	mu         sync.RWMutex
	startTime  time.Time
	metrics    map[string]*OperationMetric
	thresholds *PerformanceThresholds
	ctx        context.Context
	cancel     context.CancelFunc
}

// OperationMetric tracks performance data for specific operations
type OperationMetric struct {
	Name         string        `json:"name"`
	Count        int           `json:"count"`
	TotalTime    time.Duration `json:"total_time"`
	AverageTime  time.Duration `json:"average_time"`
	MinTime      time.Duration `json:"min_time"`
	MaxTime      time.Duration `json:"max_time"`
	LastExecuted time.Time     `json:"last_executed"`
}

// PerformanceThresholds defines acceptable performance limits
type PerformanceThresholds struct {
	ScreenshotCapture time.Duration // Default: 5s
	ReportGeneration  time.Duration // Default: 10s
	TotalCIOverhead   time.Duration // Default: 30s
}

// NewPerformanceMonitor creates a new thread-safe performance monitor
func NewPerformanceMonitor(ctx context.Context) *PerformanceMonitor {
	monitorCtx, cancel := context.WithCancel(ctx)
	return &PerformanceMonitor{
		startTime: time.Now(),
		metrics:   make(map[string]*OperationMetric),
		thresholds: &PerformanceThresholds{
			ScreenshotCapture: 5 * time.Second,
			ReportGeneration:  10 * time.Second,
			TotalCIOverhead:   30 * time.Second,
		},
		ctx:    monitorCtx,
		cancel: cancel,
	}
}

// TrackOperation measures and records the performance of an operation
func (pm *PerformanceMonitor) TrackOperation(name string, operation func() error) error {
	start := time.Now()
	err := operation()
	duration := time.Since(start)

	pm.recordMetric(name, duration)
	return err
}

// recordMetric safely records a performance metric
func (pm *PerformanceMonitor) recordMetric(name string, duration time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	metric, exists := pm.metrics[name]
	if !exists {
		metric = &OperationMetric{
			Name:    name,
			MinTime: duration,
			MaxTime: duration,
		}
		pm.metrics[name] = metric
	}

	metric.Count++
	metric.TotalTime += duration
	metric.AverageTime = metric.TotalTime / time.Duration(metric.Count)
	metric.LastExecuted = time.Now()

	if duration < metric.MinTime {
		metric.MinTime = duration
	}
	if duration > metric.MaxTime {
		metric.MaxTime = duration
	}
}

// GetMetrics returns a thread-safe copy of all metrics
func (pm *PerformanceMonitor) GetMetrics() map[string]OperationMetric {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	result := make(map[string]OperationMetric)
	for name, metric := range pm.metrics {
		result[name] = *metric // Copy the metric
	}
	return result
}

// CheckThresholds validates current performance against defined thresholds
func (pm *PerformanceMonitor) CheckThresholds() []string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var violations []string

	for name, metric := range pm.metrics {
		var threshold time.Duration
		switch name {
		case "screenshot_capture":
			threshold = pm.thresholds.ScreenshotCapture
		case "report_generation":
			threshold = pm.thresholds.ReportGeneration
		case "total_ci":
			threshold = pm.thresholds.TotalCIOverhead
		default:
			continue // Skip unknown metrics
		}

		if metric.AverageTime > threshold {
			violations = append(violations, fmt.Sprintf(
				"%s: %v average exceeds threshold %v",
				name, metric.AverageTime, threshold))
		}
	}

	return violations
}

// Stop gracefully shuts down the performance monitor
func (pm *PerformanceMonitor) Stop() {
	if pm.cancel != nil {
		pm.cancel()
	}
}

// generateHTMLReport creates an HTML report with visual elements (fallback)
func (vtl *VisualTestLogger) generateHTMLReport() string {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Visual Test Report (Fallback): %s</title>
    <style>
        body { font-family: monospace; margin: 40px; background: #f0f0f0; color: #333; }
        .header { border-bottom: 2px solid #666; padding: 20px 0; margin-bottom: 30px; }
        .event { margin: 20px 0; padding: 15px; border-left: 3px solid #666; background: #fff; }
        .timestamp { color: #666; font-size: 12px; }
        .event-type { color: #333; font-weight: bold; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Visual Test Report (Fallback Mode)</h1>
        <h2>Test: %s</h2>
        <p>Start Time: %s</p>
        <p>Total Events: %d</p>
        <p>Note: Screenshot capture not available in this environment</p>
    </div>`, vtl.TestName, vtl.TestName, vtl.StartTime.Format("2006-01-02 15:04:05"), len(vtl.Events))

	for _, event := range vtl.Events {
		html += fmt.Sprintf(`
    <div class="event">
        <div class="timestamp">%s</div>
        <div class="event-type">%s</div>
        <div>%s</div>
    </div>`, event.Timestamp.Format("15:04:05.000"), event.Type, event.Description)
	}

	html += `
</body>
</html>`

	return html
}

// generateDemoStoryboard creates a demo-focused visual storyboard (fallback)
func (vtl *VisualTestLogger) generateDemoStoryboard() string {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Demo Storyboard (Fallback): %s</title>
    <style>
        body { font-family: sans-serif; margin: 0; background: #e0e0e0; }
        .container { max-width: 1200px; margin: 0 auto; padding: 40px 20px; }
        .header { text-align: center; color: #333; margin-bottom: 50px; }
        .storyboard { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 30px; }
        .scene { background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 8px rgba(0,0,0,0.1); }
        .scene-header { background: #666; color: white; padding: 15px; font-weight: bold; }
        .scene-description { padding: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>AcoustiCalc Demo Storyboard (Fallback)</h1>
            <h2>%s</h2>
            <p>Visual content not available without GUI support</p>
        </div>
        <div class="storyboard">`, vtl.TestName, vtl.TestName)

	for i, event := range vtl.Events {
		html += fmt.Sprintf(`
            <div class="scene">
                <div class="scene-header">Scene %d: %s</div>
                <div class="scene-description">
                    <strong>Action:</strong> %s<br>
                    <strong>Time:</strong> %s<br>
                    <strong>Note:</strong> Screenshot not available
                </div>
            </div>`, i+1, event.Type, event.Description, event.Timestamp.Format("15:04:05"))
	}

	html += `
        </div>
    </div>
</body>
</html>`

	return html
}

// OptimizeScreenshots provides fallback for screenshot optimization
func OptimizeScreenshots(inputDir, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create optimization output directory: %w", err)
	}

	// In fallback mode, just copy text placeholders
	return filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".txt" {
			// Copy placeholder files
			outputPath := filepath.Join(outputDir, info.Name())
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			
			optimizedContent := fmt.Sprintf("Optimized: %s\n", string(content))
			return os.WriteFile(outputPath, []byte(optimizedContent), 0644)
		}

		return nil
	})
}