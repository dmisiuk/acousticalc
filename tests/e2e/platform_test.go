package e2e

import (
	"os/exec"
	"runtime"
	"testing"
)

func TestE2EPlatform(t *testing.T) {
	t.Run("TestPlatformSpecificBehavior", func(t *testing.T) {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "linux":
			t.Log("Testing on Linux")
			cmd = exec.Command("uname", "-a")
		case "darwin":
			t.Log("Testing on macOS")
			cmd = exec.Command("sw_vers")
		case "windows":
			t.Log("Testing on Windows")
			cmd = exec.Command("ver")
		default:
			t.Skipf("Skipping platform-specific test on %s", runtime.GOOS)
		}

		if cmd != nil {
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Failed to execute platform-specific command: %v\nOutput: %s", err, out)
			}
			t.Logf("Platform command output:\n%s", out)
		}
	})
}
