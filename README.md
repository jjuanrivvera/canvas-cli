# Canvas CLI

A powerful command-line interface for Canvas LMS, built with Go.

## Features

✨ **Core Features**
- 🔐 OAuth 2.0 with PKCE authentication (local callback + OOB fallback)
- 🔑 Secure token storage (system keyring + encrypted file fallback)
- 🌐 Multi-instance support
- ⚡ Adaptive rate limiting (respects Canvas API quotas)
- 📄 Automatic pagination handling
- 🔄 Exponential backoff retry logic
- 📊 Multiple output formats (table, JSON, YAML, CSV)
- 🎯 Canvas version detection for compatibility

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/jjuanrivvera/canvas-cli.git
cd canvas-cli

# Build the CLI
go build -o bin/canvas ./cmd/canvas

# Optionally, move to your PATH
sudo mv bin/canvas /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/jjuanrivvera/canvas-cli/cmd/canvas@latest
```

## Quick Start

### 1. Authenticate with Canvas

```bash
# Login to your Canvas instance
canvas auth login https://canvas.instructure.com

# Check authentication status
canvas auth status
```

### 2. List Your Courses

```bash
# List all courses
canvas courses list

# List active courses where you're a teacher
canvas courses list --enrollment-type teacher --enrollment-state active
```

### 3. Get Course Details

```bash
# Get details of a specific course
canvas courses get 123456

# Include additional data
canvas courses get 123456 --include syllabus_body,term
```

## Configuration

Canvas CLI stores its configuration in `~/.canvas-cli/config.yaml`:

```yaml
default_instance: myschool
instances:
  myschool:
    name: myschool
    url: https://myschool.instructure.com
    client_id: your-client-id
settings:
  default_output_format: table
  requests_per_second: 5.0
  cache_enabled: true
  cache_ttl_minutes: 15
  telemetry_enabled: false
  log_level: info
