package api

import (
	"testing"
)

func TestIsCanvasDomain_SameHost(t *testing.T) {
	tests := []struct {
		redirectURL string
		baseURL     string
		want        bool
	}{
		{
			redirectURL: "https://canvas.example.com/files/1/download",
			baseURL:     "https://canvas.example.com",
			want:        true,
		},
		{
			redirectURL: "https://s3.amazonaws.com/bucket/file.pdf",
			baseURL:     "https://canvas.example.com",
			want:        false,
		},
		{
			redirectURL: "https://canvas.example.com/api/v1/files",
			baseURL:     "https://canvas.example.com/something",
			want:        true,
		},
		{
			redirectURL: "https://other.example.com/files/1",
			baseURL:     "https://canvas.example.com",
			want:        false,
		},
	}

	for _, tt := range tests {
		got := isCanvasDomain(tt.redirectURL, tt.baseURL)
		if got != tt.want {
			t.Errorf("isCanvasDomain(%q, %q) = %v, want %v",
				tt.redirectURL, tt.baseURL, got, tt.want)
		}
	}
}

func TestIsCanvasDomain_InvalidURL(t *testing.T) {
	// An unparseable redirect URL returns false.
	got := isCanvasDomain("://not-valid", "https://canvas.example.com")
	if got != false {
		t.Error("expected false for invalid redirect URL")
	}

	// An unparseable base URL returns false.
	got = isCanvasDomain("https://canvas.example.com/files/1", "://not-valid")
	if got != false {
		t.Error("expected false for invalid base URL")
	}
}
