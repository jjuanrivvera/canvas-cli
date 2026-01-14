package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/oauth2"
)

// generateSecureState generates a cryptographically secure random state parameter
// to prevent CSRF attacks during OAuth flow
func generateSecureState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random state: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

const (
	defaultCallbackPort = 8080
	callbackPath        = "/oauth/callback"
	defaultTimeout      = 5 * time.Minute
)

// OAuthFlowConfig holds configuration for the OAuth flow
type OAuthFlowConfig struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	Mode         OAuthMode
	CallbackPort int
	Logger       *slog.Logger
}

// OAuthFlow implements the OAuth 2.0 flow with PKCE
type OAuthFlow struct {
	config       *OAuthFlowConfig
	oauth2Config *oauth2.Config
	pkce         *PKCEChallenge
	state        string // Secure random state for CSRF protection
	logger       *slog.Logger
}

// NewOAuthFlow creates a new OAuth flow
func NewOAuthFlow(config *OAuthFlowConfig) (*OAuthFlow, error) {
	if config.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	if config.ClientID == "" {
		return nil, fmt.Errorf("client ID is required")
	}

	if config.Mode == OAuthModeLocal || config.Mode == OAuthModeAuto {
		if config.CallbackPort == 0 {
			config.CallbackPort = defaultCallbackPort
		}
		if config.RedirectURL == "" {
			config.RedirectURL = fmt.Sprintf("http://localhost:%d%s", config.CallbackPort, callbackPath)
		}
	}

	if config.Logger == nil {
		config.Logger = slog.Default()
	}

	// Generate PKCE challenge
	pkce, err := GeneratePKCEChallenge()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PKCE challenge: %w", err)
	}

	// Generate secure random state for CSRF protection
	state, err := generateSecureState()
	if err != nil {
		return nil, fmt.Errorf("failed to generate state: %w", err)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.BaseURL + "/login/oauth2/auth",
			TokenURL: config.BaseURL + "/login/oauth2/token",
		},
		RedirectURL: config.RedirectURL,
		Scopes:      config.Scopes,
	}

	return &OAuthFlow{
		config:       config,
		oauth2Config: oauth2Config,
		pkce:         pkce,
		state:        state,
		logger:       config.Logger,
	}, nil
}

// Authenticate performs the OAuth flow
func (f *OAuthFlow) Authenticate(ctx context.Context) (*oauth2.Token, error) {
	switch f.config.Mode {
	case OAuthModeAuto:
		// Try local server first
		token, err := f.startLocalServer(ctx)
		if err != nil {
			f.logger.Info("Local OAuth server failed, falling back to out-of-band flow", "error", err)
			return f.startOOBFlow(ctx)
		}
		return token, nil
	case OAuthModeLocal:
		return f.startLocalServer(ctx)
	case OAuthModeOOB:
		return f.startOOBFlow(ctx)
	default:
		return nil, fmt.Errorf("unsupported OAuth mode: %v", f.config.Mode)
	}
}

