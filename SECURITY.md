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

## Security Best Practices

### Token Storage

Canvas CLI stores authentication tokens securely:

- **macOS**: Keychain (preferred)
- **Linux**: Secret Service API or encrypted file
- **Windows**: Windows Credential Manager or encrypted file

### Configuration Security

- Never commit `.canvas-cli.yaml` or any file containing tokens
- The CLI automatically adds sensitive files to `.gitignore`
- Use environment variables (`CANVAS_TOKEN`) for CI/CD environments

### API Security

- All API communication uses HTTPS
- Tokens are never logged. They are redacted in `--dry-run` output unless you explicitly opt in with the `--show-token` flag
- Rate limiting prevents accidental API abuse

## Security Scanning

This project uses automated security tools:

- **gosec**: Static analysis for security issues (blocking in CI)
- **govulncheck**: Dependency vulnerability scanning (blocking in CI)
- **CodeQL**: Go code scanning, published to GitHub code scanning
- **OpenSSF Scorecard**: supply-chain checks, published weekly from the default branch
- **Dependabot**: Automated dependency updates for Go modules, GitHub Actions, and the landing site's npm devDependencies

## Dependencies

We regularly update dependencies to patch security vulnerabilities. Run `go mod tidy` to ensure you have the latest versions.
