# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.7.0] - 2026-01-25

### Added

- **Command Aliases**: Create shortcuts for frequently used commands
  - `canvas alias set <name> "<command>"` - Create an alias
  - `canvas alias list` - List all aliases
  - `canvas alias delete <name>` - Remove an alias
  - Aliases are stored in config and expand at runtime

- **Context Management**: Set default values for common flags
  - `canvas context set <type> <id>` - Set course, assignment, user, or account context
  - `canvas context show` - Display current context
  - `canvas context clear [type]` - Clear all or specific context
  - Commands automatically use context when flags aren't provided

- **Output Filtering**: Filter and sort command output
  - `--filter <text>` - Filter results by text (case-insensitive, searches all fields)
  - `--columns <list>` - Select specific columns to display
  - `--sort <field>` - Sort by field (prefix with `-` for descending)
  - Works with all output formats (table, JSON, YAML, CSV)

- **Enhanced Dry-Run Mode**: Preview destructive operations with details
  - Delete commands show resource details before confirmation
  - Update commands show what would change
  - Works with `--dry-run` and `--force` flags

- **Curl Command Output**: See equivalent curl commands with `--dry-run`
  - Useful for debugging and learning the Canvas API
  - Token redacted by default, use `--show-token` to include

- **Aggressive Auto-Update**: Automatic update checking
  - `canvas update enable` - Enable automatic update checks
  - `canvas update disable` - Disable automatic update checks
  - `canvas update check` - Manually check for updates
  - `canvas update status` - Show update settings

### Changed

- Improved CLI UX inspired by modern tools (gh, kubectl, stripe-cli)
- Documentation updated with new feature guides

## [1.5.2] - 2026-01-14

### Added

- **Per-instance API Token Authentication**: New alternative to OAuth for simpler authentication
  - `canvas auth token set <instance>` - Configure API token for an instance
  - `canvas auth token remove <instance>` - Remove API token from an instance
  - Tokens stored in config file, can be mixed with OAuth per-instance
- **User-Agent Header**: All API requests now include `User-Agent: canvas-cli/VERSION`
  - Required by Canvas API (enforcement coming soon per Canvas changelog)
  - Includes version for debugging and analytics
- **Auth Status Improvements**: `canvas auth status` now shows authentication type (token/oauth/none)
- **Instance Helper Methods**: `HasToken()`, `HasOAuth()`, `AuthType()` for config

### Changed

- Token authentication takes precedence over OAuth when both are configured for an instance
- Improved error messages for authentication failures

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
- Architecture documentation in docs/development/architecture.md
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

### Planned
- Canvas Studio integration
- GraphQL API support

## [1.5.0] - 2026-01-14

### Added

#### 70+ New Write Commands
This release adds comprehensive write command support across all Canvas API resources:

##### Account Administration
- `canvas admins add` - Add account administrator
- `canvas admins list` - List account administrators
- `canvas admins remove` - Remove account administrator
- `canvas roles create` - Create custom role
- `canvas roles update` - Update role permissions
- `canvas roles delete` - Delete custom role
- `canvas roles list` - List account roles

##### Analytics
- `canvas analytics activity` - View course activity
- `canvas analytics assignments` - View assignment statistics
- `canvas analytics department` - View department-level analytics
- `canvas analytics students` - View student analytics
- `canvas analytics user` - View user-specific analytics

##### Assignment Groups
- `canvas assignment-groups list` - List assignment groups
- `canvas assignment-groups get` - Get assignment group details
- `canvas assignment-groups create` - Create assignment group
- `canvas assignment-groups update` - Update assignment group
- `canvas assignment-groups delete` - Delete assignment group

##### Blueprint Courses
- `canvas blueprint get` - Get blueprint details
- `canvas blueprint sync` - Sync blueprint to associated courses
- `canvas blueprint changes` - View unsynced changes
- `canvas blueprint associations list` - List associated courses
- `canvas blueprint associations add` - Add course associations
- `canvas blueprint associations remove` - Remove associations
- `canvas blueprint migrations list` - List sync history
- `canvas blueprint migrations get` - Get migration details

