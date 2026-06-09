package options

import (
	"strings"
	"testing"
)

func TestValidateRequired_Int(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		value     int
		wantErr   bool
	}{
		{
			name:      "non-zero int",
			fieldName: "position",
			value:     1,
			wantErr:   false,
		},
		{
			name:      "zero int",
			fieldName: "position",
			value:     0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRequired(tt.fieldName, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRequired() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRequired_UnknownType(t *testing.T) {
	// For any type not handled by the switch (e.g., float64), always returns nil.
	err := ValidateRequired("amount", 3.14)
	if err != nil {
		t.Errorf("ValidateRequired() with unhandled type error = %v, want nil", err)
	}
}

func TestErrInvalidValue(t *testing.T) {
	tests := []struct {
		name         string
		fieldName    string
		value        string
		validOptions []string
		wantContain  string
	}{
		{
			name:         "single valid option",
			fieldName:    "event",
			value:        "bad",
			validOptions: []string{"good"},
			wantContain:  "invalid event: bad (valid options: good)",
		},
		{
			name:         "two valid options",
			fieldName:    "mode",
			value:        "x",
			validOptions: []string{"a", "b"},
			wantContain:  "invalid mode: x (valid options: a or b)",
		},
		{
			name:         "three valid options",
			fieldName:    "color",
			value:        "purple",
			validOptions: []string{"red", "green", "blue"},
			wantContain:  "invalid color: purple (valid options: red, green or blue)",
		},
		{
			name:         "no valid options",
			fieldName:    "field",
			value:        "val",
			validOptions: []string{},
			wantContain:  "invalid field: val",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrInvalidValue(tt.fieldName, tt.value, tt.validOptions...)
			if err == nil {
				t.Fatal("ErrInvalidValue() returned nil, want error")
			}
			if !strings.Contains(err.Error(), tt.wantContain) {
				t.Errorf("ErrInvalidValue() = %q, want to contain %q", err.Error(), tt.wantContain)
			}
		})
	}
}
