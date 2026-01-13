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

## Tips

!!! tip "Use JSON for Scripts"
    Always use `-o json` when parsing output in scripts. Table format may change between versions.

!!! tip "CSV for Spreadsheets"
    Use `-o csv` when you need to import data into Excel, Google Sheets, or other spreadsheet applications.

!!! tip "Combine with grep"
    For quick filtering, combine table output with grep:
    ```bash
    canvas courses list | grep "Fall 2024"
    ```
