package updates

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

const (
	defaultCacheTTL = 6 * time.Hour
	cacheFileName   = "update_check.json"
)

// Checker checks for new versions of the application
type Checker struct {
	config    UpdateConfig
	cacheDir  string
	cachePath string
}

// NewChecker creates a new update checker
func NewChecker(config UpdateConfig, cacheDir string) *Checker {
	if config.CacheTTL == 0 {
		config.CacheTTL = defaultCacheTTL
	}

	cachePath := filepath.Join(cacheDir, cacheFileName)

	return &Checker{
		config:    config,
		cacheDir:  cacheDir,
		cachePath: cachePath,
	}
}

// Check checks for updates and returns the result
func (c *Checker) Check(ctx context.Context) (*CheckResult, error) {
	// Try to load from cache first if not forcing
	if !c.config.ForceCheck {
		if cached, err := c.loadCache(); err == nil && cached != nil {
			if time.Since(cached.CheckedAt) < c.config.CacheTTL {
				return cached, nil
			}
		}
	}

	// Perform actual check
	result, err := c.checkGitHub(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}

	// Save to cache
	if err := c.saveCache(result); err != nil {
		// Log but don't fail on cache write errors
		fmt.Fprintf(os.Stderr, "Warning: failed to cache update check: %v\n", err)
	}

	return result, nil
}

// checkGitHub checks GitHub releases for updates
func (c *Checker) checkGitHub(ctx context.Context) (*CheckResult, error) {
	currentVersion := strings.TrimPrefix(c.config.CurrentVersion, "v")
	if currentVersion == "" || currentVersion == "dev" || currentVersion == "unknown" {
		return &CheckResult{
			UpdateAvailable: false,
			CurrentVersion:  c.config.CurrentVersion,
			LatestVersion:   c.config.CurrentVersion,
			CheckedAt:       time.Now(),
		}, nil
	}

	// Parse current version
	current, err := semver.NewVersion(currentVersion)
	if err != nil {
		return nil, fmt.Errorf("invalid current version %q: %w", currentVersion, err)
	}

	// Check GitHub for latest release
	slug := c.config.Owner + "/" + c.config.Repo
	latest, found, err := selfupdate.DetectLatest(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to detect latest version: %w", err)
	}

	if !found {
		return nil, fmt.Errorf("no releases found for %s/%s", c.config.Owner, c.config.Repo)
	}

	result := &CheckResult{
		CurrentVersion: c.config.CurrentVersion,
		LatestVersion:  latest.Version.String(),
		CheckedAt:      time.Now(),
	}

	// Parse latest version
	latestVer, err := semver.NewVersion(latest.Version.String())
	if err != nil {
		return nil, fmt.Errorf("invalid latest version %q: %w", latest.Version.String(), err)
	}

	// Compare versions
	if latestVer.GreaterThan(current) {
		result.UpdateAvailable = true
		releaseDate := time.Now()
		if latest.PublishedAt != nil {
			releaseDate = *latest.PublishedAt
		}
		result.ReleaseInfo = &ReleaseInfo{
			Version:     latest.Version.String(),
			URL:         latest.URL,
			ReleaseDate: releaseDate,
			Notes:       latest.ReleaseNotes,
			AssetURL:    latest.AssetURL,
			AssetName:   latest.Name,
		}
	}

	return result, nil
}

// loadCache loads the cached check result
func (c *Checker) loadCache() (*CheckResult, error) {
	data, err := os.ReadFile(c.cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read cache: %w", err)
	}

	var result CheckResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse cache: %w", err)
	}

	return &result, nil
}

// saveCache saves the check result to cache
func (c *Checker) saveCache(result *CheckResult) error {
	// Ensure cache directory exists
	if err := os.MkdirAll(c.cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	if err := os.WriteFile(c.cachePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache: %w", err)
	}

	return nil
}

// ClearCache removes the cached check result
func (c *Checker) ClearCache() error {
	if err := os.Remove(c.cachePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to clear cache: %w", err)
	}
	return nil
}
