package batch

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

// GradeRecord represents a single grade entry from CSV
type GradeRecord struct {
	UserID       int64
	AssignmentID int64
	Grade        string
	Comment      string
	Row          int // Original row number for error reporting
}

// ReadGradesCSV reads grades from a CSV file
// Expected format: user_id,assignment_id,grade,comment
func ReadGradesCSV(filename string) ([]GradeRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Validate header
	if len(header) < 3 {
		return nil, fmt.Errorf("invalid CSV format: expected at least 3 columns (user_id, assignment_id, grade)")
	}

	// Read records
	var records []GradeRecord
	rowNum := 1 // Start at 1 (header is row 0)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row %d: %w", rowNum+1, err)
		}

		rowNum++

		// Skip empty rows
		if len(row) == 0 || row[0] == "" {
			continue
		}

		// Parse user ID
		userID, err := strconv.ParseInt(row[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid user_id at row %d: %w", rowNum, err)
		}

		// Parse assignment ID
		assignmentID, err := strconv.ParseInt(row[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid assignment_id at row %d: %w", rowNum, err)
		}

		// Get grade
		grade := row[2]

		// Get comment (optional)
		comment := ""
		if len(row) > 3 {
			comment = row[3]
		}

		records = append(records, GradeRecord{
			UserID:       userID,
			AssignmentID: assignmentID,
			Grade:        grade,
			Comment:      comment,
			Row:          rowNum,
		})
	}

	return records, nil
}

// WriteGradesCSV writes grades to a CSV file with secure permissions (0600)
func WriteGradesCSV(filename string, records []GradeRecord) error {
	// Create file with secure permissions (owner read/write only)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"user_id", "assignment_id", "grade", "comment"}); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write records
	for _, record := range records {
		row := []string{
			strconv.FormatInt(record.UserID, 10),
			strconv.FormatInt(record.AssignmentID, 10),
			record.Grade,
			record.Comment,
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

// ExportRecord represents a generic export record
type ExportRecord map[string]string

// ReadCSV reads generic CSV data
func ReadCSV(filename string) ([]ExportRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Read records
	var records []ExportRecord

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row: %w", err)
		}

		// Skip empty rows
		if len(row) == 0 {
			continue
		}

		// Create record
		record := make(ExportRecord)
		for i, value := range row {
			if i < len(header) {
				record[header[i]] = value
			}
		}

		records = append(records, record)
	}

	return records, nil
}

// WriteCSV writes generic CSV data with secure permissions (0600)
func WriteCSV(filename string, headers []string, records []ExportRecord) error {
	// Create file with secure permissions (owner read/write only)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write records
	for _, record := range records {
		row := make([]string, len(headers))
		for i, header := range headers {
			row[i] = record[header]
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}
