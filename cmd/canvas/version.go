package main

import (
	"fmt"
	"runtime"
)

// GetVersionInfo returns formatted version information
func GetVersionInfo(version, commit, buildDate string) string {
	return fmt.Sprintf("canvas-cli %s\nCommit: %s\nBuilt: %s\nGo: %s\nOS/Arch: %s/%s",
		version,
		commit,
		buildDate,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}
