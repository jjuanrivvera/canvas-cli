package cache

import (
	"testing"
	"time"
)

func TestNewMultiTierCache(t *testing.T) {
	tempDir := t.TempDir()

	multiCache, err := NewMultiTierCache(5*time.Minute, tempDir, 10*time.Minute)
	if err != nil {
		t.Fatalf("NewMultiTierCache failed: %v", err)
	}

	if multiCache == nil {
		t.Fatal("expected non-nil multi-tier cache")
	}
}

func TestMultiTierCache_GetSet(t *testing.T) {
	tempDir := t.TempDir()

	multiCache, err := NewMultiTierCache(5*time.Minute, tempDir, 10*time.Minute)
	if err != nil {
		t.Fatalf("NewMultiTierCache failed: %v", err)
	}

	// Set and get value
	key := "test_key"
	value := []byte(`"test_value"`)

	multiCache.Set(key, value)

	retrieved := multiCache.Get(key)
	if retrieved == nil {
		t.Fatal("expected non-nil value")
	}

	if string(retrieved) != string(value) {
		t.Errorf("expected '%s', got '%s'", value, retrieved)
	}
}

func TestMultiTierCache_Get_Miss(t *testing.T) {
	tempDir := t.TempDir()

	multiCache, err := NewMultiTierCache(5*time.Minute, tempDir, 10*time.Minute)
	if err != nil {
		t.Fatalf("NewMultiTierCache failed: %v", err)
	}

	// Get non-existent key
	value := multiCache.Get("nonexistent")
	if value != nil {
		t.Error("expected nil for cache miss")
	}
}

// Removed - combined with TestMultiTierCache_GetSet

func TestMultiTierCache_SetJSON(t *testing.T) {
	tempDir := t.TempDir()

	multiCache, err := NewMultiTierCache(5*time.Minute, tempDir, 10*time.Minute)
	if err != nil {
		t.Fatalf("NewMultiTierCache failed: %v", err)
	}

	type TestData struct {
		Name string
		Age  int
	}

	key := "json_key"
	data := TestData{Name: "Alice", Age: 30}

	err = multiCache.SetJSON(key, data)
	if err != nil {
		t.Fatalf("SetJSON failed: %v", err)
	}

	var retrieved TestData
	err = multiCache.GetJSON(key, &retrieved)
	if err != nil {
		t.Fatalf("GetJSON failed: %v", err)
	}

	if retrieved.Name != data.Name {
		t.Errorf("expected name '%s', got '%s'", data.Name, retrieved.Name)
	}
}

// Removed - tested via SetJSON

