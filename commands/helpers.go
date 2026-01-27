package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"github.com/jjuanrivvera/canvas-cli/internal/auth"
	"github.com/jjuanrivvera/canvas-cli/internal/cache"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
	"github.com/jjuanrivvera/canvas-cli/internal/output"
)

// getUserAgent returns the User-Agent string with version
func getUserAgent() string {
	if version != "" && version != "dev" {
		return fmt.Sprintf("canvas-cli/%s", version)
	}
	return "canvas-cli/dev"
}

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
			UserAgent:      getUserAgent(),
			MaxResults:     globalLimit,
			DryRun:         dryRun,
			ShowToken:      showToken,
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

	// Create cache if not disabled
	var apiCache cache.CacheInterface
	cacheEnabled := !noCache
	if cacheEnabled {
		apiCache = createCache()
	}

	var clientConfig api.ClientConfig

	// Check if instance has an API token configured (token auth - no OAuth required)
	if instance.HasToken() {
		clientConfig = api.ClientConfig{
			BaseURL:        instance.URL,
			Token:          instance.Token,
			RequestsPerSec: cfg.Settings.RequestsPerSecond,
			AsUserID:       asUserID,
			Cache:          apiCache,
			CacheEnabled:   cacheEnabled,
			UserAgent:      getUserAgent(),
			MaxResults:     globalLimit,
			DryRun:         dryRun,
			ShowToken:      showToken,
		}

		if verbose {
			fmt.Fprintf(os.Stderr, "Using API token authentication for %s\n", instance.Name)
		}
	} else {
		// OAuth flow - load token from store
		configDir, err := config.GetConfigDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get config directory: %w", err)
		}

		tokenStore := auth.NewFallbackTokenStore(configDir)
		token, err := tokenStore.Load(instance.Name)
		if err != nil {
			return nil, fmt.Errorf("not authenticated with %s. Run 'canvas auth login' or 'canvas auth token set' first", instance.Name)
		}

		// Create auto-refreshing token source if we have OAuth credentials
		if instance.HasOAuth() {
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
				UserAgent:      getUserAgent(),
				MaxResults:     globalLimit,
				DryRun:         dryRun,
				ShowToken:      showToken,
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
				UserAgent:      getUserAgent(),
				MaxResults:     globalLimit,
				DryRun:         dryRun,
				ShowToken:      showToken,
			}
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
	configDir, err := config.GetConfigDir()
	if err != nil {
		return cache.New(5 * time.Minute)
	}

	cacheDir := filepath.Join(configDir, "cache")

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

// printInfo prints an informational message unless --quiet is set.
// Use this for success messages and status output that should not appear
// when the user is piping or scripting.
func printInfo(format string, args ...interface{}) {
	if !quiet {
		fmt.Printf(format, args...)
	}
}

// printInfoln prints an informational line unless --quiet is set.
func printInfoln(a ...interface{}) {
	if !quiet {
		fmt.Println(a...)
	}
}

// formatOutput formats and prints data according to the global outputFormat setting.
// If outputFormat is "table" (default), it uses the custom display function if provided.
// For other formats (json, yaml, csv), it uses the output formatter.
// In table format, output is compact by default (key fields only). Use -v/--verbose for all fields.
// Applies filtering (--filter, --columns, --sort) before output.
func formatOutput(data interface{}, customTableDisplay func()) error {
	format := output.FormatType(outputFormat)

	// Apply filtering if any filtering options are set
	if hasFilteringOptions() {
		data = applyFiltering(data)
	}

	// For table format, use custom display if provided (but not when filtering)
	if format == output.FormatTable {
		if customTableDisplay != nil && !hasFilteringOptions() {
			customTableDisplay()
			return nil
		}
	}

	// For structured formats, use the formatter with verbose option
	return output.WriteWithOptions(os.Stdout, data, format, verbose)
}

