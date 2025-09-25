package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ArtifactInfo contains metadata about a visual testing artifact
type ArtifactInfo struct {
	Path      string            `json:"path"`
	Type      string            `json:"type"` // "screenshot", "report", "storyboard", "chart"
	TestName  string            `json:"test_name"`
	EventType string            `json:"event_type"`
	Category  string            `json:"category"` // "unit", "integration", "e2e"
	Timestamp time.Time         `json:"timestamp"`
	Size      int64             `json:"size_bytes"`
	Metadata  map[string]string `json:"metadata"`
}

// ArtifactManager handles visual testing artifact management
type ArtifactManager struct {
	BaseDir   string
	Artifacts []ArtifactInfo
}

// NewArtifactManager creates a new artifact manager
func NewArtifactManager(baseDir string) *ArtifactManager {
	return &ArtifactManager{
		BaseDir:   baseDir,
		Artifacts: make([]ArtifactInfo, 0),
	}
}

// ScanArtifacts scans the base directory for all visual testing artifacts
func (am *ArtifactManager) ScanArtifacts() error {
	am.Artifacts = am.Artifacts[:0] // Clear existing

	err := filepath.Walk(am.BaseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		artifact := am.classifyArtifact(path, info)
		if artifact != nil {
			am.Artifacts = append(am.Artifacts, *artifact)
		}

		return nil
	})

	return err
}

// classifyArtifact determines artifact type and extracts metadata
func (am *ArtifactManager) classifyArtifact(path string, info os.FileInfo) *ArtifactInfo {
	relPath, err := filepath.Rel(am.BaseDir, path)
	if err != nil {
		return nil
	}

	ext := strings.ToLower(filepath.Ext(path))
	fileName := strings.TrimSuffix(info.Name(), ext)

	artifact := &ArtifactInfo{
		Path:      relPath,
		Size:      info.Size(),
		Timestamp: info.ModTime(),
		Metadata:  make(map[string]string),
	}

	// Classify by file extension and location
	switch ext {
	case ".png":
		artifact.Type = "screenshot"
		am.parseScreenshotMetadata(artifact, fileName, relPath)
	case ".html":
		if strings.Contains(relPath, "reports") {
			artifact.Type = "report"
		} else if strings.Contains(relPath, "storyboard") {
			artifact.Type = "storyboard"
		} else if strings.Contains(relPath, "charts") {
			artifact.Type = "chart"
		} else {
			artifact.Type = "report"
		}
		am.parseHTMLMetadata(artifact, fileName)
	case ".json":
		artifact.Type = "metadata"
		am.parseJSONMetadata(artifact, fileName)
	case ".cast":
		artifact.Type = "recording"
		am.parseRecordingMetadata(artifact, fileName)
	default:
		return nil // Skip unsupported file types
	}

	// Extract category from path
	if strings.Contains(relPath, "unit") {
		artifact.Category = "unit"
	} else if strings.Contains(relPath, "integration") {
		artifact.Category = "integration"
	} else if strings.Contains(relPath, "e2e") {
		artifact.Category = "e2e"
	} else {
		artifact.Category = "unknown"
	}

	return artifact
}

// parseScreenshotMetadata extracts metadata from screenshot filenames
func (am *ArtifactManager) parseScreenshotMetadata(artifact *ArtifactInfo, fileName, path string) {
	// Format: testname_eventtype_timestamp.png
	parts := strings.Split(fileName, "_")

	if len(parts) >= 3 {
		// Extract test name (everything before the last two parts)
		testNameParts := parts[:len(parts)-2]
		artifact.TestName = strings.Join(testNameParts, "_")

		// Extract event type (second to last part)
		artifact.EventType = parts[len(parts)-2]

		// Extract timestamp (last part)
		timestampStr := parts[len(parts)-1]
		if timestamp, err := time.Parse("20060102_150405", timestampStr); err == nil {
			artifact.Timestamp = timestamp
		}
	}

	artifact.Metadata["format"] = "PNG"
	artifact.Metadata["color_space"] = "sRGB"
}

// parseHTMLMetadata extracts metadata from HTML report filenames
func (am *ArtifactManager) parseHTMLMetadata(artifact *ArtifactInfo, fileName string) {
	// Extract test name from filename
	if strings.Contains(fileName, "_visual_report") {
		artifact.TestName = strings.Replace(fileName, "_visual_report", "", 1)
	} else if strings.Contains(fileName, "_storyboard") {
		artifact.TestName = strings.Replace(fileName, "_storyboard", "", 1)
	} else if strings.Contains(fileName, "_timeline") {
		artifact.TestName = strings.Replace(fileName, "_timeline", "", 1)
	}

	artifact.Metadata["format"] = "HTML"
}

