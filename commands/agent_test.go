package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// --- Verb classification ---

func TestIsCanvasIrreversibleVerb(t *testing.T) {
	cases := []struct {
		verb string
		want bool
	}{
		{"delete", true},
		{"remove", true},
		{"conclude", true},
		{"reset", true},
		{"abort", true},
		{"crosslist", true},
		{"uncrosslist", true},
		{"deactivate", true},
		{"unpublish", true},
		{"unlink", true},
		{"clear", true},
		{"void", true},
		{"cancel", true},
		{"close", true},
		{"merge", true},
		{"split", true},
		// Compound names: any token in "-"-split that matches the set.
		{"bulk-delete", true},
		// Non-irreversible verbs.
		{"create", false},
		{"update", false},
		{"list", false},
		{"get", false},
		{"publish", false},
		{"grade", false},
	}
	for _, tc := range cases {
		t.Run(tc.verb, func(t *testing.T) {
			got := isCanvasIrreversibleVerb(tc.verb)
			if got != tc.want {
				t.Errorf("isCanvasIrreversibleVerb(%q) = %v, want %v", tc.verb, got, tc.want)
			}
		})
	}
}

// TestCanvasWriteVerbsViaClassification replaces TestIsCanvasWriteVerb.
// canvasWriteVerbs and isCanvasWriteVerb have been removed — writes come from
// the classifier's fail-safe default (anything not a read or irreversible).
// This test verifies that the expected write verbs land in the writes bucket
// and that reads and irreversibles do NOT.
func TestCanvasWriteVerbsViaClassification(t *testing.T) {
	root := &cobra.Command{Use: "canvas"}
	grp := &cobra.Command{Use: "testgroup"}

	// Add one command per verb we want to exercise.
	writeVerbs := []string{
		"create", "update", "publish", "grade", "bulk-grade", "upload",
		"add", "set", "move", "duplicate", "reply", "post", "send", "enroll",
		"accept", "reject", "reactivate", "star", "unstar", "archive",
		"unarchive", "subscribe", "unsubscribe", "relock", "revert",
		"associate", "sync", "restore", "mark-read", "mark-all-read",
		"complete", "dismiss", "done",
	}
	readVerbs := []string{"list", "get", "show", "me", "search", "quota", "front", "revisions"}
	irreversibleVerbsSlice := []string{"delete", "remove", "conclude", "crosslist"}

	for _, v := range writeVerbs {
		vv := v
		grp.AddCommand(&cobra.Command{Use: vv, RunE: func(_ *cobra.Command, _ []string) error { return nil }})
	}
	for _, v := range readVerbs {
		vv := v
		grp.AddCommand(&cobra.Command{Use: vv, RunE: func(_ *cobra.Command, _ []string) error { return nil }})
	}
	for _, v := range irreversibleVerbsSlice {
		vv := v
		grp.AddCommand(&cobra.Command{Use: vv, RunE: func(_ *cobra.Command, _ []string) error { return nil }})
	}
	root.AddCommand(grp)

	read, writes, irreversible := classifyCanvasCommands(root)

	for _, v := range writeVerbs {
		assertContainsPath(t, "write", writes, "testgroup "+v)
		assertNotContainsPath(t, "read (misclassified)", read, "testgroup "+v)
		assertNotContainsPath(t, "irreversible (misclassified)", irreversible, "testgroup "+v)
	}
	for _, v := range readVerbs {
		assertContainsPath(t, "read", read, "testgroup "+v)
		assertNotContainsPath(t, "write (misclassified)", writes, "testgroup "+v)
	}
	for _, v := range irreversibleVerbsSlice {
		assertContainsPath(t, "irreversible", irreversible, "testgroup "+v)
		assertNotContainsPath(t, "write (misclassified)", writes, "testgroup "+v)
	}
}

// --- Command tree classification ---

