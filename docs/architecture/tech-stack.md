# Technology Stack

## Core Technologies

| Category | Technology | Version | Purpose | Rationale |
|---------|------------|---------|---------|-----------|
| Language | Go | 1.24+ | Primary development language | Performance, cross-platform compilation, single binary |
| TUI Framework | Bubble Tea | Latest | Terminal UI framework | Modern, feature-rich, active community |
| Audio Library | Oto | Latest | Cross-platform audio | Good cross-platform support |
| Testing | Go testing package | Bundled with Go | Unit and integration testing | Standard library, comprehensive coverage |

## Development Tools

| Tool | Purpose | Notes |
|------|---------|-------|
| Go Modules | Dependency management | Standard Go tooling |
| GitHub Actions | CI/CD pipeline | Multi-platform builds |
| GoReleaser | Release automation | Cross-platform binary distribution |

## Platform Support

- **Windows 10+**: amd64, arm64
- **macOS 10.15+**: amd64, arm64 (Apple Silicon)
- **Linux**: amd64, arm64

## Build Requirements

- Go 1.24+ for all builds
- CGO enabled for audio libraries
- Static linking where possible
- UPX compression for binary size optimization

## Dependencies

### Core Dependencies
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/hajimehoshi/oto/v2` - Audio library

### Demo Infrastructure Dependencies (Epic 0)
- `github.com/asciinema/asciinema` - Terminal recording
- `github.com/charmbracelet/lipgloss` - Terminal styling for visual testing
- `github.com/disintegration/imaging` - Image processing for screenshots
- `github.com/go-vgo/robotgo` - Cross-platform screenshot capture
- Custom Go packages for visual testing and artifact management

### Development Dependencies
- Standard Go testing package
- GitHub Actions workflows
- GoReleaser configuration

### External Tools (System Dependencies)
- **ffmpeg** - Video processing and format conversion
- **asciinema** - Terminal session recording
- **ImageMagick** - Image processing and optimization
- **xdotool/xclip** - Linux input simulation (for demos)
- **screencapture** - macOS screenshot utilities
- **PowerShell** - Windows automation utilities

## Architecture Principles

- **Single Binary**: No external dependencies at runtime
- **Cross-Platform**: Consistent experience across all platforms
- **Modular Design**: Clear separation of concerns
- **Performance**: Startup time < 1 second, response time < 100ms
- **Resource Efficient**: Memory footprint < 50MB typical usage