# Canvas CLI Best Practices

This guide covers best practices for using Canvas CLI effectively and efficiently.

## Table of Contents

- [Initial Setup](#initial-setup)
- [Authentication](#authentication)
- [Working with Multiple Instances](#working-with-multiple-instances)
- [Productivity Features](#productivity-features)
- [Output and Filtering](#output-and-filtering)
- [Scripting and Automation](#scripting-and-automation)
- [Common Workflows](#common-workflows)
- [Performance Tips](#performance-tips)
- [Troubleshooting](#troubleshooting)

---

## Initial Setup

### 1. Install Shell Completion

Enable tab completion for faster command entry:

```bash
# Bash
canvas completion bash > /etc/bash_completion.d/canvas

# Zsh
canvas completion zsh > "${fpath[1]}/_canvas"

# Fish
canvas completion fish > ~/.config/fish/completions/canvas.fish
```

### 2. Run Diagnostics

Verify your installation is working correctly:

```bash
canvas doctor
```

This checks connectivity, authentication, and configuration.

### 3. Enable Auto-Updates

Stay current with the latest features and fixes:

```bash
canvas update enable
canvas update status
```

---

## Authentication

### Choose the Right Method

| Method | Use Case | Security |
|--------|----------|----------|
| OAuth (default) | Interactive use, personal accounts | High - tokens auto-refresh |
| API Token | Scripts, automation, CI/CD | Medium - store securely |

### OAuth Login (Recommended for Interactive Use)

```bash
# Login with browser-based OAuth
canvas auth login https://canvas.example.com

# Login to a named instance
canvas config add prod --url https://canvas.example.com
canvas auth login --instance prod
```

### API Token (For Automation)

```bash
# Set API token for an instance
canvas auth token set myinstance

# Use in scripts (set environment variable)
export CANVAS_TOKEN="your-token-here"
```

### Check Authentication Status

```bash
canvas auth status
```

---

## Working with Multiple Instances

### Configure Named Instances

```bash
# Add instances
canvas config add prod --url https://canvas.example.com
canvas config add staging --url https://staging.canvas.example.com
canvas config add dev --url https://dev.canvas.example.com

# List all instances
canvas config list

# Switch between instances
canvas config use prod
canvas config use staging
```

### Per-Command Instance Override

```bash
# Use a specific instance for one command
canvas courses list --instance staging
```

### Set Default Account

For admin operations, set a default account ID:

```bash
canvas config account 1
```

---

## Productivity Features

### Command Aliases

Create shortcuts for frequently used commands:

```bash
# Create aliases
canvas alias set courses "courses list"
canvas alias set hw "assignments list --course-id 123"
canvas alias set ungraded "submissions list --workflow-state submitted"
canvas alias set students "users list --enrollment-type student"

# Use aliases
canvas courses
canvas hw
canvas ungraded --course-id 456 --assignment-id 789

# Manage aliases
canvas alias list
canvas alias delete hw
```

**Best Practices for Aliases:**

- Use short, memorable names
- Include common flags you always use
- Don't include IDs that change frequently (use context instead)

### Context Management

Set default values to avoid repetitive typing:

```bash
# Set course context for a grading session
canvas context set course 12345
canvas context set assignment 67890

# Now these commands use context automatically
canvas submissions list          # Uses course 12345, assignment 67890
canvas submissions grade 111 --grade 95

# View current context
canvas context show

# Clear when switching tasks
canvas context clear
```

**Best Practices for Context:**

- Set context at the start of a focused work session
- Clear context when switching courses/tasks
- Use explicit flags when working across multiple courses
- Context + aliases = maximum efficiency

### Combine Aliases with Context

```bash
# Set up your workflow
canvas context set course 12345
canvas alias set subs "submissions list"
canvas alias set grade "submissions grade"

# Super efficient grading
canvas subs --assignment-id 456
canvas grade 111 --grade 95 --comment "Great work!"
canvas grade 222 --grade 88
```

---

## Output and Filtering

### Choose the Right Output Format

| Format | Use Case |
|--------|----------|
| `table` | Human reading in terminal (default) |
| `json` | Scripting, piping to jq, automation |
| `yaml` | Configuration files, readable structured data |
| `csv` | Spreadsheet import, Excel, Google Sheets |

```bash
# Examples
canvas courses list                    # Table for viewing
canvas courses list -o json            # JSON for scripts
canvas users list -o csv > users.csv   # CSV for spreadsheets
```

### Filter Results

```bash
# Text filter (case-insensitive, searches all fields)
canvas courses list --filter "Fall 2024"
canvas users list --course-id 123 --filter "student"

# Select specific columns
canvas assignments list --course-id 123 --columns id,name,due_at

# Sort results
canvas assignments list --course-id 123 --sort due_at      # Ascending
canvas assignments list --course-id 123 --sort -due_at     # Descending

# Combine all options
canvas assignments list --course-id 123 \
  --filter "exam" \
  --columns id,name,due_at,points_possible \
  --sort -due_at
```

### Limit Results

```bash
# Get only first 10 results
canvas courses list --limit 10

# Useful for testing or quick checks
canvas submissions list --course-id 123 --assignment-id 456 --limit 5
```

---

## Scripting and Automation

### Use JSON Output

Always use `-o json` in scripts for reliable parsing:

```bash
#!/bin/bash
# Get all course IDs
COURSES=$(canvas courses list -o json | jq -r '.[].id')

for COURSE_ID in $COURSES; do
    echo "Processing course $COURSE_ID"
    canvas assignments list --course-id "$COURSE_ID" -o json
done
```

### Dry-Run for Testing

Preview commands before executing:

```bash
# See what API calls would be made
canvas assignments create --course-id 123 --name "Test" --dry-run

# Shows curl command with redacted token
curl -X POST 'https://canvas.example.com/api/v1/courses/123/assignments' \
  -H 'Authorization: Bearer [REDACTED]' \
  ...
```

### Bulk Operations

Use CSV for bulk grading:

```bash
# Prepare grades.csv:
# student_id,grade,comment
# 123,95,Great work!
# 456,88,Good effort

canvas submissions bulk-grade \
  --course-id 123 \
  --assignment-id 456 \
  --csv grades.csv
```

### Error Handling in Scripts

```bash
#!/bin/bash
set -e  # Exit on error

# Check authentication first
if ! canvas auth status > /dev/null 2>&1; then
    echo "Not authenticated. Run: canvas auth login"
    exit 1
fi

# Proceed with operations
canvas courses list -o json
```

---

## Common Workflows

### Grading Workflow

```bash
# 1. Set context for the grading session
canvas context set course 12345
canvas context set assignment 67890

# 2. List ungraded submissions
canvas submissions list --workflow-state submitted

# 3. Grade individual submissions
canvas submissions grade 111 --grade 95 --comment "Excellent!"
canvas submissions grade 222 --grade 88

# 4. Or bulk grade from CSV
canvas submissions bulk-grade --csv grades.csv

# 5. Clear context when done
canvas context clear
```

### Course Setup Workflow

```bash
# 1. Create assignment groups
canvas assignment-groups create --course-id 123 --name "Homework" --weight 30
canvas assignment-groups create --course-id 123 --name "Exams" --weight 50
canvas assignment-groups create --course-id 123 --name "Projects" --weight 20

# 2. Create assignments
canvas assignments create --course-id 123 \
  --name "Homework 1" \
  --assignment-group-id 456 \
  --points 100 \
  --due-at "2024-09-15T23:59:00Z"

# 3. Create modules
canvas modules create --course-id 123 --name "Week 1"
canvas modules items create --course-id 123 --module-id 789 \
  --type Assignment --content-id 456
```

### User Management Workflow

```bash
# List students in a course
canvas users list --course-id 123 --enrollment-type student

# Search for a user
canvas users search --course-id 123 --search-term "john"

# Export to spreadsheet
canvas users list --course-id 123 -o csv > students.csv
```

### Course Migration Workflow

```bash
# Sync from source to destination
canvas sync course \
  --source-instance prod \
  --source-course 123 \
  --dest-instance staging \
  --dest-course 456 \
  --interactive
```

---

## Performance Tips

### Enable Caching

Caching is enabled by default. Manage it as needed:

```bash
# View cache stats
canvas cache stats

# Clear cache when data is stale
canvas cache clear

# Disable cache for one command
canvas courses list --no-cache
```

### Use Pagination Wisely

```bash
# For large datasets, use --limit to paginate manually
canvas users list --course-id 123 --limit 100

# Or let the CLI handle it (may take time for large datasets)
canvas users list --course-id 123
```

### Batch Operations

For multiple similar operations, use bulk commands when available:

```bash
# Slow: Individual grade commands
for id in 1 2 3 4 5; do
    canvas submissions grade $id --grade 100
done

# Fast: Bulk grade from CSV
canvas submissions bulk-grade --csv grades.csv
```

---

## Troubleshooting

### Run Diagnostics

```bash
canvas doctor
```

### Enable Verbose Output

```bash
canvas courses list -v
```

### Check API Calls

```bash
# See exact API request
canvas courses list --dry-run
```

### Common Issues

| Issue | Solution |
|-------|----------|
| "Not authenticated" | Run `canvas auth login` |
| "Rate limited" | Wait and retry, or reduce request frequency |
| "Course not found" | Check course ID and permissions |
| "Token expired" | Re-authenticate with `canvas auth login` |

### Clear State

```bash
# Clear all cached data
canvas cache clear

# Clear context
canvas context clear

# Re-authenticate
canvas auth logout
canvas auth login
```

---

## Quick Reference

### Essential Commands

```bash
canvas auth login                 # Authenticate
canvas courses list               # List courses
canvas assignments list -c 123    # List assignments
canvas submissions list -c 123 -a 456  # List submissions
canvas context set course 123     # Set context
canvas alias set x "command"      # Create alias
```

### Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output format (table/json/yaml/csv) |
| `--verbose` | `-v` | Show detailed output |
| `--filter` | | Filter results by text |
| `--columns` | | Select columns to display |
| `--sort` | | Sort by field (- for descending) |
| `--limit` | | Limit number of results |
| `--dry-run` | | Show curl command without executing |
| `--no-cache` | | Bypass cache |
| `--instance` | | Use specific Canvas instance |

### Getting Help

```bash
canvas --help              # General help
canvas <command> --help    # Command-specific help
canvas doctor              # Diagnose issues
```
