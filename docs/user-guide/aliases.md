# Command Aliases

Create shortcuts for frequently used commands with the alias system.

## Overview

Aliases let you define short names for long or complex commands. Instead of typing:

```bash
canvas assignments list --course-id 12345 --output json
```

You can create an alias and use:

```bash
canvas hw
```

## Creating Aliases

Use `canvas alias set` to create an alias:

```bash
canvas alias set <name> "<command>"
```

### Examples

```bash
# Simple alias for listing assignments
canvas alias set hw "assignments list --course-id 12345"

# Alias for grading submissions
canvas alias set grade "submissions grade --course-id 12345 --assignment-id 67890"

# Alias for checking course users
canvas alias set students "courses users --course-id 12345 --enrollment-type student"
```

## Using Aliases

Once created, use your alias as if it were a built-in command:

```bash
# Instead of: canvas assignments list --course-id 12345
canvas hw

# Aliases can also accept additional arguments
canvas hw --output json
```

Additional arguments are appended to the alias expansion.

## Managing Aliases

### List All Aliases

```bash
canvas alias list
```

Output:
```
Aliases:
  hw        → assignments list --course-id 12345
  grade     → submissions grade --course-id 12345 --assignment-id 67890
  students  → courses users --course-id 12345 --enrollment-type student
```

### Delete an Alias

```bash
canvas alias delete hw
```

## Alias Rules

1. **No conflicts**: Alias names cannot match built-in command names
2. **Simple expansion**: Aliases expand once (no recursive aliases)
3. **Quoted arguments**: Use quotes for multi-word values in the expansion

### Example with Quotes

```bash
# Alias with a search term containing spaces
canvas alias set find-john 'users search --query "John Smith"'
```

## Storage

Aliases are stored in your configuration file (`~/.canvas-cli/config.yaml`):

```yaml
aliases:
  hw: "assignments list --course-id 12345"
  grade: "submissions grade --course-id 12345 --assignment-id 67890"
```

## Tips

### Combine with Context

Aliases work great with [context management](context.md). Set your course context once, then use aliases that don't need `--course-id`:

```bash
# Set context
canvas context set course 12345

# Create alias without course-id (uses context)
canvas alias set hw "assignments list"

# Use it
canvas hw
```

### Alias for Different Output Formats

```bash
canvas alias set hw-json "assignments list --output json"
canvas alias set hw-csv "assignments list --output csv"
```

### Debugging Aliases

Use `--dry-run` to see what command an alias expands to:

```bash
canvas --dry-run hw
```
