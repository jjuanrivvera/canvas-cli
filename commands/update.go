package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
	"github.com/jjuanrivvera/canvas-cli/internal/output"
	"github.com/jjuanrivvera/canvas-cli/internal/updates"
	"github.com/spf13/cobra"
)

const (
	githubOwner = "jjuanrivvera"
	githubRepo  = "canvas-cli"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Manage CLI updates",
	Long: `Check for and install updates to the Canvas CLI.

Examples:
  # Check for updates
  canvas update check

  # Install the latest version
  canvas update install

  # Enable automatic update checks
  canvas update enable

  # Disable automatic update checks
  canvas update disable`,
}

// newUpdateCheckCmd creates the update check command
func newUpdateCheckCmd() *cobra.Command {
	opts := &options.UpdateCheckOptions{}

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check for available updates",
		Long:  "Check if a new version of Canvas CLI is available on GitHub.",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.NewCommandLogger(verbose)

			if err := opts.Validate(); err != nil {
				return err
			}

			logger.LogCommandStart(cmd.Context(), "update.check", map[string]interface{}{
				"force":  opts.Force,
				"format": opts.OutputFormat,
			})

			err := runUpdateCheck(cmd.Context(), opts)
			if err != nil {
				return err
			}

			logger.LogCommandComplete(cmd.Context(), "update.check", 1)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Force a fresh check, ignore cache")
	cmd.Flags().StringVarP(&opts.OutputFormat, "output", "o", "table", "Output format (table, json, yaml)")

	return cmd
}

// newUpdateInstallCmd creates the update install command
func newUpdateInstallCmd() *cobra.Command {
	opts := &options.UpdateInstallOptions{}

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install the latest version",
		Long: `Download and install the latest version of Canvas CLI.

Note: This command only works for binaries installed directly from GitHub releases.
If you installed via a package manager (Homebrew, etc.), use that to update instead.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.NewCommandLogger(verbose)

			if err := opts.Validate(); err != nil {
				return err
			}

			logger.LogCommandStart(cmd.Context(), "update.install", map[string]interface{}{
				"yes": opts.Yes,
			})

			err := runUpdateInstall(cmd.Context(), opts)
			if err != nil {
				return err
			}

			logger.LogCommandComplete(cmd.Context(), "update.install", 1)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&opts.Yes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}

// newUpdateEnableCmd creates the update enable command
func newUpdateEnableCmd() *cobra.Command {
	opts := &options.UpdateEnableOptions{
		Interval: 24, // default to 24 hours
	}

	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable automatic update checks",
		Long:  "Enable automatic checks for new versions (non-intrusive notification only).",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.NewCommandLogger(verbose)

			if err := opts.Validate(); err != nil {
				return err
			}

			logger.LogCommandStart(cmd.Context(), "update.enable", map[string]interface{}{
				"interval": opts.Interval,
			})

			err := runUpdateEnable(cmd.Context(), opts)
			if err != nil {
				return err
			}

			logger.LogCommandComplete(cmd.Context(), "update.enable", 1)
			return nil
		},
	}

	cmd.Flags().IntVar(&opts.Interval, "interval", 24, "Check interval in hours")

	return cmd
}

// newUpdateDisableCmd creates the update disable command
func newUpdateDisableCmd() *cobra.Command {
	opts := &options.UpdateDisableOptions{}

	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable automatic update checks",
		Long:  "Disable automatic checks for new versions.",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.NewCommandLogger(verbose)

			if err := opts.Validate(); err != nil {
				return err
			}

			logger.LogCommandStart(cmd.Context(), "update.disable", nil)

			err := runUpdateDisable(cmd.Context(), opts)
			if err != nil {
				return err
			}

			logger.LogCommandComplete(cmd.Context(), "update.disable", 1)
			return nil
		},
	}

	return cmd
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.AddCommand(newUpdateCheckCmd())
	updateCmd.AddCommand(newUpdateInstallCmd())
	updateCmd.AddCommand(newUpdateEnableCmd())
	updateCmd.AddCommand(newUpdateDisableCmd())
}

// runUpdateCheck checks for updates
func runUpdateCheck(ctx context.Context, opts *options.UpdateCheckOptions) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	cacheDir := filepath.Join(home, ".canvas-cli", "cache")

	checker := updates.NewChecker(updates.UpdateConfig{
		Owner:          githubOwner,
		Repo:           githubRepo,
		CurrentVersion: version,
		ForceCheck:     opts.Force,
		CacheTTL:       6 * time.Hour,
	}, cacheDir)

	result, err := checker.Check(ctx)
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	// Format output based on requested format
	format := opts.OutputFormat
	if format == "" {
		format = "table"
	}

	switch format {
	case "json", "yaml":
		// Use the output package for structured formats
		return output.WriteWithOptions(os.Stdout, result, output.FormatType(format), false)
	default:
		// Table format (human-readable)
		if result.UpdateAvailable {
			fmt.Printf("🎉 New version available!\n\n")
			fmt.Printf("Current version: %s\n", result.CurrentVersion)
			fmt.Printf("Latest version:  %s\n", result.LatestVersion)
			if result.ReleaseInfo != nil {
				fmt.Printf("Release date:    %s\n", result.ReleaseInfo.ReleaseDate.Format("2006-01-02"))
				fmt.Printf("Release URL:     %s\n", result.ReleaseInfo.URL)
			}
			fmt.Printf("\nTo update, run: canvas update install\n")
		} else {
			fmt.Printf("✓ You're running the latest version (%s)\n", result.CurrentVersion)
		}
		return nil
	}
}

// runUpdateInstall installs the latest version
func runUpdateInstall(ctx context.Context, opts *options.UpdateInstallOptions) error {
	installer := updates.NewInstaller(updates.UpdateConfig{
		Owner:          githubOwner,
		Repo:           githubRepo,
		CurrentVersion: version,
	})

	// Check if we can update
	canUpdate, reason := installer.CanUpdate()
	if !canUpdate {
		return fmt.Errorf("cannot update: %s", reason)
	}

	// Check for updates first
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	cacheDir := filepath.Join(home, ".canvas-cli", "cache")

	checker := updates.NewChecker(updates.UpdateConfig{
		Owner:          githubOwner,
		Repo:           githubRepo,
		CurrentVersion: version,
		ForceCheck:     true, // Always force check before install
	}, cacheDir)

	result, err := checker.Check(ctx)
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !result.UpdateAvailable {
		fmt.Printf("Already running the latest version (%s)\n", result.CurrentVersion)
		return nil
	}

	// Confirm with user unless --yes flag is set
	if !opts.Yes {
		fmt.Printf("Current version: %s\n", result.CurrentVersion)
		fmt.Printf("Latest version:  %s\n", result.LatestVersion)
		fmt.Printf("\nDo you want to install version %s? [y/N]: ", result.LatestVersion)

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Update cancelled")
			return nil
		}
	}

	// Perform the installation
	if err := installer.Install(ctx); err != nil {
		return fmt.Errorf("failed to install update: %w", err)
	}

	return nil
}

// runUpdateEnable enables automatic update checks
func runUpdateEnable(ctx context.Context, opts *options.UpdateEnableOptions) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	cfg.Settings.AutoUpdateCheck = true
	cfg.Settings.UpdateCheckInterval = opts.Interval

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✓ Automatic update checks enabled (interval: %d hours)\n", opts.Interval)
	return nil
}

// runUpdateDisable disables automatic update checks
func runUpdateDisable(ctx context.Context, opts *options.UpdateDisableOptions) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	cfg.Settings.AutoUpdateCheck = false

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("✓ Automatic update checks disabled")
	return nil
}
