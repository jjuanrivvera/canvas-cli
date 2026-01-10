package main

import (
	"fmt"
	"os"

	"github.com/jjuanrivvera/canvas-cli/commands"
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
	if err := commands.Execute(Version, Commit, BuildDate); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
