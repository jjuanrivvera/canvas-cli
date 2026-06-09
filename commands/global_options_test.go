package commands

import "testing"

func TestGetGlobalOptions_Defaults(t *testing.T) {
	// Save and restore global state
	origVerbose := verbose
	origQuiet := quiet
	origNoCache := noCache
	origDryRun := dryRun
	origOutputFormat := outputFormat
	defer func() {
		verbose = origVerbose
		quiet = origQuiet
		noCache = origNoCache
		dryRun = origDryRun
		outputFormat = origOutputFormat
	}()

	// Set known values
	verbose = false
	quiet = false
	noCache = false
	dryRun = false
	outputFormat = "table"

	opts := GetGlobalOptions()
	if opts == nil {
		t.Fatal("GetGlobalOptions returned nil")
	}
	if opts.Verbose != false {
		t.Errorf("expected Verbose=false, got %v", opts.Verbose)
	}
	if opts.Quiet != false {
		t.Errorf("expected Quiet=false, got %v", opts.Quiet)
	}
	if opts.NoCache != false {
		t.Errorf("expected NoCache=false, got %v", opts.NoCache)
	}
	if opts.DryRun != false {
		t.Errorf("expected DryRun=false, got %v", opts.DryRun)
	}
	if opts.OutputFormat != "table" {
		t.Errorf("expected OutputFormat=table, got %q", opts.OutputFormat)
	}
}

func TestGetGlobalOptions_WithValues(t *testing.T) {
	origVerbose := verbose
	origQuiet := quiet
	origNoCache := noCache
	origDryRun := dryRun
	origLimit := globalLimit
	origAsUser := asUserID
	origFilter := filterText
	origSort := sortField
	defer func() {
		verbose = origVerbose
		quiet = origQuiet
		noCache = origNoCache
		dryRun = origDryRun
		globalLimit = origLimit
		asUserID = origAsUser
		filterText = origFilter
		sortField = origSort
	}()

	verbose = true
	quiet = true
	noCache = true
	dryRun = true
	globalLimit = 50
	asUserID = 99
	filterText = "test"
	sortField = "-name"

	opts := GetGlobalOptions()
	if !opts.Verbose {
		t.Error("expected Verbose=true")
	}
	if !opts.Quiet {
		t.Error("expected Quiet=true")
	}
	if !opts.NoCache {
		t.Error("expected NoCache=true")
	}
	if !opts.DryRun {
		t.Error("expected DryRun=true")
	}
	if opts.Limit != 50 {
		t.Errorf("expected Limit=50, got %d", opts.Limit)
	}
	if opts.AsUserID != 99 {
		t.Errorf("expected AsUserID=99, got %d", opts.AsUserID)
	}
	if opts.FilterText != "test" {
		t.Errorf("expected FilterText=test, got %q", opts.FilterText)
	}
	if opts.SortField != "-name" {
		t.Errorf("expected SortField=-name, got %q", opts.SortField)
	}
}

func TestGetGlobalOptions_ReturnsCopy(t *testing.T) {
	// Each call should return a new struct snapshot
	a := GetGlobalOptions()
	b := GetGlobalOptions()
	if a == b {
		t.Error("expected distinct pointer each call")
	}
}
