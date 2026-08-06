package commands

import (
	"github.com/njayp/ophis"
	"github.com/spf13/cobra"
)

// mcpReadonly is toggled by "canvas mcp start --readonly". When true, the MCP
// server exposes only tools the shared classifier annotated readOnlyHint=true,
// moving the read-only boundary into the binary instead of relying on the client
// to filter (some clients — e.g. the Copilot cloud agent — do not filter at all).
//
// Because "canvas api" is deliberately unannotated it drops out under --readonly,
// while its GET-only sibling "canvas api get" (canvas_api_get) survives. See #60.
var mcpReadonly bool

// mcpToolSelector decides whether a command is exposed as an MCP tool. In
// read-only mode it admits only commands annotated readOnlyHint=true; a
// missing/false annotation (writes, and the general "api" tool) is excluded —
// the same rule read-only clients apply, enforced in the binary so it holds
// even for clients that don't filter. Outside read-only mode it admits every
// command (basic ophis safety filters still apply).
func mcpToolSelector(cmd *cobra.Command) bool {
	if mcpReadonly && cmd.Annotations[ophis.AnnotationReadOnly] != "true" {
		return false
	}
	return true
}

func init() {
	mcpCmd := ophis.Command(&ophis.Config{
		ToolNamePrefix: "canvas",
		Selectors: []ophis.Selector{
			{
				CmdSelector: mcpToolSelector,
				// Exclude sensitive inherited flags from MCP exposure.
				// show-token and config are PersistentFlags (inherited).
				InheritedFlagSelector: ophis.ExcludeFlags("show-token", "config"),
			},
		},
	})

	// --readonly is persistent so it applies to "mcp start" (and the "mcp tools"
	// preview). ophis builds its tool list at run time, after flags are parsed,
	// so the selector above observes the correct value.
	mcpCmd.PersistentFlags().BoolVar(&mcpReadonly, "readonly", false,
		`Expose only read-only tools (drops write tools and the general "api" tool; keeps "api get")`)

	rootCmd.AddCommand(mcpCmd)
}