// buildTestTree builds a minimal command tree that exercises all buckets.
// The top-level group names must NOT be in canvasLocalGroups.
func buildTestTree() *cobra.Command {
	root := &cobra.Command{Use: "canvas"}

	addLeaf := func(parent *cobra.Command, name string) {
		parent.AddCommand(&cobra.Command{
			Use:  name,
			RunE: func(_ *cobra.Command, _ []string) error { return nil },
		})
	}

	// API group — should be skipped entirely.
	apiCmd := &cobra.Command{Use: "api"}
	addLeaf(apiCmd, "GET")
	root.AddCommand(apiCmd)

	// Local utility — should be skipped.
	authCmd := &cobra.Command{Use: "auth"}
	addLeaf(authCmd, "login")
	root.AddCommand(authCmd)

	// Canvas API groups.
	coursesCmd := &cobra.Command{Use: "courses"}
	addLeaf(coursesCmd, "list")   // read
	addLeaf(coursesCmd, "get")    // read
	addLeaf(coursesCmd, "create") // write
	addLeaf(coursesCmd, "update") // write
	addLeaf(coursesCmd, "delete") // irreversible
	root.AddCommand(coursesCmd)

	assignmentsCmd := &cobra.Command{Use: "assignments"}
	addLeaf(assignmentsCmd, "list")   // read
	addLeaf(assignmentsCmd, "create") // write
	addLeaf(assignmentsCmd, "grade")  // write
	addLeaf(assignmentsCmd, "delete") // irreversible
	root.AddCommand(assignmentsCmd)

	sectionsCmd := &cobra.Command{Use: "sections"}
	addLeaf(sectionsCmd, "list")        // read
	addLeaf(sectionsCmd, "crosslist")   // irreversible
	addLeaf(sectionsCmd, "uncrosslist") // irreversible
	root.AddCommand(sectionsCmd)

	return root
}

func TestClassifyCanvasCommands_Buckets(t *testing.T) {
	root := buildTestTree()
	read, writes, irreversible := classifyCanvasCommands(root)

	// Verify approximate counts — we expect:
	// read: courses list, courses get, assignments list, sections list = 4
	// writes: courses create, courses update, assignments create, assignments grade = 4
	// irreversible: courses delete, assignments delete, sections crosslist, sections uncrosslist = 4

	if len(read) != 4 {
		t.Errorf("expected 4 read commands, got %d: %v", len(read), cliPaths(read))
	}
	if len(writes) != 4 {
		t.Errorf("expected 4 write commands, got %d: %v", len(writes), cliPaths(writes))
	}
	if len(irreversible) != 4 {
		t.Errorf("expected 4 irreversible commands, got %d: %v", len(irreversible), cliPaths(irreversible))
	}

	// Verify specific entries are in the right bucket.
	assertContainsPath(t, "read", read, "courses list")
	assertContainsPath(t, "read", read, "assignments list")
	assertContainsPath(t, "write", writes, "courses create")
	assertContainsPath(t, "write", writes, "assignments grade")
	assertContainsPath(t, "irreversible", irreversible, "courses delete")
	assertContainsPath(t, "irreversible", irreversible, "sections crosslist")
	assertContainsPath(t, "irreversible", irreversible, "sections uncrosslist")

	// Verify local/utility commands are excluded.
	assertNotContainsPath(t, "read", read, "auth login")
	assertNotContainsPath(t, "write", writes, "auth login")
	assertNotContainsPath(t, "irreversible", irreversible, "auth login")

	// Verify api command is excluded.
	assertNotContainsPath(t, "read", read, "api GET")
}

