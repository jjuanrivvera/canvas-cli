# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

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
gofmt -l .              # Check formatting (CI uses this)

# Install
make install            # Install to /usr/local/bin
make uninstall          # Remove from /usr/local/bin

# Setup
make setup-hooks        # Install git pre-commit hooks
```

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

## CI

Single workflow `.github/workflows/ci.yml` runs:
- Lint (gofmt, go vet, golangci-lint, staticcheck)
- Security (govulncheck, gosec)
- Test matrix (ubuntu/macos/windows × Go 1.21/1.22)
- Build artifacts
