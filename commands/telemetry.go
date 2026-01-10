package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jjuanrivvera/canvas-cli/internal/config"
	"github.com/spf13/cobra"
)

var telemetryCmd = &cobra.Command{
	Use:   "telemetry",
	Short: "Manage telemetry settings",
	Long: `Manage opt-in telemetry settings for Canvas CLI.

Telemetry helps improve Canvas CLI by collecting anonymous usage data
including command execution, errors, and performance metrics. No personal
information or Canvas data is ever collected.

Data collected:
  - Command usage (which commands are run)
  - Error rates and types
  - Performance metrics (command duration)
  - OS and architecture
  - Canvas CLI version

Data NOT collected:
  - Canvas credentials or tokens
  - Course content or user data
  - Personal information
  - File contents or names

All telemetry data is stored locally and never automatically transmitted.

Examples:
  # Enable telemetry
  canvas telemetry enable

  # Disable telemetry
  canvas telemetry disable

  # Check telemetry status
  canvas telemetry status

  # View collected data
  canvas telemetry show

  # Clear telemetry data
  canvas telemetry clear`,
}

var telemetryEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable telemetry collection",
	Long: `Enable anonymous usage data collection.

By enabling telemetry, you help improve Canvas CLI. All data is
collected anonymously and stored locally. You can disable telemetry
at any time or clear all collected data.`,
	RunE: runTelemetryEnable,
}

var telemetryDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable telemetry collection",
	Long:  `Disable telemetry collection. Previously collected data is preserved unless you run 'canvas telemetry clear'.`,
	RunE:  runTelemetryDisable,
}

var telemetryStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show telemetry status",
	Long:  `Display current telemetry configuration and statistics.`,
	RunE:  runTelemetryStatus,
}

var telemetryShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show collected telemetry data",
	Long:  `Display telemetry data files and their contents.`,
	RunE:  runTelemetryShow,
}

var telemetryClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all telemetry data",
	Long:  `Delete all collected telemetry data. This does not disable telemetry.`,
	RunE:  runTelemetryClear,
}

func init() {
	rootCmd.AddCommand(telemetryCmd)
	telemetryCmd.AddCommand(telemetryEnableCmd)
	telemetryCmd.AddCommand(telemetryDisableCmd)
	telemetryCmd.AddCommand(telemetryStatusCmd)
	telemetryCmd.AddCommand(telemetryShowCmd)
	telemetryCmd.AddCommand(telemetryClearCmd)
}

func runTelemetryEnable(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.Settings == nil {
		cfg.Settings = config.DefaultSettings()
	}

	cfg.Settings.TelemetryEnabled = true

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("✓ Telemetry enabled")
	fmt.Println()
	fmt.Println("Thank you for helping improve Canvas CLI!")
	fmt.Println()
	fmt.Println("What's collected:")
	fmt.Println("  • Command usage and performance")
	fmt.Println("  • Error rates and types")
	fmt.Println("  • System information (OS, architecture)")
	fmt.Println()
	fmt.Println("What's NOT collected:")
	fmt.Println("  • Canvas credentials or tokens")
	fmt.Println("  • Course content or user data")
	fmt.Println("  • Personal information")
	fmt.Println()
	fmt.Println("You can disable telemetry anytime with: canvas telemetry disable")

	return nil
}

func runTelemetryDisable(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.Settings == nil {
		cfg.Settings = config.DefaultSettings()
	}

	cfg.Settings.TelemetryEnabled = false

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("✓ Telemetry disabled")
	fmt.Println()
	fmt.Println("Note: Previously collected data is still stored.")
	fmt.Println("To remove it, run: canvas telemetry clear")

	return nil
}

func runTelemetryStatus(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	enabled := false
	if cfg.Settings != nil {
		enabled = cfg.Settings.TelemetryEnabled
	}

	fmt.Println("Telemetry Status")
	fmt.Println("================")
	fmt.Println()

	if enabled {
		fmt.Println("Status: ✓ Enabled")
	} else {
		fmt.Println("Status: ✗ Disabled")
	}

	// Check for telemetry data
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	telemetryDir := filepath.Join(home, ".canvas-cli", "telemetry")
	files, err := os.ReadDir(telemetryDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Data files: 0")
			return nil
		}
		return fmt.Errorf("failed to read telemetry directory: %w", err)
	}

	// Count event files
	eventFiles := 0
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			eventFiles++
		}
	}

	fmt.Printf("Data files: %d\n", eventFiles)
	fmt.Printf("Data directory: %s\n", telemetryDir)
	fmt.Println()

	if enabled {
		fmt.Println("To disable: canvas telemetry disable")
	} else {
		fmt.Println("To enable: canvas telemetry enable")
	}

	if eventFiles > 0 {
		fmt.Println("To view data: canvas telemetry show")
		fmt.Println("To clear data: canvas telemetry clear")
	}

	return nil
}

func runTelemetryShow(cmd *cobra.Command, args []string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	telemetryDir := filepath.Join(home, ".canvas-cli", "telemetry")
	files, err := os.ReadDir(telemetryDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No telemetry data collected yet.")
			return nil
		}
		return fmt.Errorf("failed to read telemetry directory: %w", err)
	}

	// List event files
	eventFiles := make([]string, 0)
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			eventFiles = append(eventFiles, file.Name())
		}
	}

	if len(eventFiles) == 0 {
		fmt.Println("No telemetry data collected yet.")
		return nil
	}

	fmt.Println("Telemetry Data Files")
	fmt.Println("====================")
	fmt.Println()

	for _, filename := range eventFiles {
		filepath := filepath.Join(telemetryDir, filename)
		info, err := os.Stat(filepath)
		if err != nil {
			continue
		}

		fmt.Printf("File: %s\n", filename)
		fmt.Printf("Size: %d bytes\n", info.Size())
		fmt.Printf("Modified: %s\n", info.ModTime().Format("2006-01-02 15:04:05"))
		fmt.Printf("Path: %s\n", filepath)
		fmt.Println()
	}

	fmt.Printf("Total files: %d\n", len(eventFiles))
	fmt.Printf("Directory: %s\n", telemetryDir)

	return nil
}

func runTelemetryClear(cmd *cobra.Command, args []string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	telemetryDir := filepath.Join(home, ".canvas-cli", "telemetry")

	// Check if directory exists
	if _, err := os.Stat(telemetryDir); os.IsNotExist(err) {
		fmt.Println("No telemetry data to clear.")
		return nil
	}

	files, err := os.ReadDir(telemetryDir)
	if err != nil {
		return fmt.Errorf("failed to read telemetry directory: %w", err)
	}

	// Count and remove event files
	removed := 0
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			filepath := filepath.Join(telemetryDir, file.Name())
			if err := os.Remove(filepath); err != nil {
				fmt.Printf("Warning: failed to remove %s: %v\n", file.Name(), err)
			} else {
				removed++
			}
		}
	}

	if removed == 0 {
		fmt.Println("No telemetry data to clear.")
	} else {
		fmt.Printf("✓ Cleared %d telemetry data file(s)\n", removed)
	}

	return nil
}
