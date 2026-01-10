package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jjuanrivvera/canvas-cli/internal/auth"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	authInstanceName string
	authOAuthMode    string
	authClientID     string
	authClientSecret string
)

// authCmd represents the auth command group
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication with Canvas",
	Long: `Manage authentication with Canvas LMS instances.

The auth command provides subcommands for logging in, logging out,
and checking authentication status.`,
}

// authLoginCmd represents the auth login command
var authLoginCmd = &cobra.Command{
	Use:   "login [instance-url]",
	Short: "Authenticate with a Canvas instance",
	Long: `Authenticate with a Canvas instance using OAuth 2.0 with PKCE.

The login command starts an OAuth flow to authenticate with Canvas.
By default, it will try to open a local callback server. If that fails,
it will fall back to out-of-band (manual copy-paste) mode.

OAuth credentials (client ID and client secret) are required and can be
provided via flags or entered interactively when prompted.

Examples:
  canvas auth login https://canvas.instructure.com
  canvas auth login https://myschool.instructure.com --instance myschool
  canvas auth login https://canvas.instructure.com --client-id YOUR_ID --client-secret YOUR_SECRET
  canvas auth login --mode oob  # Force out-of-band mode`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAuthLogin,
}

// authLogoutCmd represents the auth logout command
var authLogoutCmd = &cobra.Command{
	Use:   "logout [instance-name]",
	Short: "Logout from a Canvas instance",
	Long: `Logout from a Canvas instance by removing stored credentials.

If no instance name is provided, logs out from the default instance.

Examples:
  canvas auth logout
  canvas auth logout myschool`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAuthLogout,
}

// authStatusCmd represents the auth status command
var authStatusCmd = &cobra.Command{
	Use:   "status [instance-name]",
	Short: "Check authentication status",
	Long: `Check authentication status for Canvas instances.

Shows which instances are configured and authenticated.

Examples:
  canvas auth status
  canvas auth status myschool`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAuthStatus,
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authLogoutCmd)
	authCmd.AddCommand(authStatusCmd)

	// Login flags
	authLoginCmd.Flags().StringVar(&authInstanceName, "instance", "", "Instance name (defaults to hostname)")
	authLoginCmd.Flags().StringVar(&authOAuthMode, "mode", "auto", "OAuth mode: auto, local, oob")
	authLoginCmd.Flags().StringVar(&authClientID, "client-id", "", "OAuth client ID")
	authLoginCmd.Flags().StringVar(&authClientSecret, "client-secret", "", "OAuth client secret")
}

