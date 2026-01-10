# Usage Examples

Practical examples for common Canvas CLI workflows.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Course Management](#course-management)
3. [Assignment & Grading](#assignment--grading)
4. [Bulk Operations](#bulk-operations)
5. [User Management](#user-management)
6. [File Management](#file-management)
7. [Multi-Instance Workflows](#multi-instance-workflows)
8. [Automation & Scripting](#automation--scripting)

## Getting Started

### Initial Setup

```bash
# Install Canvas CLI
brew install canvas-cli

# Authenticate
canvas auth login --instance https://canvas.instructure.com

# Verify authentication
canvas auth status

# Test with a simple command
canvas courses list
```

### Enable Shell Completion

```bash
# Bash
canvas completion bash > ~/.canvas-completion.bash
echo 'source ~/.canvas-completion.bash' >> ~/.bashrc

# Zsh
canvas completion zsh > "${fpath[1]}/_canvas"

# Fish
canvas completion fish > ~/.config/fish/completions/canvas.fish
```

## Course Management

### List All Your Courses

```bash
# List active courses
canvas courses list --state available

# List as teacher
canvas courses list --enrollment-type teacher

# Include student count
canvas courses list --include total_students --format table
```

### Get Course Details

```bash
# Basic details
canvas courses get 12345

# Include enrolled teachers and students
canvas courses get 12345 --include teachers,students --format json
```

### List Course Users

```bash
# List all students in a course
canvas courses users 12345 --enrollment-type student

# Export to CSV
canvas courses users 12345 --enrollment-type student --format csv > students.csv
```

## Assignment & Grading

### List Assignments

```bash
# List all assignments for a course
canvas assignments list --course 12345

# Include submission data
canvas assignments list --course 12345 --include submission --format json
```

### View Assignment Submissions

```bash
# List submissions
canvas assignments submissions --course 12345 --assignment 67890

# Include user details
canvas assignments submissions --course 12345 --assignment 67890 --include user
```

### Grade Single Submission

```bash
# Grade with score
canvas assignments grade --course 12345 --assignment 67890 --user 11111 --score 95

# Grade with score and comment
canvas assignments grade --course 12345 --assignment 67890 --user 11111 \
  --score 95 --comment "Excellent work! Full credit."
```

### Bulk Grading from CSV

Create a CSV file (`grades.csv`):
```csv
student_id,score,comment
11111,95,Great work
22222,87,Good job
33333,92,Excellent
```

Then import:
```bash
canvas assignments grade --course 12345 --assignment 67890 --csv grades.csv
```

### Export Grades to CSV

```bash
# Export all submissions with scores
canvas assignments submissions --course 12345 --assignment 67890 \
  --include user --format csv > grades_export.csv
```

## Bulk Operations

### Grade Multiple Students

```bash
# Create grades CSV
cat > grades.csv << EOF
student_id,score,comment
11111,95,Excellent
22222,88,Very good
33333,92,Great work
44444,85,Good
EOF

# Import grades
canvas assignments grade --course 12345 --assignment 67890 --csv grades.csv
```

### Enroll Multiple Students

```bash
# Create enrollments from a list
while IFS=, read -r user_id; do
  canvas enrollments create --course 12345 --user "$user_id" --type student
done < student_ids.txt
```

### Download All Course Files

```bash
# List files and download each
canvas files list --course 12345 --format json | \
  jq -r '.[] | .id' | \
  while read file_id; do
    canvas files download "$file_id"
  done
```

## User Management

### Create Multiple Users

```bash
# From CSV
cat > users.csv << EOF
name,email,login
John Doe,john@example.com,jdoe
Jane Smith,jane@example.com,jsmith
EOF

# Import users
while IFS=, read -r name email login; do
  canvas users create --name "$name" --email "$email" --login "$login"
done < users.csv
```

### Update User Information

```bash
# Update single user
canvas users update 12345 --email newemail@example.com

# Bulk update from list
while IFS=, read -r user_id email; do
  canvas users update "$user_id" --email "$email"
done < email_updates.csv
```

### List All Students with Enrollment Status

```bash
# Get all enrollments for a course
canvas enrollments list 12345 --type student --format json | \
  jq -r '.[] | "\(.user.name),\(.user.email),\(.enrollment_state)"'
```

## File Management

### Upload Course Materials

```bash
# Upload single file
canvas files upload syllabus.pdf --course 12345 --folder "Course Documents"

# Upload multiple files
for file in *.pdf; do
  canvas files upload "$file" --course 12345 --folder "Readings"
done
```

### Organize Files in Folders

```bash
# Upload to specific folder path
canvas files upload lecture01.pdf --course 12345 --folder "Lectures/Week 1"
canvas files upload lecture02.pdf --course 12345 --folder "Lectures/Week 2"
```

### Download Student Submissions

```bash
# Get submission file URLs and download
canvas assignments submissions --course 12345 --assignment 67890 \
  --include submission --format json | \
  jq -r '.[] | select(.attachments != null) | .attachments[].url' | \
  while read url; do
    wget "$url"
  done
```

## Multi-Instance Workflows

### Set Up Multiple Instances

```bash
# Add production instance
canvas config add production --url https://canvas.instructure.com
canvas auth login --instance production

# Add staging instance
canvas config add staging --url https://staging.canvas.com
canvas auth login --instance staging

# Add self-hosted instance
canvas config add onprem --url https://canvas.company.com
canvas auth login --instance onprem
```

### Copy Course Between Instances

```bash
# Copy course from production to staging
canvas sync copy \
  --from production \
  --to staging \
  --course 12345 \
  --assignments \
  --files
```

### Compare Courses Across Instances

```bash
# List courses from production
canvas courses list --instance production --format json > prod_courses.json

# List courses from staging
canvas courses list --instance staging --format json > staging_courses.json

# Compare (requires jq)
diff <(jq -S . prod_courses.json) <(jq -S . staging_courses.json)
```

## Automation & Scripting

### Daily Grade Export

```bash
#!/bin/bash
# daily_export.sh - Export grades daily

COURSE_ID=12345
ASSIGNMENT_ID=67890
DATE=$(date +%Y-%m-%d)

canvas assignments submissions \
  --course $COURSE_ID \
  --assignment $ASSIGNMENT_ID \
  --include user \
  --format csv > "grades_${DATE}.csv"

echo "Grades exported to grades_${DATE}.csv"
```

### Automated Enrollment Sync

```bash
#!/bin/bash
# sync_enrollments.sh - Sync enrollments from SIS

SIS_EXPORT="sis_enrollments.csv"

while IFS=, read -r course_id user_id type; do
  echo "Enrolling user $user_id in course $course_id as $type"
  canvas enrollments create \
    --course "$course_id" \
    --user "$user_id" \
    --type "$type" \
    --notify
done < "$SIS_EXPORT"
```

### Monitor Assignment Submissions

```bash
#!/bin/bash
# monitor_submissions.sh - Check for new submissions

COURSE_ID=12345
ASSIGNMENT_ID=67890
LAST_CHECK=$(date -u -d '1 hour ago' +%Y-%m-%dT%H:%M:%SZ)

canvas assignments submissions \
  --course $COURSE_ID \
  --assignment $ASSIGNMENT_ID \
  --format json | \
  jq -r ".[] | select(.submitted_at > \"$LAST_CHECK\") |
    \"New submission from \(.user.name) at \(.submitted_at)\""
```

### Batch File Upload

```bash
#!/bin/bash
# batch_upload.sh - Upload all files in a directory

COURSE_ID=12345
SOURCE_DIR="./course_materials"
DEST_FOLDER="Course Materials"

find "$SOURCE_DIR" -type f | while read file; do
  echo "Uploading: $file"
  canvas files upload "$file" \
    --course $COURSE_ID \
    --folder "$DEST_FOLDER" \
    --on-duplicate rename
done

echo "Upload complete"
```

### Weekly Status Report

```bash
#!/bin/bash
# weekly_report.sh - Generate weekly course status report

COURSE_ID=12345
REPORT_FILE="weekly_report_$(date +%Y%m%d).txt"

{
  echo "Weekly Course Report - $(date)"
  echo "================================"
  echo ""

  echo "Total Students:"
  canvas courses users $COURSE_ID --enrollment-type student --format json | jq '. | length'
  echo ""

  echo "Active Assignments:"
  canvas assignments list --course $COURSE_ID --format json | jq '. | length'
  echo ""

  echo "Recent Submissions (last 7 days):"
  canvas assignments submissions --course $COURSE_ID --assignment ALL --format json | \
    jq '[.[] | select(.submitted_at != null)] | length'

} > "$REPORT_FILE"

echo "Report saved to $REPORT_FILE"
```

## Interactive REPL Mode

### Using REPL for Exploration

```bash
# Start REPL
canvas shell

# Inside REPL (no 'canvas' prefix needed):
canvas> courses list
canvas> courses get 12345
canvas> assignments list --course 12345
canvas> exit
```

### REPL with Tab Completion

```bash
# Start REPL
canvas shell

# Press TAB to autocomplete:
canvas> courses <TAB>
  get   list   users

# Autocomplete flags:
canvas> courses list --<TAB>
  --state    --enrollment-type    --include    --format
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Export Grades

on:
  schedule:
    - cron: '0 0 * * 0'  # Weekly on Sunday

jobs:
  export:
    runs-on: ubuntu-latest
    steps:
      - name: Install Canvas CLI
        run: |
          curl -sSL https://github.com/jjuanrivvera/canvas-cli/releases/download/v1.0.0/canvas_Linux_x86_64.tar.gz | tar xz
          sudo mv canvas /usr/local/bin/

      - name: Export grades
        env:
          CANVAS_URL: ${{ secrets.CANVAS_URL }}
          CANVAS_TOKEN: ${{ secrets.CANVAS_TOKEN }}
        run: |
          canvas assignments submissions \
            --course 12345 \
            --assignment 67890 \
            --format csv > grades.csv

      - name: Upload artifact
        uses: actions/upload-artifact@v3
        with:
          name: grades
          path: grades.csv
```

### GitLab CI Example

```yaml
export_grades:
  image: golang:1.21
  script:
    - go install github.com/jjuanrivvera/canvas-cli/cmd/canvas@latest
    - canvas assignments submissions --course 12345 --assignment 67890 --format csv > grades.csv
  artifacts:
    paths:
      - grades.csv
  only:
    - schedules
```

## Troubleshooting Examples

### Check Authentication

```bash
# Verify auth status
canvas auth status

# Re-authenticate if needed
canvas auth logout
canvas auth login --instance https://canvas.instructure.com
```

### Debug API Issues

```bash
# Enable debug mode
canvas --debug courses list

# Check diagnostics
canvas doctor

# Clear cache if data seems stale
canvas cache clear --all
```

### Handle Rate Limiting

```bash
# Add delays between bulk operations
for id in $(cat student_ids.txt); do
  canvas enrollments create --course 12345 --user "$id" --type student
  sleep 1  # Wait 1 second between requests
done
```

## Next Steps

- Review [Command Reference](COMMANDS.md) for complete command documentation
- Check [Authentication Guide](AUTHENTICATION.md) for OAuth setup details
- Run `canvas doctor` to verify your installation