// TestClassifyCanvasCommands_FailSafeDefault is the regression guard for the
// classifier's most important property: a verb the guard does not recognize is
// treated as a WRITE requiring approval, never silently allowed as a read. This
// keeps future commands (and non-obvious mutating verbs) gated by default.
func TestClassifyCanvasCommands_FailSafeDefault(t *testing.T) {
	root := &cobra.Command{Use: "canvas"}
	grp := &cobra.Command{Use: "widgets"}
	// A verb that exists in no set today (stands in for a future command).
	for _, v := range []string{"frobnicate", "merge", "cancel", "split", "close", "get"} {
		vv := v
		grp.AddCommand(&cobra.Command{Use: vv, RunE: func(_ *cobra.Command, _ []string) error { return nil }})
	}
	root.AddCommand(grp)

	read, writes, irreversible := classifyCanvasCommands(root)

	// Unknown verb must NOT be allowed.
	assertNotContainsPath(t, "read", read, "widgets frobnicate")
	assertContainsPath(t, "write (fail-safe)", writes, "widgets frobnicate")

	// Destructive verbs that were previously mis-bucketed must hard-block.
	for _, v := range []string{"merge", "cancel", "split", "close"} {
		assertContainsPath(t, "irreversible", irreversible, "widgets "+v)
	}

	// A genuine read still reads.
	assertContainsPath(t, "read", read, "widgets get")
}

func TestClassifyCanvasCommands_LocalGroupsExcluded(t *testing.T) {
	// Build a tree with all local groups and verify none leak through.
	root := &cobra.Command{Use: "canvas"}
	for group := range canvasLocalGroups {
		g := &cobra.Command{Use: group}
		g.AddCommand(&cobra.Command{
			Use:  "do-something",
			RunE: func(_ *cobra.Command, _ []string) error { return nil },
		})
		root.AddCommand(g)
	}
	read, writes, irreversible := classifyCanvasCommands(root)
	if len(read)+len(writes)+len(irreversible) != 0 {
		t.Errorf("expected all local groups to be excluded, but got %d total classified commands",
			len(read)+len(writes)+len(irreversible))
	}
}

// TestClassifyCanvasCommands_BulkDestructivePaths verifies that paths in
// canvasBulkDestructivePaths are hard-blocked even when their verb is a write.
func TestClassifyCanvasCommands_BulkDestructivePaths(t *testing.T) {
	root := &cobra.Command{Use: "canvas"}
	sisCmd := &cobra.Command{Use: "sis-imports"}
	sisCmd.AddCommand(&cobra.Command{
		Use:  "create", // verb is a write, but path is bulk-destructive
		RunE: func(_ *cobra.Command, _ []string) error { return nil },
	})
	sisCmd.AddCommand(&cobra.Command{
		Use:  "list",
		RunE: func(_ *cobra.Command, _ []string) error { return nil },
	})
	root.AddCommand(sisCmd)

	_, writes, irreversible := classifyCanvasCommands(root)

	// "sis-imports create" must be hard-blocked, NOT in writes.
	assertContainsPath(t, "irreversible", irreversible, "sis-imports create")
	assertNotContainsPath(t, "write (must not)", writes, "sis-imports create")

	// "sis-imports list" is an ordinary read — not affected by the override.
	read, _, _ := classifyCanvasCommands(root)
	assertContainsPath(t, "read", read, "sis-imports list")
}

// TestClassifyCanvasCommands_WriteOverridePaths verifies that paths whose leaf
// name collides with a read verb but actually write (e.g. "sync assignments",
// which shares its leaf with the "analytics assignments" read) are forced into
// the write bucket by canvasWriteOverridePaths.
func TestClassifyCanvasCommands_WriteOverridePaths(t *testing.T) {
	read, writes, _ := classifyCanvasCommands(rootCmd)

	// "sync assignments" copies assignments INTO a target course — must ask.
	assertContainsPath(t, "write (override)", writes, "sync assignments")
	assertNotContainsPath(t, "read (must not)", read, "sync assignments")

	// The genuine read that shares the leaf name is unaffected.
	assertContainsPath(t, "read", read, "analytics assignments")
	// The sibling sync command stays a write via the fail-safe default.
	assertContainsPath(t, "write", writes, "sync course")
}

// --- guardPlan ---

