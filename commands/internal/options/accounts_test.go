package options

import "testing"

func TestAccountsListOptions_Validate(t *testing.T) {
	opts := &AccountsListOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("AccountsListOptions.Validate() error = %v, want nil", err)
	}
}

func TestAccountsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AccountsGetOptions
		wantErr bool
	}{
		{
			name:    "valid account ID",
			opts:    &AccountsGetOptions{AccountID: 1},
			wantErr: false,
		},
		{
			name:    "zero account ID",
			opts:    &AccountsGetOptions{AccountID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAccountsSubAccountsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AccountsSubAccountsOptions
		wantErr bool
	}{
		{
			name:    "valid account ID",
			opts:    &AccountsSubAccountsOptions{AccountID: 1},
			wantErr: false,
		},
		{
			name:    "zero account ID",
			opts:    &AccountsSubAccountsOptions{AccountID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountsSubAccountsOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
