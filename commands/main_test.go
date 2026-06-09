package commands

import (
	"os"
	"testing"
	"time"
)

// TestMain disables the background auto-updater for the whole test binary.
//
// rootCmd's PersistentPreRun launches a real, asynchronous update check
// (RunUpdateCheckAsync). Tests that dispatch through rootCmd execute it
// repeatedly; the async check makes a network call and PersistentPostRun only
// waits 5s, so a slow check outlives the command and the next execution races
// with it on the shared updater state — failing under `go test -race`. A real
// CLI process runs rootCmd exactly once, so this only affects the harness.
//
// (Config/cache isolation is handled per-test by the setup*TestHome helpers,
// which set both HOME and USERPROFILE so os.UserHomeDir() is redirected on every
// platform. Isolating globally here is avoided because some tests, e.g. the
// doctor command, legitimately inspect the real environment.)
func TestMain(m *testing.M) {
	disableAutoUpdate = true
	// Use a near-zero retry backoff so retryable-error tests don't each sleep
	// through the default 1s/2s/4s exponential backoff (~7s/test).
	testRetryBackoff = time.Millisecond
	os.Exit(m.Run())
}
