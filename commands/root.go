package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	instanceURL  string
	outputFormat string
	verbose      bool
	noCache      bool  // Disable caching for API requests
	asUserID     int64 // Masquerading: act as another user
	version      string
	commit       string
	buildDate    string
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
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute(v, c, bd string) error {
	version = v
	commit = c
	buildDate = bd
	return rootCmd.Execute()
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
		// Find home directory
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".canvas-cli"
		configDir := home + "/.canvas-cli"
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
