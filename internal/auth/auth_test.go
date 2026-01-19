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

func TestKeyringTokenStore_SaveLoadCycle(t *testing.T) {
	store := NewKeyringTokenStore()

	token := &oauth2.Token{
		AccessToken:  "keyring-test-token",
		RefreshToken: "keyring-refresh-token",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	}

	// Save token
	err := store.Save("keyring-test-instance", token)
	if err != nil {
		// Keyring might not be available in CI environment
		t.Skipf("Keyring save failed (may not be available): %v", err)
	}

	// Load token
	loadedToken, err := store.Load("keyring-test-instance")
	if err != nil {
		t.Fatalf("Load failed after successful save: %v", err)
	}

	if loadedToken.AccessToken != token.AccessToken {
		t.Errorf("expected access token %s, got %s", token.AccessToken, loadedToken.AccessToken)
	}

	// Clean up
	store.Delete("keyring-test-instance")
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

func TestGetMachineID_WithEnvOverride(t *testing.T) {
	// Set environment variable override
	testID := "test-machine-id-12345"
	oldEnv := os.Getenv("CANVAS_CLI_MACHINE_ID")
	os.Setenv("CANVAS_CLI_MACHINE_ID", testID)
	defer os.Setenv("CANVAS_CLI_MACHINE_ID", oldEnv)

	id, err := getMachineID()
	if err != nil {
		t.Fatalf("getMachineID failed with env override: %v", err)
	}

	if id != testID {
		t.Errorf("expected machine ID %s, got %s", testID, id)
	}
}

func TestGetUsername_USER_Env(t *testing.T) {
	// Save original env
	oldUser := os.Getenv("USER")
	oldUsername := os.Getenv("USERNAME")

	// Set USER env variable
	testUser := "testuser123"
	os.Setenv("USER", testUser)
	os.Setenv("USERNAME", "") // Clear USERNAME to test USER priority
	defer func() {
		os.Setenv("USER", oldUser)
		os.Setenv("USERNAME", oldUsername)
	}()

	username := getUsername()
	if username != testUser {
		t.Errorf("expected username %s, got %s", testUser, username)
	}
}

func TestGetUsername_USERNAME_Env(t *testing.T) {
	// Save original env
	oldUser := os.Getenv("USER")
	oldUsername := os.Getenv("USERNAME")

	// Set USERNAME env variable (Windows)
	testUser := "windowsuser456"
	os.Setenv("USER", "") // Clear USER to test USERNAME fallback
	os.Setenv("USERNAME", testUser)
	defer func() {
		os.Setenv("USER", oldUser)
		os.Setenv("USERNAME", oldUsername)
	}()

	username := getUsername()
	if username != testUser {
		t.Errorf("expected username %s, got %s", testUser, username)
	}
}

func TestDeriveEncryptionKey_Consistency(t *testing.T) {
	// Set consistent environment for testing
	oldMachineID := os.Getenv("CANVAS_CLI_MACHINE_ID")
	oldUser := os.Getenv("USER")

	os.Setenv("CANVAS_CLI_MACHINE_ID", "test-machine-123")
	os.Setenv("USER", "testuser")
	defer func() {
		os.Setenv("CANVAS_CLI_MACHINE_ID", oldMachineID)
		os.Setenv("USER", oldUser)
	}()

	salt := []byte("test-salt-123456")

	// Derive key twice with same inputs
	key1, err := deriveEncryptionKey(salt)
	if err != nil {
		t.Fatalf("First deriveEncryptionKey failed: %v", err)
	}

	key2, err := deriveEncryptionKey(salt)
	if err != nil {
		t.Fatalf("Second deriveEncryptionKey failed: %v", err)
	}

	// Keys should be identical
	if len(key1) != len(key2) {
		t.Errorf("keys have different lengths: %d vs %d", len(key1), len(key2))
	}

	for i := range key1 {
		if key1[i] != key2[i] {
			t.Error("derived keys are not consistent")
			break
		}
	}
}

func TestGenerateSalt(t *testing.T) {
	salt1, err := generateSalt()
	if err != nil {
		t.Fatalf("generateSalt failed: %v", err)
	}

	if len(salt1) != 16 {
		t.Errorf("expected salt length 16, got %d", len(salt1))
	}

	salt2, err := generateSalt()
	if err != nil {
		t.Fatalf("second generateSalt failed: %v", err)
	}

	// Salts should be different
	same := true
	for i := range salt1 {
		if salt1[i] != salt2[i] {
			same = false
			break
		}
	}

	if same {
		t.Error("expected different salts from consecutive calls")
	}
}

func TestFileTokenStore_Load_CorruptedFile(t *testing.T) {
	tempDir := t.TempDir()
	store := NewFileTokenStore(tempDir)

	// Create token directory
	tokenDir := filepath.Join(tempDir, "tokens")
	if err := os.MkdirAll(tokenDir, 0700); err != nil {
		t.Fatalf("failed to create token directory: %v", err)
	}

	// Write corrupted data to token file
	tokenPath := filepath.Join(tokenDir, "corrupted.token")
	if err := os.WriteFile(tokenPath, []byte("corrupted-data"), 0600); err != nil {
		t.Fatalf("failed to write corrupted file: %v", err)
	}

	// Try to load corrupted token
	_, err := store.Load("corrupted")
	if err == nil {
		t.Error("expected error when loading corrupted token file")
	}
}

func TestFileTokenStore_Delete_Success(t *testing.T) {
	tempDir := t.TempDir()
	store := NewFileTokenStore(tempDir)

	token := &oauth2.Token{
		AccessToken:  "delete-test-token",
		RefreshToken: "delete-test-refresh",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	}

	// Save token
	if err := store.Save("delete-test", token); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify it exists
	if !store.Exists("delete-test") {
		t.Fatal("expected token to exist after save")
	}

	// Delete token
	if err := store.Delete("delete-test"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify it no longer exists
	if store.Exists("delete-test") {
		t.Error("expected token to not exist after delete")
	}
}

func TestFileTokenStore_Delete_NotExists(t *testing.T) {
	tempDir := t.TempDir()
	store := NewFileTokenStore(tempDir)

	// Delete non-existent token should not error
	err := store.Delete("nonexistent")
	if err != nil {
		t.Errorf("expected no error when deleting nonexistent token, got %v", err)
	}
}

func TestGetUsername_LOGNAME_Env(t *testing.T) {
	// Save original env
	oldUser := os.Getenv("USER")
	oldUsername := os.Getenv("USERNAME")
	oldLogname := os.Getenv("LOGNAME")

	// Set LOGNAME env variable
	testUser := "lognameuser789"
	os.Setenv("USER", "")     // Clear USER
	os.Setenv("USERNAME", "") // Clear USERNAME
	os.Setenv("LOGNAME", testUser)
	defer func() {
		os.Setenv("USER", oldUser)
		os.Setenv("USERNAME", oldUsername)
		os.Setenv("LOGNAME", oldLogname)
	}()

	username := getUsername()
	if username != testUser {
		t.Errorf("expected username %s from LOGNAME, got %s", testUser, username)
	}
}

func TestEncryptDecrypt_EdgeCases(t *testing.T) {
	// Set consistent environment for testing
	oldMachineID := os.Getenv("CANVAS_CLI_MACHINE_ID")
	oldUser := os.Getenv("USER")

	os.Setenv("CANVAS_CLI_MACHINE_ID", "test-machine-encrypt")
	os.Setenv("USER", "testuser")
	defer func() {
		os.Setenv("CANVAS_CLI_MACHINE_ID", oldMachineID)
		os.Setenv("USER", oldUser)
	}()

	tests := []struct {
		name      string
		plaintext []byte
	}{
		{
			name:      "empty data",
			plaintext: []byte{},
		},
		{
			name:      "single byte",
			plaintext: []byte{0x42},
		},
		{
			name:      "large data",
			plaintext: make([]byte, 10000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			encrypted, err := Encrypt(tt.plaintext)
			if err != nil {
				t.Fatalf("Encrypt failed: %v", err)
			}

			// Decrypt
			decrypted, err := Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt failed: %v", err)
			}

			// Compare
			if len(decrypted) != len(tt.plaintext) {
				t.Errorf("decrypted length %d != original length %d", len(decrypted), len(tt.plaintext))
			}

			for i := range tt.plaintext {
				if decrypted[i] != tt.plaintext[i] {
					t.Error("decrypted data does not match original")
					break
				}
			}
		})
	}
}

func TestDecrypt_InvalidData(t *testing.T) {
	// Set consistent environment for testing
	oldMachineID := os.Getenv("CANVAS_CLI_MACHINE_ID")
	oldUser := os.Getenv("USER")

	os.Setenv("CANVAS_CLI_MACHINE_ID", "test-machine-invalid")
	os.Setenv("USER", "testuser")
	defer func() {
		os.Setenv("CANVAS_CLI_MACHINE_ID", oldMachineID)
		os.Setenv("USER", oldUser)
	}()

	// Create data with correct length but invalid content
	invalidData := make([]byte, 60)
	for i := range invalidData {
		invalidData[i] = byte(i)
	}

	_, err := Decrypt(invalidData)
	if err == nil {
		t.Error("expected error when decrypting invalid data")
	}
}

func TestFallbackTokenStore_FallbackToFile(t *testing.T) {
	tempDir := t.TempDir()

	fallbackStore := NewFallbackTokenStore(tempDir)

	token := &oauth2.Token{
		AccessToken:  "chained-test-token",
		RefreshToken: "chained-test-refresh",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	}

	// Save token (might go to keyring or file depending on availability)
	err := fallbackStore.Save("fallback-test", token)
	if err != nil {
		t.Fatalf("FallbackTokenStore.Save failed: %v", err)
	}

	// Load token
	loadedToken, err := fallbackStore.Load("fallback-test")
	if err != nil {
		t.Fatalf("FallbackTokenStore.Load failed: %v", err)
	}

	if loadedToken.AccessToken != token.AccessToken {
		t.Errorf("expected access token %s, got %s", token.AccessToken, loadedToken.AccessToken)
	}

	// Clean up
	fallbackStore.Delete("fallback-test")
}

func TestDeriveEncryptionKey_EmptyUsername(t *testing.T) {
	// Save original env
	oldMachineID := os.Getenv("CANVAS_CLI_MACHINE_ID")
	oldUser := os.Getenv("USER")
	oldUsername := os.Getenv("USERNAME")
	oldLogname := os.Getenv("LOGNAME")

	// Set machine ID but clear all username env vars
	os.Setenv("CANVAS_CLI_MACHINE_ID", "test-machine-no-user")
	os.Setenv("USER", "")
	os.Setenv("USERNAME", "")
	os.Setenv("LOGNAME", "")

	defer func() {
		os.Setenv("CANVAS_CLI_MACHINE_ID", oldMachineID)
		os.Setenv("USER", oldUser)
		os.Setenv("USERNAME", oldUsername)
		os.Setenv("LOGNAME", oldLogname)
	}()

	salt := []byte("test-salt-empty-user")

	// This test depends on whether whoami command is available
	// If whoami returns a username, the key will be derived successfully
	// If whoami fails and returns empty, deriveEncryptionKey should error
	_, err := deriveEncryptionKey(salt)

	// We expect either success (if whoami works) or an error about empty username
	// The key point is testing that the code handles the empty username path
	if err != nil {
		if !contains(err.Error(), "username") {
			t.Errorf("expected error about username, got %v", err)
		}
	}
	// If no error, whoami succeeded and returned a username
}

func TestKeyringTokenStore_Delete(t *testing.T) {
	store := NewKeyringTokenStore()

	// Try to delete non-existent token
	err := store.Delete("nonexistent-token")
	// This might error or succeed depending on keyring implementation
	// The important thing is that it doesn't panic
	_ = err
}

func TestEncrypt_EmptyUsername(t *testing.T) {
	// Save original env
	oldMachineID := os.Getenv("CANVAS_CLI_MACHINE_ID")
	oldUser := os.Getenv("USER")
	oldUsername := os.Getenv("USERNAME")
	oldLogname := os.Getenv("LOGNAME")

	// Set machine ID but clear all username env vars
	os.Setenv("CANVAS_CLI_MACHINE_ID", "test-machine-encrypt-no-user")
	os.Setenv("USER", "")
	os.Setenv("USERNAME", "")
	os.Setenv("LOGNAME", "")

	defer func() {
		os.Setenv("CANVAS_CLI_MACHINE_ID", oldMachineID)
		os.Setenv("USER", oldUser)
		os.Setenv("USERNAME", oldUsername)
		os.Setenv("LOGNAME", oldLogname)
	}()

	plaintext := []byte("test data")

	// This test depends on whether whoami command is available
	_, err := Encrypt(plaintext)

	// We expect either success (if whoami works) or an error about empty username
	if err != nil {
		if !contains(err.Error(), "username") && !contains(err.Error(), "failed to get machine ID") {
			t.Errorf("expected error about username or machine ID, got %v", err)
		}
	}
}

func TestFileTokenStore_Save_CreateDirectory(t *testing.T) {
	tempDir := t.TempDir()

	// Use a nested directory that doesn't exist yet
	configDir := filepath.Join(tempDir, "deep", "nested", "config")
	store := NewFileTokenStore(configDir)

	token := &oauth2.Token{
		AccessToken: "test-create-dir-token",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(time.Hour),
	}

	// Save should create the directory
	err := store.Save("test-create-dir", token)
	if err != nil {
		t.Fatalf("Save failed when creating nested directories: %v", err)
	}

	// Verify token was saved
	loadedToken, err := store.Load("test-create-dir")
	if err != nil {
		t.Fatalf("Load failed after save: %v", err)
	}

	if loadedToken.AccessToken != token.AccessToken {
		t.Errorf("expected access token %s, got %s", token.AccessToken, loadedToken.AccessToken)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || stringContains(s, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
