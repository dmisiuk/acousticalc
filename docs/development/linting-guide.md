# Linting and Code Quality Guide

This guide explains the linting and code quality setup for AcoustiCalc to ensure consistent, high-quality code and prevent formatting issues from reaching PRs.

## Overview

The linting setup includes:
- **Pre-commit hooks** that run automatically before each commit
- **Makefile targets** for manual linting and formatting
- **CI integration** that enforces code quality standards
- **Comprehensive golangci-lint configuration** for thorough code analysis

## Quick Start

### 1. Install Dependencies
```bash
make install-deps
```

### 2. Install Git Hooks
```bash
make setup-hooks
```

### 3. Run Pre-commit Checks
```bash
make pre-commit
```

## Available Commands

### Makefile Commands (Root Level)

```bash
# Linting and formatting
make lint              # Run all linters
make lint-fix          # Fix linting issues automatically
make format            # Format code
make format-check      # Check if code is formatted
make vet               # Run go vet
make staticcheck       # Run staticcheck
make security-scan     # Run security scan

# Development workflow
make dev-check         # Quick development check (lint + unit tests)
make ci-check          # CI-style check (lint + all tests)
make validate          # Quick validation

# Setup
make setup-hooks       # Install git hooks
make install-deps      # Install development dependencies
```

### Detailed Commands (Tests Scripts)

```bash
cd tests/scripts
make lint             # Comprehensive linting
make pre-commit       # Pre-commit checks
make format-check     # Format validation
make vet              # Go vet analysis
make staticcheck      # Static analysis
make security-scan    # Security scanning
```

## Pre-commit Hooks

The pre-commit hook automatically runs before each commit and checks:
1. Code formatting (`make format-check`)
2. Go vet analysis (`make vet`)
3. Static analysis (`make staticcheck`)
4. Golangci-lint (if installed)
5. Unit tests

If any check fails, the commit is prevented until issues are fixed.

## CI Integration

The CI workflow now includes linting checks:
- **Format check**: Ensures code is properly formatted
- **Go vet**: Catches programming errors
- **Security scan**: Checks for vulnerabilities
- **Unit tests**: Ensures code functionality

## Configuration Files

### `.golangci.yml`
Comprehensive golangci-lint configuration including:
- Error handling checks
- Security scanning (gosec)
- Performance checks
- Code complexity analysis
- Import organization
- Naming conventions

### `Makefile` (Root)
Main development commands and workflow targets.

### `tests/scripts/Makefile`
Detailed Unix-specific testing and linting operations.

## Development Workflow

### 1. Before Committing
```bash
# Option 1: Let the pre-commit hook handle it
git commit -m "Your commit message"

# Option 2: Run checks manually
make pre-commit
git commit -m "Your commit message"
```

### 2. Before Pushing
```bash
# Full CI-style check
make ci-check

# Or push and let CI handle it
git push origin your-branch
```

### 3. Fixing Issues
```bash
# Auto-fix formatting issues
make lint-fix

# Manual formatting
make format

# Check specific linters
make vet
make staticcheck
```

## Linters Included

### Built-in Go Tools
- **go fmt**: Code formatting
- **go vet**: Static analysis
- **goimports**: Import organization

### Golangci-lint Suite
- **gosec**: Security scanning
- **staticcheck**: Static analysis
- **gocyclo**: Complexity analysis
- **nestif**: Nesting depth checks
- **prealloc**: Performance optimization
- **revive**: Code style enforcement
- **errorlint**: Error handling best practices

## Quality Gates

The following must pass for code to be committed:
1. ✅ Code formatting (go fmt)
2. ✅ Go vet analysis
3. ✅ Unit tests
4. ✅ No security vulnerabilities (govulncheck)
5. ✅ Static analysis (staticcheck, golangci-lint)

## Troubleshooting

### Common Issues

**Pre-commit hook fails:**
```bash
# Check what failed
make pre-commit

# Auto-fix what you can
make lint-fix

# Commit again
git commit -m "Your commit message"
```

**Formatting issues:**
```bash
# Auto-format code
make format

# Check formatting
make format-check
```

**Linting errors:**
```bash
# Run specific linter
make vet
make staticcheck

# Or run all linters
make lint
```

### Bypassing Hooks (Not Recommended)

If you absolutely need to bypass pre-commit checks:
```bash
git commit --no-verify -m "Your commit message"
```

**Note**: This is not recommended as it allows poor quality code to be committed.

## Best Practices

1. **Run `make dev-check` frequently** during development
2. **Fix issues as they appear** rather than letting them accumulate
3. **Use the pre-commit hooks** as your first line of defense
4. **Read the linting output** to understand the issues
5. **Configure your editor** to run formatting on save

## Editor Integration

### VS Code
Install the Go extension and add to your settings.json:
```json
{
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "file",
    "go.formatOnSave": true,
    "editor.formatOnSave": true
}
```

### Vim/Neovim
Add to your vimrc:
```vim
autocmd BufWritePre *.go :GoImports
autocmd BufWritePre *.go :GoFmt
```

## Contributing

All contributors must ensure their code passes all linting checks before submitting PRs. The CI system will automatically reject code that doesn't meet quality standards.

## Regular Updates

The linting configuration should be reviewed and updated:
- When Go version changes
- When new best practices emerge
- When security vulnerabilities are discovered
- When team coding standards evolve

For more information about specific linters, see their respective documentation:
- [golangci-lint](https://golangci-lint.run/)
- [staticcheck](https://staticcheck.io/)
- [gosec](https://github.com/securego/gosec)