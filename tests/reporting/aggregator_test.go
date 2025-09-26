package reporting

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// ReportConfig holds configuration for test reporting
type ReportConfig struct {
	OutputDir     string
	Format        string // "json", "html", "text"
	IncludeVisuals bool
	IncludeMetrics bool
	Platform      string
	Verbose       bool
}

// DefaultReportConfig returns default configuration for test reporting
func DefaultReportConfig() *ReportConfig {
	return &ReportConfig{
		OutputDir:     "tests/artifacts/reports",
		Format:        "json",
		IncludeVisuals: true,
		IncludeMetrics: true,
		Platform:      runtime.GOOS,
		Verbose:       false,
	}
}

// TestReport holds comprehensive test report data
type TestReport struct {
	Metadata    ReportMetadata    `json:"metadata"`
	Summary     ReportSummary     `json:"summary"`
	Results     []TestResult      `json:"results"`
	Platforms   []PlatformInfo    `json:"platforms"`
	Artifacts   []ArtifactInfo    `json:"artifacts"`
	Metrics     ReportMetrics     `json:"metrics"`
	Visuals     []VisualInfo      `json:"visuals"`
	GeneratedAt time.Time         `json:"generated_at"`
}

// ReportMetadata holds metadata about the test report
type ReportMetadata struct {
	Version     string    `json:"version"`
	Platform    string    `json:"platform"`
	GoVersion   string    `json:"go_version"`
	TestSuite   string    `json:"test_suite"`
	Environment string    `json:"environment"`
	GeneratedAt time.Time `json:"generated_at"`
}

// ReportSummary holds summary statistics
type ReportSummary struct {
	TotalTests    int           `json:"total_tests"`
	PassedTests   int           `json:"passed_tests"`
	FailedTests   int           `json:"failed_tests"`
	SkippedTests  int           `json:"skipped_tests"`
	TotalDuration time.Duration `json:"total_duration"`
	Coverage      float64       `json:"coverage"`
	SuccessRate   float64       `json:"success_rate"`
}

// TestResult holds individual test result data
type TestResult struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Platform    string            `json:"platform"`
	Status      string            `json:"status"`
	Duration    time.Duration     `json:"duration"`
	Error       string            `json:"error,omitempty"`
	Artifacts   []string          `json:"artifacts"`
	Metadata    map[string]string `json:"metadata"`
	Timestamp   time.Time         `json:"timestamp"`
}

// PlatformInfo holds platform-specific information
type PlatformInfo struct {
	Name        string            `json:"name"`
	Architecture string           `json:"architecture"`
	OS          string            `json:"os"`
	GoVersion   string            `json:"go_version"`
	Environment map[string]string `json:"environment"`
	TestCount   int               `json:"test_count"`
	PassRate    float64           `json:"pass_rate"`
}

// ArtifactInfo holds information about test artifacts
type ArtifactInfo struct {
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
	Platform    string    `json:"platform"`
	TestType    string    `json:"test_type"`
}

// ReportMetrics holds performance and quality metrics
type ReportMetrics struct {
	Performance PerformanceMetrics `json:"performance"`
	Quality     QualityMetrics     `json:"quality"`
	Reliability ReliabilityMetrics `json:"reliability"`
}

// PerformanceMetrics holds performance-related metrics
type PerformanceMetrics struct {
	AverageTestDuration time.Duration `json:"average_test_duration"`
	SlowestTest         string        `json:"slowest_test"`
	FastestTest         string        `json:"fastest_test"`
	TotalExecutionTime  time.Duration `json:"total_execution_time"`
	Throughput          float64       `json:"throughput"` // tests per second
}

// QualityMetrics holds quality-related metrics
type QualityMetrics struct {
	CoveragePercentage float64 `json:"coverage_percentage"`
	CodeQualityScore   float64 `json:"code_quality_score"`
	TechnicalDebt      float64 `json:"technical_debt"`
	Maintainability    float64 `json:"maintainability"`
}

// ReliabilityMetrics holds reliability-related metrics
type ReliabilityMetrics struct {
	FlakyTestCount     int     `json:"flaky_test_count"`
	StabilityScore     float64 `json:"stability_score"`
	ErrorRate          float64 `json:"error_rate"`
	RecoveryTime       float64 `json:"recovery_time"`
}

