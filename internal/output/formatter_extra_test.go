package output

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

// ------------------- helper structs -------------------

// CourseStub mimics a Canvas Course struct with a known keyFieldsMap entry
// so filterKeyFields takes the "key fields defined" path.
type CourseStub struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	CourseCode     string `json:"course_code"`
	WorkflowState  string `json:"workflow_state"`
	AccountID      int    `json:"account_id"`
	EnrollmentTerm int64  `json:"enrollment_term_id"`
}

// UserStub mimics a Canvas User struct.
type UserStub struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	SortableName string `json:"sortable_name"`
	Email        string `json:"email"`
	LoginID      string `json:"login_id"`
}

// WideStruct has more than 6 fields to exercise the ">6 fallback" path.
type WideStruct struct {
	F1 string `json:"f1"`
	F2 string `json:"f2"`
	F3 string `json:"f3"`
	F4 string `json:"f4"`
	F5 string `json:"f5"`
	F6 string `json:"f6"`
	F7 string `json:"f7"`
}

// StructWithIdentifiers exercises formatStructCompact paths.
type NamedThing struct {
	Name string
	ID   int
}

type IDOnlyThing struct {
	ID int
}

type TitleThing struct {
	Title string
}

type EmptyThing struct{}

type FirstFieldThing struct {
	Alpha string
}

