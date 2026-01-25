# Output Formats

Canvas CLI supports multiple output formats to suit different use cases.

## Available Formats

### Table (Default)

Human-readable table format, ideal for terminal viewing.

```bash
canvas courses list
```

```
ID     NAME                    CODE       TERM
123    Introduction to CS      CS101      Fall 2024
456    Data Structures         CS201      Fall 2024
789    Algorithms              CS301      Fall 2024
```

### JSON

Structured JSON output, ideal for scripting and automation.

```bash
canvas courses list --output json
```

```json
[
  {
    "id": 123,
    "name": "Introduction to CS",
    "course_code": "CS101",
    "enrollment_term_id": 1
  }
]
```

### YAML

YAML format for configuration files and human-readable structured data.

```bash
canvas courses list --output yaml
```

```yaml
- id: 123
  name: Introduction to CS
  course_code: CS101
  enrollment_term_id: 1
```

### CSV

Comma-separated values for spreadsheet import.

```bash
canvas courses list --output csv
```

```csv
id,name,course_code,enrollment_term_id
123,Introduction to CS,CS101,1
456,Data Structures,CS201,1
```

## Using Output in Scripts

### Bash with jq

```bash
# Get course IDs
canvas courses list -o json | jq '.[].id'

# Filter by name
canvas courses list -o json | jq '.[] | select(.name | contains("CS"))'
```

### Piping to Files

```bash
# Export to CSV file
canvas users list --course-id 123 -o csv > users.csv

# Export to JSON file
canvas assignments list --course-id 123 -o json > assignments.json
```

### Processing with Other Tools

```bash
# Count users
canvas users list --course-id 123 -o json | jq length

# Get specific fields
canvas users list --course-id 123 -o json | jq '.[].email'
```

## Output Filtering

Canvas CLI includes built-in filtering, column selection, and sorting capabilities.

### Filter Results

Use `--filter` to search across all fields:

```bash
# Find courses containing "CS"
canvas courses list --filter "CS"

# Find assignments with "exam" in the name
canvas assignments list --course-id 123 --filter "exam"
```

The filter is case-insensitive and searches all fields.

### Select Columns

Use `--columns` to display only specific fields:

```bash
# Show only id and name
canvas courses list --columns id,name

# Show specific assignment fields
canvas assignments list --course-id 123 --columns id,name,due_at,points_possible
```

### Sort Results

Use `--sort` to order results by a field:

```bash
# Sort by name (ascending)
canvas courses list --sort name

# Sort by due date (descending with - prefix)
canvas assignments list --course-id 123 --sort -due_at

# Sort by points
canvas assignments list --course-id 123 --sort points_possible
```

### Combining Options

All filtering options can be combined:

```bash
# Find exams, show key fields, sort by due date
canvas assignments list --course-id 123 \
  --filter "exam" \
  --columns id,name,due_at,points_possible \
  --sort -due_at
```

### Works with All Formats

Filtering works with any output format:

```bash
# Filter JSON output
canvas courses list --filter "CS" --output json

# Filter and export to CSV
canvas users list --course-id 123 --filter "student" --output csv > students.csv
```

## Tips

!!! tip "Use JSON for Scripts"
    Always use `-o json` when parsing output in scripts. Table format may change between versions.

!!! tip "CSV for Spreadsheets"
    Use `-o csv` when you need to import data into Excel, Google Sheets, or other spreadsheet applications.

!!! tip "Built-in Filtering vs jq"
    For simple filtering, use built-in `--filter`. For complex queries, use JSON output with `jq`:
    ```bash
    # Simple - use built-in
    canvas courses list --filter "Fall 2024"

    # Complex - use jq
    canvas courses list -o json | jq '.[] | select(.enrollment_term_id == 5)'
    ```
