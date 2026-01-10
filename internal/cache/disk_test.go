package cache

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewDiskCache(t *testing.T) {
	tempDir := t.TempDir()

	cache, err := NewDiskCache(tempDir, 5*time.Minute)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	if cache == nil {
		t.Fatal("expected non-nil cache")
	}

	// Verify directory was created
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Error("expected cache directory to be created")
	}
}

func TestDiskCache_SetGet(t *testing.T) {
	tempDir := t.TempDir()
	cache, err := NewDiskCache(tempDir, 5*time.Minute)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	// Set a value (must be valid JSON for RawMessage)
	key := "test_key"
	value := []byte(`"test_value"`)

	err = cache.Set(key, value)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get the value
	retrieved := cache.Get(key)
	if retrieved == nil {
		t.Fatal("expected non-nil value")
	}

	if string(retrieved) != string(value) {
		t.Errorf("expected '%s', got '%s'", value, retrieved)
	}
}

func TestDiskCache_GetNonExistent(t *testing.T) {
	tempDir := t.TempDir()
	cache, err := NewDiskCache(tempDir, 5*time.Minute)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	value := cache.Get("nonexistent")
	if value != nil {
		t.Error("expected nil for nonexistent key")
	}
}

func TestDiskCache_SetJSON(t *testing.T) {
	tempDir := t.TempDir()
	cache, err := NewDiskCache(tempDir, 5*time.Minute)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	type TestStruct struct {
		Name string
		Age  int
	}

	original := TestStruct{Name: "John", Age: 30}

	err = cache.SetJSON("json_key", original)
	if err != nil {
		t.Fatalf("SetJSON failed: %v", err)
	}

	var retrieved TestStruct
	err = cache.GetJSON("json_key", &retrieved)
	if err != nil {
		t.Fatalf("GetJSON failed: %v", err)
	}

	if retrieved.Name != original.Name {
		t.Errorf("expected name '%s', got '%s'", original.Name, retrieved.Name)
	}
	if retrieved.Age != original.Age {
		t.Errorf("expected age %d, got %d", original.Age, retrieved.Age)
	}
}

func TestDiskCache_GetJSON_NotFound(t *testing.T) {
	tempDir := t.TempDir()
	cache, err := NewDiskCache(tempDir, 5*time.Minute)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	var result map[string]string
	err = cache.GetJSON("nonexistent", &result)
	if err == nil {
		t.Error("expected error for nonexistent key")
	}
}

