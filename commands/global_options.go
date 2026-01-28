package commands

// GlobalOptions wraps all root-level flags into a single struct.
// Existing global vars are kept for backward compatibility; this struct
// provides a clean accessor for new and migrated commands.
type GlobalOptions struct {
	ConfigFile    string
	InstanceURL   string
	OutputFormat  string
	Verbose       bool
	Quiet         bool
	NoCache       bool
	AsUserID      int64
	Limit         int
	DryRun        bool
	ShowToken     bool
	FilterText    string
	FilterColumns []string
	SortField     string
}

// GetGlobalOptions returns a snapshot of the current global flag values.
func GetGlobalOptions() *GlobalOptions {
	return &GlobalOptions{
		ConfigFile:    cfgFile,
		InstanceURL:   instanceURL,
		OutputFormat:  outputFormat,
		Verbose:       verbose,
		Quiet:         quiet,
		NoCache:       noCache,
		AsUserID:      asUserID,
		Limit:         globalLimit,
		DryRun:        dryRun,
		ShowToken:     showToken,
		FilterText:    filterText,
		FilterColumns: filterColumns,
		SortField:     sortField,
	}
}
