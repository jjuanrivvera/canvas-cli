package cache

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// Cache represents an in-memory cache with TTL support
type Cache struct {
	items  map[string]*item
	mu     sync.RWMutex
	ttl    time.Duration
	stopCh chan struct{} // Channel to stop cleanup goroutine
	closed bool          // Flag to track if cache has been closed
}

// item represents a cached item with expiration
type item struct {
	value      []byte
	expiration time.Time
}

// New creates a new cache with the specified default TTL
func New(ttl time.Duration) *Cache {
	c := &Cache{
		items:  make(map[string]*item),
		ttl:    ttl,
		stopCh: make(chan struct{}),
	}

	// Start cleanup goroutine
	go c.cleanup()

	return c
}

// Close stops the cleanup goroutine and releases resources
// This should be called when the cache is no longer needed to prevent goroutine leaks
func (c *Cache) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil // Already closed
	}

	c.closed = true
	close(c.stopCh)
	return nil
}

// Get retrieves a value from the cache
// Returns nil if the key doesn't exist or has expired
func (c *Cache) Get(key string) []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil
	}

	// Check if expired
	if time.Now().After(item.expiration) {
		return nil
	}

	return item.value
}

// GetJSON retrieves and unmarshals a JSON value from the cache
func (c *Cache) GetJSON(key string, v interface{}) error {
	data := c.Get(key)
	if data == nil {
		return ErrCacheMiss
	}

	return json.Unmarshal(data, v)
}

// Set stores a value in the cache with the default TTL
func (c *Cache) Set(key string, value []byte) {
	c.SetWithTTL(key, value, c.ttl)
}

// SetJSON marshals and stores a value in the cache with the default TTL
func (c *Cache) SetJSON(key string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	c.Set(key, data)
	return nil
}

// SetWithTTL stores a value in the cache with a custom TTL
func (c *Cache) SetWithTTL(key string, value []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &item{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

// Delete removes a key from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// Clear removes all items from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*item)
}

// Has checks if a key exists and is not expired
func (c *Cache) Has(key string) bool {
	return c.Get(key) != nil
}

// Size returns the number of items in the cache (including expired ones)
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

// cleanup periodically removes expired items from the cache
func (c *Cache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.removeExpired()
		case <-c.stopCh:
			return // Stop cleanup goroutine when cache is closed
		}
	}
}

// removeExpired removes all expired items from the cache
func (c *Cache) removeExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Don't access map if cache is closed (prevents race condition)
	if c.closed {
		return
	}

	now := time.Now()
	for key, item := range c.items {
		if now.After(item.expiration) {
			delete(c.items, key)
		}
	}
}

// Stats returns cache statistics
func (c *Cache) Stats() Stats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	now := time.Now()
	var expired int

	for _, item := range c.items {
		if now.After(item.expiration) {
			expired++
		}
	}

	return Stats{
		Total:   len(c.items),
		Expired: expired,
		Active:  len(c.items) - expired,
	}
}

// Stats represents cache statistics
type Stats struct {
	Total   int // Total items in cache
	Expired int // Expired items (not yet cleaned up)
	Active  int // Active (non-expired) items
}

// ErrCacheMiss is returned when a key is not found in the cache
var ErrCacheMiss = &CacheError{message: "cache miss"}

// CacheError represents a cache-related error
type CacheError struct {
	message string
}

func (e *CacheError) Error() string {
	return e.message
}

// IsCacheMiss returns true if the error is a cache miss
// Uses errors.Is to properly handle wrapped errors
func IsCacheMiss(err error) bool {
	return errors.Is(err, ErrCacheMiss)
}

// Is implements the errors.Is interface for CacheError
func (e *CacheError) Is(target error) bool {
	t, ok := target.(*CacheError)
	if !ok {
		return false
	}
	return e.message == t.message
}
