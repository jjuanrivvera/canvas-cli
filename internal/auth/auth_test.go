package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestFileTokenStore_Save_Load(t *testing.T) {
	tempDir := t.TempDir()

	store := NewFileTokenStore(tempDir)

	token := &oauth2.Token{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	}

	// Test Save
	err := store.Save("test-instance", token)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Test Load
	loadedToken, err := store.Load("test-instance")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loadedToken.AccessToken != token.AccessToken {
		t.Errorf("expected access token '%s', got '%s'", token.AccessToken, loadedToken.AccessToken)
	}

	if loadedToken.RefreshToken != token.RefreshToken {
		t.Errorf("expected refresh token '%s', got '%s'", token.RefreshToken, loadedToken.RefreshToken)
	}
}

func TestFileTokenStore_Load_NotFound(t *testing.T) {
	tempDir := t.TempDir()
	store := NewFileTokenStore(tempDir)

	_, err := store.Load("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent token")
	}
}

func TestFileTokenStore_Delete(t *testing.T) {
	tempDir := t.TempDir()
	store := NewFileTokenStore(tempDir)

	token := &oauth2.Token{
		AccessToken: "test-token",
	}

	// Save token
	err := store.Save("test-instance", token)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify it exists
	if !store.Exists("test-instance") {
		t.Error("expected token to exist after save")
	}

	// Delete token
	err = store.Delete("test-instance")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify it doesn't exist
	if store.Exists("test-instance") {
		t.Error("expected token to not exist after delete")
	}
}

func TestFileTokenStore_Exists(t *testing.T) {
	tempDir := t.TempDir()
	store := NewFileTokenStore(tempDir)

	// Test non-existent
	if store.Exists("nonexistent") {
		t.Error("expected Exists to return false for nonexistent token")
	}

	// Save token
	token := &oauth2.Token{AccessToken: "test"}
	store.Save("test-instance", token)

	// Test exists
	if !store.Exists("test-instance") {
		t.Error("expected Exists to return true for saved token")
	}
}

func TestFallbackTokenStore_PreferKeyring(t *testing.T) {
	tempDir := t.TempDir()
	store := NewFallbackTokenStore(tempDir)

	// FallbackTokenStore should be created successfully
	if store == nil {
		t.Fatal("expected non-nil fallback store")
	}
}

func TestFallbackTokenStore_Save_Load(t *testing.T) {
	tempDir := t.TempDir()
	store := NewFallbackTokenStore(tempDir)

	token := &oauth2.Token{
		AccessToken:  "fallback-test-token",
		RefreshToken: "fallback-refresh",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	}

	// Save should succeed (falls back to file if keyring fails)
	err := store.Save("test-fallback", token)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load should succeed
	loadedToken, err := store.Load("test-fallback")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loadedToken.AccessToken != token.AccessToken {
		t.Errorf("expected access token '%s', got '%s'", token.AccessToken, loadedToken.AccessToken)
	}
}

func TestFallbackTokenStore_Delete(t *testing.T) {
	tempDir := t.TempDir()
	store := NewFallbackTokenStore(tempDir)

	token := &oauth2.Token{AccessToken: "test-delete"}

	// Save
	err := store.Save("test-delete", token)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Delete
	err = store.Delete("test-delete")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deleted
	if store.Exists("test-delete") {
		t.Error("expected token to not exist after delete")
	}
}

func TestFallbackTokenStore_Exists(t *testing.T) {
	tempDir := t.TempDir()
	store := NewFallbackTokenStore(tempDir)

	// Non-existent
	if store.Exists("nonexistent") {
		t.Error("expected false for nonexistent token")
	}

	// Save and check
	token := &oauth2.Token{AccessToken: "exists-test"}
	store.Save("exists-test", token)

	if !store.Exists("exists-test") {
		t.Error("expected true for existing token")
	}
}

func TestEncryption_Encrypt_Decrypt(t *testing.T) {
	plaintext := []byte("sensitive data to encrypt")

	// Test encryption
	ciphertext, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if len(ciphertext) == 0 {
		t.Error("expected non-empty ciphertext")
	}

	// Test decryption
	decrypted, err := Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("expected '%s', got '%s'", plaintext, decrypted)
	}
}

func TestEncryption_DecryptInvalidData(t *testing.T) {
	invalidData := []byte("not encrypted data")

	_, err := Decrypt(invalidData)
	if err == nil {
		t.Error("expected error when decrypting invalid data")
	}
}

func TestFileTokenStore_TokenFilePath(t *testing.T) {
	tempDir := t.TempDir()
	store := NewFileTokenStore(tempDir)

	token := &oauth2.Token{AccessToken: "path-test"}
	err := store.Save("test-instance", token)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists in expected location
	expectedPath := filepath.Join(tempDir, "tokens", "test-instance.token.enc")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("expected token file at %s", expectedPath)
	}
}

