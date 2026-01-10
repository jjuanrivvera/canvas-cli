package repl

import (
	"strings"
)

// ANSI color codes
const (
	colorReset   = "\033[0m"
	colorBlue    = "\033[34m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorCyan    = "\033[36m"
	colorMagenta = "\033[35m"
	colorRed     = "\033[31m"
	colorBold    = "\033[1m"
)

// Highlighter provides syntax highlighting for REPL input
type Highlighter struct {
	commands    map[string]bool
	subcommands map[string]bool
	enabled     bool
}

// NewHighlighter creates a new syntax highlighter
func NewHighlighter(enabled bool) *Highlighter {
	h := &Highlighter{
		commands:    make(map[string]bool),
		subcommands: make(map[string]bool),
		enabled:     enabled,
	}

	// Register known commands
	knownCommands := []string{
		"courses", "assignments", "users", "enrollments",
		"submissions", "accounts", "config", "auth",
		"batch", "help", "exit", "quit", "sync",
	}
	for _, cmd := range knownCommands {
		h.commands[cmd] = true
	}

	// Register known subcommands
	knownSubcommands := []string{
		"list", "get", "create", "update", "delete",
		"sub", "grade", "sync", "import", "export",
		"login", "logout", "status", "set", "show",
	}
	for _, sub := range knownSubcommands {
		h.subcommands[sub] = true
	}

	return h
}

// Highlight applies syntax highlighting to the input string
func (h *Highlighter) Highlight(input string) string {
	if !h.enabled || input == "" {
		return input
	}

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return input
	}

	// Preserve original spacing by tracking positions
	result := h.highlightTokens(parts)
	return result
}

// highlightTokens applies highlighting to individual tokens
func (h *Highlighter) highlightTokens(parts []string) string {
	var highlighted []string

	for i, part := range parts {
		switch {
		case i == 0 && h.commands[part]:
			// Command in blue
			highlighted = append(highlighted, colorBlue+colorBold+part+colorReset)
		case i == 1 && h.subcommands[part]:
			// Subcommand in cyan
			highlighted = append(highlighted, colorCyan+part+colorReset)
		case strings.HasPrefix(part, "--"):
			// Long flags in green
			if idx := strings.Index(part, "="); idx > 0 {
				flag := part[:idx]
				value := part[idx+1:]
				highlighted = append(highlighted, colorGreen+flag+colorReset+"="+colorYellow+value+colorReset)
			} else {
				highlighted = append(highlighted, colorGreen+part+colorReset)
			}
		case strings.HasPrefix(part, "-"):
			// Short flags in green
			highlighted = append(highlighted, colorGreen+part+colorReset)
		case isNumeric(part):
			// Numeric IDs in magenta
			highlighted = append(highlighted, colorMagenta+part+colorReset)
		case strings.HasPrefix(part, "\"") || strings.HasPrefix(part, "'"):
			// Quoted strings in yellow
			highlighted = append(highlighted, colorYellow+part+colorReset)
		default:
			// Regular text, no highlighting
			highlighted = append(highlighted, part)
		}
	}

	return strings.Join(highlighted, " ")
}

// SetEnabled enables or disables syntax highlighting
func (h *Highlighter) SetEnabled(enabled bool) {
	h.enabled = enabled
}

// IsEnabled returns whether syntax highlighting is enabled
func (h *Highlighter) IsEnabled() bool {
	return h.enabled
}

// AddCommand adds a command to the known commands set
func (h *Highlighter) AddCommand(cmd string) {
	h.commands[cmd] = true
}

// AddSubcommand adds a subcommand to the known subcommands set
func (h *Highlighter) AddSubcommand(sub string) {
	h.subcommands[sub] = true
}

// isNumeric checks if a string is a numeric value
func isNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// StripANSI removes ANSI escape codes from a string
func StripANSI(s string) string {
	var result strings.Builder
	inEscape := false

	for _, c := range s {
		if c == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if c == 'm' {
				inEscape = false
			}
			continue
		}
		result.WriteRune(c)
	}

	return result.String()
}
