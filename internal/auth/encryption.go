package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// PBKDF2 iterations - OWASP recommends 600,000 for SHA-256 (2023)
	// Using 100,000 as a balance between security and CLI responsiveness
	pbkdf2Iterations = 100000
	// Salt size in bytes (16 bytes = 128 bits)
	saltSize = 16
	// Key size for AES-256
	keySize = 32
)

// deriveEncryptionKey derives an encryption key from machine ID and username using PBKDF2
// The salt must be provided for decryption; for encryption, generate a new salt
func deriveEncryptionKey(salt []byte) ([]byte, error) {
	machineID, err := getMachineID()
	if err != nil {
		return nil, fmt.Errorf("failed to get machine ID: %w", err)
	}

	username := getUsername()
	if username == "" {
		return nil, fmt.Errorf("failed to get username")
	}

	// Combine machine ID and username as the password material
	combined := machineID + ":" + username

	// Use PBKDF2 with SHA-256 for proper key derivation
	// This provides key stretching and protection against rainbow table attacks
	key := pbkdf2.Key([]byte(combined), salt, pbkdf2Iterations, keySize, sha256.New)
	return key, nil
}

// generateSalt generates a cryptographically secure random salt
func generateSalt() ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

// getMachineID gets a unique machine identifier
// Returns an error if no secure machine ID can be obtained
// Note: Does NOT fall back to hostname as that is too easily guessable
func getMachineID() (string, error) {
	// Check for environment variable override (useful for testing/CI)
	if envID := os.Getenv("CANVAS_CLI_MACHINE_ID"); envID != "" {
		return envID, nil
	}

	// Try Linux machine-id files first
	if fileExists("/etc/machine-id") {
		data, err := os.ReadFile("/etc/machine-id")
		if err == nil {
			id := strings.TrimSpace(string(data))
			if id != "" {
				return id, nil
			}
		}
	}

	if fileExists("/var/lib/dbus/machine-id") {
		data, err := os.ReadFile("/var/lib/dbus/machine-id")
		if err == nil {
			id := strings.TrimSpace(string(data))
			if id != "" {
				return id, nil
			}
		}
	}

	// macOS - use IOPlatformUUID
	cmd := exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "IOPlatformUUID") {
				parts := strings.Split(line, "\"")
				if len(parts) >= 4 {
					uuid := parts[3]
					if uuid != "" {
						return uuid, nil
					}
				}
			}
		}
	}

	// Windows - try PowerShell (more reliable on modern Windows and CI)
	cmd = exec.Command("powershell", "-Command",
		"(Get-CimInstance -ClassName Win32_ComputerSystemProduct).UUID")
	output, err = cmd.Output()
	if err == nil {
		id := strings.TrimSpace(string(output))
		if id != "" && id != "FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF" {
			return id, nil
		}
	}

	// Windows fallback - try wmic (deprecated but may still work)
	cmd = exec.Command("wmic", "csproduct", "get", "UUID")
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			// Skip header line and empty lines
			if line != "" && line != "UUID" && line != "FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF" {
				return line, nil
			}
		}
	}

	// Windows fallback - try registry via PowerShell
	cmd = exec.Command("powershell", "-Command",
		"(Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Cryptography' -Name 'MachineGuid').MachineGuid")
	output, err = cmd.Output()
	if err == nil {
		id := strings.TrimSpace(string(output))
		if id != "" {
			return id, nil
		}
	}

	// Fail securely - do not fall back to guessable identifiers like hostname
	// This ensures encryption keys are based on truly unique machine identifiers
	return "", fmt.Errorf("could not obtain unique machine identifier: " +
		"ensure /etc/machine-id exists (Linux), ioreg works (macOS), or wmic/PowerShell works (Windows)")
}

// getUsername gets the current username
func getUsername() string {
	// Try USER environment variable first
	username := os.Getenv("USER")
	if username != "" {
		return username
	}

	// Try USERNAME for Windows
	username = os.Getenv("USERNAME")
	if username != "" {
		return username
	}

	// Try LOGNAME
	username = os.Getenv("LOGNAME")
	if username != "" {
		return username
	}

	// Fallback: try to get from whoami command
	cmd := exec.Command("whoami")
	output, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output))
	}

	return "unknown"
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Encrypt encrypts data using AES-GCM with a derived key
// Output format: [salt (16 bytes)][nonce (12 bytes)][ciphertext+tag]
func Encrypt(plaintext []byte) ([]byte, error) {
	// Generate a fresh salt for this encryption
	salt, err := generateSalt()
	if err != nil {
		return nil, err
	}

	// Derive key using PBKDF2 with the new salt
	key, err := deriveEncryptionKey(salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt: output is [nonce][ciphertext+tag]
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// Prepend salt to the final output: [salt][nonce][ciphertext+tag]
	result := make([]byte, saltSize+len(ciphertext))
	copy(result[:saltSize], salt)
	copy(result[saltSize:], ciphertext)

	return result, nil
}

// Decrypt decrypts data using AES-GCM with a derived key
// Input format: [salt (16 bytes)][nonce (12 bytes)][ciphertext+tag]
func Decrypt(data []byte) ([]byte, error) {
	// Validate minimum length: salt + nonce + 16-byte auth tag (empty plaintext is valid)
	// salt=16, nonce=12, tag=16, total minimum = 44 bytes
	if len(data) < saltSize+12+16 {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Extract salt from the beginning
	salt := data[:saltSize]
	encryptedData := data[saltSize:]

	// Derive key using PBKDF2 with the extracted salt
	key, err := deriveEncryptionKey(salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}
