package commands

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"github.com/jjuanrivvera/canvas-cli/internal/auth"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Canvas CLI configuration",
	Long: `Manage Canvas CLI configuration including Canvas instances and settings.

Examples:
  canvas config list                              # List all configured instances
  canvas config add prod --url https://canvas.example.com
  canvas config use prod                          # Switch to prod instance
  canvas config show                              # Show current configuration
  canvas config remove staging                    # Remove an instance`,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured Canvas instances",
	Long:  `List all configured Canvas instances with their URLs and status.`,
	RunE:  runConfigList,
}

var configAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new Canvas instance",
	Long: `Add a new Canvas instance to the configuration.

Examples:
  canvas config add production --url https://canvas.example.com
  canvas config add staging --url https://canvas-staging.example.com --description "Staging environment"`,
	Args: ExactArgsWithUsage(1, "name"),
	RunE: runConfigAdd,
}

var configUseCmd = &cobra.Command{
	Use:   "use [instance-name]",
	Short: "Switch to a different Canvas instance",
	Long: `Set the default Canvas instance to use for all commands.

Examples:
  canvas config use production
  canvas config use staging`,
	Args: ExactArgsWithUsage(1, "instance-name"),
	RunE: runConfigUse,
}

var configRemoveCmd = &cobra.Command{
	Use:   "remove [instance-name]",
	Short: "Remove a Canvas instance from configuration",
	Long: `Remove a configured Canvas instance. This will not delete any data from Canvas,
only remove the instance from your local configuration.

Examples:
  canvas config remove staging`,
	Args: ExactArgsWithUsage(1, "instance-name"),
	RunE: runConfigRemove,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration details",
	Long:  `Display the current configuration including default instance and settings.`,
	RunE:  runConfigShow,
}

var configAccountCmd = &cobra.Command{
	Use:   "account [instance-name] [account-id]",
	Short: "Set the default account ID for an instance",
	Long: `Set or auto-detect the default account ID for a Canvas instance.

The default account ID is used when API calls require an account ID but none is specified.

Examples:
  # Set account ID manually
  canvas config account production 1
  canvas config account staging 42

  # Auto-detect account ID (fetches accounts from API)
  canvas config account production --detect
  canvas config account --detect              # Uses default instance`,
	Args: cobra.MaximumNArgs(2),
	RunE: runConfigAccount,
}

// Config command flags
var (
	configURL         string
	configDescription string
	configClientID    string
	configForce       bool
	configDetect      bool
)

func init() {
	rootCmd.AddCommand(configCmd)

	// Add subcommands
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configAddCmd)
	configCmd.AddCommand(configUseCmd)
	configCmd.AddCommand(configRemoveCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configAccountCmd)

	// Flags for add command
	configAddCmd.Flags().StringVar(&configURL, "url", "", "Canvas instance URL (required)")
	configAddCmd.Flags().StringVar(&configDescription, "description", "", "Instance description")
	configAddCmd.Flags().StringVar(&configClientID, "client-id", "", "OAuth client ID")
	configAddCmd.MarkFlagRequired("url")

	// Flags for remove command
	configRemoveCmd.Flags().BoolVarP(&configForce, "force", "f", false, "Skip confirmation prompt")

	// Flags for account command
	configAccountCmd.Flags().BoolVar(&configDetect, "detect", false, "Auto-detect account ID from Canvas API")
}

func runConfigList(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	instances := cfg.ListInstances()
	if len(instances) == 0 {
		fmt.Println("No instances configured.")
		fmt.Println("\nTo add an instance:")
		fmt.Println("  canvas config add <name> --url <canvas-url>")
		return nil
	}

	// Print header
	fmt.Printf("%-15s %-45s %-10s %s\n", "NAME", "URL", "DEFAULT", "DESCRIPTION")
	fmt.Println(strings.Repeat("-", 90))

	// Print instances
	for _, instance := range instances {
		isDefault := ""
		if instance.Name == cfg.DefaultInstance {
			isDefault = "*"
		}

		description := instance.Description
		if len(description) > 20 {
			description = description[:17] + "..."
		}

		fmt.Printf("%-15s %-45s %-10s %s\n",
			instance.Name,
			instance.URL,
			isDefault,
			description,
		)
	}

	fmt.Printf("\nTotal: %d instance(s)\n", len(instances))
	return nil
}

