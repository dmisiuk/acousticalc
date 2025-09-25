package reporting

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// E2EReportData holds the aggregated data for the E2E test report.
type E2EReportData struct {
	TestRuns []TestRun
}

// TestRun represents a single E2E test execution.
type TestRun struct {
	Name      string
	Platform  string
	Recording string // Relative path to the recording file
}

// DiscoverE2EArtifacts scans the artifacts directory to find all E2E test runs.
func DiscoverE2EArtifacts(artifactsDir string) (*E2EReportData, error) {
	reportData := &E2EReportData{
		TestRuns: make([]TestRun, 0),
	}

	e2eArtifactsDir := filepath.Join(artifactsDir, "e2e")
	if _, err := os.Stat(e2eArtifactsDir); os.IsNotExist(err) {
		// It's not an error if the directory doesn't exist, just means no E2E tests were run
		return reportData, nil
	}

	// Walk through the E2E artifacts directory
	err := filepath.Walk(e2eArtifactsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a recording (.cast for Unix, .txt for Windows)
		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".cast") || strings.HasSuffix(info.Name(), ".txt")) {
			// Make path relative to the e2e artifacts dir to easily extract platform and test name
			relativePath, err := filepath.Rel(e2eArtifactsDir, path)
			if err != nil {
				return err // Or log a warning
			}

			parts := strings.Split(relativePath, string(os.PathSeparator))

			// Expected structure: <platform>/<test_name>/recordings/<filename>.cast
			if len(parts) >= 4 {
				platform := parts[0]
				testName := parts[1]

				// Make the recording path relative to the root artifacts directory for the HTML report
				reportRelativePath, err := filepath.Rel(artifactsDir, path)
				if err != nil {
					reportRelativePath = path // Fallback to full path
				}

				run := TestRun{
					Name:      testName,
					Platform:  platform,
					Recording: reportRelativePath,
				}
				reportData.TestRuns = append(reportData.TestRuns, run)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk E2E artifacts directory: %w", err)
	}

	return reportData, nil
}

// GenerateE2EReport creates an HTML report from the aggregated E2E test data.
func GenerateE2EReport(reportData *E2EReportData, outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create report output directory: %w", err)
	}

	reportPath := filepath.Join(outputDir, "e2e_report.html")
	htmlContent := generateHTMLContent(reportData)

	if err := os.WriteFile(reportPath, []byte(htmlContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write E2E report: %w", err)
	}

	return reportPath, nil
}

func generateHTMLContent(data *E2EReportData) string {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>E2E Test Report</title>
    <style>
        body { font-family: 'Monaco', 'Menlo', monospace; margin: 40px; background: #1e1e1e; color: #d4d4d4; }
        .header { border-bottom: 2px solid #32cd32; padding: 20px 0; margin-bottom: 30px; }
        table { width: 100%; border-collapse: collapse; }
        th, td { padding: 10px; border: 1px solid #444; text-align: left; }
        th { background: #2d2d2d; }
        a { color: #32cd32; }
    </style>
</head>
<body>
    <div class="header">
        <h1>E2E Test Report</h1>
        <p>Total Test Runs: ` + fmt.Sprintf("%d", len(data.TestRuns)) + `</p>
    </div>
    <table>
        <tr>
            <th>Test Name</th>
            <th>Platform</th>
            <th>Recording</th>
        </tr>`

	for _, run := range data.TestRuns {
		html += fmt.Sprintf(`
        <tr>
            <td>%s</td>
            <td>%s</td>
            <td><a href="%s" target="_blank">View Recording</a></td>
        </tr>`, run.Name, run.Platform, run.Recording)
	}

	html += `
    </table>
</body>
</html>`
	return html
}