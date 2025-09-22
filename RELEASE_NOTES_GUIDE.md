# AcoustiCalc Release Notes Guide

This guide explains how to create releases with custom release notes using tag annotation messages instead of auto-generated changelogs.

## Overview

The solution uses GoReleaser's built-in template system to extract tag annotation messages:

1. **GoReleaser Configuration**: The `.goreleaser.yml` file uses `release.header: "{{ .TagBody }}"` to extract tag annotation content and `changelog.disable: true` to prevent automatic changelog generation.
2. **Release Script**: The `release-with-tag-annotation.sh` script simplifies the release process with auto-download capability.

## Quick Start

1. **Create an annotated tag with release notes**:
   ```bash
   git tag -a v1.0.0 -m "🚀 Release v1.0.0

   ## New Features
   - Added acoustic calculation engine
   - Implemented frequency analysis
   - Added CLI interface

   ## Bug Fixes
   - Fixed calculation precision issues
   - Resolved memory leaks

   ## Breaking Changes
   - Changed API endpoint structure"
   ```

2. **Run the release script**:
   ```bash
   ./release-with-tag-annotation.sh
   ```

That's it! GoReleaser will automatically use the tag annotation as the GitHub release notes.

## How It Works

### GoReleaser Configuration

The `.goreleaser.yml` file contains:

```yaml
release:
  name_template: "AcoustiCalc {{.Version}}"
  # Use tag annotation as release notes
  header: |
    {{ .TagBody }}

changelog:
  # Disable automatic changelog generation to use tag annotation messages instead
  disable: true
```

The `{{ .TagBody }}` template variable extracts the full content of the tag annotation message.

### Release Script Features

The `release-with-tag-annotation.sh` script:

- **Auto-detects current tag**: Finds the tag pointing to the current commit
- **Validates tag annotation**: Warns if no annotation message exists
- **Auto-downloads GoReleaser**: Downloads GoReleaser if not installed
- **Displays preview**: Shows the tag annotation content before release
- **Error handling**: Provides clear error messages and validation

## Benefits

- **Rich Release Notes**: Use Markdown formatting, emojis, and detailed descriptions
- **Consistent Process**: Same workflow for all releases
- **Version Control**: Release notes are stored in Git history with the tag
- **No Duplication**: Avoid maintaining separate changelog files
- **Simplified Script**: No temporary file handling needed
- **Built-in Support**: Uses GoReleaser's native template system

## Examples

### Simple Release
```bash
git tag -a v1.0.1 -m "Bug fix release

Fixed critical issue with calculation accuracy."
./release-with-tag-annotation.sh
```

### Feature Release
```bash
git tag -a v1.1.0 -m "🎉 New Feature Release

## What's New
✅ Added real-time visualization
✅ Improved performance by 40%
✅ New export formats (PDF, CSV)

## Bug Fixes
🐛 Fixed memory leak in long calculations
🐛 Resolved UI freezing issues

## Documentation
📚 Updated API documentation
📚 Added tutorial videos

Thanks to all contributors! 🙏"
./release-with-tag-annotation.sh
```

### Breaking Changes Release
```bash
git tag -a v2.0.0 -m "⚠️ Major Release v2.0.0

## 🚨 Breaking Changes
- API endpoints now use v2 prefix
- Configuration file format changed
- Minimum Go version: 1.19

## 🆕 New Features
- Complete API redesign
- Enhanced calculation engine
- New plugin system

## 📖 Migration Guide
See MIGRATION.md for upgrade instructions.

## 🐛 Bug Fixes
- Fixed all known calculation edge cases
- Resolved performance bottlenecks

This release represents a major milestone! 🎯"
./release-with-tag-annotation.sh
```

## Troubleshooting

### No Tag Found Error
```
Error: No tag found on current commit
```
**Solution**: Create an annotated tag first:
```bash
git tag -a v1.0.0 -m 'Your release notes here'
```

### Empty Release Notes
If your GitHub release has empty notes, ensure you're using an annotated tag (not a lightweight tag):
```bash
# Wrong (lightweight tag)
git tag v1.0.0

# Correct (annotated tag)
git tag -a v1.0.0 -m 'Release notes here'
```

### GoReleaser Not Found
The script automatically downloads GoReleaser if not installed. If you prefer to install it manually:
```bash
# Using Go
go install github.com/goreleaser/goreleaser@latest

# Using Homebrew (macOS)
brew install goreleaser

# Using curl
curl -sfL https://goreleaser.com/static/run | bash
```

## Advanced Usage

### Custom Release Names
Modify the `name_template` in `.goreleaser.yml`:
```yaml
release:
  name_template: "AcoustiCalc {{.Version}} - {{.Date}}"
```

### Additional Release Content
Add more sections to the header template:
```yaml
release:
  header: |
    {{ .TagBody }}
    
    ## Download
    Choose the appropriate binary for your platform below.
    
    ## Checksums
    Verify your download using the checksums file.
```

### Testing Releases
Enable draft mode for testing:
```yaml
release:
  draft: true  # Creates draft releases for testing
```

## Solution Summary

This implementation resolves GoReleaser issue #7 by:

1. ✅ **Using tag annotation messages** instead of auto-generated changelogs
2. ✅ **Disabling automatic changelog generation** via `changelog.disable: true`
3. ✅ **Leveraging GoReleaser's template system** with `{{ .TagBody }}`
4. ✅ **Providing a simple release script** with auto-download capability
5. ✅ **Supporting rich Markdown formatting** in release notes
6. ✅ **Maintaining version-controlled release notes** in Git history

The solution is production-ready and provides a streamlined workflow for creating releases with custom, meaningful release notes.