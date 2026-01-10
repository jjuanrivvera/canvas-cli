package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/repl"
)

var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Start interactive REPL mode",
	Long: `Start an interactive Read-Eval-Print Loop (REPL) for Canvas CLI.

The REPL provides an interactive shell where you can execute Canvas commands
without typing the 'canvas' prefix. It maintains session state and command
history throughout the session.

Special REPL commands:
  history           - Show command history
  clear             - Clear the screen
  session           - Show current session state
  session set <k> <v> - Set a session variable
  session get <k>   - Get a session variable
  session clear     - Clear session state
  exit/quit         - Exit the REPL

Examples:
  # Start the REPL
  canvas repl

  # In the REPL:
  canvas> courses list
  canvas> session set course_id 12345
  canvas> assignments list --course-id 12345
  canvas> history
  canvas> exit`,
	RunE: runRepl,
}

func init() {
	rootCmd.AddCommand(replCmd)
}

func runRepl(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Create REPL instance
	r := repl.New(rootCmd)

	// Start the REPL
	if err := r.Run(ctx); err != nil {
		return fmt.Errorf("REPL error: %w", err)
	}

	return nil
}
