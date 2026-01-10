# Authentication Guide

Canvas CLI uses OAuth 2.0 for secure authentication with Canvas LMS. This guide covers everything you need to set up and manage authentication.

## Quick Start

The simplest way to authenticate:

```bash
canvas auth login --instance https://canvas.instructure.com
```

This will:
1. Open your browser to Canvas
2. Ask you to authorize the application
3. Automatically save your credentials securely

## OAuth 2.0 Setup

### Step 1: Create a Developer Key (Optional)

By default, Canvas CLI uses embedded OAuth credentials. For enhanced security or custom integrations, you can create your own:

1. Log into Canvas as an administrator
2. Navigate to **Admin > Developer Keys**
3. Click **+ Developer Key > + API Key**
4. Fill in the details:
   - **Key Name**: Canvas CLI
   - **Owner Email**: your-email@example.com
   - **Redirect URIs**: `http://localhost:8080/callback` (or your preferred port)
   - **Scopes**: Select all required scopes
5. Save and note your **Client ID** and **Client Secret**

### Step 2: Authenticate

#### Using Embedded Credentials (Easiest)

```bash
canvas auth login --instance https://canvas.instructure.com
```

#### Using Your Own Credentials

```bash
canvas auth login \
  --instance https://canvas.instructure.com \
  --client-id YOUR_CLIENT_ID \
  --client-secret YOUR_CLIENT_SECRET
```

### Step 3: Verify Authentication

```bash
# Check authentication status
canvas auth status

# Test with a simple API call
canvas courses list
```

## Authentication Methods

### Local Browser Flow (Default)

The CLI starts a local web server and opens your browser:

```bash
canvas auth login --instance https://canvas.instructure.com
```

This is the most user-friendly method and works on most systems.

### Out-of-Band (OOB) Flow

For systems without a browser or remote SSH sessions:

```bash
canvas auth login --instance https://canvas.instructure.com --oob
```

This will:
1. Display an authorization URL
2. Ask you to visit it manually
3. Prompt you to paste the authorization code

### Environment Variables (CI/CD)

For automated workflows and CI/CD pipelines, Canvas CLI supports authentication via environment variables:

```bash
# Required environment variables
export CANVAS_URL="https://canvas.instructure.com"
export CANVAS_TOKEN="your-access-token"

# Optional: Control rate limiting (default: 5.0)
export CANVAS_REQUESTS_PER_SEC="10.0"

# Commands will use these automatically
canvas courses list
```

**How to get your access token:**

1. Log into Canvas
2. Go to **Account > Settings > Approved Integrations**
3. Click **+ New Access Token**
4. Copy the token (it will only be shown once)

**Priority Order:**

Canvas CLI checks for credentials in this order:
1. **Environment variables** (CANVAS_URL + CANVAS_TOKEN) - highest priority
2. **Config file** (~/.canvas-cli/config.yaml)
3. **OAuth tokens** (stored in system keychain)

**CI/CD Example (GitHub Actions):**

```yaml
name: Canvas CI
on: [push]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run Canvas CLI
        env:
          CANVAS_URL: ${{ secrets.CANVAS_URL }}
          CANVAS_TOKEN: ${{ secrets.CANVAS_TOKEN }}
        run: |
          canvas courses list
          canvas assignments list --course-id 123
```

## Multi-Instance Management

Canvas CLI supports multiple Canvas installations:

### Add Multiple Instances

```bash
# Add production instance
canvas config add production \
  --url https://canvas.instructure.com

# Add staging instance
canvas config add staging \
  --url https://staging.canvas.instructure.com

# Add self-hosted instance
canvas config add onprem \
  --url https://canvas.company.com
```

### Switch Between Instances

```bash
# List all instances
canvas config list

# Switch to a different instance
canvas config use staging

# Use a specific instance for one command
canvas courses list --instance production
```

### Authenticate Each Instance

```bash
# Authenticate with each instance separately
canvas auth login --instance production
canvas auth login --instance staging
canvas auth login --instance onprem
```

## Credential Storage

Canvas CLI stores credentials securely using your system's native keychain:

- **macOS**: Keychain
- **Linux**: Secret Service API (gnome-keyring, kwallet)
- **Windows**: Windows Credential Manager

### Fallback Storage

If keychain access fails, credentials are stored in encrypted files:

- **Location**: `~/.canvas-cli/credentials.enc`
- **Encryption**: AES-256-GCM with machine-specific key
- **Key Derivation**: Machine ID + username

### Manual Credential Management

```bash
# View stored credentials (won't show tokens)
canvas auth status

# Logout (removes credentials)
canvas auth logout

# Logout from specific instance
canvas auth logout --instance staging

# Logout from all instances
canvas auth logout --all
```

## Token Management

### Token Refresh

Access tokens expire after 1 hour. Canvas CLI automatically refreshes them using refresh tokens.

### Token Permissions (Scopes)

Canvas CLI requests the following OAuth scopes:

- `url:GET|/api/v1/courses` - Read courses
- `url:POST|/api/v1/courses` - Create courses
- `url:PUT|/api/v1/courses/:id` - Update courses
- `url:DELETE|/api/v1/courses/:id` - Delete courses
- Similar patterns for assignments, submissions, users, files

### Manual Token Usage

If you have an access token:

```bash
# Use directly
canvas courses list --token YOUR_ACCESS_TOKEN

# Or set environment variable
export CANVAS_TOKEN="YOUR_ACCESS_TOKEN"
canvas courses list
```

## Security Best Practices

### 1. Use System Keychain

Always allow Canvas CLI to use your system keychain when prompted.

### 2. Rotate Tokens Regularly

```bash
# Logout and login again to get fresh tokens
canvas auth logout
canvas auth login
```

### 3. Limit Scope

When creating developer keys, only grant necessary scopes.

### 4. Use Environment Variables for CI/CD

For automated systems, use environment variables instead of storing credentials in code:

```yaml
# Example GitHub Actions workflow
env:
  CANVAS_URL: ${{ secrets.CANVAS_URL }}
  CANVAS_TOKEN: ${{ secrets.CANVAS_TOKEN }}
```

### 5. Keep Credentials Private

Never commit tokens or credentials to version control:

```bash
# Add to .gitignore
echo ".canvas-cli/" >> .gitignore
```

## Troubleshooting

### "Failed to open browser"

Use OOB flow:
```bash
canvas auth login --instance https://canvas.instructure.com --oob
```

### "Failed to access keychain"

Grant permission in system settings or use encrypted file storage (automatic fallback).

### "Token expired"

Tokens refresh automatically. If refresh fails:
```bash
canvas auth logout
canvas auth login
```

### "Unauthorized" or "403 Forbidden"

Check your OAuth scopes or Canvas permissions:
```bash
canvas auth status
```

### "CORS errors"

Ensure your Canvas instance allows OAuth from localhost:
- Check Developer Key redirect URIs
- Verify Canvas OAuth settings

## Advanced Configuration

### Custom Redirect Port

```bash
canvas auth login \
  --instance https://canvas.instructure.com \
  --port 9000
```

### Specify Redirect URI

```bash
canvas auth login \
  --instance https://canvas.instructure.com \
  --redirect-uri http://localhost:8080/callback
```

### Debug Authentication

```bash
canvas auth login --instance https://canvas.instructure.com --debug
```

## Next Steps

After authentication:
- [Command Reference](COMMANDS.md) - Learn available commands
- [Examples](EXAMPLES.md) - See common use cases
- Test your setup: `canvas courses list`
