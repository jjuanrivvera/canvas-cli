# Security Policy

## Supported Versions

Only the latest minor release line receives security fixes.

| Version | Supported          |
| ------- | ------------------ |
| 1.15.x  | :white_check_mark: |
| < 1.15  | :x:                |

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability, please report it responsibly.

### How to Report

1. **Do NOT** open a public GitHub issue for security vulnerabilities
2. Report it privately through [GitHub Security Advisories](https://github.com/jjuanrivvera/canvas-cli/security/advisories/new) (private vulnerability reporting)
3. Include as much detail as possible:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

### What to Expect

- **Acknowledgment**: We will acknowledge receipt within 48 hours
- **Assessment**: We will assess the vulnerability and determine its severity
- **Fix Timeline**: Critical vulnerabilities will be addressed within 7 days
- **Disclosure**: We will coordinate with you on public disclosure timing

## Supply-Chain Security

Since v1.9.0 every release is verifiable end to end:

- **Signed checksums**: `checksums.txt` is signed keylessly with
  [cosign](https://github.com/sigstore/cosign) (Sigstore) by the GitHub
  Actions release workflow — see
  [verification instructions](getting-started/installation.md#verifying-a-downloaded-release)
- **SBOMs**: every archive ships a Software Bill of Materials (`*.sbom.json`)
- **Fail-closed self-update**: `canvas update` refuses to install a binary
  whose checksums cannot be downloaded and verified
- **Distroless Docker image**: `ghcr.io/jjuanrivvera/canvas-cli` is built from
  `distroless/static` and runs as a non-root user with no shell inside

## Security Best Practices

### Token Storage

Canvas CLI stores authentication tokens securely — both OAuth tokens and
static API tokens set via `canvas auth token set`:

- **macOS**: Keychain (preferred)
- **Linux**: Secret Service API or encrypted file
- **Windows**: Windows Credential Manager or encrypted file

Tokens written to `config.yaml` by versions before 1.9.0 keep working;
re-run `canvas auth token set` to migrate them into secure storage.

### Network Security

- All API communication uses HTTPS — `http://` instance URLs are rejected
  for anything but loopback hosts
- The OAuth callback server binds to `127.0.0.1` only
- The webhook listener defaults to `127.0.0.1:8080` and warns loudly when
  started without signature verification configured

### Configuration Security

- Never commit `config.yaml` or any file containing tokens
- Configuration and history directories are created with `0700` permissions
- Use environment variables (`CANVAS_URL`/`CANVAS_TOKEN`) for CI/CD environments

### API Security

- Tokens are never logged. They are redacted in `--dry-run` output unless you
  explicitly opt in with the `--show-token` flag
- Adaptive rate limiting prevents accidental API abuse

## Security Scanning

This project uses automated security tools in CI:

- **govulncheck**: Vulnerability scanning (blocking — a reachable
  vulnerability fails the build)
- **gosec**: Static analysis for security issues (blocking in CI)
- **CodeQL**: Go code scanning, published to GitHub code scanning
- **OpenSSF Scorecard**: supply-chain checks, published weekly from the default branch
- **Dependabot**: Automated dependency updates for Go modules, GitHub Actions, and the landing site's npm devDependencies

## Dependencies

We regularly update dependencies to patch security vulnerabilities. Run `go mod tidy` to ensure you have the latest versions.
