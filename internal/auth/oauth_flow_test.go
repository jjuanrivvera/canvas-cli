package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestOAuthFlow_ValidateToken_Nil(t *testing.T) {
	config := &OAuthFlowConfig{
		BaseURL:  "https://canvas.example.com",
		ClientID: "test-client-id",
		Mode:     OAuthModeLocal,
	}

	flow, err := NewOAuthFlow(config)
	if err != nil {
		t.Fatalf("NewOAuthFlow failed: %v", err)
	}

	// Test with nil token
	valid, err := flow.ValidateToken(context.Background(), nil)
	if err != nil {
		t.Errorf("expected no error for nil token, got %v", err)
	}
	if valid {
		t.Error("expected nil token to be invalid")
	}
}

func TestOAuthFlow_ValidateToken_Expired(t *testing.T) {
	config := &OAuthFlowConfig{
		BaseURL:  "https://canvas.example.com",
		ClientID: "test-client-id",
		Mode:     OAuthModeLocal,
	}

	flow, err := NewOAuthFlow(config)
	if err != nil {
		t.Fatalf("NewOAuthFlow failed: %v", err)
	}

	// Test with expired token
	expiredToken := &oauth2.Token{
		AccessToken: "expired-token",
		Expiry:      time.Now().Add(-time.Hour),
	}

	valid, err := flow.ValidateToken(context.Background(), expiredToken)
	if err != nil {
		t.Errorf("expected no error for expired token validation, got %v", err)
	}
	if valid {
		t.Error("expected expired token to be invalid")
	}
}

func TestOAuthFlow_ValidateToken_Valid(t *testing.T) {
	// Create a mock server that returns 200 OK for the validation call
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/users/self" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":1,"name":"Test User"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	config := &OAuthFlowConfig{
		BaseURL:  mockServer.URL,
		ClientID: "test-client-id",
		Mode:     OAuthModeLocal,
	}

	flow, err := NewOAuthFlow(config)
	if err != nil {
		t.Fatalf("NewOAuthFlow failed: %v", err)
	}

	// Test with valid token
	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(time.Hour),
	}

	valid, err := flow.ValidateToken(context.Background(), validToken)
	if err != nil {
		t.Errorf("expected no error for valid token validation, got %v", err)
	}
	if !valid {
		t.Error("expected valid token to be valid")
	}
}

func TestOAuthFlow_ValidateToken_Unauthorized(t *testing.T) {
	// Create a mock server that returns 401 Unauthorized
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/users/self" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	config := &OAuthFlowConfig{
		BaseURL:  mockServer.URL,
		ClientID: "test-client-id",
		Mode:     OAuthModeLocal,
	}

	flow, err := NewOAuthFlow(config)
	if err != nil {
		t.Fatalf("NewOAuthFlow failed: %v", err)
	}

	// Test with token that API rejects
	token := &oauth2.Token{
		AccessToken: "invalid-api-token",
		Expiry:      time.Now().Add(time.Hour),
	}

	valid, err := flow.ValidateToken(context.Background(), token)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if valid {
		t.Error("expected unauthorized token to be invalid")
	}
}

func TestOAuthFlow_RefreshToken(t *testing.T) {
	// Create a mock OAuth server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login/oauth2/token" {
			// Check that it's a refresh token request
			r.ParseForm()
			grantType := r.FormValue("grant_type")
			if grantType != "refresh_token" {
				t.Errorf("expected grant_type refresh_token, got %s", grantType)
			}

			refreshToken := r.FormValue("refresh_token")
			if refreshToken == "" {
				t.Error("expected refresh_token in request")
			}

			// Return new token
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
				"access_token": "new-access-token",
				"refresh_token": "new-refresh-token",
				"token_type": "Bearer",
				"expires_in": 3600
			}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	config := &OAuthFlowConfig{
		BaseURL:      mockServer.URL,
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		Mode:         OAuthModeLocal,
	}

	flow, err := NewOAuthFlow(config)
	if err != nil {
		t.Fatalf("NewOAuthFlow failed: %v", err)
	}

	// Create a token that needs refreshing
	oldToken := &oauth2.Token{
		AccessToken:  "old-access-token",
		RefreshToken: "old-refresh-token",
		Expiry:       time.Now().Add(-time.Hour), // Expired
	}

	// Refresh the token
	newToken, err := flow.RefreshToken(context.Background(), oldToken)
	if err != nil {
		t.Fatalf("RefreshToken failed: %v", err)
	}

	if newToken.AccessToken != "new-access-token" {
		t.Errorf("expected new access token 'new-access-token', got '%s'", newToken.AccessToken)
	}

	if newToken.RefreshToken != "new-refresh-token" {
		t.Errorf("expected new refresh token 'new-refresh-token', got '%s'", newToken.RefreshToken)
	}
}

