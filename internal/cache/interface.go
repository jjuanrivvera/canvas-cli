package cache

import (
	"encoding/json"
	"time"
)

// CacheInterface defines the interface for cache implementations
type CacheInterface interface {
	Get(key string) []byte
	GetJSON(key string, v interface{}) error
	Set(key string, value []byte)
	SetJSON(key string, v interface{}) error
	SetWithTTL(key string, value []byte, ttl time.Duration)
	Delete(key string)
	Clear()
	Has(key string) bool
	Stats() Stats
}

// MultiTierCache combines in-memory and disk caching
type MultiTierCache struct {
	memory *Cache
	disk   *DiskCache
}

// NewMultiTierCache creates a new multi-tier cache
func NewMultiTierCache(memoryTTL time.Duration, diskDir string, diskTTL time.Duration) (*MultiTierCache, error) {
	disk, err := NewDiskCache(diskDir, diskTTL)
	if err != nil {
		return nil, err
	}

	return &MultiTierCache{
		memory: New(memoryTTL),
		disk:   disk,
	}, nil
}

// Get retrieves a value from the cache (memory first, then disk)
func (c *MultiTierCache) Get(key string) []byte {
	// Try memory first
	if data := c.memory.Get(key); data != nil {
		return data
	}

	// Try disk
	if data := c.disk.Get(key); data != nil {
		// Populate memory cache
		c.memory.Set(key, data)
		return data
	}

	return nil
}

// GetJSON retrieves and unmarshals a JSON value from the cache
func (c *MultiTierCache) GetJSON(key string, v interface{}) error {
	data := c.Get(key)
	if data == nil {
		return ErrCacheMiss
	}

	return json.Unmarshal(data, v)
}

// Set stores a value in both memory and disk caches
func (c *MultiTierCache) Set(key string, value []byte) {
	c.memory.Set(key, value)
	_ = c.disk.Set(key, value) // Best-effort disk cache
}

// SetJSON marshals and stores a value in both caches
func (c *MultiTierCache) SetJSON(key string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	c.Set(key, data)
	return nil
}

// SetWithTTL stores a value in both caches with a custom TTL
func (c *MultiTierCache) SetWithTTL(key string, value []byte, ttl time.Duration) {
	c.memory.SetWithTTL(key, value, ttl)
	_ = c.disk.SetWithTTL(key, value, ttl) // Best-effort disk cache
}

// Delete removes a key from both caches
func (c *MultiTierCache) Delete(key string) {
	c.memory.Delete(key)
	_ = c.disk.Delete(key) // Best-effort disk cache
}

// Clear removes all items from both caches
func (c *MultiTierCache) Clear() {
	c.memory.Clear()
	_ = c.disk.Clear() // Best-effort disk cache
}

// Has checks if a key exists in either cache
func (c *MultiTierCache) Has(key string) bool {
	return c.memory.Has(key) || c.disk.Has(key)
}

// Stats returns combined cache statistics
func (c *MultiTierCache) Stats() Stats {
	memStats := c.memory.Stats()
	diskStats, _ := c.disk.Stats()

	return Stats{
		Total:   memStats.Total + diskStats.Total,
		Expired: memStats.Expired + diskStats.Expired,
		Active:  memStats.Active + diskStats.Active,
	}
}

// MemoryStats returns memory cache statistics
func (c *MultiTierCache) MemoryStats() Stats {
	return c.memory.Stats()
}

// DiskStats returns disk cache statistics
func (c *MultiTierCache) DiskStats() (Stats, error) {
	return c.disk.Stats()
}
