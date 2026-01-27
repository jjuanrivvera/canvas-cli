package repl

import (
	"testing"

	"github.com/spf13/cobra"
)

func newTestRootCmd() *cobra.Command {
	root := &cobra.Command{Use: "canvas"}

	courses := &cobra.Command{Use: "courses", Short: "Manage courses"}
	courses.AddCommand(&cobra.Command{Use: "list", Short: "List courses"})
	courses.AddCommand(&cobra.Command{Use: "get", Short: "Get a course"})

	config := &cobra.Command{Use: "config", Short: "Manage config"}
	config.AddCommand(&cobra.Command{Use: "set", Short: "Set config"})

	root.AddCommand(courses, config)
	return root
}

func TestReadlineCompleter_Do(t *testing.T) {
	root := newTestRootCmd()
	completer := NewCompleter(root)
	rc := NewReadlineCompleter(completer)

	tests := []struct {
		name        string
		line        string
		pos         int
		wantResults bool
		wantPrefLen int
		description string
	}{
		{
			name:        "complete from empty",
			line:        "",
			pos:         0,
			wantResults: true,
			wantPrefLen: 0,
			description: "Should return top-level commands",
		},
		{
			name:        "partial command",
			line:        "co",
			pos:         2,
			wantResults: true,
			wantPrefLen: 2,
			description: "Should match courses and config",
		},
		{
			name:        "subcommand after space",
			line:        "courses ",
			pos:         8,
			wantResults: true,
			wantPrefLen: 0,
			description: "Should return subcommands of courses",
		},
		{
			name:        "partial subcommand",
			line:        "courses l",
			pos:         9,
			wantResults: true,
			wantPrefLen: 1,
			description: "Should match 'list'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, prefLen := rc.Do([]rune(tt.line), tt.pos)
			if tt.wantResults && len(results) == 0 {
				t.Errorf("expected results for %q but got none", tt.line)
			}
			if !tt.wantResults && len(results) > 0 {
				t.Errorf("expected no results for %q but got %d", tt.line, len(results))
			}
			if prefLen != tt.wantPrefLen {
				t.Errorf("prefixLen = %d, want %d", prefLen, tt.wantPrefLen)
			}
		})
	}
}

func TestReadlineCompleter_InterfaceCompliance(t *testing.T) {
	root := newTestRootCmd()
	completer := NewCompleter(root)
	rc := NewReadlineCompleter(completer)

	// Verify that the interface is satisfied (compile-time check is in the source,
	// but this is a runtime sanity check)
	if rc == nil {
		t.Fatal("expected non-nil ReadlineCompleter")
	}
}
