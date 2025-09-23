# Coding Standards

## Go Language Standards

### General Principles
- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` for consistent formatting
- Use `golint` and `go vet` for code quality
- Follow Go community best practices

### Code Organization

#### Package Structure
```go
// Package declaration
package calculator

// Import organization
import (
    // Standard library imports first
    "fmt"
    "strconv"

    // Third-party imports
    "github.com/charmbracelet/bubbletea"

    // Local imports last
    "github.com/dmisiuk/acousticalc/internal/config"
)
```

#### Naming Conventions
- **Variables**: camelCase (`calculatorEngine`, `userInput`)
- **Functions**: PascalCase for exported (`Calculate`, `ParseExpression`)
- **Functions**: camelCase for private (`parseToken`, `validateInput`)
- **Constants**: UPPER_SNAKE_CASE (`MAX_PRECISION`, `DEFAULT_TIMEOUT`)
- **Types**: PascalCase (`Calculator`, `AudioPlayer`)

### Error Handling

#### Error Creation
```go
// Use standard errors package
import "errors"

var (
    ErrInvalidExpression = errors.New("invalid mathematical expression")
    ErrDivisionByZero   = errors.New("division by zero")
)

// Or use fmt.Errorf for dynamic errors
func validateInput(input string) error {
    if len(input) == 0 {
        return fmt.Errorf("input cannot be empty")
    }
    return nil
}
```

#### Error Handling Patterns
```go
// Always check errors
result, err := calculator.Evaluate(expression)
if err != nil {
    return nil, fmt.Errorf("failed to evaluate expression: %w", err)
}

// Use early returns
func processCalculation(input string) (*Result, error) {
    if input == "" {
        return nil, ErrInvalidExpression
    }

    tokens, err := tokenize(input)
    if err != nil {
        return nil, err
    }

    // Continue processing...
}
```

### Testing Standards

#### Test File Organization
- Test files: `*_test.go`
- Test functions: `TestFunctionName`
- Benchmark functions: `BenchmarkFunctionName`

#### Test Structure
```go
func TestCalculatorAdd(t *testing.T) {
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

### Documentation

#### Package Documentation
```go
// Package calculator provides mathematical expression parsing and evaluation
// capabilities for the AcoustiCalc terminal calculator.
//
// The package supports basic arithmetic operations (addition, subtraction,
// multiplication, division) and follows standard mathematical operator
// precedence rules.
package calculator
```

#### Function Documentation
```go
// Calculate parses and evaluates a mathematical expression string.
// It returns the result as a float64 and any error encountered during parsing
// or evaluation.
//
// Supported operations: +, -, *, /, (), and decimal numbers.
// Example: Calculate("2 + 3 * 4") returns 14.0, nil
func Calculate(expression string) (float64, error) {
    // Implementation...
}
```

### Performance Guidelines

#### Memory Management
- Prefer stack allocation over heap when possible
- Use object pooling for frequently allocated objects
- Avoid unnecessary string concatenation in loops

#### Concurrency
- Use channels for communication between goroutines
- Prefer `sync.Once` for one-time initialization
- Use context for cancellation and timeouts

### Security Standards

#### Input Validation
```go
func validateExpression(expr string) error {
    // Check length limits
    if len(expr) > MaxExpressionLength {
        return ErrExpressionTooLong
    }

    // Validate characters
    for _, r := range expr {
        if !isValidCharacter(r) {
            return ErrInvalidCharacter
        }
    }

    return nil
}
```

#### Sanitization
- Always validate and sanitize user input
- Use whitelisting rather than blacklisting
- Implement appropriate length limits

## Project-Specific Standards

### Audio System
- Always handle audio initialization failures gracefully
- Provide fallback for systems without audio support
- Use appropriate audio sample rates and formats

### TUI Components
- Follow consistent event handling patterns
- Implement proper keyboard navigation
- Ensure accessibility compliance

### Configuration
- Use sensible defaults for all settings
- Validate configuration values
- Support environment variable overrides

## Quality Gates

### Pre-commit Requirements
- [ ] Code passes `go fmt`
- [ ] Code passes `go vet`
- [ ] All tests pass
- [ ] Test coverage > 80%
- [ ] No security vulnerabilities detected

### Code Review Checklist
- [ ] Follows naming conventions
- [ ] Proper error handling
- [ ] Adequate test coverage
- [ ] Documentation is clear and complete
- [ ] No performance regressions