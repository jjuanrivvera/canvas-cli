# Context Management

Set default values for common flags like `--course-id` so you don't have to type them repeatedly.

## Overview

When working with a specific course, you typically run many commands with the same `--course-id`. Context management lets you set this once and have it apply automatically:

```bash
# Without context - repetitive
canvas assignments list --course-id 12345
canvas submissions list --course-id 12345 --assignment-id 67890
canvas modules list --course-id 12345

# With context - clean and simple
canvas context set course 12345
canvas assignments list
canvas submissions list --assignment-id 67890
canvas modules list
```

## Setting Context

Use `canvas context set` to set default values:

```bash
canvas context set <type> <id>
```

### Available Context Types

| Type | Flag it replaces | Example |
|------|------------------|---------|
| `course` | `--course-id` | `canvas context set course 12345` |
| `assignment` | `--assignment-id` | `canvas context set assignment 67890` |
| `user` | `--user-id` | `canvas context set user 111` |
| `account` | `--account-id` | `canvas context set account 1` |

### Setting Multiple Values

```bash
canvas context set course 12345
canvas context set assignment 67890
```

Now commands automatically use both:
```bash
# Uses course 12345 and assignment 67890
canvas submissions list
```

## Viewing Context

See your current context settings:

```bash
canvas context show
```

Output:
```
Current context:
  course_id:     12345
  assignment_id: 67890
```

## Clearing Context

### Clear All Context

```bash
canvas context clear
```

### Clear Specific Value

```bash
canvas context clear course
canvas context clear assignment
```

## How Context Works

1. When you run a command, Canvas CLI checks if required flags are provided
2. If a flag is missing, it checks for a context value
3. Context values are used as defaults when flags aren't specified
4. Explicit flags always override context values

### Example Flow

```bash
# Set context
canvas context set course 12345

# These are equivalent:
canvas assignments list                    # Uses context (course 12345)
canvas assignments list --course-id 12345  # Explicit flag

# Override context with a different value
canvas assignments list --course-id 99999  # Uses 99999, ignores context
```

## Storage

Context is stored in your configuration file (`~/.canvas-cli/config.yaml`):

```yaml
context:
  course_id: 12345
  assignment_id: 67890
```

## Workflow Examples

### Grading Workflow

```bash
# Set up context for grading session
canvas context set course 12345
canvas context set assignment 67890

# Now grade submissions efficiently
canvas submissions list
canvas submissions grade 111 --grade 95 --comment "Great work!"
canvas submissions grade 222 --grade 88

# Done grading, clear context
canvas context clear
```

### Course Management Workflow

```bash
# Working on a specific course
canvas context set course 12345

# Run multiple commands
canvas assignments list
canvas modules list
canvas users list --enrollment-type student
canvas announcements list

# Switch to another course
canvas context set course 54321
```

### Combined with Aliases

Context and [aliases](aliases.md) work great together:

```bash
# Set course context
canvas context set course 12345

# Create aliases that use context
canvas alias set hw "assignments list"
canvas alias set grades "submissions list"

# Use them - context provides the course ID
canvas hw
canvas grades --assignment-id 67890
```

## Tips

- Set context at the start of a work session
- Clear context when switching tasks to avoid confusion
- Use `canvas context show` to verify your current context
- Explicit flags always take precedence over context
