package config

import (
	"fmt"
	"net/url"
	"strings"
)

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate default instance
	if c.DefaultInstance != "" {
		if _, exists := c.Instances[c.DefaultInstance]; !exists {
			return fmt.Errorf("default instance %q does not exist", c.DefaultInstance)
		}
	}

	// Validate each instance
	for name, instance := range c.Instances {
		if err := ValidateInstance(instance); err != nil {
			return fmt.Errorf("instance %q is invalid: %w", name, err)
		}
	}

	// Validate settings
	if c.Settings != nil {
		if err := ValidateSettings(c.Settings); err != nil {
			return fmt.Errorf("settings are invalid: %w", err)
		}
	}

	return nil
}

// ValidateInstance validates an instance configuration
func ValidateInstance(instance *Instance) error {
	if instance == nil {
		return fmt.Errorf("instance cannot be nil")
	}

	// Validate name
	if instance.Name == "" {
		return fmt.Errorf("instance name is required")
	}

	if len(instance.Name) > 100 {
		return fmt.Errorf("instance name is too long (max 100 characters)")
	}

	// Validate name doesn't contain special characters
	if strings.ContainsAny(instance.Name, "/\\:*?\"<>|") {
		return fmt.Errorf("instance name contains invalid characters")
	}

	// Validate URL
	if instance.URL == "" {
		return fmt.Errorf("instance URL is required")
	}

	parsedURL, err := url.Parse(instance.URL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL must use http or https scheme")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a host")
	}

	// Validate authentication if configured
	if instance.Token != "" {
		if err := ValidateToken(instance.Token); err != nil {
			return fmt.Errorf("invalid token: %w", err)
		}
	}

	// Validate OAuth credentials if configured
	if instance.ClientID != "" && instance.ClientSecret == "" {
		return fmt.Errorf("client_secret is required when client_id is set")
	}
	if instance.ClientSecret != "" && instance.ClientID == "" {
		return fmt.Errorf("client_id is required when client_secret is set")
	}

	if instance.ClientID != "" && len(instance.ClientID) < 10 {
		return fmt.Errorf("client_id seems too short (minimum 10 characters)")
	}

	if instance.ClientSecret != "" && len(instance.ClientSecret) < 10 {
		return fmt.Errorf("client_secret seems too short (minimum 10 characters)")
	}

	// Validate default account ID
	if instance.DefaultAccountID < 0 {
		return fmt.Errorf("default_account_id cannot be negative")
	}

	return nil
}

// ValidateToken validates an API token format
func ValidateToken(token string) error {
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	// Trim whitespace for validation
	trimmed := strings.TrimSpace(token)
	if trimmed == "" {
		return fmt.Errorf("token cannot be only whitespace")
	}

	// Canvas tokens are typically long alphanumeric strings
	// Minimum reasonable length is 20 characters
	if len(trimmed) < 20 {
		return fmt.Errorf("token seems too short (minimum 20 characters), got %d", len(trimmed))
	}

	// Maximum reasonable length (Canvas tokens are usually < 100 chars)
	if len(trimmed) > 500 {
		return fmt.Errorf("token seems too long (maximum 500 characters), got %d", len(trimmed))
	}

	// Check for common placeholder values
	lowerToken := strings.ToLower(trimmed)
	placeholders := []string{"your-token-here", "your_token_here", "replace-me", "changeme", "example", "token"}
	for _, placeholder := range placeholders {
		if lowerToken == placeholder {
			return fmt.Errorf("token appears to be a placeholder value: %q", placeholder)
		}
	}

	return nil
}

// ValidateSettings validates the settings configuration
func ValidateSettings(settings *Settings) error {
	if settings == nil {
		return fmt.Errorf("settings cannot be nil")
	}

	// Validate output format
	validFormats := map[string]bool{
		"table": true,
		"json":  true,
		"yaml":  true,
		"csv":   true,
	}

	if !validFormats[settings.DefaultOutputFormat] {
		return fmt.Errorf("invalid output format: %q (must be one of: table, json, yaml, csv)", settings.DefaultOutputFormat)
	}

	// Validate requests per second
	if settings.RequestsPerSecond <= 0 {
		return fmt.Errorf("requests_per_second must be positive")
	}

	if settings.RequestsPerSecond > 100 {
		return fmt.Errorf("requests_per_second is too high (max 100)")
	}

	// Validate cache TTL
	if settings.CacheEnabled && settings.CacheTTL < 0 {
		return fmt.Errorf("cache_ttl_minutes cannot be negative")
	}

	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[strings.ToLower(settings.LogLevel)] {
		return fmt.Errorf("invalid log level: %q (must be one of: debug, info, warn, error)", settings.LogLevel)
	}

	return nil
}

// SanitizeInstanceName sanitizes an instance name
func SanitizeInstanceName(name string) string {
	// Remove invalid characters
	name = strings.Map(func(r rune) rune {
		if strings.ContainsRune("/\\:*?\"<>|", r) {
			return '-'
		}
		return r
	}, name)

	// Trim whitespace
	name = strings.TrimSpace(name)

	// Limit length
	if len(name) > 100 {
		name = name[:100]
	}

	return name
}

// NormalizeURL normalizes a Canvas instance URL
func NormalizeURL(rawURL string) (string, error) {
	// Ensure URL has a scheme
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	// Parse URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	// Remove trailing slash
	parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/")

	// Remove any path if it's just /
	if parsedURL.Path == "/" {
		parsedURL.Path = ""
	}

	return parsedURL.String(), nil
}
