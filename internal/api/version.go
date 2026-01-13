package api

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

// CanvasVersion represents the Canvas version information
type CanvasVersion struct {
	Major int
	Minor int
	Patch int
	Raw   string
}

var versionRegex = regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)

// ParseVersion parses a version string into a CanvasVersion
func ParseVersion(version string) (*CanvasVersion, error) {
	matches := versionRegex.FindStringSubmatch(version)
	if len(matches) != 4 {
		return nil, fmt.Errorf("invalid version format: %s", version)
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	return &CanvasVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
		Raw:   version,
	}, nil
}

// IsAtLeast checks if the version is at least the specified version
func (v *CanvasVersion) IsAtLeast(major, minor, patch int) bool {
	if v.Major > major {
		return true
	}
	if v.Major < major {
		return false
	}

	// Major version is equal, check minor
	if v.Minor > minor {
		return true
	}
	if v.Minor < minor {
		return false
	}

	// Minor version is equal, check patch
	return v.Patch >= patch
}

// String returns the string representation of the version
func (v *CanvasVersion) String() string {
	return v.Raw
}

// versionCacheItem represents a cached version detection result
type versionCacheItem struct {
	Version    *CanvasVersion `json:"version"`
	Expiration time.Time      `json:"expiration"`
	Unknown    bool           `json:"unknown"` // true if version couldn't be detected
}

// getVersionCachePath returns the path to the version cache file for a specific base URL
func getVersionCachePath(baseURL string) string {
	// Get cache directory
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	cacheDir = filepath.Join(cacheDir, "canvas-cli")

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		// If we can't create the cache directory, just return the path anyway
		// The write will fail later but that's acceptable for caching
		slog.Debug("Failed to create cache directory", "error", err)
	}

	// Hash the baseURL to create a unique cache file name
	hash := md5.Sum([]byte(baseURL))
	filename := "version_" + hex.EncodeToString(hash[:]) + ".json"

	return filepath.Join(cacheDir, filename)
}

// loadCachedVersion loads the version from cache if available and not expired
func loadCachedVersion(baseURL string) (*CanvasVersion, bool, bool) {
	cachePath := getVersionCachePath(baseURL)

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, false, false
	}

	var item versionCacheItem
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, false, false
	}

	// Check if expired (24 hours)
	if time.Now().After(item.Expiration) {
		os.Remove(cachePath)
		return nil, false, false
	}

	return item.Version, true, item.Unknown
}

// saveCachedVersion saves the version detection result to cache
func saveCachedVersion(baseURL string, version *CanvasVersion, unknown bool) {
	cachePath := getVersionCachePath(baseURL)

	item := versionCacheItem{
		Version:    version,
		Expiration: time.Now().Add(24 * time.Hour),
		Unknown:    unknown,
	}

	data, err := json.Marshal(item)
	if err != nil {
		return
	}

	if err := os.WriteFile(cachePath, data, 0600); err != nil {
		slog.Debug("Failed to write version cache", "error", err)
	}
}

// DetectCanvasVersion detects the Canvas version from the API
// It caches the result for 24 hours to avoid repeated warnings
func DetectCanvasVersion(ctx context.Context, client *http.Client, baseURL string) (*CanvasVersion, error) {
	// Check cache first
	if cachedVersion, found, wasUnknown := loadCachedVersion(baseURL); found {
		// Don't log warnings for cached unknown versions to reduce noise
		if !wasUnknown {
			slog.Debug("Using cached Canvas version", "version", cachedVersion.String())
		}
		return cachedVersion, nil
	}

	// Try to get version from /api/v1/accounts endpoint
	// Canvas includes version info in the X-Canvas-Meta header
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/accounts", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %w", err)
	}
	defer resp.Body.Close()

	// Check for Canvas meta header
	canvasMeta := resp.Header.Get("X-Canvas-Meta")
	if canvasMeta != "" {
		// Parse the meta header (usually contains version info)
		var meta map[string]interface{}
		if err := json.Unmarshal([]byte(canvasMeta), &meta); err == nil {
			if version, ok := meta["version"].(string); ok {
				v, err := ParseVersion(version)
				if err == nil {
					slog.Info("Detected Canvas version", "version", v.String())
					saveCachedVersion(baseURL, v, false)
					return v, nil
				}
			}
		}
	}

	// If we can't detect the version, assume latest and cache it
	// Only warn once when first detecting, cached lookups won't show the warning
	slog.Warn("Could not detect Canvas version, assuming latest (this warning will not repeat for 24 hours)")
	unknownVersion := &CanvasVersion{
		Major: 999,
		Minor: 999,
		Patch: 999,
		Raw:   "unknown",
	}
	saveCachedVersion(baseURL, unknownVersion, true)
	return unknownVersion, nil
}

// FeatureChecker checks if a feature is available based on Canvas version
type FeatureChecker struct {
	version *CanvasVersion
	logger  *slog.Logger
}

// NewFeatureChecker creates a new FeatureChecker
func NewFeatureChecker(version *CanvasVersion) *FeatureChecker {
	return &FeatureChecker{
		version: version,
		logger:  slog.Default(),
	}
}

// SupportsFeature checks if a feature is supported based on version
func (f *FeatureChecker) SupportsFeature(feature string) bool {
	switch feature {
	case "graphql":
		// GraphQL API was introduced in Canvas 2019.x
		return f.version.IsAtLeast(2019, 0, 0)
	case "new_quizzes":
		// New Quizzes was introduced in Canvas 2020.x
		return f.version.IsAtLeast(2020, 0, 0)
	case "outcomes":
		// Outcomes API improvements in Canvas 2021.x
		return f.version.IsAtLeast(2021, 0, 0)
	case "rubrics_v2":
		// Enhanced Rubrics in Canvas 2022.x
		return f.version.IsAtLeast(2022, 0, 0)
	case "canvas_studio":
		// Canvas Studio integration in Canvas 2023.x
		return f.version.IsAtLeast(2023, 0, 0)
	default:
		// Unknown feature, assume supported
		return true
	}
}

// WarnIfUnsupported logs a warning if a feature is unsupported
func (f *FeatureChecker) WarnIfUnsupported(feature string) bool {
	supported := f.SupportsFeature(feature)
	if !supported {
		f.logger.Warn("Feature not supported in this Canvas version",
			"feature", feature,
			"version", f.version.String(),
		)
	}
	return supported
}
