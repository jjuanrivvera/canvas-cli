package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jjuanrivvera/canvas-cli/commands"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

var (
	// Version is set during build time
	Version = "dev"
	// Commit is set during build time
	Commit = "none"
	// BuildDate is set during build time
	BuildDate = "unknown"
)

func main() {
	// Expand aliases before executing commands
	expandAliases()

	if err := commands.Execute(Version, Commit, BuildDate); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// expandAliases checks if the first argument is an alias and expands it
func expandAliases() {
	if len(os.Args) < 2 {
		return
	}

	// Skip if first arg is a flag
	if strings.HasPrefix(os.Args[1], "-") {
		return
	}

	// Load config to get aliases
	cfg, err := config.Load()
	if err != nil {
		return // Silently ignore config errors
	}

	// Check if first arg is an alias
	firstArg := os.Args[1]
	expansion, exists := cfg.GetAlias(firstArg)
	if !exists {
		return
	}

	// Parse the expansion into args
	expandedArgs := parseAliasExpansion(expansion)
	if len(expandedArgs) == 0 {
		return
	}

	// Replace the alias with expanded args
	// os.Args[0] is the program name, os.Args[1] is the alias
	// os.Args[2:] are any additional args passed after the alias
	newArgs := make([]string, 0, len(os.Args)+len(expandedArgs))
	newArgs = append(newArgs, os.Args[0])
	newArgs = append(newArgs, expandedArgs...)
	newArgs = append(newArgs, os.Args[2:]...)
	os.Args = newArgs
}

// parseAliasExpansion splits an alias expansion into arguments
// Handles quoted strings properly
func parseAliasExpansion(expansion string) []string {
	var args []string
	var current strings.Builder
	inQuote := false
	quoteChar := rune(0)

	for _, r := range expansion {
		switch {
		case r == '"' || r == '\'':
			if inQuote && r == quoteChar {
				inQuote = false
				quoteChar = 0
			} else if !inQuote {
				inQuote = true
				quoteChar = r
			} else {
				current.WriteRune(r)
			}
		case r == ' ' && !inQuote:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}
