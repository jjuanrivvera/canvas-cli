# Command Reference

Complete reference for all Canvas CLI commands.

## Global Flags

These flags work with all commands:

```bash
--instance string      Canvas instance to use (default: active instance)
--format string        Output format: json, yaml, table, csv (default: table)
--output string        Write output to file instead of stdout
--debug               Enable debug logging
--no-cache            Disable caching for this request
--help                Show help for command
```

## Authentication Commands

### `canvas auth login`

Authenticate with a Canvas instance.

```bash
canvas auth login [flags]

Flags:
  --instance string        Canvas instance URL (required)
  --client-id string       OAuth client ID (optional)
  --client-secret string   OAuth client secret (optional)
  --oob                    Use out-of-band OAuth flow
  --port int               Local server port (default: 8080)
```

**Examples:**

```bash
# Interactive login
canvas auth login --instance https://canvas.instructure.com

# Out-of-band flow (for SSH sessions)
canvas auth login --instance https://canvas.instructure.com --oob

# Custom OAuth credentials
canvas auth login --instance https://canvas.example.com \
  --client-id abc123 \
  --client-secret xyz789
```

### `canvas auth logout`

Remove stored credentials.

```bash
canvas auth logout [flags]

Flags:
  --instance string   Logout from specific instance
  --all              Logout from all instances
```

### `canvas auth status`

Check authentication status.

```bash
canvas auth status [flags]

Flags:
  --instance string   Check specific instance (default: all)
```

## Configuration Commands

### `canvas config list`

List all configured instances.

```bash
canvas config list
```

### `canvas config add`

Add a new Canvas instance.

```bash
canvas config add NAME [flags]

Flags:
  --url string   Canvas instance URL (required)

Example:
  canvas config add production --url https://canvas.instructure.com
```

### `canvas config use`

Switch active instance.

```bash
canvas config use NAME

Example:
  canvas config use staging
```

### `canvas config remove`

Remove a configured instance.

```bash
canvas config remove NAME
```

### `canvas config show`

Show current configuration.

```bash
canvas config show
```

## Course Commands

### `canvas courses list`

List accessible courses.

```bash
canvas courses list [flags]

Flags:
  --state string        Filter by state: available, completed, deleted (default: all)
  --enrollment-type string   Filter by enrollment type: teacher, student, ta, observer
  --include strings     Additional data: term, teachers, students, total_students
  --per-page int        Results per page (default: 10)

Examples:
  # List all active courses
  canvas courses list --state available

  # List courses where you're a teacher
  canvas courses list --enrollment-type teacher

  # Include enrollment counts
  canvas courses list --include total_students

  # Get more results
  canvas courses list --per-page 50
```

### `canvas courses get`

Get details for a specific course.

```bash
canvas courses get COURSE_ID [flags]

Flags:
  --include strings   Additional data: term, teachers, students, syllabus_body

Example:
  canvas courses get 12345 --include teachers,students
```

### `canvas courses users`

List users enrolled in a course.

```bash
canvas courses users COURSE_ID [flags]

Flags:
  --enrollment-type strings   Filter by type: student, teacher, ta, observer
  --per-page int             Results per page (default: 10)

Example:
  canvas courses users 12345 --enrollment-type student
```

## Assignment Commands

### `canvas assignments list`

List assignments for a course.

```bash
canvas assignments list [flags]

Flags:
  --course int          Course ID (required)
  --include strings     Additional data: submission, overrides
  --per-page int        Results per page (default: 10)

Example:
  canvas assignments list --course 12345 --include submission
```

### `canvas assignments get`

Get details for a specific assignment.

```bash
canvas assignments get [flags]

Flags:
  --course int         Course ID (required)
  --assignment int     Assignment ID (required)
  --include strings    Additional data: submission, overrides

Example:
  canvas assignments get --course 12345 --assignment 67890
```

### `canvas assignments submissions`

List submissions for an assignment.

```bash
canvas assignments submissions [flags]

Flags:
  --course int         Course ID (required)
  --assignment int     Assignment ID (required)
  --include strings    Additional data: user, submission_comments
  --per-page int       Results per page (default: 10)

Example:
  canvas assignments submissions --course 12345 --assignment 67890 --include user
```

### `canvas assignments grade`

Grade a submission.

```bash
canvas assignments grade [flags]

Flags:
  --course int         Course ID (required)
  --assignment int     Assignment ID (required)
  --user int           User ID (required)
  --score float        Numeric score
  --comment string     Grading comment
  --csv string         CSV file for bulk grading

Examples:
  # Grade single submission
  canvas assignments grade --course 12345 --assignment 67890 --user 11111 --score 95

  # Grade with comment
  canvas assignments grade --course 12345 --assignment 67890 --user 11111 \
    --score 95 --comment "Great work!"

  # Bulk grade from CSV
  canvas assignments grade --course 12345 --assignment 67890 --csv grades.csv
```

## User Commands

### `canvas users me`

Get current user information.

```bash
canvas users me
```

