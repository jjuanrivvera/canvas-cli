package options

import (
	"testing"
)

func TestAliasSetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AliasSetOptions
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid options",
			opts: &AliasSetOptions{
				Name:      "ca",
				Expansion: "assignments list --course-id 123",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			opts: &AliasSetOptions{
				Name:      "",
				Expansion: "assignments list",
			},
			wantErr: true,
			errMsg:  "alias name is required",
		},
		{
			name: "empty expansion",
			opts: &AliasSetOptions{
				Name:      "ca",
				Expansion: "",
			},
			wantErr: true,
			errMsg:  "alias expansion is required",
		},
		{
			name: "name with spaces",
			opts: &AliasSetOptions{
				Name:      "my alias",
				Expansion: "assignments list",
			},
			wantErr: true,
			errMsg:  "alias name cannot contain spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AliasSetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("AliasSetOptions.Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestAliasListOptions_Validate(t *testing.T) {
	opts := &AliasListOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("AliasListOptions.Validate() error = %v, want nil", err)
	}
}

func TestAliasDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AliasDeleteOptions
		wantErr bool
	}{
		{
			name: "valid options",
			opts: &AliasDeleteOptions{
				Name: "ca",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			opts: &AliasDeleteOptions{
				Name: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AliasDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
