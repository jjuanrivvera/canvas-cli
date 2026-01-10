# Canvas CLI Documentation

Welcome to the Canvas CLI documentation. This guide will help you get started with using the Canvas CLI to automate your Canvas LMS workflows.

## Table of Contents

1. [Installation](INSTALLATION.md) - How to install Canvas CLI
2. [Authentication](AUTHENTICATION.md) - Setting up OAuth authentication
3. [Commands](COMMANDS.md) - Complete command reference
4. [Examples](EXAMPLES.md) - Common usage examples
5. [Architecture](ARCHITECTURE.md) - Internal design and components

## Quick Start

### Install

```bash
# Using Homebrew (macOS/Linux)
brew tap jjuanrivvera/canvas-cli
brew install canvas-cli

# Using Go
go install github.com/jjuanrivvera/canvas-cli/cmd/canvas@latest

# Download binary
# Visit https://github.com/jjuanrivvera/canvas-cli/releases
```

### Authenticate

```bash
# Option 1: Interactive login (OAuth)
canvas auth login --instance https://canvas.instructure.com

# Option 2: Environment variables (CI/CD)
export CANVAS_URL="https://canvas.instructure.com"
export CANVAS_TOKEN="your-access-token"
```

### Basic Usage

```bash
# List your courses
canvas courses list

# Get course details
canvas courses get 12345

# List assignments
canvas assignments list --course 12345

# Grade a submission
canvas assignments grade --course 12345 --assignment 67890 --user 11111 --score 95
```

## Features

- **Multi-instance Support**: Manage multiple Canvas installations
- **OAuth 2.0 Authentication**: Secure authentication with automatic token refresh
- **Environment Variable Auth**: CI/CD-friendly authentication via CANVAS_URL and CANVAS_TOKEN
- **Batch Operations**: Process multiple items efficiently with progress tracking
- **CSV Bulk Grading**: Import grades from spreadsheets
- **Interactive REPL**: Explore Canvas interactively
- **Shell Completion**: Tab completion for bash, zsh, fish, PowerShell
- **Smart Caching**: Reduce API calls and improve performance

## Getting Help

```bash
# Get help for any command
canvas --help
canvas courses --help
canvas assignments grade --help

# Run diagnostics
canvas doctor

# Check authentication status
canvas auth status
```

## Support

- GitHub Issues: https://github.com/jjuanrivvera/canvas-cli/issues
- Documentation: https://github.com/jjuanrivvera/canvas-cli/tree/main/docs

## License

MIT License - see [LICENSE](../LICENSE) for details
