# Canvas Test Data

This directory contains synthetic Canvas LMS test instance data used for testing.

## Purpose

These files provide realistic but completely synthetic data for:
- Unit tests
- Integration tests
- Example documentation
- Development without real Canvas instances

## Data Privacy

**Important**: All data in this directory is synthetic and contains:
- NO real student names or information
- NO real instructor information
- NO real course content
- NO personally identifiable information (PII)

## Structure

- `courses.json` - Sample course data
- `users.json` - Sample user profiles
- `assignments.json` - Sample assignment definitions
- `submissions.json` - Sample submission data

## Usage in Tests

```go
// Example test usage
func TestCourseList(t *testing.T) {
    data, err := os.ReadFile("testdata/canvas-test/courses.json")
    require.NoError(t, err)

    var courses []Course
    err = json.Unmarshal(data, &courses)
    require.NoError(t, err)

    // Use synthetic data for testing
}
```

## Updating Test Data

When updating test data:
1. Ensure all data remains synthetic
2. Do not use real Canvas data
3. Keep realistic structure but fake content
4. Update this README if structure changes
