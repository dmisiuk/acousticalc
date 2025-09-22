#!/bin/bash

# Script to release with GoReleaser using tag annotation as release notes

set -e

# Get the current tag
TAG=$(git describe --tags --exact-match 2>/dev/null || echo "")

if [ -z "$TAG" ]; then
    echo "Error: No tag found on current commit"
    exit 1
fi

echo "Creating release for tag: $TAG"

# Extract tag annotation message
git tag -l --format='%(contents)' "$TAG" > /tmp/release-notes.md

echo "Tag annotation content:"
echo "======================"
cat /tmp/release-notes.md
echo "======================"

# Run GoReleaser with the tag annotation as release notes
export PATH=/tmp/go/bin:$PATH
curl -sfL https://goreleaser.com/static/run | bash -s -- release --clean --release-notes /tmp/release-notes.md

# Clean up
rm -f /tmp/release-notes.md