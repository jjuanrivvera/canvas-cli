package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

// setupTelemetryTestHome creates an isolated HOME dir and resets the config cache.
func setupTelemetryTestHome(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("USERPROFILE", tmpDir) // Windows: os.UserHomeDir() reads %USERPROFILE%
	config.ResetCache()
	t.Cleanup(config.ResetCache)
	return tmpDir
}

func captureRunOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	buf := make([]byte, 16384)
	n, _ := r.Read(buf)
	r.Close()
	return string(buf[:n])
}

// TestTelemetryEnableCmd verifies enable writes to config and prints confirmation.
func TestTelemetryEnableCmd(t *testing.T) {
	setupTelemetryTestHome(t)

	out := captureRunOutput(func() {
		err := runTelemetryEnable(telemetryEnableCmd, []string{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "enabled") {
		t.Errorf("expected 'enabled' in output, got: %q", out)
	}
}

// TestTelemetryDisableCmd verifies disable writes to config and prints confirmation.
func TestTelemetryDisableCmd(t *testing.T) {
	setupTelemetryTestHome(t)

	out := captureRunOutput(func() {
		err := runTelemetryDisable(telemetryDisableCmd, []string{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "disabled") {
		t.Errorf("expected 'disabled' in output, got: %q", out)
	}
}

// TestTelemetryStatusCmd_Disabled verifies status when telemetry is off.
func TestTelemetryStatusCmd_Disabled(t *testing.T) {
	setupTelemetryTestHome(t)

	out := captureRunOutput(func() {
		err := runTelemetryStatus(telemetryStatusCmd, []string{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "Telemetry Status") {
		t.Errorf("expected 'Telemetry Status' header in output, got: %q", out)
	}
}

// TestTelemetryStatusCmd_Enabled verifies status shows enabled when configured.
func TestTelemetryStatusCmd_Enabled(t *testing.T) {
	setupTelemetryTestHome(t)

	// Enable telemetry first.
	_ = runTelemetryEnable(telemetryEnableCmd, []string{})
	config.ResetCache()

	out := captureRunOutput(func() {
		err := runTelemetryStatus(telemetryStatusCmd, []string{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "Enabled") {
		t.Errorf("expected 'Enabled' in output, got: %q", out)
	}
}

// TestTelemetryShowCmd_NoData verifies show output when no telemetry data exists.
func TestTelemetryShowCmd_NoData(t *testing.T) {
	setupTelemetryTestHome(t)

	out := captureRunOutput(func() {
		err := runTelemetryShow(telemetryShowCmd, []string{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "No telemetry data") {
		t.Errorf("expected 'No telemetry data' in output, got: %q", out)
	}
}

// TestTelemetryShowCmd_WithData verifies show lists files when data exists.
func TestTelemetryShowCmd_WithData(t *testing.T) {
	tmpDir := setupTelemetryTestHome(t)

	// Create a fake telemetry data file.
	telDir := filepath.Join(tmpDir, ".canvas-cli", "telemetry")
	if err := os.MkdirAll(telDir, 0755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(telDir, "event1.json"), []byte(`{"event":"test"}`), 0600)

	out := captureRunOutput(func() {
		err := runTelemetryShow(telemetryShowCmd, []string{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "event1.json") {
		t.Errorf("expected filename in output, got: %q", out)
	}
}

// TestTelemetryClearCmd_NoData verifies clear output when no data exists.
func TestTelemetryClearCmd_NoData(t *testing.T) {
	setupTelemetryTestHome(t)

	out := captureRunOutput(func() {
		err := runTelemetryClear(telemetryClearCmd, []string{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "No telemetry data") {
		t.Errorf("expected 'No telemetry data' in output, got: %q", out)
	}
}

// TestTelemetryClearCmd_WithData verifies that JSON files are removed.
func TestTelemetryClearCmd_WithData(t *testing.T) {
	tmpDir := setupTelemetryTestHome(t)

	telDir := filepath.Join(tmpDir, ".canvas-cli", "telemetry")
	if err := os.MkdirAll(telDir, 0755); err != nil {
		t.Fatal(err)
	}
	filePath := filepath.Join(telDir, "event1.json")
	_ = os.WriteFile(filePath, []byte(`{"event":"test"}`), 0600)

	out := captureRunOutput(func() {
		err := runTelemetryClear(telemetryClearCmd, []string{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "Cleared") {
		t.Errorf("expected 'Cleared' in output, got: %q", out)
	}

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("expected telemetry file to be removed")
	}
}
