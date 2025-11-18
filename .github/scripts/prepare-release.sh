#!/usr/bin/env bash

# ==============================================================================
# Script: prepare-release.sh
# Description: Prepares the repository for a semantic-release.
#              Commits the updated CHANGELOG.md to the main branch.
#
# Usage: ./prepare-release.sh <version>
# ==============================================================================

set -euo pipefail

VERSION="${1:-}"

if [[ -z "$VERSION" ]]; then
    echo "Error: Version argument is required."
    echo "Usage: ./prepare-release.sh <version>"
    exit 1
fi

echo "Preparing release for version: $VERSION"

# Remove 'v' prefix if present
VERSION="${VERSION#v}"

# Configure git if needed
if ! git config user.name >/dev/null; then
    git config user.name "github-actions[bot]"
    git config user.email "41898282+github-actions[bot]@users.noreply.github.com"
fi

# Commit and push the changelog
echo "Committing CHANGELOG.md..."
git add CHANGELOG.md
git commit -m "chore(release): ${VERSION} [skip ci]" || echo "No changes to commit"
git push origin main

echo "Preparation complete for version: $VERSION"

