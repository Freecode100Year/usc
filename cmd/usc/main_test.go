package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func runCLI(t *testing.T, args ...string) (string, int) {
	t.Helper()
	exePath := filepath.Join("..", "..", "bin", "usc.exe")
	if _, err := os.Stat(exePath); err != nil {
		exePath = "usc"
	}
	cmd := exec.Command(exePath, args...)
	out, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}
	return string(out), exitCode
}

func TestCLIHelpAndUnknown(t *testing.T) {
	out, code := runCLI(t, "--help")
	if code != 0 || !strings.Contains(out, "USC (Universal Skill Compiler)") {
		t.Fatalf("expected code 0 with usage, got code %d, out: %s", code, out)
	}

	out, code = runCLI(t, "unknown-command-test")
	if code != 1 || !strings.Contains(out, "Error: Unknown command") {
		t.Fatalf("expected code 1 for unknown command, got %d, out: %s", code, out)
	}
}

func TestCLIVerifyNonExistent(t *testing.T) {
	out, code := runCLI(t, "verify", "definitely_non_existent.usc")
	if code != 1 {
		t.Fatalf("expected exit code 1 for non-existent verify, got %d, out: %s", code, out)
	}
	if !strings.Contains(out, "VERIFICATION REJECTED") {
		t.Fatalf("expected rejection output, got: %s", out)
	}
}

func TestCLITargets(t *testing.T) {
	out, code := runCLI(t, "targets")
	if code != 0 || !strings.Contains(out, "SUPPORTED AGENT RUNTIMES") {
		t.Fatalf("expected code 0 for targets, got %d, out: %s", code, out)
	}
}

func TestCLITop(t *testing.T) {
	out, code := runCLI(t, "top")
	if code != 0 || !strings.Contains(out, "USC ZERO-TRUST SECURITY DASHBOARD") {
		t.Fatalf("expected code 0 for top, got %d, out: %s", code, out)
	}
}

func TestCLIRunNonExistent(t *testing.T) {
	out, code := runCLI(t, "run", "missing_artifact.usc")
	if code != 1 || !strings.Contains(out, "Refusing execution") {
		t.Fatalf("expected code 1 for missing artifact run, got %d, out: %s", code, out)
	}
}

func TestCLIExportNonExistent(t *testing.T) {
	out, code := runCLI(t, "export", "missing_artifact.usc", "--target", "openclaw")
	if code != 1 || !strings.Contains(out, "Export rejected") {
		t.Fatalf("expected code 1 for missing artifact export, got %d, out: %s", code, out)
	}
}

func TestCLIAnalyzeMissing(t *testing.T) {
	out, code := runCLI(t, "analyze", "missing_dir_path")
	if code != 1 || !strings.Contains(out, "Source directory not found") {
		t.Fatalf("expected code 1 for missing analyze dir, got %d, out: %s", code, out)
	}
}

func TestCLIReplayMissing(t *testing.T) {
	out, code := runCLI(t, "replay", "missing_trace.usctrace")
	if code != 1 || !strings.Contains(out, "trace file not found") {
		t.Fatalf("expected code 1 for missing replay trace, got %d, out: %s", code, out)
	}
}


