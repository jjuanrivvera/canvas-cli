package options

import "testing"

func TestAuthLoginOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AuthLoginOptions
		wantErr bool
	}{
		{
			name:    "valid with URL",
			opts:    &AuthLoginOptions{InstanceURL: "https://canvas.example.com"},
			wantErr: false,
		},
		{
			name:    "valid with instance name",
			opts:    &AuthLoginOptions{InstanceName: "myschool"},
			wantErr: false,
		},
		{
			name:    "valid with URL and oauth mode auto",
			opts:    &AuthLoginOptions{InstanceURL: "https://canvas.example.com", OAuthMode: "auto"},
			wantErr: false,
		},
		{
			name:    "valid with URL and oauth mode local",
			opts:    &AuthLoginOptions{InstanceURL: "https://canvas.example.com", OAuthMode: "local"},
			wantErr: false,
		},
		{
			name:    "valid with URL and oauth mode oob",
			opts:    &AuthLoginOptions{InstanceURL: "https://canvas.example.com", OAuthMode: "oob"},
			wantErr: false,
		},
		{
			name:    "missing both URL and instance name",
			opts:    &AuthLoginOptions{},
			wantErr: true,
		},
		{
			name:    "invalid OAuth mode",
			opts:    &AuthLoginOptions{InstanceURL: "https://canvas.example.com", OAuthMode: "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthLoginOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthLogoutOptions_Validate(t *testing.T) {
	opts := &AuthLogoutOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("AuthLogoutOptions.Validate() error = %v, want nil", err)
	}
}

func TestAuthStatusOptions_Validate(t *testing.T) {
	opts := &AuthStatusOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("AuthStatusOptions.Validate() error = %v, want nil", err)
	}
}

func TestAuthTokenSetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AuthTokenSetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AuthTokenSetOptions{InstanceName: "myschool"},
			wantErr: false,
		},
		{
			name:    "missing instance name",
			opts:    &AuthTokenSetOptions{InstanceName: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthTokenSetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthTokenRemoveOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AuthTokenRemoveOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AuthTokenRemoveOptions{InstanceName: "myschool"},
			wantErr: false,
		},
		{
			name:    "missing instance name",
			opts:    &AuthTokenRemoveOptions{InstanceName: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthTokenRemoveOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
