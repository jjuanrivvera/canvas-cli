# Configuration

Canvas CLI can be configured to work with multiple Canvas instances and customized for your workflow.

## Configuration File

Canvas CLI stores configuration in `~/.canvas-cli/config.yaml`.

```yaml
default_instance: production
instances:
  production:
    name: production
    url: https://canvas.example.com
    # OAuth credentials (optional - for auto-refresh)
    client_id: "your-client-id"
    client_secret: "your-client-secret"
  sandbox:
    name: sandbox
    url: https://canvas-sandbox.example.com
    # API token (alternative to OAuth)
    token: "7~your-api-token-here"
```

!!! note "Authentication Methods"
    - **API Token**: Stored directly in config file (set via `canvas auth token set`)
    - **OAuth tokens**: Stored securely in your system keychain (set via `canvas auth login`)

    You can mix both methods - use OAuth for production and API tokens for testing/sandbox environments.

## Managing Instances

### Add an Instance

```bash
# Add a new instance
canvas config add production --url https://canvas.example.com

# Add with description
canvas config add staging --url https://canvas-staging.example.com --description "Staging environment"
```

After adding, authenticate with OAuth:
```bash
canvas auth login --instance production
```

Or set an API token:
```bash
canvas auth token set production --token 7~your-token-here
```

### List Instances

```bash
canvas config list
```

### Switch Default Instance

```bash
canvas config use sandbox
```

### Remove an Instance

```bash
canvas config remove sandbox
```

### Show Current Configuration

```bash
canvas config show
```

## Environment Variables

Canvas CLI supports environment variables for configuration (useful for CI/CD):

| Variable | Description |
|----------|-------------|
| `CANVAS_URL` | Canvas instance URL |
| `CANVAS_TOKEN` | API access token |
| `CANVAS_REQUESTS_PER_SEC` | Rate limit (default: 5.0) |

Example:

```bash
export CANVAS_URL=https://canvas.example.com
export CANVAS_TOKEN=your-api-token
canvas courses list
```

!!! tip "Priority"
    Environment variables take precedence over the config file.

## Command-Line Overrides

You can override configuration with command-line flags:

```bash
# Override instance
canvas courses list --instance https://other-canvas.example.com

# Override output format
canvas courses list --output json

# Disable caching
canvas courses list --no-cache
```

## Multiple Instances

Canvas CLI supports working with multiple Canvas instances. This is useful for:

- Development vs. production environments
- Multiple institutions
- Testing and staging

### Switching Instances

```bash
# Use a specific instance for one command
canvas courses list --instance sandbox

# Switch default instance
canvas config use sandbox
```

### Syncing Between Instances

```bash
# Sync course 123 from production to course 456 on sandbox
canvas sync course production 123 sandbox 456
```

See the [Course Sync Tutorial](../tutorials/course-sync.md) for more details.
