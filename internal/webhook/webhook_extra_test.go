package webhook

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ─── JWT helpers ──────────────────────────────────────────────────────────────

// signJWT creates a compact JWT signed with the given RSA private key and kid.
func signJWT(t *testing.T, claims jwt.MapClaims, key *rsa.PrivateKey, kid string) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid
	s, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("signJWT: %v", err)
	}
	return s
}

// newJWKServerFor creates an httptest.Server that serves a JWK set containing
// the public key for the given private key / kid pair.
func newJWKServerFor(t *testing.T, kid string, priv *rsa.PrivateKey) *httptest.Server {
	t.Helper()
	pub := &priv.PublicKey
	entry := jwk{
		Kty: "RSA",
		Kid: kid,
		Use: "sig",
		Alg: "RS256",
		N:   base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes()),
	}
	resp := jwkResponse{Keys: []jwk{entry}}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
}

// ─── verifyJWT ───────────────────────────────────────────────────────────────

func TestListener_verifyJWT_NoJWKSet(t *testing.T) {
	l := New(&Config{Addr: ":0"})
	// jwkSet is nil
	_, err := l.verifyJWT([]byte("anything"))
	if err == nil {
		t.Error("expected error when JWK set not configured")
	}
}

func TestListener_verifyJWT_ValidToken(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	kid := "key-1"
	srv := newJWKServerFor(t, kid, priv)
	defer srv.Close()

	l := New(&Config{Addr: ":0", JWKSetURL: srv.URL})

	claims := jwt.MapClaims{
		"id":         "evt-001",
		"event_type": "assignment_created",
		"event_time": time.Now().Format(time.RFC3339),
		"body": map[string]interface{}{
			"assignment_id": float64(42),
		},
		"exp": float64(time.Now().Add(time.Hour).Unix()),
	}
	tokenStr := signJWT(t, claims, priv, kid)

	event, err := l.verifyJWT([]byte(tokenStr))
	if err != nil {
		t.Fatalf("verifyJWT: %v", err)
	}
	if event.EventType != "assignment_created" {
		t.Errorf("EventType = %q, want 'assignment_created'", event.EventType)
	}
	if event.ID != "evt-001" {
		t.Errorf("ID = %q, want 'evt-001'", event.ID)
	}
}

func TestListener_verifyJWT_WithMetadataEventName(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	kid := "key-meta"
	srv := newJWKServerFor(t, kid, priv)
	defer srv.Close()

	l := New(&Config{Addr: ":0", JWKSetURL: srv.URL})

	claims := jwt.MapClaims{
		"metadata": map[string]interface{}{
			"event_name": "grade_change",
		},
		"data": map[string]interface{}{
			"grade": "A",
		},
		"exp": float64(time.Now().Add(time.Hour).Unix()),
	}
	tokenStr := signJWT(t, claims, priv, kid)

	event, err := l.verifyJWT([]byte(tokenStr))
	if err != nil {
		t.Fatalf("verifyJWT: %v", err)
	}
	if event.EventType != "grade_change" {
		t.Errorf("EventType = %q, want 'grade_change'", event.EventType)
	}
}

func TestListener_verifyJWT_AllClaimsAsBody(t *testing.T) {
	// When neither "body" nor "data" key is present, all non-standard claims
	// become the body.
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	kid := "key-body"
	srv := newJWKServerFor(t, kid, priv)
	defer srv.Close()

	l := New(&Config{Addr: ":0", JWKSetURL: srv.URL})

	claims := jwt.MapClaims{
		"event_type":   "user_created",
		"custom_field": "custom_value",
		"exp":          float64(time.Now().Add(time.Hour).Unix()),
	}
	tokenStr := signJWT(t, claims, priv, kid)

	event, err := l.verifyJWT([]byte(tokenStr))
	if err != nil {
		t.Fatalf("verifyJWT: %v", err)
	}
	if event.Body["custom_field"] != "custom_value" {
		t.Errorf("expected custom_field in body, got %v", event.Body)
	}
}

