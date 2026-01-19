package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/auth"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

// authCmd represents the auth command group
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication with Canvas",
	Long: `Manage authentication with Canvas LMS instances.

The auth command provides subcommands for logging in, logging out,
and checking authentication status.`,
}

// authTokenCmd represents the auth token command group
var authTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage API token authentication",
	Long: `Manage API token authentication for Canvas instances.

API tokens provide a simpler alternative to OAuth for authentication.
You can generate an API token in Canvas under Account > Settings > New Access Token.

Use 'canvas auth token set' to configure token authentication for an instance.`,
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(newAuthLoginCmd())
	authCmd.AddCommand(newAuthLogoutCmd())
	authCmd.AddCommand(newAuthStatusCmd())
	authCmd.AddCommand(authTokenCmd)

	// Token subcommands
	authTokenCmd.AddCommand(newAuthTokenSetCmd())
	authTokenCmd.AddCommand(newAuthTokenRemoveCmd())
}

func newAuthLoginCmd() *cobra.Command {
	opts := &options.AuthLoginOptions{
		OAuthMode: "auto", // default value
	}

	cmd := &cobra.Command{
		Use:   "login [instance-url]",
		Short: "Authenticate with a Canvas instance",
		Long: `Authenticate with a Canvas instance using OAuth 2.0 with PKCE.

The login command starts an OAuth flow to authenticate with Canvas.
By default, it will try to open a local callback server. If that fails,
it will fall back to out-of-band (manual copy-paste) mode.

If --instance is provided and the instance exists in your config, the URL
and OAuth credentials will be loaded automatically. You can override them
with flags.

Examples:
  # Login using a configured instance (recommended)
  canvas auth login --instance prod

  # Login with URL (creates new instance)
  canvas auth login https://canvas.instructure.com

  # Login with custom OAuth credentials
  canvas auth login --instance prod --client-id YOUR_ID --client-secret YOUR_SECRET

  # Force out-of-band mode (for headless systems)
  canvas auth login --instance prod --mode oob`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.InstanceURL = args[0]
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			return runAuthLogin(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.InstanceName, "instance", "", "Instance name (defaults to hostname)")
	cmd.Flags().StringVar(&opts.OAuthMode, "mode", "auto", "OAuth mode: auto, local, oob")
	cmd.Flags().StringVar(&opts.ClientID, "client-id", "", "OAuth client ID")
	cmd.Flags().StringVar(&opts.ClientSecret, "client-secret", "", "OAuth client secret")

	return cmd
}

func newAuthLogoutCmd() *cobra.Command {
	opts := &options.AuthLogoutOptions{}

	cmd := &cobra.Command{
		Use:   "logout [instance-name]",
		Short: "Logout from a Canvas instance",
		Long: `Logout from a Canvas instance by removing stored credentials.

If no instance name is provided, logs out from the default instance.

Examples:
  canvas auth logout
  canvas auth logout myschool`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.InstanceName = args[0]
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			return runAuthLogout(cmd.Context(), opts)
		},
	}

	return cmd
}

func newAuthStatusCmd() *cobra.Command {
	opts := &options.AuthStatusOptions{}

	cmd := &cobra.Command{
		Use:   "status [instance-name]",
		Short: "Check authentication status",
		Long: `Check authentication status for Canvas instances.

Shows which instances are configured and authenticated.

Examples:
  canvas auth status
  canvas auth status myschool`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.InstanceName = args[0]
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			return runAuthStatus(cmd.Context(), opts)
		},
	}

	return cmd
}

