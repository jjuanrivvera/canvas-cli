package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

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
		RunE: runAliasSet,
	}

	return cmd
}

func runAliasSet(cmd *cobra.Command, args []string) error {
	name := args[0]
	expansion := args[1]

	// Validate alias name (no spaces, not a built-in command)
	if strings.Contains(name, " ") {
		return fmt.Errorf("alias name cannot contain spaces")
	}

	// Check if it conflicts with a built-in command
	for _, c := range rootCmd.Commands() {
		if c.Name() == name || containsString(c.Aliases, name) {
			return fmt.Errorf("cannot create alias %q: conflicts with built-in command", name)
		}
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.SetAlias(name, expansion); err != nil {
		return fmt.Errorf("failed to save alias: %w", err)
	}

	fmt.Printf("Alias %q set to: %s\n", name, expansion)
	return nil
}

func newAliasListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all aliases",
		Long:  `Display all configured command aliases.`,
		Args:  cobra.NoArgs,
		RunE:  runAliasList,
	}

	return cmd
}

func runAliasList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	aliases := cfg.ListAliases()
	if len(aliases) == 0 {
		fmt.Println("No aliases configured.")
		fmt.Println("\nCreate one with: canvas alias set <name> <command>")
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

	return formatOutput(data, nil)
}

func newAliasDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <name>",
		Aliases: []string{"rm", "remove"},
		Short:   "Delete an alias",
		Long:    `Remove a command alias.`,
		Args:    cobra.ExactArgs(1),
		RunE:    runAliasDelete,
	}

	return cmd
}

func runAliasDelete(cmd *cobra.Command, args []string) error {
	name := args[0]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.DeleteAlias(name); err != nil {
		return err
	}

	fmt.Printf("Alias %q deleted.\n", name)
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
