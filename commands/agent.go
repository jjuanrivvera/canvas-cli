package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// canvasIrreversibleVerbs are the Canvas operations that cannot be undone.
// "canvas agent guard" hard-blocks these by default; ordinary writes (create,
// update, publish, grade, …) only require approval.
var canvasIrreversibleVerbs = map[string]bool{
	"delete":      true,
	"remove":      true,
	"conclude":    true,
	"reset":       true,
	"abort":       true,
	"crosslist":   true,
	"uncrosslist": true,
	"deactivate":  true,
	"unpublish":   true,
	"unlink":      true,
	"clear":       true,
	"void":        true,
	"cancel":      true,
	"close":       true,
	"merge":       true, // user merge — data loss
	"split":       true, // reverse of merge, but lossy
}

// canvasReadVerbs is an explicit allowlist of read-only operations. Anything not
// in this set (and not irreversible) is treated as a write requiring approval —
// a fail-safe default so a verb the guard hasn't seen (a future command, or a
// mutating verb like "assign" or "bind" that isn't obviously a "create") is
// gated rather than silently allowed. Matched on the full leaf command name
// only (no token splitting), so a compound like "delete-entry" can never match
// a read token.
//
// This allowlist is hand-maintained. Misses are safe (extra approval prompts),
// not dangerous — the fail-safe default prevents silent read misclassification
// from becoming a missed write.
var canvasReadVerbs = map[string]bool{
	"get": true, "list": true, "show": true, "me": true, "search": true,
	"view": true, "feed": true, "front": true, "profile": true, "quota": true,
	"permissions": true, "tabs": true, "items": true, "members": true,
	"memberships": true, "comments": true, "replies": true, "entries": true,
	"entry-list": true, "results": true, "alignments": true, "revisions": true,
	"history": true, "history-day": true, "history-submissions": true,
	"changes": true, "issues": true, "migrators": true, "students": true,
	"recent-students": true, "missing-submissions": true, "page-views": true,
	"upcoming-events": true, "unread-count": true, "todo": true, "tracks": true,
	"runs": true, "next": true, "licenses": true, "errors": true,
	// Analytics — these are all read paths under "canvas analytics <resource>"
	"activity": true, "user": true, "assignments": true, "department": true,
	// Users sub-resources that are reads
	"courses": true, "groups": true,
	"effective-due-dates": true, "late-policy": true,
	"completed-statistics": true, "term-activity": true, "term-grades": true,
	"term-statistics": true, "list-closed": true, "list-enabled": true,
	"list-opened": true, "list-received": true, "list-sent": true,
	"list-for-account": true, "list-periods": true, "get-flag": true,
	"enabled": true, "resolve-path": true, "epub-get": true, "download": true,
	"content": true, "settings": true, "sso-settings": true, "logins": true,
	"overrides": true, "pages": true, "media": true, "collaborations": true,
	"conferences": true, "assignment-override": true,
	// Appointment group sub-reads
	"appointment-groups": true,
}

// canvasLocalGroups are top-level command groups that never call the Canvas
// API — they perform local operations only and must never be gated.
var canvasLocalGroups = map[string]bool{
	"auth":       true,
	"config":     true,
	"context":    true,
	"alias":      true,
	"cache":      true,
	"doctor":     true,
	"completion": true,
	"version":    true,
	"mcp":        true,
	"skills":     true,
	"repl":       true,
	"shell":      true,
	"update":     true,
	"telemetry":  true,
	"webhook":    true,
	"agent":      true,
	"help":       true,
}

// canvasBulkDestructivePaths lists full CLI paths (without root) that are
// ALWAYS hard-blocked regardless of their verb, because they can destroy large
// amounts of data in a single call. These supplement the verb-level blocks.
// "sis-imports create" batches a full institutional data import that can
// conclude/delete many enrollments/courses in one shot and cannot be rolled
// back safely.
var canvasBulkDestructivePaths = map[string]bool{
	"sis-imports create": true,
}

