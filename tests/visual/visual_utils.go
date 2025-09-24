package visual

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/go-vgo/robotgo"
)

// ScreenshotCapturer defines the interface for screenshot capture implementations
type ScreenshotCapturer interface {
	CaptureScreen(eventType string) (string, error)
	SetOutputDir(dir string)
	SetTestName(name string)
}

// ScreenshotCapture handles cross-platform screenshot capture
type ScreenshotCapture struct {
	OutputDir string
	TestName  string
	Timestamp time.Time
	Format    string // "png" (default)
	Quality   int    // 100 (lossless for PNG)
	capturer  ScreenshotEngine
}

// ScreenshotEngine abstracts the underlying screenshot mechanism
type ScreenshotEngine interface {
	Capture() ([]byte, error)
	GetImageData() (interface{}, error)
	GetPlatform() string
	IsAvailable() bool
}

// RobotGoEngine implements ScreenshotEngine using robotgo
type RobotGoEngine struct {
	platform string
}

func NewRobotGoEngine() *RobotGoEngine {
	return &RobotGoEngine{
		platform: runtime.GOOS,
	}
}

func (rg *RobotGoEngine) Capture() ([]byte, error) {
	bitmap := robotgo.CaptureScreen()
	if bitmap == nil {
		return nil, fmt.Errorf("failed to capture screen with robotgo")
	}
	// Return bitmap data - simplified for interface
	return []byte{}, nil
}

func (rg *RobotGoEngine) GetImageData() (interface{}, error) {
	bitmap := robotgo.CaptureScreen()
	if bitmap == nil {
		return nil, fmt.Errorf("failed to capture screen")
	}
	return robotgo.ToImage(bitmap), nil
}

func (rg *RobotGoEngine) GetPlatform() string {
	return rg.platform
}

func (rg *RobotGoEngine) IsAvailable() bool {
	// Check if robotgo is available on this platform
	return rg.platform == "darwin" || rg.platform == "linux"
}

// ScreenshotEngineFactory implements Strategy pattern for engine selection
type ScreenshotEngineFactory struct{}

