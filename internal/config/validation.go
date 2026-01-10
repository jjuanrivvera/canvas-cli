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
