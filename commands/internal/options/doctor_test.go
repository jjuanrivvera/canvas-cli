package options

import "testing"

func TestDoctorOptions_Validate(t *testing.T) {
	opts := &DoctorOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("DoctorOptions.Validate() error = %v, want nil", err)
	}
}
