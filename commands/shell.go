package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/repl"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Start interactive shell (REPL mode)",
	Long: `Start an interactive Read-Eval-Print Loop (REPL) shell for Canvas CLI.

The shell provides an interactive environment where you can execute Canvas commands
without typing the 'canvas' prefix. It maintains session state and command
history throughout the session.

Special shell commands:
  history           - Show command history
  clear             - Clear the screen
  session           - Show current session state
  session set <k> <v> - Set a session variable
  session get <k>   - Get a session variable
  session clear     - Clear session state
  exit/quit         - Exit the shell

Examples:
  # Start the interactive shell
  canvas shell

  # In the shell:
  canvas> courses list
  canvas> session set course_id 12345
  canvas> assignments list --course-id 12345
  canvas> history
  canvas> exit

Note: This command is an alias for 'canvas repl'. Both commands start the same
interactive shell environment.`,
	RunE: runShell,
}

func init() {
	rootCmd.AddCommand(shellCmd)
}

func runShell(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Create REPL instance
	r := repl.New(rootCmd)

	// Start the REPL
	if err := r.Run(ctx); err != nil {
		return fmt.Errorf("shell error: %w", err)
	}

	return nil
}
