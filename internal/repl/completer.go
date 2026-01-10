package repl

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Completer provides command completion for the REPL
type Completer struct {
	rootCmd *cobra.Command
}

// NewCompleter creates a new completer
func NewCompleter(rootCmd *cobra.Command) *Completer {
	return &Completer{
		rootCmd: rootCmd,
	}
}

// Complete returns completion suggestions for the given input
func (c *Completer) Complete(input string) []string {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return c.rootCommands()
	}

	// If the input ends with a space, we're completing the next argument
	endsWithSpace := strings.HasSuffix(input, " ")
	if endsWithSpace {
		parts = append(parts, "")
	}

	// Get the command being completed
	cmd, args, err := c.findCommand(parts)
	if err != nil {
		return nil
	}

	// If we're on the first part and no space, complete command names
	if len(args) == 0 && !endsWithSpace {
		return c.matchCommands(cmd, parts[0])
	}

	// Complete subcommands
	if len(args) == 1 && cmd.HasSubCommands() {
		prefix := ""
		if !endsWithSpace {
			prefix = args[0]
		}
		return c.matchSubcommands(cmd, prefix)
	}

	// Complete flags
	return c.matchFlags(cmd, args[len(args)-1])
}

// rootCommands returns all root-level commands
func (c *Completer) rootCommands() []string {
	commands := make([]string, 0)
	for _, cmd := range c.rootCmd.Commands() {
		if !cmd.Hidden {
			commands = append(commands, cmd.Name())
		}
	}
	// Add REPL-specific commands
	commands = append(commands, "history", "clear", "session", "exit", "quit")
	return commands
}

// matchCommands returns commands matching the given prefix
func (c *Completer) matchCommands(parent *cobra.Command, prefix string) []string {
	matches := make([]string, 0)
	for _, cmd := range parent.Commands() {
		if !cmd.Hidden && strings.HasPrefix(cmd.Name(), prefix) {
			matches = append(matches, cmd.Name())
		}
	}
	return matches
}

// matchSubcommands returns subcommands matching the given prefix
func (c *Completer) matchSubcommands(parent *cobra.Command, prefix string) []string {
	matches := make([]string, 0)
	for _, cmd := range parent.Commands() {
		if !cmd.Hidden && strings.HasPrefix(cmd.Name(), prefix) {
			matches = append(matches, cmd.Name())
		}
	}
	return matches
}

// matchFlags returns flags matching the given prefix
func (c *Completer) matchFlags(cmd *cobra.Command, prefix string) []string {
	if !strings.HasPrefix(prefix, "-") {
		return nil
	}

	matches := make([]string, 0)

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Hidden {
			return
		}

		longFlag := "--" + flag.Name
		if strings.HasPrefix(longFlag, prefix) {
			matches = append(matches, longFlag)
		}

		if flag.Shorthand != "" {
			shortFlag := "-" + flag.Shorthand
			if strings.HasPrefix(shortFlag, prefix) {
				matches = append(matches, shortFlag)
			}
		}
	})

	return matches
}

// findCommand finds the command being completed
func (c *Completer) findCommand(parts []string) (*cobra.Command, []string, error) {
	cmd := c.rootCmd
	var args []string

	for i, part := range parts {
		// Skip flags
		if strings.HasPrefix(part, "-") {
			continue
		}

		// Try to find subcommand
		found := false
		for _, subCmd := range cmd.Commands() {
			if subCmd.Name() == part {
				cmd = subCmd
				found = true
				break
			}
		}

		// If not found, treat remaining as args
		if !found {
			args = parts[i:]
			break
		}
	}

	return cmd, args, nil
}

// GetCommandHelp returns help text for a command
func (c *Completer) GetCommandHelp(cmdName string) string {
	for _, cmd := range c.rootCmd.Commands() {
		if cmd.Name() == cmdName {
			return cmd.Short
		}
	}
	return ""
}

// GetFlagHelp returns help text for a flag
func (c *Completer) GetFlagHelp(cmd *cobra.Command, flagName string) string {
	flagName = strings.TrimPrefix(flagName, "--")
	flagName = strings.TrimPrefix(flagName, "-")

	flag := cmd.Flags().Lookup(flagName)
	if flag != nil {
		return flag.Usage
	}
	return ""
}
