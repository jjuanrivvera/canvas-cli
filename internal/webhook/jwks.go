package webhook

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"
)

// CanvasDataServicesJWKURL is the official Canvas Data Services JWK endpoint
const CanvasDataServicesJWKURL = "https://8axpcl50e4.execute-api.us-east-1.amazonaws.com/main/jwks"

// JWKSet manages Canvas JWK public keys for JWT verification
type JWKSet struct {
	url       string
	keys      map[string]*rsa.PublicKey
	mu        sync.RWMutex
	lastFetch time.Time
	ttl       time.Duration
	client    *http.Client
}

// jwkResponse represents the JSON response from a JWK endpoint
type jwkResponse struct {
	Keys []jwk `json:"keys"`
}

// jwk represents a single JSON Web Key
type jwk struct {
	Kty string `json:"kty"` // Key type (RSA)
	Kid string `json:"kid"` // Key ID
	Use string `json:"use"` // Key usage (sig)
	Alg string `json:"alg"` // Algorithm (RS256)
	N   string `json:"n"`   // RSA modulus (base64url)
	E   string `json:"e"`   // RSA exponent (base64url)
}

// NewJWKSet creates a new JWK set fetcher
func NewJWKSet(url string) *JWKSet {
	return &JWKSet{
		url:  url,
		keys: make(map[string]*rsa.PublicKey),
		ttl:  1 * time.Hour, // Cache keys for 1 hour (Canvas rotates monthly)
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetKey returns the public key for the given key ID
// It will fetch/refresh keys if not cached or cache is stale
func (j *JWKSet) GetKey(kid string) (*rsa.PublicKey, error) {
	j.mu.RLock()
	key, exists := j.keys[kid]
	stale := time.Since(j.lastFetch) > j.ttl
	j.mu.RUnlock()

	// Return cached key if exists and not stale
	if exists && !stale {
		return key, nil
	}

	// Refresh keys
	if err := j.Refresh(); err != nil {
		// If we have a cached key but refresh failed, use the cached one
		if exists {
			return key, nil
		}
		return nil, fmt.Errorf("failed to refresh JWKs: %w", err)
	}

	// Look up key after refresh
	j.mu.RLock()
	key, exists = j.keys[kid]
	j.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("key with ID %q not found in JWK set", kid)
	}

	return key, nil
}

// Refresh fetches the latest JWKs from the configured URL
func (j *JWKSet) Refresh() error {
	resp, err := j.client.Get(j.url)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("JWK endpoint returned status %d", resp.StatusCode)
	}

	var jwkResp jwkResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwkResp); err != nil {
		return fmt.Errorf("failed to decode JWK response: %w", err)
	}

	// Parse and store keys
	newKeys := make(map[string]*rsa.PublicKey)
	for _, key := range jwkResp.Keys {
		if key.Kty != "RSA" {
			continue // Only support RSA keys
		}

		pubKey, err := parseRSAPublicKey(key)
		if err != nil {
			continue // Skip invalid keys
		}

		newKeys[key.Kid] = pubKey
	}

	if len(newKeys) == 0 {
		return errors.New("no valid RSA keys found in JWK set")
	}

	// Update cache
	j.mu.Lock()
	j.keys = newKeys
	j.lastFetch = time.Now()
	j.mu.Unlock()

	return nil
}

// KeyCount returns the number of cached keys
func (j *JWKSet) KeyCount() int {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return len(j.keys)
}

// parseRSAPublicKey converts a JWK to an RSA public key
func parseRSAPublicKey(key jwk) (*rsa.PublicKey, error) {
	// Decode modulus (n)
	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %w", err)
	}

	// Decode exponent (e)
	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %w", err)
	}

	// Convert exponent bytes to int
	var e int
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: e,
	}, nil
}
