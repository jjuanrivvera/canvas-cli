package update

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		name     string
		latest   string
		current  string
		expected bool
	}{
		{"newer major", "2.0.0", "1.0.0", true},
		{"newer minor", "1.1.0", "1.0.0", true},
		{"newer patch", "1.0.1", "1.0.0", true},
		{"same version", "1.0.0", "1.0.0", false},
		{"older major", "1.0.0", "2.0.0", false},
		{"older minor", "1.0.0", "1.1.0", false},
		{"older patch", "1.0.0", "1.0.1", false},
		{"with v prefix", "v1.1.0", "v1.0.0", true},
		{"mixed prefix", "1.1.0", "v1.0.0", true},
		{"with prerelease", "1.1.0-beta", "1.0.0", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNewerVersion(tt.latest, tt.current)
			if result != tt.expected {
				t.Errorf("isNewerVersion(%q, %q) = %v, expected %v", tt.latest, tt.current, result, tt.expected)
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected [3]int
	}{
		{"1.2.3", [3]int{1, 2, 3}},
		{"v1.2.3", [3]int{1, 2, 3}},
		{"1.2", [3]int{1, 2, 0}},
		{"1", [3]int{1, 0, 0}},
		{"1.2.3-beta", [3]int{1, 2, 3}},
		{"1.2.3-rc.1", [3]int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseVersion(tt.input)
			if result != tt.expected {
				t.Errorf("parseVersion(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetLatestRelease(t *testing.T) {
	// Create a mock server
	release := Release{
		TagName:    "v1.2.3",
		Name:       "v1.2.3",
		Draft:      false,
		Prerelease: false,
		Assets: []Asset{
			{
				Name:               "canvas-cli_darwin_x86_64.tar.gz",
				BrowserDownloadURL: "https://example.com/canvas-cli_darwin_x86_64.tar.gz",
			},
			{
				Name:               "checksums.txt",
				BrowserDownloadURL: "https://example.com/checksums.txt",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	// Create updater with custom HTTP client pointing to mock server
	updater := NewUpdater("1.0.0")
	updater.HTTPClient = server.Client()

	// Override the URL by using a custom transport
	originalURL := server.URL
	updater.HTTPClient.Transport = &urlRewriteTransport{
		targetURL: originalURL,
	}

	ctx := context.Background()
	result, err := updater.GetLatestRelease(ctx)
	if err != nil {
		t.Fatalf("GetLatestRelease failed: %v", err)
	}

	if result.TagName != "v1.2.3" {
		t.Errorf("Expected tag v1.2.3, got %s", result.TagName)
	}
}

// urlRewriteTransport rewrites requests to the test server
type urlRewriteTransport struct {
	targetURL string
}

func (t *urlRewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	req.URL.Host = t.targetURL[7:] // Remove "http://"
	return http.DefaultTransport.RoundTrip(req)
}

func TestFindAssets(t *testing.T) {
	release := &Release{
		Assets: []Asset{
			{Name: "canvas-cli_darwin_x86_64.tar.gz", BrowserDownloadURL: "https://example.com/darwin_amd64.tar.gz"},
			{Name: "canvas-cli_darwin_arm64.tar.gz", BrowserDownloadURL: "https://example.com/darwin_arm64.tar.gz"},
			{Name: "canvas-cli_linux_x86_64.tar.gz", BrowserDownloadURL: "https://example.com/linux_amd64.tar.gz"},
			{Name: "canvas-cli_windows_x86_64.zip", BrowserDownloadURL: "https://example.com/windows_amd64.zip"},
			{Name: "checksums.txt", BrowserDownloadURL: "https://example.com/checksums.txt"},
		},
	}

	updater := NewUpdater("1.0.0")
	binaryAsset, checksumAsset := updater.findAssets(release)

	if checksumAsset == nil {
		t.Error("Expected to find checksum asset")
	}
	if checksumAsset != nil && checksumAsset.Name != "checksums.txt" {
		t.Errorf("Expected checksums.txt, got %s", checksumAsset.Name)
	}

	// Binary asset depends on current OS/arch
	if binaryAsset != nil {
		t.Logf("Found binary asset for %s/%s: %s", runtime.GOOS, runtime.GOARCH, binaryAsset.Name)
	}
}

func TestVerifyChecksum(t *testing.T) {
	testData := []byte("test binary content")
	hash := sha256.Sum256(testData)
	expectedHash := hex.EncodeToString(hash[:])

	checksums := map[string]string{
		"test.tar.gz": expectedHash,
	}

	updater := NewUpdater("1.0.0")

	// Valid checksum
	if !updater.verifyChecksum(testData, "test.tar.gz", checksums) {
		t.Error("Expected checksum verification to pass")
	}

	// Wrong file
	if updater.verifyChecksum(testData, "other.tar.gz", checksums) {
		t.Error("Expected checksum verification to fail for unknown file")
	}

	// Wrong data
	if updater.verifyChecksum([]byte("wrong data"), "test.tar.gz", checksums) {
		t.Error("Expected checksum verification to fail for wrong data")
	}
}

func TestExtractFromTarGz(t *testing.T) {
	// Create a test tar.gz archive
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzWriter)

	binaryName := "canvas"
	if runtime.GOOS == "windows" {
		binaryName = "canvas.exe"
	}

	// Add the binary file
	binaryContent := []byte("#!/bin/bash\necho hello")
	hdr := &tar.Header{
		Name: binaryName,
		Mode: 0755,
		Size: int64(len(binaryContent)),
	}
	if err := tarWriter.WriteHeader(hdr); err != nil {
		t.Fatal(err)
	}
	if _, err := tarWriter.Write(binaryContent); err != nil {
		t.Fatal(err)
	}

	tarWriter.Close()
	gzWriter.Close()

	updater := NewUpdater("1.0.0")
	extracted, err := updater.extractFromTarGz(buf.Bytes())
	if err != nil {
		t.Fatalf("extractFromTarGz failed: %v", err)
	}

	if !bytes.Equal(extracted, binaryContent) {
		t.Error("Extracted content doesn't match original")
	}
}

func TestApplyUpdate(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "update-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a fake "current" binary
	currentBinary := filepath.Join(tmpDir, "canvas")
	if err := os.WriteFile(currentBinary, []byte("old version"), 0755); err != nil {
		t.Fatal(err)
	}

	// New binary content
	newContent := []byte("new version")

	updater := NewUpdater("1.0.0")
	updater.ExecutablePath = currentBinary

	if err := updater.applyUpdate(newContent); err != nil {
		t.Fatalf("applyUpdate failed: %v", err)
	}

	// Verify the new content
	content, err := os.ReadFile(currentBinary)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(content, newContent) {
		t.Errorf("Expected new content, got: %s", content)
	}

	// Verify backup was removed
	if _, err := os.Stat(currentBinary + ".bak"); !os.IsNotExist(err) {
		t.Error("Backup file should have been removed")
	}
}

func TestCheckAndUpdateSkipsDevVersion(t *testing.T) {
	updater := NewUpdater("dev")

	ctx := context.Background()
	result := updater.CheckAndUpdate(ctx)

	if result.Updated {
		t.Error("Should not update dev version")
	}
	if result.Error != nil {
		t.Errorf("Should not error for dev version: %v", result.Error)
	}
}
