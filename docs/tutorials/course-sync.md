# Course Sync Tutorial

Learn how to synchronize courses between Canvas instances.

## Overview

Canvas CLI's sync feature allows you to:

- Copy course content between instances (production ↔ sandbox)
- Migrate courses between semesters
- Maintain course templates

## Prerequisites

- Canvas CLI installed
- Access to both source and destination Canvas instances
- Appropriate permissions on both instances

## Step 1: Configure Multiple Instances

First, set up both Canvas instances:

```bash
# Add production instance
canvas config add production --url https://canvas.example.com --token YOUR_PROD_TOKEN

# Add sandbox instance
canvas config add sandbox --url https://canvas-sandbox.example.com --token YOUR_SANDBOX_TOKEN
```

Verify your instances:

```bash
canvas config list
```

Output:
```
NAME         URL                                  DEFAULT
production   https://canvas.example.com           *
sandbox      https://canvas-sandbox.example.com
```

## Step 2: Identify Source Course

Find the course you want to sync:

```bash
canvas courses list --instance production
```

Note the course ID you want to sync.

## Step 3: Create Destination Course

If needed, create a new course on the destination:

```bash
canvas courses create \
  --instance sandbox \
  --name "CS101 - Test Copy" \
  --code "CS101-TEST"
```

Or identify an existing course:

```bash
canvas courses list --instance sandbox
```

## Step 4: Sync Course Content

Sync the course from production to sandbox:

```bash
canvas sync course 123 \
  --from production \
  --to sandbox \
  --destination-course 456
```

### Sync Options

| Option | Description |
|--------|-------------|
| `--modules` | Sync modules and items |
| `--assignments` | Sync assignments |
| `--pages` | Sync pages |
| `--files` | Sync files |
| `--all` | Sync all content types |

Example with specific content:

```bash
canvas sync course 123 \
  --from production \
  --to sandbox \
  --destination-course 456 \
  --modules \
  --assignments
```

## Step 5: Verify Sync

Check the destination course:

```bash
canvas modules list --course-id 456 --instance sandbox
canvas assignments list --course-id 456 --instance sandbox
```

## Use Cases

### Development Testing

Sync a production course to sandbox for testing changes:

```bash
canvas sync course 123 --from production --to sandbox --all
```

### Semester Rollover

Copy a course template to create a new semester's course:

```bash
canvas sync course 100 \
  --from production \
  --to production \
  --destination-course 200 \
  --modules \
  --pages
```

### Multi-Institution Deployment

Sync course content between institutions:

```bash
# Configure second institution
canvas config add partner --url https://partner.instructure.com --token TOKEN

# Sync course
canvas sync course 123 --from production --to partner --all
```

## Tips

!!! warning "Enrollment Data"
    Sync does **not** copy enrollment or grade data. Only course structure and content are synced.

!!! tip "Incremental Sync"
    For large courses, sync specific content types incrementally rather than using `--all`.

!!! tip "Verify Before Production"
    Always sync to a sandbox first to verify the results before syncing to production.

## Troubleshooting

### Permission Denied

Ensure you have the required permissions:
- Source: Read access to course content
- Destination: Write access to create/update content

### Content Not Syncing

Some content types have dependencies:
- Module items require modules to exist first
- Assignment groups should sync before assignments

### Rate Limiting

For large syncs, Canvas CLI automatically handles rate limiting. If you see rate limit errors, the CLI will retry automatically.
