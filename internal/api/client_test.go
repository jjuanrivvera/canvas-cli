package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/time/rate"

	"github.com/jjuanrivvera/canvas-cli/internal/cache"
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

func TestClient_GetJSON_WithCache(t *testing.T) {
	var requestCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
			return
		}
		if r.URL.Path == "/api/v1/courses/123" {
			atomic.AddInt32(&requestCount, 1)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":123,"name":"Test Course"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Create cache
	testCache := cache.New(5 * time.Minute)

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
		Cache:          testCache,
		CacheEnabled:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// First request - should hit the server
	var course1 Course
	err = client.GetJSON(ctx, "/api/v1/courses/123", &course1)
	if err != nil {
		t.Fatalf("First GetJSON failed: %v", err)
	}
	if course1.ID != 123 {
		t.Errorf("Expected course ID 123, got %d", course1.ID)
	}

	firstCount := atomic.LoadInt32(&requestCount)
	if firstCount != 1 {
		t.Errorf("Expected 1 request after first call, got %d", firstCount)
	}

	// Second request - should hit cache
	var course2 Course
	err = client.GetJSON(ctx, "/api/v1/courses/123", &course2)
	if err != nil {
		t.Fatalf("Second GetJSON failed: %v", err)
	}
	if course2.ID != 123 {
		t.Errorf("Expected course ID 123, got %d", course2.ID)
	}

	secondCount := atomic.LoadInt32(&requestCount)
	if secondCount != 1 {
		t.Errorf("Expected 1 request after second call (cache hit), got %d", secondCount)
	}
}

func TestClient_GetJSON_CacheDisabled(t *testing.T) {
	var requestCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
			return
		}
		if r.URL.Path == "/api/v1/courses/456" {
			atomic.AddInt32(&requestCount, 1)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":456,"name":"Test Course"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Create cache but disable it
	testCache := cache.New(5 * time.Minute)

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
		Cache:          testCache,
		CacheEnabled:   false, // Disabled
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// First request
	var course1 Course
	err = client.GetJSON(ctx, "/api/v1/courses/456", &course1)
	if err != nil {
		t.Fatalf("First GetJSON failed: %v", err)
	}

	// Second request - should hit server again (cache disabled)
	var course2 Course
	err = client.GetJSON(ctx, "/api/v1/courses/456", &course2)
	if err != nil {
		t.Fatalf("Second GetJSON failed: %v", err)
	}

	count := atomic.LoadInt32(&requestCount)
	if count != 2 {
		t.Errorf("Expected 2 requests when cache is disabled, got %d", count)
	}
}

func TestClient_CacheHelpers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			w.Write([]byte(`[]`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	testCache := cache.New(5 * time.Minute)

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
		Cache:          testCache,
		CacheEnabled:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test IsCacheEnabled
	if !client.IsCacheEnabled() {
		t.Error("Expected cache to be enabled")
	}

	// Test SetCacheEnabled
	client.SetCacheEnabled(false)
	if client.IsCacheEnabled() {
		t.Error("Expected cache to be disabled")
	}

	client.SetCacheEnabled(true)
	if !client.IsCacheEnabled() {
		t.Error("Expected cache to be re-enabled")
	}

	// Test CacheStats (just verify it doesn't panic)
	stats := client.CacheStats()
	t.Logf("Cache stats: %+v", stats)

	// Test ClearCache (just verify it doesn't panic)
	client.ClearCache()
}

func TestClient_CacheKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			w.Write([]byte(`[]`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	testCache := cache.New(5 * time.Minute)

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
		Cache:          testCache,
		CacheEnabled:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test that cache keys are generated correctly
	key1 := client.cacheKey("/api/v1/courses/123")
	key2 := client.cacheKey("/api/v1/courses/123")
	key3 := client.cacheKey("/api/v1/courses/456")

	// Same path should generate same key
	if key1 != key2 {
		t.Errorf("Same path should generate same key: %s != %s", key1, key2)
	}

	// Different paths should generate different keys
	if key1 == key3 {
		t.Errorf("Different paths should generate different keys: %s == %s", key1, key3)
	}
}

func TestClient_CacheKey_WithMasquerade(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			w.Write([]byte(`[]`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	testCache := cache.New(5 * time.Minute)

	// Client without masquerade
	client1, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
		Cache:          testCache,
		CacheEnabled:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create client1: %v", err)
	}

	// Client with masquerade
	client2, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
		AsUserID:       999,
		Cache:          testCache,
		CacheEnabled:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create client2: %v", err)
	}

	key1 := client1.cacheKey("/api/v1/courses/123")
	key2 := client2.cacheKey("/api/v1/courses/123")

	// Keys should be different when masquerading
	if key1 == key2 {
		t.Errorf("Keys should differ when masquerading: %s == %s", key1, key2)
	}
}

func TestClient_SetQuotaTotal(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
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

	// Default quota
	defaultQuota := client.GetQuotaTotal()
	if defaultQuota != defaultQuotaTotal {
		t.Errorf("Expected default quota %f, got %f", defaultQuotaTotal, defaultQuota)
	}

	// Set custom quota
	client.SetQuotaTotal(1000.0)
	newQuota := client.GetQuotaTotal()
	if newQuota != 1000.0 {
		t.Errorf("Expected quota 1000.0, got %f", newQuota)
	}
}

func TestClient_UserAgent_Default(t *testing.T) {
	var receivedUserAgent string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedUserAgent = r.Header.Get("User-Agent")
		if r.URL.Path == "/api/v1/accounts" {
			w.Write([]byte(`[]`))
			return
		}
		if r.URL.Path == "/api/v1/courses" {
			w.Header().Set("Content-Type", "application/json")
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

	// Make a request to trigger User-Agent header
	_, err = client.Get(context.Background(), "/api/v1/courses")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	// Verify default User-Agent is set
	if receivedUserAgent != "canvas-cli" {
		t.Errorf("Expected default User-Agent 'canvas-cli', got '%s'", receivedUserAgent)
	}
}

func TestClient_UserAgent_Custom(t *testing.T) {
	var receivedUserAgent string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedUserAgent = r.Header.Get("User-Agent")
		if r.URL.Path == "/api/v1/accounts" {
			w.Write([]byte(`[]`))
			return
		}
		if r.URL.Path == "/api/v1/courses" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`[]`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	customUserAgent := "canvas-cli/v1.5.0"
	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
		UserAgent:      customUserAgent,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Make a request to trigger User-Agent header
	_, err = client.Get(context.Background(), "/api/v1/courses")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	// Verify custom User-Agent is set
	if receivedUserAgent != customUserAgent {
		t.Errorf("Expected User-Agent '%s', got '%s'", customUserAgent, receivedUserAgent)
	}
}

// BenchmarkGetAllPages_Reflection benchmarks the old reflection-based GetAllPages method
func BenchmarkGetAllPages_Reflection(b *testing.B) {
	// Create test server with paginated data
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		if page == "" || page == "1" {
			w.Header().Set("Link", `<`+r.URL.String()+`?page=2>; rel="next"`)
			w.Write([]byte(`[{"id":1,"name":"Item 1"},{"id":2,"name":"Item 2"}]`))
		} else if page == "2" {
			w.Write([]byte(`[{"id":3,"name":"Item 3"},{"id":4,"name":"Item 4"}]`))
		}
	}))
	defer server.Close()

	client, _ := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 1000, // High rate to avoid throttling in benchmarks
	})

	type Item struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var items []Item
		if err := client.GetAllPages(context.Background(), "/items", &items); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetAllPages_Generics benchmarks the new generic GetAllPages method