func TestGuardPlan_BlockedAndAsked(t *testing.T) {
	root := buildTestTree()
	_, writes, irreversible := classifyCanvasCommands(root)

	t.Run("default (irreversible blocked, writes asked)", func(t *testing.T) {
		g := canvasGuardPlan{irreversible: irreversible, writes: writes, allWrites: false}
		blocked := g.blocked()
		asked := g.asked()
		if len(blocked) != len(irreversible) {
			t.Errorf("blocked count %d != irreversible count %d", len(blocked), len(irreversible))
		}
		if len(asked) != len(writes) {
			t.Errorf("asked count %d != writes count %d", len(asked), len(writes))
		}
	})

	t.Run("--all-writes (everything blocked, nothing asked)", func(t *testing.T) {
		g := canvasGuardPlan{irreversible: irreversible, writes: writes, allWrites: true}
		blocked := g.blocked()
		asked := g.asked()
		if len(blocked) != len(irreversible)+len(writes) {
			t.Errorf("blocked count %d, want %d", len(blocked), len(irreversible)+len(writes))
		}
		if len(asked) != 0 {
			t.Errorf("asked count %d, want 0", len(asked))
		}
	})
}

func TestDistinctCanvasVerbs(t *testing.T) {
	cmds := []canvasGuardCmd{
		{verb: "delete"}, {verb: "delete"}, {verb: "remove"}, {verb: "remove"},
	}
	got := distinctCanvasVerbs(cmds)
	if len(got) != 2 {
		t.Errorf("expected 2 distinct verbs, got %d: %v", len(got), got)
	}
	if got[0] != "delete" || got[1] != "remove" {
		t.Errorf("unexpected verbs %v", got)
	}
}

// --- Emit functions (print mode) ---

