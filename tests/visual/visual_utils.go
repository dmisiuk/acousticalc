package visual

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"github.com/go-vgo/robotgo"
)

// ScreenshotCapture handles cross-platform screenshot capture
type ScreenshotCapture struct {
	OutputDir    string
	TestName     string
	Timestamp    time.Time
	Format       string // "png" (default)
	Quality      int    // 100 (lossless for PNG)
}

// NewScreenshotCapture creates a new screenshot capture instance
func NewScreenshotCapture(testName, outputDir string) *ScreenshotCapture {
	return &ScreenshotCapture{
		OutputDir: outputDir,
		TestName:  testName,
		Timestamp: time.Now(),
		Format:    "png",
		Quality:   100,
	}
}

// CaptureScreen captures a screenshot using cross-platform robotgo
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

	filepath := filepath.Join(sc.OutputDir, filename)

	// Capture screenshot using robotgo
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
	err := imaging.Save(img, filepath)
	if err != nil {
		return "", fmt.Errorf("failed to save screenshot: %w", err)
	}

	return filepath, nil
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
	}

	vtl.Events = append(vtl.Events, event)
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