func newAuthTokenSetCmd() *cobra.Command {
	opts := &options.AuthTokenSetOptions{}

	cmd := &cobra.Command{
		Use:   "set <instance-name>",
		Short: "Set API token for an instance",
		Long: `Set an API access token for a Canvas instance.

This is an alternative to OAuth authentication. Generate a token in Canvas
under Account > Settings > New Access Token.

If the instance doesn't exist, it will be created (requires --url).
If the instance exists, the token will be updated.

Examples:
  # Set token for an existing instance
  canvas auth token set myschool --token 7~abc123...

  # Create a new instance with token auth
  canvas auth token set sandbox --url https://sandbox.instructure.com --token 7~xyz789...

  # Interactive mode (prompts for token)
  canvas auth token set myschool`,
		Args: ExactArgsWithUsage(1, "instance-name"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.InstanceName = config.SanitizeInstanceName(args[0])
			if err := opts.Validate(); err != nil {
				return err
			}
			return runAuthTokenSet(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Token, "token", "", "API access token")
	cmd.Flags().StringVar(&opts.URL, "url", "", "Canvas instance URL (required for new instances)")

	return cmd
}

func newAuthTokenRemoveCmd() *cobra.Command {
	opts := &options.AuthTokenRemoveOptions{}

	cmd := &cobra.Command{
		Use:   "remove <instance-name>",
		Short: "Remove API token from an instance",
		Long: `Remove the API token from a Canvas instance configuration.

This removes token-based authentication. If the instance also has OAuth
credentials, you can still use 'canvas auth login' to authenticate.

Examples:
  canvas auth token remove myschool`,
		Args: ExactArgsWithUsage(1, "instance-name"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.InstanceName = args[0]
			if err := opts.Validate(); err != nil {
				return err
			}
			return runAuthTokenRemove(cmd.Context(), opts)
		},
	}

	return cmd
}

func runAuthLogin(ctx context.Context, opts *options.AuthLoginOptions) error {
	// Load config to check for existing instances
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	var instanceURL string
	var existingInstance *config.Instance

	// If --instance is provided, try to look it up from config
	if opts.InstanceName != "" {
		if inst, err := cfg.GetInstance(opts.InstanceName); err == nil {
			existingInstance = inst
			instanceURL = inst.URL
			// Use stored client ID/secret if not provided via flags
			if opts.ClientID == "" && inst.ClientID != "" {
				opts.ClientID = inst.ClientID
			}
			if opts.ClientSecret == "" && inst.ClientSecret != "" {
				opts.ClientSecret = inst.ClientSecret
			}
		}
	}

	// If URL provided as positional arg, use it (overrides config lookup)
	if opts.InstanceURL != "" {
		instanceURL = opts.InstanceURL
	}

	// Still no URL? Error out
	if instanceURL == "" {
		if opts.InstanceName != "" {
			return fmt.Errorf("instance %q not found in config. Either add it first with 'canvas config add' or provide the URL", opts.InstanceName)
		}
		return fmt.Errorf("instance URL is required")
	}

	// Normalize URL
	normalizedURL, err := config.NormalizeURL(instanceURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Determine instance name if not provided
	if opts.InstanceName == "" {
		opts.InstanceName = getHostnameFromURL(normalizedURL)
	}

	opts.InstanceName = config.SanitizeInstanceName(opts.InstanceName)

	// If we found an existing instance, update the URL in case it was normalized
	if existingInstance != nil {
		existingInstance.URL = normalizedURL
	}

	fmt.Printf("🔐 Logging in to Canvas instance: %s\n", normalizedURL)
	fmt.Printf("Instance name: %s\n\n", opts.InstanceName)

	// Parse OAuth mode
	var oauthMode auth.OAuthMode
	switch opts.OAuthMode {
	case "auto":
		oauthMode = auth.OAuthModeAuto
	case "local":
		oauthMode = auth.OAuthModeLocal
	case "oob":
		oauthMode = auth.OAuthModeOOB
	default:
		return fmt.Errorf("invalid OAuth mode: %s (must be auto, local, or oob)", opts.OAuthMode)
	}

	// Get or prompt for client ID
	if opts.ClientID == "" {
		fmt.Print("Enter OAuth Client ID: ")
		fmt.Scanln(&opts.ClientID)
	}

	// If client ID is provided, also require client secret for OAuth
	if opts.ClientID != "" && opts.ClientSecret == "" {
		fmt.Print("Enter OAuth Client Secret: ")
		fmt.Scanln(&opts.ClientSecret)
		if opts.ClientSecret == "" {
			return fmt.Errorf("client secret is required when using OAuth with a client ID")
		}
	}

	// Create OAuth flow
	oauthFlow, err := auth.NewOAuthFlow(&auth.OAuthFlowConfig{
		BaseURL:      normalizedURL,
		ClientID:     opts.ClientID,
		ClientSecret: opts.ClientSecret,
		Mode:         oauthMode,
	})
	if err != nil {
		return fmt.Errorf("failed to create OAuth flow: %w", err)
	}

	// Perform authentication
	authCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	token, err := oauthFlow.Authenticate(authCtx)
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
	if err := tokenStore.Save(opts.InstanceName, token); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	// Add or update instance
	instance := &config.Instance{
		Name:         opts.InstanceName,
		URL:          normalizedURL,
		ClientID:     opts.ClientID,
		ClientSecret: opts.ClientSecret,
	}

	if _, exists := cfg.Instances[opts.InstanceName]; exists {
		if err := cfg.UpdateInstance(opts.InstanceName, instance); err != nil {
			return fmt.Errorf("failed to update instance: %w", err)
		}
	} else {
		if err := cfg.AddInstance(instance); err != nil {
			return fmt.Errorf("failed to add instance: %w", err)
		}
	}

	fmt.Printf("\n✓ Successfully authenticated with %s\n", opts.InstanceName)
	fmt.Printf("Token expires: %s\n", token.Expiry.Format(time.RFC3339))

	return nil
}

func runAuthLogout(ctx context.Context, opts *options.AuthLogoutOptions) error {
	// Determine instance name
	var instanceName string
	if opts.InstanceName != "" {
		instanceName = opts.InstanceName
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

func runAuthStatus(ctx context.Context, opts *options.AuthStatusOptions) error {
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
	if opts.InstanceName != "" {
		instance, err := cfg.GetInstance(opts.InstanceName)
		if err != nil {
			return err
		}

		printInstanceStatus(instance, cfg.DefaultInstance == opts.InstanceName, tokenStore)
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

	// Check authentication type and status
	if instance.HasToken() {
		// Token-based authentication
		fmt.Printf("   Auth: API Token\n")
		fmt.Printf("   Status: ✓ Configured (token does not expire)\n")
	} else if tokenStore.Exists(instance.Name) {
		// OAuth-based authentication
		fmt.Printf("   Auth: OAuth\n")
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
		fmt.Printf("   Auth: None\n")
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

func runAuthTokenSet(ctx context.Context, opts *options.AuthTokenSetOptions) error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if instance exists
	existingInstance, _ := cfg.GetInstance(opts.InstanceName)

	// Determine URL
	var instanceURL string
	if opts.URL != "" {
		instanceURL = opts.URL
	} else if existingInstance != nil {
		instanceURL = existingInstance.URL
	} else {
		return fmt.Errorf("instance %q not found. Provide --url to create a new instance", opts.InstanceName)
	}

	// Normalize URL
	normalizedURL, err := config.NormalizeURL(instanceURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Get token (from flag or prompt)
	token := opts.Token
	if token == "" {
		fmt.Print("Enter API Access Token: ")
		fmt.Scanln(&token)
		if token == "" {
			return fmt.Errorf("API token is required")
		}
	}

	// Create or update instance
	instance := &config.Instance{
		Name:  opts.InstanceName,
		URL:   normalizedURL,
		Token: token,
	}

	// Preserve OAuth credentials if updating existing instance
	if existingInstance != nil {
		instance.ClientID = existingInstance.ClientID
		instance.ClientSecret = existingInstance.ClientSecret
		instance.Description = existingInstance.Description

		if err := cfg.UpdateInstance(opts.InstanceName, instance); err != nil {
			return fmt.Errorf("failed to update instance: %w", err)
		}
		fmt.Printf("✓ Updated API token for %s\n", opts.InstanceName)
	} else {
		if err := cfg.AddInstance(instance); err != nil {
			return fmt.Errorf("failed to add instance: %w", err)
		}
		fmt.Printf("✓ Created instance %s with API token authentication\n", opts.InstanceName)
	}

	fmt.Printf("URL: %s\n", normalizedURL)
	fmt.Printf("Auth type: token\n")

	return nil
}

func runAuthTokenRemove(ctx context.Context, opts *options.AuthTokenRemoveOptions) error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if instance exists
	instance, err := cfg.GetInstance(opts.InstanceName)
	if err != nil {
		return err
	}

	if instance.Token == "" {
		return fmt.Errorf("instance %q does not have an API token configured", opts.InstanceName)
	}

	// Confirm removal
	fmt.Printf("Are you sure you want to remove the API token from %s? (y/N): ", opts.InstanceName)
	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "y" && confirm != "Y" {
		fmt.Println("Token removal cancelled")
		return nil
	}

	// Clear token
	instance.Token = ""

	if err := cfg.UpdateInstance(opts.InstanceName, instance); err != nil {
		return fmt.Errorf("failed to update instance: %w", err)
	}

	fmt.Printf("✓ Removed API token from %s\n", opts.InstanceName)

	// Suggest next steps if no auth remains
	if !instance.HasOAuth() {
		fmt.Printf("\nNote: Instance %s now has no authentication configured.\n", opts.InstanceName)
		fmt.Println("Use 'canvas auth login' or 'canvas auth token set' to authenticate.")
	}

	return nil
}