func TestDiskCache_SetWithTTL(t *testing.T) {
	tempDir := t.TempDir()
	cache, err := NewDiskCache(tempDir, 5*time.Minute)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	key := "ttl_key"
	value := []byte(`"ttl_value"`)

	err = cache.SetWithTTL(key, value, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("SetWithTTL failed: %v", err)
	}

	// Value should exist immediately
	if !cache.Has(key) {
		t.Error("expected key to exist")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Value should be expired but file might still exist
	// (cleanup runs periodically)
	retrieved := cache.Get(key)
	if retrieved != nil {
		t.Error("expected nil for expired key")
	}
}

func TestDiskCache_Delete(t *testing.T) {
	tempDir := t.TempDir()
	cache, err := NewDiskCache(tempDir, 5*time.Minute)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	key := "delete_key"
	value := []byte(`"delete_value"`)

	err = cache.Set(key, value)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify it exists
	if !cache.Has(key) {
		t.Error("expected key to exist before delete")
	}

	// Delete it
	err = cache.Delete(key)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify it doesn't exist
	if cache.Has(key) {
		t.Error("expected key to not exist after delete")
	}
}

func TestDiskCache_Clear(t *testing.T) {
	tempDir := t.TempDir()
	cache, err := NewDiskCache(tempDir, 5*time.Minute)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	// Set multiple values
	cache.Set("key1", []byte(`"value1"`))
	cache.Set("key2", []byte(`"value2"`))
	cache.Set("key3", []byte(`"value3"`))

	// Clear cache
	err = cache.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// Verify all keys are gone
	if cache.Has("key1") || cache.Has("key2") || cache.Has("key3") {
		t.Error("expected all keys to be cleared")
	}
}

func TestDiskCache_Has(t *testing.T) {
	tempDir := t.TempDir()
	cache, err := NewDiskCache(tempDir, 5*time.Minute)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	key := "has_key"

	// Should not exist initially
	if cache.Has(key) {
		t.Error("expected key to not exist initially")
	}

	// Set the key
	cache.Set(key, []byte(`"value"`))

	// Should exist now
	if !cache.Has(key) {
		t.Error("expected key to exist after set")
	}
}

func TestDiskCache_Stats(t *testing.T) {
	tempDir := t.TempDir()
	cache, err := NewDiskCache(tempDir, 5*time.Minute)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	// Set some values
	cache.Set("key1", []byte(`"value1"`))
	cache.Set("key2", []byte(`"value2"`))

	stats, err := cache.Stats()
	if err != nil {
		t.Fatalf("Stats failed: %v", err)
	}

	if stats.Total < 2 {
		t.Errorf("expected at least 2 items, got %d", stats.Total)
	}

	if stats.Active < 2 {
		t.Errorf("expected at least 2 active items, got %d", stats.Active)
	}
}

func TestDiskCache_KeyPath(t *testing.T) {
	tempDir := t.TempDir()
	cache, err := NewDiskCache(tempDir, 5*time.Minute)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	key := "test_key"
	path := cache.keyPath(key)

	// Path should be in the cache directory
	if !filepath.IsAbs(path) {
		t.Error("expected absolute path")
	}

	// Path should be within the cache directory
	if !filepath.HasPrefix(path, tempDir) {
		t.Errorf("expected path to be in cache dir %s, got %s", tempDir, path)
	}

	// Filename should have .cache extension
	if filepath.Ext(path) != ".cache" {
		t.Errorf("expected .cache extension, got %s", filepath.Ext(path))
	}
}

func TestDiskCache_Cleanup(t *testing.T) {
	tempDir := t.TempDir()
	cache, err := NewDiskCache(tempDir, 5*time.Minute)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	// Set a key with short TTL
	key := "cleanup_key"
	cache.SetWithTTL(key, []byte(`"value"`), 50*time.Millisecond)

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Manually trigger cleanup
	cache.removeExpired()

	// Key should be gone
	if cache.Has(key) {
		t.Error("expected expired key to be removed by cleanup")
	}
}

func TestDiskCache_MultipleKeys(t *testing.T) {
	tempDir := t.TempDir()
	cache, err := NewDiskCache(tempDir, 5*time.Minute)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	// Set multiple keys
	for i := 0; i < 10; i++ {
		key := "key" + string(rune(i+'0'))
		value := []byte(`"value` + string(rune(i+'0')) + `"`)
		err := cache.Set(key, value)
		if err != nil {
			t.Fatalf("Set failed for key %s: %v", key, err)
		}
	}

	// Verify all keys exist
	for i := 0; i < 10; i++ {
		key := "key" + string(rune(i+'0'))
		if !cache.Has(key) {
			t.Errorf("expected key %s to exist", key)
		}
	}
}

func TestDiskCache_InvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	cache, err := NewDiskCache(tempDir, 5*time.Minute)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	// Set invalid JSON data directly
	key := "invalid_json"
	cache.Set(key, []byte("not valid json"))

	var result map[string]string
	err = cache.GetJSON(key, &result)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestDiskCache_EmptyDir(t *testing.T) {
	tempDir := t.TempDir()

	// Create cache with empty directory
	cache, err := NewDiskCache(tempDir, 5*time.Minute)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	stats, err := cache.Stats()
	if err != nil {
		t.Fatalf("Stats failed: %v", err)
	}

	if stats.Total != 0 {
		t.Errorf("expected 0 items in new cache, got %d", stats.Total)
	}
}

func TestDiskCache_SetJSON_MarshalError(t *testing.T) {
	tempDir := t.TempDir()
	cache, err := NewDiskCache(tempDir, 5*time.Minute)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	// Channels cannot be marshaled to JSON
	invalidData := make(chan int)
	err = cache.SetJSON("key", invalidData)
	if err == nil {
		t.Error("expected error when marshaling invalid JSON data")
	}
}

func TestDiskCache_Get_CorruptedFile(t *testing.T) {
	tempDir := t.TempDir()
	cache, err := NewDiskCache(tempDir, 5*time.Minute)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	// Write corrupted data directly to disk
	path := filepath.Join(tempDir, "corrupted.cache")
	err = os.WriteFile(path, []byte("not valid json"), 0600)
	if err != nil {
		t.Fatalf("failed to write corrupted file: %v", err)
	}

	// Try to get the corrupted key using the hash
	hash := md5.Sum([]byte("corrupted"))
	filename := hex.EncodeToString(hash[:])
	cachedPath := filepath.Join(tempDir, filename+".cache")
	os.WriteFile(cachedPath, []byte("not valid json"), 0600)

	// Get should return nil for corrupted data
	result := cache.Get("corrupted")
	if result != nil {
		t.Error("expected nil for corrupted cache file")
	}
}

func TestDiskCache_Delete_NonExistent(t *testing.T) {
	tempDir := t.TempDir()
	cache, err := NewDiskCache(tempDir, 5*time.Minute)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	// Delete a non-existent key should not error
	err = cache.Delete("nonexistent")
	if err != nil {
		t.Errorf("expected no error when deleting non-existent key, got: %v", err)
	}
}

func TestDiskCache_Clear_WithSubdirs(t *testing.T) {
	tempDir := t.TempDir()
	cache, err := NewDiskCache(tempDir, 5*time.Minute)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	// Add some cache files
	cache.Set("key1", []byte("value1"))
	cache.Set("key2", []byte("value2"))

	// Create a subdirectory
	subdir := filepath.Join(tempDir, "subdir")
	os.Mkdir(subdir, 0700)

	// Clear should skip subdirectories and only remove files
	err = cache.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// Verify files are gone but subdirectory remains
	entries, _ := os.ReadDir(tempDir)
	fileCount := 0
	subdirCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			subdirCount++
		} else {
			fileCount++
		}
	}

	if fileCount != 0 {
		t.Errorf("expected 0 files, got %d", fileCount)
	}

	if subdirCount != 1 {
		t.Errorf("expected 1 subdirectory, got %d", subdirCount)
	}
}

func TestDiskCache_Stats_WithSubdirectories(t *testing.T) {
	tempDir := t.TempDir()
	cache, err := NewDiskCache(tempDir, 5*time.Minute)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	// Add some cache files
	err = cache.Set("key1", []byte(`"value1"`))
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	err = cache.Set("key2", []byte(`"value2"`))
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Give a moment for files to be written
	time.Sleep(10 * time.Millisecond)

	// Verify files exist before stats
	if !cache.Has("key1") {
		t.Error("key1 should exist")
	}
	if !cache.Has("key2") {
		t.Error("key2 should exist")
	}

	// Create subdirectory (should be skipped in stats)
	subdir := filepath.Join(tempDir, "subdir")
	os.Mkdir(subdir, 0700)

	stats, err := cache.Stats()
	if err != nil {
		t.Fatalf("Stats failed: %v", err)
	}

	// Should count only cache files, not subdirectories
	if stats.Total < 2 {
		t.Errorf("expected at least 2 files (subdirs excluded), got %d", stats.Total)
	}

	if stats.Active < 2 {
		t.Errorf("expected at least 2 active, got %d", stats.Active)
	}
}
