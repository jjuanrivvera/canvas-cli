package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jjuanrivvera/canvas-cli/commands"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
	"github.com/jjuanrivvera/canvas-cli/internal/shellparse"
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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := commands.ExecuteContext(ctx, Version, Commit, BuildDate); err != nil {
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

	// Parse the expansion into args (quote-aware)
	expandedArgs := shellparse.Split(expansion)
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
