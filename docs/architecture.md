# AcoustiCalc Architecture Document

## 1. Introduction

This document outlines the architecture for AcoustiCalc, a revolutionary terminal calculator that combines Go's performance with an intuitive Text-based User Interface (TUI), mouse support, and audio feedback. The architecture is designed to deliver a calculator experience that's both powerful and accessible to everyday users who work in terminal environments.

## 2. High-Level Architecture

### 2.1 Technical Summary

AcoustiCalc follows a modular architecture with clear separation of concerns:
- **Calculator Engine**: Core mathematical operations and expression parsing
- **TUI Layer**: Terminal-based user interface with mouse support
- **Audio System**: Sound feedback for operations and events
- **Configuration**: Settings management and persistence
- **History/Storage**: Calculation history persistence
- **Demo Infrastructure**: Automated testing, recording, and demo generation (Epic 0)

The application is built as a single binary with no external dependencies, ensuring easy distribution across Windows, macOS, and Linux platforms.

### 2.2 Architecture Style

AcoustiCalc uses a layered architecture with the following components:
- Presentation Layer (TUI)
- Application Logic Layer (Calculator Engine)
- Audio Layer (Sound Feedback)
- Data Layer (History/Configuration)
- Infrastructure Layer (Demo/Testing/Recording)

### 2.3 Technology Stack

| Category | Technology | Version | Purpose | Rationale |
|---------|------------|---------|---------|-----------|
| Language | Go | 1.25.1+ | Primary development language | Performance, cross-platform compilation, single binary |
| TUI Framework | Bubble Tea | Latest | Terminal UI framework | Modern, feature-rich, active community |
| Audio Library | Oto | Latest | Cross-platform audio | Good cross-platform support |
| Testing | Go testing package | Bundled with Go | Unit and integration testing | Standard library, comprehensive coverage |
| **Demo Infrastructure** | | | | |
| Terminal Recording | asciinema/ttyrec | Latest | Terminal session recording | Cross-platform terminal recording |
| Video Processing | ffmpeg | Latest | Video conversion and optimization | Industry standard for video processing |
| Input Visualization | Custom Go solution | - | Keystroke/mouse overlay system | Custom implementation for demo needs |
| GitHub Automation | GitHub Actions | Latest | CI/CD and demo automation | Native GitHub integration |
| Visual Testing | Custom TUI testing | - | Screenshot and UI validation | Specialized for TUI applications |

## 3. Component Design

### 3.1 Calculator Engine
- Responsible for parsing mathematical expressions
- Performs arithmetic operations
- Handles error validation
- Provides extensibility for advanced operations

### 3.2 TUI Layer
- Implements visual calculator layout
- Handles mouse events and keyboard input
- Manages UI state and rendering
- Provides responsive design for different terminal sizes

### 3.3 Audio System
- Plays sound effects for operations
- Manages audio configuration (volume, themes)
- Handles platform-specific audio implementations
- Provides fallback mechanisms

### 3.4 Configuration
- Manages user preferences
- Handles settings persistence
- Provides default configurations
- Supports customization options

### 3.5 History/Storage
- Stores calculation history
- Manages memory functions
- Provides persistence across sessions
- Implements efficient data storage

### 3.6 Demo Infrastructure (Epic 0)
- **Testing Framework**: Enhanced unit, integration, e2e, and UI testing with visual artifacts
- **Recording System**: Terminal session recording with visual input display (keystrokes, mouse)
- **PR Automation**: Automatic demo generation and embedding in pull requests
- **GitHub Integration**: Automated issue creation, branch management, and workflow triggers
- **Visual Validation**: Screenshot comparison and UI regression testing

**Key Components:**
- `pkg/testing/visual/`: Visual testing utilities and screenshot capture
- `pkg/recording/`: Terminal recording and video generation
- `pkg/demo/`: Demo scenario management and execution
- `.github/workflows/`: GitHub Actions for demo automation
- `demos/`: Generated demo artifacts and organization

## 4. Data Models

### 4.1 Calculation
- Expression: string representation of the mathematical expression
- Result: computed result of the expression
- Timestamp: when the calculation was performed

### 4.2 Configuration
- AudioEnabled: boolean for audio feedback
- VolumeLevel: integer for volume control
- Theme: string for UI theme selection
- HistoryEnabled: boolean for history tracking

### 4.3 Demo Infrastructure Models
- **TestArtifact**: Metadata for generated screenshots, videos, and test evidence
- **DemoScenario**: Structured definition of user interaction sequences for recording
- **RecordingSession**: Session metadata including input visualization and timing data
- **PRDemoContent**: Structured data for automatic PR demo embedding

## 5. API Design

AcoustiCalc is a single binary application with no external API. All functionality is accessed through the TUI interface.

## 6. Security Considerations

- Input validation for mathematical expressions
- Secure handling of user configuration data
- No network communication, minimizing attack surface
- Safe audio file handling

## 7. Performance Requirements

### 7.1 Core Application Performance
- Startup time: < 1 second on all platforms
- Response time: < 100ms for all operations
- Memory footprint: < 50MB typical usage
- CPU usage: minimal during calculations

### 7.2 Demo Infrastructure Performance
- Demo generation: < 2 minutes additional CI time
- Video recording: Minimal impact on application performance
- Screenshot capture: < 500ms per capture
- GitHub workflow execution: < 5 minutes total including demos
- Cross-platform demo consistency: 100% reproducible results

## 8. Deployment Architecture

### 8.1 Build Process
- Cross-platform compilation using Go
- Single binary output for each platform
- UPX compression for binary size optimization
- Automated GitHub Actions for releases

### 8.2 Distribution
- GitHub Releases for binary distribution
- Platform-specific packages (Homebrew, Scoop, AUR)
- Single binary installation with no dependencies

## 9. Testing Strategy

### 9.1 Unit Testing
- Calculator engine logic testing
- Individual component functionality verification
- Edge case and error condition testing

### 9.2 Integration Testing
- TUI interaction testing
- Audio system integration verification
- Cross-platform behavior validation

### 9.3 Manual Testing
- User experience validation
- Mouse interaction testing
- Audio feedback verification

## 10. Error Handling

- Graceful degradation when audio is unavailable
- Clear error messages for invalid expressions
- Recovery mechanisms for UI issues
- Logging for debugging purposes

## 11. Monitoring and Observability

- Basic logging for error tracking
- Performance metrics collection
- User interaction analytics (opt-in)
- Crash reporting (opt-in)

## 12. Future Extensibility

- Plugin architecture for additional functions
- Custom theme support
- Unit conversion capabilities
- Scientific calculator mode