func runAuthLogin(cmd *cobra.Command, args []string) error {
	// Get instance URL
	var instanceURL string
	if len(args) > 0 {
		instanceURL = args[0]
	} else {
		return fmt.Errorf("instance URL is required")
	}

	// Normalize URL
	normalizedURL, err := config.NormalizeURL(instanceURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Determine instance name
	if authInstanceName == "" {
		// Use hostname as instance name
		authInstanceName = getHostnameFromURL(normalizedURL)
	}

	authInstanceName = config.SanitizeInstanceName(authInstanceName)

	fmt.Printf("🔐 Logging in to Canvas instance: %s\n", normalizedURL)
	fmt.Printf("Instance name: %s\n\n", authInstanceName)

	// Parse OAuth mode
	var oauthMode auth.OAuthMode
	switch authOAuthMode {
	case "auto":
		oauthMode = auth.OAuthModeAuto
	case "local":
		oauthMode = auth.OAuthModeLocal
	case "oob":
		oauthMode = auth.OAuthModeOOB
	default:
		return fmt.Errorf("invalid OAuth mode: %s (must be auto, local, or oob)", authOAuthMode)
	}

	// Get or prompt for client ID
	if authClientID == "" {
		fmt.Print("Enter OAuth Client ID: ")
		fmt.Scanln(&authClientID)
	}

	// If client ID is provided, also require client secret for OAuth
	if authClientID != "" && authClientSecret == "" {
		fmt.Print("Enter OAuth Client Secret: ")
		fmt.Scanln(&authClientSecret)
		if authClientSecret == "" {
			return fmt.Errorf("client secret is required when using OAuth with a client ID")
		}
	}

	// Create OAuth flow
	oauthFlow, err := auth.NewOAuthFlow(&auth.OAuthFlowConfig{
		BaseURL:      normalizedURL,
		ClientID:     authClientID,
		ClientSecret: authClientSecret,
		Mode:         oauthMode,
	})
	if err != nil {
		return fmt.Errorf("failed to create OAuth flow: %w", err)
	}

	// Perform authentication
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	token, err := oauthFlow.Authenticate(ctx)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Save token
	configDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	configDir = configDir + "/.canvas-cli"

	tokenStore := auth.NewFallbackTokenStore(configDir)
	if err := tokenStore.Save(authInstanceName, token); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Add or update instance
	instance := &config.Instance{
		Name:         authInstanceName,
		URL:          normalizedURL,
		ClientID:     authClientID,
		ClientSecret: authClientSecret,
	}

	if _, exists := cfg.Instances[authInstanceName]; exists {
		if err := cfg.UpdateInstance(authInstanceName, instance); err != nil {
			return fmt.Errorf("failed to update instance: %w", err)
		}
	} else {
		if err := cfg.AddInstance(instance); err != nil {
			return fmt.Errorf("failed to add instance: %w", err)
		}
	}

	fmt.Printf("\n✓ Successfully authenticated with %s\n", authInstanceName)
	fmt.Printf("Token expires: %s\n", token.Expiry.Format(time.RFC3339))

	return nil
}

func runAuthLogout(cmd *cobra.Command, args []string) error {
	// Determine instance name
	var instanceName string
	if len(args) > 0 {
		instanceName = args[0]
	} else {
		// Load config to get default instance
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if cfg.DefaultInstance == "" {
			return fmt.Errorf("no default instance configured")
		}

		instanceName = cfg.DefaultInstance
	}

	// Confirm logout
	fmt.Printf("Are you sure you want to logout from %s? (y/N): ", instanceName)
	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "y" && confirm != "Y" {
		fmt.Println("Logout cancelled")
		return nil
	}

	// Get config directory
	configDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	configDir = configDir + "/.canvas-cli"

	// Delete token
	tokenStore := auth.NewFallbackTokenStore(configDir)
	if err := tokenStore.Delete(instanceName); err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	fmt.Printf("✓ Successfully logged out from %s\n", instanceName)

	return nil
}

func runAuthStatus(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(cfg.Instances) == 0 {
		fmt.Println("No Canvas instances configured")
		fmt.Println("\nRun 'canvas auth login <instance-url>' to get started")
		return nil
	}

	// Get config directory
	configDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	configDir = configDir + "/.canvas-cli"

	tokenStore := auth.NewFallbackTokenStore(configDir)

	// Check specific instance or all instances
	if len(args) > 0 {
		instanceName := args[0]
		instance, err := cfg.GetInstance(instanceName)
		if err != nil {
			return err
		}

		printInstanceStatus(instance, cfg.DefaultInstance == instanceName, tokenStore)
	} else {
		// Show all instances
		fmt.Println("Canvas Instances:")
		fmt.Println()

		for _, instance := range cfg.ListInstances() {
			printInstanceStatus(instance, cfg.DefaultInstance == instance.Name, tokenStore)
			fmt.Println()
		}
	}

	return nil
}

func printInstanceStatus(instance *config.Instance, isDefault bool, tokenStore auth.TokenStore) {
	defaultMarker := ""
	if isDefault {
		defaultMarker = " (default)"
	}

	fmt.Printf("📌 %s%s\n", instance.Name, defaultMarker)
	fmt.Printf("   URL: %s\n", instance.URL)

	// Check authentication status
	if tokenStore.Exists(instance.Name) {
		token, err := tokenStore.Load(instance.Name)
		if err != nil {
			fmt.Printf("   Status: ❌ Error loading token\n")
			return
		}

		if token.Expiry.Before(time.Now()) {
			fmt.Printf("   Status: ⚠️  Token expired\n")
			fmt.Printf("   Expired: %s\n", token.Expiry.Format(time.RFC3339))
		} else {
			fmt.Printf("   Status: ✓ Authenticated\n")
			fmt.Printf("   Expires: %s\n", token.Expiry.Format(time.RFC3339))
		}
	} else {
		fmt.Printf("   Status: ❌ Not authenticated\n")
	}
}

func getHostnameFromURL(urlStr string) string {
	// Simple extraction - just get the hostname part
	// This is a basic implementation
	start := 0
	if idx := findIndex(urlStr, "://"); idx != -1 {
		start = idx + 3
	}

	end := len(urlStr)
	if idx := findIndexFrom(urlStr, "/", start); idx != -1 {
		end = idx
	}
	if idx := findIndexFrom(urlStr, ":", start); idx != -1 {
		if end > idx {
			end = idx
		}
	}

	hostname := urlStr[start:end]

	// Remove "www." prefix if present
	if len(hostname) > 4 && hostname[:4] == "www." {
		hostname = hostname[4:]
	}

	// Remove domain extension for cleaner name
	if idx := findIndex(hostname, "."); idx != -1 {
		hostname = hostname[:idx]
	}

	return hostname
}

func findIndex(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func findIndexFrom(s, substr string, start int) int {
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
