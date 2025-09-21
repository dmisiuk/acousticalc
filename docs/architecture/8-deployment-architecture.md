# 8. Deployment Architecture

## 8.1 Build Process
- Cross-platform compilation using Go
- Single binary output for each platform
- UPX compression for binary size optimization
- Automated GitHub Actions for releases

## 8.2 Distribution
- GitHub Releases for binary distribution
- Platform-specific packages (Homebrew, Scoop, AUR)
- Single binary installation with no dependencies
