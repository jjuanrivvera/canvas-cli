package commands

import (
	"context"
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

func init() {
	rootCmd.AddCommand(newAliasCmd())
}

func newAliasCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alias",
		Short: "Manage command aliases",
		Long: `Create, list, and delete custom command aliases.

Aliases allow you to create shortcuts for frequently used commands.
They can include flags and arguments.

Examples:
  # Create an alias
  canvas alias set ca "assignments list --course-id 123"

  # Use the alias
  canvas ca

  # Create an alias with multiple flags
  canvas alias set ungraded "assignments list --course-id 123 --bucket ungraded"

  # List all aliases
  canvas alias list

  # Delete an alias
  canvas alias delete ca`,
	}

	cmd.AddCommand(newAliasSetCmd())
	cmd.AddCommand(newAliasListCmd())
	cmd.AddCommand(newAliasDeleteCmd())

	return cmd
}

func newAliasSetCmd() *cobra.Command {
	opts := &options.AliasSetOptions{}

	cmd := &cobra.Command{
		Use:   "set <name> <expansion>",
		Short: "Create or update an alias",
		Long: `Create or update a command alias.

The alias name should be a single word. The expansion is the command
that will be executed when the alias is used.

Examples:
  canvas alias set ca "assignments list --course-id 123"
  canvas alias set myc "courses list --enrollment-type teacher"
  canvas alias set grade "submissions bulk-grade --csv"`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			opts.Expansion = args[1]
			if err := opts.Validate(); err != nil {
				return err
			}
			return runAliasSet(cmd.Context(), opts)
		},
	}

	return cmd
}

func runAliasSet(ctx context.Context, opts *options.AliasSetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "alias.set", map[string]interface{}{
		"name":      opts.Name,
		"expansion": opts.Expansion,
	})

	// Check if it conflicts with a built-in command
	for _, c := range rootCmd.Commands() {
		if c.Name() == opts.Name || containsString(c.Aliases, opts.Name) {
			err := fmt.Errorf("cannot create alias %q: conflicts with built-in command", opts.Name)
			logger.LogCommandError(ctx, "alias.set", err, nil)
			return err
		}
	}

	cfg, err := config.Load()
	if err != nil {
		logger.LogCommandError(ctx, "alias.set", err, nil)
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.SetAlias(opts.Name, opts.Expansion); err != nil {
		logger.LogCommandError(ctx, "alias.set", err, nil)
		return fmt.Errorf("failed to save alias: %w", err)
	}

	logger.LogCommandComplete(ctx, "alias.set", 1)
	fmt.Printf("Alias %q set to: %s\n", opts.Name, opts.Expansion)
	return nil
}

func newAliasListCmd() *cobra.Command {
	opts := &options.AliasListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all aliases",
		Long:  `Display all configured command aliases.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			return runAliasList(cmd.Context(), opts)
		},
	}

	return cmd
}

func runAliasList(ctx context.Context, opts *options.AliasListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "alias.list", nil)

	cfg, err := config.Load()
	if err != nil {
		logger.LogCommandError(ctx, "alias.list", err, nil)
		return fmt.Errorf("failed to load config: %w", err)
	}

	aliases := cfg.ListAliases()
	if len(aliases) == 0 {
		fmt.Println("No aliases configured.")
		fmt.Println("\nCreate one with: canvas alias set <name> <command>")
		logger.LogCommandComplete(ctx, "alias.list", 0)
		return nil
	}

	// Sort aliases by name for consistent output
	names := make([]string, 0, len(aliases))
	for name := range aliases {
		names = append(names, name)
	}
	sort.Strings(names)

	// Build output data
	type aliasEntry struct {
		Name      string `json:"name"`
		Expansion string `json:"expansion"`
	}

	data := make([]aliasEntry, 0, len(aliases))
	for _, name := range names {
		data = append(data, aliasEntry{
			Name:      name,
			Expansion: aliases[name],
		})
	}

	logger.LogCommandComplete(ctx, "alias.list", len(data))
	return formatOutput(data, nil)
}

func newAliasDeleteCmd() *cobra.Command {
	opts := &options.AliasDeleteOptions{}

	cmd := &cobra.Command{
		Use:     "delete <name>",
		Aliases: []string{"rm", "remove"},
		Short:   "Delete an alias",
		Long:    `Remove a command alias.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			if err := opts.Validate(); err != nil {
				return err
			}
			return runAliasDelete(cmd.Context(), opts)
		},
	}

	return cmd
}

func runAliasDelete(ctx context.Context, opts *options.AliasDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "alias.delete", map[string]interface{}{
		"name": opts.Name,
	})

	cfg, err := config.Load()
	if err != nil {
		logger.LogCommandError(ctx, "alias.delete", err, nil)
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.DeleteAlias(opts.Name); err != nil {
		logger.LogCommandError(ctx, "alias.delete", err, nil)
		return err
	}

	logger.LogCommandComplete(ctx, "alias.delete", 1)
	fmt.Printf("Alias %q deleted.\n", opts.Name)
	return nil
}

// containsString checks if a slice contains a string
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