### `canvas users get`

Get user details.

```bash
canvas users get USER_ID

Example:
  canvas users get 12345
```

### `canvas users create`

Create a new user.

```bash
canvas users create [flags]

Flags:
  --name string          Full name (required)
  --email string         Email address (required)
  --login string         Login ID (required)
  --sis-user-id string   SIS user ID
  --send-confirmation    Send confirmation email

Example:
  canvas users create --name "John Doe" --email john@example.com --login jdoe
```

### `canvas users update`

Update user information.

```bash
canvas users update USER_ID [flags]

Flags:
  --name string    New name
  --email string   New email

Example:
  canvas users update 12345 --email newemail@example.com
```

## Enrollment Commands

### `canvas enrollments list`

List enrollments for a course.

```bash
canvas enrollments list COURSE_ID [flags]

Flags:
  --type strings     Filter by type: student, teacher, ta, observer
  --state strings    Filter by state: active, invited, completed
  --per-page int     Results per page (default: 10)

Example:
  canvas enrollments list 12345 --type student --state active
```

### `canvas enrollments create`

Create an enrollment.

```bash
canvas enrollments create [flags]

Flags:
  --course int       Course ID (required)
  --user int         User ID (required)
  --type string      Enrollment type: student, teacher, ta, observer (required)
  --role-id int      Role ID
  --state string     Enrollment state: active, invited, completed
  --notify           Send enrollment notification

Example:
  canvas enrollments create --course 12345 --user 67890 --type student --notify
```

## File Commands

### `canvas files upload`

Upload a file.

```bash
canvas files upload FILE_PATH [flags]

Flags:
  --course int       Upload to course (specify course ID)
  --user             Upload to user files
  --folder string    Destination folder path
  --name string      Custom filename
  --on-duplicate string   Action on duplicate: overwrite, rename (default: rename)

Examples:
  # Upload to user files
  canvas files upload document.pdf --user

  # Upload to course
  canvas files upload syllabus.pdf --course 12345 --folder "Course Documents"

  # Overwrite existing file
  canvas files upload document.pdf --course 12345 --on-duplicate overwrite
```

### `canvas files download`

Download a file.

```bash
canvas files download FILE_ID [flags]

Flags:
  --output string   Output path (default: original filename)

Example:
  canvas files download 98765 --output downloaded_file.pdf
```

## Sync Commands

### `canvas sync copy`

Copy resources between Canvas instances.

```bash
canvas sync copy [flags]

Flags:
  --from string         Source instance
  --to string           Target instance
  --course int          Course ID to copy
  --assignments         Copy assignments
  --files               Copy files
  --resolve string      Conflict resolution: skip, overwrite, prompt (default: prompt)

Example:
  canvas sync copy --from production --to staging --course 12345 --assignments
```

## Utility Commands

### `canvas shell`

Start interactive REPL mode.

```bash
canvas shell

# In REPL mode, commands work without 'canvas' prefix:
canvas> courses list
canvas> assignments list --course 12345
canvas> exit
```

### `canvas doctor`

Run diagnostics.

```bash
canvas doctor

# Checks:
# - Internet connectivity
# - Canvas API reachability
# - Authentication status
# - Configuration validity
# - Keychain access
```

### `canvas cache clear`

Clear the cache.

```bash
canvas cache clear

Flags:
  --all      Clear all cached data
  --expired  Clear only expired entries
```

### `canvas cache stats`

Show cache statistics.

```bash
canvas cache stats
```

### `canvas completion`

Generate shell completion script.

```bash
canvas completion SHELL

Supported shells: bash, zsh, fish, powershell

Examples:
  canvas completion bash > /etc/bash_completion.d/canvas
  canvas completion zsh > "${fpath[1]}/_canvas"
  canvas completion fish > ~/.config/fish/completions/canvas.fish
```

### `canvas version`

Show version information.

```bash
canvas version
```

## Output Formats

All commands support multiple output formats:

### JSON (default for scripts)

```bash
canvas courses list --format json
```

### YAML

```bash
canvas courses list --format yaml
```

### Table (default for interactive)

```bash
canvas courses list --format table
```

### CSV

```bash
canvas courses list --format csv > courses.csv
```

## Exit Codes

Canvas CLI uses standard exit codes:

- `0` - Success
- `1` - General error
- `2` - Authentication error
- `3` - Not found (404)
- `4` - Permission denied (403)
- `5` - Rate limit exceeded (429)

## Environment Variables

- `CANVAS_URL` - Default Canvas instance URL
- `CANVAS_TOKEN` - Access token for authentication
- `CANVAS_CONFIG_DIR` - Configuration directory (default: `~/.canvas-cli`)
- `CANVAS_CACHE_DIR` - Cache directory
- `CANVAS_DEBUG` - Enable debug logging (set to `1` or `true`)

## Next Steps

- [Examples](EXAMPLES.md) - See practical usage examples
- [Authentication](AUTHENTICATION.md) - Learn about OAuth setup
