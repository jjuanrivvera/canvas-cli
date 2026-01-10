package repl

import (
	"strings"
	"testing"
)

func TestNewHighlighter(t *testing.T) {
	h := NewHighlighter(true)

	if h == nil {
		t.Fatal("expected non-nil highlighter")
	}

	if !h.enabled {
		t.Error("expected highlighter to be enabled")
	}

	// Check that known commands are registered
	if !h.commands["courses"] {
		t.Error("expected 'courses' to be a known command")
	}

	// Check that known subcommands are registered
	if !h.subcommands["list"] {
		t.Error("expected 'list' to be a known subcommand")
	}
}

func TestHighlighter_Highlight_Disabled(t *testing.T) {
	h := NewHighlighter(false)

	input := "courses list"
	result := h.Highlight(input)

	if result != input {
		t.Errorf("expected unchanged input when disabled, got %q", result)
	}
}

func TestHighlighter_Highlight_Empty(t *testing.T) {
	h := NewHighlighter(true)

	result := h.Highlight("")
	if result != "" {
		t.Errorf("expected empty string for empty input, got %q", result)
	}
}

func TestHighlighter_Highlight_Command(t *testing.T) {
	h := NewHighlighter(true)

	result := h.Highlight("courses")

	// Should contain ANSI blue color code
	if !strings.Contains(result, colorBlue) {
		t.Error("expected command to be highlighted in blue")
	}

	// Stripped version should match original
	if StripANSI(result) != "courses" {
		t.Errorf("stripped result should be 'courses', got %q", StripANSI(result))
	}
}

func TestHighlighter_Highlight_CommandAndSubcommand(t *testing.T) {
	h := NewHighlighter(true)

	result := h.Highlight("courses list")

	// Should contain blue for command and cyan for subcommand
	if !strings.Contains(result, colorBlue) {
		t.Error("expected command to be highlighted in blue")
	}
	if !strings.Contains(result, colorCyan) {
		t.Error("expected subcommand to be highlighted in cyan")
	}

	// Stripped version should match original
	if StripANSI(result) != "courses list" {
		t.Errorf("stripped result should be 'courses list', got %q", StripANSI(result))
	}
}

func TestHighlighter_Highlight_Flags(t *testing.T) {
	h := NewHighlighter(true)

	tests := []struct {
		name  string
		input string
	}{
		{"long flag", "courses list --account"},
		{"short flag", "courses list -a"},
		{"flag with value", "courses list --account=123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := h.Highlight(tt.input)

			// Should contain green for flags
			if !strings.Contains(result, colorGreen) {
				t.Error("expected flag to be highlighted in green")
			}

			// Stripped version should match original
			if StripANSI(result) != tt.input {
				t.Errorf("stripped result should be %q, got %q", tt.input, StripANSI(result))
			}
		})
	}
}

func TestHighlighter_Highlight_NumericIDs(t *testing.T) {
	h := NewHighlighter(true)

	result := h.Highlight("courses get 12345")

	// Should contain magenta for numeric ID
	if !strings.Contains(result, colorMagenta) {
		t.Error("expected numeric ID to be highlighted in magenta")
	}

	// Stripped version should match original
	if StripANSI(result) != "courses get 12345" {
		t.Errorf("stripped result should be 'courses get 12345', got %q", StripANSI(result))
	}
}

func TestHighlighter_Highlight_QuotedStrings(t *testing.T) {
	h := NewHighlighter(true)

	tests := []struct {
		name  string
		input string
	}{
		{"double quotes", `courses list --search "Biology"`},
		{"single quotes", `courses list --search 'Biology'`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := h.Highlight(tt.input)

			// Should contain yellow for quoted strings
			if !strings.Contains(result, colorYellow) {
				t.Error("expected quoted string to be highlighted in yellow")
			}
		})
	}
}

func TestHighlighter_SetEnabled(t *testing.T) {
	h := NewHighlighter(true)

	if !h.IsEnabled() {
		t.Error("expected highlighter to be enabled initially")
	}

	h.SetEnabled(false)
	if h.IsEnabled() {
		t.Error("expected highlighter to be disabled after SetEnabled(false)")
	}

	h.SetEnabled(true)
	if !h.IsEnabled() {
		t.Error("expected highlighter to be enabled after SetEnabled(true)")
	}
}

func TestHighlighter_AddCommand(t *testing.T) {
	h := NewHighlighter(true)

	// Custom command should not be highlighted initially
	result1 := h.Highlight("mycommand")
	if strings.Contains(result1, colorBlue) {
		t.Error("unexpected blue highlight for unknown command")
	}

	// After adding, it should be highlighted
	h.AddCommand("mycommand")
	result2 := h.Highlight("mycommand")
	if !strings.Contains(result2, colorBlue) {
		t.Error("expected blue highlight for added command")
	}
}

func TestHighlighter_AddSubcommand(t *testing.T) {
	h := NewHighlighter(true)

	// Custom subcommand should not be highlighted initially
	result1 := h.Highlight("courses mysub")
	if strings.Count(result1, colorCyan) > 0 {
		t.Error("unexpected cyan highlight for unknown subcommand")
	}

	// After adding, it should be highlighted
	h.AddSubcommand("mysub")
	result2 := h.Highlight("courses mysub")
	if !strings.Contains(result2, colorCyan) {
		t.Error("expected cyan highlight for added subcommand")
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123", true},
		{"0", true},
		{"999999", true},
		{"", false},
		{"12a34", false},
		{"abc", false},
		{"-123", false},
		{"12.34", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isNumeric(tt.input)
			if result != tt.expected {
				t.Errorf("isNumeric(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStripANSI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no ANSI codes",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "simple color",
			input:    "\033[34mblue\033[0m",
			expected: "blue",
		},
		{
			name:     "multiple colors",
			input:    "\033[34mblue\033[0m \033[32mgreen\033[0m",
			expected: "blue green",
		},
		{
			name:     "nested colors",
			input:    "\033[1m\033[34mbold blue\033[0m",
			expected: "bold blue",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripANSI(tt.input)
			if result != tt.expected {
				t.Errorf("StripANSI(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestHighlighter_ComplexInput(t *testing.T) {
	h := NewHighlighter(true)

	input := "courses list --account=123 --search \"Biology\" -v"
	result := h.Highlight(input)

	// Verify stripped output matches input
	stripped := StripANSI(result)
	if stripped != input {
		t.Errorf("stripped result should match input\ngot:      %q\nexpected: %q", stripped, input)
	}

	// Verify colors are present
	if !strings.Contains(result, colorBlue) {
		t.Error("expected blue for command")
	}
	if !strings.Contains(result, colorCyan) {
		t.Error("expected cyan for subcommand")
	}
	if !strings.Contains(result, colorGreen) {
		t.Error("expected green for flags")
	}
}

// Benchmark tests

func BenchmarkHighlighter_Highlight_Simple(b *testing.B) {
	h := NewHighlighter(true)
	input := "courses list"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = h.Highlight(input)
	}
}

func BenchmarkHighlighter_Highlight_Complex(b *testing.B) {
	h := NewHighlighter(true)
	input := "courses list --account=123 --search \"Biology\" -v --state available"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = h.Highlight(input)
	}
}

func BenchmarkStripANSI(b *testing.B) {
	input := "\033[34m\033[1mcourses\033[0m \033[36mlist\033[0m \033[32m--account\033[0m=\033[33m123\033[0m"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StripANSI(input)
	}
}
