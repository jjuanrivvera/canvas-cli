# Development

Resources for contributing to Canvas CLI development.

## Overview

Canvas CLI is built with:

- **Go 1.21+** - Primary language
- **Cobra** - CLI framework
- **Viper** - Configuration management

## Topics

<div class="grid cards" markdown>

-   :material-folder-multiple:{ .lg .middle } **Architecture**

    ---

    Learn about Canvas CLI's internal structure and design decisions

    [:octicons-arrow-right-24: Architecture](architecture.md)

-   :material-source-pull:{ .lg .middle } **Contributing**

    ---

    Guidelines for contributing code, documentation, and bug reports

    [:octicons-arrow-right-24: Contributing](contributing.md)

</div>

## Quick Start for Contributors

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional but recommended)

### Setup

```bash
# Clone the repository
git clone https://github.com/jjuanrivvera/canvas-cli.git
cd canvas-cli

# Install dependencies
go mod download

# Build
make build

# Run tests
make test
```

### Development Workflow

```bash
# Create a feature branch
git checkout -b feature/my-feature

# Make changes and run tests
make dev
make test

# Commit and push
git add .
git commit -m "feat: add my feature"
git push origin feature/my-feature
```

### Available Make Targets

| Target | Description |
|--------|-------------|
| `make build` | Build binary to `bin/canvas` |
| `make dev` | Build with fmt and vet |
| `make test` | Run all tests |
| `make test-coverage` | Run tests with coverage |
| `make lint` | Run golangci-lint |
| `make fmt` | Format code |
| `make install` | Install to `/usr/local/bin` |

## Project Structure

```
canvas-cli/
├── cmd/canvas/       # Entry point
├── commands/         # Cobra commands
├── internal/
│   ├── api/          # Canvas API client
│   ├── auth/         # Authentication
│   ├── config/       # Configuration
│   ├── cache/        # Response caching
│   ├── batch/        # Batch operations
│   └── output/       # Output formatters
├── tools/            # Development tools
└── docs/             # Documentation
```

## Code Style

Canvas CLI follows standard Go conventions:

- Run `gofmt` before committing
- Use `golangci-lint` for static analysis
- Write tests for new features
- Document exported functions

## Getting Help

- [GitHub Issues](https://github.com/jjuanrivvera/canvas-cli/issues) - Bug reports and feature requests
- [GitHub Discussions](https://github.com/jjuanrivvera/canvas-cli/discussions) - Questions and ideas
