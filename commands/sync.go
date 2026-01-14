package commands

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"github.com/jjuanrivvera/canvas-cli/internal/auth"
	"github.com/jjuanrivvera/canvas-cli/internal/batch"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize resources between Canvas instances",
	Long: `Synchronize courses, assignments, and other resources between different Canvas instances.

This is useful for:
- Migrating content between instances
- Backing up course content
- Copying course structures
- Synchronizing development and production environments`,
}

var syncAssignmentsCmd = &cobra.Command{
	Use:   "assignments <source-instance> <source-course-id> <target-instance> <target-course-id>",
	Short: "Sync assignments between instances",
	Long: `Synchronize all assignments from a source course to a target course.

The source and target can be on different Canvas instances.

Examples:
  # Sync assignments from production to staging
  canvas sync assignments prod 12345 staging 67890

  # Sync with interactive conflict resolution
  canvas sync assignments prod 12345 staging 67890 --interactive`,
	Args: ExactArgsWithUsage(4, "source-instance", "source-course-id", "target-instance", "target-course-id"),
	RunE: runSyncAssignments,
}

var syncCourseCmd = &cobra.Command{
	Use:   "course <source-instance> <source-course-id> <target-instance> <target-course-id>",
	Short: "Sync entire course between instances",
	Long: `Synchronize an entire course structure including assignments, files, and settings.

The source and target can be on different Canvas instances.

Examples:
  # Sync course from production to staging
  canvas sync course prod 12345 staging 67890

  # Sync with interactive conflict resolution
  canvas sync course prod 12345 staging 67890 --interactive`,
	Args: ExactArgsWithUsage(4, "source-instance", "source-course-id", "target-instance", "target-course-id"),
	RunE: runSyncCourse,
}

var (
	syncInteractive bool
)

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.AddCommand(syncAssignmentsCmd)
	syncCmd.AddCommand(syncCourseCmd)

	syncCmd.PersistentFlags().BoolVarP(&syncInteractive, "interactive", "i", false, "Enable interactive conflict resolution")
}

func runSyncAssignments(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	sourceInstance := args[0]
	sourceCourseID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid source course ID: %w", err)
	}

	targetInstance := args[2]
	targetCourseID, err := strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid target course ID: %w", err)
	}

	// Create API clients for both instances
	sourceClient, err := getAPIClientForInstance(sourceInstance)
	if err != nil {
		return fmt.Errorf("failed to create source client: %w", err)
	}

	targetClient, err := getAPIClientForInstance(targetInstance)
	if err != nil {
		return fmt.Errorf("failed to create target client: %w", err)
	}

	// Create sync operation
	syncOp := batch.NewSyncOperation(sourceClient, targetClient, syncInteractive)

	fmt.Printf("🔄 Syncing assignments from %s (course %d) to %s (course %d)\n\n",
		sourceInstance, sourceCourseID, targetInstance, targetCourseID)

	// Perform sync
	result, err := syncOp.SyncAssignments(ctx, sourceCourseID, targetCourseID)
	if err != nil {
		fmt.Printf("\n❌ Sync failed: %v\n", err)
		return err
	}

	// Display results
	fmt.Printf("\n✅ Sync complete!\n")
	fmt.Printf("Total assignments: %d\n", result.TotalItems)
	fmt.Printf("Synced: %d\n", result.SyncedItems)
	fmt.Printf("Skipped: %d\n", result.SkippedItems)
	fmt.Printf("Failed: %d\n", result.FailedItems)

	if len(result.Errors) > 0 {
		fmt.Println("\n⚠️  Errors:")
		for _, err := range result.Errors {
			fmt.Printf("  - %v\n", err)
		}
	}

	return nil
}

func runSyncCourse(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	sourceInstance := args[0]
	sourceCourseID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid source course ID: %w", err)
	}

	targetInstance := args[2]
	targetCourseID, err := strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid target course ID: %w", err)
	}

	// Create API clients for both instances
	sourceClient, err := getAPIClientForInstance(sourceInstance)
	if err != nil {
		return fmt.Errorf("failed to create source client: %w", err)
	}

	targetClient, err := getAPIClientForInstance(targetInstance)
	if err != nil {
		return fmt.Errorf("failed to create target client: %w", err)
	}

	// Create sync operation
	syncOp := batch.NewSyncOperation(sourceClient, targetClient, syncInteractive)

	fmt.Printf("🔄 Syncing course from %s (course %d) to %s (course %d)\n\n",
		sourceInstance, sourceCourseID, targetInstance, targetCourseID)

	// Perform sync
	err = syncOp.CopyCourse(ctx, sourceCourseID, targetCourseID)
	if err != nil {
		fmt.Printf("\n❌ Sync failed: %v\n", err)
		return err
	}

	fmt.Printf("\n✅ Course sync complete!\n")

	return nil
}

// getAPIClientForInstance creates an API client for a specific instance name
func getAPIClientForInstance(instanceName string) (*api.Client, error) {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Get instance by name
	instance, err := cfg.GetInstance(instanceName)
	if err != nil {
		return nil, fmt.Errorf("instance not found: %w", err)
	}

	// Get config directory
	configDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	configDir = configDir + "/.canvas-cli"

	// Load token
	tokenStore := auth.NewFallbackTokenStore(configDir)
	token, err := tokenStore.Load(instance.Name)
	if err != nil {
		return nil, fmt.Errorf("not authenticated with %s. Run 'canvas auth login' first", instance.Name)
	}

	// Create auto-refreshing token source if we have OAuth credentials
	var clientConfig api.ClientConfig
	if instance.ClientID != "" && instance.ClientSecret != "" {
		// Create oauth2 config for token refresh
		oauth2Config := auth.CreateOAuth2ConfigForInstance(instance.URL, instance.ClientID, instance.ClientSecret)
		tokenSource := auth.NewAutoRefreshTokenSource(oauth2Config, tokenStore, instance.Name, token)

		clientConfig = api.ClientConfig{
			BaseURL:        instance.URL,
			TokenSource:    tokenSource,
			RequestsPerSec: cfg.Settings.RequestsPerSecond,
		}
	} else {
		// Fall back to static token (no auto-refresh)
		clientConfig = api.ClientConfig{
			BaseURL:        instance.URL,
			Token:          token.AccessToken,
			RequestsPerSec: cfg.Settings.RequestsPerSecond,
		}
	}

	// Create API client
	client, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	return client, nil
}
