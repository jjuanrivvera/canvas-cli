<p align="center">
  <img src="docs/assets/images/logo.svg" alt="Canvas CLI" width="280">
</p>

<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go" alt="Go Version"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License"></a>
  <a href="https://github.com/jjuanrivvera/canvas-cli/releases"><img src="https://img.shields.io/github/v/release/jjuanrivvera/canvas-cli" alt="Release"></a>
  <a href="https://goreportcard.com/report/github.com/jjuanrivvera/canvas-cli"><img src="https://goreportcard.com/badge/github.com/jjuanrivvera/canvas-cli" alt="Go Report Card"></a>
</p>

<p align="center">
  A powerful command-line interface for <a href="https://www.instructure.com/canvas">Canvas LMS</a>, built with Go.
</p>

<p align="center">
  <a href="https://jjuanrivvera.github.io/canvas-cli/"><strong>Documentation</strong></a> ·
  <a href="https://jjuanrivvera.github.io/canvas-cli/getting-started/installation/"><strong>Installation</strong></a> ·
  <a href="https://jjuanrivvera.github.io/canvas-cli/commands/"><strong>Commands</strong></a>
</p>

---

## Features

- **Secure Authentication** - OAuth 2.0 with PKCE, system keyring integration
- **Multi-Instance** - Manage multiple Canvas instances from one CLI
- **Smart Rate Limiting** - Adaptive throttling based on API quotas
- **Multiple Outputs** - Table, JSON, YAML, and CSV formats
- **Interactive Mode** - REPL shell with command history and completion
- **280+ Commands** - Full coverage of Canvas LMS resources
- **MCP Server** - Use as an AI agent tool via Model Context Protocol

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap jjuanrivvera/canvas-cli
brew install canvas-cli
```

### Go Install

Requires Go 1.25+ (or Go 1.24+ with automatic toolchain download).

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

## MCP Server Mode

Canvas CLI can also run as an [MCP](https://modelcontextprotocol.io/) server, exposing all 253 commands as tools for AI coding agents (Claude Code, Cursor, VS Code Copilot).

```bash
# Start as STDIO MCP server
canvas mcp start

# Start as HTTP MCP server
canvas mcp stream --port 8080

# Export all tool schemas to JSON
canvas mcp tools

# Auto-configure in your editor
canvas mcp claude enable
canvas mcp vscode enable
canvas mcp cursor enable
```

The same binary, two interfaces. When used as an MCP server, each CLI command becomes an MCP tool with typed parameters derived from the command's flags. Required flags become required schema properties. All output goes through structured JSON.

Sensitive flags (`--show-token`, `--config`) are automatically excluded from MCP exposure.

> **Note for `go install` users**: MCP support requires Go 1.25+ due to the [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) dependency. Homebrew and binary downloads are not affected by this requirement.

For full setup (Claude Desktop, Claude Code CLI, Cursor, VS Code, OpenCode, Codex), auth precedence, and troubleshooting, see:

- [MCP Integration Guide](https://jjuanrivvera.github.io/canvas-cli/user-guide/mcp/)

## Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[MIT License](LICENSE)
