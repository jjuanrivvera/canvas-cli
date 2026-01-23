package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/config"
	"github.com/jjuanrivvera/canvas-cli/internal/update"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for and install updates",
	Long: `Check for new versions of canvas-cli and install them automatically.

The CLI automatically checks for updates in the background. You can use this
command to check immediately or configure auto-update behavior.

Examples:
  canvas update              # Check and install updates now
  canvas update --check      # Check for updates without installing
  canvas update --disable    # Disable auto-updates
  canvas update --enable     # Enable auto-updates
  canvas update --status     # Show auto-update status`,
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.AddCommand(newUpdateCheckCmd())
	updateCmd.AddCommand(newUpdateEnableCmd())
	updateCmd.AddCommand(newUpdateDisableCmd())
	updateCmd.AddCommand(newUpdateStatusCmd())

	// Default action for `canvas update` (no subcommand)
	updateCmd.RunE = runUpdateNow
}

func newUpdateCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Check for updates without installing",
		RunE:  runUpdateCheck,
	}
}

func newUpdateEnableCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable automatic updates",
		RunE:  runUpdateEnable,
	}

	cmd.Flags().Int("interval", 60, "Check interval in minutes (default 60)")
	return cmd
}

func newUpdateDisableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disable",
		Short: "Disable automatic updates",
		RunE:  runUpdateDisable,
	}
}

func newUpdateStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show auto-update status and configuration",
		RunE:  runUpdateStatus,
	}
}

func runUpdateNow(cmd *cobra.Command, args []string) error {
	fmt.Println("Checking for updates...")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	updater := update.NewUpdater(version)
	release, err := updater.GetLatestRelease(ctx)
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	latestVersion := release.TagName
	fmt.Printf("Current version: %s\n", version)
	fmt.Printf("Latest version:  %s\n", latestVersion)

	// Check if update is needed
	if !needsUpdate(version, latestVersion) {
		fmt.Println("\nYou're already running the latest version!")
		return nil
	}

	fmt.Println("\nDownloading and installing update...")

	result := updater.CheckAndUpdate(ctx)
	if result.Error != nil {
		return fmt.Errorf("update failed: %w", result.Error)
	}

	if result.Updated {
		fmt.Printf("\n\033[32m✓ Successfully updated from v%s to v%s\033[0m\n", result.FromVersion, result.ToVersion)
		fmt.Println("  Restart the CLI to use the new version.")
	}

	return nil
}

func runUpdateCheck(cmd *cobra.Command, args []string) error {
	fmt.Println("Checking for updates...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	updater := update.NewUpdater(version)
	release, err := updater.GetLatestRelease(ctx)
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	latestVersion := release.TagName
	fmt.Printf("Current version: %s\n", version)
	fmt.Printf("Latest version:  %s\n", latestVersion)

	if needsUpdate(version, latestVersion) {
		fmt.Printf("\n\033[33mA new version is available!\033[0m\n")
		fmt.Println("Run 'canvas update' to install it.")
	} else {
		fmt.Println("\nYou're already running the latest version!")
	}

	return nil
}

func runUpdateEnable(cmd *cobra.Command, args []string) error {
	interval, _ := cmd.Flags().GetInt("interval")
	if interval < 1 {
		interval = 60
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.Settings == nil {
		cfg.Settings = config.DefaultSettings()
	}

	cfg.Settings.AutoUpdateEnabled = true
	cfg.Settings.AutoUpdateIntervalMin = interval

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Auto-updates enabled (checking every %d minutes)\n", interval)
	return nil
}

func runUpdateDisable(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.Settings == nil {
		cfg.Settings = config.DefaultSettings()
	}

	cfg.Settings.AutoUpdateEnabled = false

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("Auto-updates disabled")
	fmt.Println("Run 'canvas update enable' to re-enable automatic updates.")
	return nil
}

func runUpdateStatus(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("Auto-Update Status")
	fmt.Println("==================")
	fmt.Printf("Current version: %s\n", version)

	if cfg.Settings != nil {
		status := "disabled"
		if cfg.Settings.AutoUpdateEnabled {
			status = "enabled"
		}
		fmt.Printf("Auto-update: %s\n", status)
		fmt.Printf("Check interval: %d minutes\n", cfg.Settings.AutoUpdateIntervalMin)
	} else {
		fmt.Println("Auto-update: enabled (default)")
		fmt.Println("Check interval: 60 minutes (default)")
	}

	// Show last check info
	stateManager, err := update.NewStateManager()
	if err == nil {
		state, err := stateManager.Load()
		if err == nil && !state.LastCheckTime.IsZero() {
			fmt.Printf("\nLast check: %s\n", state.LastCheckTime.Format(time.RFC3339))
			if state.LastError != "" && !state.LastErrorTime.IsZero() {
				fmt.Printf("Last error: %s (%s)\n", state.LastError, state.LastErrorTime.Format(time.RFC3339))
			}
		}
	}

	return nil
}

// needsUpdate compares versions to determine if an update is needed
func needsUpdate(current, latest string) bool {
	// Strip 'v' prefix for comparison
	current = trimVersionPrefix(current)
	latest = trimVersionPrefix(latest)

	// Skip dev versions
	if current == "dev" || current == "" {
		return false
	}

	return current != latest
}

func trimVersionPrefix(v string) string {
	if len(v) > 0 && v[0] == 'v' {
		return v[1:]
	}
	return v
}