func TestListener_verifyJWT_InvalidToken(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	kid := "key-invalid"
	srv := newJWKServerFor(t, kid, priv)
	defer srv.Close()

	l := New(&Config{Addr: ":0", JWKSetURL: srv.URL})

	// Corrupt the token signature
	tokenStr := signJWT(t, jwt.MapClaims{"exp": float64(time.Now().Add(time.Hour).Unix())}, priv, kid)
	parts := strings.Split(tokenStr, ".")
	parts[2] = "invalidsignature"
	badToken := strings.Join(parts, ".")

	_, err = l.verifyJWT([]byte(badToken))
	if err == nil {
		t.Error("expected error for token with bad signature")
	}
}

func TestListener_verifyJWT_WrongAlg(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	kid := "key-alg"
	srv := newJWKServerFor(t, kid, priv)
	defer srv.Close()

	l := New(&Config{Addr: ":0", JWKSetURL: srv.URL})

	// Sign with HMAC instead of RSA — key function should reject it
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": float64(time.Now().Add(time.Hour).Unix()),
	})
	token.Header["kid"] = kid
	tokenStr, _ := token.SignedString([]byte("hmac-secret"))

	_, err = l.verifyJWT([]byte(tokenStr))
	if err == nil {
		t.Error("expected error for wrong signing algorithm")
	}
}

func TestListener_verifyJWT_MissingKID(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	// Server that serves a valid JWK set but the token won't have a kid
	srv := newJWKServerFor(t, "some-key", priv)
	defer srv.Close()

	l := New(&Config{Addr: ":0", JWKSetURL: srv.URL})

	// Create token without kid in header
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"exp": float64(time.Now().Add(time.Hour).Unix()),
	})
	// Do NOT set token.Header["kid"]
	tokenStr, _ := token.SignedString(priv)

	_, err = l.verifyJWT([]byte(tokenStr))
	if err == nil {
		t.Error("expected error when kid is missing from JWT header")
	}
}

// ─── parseJWTClaims ───────────────────────────────────────────────────────────

func TestListener_parseJWTClaims_Valid(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	claims := jwt.MapClaims{
		"id":         "parse-001",
		"event_type": "submission_created",
		"event_time": time.Now().Format(time.RFC3339),
		"body":       map[string]interface{}{"student": "alice"},
		"exp":        float64(time.Now().Add(time.Hour).Unix()),
	}
	tokenStr := signJWT(t, claims, priv, "kid-x")

	l := New(&Config{Addr: ":0"})
	event, err := l.parseJWTClaims([]byte(tokenStr))
	if err != nil {
		t.Fatalf("parseJWTClaims: %v", err)
	}
	if event.EventType != "submission_created" {
		t.Errorf("EventType = %q, want 'submission_created'", event.EventType)
	}
	if event.ID != "parse-001" {
		t.Errorf("ID = %q, want 'parse-001'", event.ID)
	}
	if event.Body["student"] != "alice" {
		t.Errorf("expected body.student='alice', got %v", event.Body)
	}
}

func TestListener_parseJWTClaims_WithMetadata(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	claims := jwt.MapClaims{
		"metadata": map[string]interface{}{
			"event_name": "enrollment_created",
		},
		"data": map[string]interface{}{"user_id": float64(7)},
		"exp":  float64(time.Now().Add(time.Hour).Unix()),
	}
	tokenStr := signJWT(t, claims, priv, "k")

	l := New(&Config{Addr: ":0"})
	event, err := l.parseJWTClaims([]byte(tokenStr))
	if err != nil {
		t.Fatalf("parseJWTClaims: %v", err)
	}
	if event.EventType != "enrollment_created" {
		t.Errorf("EventType = %q, want 'enrollment_created'", event.EventType)
	}
	if event.Body["user_id"] != float64(7) {
		t.Errorf("expected body.user_id=7, got %v", event.Body)
	}
}