##### Content Migrations
- `canvas content-migrations list` - List migrations
- `canvas content-migrations get` - Get migration details
- `canvas content-migrations create` - Start content migration
- `canvas content-migrations issues` - View migration issues

##### Conversations (Inbox)
- `canvas conversations list` - List conversations
- `canvas conversations get` - Get conversation details
- `canvas conversations create` - Create new conversation
- `canvas conversations reply` - Reply to conversation
- `canvas conversations forward` - Forward conversation
- `canvas conversations add-recipients` - Add recipients
- `canvas conversations mark-read` - Mark as read
- `canvas conversations mark-unread` - Mark as unread
- `canvas conversations archive` - Archive conversation
- `canvas conversations unarchive` - Unarchive conversation
- `canvas conversations star` - Star conversation
- `canvas conversations unstar` - Unstar conversation
- `canvas conversations delete` - Delete conversation
- `canvas conversations batch-update` - Bulk update conversations

##### Courses
- `canvas courses create` - Create new course
- `canvas courses update` - Update course
- `canvas courses delete` - Delete/conclude course

##### External Tools (LTI)
- `canvas external-tools list` - List external tools
- `canvas external-tools get` - Get tool details
- `canvas external-tools create` - Create external tool
- `canvas external-tools update` - Update external tool
- `canvas external-tools delete` - Delete external tool
- `canvas external-tools sessionless-launch` - Get sessionless launch URL

##### Grades & Gradebook
- `canvas grades summary` - View grade summary
- `canvas grades history` - View grade history
- `canvas grades bulk-update` - Bulk update grades
- `canvas grades final` - Get final grades
- `canvas grades current` - Get current grades

##### Groups
- `canvas groups list` - List groups
- `canvas groups get` - Get group details
- `canvas groups create` - Create group
- `canvas groups update` - Update group
- `canvas groups delete` - Delete group
- `canvas groups users` - List group members
- `canvas groups invite` - Invite users to group
- `canvas groups join` - Join a group
- `canvas groups leave` - Leave a group
- `canvas groups categories list` - List group categories
- `canvas groups categories create` - Create category
- `canvas groups categories update` - Update category
- `canvas groups categories delete` - Delete category

##### Learning Outcomes
- `canvas outcomes list` - List outcomes
- `canvas outcomes get` - Get outcome details
- `canvas outcomes create` - Create learning outcome
- `canvas outcomes update` - Update outcome
- `canvas outcomes delete` - Delete outcome
- `canvas outcomes groups list` - List outcome groups
- `canvas outcomes groups get` - Get group details
- `canvas outcomes groups create` - Create outcome group
- `canvas outcomes groups update` - Update group
- `canvas outcomes groups delete` - Delete group
- `canvas outcomes import` - Import outcomes
- `canvas outcomes alignments` - View outcome alignments
- `canvas outcomes results` - View outcome results

##### Assignment Overrides
- `canvas overrides list` - List assignment overrides
- `canvas overrides get` - Get override details
- `canvas overrides create` - Create date/student override
- `canvas overrides update` - Update override
- `canvas overrides delete` - Delete override
- `canvas overrides batch-create` - Bulk create overrides
- `canvas overrides batch-update` - Bulk update overrides

##### Peer Reviews
- `canvas peer-reviews list` - List peer reviews
- `canvas peer-reviews create` - Assign peer review
- `canvas peer-reviews delete` - Remove peer review assignment

##### Quizzes (Classic Quizzes)
- `canvas quizzes list` - List quizzes
- `canvas quizzes get` - Get quiz details
- `canvas quizzes create` - Create quiz
- `canvas quizzes update` - Update quiz
- `canvas quizzes delete` - Delete quiz
- `canvas quizzes reorder` - Reorder quiz questions
- `canvas quizzes validate-token` - Validate access code
- `canvas quizzes questions list` - List quiz questions
- `canvas quizzes questions get` - Get question details
- `canvas quizzes questions create` - Create question
- `canvas quizzes questions update` - Update question
- `canvas quizzes questions delete` - Delete question
- `canvas quizzes submissions list` - List quiz submissions
- `canvas quizzes submissions get` - Get submission details
- `canvas quizzes submissions start` - Start quiz attempt
- `canvas quizzes submissions complete` - Complete quiz attempt

