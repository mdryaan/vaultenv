# Contributing to vaultenv

Thank you for your interest in contributing to vaultenv! This document outlines the process for contributing to the project.

## Prerequisites

- Go 1.22 or later
- Git
- `golangci-lint` (for linting)
- A terminal with a Unix-like environment (Linux, macOS, WSL2)

## Getting Started

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork:
   ```bash
   git clone https://github.com/<your-username>/vaultenv.git
   cd vaultenv
   ```
3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/mdryaan/vaultenv.git
   ```

### Create a Branch

Always work from a feature branch, never directly from `main`:

```bash
git checkout -b feat/your-feature-name
```

**Branch naming conventions:**

| Prefix | Purpose |
|--------|---------|
| `feat/` | New feature |
| `fix/` | Bug fix |
| `chore/` | Maintenance, deps, tooling |
| `test/` | Test additions or fixes |
| `docs/` | Documentation only |
| `refactor/` | Code restructuring without behavior change |

## Commit Message Format

Follow the Conventional Commits specification:

```
<type>(<scope>): <description>

[optional body]
```

**Types:** `feat`, `fix`, `docs`, `chore`, `test`, `refactor`, `perf`

**Examples:**

```
feat(vault): add rotate command for password rotation

fix(crypto): handle empty plaintext in AES encrypt

test(dotenv): add table-driven parser tests

docs: update README with CI/CD usage section
```

## Security Contributions

If you discover a security vulnerability, **do not open a public issue**. Instead:

1. Email the maintainer privately with details
2. Allow time for a patch to be prepared and released before public disclosure
3. We follow responsible disclosure practices

For security-related code changes:
- All crypto code must have corresponding test vectors
- Password bytes must be zeroed after use
- Never log or print secret values
- Sensitive data must never be written to temporary files without proper cleanup

## Running Tests

```bash
make test
```

Or run a specific package:

```bash
go test ./internal/crypto/... -v
```

Check coverage:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Manual Testing Steps

Before submitting a PR, verify the following manually:

```bash
# Build
make build

# Initialize a vault
./bin/vaultenv init

# Set secrets
./bin/vaultenv set DATABASE_URL postgres://localhost/myapp
./bin/vaultenv set API_KEY --generate
./bin/vaultenv set REDIS_URL redis://localhost --tags production,backend

# Retrieve secrets
./bin/vaultenv get DATABASE_URL
./bin/vaultenv get API_KEY --mask

# List secrets
./bin/vaultenv list
./bin/vaultenv list --show-values
./bin/vaultenv list --tags production
./bin/vaultenv list --output json

# Export/import
./bin/vaultenv export --output .env.test
./bin/vaultenv import .env.test --dry-run

# Rotate password
./bin/vaultenv rotate

# Shell completion
./bin/vaultenv completion bash
```

## Linting

```bash
make lint
```

Or format code:

```bash
make fmt
```

## Pull Request Checklist

Before submitting your PR, ensure:

- [ ] All tests pass: `make test`
- [ ] Code is formatted: `make fmt`
- [ ] No lint errors: `make lint`
- [ ] New features include tests
- [ ] Crypto changes include test vectors
- [ ] No secret values appear in logs or error messages
- [ ] Password byte slices are zeroed after use
- [ ] Commit messages follow Conventional Commits format
- [ ] PR description explains the motivation and approach

## Code Style

- Follow standard Go idioms and conventions
- Prefer table-driven tests
- Use `require` for fatal assertions, `assert` for non-fatal in tests
- Keep functions focused and short
- Error messages should be lowercase and not end with punctuation
- Wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