func TestGeneratePKCEChallenge(t *testing.T) {
	pkce, err := GeneratePKCEChallenge()
	if err != nil {
		t.Fatalf("GeneratePKCEChallenge failed: %v", err)
	}

	if len(pkce.Verifier) == 0 {
		t.Error("expected non-empty verifier")
	}

	if len(pkce.Challenge) == 0 {
		t.Error("expected non-empty challenge")
	}

	// Verifier should be different from challenge
	if pkce.Verifier == pkce.Challenge {
		t.Error("verifier and challenge should be different")
	}

	// Challenge should be base64 encoded
	if len(pkce.Challenge) < 40 {
		t.Error("challenge seems too short")
	}

	// Method should be S256
	if pkce.Method != "S256" {
		t.Errorf("expected method 'S256', got '%s'", pkce.Method)
	}
}

func TestGeneratePKCEChallenge_Uniqueness(t *testing.T) {
	pkce1, err := GeneratePKCEChallenge()
	if err != nil {
		t.Fatalf("First generation failed: %v", err)
	}

	pkce2, err := GeneratePKCEChallenge()
	if err != nil {
		t.Fatalf("Second generation failed: %v", err)
	}

	// Each generation should produce unique values
	if pkce1.Verifier == pkce2.Verifier {
		t.Error("expected different verifiers")
	}

	if pkce1.Challenge == pkce2.Challenge {
		t.Error("expected different challenges")
	}
}

func TestNewOAuthFlow(t *testing.T) {
	config := &OAuthFlowConfig{
		BaseURL:  "https://canvas.example.com",
		ClientID: "test-client-id",
		Mode:     OAuthModeLocal,
	}

	flow, err := NewOAuthFlow(config)
	if err != nil {
		t.Fatalf("NewOAuthFlow failed: %v", err)
	}

	if flow == nil {
		t.Fatal("expected non-nil flow")
	}

	if flow.config.ClientID != "test-client-id" {
		t.Errorf("expected client ID 'test-client-id', got '%s'", flow.config.ClientID)
	}
}

func TestNewOAuthFlow_MissingBaseURL(t *testing.T) {
	config := &OAuthFlowConfig{
		ClientID: "test-client-id",
		Mode:     OAuthModeLocal,
	}

	_, err := NewOAuthFlow(config)
	if err == nil {
		t.Error("expected error when base URL is missing")
	}
}

func TestNewOAuthFlow_MissingClientID(t *testing.T) {
	config := &OAuthFlowConfig{
		BaseURL: "https://canvas.example.com",
		Mode:    OAuthModeLocal,
	}

	_, err := NewOAuthFlow(config)
	if err == nil {
		t.Error("expected error when client ID is missing")
	}
}

func TestToken_IsExpired(t *testing.T) {
	// Expired token
	expiredToken := &oauth2.Token{
		AccessToken: "test",
		Expiry:      time.Now().Add(-time.Hour),
	}

	if expiredToken.Valid() {
		t.Error("expected expired token to be invalid")
	}

	// Valid token
	validToken := &oauth2.Token{
		AccessToken: "test",
		Expiry:      time.Now().Add(time.Hour),
	}

	if !validToken.Valid() {
		t.Error("expected valid token to be valid")
	}
}

func TestFileTokenStore_MultipleInstances(t *testing.T) {
	tempDir := t.TempDir()
	store := NewFileTokenStore(tempDir)

	// Save tokens for multiple instances
	token1 := &oauth2.Token{AccessToken: "instance1-token"}
	token2 := &oauth2.Token{AccessToken: "instance2-token"}
	token3 := &oauth2.Token{AccessToken: "instance3-token"}

	store.Save("instance1", token1)
	store.Save("instance2", token2)
	store.Save("instance3", token3)

	// Verify all exist
	if !store.Exists("instance1") || !store.Exists("instance2") || !store.Exists("instance3") {
		t.Error("expected all instances to exist")
	}

	// Load and verify each
	loaded1, _ := store.Load("instance1")
	loaded2, _ := store.Load("instance2")
	loaded3, _ := store.Load("instance3")

	if loaded1.AccessToken != "instance1-token" {
		t.Error("instance1 token mismatch")
	}
	if loaded2.AccessToken != "instance2-token" {
		t.Error("instance2 token mismatch")
	}
	if loaded3.AccessToken != "instance3-token" {
		t.Error("instance3 token mismatch")
	}
}

func TestOAuthMode_String(t *testing.T) {
	tests := []struct {
		mode     OAuthMode
		expected string
	}{
		{OAuthModeAuto, "auto"},
		{OAuthModeLocal, "local"},
		{OAuthModeOOB, "oob"},
		{OAuthMode(999), "unknown"},
	}

	for _, tt := range tests {
		result := tt.mode.String()
		if result != tt.expected {
			t.Errorf("expected '%s', got '%s'", tt.expected, result)
		}
	}
}

