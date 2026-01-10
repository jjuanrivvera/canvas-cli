package config

import (
	"strings"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: &Config{
				DefaultInstance: "prod",
				Instances: map[string]*Instance{
					"prod": {
						Name: "prod",
						URL:  "https://canvas.instructure.com",
					},
				},
				Settings: DefaultSettings(),
			},
			wantErr: false,
		},
		{
			name: "valid config without default instance",
			config: &Config{
				Instances: map[string]*Instance{
					"prod": {
						Name: "prod",
						URL:  "https://canvas.instructure.com",
					},
				},
				Settings: DefaultSettings(),
			},
			wantErr: false,
		},
		{
			name: "default instance does not exist",
			config: &Config{
				DefaultInstance: "nonexistent",
				Instances: map[string]*Instance{
					"prod": {
						Name: "prod",
						URL:  "https://canvas.instructure.com",
					},
				},
				Settings: DefaultSettings(),
			},
			wantErr: true,
			errMsg:  "default instance \"nonexistent\" does not exist",
		},
		{
			name: "invalid instance",
			config: &Config{
				Instances: map[string]*Instance{
					"bad": {
						Name: "",
						URL:  "https://canvas.instructure.com",
					},
				},
				Settings: DefaultSettings(),
			},
			wantErr: true,
			errMsg:  "instance \"bad\" is invalid",
		},
		{
			name: "invalid settings",
			config: &Config{
				Instances: map[string]*Instance{
					"prod": {
						Name: "prod",
						URL:  "https://canvas.instructure.com",
					},
				},
				Settings: &Settings{
					DefaultOutputFormat: "invalid",
					RequestsPerSecond:   10,
					CacheTTL:            60,
					LogLevel:            "info",
				},
			},
			wantErr: true,
			errMsg:  "settings are invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateInstance(t *testing.T) {
	tests := []struct {
		name     string
		instance *Instance
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid instance",
			instance: &Instance{
				Name: "production",
				URL:  "https://canvas.instructure.com",
			},
			wantErr: false,
		},
		{
			name: "valid instance with http",
			instance: &Instance{
				Name: "local",
				URL:  "http://localhost:3000",
			},
			wantErr: false,
		},
		{
			name:     "nil instance",
			instance: nil,
			wantErr:  true,
			errMsg:   "instance cannot be nil",
		},
		{
			name: "empty name",
			instance: &Instance{
				Name: "",
				URL:  "https://canvas.instructure.com",
			},
			wantErr: true,
			errMsg:  "instance name is required",
		},
		{
			name: "name too long",
			instance: &Instance{
				Name: strings.Repeat("a", 101),
				URL:  "https://canvas.instructure.com",
			},
			wantErr: true,
			errMsg:  "instance name is too long",
		},
		{
			name: "name with invalid characters - slash",
			instance: &Instance{
				Name: "prod/test",
				URL:  "https://canvas.instructure.com",
			},
			wantErr: true,
			errMsg:  "instance name contains invalid characters",
		},
		{
			name: "name with invalid characters - backslash",
			instance: &Instance{
				Name: "prod\\test",
				URL:  "https://canvas.instructure.com",
			},
			wantErr: true,
			errMsg:  "instance name contains invalid characters",
		},
		{
			name: "name with invalid characters - colon",
			instance: &Instance{
				Name: "prod:test",
				URL:  "https://canvas.instructure.com",
			},
			wantErr: true,
			errMsg:  "instance name contains invalid characters",
		},
		{
			name: "name with invalid characters - asterisk",
			instance: &Instance{
				Name: "prod*test",
				URL:  "https://canvas.instructure.com",
			},
			wantErr: true,
			errMsg:  "instance name contains invalid characters",
		},
		{
			name: "name with invalid characters - question mark",
			instance: &Instance{
				Name: "prod?test",
				URL:  "https://canvas.instructure.com",
			},
			wantErr: true,
			errMsg:  "instance name contains invalid characters",
		},
		{
			name: "name with invalid characters - quotes",
			instance: &Instance{
				Name: "prod\"test",
				URL:  "https://canvas.instructure.com",
			},
			wantErr: true,
			errMsg:  "instance name contains invalid characters",
		},
		{
			name: "name with invalid characters - angle brackets",
			instance: &Instance{
				Name: "prod<test>",
				URL:  "https://canvas.instructure.com",
			},
			wantErr: true,
			errMsg:  "instance name contains invalid characters",
		},
		{
			name: "name with invalid characters - pipe",
			instance: &Instance{
				Name: "prod|test",
				URL:  "https://canvas.instructure.com",
			},
			wantErr: true,
			errMsg:  "instance name contains invalid characters",
		},
		{
			name: "empty URL",
			instance: &Instance{
				Name: "prod",
				URL:  "",
			},
			wantErr: true,
			errMsg:  "instance URL is required",
		},
		{
			name: "invalid URL",
			instance: &Instance{
				Name: "prod",
				URL:  "://invalid",
			},
			wantErr: true,
			errMsg:  "invalid URL",
		},
		{
			name: "URL with invalid scheme - ftp",
			instance: &Instance{
				Name: "prod",
				URL:  "ftp://canvas.instructure.com",
			},
			wantErr: true,
			errMsg:  "URL must use http or https scheme",
		},
		{
			name: "URL without host",
			instance: &Instance{
				Name: "prod",
				URL:  "https://",
			},
			wantErr: true,
			errMsg:  "URL must have a host",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateInstance(tt.instance)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateSettings(t *testing.T) {
	tests := []struct {
		name     string
		settings *Settings
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid settings - table",
			settings: DefaultSettings(),
			wantErr:  false,
		},
		{
			name: "valid settings - json",
			settings: &Settings{
				DefaultOutputFormat: "json",
				RequestsPerSecond:   10,
				CacheTTL:            60,
				CacheEnabled:        true,
				LogLevel:            "info",
			},
			wantErr: false,
		},
		{
			name: "valid settings - yaml",
			settings: &Settings{
				DefaultOutputFormat: "yaml",
				RequestsPerSecond:   10,
				CacheTTL:            60,
				CacheEnabled:        true,
				LogLevel:            "info",
			},
			wantErr: false,
		},
		{
			name: "valid settings - csv",
			settings: &Settings{
				DefaultOutputFormat: "csv",
				RequestsPerSecond:   10,
				CacheTTL:            60,
				CacheEnabled:        true,
				LogLevel:            "info",
			},
			wantErr: false,
		},
		{
			name: "valid settings - debug log level",
			settings: &Settings{
				DefaultOutputFormat: "table",
				RequestsPerSecond:   10,
				CacheTTL:            60,
				CacheEnabled:        true,
				LogLevel:            "debug",
			},
			wantErr: false,
		},
		{
			name: "valid settings - warn log level",
			settings: &Settings{
				DefaultOutputFormat: "table",
				RequestsPerSecond:   10,
				CacheTTL:            60,
				CacheEnabled:        true,
				LogLevel:            "warn",
			},
			wantErr: false,
		},
		{
			name: "valid settings - error log level",
			settings: &Settings{
				DefaultOutputFormat: "table",
				RequestsPerSecond:   10,
				CacheTTL:            60,
				CacheEnabled:        true,
				LogLevel:            "error",
			},
			wantErr: false,
		},
		{
			name: "valid settings - uppercase log level",
			settings: &Settings{
				DefaultOutputFormat: "table",
				RequestsPerSecond:   10,
				CacheTTL:            60,
				CacheEnabled:        true,
				LogLevel:            "INFO",
			},
			wantErr: false,
		},
		{
			name: "valid settings - cache disabled with negative TTL",
			settings: &Settings{
				DefaultOutputFormat: "table",
				RequestsPerSecond:   10,
				CacheTTL:            -1,
				CacheEnabled:        false,
				LogLevel:            "info",
			},
			wantErr: false,
		},
		{
			name:     "nil settings",
			settings: nil,
			wantErr:  true,
			errMsg:   "settings cannot be nil",
		},
		{
			name: "invalid output format",
			settings: &Settings{
				DefaultOutputFormat: "xml",
				RequestsPerSecond:   10,
				CacheTTL:            60,
				CacheEnabled:        true,
				LogLevel:            "info",
			},
			wantErr: true,
			errMsg:  "invalid output format",
		},
		{
			name: "requests per second zero",
			settings: &Settings{
				DefaultOutputFormat: "table",
				RequestsPerSecond:   0,
				CacheTTL:            60,
				CacheEnabled:        true,
				LogLevel:            "info",
			},
			wantErr: true,
			errMsg:  "requests_per_second must be positive",
		},
		{
			name: "requests per second negative",
			settings: &Settings{
				DefaultOutputFormat: "table",
				RequestsPerSecond:   -1,
				CacheTTL:            60,
				CacheEnabled:        true,
				LogLevel:            "info",
			},
			wantErr: true,
			errMsg:  "requests_per_second must be positive",
		},
		{
			name: "requests per second too high",
			settings: &Settings{
				DefaultOutputFormat: "table",
				RequestsPerSecond:   101,
				CacheTTL:            60,
				CacheEnabled:        true,
				LogLevel:            "info",
			},
			wantErr: true,
			errMsg:  "requests_per_second is too high",
		},
		{
			name: "cache enabled with negative TTL",
			settings: &Settings{
				DefaultOutputFormat: "table",
				RequestsPerSecond:   10,
				CacheTTL:            -1,
				CacheEnabled:        true,
				LogLevel:            "info",
			},
			wantErr: true,
			errMsg:  "cache_ttl_minutes cannot be negative",
		},
		{
			name: "invalid log level",
			settings: &Settings{
				DefaultOutputFormat: "table",
				RequestsPerSecond:   10,
				CacheTTL:            60,
				CacheEnabled:        true,
				LogLevel:            "trace",
			},
			wantErr: true,
			errMsg:  "invalid log level",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSettings(tt.settings)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestSanitizeInstanceName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "valid name unchanged",
			input: "production",
			want:  "production",
		},
		{
			name:  "removes slash",
			input: "prod/test",
			want:  "prod-test",
		},
		{
			name:  "removes backslash",
			input: "prod\\test",
			want:  "prod-test",
		},
		{
			name:  "removes colon",
			input: "prod:test",
			want:  "prod-test",
		},
		{
			name:  "removes asterisk",
			input: "prod*test",
			want:  "prod-test",
		},
		{
			name:  "removes question mark",
			input: "prod?test",
			want:  "prod-test",
		},
		{
			name:  "removes quotes",
			input: "prod\"test",
			want:  "prod-test",
		},
		{
			name:  "removes angle brackets",
			input: "prod<test>",
			want:  "prod-test-",
		},
		{
			name:  "removes pipe",
			input: "prod|test",
			want:  "prod-test",
		},
		{
			name:  "removes multiple invalid characters",
			input: "prod/test:123",
			want:  "prod-test-123",
		},
		{
			name:  "trims whitespace",
			input: "  production  ",
			want:  "production",
		},
		{
			name:  "limits length",
			input: strings.Repeat("a", 150),
			want:  strings.Repeat("a", 100),
		},
		{
			name:  "trims and sanitizes",
			input: "  prod/test  ",
			want:  "prod-test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeInstanceName(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeInstanceName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "adds https scheme",
			input:   "canvas.instructure.com",
			want:    "https://canvas.instructure.com",
			wantErr: false,
		},
		{
			name:    "keeps https scheme",
			input:   "https://canvas.instructure.com",
			want:    "https://canvas.instructure.com",
			wantErr: false,
		},
		{
			name:    "keeps http scheme",
			input:   "http://localhost:3000",
			want:    "http://localhost:3000",
			wantErr: false,
		},
		{
			name:    "removes trailing slash",
			input:   "https://canvas.instructure.com/",
			want:    "https://canvas.instructure.com",
			wantErr: false,
		},
		{
			name:    "removes path with trailing slash",
			input:   "https://canvas.instructure.com/api/",
			want:    "https://canvas.instructure.com/api",
			wantErr: false,
		},
		{
			name:    "keeps path without trailing slash",
			input:   "https://canvas.instructure.com/api",
			want:    "https://canvas.instructure.com/api",
			wantErr: false,
		},
		{
			name:    "removes single slash path",
			input:   "https://canvas.instructure.com/",
			want:    "https://canvas.instructure.com",
			wantErr: false,
		},
		{
			name:    "handles URL with port",
			input:   "localhost:3000",
			want:    "https://localhost:3000",
			wantErr: false,
		},
		{
			name:    "handles URL with http and port",
			input:   "http://localhost:3000",
			want:    "http://localhost:3000",
			wantErr: false,
		},
		{
			name:    "handles URL starting with colon-slash",
			input:   "://invalid",
			want:    "https://://invalid",
			wantErr: false,
		},
		{
			name:    "handles subdomain",
			input:   "test.instructure.com",
			want:    "https://test.instructure.com",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeURL(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if got != tt.want {
					t.Errorf("NormalizeURL(%q) = %q, want %q", tt.input, got, tt.want)
				}
			}
		})
	}
}