// VisualInfo holds information about visual test artifacts
type VisualInfo struct {
	Name        string    `json:"name"`
	Type        string    `json:"type"` // "screenshot", "recording", "chart"
	Path        string    `json:"path"`
	Platform    string    `json:"platform"`
	TestType    string    `json:"test_type"`
	CreatedAt   time.Time `json:"created_at"`
	Size        int64     `json:"size"`
	Description string    `json:"description"`
}

// ReportAggregator manages test report aggregation
type ReportAggregator struct {
	config *ReportConfig
	report *TestReport
}

// NewReportAggregator creates a new report aggregator
func NewReportAggregator(config *ReportConfig) *ReportAggregator {
	if config == nil {
		config = DefaultReportConfig()
	}
	
	// Ensure output directory exists
	os.MkdirAll(config.OutputDir, 0755)
	
	return &ReportAggregator{
		config: config,
		report: &TestReport{
			Metadata: ReportMetadata{
				Version:     "1.0.0",
				Platform:    config.Platform,
				GoVersion:   runtime.Version(),
				TestSuite:   "E2E Cross-Platform Testing",
				Environment: "CI/CD",
				GeneratedAt: time.Now(),
			},
			Results:   make([]TestResult, 0),
			Platforms: make([]PlatformInfo, 0),
			Artifacts: make([]ArtifactInfo, 0),
			Visuals:   make([]VisualInfo, 0),
			GeneratedAt: time.Now(),
		},
	}
}

// TestReportAggregation tests report aggregation functionality
func TestReportAggregation(t *testing.T) {
	config := DefaultReportConfig()
	config.Verbose = true
	aggregator := NewReportAggregator(config)
	
	t.Run("CreateAggregator", func(t *testing.T) {
		if aggregator == nil {
			t.Error("Report aggregator is nil")
		}
		
		if aggregator.config == nil {
			t.Error("Report aggregator config is nil")
		}
		
		if aggregator.report == nil {
			t.Error("Report aggregator report is nil")
		}
		
		// Verify metadata
		metadata := aggregator.report.Metadata
		if metadata.Version == "" {
			t.Error("Report metadata version is empty")
		}
		
		if metadata.Platform == "" {
			t.Error("Report metadata platform is empty")
		}
		
		if metadata.GoVersion == "" {
			t.Error("Report metadata Go version is empty")
		}
	})
	
	t.Run("AddTestResult", func(t *testing.T) {
		testResult := TestResult{
			Name:      "test_example",
			Type:      "unit",
			Platform:  config.Platform,
			Status:    "passed",
			Duration:  100 * time.Millisecond,
			Artifacts: []string{"test_output.txt"},
			Metadata: map[string]string{
				"test_file": "example_test.go",
				"function":  "TestExample",
			},
			Timestamp: time.Now(),
		}
		
		aggregator.AddTestResult(testResult)
		
		if len(aggregator.report.Results) != 1 {
			t.Errorf("Expected 1 test result, got %d", len(aggregator.report.Results))
		}
		
		addedResult := aggregator.report.Results[0]
		if addedResult.Name != testResult.Name {
			t.Error("Added test result name does not match")
		}
		
		if addedResult.Status != testResult.Status {
			t.Error("Added test result status does not match")
		}
	})
	
	t.Run("AddPlatformInfo", func(t *testing.T) {
		platformInfo := PlatformInfo{
			Name:        config.Platform,
			Architecture: runtime.GOARCH,
			OS:          runtime.GOOS,
			GoVersion:   runtime.Version(),
			Environment: map[string]string{
				"GOPATH": os.Getenv("GOPATH"),
				"GOROOT": os.Getenv("GOROOT"),
			},
			TestCount: 1,
			PassRate:  100.0,
		}
		
		aggregator.AddPlatformInfo(platformInfo)
		
		if len(aggregator.report.Platforms) != 1 {
			t.Errorf("Expected 1 platform info, got %d", len(aggregator.report.Platforms))
		}
		
		addedPlatform := aggregator.report.Platforms[0]
		if addedPlatform.Name != platformInfo.Name {
			t.Error("Added platform info name does not match")
		}
	})
	
	t.Run("AddArtifactInfo", func(t *testing.T) {
		artifactInfo := ArtifactInfo{
			Name:      "test_artifact.txt",
			Type:      "output",
			Path:      filepath.Join(config.OutputDir, "test_artifact.txt"),
			Size:      1024,
			CreatedAt: time.Now(),
			Platform:  config.Platform,
			TestType:  "unit",
		}
		
		aggregator.AddArtifactInfo(artifactInfo)
		
		if len(aggregator.report.Artifacts) != 1 {
			t.Errorf("Expected 1 artifact info, got %d", len(aggregator.report.Artifacts))
		}
		
		addedArtifact := aggregator.report.Artifacts[0]
		if addedArtifact.Name != artifactInfo.Name {
			t.Error("Added artifact info name does not match")
		}
	})
	
	t.Run("AddVisualInfo", func(t *testing.T) {
		visualInfo := VisualInfo{
			Name:        "test_screenshot.png",
			Type:        "screenshot",
			Path:        filepath.Join(config.OutputDir, "test_screenshot.png"),
			Platform:    config.Platform,
			TestType:    "e2e",
			CreatedAt:   time.Now(),
			Size:        2048,
			Description: "Test execution screenshot",
		}
		
		aggregator.AddVisualInfo(visualInfo)
		
		if len(aggregator.report.Visuals) != 1 {
			t.Errorf("Expected 1 visual info, got %d", len(aggregator.report.Visuals))
		}
		
		addedVisual := aggregator.report.Visuals[0]
		if addedVisual.Name != visualInfo.Name {
			t.Error("Added visual info name does not match")
		}
	})
}

