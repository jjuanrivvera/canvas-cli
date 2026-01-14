package output

import (
	"bytes"
	"strings"
	"testing"
)

type TestStruct struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func TestNewFormatter(t *testing.T) {
	tests := []struct {
		format      FormatType
		shouldError bool
	}{
		{FormatJSON, false},
		{FormatYAML, false},
		{FormatCSV, false},
		{FormatTable, false},
		{"invalid", true},
	}

	for _, tt := range tests {
		formatter, err := NewFormatter(tt.format)

		if tt.shouldError {
			if err == nil {
				t.Errorf("expected error for format %s", tt.format)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for format %s: %v", tt.format, err)
			}
			if formatter == nil {
				t.Errorf("expected non-nil formatter for format %s", tt.format)
			}
		}
	}
}

func TestJSONFormatter_Format(t *testing.T) {
	formatter := &JSONFormatter{}

	data := TestStruct{
		Name:  "John Doe",
		Age:   30,
		Email: "john@example.com",
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	if !strings.Contains(output, "John Doe") {
		t.Error("expected output to contain 'John Doe'")
	}

	if !strings.Contains(output, "john@example.com") {
		t.Error("expected output to contain 'john@example.com'")
	}

	if !strings.Contains(output, "\"age\": 30") {
		t.Error("expected output to contain age")
	}
}

func TestJSONFormatter_FormatSlice(t *testing.T) {
	formatter := &JSONFormatter{}

	data := []TestStruct{
		{Name: "John", Age: 30, Email: "john@example.com"},
		{Name: "Jane", Age: 25, Email: "jane@example.com"},
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	if !strings.Contains(output, "John") {
		t.Error("expected output to contain 'John'")
	}

	if !strings.Contains(output, "Jane") {
		t.Error("expected output to contain 'Jane'")
	}
}

func TestYAMLFormatter_Format(t *testing.T) {
	formatter := &YAMLFormatter{}

	data := TestStruct{
		Name:  "John Doe",
		Age:   30,
		Email: "john@example.com",
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	if !strings.Contains(output, "name: John Doe") {
		t.Error("expected output to contain 'name: John Doe'")
	}

	if !strings.Contains(output, "age: 30") {
		t.Error("expected output to contain 'age: 30'")
	}
}

func TestCSVFormatter_Format(t *testing.T) {
	formatter := &CSVFormatter{}

	data := []TestStruct{
		{Name: "John", Age: 30, Email: "john@example.com"},
		{Name: "Jane", Age: 25, Email: "jane@example.com"},
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 3 { // Header + 2 rows
		t.Errorf("expected 3 lines, got %d", len(lines))
	}

	// Check header
	if !strings.Contains(lines[0], "name") {
		t.Error("expected header to contain 'name'")
	}

	// Check data
	if !strings.Contains(output, "John") {
		t.Error("expected output to contain 'John'")
	}

	if !strings.Contains(output, "Jane") {
		t.Error("expected output to contain 'Jane'")
	}
}

func TestCSVFormatter_FormatEmpty(t *testing.T) {
	formatter := &CSVFormatter{}

	data := []TestStruct{}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	if output != "" {
		t.Errorf("expected empty output, got '%s'", output)
	}
}

func TestTableFormatter_Format(t *testing.T) {
	formatter := &TableFormatter{}

	data := []TestStruct{
		{Name: "John", Age: 30, Email: "john@example.com"},
		{Name: "Jane", Age: 25, Email: "jane@example.com"},
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Check for table borders
	if !strings.Contains(output, "┌") || !strings.Contains(output, "┐") {
		t.Error("expected table to have top border")
	}

	if !strings.Contains(output, "└") || !strings.Contains(output, "┘") {
		t.Error("expected table to have bottom border")
	}

	// Check for data
	if !strings.Contains(output, "John") {
		t.Error("expected output to contain 'John'")
	}

	if !strings.Contains(output, "Jane") {
		t.Error("expected output to contain 'Jane'")
	}
}

func TestTableFormatter_FormatEmpty(t *testing.T) {
	formatter := &TableFormatter{}

	data := []TestStruct{}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	if output != "No data" {
		t.Errorf("expected 'No data', got '%s'", output)
	}
}

func TestTableFormatter_FormatSingleItem(t *testing.T) {
	formatter := &TableFormatter{}

	data := TestStruct{Name: "John", Age: 30, Email: "john@example.com"}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	if !strings.Contains(output, "John") {
		t.Error("expected output to contain 'John'")
	}
}

func TestToSlice(t *testing.T) {
	// Test with slice
	slice := []string{"a", "b", "c"}
	result := toSlice(slice)

	if len(result) != 3 {
		t.Errorf("expected length 3, got %d", len(result))
	}

	// Test with single item
	single := "test"
	result = toSlice(single)

	if len(result) != 1 {
		t.Errorf("expected length 1, got %d", len(result))
	}

	if result[0] != "test" {
		t.Errorf("expected 'test', got '%v'", result[0])
	}
}

func TestGetHeaders(t *testing.T) {
	// Test with struct
	data := TestStruct{Name: "John", Age: 30, Email: "john@example.com"}
	headers := getHeaders(data)

	if len(headers) != 3 {
		t.Errorf("expected 3 headers, got %d", len(headers))
	}

	// Headers should use json tags
	expectedHeaders := []string{"name", "age", "email"}
	for i, expected := range expectedHeaders {
		if headers[i] != expected {
			t.Errorf("expected header '%s', got '%s'", expected, headers[i])
		}
	}

	// Test with map
	mapData := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	headers = getHeaders(mapData)

	if len(headers) != 2 {
		t.Errorf("expected 2 headers, got %d", len(headers))
	}
}

func TestGetRow(t *testing.T) {
	data := TestStruct{Name: "John", Age: 30, Email: "john@example.com"}
	headers := []string{"name", "age", "email"}

	row := getRow(data, headers)

	if len(row) != 3 {
		t.Errorf("expected 3 values, got %d", len(row))
	}

	if row[0] != "John" {
		t.Errorf("expected 'John', got '%s'", row[0])
	}

	if row[1] != "30" {
		t.Errorf("expected '30', got '%s'", row[1])
	}

	if row[2] != "john@example.com" {
		t.Errorf("expected 'john@example.com', got '%s'", row[2])
	}
}

func TestFormatValue(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{nil, ""},
		{"test", "test"},
		{123, "123"},
		{45.67, "45.67"},
		{true, "true"},
		{false, "false"},
		{[]string{"a", "b", "c"}, "[a, b, c]"},
		{[]int{}, ""},
	}

	for _, tt := range tests {
		result := formatValue(tt.input)
		if result != tt.expected {
			t.Errorf("formatValue(%v) = '%s', expected '%s'", tt.input, result, tt.expected)
		}
	}
}

func TestWrite(t *testing.T) {
	data := []TestStruct{
		{Name: "John", Age: 30, Email: "john@example.com"},
	}

	// Test JSON
	var buf bytes.Buffer
	err := Write(&buf, data, FormatJSON)
	if err != nil {
		t.Fatalf("Write JSON failed: %v", err)
	}

	if !strings.Contains(buf.String(), "John") {
		t.Error("expected output to contain 'John'")
	}

	// Test YAML
	buf.Reset()
	err = Write(&buf, data, FormatYAML)
	if err != nil {
		t.Fatalf("Write YAML failed: %v", err)
	}

	// Test CSV
	buf.Reset()
	err = Write(&buf, data, FormatCSV)
	if err != nil {
		t.Fatalf("Write CSV failed: %v", err)
	}

	// Test Table
	buf.Reset()
	err = Write(&buf, data, FormatTable)
	if err != nil {
		t.Fatalf("Write Table failed: %v", err)
	}

	// Test invalid format
	buf.Reset()
	err = Write(&buf, data, "invalid")
	if err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestFindField(t *testing.T) {
	// This is an internal function, but we can test it indirectly through getRow
	data := TestStruct{Name: "John", Age: 30, Email: "john@example.com"}

	// Test finding by json tag
	headers := []string{"name"}
	row := getRow(data, headers)

	if len(row) != 1 || row[0] != "John" {
		t.Error("failed to find field by json tag")
	}

	// Test with non-existent field
	headers = []string{"nonexistent"}
	row = getRow(data, headers)

	if len(row) != 1 || row[0] != "" {
		t.Error("expected empty string for non-existent field")
	}
}

func TestCSVFormatter_FormatMap(t *testing.T) {
	formatter := &CSVFormatter{}

	data := []map[string]interface{}{
		{"name": "John", "age": 30},
		{"name": "Jane", "age": 25},
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	if !strings.Contains(output, "John") {
		t.Error("expected output to contain 'John'")
	}
}

func TestTableFormatter_FormatMap(t *testing.T) {
	formatter := &TableFormatter{}

	data := []map[string]interface{}{
		{"name": "John", "age": 30},
		{"name": "Jane", "age": 25},
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	if !strings.Contains(output, "John") {
		t.Error("expected output to contain 'John'")
	}
}

func TestFormatValue_ComplexTypes(t *testing.T) {
	// Test with nested slice
	nested := [][]int{{1, 2}, {3, 4}}
	result := formatValue(nested)

	if !strings.Contains(result, "1") {
		t.Error("expected formatted value to contain nested data")
	}

	// Test with uint types
	var u uint = 42
	result = formatValue(u)
	if result != "42" {
		t.Errorf("expected '42', got '%s'", result)
	}

	// Test with float
	f := 3.14159
	result = formatValue(f)
	if result != "3.14" {
		t.Errorf("expected '3.14', got '%s'", result)
	}
}

func TestGetHeaders_Pointer(t *testing.T) {
	data := &TestStruct{Name: "John", Age: 30, Email: "john@example.com"}
	headers := getHeaders(data)

	if len(headers) != 3 {
		t.Errorf("expected 3 headers from pointer, got %d", len(headers))
	}
}

func TestGetRow_Pointer(t *testing.T) {
	data := &TestStruct{Name: "John", Age: 30, Email: "john@example.com"}
	headers := []string{"name", "age", "email"}

	row := getRow(data, headers)

	if len(row) != 3 {
		t.Errorf("expected 3 values from pointer, got %d", len(row))
	}

	if row[0] != "John" {
		t.Error("failed to get row from pointer")
	}
}

func TestTableFormatter_ColumnWidthCalculation(t *testing.T) {
	formatter := &TableFormatter{}

	data := []TestStruct{
		{Name: "A", Age: 1, Email: "a@example.com"},
		{Name: "VeryLongName", Age: 99999, Email: "verylongemail@example.com"},
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Verify the long values are present
	if !strings.Contains(output, "VeryLongName") {
		t.Error("expected output to contain 'VeryLongName'")
	}

	if !strings.Contains(output, "verylongemail@example.com") {
		t.Error("expected output to contain long email")
	}
}
