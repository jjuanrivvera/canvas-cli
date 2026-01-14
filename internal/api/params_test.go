package api

import (
	"testing"
)

func TestParamsBuilder_Set(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    interface{}
		expected bool
	}{
		{"non-empty string", "name", "test", true},
		{"empty string", "name", "", false},
		{"non-zero int", "count", 5, true},
		{"zero int", "count", 0, false},
		{"non-zero int64", "id", int64(123), true},
		{"zero int64", "id", int64(0), false},
		{"non-zero float", "points", 10.5, true},
		{"zero float", "points", 0.0, false},
		{"non-nil pointer", "value", ptrInt(5), true},
		{"nil pointer", "value", (*int)(nil), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewParamsBuilder()
			b.Set(tt.key, tt.value)
			params := b.Build()

			_, exists := params[tt.key]
			if exists != tt.expected {
				t.Errorf("expected key %q exists=%v, got %v", tt.key, tt.expected, exists)
			}
		})
	}
}

func TestParamsBuilder_SetBool(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"true bool", true, true},
		{"false bool", false, false},
		{"non-nil true ptr", ptrBool(true), true},
		{"non-nil false ptr", ptrBool(false), true}, // Should set to false
		{"nil bool ptr", (*bool)(nil), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewParamsBuilder()
			b.SetBool("flag", tt.value)
			params := b.Build()

			_, exists := params["flag"]
			if exists != tt.expected {
				t.Errorf("expected key exists=%v, got %v", tt.expected, exists)
			}
		})
	}
}

func TestParamsBuilder_SetSlice(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"non-empty slice", []string{"a", "b"}, true},
		{"empty slice", []string{}, false},
		{"nil slice", []string(nil), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewParamsBuilder()
			b.SetSlice("items", tt.value)
			params := b.Build()

			_, exists := params["items"]
			if exists != tt.expected {
				t.Errorf("expected key exists=%v, got %v", tt.expected, exists)
			}
		})
	}
}

func TestParamsBuilder_SetMap(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"non-empty map", map[string]string{"key": "value"}, true},
		{"empty map", map[string]string{}, false},
		{"nil map", map[string]string(nil), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewParamsBuilder()
			b.SetMap("data", tt.value)
			params := b.Build()

			_, exists := params["data"]
			if exists != tt.expected {
				t.Errorf("expected key exists=%v, got %v", tt.expected, exists)
			}
		})
	}
}

func TestParamsBuilder_WithPrefix(t *testing.T) {
	b := NewParamsBuilder().WithPrefix("assignment")
	b.Set("name", "Test Assignment")
	b.Set("points_possible", 100)

	params := b.Build()

	if _, ok := params["assignment[name]"]; !ok {
		t.Error("expected prefixed key 'assignment[name]' to exist")
	}
	if _, ok := params["assignment[points_possible]"]; !ok {
		t.Error("expected prefixed key 'assignment[points_possible]' to exist")
	}
}

func TestURLParamsBuilder(t *testing.T) {
	b := NewURLParamsBuilder()
	b.Set("search_term", "test")
	b.Set("empty", "")
	b.SetInt("per_page", 50)
	b.SetInt("zero", 0)
	b.SetInt64("course_id", int64(12345))
	b.SetBool("include_deleted", true)
	b.SetBoolPtr("published", ptrBool(true))
	b.SetBoolPtr("nil_bool", nil)
	b.AddSlice("include[]", []string{"items", "content"})

	values := b.Build()

	// Check expected values
	if v := values.Get("search_term"); v != "test" {
		t.Errorf("expected search_term=test, got %s", v)
	}
	if v := values.Get("empty"); v != "" {
		t.Errorf("expected empty to not be set, got %s", v)
	}
	if v := values.Get("per_page"); v != "50" {
		t.Errorf("expected per_page=50, got %s", v)
	}
	if v := values.Get("zero"); v != "" {
		t.Errorf("expected zero to not be set, got %s", v)
	}
	if v := values.Get("course_id"); v != "12345" {
		t.Errorf("expected course_id=12345, got %s", v)
	}
	if v := values.Get("include_deleted"); v != "true" {
		t.Errorf("expected include_deleted=true, got %s", v)
	}
	if v := values.Get("published"); v != "true" {
		t.Errorf("expected published=true, got %s", v)
	}
	if v := values.Get("nil_bool"); v != "" {
		t.Errorf("expected nil_bool to not be set, got %s", v)
	}
	if values["include[]"] == nil || len(values["include[]"]) != 2 {
		t.Error("expected include[] to have 2 values")
	}
}

// Helper functions
func ptrInt(i int) *int {
	return &i
}

func ptrBool(b bool) *bool {
	return &b
}
