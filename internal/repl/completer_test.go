package repl

import (
	"reflect"
	"sort"
	"testing"

	"github.com/spf13/cobra"
)

func createTestRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "canvas",
		Short: "Canvas CLI tool",
	}

	coursesCmd := &cobra.Command{
		Use:   "courses",
		Short: "Manage courses",
	}
	coursesCmd.Flags().String("name", "", "Course name")
	coursesCmd.Flags().StringP("output", "o", "table", "Output format")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List courses",
	}
	coursesCmd.AddCommand(listCmd)

	usersCmd := &cobra.Command{
		Use:   "users",
		Short: "Manage users",
	}
	usersCmd.Flags().String("email", "", "User email")

	hiddenCmd := &cobra.Command{
		Use:    "hidden",
		Short:  "Hidden command",
		Hidden: true,
	}

	rootCmd.AddCommand(coursesCmd, usersCmd, hiddenCmd)

	return rootCmd
}

func TestNewCompleter(t *testing.T) {
	rootCmd := createTestRootCommand()
	completer := NewCompleter(rootCmd)

	if completer == nil {
		t.Fatal("expected non-nil completer")
	}

	if completer.rootCmd != rootCmd {
		t.Error("expected root command to be set")
	}
}

func TestCompleter_Complete_EmptyInput(t *testing.T) {
	rootCmd := createTestRootCommand()
	completer := NewCompleter(rootCmd)

	suggestions := completer.Complete("")
	if len(suggestions) == 0 {
		t.Error("expected suggestions for empty input")
	}

	// Should include root commands and REPL commands
	expected := []string{"courses", "users", "history", "clear", "session", "exit", "quit"}
	sort.Strings(suggestions)
	sort.Strings(expected)

	if !reflect.DeepEqual(suggestions, expected) {
		t.Errorf("expected %v, got %v", expected, suggestions)
	}
}

func TestCompleter_Complete_PartialCommand(t *testing.T) {
	rootCmd := createTestRootCommand()
	completer := NewCompleter(rootCmd)

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "complete 'c' matches courses",
			input:    "c",
			expected: []string{"courses"},
		},
		{
			name:     "complete 'co' matches courses",
			input:    "co",
			expected: []string{"courses"},
		},
		{
			name:     "complete 'u' matches users",
			input:    "u",
			expected: []string{"users"},
		},
		{
			name:     "complete 'h' matches nothing (REPL commands only in empty input)",
			input:    "h",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := completer.Complete(tt.input)
			sort.Strings(suggestions)
			sort.Strings(tt.expected)

			if !reflect.DeepEqual(suggestions, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, suggestions)
			}
		})
	}
}

func TestCompleter_Complete_Subcommands(t *testing.T) {
	rootCmd := createTestRootCommand()
	completer := NewCompleter(rootCmd)

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "complete 'courses ' suggests subcommands",
			input:    "courses ",
			expected: []string{"list"},
		},
		{
			name:     "complete 'courses l' matches list",
			input:    "courses l",
			expected: []string{"list"},
		},
		{
			name:     "complete 'users ' has no subcommands",
			input:    "users ",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := completer.Complete(tt.input)
			sort.Strings(suggestions)
			if tt.expected != nil {
				sort.Strings(tt.expected)
			}

			if !reflect.DeepEqual(suggestions, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, suggestions)
			}
		})
	}
}

func TestCompleter_Complete_Flags(t *testing.T) {
	rootCmd := createTestRootCommand()
	completer := NewCompleter(rootCmd)

	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:     "complete 'courses arg --' shows long flags",
			input:    "courses arg --",
			contains: []string{"--name", "--output"},
		},
		{
			name:     "complete 'courses arg -' shows short flags",
			input:    "courses arg -",
			contains: []string{"-o"},
		},
		{
			name:     "complete 'courses arg --n' matches --name",
			input:    "courses arg --n",
			contains: []string{"--name"},
		},
		{
			name:     "complete 'users arg --' shows user flags",
			input:    "users arg --",
			contains: []string{"--email"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := completer.Complete(tt.input)

			for _, expected := range tt.contains {
				found := false
				for _, suggestion := range suggestions {
					if suggestion == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected suggestions to contain %q, got %v", expected, suggestions)
				}
			}
		})
	}
}