##### Rubrics
- `canvas rubrics list` - List rubrics
- `canvas rubrics get` - Get rubric details
- `canvas rubrics create` - Create rubric
- `canvas rubrics update` - Update rubric
- `canvas rubrics delete` - Delete rubric
- `canvas rubrics associations` - List rubric associations
- `canvas rubrics associate` - Associate rubric with assignment
- `canvas rubrics assessments` - View rubric assessments

##### Course Sections
- `canvas sections list` - List course sections
- `canvas sections get` - Get section details
- `canvas sections create` - Create section
- `canvas sections update` - Update section
- `canvas sections delete` - Delete section
- `canvas sections crosslist` - Cross-list section
- `canvas sections uncrosslist` - Remove cross-listing

##### SIS Imports
- `canvas sis-imports list` - List import history
- `canvas sis-imports get` - Get import details
- `canvas sis-imports create` - Start SIS import
- `canvas sis-imports abort` - Abort running import
- `canvas sis-imports restore` - Restore deleted items
- `canvas sis-imports errors` - View import errors

##### Raw API Access
- `canvas api` - Make raw API requests to any Canvas endpoint

##### Modules Improvements
- `canvas modules publish` - Publish module (convenience)
- `canvas modules unpublish` - Unpublish module (convenience)
- `canvas modules items update` - Update module item (was missing)

##### Enrollments Improvements
- `canvas enrollments create` - Create enrollment
- `canvas enrollments update` - Update enrollment state
- `canvas enrollments delete` - Delete/deactivate enrollment
- `canvas enrollments accept` - Accept enrollment invitation
- `canvas enrollments reject` - Reject enrollment invitation
- `canvas enrollments reactivate` - Reactivate enrollment

##### Submissions Improvements
- `canvas submissions update` - Update submission
- `canvas submissions summary` - Get submission summary

#### Webhook JWT Verification (Canvas Data Services)
- **JWT verification support**: Use `--canvas-data-services` flag for Instructure-hosted Canvas instances that use Canvas Data Services
- **Custom JWK endpoints**: Use `--jwks-url` for custom JWK endpoints
- **Automatic JWK caching**: Public keys are cached for 1 hour and refreshed automatically
- **Fallback mode**: Both JWT and HMAC verification can be enabled simultaneously

### Fixed

#### UX Improvements
- **JSON output for write commands**: All create/update/delete commands now properly support `-o json` output format
- **Rubrics response parsing**: Fixed issue where rubrics were wrapped in `{rubric: {...}}` envelope
- **Conversations JSON keys**: Fixed duplicate array bracket suffix `[]` in JSON request keys
- **Zero date display**: Now shows "Not set" instead of "0001-01-01 00:00:00" for unset dates
- **Empty collections**: Hidden in output instead of showing `map[]` or `[]`
- **404 error messages**: Now include descriptive text explaining what resource was not found
- **Nested struct display**: New `formatStructCompact()` for clean display of complex nested structures

### Changed
- External tools delete now requires `--force` flag for confirmation
- Courses create now accepts `--account` as alias for `--account-id`

## [1.4.0] - 2026-01-13

### Added

#### Authentication Improvements
- **Automatic OAuth Token Refresh**: Access tokens are now automatically refreshed using refresh tokens when they expire, eliminating the need for manual re-authentication
- **Instance Config Lookup**: `canvas auth login --instance <name>` now automatically loads the URL and OAuth credentials from your config file
- **Positional Instance Name**: `canvas config add` now accepts instance name as a positional argument: `canvas config add production --url https://canvas.example.com`

#### Table Output Improvements
- **Compact Table Output**: Default table output now shows only key fields for cleaner display
- **Verbose Mode**: Use `-v/--verbose` flag to see all fields in table output
- **Improved Field Selection**: Key fields are optimized for each resource type (Course, User, Assignment, etc.)
- **Instance Name Support**: The `--instance` flag now accepts instance names (not just URLs)

### Changed
- `canvas config add <name> --url <url>` syntax replaces `canvas config add --name <name> --url <url>`
- Table formatter now uses structured formatters instead of custom display functions
- Removed "Found X items:" messages in compact (non-verbose) mode