// AddTestResult adds a test result to the report
func (ra *ReportAggregator) AddTestResult(result TestResult) {
	ra.report.Results = append(ra.report.Results, result)
}

// AddPlatformInfo adds platform information to the report
func (ra *ReportAggregator) AddPlatformInfo(info PlatformInfo) {
	ra.report.Platforms = append(ra.report.Platforms, info)
}

// AddArtifactInfo adds artifact information to the report
func (ra *ReportAggregator) AddArtifactInfo(info ArtifactInfo) {
	ra.report.Artifacts = append(ra.report.Artifacts, info)
}

// AddVisualInfo adds visual information to the report
func (ra *ReportAggregator) AddVisualInfo(info VisualInfo) {
	ra.report.Visuals = append(ra.report.Visuals, info)
}

// TestReportGeneration tests report generation functionality
func TestReportGeneration(t *testing.T) {
	config := DefaultReportConfig()
	aggregator := NewReportAggregator(config)
	
	// Add some test data
	testResult := TestResult{
		Name:      "test_example",
		Type:      "unit",
		Platform:  config.Platform,
		Status:    "passed",
		Duration:  100 * time.Millisecond,
		Artifacts: []string{"test_output.txt"},
		Metadata: map[string]string{
			"test_file": "example_test.go",
		},
		Timestamp: time.Now(),
	}
	aggregator.AddTestResult(testResult)
	
	platformInfo := PlatformInfo{
		Name:        config.Platform,
		Architecture: runtime.GOARCH,
		OS:          runtime.GOOS,
		GoVersion:   runtime.Version(),
		TestCount:   1,
		PassRate:    100.0,
	}
	aggregator.AddPlatformInfo(platformInfo)
	
	t.Run("GenerateJSONReport", func(t *testing.T) {
		config.Format = "json"
		reportFile := filepath.Join(config.OutputDir, "test_report.json")
		
		err := aggregator.GenerateReport(reportFile)
		if err != nil {
			t.Errorf("Failed to generate JSON report: %v", err)
		}
		
		// Verify report file was created
		if _, err := os.Stat(reportFile); os.IsNotExist(err) {
			t.Error("JSON report file was not created")
		}
		
		// Verify report content
		content, err := os.ReadFile(reportFile)
		if err != nil {
			t.Errorf("Failed to read JSON report: %v", err)
		}
		
		var report TestReport
		err = json.Unmarshal(content, &report)
		if err != nil {
			t.Errorf("Failed to unmarshal JSON report: %v", err)
		}
		
		if report.Metadata.Version == "" {
			t.Error("JSON report metadata version is empty")
		}
		
		if len(report.Results) != 1 {
			t.Errorf("JSON report should have 1 result, got %d", len(report.Results))
		}
		
		// Clean up
		os.Remove(reportFile)
	})
	
	t.Run("GenerateTextReport", func(t *testing.T) {
		config.Format = "text"
		reportFile := filepath.Join(config.OutputDir, "test_report.txt")
		
		err := aggregator.GenerateReport(reportFile)
		if err != nil {
			t.Errorf("Failed to generate text report: %v", err)
		}
		
		// Verify report file was created
		if _, err := os.Stat(reportFile); os.IsNotExist(err) {
			t.Error("Text report file was not created")
		}
		
		// Verify report content
		content, err := os.ReadFile(reportFile)
		if err != nil {
			t.Errorf("Failed to read text report: %v", err)
		}
		
		contentStr := string(content)
		if !strings.Contains(contentStr, "Test Report") {
			t.Error("Text report does not contain expected header")
		}
		
		if !strings.Contains(contentStr, "test_example") {
			t.Error("Text report does not contain test result")
		}
		
		// Clean up
		os.Remove(reportFile)
	})
	
	t.Run("GenerateHTMLReport", func(t *testing.T) {
		config.Format = "html"
		reportFile := filepath.Join(config.OutputDir, "test_report.html")
		
		err := aggregator.GenerateReport(reportFile)
		if err != nil {
			t.Errorf("Failed to generate HTML report: %v", err)
		}
		
		// Verify report file was created
		if _, err := os.Stat(reportFile); os.IsNotExist(err) {
			t.Error("HTML report file was not created")
		}
		
		// Verify report content
		content, err := os.ReadFile(reportFile)
		if err != nil {
			t.Errorf("Failed to read HTML report: %v", err)
		}
		
		contentStr := string(content)
		if !strings.Contains(contentStr, "<html>") {
			t.Error("HTML report does not contain HTML structure")
		}
		
		if !strings.Contains(contentStr, "Test Report") {
			t.Error("HTML report does not contain expected header")
		}
		
		// Clean up
		os.Remove(reportFile)
	})
}

