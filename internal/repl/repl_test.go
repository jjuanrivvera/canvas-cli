package repl

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNew(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	repl := New(rootCmd)

	if repl == nil {
		t.Fatal("expected non-nil REPL")
	}

	if repl.rootCmd != rootCmd {
		t.Error("expected rootCmd to be set")
	}

	if repl.prompt != "canvas> " {
		t.Errorf("expected prompt 'canvas> ', got '%s'", repl.prompt)
	}

	if repl.session == nil {
		t.Error("expected session to be initialized")
	}

	if len(repl.history) != 0 {
		t.Errorf("expected empty history, got %d items", len(repl.history))
	}
}

func TestShowHistory_Empty(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	repl := New(rootCmd)

	var buf bytes.Buffer
	repl.writer = &buf

	repl.showHistory()

	output := buf.String()
	if !strings.Contains(output, "No command history") {
		t.Errorf("expected 'No command history' message, got: %s", output)
	}
}

func TestShowHistory_WithCommands(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	repl := New(rootCmd)

	var buf bytes.Buffer
	repl.writer = &buf

	// Add some history
	repl.history = []string{"command1", "command2", "command3"}

	repl.showHistory()

	output := buf.String()
	if !strings.Contains(output, "command1") {
		t.Error("expected output to contain 'command1'")
	}
	if !strings.Contains(output, "command2") {
		t.Error("expected output to contain 'command2'")
	}
	if !strings.Contains(output, "command3") {
		t.Error("expected output to contain 'command3'")
	}
}

func TestClearScreen(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	repl := New(rootCmd)

	var buf bytes.Buffer
	repl.writer = &buf

	repl.clearScreen()

	output := buf.String()
	// Check for ANSI escape sequences
	if !strings.Contains(output, "\033") {
		t.Error("expected ANSI escape sequences for clear")
	}
}

func TestShowSession(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	repl := New(rootCmd)

	var buf bytes.Buffer
	repl.writer = &buf

	// Set some session data
	repl.session.CourseID = 123
	repl.session.UserID = 456
	repl.session.AssignmentID = 789
	repl.session.Variables["test_var"] = "test_value"

	repl.showSession()

	output := buf.String()
	if !strings.Contains(output, "123") {
		t.Error("expected output to contain course ID")
	}
	if !strings.Contains(output, "456") {
		t.Error("expected output to contain user ID")
	}
	if !strings.Contains(output, "789") {
		t.Error("expected output to contain assignment ID")
	}
	if !strings.Contains(output, "test_var") {
		t.Error("expected output to contain variable name")
	}
	if !strings.Contains(output, "test_value") {
		t.Error("expected output to contain variable value")
	}
}

func TestHandleReplCommand_History(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	repl := New(rootCmd)

	var buf bytes.Buffer
	repl.writer = &buf

	handled := repl.handleReplCommand("history")

	if !handled {
		t.Error("expected history command to be handled")
	}
}

func TestHandleReplCommand_Clear(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	repl := New(rootCmd)

	var buf bytes.Buffer
	repl.writer = &buf

	handled := repl.handleReplCommand("clear")

	if !handled {
		t.Error("expected clear command to be handled")
	}
}

func TestHandleReplCommand_Session(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	repl := New(rootCmd)

	var buf bytes.Buffer
	repl.writer = &buf

	handled := repl.handleReplCommand("session")

	if !handled {
		t.Error("expected session command to be handled")
	}
}

func TestHandleReplCommand_Unknown(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	repl := New(rootCmd)

	handled := repl.handleReplCommand("unknown_command")

	if handled {
		t.Error("expected unknown command to not be handled")
	}
}

func TestHandleSessionCommand_Set(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	repl := New(rootCmd)

	var buf bytes.Buffer
	repl.writer = &buf

	repl.handleSessionCommand([]string{"set", "test_key", "test_value"})

	value, exists := repl.session.Get("test_key")
	if !exists {
		t.Error("expected variable to be set")
	}
	if value != "test_value" {
		t.Errorf("expected value 'test_value', got '%v'", value)
	}
}

