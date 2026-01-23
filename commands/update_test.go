package commands

import "testing"

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		name     string
		latest   string
		current  string
		expected bool
	}{
		{"newer major", "2.0.0", "1.0.0", true},
		{"newer minor", "1.1.0", "1.0.0", true},
		{"newer patch", "1.0.1", "1.0.0", true},
		{"same version", "1.0.0", "1.0.0", false},
		{"older major", "1.0.0", "2.0.0", false},
		{"older minor", "1.0.0", "1.1.0", false},
		{"older patch", "1.0.0", "1.0.1", false},
		{"with v prefix latest", "v1.1.0", "1.0.0", true},
		{"with v prefix current", "1.1.0", "v1.0.0", true},
		{"with v prefix both", "v1.1.0", "v1.0.0", true},
		{"with prerelease", "1.1.0-beta", "1.0.0", true},
		{"dev version current", "1.1.0", "dev", false},
		{"empty current", "1.1.0", "", false},
		// Edge case: downgrade protection
		{"downgrade protection major", "1.0.0", "2.0.0", false},
		{"downgrade protection minor", "1.0.0", "1.1.0", false},
		{"downgrade protection patch", "1.0.0", "1.0.1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNewerVersion(tt.latest, tt.current)
			if result != tt.expected {
				t.Errorf("isNewerVersion(%q, %q) = %v, expected %v", tt.latest, tt.current, result, tt.expected)
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected [3]int
	}{
		{"1.2.3", [3]int{1, 2, 3}},
		{"v1.2.3", [3]int{1, 2, 3}},
		{"1.2", [3]int{1, 2, 0}},
		{"1", [3]int{1, 0, 0}},
		{"1.2.3-beta", [3]int{1, 2, 3}},
		{"1.2.3-rc.1", [3]int{1, 2, 3}},
		{"0.0.0", [3]int{0, 0, 0}},
		{"10.20.30", [3]int{10, 20, 30}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseVersion(tt.input)
			if result != tt.expected {
				t.Errorf("parseVersion(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
