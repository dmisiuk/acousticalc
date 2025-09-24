package visual

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPerformanceDashboard(t *testing.T) {
	t.Run("create_dashboard", func(t *testing.T) {
		dashboard := NewPerformanceDashboard()

		if dashboard == nil {
			t.Fatal("Dashboard should not be nil")
		}
		if dashboard.Summary.PlatformBreakdown == nil {
			t.Error("PlatformBreakdown should be initialized")
		}
		if dashboard.Summary.TotalReports != 0 {
			t.Error("Initial total reports should be 0")
		}
	})

	t.Run("load_reports", func(t *testing.T) {
		// Create temporary directory with test reports
		tempDir := filepath.Join(os.TempDir(), "test_dashboard_reports")
		defer os.RemoveAll(tempDir)
		os.MkdirAll(tempDir, 0755)

		// Create test reports
		testReports := []CIPerformanceMonitor{
			{
				Platform:        "darwin/arm64",
				TotalDuration:   15 * time.Second,
				ScreenshotCount: 5,
				ArtifactCount:   3,
				StartTime:       time.Now().Add(-1 * time.Hour),
				EndTime:         time.Now().Add(-1 * time.Hour).Add(15 * time.Second),
				IsCI:            true,
				Thresholds: PerformanceThreshold{
					MaxCIOverhead: 30 * time.Second,
				},
			},
			{
				Platform:        "windows/amd64",
				TotalDuration:   25 * time.Second,
				ScreenshotCount: 8,
				ArtifactCount:   5,
				StartTime:       time.Now().Add(-2 * time.Hour),
				EndTime:         time.Now().Add(-2 * time.Hour).Add(25 * time.Second),
				IsCI:            true,
				Thresholds: PerformanceThreshold{
					MaxCIOverhead: 30 * time.Second,
				},
			},
			{
				Platform:        "linux/amd64",
				TotalDuration:   35 * time.Second, // Exceeds threshold
				ScreenshotCount: 10,
				ArtifactCount:   7,
				StartTime:       time.Now().Add(-3 * time.Hour),
				EndTime:         time.Now().Add(-3 * time.Hour).Add(35 * time.Second),
				IsCI:            false,
				Thresholds: PerformanceThreshold{
					MaxCIOverhead: 30 * time.Second,
				},
			},
		}

		// Save test reports
		for i, report := range testReports {
			filename := filepath.Join(tempDir, fmt.Sprintf("ci_performance_%s_%d.json", report.Platform, i))
			data, _ := json.MarshalIndent(report, "", "  ")
			os.WriteFile(filename, data, 0644)
		}

		// Load reports into dashboard
		dashboard := NewPerformanceDashboard()
		err := dashboard.LoadReports(tempDir)

		if err != nil {
			t.Fatalf("Failed to load reports: %v", err)
		}

		// Verify loaded data
		if len(dashboard.Reports) != 3 {
			t.Errorf("Expected 3 reports, got %d", len(dashboard.Reports))
		}

		if dashboard.Summary.TotalReports != 3 {
			t.Errorf("Expected 3 total reports in summary, got %d", dashboard.Summary.TotalReports)
		}

		// Check platform breakdown
		expected := map[string]int{
			"darwin/arm64":  1,
			"windows/amd64": 1,
			"linux/amd64":   1,
		}

		for platform, count := range expected {
			if dashboard.Summary.PlatformBreakdown[platform] != count {
				t.Errorf("Expected %d reports for %s, got %d",
					count, platform, dashboard.Summary.PlatformBreakdown[platform])
			}
		}

		// Check threshold violations (1 report exceeds threshold)
		if dashboard.Summary.ThresholdViolations != 1 {
			t.Errorf("Expected 1 threshold violation, got %d", dashboard.Summary.ThresholdViolations)
		}

		// Verify average calculation
		expectedAvg := (15 + 25 + 35) * time.Second / 3
		if dashboard.Summary.AverageCI_Time != expectedAvg {
			t.Errorf("Expected average %v, got %v", expectedAvg, dashboard.Summary.AverageCI_Time)
		}

		// Verify fastest/slowest
		if dashboard.Summary.FastestCI_Time != 15*time.Second {
			t.Errorf("Expected fastest 15s, got %v", dashboard.Summary.FastestCI_Time)
		}
		if dashboard.Summary.SlowestCI_Time != 35*time.Second {
			t.Errorf("Expected slowest 35s, got %v", dashboard.Summary.SlowestCI_Time)
		}
	})

	t.Run("generate_html", func(t *testing.T) {
		dashboard := NewPerformanceDashboard()

		// Add test data
		dashboard.Reports = []CIPerformanceMonitor{
			{
				Platform:        "darwin/arm64",
				TotalDuration:   20 * time.Second,
				ScreenshotCount: 3,
				ArtifactCount:   2,
				EndTime:         time.Now(),
				IsCI:            true,
				Thresholds: PerformanceThreshold{
					MaxCIOverhead: 30 * time.Second,
				},
			},
		}
		dashboard.calculateSummary()

		// Generate HTML
		tempFile := filepath.Join(os.TempDir(), "test_dashboard.html")
		defer os.Remove(tempFile)

		err := dashboard.GenerateHTML(tempFile)
		if err != nil {
			t.Fatalf("Failed to generate HTML: %v", err)
		}

		// Verify HTML file was created
		if _, err := os.Stat(tempFile); os.IsNotExist(err) {
			t.Error("HTML file was not created")
		}

		// Read and verify HTML content
		content, err := os.ReadFile(tempFile)
		if err != nil {
			t.Fatalf("Failed to read HTML file: %v", err)
		}

		htmlContent := string(content)

		// Check for expected content
		expectedElements := []string{
			"Visual Testing Performance Dashboard",
			"darwin/arm64",
			"20s", // Duration formatting
			"1",   // Total reports
		}

		for _, expected := range expectedElements {
			if !strings.Contains(htmlContent, expected) {
				t.Errorf("HTML content missing expected element: %s", expected)
			}
		}

		// Verify it's valid HTML structure
		if !strings.Contains(htmlContent, "<!DOCTYPE html>") {
			t.Error("HTML should have DOCTYPE declaration")
		}
		if !strings.Contains(htmlContent, "<html") {
			t.Error("HTML should have html tag")
		}
		if !strings.Contains(htmlContent, "</html>") {
			t.Error("HTML should have closing html tag")
		}
	})

	t.Run("save_json", func(t *testing.T) {
		dashboard := NewPerformanceDashboard()
		dashboard.Reports = []CIPerformanceMonitor{
			{
				Platform:      "test/platform",
				TotalDuration: 10 * time.Second,
				EndTime:       time.Now(),
			},
		}
		dashboard.calculateSummary()

		tempFile := filepath.Join(os.TempDir(), "test_dashboard.json")
		defer os.Remove(tempFile)

		err := dashboard.SaveJSON(tempFile)
		if err != nil {
			t.Fatalf("Failed to save JSON: %v", err)
		}

		// Verify JSON file was created and is valid
		data, err := os.ReadFile(tempFile)
		if err != nil {
			t.Fatalf("Failed to read JSON file: %v", err)
		}

		var loadedDashboard PerformanceDashboard
		if err := json.Unmarshal(data, &loadedDashboard); err != nil {
			t.Fatalf("JSON is not valid: %v", err)
		}

		if len(loadedDashboard.Reports) != 1 {
			t.Errorf("Expected 1 report in loaded JSON, got %d", len(loadedDashboard.Reports))
		}
	})

	t.Run("generate_dashboard_integration", func(t *testing.T) {
		// Create test reports directory
		reportsDir := filepath.Join(os.TempDir(), "integration_reports")
		outputDir := filepath.Join(os.TempDir(), "integration_output")
		defer func() {
			os.RemoveAll(reportsDir)
			os.RemoveAll(outputDir)
		}()

		os.MkdirAll(reportsDir, 0755)

		// Create test report
		report := CIPerformanceMonitor{
			Platform:        "integration/test",
			TotalDuration:   12 * time.Second,
			ScreenshotCount: 2,
			ArtifactCount:   1,
			EndTime:         time.Now(),
			IsCI:            true,
		}

		reportData, _ := json.MarshalIndent(report, "", "  ")
		reportFile := filepath.Join(reportsDir, "ci_performance_integration_test_123.json")
		os.WriteFile(reportFile, reportData, 0644)

		// Generate dashboard
		err := GenerateDashboard(reportsDir, outputDir)
		if err != nil {
			t.Fatalf("GenerateDashboard failed: %v", err)
		}

		// Verify output files
		htmlFile := filepath.Join(outputDir, "performance_dashboard.html")
		jsonFile := filepath.Join(outputDir, "performance_dashboard.json")

		if _, err := os.Stat(htmlFile); os.IsNotExist(err) {
			t.Error("HTML dashboard file was not created")
		}
		if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
			t.Error("JSON dashboard file was not created")
		}
	})
}

