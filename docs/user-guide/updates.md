# Auto-Updates

Canvas CLI includes built-in support for checking and installing updates directly from GitHub releases.

## Checking for Updates

To check if a new version is available:

```bash
canvas update check
```

This will display whether an update is available and show version information.

### Output Formats

You can get the update check results in different formats:

```bash
# JSON format
canvas update check --output json

# YAML format
canvas update check --output yaml
```

### Forcing a Fresh Check

Update checks are cached for 6 hours to avoid rate limits. To force a fresh check:

```bash
canvas update check --force
```

## Installing Updates

To install the latest version:

```bash
canvas update install
```

The command will:

1. Check for available updates
2. Ask for confirmation (unless `--yes` is used)
3. Download the latest release
4. Replace the current binary
5. Prompt you to restart

### Skip Confirmation

To install without prompting:

```bash
canvas update install --yes
```

## Automatic Update Checks

Canvas CLI can automatically check for updates in the background (non-intrusive notifications only).

### Enable Automatic Checks

```bash
# Enable with default 24-hour interval
canvas update enable

# Enable with custom interval (in hours)
canvas update enable --interval 12
```

### Disable Automatic Checks

```bash
canvas update disable
```

### Configuration

Automatic update settings are stored in your config file (`~/.canvas-cli/config.yaml`):

```yaml
settings:
  auto_update_check: true
  update_check_interval_hours: 24
```

## Installation Method Compatibility

The `canvas update install` command only works for binaries installed directly from GitHub releases.

### Homebrew Users

If you installed via Homebrew, update using:

```bash
brew upgrade canvas-cli
```

### Go Install Users

If you installed via `go install`, update using:

```bash
go install github.com/jjuanrivvera/canvas-cli/cmd/canvas@latest
```

### Package Manager Detection

Canvas CLI automatically detects common installation methods and will inform you if updates should be performed through your package manager instead:

- Homebrew (`/usr/local/Cellar/`)
- Snap (`/snap/`)
- Flatpak (`/flatpak/`)

## Examples

### Check for Updates Weekly

```bash
# Check manually
canvas update check

# Enable automatic checks with weekly interval
canvas update enable --interval 168
```

### Automated Update in CI/CD

```bash
# Check and install if available (non-interactive)
canvas update check --output json
canvas update install --yes
```

### Get Update Information as JSON

```bash
canvas update check --output json | jq '.update_available'
```

This outputs:

```json
{
  "update_available": true,
  "current_version": "v1.5.0",
  "latest_version": "v1.6.0",
  "release_info": {
    "version": "v1.6.0",
    "url": "https://github.com/jjuanrivvera/canvas-cli/releases/tag/v1.6.0",
    "release_date": "2024-01-15T10:30:00Z",
    "notes": "Release notes...",
    "asset_url": "https://github.com/jjuanrivvera/canvas-cli/releases/download/v1.6.0/canvas-cli_1.6.0_darwin_amd64.tar.gz",
    "asset_name": "canvas-cli_1.6.0_darwin_amd64.tar.gz"
  },
  "checked_at": "2024-01-15T15:45:00Z"
}
```

## Troubleshooting

### Update Check Fails

If update checks fail, verify:

1. **Internet connectivity**: Ensure you can reach `github.com`
2. **Rate limits**: GitHub API has rate limits; try again later or use `--force` less frequently
3. **Proxy settings**: Configure proxy environment variables if needed

### Permission Denied During Install

If you get a permission error during installation:

1. Ensure you have write permissions to the binary location
2. On Unix systems, you may need to run with appropriate permissions
3. Consider reinstalling to a user-writable location

### Development Version Detection

If you're running a development build (version "dev"), update checks will not find new versions. Install a released version first:

```bash
# Install latest release
go install github.com/jjuanrivvera/canvas-cli/cmd/canvas@latest

# Or download from releases
# https://github.com/jjuanrivvera/canvas-cli/releases
```

## Security

### Update Verification

Canvas CLI uses the official GitHub selfupdate library which:

- Verifies checksums of downloaded binaries
- Uses HTTPS for all communications
- Validates release authenticity through GitHub's API

### Update Source

Updates are only fetched from the official repository:

- **Owner**: `jjuanrivvera`
- **Repository**: `canvas-cli`
- **URL**: https://github.com/jjuanrivvera/canvas-cli

## Cache Management

Update check results are cached in `~/.canvas-cli/cache/update_check.json` for 6 hours.

To clear the update cache:

```bash
canvas cache clear
```

Or manually delete the cache file:

```bash
rm ~/.canvas-cli/cache/update_check.json
```