func TestEmitCanvasClaudeCode_PrintsWithoutWrite(t *testing.T) {
	root := buildTestTree()
	_, writes, irreversible := classifyCanvasCommands(root)
	g := canvasGuardPlan{irreversible: irreversible, writes: writes, allWrites: false}

	cmd := &cobra.Command{Use: "canvas"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := emitCanvasClaudeCode(cmd, g, false); err != nil {
		t.Fatalf("emitCanvasClaudeCode: %v", err)
	}
	out := buf.String()

	// Must include the settings path header.
	if !strings.Contains(out, ".claude/settings.json") {
		t.Error("expected .claude/settings.json reference in output")
	}
	// Must include the hook script path.
	if !strings.Contains(out, ".claude/hooks/canvas-guard.sh") {
		t.Error("expected hook script path in output")
	}
	// Must contain an exact deny rule for a known irreversible command.
	if !strings.Contains(out, "Bash(canvas courses delete:*)") {
		t.Error("expected exact deny rule 'Bash(canvas courses delete:*)' in output")
	}
	// Must contain an exact ask rule for a known write command.
	if !strings.Contains(out, "Bash(canvas courses create:*)") {
		t.Error("expected exact ask rule 'Bash(canvas courses create:*)' in output")
	}
	// Must NOT have created any files (print mode only).
}

func TestEmitCanvasClaudeCode_DenyRulesExactPerCommand(t *testing.T) {
	root := buildTestTree()
	_, writes, irreversible := classifyCanvasCommands(root)
	g := canvasGuardPlan{irreversible: irreversible, writes: writes, allWrites: false}

	cmd := &cobra.Command{Use: "canvas"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := emitCanvasClaudeCode(cmd, g, false); err != nil {
		t.Fatalf("emitCanvasClaudeCode: %v", err)
	}
	out := buf.String()

	for _, gc := range irreversible {
		// Each blocked command must have an exact Bash rule.
		bashRule := "Bash(canvas " + gc.cli + ":*)"
		if !strings.Contains(out, bashRule) {
			t.Errorf("expected exact deny Bash rule %q in output", bashRule)
		}
		// Each blocked command must have an exact MCP tool name.
		mcpRule := "mcp__canvas__" + gc.tool
		if !strings.Contains(out, mcpRule) {
			t.Errorf("expected exact deny MCP rule %q in output", mcpRule)
		}
	}

	// Ensure NO regex-style MCP rules are present (the old pattern was mcp__.*canvas.*_verb).
	if strings.Contains(out, "mcp__.*canvas.*_") {
		t.Error("output must not contain regex-style MCP rule (mcp__.*canvas.*_)")
	}
}

func TestEmitCanvasClaudeCode_AskRulesExactPerCommand(t *testing.T) {
	root := buildTestTree()
	_, writes, irreversible := classifyCanvasCommands(root)
	g := canvasGuardPlan{irreversible: irreversible, writes: writes, allWrites: false}

	cmd := &cobra.Command{Use: "canvas"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := emitCanvasClaudeCode(cmd, g, false); err != nil {
		t.Fatalf("emitCanvasClaudeCode: %v", err)
	}
	out := buf.String()

	for _, gc := range writes {
		bashRule := "Bash(canvas " + gc.cli + ":*)"
		if !strings.Contains(out, bashRule) {
			t.Errorf("expected exact ask Bash rule %q in output", bashRule)
		}
		mcpRule := "mcp__canvas__" + gc.tool
		if !strings.Contains(out, mcpRule) {
			t.Errorf("expected exact ask MCP rule %q in output", mcpRule)
		}
	}
}

func TestEmitCanvasClaudeCode_AllWritesMode(t *testing.T) {
	root := buildTestTree()
	_, writes, irreversible := classifyCanvasCommands(root)
	g := canvasGuardPlan{irreversible: irreversible, writes: writes, allWrites: true}

	cmd := &cobra.Command{Use: "canvas"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := emitCanvasClaudeCode(cmd, g, false); err != nil {
		t.Fatalf("emitCanvasClaudeCode: %v", err)
	}
	out := buf.String()

	// All writes should now appear in deny, not ask.
	for _, gc := range writes {
		denyRule := "Bash(canvas " + gc.cli + ":*)"
		if !strings.Contains(out, denyRule) {
			t.Errorf("--all-writes: expected %q in deny section", denyRule)
		}
	}
	// The ask section should be an empty JSON array when --all-writes is active.
	if !strings.Contains(out, `"ask": []`) {
		t.Error("--all-writes: expected empty ask array in JSON output")
	}
}

// TestEmitCanvasClaudeCode_PATCHInDenyRules verifies PATCH is in the api deny rules.
func TestEmitCanvasClaudeCode_PATCHInDenyRules(t *testing.T) {
	root := buildTestTree()
	_, writes, irreversible := classifyCanvasCommands(root)
	g := canvasGuardPlan{irreversible: irreversible, writes: writes, allWrites: false}

	cmd := &cobra.Command{Use: "canvas"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := emitCanvasClaudeCode(cmd, g, false); err != nil {
		t.Fatalf("emitCanvasClaudeCode: %v", err)
	}
	out := buf.String()

	for _, method := range []string{"DELETE", "PUT", "POST", "PATCH"} {
		rule := "Bash(canvas api " + method + ":*)"
		if !strings.Contains(out, rule) {
			t.Errorf("expected deny rule %q for raw-api method in output", rule)
		}
	}
}

func TestEmitCanvasCodex_PrintsWithoutWrite(t *testing.T) {
	root := buildTestTree()
	_, writes, irreversible := classifyCanvasCommands(root)
	g := canvasGuardPlan{irreversible: irreversible, writes: writes, allWrites: false}

	cmd := &cobra.Command{Use: "canvas"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := emitCanvasCodex(cmd, g, false); err != nil {
		t.Fatalf("emitCanvasCodex: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "config.toml") {
		t.Error("expected config.toml reference in codex output")
	}
	if !strings.Contains(out, "read-only") {
		t.Error("expected read-only sandbox mode in codex output")
	}
	if !strings.Contains(out, "untrusted") {
		t.Error("expected untrusted approval policy in codex output")
	}
	// The output must honestly document that Codex only approval-gates.
	if !strings.Contains(out, "approval") {
		t.Error("expected approval-gate documentation in codex output")
	}
}

func TestEmitCanvasOpenCode_PrintsWithoutWrite(t *testing.T) {
	root := buildTestTree()
	_, writes, irreversible := classifyCanvasCommands(root)
	g := canvasGuardPlan{irreversible: irreversible, writes: writes, allWrites: false}

	cmd := &cobra.Command{Use: "canvas"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := emitCanvasOpenCode(cmd, g, false); err != nil {
		t.Fatalf("emitCanvasOpenCode: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "opencode.json") {
		t.Error("expected opencode.json reference in opencode output")
	}
	if !strings.Contains(out, "deny") {
		t.Error("expected 'deny' in opencode output")
	}
	// Must use exact per-command rules, not trailing-glob verb rules.
	if !strings.Contains(out, "canvas courses delete") {
		t.Error("expected exact 'canvas courses delete' in opencode deny rules")
	}
	// Must NOT use glob verb patterns like "canvas * delete*".
	if strings.Contains(out, "canvas * delete") {
		t.Error("opencode output must not use glob verb patterns like 'canvas * delete'")
	}
}

// --- Write-path tests ---

// TestEmitCanvasClaudeCode_WriteCreatesFiles verifies that --write installs
// the hook and settings.json under the project root, respects never-overwrite,
// and resolves the root correctly given a .git marker.
func TestEmitCanvasClaudeCode_WriteCreatesFiles(t *testing.T) {
	// Set up a temp directory as a fake project root with a .git dir so
	// findProjectRoot walks up to it correctly.
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.Mkdir(gitDir, 0o750); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	// Change CWD into a sub-directory of tmpDir to verify root-walking.
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0o750); err != nil {
		t.Fatalf("mkdir subdir: %v", err)
	}
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(subDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	root := buildTestTree()
	_, writes, irreversible := classifyCanvasCommands(root)
	g := canvasGuardPlan{irreversible: irreversible, writes: writes, allWrites: false}

	cmd := &cobra.Command{Use: "canvas"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := emitCanvasClaudeCode(cmd, g, true); err != nil {
		t.Fatalf("emitCanvasClaudeCode(write=true): %v", err)
	}

	// Both files must be written under the project root (not the subdir).
	hookFile := filepath.Join(tmpDir, ".claude", "hooks", "canvas-guard.sh")
	settingsFile := filepath.Join(tmpDir, ".claude", "settings.json")

	if _, err := os.Stat(hookFile); os.IsNotExist(err) {
		t.Errorf("hook script was not created at %s", hookFile)
	}
	if _, err := os.Stat(settingsFile); os.IsNotExist(err) {
		t.Errorf("settings.json was not created at %s", settingsFile)
	}

	// Second call must NOT overwrite existing files.
	firstContent, _ := os.ReadFile(settingsFile)
	var buf2 bytes.Buffer
	cmd2 := &cobra.Command{Use: "canvas"}
	cmd2.SetOut(&buf2)
	if err := emitCanvasClaudeCode(cmd2, g, true); err != nil {
		t.Fatalf("second emitCanvasClaudeCode(write=true): %v", err)
	}
	secondContent, _ := os.ReadFile(settingsFile)
	if string(firstContent) != string(secondContent) {
		t.Error("second write call must not overwrite existing settings.json")
	}
	// Second call output should mention "already exists".
	if !strings.Contains(buf2.String(), "already exists") {
		t.Error("second write call should print 'already exists' for existing files")
	}

	// Hook script must be executable. Windows has no Unix executable bit
	// (os.Stat reports -rw-rw-rw-), so this check only applies on POSIX hosts.
	if runtime.GOOS != "windows" {
		fi, err := os.Stat(hookFile)
		if err != nil {
			t.Fatalf("stat hook file: %v", err)
		}
		if fi.Mode()&0o111 == 0 {
			t.Errorf("hook script must be executable, got mode %v", fi.Mode())
		}
	}
}

// --- newAgentGuardCmd (command wiring) ---

func TestNewAgentGuardCmd_MissingHost(t *testing.T) {
	// Run the guard subcommand via a fresh root so classifyCanvasCommands
	// doesn't walk the real global tree. We invoke via the parent root so
	// cobra propagates the RunE error back through Execute().
	fresh := &cobra.Command{Use: "canvas", SilenceErrors: true, SilenceUsage: true}
	guard := newAgentGuardCmd()
	fresh.AddCommand(guard)
	var buf bytes.Buffer
	fresh.SetOut(&buf)
	fresh.SetErr(&buf)
	fresh.SetArgs([]string{"guard"}) // no --host
	err := fresh.Execute()
	if err == nil {
		t.Error("expected error when --host is not provided")
	}
}

func TestNewAgentGuardCmd_UnknownHost(t *testing.T) {
	fresh := &cobra.Command{Use: "canvas", SilenceErrors: true, SilenceUsage: true}
	guard := newAgentGuardCmd()
	fresh.AddCommand(guard)
	var buf bytes.Buffer
	fresh.SetOut(&buf)
	fresh.SetErr(&buf)
	fresh.SetArgs([]string{"guard", "--host", "notahost"})
	err := fresh.Execute()
	if err == nil {
		t.Error("expected error for unknown host")
	}
	if !strings.Contains(err.Error(), "notahost") {
		t.Errorf("error should mention the unknown host, got: %v", err)
	}
}

// --- Hook script content ---

func TestCanvasHookScript_ContainsBlockedPaths(t *testing.T) {
	cmds := []canvasGuardCmd{
		{cli: "courses delete", tool: "canvas_courses_delete", verb: "delete"},
		{cli: "sections crosslist", tool: "canvas_sections_crosslist", verb: "crosslist"},
	}
	script := canvasHookScript(cmds)

	if !strings.Contains(script, "courses delete") {
		t.Error("hook script should contain cli path 'courses delete'")
	}
	if !strings.Contains(script, "sections crosslist") {
		t.Error("hook script should contain cli path 'sections crosslist'")
	}
	if !strings.Contains(script, "mcp__canvas__canvas_courses_delete") {
		t.Error("hook script should contain exact MCP tool name")
	}
	if !strings.Contains(script, "canvas agent guard") {
		t.Error("hook script should reference canvas agent guard")
	}
	if !strings.Contains(script, "canvas api") {
		t.Error("hook script should mention canvas api escape hatch")
	}
	if !strings.Contains(script, "DELETE|PUT|POST|PATCH") {
		t.Error("hook script should handle raw HTTP methods including PATCH")
	}
	// Verify the script is valid bash (starts with shebang).
	if !strings.HasPrefix(script, "#!/usr/bin/env bash") {
		t.Error("hook script should start with bash shebang")
	}
	// Must NOT use bare verb regex matching anywhere in the jq branch.
	if strings.Contains(script, `grep -qiE "\bcanvas\b.*\b`) {
		t.Error("hook script must not use bare-verb anywhere-on-line grep")
	}
}

func TestCanvasHookScript_NoRegexVerbStyle(t *testing.T) {
	// Verify that the hook uses path-anchored matching, not bare-verb grep.
	cmds := []canvasGuardCmd{
		{cli: "courses delete", tool: "canvas_courses_delete", verb: "delete"},
	}
	script := canvasHookScript(cmds)

	// The old verb-group pattern should be absent.
	if strings.Contains(script, "verbs=") {
		t.Error("hook script must not contain bare 'verbs=' variable")
	}
	// Exact MCP tool list should be present.
	if !strings.Contains(script, "blocked_tools=") {
		t.Error("hook script must contain blocked_tools array")
	}
	// Exact cli path list should be present.
	if !strings.Contains(script, "blocked_cmds=") {
		t.Error("hook script must contain blocked_cmds array")
	}
}

// --- RealTree smoke test (against actual rootCmd) ---

func TestClassifyCanvasCommands_RealTree_SanityCheck(t *testing.T) {
	read, writes, irreversible := classifyCanvasCommands(rootCmd)

	total := len(read) + len(writes) + len(irreversible)
	if total == 0 {
		t.Fatal("expected at least some classified commands from the real tree")
	}

	// Verify expected commands exist in the correct buckets.
	assertContainsPath(t, "read", read, "courses list")
	assertContainsPath(t, "write", writes, "courses create")
	assertContainsPath(t, "irreversible", irreversible, "courses delete")

	// Verify that cancel/close/merge/split are hard-blocked.
	for _, path := range []string{"users merge", "sections merge"} {
		// These may or may not exist; if they do, they must be in irreversible.
		for _, gc := range append(read, writes...) {
			if gc.cli == path {
				t.Errorf("command %q must be in irreversible, found in read/writes", path)
			}
		}
	}

	// Check counts are in a sane range (these will grow as the CLI grows).
	if len(irreversible) < 5 {
		t.Errorf("expected at least 5 irreversible commands, got %d", len(irreversible))
	}
	if len(writes) < 10 {
		t.Errorf("expected at least 10 write commands, got %d", len(writes))
	}
	if len(read) < 10 {
		t.Errorf("expected at least 10 read commands, got %d", len(read))
	}

	// Verify local utility groups are excluded.
	for _, gc := range append(append(read, writes...), irreversible...) {
		parts := strings.SplitN(gc.cli, " ", 2)
		if len(parts) > 0 && canvasLocalGroups[parts[0]] {
			t.Errorf("local group command leaked through: %q", gc.cli)
		}
	}

	// Verify "api" command is excluded.
	for _, gc := range append(append(read, writes...), irreversible...) {
		if strings.HasPrefix(gc.cli, "api ") || gc.cli == "api" {
			t.Errorf("raw api command leaked through: %q", gc.cli)
		}
	}

	t.Logf("Real tree: %d read, %d write, %d irreversible (total %d)", len(read), len(writes), len(irreversible), total)
}

// --- Helpers ---

func cliPaths(cs []canvasGuardCmd) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.cli
	}
	return out
}

func assertContainsPath(t *testing.T, bucket string, cs []canvasGuardCmd, path string) {
	t.Helper()
	for _, c := range cs {
		if c.cli == path {
			return
		}
	}
	t.Errorf("%s bucket: expected %q but not found; have: %v", bucket, path, cliPaths(cs))
}

func assertNotContainsPath(t *testing.T, bucket string, cs []canvasGuardCmd, path string) {
	t.Helper()
	for _, c := range cs {
		if c.cli == path {
			t.Errorf("%s bucket: expected %q to be excluded, but it was present", bucket, path)
			return
		}
	}
}

// TestClassifyCanvasCommand_APIGetIsRead pins the #60 exception: the GET-only
// "api get" sibling classifies as a read (so it can advertise readOnlyHint and
// survive read-only MCP clients), while the general "api" escape hatch stays
// skipped/unannotated.
func TestClassifyCanvasCommand_APIGetIsRead(t *testing.T) {
	root := rootCmd

	apiCmd := findCmd(root, "api")
	if apiCmd == nil {
		t.Fatal(`"canvas api" not found`)
	}
	if class, _ := classifyCanvasCommand(root, apiCmd); class != canvasClassSkip {
		t.Errorf(`"canvas api" must be canvasClassSkip, got %v`, class)
	}

	apiGet := findCmd(root, "api get")
	if apiGet == nil {
		t.Fatal(`"canvas api get" not found`)
	}
	if class, _ := classifyCanvasCommand(root, apiGet); class != canvasClassRead {
		t.Errorf(`"canvas api get" must be canvasClassRead, got %v`, class)
	}
}
