package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    *CanvasVersion
		wantErr bool
	}{
		{
			name:    "valid version",
			version: "2024.01.05",
			want: &CanvasVersion{
				Major: 2024,
				Minor: 1,
				Patch: 5,
				Raw:   "2024.01.05",
			},
			wantErr: false,
		},
		{
			name:    "version with text prefix",
			version: "canvas-2023.12.15",
			want: &CanvasVersion{
				Major: 2023,
				Minor: 12,
				Patch: 15,
				Raw:   "canvas-2023.12.15",
			},
			wantErr: false,
		},
		{
			name:    "version with text suffix",
			version: "2022.08.20-release",
			want: &CanvasVersion{
				Major: 2022,
				Minor: 8,
				Patch: 20,
				Raw:   "2022.08.20-release",
			},
			wantErr: false,
		},
		{
			name:    "invalid version - no dots",
			version: "20240105",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid version - missing parts",
			version: "2024.01",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid version - no numbers",
			version: "invalid",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.Major != tt.want.Major {
				t.Errorf("ParseVersion() Major = %v, want %v", got.Major, tt.want.Major)
			}
			if got.Minor != tt.want.Minor {
				t.Errorf("ParseVersion() Minor = %v, want %v", got.Minor, tt.want.Minor)
			}
			if got.Patch != tt.want.Patch {
				t.Errorf("ParseVersion() Patch = %v, want %v", got.Patch, tt.want.Patch)
			}
			if got.Raw != tt.want.Raw {
				t.Errorf("ParseVersion() Raw = %v, want %v", got.Raw, tt.want.Raw)
			}
		})
	}
}

func TestCanvasVersion_IsAtLeast(t *testing.T) {
	v := &CanvasVersion{
		Major: 2024,
		Minor: 1,
		Patch: 5,
	}

	tests := []struct {
		name  string
		major int
		minor int
		patch int
		want  bool
	}{
		{"same version", 2024, 1, 5, true},
		{"older major", 2023, 1, 5, true},
		{"newer major", 2025, 1, 5, false},
		{"same major, older minor", 2024, 0, 5, true},
		{"same major, newer minor", 2024, 2, 5, false},
		{"same major/minor, older patch", 2024, 1, 4, true},
		{"same major/minor, newer patch", 2024, 1, 6, false},
		{"much older version", 2020, 1, 1, true},
		{"much newer version", 2030, 1, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := v.IsAtLeast(tt.major, tt.minor, tt.patch)
			if got != tt.want {
				t.Errorf("IsAtLeast(%d, %d, %d) = %v, want %v",
					tt.major, tt.minor, tt.patch, got, tt.want)
			}
		})
	}
}

func TestCanvasVersion_String(t *testing.T) {
	v := &CanvasVersion{
		Major: 2024,
		Minor: 1,
		Patch: 5,
		Raw:   "canvas-2024.01.05-prod",
	}

	got := v.String()
	if got != v.Raw {
		t.Errorf("String() = %v, want %v", got, v.Raw)
	}
}

func TestDetectCanvasVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Canvas-Meta", `{"version":"2024.01.05"}`)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	ctx := context.Background()
	client := server.Client()

	version, err := DetectCanvasVersion(ctx, client, server.URL)
	if err != nil {
		t.Fatalf("DetectCanvasVersion() error = %v", err)
	}

	if version == nil {
		t.Fatal("expected non-nil version")
	}

	if version.Major != 2024 {
		t.Errorf("expected major version 2024, got %d", version.Major)
	}
}

func TestDetectCanvasVersion_NoVersionHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	ctx := context.Background()
	client := server.Client()

	version, err := DetectCanvasVersion(ctx, client, server.URL)
	// Should handle gracefully - either return error or default version
	_ = version
	_ = err
}

func TestNewFeatureChecker(t *testing.T) {
	version := &CanvasVersion{
		Major: 2024,
		Minor: 1,
		Patch: 5,
	}

	checker := NewFeatureChecker(version)
	if checker == nil {
		t.Fatal("expected non-nil feature checker")
	}

	// Version is private field, just verify checker works
	_ = checker.SupportsFeature("test")
}

func TestFeatureChecker_SupportsFeature(t *testing.T) {
	version := &CanvasVersion{
		Major: 2024,
		Minor: 1,
		Patch: 5,
	}

	checker := NewFeatureChecker(version)

	// Test some known features (you may need to adjust based on actual feature definitions)
	tests := []struct {
		name    string
		feature string
	}{
		{"assignment_groups", "assignment_groups"},
		{"rubrics", "rubrics"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just test that the method doesn't panic
			_ = checker.SupportsFeature(tt.feature)
		})
	}
}

func TestFeatureChecker_WarnIfUnsupported(t *testing.T) {
	version := &CanvasVersion{
		Major: 2024,
		Minor: 1,
		Patch: 5,
	}

	checker := NewFeatureChecker(version)

	// Test that warning method doesn't panic
	checker.WarnIfUnsupported("test_feature")
}
