# Canvas CLI

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/jjuanrivvera/canvas-cli)](https://github.com/jjuanrivvera/canvas-cli/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/jjuanrivvera/canvas-cli)](https://goreportcard.com/report/github.com/jjuanrivvera/canvas-cli)

A powerful command-line interface for [Canvas LMS](https://www.instructure.com/canvas), built with Go.

**[Documentation](https://jjuanrivvera.github.io/canvas-cli/)** | **[Installation](https://jjuanrivvera.github.io/canvas-cli/getting-started/installation/)** | **[Commands](https://jjuanrivvera.github.io/canvas-cli/commands/)**

## Features

- **Secure Authentication** - OAuth 2.0 with PKCE, system keyring integration
- **Multi-Instance** - Manage multiple Canvas instances from one CLI
- **Smart Rate Limiting** - Adaptive throttling based on API quotas
- **Multiple Outputs** - Table, JSON, YAML, and CSV formats
- **Interactive Mode** - REPL shell with command history and completion
- **89+ Commands** - Full coverage of Canvas LMS resources

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap jjuanrivvera/canvas-cli
brew install canvas-cli
```

### Go Install

```bash
go install github.com/jjuanrivvera/canvas-cli/cmd/canvas@latest
```

### Binary Download

Download from [GitHub Releases](https://github.com/jjuanrivvera/canvas-cli/releases).

## Quick Start

```bash
# Authenticate with your Canvas instance
canvas auth login https://your-school.instructure.com

# List your courses
canvas courses list

# Get assignments for a course
canvas assignments list <course-id>

# Start interactive mode
canvas shell
```

## Command Overview

| Category | Commands |
|----------|----------|
| **Auth** | `login`, `logout`, `status` |
| **Courses** | `list`, `get`, `users` |
| **Assignments** | `list`, `get`, `create`, `update`, `bulk-update` |
| **Submissions** | `list`, `get`, `grade`, `bulk-grade` |
| **Users** | `me`, `list`, `get`, `create`, `update` |
| **Modules** | `list`, `get`, `create`, `update`, `delete`, `items` |
| **Pages** | `list`, `get`, `create`, `update`, `delete`, `front` |
| **Discussions** | `list`, `get`, `create`, `entries`, `post`, `reply` |
| **Announcements** | `list`, `get`, `create`, `update`, `delete` |
| **Calendar** | `list`, `get`, `create`, `update`, `delete` |
| **Planner** | `items`, `notes`, `complete`, `dismiss` |
| **Files** | `upload`, `download` |
| **Utilities** | `shell`, `doctor`, `webhook`, `version` |

See [full command reference](https://jjuanrivvera.github.io/canvas-cli/commands/) for all options and flags.

## Configuration

```yaml
# ~/.canvas-cli/config.yaml
default_instance: myschool
instances:
  myschool:
    url: https://myschool.instructure.com
    client_id: your-client-id
settings:
  default_output_format: table
  cache_enabled: true
```

See [Authentication Guide](https://jjuanrivvera.github.io/canvas-cli/getting-started/authentication/) for detailed setup.

## Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[MIT License](LICENSE)
