package commands

import "testing"

func TestRootQuietFlagHasNoShorthand(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("quiet")
	if flag == nil {
		t.Fatal("expected quiet flag to exist on root command")
		return
	}

	if flag.Shorthand != "" {
		t.Fatalf("expected --quiet to have no shorthand, got %q", flag.Shorthand)
	}
}

func TestAPIQueryFlagKeepsQShorthand(t *testing.T) {
	flag := apiCmd.Flags().Lookup("query")
	if flag == nil {
		t.Fatal("expected query flag to exist on api command")
		return
	}

	if flag.Shorthand != "q" {
		t.Fatalf("expected --query shorthand to be q, got %q", flag.Shorthand)
	}
}
