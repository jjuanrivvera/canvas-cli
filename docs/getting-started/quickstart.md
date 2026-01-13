# Quick Start

Get up and running with Canvas CLI in 5 minutes.

## Prerequisites

- Canvas CLI installed ([Installation Guide](installation.md))
- Canvas LMS account with API access

## Step 1: Authenticate

First, authenticate with your Canvas instance:

```bash
canvas auth login
```

This will open your browser for OAuth authentication. After authorizing, you'll be redirected back to the CLI.

!!! tip "API Token Alternative"
    If OAuth isn't available, you can use an API token:
    ```bash
    canvas config add --name mycanvas --url https://canvas.example.com --token YOUR_API_TOKEN
    ```

## Step 2: Verify Authentication

Check that you're authenticated:

```bash
canvas auth status
```

You should see your user information and the Canvas instance URL.

## Step 3: List Your Courses

```bash
canvas courses list
```

This shows all courses you have access to. Note the course IDs for use in other commands.

## Step 4: Explore Course Content

```bash
# List assignments
canvas assignments list --course-id 123

# List users in a course
canvas users list --course-id 123

# List modules
canvas modules list --course-id 123
```

## Step 5: Try Different Output Formats

Canvas CLI supports multiple output formats:

=== "Table (default)"

    ```bash
    canvas courses list
    ```

=== "JSON"

    ```bash
    canvas courses list --output json
    ```

=== "YAML"

    ```bash
    canvas courses list --output yaml
    ```

=== "CSV"

    ```bash
    canvas courses list --output csv
    ```

## Common Commands

| Task | Command |
|------|---------|
| List courses | `canvas courses list` |
| Get course details | `canvas courses get 123` |
| List assignments | `canvas assignments list --course-id 123` |
| Grade submission | `canvas submissions grade 456 --course-id 123 --score 95` |
| List users | `canvas users list --course-id 123` |
| Upload file | `canvas files upload ./file.pdf --course-id 123` |

## Next Steps

- [User Guide](../user-guide/index.md) - Learn more about configuration and features
- [Command Reference](../commands/index.md) - Complete command documentation
- [Tutorials](../tutorials/index.md) - Step-by-step guides for common tasks
