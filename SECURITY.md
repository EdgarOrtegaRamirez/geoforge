# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |

## Reporting a Vulnerability

If you discover a security vulnerability, please report it responsibly:

1. **Do NOT** open a public GitHub issue
2. Email the maintainer directly
3. Include a description of the vulnerability
4. Include steps to reproduce if possible

## Security Considerations

- GeoForge processes untrusted GeoJSON input — all parsing is done with care
- No network operations (pure file I/O)
- No command injection vectors
- No eval or dynamic code execution
- All file operations use safe path handling