func TestMultiTierCache_SetWithTTL(t *testing.T) {
	tempDir := t.TempDir()

	multiCache, err := NewMultiTierCache(5*time.Minute, tempDir, 10*time.Minute)
	if err != nil {
		t.Fatalf("NewMultiTierCache failed: %v", err)
	}

	key := "ttl_key"
	value := []byte(`"ttl_value"`)

	multiCache.SetWithTTL(key, value, 100*time.Millisecond)

	// Value should exist immediately
	if !multiCache.Has(key) {
		t.Error("expected key to exist")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Value should be expired
	retrieved := multiCache.Get(key)
	if retrieved != nil {
		t.Error("expected nil for expired key")
	}
}

func TestMultiTierCache_Delete(t *testing.T) {
	tempDir := t.TempDir()

	multiCache, err := NewMultiTierCache(5*time.Minute, tempDir, 10*time.Minute)
	if err != nil {
		t.Fatalf("NewMultiTierCache failed: %v", err)
	}

	key := "delete_key"
	value := []byte(`"delete_value"`)

	multiCache.Set(key, value)

	// Verify it exists in both
	if !multiCache.Has(key) {
		t.Error("expected key to exist before delete")
	}

	multiCache.Delete(key)

	// Verify it's gone
	if multiCache.Has(key) {
		t.Error("expected key to not exist after delete")
	}
}

func TestMultiTierCache_Clear(t *testing.T) {
	tempDir := t.TempDir()

	multiCache, err := NewMultiTierCache(5*time.Minute, tempDir, 10*time.Minute)
	if err != nil {
		t.Fatalf("NewMultiTierCache failed: %v", err)
	}

	// Set multiple keys
	multiCache.Set("key1", []byte(`"value1"`))
	multiCache.Set("key2", []byte(`"value2"`))
	multiCache.Set("key3", []byte(`"value3"`))

	multiCache.Clear()

	// Verify all keys are gone
	if multiCache.Has("key1") || multiCache.Has("key2") || multiCache.Has("key3") {
		t.Error("expected all keys to be cleared")
	}
}

func TestMultiTierCache_Has(t *testing.T) {
	tempDir := t.TempDir()

	multiCache, err := NewMultiTierCache(5*time.Minute, tempDir, 10*time.Minute)
	if err != nil {
		t.Fatalf("NewMultiTierCache failed: %v", err)
	}

	key := "has_key"

	// Should not exist initially
	if multiCache.Has(key) {
		t.Error("expected key to not exist initially")
	}

	// Set value
	multiCache.Set(key, []byte(`"value"`))

	// Should exist now
	if !multiCache.Has(key) {
		t.Error("expected key to exist after set")
	}
}

func TestMultiTierCache_Stats(t *testing.T) {
	tempDir := t.TempDir()

	multiCache, err := NewMultiTierCache(5*time.Minute, tempDir, 10*time.Minute)
	if err != nil {
		t.Fatalf("NewMultiTierCache failed: %v", err)
	}

	// Set some values
	multiCache.Set("key1", []byte(`"value1"`))
	multiCache.Set("key2", []byte(`"value2"`))

	stats := multiCache.Stats()

	if stats.Total < 2 {
		t.Errorf("expected at least 2 total items, got %d", stats.Total)
	}
}

func TestMultiTierCache_MemoryStats(t *testing.T) {
	tempDir := t.TempDir()

	multiCache, err := NewMultiTierCache(5*time.Minute, tempDir, 10*time.Minute)
	if err != nil {
		t.Fatalf("NewMultiTierCache failed: %v", err)
	}

	// Set some values
	multiCache.Set("key1", []byte(`"value1"`))

	stats := multiCache.MemoryStats()
	if stats.Total < 1 {
		t.Errorf("expected at least 1 item in memory, got %d", stats.Total)
	}
}

func TestMultiTierCache_DiskStats(t *testing.T) {
	tempDir := t.TempDir()

	multiCache, err := NewMultiTierCache(5*time.Minute, tempDir, 10*time.Minute)
	if err != nil {
		t.Fatalf("NewMultiTierCache failed: %v", err)
	}

	// Set some values
	multiCache.Set("key1", []byte(`"value1"`))

	stats, err := multiCache.DiskStats()
	if err != nil {
		t.Fatalf("DiskStats failed: %v", err)
	}

	if stats.Total < 1 {
		t.Errorf("expected at least 1 item on disk, got %d", stats.Total)
	}
}

func TestNewMultiTierCache_DiskError(t *testing.T) {
	// Use an invalid directory path to trigger NewDiskCache error
	invalidPath := "/dev/null/invalid/path"

	_, err := NewMultiTierCache(5*time.Minute, invalidPath, 10*time.Minute)
	if err == nil {
		t.Error("expected error when disk cache creation fails")
	}
}

func TestMultiTierCache_GetJSON_Miss(t *testing.T) {
	tempDir := t.TempDir()

	multiCache, err := NewMultiTierCache(5*time.Minute, tempDir, 10*time.Minute)
	if err != nil {
		t.Fatalf("NewMultiTierCache failed: %v", err)
	}

	type TestData struct {
		Name string
	}

	var data TestData
	err = multiCache.GetJSON("nonexistent", &data)
	if err != ErrCacheMiss {
		t.Errorf("expected ErrCacheMiss, got %v", err)
	}
}

func TestMultiTierCache_GetJSON_UnmarshalError(t *testing.T) {
	tempDir := t.TempDir()

	multiCache, err := NewMultiTierCache(5*time.Minute, tempDir, 10*time.Minute)
	if err != nil {
		t.Fatalf("NewMultiTierCache failed: %v", err)
	}

	// Set invalid JSON
	multiCache.Set("bad_json", []byte("not valid json"))

	type TestData struct {
		Name string
	}

	var data TestData
	err = multiCache.GetJSON("bad_json", &data)
	if err == nil {
		t.Error("expected error when unmarshaling invalid JSON")
	}
}

func TestMultiTierCache_SetJSON_MarshalError(t *testing.T) {
	tempDir := t.TempDir()

	multiCache, err := NewMultiTierCache(5*time.Minute, tempDir, 10*time.Minute)
	if err != nil {
		t.Fatalf("NewMultiTierCache failed: %v", err)
	}

	// Channels cannot be marshaled to JSON
	invalidData := make(chan int)
	err = multiCache.SetJSON("key", invalidData)
	if err == nil {
		t.Error("expected error when marshaling invalid JSON data")
	}
}
