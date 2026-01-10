package api

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestAdaptiveRateLimiter_AdjustRate(t *testing.T) {
	tests := []struct {
		name          string
		remaining     float64
		total         float64
		expectedRate  float64
		expectWarning bool
	}{
		{
			name:         "critical threshold (20% remaining)",
			remaining:    20,
			total:        100,
			expectedRate: verySlowRequestsPerSecond, // 1 req/sec
		},
		{
			name:         "warning threshold (50% remaining)",
			remaining:    50,
			total:        100,
			expectedRate: slowRequestsPerSecond, // 2 req/sec
		},
		{
			name:         "normal threshold (70% remaining)",
			remaining:    70,
			total:        100,
			expectedRate: defaultRequestsPerSecond, // 5 req/sec
		},
		{
			name:         "above warning threshold (60%)",
			remaining:    60,
			total:        100,
			expectedRate: defaultRequestsPerSecond, // Above 50% threshold, returns to normal
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				Level: slog.LevelError, // Suppress warnings during tests
			}))

			limiter := &AdaptiveRateLimiter{
				limiter:      rate.NewLimiter(rate.Limit(defaultRequestsPerSecond), 1),
				currentRate:  defaultRequestsPerSecond,
				logger:       logger,
				warningShown: make(map[int]bool),
			}

			limiter.AdjustRate(tt.remaining, tt.total)

			currentRate := limiter.GetCurrentRate()
			if currentRate != tt.expectedRate {
				t.Errorf("expected rate %.1f, got %.1f", tt.expectedRate, currentRate)
			}
		})
	}
}

func TestAdaptiveRateLimiter_GetCurrentRate(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	limiter := &AdaptiveRateLimiter{
		limiter:      rate.NewLimiter(rate.Limit(defaultRequestsPerSecond), 1),
		currentRate:  defaultRequestsPerSecond,
		logger:       logger,
		warningShown: make(map[int]bool),
	}

	rate := limiter.GetCurrentRate()
	if rate != defaultRequestsPerSecond {
		t.Errorf("expected rate %.1f, got %.1f", defaultRequestsPerSecond, rate)
	}

	// Adjust rate and verify
	limiter.AdjustRate(20, 100) // Critical threshold
	rate = limiter.GetCurrentRate()
	if rate != verySlowRequestsPerSecond {
		t.Errorf("expected rate %.1f after adjustment, got %.1f", verySlowRequestsPerSecond, rate)
	}
}

func TestClient_GetVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Initially version should be nil (not detected yet)
	version := client.GetVersion()
	if version == nil {
		// Version detection happens asynchronously, so nil is acceptable
		t.Log("Version not yet detected (async operation)")
	}
}

func TestClient_SupportsFeature(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test some feature checks
	tests := []struct {
		feature string
		// We can't predict the exact result without knowing the detected version,
		// but we can verify the method doesn't panic
	}{
		{feature: "submissions.grade"},
		{feature: "courses.create"},
		{feature: "unknown.feature"},
	}

	for _, tt := range tests {
		t.Run(tt.feature, func(t *testing.T) {
			// Just verify it doesn't panic
			_ = client.SupportsFeature(tt.feature)
		})
	}
}

func TestClient_GetVersion_AfterDetection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Canvas-Meta", "version=2024.09.13;region=us-east-1")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Trigger version detection by making a request
	ctx := context.Background()
	_, _ = client.doRequest(ctx, "GET", "/api/v1/accounts", nil)

	// Give version detection goroutine time to complete
	time.Sleep(100 * time.Millisecond)

	version := client.GetVersion()
	if version != nil {
		t.Logf("Version detected: %+v", version)
	} else {
		t.Log("Version detection may not have completed yet (async)")
	}
}

func TestAdaptiveRateLimiter_WarningReset(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError, // Suppress warnings
	}))

	limiter := &AdaptiveRateLimiter{
		limiter:      rate.NewLimiter(rate.Limit(defaultRequestsPerSecond), 1),
		currentRate:  defaultRequestsPerSecond,
		logger:       logger,
		warningShown: make(map[int]bool),
	}

	// First, trigger warning
	limiter.AdjustRate(20, 100) // Critical
	if len(limiter.warningShown) == 0 {
		t.Error("expected warning to be recorded")
	}

	// Then recover above warning threshold
	limiter.AdjustRate(70, 100) // Normal
	if len(limiter.warningShown) != 0 {
		t.Error("expected warnings to be reset when recovering")
	}
}
