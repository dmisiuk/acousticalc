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

# Extract tag annotation message for display
TAG_MESSAGE=$(git tag -l --format='%(contents)' "$TAG")

# Check if tag has annotation content
if [ -z "$TAG_MESSAGE" ]; then
    echo "Warning: Tag $TAG has no annotation message"
    echo "GoReleaser will create a default release."
else
    echo "Tag annotation content:"
    echo "======================"
    echo "$TAG_MESSAGE"
    echo "======================"
fi

# Run GoReleaser (it will use the tag annotation via the .goreleaser.yml header configuration)
echo "Running GoReleaser..."

# Check if goreleaser is available, if not download it
if ! command -v goreleaser &> /dev/null; then
    echo "GoReleaser not found, downloading..."
    curl -sfL https://goreleaser.com/static/run | bash -s -- release --clean
else
    goreleaser release --clean
fi

echo "Release completed successfully!"