// Package update provides auto-update functionality for the CLI.
package update

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	// GitHubOwner is the repository owner
	GitHubOwner = "jjuanrivvera"
	// GitHubRepo is the repository name
	GitHubRepo = "canvas-cli"
	// BinaryName is the name of the binary to update
	BinaryName = "canvas"
	// DefaultCheckInterval is the default interval between update checks
	DefaultCheckInterval = 1 * time.Hour
)

// Release represents a GitHub release
type Release struct {
	TagName    string  `json:"tag_name"`
	Name       string  `json:"name"`
	Draft      bool    `json:"draft"`
	Prerelease bool    `json:"prerelease"`
	Assets     []Asset `json:"assets"`
	Body       string  `json:"body"`
}

// Asset represents a release asset
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// UpdateResult contains the result of an update operation
type UpdateResult struct {
	Updated     bool
	FromVersion string
	ToVersion   string
	Error       error
}

// Updater handles checking and applying updates
type Updater struct {
	CurrentVersion string
	HTTPClient     *http.Client
	// For testing - allows overriding the executable path
	ExecutablePath string
}

// NewUpdater creates a new Updater instance
func NewUpdater(currentVersion string) *Updater {
	return &Updater{
		CurrentVersion: currentVersion,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CheckAndUpdate checks for updates and applies them if available
func (u *Updater) CheckAndUpdate(ctx context.Context) *UpdateResult {
	result := &UpdateResult{
		FromVersion: u.CurrentVersion,
	}

	// Skip if running dev version
	if u.CurrentVersion == "dev" || u.CurrentVersion == "" {
		return result
	}

	// Get latest release
	release, err := u.GetLatestRelease(ctx)
	if err != nil {
		result.Error = fmt.Errorf("failed to check for updates: %w", err)
		return result
	}

	// Compare versions
	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := strings.TrimPrefix(u.CurrentVersion, "v")

	if !isNewerVersion(latestVersion, currentVersion) {
		return result // No update needed
	}

	result.ToVersion = latestVersion

	// Find the appropriate asset
	asset, checksumAsset := u.findAssets(release)
	if asset == nil {
		result.Error = fmt.Errorf("no compatible binary found for %s/%s", runtime.GOOS, runtime.GOARCH)
		return result
	}

	// Download the binary
	binary, err := u.downloadAsset(ctx, asset)
	if err != nil {
		result.Error = fmt.Errorf("failed to download update: %w", err)
		return result
	}

	// Verify checksum if available
	if checksumAsset != nil {
		checksums, err := u.downloadChecksums(ctx, checksumAsset)
		if err != nil {
			result.Error = fmt.Errorf("failed to download checksums: %w", err)
			return result
		}

		if !u.verifyChecksum(binary, asset.Name, checksums) {
			result.Error = fmt.Errorf("checksum verification failed")
			return result
		}
	}

	// Extract binary from archive
	extractedBinary, err := u.extractBinary(binary, asset.Name)
	if err != nil {
		result.Error = fmt.Errorf("failed to extract binary: %w", err)
		return result
	}

	// Apply the update
	if err := u.applyUpdate(extractedBinary); err != nil {
		result.Error = fmt.Errorf("failed to apply update: %w", err)
		return result
	}

	result.Updated = true
	return result
}

// GetLatestRelease fetches the latest release from GitHub
func (u *Updater) GetLatestRelease(ctx context.Context) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", GitHubOwner, GitHubRepo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "canvas-cli-updater")

	resp, err := u.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// findAssets finds the appropriate binary and checksum assets for the current platform
func (u *Updater) findAssets(release *Release) (*Asset, *Asset) {
	// Map GOARCH to archive naming convention
	archName := runtime.GOARCH
	if archName == "amd64" {
		archName = "x86_64"
	}

	// Determine expected archive name pattern
	// Format: canvas-cli_<os>_<arch>.tar.gz (or .zip for windows)
	ext := ".tar.gz"
	if runtime.GOOS == "windows" {
		ext = ".zip"
	}

	expectedName := fmt.Sprintf("canvas-cli_%s_%s%s", runtime.GOOS, archName, ext)

	var binaryAsset, checksumAsset *Asset

	for i := range release.Assets {
		asset := &release.Assets[i]
		if asset.Name == expectedName {
			binaryAsset = asset
		}
		if asset.Name == "checksums.txt" {
			checksumAsset = asset
		}
	}

	return binaryAsset, checksumAsset
}

// downloadAsset downloads a release asset
func (u *Updater) downloadAsset(ctx context.Context, asset *Asset) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", asset.BrowserDownloadURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "canvas-cli-updater")

	resp, err := u.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download asset: status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// downloadChecksums downloads and parses the checksums file
func (u *Updater) downloadChecksums(ctx context.Context, asset *Asset) (map[string]string, error) {
	data, err := u.downloadAsset(ctx, asset)
	if err != nil {
		return nil, err
	}

	checksums := make(map[string]string)
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Format: <checksum>  <filename>
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			checksums[parts[1]] = parts[0]
		}
	}

	return checksums, nil
}

