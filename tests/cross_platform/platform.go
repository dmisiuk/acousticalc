package cross_platform

import (
	"runtime"
	"strings"
)

// Platform returns the canonical name of the platform the tests are running on.
//
// The Story 0.2.3 acceptance criteria requires cross-platform validation and
// reporting. Centralising the platform string ensures consistent labelling
// across the E2E framework, recorder, and reporters.
func Platform() string {
	return runtime.GOOS
}

// NormalizeNewlines converts Windows style CRLF sequences into the newline
// representation used on Unix-like systems. This allows assertions written in
// a platform-agnostic fashion to succeed regardless of the host operating
// system the GitHub Actions job is running on.
func NormalizeNewlines(s string) string {
	if runtime.GOOS != "windows" {
		return s
	}
	return strings.ReplaceAll(s, "\r\n", "\n")
}

// PathWithExecutableSuffix appends the Windows executable suffix when running
// on Windows. Other platforms return the path unchanged.
func PathWithExecutableSuffix(path string) string {
	if runtime.GOOS == "windows" {
		return path + ".exe"
	}
	return path
}
