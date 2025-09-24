package visual

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// PerformanceDashboard generates HTML reports for CI performance monitoring
type PerformanceDashboard struct {
	Reports []CIPerformanceMonitor `json:"reports"`
	Summary DashboardSummary       `json:"summary"`
}

// DashboardSummary contains aggregated performance metrics
type DashboardSummary struct {
	TotalReports        int            `json:"total_reports"`
	AverageCI_Time      time.Duration  `json:"average_ci_time"`
	FastestCI_Time      time.Duration  `json:"fastest_ci_time"`
	SlowestCI_Time      time.Duration  `json:"slowest_ci_time"`
	PlatformBreakdown   map[string]int `json:"platform_breakdown"`
	ThresholdViolations int            `json:"threshold_violations"`
	LastUpdated         time.Time      `json:"last_updated"`
}

const dashboardTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Visual Testing Performance Dashboard</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { background: white; padding: 20px; border-radius: 8px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .metrics { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-bottom: 20px; }
        .metric { background: white; padding: 20px; border-radius: 8px; text-align: center; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .metric-value { font-size: 2em; font-weight: bold; color: #2563eb; }
        .metric-label { color: #6b7280; margin-top: 5px; }
        .chart-container { background: white; padding: 20px; border-radius: 8px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .table-container { background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        table { width: 100%; border-collapse: collapse; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #e5e7eb; }
        th { background: #f9fafb; font-weight: 600; }
        .status-pass { color: #059669; }
        .status-fail { color: #dc2626; }
        .platform-tag { padding: 4px 8px; background: #e5e7eb; border-radius: 4px; font-size: 0.875em; }
        .violation { background: #fef2f2; }
        .chart { height: 300px; background: #f9fafb; border-radius: 4px; display: flex; align-items: center; justify-content: center; }
        .timestamp { color: #6b7280; font-size: 0.875em; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🧪 Visual Testing Performance Dashboard</h1>
            <p>Real-time monitoring of visual testing performance across platforms</p>
            <div class="timestamp">Last Updated: {{.Summary.LastUpdated.Format "2006-01-02 15:04:05 UTC"}}</div>
        </div>

        <div class="metrics">
            <div class="metric">
                <div class="metric-value">{{.Summary.TotalReports}}</div>
                <div class="metric-label">Total Test Runs</div>
            </div>
            <div class="metric">
                <div class="metric-value">{{formatDuration .Summary.AverageCI_Time}}</div>
                <div class="metric-label">Average CI Time</div>
            </div>
            <div class="metric">
                <div class="metric-value">{{formatDuration .Summary.FastestCI_Time}}</div>
                <div class="metric-label">Fastest Run</div>
            </div>
            <div class="metric">
                <div class="metric-value">{{formatDuration .Summary.SlowestCI_Time}}</div>
                <div class="metric-label">Slowest Run</div>
            </div>
            <div class="metric {{if gt .Summary.ThresholdViolations 0}}violation{{end}}">
                <div class="metric-value">{{.Summary.ThresholdViolations}}</div>
                <div class="metric-label">Threshold Violations</div>
            </div>
        </div>

        <div class="chart-container">
            <h2>Platform Distribution</h2>
            <div class="chart">
                {{range $platform, $count := .Summary.PlatformBreakdown}}
                    <div style="margin: 10px;">
                        <span class="platform-tag">{{$platform}}: {{$count}}</span>
                    </div>
                {{end}}
            </div>
        </div>

        <div class="table-container">
            <table>
                <thead>
                    <tr>
                        <th>Timestamp</th>
                        <th>Platform</th>
                        <th>Duration</th>
                        <th>Screenshots</th>
                        <th>Artifacts</th>
                        <th>Status</th>
                        <th>CI Environment</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Reports}}
                    <tr {{if eq .getStatus "FAIL"}}class="violation"{{end}}>
                        <td>{{.EndTime.Format "01-02 15:04"}}</td>
                        <td><span class="platform-tag">{{.Platform}}</span></td>
                        <td>{{formatDuration .TotalDuration}}</td>
                        <td>{{.ScreenshotCount}}</td>
                        <td>{{.ArtifactCount}}</td>
                        <td class="{{if eq .getStatus "PASS"}}status-pass{{else}}status-fail{{end}}">
                            {{.getStatus}}
                        </td>
                        <td>{{if .IsCI}}Yes{{else}}Local{{end}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>
    </div>
</body>
</html>
`

// NewPerformanceDashboard creates a new dashboard instance
func NewPerformanceDashboard() *PerformanceDashboard {
	return &PerformanceDashboard{
		Reports: make([]CIPerformanceMonitor, 0),
		Summary: DashboardSummary{
			PlatformBreakdown: make(map[string]int),
			LastUpdated:       time.Now(),
		},
	}
}

// LoadReports loads performance reports from a directory
func (d *PerformanceDashboard) LoadReports(reportsDir string) error {
	files, err := filepath.Glob(filepath.Join(reportsDir, "ci_performance_*.json"))
	if err != nil {
		return fmt.Errorf("failed to list report files: %w", err)
	}

	d.Reports = make([]CIPerformanceMonitor, 0, len(files))

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue // Skip corrupted files
		}

		var report CIPerformanceMonitor
		if err := json.Unmarshal(data, &report); err != nil {
			continue // Skip invalid JSON
		}

		d.Reports = append(d.Reports, report)
	}

	// Sort reports by timestamp (newest first)
	sort.Slice(d.Reports, func(i, j int) bool {
		return d.Reports[i].EndTime.After(d.Reports[j].EndTime)
	})

	d.calculateSummary()
	return nil
}

// calculateSummary computes dashboard summary metrics
func (d *PerformanceDashboard) calculateSummary() {
	if len(d.Reports) == 0 {
		return
	}

	d.Summary.TotalReports = len(d.Reports)
	d.Summary.LastUpdated = time.Now()

	var totalDuration time.Duration
	d.Summary.FastestCI_Time = d.Reports[0].TotalDuration
	d.Summary.SlowestCI_Time = d.Reports[0].TotalDuration

	for _, report := range d.Reports {
		totalDuration += report.TotalDuration

		if report.TotalDuration < d.Summary.FastestCI_Time {
			d.Summary.FastestCI_Time = report.TotalDuration
		}
		if report.TotalDuration > d.Summary.SlowestCI_Time {
			d.Summary.SlowestCI_Time = report.TotalDuration
		}

		// Count platform distribution
		d.Summary.PlatformBreakdown[report.Platform]++

		// Count threshold violations
		if report.getStatus() == "FAIL" {
			d.Summary.ThresholdViolations++
		}
	}

	d.Summary.AverageCI_Time = totalDuration / time.Duration(len(d.Reports))
}

// GenerateHTML creates an HTML dashboard report
func (d *PerformanceDashboard) GenerateHTML(outputPath string) error {
	tmpl := template.Must(template.New("dashboard").Funcs(template.FuncMap{
		"formatDuration": func(dur time.Duration) string {
			if dur < time.Second {
				return fmt.Sprintf("%dms", dur.Milliseconds())
			}
			return dur.Round(time.Millisecond).String()
		},
	}).Parse(dashboardTemplate))

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create dashboard file: %w", err)
	}
	defer file.Close()

	return tmpl.Execute(file, d)
}

// SaveJSON saves dashboard data as JSON
func (d *PerformanceDashboard) SaveJSON(outputPath string) error {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal dashboard data: %w", err)
	}

	return os.WriteFile(outputPath, data, 0644)
}

// GenerateDashboard creates a complete performance dashboard
func GenerateDashboard(reportsDir, outputDir string) error {
	dashboard := NewPerformanceDashboard()

	if err := dashboard.LoadReports(reportsDir); err != nil {
		return fmt.Errorf("failed to load reports: %w", err)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate HTML dashboard
	htmlPath := filepath.Join(outputDir, "performance_dashboard.html")
	if err := dashboard.GenerateHTML(htmlPath); err != nil {
		return fmt.Errorf("failed to generate HTML dashboard: %w", err)
	}

	// Save JSON data
	jsonPath := filepath.Join(outputDir, "performance_dashboard.json")
	if err := dashboard.SaveJSON(jsonPath); err != nil {
		return fmt.Errorf("failed to save JSON data: %w", err)
	}

	fmt.Printf("Performance dashboard generated:\n")
	fmt.Printf("  HTML: %s\n", htmlPath)
	fmt.Printf("  JSON: %s\n", jsonPath)
	fmt.Printf("  Reports analyzed: %d\n", dashboard.Summary.TotalReports)
	if dashboard.Summary.ThresholdViolations > 0 {
		fmt.Printf("  ⚠️ Threshold violations: %d\n", dashboard.Summary.ThresholdViolations)
	}

	return nil
}
