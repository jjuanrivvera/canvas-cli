package api

import "testing"

func TestValidatePositiveID_Valid(t *testing.T) {
	if err := ValidatePositiveID(1, "course_id"); err != nil {
		t.Errorf("expected no error for positive ID, got %v", err)
	}
	if err := ValidatePositiveID(999, "assignment_id"); err != nil {
		t.Errorf("expected no error for large positive ID, got %v", err)
	}
}

func TestValidatePositiveID_Zero(t *testing.T) {
	if err := ValidatePositiveID(0, "course_id"); err == nil {
		t.Error("expected error for zero ID")
	}
}

func TestValidatePositiveID_Negative(t *testing.T) {
	if err := ValidatePositiveID(-1, "course_id"); err == nil {
		t.Error("expected error for negative ID")
	}
}

func TestValidateNonEmpty_Valid(t *testing.T) {
	if err := ValidateNonEmpty("hello", "name"); err != nil {
		t.Errorf("expected no error for non-empty string, got %v", err)
	}
}

func TestValidateNonEmpty_Empty(t *testing.T) {
	if err := ValidateNonEmpty("", "name"); err == nil {
		t.Error("expected error for empty string")
	}
}

func TestValidateNotNil_NonNil(t *testing.T) {
	v := "something"
	if err := ValidateNotNil(&v, "param"); err != nil {
		t.Errorf("expected no error for non-nil pointer, got %v", err)
	}
}

func TestValidateNotNil_Nil(t *testing.T) {
	if err := ValidateNotNil(nil, "param"); err == nil {
		t.Error("expected error for nil pointer")
	}
}