// startLocalServer starts a local HTTP server for OAuth callback
func (f *OAuthFlow) startLocalServer(ctx context.Context) (*oauth2.Token, error) {
	// Create channel for result
	resultChan := make(chan *oauth2.Token, 1)
	errChan := make(chan error, 1)

	// Create HTTP server
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", f.config.CallbackPort),
		Handler: mux,
	}

	// Handle callback
	mux.HandleFunc(callbackPath, func(w http.ResponseWriter, r *http.Request) {
		// Validate state parameter to prevent CSRF attacks
		returnedState := r.URL.Query().Get("state")
		if returnedState != f.state {
			err := fmt.Errorf("invalid state parameter - possible CSRF attack")
			errChan <- err
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			err := fmt.Errorf("authorization code not found in callback")
			errChan <- err
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Exchange code for token with PKCE verifier
		// Use request context for proper cancellation propagation
		token, err := f.oauth2Config.Exchange(
			r.Context(),
			code,
			oauth2.SetAuthURLParam("code_verifier", f.pkce.Verifier),
		)
		if err != nil {
			errChan <- fmt.Errorf("failed to exchange code for token: %w", err)
			http.Error(w, "Authentication failed", http.StatusInternalServerError)
			return
		}

		resultChan <- token
		fmt.Fprintf(w, `
			<html>
				<head><title>Authentication Successful</title></head>
				<body>
					<h1>✓ Authentication Successful</h1>
					<p>You can close this window and return to the CLI.</p>
				</body>
			</html>
		`)
	})

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("failed to start local server: %w", err)
		}
	}()

	// Ensure server is shut down
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)
	}()

	// Generate authorization URL with secure state
	authURL := f.oauth2Config.AuthCodeURL(
		f.state,
		oauth2.SetAuthURLParam("code_challenge", f.pkce.Challenge),
		oauth2.SetAuthURLParam("code_challenge_method", f.pkce.Method),
	)

	f.logger.Info("Opening browser for authentication...")
	f.logger.Info("If browser doesn't open, visit this URL:", "url", authURL)
	fmt.Printf("\n🔐 Opening browser for Canvas authentication...\n")
	fmt.Printf("If your browser doesn't open automatically, visit:\n%s\n\n", authURL)

	// Try to open browser (best effort)
	openBrowser(authURL)

	// Wait for result or timeout
	select {
	case token := <-resultChan:
		return token, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(defaultTimeout):
		return nil, fmt.Errorf("authentication timed out after %v", defaultTimeout)
	}
}

// startOOBFlow starts the out-of-band flow (manual copy-paste)
func (f *OAuthFlow) startOOBFlow(ctx context.Context) (*oauth2.Token, error) {
	// For OOB, we use a special redirect URI
	f.oauth2Config.RedirectURL = "urn:ietf:wg:oauth:2.0:oob"

	// Generate authorization URL with secure state
	// Note: OOB flow can't validate state server-side, but we include it for consistency
	authURL := f.oauth2Config.AuthCodeURL(
		f.state,
		oauth2.SetAuthURLParam("code_challenge", f.pkce.Challenge),
		oauth2.SetAuthURLParam("code_challenge_method", f.pkce.Method),
	)

	fmt.Printf("\n🔐 Canvas OAuth Authentication (Out-of-Band Mode)\n\n")
	fmt.Printf("1. Visit this URL in your browser:\n%s\n\n", authURL)
	fmt.Printf("2. Authorize the application\n")
	fmt.Printf("3. Copy the authorization code from the page\n")
	fmt.Printf("4. Paste the code here: ")

	// Read authorization code from stdin
	var code string
	fmt.Scanln(&code)

	if code == "" {
		return nil, fmt.Errorf("authorization code is required")
	}

	// Exchange code for token with PKCE verifier
	token, err := f.oauth2Config.Exchange(
		ctx,
		code,
		oauth2.SetAuthURLParam("code_verifier", f.pkce.Verifier),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	return token, nil
}

// RefreshToken refreshes an OAuth token
func (f *OAuthFlow) RefreshToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	tokenSource := f.oauth2Config.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	return newToken, nil
}

// ValidateToken checks if a token is valid by making a test API call
func (f *OAuthFlow) ValidateToken(ctx context.Context, token *oauth2.Token) (bool, error) {
	if token == nil {
		return false, nil
	}

	// Check if token is expired
	if token.Expiry.Before(time.Now()) {
		return false, nil
	}

	// Make a test API call to verify the token works
	client := f.oauth2Config.Client(ctx, token)
	resp, err := client.Get(f.config.BaseURL + "/api/v1/users/self")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// openBrowser attempts to open the URL in the default browser
// This is a best-effort attempt - errors are silently ignored
// The user can always use the printed URL if this fails
func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		// Try xdg-open first, fall back to sensible-browser
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		// Unknown OS, silently fail
		return
	}

	// Run in background, ignore errors (best effort)
	_ = cmd.Start()
}
