package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

// mockTokenStore implements TokenStore for testing
type mockTokenStore struct {
	tokens map[string]*oauth2.Token
}

func newMockTokenStore() *mockTokenStore {
	return &mockTokenStore{tokens: make(map[string]*oauth2.Token)}
}

func (m *mockTokenStore) Save(instanceName string, token *oauth2.Token) error {
	m.tokens[instanceName] = token
	return nil
}

func (m *mockTokenStore) Load(instanceName string) (*oauth2.Token, error) {
	if token, ok := m.tokens[instanceName]; ok {
		return token, nil
	}
	return nil, nil
}

func (m *mockTokenStore) Delete(instanceName string) error {
	delete(m.tokens, instanceName)
	return nil
}

func (m *mockTokenStore) Exists(instanceName string) bool {
	_, ok := m.tokens[instanceName]
	return ok
}

func TestAutoRefreshTokenSource_ValidToken(t *testing.T) {
	// Create a valid token that expires in 1 hour
	token := &oauth2.Token{
		AccessToken:  "valid-access-token",
		RefreshToken: "refresh-token",
		Expiry:       time.Now().Add(1 * time.Hour),
	}

	store := newMockTokenStore()
	oauth2Config := &oauth2.Config{}

	ts := NewAutoRefreshTokenSource(oauth2Config, store, "test-instance", token)

	// Get token - should return the same token without refresh
	result, err := ts.Token()
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}

	if result.AccessToken != "valid-access-token" {
		t.Errorf("Token() AccessToken = %v, want %v", result.AccessToken, "valid-access-token")
	}
}

func TestAutoRefreshTokenSource_GetAccessToken(t *testing.T) {
	token := &oauth2.Token{
		AccessToken:  "test-access-token",
		RefreshToken: "refresh-token",
		Expiry:       time.Now().Add(1 * time.Hour),
	}

	store := newMockTokenStore()
	oauth2Config := &oauth2.Config{}

	ts := NewAutoRefreshTokenSource(oauth2Config, store, "test-instance", token)

	// Test GetAccessToken convenience method
	accessToken, err := ts.GetAccessToken()
	if err != nil {
		t.Fatalf("GetAccessToken() error = %v", err)
	}

	if accessToken != "test-access-token" {
		t.Errorf("GetAccessToken() = %v, want %v", accessToken, "test-access-token")
	}
}

func TestAutoRefreshTokenSource_IsExpired_ValidToken(t *testing.T) {
	token := &oauth2.Token{
		AccessToken:  "test-access-token",
		RefreshToken: "refresh-token",
		Expiry:       time.Now().Add(1 * time.Hour),
	}

	store := newMockTokenStore()
	oauth2Config := &oauth2.Config{}

	ts := NewAutoRefreshTokenSource(oauth2Config, store, "test-instance", token)

	if ts.IsExpired() {
		t.Errorf("IsExpired() = true for token expiring in 1 hour")
	}
}

func TestAutoRefreshTokenSource_IsExpired_ExpiredToken(t *testing.T) {
	token := &oauth2.Token{
		AccessToken:  "test-access-token",
		RefreshToken: "refresh-token",
		Expiry:       time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
	}

	store := newMockTokenStore()
	oauth2Config := &oauth2.Config{}

	ts := NewAutoRefreshTokenSource(oauth2Config, store, "test-instance", token)

	if !ts.IsExpired() {
		t.Errorf("IsExpired() = false for token expired 1 hour ago")
	}
}

func TestAutoRefreshTokenSource_IsExpired_NearExpiry(t *testing.T) {
	// Token expiring within RefreshBuffer should be considered expired
	token := &oauth2.Token{
		AccessToken:  "test-access-token",
		RefreshToken: "refresh-token",
		Expiry:       time.Now().Add(2 * time.Minute), // Expires in 2 minutes (< 5 min buffer)
	}

	store := newMockTokenStore()
	oauth2Config := &oauth2.Config{}

	ts := NewAutoRefreshTokenSource(oauth2Config, store, "test-instance", token)

	if !ts.IsExpired() {
		t.Errorf("IsExpired() = false for token expiring in 2 minutes (within buffer)")
	}
}

func TestAutoRefreshTokenSource_NoRefreshToken(t *testing.T) {
	// Create an expired token with no refresh token
	token := &oauth2.Token{
		AccessToken:  "expired-access-token",
		RefreshToken: "", // No refresh token
		Expiry:       time.Now().Add(-1 * time.Hour),
	}

	store := newMockTokenStore()
	oauth2Config := &oauth2.Config{}

	ts := NewAutoRefreshTokenSource(oauth2Config, store, "test-instance", token)

	// Should error when trying to refresh
	_, err := ts.Token()
	if err == nil {
		t.Errorf("Token() should error when expired with no refresh token")
	}
}

func TestAutoRefreshTokenSource_RefreshAndSave(t *testing.T) {
	// Create a mock OAuth2 server that returns a new token
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return a new token
		response := map[string]interface{}{
			"access_token":  "new-access-token",
			"token_type":    "Bearer",
			"expires_in":    3600,
			"refresh_token": "new-refresh-token",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create an expired token
	token := &oauth2.Token{
		AccessToken:  "old-access-token",
		RefreshToken: "old-refresh-token",
		Expiry:       time.Now().Add(-1 * time.Hour), // Expired
	}

	store := newMockTokenStore()
	oauth2Config := &oauth2.Config{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		Endpoint: oauth2.Endpoint{
			TokenURL: server.URL,
		},
	}

	ts := NewAutoRefreshTokenSource(oauth2Config, store, "test-instance", token)

	// Get token - should trigger refresh
	result, err := ts.Token()
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}

	// Check that we got the new token
	if result.AccessToken != "new-access-token" {
		t.Errorf("Token() AccessToken = %v, want %v", result.AccessToken, "new-access-token")
	}

	// Check that the token was saved to storage
	savedToken, _ := store.Load("test-instance")
	if savedToken == nil {
		t.Errorf("Token was not saved to storage")
	} else if savedToken.AccessToken != "new-access-token" {
		t.Errorf("Saved token AccessToken = %v, want %v", savedToken.AccessToken, "new-access-token")
	}
}

func TestCreateOAuth2ConfigForInstance(t *testing.T) {
	config := CreateOAuth2ConfigForInstance(
		"https://canvas.example.com",
		"client-123",
		"secret-456",
	)

	if config.ClientID != "client-123" {
		t.Errorf("ClientID = %v, want %v", config.ClientID, "client-123")
	}

	if config.ClientSecret != "secret-456" {
		t.Errorf("ClientSecret = %v, want %v", config.ClientSecret, "secret-456")
	}

	expectedAuthURL := "https://canvas.example.com/login/oauth2/auth"
	if config.Endpoint.AuthURL != expectedAuthURL {
		t.Errorf("AuthURL = %v, want %v", config.Endpoint.AuthURL, expectedAuthURL)
	}

	expectedTokenURL := "https://canvas.example.com/login/oauth2/token"
	if config.Endpoint.TokenURL != expectedTokenURL {
		t.Errorf("TokenURL = %v, want %v", config.Endpoint.TokenURL, expectedTokenURL)
	}
}
