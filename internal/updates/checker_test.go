package updates

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewChecker(t *testing.T) {
	tmpDir := t.TempDir()

	config := UpdateConfig{
		Owner:          "jjuanrivvera",
		Repo:           "canvas-cli",
		CurrentVersion: "v1.0.0",
		ForceCheck:     false,
		CacheTTL:       6 * time.Hour,
	}

	checker := NewChecker(config, tmpDir)

	if checker == nil {
		t.Fatal("NewChecker returned nil")
	}

	if checker.config.Owner != "jjuanrivvera" {
		t.Errorf("Expected owner to be 'jjuanrivvera', got %s", checker.config.Owner)
	}

	if checker.config.CacheTTL != 6*time.Hour {
		t.Errorf("Expected CacheTTL to be 6h, got %v", checker.config.CacheTTL)
	}
}

func TestNewCheckerDefaultTTL(t *testing.T) {
	tmpDir := t.TempDir()

	config := UpdateConfig{
		Owner:          "jjuanrivvera",
		Repo:           "canvas-cli",
		CurrentVersion: "v1.0.0",
		// CacheTTL not set
	}

	checker := NewChecker(config, tmpDir)

	if checker.config.CacheTTL != defaultCacheTTL {
		t.Errorf("Expected default CacheTTL to be %v, got %v", defaultCacheTTL, checker.config.CacheTTL)
	}
}

func TestCheckDevVersion(t *testing.T) {
	tmpDir := t.TempDir()

	config := UpdateConfig{
		Owner:          "jjuanrivvera",
		Repo:           "canvas-cli",
		CurrentVersion: "dev",
		ForceCheck:     false,
	}

	checker := NewChecker(config, tmpDir)
	ctx := context.Background()

	result, err := checker.Check(ctx)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	if result.UpdateAvailable {
		t.Error("Development version should not have updates available")
	}

	if result.CurrentVersion != "dev" {
		t.Errorf("Expected current version to be 'dev', got %s", result.CurrentVersion)
	}
}

func TestCacheSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()

	config := UpdateConfig{
		Owner:          "jjuanrivvera",
		Repo:           "canvas-cli",
		CurrentVersion: "v1.0.0",
	}

	checker := NewChecker(config, tmpDir)

	// Create a test result
	result := &CheckResult{
		UpdateAvailable: true,
		CurrentVersion:  "v1.0.0",
		LatestVersion:   "v1.1.0",
		CheckedAt:       time.Now(),
		ReleaseInfo: &ReleaseInfo{
			Version:     "v1.1.0",
			URL:         "https://github.com/jjuanrivvera/canvas-cli/releases/tag/v1.1.0",
			ReleaseDate: time.Now(),
			Notes:       "Test release",
		},
	}

	// Save to cache
	if err := checker.saveCache(result); err != nil {
		t.Fatalf("saveCache failed: %v", err)
	}

	// Load from cache
	loaded, err := checker.loadCache()
	if err != nil {
		t.Fatalf("loadCache failed: %v", err)
	}

	if loaded == nil {
		t.Fatal("loadCache returned nil")
	}

	if loaded.UpdateAvailable != result.UpdateAvailable {
		t.Errorf("Expected UpdateAvailable to be %v, got %v", result.UpdateAvailable, loaded.UpdateAvailable)
	}

	if loaded.CurrentVersion != result.CurrentVersion {
		t.Errorf("Expected CurrentVersion to be %s, got %s", result.CurrentVersion, loaded.CurrentVersion)
	}

	if loaded.LatestVersion != result.LatestVersion {
		t.Errorf("Expected LatestVersion to be %s, got %s", result.LatestVersion, loaded.LatestVersion)
	}
}

func TestClearCache(t *testing.T) {
	tmpDir := t.TempDir()

	config := UpdateConfig{
		Owner:          "jjuanrivvera",
		Repo:           "canvas-cli",
		CurrentVersion: "v1.0.0",
	}

	checker := NewChecker(config, tmpDir)

	// Create a test result and save it
	result := &CheckResult{
		UpdateAvailable: true,
		CurrentVersion:  "v1.0.0",
		LatestVersion:   "v1.1.0",
		CheckedAt:       time.Now(),
	}

	if err := checker.saveCache(result); err != nil {
		t.Fatalf("saveCache failed: %v", err)
	}

	// Verify cache file exists
	cachePath := filepath.Join(tmpDir, cacheFileName)
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		t.Fatal("Cache file was not created")
	}

	// Clear cache
	if err := checker.ClearCache(); err != nil {
		t.Fatalf("ClearCache failed: %v", err)
	}

	// Verify cache file is removed
	if _, err := os.Stat(cachePath); !os.IsNotExist(err) {
		t.Error("Cache file was not removed")
	}
}

func TestLoadCacheNonExistent(t *testing.T) {
	tmpDir := t.TempDir()

	config := UpdateConfig{
		Owner:          "jjuanrivvera",
		Repo:           "canvas-cli",
		CurrentVersion: "v1.0.0",
	}

	checker := NewChecker(config, tmpDir)

	// Try to load from non-existent cache
	result, err := checker.loadCache()
	if err != nil {
		t.Errorf("loadCache should not return error for non-existent cache, got: %v", err)
	}

	if result != nil {
		t.Error("loadCache should return nil for non-existent cache")
	}
}