// canvasWriteOverridePaths lists full CLI paths whose leaf name collides with
// a read verb but which actually write. "sync assignments" ends in
// "assignments" (a read verb via "analytics assignments") yet copies
// assignments INTO the target course. Checked before the read allowlist so
// the collision cannot silently allow a write.
var canvasWriteOverridePaths = map[string]bool{
	"sync assignments": true,
}

// canvasGuardCmd represents one classified Canvas operation.
type canvasGuardCmd struct {
	cli  string // CLI path without root, e.g. "courses delete"
	tool string // MCP tool name, e.g. "canvas_courses_delete"
	verb string // leaf command name, e.g. "delete"
}

func init() {
	agentCmd := &cobra.Command{
		Use:   "agent",
		Short: "Helpers for running canvas under an AI agent",
		Long:  "Helpers for running canvas under an AI agent (Claude Code, Codex, OpenCode, …).",
	}
	agentCmd.AddCommand(newAgentGuardCmd())
	rootCmd.AddCommand(agentCmd)
}

func newAgentGuardCmd() *cobra.Command {
	var host string
	var allWrites bool
	var write bool

	cmd := &cobra.Command{
		Use:   "guard",
		Short: "Generate agent-safety config that blocks destructive canvas operations",
		Long: `guard generates the permission rules and hooks that stop an AI agent from
running destructive canvas operations, derived from the live command tree so the
list is always complete.

By default it hard-blocks the irreversible actions (delete, remove, conclude,
crosslist, deactivate, unpublish, cancel, close, merge, split, reset, ...) and
makes everything else that isn't a known read (create, update, grade, and any
verb the guard doesn't recognize) require approval. Classification is fail-safe:
only an explicit allowlist of read verbs (get, list, show, ...) stays allowed,
so a new or non-obvious mutating command is gated rather than slipping through.
Pass --all-writes to block writes too.

IMPORTANT: the "canvas api" escape hatch can issue any HTTP verb. The guard
blocks "canvas api DELETE/PUT/POST/PATCH" patterns on the Bash surface but
cannot enumerate arbitrary path arguments. For a hard guarantee, run the agent
MCP-only (no Bash tool) or inside a read-only sandbox.

Output is printed for review by default; pass --write to install it. See the
Agent Safety guide: https://jjuanrivvera.github.io/canvas-cli/user-guide/agent-safety/`,
		Args: cobra.NoArgs,
		Example: "  canvas agent guard --host claude-code\n" +
			"  canvas agent guard --host codex\n" +
			"  canvas agent guard --host opencode --all-writes\n" +
			"  canvas agent guard --host claude-code --write",
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, writes, irreversible := classifyCanvasCommands(rootCmd)
			g := canvasGuardPlan{
				irreversible: irreversible,
				writes:       writes,
				allWrites:    allWrites,
			}
			switch host {
			case "claude-code", "claude":
				return emitCanvasClaudeCode(cmd, g, write)
			case "codex":
				return emitCanvasCodex(cmd, g, write)
			case "opencode":
				return emitCanvasOpenCode(cmd, g, write)
			case "":
				return fmt.Errorf("--host is required (claude-code, codex, or opencode)")
			default:
				return fmt.Errorf("unknown host %q (use claude-code, codex, or opencode)", host)
			}
		},
	}
	cmd.Flags().StringVar(&host, "host", "", "Target agent host: claude-code, codex, opencode")
	cmd.Flags().BoolVar(&allWrites, "all-writes", false, "Also block create/update/grade (default: those require approval)")
	cmd.Flags().BoolVar(&write, "write", false, "Write the config/hook files instead of printing them")
	_ = cmd.RegisterFlagCompletionFunc("host", func(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		vals := canvasGuardHosts()
		out := make([]string, 0, len(vals))
		for _, v := range vals {
			if strings.HasPrefix(v, toComplete) {
				out = append(out, v)
			}
		}
		return out, cobra.ShellCompDirectiveNoFileComp
	})
	return cmd
}

