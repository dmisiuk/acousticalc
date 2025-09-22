# AcoustiCalc Release Notes Guide

This guide explains how to create releases with custom release notes using tag annotation messages instead of auto-generated changelogs.

## Overview

The solution uses GoReleaser's built-in template system with GitHub Actions automation:

1. **GoReleaser Configuration**: The `.goreleaser.yml` file uses `release.header: "{{ .TagBody }}"` to extract tag annotation content and `changelog.disable: true` to prevent automatic changelog generation.
2. **GitHub Actions Workflow**: The `.github/workflows/release.yml` automatically runs GoReleaser when tags are pushed to the repository.

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

2. **Push the tag to GitHub**:
   ```bash
   git push origin v1.0.0
   ```

That's it! GitHub Actions will automatically run GoReleaser and create a release using the tag annotation as the release notes.

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

### GitHub Actions Workflow

The `.github/workflows/release.yml` workflow:

- **Triggers on tag push**: Automatically runs when tags matching `v*` are pushed
- **Sets up Go environment**: Configures Go 1.21 for building
- **Runs GoReleaser**: Uses the official GoReleaser action with our configuration
- **Creates GitHub release**: Automatically publishes release with tag annotation content
- **Uploads artifacts**: Builds and uploads binaries for all platforms

## Benefits

- **Rich Release Notes**: Use Markdown formatting, emojis, and detailed descriptions
- **Automated Process**: No manual intervention needed - just push tags
- **Version Control**: Release notes are stored in Git history with the tag
- **No Duplication**: Avoid maintaining separate changelog files
- **CI/CD Integration**: Fully automated via GitHub Actions
- **Built-in Support**: Uses GoReleaser's native template system

## Examples

### Simple Release
```bash
git tag -a v1.0.1 -m "Bug fix release

Fixed critical issue with calculation accuracy."
git push origin v1.0.1
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
git push origin v1.1.0
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
git push origin v2.0.0
```

## Troubleshooting

### GitHub Actions Workflow Not Triggering
If the release workflow doesn't run when you push a tag:
1. Ensure your tag follows the `v*` pattern (e.g., `v1.0.0`, `v2.1.3`)
2. Check that the tag was pushed to the repository: `git push origin v1.0.0`
3. Verify the workflow file exists at `.github/workflows/release.yml`

### Empty Release Notes
If your GitHub release has empty notes, ensure you're using an annotated tag (not a lightweight tag):
```bash
# Wrong (lightweight tag)
git tag v1.0.0

# Correct (annotated tag)
git tag -a v1.0.0 -m 'Release notes here'
```

### Workflow Fails to Build
Check the GitHub Actions logs for build errors:
1. Go to your repository on GitHub
2. Click the "Actions" tab
3. Find the failed workflow run
4. Check the logs for specific error messages

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
4. ✅ **Automating releases via GitHub Actions** - no manual intervention needed
5. ✅ **Supporting rich Markdown formatting** in release notes
6. ✅ **Maintaining version-controlled release notes** in Git history

The solution is production-ready and provides a fully automated CI/CD workflow for creating releases with custom, meaningful release notes. Simply push an annotated tag and GitHub Actions handles the rest!