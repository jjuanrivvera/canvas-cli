# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-01-09

### Added

#### Core Functionality
- OAuth 2.0 authentication with PKCE support
- Local callback server mode for OAuth flow
- Out-of-band (OOB) OAuth flow fallback for SSH/remote environments
- Secure token storage using system keyring (macOS Keychain, Windows Credential Manager, Linux Secret Service)
- Encrypted file storage fallback with AES-256-GCM encryption
- User-derived encryption keys from machine ID + username
- Multi-instance configuration management
- Canvas version detection and compatibility handling

#### API Client Features
- Comprehensive Canvas LMS API client
- Adaptive rate limiting (5 req/sec → 2 req/sec → 1 req/sec based on quota)
- Automatic pagination handling for large result sets
- Exponential backoff retry logic with 3 max retries
- Data normalization for consistent API responses
- Custom error types with helpful suggestions and documentation links
- Request/response logging with --debug flag

#### Commands - Authentication
- `canvas auth login` - OAuth 2.0 authentication flow
- `canvas auth logout` - Logout and clear credentials
- `canvas auth status` - Check authentication status

#### Commands - Courses
- `canvas courses list` - List courses with filtering options
- `canvas courses get` - Get course details with includes
- `canvas courses users` - List users in a course

#### Commands - Assignments
- `canvas assignments list` - List assignments in a course
- `canvas assignments get` - Get assignment details
- `canvas assignments create` - Create new assignment with full parameter support
- `canvas assignments update` - Update assignment with pointer types for optional fields
- `canvas assignments bulk-update` - Bulk update multiple assignments

#### Commands - Users
- `canvas users me` - Get current authenticated user
- `canvas users list` - List users with filtering
- `canvas users get` - Get user details
- `canvas users create` - Create new user with pseudonym and communication channel
- `canvas users update` - Update user with avatar support

#### Commands - Enrollments
- `canvas enrollments list` - List enrollments in course/section
- `canvas enrollments create` - Create new enrollment

#### Commands - Submissions
- `canvas submissions list` - List submissions for assignment
- `canvas submissions get` - Get submission details
- `canvas submissions grade` - Grade individual submission
- `canvas submissions bulk-grade` - Bulk grade from CSV

#### Commands - Files
- `canvas files upload` - Upload files with progress tracking
- `canvas files download` - Download files with resumable support

#### Advanced Features
- **REPL Mode**: Interactive shell with command history, tab completion, and syntax highlighting
- **Smart Caching**: TTL-based caching (courses: 15min, users: 5min, assignments: 10min)
- **Batch Operations**: Concurrent processing with progress bars and error collection
- **Webhook Listener**: Real-time webhook event handling with signature verification
- **Diagnostics**: `canvas doctor` command for health checks and troubleshooting
- **Telemetry**: Opt-in anonymous usage tracking for feature prioritization

#### Output Formats
- Table format (ASCII tables with proper truncation)
- JSON format (structured output)
- YAML format (human-readable)
- CSV format (for data export)

#### Developer Features
- Comprehensive test suite with 90% coverage
- HTTP request/response recording for tests
- Mock Canvas API server for testing
- Synthetic test data (no PII in test fixtures)
- Race condition detection in tests
- CI/CD ready with stable exit codes

### Testing
- 90% test coverage for core functionality (89.9% weighted average)
- 8 out of 9 packages at 90%+ coverage
- All tests passing (100% pass rate)
- No race conditions detected
- Comprehensive parameter testing for all API operations
- Edge case coverage for error scenarios
- Mock HTTP server testing with httptest

### Security
- OAuth 2.0 with PKCE (Proof Key for Code Exchange)
- Secure credential storage with system keyring integration
- AES-256-GCM encryption for file-based token storage
- User-derived encryption keys (never stored)
- Webhook signature verification with HMAC-SHA256
- No hardcoded credentials
- No sensitive data in logs or cache

### Performance
- Adaptive rate limiting respects Canvas API quotas
- Smart caching reduces redundant API calls
- Concurrent batch operations (5 concurrent by default)
- Automatic pagination for large datasets
- Efficient memory usage (<100MB for 10,000 cached items)
- Progress indicators for operations >3 seconds

### Documentation
- Comprehensive README with quick start guide
- SPECIFICATION.md with complete architecture documentation
- CONTRIBUTING.md with development guidelines
- PROJECT_STATUS.md tracking implementation progress
- COVERAGE_REPORT.md with detailed test coverage metrics
- Inline code documentation with examples

### Infrastructure
- Go 1.21+ support with modern stdlib features (log/slog)
- Cross-platform support (macOS, Linux, Windows)
- Cobra CLI framework for command structure
- Viper for configuration management
- Standard Go project layout

## [Unreleased]

### Planned for v1.1+
- Canvas Studio integration
- Quizzes module
- Announcements and Discussions
- GraphQL API support
- Shell completion generation
- Man page generation

## Version History

- **1.0.0** (2026-01-09) - Initial production release
  - All v1.0 scope features complete
  - 90% test coverage achieved
  - Production-ready and stable

---

For more details on each change, see the [commit history](https://github.com/jjuanrivvera/canvas-cli/commits/main).

For planned features and roadmap, see [SPECIFICATION.md](SPECIFICATION.md).