// canvasGuardPlan is the classified set of operations the generators convert
// into host-specific config.
type canvasGuardPlan struct {
	irreversible []canvasGuardCmd
	writes       []canvasGuardCmd
	allWrites    bool
}

// blocked returns the operations to hard-block (irreversible, plus writes when
// --all-writes is set).
func (g canvasGuardPlan) blocked() []canvasGuardCmd {
	if g.allWrites {
		return append(append([]canvasGuardCmd{}, g.irreversible...), g.writes...)
	}
	return g.irreversible
}

// asked returns the operations that require approval (writes, unless
// --all-writes folds them into the hard-block set).
func (g canvasGuardPlan) asked() []canvasGuardCmd {
	if g.allWrites {
		return nil
	}
	return g.writes
}

// isCanvasIrreversibleVerb reports whether a command name contains an
// irreversible verb, handling compound names like "bulk-grade" by splitting on
// "-" (matching alegra's isIrreversibleVerb approach).
func isCanvasIrreversibleVerb(name string) bool {
	for _, tok := range strings.Split(name, "-") {
		if canvasIrreversibleVerbs[tok] {
			return true
		}
	}
	return false
}

// isCanvasReadVerb reports whether a command name is a known read-only verb.
// Matched on the full leaf name only — unknown verbs are deliberately NOT reads.
func isCanvasReadVerb(name string) bool {
	return canvasReadVerbs[name]
}

// topLevelGroup returns the top-level command group name for a given command by
// walking up the parent chain.
func topLevelGroup(c *cobra.Command) string {
	// Walk up to find the child of root
	cur := c
	for cur.Parent() != nil && cur.Parent().Parent() != nil {
		cur = cur.Parent()
	}
	if cur.Parent() == nil {
		// c itself is root
		return ""
	}
	return cur.Name()
}

// canvasClass buckets a single command for both "canvas agent guard" and the
// MCP tool annotations emitted by applyMCPAnnotations.
type canvasClass int

const (
	// canvasClassSkip marks commands that are neither gated nor annotated:
	// non-runnable parents, hidden/help commands, and the "canvas api"
	// raw-escape hatch, which can issue any HTTP verb and so must never
	// advertise readOnlyHint.
	canvasClassSkip canvasClass = iota
	// canvasClassLocal marks commands under local/utility top-level groups
	// (auth, config, cache, …) that never call the Canvas API. The guard does
	// not gate them; annotations handle them via canvasLocalReadPaths.
	canvasClassLocal
	canvasClassRead
	canvasClassWrite
	canvasClassIrreversible
)

// classifyCanvasCommand buckets one command. It is the single source of truth
// shared by "canvas agent guard" and the MCP readOnlyHint annotations, so a
// command can never be gated as a write by the guard while simultaneously
// advertising itself as read-only to an MCP client.
//
// Classification is purely by verb (leaf command name) — the cobra annotations
// flow the other way, derived from this function rather than feeding it.
func classifyCanvasCommand(root, sub *cobra.Command) (canvasClass, canvasGuardCmd) {
	if !sub.Runnable() || sub.Hidden || sub.Name() == "help" {
		return canvasClassSkip, canvasGuardCmd{}
	}

	gc := canvasGuardCmd{
		cli:  strings.TrimPrefix(sub.CommandPath(), root.Name()+" "),
		tool: strings.ReplaceAll(sub.CommandPath(), " ", "_"),
		verb: sub.Name(),
	}

	group := topLevelGroup(sub)

	// The "api" raw-escape command is documented separately in the hook script
	// header and cannot be safely classified by verb.
	if group == "api" || sub.Name() == "api" {
		// "canvas api get" is a GET-only sibling that never mutates state, so it
		// is a genuine read for both the guard and the MCP readOnlyHint. The
		// general "canvas api" escape hatch (any HTTP verb) stays skipped.
		if group == "api" && sub.Name() == "get" {
			return canvasClassRead, gc
		}
		return canvasClassSkip, gc
	}

	// Local/utility groups never hit the Canvas API.
	if canvasLocalGroups[group] {
		return canvasClassLocal, gc
	}

	// Fail-safe ordering: hard-block bulk-destructive paths and irreversible
	// verbs, allow only explicit reads, and treat everything else — known
	// writes AND any verb the guard doesn't recognize — as a write.
	switch {
	case canvasBulkDestructivePaths[gc.cli]:
		// Bulk-destructive full-path override: always hard-block regardless of
		// verb classification (e.g. "sis-imports create" batches a full
		// institutional import that can't be safely rolled back).
		return canvasClassIrreversible, gc
	case isCanvasIrreversibleVerb(sub.Name()):
		return canvasClassIrreversible, gc
	case canvasWriteOverridePaths[gc.cli]:
		return canvasClassWrite, gc
	case isCanvasReadVerb(sub.Name()):
		return canvasClassRead, gc
	default:
		return canvasClassWrite, gc
	}
}