func (sef *ScreenshotEngineFactory) CreateEngine() ScreenshotEngine {
	engine := NewRobotGoEngine()
	if engine.IsAvailable() {
		return engine
	}

	// Fallback to a mock engine for unsupported platforms
	return &MockScreenshotEngine{}
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

// NewScreenshotCapture creates a new screenshot capture instance
func NewScreenshotCapture(testName, outputDir string) *ScreenshotCapture {
	factory := &ScreenshotEngineFactory{}
	return &ScreenshotCapture{
		OutputDir: outputDir,
		TestName:  testName,
		Timestamp: time.Now(),
		Format:    "png",
		Quality:   100,
		capturer:  factory.CreateEngine(),
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

// CaptureScreen captures a screenshot using the configured engine
func (sc *ScreenshotCapture) CaptureScreen(eventType string) (string, error) {
	// Ensure output directory exists
	if err := os.MkdirAll(sc.OutputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate filename with timestamp and event type
	filename := fmt.Sprintf("%s_%s_%s.png",
		sc.TestName,
		eventType,
		sc.Timestamp.Format("20060102_150405"))

	filePath := filepath.Join(sc.OutputDir, filename)

	// Capture screenshot using robotgo (maintaining proven functionality)
	// Future: This can be abstracted through the capturer interface when needed
	bitmap := robotgo.CaptureScreen()
	if bitmap == nil {
		return "", fmt.Errorf("failed to capture screen")
	}

	// Convert robotgo bitmap to standard image
	img := robotgo.ToImage(bitmap)
	if img == nil {
		return "", fmt.Errorf("failed to convert bitmap to image")
	}

	// Save with PNG format for lossless compression
	if err := imaging.Save(img, filePath); err != nil {
		return "", fmt.Errorf("failed to save screenshot: %w", err)
	}

	return filePath, nil
}

// CaptureTestEvent captures screenshot for specific test events
func (sc *ScreenshotCapture) CaptureTestEvent(t interface{}, eventType string) string {
	// Note: Using interface{} instead of *testing.T to avoid import cycle

	filepath, err := sc.CaptureScreen(eventType)
	if err != nil {
		// In real testing scenarios, this would log to t.Logf
		fmt.Printf("Warning: Screenshot capture failed for %s: %v\n", eventType, err)
		return ""
	}

	fmt.Printf("Screenshot captured: %s\n", filepath)
	return filepath
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

// VisualTestLogger handles visual logging and artifact creation
type VisualTestLogger struct {
	TestName    string
	OutputDir   string
	Screenshots []string
	Events      []VisualEvent
	StartTime   time.Time
	observers   []VisualTestObserver
	mu          sync.RWMutex
}

// VisualTestObserver interface for event notifications
type VisualTestObserver interface {
	OnEvent(event VisualEvent)
	OnScreenshot(capturedPath string, eventType string)
	OnTestComplete(logger *VisualTestLogger)
}

// ScreenshotCaptureObserver implements observer for screenshot events
type ScreenshotCaptureObserver struct {
	callback func(event VisualEvent)
}

func (sco *ScreenshotCaptureObserver) OnEvent(event VisualEvent) {
	if sco.callback != nil {
		sco.callback(event)
	}
}

func (sco *ScreenshotCaptureObserver) OnScreenshot(capturedPath string, eventType string) {
	// Default implementation - can be overridden
	fmt.Printf("Screenshot captured: %s for event: %s\n", capturedPath, eventType)
}

func (sco *ScreenshotCaptureObserver) OnTestComplete(logger *VisualTestLogger) {
	// Default implementation
	duration := time.Since(logger.StartTime)
	fmt.Printf("Test completed: %s, Duration: %v, Screenshots: %d\n",
		logger.TestName, duration, len(logger.Screenshots))
}

// VisualEvent represents a test event with visual context
type VisualEvent struct {
	Type        VisualTestEvent
	Timestamp   time.Time
	Description string
	Screenshot  string
	Metadata    map[string]interface{}
}

// NewVisualTestLogger creates a new visual test logger
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

// notifyObservers notifies all observers of an event
func (vtl *VisualTestLogger) notifyObservers(event VisualEvent) {
	vtl.mu.RLock()
	defer vtl.mu.RUnlock()

	for _, observer := range vtl.observers {
		observer.OnEvent(event)
	}
}

// notifyScreenshotObservers notifies all observers of a screenshot capture
func (vtl *VisualTestLogger) notifyScreenshotObservers(capturedPath string, eventType string) {
	vtl.mu.RLock()
	defer vtl.mu.RUnlock()

	for _, observer := range vtl.observers {
		observer.OnScreenshot(capturedPath, eventType)
	}
}

// notifyCompletionObservers notifies all observers of test completion
func (vtl *VisualTestLogger) notifyCompletionObservers() {
	vtl.mu.RLock()
	defer vtl.mu.RUnlock()

	for _, observer := range vtl.observers {
		observer.OnTestComplete(vtl)
	}
}

// LogEvent logs a visual test event with optional screenshot
func (vtl *VisualTestLogger) LogEvent(eventType VisualTestEvent, description string, metadata map[string]interface{}) {
	event := VisualEvent{
		Type:        eventType,
		Timestamp:   time.Now(),
		Description: description,
		Metadata:    metadata,
	}

	// Capture screenshot for the event
	capture := NewScreenshotCapture(vtl.TestName, vtl.OutputDir)
	if screenshot, err := capture.CaptureScreen(string(eventType)); err == nil {
		event.Screenshot = screenshot
		vtl.Screenshots = append(vtl.Screenshots, screenshot)
		vtl.notifyScreenshotObservers(screenshot, string(eventType))
	}

	vtl.Events = append(vtl.Events, event)
	vtl.notifyObservers(event)
}

// GenerateVisualReport generates a visual test execution report
func (vtl *VisualTestLogger) GenerateVisualReport() error {
	reportPath := filepath.Join(vtl.OutputDir, fmt.Sprintf("%s_visual_report.html", vtl.TestName))

	// Create HTML report with visual timeline
	htmlContent := vtl.generateHTMLReport()

	if err := os.WriteFile(reportPath, []byte(htmlContent), 0644); err != nil {
		return fmt.Errorf("failed to write visual report: %w", err)
	}

	return nil
}

// generateHTMLReport creates an HTML report with visual elements
func (vtl *VisualTestLogger) generateHTMLReport() string {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Visual Test Report: ` + vtl.TestName + `</title>
    <style>
        body { font-family: 'Monaco', 'Menlo', monospace; margin: 40px; background: #1e1e1e; color: #d4d4d4; }
        .header { border-bottom: 2px solid #32cd32; padding: 20px 0; margin-bottom: 30px; }
        .event { margin: 20px 0; padding: 15px; border-left: 3px solid #32cd32; background: #2d2d2d; }
        .screenshot { max-width: 300px; border: 1px solid #444; margin: 10px 0; }
        .timestamp { color: #808080; font-size: 12px; }
        .event-type { color: #32cd32; font-weight: bold; }
        .metadata { background: #3c3c3c; padding: 10px; margin: 10px 0; font-size: 12px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Visual Test Report</h1>
        <h2>Test: ` + vtl.TestName + `</h2>
        <p>Start Time: ` + vtl.StartTime.Format("2006-01-02 15:04:05") + `</p>
        <p>Total Events: ` + fmt.Sprintf("%d", len(vtl.Events)) + `</p>
        <p>Screenshots Captured: ` + fmt.Sprintf("%d", len(vtl.Screenshots)) + `</p>
    </div>`

	for _, event := range vtl.Events {
		html += `
    <div class="event">
        <div class="timestamp">` + event.Timestamp.Format("15:04:05.000") + `</div>
        <div class="event-type">` + string(event.Type) + `</div>
        <div>` + event.Description + `</div>`

		if event.Screenshot != "" {
			html += `<img src="` + filepath.Base(event.Screenshot) + `" class="screenshot" alt="Screenshot for ` + string(event.Type) + `">`
		}

		if len(event.Metadata) > 0 {
			html += `<div class="metadata"><strong>Metadata:</strong><br>`
			for key, value := range event.Metadata {
				html += fmt.Sprintf("%s: %v<br>", key, value)
			}
			html += `</div>`
		}

		html += `</div>`
	}

	html += `
</body>
</html>`

	return html
}

// CreateDemoStoryboard creates a visual storyboard for demo content
func (vtl *VisualTestLogger) CreateDemoStoryboard() error {
	storyboardDir := filepath.Join(vtl.OutputDir, "../demo_content/storyboards")
	if err := os.MkdirAll(storyboardDir, 0755); err != nil {
		return fmt.Errorf("failed to create storyboard directory: %w", err)
	}

	storyboardPath := filepath.Join(storyboardDir, fmt.Sprintf("%s_storyboard.html", vtl.TestName))

	// Create demo storyboard HTML
	storyboard := vtl.generateDemoStoryboard()

	if err := os.WriteFile(storyboardPath, []byte(storyboard), 0644); err != nil {
		return fmt.Errorf("failed to write demo storyboard: %w", err)
	}

	return nil
}

// generateDemoStoryboard creates a demo-focused visual storyboard
func (vtl *VisualTestLogger) generateDemoStoryboard() string {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Demo Storyboard: ` + vtl.TestName + `</title>
    <style>
        body { font-family: 'San Francisco', -apple-system, sans-serif; margin: 0; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
        .container { max-width: 1200px; margin: 0 auto; padding: 40px 20px; }
        .header { text-align: center; color: white; margin-bottom: 50px; }
        .storyboard { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 30px; }
        .scene { background: white; border-radius: 12px; overflow: hidden; box-shadow: 0 10px 30px rgba(0,0,0,0.3); }
        .scene-header { background: #32cd32; color: white; padding: 15px; font-weight: bold; }
        .scene-image { width: 100%; height: 200px; object-fit: cover; }
        .scene-description { padding: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>AcoustiCalc Demo Storyboard</h1>
            <h2>` + vtl.TestName + `</h2>
            <p>Professional Demo Content Generation</p>
        </div>
        <div class="storyboard">`

	for i, event := range vtl.Events {
		if event.Screenshot != "" {
			html += `
            <div class="scene">
                <div class="scene-header">Scene ` + fmt.Sprintf("%d", i+1) + `: ` + string(event.Type) + `</div>
                <img src="../../../screenshots/unit/` + filepath.Base(event.Screenshot) + `" class="scene-image" alt="` + event.Description + `">
                <div class="scene-description">
                    <strong>Action:</strong> ` + event.Description + `<br>
                    <strong>Time:</strong> ` + event.Timestamp.Format("15:04:05") + `
                </div>
            </div>`
		}
	}

	html += `
        </div>
    </div>
</body>
</html>`

	return html
}

// Complete marks the test as complete and notifies observers
func (vtl *VisualTestLogger) Complete() {
	vtl.notifyCompletionObservers()
}

// PerformanceMonitor provides thread-safe performance tracking for visual tests
type PerformanceMonitor struct {
	mu           sync.RWMutex
	startTime    time.Time
	metrics      map[string]*OperationMetric
	thresholds   *PerformanceThresholds
	ctx          context.Context
	cancel       context.CancelFunc
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

// OptimizeScreenshots optimizes screenshots for demo use
func OptimizeScreenshots(inputDir, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create optimization output directory: %w", err)
	}

	return filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".png" {
			// Load image
			img, err := imaging.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open image %s: %w", path, err)
			}

			// Optimize: resize if too large, maintain aspect ratio
			const maxWidth = 1920
			const maxHeight = 1080

			bounds := img.Bounds()
			if bounds.Max.X > maxWidth || bounds.Max.Y > maxHeight {
				img = imaging.Fit(img, maxWidth, maxHeight, imaging.Lanczos)
			}

			// Save optimized version
			outputPath := filepath.Join(outputDir, info.Name())
			if err := imaging.Save(img, outputPath); err != nil {
				return fmt.Errorf("failed to save optimized image %s: %w", outputPath, err)
			}
		}

		return nil
	})
}
