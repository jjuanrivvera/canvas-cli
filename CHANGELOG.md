# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

This file is the single source of truth; `docs/changelog.md` is a copy kept in
sync by `make docs-gen` and the documentation workflow.

## [Unreleased]

### Added

- `auth login --public-client`: secret-less OAuth for Canvas developer keys provisioned with
  `client_type = "public"` — the token exchange and refresh are protected by PKCE only, and the
  request omits `client_secret` entirely. Canvas validates PKCE on hosted instances since the
  March 2026 release; public-client keys must currently be provisioned by Instructure (they are
  not self-service in the Developer Keys UI). Instances persist this as `public_client: true`
  in `config.yaml`. ([#48](https://github.com/jjuanrivvera/canvas-cli/issues/48),
  [#51](https://github.com/jjuanrivvera/canvas-cli/issues/51))

### Planned

- Canvas Studio integration
- GraphQL API support

## [1.10.5] - 2026-07-11

### Security

- `auth login` now reads secrets with a **hidden** prompt instead of `fmt.Scanln` — the API
  access token, the OAuth client secret, and the OAuth authorization code. Previously they were
  echoed to the terminal in plaintext (landing in scrollback) and `fmt.Scanln` could stall on a
  long pasted token. Non-secret inputs (client ID, y/n confirmations) are unchanged.

## [1.10.4] - 2026-07-11

### Added

- One-line install script for macOS/Linux (checksum-verified):
  `curl -fsSL https://raw.githubusercontent.com/jjuanrivvera/canvas-cli/main/install.sh | sh`.

### Security

- Build with Go 1.25.12 to clear GO-2026-5856 (privacy leak in `crypto/tls`
  Encrypted Client Hello).

## [1.10.3] - 2026-07-10

### Added

- Scoop (Windows) packaging via `jjuanrivvera/scoop-canvas-cli`, documented in the
  README alongside Homebrew, Go, Docker, and binary installs.

## [1.10.2] - 2026-07-02

### Fixed

- **`agent guard` hook hardening.** The generated PreToolUse hook missed several
  bypasses: path-invoked binaries (`./bin/canvas`, `/usr/local/bin/canvas`) were
  not matched; a shell separator glued to a no-arg irreversible verb
  (`canvas favorites courses reset;true`) slipped the trailing boundary; and the
  no-jq fallback could fail open because the compact JSON payload glues the
  command to its key. All three are fixed, and `sync assignments` (whose leaf
  collides with the `analytics assignments` read verb) is now correctly gated as
  a write. A different binary that merely ends in `canvas` is still not matched.
- Cleared pre-existing lint debt surfaced by the current golangci-lint
  (`reflect.Ptr` → `reflect.Pointer`; `WriteString(fmt.Sprintf(...))` →
  `fmt.Fprintf(...)` in speccheck). No behavior change.

## [1.10.1] - 2026-07-02

### Fixed

- **Documentation audited against the shipped CLI**: repaired examples that no
  longer ran — `submissions grade-batch` → `bulk-grade --csv-file` with the
  correct `user_id,assignment_id,score,comment` CSV columns, `submissions
  grade` flag usage (`--user-id`/`--score`/`--posted-grade`, no positional id),
  `enrollments create --user-id`/`--type`, positional `sync course` and
  `users search`, `assignments create --group-id`, and `config account
  <instance> <account-id>`. Corrected stale numbers (93 command groups — the
  v1.10.0 note below originally overcounted 46→93 as 52→98 — 530+ MCP tools,
  the 80% CI coverage gate) and removed the nonexistent
  `CANVAS_OUTPUT`/`CANVAS_NO_CACHE` env vars from docs and the bundled skill
  (`CANVAS_CLI_MACHINE_ID` is documented instead).
- **`canvas context` help text** no longer claims every command consumes the
  stored context: only `assignments list`/`assignments get` read the course
  context today. The user guide, best practices, and bundled agent skill now
  document the same behavior.
- **Agent-guard docs**: dropped `status` from the documented read allowlist to
  match the code's `canvasReadVerbs`; the Agent Safety guide is now linked
  from the docs landing pages.
- **Contributor docs** reflect the current repo: no `pkg/` directory or VCR
  cassettes, `make check` and the spec targets documented, architecture page
  updated with the real package layout and the spec-compliance harness.

## [1.10.0] - 2026-06-23

### Added

- **Agent safety guard** (`canvas agent guard --host claude-code|codex|opencode`):
  generates AI-agent safety config — permission rules plus a `PreToolUse` hook —
  that hard-blocks irreversible Canvas operations (delete, conclude, crosslist,
  cancel, close, merge, split, reset, …) and requires approval for ordinary
  writes, derived from the live command tree. Classification is fail-safe: only
  an explicit read allowlist stays allowed, so unrecognized/future verbs are
  gated rather than slipping through. Rules are emitted as exact command paths
  and exact MCP tool names; the hook matches the subcommand at the command
  position (no false positives from verb words in arguments) and covers the
  `canvas api` raw escape hatch (incl. PATCH). `--write` installs the config
  under the project root without overwriting existing files. See the
  [Agent Safety guide](https://jjuanrivvera.github.io/canvas-cli/user-guide/agent-safety/).
- **API coverage raised to 80%** (876 of 1086 documented endpoint patterns,
  method-aware): added ~100 service-layer endpoints across assignments,
  quizzes, modules, outcomes, outcome-imports, courses, sections,
  content-migrations, files, folders, pages, discussions, blueprint, accounts,
  and assignment-groups, each validated by the contract test and covered with
  assertion-rich tests (method + path + request body + parsed response fields).
- **Accurate, method-aware coverage measurement**: the `speccheck` coverage
  harvester now pairs each path with its HTTP verb and resolves context-path
  helper functions, so coverage is counted at (method, path) granularity on
  both sides (previously mismatched-granularity).
- **API spec-compliance harness**: CLI endpoint paths are now validated in CI
  against Canvas's official API spec (Swagger 1.2), committed under
  `testdata/spec/`. A network-free contract test fails the build on any path
  Canvas doesn't document. `make spec-sync` refreshes the manifest from a live
  Canvas host; `make spec-coverage` reports the gap. See
  [API Coverage](https://jjuanrivvera.github.io/canvas-cli/development/api-coverage/).
- **Major API coverage expansion**: 31% → 67% of Canvas's documented endpoints
  (46 → 93 command groups). New command groups include `polls`,
  `appointment-groups`, `folders`, `favorites`, `bookmarks`, `course-nicknames`,
  `observees`, `comm-channels`, `content-shares`, `audit`, `media`,
  `conferences`, `collaborations`, `eportfolios`, `brand`, `jwts`, `progress`,
  `history`, account administration (`auth-providers`, `csp-settings`,
  `account-reports`, `enrollment-terms`, `developer-keys`,
  `grading-period-sets`), grading (`grading-periods`, `grading-standards`,
  `rubric-associations`, `live-assessments`), and content management
  (`content-exports`, `blackout-dates`, `course-pacing`).
- **Full quizzes surface**: reports, statistics, extensions, IP filters,
  question groups, and submission questions.
- **Multi-context support**: discussions, pages, files, folders, and
  content-migrations now work under group and user contexts via `--group-id`
  and `--user-id`, in addition to courses.

### Fixed

- **Response envelope parsing**: many newly added endpoints decoded Canvas's
  named-array envelopes (e.g. `{"polls":[...]}`, `{"events":[...]}`) into bare
  Go types and would fail against a live Canvas; all corrected and now verified
  by shape-asserting tests (polls, audit logs, account calendars, SIS imports,
  quiz submission questions, notification preferences).
- **Request construction**: `csp-settings` domain removal no longer discards
  its argument; `files set-usage-rights` no longer drops all but the last ID;
  group file uploads target `/groups/:id/files` instead of the folder endpoint;
  account/grading param bodies use proper nested JSON instead of unparsed
  bracketed keys; a deliberate `false` for weighted grading periods is no
  longer dropped.
- **Field/return types**: `course-pacing` root-account field and
  `user-features` enabled-list return type corrected.
- Earlier path fixes surfaced by the harness: content-migrations selective
  import (`/selective_data`), last-attended (user-scoped), and quiz IP filters
  (require `:quiz_id`).

### Changed

- CI coverage accounting is HTTP-method-aware (a path isn't counted implemented
  just because the CLI has a different verb on it).
- AGENTS.md documents the spec-compliance + coverage workflow; the bundled
  agent skill's command map covers the expanded surface.

## [1.9.1] - 2026-06-11

### Security

- The OAuth callback server and webhook listener now set a read-header
  timeout (slowloris hardening)
- gosec static analysis is now a blocking CI gate (283-finding backlog
  resolved: real fixes plus justified suppressions)
- Self-update state directory permissions tightened

### Changed

- Releases are built with current GitHub Actions runtimes (Node 24);
  cosign stays on the v2 line so published verification instructions
  keep working
- Dependabot now keeps GitHub Actions and Go dependencies current

### Internal

- New binary-level integration test suite (`make test-integration`)
  exercising the compiled CLI against a mock Canvas server
- New `make check` target running every CI gate locally
- Removed the misleading legacy command test framework; deduplicated
  ~200 test client constructions behind a shared helper

## [1.9.0] - 2026-06-10

### Added

- **Docker image**: releases now publish `ghcr.io/jjuanrivvera/canvas-cli` (distroless, multi-tag)
- **Signed releases**: checksums are signed keylessly with cosign (Sigstore) and archives ship SBOMs
- **AI agent skill**: bundled skill for Claude Code, Cursor, and other agents — install via `canvas skills install`, `npx skills add jjuanrivvera/canvas-cli`, or the Claude Code plugin marketplace
- Confirmation prompt and `--force` flag for `submissions delete-comment`
- `--json-file`/`--csv-file` input flags (the ambiguous `--json`/`--csv` input flags are deprecated but still work)
- `doctor` now honors the global `-o json` output flag

### Fixed

- **Batch assignment sync no longer panics** on courses with assignments
- Retried POST/PUT requests resend the full body instead of an empty one
- `assignments bulk-update` sends `assignment_ids` in a format Canvas accepts
- User assignments endpoint requested the wrong path (user ID was used as course ID)
- `bulk-grade` progress ID is now parsed from the Canvas response
- Ctrl+C now cancels in-flight course validation and sync operations
- Rate-limit bookkeeping no longer races when quota is updated concurrently
- DELETE responses are drained and closed, restoring HTTP connection reuse
- Pagination is guarded against servers that repeat the same `next` link
- 429 responses now honor the `Retry-After` header
- File upload confirmation works for OAuth (auto-refreshed) sessions

### Security

- Static API tokens from `auth token set` are stored in the OS keyring/encrypted store instead of plaintext `config.yaml` (existing configs keep working)
- File downloads sanitize server-supplied filenames (path-traversal protection)
- Self-update aborts when a release is missing `checksums.txt` (fail closed)
- `http://` instance URLs are rejected for non-loopback hosts
- Webhook listener defaults to `127.0.0.1` and warns when signature verification is not configured
- OAuth callback server binds to loopback only
- Config and REPL history directories are created with `0700` permissions

### Changed

- Delete confirmations are unified: all destructive commands accept `y`/`yes` and honor `--dry-run`
- `shell` is now an alias of `repl` (single REPL command)
- Remaining commands migrated to the options-struct pattern (`api`, `cache`, `sync`, `telemetry`, `repl`, `completion`)
- CI pins Go via `go.mod`, blocks on `govulncheck`, and pins security scanners
- A warning is printed when `CANVAS_URL`/`CANVAS_TOKEN` env vars override an explicit `--instance` flag

## [1.8.1] - 2026-06-09

### Changed

- **Test Quality**: Raised overall test coverage to ~82% and added a CI coverage gate (#31)
- **CI Stability**: Stabilized cross-platform CI runs (ubuntu/macos/windows) (#31)

### Documentation

- Improved MCP setup and auth/environment documentation (#29)
- Added CI, pkg.go.dev, codecov, and Ask DeepWiki badges (#30)

## [1.8.0] - 2026-04-29

### Added

- **MCP Server Mode**: Run Canvas CLI as a Model Context Protocol server for AI tools and editors.
  - Added `canvas mcp` command group for server startup, tool export, and editor integration.
  - Added MCP command discovery and schema generation via `ophis`.
- **CLI Runtime Utilities**: Added shared terminal and parsing utilities.
  - Added signal-aware command execution via `ExecuteContext`.
  - Added shell-style alias parsing support with `internal/shellparse`.

### Changed

- **Toolchain Baseline**: Updated minimum Go version to 1.25.0 and aligned CI accordingly.
- **Structured Output Behavior**: Improved command output helpers to better preserve parseable JSON/YAML/CSV in scripts.

### Fixed

- **Canvas Soft Error Detection**: Detect API error bodies returned with HTTP 200 and surface them as errors instead of successful responses.
- **Safety for API Command Flags**: Removed global `-q` shorthand from `--quiet` to avoid collision with `canvas api --query/-q`.

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

## [1.6.1] - 2026-01-19

### Fixed

- Changed command lifecycle logs to DEBUG level to keep normal output clean

## [1.6.0] - 2026-01-19

### Added

- **Command Infrastructure Packages** (#18)
  - `commands/internal/options` package with an option-struct validation framework
  - `commands/internal/logging` package with structured command logging
  - `commands/internal/testing` package for command integration tests

### Changed

- Began migrating commands from package-level flag variables to the options-struct pattern

## [1.5.3] - 2026-01-15

### Added

- **Default Account ID**: Configure a default account so account-scoped commands don't require `--account-id` every time (#17)
- **Global `--limit` Flag**: Limit the number of results for any list operation (#17)

### Fixed

- Assorted Canvas API request fixes (#17)

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

## [1.5.1] - 2026-01-14

### Fixed

- Wrapped command examples in code blocks in the generated CLI reference documentation

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

## [1.3.1] - 2026-01-13

### Fixed

- Corrected ldflags variable names so released binaries report the right version (`main.Version`, `main.Commit`, `main.BuildDate`)

## [1.3.0] - 2026-01-13

### Added

- **GoReleaser Configuration**: Automated multi-platform release builds with checksums (#7)
- **Homebrew Tap**: Install via `brew tap jjuanrivvera/canvas-cli && brew install canvas-cli` (#7)

### Fixed

- Corrected import grouping for the goimports linter
- Removed broken SPECIFICATION.md link from the changelog
- Removed conflicting `goreleaser.yml` config file

## [1.2.1] - 2026-01-13

### Fixed

- Addressed lint and security issues from PR review
- Resolved all testing report findings

## [1.2.0] - 2026-01-11

### Added

- **Input Validation**: Validation for command inputs with helpful error messages
- **Enhanced Linting**: Added errcheck and unconvert linters with proper exclusions
- **Pre-commit Hook**: Added `.githooks/pre-commit` running fmt, vet, lint, and short tests (`make setup-hooks`)
- Implemented missing features and resolved open issues from the initial release

### Documentation

- Clarified branching strategy and release process
- Updated AGENTS.md and CONTRIBUTING.md with hooks info

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
- Inline code documentation with examples

### Infrastructure
- Cross-platform support (macOS, Linux, Windows)
- Cobra CLI framework for command structure
- Viper for configuration management
- Standard Go project layout

---

[Unreleased]: https://github.com/jjuanrivvera/canvas-cli/compare/v1.9.1...HEAD
[1.9.1]: https://github.com/jjuanrivvera/canvas-cli/compare/v1.9.0...v1.9.1
[1.9.0]: https://github.com/jjuanrivvera/canvas-cli/compare/v1.8.1...v1.9.0
[1.8.1]: https://github.com/jjuanrivvera/canvas-cli/compare/v1.8.0...v1.8.1
[1.8.0]: https://github.com/jjuanrivvera/canvas-cli/compare/v1.7.0...v1.8.0
[1.7.0]: https://github.com/jjuanrivvera/canvas-cli/compare/v1.6.1...v1.7.0
[1.6.1]: https://github.com/jjuanrivvera/canvas-cli/compare/v1.6.0...v1.6.1
[1.6.0]: https://github.com/jjuanrivvera/canvas-cli/compare/v1.5.3...v1.6.0
[1.5.3]: https://github.com/jjuanrivvera/canvas-cli/compare/v1.5.2...v1.5.3
[1.5.2]: https://github.com/jjuanrivvera/canvas-cli/compare/v1.5.1...v1.5.2
[1.5.1]: https://github.com/jjuanrivvera/canvas-cli/compare/v1.5.0...v1.5.1
[1.5.0]: https://github.com/jjuanrivvera/canvas-cli/compare/v1.4.0...v1.5.0
[1.4.0]: https://github.com/jjuanrivvera/canvas-cli/compare/v1.3.1...v1.4.0
[1.3.1]: https://github.com/jjuanrivvera/canvas-cli/compare/v1.3.0...v1.3.1
[1.3.0]: https://github.com/jjuanrivvera/canvas-cli/compare/v1.2.1...v1.3.0
[1.2.1]: https://github.com/jjuanrivvera/canvas-cli/compare/v1.2.0...v1.2.1
[1.2.0]: https://github.com/jjuanrivvera/canvas-cli/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/jjuanrivvera/canvas-cli/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/jjuanrivvera/canvas-cli/releases/tag/v1.0.0
