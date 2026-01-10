# Contributing to Canvas CLI

Thank you for your interest in contributing to Canvas CLI! This document provides guidelines and instructions for contributing.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, but recommended)

### Setting Up Your Development Environment

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/canvas-cli.git
   cd canvas-cli
   ```

3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/jjuanrivvera/canvas-cli.git
   ```

4. Install dependencies:
   ```bash
   make deps
   ```

5. Build the project:
   ```bash
   make build
   ```

## Development Workflow

### Branch Strategy

We follow a simplified Git Flow:

```
main ─────●─────●─────●─────●───── (stable releases)
           \         /
            \       /
develop ─────●───●───●───●─────── (integration branch)
              \     /
               \   /
feature/xyz ────●─●─────────────── (your work)
```

**Main Branches:**
- `main` - Production-ready code, tagged releases only
- `develop` - Integration branch for features (PR target)

**Supporting Branches:**
- `feature/*` - New features → merge to `develop`
- `fix/*` - Bug fixes → merge to `develop`
- `hotfix/*` - Urgent production fixes → merge to `main` and `develop`
- `release/*` - Release preparation → merge to `main` and `develop`

### Creating a Branch

Always create a new branch from `develop`:

```bash
git checkout develop
git pull upstream develop
git checkout -b feature/your-feature-name
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `hotfix/` - Urgent production fixes
- `release/` - Release preparation
- `docs/` - Documentation changes
- `refactor/` - Code refactoring
- `test/` - Adding or updating tests

### Making Changes

1. Make your changes in your branch
2. Follow the code style guidelines (see below)
3. Add tests for new functionality
4. Update documentation as needed

### Testing

Run tests before committing:

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test -v ./internal/api/...
```

### Code Style

- Follow standard Go conventions
- Use `gofmt` to format your code:
  ```bash
  make fmt
  ```

- Run the linter:
  ```bash
  make lint
  ```

- Run `go vet`:
  ```bash
  make vet
  ```

### Committing Changes

We follow conventional commit messages:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Example:
```
feat(api): add assignments CRUD operations

Implement Create, Read, Update, and Delete operations for Canvas
assignments with proper error handling and pagination support.

Closes #123
```

### Submitting a Pull Request

1. Push your branch to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. Open a pull request on GitHub

3. Ensure your PR:
   - Has a clear title and description
   - References any related issues
   - Includes tests for new functionality
   - Passes all CI checks
   - Updates documentation as needed

## Code Organization

```
canvas-cli/
├── cmd/canvas/           # Main entry point
├── commands/             # Cobra CLI commands
├── internal/             # Private packages
│   ├── api/             # Canvas API client
│   ├── auth/            # Authentication & token storage
│   ├── config/          # Configuration management
│   ├── cache/           # Caching system
│   ├── batch/           # Batch operations
│   ├── output/          # Output formatters
│   └── ...
├── pkg/                 # Public packages (reusable by other projects)
└── test/                # Test fixtures and integration tests
```

### Package Guidelines

- **internal/**: Private packages that shouldn't be imported by external projects
- **pkg/**: Public, reusable packages
- **commands/**: CLI command implementations
- Each package should have a clear, single responsibility

## Adding New Features

### Adding a New Command

1. Create a new file in `commands/` (e.g., `commands/users.go`)
2. Define your command using Cobra:
   ```go
   var usersCmd = &cobra.Command{
       Use:   "users",
       Short: "Manage Canvas users",
       RunE:  runUsers,
   }

   func init() {
       rootCmd.AddCommand(usersCmd)
   }
   ```

3. Implement the command logic
4. Add tests in `commands/users_test.go`
5. Update documentation

### Adding API Endpoints

1. Add types to `internal/api/types.go` if needed
2. Create or update the service file (e.g., `internal/api/users.go`)
3. Implement methods following existing patterns
4. Add data normalization in `internal/api/normalize.go`
5. Add tests

### Adding a New Output Format

1. Implement the `Formatter` interface in `internal/output/`
2. Register the format in `NewFormatter()`
3. Add tests for the new format
4. Update documentation

## Testing Guidelines

### Unit Tests

- Test file names should match source files with `_test.go` suffix
- Use table-driven tests when appropriate
- Mock external dependencies
- Aim for 90%+ code coverage

Example:
```go
func TestNewClient(t *testing.T) {
    tests := []struct {
        name    string
        config  ClientConfig
        wantErr bool
    }{
        {
            name: "valid config",
            config: ClientConfig{
                BaseURL: "https://canvas.example.com",
                Token:   "test-token",
            },
            wantErr: false,
        },
        // Add more test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := NewClient(tt.config)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Integration Tests

- Use VCR cassettes to record/replay HTTP interactions
- Never commit real credentials or PII
- Use synthetic test data

## Documentation

### Code Documentation

- Add package documentation to the first file in each package
- Document exported functions, types, and constants
- Use complete sentences in comments
- Follow Go doc conventions

Example:
```go
// Client is the Canvas API client with adaptive rate limiting
// and automatic retry logic. It implements the Canvas LMS REST API
// specification.
type Client struct {
    // ...
}

// NewClient creates a new Canvas API client with the given configuration.
// It automatically detects the Canvas version and configures appropriate
// rate limiting based on the instance's quota.
func NewClient(config ClientConfig) (*Client, error) {
    // ...
}
```

### User Documentation

- Update `README.md` for user-facing changes
- Add examples for new commands
- Document new configuration options
- Update the help text in commands

## Release Process

Releases are automated through GitHub Actions:

1. Update version in appropriate files
2. Create and push a version tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

3. GitHub Actions will:
   - Run tests
   - Build binaries for all platforms
   - Create a GitHub release
   - Upload binaries and checksums

## Getting Help

- 📖 Read the [README](README.md) and [SPECIFICATION](SPECIFICATION.md)
- 💬 Ask questions in GitHub Discussions
- 🐛 Report bugs in GitHub Issues
- 📧 Contact maintainers

## Code of Conduct

### Our Standards

- Be respectful and inclusive
- Welcome newcomers
- Provide constructive feedback
- Focus on what's best for the community
- Show empathy towards others

### Our Responsibilities

Maintainers are responsible for:
- Clarifying standards of acceptable behavior
- Taking appropriate action in response to unacceptable behavior
- Moderating comments, commits, code, and other contributions

## License

By contributing to Canvas CLI, you agree that your contributions will be licensed under the MIT License.

## Recognition

Contributors will be recognized in:
- GitHub contributors page
- Release notes
- Project README (for significant contributions)

Thank you for contributing to Canvas CLI! 🎉