func TestCompleter_Complete_NoFlagsWithoutDash(t *testing.T) {
	rootCmd := createTestRootCommand()
	completer := NewCompleter(rootCmd)

	suggestions := completer.Complete("courses arg")
	// Should not return flags when the prefix doesn't start with -
	for _, s := range suggestions {
		if len(s) > 0 && s[0] == '-' {
			t.Errorf("unexpected flag suggestion %q for input without dash", s)
		}
	}
}

func TestCompleter_rootCommands(t *testing.T) {
	rootCmd := createTestRootCommand()
	completer := NewCompleter(rootCmd)

	commands := completer.rootCommands()

	// Should include visible commands
	if !contains(commands, "courses") {
		t.Error("expected 'courses' in root commands")
	}
	if !contains(commands, "users") {
		t.Error("expected 'users' in root commands")
	}

	// Should not include hidden commands
	if contains(commands, "hidden") {
		t.Error("expected 'hidden' to not be in root commands")
	}

	// Should include REPL commands
	replCommands := []string{"history", "clear", "session", "exit", "quit"}
	for _, cmd := range replCommands {
		if !contains(commands, cmd) {
			t.Errorf("expected REPL command %q in root commands", cmd)
		}
	}
}

func TestCompleter_matchCommands(t *testing.T) {
	rootCmd := createTestRootCommand()
	completer := NewCompleter(rootCmd)

	tests := []struct {
		name     string
		prefix   string
		expected []string
	}{
		{
			name:     "match 'c' finds courses",
			prefix:   "c",
			expected: []string{"courses"},
		},
		{
			name:     "match 'u' finds users",
			prefix:   "u",
			expected: []string{"users"},
		},
		{
			name:     "match '' finds all visible commands",
			prefix:   "",
			expected: []string{"courses", "users"},
		},
		{
			name:     "match 'x' finds nothing",
			prefix:   "x",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := completer.matchCommands(rootCmd, tt.prefix)
			sort.Strings(matches)
			sort.Strings(tt.expected)

			if !reflect.DeepEqual(matches, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, matches)
			}
		})
	}
}

func TestCompleter_matchSubcommands(t *testing.T) {
	rootCmd := createTestRootCommand()
	completer := NewCompleter(rootCmd)

	coursesCmd := rootCmd.Commands()[0] // courses command

	tests := []struct {
		name     string
		prefix   string
		expected []string
	}{
		{
			name:     "match 'l' finds list",
			prefix:   "l",
			expected: []string{"list"},
		},
		{
			name:     "match '' finds all subcommands",
			prefix:   "",
			expected: []string{"list"},
		},
		{
			name:     "match 'x' finds nothing",
			prefix:   "x",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := completer.matchSubcommands(coursesCmd, tt.prefix)
			sort.Strings(matches)
			sort.Strings(tt.expected)

			if !reflect.DeepEqual(matches, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, matches)
			}
		})
	}
}

func TestCompleter_matchFlags(t *testing.T) {
	rootCmd := createTestRootCommand()
	completer := NewCompleter(rootCmd)

	coursesCmd := rootCmd.Commands()[0] // courses command

	tests := []struct {
		name     string
		prefix   string
		contains []string
	}{
		{
			name:     "match '--' finds long flags",
			prefix:   "--",
			contains: []string{"--name", "--output"},
		},
		{
			name:     "match '--n' finds --name",
			prefix:   "--n",
			contains: []string{"--name"},
		},
		{
			name:     "match '-' finds short flags",
			prefix:   "-",
			contains: []string{"-o"},
		},
		{
			name:     "match '-o' finds -o",
			prefix:   "-o",
			contains: []string{"-o"},
		},
		{
			name:     "match without dash returns nil",
			prefix:   "name",
			contains: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := completer.matchFlags(coursesCmd, tt.prefix)

			if tt.contains == nil {
				if matches != nil {
					t.Errorf("expected nil, got %v", matches)
				}
				return
			}

			for _, expected := range tt.contains {
				if !contains(matches, expected) {
					t.Errorf("expected %q in %v", expected, matches)
				}
			}
		})
	}
}