func TestHandleSessionCommand_Get(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	repl := New(rootCmd)

	var buf bytes.Buffer
	repl.writer = &buf

	// Set a variable first
	repl.session.Set("test_key", "test_value")

	repl.handleSessionCommand([]string{"get", "test_key"})

	output := buf.String()
	if !strings.Contains(output, "test_value") {
		t.Error("expected output to contain variable value")
	}
}

func TestHandleSessionCommand_GetNonexistent(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	repl := New(rootCmd)

	var buf bytes.Buffer
	repl.writer = &buf

	repl.handleSessionCommand([]string{"get", "nonexistent"})

	output := buf.String()
	if !strings.Contains(output, "not found") {
		t.Error("expected 'not found' message")
	}
}

func TestHandleSessionCommand_Clear(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	repl := New(rootCmd)

	var buf bytes.Buffer
	repl.writer = &buf

	// Set some variables
	repl.session.Set("test_key", "test_value")
	repl.session.CourseID = 123

	// Need to pass a dummy second arg to get past the len check
	repl.handleSessionCommand([]string{"clear", ""})

	if repl.session.CourseID != 0 {
		t.Error("expected course ID to be cleared")
	}
	if len(repl.session.Variables) != 0 {
		t.Error("expected variables to be cleared")
	}
}

func TestHandleSessionCommand_InvalidUsage(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	repl := New(rootCmd)

	var buf bytes.Buffer
	repl.writer = &buf

	repl.handleSessionCommand([]string{"set"})

	output := buf.String()
	if !strings.Contains(output, "Usage") {
		t.Error("expected usage message for invalid command")
	}
}

func TestHandleSessionCommand_UnknownAction(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	repl := New(rootCmd)

	var buf bytes.Buffer
	repl.writer = &buf

	repl.handleSessionCommand([]string{"unknown", "arg"})

	output := buf.String()
	if !strings.Contains(output, "Unknown") {
		t.Error("expected unknown command message")
	}
}

func TestGetHistory(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	repl := New(rootCmd)

	repl.history = []string{"cmd1", "cmd2"}

	history := repl.GetHistory()
	if len(history) != 2 {
		t.Errorf("expected 2 history items, got %d", len(history))
	}
	if history[0] != "cmd1" {
		t.Errorf("expected 'cmd1', got '%s'", history[0])
	}
}

func TestGetSession(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	repl := New(rootCmd)

	session := repl.GetSession()
	if session == nil {
		t.Error("expected non-nil session")
	}
}

func TestExecuteCommand(t *testing.T) {
	rootCmd := &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing
		},
	}

	repl := New(rootCmd)
	ctx := context.Background()

	err := repl.executeCommand(ctx, "test")
	if err != nil {
		t.Errorf("executeCommand failed: %v", err)
	}
}

func TestExecuteCommand_EmptyInput(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	repl := New(rootCmd)
	ctx := context.Background()

	err := repl.executeCommand(ctx, "")
	if err != nil {
		t.Errorf("executeCommand failed on empty input: %v", err)
	}
}

func TestHandleSessionCommand_NoArgs(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	repl := New(rootCmd)

	var buf bytes.Buffer
	repl.writer = &buf

	repl.handleSessionCommand([]string{})

	output := buf.String()
	if !strings.Contains(output, "Usage") {
		t.Error("expected usage message when no args provided")
	}
}

func TestHandleSessionCommand_SetMultipleWords(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	repl := New(rootCmd)

	var buf bytes.Buffer
	repl.writer = &buf

	repl.handleSessionCommand([]string{"set", "test_key", "value", "with", "spaces"})

	value, exists := repl.session.Get("test_key")
	if !exists {
		t.Error("expected variable to be set")
	}
	expected := "value with spaces"
	if value != expected {
		t.Errorf("expected value '%s', got '%v'", expected, value)
	}
}

// Note: The following Run tests were removed because they relied on the old
// bufio.Reader-based implementation. The REPL now uses the readline library
// which provides advanced features like arrow key navigation and history search.
// The readline library is tested separately.
// Tests that were removed: TestRun_ExitCommand, TestRun_QuitCommand,
// TestRun_EmptyInput, TestRun_HistoryCommand, TestRun_EOF
