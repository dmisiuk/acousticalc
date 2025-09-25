package reporting

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// E2ETestResult represents the result of a single E2E test.
type E2ETestResult struct {
	TestName  string
	Passed    bool
	Recording string // Path to the recording file
}

// GenerateE2EReport generates an HTML report for a set of E2E test results.
func GenerateE2EReport(testName, outputDir string, results []E2ETestResult) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	filename := fmt.Sprintf("%s_e2e_report_%s.html",
		testName,
		time.Now().Format("20060102_150405"))
	filePath := filepath.Join(outputDir, filename)

	html := `<!DOCTYPE html>
<html>
<head>
    <title>E2E Test Report: ` + testName + `</title>
    <style>
        body { font-family: 'Monaco', 'Menlo', monospace; margin: 40px; background: #1e1e1e; color: #d4d4d4; }
        .header { border-bottom: 2px solid #32cd32; padding: 20px 0; margin-bottom: 30px; }
        .test { margin: 20px 0; padding: 15px; border-left: 3px solid; }
		.passed { border-left-color: #32cd32; }
		.failed { border-left-color: #cd3232; }
        .recording-link { margin-top: 10px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>E2E Test Report</h1>
        <h2>Test Suite: ` + testName + `</h2>
    </div>`

	for _, result := range results {
		statusClass := "passed"
		statusText := "PASSED"
		if !result.Passed {
			statusClass = "failed"
			statusText = "FAILED"
		}

		html += fmt.Sprintf(`
    <div class="test %s">
        <div><strong>Test:</strong> %s</div>
        <div><strong>Status:</strong> %s</div>`, statusClass, result.TestName, statusText)

		if result.Recording != "" {
			html += fmt.Sprintf(`
        <div class="recording-link">
            <a href="%s" target="_blank">View Recording</a>
        </div>`, result.Recording)
		}

		html += `</div>`
	}

	html += `
</body>
</html>`

	if err := os.WriteFile(filePath, []byte(html), 0644); err != nil {
		return "", fmt.Errorf("failed to write E2E report: %w", err)
	}

	return filePath, nil
}