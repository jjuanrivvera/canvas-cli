package commands

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"github.com/jjuanrivvera/canvas-cli/internal/auth"
	"github.com/jjuanrivvera/canvas-cli/internal/cache"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
	"github.com/jjuanrivvera/canvas-cli/internal/output"
)

// getAPIClient creates an API client for the default or specified instance
func getAPIClient() (*api.Client, error) {
	// Check for environment variable authentication (CI/CD support)
	envURL := os.Getenv("CANVAS_URL")
	envToken := os.Getenv("CANVAS_TOKEN")

	if envURL != "" && envToken != "" {
		// Use environment variables for authentication
		requestsPerSec := 5.0 // Default
		if envRPS := os.Getenv("CANVAS_REQUESTS_PER_SEC"); envRPS != "" {
			fmt.Sscanf(envRPS, "%f", &requestsPerSec)
		}

		// Create cache if not disabled
		var apiCache cache.CacheInterface
		cacheEnabled := !noCache
		if cacheEnabled {
			apiCache = createCache()
		}

		client, err := api.NewClient(api.ClientConfig{
			BaseURL:        envURL,
			Token:          envToken,
			RequestsPerSec: requestsPerSec,
			AsUserID:       asUserID,
			Cache:          apiCache,
			CacheEnabled:   cacheEnabled,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create API client from environment: %w", err)
		}

		if verbose {
			fmt.Fprintln(os.Stderr, "Using Canvas credentials from environment variables")
			if cacheEnabled {
				fmt.Fprintln(os.Stderr, "Response caching enabled")
			}
		}

		return client, nil
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Get instance
	var instance *config.Instance
	if instanceURL != "" {
		// Find instance by name or URL
		for _, inst := range cfg.Instances {
			if inst.Name == instanceURL || inst.URL == instanceURL {
				instance = inst
				break
			}
		}
		if instance == nil {
			return nil, fmt.Errorf("no instance found with name or URL: %s. Use 'canvas auth list' to see configured instances", instanceURL)
		}
	} else {
		// Use default instance
		instance, err = cfg.GetDefaultInstance()
		if err != nil {
			return nil, fmt.Errorf("failed to get default instance: %w", err)
		}
	}

	// Get config directory
	configDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	configDir = configDir + "/.canvas-cli"

	// Load token
	tokenStore := auth.NewFallbackTokenStore(configDir)
	token, err := tokenStore.Load(instance.Name)
	if err != nil {
		return nil, fmt.Errorf("not authenticated with %s. Run 'canvas auth login' first", instance.Name)
	}

	// Create cache if not disabled
	var apiCache cache.CacheInterface
	cacheEnabled := !noCache
	if cacheEnabled {
		apiCache = createCache()
	}

	// Create auto-refreshing token source if we have OAuth credentials
	var clientConfig api.ClientConfig
	if instance.ClientID != "" && instance.ClientSecret != "" {
		// Create oauth2 config for token refresh
		oauth2Config := auth.CreateOAuth2ConfigForInstance(instance.URL, instance.ClientID, instance.ClientSecret)
		tokenSource := auth.NewAutoRefreshTokenSource(oauth2Config, tokenStore, instance.Name, token)

		clientConfig = api.ClientConfig{
			BaseURL:        instance.URL,
			TokenSource:    tokenSource,
			RequestsPerSec: cfg.Settings.RequestsPerSecond,
			AsUserID:       asUserID,
			Cache:          apiCache,
			CacheEnabled:   cacheEnabled,
		}
	} else {
		// Fall back to static token (no auto-refresh)
		clientConfig = api.ClientConfig{
			BaseURL:        instance.URL,
			Token:          token.AccessToken,
			RequestsPerSec: cfg.Settings.RequestsPerSecond,
			AsUserID:       asUserID,
			Cache:          apiCache,
			CacheEnabled:   cacheEnabled,
		}
	}

	// Create API client
	client, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	// Show masquerading warning if active
	if asUserID > 0 && verbose {
		fmt.Fprintf(os.Stderr, "WARNING: Masquerading as user %d. All actions will be recorded in the audit log.\n", asUserID)
	}

	if verbose && cacheEnabled {
		fmt.Fprintln(os.Stderr, "Response caching enabled")
	}

	return client, nil
}

// createCache creates a multi-tier cache for API responses
func createCache() cache.CacheInterface {
	// Get cache directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fall back to memory-only cache
		return cache.New(5 * time.Minute)
	}

	cacheDir := homeDir + "/.canvas-cli/cache"

	// Create multi-tier cache (memory + disk) with 5 minute TTL
	multiCache, err := cache.NewMultiTierCache(
		5*time.Minute, // Memory TTL
		cacheDir,
		5*time.Minute, // Disk TTL
	)
	if err != nil {
		// Fall back to memory-only cache if disk cache fails
		return cache.New(5 * time.Minute)
	}

	return multiCache
}

// getConfig loads the configuration
func getConfig() (*config.Config, error) {
	return config.Load()
}

// printVerbose prints a message only in verbose mode
func printVerbose(format string, args ...interface{}) {
	if verbose {
		fmt.Printf(format, args...)
	}
}

// formatOutput formats and prints data according to the global outputFormat setting.
// If outputFormat is "table" (default), it uses the custom display function if provided.
// For other formats (json, yaml, csv), it uses the output formatter.
// In table format, output is compact by default (key fields only). Use -v/--verbose for all fields.
func formatOutput(data interface{}, customTableDisplay func()) error {
	format := output.FormatType(outputFormat)

	// For table format, use custom display if provided
	if format == output.FormatTable {
		if customTableDisplay != nil {
			customTableDisplay()
			return nil
		}
	}

	// For structured formats, use the formatter with verbose option
	return output.WriteWithOptions(os.Stdout, data, format, verbose)
}

// validateCourseID checks if a course ID exists and returns a user-friendly error if not.
// Returns the course object if found, nil otherwise (for optional use of course data).
func validateCourseID(client *api.Client, courseID int64) (*api.Course, error) {
	if courseID <= 0 {
		return nil, fmt.Errorf("invalid course ID: %d", courseID)
	}

	coursesService := api.NewCoursesService(client)
	ctx := context.Background()

	course, err := coursesService.Get(ctx, courseID, nil)
	if err != nil {
		// Check for common error patterns and provide helpful messages
		errStr := err.Error()
		if strings.Contains(errStr, "404") || strings.Contains(errStr, "not found") {
			return nil, fmt.Errorf("course with ID %d not found. Use 'canvas courses list' to see available courses", courseID)
		}
		if strings.Contains(errStr, "401") || strings.Contains(errStr, "unauthorized") {
			return nil, fmt.Errorf("you are not authorized to access course %d. Check your permissions or authentication", courseID)
		}
		if strings.Contains(errStr, "403") || strings.Contains(errStr, "forbidden") {
			return nil, fmt.Errorf("access to course %d is forbidden. You may not have the required permissions", courseID)
		}
		return nil, fmt.Errorf("failed to verify course %d: %w", courseID, err)
	}

	return course, nil
}
