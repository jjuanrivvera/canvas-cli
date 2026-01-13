# Bulk Grading Tutorial

Learn how to grade multiple submissions efficiently using Canvas CLI.

## Overview

Canvas CLI provides powerful bulk grading capabilities that allow you to:

- Export submissions to CSV for offline grading
- Import grades from CSV files
- Grade with comments and rubric scores

## Prerequisites

- Canvas CLI installed and authenticated
- Course ID and Assignment ID
- Appropriate grading permissions

## Step 1: List Assignments

First, find the assignment you want to grade:

```bash
canvas assignments list --course-id 123
```

Output:
```
ID     NAME                    DUE DATE              POINTS
456    Midterm Exam            2024-10-15 23:59      100
789    Final Project           2024-12-01 23:59      200
```

## Step 2: Export Submissions

Export all submissions for an assignment:

```bash
canvas submissions list --course-id 123 --assignment-id 456 -o csv > submissions.csv
```

The CSV file will contain:
```csv
user_id,user_name,submission_type,submitted_at,score,grade
1001,John Doe,online_upload,2024-10-14T10:30:00Z,,
1002,Jane Smith,online_upload,2024-10-15T08:45:00Z,,
```

## Step 3: Grade Offline

Open the CSV in your spreadsheet application and add scores:

```csv
user_id,user_name,submission_type,submitted_at,score,grade,comment
1001,John Doe,online_upload,2024-10-14T10:30:00Z,85,B,Good work!
1002,Jane Smith,online_upload,2024-10-15T08:45:00Z,92,A-,Excellent analysis
```

## Step 4: Import Grades

Create a grades file with the required columns:

```csv
user_id,score,comment
1001,85,Good work!
1002,92,Excellent analysis
```

Import the grades:

```bash
canvas submissions grade-batch \
  --course-id 123 \
  --assignment-id 456 \
  --file grades.csv
```

## Step 5: Verify Grades

Check that grades were applied:

```bash
canvas submissions list --course-id 123 --assignment-id 456
```

## Advanced: Grading with Rubrics

If your assignment uses a rubric, you can include rubric scores:

```csv
user_id,score,rubric_assessment
1001,85,{"criterion_1":{"points":20},"criterion_2":{"points":15}}
```

## Tips

!!! tip "Batch Size"
    For large classes, grades are processed in batches automatically. The CLI handles rate limiting.

!!! warning "Grade Overwrites"
    Importing grades will overwrite existing grades. Export current grades first as a backup.

!!! tip "Dry Run"
    Test your CSV format by grading a single submission first:
    ```bash
    canvas submissions grade 1001 \
      --course-id 123 \
      --assignment-id 456 \
      --score 85 \
      --comment "Good work!"
    ```

## Complete Workflow Script

```bash
#!/bin/bash
COURSE_ID=123
ASSIGNMENT_ID=456

# Export submissions
canvas submissions list \
  --course-id $COURSE_ID \
  --assignment-id $ASSIGNMENT_ID \
  -o csv > submissions.csv

echo "Edit submissions.csv and add grades, then press Enter"
read

# Import grades
canvas submissions grade-batch \
  --course-id $COURSE_ID \
  --assignment-id $ASSIGNMENT_ID \
  --file grades.csv

echo "Grades imported successfully!"
```
