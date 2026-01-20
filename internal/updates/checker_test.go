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

func TestCheckUnknownVersion(t *testing.T) {
	tmpDir := t.TempDir()

	config := UpdateConfig{
		Owner:          "jjuanrivvera",
		Repo:           "canvas-cli",
		CurrentVersion: "unknown",
		ForceCheck:     false,
	}

	checker := NewChecker(config, tmpDir)
	ctx := context.Background()

	result, err := checker.Check(ctx)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	if result.UpdateAvailable {
		t.Error("Unknown version should not have updates available")
	}

	if result.CurrentVersion != "unknown" {
		t.Errorf("Expected current version to be 'unknown', got %s", result.CurrentVersion)
	}
}

func TestCheckEmptyVersion(t *testing.T) {
	tmpDir := t.TempDir()

	config := UpdateConfig{
		Owner:          "jjuanrivvera",
		Repo:           "canvas-cli",
		CurrentVersion: "",
		ForceCheck:     false,
	}

	checker := NewChecker(config, tmpDir)
	ctx := context.Background()

	result, err := checker.Check(ctx)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	if result.UpdateAvailable {
		t.Error("Empty version should not have updates available")
	}
}

func TestCheckWithExpiredCache(t *testing.T) {
	tmpDir := t.TempDir()

	config := UpdateConfig{
		Owner:          "jjuanrivvera",
		Repo:           "canvas-cli",
		CurrentVersion: "dev",
		CacheTTL:       1 * time.Millisecond,
	}

	checker := NewChecker(config, tmpDir)

	// Create an old cached result
	oldResult := &CheckResult{
		UpdateAvailable: true,
		CurrentVersion:  "v1.0.0",
		LatestVersion:   "v1.1.0",
		CheckedAt:       time.Now().Add(-2 * time.Millisecond),
	}

	if err := checker.saveCache(oldResult); err != nil {
		t.Fatalf("saveCache failed: %v", err)
	}

	// Wait for cache to expire
	time.Sleep(2 * time.Millisecond)

	ctx := context.Background()
	result, err := checker.Check(ctx)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	// Should get a fresh result (not from cache)
	// For dev version, we expect no updates available
	if result.UpdateAvailable {
		t.Error("Should not use expired cache")
	}
}

func TestCheckWithValidCache(t *testing.T) {
	tmpDir := t.TempDir()

	config := UpdateConfig{
		Owner:          "jjuanrivvera",
		Repo:           "canvas-cli",
		CurrentVersion: "v1.0.0",
		CacheTTL:       1 * time.Hour,
	}

	checker := NewChecker(config, tmpDir)

	// Create a recent cached result
	cachedResult := &CheckResult{
		UpdateAvailable: true,
		CurrentVersion:  "v1.0.0",
		LatestVersion:   "v1.1.0",
		CheckedAt:       time.Now(),
	}

	if err := checker.saveCache(cachedResult); err != nil {
		t.Fatalf("saveCache failed: %v", err)
	}

	ctx := context.Background()
	result, err := checker.Check(ctx)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	// Should use cached result
	if result.LatestVersion != "v1.1.0" {
		t.Errorf("Expected to use cached result with version v1.1.0, got %s", result.LatestVersion)
	}
}

func TestCheckForceCheckIgnoresCache(t *testing.T) {
	tmpDir := t.TempDir()

	config := UpdateConfig{
		Owner:          "jjuanrivvera",
		Repo:           "canvas-cli",
		CurrentVersion: "dev",
		ForceCheck:     true,
		CacheTTL:       1 * time.Hour,
	}

	checker := NewChecker(config, tmpDir)

	// Create a recent cached result
	cachedResult := &CheckResult{
		UpdateAvailable: true,
		CurrentVersion:  "v1.0.0",
		LatestVersion:   "v1.1.0",
		CheckedAt:       time.Now(),
	}

	if err := checker.saveCache(cachedResult); err != nil {
		t.Fatalf("saveCache failed: %v", err)
	}

	ctx := context.Background()
	result, err := checker.Check(ctx)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	// Should NOT use cache when ForceCheck is true
	// For dev version, we expect the current version to be "dev"
	if result.CurrentVersion != "dev" {
		t.Errorf("Expected fresh check with version 'dev', got %s", result.CurrentVersion)
	}
}

func TestLoadCacheInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()

	config := UpdateConfig{
		Owner:          "jjuanrivvera",
		Repo:           "canvas-cli",
		CurrentVersion: "v1.0.0",
	}

	checker := NewChecker(config, tmpDir)

	// Write invalid JSON to cache file
	cachePath := filepath.Join(tmpDir, cacheFileName)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("Failed to create cache dir: %v", err)
	}
	if err := os.WriteFile(cachePath, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("Failed to write invalid cache: %v", err)
	}

	// Try to load invalid cache
	result, err := checker.loadCache()
	if err == nil {
		t.Error("Expected error when loading invalid JSON")
	}
	if result != nil {
		t.Error("Should not return result for invalid cache")
	}
}

func TestClearCacheNonExistent(t *testing.T) {
	tmpDir := t.TempDir()

	config := UpdateConfig{
		Owner:          "jjuanrivvera",
		Repo:           "canvas-cli",
		CurrentVersion: "v1.0.0",
	}

	checker := NewChecker(config, tmpDir)

	// Try to clear non-existent cache (should not error)
	if err := checker.ClearCache(); err != nil {
		t.Errorf("ClearCache should not error for non-existent cache: %v", err)
	}
}

func TestSaveCacheCreatesDirectory(t *testing.T) {
	// Use a non-existent directory
	tmpBase := t.TempDir()
	cacheDir := filepath.Join(tmpBase, "subdir", "cache")

	config := UpdateConfig{
		Owner:          "jjuanrivvera",
		Repo:           "canvas-cli",
		CurrentVersion: "v1.0.0",
	}

	checker := NewChecker(config, cacheDir)

	result := &CheckResult{
		UpdateAvailable: false,
		CurrentVersion:  "v1.0.0",
		LatestVersion:   "v1.0.0",
		CheckedAt:       time.Now(),
	}

	// Should create the directory and save successfully
	if err := checker.saveCache(result); err != nil {
		t.Fatalf("saveCache should create directory: %v", err)
	}

	// Verify the file was created
	cachePath := filepath.Join(cacheDir, cacheFileName)
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		t.Error("Cache file was not created")
	}
}
