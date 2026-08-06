package commands

import (
	"strings"
	"testing"

	"github.com/njayp/ophis"
	"github.com/spf13/cobra"
)

// findCmd resolves a space-separated CLI path (without the root name) to its
// command, or nil.
func findCmd(root *cobra.Command, path string) *cobra.Command {
	cur := root
	for _, part := range strings.Fields(path) {
		var next *cobra.Command
		for _, sub := range cur.Commands() {
			if sub.Name() == part {
				next = sub
				break
			}
		}
		if next == nil {
			return nil
		}
		cur = next
	}
	return cur
}

func annotationOf(cmd *cobra.Command, key string) (string, bool) {
	if cmd == nil || cmd.Annotations == nil {
		return "", false
	}
	v, ok := cmd.Annotations[key]
	return v, ok
}

func assertReadOnly(t *testing.T, root *cobra.Command, path, want string) {
	t.Helper()
	cmd := findCmd(root, path)
	if cmd == nil {
		t.Fatalf("command %q not found in tree", path)
	}
	got, ok := annotationOf(cmd, ophis.AnnotationReadOnly)
	if want == "" {
		if ok {
			t.Errorf("%q: expected no readOnlyHint, got %q", path, got)
		}
		return
	}
	if !ok {
		t.Errorf("%q: expected readOnlyHint=%s, got no annotation", path, want)
		return
	}
	if got != want {
		t.Errorf("%q: expected readOnlyHint=%s, got %s", path, want, got)
	}
}

// TestApplyMCPAnnotations_MatchesGuardClassification is the anti-drift guard: a
// command must never advertise readOnlyHint:true over MCP while "canvas agent
// guard" gates it as a write or hard-blocks it. Both derive from
// classifyCanvasCommand, and this asserts they cannot diverge on the real tree.
func TestApplyMCPAnnotations_MatchesGuardClassification(t *testing.T) {
	root := rootCmd
	applyMCPAnnotations(root)

	read, writes, irreversible := classifyCanvasCommands(root)

	if len(read) == 0 || len(writes) == 0 || len(irreversible) == 0 {
		t.Fatalf("real tree produced empty buckets: read=%d writes=%d irreversible=%d",
			len(read), len(writes), len(irreversible))
	}

	for _, c := range read {
		assertReadOnly(t, root, c.cli, "true")
	}
	for _, c := range writes {
		assertReadOnly(t, root, c.cli, "false")
		// Ordinary writes must not claim to be non-destructive; the hint is
		// left unset so clients apply the conservative spec default.
		if v, ok := annotationOf(findCmd(root, c.cli), ophis.AnnotationDestructive); ok {
			t.Errorf("%q: write should not set destructiveHint, got %q", c.cli, v)
		}
	}
	for _, c := range irreversible {
		assertReadOnly(t, root, c.cli, "false")
		if v, _ := annotationOf(findCmd(root, c.cli), ophis.AnnotationDestructive); v != "true" {
			t.Errorf("%q: irreversible should set destructiveHint=true, got %q", c.cli, v)
		}
	}
}

// TestApplyMCPAnnotations_APIEscapeHatchUnannotated pins the one deliberate
// exception: "canvas api" can issue any HTTP verb, so it must stay unannotated
// and therefore be filtered out by read-only MCP clients.
func TestApplyMCPAnnotations_APIEscapeHatchUnannotated(t *testing.T) {
	root := rootCmd
	applyMCPAnnotations(root)

	apiCmd := findCmd(root, "api")
	if apiCmd == nil {
		t.Fatal(`"canvas api" not found in tree`)
	}
	if _, ok := annotationOf(apiCmd, ophis.AnnotationReadOnly); ok {
		t.Error(`"canvas api" must not carry a readOnlyHint`)
	}
	// Every api subcommand stays unannotated EXCEPT the GET-only "api get"
	// sibling, which is a genuine read and must advertise readOnlyHint=true (#60).
	for _, sub := range apiCmd.Commands() {
		v, ok := annotationOf(sub, ophis.AnnotationReadOnly)
		if sub.Name() == "get" {
			if !ok || v != "true" {
				t.Errorf("%q must carry readOnlyHint=true, got %q (present=%v)", sub.CommandPath(), v, ok)
			}
			continue
		}
		if ok {
			t.Errorf("%q must not carry a readOnlyHint", sub.CommandPath())
		}
	}
}

// TestApplyMCPAnnotations_LocalGroups verifies the hand-maintained local
// allowlist: read-only local commands are annotated, and local commands that
// mutate state stay unannotated even when their leaf name looks inert.
func TestApplyMCPAnnotations_LocalGroups(t *testing.T) {
	root := rootCmd
	applyMCPAnnotations(root)

	for _, path := range []string{
		"version", "doctor", "auth status", "config show", "config list",
		"cache stats", "telemetry status", "update check", "webhook events",
	} {
		assertReadOnly(t, root, path, "true")
	}

	// These mutate local state and must be filtered out by read-only clients.
	// "config account" is the trap: it sets the default account ID.
	for _, path := range []string{
		"auth login", "auth logout", "config account", "config use",
		"context set", "cache clear", "skills install", "telemetry enable",
		"update", "webhook listen", "repl",
	} {
		assertReadOnly(t, root, path, "")
	}
}

