package commands

import (
	"fmt"
	"os"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"github.com/jjuanrivvera/canvas-cli/internal/auth"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
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

		client, err := api.NewClient(api.ClientConfig{
			BaseURL:        envURL,
			Token:          envToken,
			RequestsPerSec: requestsPerSec,
			AsUserID:       asUserID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create API client from environment: %w", err)
		}

		if verbose {
			fmt.Fprintln(os.Stderr, "Using Canvas credentials from environment variables")
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
		// Find instance by URL
		for _, inst := range cfg.Instances {
			if inst.URL == instanceURL {
				instance = inst
				break
			}
		}
		if instance == nil {
			return nil, fmt.Errorf("no instance found with URL: %s", instanceURL)
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

	// Create API client
	client, err := api.NewClient(api.ClientConfig{
		BaseURL:        instance.URL,
		Token:          token.AccessToken,
		RequestsPerSec: cfg.Settings.RequestsPerSecond,
		AsUserID:       asUserID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	// Show masquerading warning if active
	if asUserID > 0 && verbose {
		fmt.Fprintf(os.Stderr, "WARNING: Masquerading as user %d. All actions will be recorded in the audit log.\n", asUserID)
	}

	return client, nil
}

// getConfig loads the configuration
func getConfig() (*config.Config, error) {
	return config.Load()
}
