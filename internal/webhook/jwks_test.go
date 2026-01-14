package webhook

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// generateTestRSAKey generates a test RSA key pair
func generateTestRSAKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}
	return key
}

// createTestJWKServer creates a test server that serves JWKs
func createTestJWKServer(t *testing.T, keys map[string]*rsa.PublicKey) *httptest.Server {
	t.Helper()

	jwks := jwkResponse{Keys: make([]jwk, 0, len(keys))}
	for kid, key := range keys {
		jwks.Keys = append(jwks.Keys, jwk{
			Kty: "RSA",
			Kid: kid,
			Use: "sig",
			Alg: "RS256",
			N:   base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
			E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes()),
		})
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(jwks); err != nil {
			t.Logf("Failed to encode JWKs: %v", err)
		}
	}))
}

func TestNewJWKSet(t *testing.T) {
	url := "https://example.com/jwks"
	jwkSet := NewJWKSet(url)

	if jwkSet.url != url {
		t.Errorf("Expected URL %q, got %q", url, jwkSet.url)
	}

	if jwkSet.keys == nil {
		t.Error("Expected keys map to be initialized")
	}

	if jwkSet.ttl != 1*time.Hour {
		t.Errorf("Expected TTL of 1 hour, got %v", jwkSet.ttl)
	}

	if jwkSet.client == nil {
		t.Error("Expected HTTP client to be initialized")
	}
}

func TestJWKSet_Refresh(t *testing.T) {
	// Generate test keys
	key1 := generateTestRSAKey(t)
	key2 := generateTestRSAKey(t)

	keys := map[string]*rsa.PublicKey{
		"key-1": &key1.PublicKey,
		"key-2": &key2.PublicKey,
	}

	server := createTestJWKServer(t, keys)
	defer server.Close()

	jwkSet := NewJWKSet(server.URL)

	// Refresh should succeed
	err := jwkSet.Refresh()
	if err != nil {
		t.Fatalf("Refresh failed: %v", err)
	}

	// Should have 2 keys
	if jwkSet.KeyCount() != 2 {
		t.Errorf("Expected 2 keys, got %d", jwkSet.KeyCount())
	}

	// lastFetch should be updated
	if jwkSet.lastFetch.IsZero() {
		t.Error("lastFetch should be updated after refresh")
	}
}

func TestJWKSet_GetKey(t *testing.T) {
	// Generate test key
	privateKey := generateTestRSAKey(t)

	keys := map[string]*rsa.PublicKey{
		"test-key-id": &privateKey.PublicKey,
	}

	server := createTestJWKServer(t, keys)
	defer server.Close()

	jwkSet := NewJWKSet(server.URL)

	// Get key should succeed
	pubKey, err := jwkSet.GetKey("test-key-id")
	if err != nil {
		t.Fatalf("GetKey failed: %v", err)
	}

	if pubKey == nil {
		t.Fatal("Expected public key, got nil")
	}

	// Verify it's the same key
	if pubKey.N.Cmp(privateKey.PublicKey.N) != 0 {
		t.Error("Retrieved key doesn't match original")
	}
}

func TestJWKSet_GetKey_NotFound(t *testing.T) {
	// Generate test key
	privateKey := generateTestRSAKey(t)

	keys := map[string]*rsa.PublicKey{
		"existing-key": &privateKey.PublicKey,
	}

	server := createTestJWKServer(t, keys)
	defer server.Close()

	jwkSet := NewJWKSet(server.URL)

	// Get non-existent key should fail
	_, err := jwkSet.GetKey("non-existent-key")
	if err == nil {
		t.Error("Expected error for non-existent key")
	}
}

func TestJWKSet_GetKey_CacheHit(t *testing.T) {
	// Generate test key
	privateKey := generateTestRSAKey(t)

	keys := map[string]*rsa.PublicKey{
		"cached-key": &privateKey.PublicKey,
	}

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		jwks := jwkResponse{Keys: make([]jwk, 0, len(keys))}
		for kid, key := range keys {
			jwks.Keys = append(jwks.Keys, jwk{
				Kty: "RSA",
				Kid: kid,
				Use: "sig",
				Alg: "RS256",
				N:   base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
				E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes()),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jwks)
	}))
	defer server.Close()

	jwkSet := NewJWKSet(server.URL)

	// First call should fetch
	_, err := jwkSet.GetKey("cached-key")
	if err != nil {
		t.Fatalf("First GetKey failed: %v", err)
	}

	if requestCount != 1 {
		t.Errorf("Expected 1 request, got %d", requestCount)
	}

	// Second call should use cache
	_, err = jwkSet.GetKey("cached-key")
	if err != nil {
		t.Fatalf("Second GetKey failed: %v", err)
	}

	if requestCount != 1 {
		t.Errorf("Expected 1 request (cached), got %d", requestCount)
	}
}

func TestJWKSet_Refresh_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	jwkSet := NewJWKSet(server.URL)

	err := jwkSet.Refresh()
	if err == nil {
		t.Error("Expected error for server error response")
	}
}

func TestJWKSet_Refresh_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	jwkSet := NewJWKSet(server.URL)

	err := jwkSet.Refresh()
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestJWKSet_Refresh_NoRSAKeys(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Return empty keys array
		_ = json.NewEncoder(w).Encode(jwkResponse{Keys: []jwk{}})
	}))
	defer server.Close()

	jwkSet := NewJWKSet(server.URL)

	err := jwkSet.Refresh()
	if err == nil {
		t.Error("Expected error when no RSA keys found")
	}
}

func TestParseRSAPublicKey(t *testing.T) {
	// Generate a real RSA key for testing
	privateKey := generateTestRSAKey(t)
	pubKey := &privateKey.PublicKey

	// Create JWK from the public key
	testJWK := jwk{
		Kty: "RSA",
		Kid: "test",
		Use: "sig",
		Alg: "RS256",
		N:   base64.RawURLEncoding.EncodeToString(pubKey.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pubKey.E)).Bytes()),
	}

	// Parse it back
	parsedKey, err := parseRSAPublicKey(testJWK)
	if err != nil {
		t.Fatalf("parseRSAPublicKey failed: %v", err)
	}

	// Verify the parsed key matches
	if parsedKey.N.Cmp(pubKey.N) != 0 {
		t.Error("Parsed modulus doesn't match original")
	}

	if parsedKey.E != pubKey.E {
		t.Errorf("Parsed exponent %d doesn't match original %d", parsedKey.E, pubKey.E)
	}
}

func TestParseRSAPublicKey_InvalidModulus(t *testing.T) {
	testJWK := jwk{
		Kty: "RSA",
		Kid: "test",
		N:   "!!!invalid-base64!!!",
		E:   "AQAB",
	}

	_, err := parseRSAPublicKey(testJWK)
	if err == nil {
		t.Error("Expected error for invalid modulus encoding")
	}
}

func TestParseRSAPublicKey_InvalidExponent(t *testing.T) {
	testJWK := jwk{
		Kty: "RSA",
		Kid: "test",
		N:   base64.RawURLEncoding.EncodeToString([]byte("valid")),
		E:   "!!!invalid-base64!!!",
	}

	_, err := parseRSAPublicKey(testJWK)
	if err == nil {
		t.Error("Expected error for invalid exponent encoding")
	}
}