// GenerateReport generates a test report in the specified format
func (ra *ReportAggregator) GenerateReport(filename string) error {
	// Calculate summary
	ra.calculateSummary()
	
	// Calculate metrics
	ra.calculateMetrics()
	
	// Generate report based on format
	switch ra.config.Format {
	case "json":
		return ra.generateJSONReport(filename)
	case "text":
		return ra.generateTextReport(filename)
	case "html":
		return ra.generateHTMLReport(filename)
	default:
		return fmt.Errorf("unsupported report format: %s", ra.config.Format)
	}
}

// calculateSummary calculates summary statistics
func (ra *ReportAggregator) calculateSummary() {
	summary := ReportSummary{}
	
	for _, result := range ra.report.Results {
		summary.TotalTests++
		summary.TotalDuration += result.Duration
		
		switch result.Status {
		case "passed":
			summary.PassedTests++
		case "failed":
			summary.FailedTests++
		case "skipped":
			summary.SkippedTests++
		}
	}
	
	if summary.TotalTests > 0 {
		summary.SuccessRate = float64(summary.PassedTests) / float64(summary.TotalTests) * 100
	}
	
	ra.report.Summary = summary
}

// calculateMetrics calculates performance and quality metrics
func (ra *ReportAggregator) calculateMetrics() {
	metrics := ReportMetrics{}
	
	// Calculate performance metrics
	if len(ra.report.Results) > 0 {
		var totalDuration time.Duration
		var slowestTest, fastestTest string
		var slowestDuration, fastestDuration time.Duration
		
		for _, result := range ra.report.Results {
			totalDuration += result.Duration
			
			if result.Duration > slowestDuration {
				slowestDuration = result.Duration
				slowestTest = result.Name
			}
			
			if fastestDuration == 0 || result.Duration < fastestDuration {
				fastestDuration = result.Duration
				fastestTest = result.Name
			}
		}
		
		metrics.Performance = PerformanceMetrics{
			AverageTestDuration: totalDuration / time.Duration(len(ra.report.Results)),
			SlowestTest:         slowestTest,
			FastestTest:         fastestTest,
			TotalExecutionTime:  totalDuration,
			Throughput:          float64(len(ra.report.Results)) / totalDuration.Seconds(),
		}
	}
	
	// Calculate quality metrics (simplified)
	metrics.Quality = QualityMetrics{
		CoveragePercentage: ra.report.Summary.SuccessRate,
		CodeQualityScore:   ra.report.Summary.SuccessRate,
		TechnicalDebt:      0.0,
		Maintainability:    ra.report.Summary.SuccessRate,
	}
	
	// Calculate reliability metrics (simplified)
	metrics.Reliability = ReliabilityMetrics{
		FlakyTestCount: 0,
		StabilityScore: ra.report.Summary.SuccessRate,
		ErrorRate:      100.0 - ra.report.Summary.SuccessRate,
		RecoveryTime:   0.0,
	}
	
	ra.report.Metrics = metrics
}

