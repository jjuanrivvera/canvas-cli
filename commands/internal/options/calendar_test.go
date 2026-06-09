package options

import "testing"

func TestCalendarListOptions_Validate(t *testing.T) {
	opts := &CalendarListOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("CalendarListOptions.Validate() error = %v, want nil", err)
	}
}

func TestCalendarGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *CalendarGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &CalendarGetOptions{EventID: 1},
			wantErr: false,
		},
		{
			name:    "missing event ID",
			opts:    &CalendarGetOptions{EventID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CalendarGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCalendarCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *CalendarCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &CalendarCreateOptions{ContextCode: "course_1", Title: "Test Event"},
			wantErr: false,
		},
		{
			name:    "missing context code",
			opts:    &CalendarCreateOptions{ContextCode: "", Title: "Test Event"},
			wantErr: true,
		},
		{
			name:    "missing title",
			opts:    &CalendarCreateOptions{ContextCode: "course_1", Title: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CalendarCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCalendarUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *CalendarUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &CalendarUpdateOptions{EventID: 1},
			wantErr: false,
		},
		{
			name:    "missing event ID",
			opts:    &CalendarUpdateOptions{EventID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CalendarUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCalendarDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *CalendarDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &CalendarDeleteOptions{EventID: 1},
			wantErr: false,
		},
		{
			name:    "missing event ID",
			opts:    &CalendarDeleteOptions{EventID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CalendarDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCalendarReserveOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *CalendarReserveOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &CalendarReserveOptions{EventID: 1},
			wantErr: false,
		},
		{
			name:    "missing event ID",
			opts:    &CalendarReserveOptions{EventID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CalendarReserveOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
