package commands

import (
	"os"
	"testing"
)

// TestMain disables the background auto-updater for the whole test binary.
//
// rootCmd's PersistentPreRun launches a real, asynchronous update check
// (RunUpdateCheckAsync). Tests that dispatch through rootCmd execute it
// repeatedly; because the async check makes a network call and PersistentPostRun
// only waits 5s, a slow check outlives the command and the next execution races
// with it on the shared updater state — which fails under `go test -race`. A
// real CLI process runs rootCmd exactly once, so this only affects the harness.
func TestMain(m *testing.M) {
	disableAutoUpdate = true
	os.Exit(m.Run())
}