func runConfigAdd(_ *cobra.Command, args []string) error {
	instanceName := args[0]

	// Validate URL
	parsedURL, err := url.Parse(configURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme == "" {
		configURL = "https://" + configURL
		parsedURL, _ = url.Parse(configURL)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL must use http or https scheme")
	}

	// Remove trailing slash
	configURL = strings.TrimSuffix(configURL, "/")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	instance := &config.Instance{
		Name:        instanceName,
		URL:         configURL,
		Description: configDescription,
		ClientID:    configClientID,
	}

	if err := cfg.AddInstance(instance); err != nil {
		return fmt.Errorf("failed to add instance: %w", err)
	}

	fmt.Printf("Instance %q added successfully.\n", instanceName)

	if cfg.DefaultInstance == instanceName {
		fmt.Printf("Set as default instance.\n")
	}

	fmt.Println("\nNext step: Authenticate with this instance:")
	fmt.Printf("  canvas auth login --instance %s\n", instanceName)

	return nil
}

func runConfigUse(_ *cobra.Command, args []string) error {
	instanceName := args[0]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.SetDefaultInstance(instanceName); err != nil {
		return fmt.Errorf("failed to set default instance: %w", err)
	}

	fmt.Printf("Switched to instance %q.\n", instanceName)
	return nil
}

func runConfigRemove(_ *cobra.Command, args []string) error {
	instanceName := args[0]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if instance exists
	if _, err := cfg.GetInstance(instanceName); err != nil {
		return fmt.Errorf("instance %q not found", instanceName)
	}

	// Confirm removal unless --force is used
	if !configForce {
		fmt.Printf("Are you sure you want to remove instance %q? [y/N]: ", instanceName)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	if err := cfg.RemoveInstance(instanceName); err != nil {
		return fmt.Errorf("failed to remove instance: %w", err)
	}

	fmt.Printf("Instance %q removed.\n", instanceName)

	if cfg.DefaultInstance != "" {
		fmt.Printf("New default instance: %s\n", cfg.DefaultInstance)
	}

	return nil
}

func runConfigShow(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	configPath, _ := config.GetConfigPath()

	fmt.Println("Canvas CLI Configuration")
	fmt.Println(strings.Repeat("=", 50))

	fmt.Printf("\nConfig file: %s\n", configPath)
	fmt.Printf("Default instance: %s\n", valueOrNone(cfg.DefaultInstance))

	fmt.Println("\nInstances:")
	if len(cfg.Instances) == 0 {
		fmt.Println("  (none)")
	} else {
		for name, instance := range cfg.Instances {
			marker := "  "
			if name == cfg.DefaultInstance {
				marker = "* "
			}
			fmt.Printf("%s%s: %s\n", marker, name, instance.URL)
			if instance.Description != "" {
				fmt.Printf("    Description: %s\n", instance.Description)
			}
			if instance.ClientID != "" {
				fmt.Printf("    Client ID: %s\n", instance.ClientID)
			}
			if instance.HasDefaultAccountID() {
				fmt.Printf("    Default Account ID: %d\n", instance.DefaultAccountID)
			}
		}
	}

	fmt.Println("\nSettings:")
	if cfg.Settings != nil {
		fmt.Printf("  Output format: %s\n", cfg.Settings.DefaultOutputFormat)
		fmt.Printf("  Requests/sec: %.1f\n", cfg.Settings.RequestsPerSecond)
		fmt.Printf("  Cache enabled: %t\n", cfg.Settings.CacheEnabled)
		fmt.Printf("  Cache TTL: %d minutes\n", cfg.Settings.CacheTTL)
		fmt.Printf("  Telemetry: %t\n", cfg.Settings.TelemetryEnabled)
		fmt.Printf("  Log level: %s\n", cfg.Settings.LogLevel)
	}

	return nil
}

func valueOrNone(s string) string {
	if s == "" {
		return "(none)"
	}
	return s
}

func runConfigAccount(_ *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Determine instance name
	var instanceName string
	if len(args) >= 1 {
		instanceName = args[0]
	} else {
		// Use default instance
		if cfg.DefaultInstance == "" {
			return fmt.Errorf("no instance specified and no default instance configured")
		}
		instanceName = cfg.DefaultInstance
	}

	// Verify instance exists
	instance, err := cfg.GetInstance(instanceName)
	if err != nil {
		return fmt.Errorf("instance %q not found", instanceName)
	}

	// Check for detect mode
	if configDetect {
		return runConfigAccountDetect(cfg, instance)
	}

	// Manual mode - requires account ID argument
	if len(args) < 2 {
		// Show current setting if no account ID provided
		if instance.HasDefaultAccountID() {
			fmt.Printf("Instance %q default account ID: %d\n", instanceName, instance.DefaultAccountID)
		} else {
			fmt.Printf("Instance %q has no default account ID configured.\n", instanceName)
			fmt.Println("\nTo set an account ID:")
			fmt.Printf("  canvas config account %s <account-id>\n", instanceName)
			fmt.Printf("  canvas config account %s --detect\n", instanceName)
		}
		return nil
	}

	// Parse account ID
	accountID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid account ID: %s", args[1])
	}

	if accountID <= 0 {
		return fmt.Errorf("account ID must be a positive number")
	}

	// Set the account ID
	if err := cfg.SetDefaultAccountID(instanceName, accountID); err != nil {
		return fmt.Errorf("failed to set default account ID: %w", err)
	}

	fmt.Printf("Default account ID for %q set to %d.\n", instanceName, accountID)
	return nil
}

func runConfigAccountDetect(cfg *config.Config, instance *config.Instance) error {
	fmt.Printf("Detecting accounts for instance %q...\n", instance.Name)

	// Create API client for this instance
	client, err := getAPIClientForInstanceByName(instance.Name)
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	// Fetch accounts with a longer timeout for slow connections
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	accountsService := api.NewAccountsService(client)
	accounts, err := accountsService.List(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch accounts: %w", err)
	}

	if len(accounts) == 0 {
		return fmt.Errorf("no accounts found. You may not have permission to view any accounts")
	}

	var selectedAccountID int64

	if len(accounts) == 1 {
		// Single account - use it automatically
		selectedAccountID = accounts[0].ID
		fmt.Printf("Found 1 account: %s (ID: %d)\n", accounts[0].Name, accounts[0].ID)
	} else {
		// Multiple accounts - prompt for selection
		fmt.Printf("\nFound %d accounts:\n\n", len(accounts))
		fmt.Printf("  %-6s %-40s %s\n", "NUM", "NAME", "ID")
		fmt.Println("  " + strings.Repeat("-", 60))

		for i, account := range accounts {
			name := account.Name
			if len(name) > 38 {
				name = name[:35] + "..."
			}
			fmt.Printf("  %-6d %-40s %d\n", i+1, name, account.ID)
		}

		fmt.Print("\nSelect account number: ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		response = strings.TrimSpace(response)
		selection, err := strconv.Atoi(response)
		if err != nil || selection < 1 || selection > len(accounts) {
			return fmt.Errorf("invalid selection: %s", response)
		}

		selectedAccountID = accounts[selection-1].ID
	}

	// Save the selected account ID
	if err := cfg.SetDefaultAccountID(instance.Name, selectedAccountID); err != nil {
		return fmt.Errorf("failed to save account ID: %w", err)
	}

	fmt.Printf("\nDefault account ID for %q set to %d.\n", instance.Name, selectedAccountID)
	return nil
}

// getAPIClientForInstanceByName creates an API client for a specific instance by name
// This version supports both token auth and OAuth, with a longer timeout for interactive operations
func getAPIClientForInstanceByName(instanceName string) (*api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	instance, err := cfg.GetInstance(instanceName)
	if err != nil {
		return nil, fmt.Errorf("instance not found: %w", err)
	}

	var clientConfig api.ClientConfig

	// Check if instance has an API token configured
	if instance.HasToken() {
		clientConfig = api.ClientConfig{
			BaseURL:        instance.URL,
			Token:          instance.Token,
			RequestsPerSec: cfg.Settings.RequestsPerSecond,
			UserAgent:      getUserAgent(),
			Timeout:        60 * time.Second, // Longer timeout for interactive operations
		}
	} else {
		// OAuth flow - load token from store
		configDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		configDir = configDir + "/.canvas-cli"

		tokenStore := auth.NewFallbackTokenStore(configDir)
		token, err := tokenStore.Load(instance.Name)
		if err != nil {
			return nil, fmt.Errorf("not authenticated with %s. Run 'canvas auth login' or 'canvas auth token set' first", instance.Name)
		}

		if instance.HasOAuth() {
			oauth2Config := auth.CreateOAuth2ConfigForInstance(instance.URL, instance.ClientID, instance.ClientSecret)
			tokenSource := auth.NewAutoRefreshTokenSource(oauth2Config, tokenStore, instance.Name, token)

			clientConfig = api.ClientConfig{
				BaseURL:        instance.URL,
				TokenSource:    tokenSource,
				RequestsPerSec: cfg.Settings.RequestsPerSecond,
				UserAgent:      getUserAgent(),
				Timeout:        60 * time.Second, // Longer timeout for interactive operations
			}
		} else {
			clientConfig = api.ClientConfig{
				BaseURL:        instance.URL,
				Token:          token.AccessToken,
				RequestsPerSec: cfg.Settings.RequestsPerSecond,
				UserAgent:      getUserAgent(),
				Timeout:        60 * time.Second, // Longer timeout for interactive operations
			}
		}
	}

	return api.NewClient(clientConfig)
}