func TestListener_parseJWTClaims_AllClaimsBody(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	claims := jwt.MapClaims{
		"event_type": "course_created",
		"extra":      "value",
		"exp":        float64(time.Now().Add(time.Hour).Unix()),
	}
	tokenStr := signJWT(t, claims, priv, "k2")

	l := New(&Config{Addr: ":0"})
	event, err := l.parseJWTClaims([]byte(tokenStr))
	if err != nil {
		t.Fatalf("parseJWTClaims: %v", err)
	}
	if event.Body["extra"] != "value" {
		t.Errorf("expected body.extra='value', got %v", event.Body)
	}
	// Standard JWT claims must be excluded from body
	for _, std := range []string{"iss", "sub", "aud", "exp", "nbf", "iat", "jti"} {
		if _, ok := event.Body[std]; ok {
			t.Errorf("standard claim %q should not appear in event.Body", std)
		}
	}
}

func TestListener_parseJWTClaims_InvalidToken(t *testing.T) {
	l := New(&Config{Addr: ":0"})
	_, err := l.parseJWTClaims([]byte("not.a.jwt"))
	if err == nil {
		t.Error("expected error for malformed JWT")
	}
}

// ─── handleWebhook: JWT body paths ───────────────────────────────────────────

func TestListener_handleWebhook_QuotedJWTBody(t *testing.T) {
	// Body arrives as a quoted JWT string (Canvas Data Services wraps JWT in quotes).
	// The listener has no JWKSet and no HMAC secret, so it uses the quote-strip +
	// parseJWTClaims path (no verification, just claim extraction).
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)

	l := New(&Config{Addr: ":0"}) // no JWKSetURL, no Secret

	var handledEvent *Event
	l.On("assignment_created", func(ctx context.Context, event *Event) error {
		handledEvent = event
		return nil
	})

	claims := jwt.MapClaims{
		"event_type": "assignment_created",
		"id":         "q-001",
		"exp":        float64(time.Now().Add(time.Hour).Unix()),
	}
	tokenStr := signJWT(t, claims, priv, "k")
	// Wrap in quotes as Canvas Data Services does
	quotedBody := fmt.Sprintf("%q", tokenStr)

	req := httptest.NewRequest("POST", "/webhook", bytes.NewBufferString(quotedBody))
	w := httptest.NewRecorder()
	l.handleWebhook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
	_ = handledEvent
}

