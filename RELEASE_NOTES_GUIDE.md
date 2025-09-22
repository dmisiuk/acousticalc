# Release Notes with Tag Annotations

This project is configured to use Git tag annotation messages as GitHub release notes instead of auto-generated changelogs.

## How It Works

1. **GoReleaser Configuration**: The `.goreleaser.yml` file has `changelog.disable: true` to prevent automatic changelog generation.

2. **Release Script**: The `release-with-tag-annotation.sh` script extracts the tag annotation message and passes it to GoReleaser using the `--release-notes` flag.

3. **Tag Annotations**: When you create an annotated tag, the message becomes the release notes.

## Usage

### Creating a Release

1. **Create an annotated tag with your release notes**:
   ```bash
   git tag -a v1.2.0 -m "🎉 Version 1.2.0 Release

   ## What's New
   - Added new acoustic calculation algorithms
   - Improved performance by 25%
   - Fixed critical bug in frequency analysis

   ## Breaking Changes
   - Renamed `calculate()` to `computeAcoustics()`

   ## Installation
   Download the appropriate binary for your platform from the release assets below.

   ## Contributors
   Thanks to all contributors who made this release possible!"
   ```

2. **Push the tag**:
   ```bash
   git push origin v1.2.0
   ```

3. **Run the release script**:
   ```bash
   ./release-with-tag-annotation.sh
   ```

### Alternative: Manual GoReleaser

If you prefer to run GoReleaser manually:

```bash
# Extract tag annotation to a file
git tag -l --format='%(contents)' $(git describe --tags --exact-match) > release-notes.md

# Run GoReleaser with custom release notes
goreleaser release --clean --release-notes release-notes.md

# Clean up
rm release-notes.md
```

## Benefits

- **Rich Release Notes**: Use Markdown formatting, emojis, and detailed descriptions
- **Consistent Process**: Same workflow for all releases
- **Version Control**: Release notes are stored in Git history with the tag
- **No Duplication**: Avoid maintaining separate changelog files

## Examples

### Simple Release
```bash
git tag -a v1.0.1 -m "Bug fix release

Fixed issue with division by zero in acoustic calculations."
```

### Feature Release
```bash
git tag -a v1.1.0 -m "🚀 New Features Release

## New Features
- Added support for multiple audio formats
- Implemented real-time frequency analysis
- Added command-line progress indicators

## Improvements
- 30% faster calculation performance
- Reduced memory usage by 15%
- Better error messages

## Bug Fixes
- Fixed crash when processing large files
- Corrected frequency range validation

Download the latest version from the assets below!"
```

## Tips

1. **Use Markdown**: Tag annotations support full Markdown formatting
2. **Include Emojis**: Make your release notes more engaging
3. **Structure Content**: Use headers, lists, and sections for clarity
4. **Mention Breaking Changes**: Always highlight breaking changes
5. **Thank Contributors**: Acknowledge community contributions

## Troubleshooting

- **Empty Release Notes**: If a tag has no annotation, the script creates default release notes
- **Missing Tag**: The script will fail if run on a commit without a tag
- **GoReleaser Errors**: Check that all required environment variables (like `GITHUB_TOKEN`) are set