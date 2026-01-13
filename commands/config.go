package commands

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"

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
	Args: cobra.ExactArgs(1),
	RunE: runConfigAdd,
}

var configUseCmd = &cobra.Command{
	Use:   "use [instance-name]",
	Short: "Switch to a different Canvas instance",
	Long: `Set the default Canvas instance to use for all commands.

Examples:
  canvas config use production
  canvas config use staging`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigUse,
}

var configRemoveCmd = &cobra.Command{
	Use:   "remove [instance-name]",
	Short: "Remove a Canvas instance from configuration",
	Long: `Remove a configured Canvas instance. This will not delete any data from Canvas,
only remove the instance from your local configuration.

Examples:
  canvas config remove staging`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigRemove,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration details",
	Long:  `Display the current configuration including default instance and settings.`,
	RunE:  runConfigShow,
}

// Config command flags
var (
	configURL         string
	configDescription string
	configClientID    string
	configForce       bool
)

func init() {
	rootCmd.AddCommand(configCmd)

	// Add subcommands
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configAddCmd)
	configCmd.AddCommand(configUseCmd)
	configCmd.AddCommand(configRemoveCmd)
	configCmd.AddCommand(configShowCmd)

	// Flags for add command
	configAddCmd.Flags().StringVar(&configURL, "url", "", "Canvas instance URL (required)")
	configAddCmd.Flags().StringVar(&configDescription, "description", "", "Instance description")
	configAddCmd.Flags().StringVar(&configClientID, "client-id", "", "OAuth client ID")
	configAddCmd.MarkFlagRequired("url")

	// Flags for remove command
	configRemoveCmd.Flags().BoolVarP(&configForce, "force", "f", false, "Skip confirmation prompt")
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
