package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jjuanrivvera/canvas-cli/internal/cache"
)

// cacheTestItem is a minimal decode target for the cache-safety tests below,
// kept independent of the full Course model so the fixtures stay small and
// obvious.
type cacheTestItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// TestClient_DryRun_DoesNotWriteToCache reproduces bug 1: a dry run used to
// write its placeholder body under the real cache key, because every cache
// gate in client.go only checked "cacheEnabled && cache != nil" with no
// dry-run check. That let a later real (non-dry-run) request within the
// cache TTL be served the dry run's empty rehearsal data instead of hitting
// Canvas.
func TestClient_DryRun_DoesNotWriteToCache(t *testing.T) {
	testCache := cache.New(5 * time.Minute)

	client, err := NewClient(ClientConfig{
		BaseURL:        "https://canvas.example.com",
		Token:          "test-token",
		RequestsPerSec: 10,
		Cache:          testCache,
		CacheEnabled:   true,
		DryRun:         true,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	var items []cacheTestItem
	if err := client.GetJSON(context.Background(), "/api/v1/courses", &items); err != nil {
		t.Fatalf("GetJSON in dry-run mode: %v", err)
	}

	key := client.cacheKey("/api/v1/courses")
	if testCache.Has(key) {
		t.Fatal("dry run must not write a cache entry, but the placeholder response was cached")
	}
}

// TestClient_DryRun_DoesNotReadFromCache reproduces the read side of bug 1:
// a dry run must never be served real data back from an earlier cache entry
// either, because a dry run's whole point is to show what request WOULD be
// sent without actually depending on live (or previously live) data.
func TestClient_DryRun_DoesNotReadFromCache(t *testing.T) {
	testCache := cache.New(5 * time.Minute)

	client, err := NewClient(ClientConfig{
		BaseURL:        "https://canvas.example.com",
		Token:          "test-token",
		RequestsPerSec: 10,
		Cache:          testCache,
		CacheEnabled:   true,
		DryRun:         true,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	// Pre-populate the cache as if an earlier, real (non-dry-run) request
	// had already cached a real Canvas response under this exact key.
	key := client.cacheKey("/api/v1/courses/123")
	real := cacheTestItem{ID: 123, Name: "Real Course"}
	realJSON, err := json.Marshal(real)
	if err != nil {
		t.Fatalf("marshal fixture: %v", err)
	}
	testCache.Set(key, realJSON)

	var got cacheTestItem
	if err := client.GetJSON(context.Background(), "/api/v1/courses/123", &got); err != nil {
		t.Fatalf("GetJSON in dry-run mode: %v", err)
	}

	if got.ID == real.ID {
		t.Fatalf("dry run served real cached data back (got %+v); it must hit the dry-run placeholder instead", got)
	}
}

// TestClient_GetAllPages_UndecodableCachedList_LiveFetchIsClean reproduces
// bug 2 on the paginated (list) path. encoding/json keeps decoding a slice
// element by element even after one element fails, so a cache entry with one
// bad element still leaves the earlier and later elements populated. The old
// code decoded straight into the caller's result slice, so a failed cache
// read left it holding that wreckage, and the live fetch that followed
// appended its rows onto it instead of starting clean.
func TestClient_GetAllPages_UndecodableCachedList_LiveFetchIsClean(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"id":1,"name":"Live A"},{"id":2,"name":"Live B"}]`))
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
		t.Fatalf("NewClient: %v", err)
	}

	// The middle element's "id" is a string, so it fails to decode into the
	// int64 field, but encoding/json still decodes the elements before and
	// after it.
	key := client.cacheKey("pages:/api/v1/courses")
	testCache.Set(key, []byte(`[{"id":901,"name":"Stale A"},{"id":"bad","name":"Stale B"},{"id":903,"name":"Stale C"}]`))

	var got []cacheTestItem
	if err := client.GetAllPages(context.Background(), "/api/v1/courses", &got); err != nil {
		t.Fatalf("GetAllPages: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected exactly the 2 live rows, got %d: %+v", len(got), got)
	}
	for _, item := range got {
		if item.ID == 901 || item.ID == 903 {
			t.Fatalf("live fetch result still contains stale rows from the poisoned cache read: %+v", got)
		}
	}
}

// TestClient_GetJSON_UndecodableCachedObject_LiveFetchIsClean reproduces bug
// 2 on the single-object path. A cached object whose "id" fails to decode
// still leaves "name" populated from the cache, and the old code decoded
// straight into the caller's result. When the live response for the retry
// doesn't happen to mention "name" at all, the stale cached value survives
// the live decode untouched — a fabricated field the server never sent.
func TestClient_GetJSON_UndecodableCachedObject_LiveFetchIsClean(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		// The live response never mentions "name".
		w.Write([]byte(`{"id":55}`))
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
		t.Fatalf("NewClient: %v", err)
	}

	key := client.cacheKey("/api/v1/courses/55")
	testCache.Set(key, []byte(`{"id":"bad","name":"StaleName"}`))

	var got cacheTestItem
	if err := client.GetJSON(context.Background(), "/api/v1/courses/55", &got); err != nil {
		t.Fatalf("GetJSON: %v", err)
	}

	if got.ID != 55 {
		t.Fatalf("expected the live id 55, got %d", got.ID)
	}
	if got.Name != "" {
		t.Fatalf("live fetch result contains a fabricated field from the poisoned cache read: %+v", got)
	}
}

// TestClient_NonDryRun_CachingStillWorks guards against the fix being
// "just disable the cache": a normal (non-dry-run) client must still serve
// its second request for the same path from the cache without hitting the
// server again.
func TestClient_NonDryRun_CachingStillWorks(t *testing.T) {
	var calls int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		atomic.AddInt32(&calls, 1)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":7,"name":"Cached Course"}`))
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
		t.Fatalf("NewClient: %v", err)
	}

	ctx := context.Background()

	var first cacheTestItem
	if err := client.GetJSON(ctx, "/api/v1/courses/7", &first); err != nil {
		t.Fatalf("first GetJSON: %v", err)
	}

	var second cacheTestItem
	if err := client.GetJSON(ctx, "/api/v1/courses/7", &second); err != nil {
		t.Fatalf("second GetJSON: %v", err)
	}

	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("expected exactly 1 live request (the second call should be a cache hit), got %d", got)
	}
	if second.ID != 7 || second.Name != "Cached Course" {
		t.Fatalf("cache hit returned unexpected data: %+v", second)
	}

	key := client.cacheKey("/api/v1/courses/7")
	if !testCache.Has(key) {
		t.Fatal("expected the live response to have been written to the cache")
	}
}
