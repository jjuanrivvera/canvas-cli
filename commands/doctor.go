package commands

import (
	"fmt"
	"strings"

	"github.com/jjuanrivvera/canvas-cli/internal/diagnostics"
	"github.com/spf13/cobra"
)

var (
	doctorVerbose bool
	doctorJSON    bool
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run system diagnostics",
	Long: `Run diagnostic checks to verify Canvas CLI configuration and connectivity.

The doctor command performs the following checks:
  - Environment (OS, architecture, Go version)
  - Configuration (base URL, access token)
  - Connectivity (network connection to Canvas)
  - Authentication (API token validation)
  - API Access (Canvas API availability)
  - Disk Space (cache directory availability)
  - Permissions (file and directory permissions)

Examples:
  # Run all diagnostic checks
  canvas doctor

  # Run with verbose output
  canvas doctor --verbose

  # Output results as JSON
  canvas doctor --json`,
	RunE: runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)

	doctorCmd.Flags().BoolVarP(&doctorVerbose, "verbose", "v", false, "Show detailed output")
	doctorCmd.Flags().BoolVar(&doctorJSON, "json", false, "Output results as JSON")
}

func runDoctor(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Load config
	cfg, err := getConfig()
	if err != nil {
		// Config is optional, continue without it
		cfg = nil
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		// Client is optional, continue without it
		client = nil
	}

	// Create doctor instance
	doctor := diagnostics.New(cfg, client)

	// Run diagnostics
	fmt.Println("Running diagnostics...")
	report, err := doctor.Run(ctx)
	if err != nil {
		return fmt.Errorf("diagnostic error: %w", err)
	}

	// Print report
	if doctorJSON {
		printReportJSON(report)
	} else {
		printReportHuman(report)
	}

	// Exit with error if any checks failed
	if !report.IsHealthy() {
		return fmt.Errorf("diagnostics failed")
	}

	return nil
}

func printReportHuman(report *diagnostics.Report) {
	// Print header
	fmt.Println("Diagnostic Report")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	// Print checks
	for _, check := range report.Checks {
		printCheck(check)
	}

	// Print summary
	fmt.Println()
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("Summary: %s\n", report.Summary())
	fmt.Printf("Duration: %s\n", report.Duration)

	// Print result
	fmt.Println()
	if report.IsHealthy() {
		fmt.Println("✓ All checks passed!")
	} else {
		fmt.Printf("✗ %d check(s) failed\n", report.FailCount)
	}
}

func printCheck(check diagnostics.Check) {
	// Status icon
	icon := getStatusIcon(check.Status)

	// Print check name and status
	fmt.Printf("%s [%s] %s\n", icon, check.Status, check.Name)

	// Print description
	if doctorVerbose {
		fmt.Printf("  Description: %s\n", check.Description)
	}

	// Print message
	if check.Message != "" {
		fmt.Printf("  Message: %s\n", check.Message)
	}

	// Print error
	if check.Error != nil && doctorVerbose {
		fmt.Printf("  Error: %v\n", check.Error)
	}

	// Print duration
	if doctorVerbose {
		fmt.Printf("  Duration: %s\n", check.Duration)
	}

	fmt.Println()
}

func getStatusIcon(status diagnostics.CheckStatus) string {
	switch status {
	case diagnostics.StatusPass:
		return "✓"
	case diagnostics.StatusFail:
		return "✗"
	case diagnostics.StatusWarning:
		return "⚠"
	case diagnostics.StatusSkipped:
		return "○"
	default:
		return "?"
	}
}

func printReportJSON(report *diagnostics.Report) {
	// Create a simple JSON structure
	fmt.Println("{")
	fmt.Printf("  \"duration\": \"%s\",\n", report.Duration)
	fmt.Printf("  \"summary\": \"%s\",\n", report.Summary())
	fmt.Printf("  \"healthy\": %t,\n", report.IsHealthy())
	fmt.Printf("  \"checks\": [\n")

	for i, check := range report.Checks {
		fmt.Println("    {")
		fmt.Printf("      \"name\": \"%s\",\n", check.Name)
		fmt.Printf("      \"description\": \"%s\",\n", check.Description)
		fmt.Printf("      \"status\": \"%s\",\n", check.Status)
		fmt.Printf("      \"message\": \"%s\",\n", check.Message)
		fmt.Printf("      \"duration\": \"%s\"\n", check.Duration)

		if i < len(report.Checks)-1 {
			fmt.Println("    },")
		} else {
			fmt.Println("    }")
		}
	}

	fmt.Println("  ]")
	fmt.Println("}")
}