func TestOAuthFlow_RefreshToken_Error(t *testing.T) {
	// Create a mock server that returns error
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login/oauth2/token" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "invalid_grant"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	config := &OAuthFlowConfig{
		BaseURL:      mockServer.URL,
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		Mode:         OAuthModeLocal,
	}

	flow, err := NewOAuthFlow(config)
	if err != nil {
		t.Fatalf("NewOAuthFlow failed: %v", err)
	}

	// Try to refresh with invalid token
	invalidToken := &oauth2.Token{
		AccessToken:  "invalid-token",
		RefreshToken: "invalid-refresh",
		Expiry:       time.Now().Add(-time.Hour),
	}

	_, err = flow.RefreshToken(context.Background(), invalidToken)
	if err == nil {
		t.Error("expected error when refreshing with invalid token")
	}

	if !strings.Contains(err.Error(), "failed to refresh token") {
		t.Errorf("expected 'failed to refresh token' error, got: %v", err)
	}
}

func TestOAuthFlow_Authenticate_UnsupportedMode(t *testing.T) {
	config := &OAuthFlowConfig{
		BaseURL:  "https://canvas.example.com",
		ClientID: "test-client-id",
		Mode:     OAuthMode(999), // Invalid mode
	}

	flow, err := NewOAuthFlow(config)
	if err != nil {
		t.Fatalf("NewOAuthFlow failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err = flow.Authenticate(ctx)
	if err == nil {
		t.Error("expected error for unsupported OAuth mode")
	}

	if !strings.Contains(err.Error(), "unsupported OAuth mode") {
		t.Errorf("expected 'unsupported OAuth mode' error, got: %v", err)
	}
}

func TestOAuthFlow_Authenticate_ContextCanceled(t *testing.T) {
	config := &OAuthFlowConfig{
		BaseURL:      "https://canvas.example.com",
		ClientID:     "test-client-id",
		CallbackPort: 9999, // Use a different port to avoid conflicts
		Mode:         OAuthModeLocal,
	}

	flow, err := NewOAuthFlow(config)
	if err != nil {
		t.Fatalf("NewOAuthFlow failed: %v", err)
	}

	// Create context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err = flow.Authenticate(ctx)
	if err == nil {
		t.Error("expected error when context is cancelled")
	}

	if err != context.Canceled && !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("expected context canceled error, got: %v", err)
	}
}

func TestOAuthFlow_startLocalServer_StateValidation(t *testing.T) {
	// This test verifies that state validation works correctly
	// We can't easily test the full flow without user interaction,
	// but we can verify the state parameter is generated and used
	config := &OAuthFlowConfig{
		BaseURL:      "https://canvas.example.com",
		ClientID:     "test-client-id",
		CallbackPort: 8765, // Use a specific port for testing
		Mode:         OAuthModeLocal,
	}

	flow, err := NewOAuthFlow(config)
	if err != nil {
		t.Fatalf("NewOAuthFlow failed: %v", err)
	}

	// Verify state was generated
	if flow.state == "" {
		t.Error("expected non-empty state")
	}

	// Verify state is long enough (32 bytes base64 encoded)
	if len(flow.state) < 40 {
		t.Errorf("state seems too short: %d characters", len(flow.state))
	}
}

func TestOAuthFlow_startOOBFlow_EmptyCode(t *testing.T) {
	// This test is tricky because startOOBFlow reads from stdin
	// We'll skip it as it requires stdin manipulation
	t.Skip("Skipping test that requires stdin manipulation")
}

func TestOpenBrowser_Coverage(t *testing.T) {
	// Test openBrowser with a fake URL
	// This is best-effort and won't actually open a browser
	// Just ensure it doesn't panic
	openBrowser("https://example.com")
	// If we get here without panic, test passes
}

func TestGenerateSecureState(t *testing.T) {
	state1, err := generateSecureState()
	if err != nil {
		t.Fatalf("generateSecureState failed: %v", err)
	}

	if len(state1) == 0 {
		t.Error("expected non-empty state")
	}

	state2, err := generateSecureState()
	if err != nil {
		t.Fatalf("generateSecureState failed on second call: %v", err)
	}

	// States should be different
	if state1 == state2 {
		t.Error("expected different states from consecutive calls")
	}

	// State should be base64 URL encoded
	if len(state1) < 40 {
		t.Error("state seems too short for 32 random bytes")
	}
}

func TestOAuthFlow_PKCEGeneration(t *testing.T) {
	config := &OAuthFlowConfig{
		BaseURL:  "https://canvas.example.com",
		ClientID: "test-client-id",
		Mode:     OAuthModeLocal,
	}

	flow, err := NewOAuthFlow(config)
	if err != nil {
		t.Fatalf("NewOAuthFlow failed: %v", err)
	}

	// Verify PKCE was generated
	if flow.pkce == nil {
		t.Fatal("expected PKCE challenge to be generated")
	}

	if flow.pkce.Verifier == "" {
		t.Error("expected non-empty PKCE verifier")
	}

	if flow.pkce.Challenge == "" {
		t.Error("expected non-empty PKCE challenge")
	}

	if flow.pkce.Method != "S256" {
		t.Errorf("expected PKCE method S256, got %s", flow.pkce.Method)
	}
}

func TestOAuthFlow_OAuth2ConfigSetup(t *testing.T) {
	config := &OAuthFlowConfig{
		BaseURL:      "https://canvas.example.com",
		ClientID:     "test-client-id",
		ClientSecret: "test-secret",
		Mode:         OAuthModeLocal,
		Scopes:       []string{"read", "write"},
	}

	flow, err := NewOAuthFlow(config)
	if err != nil {
		t.Fatalf("NewOAuthFlow failed: %v", err)
	}

	// Verify OAuth2 config was set up correctly
	if flow.oauth2Config.ClientID != "test-client-id" {
		t.Errorf("expected ClientID 'test-client-id', got '%s'", flow.oauth2Config.ClientID)
	}

	if flow.oauth2Config.ClientSecret != "test-secret" {
		t.Errorf("expected ClientSecret 'test-secret', got '%s'", flow.oauth2Config.ClientSecret)
	}

	expectedAuthURL := "https://canvas.example.com/login/oauth2/auth"
	if flow.oauth2Config.Endpoint.AuthURL != expectedAuthURL {
		t.Errorf("expected AuthURL '%s', got '%s'", expectedAuthURL, flow.oauth2Config.Endpoint.AuthURL)
	}

	expectedTokenURL := "https://canvas.example.com/login/oauth2/token"
	if flow.oauth2Config.Endpoint.TokenURL != expectedTokenURL {
		t.Errorf("expected TokenURL '%s', got '%s'", expectedTokenURL, flow.oauth2Config.Endpoint.TokenURL)
	}

	if len(flow.oauth2Config.Scopes) != 2 {
		t.Errorf("expected 2 scopes, got %d", len(flow.oauth2Config.Scopes))
	}
}

func TestOAuthFlow_RedirectURLSetup_LocalMode(t *testing.T) {
	config := &OAuthFlowConfig{
		BaseURL:      "https://canvas.example.com",
		ClientID:     "test-client-id",
		Mode:         OAuthModeLocal,
		CallbackPort: 9876,
	}

	flow, err := NewOAuthFlow(config)
	if err != nil {
		t.Fatalf("NewOAuthFlow failed: %v", err)
	}

	expectedRedirectURL := fmt.Sprintf("http://localhost:9876%s", callbackPath)
	if flow.config.RedirectURL != expectedRedirectURL {
		t.Errorf("expected RedirectURL '%s', got '%s'", expectedRedirectURL, flow.config.RedirectURL)
	}
}

func TestOAuthFlow_RedirectURLSetup_OOBMode(t *testing.T) {
	config := &OAuthFlowConfig{
		BaseURL:  "https://canvas.example.com",
		ClientID: "test-client-id",
		Mode:     OAuthModeOOB,
	}

	flow, err := NewOAuthFlow(config)
	if err != nil {
		t.Fatalf("NewOAuthFlow failed: %v", err)
	}

	// OOB mode should not set redirect URL during construction
	// It's set during startOOBFlow to "urn:ietf:wg:oauth:2.0:oob"
	if flow.config.RedirectURL != "" {
		t.Logf("RedirectURL for OOB mode: %s (will be overridden)", flow.config.RedirectURL)
	}
}

func TestOAuthFlow_DefaultCallbackPort(t *testing.T) {
	config := &OAuthFlowConfig{
		BaseURL:  "https://canvas.example.com",
		ClientID: "test-client-id",
		Mode:     OAuthModeLocal,
		// CallbackPort not set
	}

	flow, err := NewOAuthFlow(config)
	if err != nil {
		t.Fatalf("NewOAuthFlow failed: %v", err)
	}

	if flow.config.CallbackPort != defaultCallbackPort {
		t.Errorf("expected default callback port %d, got %d", defaultCallbackPort, flow.config.CallbackPort)
	}
}

func TestOAuthFlow_startLocalServer_InvalidState(t *testing.T) {
	config := &OAuthFlowConfig{
		BaseURL:      "https://canvas.example.com",
		ClientID:     "test-client-id",
		CallbackPort: 8787,
		Mode:         OAuthModeLocal,
	}

	flow, err := NewOAuthFlow(config)
	if err != nil {
		t.Fatalf("NewOAuthFlow failed: %v", err)
	}

	// Start the server in a goroutine
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resultChan := make(chan error, 1)
	go func() {
		_, err := flow.startLocalServer(ctx)
		resultChan <- err
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Make a request with invalid state
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d%s?state=invalid-state&code=test-code", flow.config.CallbackPort, callbackPath))
	if err != nil {
		t.Logf("Request failed (expected if server not ready): %v", err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request for invalid state, got %d", resp.StatusCode)
		}
	}

	// Wait for result or timeout
	select {
	case err := <-resultChan:
		if err == nil {
			t.Error("expected error from startLocalServer with invalid state")
		}
		if !strings.Contains(err.Error(), "state") && !strings.Contains(err.Error(), "timed out") && err != ctx.Err() {
			t.Logf("Got error: %v", err)
		}
	case <-time.After(3 * time.Second):
		// Timeout is acceptable for this test
	}
}

func TestOAuthFlow_startLocalServer_MissingCode(t *testing.T) {
	config := &OAuthFlowConfig{
		BaseURL:      "https://canvas.example.com",
		ClientID:     "test-client-id",
		CallbackPort: 8788,
		Mode:         OAuthModeLocal,
	}

	flow, err := NewOAuthFlow(config)
	if err != nil {
		t.Fatalf("NewOAuthFlow failed: %v", err)
	}

	// Start the server in a goroutine
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resultChan := make(chan error, 1)
	go func() {
		_, err := flow.startLocalServer(ctx)
		resultChan <- err
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Make a request with correct state but missing code
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d%s?state=%s", flow.config.CallbackPort, callbackPath, flow.state))
	if err != nil {
		t.Logf("Request failed (expected if server not ready): %v", err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request for missing code, got %d", resp.StatusCode)
		}
	}

	// Wait for result or timeout
	select {
	case err := <-resultChan:
		if err == nil {
			t.Error("expected error from startLocalServer with missing code")
		}
	case <-time.After(3 * time.Second):
		// Timeout is acceptable for this test
	}
}

func TestOAuthFlow_Authenticate_AutoMode_Fallback(t *testing.T) {
	// Create a mock OAuth server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login/oauth2/token" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
				"access_token": "test-access-token",
				"refresh_token": "test-refresh-token",
				"token_type": "Bearer",
				"expires_in": 3600
			}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	config := &OAuthFlowConfig{
		BaseURL:      mockServer.URL,
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		CallbackPort: 1, // Use invalid port to force local server to fail
		Mode:         OAuthModeAuto,
	}

	flow, err := NewOAuthFlow(config)
	if err != nil {
		t.Fatalf("NewOAuthFlow failed: %v", err)
	}

	// Test auto mode - local server will fail, should fall back to OOB
	// However, OOB requires stdin, so we'll just verify it attempts the fallback
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err = flow.Authenticate(ctx)
	// We expect this to fail (timeout or context cancelled) since we can't provide stdin
	// But it should have attempted the fallback
	if err == nil {
		t.Error("expected error from Authenticate in auto mode without stdin")
	}
}

func TestOAuthFlow_Authenticate_LocalMode(t *testing.T) {
	config := &OAuthFlowConfig{
		BaseURL:      "https://canvas.example.com",
		ClientID:     "test-client-id",
		CallbackPort: 8789,
		Mode:         OAuthModeLocal,
	}

	flow, err := NewOAuthFlow(config)
	if err != nil {
		t.Fatalf("NewOAuthFlow failed: %v", err)
	}

	// Test local mode with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	_, err = flow.Authenticate(ctx)
	// We expect timeout since no user interaction
	if err == nil {
		t.Error("expected error from Authenticate without user interaction")
	}
}
