package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

func init() {
	rootCmd.AddCommand(newContextCmd())
}

func newContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Manage working context (course, assignment, user IDs)",
		Long: `Manage the current working context for Canvas CLI.

Context allows you to set default values for course_id, assignment_id, user_id,
and account_id that will be used automatically when those flags are not provided.

Examples:
  # Set the current course
  canvas context set course 123

  # Set multiple context values
  canvas context set course 123
  canvas context set assignment 456

  # Now commands automatically use context
  canvas assignments list  # uses course_id 123
  canvas submissions list  # uses course_id 123 and assignment_id 456

  # Show current context
  canvas context show

  # Clear all context
  canvas context clear

  # Clear specific context value
  canvas context clear course`,
	}

	cmd.AddCommand(newContextSetCmd())
	cmd.AddCommand(newContextShowCmd())
	cmd.AddCommand(newContextClearCmd())

	return cmd
}

func newContextSetCmd() *cobra.Command {
	opts := &options.ContextSetOptions{}

	cmd := &cobra.Command{
		Use:   "set <type> <id>",
		Short: "Set a context value",
		Long: `Set a context value that will be used as default for commands.

Valid types:
  course      - Course ID (used for --course-id)
  assignment  - Assignment ID (used for --assignment-id)
  user        - User ID (used for --user-id)
  account     - Account ID (used for --account-id)

Examples:
  canvas context set course 123
  canvas context set assignment 456`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Type = args[0]
			var id int64
			if _, err := fmt.Sscanf(args[1], "%d", &id); err != nil {
				return fmt.Errorf("invalid ID %q: must be a number", args[1])
			}
			opts.ID = id
			if err := opts.Validate(); err != nil {
				return err
			}
			return runContextSet(cmd.Context(), opts)
		},
	}

	return cmd
}

func runContextSet(ctx context.Context, opts *options.ContextSetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "context.set", map[string]interface{}{
		"type": opts.Type,
		"id":   opts.ID,
	})

	cfg, err := config.Load()
	if err != nil {
		logger.LogCommandError(ctx, "context.set", err, nil)
		return fmt.Errorf("failed to load config: %w", err)
	}

	ctxVal := cfg.GetContext()
	contextType := opts.Type

	switch contextType {
	case "course", "course_id", "course-id":
		ctxVal.CourseID = opts.ID
		contextType = "course"
	case "assignment", "assignment_id", "assignment-id":
		ctxVal.AssignmentID = opts.ID
		contextType = "assignment"
	case "user", "user_id", "user-id":
		ctxVal.UserID = opts.ID
		contextType = "user"
	case "account", "account_id", "account-id":
		ctxVal.AccountID = opts.ID
		contextType = "account"
	}

	if err := cfg.SetContext(ctxVal); err != nil {
		logger.LogCommandError(ctx, "context.set", err, nil)
		return fmt.Errorf("failed to save context: %w", err)
	}

	logger.LogCommandComplete(ctx, "context.set", 1)
	fmt.Printf("Context %s set to %d\n", contextType, opts.ID)
	return nil
}

