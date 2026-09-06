package commands

import "testing"

// withGlobalDryRun sets the dryRun global (the --dry-run persistent flag)
// and restores it via t.Cleanup. Command constructors under test are not
// attached to rootCmd, so the flag itself is not available to them.
func withGlobalDryRun(t *testing.T, value bool) {
	t.Helper()
	orig := dryRun
	dryRun = value
	t.Cleanup(func() { dryRun = orig })
}