func TestDashboardSummaryCalculation(t *testing.T) {
	t.Run("empty_reports", func(t *testing.T) {
		dashboard := NewPerformanceDashboard()
		dashboard.calculateSummary()

		if dashboard.Summary.TotalReports != 0 {
			t.Errorf("Expected 0 reports, got %d", dashboard.Summary.TotalReports)
		}
		if dashboard.Summary.AverageCI_Time != 0 {
			t.Errorf("Expected 0 average time, got %v", dashboard.Summary.AverageCI_Time)
		}
	})

	t.Run("single_report", func(t *testing.T) {
		dashboard := NewPerformanceDashboard()
		dashboard.Reports = []CIPerformanceMonitor{
			{
				Platform:      "single/test",
				TotalDuration: 20 * time.Second,
			},
		}

		dashboard.calculateSummary()

		if dashboard.Summary.TotalReports != 1 {
			t.Errorf("Expected 1 report, got %d", dashboard.Summary.TotalReports)
		}
		if dashboard.Summary.AverageCI_Time != 20*time.Second {
			t.Errorf("Expected 20s average, got %v", dashboard.Summary.AverageCI_Time)
		}
		if dashboard.Summary.FastestCI_Time != 20*time.Second {
			t.Errorf("Expected 20s fastest, got %v", dashboard.Summary.FastestCI_Time)
		}
		if dashboard.Summary.SlowestCI_Time != 20*time.Second {
			t.Errorf("Expected 20s slowest, got %v", dashboard.Summary.SlowestCI_Time)
		}
	})

	t.Run("multiple_platforms", func(t *testing.T) {
		dashboard := NewPerformanceDashboard()
		dashboard.Reports = []CIPerformanceMonitor{
			{Platform: "darwin/arm64", TotalDuration: 10 * time.Second},
			{Platform: "darwin/arm64", TotalDuration: 15 * time.Second},
			{Platform: "linux/amd64", TotalDuration: 20 * time.Second},
			{Platform: "windows/amd64", TotalDuration: 25 * time.Second},
		}

		dashboard.calculateSummary()

		expected := map[string]int{
			"darwin/arm64":  2,
			"linux/amd64":   1,
			"windows/amd64": 1,
		}

		for platform, expectedCount := range expected {
			if actual := dashboard.Summary.PlatformBreakdown[platform]; actual != expectedCount {
				t.Errorf("Platform %s: expected %d, got %d", platform, expectedCount, actual)
			}
		}

		// Check min/max
		if dashboard.Summary.FastestCI_Time != 10*time.Second {
			t.Errorf("Expected fastest 10s, got %v", dashboard.Summary.FastestCI_Time)
		}
		if dashboard.Summary.SlowestCI_Time != 25*time.Second {
			t.Errorf("Expected slowest 25s, got %v", dashboard.Summary.SlowestCI_Time)
		}
	})
}

func BenchmarkDashboard(b *testing.B) {
	// Create test data
	reports := make([]CIPerformanceMonitor, 100)
	for i := 0; i < 100; i++ {
		reports[i] = CIPerformanceMonitor{
			Platform:      fmt.Sprintf("test/platform%d", i%5),
			TotalDuration: time.Duration(i) * time.Second,
			EndTime:       time.Now().Add(-time.Duration(i) * time.Hour),
		}
	}

	b.Run("calculate_summary", func(b *testing.B) {
		dashboard := NewPerformanceDashboard()
		dashboard.Reports = reports

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dashboard.calculateSummary()
		}
	})

	b.Run("generate_html", func(b *testing.B) {
		dashboard := NewPerformanceDashboard()
		dashboard.Reports = reports
		dashboard.calculateSummary()

		tempFile := filepath.Join(os.TempDir(), "bench_dashboard.html")
		defer os.Remove(tempFile)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dashboard.GenerateHTML(tempFile)
		}
	})
}
