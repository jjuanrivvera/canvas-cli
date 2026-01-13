package auth

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

// RefreshBuffer is the time before expiry when we proactively refresh
// This prevents requests from failing due to token expiry during the request
const RefreshBuffer = 5 * time.Minute

// AutoRefreshTokenSource wraps an oauth2.TokenSource with automatic persistence
// of refreshed tokens. It proactively refreshes tokens before they expire.
type AutoRefreshTokenSource struct {
	oauth2Config *oauth2.Config
	store        TokenStore
	instanceName string
	token        *oauth2.Token
	mu           sync.Mutex
}

// NewAutoRefreshTokenSource creates a token source that automatically refreshes
// expired tokens and saves them to the provided store.
func NewAutoRefreshTokenSource(
	oauth2Config *oauth2.Config,
	store TokenStore,
	instanceName string,
	token *oauth2.Token,
) *AutoRefreshTokenSource {
	return &AutoRefreshTokenSource{
		oauth2Config: oauth2Config,
		store:        store,
		instanceName: instanceName,
		token:        token,
	}
}

// Token returns a valid token, refreshing if necessary.
// If the token is refreshed, it is automatically saved to storage.
func (s *AutoRefreshTokenSource) Token() (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if token is still valid with buffer time
	if s.token.Valid() && time.Until(s.token.Expiry) > RefreshBuffer {
		return s.token, nil
	}

	// Token is expired or about to expire, refresh it
	if s.token.RefreshToken == "" {
		return nil, fmt.Errorf("token expired and no refresh token available. Run 'canvas auth login' to re-authenticate")
	}

	// Create a token source for refresh
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ts := s.oauth2Config.TokenSource(ctx, s.token)
	newToken, err := ts.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w. Run 'canvas auth login' to re-authenticate", err)
	}

	// Save the refreshed token
	if err := s.store.Save(s.instanceName, newToken); err != nil {
		// Log but don't fail - token is still valid for this session
		fmt.Printf("Warning: failed to save refreshed token: %v\n", err)
	}

	s.token = newToken
	return newToken, nil
}

// GetAccessToken is a convenience method that returns just the access token string.
// This is useful when you only need the bearer token for API calls.
func (s *AutoRefreshTokenSource) GetAccessToken() (string, error) {
	token, err := s.Token()
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

// IsExpired checks if the current token is expired or about to expire.
func (s *AutoRefreshTokenSource) IsExpired() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return !s.token.Valid() || time.Until(s.token.Expiry) <= RefreshBuffer
}

// CreateOAuth2ConfigForInstance creates an oauth2.Config for token refresh operations.
// This is used when we need to refresh a token but don't have the full OAuthFlow.
func CreateOAuth2ConfigForInstance(baseURL, clientID, clientSecret string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  baseURL + "/login/oauth2/auth",
			TokenURL: baseURL + "/login/oauth2/token",
		},
	}
}
