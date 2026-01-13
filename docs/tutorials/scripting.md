# Scripting & Automation Tutorial

Learn how to automate Canvas tasks with shell scripts.

## Overview

Canvas CLI's JSON output makes it perfect for automation:

- Integrate with shell scripts and cron jobs
- Process data with `jq`
- Build custom workflows

## Prerequisites

- Canvas CLI installed and authenticated
- Basic shell scripting knowledge
- `jq` installed (for JSON processing)

## Basic Scripting

### Using JSON Output

Always use `-o json` for scripting:

```bash
# Get all course IDs
canvas courses list -o json | jq '.[].id'

# Get courses as array
courses=$(canvas courses list -o json | jq -r '.[].id')
for course in $courses; do
  echo "Processing course: $course"
done
```

### Filtering with jq

```bash
# Find active courses
canvas courses list -o json | jq '.[] | select(.workflow_state == "available")'

# Get course names containing "CS"
canvas courses list -o json | jq '.[] | select(.name | contains("CS")) | .name'

# Count enrollments
canvas users list --course-id 123 -o json | jq length
```

## Example Scripts

### Export All Grades

Export grades for all assignments in a course:

```bash
#!/bin/bash
COURSE_ID=$1

if [ -z "$COURSE_ID" ]; then
  echo "Usage: $0 <course_id>"
  exit 1
fi

# Get all assignments
assignments=$(canvas assignments list --course-id $COURSE_ID -o json | jq -r '.[].id')

# Create output directory
mkdir -p grades/$COURSE_ID

# Export each assignment's submissions
for assignment in $assignments; do
  echo "Exporting assignment $assignment..."
  canvas submissions list \
    --course-id $COURSE_ID \
    --assignment-id $assignment \
    -o csv > "grades/$COURSE_ID/assignment_$assignment.csv"
done

echo "Done! Grades exported to grades/$COURSE_ID/"
```

### Bulk User Enrollment

Enroll users from a CSV file:

```bash
#!/bin/bash
COURSE_ID=$1
CSV_FILE=$2

if [ -z "$COURSE_ID" ] || [ -z "$CSV_FILE" ]; then
  echo "Usage: $0 <course_id> <csv_file>"
  exit 1
fi

# Skip header and process each line
tail -n +2 "$CSV_FILE" | while IFS=, read -r email role; do
  echo "Enrolling $email as $role..."
  canvas enrollments create \
    --course-id $COURSE_ID \
    --user-email "$email" \
    --role "$role"
done
```

### Course Health Check

Check course configuration and report issues:

```bash
#!/bin/bash
COURSE_ID=$1

echo "=== Course Health Check ==="
echo ""

# Get course info
course=$(canvas courses get $COURSE_ID -o json)
echo "Course: $(echo $course | jq -r '.name')"
echo "State: $(echo $course | jq -r '.workflow_state')"
echo ""

# Check modules
module_count=$(canvas modules list --course-id $COURSE_ID -o json | jq length)
echo "Modules: $module_count"

# Check assignments
assignment_count=$(canvas assignments list --course-id $COURSE_ID -o json | jq length)
echo "Assignments: $assignment_count"

# Check unpublished items
unpublished=$(canvas modules list --course-id $COURSE_ID -o json | jq '[.[] | select(.published == false)] | length')
echo "Unpublished modules: $unpublished"

# Check for missing due dates
no_due_date=$(canvas assignments list --course-id $COURSE_ID -o json | jq '[.[] | select(.due_at == null)] | length')
echo "Assignments without due date: $no_due_date"
```

### Automated Backup

Back up course content daily:

```bash
#!/bin/bash
# Add to crontab: 0 2 * * * /path/to/backup.sh

BACKUP_DIR="/backups/canvas"
DATE=$(date +%Y-%m-%d)

# Get all courses
courses=$(canvas courses list -o json | jq -r '.[].id')

for course_id in $courses; do
  course_name=$(canvas courses get $course_id -o json | jq -r '.name' | tr ' ' '_')
  output_dir="$BACKUP_DIR/$DATE/$course_name"
  mkdir -p "$output_dir"

  # Export course data
  canvas courses get $course_id -o json > "$output_dir/course.json"
  canvas modules list --course-id $course_id -o json > "$output_dir/modules.json"
  canvas assignments list --course-id $course_id -o json > "$output_dir/assignments.json"
  canvas pages list --course-id $course_id -o json > "$output_dir/pages.json"

  echo "Backed up: $course_name"
done

# Cleanup old backups (keep 30 days)
find $BACKUP_DIR -type d -mtime +30 -exec rm -rf {} \;
```

## Advanced Patterns

### Parallel Processing

Process multiple courses in parallel:

```bash
#!/bin/bash
# Process courses in parallel (max 4 at a time)
canvas courses list -o json | jq -r '.[].id' | \
  xargs -P 4 -I {} bash -c 'process_course {}'

process_course() {
  course_id=$1
  canvas assignments list --course-id $course_id -o json > "course_$course_id.json"
}
export -f process_course
```

### Error Handling

Robust error handling in scripts:

```bash
#!/bin/bash
set -e  # Exit on error

handle_error() {
  echo "Error on line $1"
  exit 1
}
trap 'handle_error $LINENO' ERR

# Check if canvas CLI is available
if ! command -v canvas &> /dev/null; then
  echo "Canvas CLI not found"
  exit 1
fi

# Verify authentication
if ! canvas auth status &> /dev/null; then
  echo "Not authenticated. Run: canvas auth login"
  exit 1
fi

# Your script logic here
canvas courses list -o json
```

### Environment-Based Configuration

Use different instances based on environment:

```bash
#!/bin/bash
# Set instance based on environment
case "$CANVAS_ENV" in
  production)
    INSTANCE="production"
    ;;
  staging)
    INSTANCE="sandbox"
    ;;
  *)
    INSTANCE="sandbox"  # Default to sandbox
    ;;
esac

canvas courses list --instance $INSTANCE -o json
```

## Tips

!!! tip "Test in Sandbox"
    Always test scripts against a sandbox instance before running on production.

!!! tip "Rate Limiting"
    Canvas CLI handles rate limiting automatically, but for large batch operations, consider adding delays.

!!! tip "Logging"
    Add timestamps and logging for long-running scripts:
    ```bash
    log() {
      echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
    }
    log "Starting process..."
    ```

!!! warning "Credentials"
    Never hardcode tokens in scripts. Use environment variables or the config file.
