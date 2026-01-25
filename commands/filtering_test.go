package commands

import (
	"reflect"
	"testing"
)

func TestToString(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "string value",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "integer as float64",
			input:    float64(123),
			expected: "123",
		},
		{
			name:     "decimal float64",
			input:    float64(123.45),
			expected: "123.45",
		},
		{
			name:     "boolean true",
			input:    true,
			expected: "true",
		},
		{
			name:     "boolean false",
			input:    false,
			expected: "false",
		},
		{
			name:     "nil value",
			input:    nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toString(tt.input)
			if result != tt.expected {
				t.Errorf("toString(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStructToMap(t *testing.T) {
	type testStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	input := testStruct{Name: "test", Value: 42}
	result := structToMap(input)

	if result == nil {
		t.Fatal("structToMap() returned nil")
	}

	if result["name"] != "test" {
		t.Errorf("structToMap().name = %v, want test", result["name"])
	}

	// JSON numbers are float64
	if result["value"] != float64(42) {
		t.Errorf("structToMap().value = %v, want 42", result["value"])
	}
}

func TestStructToMap_AlreadyMap(t *testing.T) {
	input := map[string]interface{}{"name": "test", "value": 42}
	result := structToMap(input)

	if !reflect.DeepEqual(result, input) {
		t.Errorf("structToMap() = %v, want %v", result, input)
	}
}

func TestToMapSlice(t *testing.T) {
	type testStruct struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	input := []testStruct{
		{ID: 1, Name: "first"},
		{ID: 2, Name: "second"},
	}

	result := toMapSlice(input)

	if len(result) != 2 {
		t.Fatalf("toMapSlice() len = %d, want 2", len(result))
	}

	if result[0]["name"] != "first" {
		t.Errorf("toMapSlice()[0].name = %v, want first", result[0]["name"])
	}

	if result[1]["name"] != "second" {
		t.Errorf("toMapSlice()[1].name = %v, want second", result[1]["name"])
	}
}

func TestToMapSlice_NonSlice(t *testing.T) {
	result := toMapSlice("not a slice")

	if result != nil {
		t.Errorf("toMapSlice(non-slice) = %v, want nil", result)
	}
}

func TestFilterByText(t *testing.T) {
	items := []map[string]interface{}{
		{"name": "Alice", "email": "alice@example.com"},
		{"name": "Bob", "email": "bob@example.com"},
		{"name": "Charlie", "email": "charlie@test.com"},
	}

	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{
			name:     "filter by name",
			text:     "alice",
			expected: 1,
		},
		{
			name:     "filter by email domain",
			text:     "example.com",
			expected: 2,
		},
		{
			name:     "no match",
			text:     "xyz",
			expected: 0,
		},
		{
			name:     "case insensitive",
			text:     "ALICE",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterByText(items, tt.text)
			if len(result) != tt.expected {
				t.Errorf("filterByText() len = %d, want %d", len(result), tt.expected)
			}
		})
	}
}