func TestKeyringTokenStore_Creation(t *testing.T) {
	store := NewKeyringTokenStore()
	if store == nil {
		t.Fatal("expected non-nil keyring store")
	}
}

func TestFileTokenStore_Creation(t *testing.T) {
	tempDir := t.TempDir()
	store := NewFileTokenStore(tempDir)
	if store == nil {
		t.Fatal("expected non-nil file store")
	}
	if store.configDir != tempDir {
		t.Errorf("expected configDir '%s', got '%s'", tempDir, store.configDir)
	}
}

func TestFallbackTokenStore_Creation(t *testing.T) {
	tempDir := t.TempDir()
	store := NewFallbackTokenStore(tempDir)
	if store == nil {
		t.Fatal("expected non-nil fallback store")
	}
	if store.keyring == nil {
		t.Error("expected keyring to be initialized")
	}
	if store.file == nil {
		t.Error("expected file store to be initialized")
	}
}

func TestEncryption_EmptyData(t *testing.T) {
	encrypted, err := Encrypt([]byte{})
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if len(decrypted) != 0 {
		t.Errorf("expected empty data, got %d bytes", len(decrypted))
	}
}

func TestEncryption_LargeData(t *testing.T) {
	// Test with 1MB of data
	largeData := make([]byte, 1024*1024)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	encrypted, err := Encrypt(largeData)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if len(decrypted) != len(largeData) {
		t.Errorf("expected %d bytes, got %d bytes", len(largeData), len(decrypted))
	}

	// Check first and last bytes
	if decrypted[0] != largeData[0] || decrypted[len(largeData)-1] != largeData[len(largeData)-1] {
		t.Error("decrypted data doesn't match original")
	}
}

func TestDecrypt_TooShort(t *testing.T) {
	_, err := Decrypt([]byte("short"))
	if err == nil {
		t.Error("expected error for data too short")
	}
}

func TestDecrypt_CorruptedNonce(t *testing.T) {
	data := []byte("test data")
	encrypted, err := Encrypt(data)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Corrupt the nonce
	if len(encrypted) > 12 {
		encrypted[5] ^= 0xFF
	}

	_, err = Decrypt(encrypted)
	if err == nil {
		t.Error("expected error for corrupted nonce")
	}
}

func TestGeneratePKCEChallenge_ErrorCases(t *testing.T) {
	challenge, err := GeneratePKCEChallenge()
	if err != nil {
		t.Fatalf("GeneratePKCEChallenge failed: %v", err)
	}

	if challenge.Verifier == "" {
		t.Error("expected non-empty verifier")
	}

	if challenge.Challenge == "" {
		t.Error("expected non-empty challenge")
	}

	if challenge.Method != "S256" {
		t.Errorf("expected method S256, got %s", challenge.Method)
	}

	// Verify challenge is base64url encoded
	if len(challenge.Challenge) < 43 {
		t.Error("challenge seems too short for base64url encoded SHA256")
	}
}

func TestKeyringTokenStore_ErrorHandling(t *testing.T) {
	store := NewKeyringTokenStore()

	// Test Load with non-existent token
	_, err := store.Load("nonexistent-instance")
	if err == nil {
		t.Error("expected error for non-existent token")
	}

	// Test Delete with non-existent token (should not error)
	err = store.Delete("nonexistent-instance")
	// Delete is allowed to succeed even if item doesn't exist
	if err != nil {
		// Some keyring implementations may return error, that's OK
		t.Logf("Delete returned error (acceptable): %v", err)
	}

	// Test Exists with non-existent token
	exists := store.Exists("nonexistent-instance")
	if exists {
		t.Error("expected false for non-existent token")
	}
}

func TestFileTokenStore_ErrorHandling(t *testing.T) {
	tempDir := t.TempDir()
	store := NewFileTokenStore(tempDir)

	// Test Load with non-existent file
	_, err := store.Load("nonexistent-instance-test")
	if err == nil {
		t.Error("expected error for non-existent file")
	}

	// Test Exists with non-existent file
	exists := store.Exists("nonexistent-instance-test")
	if exists {
		t.Error("expected false for non-existent file")
	}
}

func TestFallbackTokenStore_ErrorHandling(t *testing.T) {
	tempDir := t.TempDir()
	store := NewFallbackTokenStore(tempDir)

	// Test Load with non-existent token
	_, err := store.Load("nonexistent-fallback-test")
	if err == nil {
		t.Error("expected error for non-existent token")
	}

	// Test Exists with non-existent token
	exists := store.Exists("nonexistent-fallback-test")
	if exists {
		t.Error("expected false for non-existent token")
	}
}
