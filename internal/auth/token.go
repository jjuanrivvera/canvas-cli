package auth

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

const (
	keyringService = "canvas-cli"
	tokenFileName  = "token.enc"
)

// TokenStore defines the interface for token storage
type TokenStore interface {
	// Save saves a token to storage
	Save(instanceName string, token *oauth2.Token) error

	// Load loads a token from storage
	Load(instanceName string) (*oauth2.Token, error)

	// Delete deletes a token from storage
	Delete(instanceName string) error

	// Exists checks if a token exists for the instance
	Exists(instanceName string) bool
}

// KeyringTokenStore implements TokenStore using system keyring
type KeyringTokenStore struct {
	logger *slog.Logger
}

// NewKeyringTokenStore creates a new keyring token store
func NewKeyringTokenStore() *KeyringTokenStore {
	return &KeyringTokenStore{
		logger: slog.Default(),
	}
}

// Save saves a token to the keyring
func (s *KeyringTokenStore) Save(instanceName string, token *oauth2.Token) error {
	// Serialize token to JSON
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// Save to keyring
	err = keyring.Set(keyringService, instanceName, string(data))
	if err != nil {
		return fmt.Errorf("failed to save token to keyring: %w", err)
	}

	s.logger.Debug("Token saved to keyring", "instance", instanceName)
	return nil
}

// Load loads a token from the keyring
func (s *KeyringTokenStore) Load(instanceName string) (*oauth2.Token, error) {
	// Load from keyring
	data, err := keyring.Get(keyringService, instanceName)
	if err != nil {
		return nil, fmt.Errorf("failed to load token from keyring: %w", err)
	}

	// Deserialize token from JSON
	var token oauth2.Token
	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	s.logger.Debug("Token loaded from keyring", "instance", instanceName)
	return &token, nil
}

// Delete deletes a token from the keyring
func (s *KeyringTokenStore) Delete(instanceName string) error {
	err := keyring.Delete(keyringService, instanceName)
	if err != nil {
		return fmt.Errorf("failed to delete token from keyring: %w", err)
	}

	s.logger.Debug("Token deleted from keyring", "instance", instanceName)
	return nil
}

// Exists checks if a token exists in the keyring
func (s *KeyringTokenStore) Exists(instanceName string) bool {
	_, err := keyring.Get(keyringService, instanceName)
	return err == nil
}

// FileTokenStore implements TokenStore using encrypted files
type FileTokenStore struct {
	configDir string
	logger    *slog.Logger
}

// NewFileTokenStore creates a new file token store
func NewFileTokenStore(configDir string) *FileTokenStore {
	return &FileTokenStore{
		configDir: configDir,
		logger:    slog.Default(),
	}
}

// Save saves a token to an encrypted file
func (s *FileTokenStore) Save(instanceName string, token *oauth2.Token) error {
	// Ensure config directory exists
	tokenDir := filepath.Join(s.configDir, "tokens")
	if err := os.MkdirAll(tokenDir, 0700); err != nil {
		return fmt.Errorf("failed to create token directory: %w", err)
	}

	// Serialize token to JSON
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// Encrypt token
	encrypted, err := Encrypt(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt token: %w", err)
	}

	// Write to file
	tokenPath := filepath.Join(tokenDir, instanceName+"."+tokenFileName)
	if err := os.WriteFile(tokenPath, encrypted, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	s.logger.Debug("Token saved to file", "instance", instanceName, "path", tokenPath)
	return nil
}

// Load loads a token from an encrypted file
func (s *FileTokenStore) Load(instanceName string) (*oauth2.Token, error) {
	tokenPath := filepath.Join(s.configDir, "tokens", instanceName+"."+tokenFileName)

	// Read encrypted file
	encrypted, err := os.ReadFile(tokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}

	// Decrypt token
	data, err := Decrypt(encrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token: %w", err)
	}

	// Deserialize token from JSON
	var token oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	s.logger.Debug("Token loaded from file", "instance", instanceName, "path", tokenPath)
	return &token, nil
}

// Delete deletes a token file
func (s *FileTokenStore) Delete(instanceName string) error {
	tokenPath := filepath.Join(s.configDir, "tokens", instanceName+"."+tokenFileName)

	if err := os.Remove(tokenPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete token file: %w", err)
	}

	s.logger.Debug("Token deleted from file", "instance", instanceName, "path", tokenPath)
	return nil
}

// Exists checks if a token file exists
func (s *FileTokenStore) Exists(instanceName string) bool {
	tokenPath := filepath.Join(s.configDir, "tokens", instanceName+"."+tokenFileName)
	_, err := os.Stat(tokenPath)
	return err == nil
}

// FallbackTokenStore implements TokenStore with keyring primary and file fallback
type FallbackTokenStore struct {
	keyring *KeyringTokenStore
	file    *FileTokenStore
	logger  *slog.Logger
}

// NewFallbackTokenStore creates a new fallback token store
func NewFallbackTokenStore(configDir string) *FallbackTokenStore {
	return &FallbackTokenStore{
		keyring: NewKeyringTokenStore(),
		file:    NewFileTokenStore(configDir),
		logger:  slog.Default(),
	}
}

// Save saves a token, trying keyring first, then falling back to file
func (s *FallbackTokenStore) Save(instanceName string, token *oauth2.Token) error {
	// Try keyring first
	err := s.keyring.Save(instanceName, token)
	if err == nil {
		return nil
	}

	// Keyring failed, try file storage
	s.logger.Warn("Keyring save failed, using file storage", "error", err)
	return s.file.Save(instanceName, token)
}

// Load loads a token, trying keyring first, then falling back to file
func (s *FallbackTokenStore) Load(instanceName string) (*oauth2.Token, error) {
	// Try keyring first
	token, err := s.keyring.Load(instanceName)
	if err == nil {
		return token, nil
	}

	// Keyring failed, try file storage
	s.logger.Debug("Keyring load failed, trying file storage", "error", err)
	return s.file.Load(instanceName)
}

// Delete deletes a token from both storages
func (s *FallbackTokenStore) Delete(instanceName string) error {
	// Try to delete from both storages (ignore errors)
	keyringErr := s.keyring.Delete(instanceName)
	fileErr := s.file.Delete(instanceName)

	if keyringErr != nil && fileErr != nil {
		return fmt.Errorf("failed to delete from both storages: keyring=%v, file=%v", keyringErr, fileErr)
	}

	return nil
}

// Exists checks if a token exists in either storage
func (s *FallbackTokenStore) Exists(instanceName string) bool {
	return s.keyring.Exists(instanceName) || s.file.Exists(instanceName)
}