func TestItemContainsText(t *testing.T) {
	item := map[string]interface{}{
		"name":  "Alice",
		"id":    float64(123),
		"email": nil,
	}

	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{
			name:     "contains in string field",
			text:     "alice",
			expected: true,
		},
		{
			name:     "contains in numeric field",
			text:     "123",
			expected: true,
		},
		{
			name:     "does not contain",
			text:     "xyz",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := itemContainsText(item, tt.text)
			if result != tt.expected {
				t.Errorf("itemContainsText() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSortByField(t *testing.T) {
	items := []map[string]interface{}{
		{"name": "Charlie", "id": float64(3)},
		{"name": "Alice", "id": float64(1)},
		{"name": "Bob", "id": float64(2)},
	}

	t.Run("sort by name ascending", func(t *testing.T) {
		input := make([]map[string]interface{}, len(items))
		copy(input, items)

		result := sortByField(input, "name")

		if result[0]["name"] != "Alice" {
			t.Errorf("sortByField()[0].name = %v, want Alice", result[0]["name"])
		}
		if result[1]["name"] != "Bob" {
			t.Errorf("sortByField()[1].name = %v, want Bob", result[1]["name"])
		}
		if result[2]["name"] != "Charlie" {
			t.Errorf("sortByField()[2].name = %v, want Charlie", result[2]["name"])
		}
	})

	t.Run("sort by name descending", func(t *testing.T) {
		input := make([]map[string]interface{}, len(items))
		copy(input, items)

		result := sortByField(input, "-name")

		if result[0]["name"] != "Charlie" {
			t.Errorf("sortByField()[0].name = %v, want Charlie", result[0]["name"])
		}
		if result[2]["name"] != "Alice" {
			t.Errorf("sortByField()[2].name = %v, want Alice", result[2]["name"])
		}
	})

	t.Run("sort by numeric field", func(t *testing.T) {
		input := make([]map[string]interface{}, len(items))
		copy(input, items)

		result := sortByField(input, "id")

		if result[0]["id"] != float64(1) {
			t.Errorf("sortByField()[0].id = %v, want 1", result[0]["id"])
		}
		if result[2]["id"] != float64(3) {
			t.Errorf("sortByField()[2].id = %v, want 3", result[2]["id"])
		}
	})
}

func TestSelectColumns(t *testing.T) {
	items := []map[string]interface{}{
		{"name": "Alice", "email": "alice@example.com", "id": float64(1)},
		{"name": "Bob", "email": "bob@example.com", "id": float64(2)},
	}

	result := selectColumns(items, []string{"name", "id"})

	if len(result) != 2 {
		t.Fatalf("selectColumns() len = %d, want 2", len(result))
	}

	// Check first item has only selected columns
	if _, exists := result[0]["email"]; exists {
		t.Error("selectColumns() should not include email column")
	}

	if result[0]["name"] != "Alice" {
		t.Errorf("selectColumns()[0].name = %v, want Alice", result[0]["name"])
	}

	if result[0]["id"] != float64(1) {
		t.Errorf("selectColumns()[0].id = %v, want 1", result[0]["id"])
	}
}

func TestCompareValues(t *testing.T) {
	tests := []struct {
		name     string
		a        interface{}
		b        interface{}
		expected int
	}{
		{
			name:     "both nil",
			a:        nil,
			b:        nil,
			expected: 0,
		},
		{
			name:     "a nil",
			a:        nil,
			b:        "test",
			expected: -1,
		},
		{
			name:     "b nil",
			a:        "test",
			b:        nil,
			expected: 1,
		},
		{
			name:     "numeric less",
			a:        float64(1),
			b:        float64(2),
			expected: -1,
		},
		{
			name:     "numeric equal",
			a:        float64(1),
			b:        float64(1),
			expected: 0,
		},
		{
			name:     "numeric greater",
			a:        float64(2),
			b:        float64(1),
			expected: 1,
		},
		{
			name:     "string less",
			a:        "alice",
			b:        "bob",
			expected: -1,
		},
		{
			name:     "string equal",
			a:        "alice",
			b:        "alice",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareValues(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("compareValues(%v, %v) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected float64
		ok       bool
	}{
		{
			name:     "float64",
			input:    float64(1.5),
			expected: 1.5,
			ok:       true,
		},
		{
			name:     "float32",
			input:    float32(1.5),
			expected: 1.5,
			ok:       true,
		},
		{
			name:     "int",
			input:    1,
			expected: 1,
			ok:       true,
		},
		{
			name:     "int64",
			input:    int64(1),
			expected: 1,
			ok:       true,
		},
		{
			name:     "string",
			input:    "not a number",
			expected: 0,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := toFloat64(tt.input)
			if ok != tt.ok {
				t.Errorf("toFloat64(%v) ok = %v, want %v", tt.input, ok, tt.ok)
			}
			if ok && result != tt.expected {
				t.Errorf("toFloat64(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestHasFilteringOptions(t *testing.T) {
	// Save original values
	origFilterText := filterText
	origFilterColumns := filterColumns
	origSortField := sortField

	defer func() {
		// Restore original values
		filterText = origFilterText
		filterColumns = origFilterColumns
		sortField = origSortField
	}()

	t.Run("no options set", func(t *testing.T) {
		filterText = ""
		filterColumns = nil
		sortField = ""

		if hasFilteringOptions() {
			t.Error("hasFilteringOptions() = true, want false")
		}
	})

	t.Run("filter text set", func(t *testing.T) {
		filterText = "test"
		filterColumns = nil
		sortField = ""

		if !hasFilteringOptions() {
			t.Error("hasFilteringOptions() = false, want true")
		}
	})

	t.Run("filter columns set", func(t *testing.T) {
		filterText = ""
		filterColumns = []string{"name"}
		sortField = ""

		if !hasFilteringOptions() {
			t.Error("hasFilteringOptions() = false, want true")
		}
	})

	t.Run("sort field set", func(t *testing.T) {
		filterText = ""
		filterColumns = nil
		sortField = "name"

		if !hasFilteringOptions() {
			t.Error("hasFilteringOptions() = false, want true")
		}
	})
}
