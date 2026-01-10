package auth

import (
	"context"

	"golang.org/x/oauth2"
)

// Provider defines the interface for authentication providers
type Provider interface {
	// Authenticate performs the authentication flow and returns a token
	Authenticate(ctx context.Context) (*oauth2.Token, error)

	// RefreshToken refreshes an existing token
	RefreshToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error)

	// ValidateToken checks if a token is valid
	ValidateToken(ctx context.Context, token *oauth2.Token) (bool, error)
}

// OAuthMode determines the OAuth flow to use
type OAuthMode int

const (
	// OAuthModeAuto tries local server first, falls back to OOB
	OAuthModeAuto OAuthMode = iota
	// OAuthModeLocal uses local callback server
	OAuthModeLocal
	// OAuthModeOOB uses out-of-band flow (manual copy-paste)
	OAuthModeOOB
)

// String returns the string representation of the OAuth mode
func (m OAuthMode) String() string {
	switch m {
	case OAuthModeAuto:
		return "auto"
	case OAuthModeLocal:
		return "local"
	case OAuthModeOOB:
		return "oob"
	default:
		return "unknown"
	}
}
