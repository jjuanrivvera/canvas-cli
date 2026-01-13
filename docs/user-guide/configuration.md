# Configuration

Canvas CLI can be configured to work with multiple Canvas instances and customized for your workflow.

## Configuration File

Canvas CLI stores configuration in `~/.canvas-cli/config.yaml`.

```yaml
default_instance: production
instances:
  production:
    url: https://canvas.example.com
    token: your-api-token
  sandbox:
    url: https://canvas-sandbox.example.com
    token: sandbox-token
```

## Managing Instances

### Add an Instance

```bash
# Add with API token
canvas config add production https://canvas.example.com --token YOUR_TOKEN

# Add without token (will prompt for OAuth)
canvas config add production https://canvas.example.com
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

Canvas CLI supports environment variables for configuration:

| Variable | Description |
|----------|-------------|
| `CANVAS_INSTANCE` | Canvas instance URL |
| `CANVAS_TOKEN` | API token |
| `CANVAS_OUTPUT` | Default output format |
| `CANVAS_NO_CACHE` | Disable caching (true/false) |

Example:

```bash
export CANVAS_INSTANCE=https://canvas.example.com
export CANVAS_TOKEN=your-api-token
canvas courses list
```

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
# Sync a course from production to sandbox
canvas sync course 123 --from production --to sandbox
```
