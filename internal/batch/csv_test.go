package batch

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadGradesCSV(t *testing.T) {
	// Create temporary CSV file
	tempDir := t.TempDir()
	csvPath := filepath.Join(tempDir, "grades.csv")

	csvContent := `user_id,assignment_id,grade,comment
123,456,95,Great work
789,456,87,Good job
456,456,78,Needs improvement`

	err := os.WriteFile(csvPath, []byte(csvContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	records, err := ReadGradesCSV(csvPath)
	if err != nil {
		t.Fatalf("ReadGradesCSV failed: %v", err)
	}

	if len(records) != 3 {
		t.Errorf("expected 3 records, got %d", len(records))
	}

	// Check first record
	if records[0].UserID != 123 {
		t.Errorf("expected UserID 123, got %d", records[0].UserID)
	}

	if records[0].AssignmentID != 456 {
		t.Errorf("expected AssignmentID 456, got %d", records[0].AssignmentID)
	}

	if records[0].Grade != "95" {
		t.Errorf("expected Grade '95', got '%s'", records[0].Grade)
	}

	if records[0].Comment != "Great work" {
		t.Errorf("expected Comment 'Great work', got '%s'", records[0].Comment)
	}
}

func TestReadGradesCSV_InvalidFile(t *testing.T) {
	_, err := ReadGradesCSV("/nonexistent/path/file.csv")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestReadGradesCSV_InvalidFormat(t *testing.T) {
	tempDir := t.TempDir()
	csvPath := filepath.Join(tempDir, "invalid.csv")

	// CSV with only 2 columns (need at least 3)
	csvContent := `col1,col2
value1,value2`

	err := os.WriteFile(csvPath, []byte(csvContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	_, err = ReadGradesCSV(csvPath)
	if err == nil {
		t.Error("expected error for invalid CSV format")
	}
}

func TestWriteGradesCSV(t *testing.T) {
	tempDir := t.TempDir()
	csvPath := filepath.Join(tempDir, "output.csv")

	records := []GradeRecord{
		{UserID: 123, AssignmentID: 456, Grade: "95", Comment: "Great"},
		{UserID: 789, AssignmentID: 456, Grade: "87", Comment: "Good"},
	}

	err := WriteGradesCSV(csvPath, records)
	if err != nil {
		t.Fatalf("WriteGradesCSV failed: %v", err)
	}

	// Read back and verify
	readRecords, err := ReadGradesCSV(csvPath)
	if err != nil {
		t.Fatalf("Failed to read back CSV: %v", err)
	}

	if len(readRecords) != 2 {
		t.Errorf("expected 2 records, got %d", len(readRecords))
	}

	if readRecords[0].UserID != 123 {
		t.Errorf("expected UserID 123, got %d", readRecords[0].UserID)
	}
}

func TestReadCSV(t *testing.T) {
	tempDir := t.TempDir()
	csvPath := filepath.Join(tempDir, "data.csv")

	csvContent := `name,age,city
John,30,NYC
Jane,25,LA
Bob,35,SF`

	err := os.WriteFile(csvPath, []byte(csvContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	records, err := ReadCSV(csvPath)
	if err != nil {
		t.Fatalf("ReadCSV failed: %v", err)
	}

	if len(records) != 3 {
		t.Errorf("expected 3 records, got %d", len(records))
	}

	if records[0]["name"] != "John" {
		t.Errorf("expected name 'John', got '%s'", records[0]["name"])
	}

	if records[0]["age"] != "30" {
		t.Errorf("expected age '30', got '%s'", records[0]["age"])
	}

	if records[0]["city"] != "NYC" {
		t.Errorf("expected city 'NYC', got '%s'", records[0]["city"])
	}
}

func TestWriteCSV(t *testing.T) {
	tempDir := t.TempDir()
	csvPath := filepath.Join(tempDir, "output.csv")

	headers := []string{"name", "age", "city"}
	records := []ExportRecord{
		{"name": "John", "age": "30", "city": "NYC"},
		{"name": "Jane", "age": "25", "city": "LA"},
	}

	err := WriteCSV(csvPath, headers, records)
	if err != nil {
		t.Fatalf("WriteCSV failed: %v", err)
	}

	// Read back and verify
	readRecords, err := ReadCSV(csvPath)
	if err != nil {
		t.Fatalf("Failed to read back CSV: %v", err)
	}

	if len(readRecords) != 2 {
		t.Errorf("expected 2 records, got %d", len(readRecords))
	}

	if readRecords[0]["name"] != "John" {
		t.Errorf("expected name 'John', got '%s'", readRecords[0]["name"])
	}
}

func TestReadGradesCSV_WithoutComment(t *testing.T) {
	tempDir := t.TempDir()
	csvPath := filepath.Join(tempDir, "grades_no_comment.csv")

	csvContent := `user_id,assignment_id,grade
123,456,95
789,456,87`

	err := os.WriteFile(csvPath, []byte(csvContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	records, err := ReadGradesCSV(csvPath)
	if err != nil {
		t.Fatalf("ReadGradesCSV failed: %v", err)
	}

	if len(records) != 2 {
		t.Errorf("expected 2 records, got %d", len(records))
	}

	// Comment should be empty
	if records[0].Comment != "" {
		t.Errorf("expected empty comment, got '%s'", records[0].Comment)
	}
}

func TestReadGradesCSV_InvalidUserID(t *testing.T) {
	tempDir := t.TempDir()
	csvPath := filepath.Join(tempDir, "invalid_user.csv")

	csvContent := `user_id,assignment_id,grade
invalid,456,95`

	err := os.WriteFile(csvPath, []byte(csvContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	_, err = ReadGradesCSV(csvPath)
	if err == nil {
		t.Error("expected error for invalid user_id")
	}
}

func TestReadGradesCSV_InvalidAssignmentID(t *testing.T) {
	tempDir := t.TempDir()
	csvPath := filepath.Join(tempDir, "invalid_assignment.csv")

	csvContent := `user_id,assignment_id,grade
123,invalid,95`

	err := os.WriteFile(csvPath, []byte(csvContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	_, err = ReadGradesCSV(csvPath)
	if err == nil {
		t.Error("expected error for invalid assignment_id")
	}
}

func TestReadGradesCSV_EmptyRows(t *testing.T) {
	tempDir := t.TempDir()
	csvPath := filepath.Join(tempDir, "empty_rows.csv")

	csvContent := `user_id,assignment_id,grade
123,456,95

789,456,87`

	err := os.WriteFile(csvPath, []byte(csvContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	records, err := ReadGradesCSV(csvPath)
	if err != nil {
		t.Fatalf("ReadGradesCSV failed: %v", err)
	}

	// Should skip empty row
	if len(records) != 2 {
		t.Errorf("expected 2 records (empty row skipped), got %d", len(records))
	}
}

func TestWriteGradesCSV_EmptyRecords(t *testing.T) {
	tempDir := t.TempDir()
	csvPath := filepath.Join(tempDir, "empty.csv")

	records := []GradeRecord{}

	err := WriteGradesCSV(csvPath, records)
	if err != nil {
		t.Fatalf("WriteGradesCSV failed: %v", err)
	}

	// File should exist with just header
	content, err := os.ReadFile(csvPath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if len(content) == 0 {
		t.Error("expected file to contain header")
	}
}

func TestWriteGradesCSV_InvalidPath(t *testing.T) {
	records := []GradeRecord{
		{UserID: 123, AssignmentID: 456, Grade: "95"},
	}

	err := WriteGradesCSV("/nonexistent/path/file.csv", records)
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestReadCSV_InvalidFile(t *testing.T) {
	_, err := ReadCSV("/nonexistent/path/file.csv")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestReadCSV_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	csvPath := filepath.Join(tempDir, "empty.csv")

	err := os.WriteFile(csvPath, []byte(""), 0600)
	if err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	_, err = ReadCSV(csvPath)
	if err == nil {
		t.Error("expected error for empty file (no header)")
	}
}

func TestReadCSV_EmptyRows(t *testing.T) {
	tempDir := t.TempDir()
	csvPath := filepath.Join(tempDir, "with_empty.csv")

	csvContent := `name,age
John,30

Jane,25`

	err := os.WriteFile(csvPath, []byte(csvContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	records, err := ReadCSV(csvPath)
	if err != nil {
		t.Fatalf("ReadCSV failed: %v", err)
	}

	// Should skip empty row
	if len(records) != 2 {
		t.Errorf("expected 2 records (empty row skipped), got %d", len(records))
	}
}

func TestWriteCSV_EmptyHeaders(t *testing.T) {
	tempDir := t.TempDir()
	csvPath := filepath.Join(tempDir, "no_headers.csv")

	headers := []string{}
	records := []ExportRecord{
		{"name": "John"},
	}

	err := WriteCSV(csvPath, headers, records)
	if err != nil {
		t.Fatalf("WriteCSV failed: %v", err)
	}

	// Should create file with empty header row
	content, err := os.ReadFile(csvPath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if len(content) == 0 {
		t.Error("expected file to exist")
	}
}

func TestWriteCSV_EmptyRecords(t *testing.T) {
	tempDir := t.TempDir()
	csvPath := filepath.Join(tempDir, "no_records.csv")

	headers := []string{"name", "age"}
	records := []ExportRecord{}

	err := WriteCSV(csvPath, headers, records)
	if err != nil {
		t.Fatalf("WriteCSV failed: %v", err)
	}

	// Should create file with just header
	content, err := os.ReadFile(csvPath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if len(content) == 0 {
		t.Error("expected file to contain header")
	}
}

func TestWriteCSV_InvalidPath(t *testing.T) {
	headers := []string{"name"}
	records := []ExportRecord{
		{"name": "John"},
	}

	err := WriteCSV("/nonexistent/path/file.csv", headers, records)
	if err == nil {
		t.Error("expected error for invalid path")
	}
}
