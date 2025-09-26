package e2e

import (
	"os/exec"
	"testing"
	"time"
)

func TestE2EPerformance(t *testing.T) {
	t.Run("TestStartupTime", func(t *testing.T) {
		maxStartupTime := 2 * time.Second
		startTime := time.Now()

		cmd := exec.Command("go", "run", "main.go", "--version")
		err := cmd.Run()
		if err != nil {
			t.Fatalf("Failed to execute command: %v", err)
		}

		startupTime := time.Since(startTime)
		if startupTime > maxStartupTime {
			t.Errorf("Startup time exceeded limit: got %v, want <= %v", startupTime, maxStartupTime)
		}
		t.Logf("Startup time: %v", startupTime)
	})

	t.Run("TestCalculationLatency", func(t *testing.T) {
		maxLatency := 500 * time.Millisecond
		expression := "sqrt(144) * (2 + 3)"
		startTime := time.Now()

		cmd := exec.Command("go", "run", "main.go", "-e", expression)
		err := cmd.Run()
		if err != nil {
			t.Fatalf("Failed to execute calculation: %v", err)
		}

		latency := time.Since(startTime)
		if latency > maxLatency {
			t.Errorf("Calculation latency exceeded limit: got %v, want <= %v", latency, maxLatency)
		}
		t.Logf("Calculation latency for '%s': %v", expression, latency)
	})
}