package options

import "testing"

func TestConfigListOptions_Validate(t *testing.T) {
	opts := &ConfigListOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("ConfigListOptions.Validate() error = %v, want nil", err)
	}
}

func TestConfigAddOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ConfigAddOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ConfigAddOptions{Name: "myschool", URL: "https://canvas.example.com"},
			wantErr: false,
		},
		{
			name:    "missing name",
			opts:    &ConfigAddOptions{Name: "", URL: "https://canvas.example.com"},
			wantErr: true,
		},
		{
			name:    "missing URL",
			opts:    &ConfigAddOptions{Name: "myschool", URL: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigAddOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigUseOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ConfigUseOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ConfigUseOptions{InstanceName: "myschool"},
			wantErr: false,
		},
		{
			name:    "missing instance name",
			opts:    &ConfigUseOptions{InstanceName: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigUseOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigRemoveOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ConfigRemoveOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ConfigRemoveOptions{InstanceName: "myschool"},
			wantErr: false,
		},
		{
			name:    "missing instance name",
			opts:    &ConfigRemoveOptions{InstanceName: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigRemoveOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigShowOptions_Validate(t *testing.T) {
	opts := &ConfigShowOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("ConfigShowOptions.Validate() error = %v, want nil", err)
	}
}

func TestConfigAccountOptions_Validate(t *testing.T) {
	opts := &ConfigAccountOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("ConfigAccountOptions.Validate() error = %v, want nil", err)
	}
}