func BenchmarkGetAllPages_Generics(b *testing.B) {
	// Create test server with paginated data
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		if page == "" || page == "1" {
			w.Header().Set("Link", `<`+r.URL.String()+`?page=2>; rel="next"`)
			w.Write([]byte(`[{"id":1,"name":"Item 1"},{"id":2,"name":"Item 2"}]`))
		} else if page == "2" {
			w.Write([]byte(`[{"id":3,"name":"Item 3"},{"id":4,"name":"Item 4"}]`))
		}
	}))
	defer server.Close()

	client, _ := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 1000, // High rate to avoid throttling in benchmarks
	})

	type Item struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := GetAllPagesGeneric[Item](client, context.Background(), "/items"); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetAllPages_LargeDataset_Reflection benchmarks reflection with larger dataset
func BenchmarkGetAllPages_LargeDataset_Reflection(b *testing.B) {
	// Create test server with 10 pages of 100 items each (1000 total)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		pageNum := 1
		if page != "" {
			_, _ = fmt.Sscanf(page, "%d", &pageNum)
		}

		if pageNum < 10 {
			w.Header().Set("Link", fmt.Sprintf(`<%s?page=%d>; rel="next"`, r.URL.String(), pageNum+1))
		}

		// Generate 100 items per page
		items := make([]map[string]interface{}, 100)
		for i := 0; i < 100; i++ {
			items[i] = map[string]interface{}{
				"id":   (pageNum-1)*100 + i + 1,
				"name": fmt.Sprintf("Item %d", (pageNum-1)*100+i+1),
			}
		}
		json.NewEncoder(w).Encode(items)
	}))
	defer server.Close()

	client, _ := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 1000,
	})

	type Item struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var items []Item
		if err := client.GetAllPages(context.Background(), "/items", &items); err != nil {
			b.Fatal(err)
		}
		if len(items) != 1000 {
			b.Fatalf("Expected 1000 items, got %d", len(items))
		}
	}
}

// BenchmarkGetAllPages_LargeDataset_Generics benchmarks generics with larger dataset
func BenchmarkGetAllPages_LargeDataset_Generics(b *testing.B) {
	// Create test server with 10 pages of 100 items each (1000 total)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		pageNum := 1
		if page != "" {
			_, _ = fmt.Sscanf(page, "%d", &pageNum)
		}

		if pageNum < 10 {
			w.Header().Set("Link", fmt.Sprintf(`<%s?page=%d>; rel="next"`, r.URL.String(), pageNum+1))
		}

		// Generate 100 items per page
		items := make([]map[string]interface{}, 100)
		for i := 0; i < 100; i++ {
			items[i] = map[string]interface{}{
				"id":   (pageNum-1)*100 + i + 1,
				"name": fmt.Sprintf("Item %d", (pageNum-1)*100+i+1),
			}
		}
		json.NewEncoder(w).Encode(items)
	}))
	defer server.Close()

	client, _ := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 1000,
	})

	type Item struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		items, err := GetAllPagesGeneric[Item](client, context.Background(), "/items")
		if err != nil {
			b.Fatal(err)
		}
		if len(items) != 1000 {
			b.Fatalf("Expected 1000 items, got %d", len(items))
		}
	}
}
