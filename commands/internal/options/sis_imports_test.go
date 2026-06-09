package options

import "testing"

func TestSISImportsListOptions_Validate(t *testing.T) {
	opts := &SISImportsListOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("SISImportsListOptions.Validate() error = %v, want nil", err)
	}
}

func TestSISImportsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *SISImportsGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &SISImportsGetOptions{AccountID: 1, ImportID: 2},
			wantErr: false,
		},
		{
			name:    "missing account ID",
			opts:    &SISImportsGetOptions{AccountID: 0, ImportID: 2},
			wantErr: true,
		},
		{
			name:    "missing import ID",
			opts:    &SISImportsGetOptions{AccountID: 1, ImportID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SISImportsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSISImportsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *SISImportsCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &SISImportsCreateOptions{AccountID: 1, FilePath: "/tmp/data.csv"},
			wantErr: false,
		},
		{
			name:    "missing account ID",
			opts:    &SISImportsCreateOptions{AccountID: 0, FilePath: "/tmp/data.csv"},
			wantErr: true,
		},
		{
			name:    "missing file path",
			opts:    &SISImportsCreateOptions{AccountID: 1, FilePath: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SISImportsCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSISImportsAbortOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *SISImportsAbortOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &SISImportsAbortOptions{AccountID: 1, ImportID: 2},
			wantErr: false,
		},
		{
			name:    "missing account ID",
			opts:    &SISImportsAbortOptions{AccountID: 0, ImportID: 2},
			wantErr: true,
		},
		{
			name:    "missing import ID",
			opts:    &SISImportsAbortOptions{AccountID: 1, ImportID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SISImportsAbortOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSISImportsRestoreOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *SISImportsRestoreOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &SISImportsRestoreOptions{AccountID: 1, ImportID: 2},
			wantErr: false,
		},
		{
			name:    "missing account ID",
			opts:    &SISImportsRestoreOptions{AccountID: 0, ImportID: 2},
			wantErr: true,
		},
		{
			name:    "missing import ID",
			opts:    &SISImportsRestoreOptions{AccountID: 1, ImportID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SISImportsRestoreOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSISImportsErrorsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *SISImportsErrorsOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &SISImportsErrorsOptions{AccountID: 1, ImportID: 2},
			wantErr: false,
		},
		{
			name:    "missing account ID",
			opts:    &SISImportsErrorsOptions{AccountID: 0, ImportID: 2},
			wantErr: true,
		},
		{
			name:    "missing import ID",
			opts:    &SISImportsErrorsOptions{AccountID: 1, ImportID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SISImportsErrorsOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
