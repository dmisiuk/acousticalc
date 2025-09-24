package visual

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ArtifactGeneratorInterface defines the contract for artifact generation
type ArtifactGeneratorInterface interface {
	GenerateComprehensiveArtifacts(logger *VisualTestLogger) error
	CreateDirectoryStructure() error
	GenerateMetadataFiles() error
	GetArtifactSummary() ArtifactSummary
}

// ArtifactGenerator handles creation and management of visual test artifacts
type ArtifactGenerator struct {
	TestName           string
	OutputBaseDir      string
	Metadata           ArtifactMetadata
	Screenshots        []ScreenshotInfo
	Reports            []ReportInfo
	mu                 sync.RWMutex // Thread safety
	ctx                context.Context
	cancel             context.CancelFunc
	performanceMonitor *PerformanceMonitor
}

// ArtifactSummary provides a summary of generated artifacts
type ArtifactSummary struct {
	TotalArtifacts     int                        `json:"total_artifacts"`
	TotalSizeBytes     int64                      `json:"total_size_bytes"`
	GenerationTime     time.Duration              `json:"generation_time"`
	ArtifactTypes      []string                   `json:"artifact_types"`
	PerformanceMetrics map[string]OperationMetric `json:"performance_metrics"`
}

// ArtifactMetadata contains metadata about generated artifacts
type ArtifactMetadata struct {
	TestName     string    `json:"test_name"`
	Platform     string    `json:"platform"`
	Architecture string    `json:"architecture"`
	Timestamp    time.Time `json:"timestamp"`
	Version      string    `json:"version"`
	CoverageInfo string    `json:"coverage_info"`
	Screenshots  int       `json:"screenshot_count"`
	Reports      int       `json:"report_count"`
	TotalSize    int64     `json:"total_size_bytes"`
}

// ScreenshotInfo contains information about individual screenshots
type ScreenshotInfo struct {
	Filename   string            `json:"filename"`
	EventType  string            `json:"event_type"`
	Timestamp  time.Time         `json:"timestamp"`
	Size       int64             `json:"size_bytes"`
	Dimensions string            `json:"dimensions"`
	Metadata   map[string]string `json:"metadata"`
}

// ReportInfo contains information about generated reports
type ReportInfo struct {
	Filename    string    `json:"filename"`
	Type        string    `json:"type"` // "visual_report", "demo_storyboard", "coverage_chart"
	Timestamp   time.Time `json:"timestamp"`
	Size        int64     `json:"size_bytes"`
	Description string    `json:"description"`
}

// NewArtifactGenerator creates a new artifact generator with context
func NewArtifactGenerator(testName, outputBaseDir string) *ArtifactGenerator {
	ctx, cancel := context.WithCancel(context.Background())
	return &ArtifactGenerator{
		TestName:      testName,
		OutputBaseDir: outputBaseDir,
		Metadata: ArtifactMetadata{
			TestName:     testName,
			Platform:     "cross-platform",
			Architecture: "universal",
			Timestamp:    time.Now(),
			Version:      "0.2.2",
		},
		Screenshots:        make([]ScreenshotInfo, 0),
		Reports:            make([]ReportInfo, 0),
		ctx:                ctx,
		cancel:             cancel,
		performanceMonitor: NewPerformanceMonitor(ctx),
	}
}

// NewArtifactGeneratorWithContext creates a new artifact generator with specific context
func NewArtifactGeneratorWithContext(ctx context.Context, testName, outputBaseDir string) *ArtifactGenerator {
	childCtx, cancel := context.WithCancel(ctx)
	return &ArtifactGenerator{
		TestName:      testName,
		OutputBaseDir: outputBaseDir,
		Metadata: ArtifactMetadata{
			TestName:     testName,
			Platform:     "cross-platform",
			Architecture: "universal",
			Timestamp:    time.Now(),
			Version:      "0.2.2",
		},
		Screenshots:        make([]ScreenshotInfo, 0),
		Reports:            make([]ReportInfo, 0),
		ctx:                childCtx,
		cancel:             cancel,
		performanceMonitor: NewPerformanceMonitor(childCtx),
	}
}