// generateJSONReport generates a JSON report
func (ra *ReportAggregator) generateJSONReport(filename string) error {
	data, err := json.MarshalIndent(ra.report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report to JSON: %v", err)
	}
	
	return os.WriteFile(filename, data, 0644)
}

// generateTextReport generates a text report
func (ra *ReportAggregator) generateTextReport(filename string) error {
	report := fmt.Sprintf(`Test Report
Generated at: %s
Platform: %s
Go Version: %s

Summary:
- Total Tests: %d
- Passed: %d
- Failed: %d
- Skipped: %d
- Total Duration: %v
- Success Rate: %.2f%%

Platforms:
`, ra.report.Metadata.GeneratedAt.Format(time.RFC3339),
		ra.report.Metadata.Platform,
		ra.report.Metadata.GoVersion,
		ra.report.Summary.TotalTests,
		ra.report.Summary.PassedTests,
		ra.report.Summary.FailedTests,
		ra.report.Summary.SkippedTests,
		ra.report.Summary.TotalDuration,
		ra.report.Summary.SuccessRate)
	
	for _, platform := range ra.report.Platforms {
		report += fmt.Sprintf("- %s (%s): %d tests, %.2f%% pass rate\n",
			platform.Name, platform.Architecture, platform.TestCount, platform.PassRate)
	}
	
	report += "\nTest Results:\n"
	for _, result := range ra.report.Results {
		report += fmt.Sprintf("- %s (%s): %s (%v)\n",
			result.Name, result.Type, result.Status, result.Duration)
		if result.Error != "" {
			report += fmt.Sprintf("  Error: %s\n", result.Error)
		}
	}
	
	report += "\nArtifacts:\n"
	for _, artifact := range ra.report.Artifacts {
		report += fmt.Sprintf("- %s (%s): %s\n",
			artifact.Name, artifact.Type, artifact.Path)
	}
	
	report += "\nVisuals:\n"
	for _, visual := range ra.report.Visuals {
		report += fmt.Sprintf("- %s (%s): %s\n",
			visual.Name, visual.Type, visual.Path)
	}
	
	return os.WriteFile(filename, []byte(report), 0644)
}

