package repl

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// REPL represents a Read-Eval-Print Loop session
type REPL struct {
	rootCmd *cobra.Command
	reader  *bufio.Reader
	writer  io.Writer
	history []string
	prompt  string
	session *Session
}

// New creates a new REPL instance
func New(rootCmd *cobra.Command) *REPL {
	return &REPL{
		rootCmd: rootCmd,
		reader:  bufio.NewReader(os.Stdin),
		writer:  os.Stdout,
		history: make([]string, 0),
		prompt:  "canvas> ",
		session: NewSession(),
	}
}

// Run starts the REPL loop
func (r *REPL) Run(ctx context.Context) error {
	fmt.Fprintf(r.writer, "Canvas CLI Interactive Mode\n")
	fmt.Fprintf(r.writer, "Type 'help' for available commands, 'exit' to quit\n\n")

	for {
		// Print prompt
		fmt.Fprint(r.writer, r.prompt)

		// Read input
		input, err := r.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Fprintln(r.writer, "\nGoodbye!")
				return nil
			}
			return fmt.Errorf("failed to read input: %w", err)
		}

		// Trim whitespace
		input = strings.TrimSpace(input)

		// Skip empty lines
		if input == "" {
			continue
		}

		// Add to history
		r.history = append(r.history, input)

		// Check for exit
		if input == "exit" || input == "quit" {
			fmt.Fprintln(r.writer, "Goodbye!")
			return nil
		}

		// Check for special REPL commands
		if r.handleReplCommand(input) {
			continue
		}

		// Execute command
		if err := r.executeCommand(ctx, input); err != nil {
			fmt.Fprintf(r.writer, "Error: %v\n", err)
		}
	}
}

// handleReplCommand handles REPL-specific commands
func (r *REPL) handleReplCommand(input string) bool {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return false
	}

	cmd := parts[0]

	switch cmd {
	case "history":
		r.showHistory()
		return true
	case "clear":
		r.clearScreen()
		return true
	case "session":
		if len(parts) > 1 {
			r.handleSessionCommand(parts[1:])
		} else {
			r.showSession()
		}
		return true
	default:
		return false
	}
}

// executeCommand executes a CLI command
func (r *REPL) executeCommand(ctx context.Context, input string) error {
	// Parse the input into arguments
	args := strings.Fields(input)
	if len(args) == 0 {
		return nil
	}

	// Create a new command instance
	cmd := r.rootCmd

	// Set the arguments
	cmd.SetArgs(args)

	// Execute the command
	return cmd.ExecuteContext(ctx)
}

// showHistory displays the command history
func (r *REPL) showHistory() {
	if len(r.history) == 0 {
		fmt.Fprintln(r.writer, "No command history")
		return
	}

	fmt.Fprintln(r.writer, "Command History:")
	for i, cmd := range r.history {
		fmt.Fprintf(r.writer, "%4d  %s\n", i+1, cmd)
	}
}

// clearScreen clears the terminal screen
func (r *REPL) clearScreen() {
	fmt.Fprint(r.writer, "\033[2J\033[H")
}

// showSession displays the current session state
func (r *REPL) showSession() {
	fmt.Fprintln(r.writer, "Current Session:")
	fmt.Fprintf(r.writer, "  Course ID: %d\n", r.session.CourseID)
	fmt.Fprintf(r.writer, "  User ID: %d\n", r.session.UserID)
	fmt.Fprintf(r.writer, "  Assignment ID: %d\n", r.session.AssignmentID)

	if len(r.session.Variables) > 0 {
		fmt.Fprintln(r.writer, "  Variables:")
		for k, v := range r.session.Variables {
			fmt.Fprintf(r.writer, "    %s = %v\n", k, v)
		}
	}
}

// handleSessionCommand handles session-related commands
func (r *REPL) handleSessionCommand(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(r.writer, "Usage: session set <key> <value>")
		fmt.Fprintln(r.writer, "       session get <key>")
		fmt.Fprintln(r.writer, "       session clear")
		return
	}

	action := args[0]

	switch action {
	case "set":
		if len(args) < 3 {
			fmt.Fprintln(r.writer, "Usage: session set <key> <value>")
			return
		}
		key := args[1]
		value := strings.Join(args[2:], " ")
		r.session.Set(key, value)
		fmt.Fprintf(r.writer, "Set %s = %s\n", key, value)

	case "get":
		key := args[1]
		value, exists := r.session.Get(key)
		if exists {
			fmt.Fprintf(r.writer, "%s = %v\n", key, value)
		} else {
			fmt.Fprintf(r.writer, "Variable '%s' not found\n", key)
		}

	case "clear":
		r.session.Clear()
		fmt.Fprintln(r.writer, "Session cleared")

	default:
		fmt.Fprintf(r.writer, "Unknown session command: %s\n", action)
	}
}

// GetHistory returns the command history
func (r *REPL) GetHistory() []string {
	return r.history
}

// GetSession returns the current session
func (r *REPL) GetSession() *Session {
	return r.session
}
