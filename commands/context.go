package commands

import (
	"fmt"

	"github.com/spf13/cobra"

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
		RunE: runContextSet,
	}

	return cmd
}

func runContextSet(cmd *cobra.Command, args []string) error {
	contextType := args[0]
	idStr := args[1]

	var id int64
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		return fmt.Errorf("invalid ID %q: must be a number", idStr)
	}

	if id <= 0 {
		return fmt.Errorf("invalid ID %d: must be a positive number", id)
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	ctx := cfg.GetContext()

	switch contextType {
	case "course", "course_id", "course-id":
		ctx.CourseID = id
		contextType = "course"
	case "assignment", "assignment_id", "assignment-id":
		ctx.AssignmentID = id
		contextType = "assignment"
	case "user", "user_id", "user-id":
		ctx.UserID = id
		contextType = "user"
	case "account", "account_id", "account-id":
		ctx.AccountID = id
		contextType = "account"
	default:
		return fmt.Errorf("unknown context type %q. Valid types: course, assignment, user, account", contextType)
	}

	if err := cfg.SetContext(ctx); err != nil {
		return fmt.Errorf("failed to save context: %w", err)
	}

	fmt.Printf("Context %s set to %d\n", contextType, id)
	return nil
}

func newContextShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show current context",
		Long:  `Display the current working context values.`,
		Args:  cobra.NoArgs,
		RunE:  runContextShow,
	}

	return cmd
}

func runContextShow(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	ctx := cfg.GetContext()

	// Check if any context is set
	if ctx.CourseID == 0 && ctx.AssignmentID == 0 && ctx.UserID == 0 && ctx.AccountID == 0 {
		fmt.Println("No context set.")
		fmt.Println("\nSet context with: canvas context set <type> <id>")
		fmt.Println("Valid types: course, assignment, user, account")
		return nil
	}

	fmt.Println("Current context:")
	if ctx.CourseID > 0 {
		fmt.Printf("  course_id:     %d\n", ctx.CourseID)
	}
	if ctx.AssignmentID > 0 {
		fmt.Printf("  assignment_id: %d\n", ctx.AssignmentID)
	}
	if ctx.UserID > 0 {
		fmt.Printf("  user_id:       %d\n", ctx.UserID)
	}
	if ctx.AccountID > 0 {
		fmt.Printf("  account_id:    %d\n", ctx.AccountID)
	}

	return nil
}

func newContextClearCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear [type]",
		Short: "Clear context values",
		Long: `Clear all context values or a specific context type.

Examples:
  canvas context clear           # Clear all context
  canvas context clear course    # Clear only course context`,
		Args: cobra.MaximumNArgs(1),
		RunE: runContextClear,
	}

	return cmd
}

func runContextClear(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(args) == 0 {
		// Clear all context
		if err := cfg.ClearContext(); err != nil {
			return fmt.Errorf("failed to clear context: %w", err)
		}
		fmt.Println("Context cleared.")
		return nil
	}

	// Clear specific context type
	contextType := args[0]
	ctx := cfg.GetContext()

	switch contextType {
	case "course", "course_id", "course-id":
		ctx.CourseID = 0
		contextType = "course"
	case "assignment", "assignment_id", "assignment-id":
		ctx.AssignmentID = 0
		contextType = "assignment"
	case "user", "user_id", "user-id":
		ctx.UserID = 0
		contextType = "user"
	case "account", "account_id", "account-id":
		ctx.AccountID = 0
		contextType = "account"
	default:
		return fmt.Errorf("unknown context type %q. Valid types: course, assignment, user, account", contextType)
	}

	if err := cfg.SetContext(ctx); err != nil {
		return fmt.Errorf("failed to save context: %w", err)
	}

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
