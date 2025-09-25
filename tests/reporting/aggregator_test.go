package reporting

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestResult represents a test result from any test type
type TestResult struct {
	Name      string        `json:"name"`
	Type      string        `json:"type"` // "unit", "integration", "e2e", "visual", "cross_platform"
	Platform  string        `json:"platform"`
	Status    string        `json:"status"` // "pass", "fail", "skip", "error"
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
	Coverage  float64       `json:"coverage,omitempty"`
	Artifacts []string      `json:"artifacts,omitempty"`
	Error     string        `json:"error,omitempty"`
	Metadata  interface{}   `json:"metadata,omitempty"`
}

// PlatformTestResults represents test results aggregated by platform
type PlatformTestResults struct {
	Platform  string        `json:"platform"`
	Results   []TestResult  `json:"results"`
	Summary   TestSummary   `json:"summary"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
}

// TestSummary provides a summary of test results
type TestSummary struct {
	Total    int           `json:"total"`
	Passed   int           `json:"passed"`
	Failed   int           `json:"failed"`
	Skipped  int           `json:"skipped"`
	Errors   int           `json:"errors"`
	PassRate float64       `json:"pass_rate"`
	Coverage float64       `json:"coverage,omitempty"`
	Duration time.Duration `json:"duration"`
}

// CrossPlatformReport represents a comprehensive cross-platform test report
type CrossPlatformReport struct {
	GeneratedAt           time.Time             `json:"generated_at"`
	TestRunID             string                `json:"test_run_id"`
	Platforms             []PlatformTestResults `json:"platforms"`
	GlobalSummary         TestSummary           `json:"global_summary"`
	CrossPlatformAnalysis CrossPlatformAnalysis `json:"cross_platform_analysis"`
	Recommendations       []string              `json:"recommendations"`
}

// CrossPlatformAnalysis provides analysis across platforms
type CrossPlatformAnalysis struct {
	ConsistencyScore    float64                `json:"consistency_score"`
	PlatformDifferences []PlatformDifference   `json:"platform_differences"`
	PerformanceMetrics  map[string]interface{} `json:"performance_metrics"`
	ReliabilityScore    float64                `json:"reliability_score"`
}

// PlatformDifference represents differences between platforms
type PlatformDifference struct {
	TestName   string   `json:"test_name"`
	TestType   string   `json:"test_type"`
	Platforms  []string `json:"platforms"`
	Difference string   `json:"difference"`
	Severity   string   `json:"severity"` // "low", "medium", "high", "critical"
	Impact     string   `json:"impact"`
}

// TestResultAggregator aggregates test results from multiple sources
type TestResultAggregator struct {
	baseDir   string
	platforms map[string]*PlatformTestResults
	results   map[string][]TestResult
	startTime time.Time
}

// NewTestResultAggregator creates a new test result aggregator
func NewTestResultAggregator(baseDir string) *TestResultAggregator {
	return &TestResultAggregator{
		baseDir:   baseDir,
		platforms: make(map[string]*PlatformTestResults),
		results:   make(map[string][]TestResult),
		startTime: time.Now(),
	}
}

// AddTestResult adds a test result to the aggregator
func (tra *TestResultAggregator) AddTestResult(result TestResult) error {
	// Initialize platform results if not exists
	if _, exists := tra.platforms[result.Platform]; !exists {
		tra.platforms[result.Platform] = &PlatformTestResults{
			Platform:  result.Platform,
			Results:   make([]TestResult, 0),
			StartTime: time.Now(),
		}
	}

	// Add result to platform
	tra.platforms[result.Platform].Results = append(tra.platforms[result.Platform].Results, result)

	// Add to type-based results
	tra.results[result.Type] = append(tra.results[result.Type], result)

	return nil
}

// AggregateFromArtifacts aggregates test results from artifact files
func (tra *TestResultAggregator) AggregateFromArtifacts() error {
	// Scan for test result files
	artifactDir := filepath.Join(tra.baseDir, "tests", "artifacts")

	// Look for JSON result files
	resultFiles := []string{
		filepath.Join(artifactDir, "e2e", "*.json"),
		filepath.Join(artifactDir, "cross_platform", "*.json"),
		filepath.Join(artifactDir, "visual", "*.json"),
	}

	for _, pattern := range resultFiles {
		files, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("failed to glob pattern %s: %w", pattern, err)
		}

		for _, file := range files {
			if err := tra.processResultFile(file); err != nil {
				return fmt.Errorf("failed to process result file %s: %w", file, err)
			}
		}
	}

	// Generate summaries for each platform
	for platform, platformResults := range tra.platforms {
		platformResults.Summary = tra.generateSummary(platformResults.Results)
		platformResults.EndTime = time.Now()
		platformResults.Duration = platformResults.EndTime.Sub(platformResults.StartTime)
	}

	return nil
}

// processResultFile processes a single result file
func (tra *TestResultAggregator) processResultFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Try to parse as JSON array of results
	var results []TestResult
	if err := json.Unmarshal(data, &results); err != nil {
		// If that fails, try to parse as single result
		var singleResult TestResult
		if err := json.Unmarshal(data, &singleResult); err == nil {
			results = []TestResult{singleResult}
		} else {
			return fmt.Errorf("failed to parse JSON from %s: %w", filePath, err)
		}
	}

	// Add all results
	for _, result := range results {
		if result.Platform == "" {
			// Infer platform from file path if not specified
			result.Platform = tra.inferPlatformFromPath(filePath)
		}
		if err := tra.AddTestResult(result); err != nil {
			return fmt.Errorf("failed to add result: %w", err)
		}
	}

	return nil
}

// inferPlatformFromPath infers platform from file path
func (tra *TestResultAggregator) inferPlatformFromPath(filePath string) string {
	// Simple platform inference from path
	if contains(filePath, "linux") || contains(filePath, "ubuntu") {
		return "linux"
	}
	if contains(filePath, "darwin") || contains(filePath, "macos") {
		return "darwin"
	}
	if contains(filePath, "windows") {
		return "windows"
	}
	return "unknown"
}

// generateSummary generates a summary from test results
func (tra *TestResultAggregator) generateSummary(results []TestResult) TestSummary {
	summary := TestSummary{
		Total:   len(results),
		Started: time.Now(),
	}

	var totalCoverage float64
	var coverageCount int
	var totalDuration time.Duration

	for _, result := range results {
		switch result.Status {
		case "pass":
			summary.Passed++
		case "fail":
			summary.Failed++
		case "skip":
			summary.Skipped++
		case "error":
			summary.Errors++
		}

		if result.Coverage > 0 {
			totalCoverage += result.Coverage
			coverageCount++
		}

		totalDuration += result.Duration
	}

	if summary.Total > 0 {
		summary.PassRate = float64(summary.Passed) / float64(summary.Total) * 100
	}

	if coverageCount > 0 {
		summary.Coverage = totalCoverage / float64(coverageCount)
	}

	summary.Duration = totalDuration

	return summary
}

// GenerateCrossPlatformReport generates a comprehensive cross-platform report
func (tra *TestResultAggregator) GenerateCrossPlatformReport() (*CrossPlatformReport, error) {
	report := &CrossPlatformReport{
		GeneratedAt: time.Now(),
		TestRunID:   fmt.Sprintf("run_%d", time.Now().Unix()),
	}

	// Convert platforms map to slice
	for _, platformResults := range tra.platforms {
		report.Platforms = append(report.Platforms, *platformResults)
	}

	// Sort platforms by name
	sort.Slice(report.Platforms, func(i, j int) bool {
		return report.Platforms[i].Platform < report.Platforms[j].Platform
	})

	// Generate global summary
	report.GlobalSummary = tra.generateGlobalSummary()

	// Perform cross-platform analysis
	report.CrossPlatformAnalysis = tra.performCrossPlatformAnalysis()

	// Generate recommendations
	report.Recommendations = tra.generateRecommendations()

	return report, nil
}

// generateGlobalSummary generates a global summary across all platforms
func (tra *TestResultAggregator) generateGlobalSummary() TestSummary {
	var allResults []TestResult
	for _, platformResults := range tra.platforms {
		allResults = append(allResults, platformResults.Results...)
	}
	return tra.generateSummary(allResults)
}

// performCrossPlatformAnalysis performs analysis across platforms
func (tra *TestResultAggregator) performCrossPlatformAnalysis() CrossPlatformAnalysis {
	analysis := CrossPlatformAnalysis{
		PerformanceMetrics: make(map[string]interface{}),
	}

	// Calculate consistency score
	analysis.ConsistencyScore = tra.calculateConsistencyScore()

	// Find platform differences
	analysis.PlatformDifferences = tra.findPlatformDifferences()

	// Calculate reliability score
	analysis.ReliabilityScore = tra.calculateReliabilityScore()

	// Add performance metrics
	analysis.PerformanceMetrics = tra.gatherPerformanceMetrics()

	return analysis
}

// calculateConsistencyScore calculates how consistent test results are across platforms
func (tra *TestResultAggregator) calculateConsistencyScore() float64 {
	if len(tra.platforms) < 2 {
		return 100.0 // Perfect consistency if only one platform
	}

	testResults := make(map[string]map[string]string) // test_name -> platform -> status

	// Collect results for each test across platforms
	for platform, platformData := range tra.platforms {
		for _, result := range platformData.Results {
			if _, exists := testResults[result.Name]; !exists {
				testResults[result.Name] = make(map[string]string)
			}
			testResults[result.Name][platform] = result.Status
		}
	}

	// Calculate consistency
	consistentTests := 0
	totalTests := len(testResults)

	for _, platformResults := range testResults {
		if len(platformResults) < 2 {
			continue // Skip tests that don't run on multiple platforms
		}

		// Check if all platforms have the same status
		statuses := make(map[string]int)
		for _, status := range platformResults {
			statuses[status]++
		}

		// If all platforms have the same status, it's consistent
		if len(statuses) == 1 {
			consistentTests++
		}
	}

	if totalTests == 0 {
		return 100.0
	}

	return float64(consistentTests) / float64(totalTests) * 100
}

// findPlatformDifferences finds differences between platforms
func (tra *TestResultAggregator) findPlatformDifferences() []PlatformDifference {
	var differences []PlatformDifference

	testResults := make(map[string]map[string]TestResult)

	// Collect results for each test across platforms
	for platform, platformData := range tra.platforms {
		for _, result := range platformData.Results {
			if _, exists := testResults[result.Name]; !exists {
				testResults[result.Name] = make(map[string]TestResult)
			}
			testResults[result.Name][platform] = result
		}
	}

	// Find differences
	for testName, platformResults := range testResults {
		if len(platformResults) < 2 {
			continue // Skip tests that don't run on multiple platforms
		}

		// Check for status differences
		statuses := make(map[string][]string)
		durations := make(map[string]time.Duration)

		for platform, result := range platformResults {
			statuses[result.Status] = append(statuses[result.Status], platform)
			durations[platform] = result.Duration
		}

		// If there are different statuses, create a difference entry
		if len(statuses) > 1 {
			var allPlatforms []string
			for _, platforms := range statuses {
				allPlatforms = append(allPlatforms, platforms...)
			}

			difference := PlatformDifference{
				TestName:   testName,
				TestType:   "unknown", // Would need to be inferred
				Platforms:  allPlatforms,
				Difference: fmt.Sprintf("Status mismatch: %v", statuses),
				Severity:   "medium",
				Impact:     "Test behavior differs between platforms",
			}

			// Determine severity based on failure types
			if _, hasFailures := statuses["fail"]; hasFailures {
				difference.Severity = "high"
				difference.Impact = "Test failures on some platforms may indicate platform-specific bugs"
			}

			differences = append(differences, difference)
		}
	}

	return differences
}

// calculateReliabilityScore calculates overall reliability score
func (tra *TestResultAggregator) calculateReliabilityScore() float64 {
	var totalTests, passedTests int

	for _, platformData := range tra.platforms {
		for _, result := range platformData.Results {
			totalTests++
			if result.Status == "pass" {
				passedTests++
			}
		}
	}

	if totalTests == 0 {
		return 100.0
	}

	return float64(passedTests) / float64(totalTests) * 100
}

// gatherPerformanceMetrics gathers performance metrics across platforms
func (tra *TestResultAggregator) gatherPerformanceMetrics() map[string]interface{} {
	metrics := make(map[string]interface{})

	// Calculate average durations by platform
	platformDurations := make(map[string][]time.Duration)
	for platform, platformData := range tra.platforms {
		for _, result := range platformData.Results {
			platformDurations[platform] = append(platformDurations[platform], result.Duration)
		}
	}

	avgDurations := make(map[string]interface{})
	for platform, durations := range platformDurations {
		if len(durations) > 0 {
			var total time.Duration
			for _, d := range durations {
				total += d
			}
			avg := total / time.Duration(len(durations))
			avgDurations[platform] = avg.String()
		}
	}
	metrics["average_durations"] = avgDurations

	// Add test counts by type
	typeCounts := make(map[string]int)
	for _, platformData := range tra.platforms {
		for _, result := range platformData.Results {
			typeCounts[result.Type]++
		}
	}
	metrics["test_counts_by_type"] = typeCounts

	return metrics
}

// generateRecommendations generates recommendations based on test results
func (tra *TestResultAggregator) generateRecommendations() []string {
	var recommendations []string

	// Analyze consistency
	if tra.calculateConsistencyScore() < 90 {
		recommendations = append(recommendations,
			"Low consistency score detected. Investigate platform-specific differences and improve cross-platform compatibility.")
	}

	// Analyze reliability
	if tra.calculateReliabilityScore() < 95 {
		recommendations = append(recommendations,
			"Reliability score below target. Focus on fixing failing tests and improving error handling.")
	}

	// Check for platform differences
	differences := tra.findPlatformDifferences()
	if len(differences) > 0 {
		highSeverity := 0
		for _, diff := range differences {
			if diff.Severity == "high" || diff.Severity == "critical" {
				highSeverity++
			}
		}
		if highSeverity > 0 {
			recommendations = append(recommendations,
				fmt.Sprintf("Found %d high-severity platform differences. Prioritize fixing these issues.", highSeverity))
		}
	}

	// Add general recommendations
	recommendations = append(recommendations,
		"Consider adding more E2E tests to improve coverage of user workflows.")
	recommendations = append(recommendations,
		"Implement automated regression testing based on these cross-platform results.")

	return recommendations
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				findSubstring(s, substr))))
}

// findSubstring finds a substring in a string
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
