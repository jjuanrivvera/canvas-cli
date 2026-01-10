package cache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// DiskCache represents a disk-based cache with TTL support
type DiskCache struct {
	dir string
	ttl time.Duration
}

// NewDiskCache creates a new disk-based cache
func NewDiskCache(dir string, ttl time.Duration) (*DiskCache, error) {
	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	c := &DiskCache{
		dir: dir,
		ttl: ttl,
	}

	// Clean up expired files on startup
	go c.cleanup()

	return c, nil
}

// diskItem represents a cached item stored on disk
type diskItem struct {
	Value      json.RawMessage `json:"value"`
	Expiration time.Time       `json:"expiration"`
}

// Get retrieves a value from the disk cache
func (c *DiskCache) Get(key string) []byte {
	path := c.keyPath(key)

	// Read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	// Unmarshal the item
	var item diskItem
	if err := json.Unmarshal(data, &item); err != nil {
		return nil
	}

	// Check if expired
	if time.Now().After(item.Expiration) {
		// Delete expired file
		os.Remove(path)
		return nil
	}

	return []byte(item.Value)
}

// GetJSON retrieves and unmarshals a JSON value from the disk cache
func (c *DiskCache) GetJSON(key string, v interface{}) error {
	data := c.Get(key)
	if data == nil {
		return ErrCacheMiss
	}

	return json.Unmarshal(data, v)
}

// Set stores a value in the disk cache with the default TTL
func (c *DiskCache) Set(key string, value []byte) error {
	return c.SetWithTTL(key, value, c.ttl)
}

// SetJSON marshals and stores a value in the disk cache with the default TTL
func (c *DiskCache) SetJSON(key string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return c.Set(key, data)
}

// SetWithTTL stores a value in the disk cache with a custom TTL
func (c *DiskCache) SetWithTTL(key string, value []byte, ttl time.Duration) error {
	path := c.keyPath(key)

	// Create the item
	item := diskItem{
		Value:      json.RawMessage(value),
		Expiration: time.Now().Add(ttl),
	}

	// Marshal the item
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal cache item: %w", err)
	}

	// Write to disk
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

// Delete removes a key from the disk cache
func (c *DiskCache) Delete(key string) error {
	path := c.keyPath(key)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Clear removes all items from the disk cache
func (c *DiskCache) Clear() error {
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			os.Remove(filepath.Join(c.dir, entry.Name()))
		}
	}

	return nil
}

// Has checks if a key exists and is not expired
func (c *DiskCache) Has(key string) bool {
	return c.Get(key) != nil
}

// cleanup periodically removes expired items from the disk cache
func (c *DiskCache) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.removeExpired()
	}
}

// removeExpired removes all expired items from the disk cache
func (c *DiskCache) removeExpired() {
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return
	}

	now := time.Now()

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(c.dir, entry.Name())

		// Read the file
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		// Unmarshal the item
		var item diskItem
		if err := json.Unmarshal(data, &item); err != nil {
			// Remove corrupted files
			os.Remove(path)
			continue
		}

		// Delete if expired
		if now.After(item.Expiration) {
			os.Remove(path)
		}
	}
}

// keyPath converts a cache key to a file path
func (c *DiskCache) keyPath(key string) string {
	// Hash the key to create a valid filename
	hash := md5.Sum([]byte(key))
	filename := hex.EncodeToString(hash[:])
	return filepath.Join(c.dir, filename+".cache")
}

// Stats returns disk cache statistics
func (c *DiskCache) Stats() (Stats, error) {
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return Stats{}, err
	}

	now := time.Now()
	var total, expired, active int

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		total++

		path := filepath.Join(c.dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var item diskItem
		if err := json.Unmarshal(data, &item); err != nil {
			continue
		}

		if now.After(item.Expiration) {
			expired++
		} else {
			active++
		}
	}

	return Stats{
		Total:   total,
		Expired: expired,
		Active:  active,
	}, nil
}