// GenerateComprehensiveArtifacts creates all visual testing artifacts
func (ag *ArtifactGenerator) GenerateComprehensiveArtifacts(logger *VisualTestLogger) error {
	// Create all required directories
	if err := ag.createDirectoryStructure(); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// Generate visual report
	if err := ag.generateEnhancedVisualReport(logger); err != nil {
		return fmt.Errorf("failed to generate visual report: %w", err)
	}

	// Generate demo storyboard
	if err := ag.generateProfessionalStoryboard(logger); err != nil {
		return fmt.Errorf("failed to generate demo storyboard: %w", err)
	}

	// Generate test timeline
	if err := ag.generateTestTimeline(logger); err != nil {
		return fmt.Errorf("failed to generate test timeline: %w", err)
	}

	// Generate coverage visualization
	if err := ag.generateCoverageVisualization(); err != nil {
		return fmt.Errorf("failed to generate coverage visualization: %w", err)
	}

	// Generate metadata files
	if err := ag.generateMetadataFiles(); err != nil {
		return fmt.Errorf("failed to generate metadata: %w", err)
	}

	return nil
}

// createDirectoryStructure ensures all required directories exist
func (ag *ArtifactGenerator) createDirectoryStructure() error {
	dirs := []string{
		filepath.Join(ag.OutputBaseDir, "reports"),
		filepath.Join(ag.OutputBaseDir, "demo_content/storyboards"),
		filepath.Join(ag.OutputBaseDir, "demo_content/assets"),
		filepath.Join(ag.OutputBaseDir, "demo_content/metadata"),
		filepath.Join(ag.OutputBaseDir, "charts"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateEnhancedVisualReport creates a comprehensive HTML report
func (ag *ArtifactGenerator) generateEnhancedVisualReport(logger *VisualTestLogger) error {
	reportPath := filepath.Join(ag.OutputBaseDir, "reports", fmt.Sprintf("%s_comprehensive_report.html", ag.TestName))

	html := ag.generateEnhancedHTMLReport(logger)

	if err := os.WriteFile(reportPath, []byte(html), 0644); err != nil {
		return err
	}

	// Track report info
	if info, err := os.Stat(reportPath); err == nil {
		ag.addReport(ReportInfo{
			Filename:    filepath.Base(reportPath),
			Type:        "visual_report",
			Timestamp:   time.Now(),
			Size:        info.Size(),
			Description: "Comprehensive visual test execution report",
		})
	}

	return nil
}

// generateEnhancedHTMLReport creates a professional HTML report
func (ag *ArtifactGenerator) generateEnhancedHTMLReport(logger *VisualTestLogger) string {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AcoustiCalc Visual Testing Report - ` + ag.TestName + `</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'SF Pro Display', -apple-system, BlinkMacSystemFont, sans-serif;
            background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
            color: #ffffff; line-height: 1.6; min-height: 100vh;
        }
        .container { max-width: 1400px; margin: 0 auto; padding: 40px 20px; }
        .header { text-align: center; margin-bottom: 50px; }
        .header h1 { font-size: 3em; font-weight: 300; margin-bottom: 10px; }
        .header h2 { font-size: 1.5em; color: #32cd32; margin-bottom: 20px; }
        .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 20px; margin: 40px 0; }
        .stat-card {
            background: rgba(255,255,255,0.1); backdrop-filter: blur(10px);
            padding: 30px; border-radius: 15px; text-align: center; border: 1px solid rgba(255,255,255,0.2);
        }
        .stat-number { font-size: 2.5em; font-weight: bold; color: #32cd32; }
        .stat-label { font-size: 1.1em; margin-top: 10px; opacity: 0.9; }
        .events-timeline { margin: 50px 0; }
        .timeline-header { font-size: 2em; text-align: center; margin-bottom: 40px; }
        .event-item {
            background: rgba(255,255,255,0.05); margin: 20px 0; padding: 25px;
            border-radius: 12px; border-left: 4px solid #32cd32;
            display: grid; grid-template-columns: 200px 1fr; gap: 30px; align-items: center;
        }
        .event-meta { text-align: center; }
        .event-time { font-size: 1.1em; color: #32cd32; font-weight: bold; }
        .event-type { font-size: 0.9em; opacity: 0.8; margin-top: 5px; }
        .event-content { }
        .event-description { font-size: 1.1em; margin-bottom: 15px; }
        .screenshot { max-width: 400px; border-radius: 8px; box-shadow: 0 10px 30px rgba(0,0,0,0.3); }
        .metadata {
            background: rgba(0,0,0,0.2); padding: 15px; border-radius: 8px;
            margin-top: 15px; font-size: 0.9em; font-family: 'Monaco', monospace;
        }
        .footer { text-align: center; margin-top: 80px; opacity: 0.7; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>AcoustiCalc Visual Testing</h1>
            <h2>` + ag.TestName + `</h2>
            <p>Generated on ` + time.Now().Format("January 2, 2006 at 15:04:05") + `</p>
        </div>

        <div class="stats">
            <div class="stat-card">
                <div class="stat-number">` + fmt.Sprintf("%d", len(logger.Events)) + `</div>
                <div class="stat-label">Visual Events Captured</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">` + fmt.Sprintf("%d", len(logger.Screenshots)) + `</div>
                <div class="stat-label">Screenshots Generated</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">>95%</div>
                <div class="stat-label">Coverage Target</div>
            </div>
            <div class="stat-card">
                <div class="stat-number"><30s</div>
                <div class="stat-label">CI Time Constraint</div>
            </div>
        </div>

        <div class="events-timeline">
            <div class="timeline-header">Test Execution Timeline</div>`

	// Add events to timeline
	for _, event := range logger.Events {
		html += `
            <div class="event-item">
                <div class="event-meta">
                    <div class="event-time">` + event.Timestamp.Format("15:04:05.000") + `</div>
                    <div class="event-type">` + string(event.Type) + `</div>
                </div>
                <div class="event-content">
                    <div class="event-description">` + event.Description + `</div>`

		if event.Screenshot != "" {
			html += `<img src="../screenshots/unit/` + filepath.Base(event.Screenshot) + `" class="screenshot" alt="Screenshot for ` + string(event.Type) + `">`
		}

		if len(event.Metadata) > 0 {
			html += `<div class="metadata"><strong>Event Metadata:</strong><br>`
			for key, value := range event.Metadata {
				html += fmt.Sprintf("%s: %v<br>", key, value)
			}
			html += `</div>`
		}

		html += `</div></div>`
	}

	html += `
        </div>

        <div class="footer">
            <p>Generated by AcoustiCalc Visual Testing Framework v0.2.2</p>
            <p>Cross-Platform Visual Evidence & Demo Content Generation</p>
        </div>
    </div>
</body>
</html>`

	return html
}

// generateProfessionalStoryboard creates a demo-quality storyboard
func (ag *ArtifactGenerator) generateProfessionalStoryboard(logger *VisualTestLogger) error {
	storyboardPath := filepath.Join(ag.OutputBaseDir, "demo_content/storyboards", fmt.Sprintf("%s_professional_storyboard.html", ag.TestName))

	html := ag.generateProfessionalStoryboardHTML(logger)

	if err := os.WriteFile(storyboardPath, []byte(html), 0644); err != nil {
		return err
	}

	// Track report info
	if info, err := os.Stat(storyboardPath); err == nil {
		ag.addReport(ReportInfo{
			Filename:    filepath.Base(storyboardPath),
			Type:        "demo_storyboard",
			Timestamp:   time.Now(),
			Size:        info.Size(),
			Description: "Professional demo storyboard for marketing use",
		})
	}

	return nil
}

// generateProfessionalStoryboardHTML creates a marketing-grade storyboard
func (ag *ArtifactGenerator) generateProfessionalStoryboardHTML(logger *VisualTestLogger) string {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AcoustiCalc Demo Storyboard</title>
    <style>
        body {
            margin: 0; font-family: 'Helvetica Neue', Arial, sans-serif;
            background: linear-gradient(45deg, #FF6B6B, #4ECDC4, #45B7D1, #96CEB4, #FFEAA7);
            background-size: 300% 300%; animation: gradientShift 8s ease infinite;
        }
        @keyframes gradientShift { 0%, 100% { background-position: 0% 50%; } 50% { background-position: 100% 50%; } }
        .hero { height: 100vh; display: flex; align-items: center; justify-content: center; text-align: center; color: white; }
        .hero h1 { font-size: 4em; font-weight: 100; margin-bottom: 20px; text-shadow: 2px 2px 4px rgba(0,0,0,0.3); }
        .hero p { font-size: 1.5em; opacity: 0.9; }
        .storyboard { background: white; padding: 80px 0; }
        .container { max-width: 1200px; margin: 0 auto; padding: 0 40px; }
        .section-title { font-size: 3em; text-align: center; margin-bottom: 60px; color: #2c3e50; }
        .scenes { display: grid; grid-template-columns: repeat(auto-fit, minmax(350px, 1fr)); gap: 40px; }
        .scene {
            background: #f8f9fa; border-radius: 20px; overflow: hidden;
            box-shadow: 0 20px 60px rgba(0,0,0,0.1); transition: transform 0.3s ease;
        }
        .scene:hover { transform: translateY(-10px); }
        .scene-number {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white; padding: 20px; font-size: 1.2em; font-weight: bold; text-align: center;
        }
        .scene-image { width: 100%; height: 250px; object-fit: cover; }
        .scene-info { padding: 30px; }
        .scene-title { font-size: 1.4em; font-weight: bold; margin-bottom: 15px; color: #2c3e50; }
        .scene-description { color: #7f8c8d; line-height: 1.6; }
    </style>
</head>
<body>
    <div class="hero">
        <div>
            <h1>AcoustiCalc</h1>
            <p>Terminal Calculator with Audio Feedback</p>
            <p style="font-size: 1.2em; margin-top: 30px;">Professional Demo Storyboard</p>
        </div>
    </div>

    <div class="storyboard">
        <div class="container">
            <h2 class="section-title">Demo Scenes</h2>
            <div class="scenes">`

	// Add demo scenes from visual events
	for i, event := range logger.Events {
		if event.Screenshot != "" {
			html += fmt.Sprintf(`
                <div class="scene">
                    <div class="scene-number">Scene %d</div>
                    <img src="../../screenshots/unit/%s" class="scene-image" alt="%s">
                    <div class="scene-info">
                        <div class="scene-title">%s</div>
                        <div class="scene-description">%s<br><small>Captured at %s</small></div>
                    </div>
                </div>`,
				i+1,
				filepath.Base(event.Screenshot),
				event.Description,
				string(event.Type),
				event.Description,
				event.Timestamp.Format("15:04:05"))
		}
	}

	html += `
            </div>
        </div>
    </div>
</body>
</html>`

	return html
}

// generateTestTimeline creates a visual timeline of test execution
func (ag *ArtifactGenerator) generateTestTimeline(logger *VisualTestLogger) error {
	timelinePath := filepath.Join(ag.OutputBaseDir, "charts", fmt.Sprintf("%s_timeline.html", ag.TestName))

	html := ag.generateTimelineHTML(logger)

	if err := os.WriteFile(timelinePath, []byte(html), 0644); err != nil {
		return err
	}

	return nil
}

// generateTimelineHTML creates an interactive timeline
func (ag *ArtifactGenerator) generateTimelineHTML(logger *VisualTestLogger) string {
	// Simple timeline visualization
	return fmt.Sprintf(`<!DOCTYPE html>
<html><head><title>Test Timeline</title></head>
<body><h1>Test Timeline for %s</h1>
<p>Test Duration: %v</p>
<p>Total Events: %d</p>
</body></html>`,
		ag.TestName,
		time.Since(logger.StartTime),
		len(logger.Events))
}

// generateCoverageVisualization creates coverage charts
func (ag *ArtifactGenerator) generateCoverageVisualization() error {
	chartPath := filepath.Join(ag.OutputBaseDir, "charts", fmt.Sprintf("%s_coverage.html", ag.TestName))

	html := `<!DOCTYPE html>
<html><head><title>Coverage Visualization</title></head>
<body><h1>Visual Testing Coverage</h1>
<p>Target: >95% test coverage with visual evidence</p>
<p>Status: Implementation in progress</p>
</body></html>`

	return os.WriteFile(chartPath, []byte(html), 0644)
}

// generateMetadataFiles creates JSON metadata for all artifacts
func (ag *ArtifactGenerator) generateMetadataFiles() error {
	metadataPath := filepath.Join(ag.OutputBaseDir, "demo_content/metadata", fmt.Sprintf("%s_metadata.json", ag.TestName))

	// Update metadata totals
	ag.Metadata.Screenshots = len(ag.Screenshots)
	ag.Metadata.Reports = len(ag.Reports)

	metadataJSON, err := json.MarshalIndent(ag.Metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	return os.WriteFile(metadataPath, metadataJSON, 0644)
}

// GetArtifactSummary returns a thread-safe summary of generated artifacts
func (ag *ArtifactGenerator) GetArtifactSummary() ArtifactSummary {
	ag.mu.RLock()
	defer ag.mu.RUnlock()

	var totalSize int64
	artifactTypes := make(map[string]bool)

	for _, report := range ag.Reports {
		totalSize += report.Size
		artifactTypes[report.Type] = true
	}

	for _, screenshot := range ag.Screenshots {
		totalSize += screenshot.Size
		artifactTypes["screenshot"] = true
	}

	typeList := make([]string, 0, len(artifactTypes))
	for artType := range artifactTypes {
		typeList = append(typeList, artType)
	}

	var performanceMetrics map[string]OperationMetric
	if ag.performanceMonitor != nil {
		performanceMetrics = ag.performanceMonitor.GetMetrics()
	}

	return ArtifactSummary{
		TotalArtifacts:     len(ag.Reports) + len(ag.Screenshots),
		TotalSizeBytes:     totalSize,
		GenerationTime:     time.Since(ag.Metadata.Timestamp),
		ArtifactTypes:      typeList,
		PerformanceMetrics: performanceMetrics,
	}
}

// CreateDirectoryStructure ensures all required directories exist (public method)
func (ag *ArtifactGenerator) CreateDirectoryStructure() error {
	return ag.createDirectoryStructure()
}

// GenerateMetadataFiles creates JSON metadata for all artifacts (public method)
func (ag *ArtifactGenerator) GenerateMetadataFiles() error {
	return ag.generateMetadataFiles()
}

// Close gracefully shuts down the artifact generator
func (ag *ArtifactGenerator) Close() error {
	if ag.cancel != nil {
		ag.cancel()
	}
	if ag.performanceMonitor != nil {
		ag.performanceMonitor.Stop()
	}
	return nil
}

// addReport safely adds a report to the collection
func (ag *ArtifactGenerator) addReport(report ReportInfo) {
	ag.mu.Lock()
	defer ag.mu.Unlock()
	ag.Reports = append(ag.Reports, report)
}

// addScreenshot safely adds a screenshot to the collection
func (ag *ArtifactGenerator) addScreenshot(screenshot ScreenshotInfo) {
	ag.mu.Lock()
	defer ag.mu.Unlock()
	ag.Screenshots = append(ag.Screenshots, screenshot)
}
