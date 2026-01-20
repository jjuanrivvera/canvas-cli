package updates

import (
	"time"
)

// ReleaseInfo contains information about a GitHub release
type ReleaseInfo struct {
	Version     string    `json:"version"`
	URL         string    `json:"url"`
	ReleaseDate time.Time `json:"release_date"`
	Notes       string    `json:"notes"`
	AssetURL    string    `json:"asset_url"`
	AssetName   string    `json:"asset_name"`
}

// CheckResult represents the result of checking for updates
type CheckResult struct {
	UpdateAvailable bool         `json:"update_available"`
	CurrentVersion  string       `json:"current_version"`
	LatestVersion   string       `json:"latest_version"`
	ReleaseInfo     *ReleaseInfo `json:"release_info,omitempty"`
	CheckedAt       time.Time    `json:"checked_at"`
}

// UpdateConfig holds configuration for the update checker
type UpdateConfig struct {
	// Repository owner (e.g., "jjuanrivvera")
	Owner string
	// Repository name (e.g., "canvas-cli")
	Repo string
	// Current version of the application
	CurrentVersion string
	// Force check even if cache is valid
	ForceCheck bool
	// Cache TTL duration
	CacheTTL time.Duration
}
