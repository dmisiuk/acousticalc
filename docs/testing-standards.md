# AcoustiCalc Testing Standards and Guidelines

## Overview

This document provides comprehensive testing standards and guidelines for the AcoustiCalc project. It defines the testing philosophy, organization structure, tools, and best practices for maintaining high-quality test coverage and reliable software.

## Testing Philosophy

### Core Principles

1. **Test-Driven Development (TDD)**: Write tests before or alongside production code
2. **Cross-Platform Approach**: Full platform equality for Windows, macOS, and Linux with comprehensive testing coverage
3. **Comprehensive Coverage**: Maintain >80% test coverage with meaningful tests
4. **Performance Awareness**: Include performance benchmarks for critical operations
5. **Integration Focus**: Test component interactions, not just isolated units
6. **Automation First**: Automate all testing processes through CI/CD pipelines

### Testing Pyramid

```
       E2E Tests (5%)
      /             \
   Integration Tests (25%)
  /                   \
Unit Tests (70%)
```

- **Unit Tests (70%)**: Fast, isolated tests of individual functions and methods
- **Integration Tests (25%)**: Tests of component interactions and integration points
- **E2E Tests (5%)**: Full workflow tests that simulate real user scenarios from start to finish across all supported platforms.

## Test Organization Structure

### Directory Layout

```
tests/
├── unit/                          # Unit tests for individual components
│   ├── calculator_test.go         # Core calculator functionality
│   ├── coverage_enhanced_test.go # Enhanced coverage reporting
│   └── benchmark_test.go         # Performance benchmarks
├── integration/                   # Integration and component interaction tests
│   ├── component_interaction_test.go # Component interaction scenarios
│   ├── integration_fixtures.go   # Test fixtures and mock objects
│   └── isolation_test.go         # Test isolation and cleanup procedures
├── e2e/                          # End-to-end tests simulating user workflows
│   ├── workflow_test.go          # Core user journey tests
│   ├── platform_test.go          # Platform-specific E2E tests
│   └── e2e_utils.go              # E2E test utilities and orchestration
├── recording/                    # Terminal recording utilities and tests
│   ├── capture_utils.go          # Core recording functionality
│   └── capture_test.go           # Tests for recording utilities
├── cross_platform/               # Cross-platform compatibility tests
│   └── platform_compat_test.go   # Tests for platform-specific behavior
├── reporting/                    # Test reporting generation and utilities
│   ├── reporting_utils.go        # Report generation logic
│   └── aggregator_test.go        # Tests for report data aggregation
├── artifacts/                     # Generated test evidence and reports
│   ├── coverage/                  # Coverage reports and historical data
│   ├── e2e/                       # E2E test recordings and reports
│   ├── reports/                   # Other test execution reports
│   └── platform_results/          # Platform-specific test results
└── scripts/                       # Test execution utilities
    ├── unix_test_tools.sh         # Unix-specific testing tools
    ├── unix_env_validate.sh       # Environment validation
    └── Makefile                   # Unix Makefile for test operations
```

### Test Naming Conventions

#### File Naming
- **Test Files**: `{component}_test.go`
- **Integration Tests**: `{integration}_test.go`
- **Benchmark Tests**: `{component}_bench_test.go`
- **Coverage Tests**: `{component}_coverage_test.go`

#### Test Function Naming
- **Unit Tests**: `Test{FunctionName}` or `Test{Component}_{Scenario}`
- **Integration Tests**: `Test{Integration}_{Scenario}`
- **Benchmark Tests**: `Benchmark{Operation}`
- **Table Tests**: Use descriptive subtest names

#### Example Test Structure
```go
func TestCalculatorAddition(t *testing.T) {
    tests := []struct {
        name     string
        a, b     float64
        expected float64
    }{
        {"positive numbers", 2.0, 3.0, 5.0},
        {"negative numbers", -2.0, -3.0, -5.0},
        {"mixed signs", -2.0, 3.0, 1.0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%v, %v) = %v, want %v",
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### End-to-End (E2E) Tests

**Purpose**: Validate complete user workflows from start to finish across all supported platforms.

**Characteristics**:
- Slower execution time, as they involve building and running the actual application binary.
- Test the application as a whole, including its interaction with the operating system.
- Provide terminal recordings as artifacts for visual validation and debugging.

**Best Practices**:
```go
// Good E2E test
func TestE2EWorkflow(t *testing.T) {
	t.Run("TestSimpleAddition", func(t *testing.T) {
		run := NewE2ETestRun(t)

		expression := "2 + 2"
		expected := "4"

		// Record the command
		run.RecordCommand(expression)

		// Run the command for assertion
		output := run.RunCommand(expression)

		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s', but got '%s'", expected, output)
		}
	})
}
```

## Test Types and Categories

### Unit Tests

**Purpose**: Test individual functions and methods in isolation

**Characteristics**:
- Fast execution (< 1 second per test)
- No external dependencies
- Mock external services and dependencies
- Test all code paths and edge cases

**Best Practices**:
```go
// Good unit test
func TestEvaluateBasicExpression(t *testing.T) {
    result, err := Evaluate("2 + 3")
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    if result != 5.0 {
        t.Errorf("Expected 5.0, got %v", result)
    }
}

