#!/bin/bash

# Script to release with GoReleaser using tag annotation as release notes
# This script extracts the tag annotation message and uses it as the release notes
# instead of the auto-generated changelog.

set -e

# Get the current tag
TAG=$(git describe --tags --exact-match 2>/dev/null || echo "")

if [ -z "$TAG" ]; then
    echo "Error: No tag found on current commit"
    echo "Please create an annotated tag first:"
    echo "  git tag -a v1.0.0 -m 'Your release notes here'"
    exit 1
fi

echo "Creating release for tag: $TAG"

# Create temporary file for release notes
RELEASE_NOTES_FILE=$(mktemp)

# Extract tag annotation message
git tag -l --format='%(contents)' "$TAG" > "$RELEASE_NOTES_FILE"

# Check if tag has annotation content
if [ ! -s "$RELEASE_NOTES_FILE" ]; then
    echo "Warning: Tag $TAG has no annotation message"
    echo "Creating default release notes..."
    echo "Release $TAG" > "$RELEASE_NOTES_FILE"
    echo "" >> "$RELEASE_NOTES_FILE"
    echo "This release was created automatically." >> "$RELEASE_NOTES_FILE"
fi

echo "Release notes content:"
echo "======================"
cat "$RELEASE_NOTES_FILE"
echo "======================"

# Run GoReleaser with the tag annotation as release notes
echo "Running GoReleaser..."
goreleaser release --clean --release-notes "$RELEASE_NOTES_FILE"

# Clean up
rm -f "$RELEASE_NOTES_FILE"

echo "Release completed successfully!"