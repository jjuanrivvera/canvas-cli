# Environment Variables

Canvas CLI can be configured using environment variables.

## Available Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `CANVAS_URL` | Canvas instance URL (env auth mode) | From config |
| `CANVAS_TOKEN` | API access token | From config/keyring |
| `CANVAS_REQUESTS_PER_SEC` | Request rate limit for env auth mode | `5.0` |
| `CANVAS_OUTPUT` | Default output format | `table` |
| `CANVAS_NO_CACHE` | Disable response caching | `false` |

## Usage

### Temporary Override

Set for a single command:

```bash
CANVAS_OUTPUT=json canvas courses list
```

### Session Override

Set for the current shell session:

```bash
export CANVAS_URL=https://canvas.example.com
export CANVAS_TOKEN=your-api-token
canvas courses list
```

### Permanent Configuration

Add to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.):

```bash
# Canvas CLI configuration
export CANVAS_URL=https://canvas.example.com
export CANVAS_TOKEN=your-api-token
export CANVAS_OUTPUT=json
```

## Variable Details

### CANVAS_URL

The Canvas LMS instance URL.

```bash
export CANVAS_URL=https://canvas.instructure.com
```

!!! note "Priority"
    When both `CANVAS_URL` and `CANVAS_TOKEN` are set, CLI commands use env authentication directly.
    In that mode, `--instance` and `default_instance` are not used.

### CANVAS_TOKEN

Canvas API access token. Generate from Canvas Account Settings.

```bash
export CANVAS_TOKEN=7~AbCdEfGhIjKlMnOpQrStUvWxYz123456789
```

### CANVAS_REQUESTS_PER_SEC

Request rate limit used when env auth mode is active.

```bash
export CANVAS_REQUESTS_PER_SEC=10.0
```

!!! warning "Security"
    Avoid setting tokens in shared environments. Consider using the config file or keyring instead.

### CANVAS_OUTPUT

Default output format for commands.

| Value | Description |
|-------|-------------|
| `table` | Human-readable table (default) |
| `json` | JSON format |
| `yaml` | YAML format |
| `csv` | CSV format |

```bash
export CANVAS_OUTPUT=json
```

### CANVAS_NO_CACHE

Disable API response caching.

```bash
export CANVAS_NO_CACHE=true
```

!!! tip "When to Disable"
    Disable caching when you need real-time data or are debugging.

## Precedence Order

Configuration is applied in this order (highest precedence first):

1. If both `CANVAS_URL` and `CANVAS_TOKEN` are set, env auth mode is used
2. Otherwise, command-line flags (`--instance`, `--output`, etc.)
3. Otherwise, configuration file (`~/.canvas-cli/config.yaml`)
4. Built-in defaults

!!! note "MCP clients"
    This same precedence applies to MCP tool calls, because MCP executes the same command logic.
    Each MCP client runs its own server process and therefore may have different environment variables.

## Example Configurations

### Development Environment

```bash
# Use sandbox instance
export CANVAS_URL=https://canvas-sandbox.example.com
export CANVAS_TOKEN=sandbox-token
export CANVAS_NO_CACHE=true
```

### CI/CD Pipeline

```yaml
# GitHub Actions example
env:
  CANVAS_URL: ${{ secrets.CANVAS_URL }}
  CANVAS_TOKEN: ${{ secrets.CANVAS_TOKEN }}
  CANVAS_OUTPUT: json
```

### Multi-Instance Setup

```bash
# Function to switch instances
canvas-prod() {
  export CANVAS_URL=https://canvas.example.com
  export CANVAS_TOKEN=$CANVAS_PROD_TOKEN
}

canvas-sandbox() {
  export CANVAS_URL=https://canvas-sandbox.example.com
  export CANVAS_TOKEN=$CANVAS_SANDBOX_TOKEN
}
```

## Troubleshooting

### Variable Not Working

1. Verify the variable is set:
   ```bash
   echo $CANVAS_URL
   ```

2. Check for typos in variable names

3. Ensure you've sourced your profile:
   ```bash
   source ~/.bashrc
   ```

### Token Security

If your token is exposed:

1. Immediately regenerate it in Canvas Account Settings
2. Update your configuration
3. Consider using the config file with restrictive permissions:
   ```bash
   chmod 600 ~/.canvas-cli/config.yaml
   ```
