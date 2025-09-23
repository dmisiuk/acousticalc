# Source Tree Structure

## Overview

AcoustiCalc follows a standard Go project layout with clear separation of concerns and modular organization.

## Directory Structure

```
acousticalc/
├── cmd/                          # Command line applications
│   └── acousticalc/             # Main application entry point
│       ├── main.go              # Application main function
│       └── main_integration_test.go # Integration tests
├── pkg/                         # Public packages (importable by external projects)
│   └── calculator/              # Core calculator engine
│       ├── calculator.go        # Main calculator implementation
│       ├── calculator_test.go   # Unit tests
│       └── calculator_coverage_test.go # Coverage-focused tests
├── internal/                    # Private packages (not importable externally)
│   ├── ui/                      # TUI components and interfaces
│   ├── audio/                   # Audio system implementation
│   ├── config/                  # Configuration management
│   └── history/                 # Calculation history storage
├── docs/                        # Documentation
│   ├── architecture/            # Technical documentation
│   ├── prd/                     # Product requirements (sharded)
│   ├── epics/                   # Epic definitions
│   ├── stories/                 # User stories
│   └── qa/                      # Quality assurance documentation
├── assets/                      # Static assets
│   └── logo.png                 # Project logo
├── web-bundles/                 # B-Mad method configurations
│   ├── agents/                  # Agent definitions
│   ├── expansion-packs/         # Additional capabilities
│   └── teams/                   # Team configurations
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── README.md                    # Project overview and setup
├── LICENSE                      # Project license
└── RELEASE_NOTES_GUIDE.md      # Release documentation guide
```

## Package Organization

### `/cmd/acousticalc`
**Purpose**: Main application entry point
- Contains the main() function
- Handles command-line argument parsing
- Orchestrates application components
- Integration tests for end-to-end functionality

### `/pkg/calculator`
**Purpose**: Core calculation engine (public API)
- Mathematical expression parsing
- Arithmetic operations implementation
- Error handling for invalid expressions
- Public interface for calculator functionality

### `/internal/ui` (Planned)
**Purpose**: Terminal User Interface components
- Bubble Tea TUI implementation
- Mouse event handling
- Keyboard navigation
- Visual feedback and animations

### `/internal/audio` (Planned)
**Purpose**: Audio feedback system
- Sound effect playback using Oto library
- Volume control and mute functionality
- Cross-platform audio handling
- Audio configuration management

### `/internal/config` (Planned)
**Purpose**: Configuration management
- Settings persistence
- User preferences handling
- Environment variable support
- Default configuration values

### `/internal/history` (Planned)
**Purpose**: Calculation history
- History storage and retrieval
- Memory functions (store, recall)
- Session persistence
- History search and filtering

## File Naming Conventions

### Go Source Files
- **Main files**: `main.go`
- **Implementation files**: `{package_name}.go`
- **Test files**: `{package_name}_test.go`
- **Benchmark files**: `{package_name}_bench_test.go`
- **Integration tests**: `{component}_integration_test.go`

### Documentation Files
- **Markdown**: Use `.md` extension
- **Architecture docs**: Numbered prefix (e.g., `1-introduction.md`)
- **Stories**: Format `{epic}.{story}.story.md`
- **Epics**: Format `{number}.epic.md`

## Import Path Structure

```go
// External imports (calculator engine)
import "github.com/dmisiuk/acousticalc/pkg/calculator"

// Internal imports (not available to external projects)
import "github.com/dmisiuk/acousticalc/internal/ui"
import "github.com/dmisiuk/acousticalc/internal/audio"
import "github.com/dmisiuk/acousticalc/internal/config"
```

## Build Artifacts

### Generated Files
- `acousticalc` - Main binary executable
- `coverage.out` - Test coverage data
- `go.sum` - Module checksum verification

### Distribution
- Cross-platform binaries generated via GoReleaser
- Single binary with no external dependencies
- Static linking for maximum portability

## Development Workflow

### Local Development
1. **Source**: Edit code in appropriate package directories
2. **Test**: Run tests with `go test ./...`
3. **Build**: Build with `go build ./cmd/acousticalc`
4. **Run**: Execute with `./acousticalc`

### Testing Structure
- **Unit tests**: Package-level functionality testing
- **Integration tests**: Cross-package interaction testing
- **Coverage tests**: Ensure adequate test coverage
- **Benchmark tests**: Performance measurement

## Code Organization Principles

### Package Responsibilities
- **Single Responsibility**: Each package has one clear purpose
- **Dependency Direction**: Dependencies flow inward (cmd → internal → pkg)
- **Interface Segregation**: Small, focused interfaces
- **Testability**: All packages easily unit testable

### File Organization
- Keep related functionality together
- Separate test files from implementation
- Use descriptive file names
- Maintain consistent naming patterns

## Future Growth

### Planned Additions
- `/internal/plugins` - Plugin architecture support
- `/internal/themes` - UI theming system
- `/web` - Web interface components (if needed)
- `/scripts` - Build and deployment scripts

### Scalability Considerations
- Modular design allows for easy feature addition
- Clear package boundaries enable team development
- Public API in `/pkg` supports external integrations
- Internal packages protect implementation details