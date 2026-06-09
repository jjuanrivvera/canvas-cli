package batch

import (
	"os"
	"path/filepath"
	"testing"
)

// TestReadGradesCSV_EmptyFile verifies that ReadGradesCSV returns an error
// when the file is completely empty (no header row).
func TestReadGradesCSV_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "empty.csv")
	if err := os.WriteFile(p, []byte(""), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := ReadGradesCSV(p)
	if err == nil {
		t.Error("expected error for empty grades CSV file")
	}
}

// TestReadGradesCSV_MalformedRow verifies that ReadGradesCSV returns an error
// when a data row contains invalid CSV (unclosed quoted field).
func TestReadGradesCSV_MalformedRow(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "malformed.csv")

	// Valid header + a row with an unclosed quoted field to force csv.Reader error.
	content := "user_id,assignment_id,grade\n\"unclosed,456,95\n"
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := ReadGradesCSV(p)
	if err == nil {
		t.Error("expected error for malformed CSV row in grades file")
	}
}

// TestReadCSV_MalformedRow verifies that ReadCSV returns an error
// when a data row contains invalid CSV (unclosed quoted field).
func TestReadCSV_MalformedRow(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "malformed.csv")

	// Valid header + a row with an unclosed quoted field.
	content := "name,age\n\"unclosed,30\n"
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := ReadCSV(p)
	if err == nil {
		t.Error("expected error for malformed CSV row")
	}
}

// TestReadCSV_RowWithAllEmptyFields confirms ReadCSV does NOT skip rows whose
// first field is an empty string (only ReadGradesCSV does that).
func TestReadCSV_RowWithAllEmptyFields(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "empty_fields.csv")

	// A header followed by rows where first field is empty — exercises the
	// skip-empty-row logic in ReadGradesCSV (row[0] == "").
	content := "name,age\n,30\nJane,25\n"
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	records, err := ReadCSV(p)
	if err != nil {
		t.Fatalf("ReadCSV failed: %v", err)
	}
	// ReadCSV only skips len(row)==0, not empty string rows. Both rows are present.
	if len(records) != 2 {
		t.Errorf("expected 2 records, got %d", len(records))
	}
}

// TestReadGradesCSV_EmptyFirstField verifies that ReadGradesCSV skips rows
// where the first field (user_id) is an empty string.
// A CSV row ",456,95" has row[0]=="" which triggers the continue branch.
func TestReadGradesCSV_EmptyFirstField(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "empty_first.csv")

	// Row with empty user_id → row[0] == "" → should be skipped.
	content := "user_id,assignment_id,grade\n,456,95\n123,456,87\n"
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	records, err := ReadGradesCSV(p)
	if err != nil {
		t.Fatalf("ReadGradesCSV failed: %v", err)
	}
	// Only the valid row (123,456,87) should be returned; the empty-user_id row is skipped.
	if len(records) != 1 {
		t.Errorf("expected 1 record (empty first field skipped), got %d", len(records))
	}
	if records[0].UserID != 123 {
		t.Errorf("expected UserID 123, got %d", records[0].UserID)
	}
}