// parseJSONMetadata extracts metadata from JSON metadata files
func (am *ArtifactManager) parseJSONMetadata(artifact *ArtifactInfo, fileName string) {
	if strings.Contains(fileName, "_metadata") {
		artifact.TestName = strings.Replace(fileName, "_metadata", "", 1)
	}
	artifact.Metadata["format"] = "JSON"
}

// parseRecordingMetadata extracts metadata from recording filenames
func (am *ArtifactManager) parseRecordingMetadata(artifact *ArtifactInfo, fileName string) {
	// Format: testname_recording_timestamp.cast
	parts := strings.Split(fileName, "_")

	if len(parts) >= 3 {
		// Extract test name (everything before the last two parts)
		testNameParts := parts[:len(parts)-2]
		artifact.TestName = strings.Join(testNameParts, "_")

		// The event type is always "recording"
		artifact.EventType = "recording"

		// Extract timestamp (last part)
		timestampStr := strings.TrimSuffix(parts[len(parts)-1], ".cast")
		if timestamp, err := time.Parse("20060102_150405", timestampStr); err == nil {
			artifact.Timestamp = timestamp
		}
	}

	artifact.Metadata["format"] = "asciinema v2"
}

// ListArtifacts lists all artifacts with optional filtering
func (am *ArtifactManager) ListArtifacts(artifactType, category, testName string) []ArtifactInfo {
	filtered := make([]ArtifactInfo, 0)

	for _, artifact := range am.Artifacts {
		// Apply filters
		if artifactType != "" && artifact.Type != artifactType {
			continue
		}
		if category != "" && artifact.Category != category {
			continue
		}
		if testName != "" && !strings.Contains(artifact.TestName, testName) {
			continue
		}

		filtered = append(filtered, artifact)
	}

	// Sort by timestamp (newest first)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp.After(filtered[j].Timestamp)
	})

	return filtered
}

// GetSummary returns a summary of all artifacts
func (am *ArtifactManager) GetSummary() map[string]interface{} {
	summary := map[string]interface{}{
		"total_artifacts": len(am.Artifacts),
		"total_size":      int64(0),
		"by_type":         make(map[string]int),
		"by_category":     make(map[string]int),
		"by_test":         make(map[string]int),
	}

	byType := summary["by_type"].(map[string]int)
	byCategory := summary["by_category"].(map[string]int)
	byTest := summary["by_test"].(map[string]int)

	for _, artifact := range am.Artifacts {
		summary["total_size"] = summary["total_size"].(int64) + artifact.Size
		byType[artifact.Type]++
		byCategory[artifact.Category]++
		if artifact.TestName != "" {
			byTest[artifact.TestName]++
		}
	}

	return summary
}

// CleanupArtifacts removes artifacts older than specified days
func (am *ArtifactManager) CleanupArtifacts(olderThanDays int, dryRun bool) ([]string, error) {
	cutoffTime := time.Now().AddDate(0, 0, -olderThanDays)
	toDelete := make([]string, 0)

	for _, artifact := range am.Artifacts {
		if artifact.Timestamp.Before(cutoffTime) {
			fullPath := filepath.Join(am.BaseDir, artifact.Path)
			toDelete = append(toDelete, fullPath)

			if !dryRun {
				if err := os.Remove(fullPath); err != nil {
					return toDelete, fmt.Errorf("failed to delete %s: %w", fullPath, err)
				}
			}
		}
	}

	return toDelete, nil
}