func newContextShowCmd() *cobra.Command {
	opts := &options.ContextShowOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show current context",
		Long:  `Display the current working context values.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			return runContextShow(cmd.Context(), opts)
		},
	}

	return cmd
}

func runContextShow(ctx context.Context, opts *options.ContextShowOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "context.show", nil)

	cfg, err := config.Load()
	if err != nil {
		logger.LogCommandError(ctx, "context.show", err, nil)
		return fmt.Errorf("failed to load config: %w", err)
	}

	ctxVal := cfg.GetContext()

	// Check if any context is set
	if ctxVal.CourseID == 0 && ctxVal.AssignmentID == 0 && ctxVal.UserID == 0 && ctxVal.AccountID == 0 {
		fmt.Println("No context set.")
		fmt.Println("\nSet context with: canvas context set <type> <id>")
		fmt.Println("Valid types: course, assignment, user, account")
		logger.LogCommandComplete(ctx, "context.show", 0)
		return nil
	}

	count := 0
	fmt.Println("Current context:")
	if ctxVal.CourseID > 0 {
		fmt.Printf("  course_id:     %d\n", ctxVal.CourseID)
		count++
	}
	if ctxVal.AssignmentID > 0 {
		fmt.Printf("  assignment_id: %d\n", ctxVal.AssignmentID)
		count++
	}
	if ctxVal.UserID > 0 {
		fmt.Printf("  user_id:       %d\n", ctxVal.UserID)
		count++
	}
	if ctxVal.AccountID > 0 {
		fmt.Printf("  account_id:    %d\n", ctxVal.AccountID)
		count++
	}

	logger.LogCommandComplete(ctx, "context.show", count)
	return nil
}

func newContextClearCmd() *cobra.Command {
	opts := &options.ContextClearOptions{}

	cmd := &cobra.Command{
		Use:   "clear [type]",
		Short: "Clear context values",
		Long: `Clear all context values or a specific context type.

Examples:
  canvas context clear           # Clear all context
  canvas context clear course    # Clear only course context`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Type = args[0]
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			return runContextClear(cmd.Context(), opts)
		},
	}

	return cmd
}

func runContextClear(ctx context.Context, opts *options.ContextClearOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "context.clear", map[string]interface{}{
		"type": opts.Type,
	})

	cfg, err := config.Load()
	if err != nil {
		logger.LogCommandError(ctx, "context.clear", err, nil)
		return fmt.Errorf("failed to load config: %w", err)
	}

	if opts.Type == "" {
		// Clear all context
		if err := cfg.ClearContext(); err != nil {
			logger.LogCommandError(ctx, "context.clear", err, nil)
			return fmt.Errorf("failed to clear context: %w", err)
		}
		logger.LogCommandComplete(ctx, "context.clear", 1)
		fmt.Println("Context cleared.")
		return nil
	}

	// Clear specific context type
	contextType := opts.Type
	ctxVal := cfg.GetContext()

	switch contextType {
	case "course", "course_id", "course-id":
		ctxVal.CourseID = 0
		contextType = "course"
	case "assignment", "assignment_id", "assignment-id":
		ctxVal.AssignmentID = 0
		contextType = "assignment"
	case "user", "user_id", "user-id":
		ctxVal.UserID = 0
		contextType = "user"
	case "account", "account_id", "account-id":
		ctxVal.AccountID = 0
		contextType = "account"
	}

	if err := cfg.SetContext(ctxVal); err != nil {
		logger.LogCommandError(ctx, "context.clear", err, nil)
		return fmt.Errorf("failed to save context: %w", err)
	}

	logger.LogCommandComplete(ctx, "context.clear", 1)
	fmt.Printf("Context %s cleared.\n", contextType)
	return nil
}

// GetContextCourseID returns the course ID from context if set and flag is not provided
func GetContextCourseID(flagValue int64) int64 {
	if flagValue != 0 {
		return flagValue
	}
	cfg, err := config.Load()
	if err != nil {
		return 0
	}
	return cfg.GetContext().CourseID
}

// GetContextAssignmentID returns the assignment ID from context if set and flag is not provided
func GetContextAssignmentID(flagValue int64) int64 {
	if flagValue != 0 {
		return flagValue
	}
	cfg, err := config.Load()
	if err != nil {
		return 0
	}
	return cfg.GetContext().AssignmentID
}

// GetContextUserID returns the user ID from context if set and flag is not provided
func GetContextUserID(flagValue int64) int64 {
	if flagValue != 0 {
		return flagValue
	}
	cfg, err := config.Load()
	if err != nil {
		return 0
	}
	return cfg.GetContext().UserID
}

// GetContextAccountID returns the account ID from context if set and flag is not provided
func GetContextAccountID(flagValue int64) int64 {
	if flagValue != 0 {
		return flagValue
	}
	cfg, err := config.Load()
	if err != nil {
		return 0
	}
	return cfg.GetContext().AccountID
}
