package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
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

// DetectCanvasVersion detects the Canvas version from the API
func DetectCanvasVersion(ctx context.Context, client *http.Client, baseURL string) (*CanvasVersion, error) {
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
					return v, nil
				}
			}
		}
	}

	// If we can't detect the version, assume latest
	slog.Warn("Could not detect Canvas version, assuming latest")
	return &CanvasVersion{
		Major: 999,
		Minor: 999,
		Patch: 999,
		Raw:   "unknown",
	}, nil
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