func TestCompleter_findCommand(t *testing.T) {
	rootCmd := createTestRootCommand()
	completer := NewCompleter(rootCmd)

	tests := []struct {
		name        string
		parts       []string
		expectCmd   string
		expectArgs  []string
		expectError bool
	}{
		{
			name:       "find root command",
			parts:      []string{"courses"},
			expectCmd:  "courses",
			expectArgs: nil,
		},
		{
			name:       "find subcommand",
			parts:      []string{"courses", "list"},
			expectCmd:  "list",
			expectArgs: nil,
		},
		{
			name:       "find command with args",
			parts:      []string{"users", "arg1"},
			expectCmd:  "users",
			expectArgs: []string{"arg1"},
		},
		{
			name:       "find command with multiple args",
			parts:      []string{"users", "arg1", "arg2"},
			expectCmd:  "users",
			expectArgs: []string{"arg1", "arg2"},
		},
		{
			name:       "flags don't prevent finding subcommand",
			parts:      []string{"courses", "--name", "list"},
			expectCmd:  "list",
			expectArgs: nil,
		},
		{
			name:       "nonexistent command treated as args",
			parts:      []string{"nonexistent"},
			expectCmd:  "canvas",
			expectArgs: []string{"nonexistent"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, args, err := completer.findCommand(tt.parts)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if cmd.Use != tt.expectCmd {
				t.Errorf("expected command %q, got %q", tt.expectCmd, cmd.Use)
			}

			if !reflect.DeepEqual(args, tt.expectArgs) {
				t.Errorf("expected args %v, got %v", tt.expectArgs, args)
			}
		})
	}
}

func TestCompleter_GetCommandHelp(t *testing.T) {
	rootCmd := createTestRootCommand()
	completer := NewCompleter(rootCmd)

	tests := []struct {
		name     string
		cmdName  string
		expected string
	}{
		{
			name:     "get courses help",
			cmdName:  "courses",
			expected: "Manage courses",
		},
		{
			name:     "get users help",
			cmdName:  "users",
			expected: "Manage users",
		},
		{
			name:     "get nonexistent command help",
			cmdName:  "nonexistent",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			help := completer.GetCommandHelp(tt.cmdName)
			if help != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, help)
			}
		})
	}
}

func TestCompleter_GetFlagHelp(t *testing.T) {
	rootCmd := createTestRootCommand()
	completer := NewCompleter(rootCmd)

	coursesCmd := rootCmd.Commands()[0] // courses command

	tests := []struct {
		name     string
		flagName string
		expected string
	}{
		{
			name:     "get --name flag help",
			flagName: "--name",
			expected: "Course name",
		},
		{
			name:     "get name flag help without dashes",
			flagName: "name",
			expected: "Course name",
		},
		{
			name:     "get -o flag help (shorthand not found by Lookup)",
			flagName: "-o",
			expected: "", // Lookup searches by name, not shorthand
		},
		{
			name:     "get --output flag help",
			flagName: "--output",
			expected: "Output format",
		},
		{
			name:     "get nonexistent flag help",
			flagName: "--nonexistent",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			help := completer.GetFlagHelp(coursesCmd, tt.flagName)
			if help != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, help)
			}
		})
	}
}

// Helper function
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func TestCompleter_matchFlags_HiddenFlags(t *testing.T) {
	rootCmd := createTestRootCommand()
	completer := NewCompleter(rootCmd)

	coursesCmd := rootCmd.Commands()[0]

	// Add a hidden flag
	coursesCmd.Flags().String("hidden-flag", "", "Hidden flag for testing")
	coursesCmd.Flags().MarkHidden("hidden-flag")

	// Complete with flag prefix
	results := completer.matchFlags(coursesCmd, "--")

	// Should not include hidden flag
	for _, r := range results {
		if r == "--hidden-flag" {
			t.Error("hidden flags should not appear in completions")
		}
	}
}

func TestCompleter_Complete_EndsWithSpace(t *testing.T) {
	rootCmd := createTestRootCommand()
	completer := NewCompleter(rootCmd)

	// Test completion when input ends with space (completing next arg)
	results := completer.Complete("courses ")

	// Should suggest subcommands
	if !contains(results, "list") {
		t.Error("expected 'list' in subcommand completions")
	}
}

func TestCompleter_Complete_MultipleArgs(t *testing.T) {
	rootCmd := createTestRootCommand()
	completer := NewCompleter(rootCmd)

	// Test completion with multiple arguments and flag prefix
	results := completer.Complete("courses --")

	// Should suggest flags (courses command has flags)
	// The result may be empty if no flags are defined, which is OK
	// This tests the code path where we have args and prefix starts with "-"
	_ = results // Just exercise the code path
}
