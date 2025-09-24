# AcoustiCalc Testing Framework Quick Start

This guide provides a quick introduction to using the AcoustiCalc testing framework.

## Prerequisites

### System Requirements
- **Operating System**: Unix-like system (Linux, macOS)
- **Go**: Version 1.21 or higher
- **Required Tools**: git, make, bc, find, xargs, mktemp, timeout

### Quick Setup
```bash
# Clone the repository
git clone <repository-url>
cd acousticalc

# Validate your environment
./tests/scripts/unix_env_validate.sh

# Install Go dependencies
go mod download
go mod tidy
```

## Quick Test Execution

### Run All Tests
```bash
# Using Unix tools script
./tests/scripts/unix_test_tools.sh run

# Using Makefile
make test

# Quick development cycle
make dev-test
```

### Run Specific Test Types
```bash
# Unit tests only
make test-unit

# Integration tests only
make test-integration

# Benchmarks only
make test-benchmark

# Coverage reports only
make test-coverage
```

## Development Workflow

### 1. Environment Validation
```bash
# Check if your environment is properly configured
make validate

# Show current configuration
make show-config
```

### 2. Watch Mode for Development
```bash
# Watch file changes and run tests automatically
make test-watch

# Or use the Unix tools script
./tests/scripts/unix_test_tools.sh watch
```

### 3. Quick Test Cycles
```bash
# Quick unit test run
make quick

# Quick test with coverage
make quick-coverage

# Quick benchmark run
make quick-bench
```

## Test Organization

### Understanding the Test Structure
```
tests/
├── unit/           # Fast, isolated unit tests
├── integration/    # Component interaction tests
├── e2e/           # End-to-end tests (future)
└── artifacts/      # Test reports and coverage
```

### Running Tests by Category
```bash
# Run unit tests with coverage
cd tests/unit && go test -v -cover

# Run integration tests
cd tests/integration && go test -v

# Run tests in parallel
cd tests/unit && go test -parallel=4
```

## Coverage and Reports

### Generate Coverage Reports
```bash
# HTML coverage report
make coverage-html

# Coverage summary
make coverage-summary

# Combined coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### View Coverage Reports
```bash
# Open HTML coverage report
open tests/artifacts/coverage/coverage.html

# View coverage summary
cat tests/artifacts/coverage/coverage_summary.txt
```

## Performance Testing

### Run Benchmarks
```bash
# Basic benchmarks
make test-benchmark

# Detailed benchmarks with analysis
make benchmark-detailed

# Run benchmarks with custom settings
cd tests/unit && go test -bench=. -benchmem -benchtime=5s
```

### Performance Monitoring
```bash
# Monitor performance during tests
make monitor-performance

# Check system resources
make show-config
```

## Writing Tests

### Unit Test Example
```go
func TestCalculatorAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     float64
        expected float64
    }{
        {"positive", 2, 3, 5},
        {"negative", -2, -3, -5},
        {"mixed", -2, 3, 1},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### Integration Test Example
```go
func TestComponentIntegration(t *testing.T) {
    // Setup
    mockCalc := NewMockCalculator()
    dataProvider := NewTestDataProvider()

    // Test
    result, err := mockCalc.Evaluate("2 + 3")
    if err != nil {
        t.Errorf("Integration test failed: %v", err)
    }
    if result != 5 {
        t.Errorf("Expected 5, got %v", result)
    }
}
```

## Debugging Tests

### Verbose Output
```bash
# Run tests with verbose output
make test-unit VERBOSE=true

# Or use the Unix tools script
./tests/scripts/unix_test_tools.sh run --verbose
```

### Debug Coverage
```bash
# Debug coverage generation
make debug-coverage

# Check coverage files
ls -la tests/artifacts/coverage/
```

### Debug Performance
```bash
# Run tests with debug output
make debug-test

# Monitor performance
make monitor-performance
```

## CI/CD Integration

### GitHub Actions
The testing framework is integrated with GitHub Actions and automatically:
- Runs tests on Linux, macOS, and Windows
- Generates coverage reports (Unix only)
- Runs benchmarks (Linux only)
- Uploads test artifacts
- Checks coverage thresholds

### Quality Gates
- All tests must pass on all platforms
- Coverage must be ≥80%
- No security vulnerabilities
- Code must be properly formatted

## Common Issues and Solutions

### Environment Issues
```bash
# Validate environment
make validate

# Check dependencies
make check-deps

# Setup completion
make setup-completion
```

### Test Failures
```bash
# Run tests with verbose output
make debug-test

# Check for race conditions
go test -race ./...

# Run specific test
go test -run TestSpecificFunction ./...
```

### Coverage Issues
```bash
# Check coverage generation
make debug-coverage

# Generate coverage manually
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## Advanced Usage

### Custom Test Execution
```bash
# Run tests with custom parallelism
PARALLEL_JOBS=8 make test

# Disable coverage
COVERAGE_ENABLED=false make test

# Run tests with custom timeout
go test -timeout=60s ./...
```

### Performance Tuning
```bash
# Optimize for your system
export GOMAXPROCS=$(nproc)
export PARALLEL_JOBS=$(nproc)

# Apply Unix optimizations
make unix-optimize
```

## Getting Help

### Built-in Help
```bash
# Unix tools help
./tests/scripts/unix_test_tools.sh help

# Makefile help
make help

# Show configuration
make show-config
```

### Documentation
- **Full Standards**: [Testing Standards](./testing-standards.md)
- **Architecture**: [Architecture Documentation](./architecture/)
- **Go Documentation**: [Go Testing](https://golang.org/pkg/testing/)

### Troubleshooting
1. Check environment: `make validate`
2. Check dependencies: `make check-deps`
3. Run with verbose output: `make debug-test`
4. Check logs in `tests/artifacts/`

## Next Steps

1. **Run the validation**: `make validate`
2. **Execute all tests**: `make test`
3. **Set up watch mode**: `make test-watch`
4. **Explore the documentation**: Read the full [Testing Standards](./testing-standards.md)

---

**This quick start guide covers the essentials. For detailed information, refer to the full Testing Standards document.**