// Test with table-driven approach
func TestEvaluateVariousExpressions(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected float64
        wantErr  bool
    }{
        {"simple addition", "2 + 3", 5.0, false},
        {"division by zero", "10 / 0", 0, true},
        {"decimal numbers", "3.5 + 2.1", 5.6, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Evaluate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && result != tt.expected {
                t.Errorf("Evaluate() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Integration Tests

**Purpose**: Test interactions between components and integration points

**Characteristics**:
- Medium execution time (1-10 seconds)
- Test component interactions
- Use real implementations where possible
- Mock external services only

**Best Practices**:
```go
func TestCalculatorMockIntegration(t *testing.T) {
    // Setup
    mockCalc := NewMockCalculator()
    dataProvider := NewTestDataProvider()

    // Configure mock
    mockCalc.SetResult("2 + 3", 5.0)

    // Test integration
    result, err := mockCalc.Evaluate("2 + 3")
    if err != nil {
        t.Errorf("Mock calculator failed: %v", err)
    }
    if result != 5.0 {
        t.Errorf("Expected 5.0, got %v", result)
    }

    // Test data provider integration
    validExprs := dataProvider.GetValidExpressions()
    if len(validExprs) == 0 {
        t.Error("Data provider returned no valid expressions")
    }
}
```

### Benchmark Tests

**Purpose**: Measure performance of critical operations

**Characteristics**:
- Measure execution time and memory allocation
- Run multiple iterations for accuracy
- Provide baseline for performance regression detection

**Best Practices**:
```go
func BenchmarkCalculatorOperations(b *testing.B) {
    tests := []struct {
        name string
        expr string
    }{
        {"simple", "2 + 3"},
        {"complex", "2 * (3 + 4) - 5 / 2"},
        {"decimal", "3.14159 * 2.71828"},
    }

    for _, tt := range tests {
        b.Run(tt.name, func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                _, err := Evaluate(tt.expr)
                if err != nil {
                    b.Fatalf("Benchmark failed: %v", err)
                }
            }
        })
    }
}
```

## Coverage Requirements

### Coverage Thresholds

- **Overall Coverage**: Minimum 80%
- **Critical Components**: Minimum 90%
- **Core Business Logic**: Minimum 95%

### Coverage Reporting

Generate comprehensive coverage reports:
```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./...

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# Generate coverage summary
go tool cover -func=coverage.out
```

### Coverage Exclusions

The following may be excluded from coverage calculations:
- Generated code
- Test helper functions
- Platform-specific code (with appropriate build tags)
- Error handling paths that are difficult to test

## Cross-Platform Testing Approach

### Platform-Specific Optimizations

1. **Parallel Test Execution**: Use all available CPU cores across Windows, macOS, and Linux
2. **File Watching**: Automatic test execution on file changes with platform-aware implementations
3. **Shell Integration**: PowerShell, Bash, and Zsh scripts with cross-platform Makefile
4. **Performance Monitoring**: Platform-specific performance tools for comprehensive coverage

### Unix Test Tools

#### Unix Test Tools Script
```bash
# Run all tests with Unix optimizations
./tests/scripts/unix_test_tools.sh run --parallel 8 --verbose

# Watch file changes and run tests automatically
./tests/scripts/unix_test_tools.sh watch

# Validate Unix environment
./tests/scripts/unix_env_validate.sh
```

#### Makefile Integration
```bash
# Run all tests
make test

# Run specific test types
make test-unit
make test-integration
make test-benchmark

# Generate coverage reports
make coverage-html
make coverage-summary

# Quick development cycle
make dev-test
```

### Environment Validation

The Unix environment validation script checks:
- System compatibility (Linux, macOS)
- Go version and environment
- Required dependencies (git, make, bc, etc.)
- Performance settings (file descriptors, process limits)
- File system characteristics
- Network connectivity

## Cross-Platform Testing Strategy

### Platform Matrix

| Platform | Test Coverage | Coverage Reports | Benchmarks | Notes |
|----------|---------------|-----------------|------------|-------|
| Linux (Ubuntu) | ✅ Full | ✅ HTML + Text | ✅ Detailed | Primary platform |
| macOS | ✅ Full | ✅ HTML + Text | ✅ Basic | Development platform |
| Windows | ✅ Full | ❌ Limited | ❌ Limited | Full CI validation including E2E |

### CI Pipeline

The GitHub Actions CI pipeline provides:
- **Multi-platform testing**: Ubuntu, macOS, Windows
- **Coverage reporting**: HTML and text reports (Unix only)
- **Performance benchmarks**: Detailed benchmark analysis (Ubuntu only)
- **Security scanning**: Vulnerability detection
- **Artifact management**: Automatic upload and retention

### Platform-Specific Considerations

#### Linux
- Primary CI and development platform
- Full feature support including coverage and benchmarks
- Performance monitoring and optimization

#### macOS
- Development platform with full feature support
- File watching and development tools integration
- Basic benchmark support

#### Windows
- CI validation for compatibility
- Limited feature set (no coverage, no benchmarks)
- Focus on core functionality testing

## Test Data Management

### Test Fixtures

Organize test fixtures in the `integration` package:
```go
// integration_fixtures.go
type TestDataProvider struct {
    validExpressions   []string
    invalidExpressions []string
    complexExpressions []string
}

func NewTestDataProvider() *TestDataProvider {
    return &TestDataProvider{
        validExpressions: []string{
            "2 + 3", "10 - 4", "3 * 4", "15 / 3",
        },
        // ... more test data
    }
}
```

### Mock Objects

Use mock objects for external dependencies:
```go
type MockCalculator struct {
    results map[string]float64
    errors  map[string]error
}

func (m *MockCalculator) SetResult(expression string, result float64) {
    m.results[expression] = result
}

func (m *MockCalculator) Evaluate(expression string) (float64, error) {
    if result, exists := m.results[expression]; exists {
        return result, nil
    }
    if err, exists := m.errors[expression]; exists {
        return 0, err
    }
    return 0, fmt.Errorf("mock: unknown expression: %s", expression)
}
```

## Performance Testing

### Benchmark Standards

1. **Benchmark Naming**: Use descriptive names for benchmark operations
2. **Benchmark Duration**: Run benchmarks for at least 1 second
3. **Memory Allocation**: Track and minimize memory allocations
4. **Performance Regression**: Establish baselines and monitor for regressions

### Performance Monitoring

The framework provides:
- **Execution time** measurement
- **Memory allocation** tracking
- **CPU usage** monitoring
- **File descriptor** usage tracking

## CI/CD Integration

### GitHub Actions Features

1. **Multi-platform matrix testing**
2. **Coverage generation and threshold checking**
3. **Benchmark execution and reporting**
4. **Security vulnerability scanning**
5. **Test artifact management**
6. **Automated notifications**

### Quality Gates

- **All tests must pass** on all platforms
- **Coverage threshold** must be met (80% minimum)
- **No security vulnerabilities** detected
- **Code formatting** must be correct (go fmt)
- **No vet errors** (go vet)

## Development Workflow

### Local Development

1. **Setup environment**: Run Unix environment validation
2. **Write tests**: Follow TDD approach
3. **Run tests**: Use Unix tools or Makefile
4. **Check coverage**: Generate coverage reports
5. **Performance testing**: Run benchmarks
6. **Commit changes**: Ensure all tests pass

### Quick Commands

```bash
# Validate environment
make validate

# Quick test cycle
make dev-test

# Full test suite
make test-all

# Watch mode for development
make test-watch

# Coverage reports
make coverage-html
```

### Git Integration

- **Pre-commit hooks**: Run tests before committing
- **Branch protection**: Require test passing for merges
- **Pull requests**: Automated test execution
- **Issue tracking**: Test failures linked to issues

## Troubleshooting

### Common Issues

#### Test Failures
1. **Environment issues**: Run environment validation
2. **Dependency issues**: Check Go modules and dependencies
3. **Concurrency issues**: Check for race conditions
4. **Performance issues**: Check system resources and settings

#### Coverage Issues
1. **Low coverage**: Add tests for uncovered code paths
2. **Coverage generation**: Ensure Unix environment
3. **Coverage threshold**: Adjust threshold or improve coverage

#### Performance Issues
1. **Slow tests**: Use parallel execution
2. **Memory issues**: Check for memory leaks
3. **Resource issues**: Monitor system resources

### Debug Tools

- **Verbose output**: Use `--verbose` flag for detailed logging
- **Coverage debugging**: Use `make debug-coverage`
- **Performance debugging**: Use `make monitor-performance`
- **Test isolation**: Check isolation tests for concurrency issues

## Future Enhancements

### Planned Features

1. **Visual Testing**: Screenshot capture and comparison (✅ **Completed** in Story 0.2.2)
2. **E2E Testing**: Complete workflow testing (✅ **Completed** in Story 0.2.3)
3. **Performance Regression**: Historical performance tracking
4. **Test Analytics**: Test execution analytics and insights
5. **Advanced CI/CD**: Enhanced CI/CD features and integrations

### Tooling Enhancements

1. **Enhanced Reporting**: More detailed test reports and analytics
2. **IDE Integration**: Better IDE integration for testing
3. **Mobile Testing**: Mobile platform testing support
4. **Cloud Testing**: Cloud-based testing infrastructure

## Conclusion

This testing framework provides a comprehensive, cross-platform approach to testing the AcoustiCalc project with equal support for Windows, macOS, and Linux. By following these standards and guidelines, developers can ensure high-quality, reliable software with excellent test coverage and performance characteristics across all supported platforms.

The framework is designed to be extensible, maintainable, and aligned with modern testing best practices. It provides both the structure and tools needed for effective testing throughout the development lifecycle.

---

**Document Version**: 1.0
**Last Updated**: 2025-09-24
**Next Review**: 2025-12-24