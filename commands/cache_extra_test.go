package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// setupCacheTestHome redirects HOME to a temp dir and resets config cache.
func setupCacheTestHome(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("USERPROFILE", tmpDir) // Windows: os.UserHomeDir() reads %USERPROFILE%
	return tmpDir
}

// TestCacheStatsCmd_NoCacheDir verifies stats output when no cache dir exists.
func TestCacheStatsCmd_NoCacheDir(t *testing.T) {
	setupCacheTestHome(t)

	var output strings.Builder
	cmd := cacheStatsCmd
	cmd.SetOut(&output)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runCacheStats(cmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	r.Close()

	out := string(buf[:n])

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "empty") {
		t.Errorf("expected 'empty' in output, got: %q", out)
	}
}

// TestCacheClearCmd_NoCacheDir verifies clear output when no cache dir exists.
func TestCacheClearCmd_NoCacheDir(t *testing.T) {
	setupCacheTestHome(t)

	var buf strings.Builder
	cmd := cacheClearCmd
	cmd.SetOut(&buf)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runCacheClear(cmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	outBuf := make([]byte, 4096)
	n, _ := r.Read(outBuf)
	r.Close()

	out := string(outBuf[:n])

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "already empty") {
		t.Errorf("expected 'already empty' in output, got: %q", out)
	}
}

// TestFormatBytes tests the formatBytes helper.
func TestFormatBytes(t *testing.T) {
	cases := []struct {
		input    int64
		expected string
	}{
		{0, "0 bytes"},
		{100, "100 bytes"},
		{1024, "1.00 KB"},
		{2 * 1024, "2.00 KB"},
		{1024 * 1024, "1.00 MB"},
		{3 * 1024 * 1024, "3.00 MB"},
		{1024 * 1024 * 1024, "1.00 GB"},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%d_bytes", tc.input), func(t *testing.T) {
			got := formatBytes(tc.input)
			if got != tc.expected {
				t.Errorf("formatBytes(%d) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

// TestIsExpiredCacheFile tests expiration detection helpers.
func TestIsExpiredCacheFile(t *testing.T) {
	expired := fmt.Sprintf(`{"expiration": "%s"}`, time.Now().Add(-time.Hour).Format(time.RFC3339))
	active := fmt.Sprintf(`{"expiration": "%s"}`, time.Now().Add(time.Hour).Format(time.RFC3339))
	invalid := `not-json`

	if !isExpiredCacheFile([]byte(expired)) {
		t.Error("expected expired entry to be detected as expired")
	}
	if isExpiredCacheFile([]byte(active)) {
		t.Error("expected active entry not to be detected as expired")
	}
	if isExpiredCacheFile([]byte(invalid)) {
		t.Error("expected invalid JSON to return false")
	}
}

// TestGetCacheDirSize tests directory size calculation.
func TestGetCacheDirSize(t *testing.T) {
	tmpDir := t.TempDir()

	// Write a file with known content
	content := []byte("hello world")
	path := filepath.Join(tmpDir, "test.json")
	if err := os.WriteFile(path, content, 0600); err != nil {
		t.Fatal(err)
	}

	size, err := getCacheDirSize(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if size != int64(len(content)) {
		t.Errorf("expected size %d, got %d", len(content), size)
	}
}

// TestClearExpiredEntries verifies expired entries are removed and active ones kept.
func TestClearExpiredEntries(t *testing.T) {
	tmpDir := t.TempDir()

	expiredContent := fmt.Sprintf(`{"expiration": "%s", "value": "old"}`, time.Now().Add(-time.Hour).Format(time.RFC3339))
	activeContent := fmt.Sprintf(`{"expiration": "%s", "value": "new"}`, time.Now().Add(time.Hour).Format(time.RFC3339))

	expiredPath := filepath.Join(tmpDir, "expired.json")
	activePath := filepath.Join(tmpDir, "active.json")

	_ = os.WriteFile(expiredPath, []byte(expiredContent), 0600)
	_ = os.WriteFile(activePath, []byte(activeContent), 0600)

	if err := clearExpiredEntries(tmpDir); err != nil {
		t.Fatalf("clearExpiredEntries error: %v", err)
	}

	if _, err := os.Stat(expiredPath); !os.IsNotExist(err) {
		t.Error("expected expired entry to be removed")
	}
	if _, err := os.Stat(activePath); os.IsNotExist(err) {
		t.Error("expected active entry to still exist")
	}
}

// TestCacheStatsCmd_WithEntries exercises stats when entries exist.
func TestCacheStatsCmd_WithEntries(t *testing.T) {
	tmpDir := setupCacheTestHome(t)

	// Create a cache dir with one active file so DiskCache can be opened.
	cacheDir := filepath.Join(tmpDir, ".canvas-cli", "cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Write an active cache file in the expected DiskCache format.
	active := fmt.Sprintf(`{"key":"k1","value":"v","expiration":"%s"}`, time.Now().Add(time.Hour).Format(time.RFC3339))
	_ = os.WriteFile(filepath.Join(cacheDir, "entry1.json"), []byte(active), 0600)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runCacheStats(cacheStatsCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	buf := make([]byte, 8192)
	n, _ := r.Read(buf)
	r.Close()

	out := string(buf[:n])

	if err != nil {
		t.Fatalf("expected no error, got: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "Cache Statistics") {
		t.Errorf("expected 'Cache Statistics' in output, got: %q", out)
	}
}
