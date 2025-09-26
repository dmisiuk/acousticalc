package reporting

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AggregateResults scans a directory for test artifacts and aggregates them into a single report.
func AggregateResults(artifactsDir, finalReportPath string) error {
	var report strings.Builder
	report.WriteString("Aggregated Test Report\n")
	report.WriteString("======================\n\n")

	err := filepath.Walk(artifactsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".xml") || strings.HasSuffix(info.Name(), ".txt")) {
			content, readErr := os.ReadFile(path)
			if readErr != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not read artifact file %s: %v\n", path, readErr)
				return nil // Continue walking
			}
			report.WriteString(fmt.Sprintf("--- Artifact: %s ---\n", path))
			report.WriteString(string(content))
			report.WriteString("\n\n")
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking artifacts directory: %w", err)
	}

	// Ensure the report directory exists
	reportDir := filepath.Dir(finalReportPath)
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		return fmt.Errorf("failed to create report directory: %w", err)
	}

	return os.WriteFile(finalReportPath, []byte(report.String()), 0644)
}

// GenerateHTMLReport generates an HTML report from aggregated data.
// This is a placeholder for a more sophisticated HTML generation logic.
func GenerateHTMLReport(reportData string, htmlPath string) error {
	htmlContent := `
<!DOCTYPE html>
<html>
<head>
<title>Test Report</title>
<style>
  body { font-family: sans-serif; }
  pre { background-color: #f4f4f4; padding: 1em; border: 1px solid #ddd; }
</style>
</head>
<body>
<h1>Aggregated Test Report</h1>
<pre>` + reportData + `</pre>
</body>
</html>`

	return os.WriteFile(htmlPath, []byte(htmlContent), 0644)
}
