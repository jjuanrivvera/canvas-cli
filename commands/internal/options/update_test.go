package options

import "testing"

func TestUpdateOptions_Validate(t *testing.T) {
	opts := &UpdateOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("UpdateOptions.Validate() error = %v, want nil", err)
	}
}

func TestUpdateCheckOptions_Validate(t *testing.T) {
	opts := &UpdateCheckOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("UpdateCheckOptions.Validate() error = %v, want nil", err)
	}
}

func TestUpdateEnableOptions_Validate(t *testing.T) {
	tests := []struct {
		name         string
		opts         *UpdateEnableOptions
		wantErr      bool
		wantInterval int
	}{
		{
			name:         "valid positive interval",
			opts:         &UpdateEnableOptions{Interval: 30},
			wantErr:      false,
			wantInterval: 30,
		},
		{
			name:         "zero interval defaults to 60",
			opts:         &UpdateEnableOptions{Interval: 0},
			wantErr:      false,
			wantInterval: 60,
		},
		{
			name:         "negative interval defaults to 60",
			opts:         &UpdateEnableOptions{Interval: -5},
			wantErr:      false,
			wantInterval: 60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateEnableOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.opts.Interval != tt.wantInterval {
				t.Errorf("UpdateEnableOptions.Validate() Interval = %d, want %d", tt.opts.Interval, tt.wantInterval)
			}
		})
	}
}

func TestUpdateDisableOptions_Validate(t *testing.T) {
	opts := &UpdateDisableOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("UpdateDisableOptions.Validate() error = %v, want nil", err)
	}
}

func TestUpdateStatusOptions_Validate(t *testing.T) {
	opts := &UpdateStatusOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("UpdateStatusOptions.Validate() error = %v, want nil", err)
	}
}