// formatSuccessOutput prints a success message (only in table format) and outputs the data.
// For JSON/YAML/CSV, the success message is omitted and only the raw data is output.
// This enables scripting with structured output formats.
// Applies filtering (--filter, --columns, --sort) before output.
func formatSuccessOutput(data interface{}, successMessage string) error {
	format := output.FormatType(outputFormat)

	// Apply filtering if any filtering options are set
	if hasFilteringOptions() {
		data = applyFiltering(data)
	}

	// Only print success message for table format and when not quiet
	if format == output.FormatTable && successMessage != "" && !quiet {
		fmt.Println(successMessage)
	}

	return output.WriteWithOptions(os.Stdout, data, format, verbose)
}

// formatEmptyOrOutput handles the case when a list might be empty.
// For JSON/YAML output, it always outputs valid structured data ([] for empty).
// For table output, it prints a user-friendly message when empty.
// Applies filtering (--filter, --columns, --sort) before output.
func formatEmptyOrOutput(data interface{}, emptyMessage string) error {
	format := output.FormatType(outputFormat)

	// Apply filtering if any filtering options are set
	if hasFilteringOptions() {
		data = applyFiltering(data)
	}

	// For structured formats (JSON, YAML, CSV), always output the data
	// even if empty - the formatter will output [] for empty slices
	if format == output.FormatJSON || format == output.FormatYAML || format == output.FormatCSV {
		return output.WriteWithOptions(os.Stdout, data, format, verbose)
	}

	// For table format, show user-friendly message when empty
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice && v.Len() == 0 {
		if !quiet {
			fmt.Println(emptyMessage)
		}
		return nil
	}

	return output.WriteWithOptions(os.Stdout, data, format, verbose)
}

// ExactArgsWithUsage returns a cobra.PositionalArgs that validates exact arg count
// with a descriptive error message showing what arguments are expected
func ExactArgsWithUsage(n int, argNames ...string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			if len(argNames) > 0 {
				expected := strings.Join(argNames, "> <")
				if len(args) < n {
					return fmt.Errorf("missing required argument(s)\n\nUsage:\n  %s <%s>\n\nRun '%s --help' for more information", cmd.CommandPath(), expected, cmd.CommandPath())
				}
				return fmt.Errorf("too many arguments provided\n\nUsage:\n  %s <%s>\n\nRun '%s --help' for more information", cmd.CommandPath(), expected, cmd.CommandPath())
			}
			// Fallback to default message
			return fmt.Errorf("accepts %d arg(s), received %d", n, len(args))
		}
		return nil
	}
}

// MinArgsWithUsage returns a cobra.PositionalArgs that validates minimum arg count
// with a descriptive error message
func MinArgsWithUsage(n int, argNames ...string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < n {
			if len(argNames) > 0 {
				expected := strings.Join(argNames, "> <")
				return fmt.Errorf("missing required argument(s)\n\nUsage:\n  %s <%s>\n\nRun '%s --help' for more information", cmd.CommandPath(), expected, cmd.CommandPath())
			}
			return fmt.Errorf("requires at least %d arg(s), received %d", n, len(args))
		}
		return nil
	}
}

// getDefaultAccountID returns the default account ID for the current instance.
// Returns 0 if no default account ID is configured.
func getDefaultAccountID() (int64, error) {
	cfg, err := config.Load()
	if err != nil {
		return 0, fmt.Errorf("failed to load config: %w", err)
	}

	instance, err := cfg.GetDefaultInstance()
	if err != nil {
		return 0, fmt.Errorf("failed to get default instance: %w", err)
	}

	return instance.DefaultAccountID, nil
}

// resolveAccountID returns the provided account ID if non-zero, otherwise returns the default account ID.
// Returns an error with helpful guidance if no account ID is available.
func resolveAccountID(providedID int64, context string) (int64, error) {
	if providedID != 0 {
		return providedID, nil
	}

	defaultID, err := getDefaultAccountID()
	if err != nil || defaultID == 0 {
		return 0, fmt.Errorf("--account-id is required (no default account configured). Use 'canvas config account --detect' to set one")
	}

	printVerbose("Using default account ID: %d for %s\n", defaultID, context)
	return defaultID, nil
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