// ExportMetadata exports artifact metadata to JSON file
func (am *ArtifactManager) ExportMetadata(outputPath string) error {
	data := map[string]interface{}{
		"generated_at": time.Now(),
		"base_dir":     am.BaseDir,
		"summary":      am.GetSummary(),
		"artifacts":    am.Artifacts,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	return os.WriteFile(outputPath, jsonData, 0644)
}

// formatBytes formats byte count in human readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func main() {
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	baseDir := "tests/artifacts"
	if envDir := os.Getenv("ARTIFACT_DIR"); envDir != "" {
		baseDir = envDir
	}

	manager := NewArtifactManager(baseDir)

	switch os.Args[1] {
	case "scan":
		fmt.Printf("Scanning artifacts in: %s\n", baseDir)
		if err := manager.ScanArtifacts(); err != nil {
			fmt.Fprintf(os.Stderr, "Error scanning artifacts: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Found %d artifacts\n", len(manager.Artifacts))

	case "list":
		if err := manager.ScanArtifacts(); err != nil {
			fmt.Fprintf(os.Stderr, "Error scanning artifacts: %v\n", err)
			os.Exit(1)
		}

		// Parse filters
		var artifactType, category, testName string
		for i := 2; i < len(os.Args); i += 2 {
			if i+1 < len(os.Args) {
				switch os.Args[i] {
				case "--type":
					artifactType = os.Args[i+1]
				case "--category":
					category = os.Args[i+1]
				case "--test":
					testName = os.Args[i+1]
				}
			}
		}

		artifacts := manager.ListArtifacts(artifactType, category, testName)
		fmt.Printf("%-20s %-12s %-10s %-15s %-12s %s\n", "Test Name", "Type", "Category", "Event", "Size", "Path")
		fmt.Println(strings.Repeat("-", 90))

		for _, artifact := range artifacts {
			fmt.Printf("%-20s %-12s %-10s %-15s %-12s %s\n",
				truncate(artifact.TestName, 20),
				artifact.Type,
				artifact.Category,
				truncate(artifact.EventType, 15),
				formatBytes(artifact.Size),
				artifact.Path)
		}

	case "summary":
		if err := manager.ScanArtifacts(); err != nil {
			fmt.Fprintf(os.Stderr, "Error scanning artifacts: %v\n", err)
			os.Exit(1)
		}

		summary := manager.GetSummary()
		fmt.Printf("Artifact Summary:\n")
		fmt.Printf("  Total Artifacts: %d\n", summary["total_artifacts"])
		fmt.Printf("  Total Size: %s\n", formatBytes(summary["total_size"].(int64)))
		fmt.Printf("\n")

		fmt.Printf("By Type:\n")
		for t, count := range summary["by_type"].(map[string]int) {
			fmt.Printf("  %s: %d\n", t, count)
		}
		fmt.Printf("\n")

		fmt.Printf("By Category:\n")
		for cat, count := range summary["by_category"].(map[string]int) {
			fmt.Printf("  %s: %d\n", cat, count)
		}

	case "cleanup":
		days := 30
		dryRun := true

		// Parse options
		for i := 2; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "--days":
				if i+1 < len(os.Args) {
					fmt.Sscanf(os.Args[i+1], "%d", &days)
					i++
				}
			case "--execute":
				dryRun = false
			}
		}

		if err := manager.ScanArtifacts(); err != nil {
			fmt.Fprintf(os.Stderr, "Error scanning artifacts: %v\n", err)
			os.Exit(1)
		}

		toDelete, err := manager.CleanupArtifacts(days, dryRun)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error during cleanup: %v\n", err)
			os.Exit(1)
		}

		if dryRun {
			fmt.Printf("DRY RUN: Would delete %d artifacts older than %d days:\n", len(toDelete), days)
		} else {
			fmt.Printf("Deleted %d artifacts older than %d days:\n", len(toDelete), days)
		}

		for _, path := range toDelete {
			fmt.Printf("  %s\n", path)
		}

	case "export":
		outputPath := "artifact_metadata.json"
		if len(os.Args) > 2 {
			outputPath = os.Args[2]
		}

		if err := manager.ScanArtifacts(); err != nil {
			fmt.Fprintf(os.Stderr, "Error scanning artifacts: %v\n", err)
			os.Exit(1)
		}

		if err := manager.ExportMetadata(outputPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error exporting metadata: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Metadata exported to: %s\n", outputPath)

	default:
		showHelp()
		os.Exit(1)
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func showHelp() {
	fmt.Printf(`AcoustiCalc Artifact Manager

Usage: %s <command> [options]

Commands:
  scan                           Scan artifact directory
  list [--type T] [--category C] [--test T]  List artifacts with filters
  summary                        Show artifact summary
  cleanup [--days N] [--execute]  Clean up old artifacts (default: dry-run)
  export [filename]              Export metadata to JSON

Environment Variables:
  ARTIFACT_DIR                   Base artifact directory (default: tests/artifacts)

Examples:
  %s scan
  %s list --type screenshot --category unit
  %s summary
  %s cleanup --days 7 --execute
  %s export artifacts.json

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