// TestCanvasLocalReadPaths_ResolveToRealCommands stops the allowlist going
// stale: an entry that no longer matches a command (renamed or removed) would
// silently stop annotating it.
func TestCanvasLocalReadPaths_ResolveToRealCommands(t *testing.T) {
	for path := range canvasLocalReadPaths {
		cmd := findCmd(rootCmd, path)
		if cmd == nil {
			t.Errorf("canvasLocalReadPaths entry %q matches no command", path)
			continue
		}
		if !cmd.Runnable() {
			t.Errorf("canvasLocalReadPaths entry %q is not runnable", path)
		}
		if class, _ := classifyCanvasCommand(rootCmd, cmd); class != canvasClassLocal {
			t.Errorf("canvasLocalReadPaths entry %q classifies as %v, want local", path, class)
		}
	}
}

func TestApplyMCPAnnotations_TestTree(t *testing.T) {
	root := buildTestTree()
	applyMCPAnnotations(root)

	assertReadOnly(t, root, "courses list", "true")
	assertReadOnly(t, root, "courses get", "true")
	assertReadOnly(t, root, "courses create", "false")
	assertReadOnly(t, root, "courses delete", "false")
	assertReadOnly(t, root, "api GET", "")
	assertReadOnly(t, root, "auth login", "")

	if v, _ := annotationOf(findCmd(root, "courses delete"), ophis.AnnotationDestructive); v != "true" {
		t.Errorf(`"courses delete" destructiveHint = %q, want "true"`, v)
	}
	if v, _ := annotationOf(findCmd(root, "courses list"), ophis.AnnotationOpenWorld); v != "true" {
		t.Errorf(`"courses list" openWorldHint = %q, want "true"`, v)
	}
}

// TestApplyMCPAnnotations_UnknownVerbNotReadOnly mirrors the guard's fail-safe
// default: a verb the classifier does not recognize must never be advertised
// as read-only.
func TestApplyMCPAnnotations_UnknownVerbNotReadOnly(t *testing.T) {
	root := &cobra.Command{Use: "canvas"}
	grp := &cobra.Command{Use: "widgets"}
	grp.AddCommand(&cobra.Command{
		Use:  "frobnicate",
		RunE: func(_ *cobra.Command, _ []string) error { return nil },
	})
	root.AddCommand(grp)

	applyMCPAnnotations(root)

	assertReadOnly(t, root, "widgets frobnicate", "false")
}

func TestSetMCPAnnotations_PreservesUnrelatedKeys(t *testing.T) {
	cmd := &cobra.Command{Use: "thing", Annotations: map[string]string{"custom": "kept"}}

	setMCPAnnotations(cmd, true, false, false)

	if cmd.Annotations["custom"] != "kept" {
		t.Errorf("unrelated annotation clobbered: %v", cmd.Annotations)
	}
	if cmd.Annotations[ophis.AnnotationReadOnly] != "true" {
		t.Errorf("readOnlyHint = %q, want true", cmd.Annotations[ophis.AnnotationReadOnly])
	}
	if cmd.Annotations[ophis.AnnotationOpenWorld] != "false" {
		t.Errorf("openWorldHint = %q, want false", cmd.Annotations[ophis.AnnotationOpenWorld])
	}
	if _, ok := cmd.Annotations[ophis.AnnotationDestructive]; ok {
		t.Error("destructiveHint should be omitted when false")
	}
}

func TestSetMCPAnnotations_NilMap(t *testing.T) {
	cmd := &cobra.Command{Use: "thing"}

	setMCPAnnotations(cmd, false, true, true)

	if cmd.Annotations[ophis.AnnotationReadOnly] != "false" {
		t.Errorf("readOnlyHint = %q, want false", cmd.Annotations[ophis.AnnotationReadOnly])
	}
	if cmd.Annotations[ophis.AnnotationDestructive] != "true" {
		t.Errorf("destructiveHint = %q, want true", cmd.Annotations[ophis.AnnotationDestructive])
	}
}

// TestApplyMCPAnnotations_Idempotent guards against Execute and ExecuteContext
// both running the walk, and against tests dispatching through rootCmd repeatedly.
func TestApplyMCPAnnotations_Idempotent(t *testing.T) {
	root := buildTestTree()

	applyMCPAnnotations(root)
	first := map[string]string{}
	walkCanvasCommands(root, func(c *cobra.Command) {
		if v, ok := annotationOf(c, ophis.AnnotationReadOnly); ok {
			first[c.CommandPath()] = v
		}
	})

	applyMCPAnnotations(root)
	second := map[string]string{}
	walkCanvasCommands(root, func(c *cobra.Command) {
		if v, ok := annotationOf(c, ophis.AnnotationReadOnly); ok {
			second[c.CommandPath()] = v
		}
	})

	if len(first) != len(second) {
		t.Fatalf("annotation count changed: %d then %d", len(first), len(second))
	}
	for path, v := range first {
		if second[path] != v {
			t.Errorf("%q: %q then %q", path, v, second[path])
		}
	}
}
