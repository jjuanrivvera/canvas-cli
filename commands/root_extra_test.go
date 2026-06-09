package commands

import (
	"strings"
	"testing"
)

func TestGetRootCmd_NotNil(t *testing.T) {
	cmd := GetRootCmd()
	if cmd == nil {
		t.Fatal("GetRootCmd returned nil")
	}
	if cmd.Use != "canvas" {
		t.Errorf("expected Use=canvas, got %q", cmd.Use)
	}
}

func TestRootCmd_HasExpectedSubcommands(t *testing.T) {
	cmd := GetRootCmd()
	names := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		names[sub.Name()] = true
	}

	expected := []string{"courses", "assignments", "auth", "config", "context", "alias", "completion", "update"}
	for _, e := range expected {
		if !names[e] {
			t.Errorf("expected subcommand %q not found in root", e)
		}
	}
}

func TestRootCmd_HasPersistentFlags(t *testing.T) {
	cmd := GetRootCmd()
	flags := []string{"output", "verbose", "no-cache", "limit", "dry-run", "quiet", "filter", "columns", "sort"}
	for _, f := range flags {
		if cmd.PersistentFlags().Lookup(f) == nil {
			t.Errorf("expected persistent flag %q not found", f)
		}
	}
}

func TestGetHostnameFromURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://canvas.instructure.com", "canvas"},
		{"https://myschool.edu/canvas", "myschool"},
		{"https://www.example.com", "example"},
		{"http://localhost:8080", "localhost"},
		{"https://sub.domain.edu", "sub"},
		{"canvas.instructure.com", "canvas"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := getHostnameFromURL(tt.input)
			if got != tt.expected {
				t.Errorf("getHostnameFromURL(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFindIndex(t *testing.T) {
	tests := []struct {
		s, substr string
		expected  int
	}{
		{"hello://world", "://", 5},
		{"no match here", "xyz", -1},
		{"abcabc", "bc", 1},
		{"", "x", -1},
	}

	for _, tt := range tests {
		got := findIndex(tt.s, tt.substr)
		if got != tt.expected {
			t.Errorf("findIndex(%q, %q) = %d, want %d", tt.s, tt.substr, got, tt.expected)
		}
	}
}

func TestFindIndexFrom(t *testing.T) {
	tests := []struct {
		s, substr string
		start     int
		expected  int
	}{
		{"http://foo/bar", "/", 7, 10},
		{"http://foo/bar", "/", 0, 5},
		{"http://foo/bar", "xyz", 0, -1},
	}

	for _, tt := range tests {
		got := findIndexFrom(tt.s, tt.substr, tt.start)
		if got != tt.expected {
			t.Errorf("findIndexFrom(%q, %q, %d) = %d, want %d", tt.s, tt.substr, tt.start, got, tt.expected)
		}
	}
}

func TestContainsString(t *testing.T) {
	tests := []struct {
		slice    []string
		s        string
		expected bool
	}{
		{[]string{"foo", "bar", "baz"}, "bar", true},
		{[]string{"foo", "bar", "baz"}, "qux", false},
		{nil, "anything", false},
		{[]string{}, "x", false},
	}

	for _, tt := range tests {
		got := containsString(tt.slice, tt.s)
		if got != tt.expected {
			t.Errorf("containsString(%v, %q) = %v, want %v", tt.slice, tt.s, got, tt.expected)
		}
	}
}

func TestExecute_SetsVersion(t *testing.T) {
	// Execute with a version string and no subcommand — it prints help and returns nil
	// We use --help to avoid triggering interactive login
	old := rootCmd.Args
	rootCmd.SetArgs([]string{"--help"})
	defer func() {
		rootCmd.Args = old
		rootCmd.SetArgs(nil)
	}()

	err := Execute("1.2.3-test", "abc123", "2024-01-01")
	// --help exits with nil
	if err != nil && !strings.Contains(err.Error(), "help") {
		t.Logf("Execute returned (expected for help): %v", err)
	}
	if rootCmd.Version != "1.2.3-test" {
		t.Errorf("expected version 1.2.3-test, got %q", rootCmd.Version)
	}
}
