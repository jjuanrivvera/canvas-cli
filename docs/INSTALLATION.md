# Installation Guide

This guide covers all the ways to install Canvas CLI on your system.

## System Requirements

- **Operating Systems**: macOS, Linux, Windows
- **Architecture**: amd64, arm64
- **Go Version** (if building from source): 1.21 or later

## Installation Methods

### Method 1: Homebrew (macOS/Linux) - Recommended

```bash
# Add the Canvas CLI tap
brew tap jjuanrivvera/canvas-cli

# Install Canvas CLI
brew install canvas-cli

# Verify installation
canvas version
```

#### Upgrading

```bash
brew upgrade canvas-cli
```

### Method 2: Using Go

If you have Go installed:

```bash
# Install latest version
go install github.com/jjuanrivvera/canvas-cli/cmd/canvas@latest

# Install specific version
go install github.com/jjuanrivvera/canvas-cli/cmd/canvas@v1.0.0

# Verify installation
canvas version
```

### Method 3: Download Binary (All Platforms)

1. Visit the [Releases page](https://github.com/jjuanrivvera/canvas-cli/releases)
2. Download the appropriate binary for your platform:
   - **macOS (Intel)**: `canvas_v1.0.0_Darwin_x86_64.tar.gz`
   - **macOS (Apple Silicon)**: `canvas_v1.0.0_Darwin_arm64.tar.gz`
   - **Linux (64-bit)**: `canvas_v1.0.0_Linux_x86_64.tar.gz`
   - **Windows (64-bit)**: `canvas_v1.0.0_Windows_x86_64.zip`

3. Extract the archive:
   ```bash
   # macOS/Linux
   tar -xzf canvas_*.tar.gz

   # Windows - use your preferred extraction tool
   ```

4. Move the binary to your PATH:
   ```bash
   # macOS/Linux
   sudo mv canvas /usr/local/bin/

   # Windows - add to PATH or move to C:\Windows\System32\
   ```

5. Verify installation:
   ```bash
   canvas version
   ```

### Method 4: Docker

```bash
# Run Canvas CLI in Docker
docker run ghcr.io/jjuanrivvera/canvas-cli:latest version

# Create an alias for easier use
alias canvas='docker run -it --rm -v ~/.canvas-cli:/root/.canvas-cli ghcr.io/jjuanrivvera/canvas-cli:latest'

# Now use as normal
canvas courses list
```

### Method 5: Build from Source

```bash
# Clone the repository
git clone https://github.com/jjuanrivvera/canvas-cli.git
cd canvas-cli

# Build
make build

# Install
make install

# Verify
canvas version
```

## Shell Completion

Enable tab completion for your shell:

### Bash

```bash
# Generate completion script
canvas completion bash > /etc/bash_completion.d/canvas

# Or for user-level installation
canvas completion bash > ~/.canvas-completion.bash
echo 'source ~/.canvas-completion.bash' >> ~/.bashrc
```

### Zsh

```bash
# Generate completion script
canvas completion zsh > "${fpath[1]}/_canvas"

# Reload completions
autoload -U compinit && compinit
```

### Fish

```bash
# Generate completion script
canvas completion fish > ~/.config/fish/completions/canvas.fish
```

### PowerShell

```powershell
# Generate completion script
canvas completion powershell | Out-String | Invoke-Expression

# Add to profile for persistence
canvas completion powershell >> $PROFILE
```

## Verify Installation

After installation, verify everything is working:

```bash
# Check version
canvas version

# Run diagnostics
canvas doctor

# Test authentication (will prompt for credentials)
canvas auth login --instance https://canvas.instructure.com
```

## Troubleshooting

### Command not found

If you get `command not found`, ensure the installation directory is in your PATH:

```bash
# Check PATH
echo $PATH

# Add to PATH (macOS/Linux)
export PATH="$PATH:/usr/local/bin"

# Make permanent by adding to ~/.bashrc or ~/.zshrc
echo 'export PATH="$PATH:/usr/local/bin"' >> ~/.bashrc
```

### Permission denied

If you get permission errors:

```bash
# Make binary executable
chmod +x /path/to/canvas

# Or use sudo for installation
sudo mv canvas /usr/local/bin/
```

### macOS Security Warning

On macOS, you may need to allow the app in System Preferences:

1. Try to run `canvas version`
2. Go to **System Preferences > Security & Privacy**
3. Click **"Allow Anyway"** for Canvas CLI
4. Run the command again

## Updating

### Homebrew

```bash
brew upgrade canvas-cli
```

### Go

```bash
go install github.com/jjuanrivvera/canvas-cli/cmd/canvas@latest
```

### Binary

Download and replace the binary with the latest version from the releases page.

### Docker

```bash
docker pull ghcr.io/jjuanrivvera/canvas-cli:latest
```

## Uninstalling

### Homebrew

```bash
brew uninstall canvas-cli
brew untap jjuanrivvera/canvas-cli
```

### Go

```bash
rm $(which canvas)
```

### Binary

```bash
# Find and remove the binary
sudo rm /usr/local/bin/canvas

# Remove configuration and cache
rm -rf ~/.canvas-cli
```

### Docker

```bash
docker rmi ghcr.io/jjuanrivvera/canvas-cli:latest
```

## Next Steps

After installation, continue with:
- [Authentication Guide](AUTHENTICATION.md) - Set up OAuth
- [Command Reference](COMMANDS.md) - Learn available commands
- [Examples](EXAMPLES.md) - See common use cases
