package auth

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

// ─── FileTokenStore.Save: write error ─────────────────────────────────────────
// Makes the tokens directory read-only so os.WriteFile inside Save fails,
// covering the "failed to write token file" error return (token.go line ~136).

func TestFileTokenStore_Save_WriteError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("chmod not reliable on Windows")
	}
	dir := t.TempDir()
	store := NewFileTokenStore(dir)

	// Pre-create the tokens directory so MkdirAll inside Save succeeds,
	// then revoke write permission so WriteFile fails.
	tokenDir := filepath.Join(dir, "tokens")
	if err := os.MkdirAll(tokenDir, 0700); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.Chmod(tokenDir, 0555); err != nil {
		t.Fatalf("Chmod: %v", err)
	}
	defer os.Chmod(tokenDir, 0700) //nolint:errcheck

	token := &oauth2.Token{AccessToken: "will-fail"}
	err := store.Save("will-fail", token)
	if err == nil {
		t.Error("expected error when token directory is not writable")
	}
}

// ─── FileTokenStore.Load: decrypt failure ────────────────────────────────────
// Writes a file whose length satisfies the minimum Decrypt check (≥44 bytes)
// but whose content is not valid ciphertext, covering the "failed to decrypt
// token" error return (token.go line ~156).

func TestFileTokenStore_Load_DecryptError(t *testing.T) {
	dir := t.TempDir()
	store := NewFileTokenStore(dir)

	tokenDir := filepath.Join(dir, "tokens")
	if err := os.MkdirAll(tokenDir, 0700); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	// 44 bytes of sequential garbage — passes the minimum-length check in Decrypt
	// (saltSize=16 + nonceSize=12 + tag=16) but fails AES-GCM decryption.
	garbage := make([]byte, 44)
	for i := range garbage {
		garbage[i] = byte(i)
	}
	tokenPath := filepath.Join(tokenDir, "decrypt-fail.token.enc")
	if err := os.WriteFile(tokenPath, garbage, 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := store.Load("decrypt-fail")
	if err == nil {
		t.Error("expected error when token file cannot be decrypted")
	}
}

// ─── FallbackTokenStore.Save: behavior verification ──────────────────────────
// Exercises the FallbackTokenStore.Save path by verifying that a saved token
// is always loadable regardless of which backend stored it.
// Note: the specific "file fallback" path (lines 213-214) is only reachable
// when the system keyring is unavailable (e.g. headless CI). In environments
// where the keyring works, those lines are covered by the keyring-success path.

func TestFallbackTokenStore_Save_Roundtrip(t *testing.T) {
	dir := t.TempDir()
	fb := NewFallbackTokenStore(dir)

	token := &oauth2.Token{
		AccessToken:  "roundtrip-token",
		RefreshToken: "roundtrip-refresh",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	}

	if err := fb.Save("roundtrip-inst", token); err != nil {
		t.Fatalf("FallbackTokenStore.Save: %v", err)
	}

	loaded, err := fb.Load("roundtrip-inst")
	if err != nil {
		t.Fatalf("FallbackTokenStore.Load after Save: %v", err)
	}
	if loaded.AccessToken != token.AccessToken {
		t.Errorf("AccessToken mismatch: got %q, want %q", loaded.AccessToken, token.AccessToken)
	}

	// Cleanup.
	_ = fb.Delete("roundtrip-inst")
}

// ─── AutoRefreshTokenSource.Token: server rejects refresh ────────────────────
// Uses a mock server that always returns HTTP 400 so token refresh fails,
// covering token_source.go lines 65-67 (refresh error path).

func TestAutoRefreshTokenSource_Token_RefreshRejectedByServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid_grant"}`))
	}))
	defer server.Close()

	expiredToken := &oauth2.Token{
		AccessToken:  "old-access",
		RefreshToken: "old-refresh",              // non-empty so we bypass the "no refresh token" check
		Expiry:       time.Now().Add(-time.Hour), // expired
	}

	oauth2Config := &oauth2.Config{
		ClientID:     "client",
		ClientSecret: "secret",
		Endpoint:     oauth2.Endpoint{TokenURL: server.URL},
	}

	store := newMockTokenStore()
	ts := NewAutoRefreshTokenSource(oauth2Config, store, "rejected-inst", expiredToken)

	_, err := ts.Token()
	if err == nil {
		t.Fatal("expected error when the token server rejects the refresh request")
	}
}