```

## Authentication

Canvas CLI supports OAuth 2.0 with PKCE for secure authentication:

### Local Callback Mode (Default)

The CLI starts a local HTTP server to receive the OAuth callback:

```bash
canvas auth login https://canvas.instructure.com
```

### Out-of-Band Mode

For environments where local servers aren't available (e.g., SSH sessions):

```bash
canvas auth login https://canvas.instructure.com --mode oob
```

### Token Storage

Tokens are stored securely using:
1. **System Keyring** (primary): macOS Keychain, Windows Credential Manager, Linux Secret Service
2. **Encrypted File** (fallback): AES-256-GCM encrypted with machine ID + username

## Available Commands

### Authentication

```bash
canvas auth login <instance-url>     # Authenticate with Canvas
canvas auth logout [instance]        # Logout from Canvas
canvas auth status [instance]        # Check authentication status
```

### Courses

```bash
canvas courses list                  # List all courses
canvas courses get <id>              # Get course details
canvas courses users <id>            # List course users
```

### Assignments

```bash
canvas assignments list <course-id>             # List assignments
canvas assignments get <course-id> <id>         # Get assignment details
canvas assignments create <course-id>           # Create assignment
canvas assignments update <course-id> <id>      # Update assignment
canvas assignments bulk-update <course-id>      # Bulk update assignments
```

### Users

```bash
canvas users me                      # Get current user info
canvas users list <course-id>        # List users in a course
canvas users get <id>                # Get user details
canvas users create                  # Create a new user
canvas users update <id>             # Update user
```

### Enrollments

```bash
canvas enrollments list <course-id>  # List enrollments
canvas enrollments create <course-id> # Create enrollment
```

### Submissions

```bash
canvas submissions list <course-id> <assignment-id>        # List submissions
canvas submissions get <course-id> <assignment-id> <id>    # Get submission
canvas submissions grade <course-id> <assignment-id> <id>  # Grade submission
```

### Files

```bash
canvas files upload <path>           # Upload file
canvas files download <id> <path>    # Download file
```

### Utilities

```bash
canvas shell                         # Start interactive REPL mode
canvas doctor                        # Run diagnostics
canvas webhook listen                # Start webhook listener
canvas version                       # Show version info
```

### Global Flags

```bash
--instance string    # Canvas instance URL (overrides config)
--output string      # Output format: table, json, yaml, csv (default "table")
--verbose            # Enable verbose output
--config string      # Config file (default is $HOME/.canvas-cli/config.yaml)
```

## Development Status

**v1.0 - Production Ready** ✅

All core features have been implemented and tested with 90% test coverage:

- ✅ OAuth 2.0 with PKCE authentication (local + OOB modes)
- ✅ Secure token storage (keyring + encrypted file fallback)
- ✅ Multi-instance configuration management
- ✅ Core API client with adaptive rate limiting
- ✅ Automatic pagination handling
- ✅ Retry logic with exponential backoff
- ✅ Canvas version detection and compatibility
- ✅ Data normalization for consistent API responses
- ✅ Auth commands (login, logout, status)
- ✅ Courses commands (list, get, users)
- ✅ Assignments commands (list, get, create, update, bulk update)
- ✅ Users commands (list, get, create, update, me)
- ✅ Enrollments commands (list, create)
- ✅ Submissions commands (list, get, grade, bulk grade)
- ✅ Files commands (upload with resumable support, download)
- ✅ Batch operations with progress tracking
- ✅ REPL/Shell mode for interactive use
- ✅ Multiple output formatters (table, JSON, YAML, CSV)
- ✅ Smart caching system with TTL
- ✅ Webhook listener for real-time events
- ✅ Diagnostics tools (doctor command)
- ✅ Telemetry (opt-in anonymous usage tracking)
- ✅ Comprehensive test suite (90% coverage)

**Deferred to v1.1+:**
- Canvas Studio integration
- Quizzes module
- Announcements and Discussions

## Architecture

```
canvas-cli/
├── cmd/canvas/           # Main entry point
├── commands/             # Cobra commands
├── internal/
│   ├── api/             # Canvas API client
│   ├── auth/            # OAuth & token storage
│   ├── config/          # Configuration management
│   ├── cache/           # Caching system
│   ├── batch/           # Batch operations
│   ├── repl/            # Interactive REPL
│   ├── output/          # Output formatters
│   └── ...
├── pkg/                 # Public packages
└── test/                # Tests
```

## Technology Stack

- **Language**: Go 1.21+
- **CLI Framework**: Cobra
- **Configuration**: Viper
- **OAuth**: golang.org/x/oauth2
- **Rate Limiting**: golang.org/x/time/rate
- **Keyring**: zalando/go-keyring
- **Logging**: log/slog (stdlib)

## API Rate Limiting

Canvas CLI implements adaptive rate limiting:

- **Default**: 5 requests/second
- **50% quota**: Slows to 2 requests/second (warning)
- **20% quota**: Slows to 1 request/second (critical)

The CLI automatically adjusts its rate based on Canvas API quota headers.

## Security

- **OAuth 2.0 with PKCE**: Industry-standard authentication flow
- **Secure Token Storage**: System keyring with encrypted file fallback
- **AES-256-GCM Encryption**: User-derived encryption keys
- **No Hardcoded Credentials**: All credentials are user-provided

## Documentation

For comprehensive guides and examples, see the [docs/](docs/) directory:

- **[Installation Guide](docs/INSTALLATION.md)** - Detailed installation instructions for all platforms
- **[Authentication Guide](docs/AUTHENTICATION.md)** - OAuth setup, security, and multi-instance management
- **[Command Reference](docs/COMMANDS.md)** - Complete command documentation with all flags and options
- **[Usage Examples](docs/EXAMPLES.md)** - Practical examples for common workflows and automation

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

[MIT License](LICENSE)

## Acknowledgments

Built with ❤️ for the Canvas LMS community.

Based on the [Canvas LMS REST API](https://canvas.instructure.com/doc/api/).
