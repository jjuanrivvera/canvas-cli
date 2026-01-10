package cache

import (
	"testing"
	"time"
)

func TestCache_SetGet(t *testing.T) {
	cache := New(5 * time.Minute)

	// Test basic set and get
	cache.Set("key1", []byte("value1"))

	value := cache.Get("key1")
	if value == nil {
		t.Fatal("expected value, got nil")
	}

	if string(value) != "value1" {
		t.Errorf("expected 'value1', got '%s'", string(value))
	}
}

func TestCache_GetNonExistent(t *testing.T) {
	cache := New(5 * time.Minute)

	value := cache.Get("nonexistent")
	if value != nil {
		t.Errorf("expected nil for non-existent key, got %v", value)
	}
}

func TestCache_SetWithTTL(t *testing.T) {
	cache := New(5 * time.Minute)

	// Set with very short TTL
	cache.SetWithTTL("shortlived", []byte("value"), 50*time.Millisecond)

	// Should exist immediately
	value := cache.Get("shortlived")
	if value == nil {
		t.Fatal("expected value immediately after set")
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Should be expired
	value = cache.Get("shortlived")
	if value != nil {
		t.Error("expected nil after TTL expiration")
	}
}

func TestCache_Delete(t *testing.T) {
	cache := New(5 * time.Minute)

	cache.Set("key1", []byte("value1"))
	cache.Delete("key1")

	value := cache.Get("key1")
	if value != nil {
		t.Error("expected nil after delete")
	}
}

func TestCache_Clear(t *testing.T) {
	cache := New(5 * time.Minute)

	cache.Set("key1", []byte("value1"))
	cache.Set("key2", []byte("value2"))
	cache.Set("key3", []byte("value3"))

	cache.Clear()

	if cache.Get("key1") != nil {
		t.Error("expected nil after clear")
	}
	if cache.Get("key2") != nil {
		t.Error("expected nil after clear")
	}
	if cache.Get("key3") != nil {
		t.Error("expected nil after clear")
	}
}

func TestCache_Has(t *testing.T) {
	cache := New(5 * time.Minute)

	cache.Set("key1", []byte("value1"))

	if !cache.Has("key1") {
		t.Error("expected Has to return true for existing key")
	}

	if cache.Has("nonexistent") {
		t.Error("expected Has to return false for non-existent key")
	}
}

func TestCache_Stats(t *testing.T) {
	cache := New(5 * time.Minute)

	cache.Set("key1", []byte("value1"))
	cache.Set("key2", []byte("value2"))

	stats := cache.Stats()

	if stats.Total != 2 {
		t.Errorf("expected total 2, got %d", stats.Total)
	}

	if stats.Active != 2 {
		t.Errorf("expected active 2, got %d", stats.Active)
	}
}

func TestCache_GetJSON(t *testing.T) {
	cache := New(5 * time.Minute)

	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	original := TestStruct{Name: "John", Age: 30}

	err := cache.SetJSON("test", original)
	if err != nil {
		t.Fatalf("SetJSON failed: %v", err)
	}

	var retrieved TestStruct
	err = cache.GetJSON("test", &retrieved)
	if err != nil {
		t.Fatalf("GetJSON failed: %v", err)
	}

	if retrieved.Name != original.Name {
		t.Errorf("expected name %s, got %s", original.Name, retrieved.Name)
	}

	if retrieved.Age != original.Age {
		t.Errorf("expected age %d, got %d", original.Age, retrieved.Age)
	}
}

func TestCache_GetJSON_NotFound(t *testing.T) {
	cache := New(5 * time.Minute)

	var data map[string]string
	err := cache.GetJSON("nonexistent", &data)

	if err != ErrCacheMiss {
		t.Errorf("expected ErrCacheMiss, got %v", err)
	}
}

func TestCache_ExpiredCleanup(t *testing.T) {
	cache := New(5 * time.Minute)

	// Set multiple items with very short TTL
	for i := 0; i < 10; i++ {
		key := "key" + string(rune('0'+i))
		cache.SetWithTTL(key, []byte("value"), 50*time.Millisecond)
	}

	stats := cache.Stats()
	if stats.Active != 10 {
		t.Errorf("expected 10 active items, got %d", stats.Active)
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Try to get items - this will trigger expiration check
	for i := 0; i < 10; i++ {
		key := "key" + string(rune('0'+i))
		value := cache.Get(key)
		if value != nil {
			t.Errorf("expected expired item to return nil")
		}
	}

	stats = cache.Stats()
	if stats.Active != 0 {
		t.Errorf("expected 0 active items after expiration, got %d", stats.Active)
	}
}

func TestCache_Size(t *testing.T) {
	cache := New(5 * time.Minute)

	// Test empty cache
	if cache.Size() != 0 {
		t.Errorf("expected size 0 for empty cache, got %d", cache.Size())
	}

	// Add some items
	cache.Set("key1", []byte("value1"))
	cache.Set("key2", []byte("value2"))
	cache.Set("key3", []byte("value3"))

	if cache.Size() != 3 {
		t.Errorf("expected size 3, got %d", cache.Size())
	}

	// Size should include expired items (until they're cleaned up)
	cache.SetWithTTL("expired", []byte("value"), 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)

	// Size still includes expired item (it's in the map)
	size := cache.Size()
	if size != 4 {
		t.Errorf("expected size 4 (including expired), got %d", size)
	}

	// Delete an item
	cache.Delete("key1")
	if cache.Size() != 3 {
		t.Errorf("expected size 3 after delete, got %d", cache.Size())
	}

	// Clear all
	cache.Clear()
	if cache.Size() != 0 {
		t.Errorf("expected size 0 after clear, got %d", cache.Size())
	}
}

func TestCache_RemoveExpired(t *testing.T) {
	cache := New(5 * time.Minute)

	// Add some non-expired items
	cache.Set("permanent1", []byte("value1"))
	cache.Set("permanent2", []byte("value2"))

	// Add items with very short TTL
	cache.SetWithTTL("short1", []byte("value1"), 1*time.Millisecond)
	cache.SetWithTTL("short2", []byte("value2"), 1*time.Millisecond)
	cache.SetWithTTL("short3", []byte("value3"), 1*time.Millisecond)

	// Wait for short-lived items to expire
	time.Sleep(5 * time.Millisecond)

	// Before cleanup, size includes all items
	sizeBefore := cache.Size()
	if sizeBefore != 5 {
		t.Errorf("expected size 5 before cleanup, got %d", sizeBefore)
	}

	// Call removeExpired (it's called internally by cleanup, but we test it directly)
	cache.removeExpired()

	// After cleanup, only permanent items remain
	sizeAfter := cache.Size()
	if sizeAfter != 2 {
		t.Errorf("expected size 2 after cleanup, got %d", sizeAfter)
	}

	// Verify permanent items still exist
	if cache.Get("permanent1") == nil {
		t.Error("expected permanent1 to still exist")
	}
	if cache.Get("permanent2") == nil {
		t.Error("expected permanent2 to still exist")
	}

	// Verify expired items are gone
	if cache.Get("short1") != nil {
		t.Error("expected short1 to be removed")
	}
}

func TestCacheError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *CacheError
		expected string
	}{
		{
			name:     "cache miss error",
			err:      &CacheError{message: "cache miss"},
			expected: "cache miss",
		},
		{
			name:     "custom error message",
			err:      &CacheError{message: "custom error"},
			expected: "custom error",
		},
		{
			name:     "empty message",
			err:      &CacheError{message: ""},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestIsCacheMiss(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "cache miss error",
			err:      &CacheError{message: "cache miss"},
			expected: true,
		},
		{
			name:     "ErrCacheMiss constant",
			err:      ErrCacheMiss,
			expected: true,
		},
		{
			name:     "regular error",
			err:      &CacheError{message: "some other error"},
			expected: true,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "non-cache error",
			err:      &time.ParseError{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsCacheMiss(tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCache_SetJSON_MarshalError(t *testing.T) {
	cache := New(5 * time.Minute)

	// Channels cannot be marshaled to JSON
	invalidData := make(chan int)
	err := cache.SetJSON("key", invalidData)
	if err == nil {
		t.Error("expected error when marshaling invalid JSON data")
	}
}

func TestCache_Close(t *testing.T) {
	cache := New(5 * time.Minute)
	cache.Set("key", []byte("value"))

	// Close should not error
	err := cache.Close()
	if err != nil {
		t.Errorf("expected nil error from Close, got %v", err)
	}

	// Closing again should not error (idempotent)
	err = cache.Close()
	if err != nil {
		t.Errorf("expected nil error from second Close, got %v", err)
	}
}

// Benchmark tests

func BenchmarkCache_Get(b *testing.B) {
	c := New(time.Hour)
	c.Set("key", []byte("value"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get("key")
	}
}

func BenchmarkCache_Set(b *testing.B) {
	c := New(time.Hour)
	value := []byte("test value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set("key", value)
	}
}

func BenchmarkCache_SetGet(b *testing.B) {
	c := New(time.Hour)
	value := []byte("test value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + string(rune('0'+i%10))
		c.Set(key, value)
		c.Get(key)
	}
}

func BenchmarkCache_GetJSON(b *testing.B) {
	c := New(time.Hour)
	type TestData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	_ = c.SetJSON("key", TestData{Name: "test", Value: 42})

	var result TestData
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.GetJSON("key", &result)
	}
}

func BenchmarkCache_SetJSON(b *testing.B) {
	c := New(time.Hour)
	type TestData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	data := TestData{Name: "test", Value: 42}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.SetJSON("key", data)
	}
}

func BenchmarkCache_Has(b *testing.B) {
	c := New(time.Hour)
	c.Set("key", []byte("value"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Has("key")
	}
}

func BenchmarkCache_ConcurrentAccess(b *testing.B) {
	c := New(time.Hour)
	c.Set("key", []byte("value"))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Get("key")
		}
	})
}
