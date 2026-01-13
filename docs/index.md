# Canvas CLI

A powerful command-line interface for Canvas LMS.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

## Features

- **Comprehensive API Coverage** - Access courses, assignments, submissions, users, modules, pages, and more
- **Multiple Output Formats** - Table, JSON, YAML, and CSV output
- **Bulk Operations** - Grade submissions in bulk from CSV files
- **Course Synchronization** - Sync content between Canvas instances
- **Intelligent Caching** - Fast responses with automatic cache invalidation
- **Secure Authentication** - OAuth 2.0 with PKCE flow

## Quick Start

=== "macOS (Homebrew)"

    ```bash
    brew tap jjuanrivvera/canvas-cli
    brew install canvas-cli
    ```

=== "Go Install"

    ```bash
    go install github.com/jjuanrivvera/canvas-cli/cmd/canvas@latest
    ```

=== "Binary Download"

    Download the latest release from [GitHub Releases](https://github.com/jjuanrivvera/canvas-cli/releases).

Then authenticate with Canvas:

```bash
canvas auth login
```

## Example Usage

```bash
# List your courses
canvas courses list

# List assignments for a course
canvas assignments list --course-id 123

# Grade a submission
canvas submissions grade 456 --course-id 123 --score 95

# Bulk grade from CSV
canvas submissions bulk-grade --course-id 123 --csv grades.csv

# Export data as JSON
canvas users list --course-id 123 --output json
```

## Documentation

<div class="grid cards" markdown>

-   :material-rocket-launch:{ .lg .middle } **Getting Started**

    ---

    Install Canvas CLI and authenticate with your Canvas instance

    [:octicons-arrow-right-24: Installation](getting-started/installation.md)

-   :material-book-open-variant:{ .lg .middle } **User Guide**

    ---

    Learn how to configure and use Canvas CLI effectively

    [:octicons-arrow-right-24: User Guide](user-guide/index.md)

-   :material-console:{ .lg .middle } **Command Reference**

    ---

    Complete reference for all available commands

    [:octicons-arrow-right-24: Commands](commands/index.md)

-   :material-school:{ .lg .middle } **Tutorials**

    ---

    Step-by-step guides for common workflows

    [:octicons-arrow-right-24: Tutorials](tutorials/index.md)

</div>

## Support

- **Issues**: [GitHub Issues](https://github.com/jjuanrivvera/canvas-cli/issues)
- **Discussions**: [GitHub Discussions](https://github.com/jjuanrivvera/canvas-cli/discussions)

## License

Canvas CLI is released under the [MIT License](https://opensource.org/licenses/MIT).
