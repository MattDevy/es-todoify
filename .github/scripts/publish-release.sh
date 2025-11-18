#!/usr/bin/env bash

# ==============================================================================
# Script: publish-release.sh
# Description: Publishes the release using the existing release.sh script
#              and then bumps the version on main branch.
#
# Usage: ./publish-release.sh <version> <release_notes>
# Env Vars: GITHUB_TOKEN (Required for gh cli)
# ==============================================================================

set -euo pipefail

VERSION="${1:-}"
RELEASE_NOTES="${2:-}"

if [[ -z "$VERSION" ]]; then
    echo "Error: Version argument is required."
    echo "Usage: ./publish-release.sh <version> <release_notes>"
    exit 1
fi

# Remove 'v' prefix if present for version parsing
VERSION_NUM="${VERSION#v}"

echo "Publishing release for version: $VERSION"

# Run the existing release script
echo "Creating detached release..."
./.github/scripts/release.sh "$VERSION_NUM" "$RELEASE_NOTES"

# Parse version components
IFS='.' read -r MAJOR MINOR PATCH <<< "$VERSION_NUM"

# Checkout main branch to bump version
echo "Checking out main branch..."
git checkout main
git pull origin main

# Bump to next snapshot version
echo "Bumping to next snapshot version..."
./.github/scripts/bump-version.sh "$MAJOR" "$MINOR" "$PATCH"

# Commit and push the version bump
echo "Committing version bump..."
NEXT_PATCH=$((PATCH + 1))
NEXT_VERSION="${MAJOR}.${MINOR}.${NEXT_PATCH}-SNAPSHOT"

git config user.name "github-actions[bot]"
git config user.email "41898282+github-actions[bot]@users.noreply.github.com"

git add version.go
git commit -m "chore: bump version to ${NEXT_VERSION} [skip ci]"
git push origin main

echo "Release published and version bumped successfully!"

