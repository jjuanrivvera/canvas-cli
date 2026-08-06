package commands

import (
	"testing"

	"github.com/njayp/ophis"
	"github.com/spf13/cobra"
)

// TestMCPToolSelector covers the --readonly filter (#60): outside read-only mode
// every command is admitted; inside it, only commands annotated readOnlyHint=true
// survive, which drops writes and the unannotated "api" tool while keeping
// "api get".
func TestMCPToolSelector(t *testing.T) {
	readCmd := &cobra.Command{Use: "get", Annotations: map[string]string{ophis.AnnotationReadOnly: "true"}}
	writeCmd := &cobra.Command{Use: "create", Annotations: map[string]string{ophis.AnnotationReadOnly: "false"}}
	unannotated := &cobra.Command{Use: "api"} // the escape hatch: no annotations

	orig := mcpReadonly
	defer func() { mcpReadonly = orig }()

	mcpReadonly = false
	for _, c := range []*cobra.Command{readCmd, writeCmd, unannotated} {
		if !mcpToolSelector(c) {
			t.Errorf("normal mode: %q should be admitted", c.Name())
		}
	}

	mcpReadonly = true
	if !mcpToolSelector(readCmd) {
		t.Error("readonly mode: read-only command must be admitted")
	}
	if mcpToolSelector(writeCmd) {
		t.Error("readonly mode: write command must be excluded")
	}
	if mcpToolSelector(unannotated) {
		t.Error(`readonly mode: unannotated "api" command must be excluded`)
	}
}