func TestListener_handleWebhook_RawJWTBody(t *testing.T) {
	// Body is an unverified JWT (3 dot-segments, not JSON): parseJWTClaims path
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	kid := "raw-jwt-key"
	// No JWKSetURL → listener has no jwkSet, no HMAC secret
	l := New(&Config{Addr: ":0"})

	l.On("course_created", func(ctx context.Context, event *Event) error {
		return nil
	})

	claims := jwt.MapClaims{
		"event_type": "course_created",
		"id":         "r-001",
		"exp":        float64(time.Now().Add(time.Hour).Unix()),
	}
	tokenStr := signJWT(t, claims, priv, kid)

	req := httptest.NewRequest("POST", "/webhook", bytes.NewBufferString(tokenStr))
	w := httptest.NewRecorder()
	l.handleWebhook(w, req)

	// Should succeed (no verification required) with 200
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for unverified JWT body, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestListener_handleWebhook_JWTVerified_WithJWKSet(t *testing.T) {
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	kid := "verified-key"
	srv := newJWKServerFor(t, kid, priv)
	defer srv.Close()

	l := New(&Config{Addr: ":0", JWKSetURL: srv.URL})

	var received *Event
	l.On("user_created", func(ctx context.Context, event *Event) error {
		received = event
		return nil
	})

	claims := jwt.MapClaims{
		"event_type": "user_created",
		"id":         "v-001",
		"exp":        float64(time.Now().Add(time.Hour).Unix()),
	}
	tokenStr := signJWT(t, claims, priv, kid)

	req := httptest.NewRequest("POST", "/webhook", bytes.NewBufferString(tokenStr))
	w := httptest.NewRecorder()
	l.handleWebhook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
	_ = received
}

func TestListener_handleWebhook_LongRawBody(t *testing.T) {
	// Exercises the >200-char truncation logging path.
	l := New(&Config{Addr: ":0"})

	// Build a JSON body with an event_type field whose total JSON exceeds 200 chars
	longName := strings.Repeat("x", 200)
	body := fmt.Sprintf(`{"event_type":"assignment_created","long_field":"%s"}`, longName)

	req := httptest.NewRequest("POST", "/webhook", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	l.handleWebhook(w, req)

	// No handler registered for this event type → 200
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestListener_handleWebhook_HandlerError(t *testing.T) {
	l := New(&Config{Addr: ":0"})
	l.On("assignment_created", func(ctx context.Context, _ *Event) error {
		return fmt.Errorf("handler exploded")
	})

	body := `{"event_type":"assignment_created"}`
	req := httptest.NewRequest("POST", "/webhook", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	l.handleWebhook(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 when handler returns error, got %d", w.Code)
	}
}

// ─── Shutdown on nil server --------------------------------------------------

func TestListener_Shutdown_NilServer(t *testing.T) {
	l := New(&Config{Addr: ":0"})
	// server is nil (never started)
	err := l.Shutdown(context.TODO())
	if err != nil {
		t.Errorf("expected no error when shutting down un-started listener, got %v", err)
	}
}

// ─── New with JWKSetURL ------------------------------------------------------

func TestNew_WithJWKSetURL(t *testing.T) {
	l := New(&Config{
		Addr:      ":0",
		JWKSetURL: "https://example.com/jwks",
	})
	if l.jwkSet == nil {
		t.Error("expected jwkSet to be initialized when JWKSetURL is provided")
	}
	if l.jwkSet.url != "https://example.com/jwks" {
		t.Errorf("jwkSet.url = %q, want 'https://example.com/jwks'", l.jwkSet.url)
	}
}

// ─── JWKSet: GetKey with stale cache returns cached key on refresh failure ──

func TestJWKSet_GetKey_StaleCache_RefreshFails_UsesCachedKey(t *testing.T) {
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	kid := "stale-key"

	// Server initially works
	alive := true
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !alive {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		entry := jwk{
			Kty: "RSA", Kid: kid, Use: "sig", Alg: "RS256",
			N: base64.RawURLEncoding.EncodeToString(priv.PublicKey.N.Bytes()),
			E: base64.RawURLEncoding.EncodeToString(big.NewInt(int64(priv.PublicKey.E)).Bytes()),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jwkResponse{Keys: []jwk{entry}})
	}))
	defer srv.Close()

	set := NewJWKSet(srv.URL)
	set.ttl = 1 * time.Nanosecond // expire immediately after first fetch

	// Prime the cache
	key1, err := set.GetKey(kid)
	if err != nil {
		t.Fatalf("initial GetKey: %v", err)
	}
	if key1 == nil {
		t.Fatal("expected non-nil key")
	}

	// Kill the server
	alive = false

	// Cache is stale; refresh will fail; should return the previously cached key
	key2, err := set.GetKey(kid)
	if err != nil {
		t.Fatalf("GetKey after stale cache + refresh failure: %v", err)
	}
	if key2.N.Cmp(priv.PublicKey.N) != 0 {
		t.Error("returned key does not match original")
	}
}

// ─── JWKSet: Refresh with non-RSA key skipped ─────────────────────────────

func TestJWKSet_Refresh_NonRSAKeySkipped(t *testing.T) {
	// Return one non-RSA key (EC) and one valid RSA key.
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := jwkResponse{
			Keys: []jwk{
				// EC key (non-RSA) — should be skipped
				{Kty: "EC", Kid: "ec-key", Use: "sig"},
				// RSA key — should be accepted
				{
					Kty: "RSA", Kid: "rsa-key", Use: "sig", Alg: "RS256",
					N: base64.RawURLEncoding.EncodeToString(priv.PublicKey.N.Bytes()),
					E: base64.RawURLEncoding.EncodeToString(big.NewInt(int64(priv.PublicKey.E)).Bytes()),
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	set := NewJWKSet(srv.URL)
	if err := set.Refresh(); err != nil {
		t.Fatalf("Refresh: %v", err)
	}
	if set.KeyCount() != 1 {
		t.Errorf("expected 1 RSA key, got %d", set.KeyCount())
	}
}