// generateHTMLReport generates an HTML report
func (ra *ReportAggregator) generateHTMLReport(filename string) error {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Test Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background-color: #f0f0f0; padding: 20px; border-radius: 5px; }
        .summary { background-color: #e8f4f8; padding: 15px; margin: 10px 0; border-radius: 5px; }
        .results { margin: 20px 0; }
        .result { padding: 10px; margin: 5px 0; border-left: 4px solid #ccc; }
        .result.passed { border-left-color: #4CAF50; }
        .result.failed { border-left-color: #f44336; }
        .result.skipped { border-left-color: #ff9800; }
        .platforms { display: flex; flex-wrap: wrap; }
        .platform { background-color: #f9f9f9; padding: 10px; margin: 5px; border-radius: 5px; }
        .artifacts { margin: 20px 0; }
        .artifact { padding: 5px; margin: 2px 0; background-color: #f5f5f5; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Test Report</h1>
        <p>Generated at: %s</p>
        <p>Platform: %s</p>
        <p>Go Version: %s</p>
    </div>
    
    <div class="summary">
        <h2>Summary</h2>
        <p>Total Tests: %d</p>
        <p>Passed: %d</p>
        <p>Failed: %d</p>
        <p>Skipped: %d</p>
        <p>Total Duration: %v</p>
        <p>Success Rate: %.2f%%</p>
    </div>
    
    <div class="platforms">
        <h2>Platforms</h2>`, ra.report.Metadata.GeneratedAt.Format(time.RFC3339),
		ra.report.Metadata.Platform,
		ra.report.Metadata.GoVersion,
		ra.report.Summary.TotalTests,
		ra.report.Summary.PassedTests,
		ra.report.Summary.FailedTests,
		ra.report.Summary.SkippedTests,
		ra.report.Summary.TotalDuration,
		ra.report.Summary.SuccessRate)
	
	for _, platform := range ra.report.Platforms {
		html += fmt.Sprintf(`
        <div class="platform">
            <h3>%s</h3>
            <p>Architecture: %s</p>
            <p>Tests: %d</p>
            <p>Pass Rate: %.2f%%</p>
        </div>`, platform.Name, platform.Architecture, platform.TestCount, platform.PassRate)
	}
	
	html += `
    </div>
    
    <div class="results">
        <h2>Test Results</h2>`
	
	for _, result := range ra.report.Results {
		html += fmt.Sprintf(`
        <div class="result %s">
            <h3>%s</h3>
            <p>Type: %s</p>
            <p>Platform: %s</p>
            <p>Duration: %v</p>`, result.Status, result.Name, result.Type, result.Platform, result.Duration)
		
		if result.Error != "" {
			html += fmt.Sprintf(`<p>Error: %s</p>`, result.Error)
		}
		
		html += `</div>`
	}
	
	html += `
    </div>
    
    <div class="artifacts">
        <h2>Artifacts</h2>`
	
	for _, artifact := range ra.report.Artifacts {
		html += fmt.Sprintf(`
        <div class="artifact">
            <strong>%s</strong> (%s) - %s
        </div>`, artifact.Name, artifact.Type, artifact.Path)
	}
	
	html += `
    </div>
    
    <div class="artifacts">
        <h2>Visuals</h2>`
	
	for _, visual := range ra.report.Visuals {
		html += fmt.Sprintf(`
        <div class="artifact">
            <strong>%s</strong> (%s) - %s
        </div>`, visual.Name, visual.Type, visual.Path)
	}
	
	html += `
    </div>
</body>
</html>`
	
	return os.WriteFile(filename, []byte(html), 0644)
}

// TestReportAggregationIntegration tests integration between report components
func TestReportAggregationIntegration(t *testing.T) {
	config := DefaultReportConfig()
	config.Verbose = true
	aggregator := NewReportAggregator(config)
	
	// Add comprehensive test data
	testResults := []TestResult{
		{
			Name:      "unit_test_1",
			Type:      "unit",
			Platform:  config.Platform,
			Status:    "passed",
			Duration:  50 * time.Millisecond,
			Artifacts: []string{"unit_output.txt"},
			Metadata: map[string]string{
				"test_file": "unit_test.go",
			},
			Timestamp: time.Now(),
		},
		{
			Name:      "integration_test_1",
			Type:      "integration",
			Platform:  config.Platform,
			Status:    "passed",
			Duration:  200 * time.Millisecond,
			Artifacts: []string{"integration_output.txt"},
			Metadata: map[string]string{
				"test_file": "integration_test.go",
			},
			Timestamp: time.Now(),
		},
		{
			Name:      "e2e_test_1",
			Type:      "e2e",
			Platform:  config.Platform,
			Status:    "failed",
			Duration:  500 * time.Millisecond,
			Error:     "E2E test failed due to timeout",
			Artifacts: []string{"e2e_output.txt", "e2e_recording.cast"},
			Metadata: map[string]string{
				"test_file": "e2e_test.go",
			},
			Timestamp: time.Now(),
		},
	}
	
	for _, result := range testResults {
		aggregator.AddTestResult(result)
	}
	
	// Add platform information
	platformInfo := PlatformInfo{
		Name:        config.Platform,
		Architecture: runtime.GOARCH,
		OS:          runtime.GOOS,
		GoVersion:   runtime.Version(),
		TestCount:   len(testResults),
		PassRate:    66.67, // 2 out of 3 tests passed
	}
	aggregator.AddPlatformInfo(platformInfo)
	
	// Add artifact information
	artifactInfo := ArtifactInfo{
		Name:      "test_artifact.txt",
		Type:      "output",
		Path:      filepath.Join(config.OutputDir, "test_artifact.txt"),
		Size:      1024,
		CreatedAt: time.Now(),
		Platform:  config.Platform,
		TestType:  "unit",
	}
	aggregator.AddArtifactInfo(artifactInfo)
	
	// Add visual information
	visualInfo := VisualInfo{
		Name:        "test_screenshot.png",
		Type:        "screenshot",
		Path:        filepath.Join(config.OutputDir, "test_screenshot.png"),
		Platform:    config.Platform,
		TestType:    "e2e",
		CreatedAt:   time.Now(),
		Size:        2048,
		Description: "E2E test execution screenshot",
	}
	aggregator.AddVisualInfo(visualInfo)
	
	t.Run("ComprehensiveReportGeneration", func(t *testing.T) {
		// Generate all report formats
		formats := []string{"json", "text", "html"}
		
		for _, format := range formats {
			config.Format = format
			reportFile := filepath.Join(config.OutputDir, fmt.Sprintf("comprehensive_report.%s", format))
			
			err := aggregator.GenerateReport(reportFile)
			if err != nil {
				t.Errorf("Failed to generate %s report: %v", format, err)
			}
			
			// Verify report file was created
			if _, err := os.Stat(reportFile); os.IsNotExist(err) {
				t.Errorf("%s report file was not created", format)
			}
			
			// Clean up
			os.Remove(reportFile)
		}
	})
	
	t.Run("SummaryCalculation", func(t *testing.T) {
		// Verify summary calculation
		aggregator.calculateSummary()
		summary := aggregator.report.Summary
		
		if summary.TotalTests != 3 {
			t.Errorf("Expected 3 total tests, got %d", summary.TotalTests)
		}
		
		if summary.PassedTests != 2 {
			t.Errorf("Expected 2 passed tests, got %d", summary.PassedTests)
		}
		
		if summary.FailedTests != 1 {
			t.Errorf("Expected 1 failed test, got %d", summary.FailedTests)
		}
		
		if summary.SkippedTests != 0 {
			t.Errorf("Expected 0 skipped tests, got %d", summary.SkippedTests)
		}
		
		expectedSuccessRate := 66.67
		if summary.SuccessRate < expectedSuccessRate-0.01 || summary.SuccessRate > expectedSuccessRate+0.01 {
			t.Errorf("Expected success rate around %.2f%%, got %.2f%%", expectedSuccessRate, summary.SuccessRate)
		}
	})
	
	t.Run("MetricsCalculation", func(t *testing.T) {
		// Verify metrics calculation
		aggregator.calculateMetrics()
		metrics := aggregator.report.Metrics
		
		// Check performance metrics
		if metrics.Performance.AverageTestDuration <= 0 {
			t.Error("Average test duration should be positive")
		}
		
		if metrics.Performance.SlowestTest == "" {
			t.Error("Slowest test should be identified")
		}
		
		if metrics.Performance.FastestTest == "" {
			t.Error("Fastest test should be identified")
		}
		
		// Check quality metrics
		if metrics.Quality.CoveragePercentage != metrics.Quality.CodeQualityScore {
			t.Error("Coverage percentage should match code quality score")
		}
		
		// Check reliability metrics
		if metrics.Reliability.StabilityScore != aggregator.report.Summary.SuccessRate {
			t.Error("Stability score should match success rate")
		}
	})
}

// BenchmarkReportOperations benchmarks report operations
func BenchmarkReportOperations(b *testing.B) {
	config := DefaultReportConfig()
	aggregator := NewReportAggregator(config)
	
	// Add test data
	testResult := TestResult{
		Name:      "benchmark_test",
		Type:      "unit",
		Platform:  config.Platform,
		Status:    "passed",
		Duration:  100 * time.Millisecond,
		Artifacts: []string{"benchmark_output.txt"},
		Metadata: map[string]string{
			"test_file": "benchmark_test.go",
		},
		Timestamp: time.Now(),
	}
	aggregator.AddTestResult(testResult)
	
	b.Run("AddTestResult", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result := testResult
			result.Name = fmt.Sprintf("test_%d", i)
			aggregator.AddTestResult(result)
		}
	})
	
	b.Run("CalculateSummary", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			aggregator.calculateSummary()
		}
	})
	
	b.Run("CalculateMetrics", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			aggregator.calculateMetrics()
		}
	})
	
	b.Run("GenerateJSONReport", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reportFile := filepath.Join(config.OutputDir, fmt.Sprintf("benchmark_report_%d.json", i))
			err := aggregator.GenerateReport(reportFile)
			if err != nil {
				b.Fatalf("Failed to generate report: %v", err)
			}
			os.Remove(reportFile)
		}
	})
}