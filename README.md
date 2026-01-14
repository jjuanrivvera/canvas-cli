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
- **280+ Commands** - Full coverage of Canvas LMS resources

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
| **Courses** | `list`, `get`, `create`, `update`, `delete` |
| **Assignments** | `list`, `get`, `create`, `update`, `delete`, `bulk-update` |
| **Submissions** | `list`, `get`, `grade`, `bulk-grade`, `comments` |
| **Users** | `me`, `list`, `get`, `create`, `update` |
| **Enrollments** | `list`, `get`, `create`, `update`, `delete`, `accept` |
| **Modules** | `list`, `get`, `create`, `update`, `delete`, `publish`, `items` |
| **Pages** | `list`, `get`, `create`, `update`, `delete`, `front`, `revisions` |
| **Discussions** | `list`, `get`, `create`, `entries`, `post`, `reply`, `subscribe` |
| **Announcements** | `list`, `get`, `create`, `update`, `delete` |
| **Quizzes** | `list`, `get`, `create`, `update`, `delete`, `questions`, `submissions` |
| **Grades** | `summary`, `history`, `bulk-update`, `final`, `current` |
| **Groups** | `list`, `get`, `create`, `update`, `delete`, `users`, `categories` |
| **Outcomes** | `list`, `get`, `create`, `update`, `delete`, `groups`, `results` |
| **Rubrics** | `list`, `get`, `create`, `update`, `delete`, `associate` |
| **Conversations** | `list`, `get`, `create`, `reply`, `archive`, `star`, `batch-update` |
| **Calendar** | `list`, `get`, `create`, `update`, `delete`, `reserve` |
| **Files** | `list`, `get`, `upload`, `download`, `delete` |
| **Sections** | `list`, `get`, `create`, `update`, `delete`, `crosslist` |
| **Admin** | `admins`, `roles`, `analytics`, `blueprint`, `sis-imports` |
| **Utilities** | `shell`, `doctor`, `webhook`, `api`, `version` |

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
