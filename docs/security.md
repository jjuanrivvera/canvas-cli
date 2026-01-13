# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.1.x   | :white_check_mark: |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability, please report it responsibly.

### How to Report

1. **Do NOT** open a public GitHub issue for security vulnerabilities
2. Email security concerns to the maintainers via GitHub's private vulnerability reporting
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
- Tokens are never logged or displayed
- Rate limiting prevents accidental API abuse

## Security Scanning

This project uses automated security tools:

- **gosec**: Static analysis for security issues
- **govulncheck**: Dependency vulnerability scanning
- **Dependabot**: Automated dependency updates

## Dependencies

We regularly update dependencies to patch security vulnerabilities. Run `go mod tidy` to ensure you have the latest versions.