// verifyChecksum verifies the downloaded binary against the checksum
func (u *Updater) verifyChecksum(data []byte, filename string, checksums map[string]string) bool {
	expectedHash, ok := checksums[filename]
	if !ok {
		return false
	}

	hash := sha256.Sum256(data)
	actualHash := hex.EncodeToString(hash[:])

	return strings.EqualFold(expectedHash, actualHash)
}

// extractBinary extracts the binary from the archive
func (u *Updater) extractBinary(archive []byte, archiveName string) ([]byte, error) {
	if strings.HasSuffix(archiveName, ".zip") {
		return u.extractFromZip(archive)
	}
	return u.extractFromTarGz(archive)
}

// extractFromTarGz extracts the binary from a tar.gz archive
func (u *Updater) extractFromTarGz(archive []byte) ([]byte, error) {
	gzReader, err := gzip.NewReader(bytes.NewReader(archive))
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	binaryName := BinaryName
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Check if this is the binary we're looking for
		if filepath.Base(header.Name) == binaryName && header.Typeflag == tar.TypeReg {
			return io.ReadAll(tarReader)
		}
	}

	return nil, fmt.Errorf("binary %s not found in archive", binaryName)
}

// extractFromZip extracts the binary from a zip archive
func (u *Updater) extractFromZip(archive []byte) ([]byte, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	if err != nil {
		return nil, err
	}

	binaryName := BinaryName
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	for _, file := range zipReader.File {
		if filepath.Base(file.Name) == binaryName {
			rc, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}

	return nil, fmt.Errorf("binary %s not found in archive", binaryName)
}

// applyUpdate replaces the current binary with the new one
func (u *Updater) applyUpdate(newBinary []byte) error {
	execPath := u.ExecutablePath
	if execPath == "" {
		var err error
		execPath, err = os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %w", err)
		}

		// Resolve symlinks to get the actual binary path
		execPath, err = filepath.EvalSymlinks(execPath)
		if err != nil {
			return fmt.Errorf("failed to resolve executable path: %w", err)
		}
	}

	// Create a temporary file in the same directory (for atomic rename)
	execDir := filepath.Dir(execPath)
	tmpFile, err := os.CreateTemp(execDir, "canvas-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Write the new binary
	if _, err := tmpFile.Write(newBinary); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("failed to write new binary: %w", err)
	}
	tmpFile.Close()

	// Make it executable
	if err := os.Chmod(tmpPath, 0755); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// Backup the old binary
	backupPath := execPath + ".bak"
	if err := os.Rename(execPath, backupPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to backup old binary: %w", err)
	}

	// Move new binary into place
	if err := os.Rename(tmpPath, execPath); err != nil {
		// Try to restore the backup (ignore restore error - best effort)
		_ = os.Rename(backupPath, execPath)
		os.Remove(tmpPath)
		return fmt.Errorf("failed to install new binary: %w", err)
	}

	// Remove backup
	os.Remove(backupPath)

	return nil
}

// isNewerVersion compares two semver versions
// Returns true if latest is newer than current
func isNewerVersion(latest, current string) bool {
	latestParts := parseVersion(latest)
	currentParts := parseVersion(current)

	for i := 0; i < 3; i++ {
		if latestParts[i] > currentParts[i] {
			return true
		}
		if latestParts[i] < currentParts[i] {
			return false
		}
	}

	return false
}

// parseVersion parses a semver string into [major, minor, patch]
func parseVersion(v string) [3]int {
	v = strings.TrimPrefix(v, "v")
	parts := strings.Split(v, ".")

	var result [3]int
	for i := 0; i < 3 && i < len(parts); i++ {
		// Strip any pre-release suffix
		numStr := strings.Split(parts[i], "-")[0]
		fmt.Sscanf(numStr, "%d", &result[i])
	}

	return result
}
