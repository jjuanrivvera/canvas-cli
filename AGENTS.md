# AGENTS.md

This file provides guidance to AI agents (Claude Code, Cursor, Copilot, etc.) when working with code in this repository.

## Build & Development Commands

```bash
# Build
make build              # Build binary to bin/canvas
make dev                # Build with fmt and vet

# Test
make test               # Run all tests
make test-coverage      # Run tests with coverage report
go test -v ./internal/api/...  # Run specific package tests
go test -v -run TestName ./... # Run single test

# Lint & Format
make fmt                # Format code
make lint               # Run golangci-lint
make vet                # Run go vet

# Install
make install            # Install to /usr/local/bin
make uninstall          # Remove from /usr/local/bin

# Setup
make setup-hooks        # Install git pre-commit hooks
```

## Pre-commit Hook

Run `make setup-hooks` to enable. Runs automatically on each commit:
- `gofmt` - formatting check
- `golangci-lint` - comprehensive linting (if installed)
- `go vet` - static analysis
- `go test -short` - quick test pass

## Architecture

Canvas CLI is a Go CLI for Canvas LMS API, built with Cobra/Viper.

### Project Structure

```
cmd/canvas/     → Entry point (main.go)
commands/       → Cobra command definitions (one file per resource)
internal/
  api/          → Canvas API client + service layer (Client, *Service structs)
  auth/         → OAuth 2.0 + PKCE, token storage (keyring/encrypted file)
  config/       → Viper-based configuration management
  cache/        → Response caching with TTL
  batch/        → Concurrent batch operations (worker pool)
  output/       → Formatters (table, JSON, YAML, CSV)
  repl/         → Interactive shell
.ai/            → Canvas LMS API documentation (gitignored)
```

### Key Patterns

**Service Layer**: Each Canvas resource has a service in `internal/api/`:
```go
type ModulesService struct { client *Client }
func NewModulesService(client *Client) *ModulesService
```

**Command Pattern**: Commands in `commands/` follow this structure:
```go
var resourceCmd = &cobra.Command{...}
func init() { rootCmd.AddCommand(resourceCmd) }
```

**API Client**: `internal/api/client.go` provides `HTTPClient` interface with:
- Automatic pagination (`GetAllPages`)
- Adaptive rate limiting based on Canvas quota headers
- Exponential backoff retry

### Testing

Tests use `httptest.NewServer` for mock HTTP servers. Service tests follow pattern:
```go
func TestServiceMethod(t *testing.T) {
    server := httptest.NewServer(...)
    client := &Client{BaseURL: server.URL, ...}
    service := NewXxxService(client)
    // test service methods
}
```

Use `t.Fatal()` (not `t.Error()`) when nil checks would cause subsequent panics.

## Branching & Release

### Branch Model (Simplified Git Flow)

```
main     ──●─────────────────●──────► (tagged releases)
           │                 ↑↓
develop  ──●───●───●───●─────●──────► (integration)
               ↑       ↑
feature/*  ────┘       │
fix/*  ────────────────┘
```

| Branch | Purpose | Merges To |
|--------|---------|-----------|
| `main` | Tagged releases only | - |
| `develop` | Integration (PR target) | `main` on release |
| `feature/*` | New features | `develop` |
| `fix/*` | Bug fixes | `develop` |
| `hotfix/*` | Urgent fixes | `main` AND `develop` |

### When develop syncs with main

1. **After a release**: Merge `main` back to `develop` to capture release commits
2. **After a hotfix**: Hotfix merges to both `main` and `develop`

### Release Process

```bash
# 1. Merge develop to main
git checkout main && git merge develop

# 2. Tag and push
git tag -a v1.x.x -m "Release v1.x.x"
git push origin main --tags

# 3. Sync main back to develop
git checkout develop && git merge main
git push origin develop
```

GitHub Actions automatically builds binaries and creates the release on tag push.

## CI

Single workflow `.github/workflows/ci.yml` runs:
- Lint (gofmt, go vet, golangci-lint, staticcheck)
- Security (govulncheck, gosec)
- Test matrix (ubuntu/macos/windows × Go 1.21/1.22)
- Build artifacts

## Documentation

Documentation is built with MkDocs Material and deployed to GitHub Pages.

**Live site**: https://jjuanrivvera.github.io/canvas-cli/

### Local Development

```bash
# Install dependencies
pip install mkdocs-material mkdocs-git-revision-date-localized-plugin

# Generate CLI reference and serve locally
go run ./tools/gendocs/main.go
mkdocs serve
```

### Deployment

Documentation auto-deploys on push to `main` when `docs/**` or `mkdocs.yml` changes.

**Manual trigger** (via GitHub UI):
1. Go to Actions → Documentation workflow
2. Click "Run workflow"

**Manual trigger** (via CLI):
```bash
gh workflow run docs.yml
```

**If deployment gets stuck** in "queued" status:
```bash
# Force a Pages build
gh api -X POST repos/jjuanrivvera/canvas-cli/pages/builds

# Check status
gh api repos/jjuanrivvera/canvas-cli/pages --jq '.status'
```

## Releases

Releases use GoReleaser and auto-publish to GitHub Releases + Homebrew tap.

### Creating a Release

```bash
# 1. Ensure main is up to date
git checkout main && git merge develop

# 2. Create and push tag
git tag -a v1.x.x -m "Release v1.x.x"
git push origin main --tags

# 3. Sync develop
git checkout develop && git merge main
git push origin develop
```

GoReleaser automatically:
- Builds binaries for linux/darwin/windows (amd64/arm64)
- Creates GitHub release with changelog
- Updates Homebrew formula in `jjuanrivvera/homebrew-canvas-cli`

### Homebrew Tap

The formula is at: https://github.com/jjuanrivvera/homebrew-canvas-cli

**Required secret**: `HOMEBREW_TAP_TOKEN` - a PAT with `repo` scope for the tap repository