// walkCanvasCommands visits every command below root, leaves first.
func walkCanvasCommands(root *cobra.Command, visit func(*cobra.Command)) {
	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		for _, sub := range c.Commands() {
			// Recurse first so we always visit leaves.
			walk(sub)
			visit(sub)
		}
	}
	walk(root)
}

// classifyCanvasCommands walks the command tree and buckets every runnable
// leaf command into read, write, or irreversible. Commands under local/utility
// top-level groups (auth, config, cache, agent, etc.) are excluded entirely —
// they never call the Canvas API.
func classifyCanvasCommands(root *cobra.Command) (read, writes, irreversible []canvasGuardCmd) {
	walkCanvasCommands(root, func(sub *cobra.Command) {
		class, gc := classifyCanvasCommand(root, sub)
		switch class {
		case canvasClassRead:
			read = append(read, gc)
		case canvasClassWrite:
			writes = append(writes, gc)
		case canvasClassIrreversible:
			irreversible = append(irreversible, gc)
		case canvasClassSkip, canvasClassLocal:
			// Not gated by the guard.
		}
	})
	sortCanvasGuard(read)
	sortCanvasGuard(writes)
	sortCanvasGuard(irreversible)
	return read, writes, irreversible
}

func sortCanvasGuard(cs []canvasGuardCmd) {
	sort.Slice(cs, func(i, j int) bool { return cs[i].tool < cs[j].tool })
}

func distinctCanvasVerbs(cs []canvasGuardCmd) []string {
	seen := map[string]bool{}
	var out []string
	for _, c := range cs {
		if !seen[c.verb] {
			seen[c.verb] = true
			out = append(out, c.verb)
		}
	}
	sort.Strings(out)
	return out
}

// canvasWriteOrPrint either writes content to path (creating parent dirs) when
// write is set and the file does not already exist, or prints it to the
// command's output with a header. It never overwrites an existing file.
func canvasWriteOrPrint(cmd *cobra.Command, write bool, path, content string, perm os.FileMode) error {
	out := cmd.OutOrStdout()
	if !write {
		fmt.Fprintf(out, "# ----- %s -----\n%s\n", path, content)
		return nil
	}
	if _, err := os.Stat(path); err == nil {
		fmt.Fprintf(out, "# %s already exists — review and merge manually:\n%s\n", path, content)
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(content), perm); err != nil { // #nosec G306 -- hook script needs 0o755
		return err
	}
	fmt.Fprintf(out, "wrote %s\n", path)
	return nil
}

// findProjectRoot walks up from the current working directory looking for a
// .git directory, which marks the project root. Falls back to CWD with a
// warning printed to cmd's output if no .git is found.
func findProjectRoot(cmd *cobra.Command) string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	dir := cwd
	for {
		if fi, err := os.Stat(filepath.Join(dir, ".git")); err == nil && fi.IsDir() {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root without finding .git.
			break
		}
		dir = parent
	}
	fmt.Fprintf(cmd.OutOrStdout(),
		"# warning: no .git directory found; using current directory as project root\n")
	return cwd
}