// PointerStruct has a *time.Time field to exercise nil-pointer and zero-time paths.
type TimedStruct struct {
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// ------------------- JSONFormatter -------------------

func TestJSONFormatter_NilData(t *testing.T) {
	f := &JSONFormatter{}
	out, err := f.Format(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "null" {
		t.Errorf("expected 'null', got %q", out)
	}
}

func TestJSONFormatter_NilSlice(t *testing.T) {
	f := &JSONFormatter{}
	var s []TestStruct
	out, err := f.Format(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "[]" {
		t.Errorf("expected '[]', got %q", out)
	}
}

// ------------------- YAMLFormatter -------------------

func TestYAMLFormatter_NilData(t *testing.T) {
	f := &YAMLFormatter{}
	out, err := f.Format(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "[]\n" {
		t.Errorf("expected '[]\n', got %q", out)
	}
}

func TestYAMLFormatter_NilSlice(t *testing.T) {
	f := &YAMLFormatter{}
	var s []TestStruct
	out, err := f.Format(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "[]\n" {
		t.Errorf("expected '[]\n', got %q", out)
	}
}

func TestYAMLFormatter_Slice(t *testing.T) {
	f := &YAMLFormatter{}
	data := []TestStruct{
		{Name: "Alice", Age: 30, Email: "alice@example.com"},
	}
	out, err := f.Format(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Alice") {
		t.Errorf("expected 'Alice' in YAML output, got %s", out)
	}
}

// ------------------- CSVFormatter -------------------

func TestCSVFormatter_SingleItem(t *testing.T) {
	f := &CSVFormatter{}
	data := TestStruct{Name: "Bob", Age: 25, Email: "bob@example.com"}
	out, err := f.Format(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Bob") {
		t.Errorf("expected 'Bob' in CSV output, got %s", out)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 { // header + 1 data row
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

// ------------------- TableFormatter -------------------

func TestTableFormatter_Verbose_AllHeaders(t *testing.T) {
	// Verbose mode should NOT filter fields.
	f := &TableFormatter{Verbose: true}
	data := []TestStruct{
		{Name: "Carol", Age: 40, Email: "carol@example.com"},
	}
	out, err := f.Format(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// All three headers should appear.
	for _, h := range []string{"name", "age", "email"} {
		if !strings.Contains(out, h) {
			t.Errorf("expected header %q in verbose output, got:\n%s", h, out)
		}
	}
}

func TestTableFormatter_FilterKeyFields_KnownType(t *testing.T) {
	// Register "CourseStub" in keyFieldsMap so we can exercise the "found" path.
	// CourseStub has json tags matching the keys defined in keyFieldsMap for "Course".
	// We test via the Format method using a custom map entry.
	// Instead, directly call filterKeyFields with a synthetic type via a struct
	// that matches a keyFieldsMap entry.  We add a temporary entry to test it.
	oldCourse := keyFieldsMap["CourseStub"]
	keyFieldsMap["CourseStub"] = []string{"id", "name", "course_code"}
	defer func() {
		if oldCourse == nil {
			delete(keyFieldsMap, "CourseStub")
		} else {
			keyFieldsMap["CourseStub"] = oldCourse
		}
	}()

	f := &TableFormatter{Verbose: false}
	item := CourseStub{ID: 1, Name: "Intro", CourseCode: "CS101"}
	allHeaders := getHeaders(item)
	filtered := f.filterKeyFields(item, allHeaders)

	// Only the intersection of known keys and actual headers should appear.
	for _, h := range filtered {
		if h != "id" && h != "name" && h != "course_code" {
			t.Errorf("unexpected header %q after filtering", h)
		}
	}
}

func TestTableFormatter_FilterKeyFields_UnknownTypeMoreThan6(t *testing.T) {
	f := &TableFormatter{Verbose: false}
	item := WideStruct{F1: "a", F2: "b", F3: "c", F4: "d", F5: "e", F6: "f", F7: "g"}
	allHeaders := getHeaders(item)
	filtered := f.filterKeyFields(item, allHeaders)
	if len(filtered) != 6 {
		t.Errorf("expected 6 headers for unknown type with 7 fields, got %d", len(filtered))
	}
}

func TestTableFormatter_FilterKeyFields_KnownTypeNoMatch(t *testing.T) {
	// A type that IS in keyFieldsMap but whose actual json tags don't match the
	// registered keys → fallback to first 6 headers.
	oldEntry := keyFieldsMap["WideStruct"]
	keyFieldsMap["WideStruct"] = []string{"nonexistent_field_1", "nonexistent_field_2"}
	defer func() {
		if oldEntry == nil {
			delete(keyFieldsMap, "WideStruct")
		} else {
			keyFieldsMap["WideStruct"] = oldEntry
		}
	}()

	f := &TableFormatter{Verbose: false}
	item := WideStruct{}
	allHeaders := getHeaders(item)
	filtered := f.filterKeyFields(item, allHeaders)
	// No match → fallback → 6 (since WideStruct has 7 fields)
	if len(filtered) != 6 {
		t.Errorf("expected 6 headers (fallback), got %d", len(filtered))
	}
}

func TestTableFormatter_FilterKeyFields_PointerItem(t *testing.T) {
	f := &TableFormatter{Verbose: false}
	item := &TestStruct{Name: "X", Age: 1, Email: "x@x.com"}
	allHeaders := getHeaders(item)
	// TestStruct has 3 fields (<= 6), no entry in keyFieldsMap → returns all
	filtered := f.filterKeyFields(item, allHeaders)
	if len(filtered) != 3 {
		t.Errorf("expected 3 headers for pointer item, got %d", len(filtered))
	}
}

// ------------------- formatValue edge cases -------------------

func TestFormatValue_ZeroTime(t *testing.T) {
	var zero time.Time
	result := formatValue(zero)
	if result != "Not set" {
		t.Errorf("expected 'Not set' for zero time, got %q", result)
	}
}

func TestFormatValue_NonZeroTime(t *testing.T) {
	ts := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	result := formatValue(ts)
	if result != "2024-01-15 10:30" {
		t.Errorf("expected '2024-01-15 10:30', got %q", result)
	}
}

func TestFormatValue_PtrTime_Nil(t *testing.T) {
	var ts *time.Time
	result := formatValue(ts)
	if result != "Not set" {
		t.Errorf("expected 'Not set' for nil *time.Time, got %q", result)
	}
}

func TestFormatValue_PtrTime_Zero(t *testing.T) {
	zero := time.Time{}
	result := formatValue(&zero)
	if result != "Not set" {
		t.Errorf("expected 'Not set' for zero *time.Time, got %q", result)
	}
}

func TestFormatValue_PtrTime_NonZero(t *testing.T) {
	ts := time.Date(2025, 6, 8, 0, 0, 0, 0, time.UTC)
	result := formatValue(&ts)
	if result != "2025-06-08 00:00" {
		t.Errorf("expected '2025-06-08 00:00', got %q", result)
	}
}

func TestFormatValue_NilPointer(t *testing.T) {
	var p *TestStruct
	result := formatValue(p)
	// nil pointer → empty string
	if result != "" {
		t.Errorf("expected '' for nil pointer, got %q", result)
	}
}

func TestFormatValue_NonEmptyMap(t *testing.T) {
	m := map[string]int{"a": 1}
	result := formatValue(m)
	if !strings.Contains(result, "a") || !strings.Contains(result, "1") {
		t.Errorf("expected map content in output, got %q", result)
	}
	if !strings.HasPrefix(result, "{") || !strings.HasSuffix(result, "}") {
		t.Errorf("expected braces around map, got %q", result)
	}
}

func TestFormatValue_EmptyMap(t *testing.T) {
	m := map[string]int{}
	result := formatValue(m)
	if result != "" {
		t.Errorf("expected '' for empty map, got %q", result)
	}
}

func TestFormatValue_Int8_Int16_Int32(t *testing.T) {
	tests := []struct {
		val      interface{}
		expected string
	}{
		{int8(10), "10"},
		{int16(1000), "1000"},
		{int32(100000), "100000"},
		{uint8(255), "255"},
		{uint16(65535), "65535"},
		{uint32(4294967295), "4294967295"},
		{float32(1.5), "1.50"},
	}
	for _, tt := range tests {
		result := formatValue(tt.val)
		if result != tt.expected {
			t.Errorf("formatValue(%v) = %q, want %q", tt.val, result, tt.expected)
		}
	}
}

func TestFormatValue_NestedSlice(t *testing.T) {
	// Slice of structs — exercises the struct branch inside formatValue
	data := []NamedThing{{Name: "Alice", ID: 1}, {Name: "Bob", ID: 2}}
	result := formatValue(data)
	if !strings.Contains(result, "Alice") {
		t.Errorf("expected 'Alice' in nested slice output, got %q", result)
	}
}

// ------------------- formatStructCompact -------------------

func TestFormatStructCompact_NameAndID(t *testing.T) {
	item := NamedThing{Name: "Widget", ID: 42}
	out := formatValue(item)
	// "name" field found first; just shows the value
	if !strings.Contains(out, "Widget") {
		t.Errorf("expected 'Widget' in output, got %q", out)
	}
}

func TestFormatStructCompact_IDOnly(t *testing.T) {
	item := IDOnlyThing{ID: 7}
	out := formatValue(item)
	// "ID" found, formatted as "Type(ID)"
	if !strings.Contains(out, "IDOnlyThing") || !strings.Contains(out, "7") {
		t.Errorf("expected 'IDOnlyThing(7)' in output, got %q", out)
	}
}

func TestFormatStructCompact_TitleOnly(t *testing.T) {
	item := TitleThing{Title: "Hello World"}
	out := formatValue(item)
	if !strings.Contains(out, "Hello World") {
		t.Errorf("expected 'Hello World' in output, got %q", out)
	}
}

func TestFormatStructCompact_Empty(t *testing.T) {
	item := EmptyThing{}
	out := formatValue(item)
	if !strings.Contains(out, "EmptyThing") {
		t.Errorf("expected 'EmptyThing{}' in output, got %q", out)
	}
}

func TestFormatStructCompact_FirstField(t *testing.T) {
	item := FirstFieldThing{Alpha: "first"}
	out := formatValue(item)
	// No name/title/id → falls through to first-non-zero-field fallback
	if !strings.Contains(out, "first") {
		t.Errorf("expected 'first' in output, got %q", out)
	}
}

// ------------------- Ptr-to-struct branch in formatValue -------------------

func TestFormatValue_PtrToStruct_NonNil(t *testing.T) {
	item := &NamedThing{Name: "PtrThing", ID: 99}
	out := formatValue(item)
	if !strings.Contains(out, "PtrThing") {
		t.Errorf("expected 'PtrThing' in output, got %q", out)
	}
}

// ------------------- WriteWithOptions -------------------

func TestWriteWithOptions_Verbose(t *testing.T) {
	data := []TestStruct{
		{Name: "Dave", Age: 22, Email: "dave@example.com"},
	}
	var buf bytes.Buffer
	err := WriteWithOptions(&buf, data, FormatTable, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Dave") {
		t.Errorf("expected 'Dave' in verbose table output, got %s", out)
	}
}

func TestWriteWithOptions_InvalidFormat(t *testing.T) {
	var buf bytes.Buffer
	err := WriteWithOptions(&buf, nil, "xml", false)
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

// ------------------- time fields in table / CSV -------------------

func TestTableFormatter_TimeFields(t *testing.T) {
	now := time.Date(2024, 3, 10, 9, 0, 0, 0, time.UTC)
	data := []TimedStruct{
		{Name: "event1", CreatedAt: now, UpdatedAt: nil},
		{Name: "event2", CreatedAt: time.Time{}, UpdatedAt: &now},
	}

	f := &TableFormatter{Verbose: true}
	out, err := f.Format(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "2024-03-10 09:00") {
		t.Errorf("expected formatted time in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Not set") {
		t.Errorf("expected 'Not set' for zero time in output, got:\n%s", out)
	}
}

// ------------------- getHeaders and getRow with pointer struct -------------------

func TestGetRow_MapData(t *testing.T) {
	data := map[string]interface{}{"key1": "val1", "key2": 99}
	headers := getHeaders(data)

	row := getRow(data, headers)
	if len(row) != len(headers) {
		t.Errorf("row length %d != headers length %d", len(row), len(headers))
	}
}

func TestGetHeaders_Unsupported(t *testing.T) {
	headers := getHeaders(42) // int is not struct or map
	if len(headers) != 0 {
		t.Errorf("expected empty headers for int, got %v", headers)
	}
}

// ------------------- NewFormatterWithOptions -------------------

func TestNewFormatterWithOptions_Verbose(t *testing.T) {
	f, err := NewFormatterWithOptions(FormatTable, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tf, ok := f.(*TableFormatter)
	if !ok {
		t.Fatal("expected *TableFormatter")
	}
	if !tf.Verbose {
		t.Error("expected Verbose to be true")
	}
}