// ─── OAuthFlow.ValidateToken: HTTP request fails ──────────────────────────────
// Creates a server that immediately closes the connection, so the GET request
// errors out, covering oauth.go line 292-294 (network error in ValidateToken).

func TestOAuthFlow_ValidateToken_NetworkError(t *testing.T) {
	// Start a server then immediately close it so any connection attempt fails.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srvURL := srv.URL
	srv.Close() // close immediately — subsequent requests will fail with connection refused

	cfg := &OAuthFlowConfig{
		BaseURL:  srvURL,
		ClientID: "test-client",
		Mode:     OAuthModeLocal,
	}

	flow, err := NewOAuthFlow(cfg)
	if err != nil {
		t.Fatalf("NewOAuthFlow: %v", err)
	}

	validToken := &oauth2.Token{
		AccessToken: "valid-looking-token",
		Expiry:      time.Now().Add(time.Hour),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = flow.ValidateToken(ctx, validToken)
	if err == nil {
		t.Error("expected error when HTTP connection is refused")
	}
}

// ─── OAuthFlow.startLocalServer: successful code exchange ─────────────────────
// Starts a real local server on a fixed port, drives the full callback path by
// sending a GET with the correct state + code, and verifies the returned token.
// Covers oauth.go lines 163-183 (exchange success + response body).

func TestOAuthFlow_startLocalServer_SuccessfulExchange(t *testing.T) {
	// Mock token-exchange endpoint.
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"access_token":  "exchanged-token",
			"token_type":    "Bearer",
			"expires_in":    3600,
			"refresh_token": "exchanged-refresh",
		}
		w.Header().Set("Content-Type", "application/json")
		if encErr := json.NewEncoder(w).Encode(resp); encErr != nil {
			t.Errorf("encode: %v", encErr)
		}
	}))
	defer tokenServer.Close()

	const localPort = 29876
	cfg := &OAuthFlowConfig{
		BaseURL:      tokenServer.URL,
		ClientID:     "client-abc",
		ClientSecret: "secret-xyz",
		CallbackPort: localPort,
		Mode:         OAuthModeLocal,
	}

	flow, err := NewOAuthFlow(cfg)
	if err != nil {
		t.Fatalf("NewOAuthFlow: %v", err)
	}
	// Point token exchange at the mock server.
	flow.oauth2Config.Endpoint.TokenURL = tokenServer.URL

	type result struct {
		token *oauth2.Token
		err   error
	}
	ch := make(chan result, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		tok, e := flow.startLocalServer(ctx)
		ch <- result{tok, e}
	}()

	// Poll until the local server is ready, then send the callback request.
	callbackURL := "http://localhost:29876" + callbackPath +
		"?state=" + flow.state + "&code=auth-code-123"
	var sent bool
	for i := 0; i < 100; i++ {
		time.Sleep(50 * time.Millisecond)
		resp, e := http.Get(callbackURL)
		if e == nil {
			resp.Body.Close()
			sent = true
			break
		}
	}
	if !sent {
		cancel()
		t.Fatal("local OAuth server never became ready")
	}

	res := <-ch
	if res.err != nil {
		t.Fatalf("startLocalServer returned error: %v", res.err)
	}
	if res.token == nil {
		t.Fatal("expected non-nil token from startLocalServer")
	}
	if res.token.AccessToken != "exchanged-token" {
		t.Errorf("AccessToken = %q, want 'exchanged-token'", res.token.AccessToken)
	}
}

// ─── OAuthFlow.Authenticate: Auto mode success path ─────────────────────────
// When the local server succeeds in Auto mode, Authenticate returns the token
// via the `return token, nil` branch (oauth.go line 119).