### Fixed
- Pre-commit hook now includes golangci-lint for catching lint issues before push
- Removed unused display functions that were causing lint warnings
- Documentation updated to reflect correct CLI syntax and behavior

### Developer Experience
- **Pre-commit Linting**: Added golangci-lint to pre-commit hook for early lint error detection
- **Documentation Accuracy**: Fixed documentation to match actual CLI behavior (sync command syntax, environment variables, flags)

## [1.1.0] - 2026-01-10

### Added

#### Commands - Modules
- `canvas modules list` - List modules in a course
- `canvas modules get` - Get module details
- `canvas modules create` - Create new module
- `canvas modules update` - Update module
- `canvas modules delete` - Delete module
- `canvas modules relock` - Relock module progressions
- `canvas modules items` - List items in a module
- `canvas modules items get` - Get module item details
- `canvas modules items create` - Create module item
- `canvas modules items update` - Update module item
- `canvas modules items delete` - Delete module item
- `canvas modules items done` - Mark module item as done
- `canvas modules items not-done` - Mark module item as not done

#### Commands - Pages
- `canvas pages list` - List wiki pages in a course
- `canvas pages get` - Get page by URL or ID
- `canvas pages front` - Get front page
- `canvas pages create` - Create new page
- `canvas pages update` - Update existing page
- `canvas pages delete` - Delete page
- `canvas pages duplicate` - Duplicate page
- `canvas pages revisions` - List page revisions
- `canvas pages revert` - Revert to specific revision

#### Commands - Discussions
- `canvas discussions list` - List discussion topics
- `canvas discussions get` - Get discussion details
- `canvas discussions create` - Create new discussion
- `canvas discussions update` - Update discussion
- `canvas discussions delete` - Delete discussion
- `canvas discussions entries` - List discussion entries
- `canvas discussions post` - Post new entry
- `canvas discussions reply` - Reply to entry
- `canvas discussions subscribe` - Subscribe to topic
- `canvas discussions unsubscribe` - Unsubscribe from topic

#### Commands - Announcements
- `canvas announcements list` - List course announcements
- `canvas announcements get` - Get announcement details
- `canvas announcements create` - Create new announcement
- `canvas announcements update` - Update announcement
- `canvas announcements delete` - Delete announcement

#### Commands - Calendar
- `canvas calendar list` - List calendar events
- `canvas calendar get` - Get event details
- `canvas calendar create` - Create new event
- `canvas calendar update` - Update event
- `canvas calendar delete` - Delete event
- `canvas calendar reserve` - Reserve time slot

#### Commands - Planner
- `canvas planner items` - List planner items
- `canvas planner notes list` - List planner notes
- `canvas planner notes get` - Get note details
- `canvas planner notes create` - Create planner note
- `canvas planner notes update` - Update note
- `canvas planner notes delete` - Delete note
- `canvas planner complete` - Mark item as complete
- `canvas planner dismiss` - Dismiss item from planner
- `canvas planner overrides` - List planner overrides

### Testing
- Added comprehensive tests for all new API services
- Tests for Modules, Pages, Discussions, Announcements, Calendar, and Planner
- All tests passing with consistent patterns

## Version History

- **1.5.0** (2026-01-14) - Major feature release
  - 70+ new write commands across all Canvas API resources
  - JWT verification for Canvas Data Services webhooks
  - Comprehensive UX improvements and bug fixes
- **1.4.0** (2026-01-13) - Feature release
  - Automatic OAuth token refresh
  - Instance config lookup for auth login
  - Improved table output with verbose mode
- **1.1.0** (2026-01-10) - Feature release
  - Modules, Pages, Discussions, Announcements, Calendar, Planner commands
  - Comprehensive API coverage for course content management
- **1.0.0** (2026-01-09) - Initial production release
  - All v1.0 scope features complete
  - 90% test coverage achieved
  - Production-ready and stable

---

For more details on each change, see the [commit history](https://github.com/jjuanrivvera/canvas-cli/commits/main).

For planned features and roadmap, see the [Unreleased](#unreleased) section above.
