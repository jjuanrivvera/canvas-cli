package auth

// Tests for public-client (secret-less) OAuth. Canvas developer keys with
// client_type = "public" exchange and refresh tokens with PKCE only, and the
// server only takes the public-client path when client_secret is absent from
// the request: an empty client_secret param or a Basic auth header with an
// empty password would still be treated as a (failed) confidential exchange.

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestCanvasAuthStyle(t *testing.T) {
	if got := canvasAuthStyle(""); got != oauth2.AuthStyleInParams {
		t.Errorf("canvasAuthStyle(\"\") = %v, want AuthStyleInParams", got)
	}
	if got := canvasAuthStyle("some-secret"); got != oauth2.AuthStyleAutoDetect {
		t.Errorf("canvasAuthStyle(secret) = %v, want AuthStyleAutoDetect", got)
	}
}

// newTokenRecorder returns a test server that records the form and auth
// header of requests to /login/oauth2/token and responds with a valid token.
func newTokenRecorder(t *testing.T, gotForm *url.Values, gotAuthHeader *string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/login/oauth2/token" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err := r.ParseForm(); err != nil {
			t.Errorf("ParseForm failed: %v", err)
		}
		*gotForm = r.Form
		*gotAuthHeader = r.Header.Get("Authorization")

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"access_token": "public-access-token",
			"refresh_token": "public-refresh-token",
			"token_type": "Bearer",
			"expires_in": 7200
		}`)
	}))
}

func TestOAuthFlow_PublicClient_ExchangeOmitsSecret(t *testing.T) {
	var gotForm url.Values
	var gotAuthHeader string
	mockServer := newTokenRecorder(t, &gotForm, &gotAuthHeader)
	defer mockServer.Close()

	config := &OAuthFlowConfig{
		BaseURL:      mockServer.URL,
		ClientID:     "public-client-id",
		CallbackPort: 18790,
		Mode:         OAuthModeLocal,
	}

	flow, err := NewOAuthFlow(config)
	if err != nil {
		t.Fatalf("NewOAuthFlow failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resultChan := make(chan *oauth2.Token, 1)
	errChan := make(chan error, 1)
	go func() {
		token, err := flow.startLocalServer(ctx)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- token
	}()

	// Give the callback server time to start, then complete the flow with a
	// valid state so the code is exchanged against the mock token endpoint.
	time.Sleep(100 * time.Millisecond)
	resp, err := http.Get(fmt.Sprintf(
		"http://localhost:%d%s?state=%s&code=test-auth-code",
		config.CallbackPort, callbackPath, url.QueryEscape(flow.state),
	))
	if err != nil {
		t.Fatalf("callback request failed: %v", err)
	}
	resp.Body.Close()

	select {
	case token := <-resultChan:
		if token.AccessToken != "public-access-token" {
			t.Errorf("access token = %q, want %q", token.AccessToken, "public-access-token")
		}
	case err := <-errChan:
		t.Fatalf("startLocalServer failed: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for token exchange")
	}

	if got := gotForm.Get("client_id"); got != "public-client-id" {
		t.Errorf("client_id = %q, want %q", got, "public-client-id")
	}
	if got := gotForm.Get("code_verifier"); got != flow.pkce.Verifier {
		t.Errorf("code_verifier = %q, want PKCE verifier %q", got, flow.pkce.Verifier)
	}
	if _, present := gotForm["client_secret"]; present {
		t.Errorf("client_secret must be omitted for public clients, got %q", gotForm.Get("client_secret"))
	}
	if gotAuthHeader != "" {
		t.Errorf("Authorization header must be absent for public clients, got %q", gotAuthHeader)
	}
}

func TestOAuthFlow_PublicClient_RefreshOmitsSecret(t *testing.T) {
	var gotForm url.Values
	var gotAuthHeader string
	mockServer := newTokenRecorder(t, &gotForm, &gotAuthHeader)
	defer mockServer.Close()

	oauth2Config := CreateOAuth2ConfigForInstance(mockServer.URL, "public-client-id", "")

	expired := &oauth2.Token{
		AccessToken:  "old-access-token",
		RefreshToken: "old-refresh-token",
		Expiry:       time.Now().Add(-time.Hour),
	}

	newToken, err := oauth2Config.TokenSource(context.Background(), expired).Token()
	if err != nil {
		t.Fatalf("refresh failed: %v", err)
	}
	if newToken.AccessToken != "public-access-token" {
		t.Errorf("access token = %q, want %q", newToken.AccessToken, "public-access-token")
	}

	if got := gotForm.Get("grant_type"); got != "refresh_token" {
		t.Errorf("grant_type = %q, want refresh_token", got)
	}
	if got := gotForm.Get("client_id"); got != "public-client-id" {
		t.Errorf("client_id = %q, want %q", got, "public-client-id")
	}
	if _, present := gotForm["client_secret"]; present {
		t.Errorf("client_secret must be omitted for public-client refresh, got %q", gotForm.Get("client_secret"))
	}
	if gotAuthHeader != "" {
		t.Errorf("Authorization header must be absent for public-client refresh, got %q", gotAuthHeader)
	}
}