func TestOAuthFlow_Authenticate_AutoMode_LocalSucceeds(t *testing.T) {
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"access_token":  "auto-mode-token",
			"token_type":    "Bearer",
			"expires_in":    3600,
			"refresh_token": "auto-mode-refresh",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer tokenServer.Close()

	const localPort = 29878
	cfg := &OAuthFlowConfig{
		BaseURL:      tokenServer.URL,
		ClientID:     "client-auto",
		ClientSecret: "secret-auto",
		CallbackPort: localPort,
		Mode:         OAuthModeAuto, // Auto mode: tries local server first
	}

	flow, err := NewOAuthFlow(cfg)
	if err != nil {
		t.Fatalf("NewOAuthFlow: %v", err)
	}
	flow.oauth2Config.Endpoint.TokenURL = tokenServer.URL

	type result struct {
		token *oauth2.Token
		err   error
	}
	ch := make(chan result, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		tok, e := flow.Authenticate(ctx)
		ch <- result{tok, e}
	}()

	// Drive the callback once the local server is ready.
	callbackURL := "http://localhost:29878" + callbackPath +
		"?state=" + flow.state + "&code=auto-code-456"
	var sent bool
	for i := 0; i < 100; i++ {
		time.Sleep(50 * time.Millisecond)
		resp, e := http.Get(callbackURL)
		if e == nil {
			resp.Body.Close()
			sent = true
			break
		}
	}
	if !sent {
		cancel()
		t.Fatal("local OAuth server never became ready")
	}

	res := <-ch
	if res.err != nil {
		t.Fatalf("Authenticate (auto mode): %v", res.err)
	}
	if res.token == nil {
		t.Fatal("expected non-nil token from Authenticate (auto mode)")
	}
	if res.token.AccessToken != "auto-mode-token" {
		t.Errorf("AccessToken = %q, want 'auto-mode-token'", res.token.AccessToken)
	}
}

// ─── OAuthFlow.startLocalServer: port already in use ────────────────────────
// Binds a socket on the target port first, then starts startLocalServer with
// the same port. ListenAndServe fails immediately, sending the error to errChan
// (oauth.go line 188-190).

func TestOAuthFlow_startLocalServer_PortInUse(t *testing.T) {
	// Bind a TCP listener on the port before startLocalServer tries it.
	occupied, err := net.Listen("tcp", "127.0.0.1:29880")
	if err != nil {
		t.Skipf("could not pre-bind port 29880: %v", err)
	}
	defer occupied.Close()

	cfg := &OAuthFlowConfig{
		BaseURL:      "https://canvas.example.com",
		ClientID:     "port-in-use-client",
		CallbackPort: 29880,
		Mode:         OAuthModeLocal,
	}

	flow, err := NewOAuthFlow(cfg)
	if err != nil {
		t.Fatalf("NewOAuthFlow: %v", err)
	}

	// Long-lived context so ctx.Done() doesn't race with errChan.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = flow.startLocalServer(ctx)
	if err == nil {
		t.Error("expected error when port is already in use")
	}
}

// ─── OAuthFlow.startLocalServer: exchange failure ────────────────────────────
// Drives the callback with a valid state + code but the token server returns
// HTTP 400, covering oauth.go lines 168-172 (exchange error path).

func TestOAuthFlow_startLocalServer_ExchangeFailure(t *testing.T) {
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid_grant"}`))
	}))
	defer tokenServer.Close()

	const localPort = 29877
	cfg := &OAuthFlowConfig{
		BaseURL:      tokenServer.URL,
		ClientID:     "client-abc",
		ClientSecret: "secret-xyz",
		CallbackPort: localPort,
		Mode:         OAuthModeLocal,
	}

	flow, err := NewOAuthFlow(cfg)
	if err != nil {
		t.Fatalf("NewOAuthFlow: %v", err)
	}
	flow.oauth2Config.Endpoint.TokenURL = tokenServer.URL

	type result struct {
		token *oauth2.Token
		err   error
	}
	ch := make(chan result, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		tok, e := flow.startLocalServer(ctx)
		ch <- result{tok, e}
	}()

	// Wait for server to start, then trigger the exchange failure.
	callbackURL := "http://localhost:29877" + callbackPath +
		"?state=" + flow.state + "&code=bad-code"
	var sent bool
	for i := 0; i < 100; i++ {
		time.Sleep(50 * time.Millisecond)
		resp, e := http.Get(callbackURL)
		if e == nil {
			resp.Body.Close()
			sent = true
			break
		}
	}
	if !sent {
		cancel()
		t.Fatal("local OAuth server never became ready")
	}

	res := <-ch
	if res.err == nil {
		t.Fatal("expected error when token exchange fails")
	}
}
