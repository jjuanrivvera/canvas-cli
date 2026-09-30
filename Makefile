.PHONY: help verify cover-check build test test-integration clean install uninstall fmt lint vet run deps setup-hooks docs-gen docs-serve docs-build docs-deploy spec-sync spec-coverage

# Variables
BINARY_NAME=canvas
VERSION?=dev
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE)"

# The coverage floor. It is the same number as .github/workflows/ci.yml, and it
# is a ratchet: raise it as coverage grows, never lower it.
COVER_MIN?=80

# Default target
help:
	@echo "Canvas CLI - Makefile targets:"
	@echo ""
	@echo "  make build        - Build the binary"
	@echo "  make install      - Install the binary to /usr/local/bin"
	@echo "  make uninstall    - Remove the binary from /usr/local/bin"
	@echo "  make verify       - The gate: everything check runs, plus the coverage floor"
	@echo "  make check        - Run everything CI runs: fmt-check, vet, lint, gosec, tests, integration"
	@echo "  make test         - Run tests"
	@echo "  make test-integration - Run binary-level integration tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Run linter"
	@echo "  make vet          - Run go vet"
	@echo "  make run          - Build and run the CLI"
	@echo "  make deps         - Download dependencies"
	@echo "  make release      - Build snapshot release for all platforms (GoReleaser)"
	@echo "  make setup-hooks  - Install git pre-commit hooks"
	@echo ""
	@echo "Documentation:"
	@echo "  make docs-gen     - Generate CLI reference documentation"
	@echo "  make docs-serve   - Serve documentation locally"
	@echo "  make docs-build   - Build documentation site"
	@echo "  make docs-deploy  - Deploy documentation to GitHub Pages"
	@echo ""
	@echo "Spec compliance:"
	@echo "  make spec-sync     - Fetch official Canvas Swagger and regenerate testdata/spec/canvas_endpoints.json"
	@echo "  make spec-sync CANVAS_SPEC_HOST=https://myschool.instructure.com  - Use a specific Canvas host"
	@echo "  make spec-coverage - Print official Canvas API coverage gap report (network-free)"
	@echo ""

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/canvas
	@echo "✓ Build complete: bin/$(BINARY_NAME)"

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp bin/$(BINARY_NAME) /usr/local/bin/
	@echo "✓ Installed: /usr/local/bin/$(BINARY_NAME)"

# Uninstall the binary
uninstall:
	@echo "Removing $(BINARY_NAME) from /usr/local/bin..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✓ Uninstalled"

# The fleet gate. A change is done when `make verify` exits 0: everything CI
# enforces, including the coverage floor, which `check` alone does not measure.
verify: check cover-check
	@echo "✓ verify"

cover-check:
	@echo "Measuring coverage..."
	@go test -coverprofile=coverage.out ./... > /dev/null
	@total=$$(go tool cover -func=coverage.out | awk '/^total:/ {print $$3}' | tr -d '%'); \
	awk -v t="$$total" -v min="$(COVER_MIN)" 'BEGIN { if (t+0 < min+0) { printf "✗ coverage %.1f%% < %s%%\n", t, min; exit 1 } printf "✓ coverage %.1f%% ≥ %s%%\n", t, min }'

# Run everything CI runs, locally
check: vet lint
	@echo "Checking formatting..."
	@unformatted=$$(gofmt -l . | grep -v '^.claude' || true); \
	if [ -n "$$unformatted" ]; then echo "✗ Needs gofmt:"; echo "$$unformatted"; exit 1; fi
	@echo "✓ Formatting clean"
	@if command -v gosec > /dev/null; then \
		echo "Running gosec..."; gosec -quiet ./... && echo "✓ Gosec clean"; \
	else echo "⚠ gosec not installed — skipping (CI enforces it)"; fi
	@echo "Running tests (-race)..."
	@go test -race ./...
	@echo "✓ Tests pass"
	@$(MAKE) --no-print-directory test-integration
	@echo ""
	@echo "✓ All checks passed"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run binary-level integration tests (requires compiled binary; skipped by default)
test-integration:
	@echo "Running integration tests..."
	@go test -tags integration -v -timeout 5m ./test/integration/

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "✓ Coverage report: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf dist/
	@rm -f coverage.txt coverage.html
	@echo "✓ Clean complete"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Format complete"

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install: https://golangci-lint.run/usage/install/" && exit 1)
	@golangci-lint run
	@echo "✓ Lint complete"

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "✓ Vet complete"

# Build and run
run: build
	@./bin/$(BINARY_NAME)

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies downloaded"

# Snapshot release build via GoReleaser (same pipeline as tagged releases:
# ldflags, archives, checksums) without publishing anything.
release:
	@which goreleaser > /dev/null || (echo "goreleaser not installed. Install: https://goreleaser.com/install/" && exit 1)
	@echo "Building snapshot release with GoReleaser..."
	@goreleaser release --snapshot --clean --skip=sign,sbom,docker
	@echo "✓ Snapshot release built in dist/"

# Development build with verbose output
dev: fmt vet
	@go build -v $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/canvas

# Setup git hooks
setup-hooks:
	@echo "Setting up git hooks..."
	@git config core.hooksPath .githooks
	@echo "✓ Git hooks installed (.githooks/pre-commit)"

# Documentation targets
docs-gen:
	@echo "Generating CLI reference documentation..."
	@go run ./tools/gendocs/main.go
	@cp CHANGELOG.md docs/changelog.md
	@echo "✓ CLI docs generated in docs/commands/ (changelog synced)"

docs-serve: docs-gen
	@echo "Serving documentation at http://localhost:8000..."
	@mkdocs serve

docs-build: docs-gen
	@echo "Building documentation site..."
	@mkdocs build --strict
	@echo "✓ Documentation built in site/"

docs-deploy: docs-gen
	@echo "Deploying documentation to GitHub Pages..."
	@mkdocs gh-deploy --force
	@echo "✓ Documentation deployed"

# Spec compliance targets.
# Source of truth: official Canvas Swagger 1.2 API (not .ai/canvas-lms-docs).
# The default host is https://learn.canvas.net; override with -host or $CANVAS_SPEC_HOST.
# canvas.instructure.com returns 503 from datacenter IPs — use a real institutional host.
spec-sync:
	@echo "Fetching official Canvas Swagger and regenerating spec manifest..."
	@go run ./tools/speccheck -sync
	@echo "✓ Spec manifest updated: testdata/spec/canvas_endpoints.json"

spec-coverage:
	@echo "Running spec coverage analysis (network-free; reads committed manifest)..."
	@go run ./tools/speccheck -coverage
