package updates

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

// Installer handles installing updates
type Installer struct {
	config UpdateConfig
}

// NewInstaller creates a new update installer
func NewInstaller(config UpdateConfig) *Installer {
	return &Installer{
		config: config,
	}
}

// Install downloads and installs the latest version
func (i *Installer) Install(ctx context.Context) error {
	currentVersion := strings.TrimPrefix(i.config.CurrentVersion, "v")
	if currentVersion == "" || currentVersion == "dev" || currentVersion == "unknown" {
		return fmt.Errorf("cannot update from development version, please install a released version first")
	}

	// Detect latest release
	slug := i.config.Owner + "/" + i.config.Repo
	latest, found, err := selfupdate.DetectLatest(slug)
	if err != nil {
		return fmt.Errorf("failed to detect latest version: %w", err)
	}

	if !found {
		return fmt.Errorf("no releases found for %s/%s", i.config.Owner, i.config.Repo)
	}

	// Get current executable path
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Check if update is needed
	if latest.Version.String() != currentVersion {
		fmt.Printf("Downloading version %s...\n", latest.Version.String())

		// Perform the update
		if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
			return fmt.Errorf("failed to update binary: %w", err)
		}

		fmt.Printf("Successfully updated to version %s\n", latest.Version.String())
		fmt.Println("Restart the application to use the new version")
		return nil
	}

	fmt.Printf("Already running the latest version (%s)\n", currentVersion)
	return nil
}

// CanUpdate checks if the current binary can be updated
// (e.g., not installed via package manager)
func (i *Installer) CanUpdate() (bool, string) {
	exe, err := os.Executable()
	if err != nil {
		return false, "cannot determine executable path"
	}

	// Check if we have write permission to the executable
	info, err := os.Stat(exe)
	if err != nil {
		return false, "cannot stat executable"
	}

	// On Unix, check if we can write to the file
	if info.Mode().Perm()&0200 == 0 {
		return false, "no write permission to executable"
	}

	// Check for common package manager installation paths
	// These typically shouldn't be updated via self-update
	if strings.Contains(exe, "/usr/local/Cellar/") {
		return false, "installed via Homebrew, use 'brew upgrade canvas-cli' instead"
	}

	if strings.Contains(exe, "/snap/") {
		return false, "installed via Snap, use snap refresh instead"
	}

	if strings.Contains(exe, "/flatpak/") {
		return false, "installed via Flatpak, use flatpak update instead"
	}

	return true, ""
}
