package api

import (
	"testing"
)

func TestFeatureChecker_SupportsFeature_AllCases(t *testing.T) {
	tests := []struct {
		name     string
		version  *CanvasVersion
		feature  string
		expected bool
	}{
		// graphql: requires >= 2019.0.0
		{
			name:     "graphql supported (2024)",
			version:  &CanvasVersion{Major: 2024, Minor: 9, Patch: 0},
			feature:  "graphql",
			expected: true,
		},
		{
			name:     "graphql not supported (2018)",
			version:  &CanvasVersion{Major: 2018, Minor: 12, Patch: 0},
			feature:  "graphql",
			expected: false,
		},
		// new_quizzes: requires >= 2020.0.0
		{
			name:     "new_quizzes supported (2021)",
			version:  &CanvasVersion{Major: 2021, Minor: 1, Patch: 0},
			feature:  "new_quizzes",
			expected: true,
		},
		{
			name:     "new_quizzes not supported (2019)",
			version:  &CanvasVersion{Major: 2019, Minor: 11, Patch: 0},
			feature:  "new_quizzes",
			expected: false,
		},
		// outcomes: requires >= 2021.0.0
		{
			name:     "outcomes supported (2022)",
			version:  &CanvasVersion{Major: 2022, Minor: 0, Patch: 0},
			feature:  "outcomes",
			expected: true,
		},
		{
			name:     "outcomes not supported (2020)",
			version:  &CanvasVersion{Major: 2020, Minor: 6, Patch: 0},
			feature:  "outcomes",
			expected: false,
		},
		// rubrics_v2: requires >= 2022.0.0
		{
			name:     "rubrics_v2 supported (2023)",
			version:  &CanvasVersion{Major: 2023, Minor: 0, Patch: 0},
			feature:  "rubrics_v2",
			expected: true,
		},
		{
			name:     "rubrics_v2 not supported (2021)",
			version:  &CanvasVersion{Major: 2021, Minor: 9, Patch: 5},
			feature:  "rubrics_v2",
			expected: false,
		},
		// canvas_studio: requires >= 2023.0.0
		{
			name:     "canvas_studio supported (2024)",
			version:  &CanvasVersion{Major: 2024, Minor: 0, Patch: 0},
			feature:  "canvas_studio",
			expected: true,
		},
		{
			name:     "canvas_studio not supported (2022)",
			version:  &CanvasVersion{Major: 2022, Minor: 12, Patch: 0},
			feature:  "canvas_studio",
			expected: false,
		},
		// unknown feature defaults to true
		{
			name:     "unknown feature returns true",
			version:  &CanvasVersion{Major: 2018, Minor: 0, Patch: 0},
			feature:  "totally.unknown.feature",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := NewFeatureChecker(tt.version)
			got := checker.SupportsFeature(tt.feature)
			if got != tt.expected {
				t.Errorf("SupportsFeature(%q) = %v, want %v", tt.feature, got, tt.expected)
			}
		})
	}
}

func TestFeatureChecker_WarnIfUnsupported_Supported(t *testing.T) {
	version := &CanvasVersion{Major: 2024, Minor: 0, Patch: 0}
	checker := NewFeatureChecker(version)
	result := checker.WarnIfUnsupported("graphql")
	if !result {
		t.Error("expected true for supported feature")
	}
}

func TestFeatureChecker_WarnIfUnsupported_Unsupported(t *testing.T) {
	// Use an old Canvas version that doesn't support graphql
	version := &CanvasVersion{Major: 2018, Minor: 0, Patch: 0, Raw: "2018.01.01"}
	checker := NewFeatureChecker(version)
	result := checker.WarnIfUnsupported("graphql")
	if result {
		t.Error("expected false for unsupported feature")
	}
}

func TestFeatureChecker_WarnIfUnsupported_AllFeatures_OldVersion(t *testing.T) {
	// Test WarnIfUnsupported covers the warning log path for all feature types
	old := &CanvasVersion{Major: 2018, Minor: 1, Patch: 0, Raw: "2018.01.15"}
	checker := NewFeatureChecker(old)

	features := []string{"graphql", "new_quizzes", "outcomes", "rubrics_v2", "canvas_studio"}
	for _, f := range features {
		result := checker.WarnIfUnsupported(f)
		if result {
			t.Errorf("expected false (unsupported) for feature %q on old version", f)
		}
	}
}

func TestParseVersion_Valid(t *testing.T) {
	v, err := ParseVersion("2024.09.13")
	if err != nil {
		t.Fatalf("ParseVersion: %v", err)
	}
	if v.Major != 2024 || v.Minor != 9 || v.Patch != 13 {
		t.Errorf("unexpected version: %+v", v)
	}
}

func TestParseVersion_Invalid(t *testing.T) {
	_, err := ParseVersion("not-a-version")
	if err == nil {
		t.Error("expected error for invalid version string")
	}
}

func TestSaveCachedVersion_AndLoad(t *testing.T) {
	// Use a unique test URL to avoid polluting real cache
	testURL := "https://test-canvas-version-cache.example.invalid"

	version := &CanvasVersion{Major: 2024, Minor: 9, Patch: 1, Raw: "2024.09.01"}
	saveCachedVersion(testURL, version, false)

	loaded, found, wasUnknown := loadCachedVersion(testURL)
	if !found {
		t.Fatal("expected cached version to be found")
	}
	if wasUnknown {
		t.Error("expected wasUnknown=false")
	}
	if loaded.Major != 2024 {
		t.Errorf("expected major=2024, got %d", loaded.Major)
	}
}

func TestSaveCachedVersion_Unknown(t *testing.T) {
	testURL := "https://test-canvas-version-unknown.example.invalid"

	unknown := &CanvasVersion{Major: 999, Minor: 999, Patch: 999, Raw: "unknown"}
	saveCachedVersion(testURL, unknown, true)

	loaded, found, wasUnknown := loadCachedVersion(testURL)
	if !found {
		t.Fatal("expected cached version to be found")
	}
	if !wasUnknown {
		t.Error("expected wasUnknown=true")
	}
	if loaded.Major != 999 {
		t.Errorf("expected major=999, got %d", loaded.Major)
	}
}
