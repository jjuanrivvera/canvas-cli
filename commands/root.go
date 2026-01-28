package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jjuanrivvera/canvas-cli/internal/config"
	"github.com/jjuanrivvera/canvas-cli/internal/update"
)

var (
	cfgFile      string
	instanceURL  string
	outputFormat string
	verbose      bool
	noCache      bool  // Disable caching for API requests
	asUserID     int64 // Masquerading: act as another user
	globalLimit  int   // Global limit for list operations
	dryRun       bool  // Print curl commands instead of executing
	showToken    bool  // Show actual token in dry-run output
	quiet        bool  // Suppress informational messages
	version      string
	commit       string
	buildDate    string

	// Output filtering flags
	filterText    string   // Filter results by text (substring match)
	filterColumns []string // Select specific columns to display
	sortField     string   // Sort results by field

	// Auto-updater instance
	autoUpdater *update.AutoUpdater
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "canvas",
	Short: "Canvas LMS CLI - Interact with Canvas from the command line",
	Long: `canvas-cli is a powerful command-line interface for Canvas LMS.
It provides comprehensive access to Canvas API features including courses,
assignments, users, submissions, and more.

Examples:
  canvas auth login                                              # Authenticate with Canvas
  canvas courses list                                            # List all courses
  canvas assignments list --course-id 123                        # List assignments for a course
  canvas submissions bulk-grade --course-id 123 --csv grades.csv # Bulk grade from CSV`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize and run auto-updater asynchronously
		initAutoUpdater()
		if autoUpdater != nil {
			autoUpdater.RunUpdateCheckAsync(context.Background())
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// Print any update notifications after command completes
		if autoUpdater != nil {
			// Wait for async update check to complete (with timeout to avoid blocking forever)
			autoUpdater.WaitForCompletion(5 * time.Second)
			autoUpdater.PrintNotifications()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute(v, c, bd string) error {
	version = v
	commit = c
	buildDate = bd

	// Set version information
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("canvas-cli version {{.Version}}\n")

	return rootCmd.Execute()
}

// ExecuteContext is like Execute but accepts a context for signal handling.
func ExecuteContext(ctx context.Context, v, c, bd string) error {
	version = v
	commit = c
	buildDate = bd

	rootCmd.Version = version
	rootCmd.SetVersionTemplate("canvas-cli version {{.Version}}\n")

	return rootCmd.ExecuteContext(ctx)
}

// GetRootCmd returns the root command for documentation generation
func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.canvas-cli/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&instanceURL, "instance", "", "Canvas instance URL (overrides config)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table, json, yaml, csv")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().Int64Var(&asUserID, "as-user", 0, "Masquerade as another user (admin feature, requires permission)")

	rootCmd.PersistentFlags().BoolVar(&noCache, "no-cache", false, "Disable caching of API responses")
	rootCmd.PersistentFlags().IntVar(&globalLimit, "limit", 0, "Limit number of results for list operations (0 = unlimited)")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Print curl commands instead of executing requests")
	rootCmd.PersistentFlags().BoolVar(&showToken, "show-token", false, "Show actual token in dry-run output (default: redacted)")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Suppress informational messages, only output data and errors")

	// Output filtering flags
	rootCmd.PersistentFlags().StringVar(&filterText, "filter", "", "Filter results by text (case-insensitive substring match)")
	rootCmd.PersistentFlags().StringSliceVar(&filterColumns, "columns", nil, "Select specific columns to display (comma-separated)")
	rootCmd.PersistentFlags().StringVar(&sortField, "sort", "", "Sort results by field (prefix with - for descending, e.g., -name)")

	// Bind flags to viper
	viper.BindPFlag("instance", rootCmd.PersistentFlags().Lookup("instance"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("as-user", rootCmd.PersistentFlags().Lookup("as-user"))
	viper.BindPFlag("no-cache", rootCmd.PersistentFlags().Lookup("no-cache"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		configDir, err := config.GetConfigDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// Environment variables
	viper.SetEnvPrefix("CANVAS")
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// initAutoUpdater initializes the auto-updater based on config settings
func initAutoUpdater() {
	// Skip if already initialized
	if autoUpdater != nil {
		return
	}

	// Load config to check if auto-update is enabled
	cfg, err := config.Load()
	if err != nil {
		// If config fails, use default settings (auto-update enabled)
		cfg = &config.Config{Settings: config.DefaultSettings()}
	}

	// Default to enabled if settings are nil
	enabled := true
	interval := 60 * time.Minute

	if cfg.Settings != nil {
		enabled = cfg.Settings.AutoUpdateEnabled
		if cfg.Settings.AutoUpdateIntervalMin > 0 {
			interval = time.Duration(cfg.Settings.AutoUpdateIntervalMin) * time.Minute
		}
	}

	// Create the auto-updater
	autoUpdater, err = update.NewAutoUpdater(version, enabled, interval)
	if err != nil {
		// Silently ignore errors - auto-update is optional
		return
	}
}
