package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
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
  canvas update check        # Check for updates without installing
  canvas update disable      # Disable auto-updates
  canvas update enable       # Enable auto-updates
  canvas update status       # Show auto-update status`,
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.AddCommand(newUpdateCheckCmd())
	updateCmd.AddCommand(newUpdateEnableCmd())
	updateCmd.AddCommand(newUpdateDisableCmd())
	updateCmd.AddCommand(newUpdateStatusCmd())

	// Default action for `canvas update` (no subcommand)
	updateCmd.RunE = func(cmd *cobra.Command, args []string) error {
		opts := &options.UpdateOptions{}
		if err := opts.Validate(); err != nil {
			return err
		}
		return runUpdateNow(cmd.Context(), opts)
	}
}

func newUpdateCheckCmd() *cobra.Command {
	opts := &options.UpdateCheckOptions{}

	return &cobra.Command{
		Use:   "check",
		Short: "Check for updates without installing",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			return runUpdateCheck(cmd.Context(), opts)
		},
	}
}

func newUpdateEnableCmd() *cobra.Command {
	opts := &options.UpdateEnableOptions{}

	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable automatic updates",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			return runUpdateEnable(cmd.Context(), opts)
		},
	}

	cmd.Flags().IntVar(&opts.Interval, "interval", 60, "Check interval in minutes (default 60)")
	return cmd
}

func newUpdateDisableCmd() *cobra.Command {
	opts := &options.UpdateDisableOptions{}

	return &cobra.Command{
		Use:   "disable",
		Short: "Disable automatic updates",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			return runUpdateDisable(cmd.Context(), opts)
		},
	}
}

func newUpdateStatusCmd() *cobra.Command {
	opts := &options.UpdateStatusOptions{}

	return &cobra.Command{
		Use:   "status",
		Short: "Show auto-update status and configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			return runUpdateStatus(cmd.Context(), opts)
		},
	}
}

func runUpdateNow(ctx context.Context, opts *options.UpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "update.now", map[string]interface{}{
		"current_version": version,
	})

	fmt.Println("Checking for updates...")

	updateCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	updater := update.NewUpdater(version)
	release, err := updater.GetLatestRelease(updateCtx)
	if err != nil {
		logger.LogCommandError(ctx, "update.now", err, map[string]interface{}{})
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	latestVersion := release.TagName
	fmt.Printf("Current version: %s\n", version)
	fmt.Printf("Latest version:  %s\n", latestVersion)

	// Check if update is needed using semantic version comparison
	if !isNewerVersion(latestVersion, version) {
		fmt.Println("\nYou're already running the latest version!")
		logger.LogCommandComplete(ctx, "update.now", 0)
		return nil
	}

	fmt.Println("\nDownloading and installing update...")

	result := updater.CheckAndUpdate(updateCtx)
	if result.Error != nil {
		logger.LogCommandError(ctx, "update.now", result.Error, map[string]interface{}{
			"from_version": result.FromVersion,
			"to_version":   result.ToVersion,
		})
		return fmt.Errorf("update failed: %w", result.Error)
	}

	if result.Updated {
		fmt.Printf("\n\033[32m✓ Successfully updated from v%s to v%s\033[0m\n", result.FromVersion, result.ToVersion)
		fmt.Println("  Restart the CLI to use the new version.")
	}

	logger.LogCommandComplete(ctx, "update.now", 1)
	return nil
}

func runUpdateCheck(ctx context.Context, opts *options.UpdateCheckOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "update.check", map[string]interface{}{
		"current_version": version,
	})

	fmt.Println("Checking for updates...")

	checkCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	updater := update.NewUpdater(version)
	release, err := updater.GetLatestRelease(checkCtx)
	if err != nil {
		logger.LogCommandError(ctx, "update.check", err, map[string]interface{}{})
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	latestVersion := release.TagName
	fmt.Printf("Current version: %s\n", version)
	fmt.Printf("Latest version:  %s\n", latestVersion)

	if isNewerVersion(latestVersion, version) {
		fmt.Printf("\n\033[33mA new version is available!\033[0m\n")
		fmt.Println("Run 'canvas update' to install it.")
	} else {
		fmt.Println("\nYou're already running the latest version!")
	}

	logger.LogCommandComplete(ctx, "update.check", 1)
	return nil
}

func runUpdateEnable(ctx context.Context, opts *options.UpdateEnableOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "update.enable", map[string]interface{}{
		"interval": opts.Interval,
	})

	cfg, err := config.Load()
	if err != nil {
		logger.LogCommandError(ctx, "update.enable", err, map[string]interface{}{})
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.Settings == nil {
		cfg.Settings = config.DefaultSettings()
	}

	cfg.Settings.AutoUpdateEnabled = true
	cfg.Settings.AutoUpdateIntervalMin = opts.Interval

	if err := cfg.Save(); err != nil {
		logger.LogCommandError(ctx, "update.enable", err, map[string]interface{}{})
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Auto-updates enabled (checking every %d minutes)\n", opts.Interval)
	logger.LogCommandComplete(ctx, "update.enable", 1)
	return nil
}

func runUpdateDisable(ctx context.Context, opts *options.UpdateDisableOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "update.disable", map[string]interface{}{})

	cfg, err := config.Load()
	if err != nil {
		logger.LogCommandError(ctx, "update.disable", err, map[string]interface{}{})
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.Settings == nil {
		cfg.Settings = config.DefaultSettings()
	}

	cfg.Settings.AutoUpdateEnabled = false

	if err := cfg.Save(); err != nil {
		logger.LogCommandError(ctx, "update.disable", err, map[string]interface{}{})
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("Auto-updates disabled")
	fmt.Println("Run 'canvas update enable' to re-enable automatic updates.")
	logger.LogCommandComplete(ctx, "update.disable", 1)
	return nil
}

func runUpdateStatus(ctx context.Context, opts *options.UpdateStatusOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "update.status", map[string]interface{}{})

	cfg, err := config.Load()
	if err != nil {
		logger.LogCommandError(ctx, "update.status", err, map[string]interface{}{})
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

	logger.LogCommandComplete(ctx, "update.status", 1)
	return nil
}

// isNewerVersion compares two semver versions
// Returns true if latest is newer than current
func isNewerVersion(latest, current string) bool {
	// Strip 'v' prefix for comparison
	latest = strings.TrimPrefix(latest, "v")
	current = strings.TrimPrefix(current, "v")

	// Skip dev versions
	if current == "dev" || current == "" {
		return false
	}

	latestParts := parseVersion(latest)
	currentParts := parseVersion(current)

	for i := 0; i < 3; i++ {
		if latestParts[i] > currentParts[i] {
			return true
		}
		if latestParts[i] < currentParts[i] {
			return false
		}
	}

	return false
}

// parseVersion parses a semver string into [major, minor, patch]
func parseVersion(v string) [3]int {
	v = strings.TrimPrefix(v, "v")
	parts := strings.Split(v, ".")

	var result [3]int
	for i := 0; i < 3 && i < len(parts); i++ {
		// Strip any pre-release suffix
		numStr := strings.Split(parts[i], "-")[0]
		fmt.Sscanf(numStr, "%d", &result[i])
	}

	return result
}
