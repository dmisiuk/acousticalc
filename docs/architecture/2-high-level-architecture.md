# 2. High-Level Architecture

## 2.1 Technical Summary

AcoustiCalc follows a modular architecture with clear separation of concerns:
- **Calculator Engine**: Core mathematical operations and expression parsing
- **TUI Layer**: Terminal-based user interface with mouse support
- **Audio System**: Sound feedback for operations and events
- **Configuration**: Settings management and persistence
- **History/Storage**: Calculation history persistence

The application is built as a single binary with no external dependencies, ensuring easy distribution across Windows, macOS, and Linux platforms.

## 2.2 Architecture Style

AcoustiCalc uses a layered architecture with the following components:
- Presentation Layer (TUI)
- Application Logic Layer (Calculator Engine)
- Audio Layer (Sound Feedback)
- Data Layer (History/Configuration)

## 2.3 Technology Stack

| Category | Technology | Version | Purpose | Rationale |
|---------|------------|---------|---------|-----------|
| Language | Go | 1.25.1+ | Primary development language | Performance, cross-platform compilation, single binary |
| TUI Framework | Bubble Tea | Latest | Terminal UI framework | Modern, feature-rich, active community |
| Audio Library | Oto | Latest | Cross-platform audio | Good cross-platform support |
| Testing | Go testing package | Bundled with Go | Unit and integration testing | Standard library, comprehensive coverage |
