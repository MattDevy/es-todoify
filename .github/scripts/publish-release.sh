#!/usr/bin/env bash

# ==============================================================================
# Script: publish-release.sh
# Description: Publishes the release using the existing release.sh script
#              and prepares version bump for PR creation.
#
# Usage: RELEASE_NOTES="..." ./publish-release.sh <version>
# Env Vars: 
#   - GITHUB_TOKEN (Required for gh cli)
#   - RELEASE_NOTES (Release notes content)
#   - GIT_USER_NAME (Optional: Git committer name, default: "github-actions[bot]")
#   - GIT_USER_EMAIL (Optional: Git committer email, default: "41898282+github-actions[bot]@users.noreply.github.com")
#   - RELEASE_BRANCH (Optional: Branch to release from, default: "main")
#
# Note: This script prepares the version bump but does NOT create the PR.
#       Use peter-evans/create-pull-request action in the workflow to create PR.
# ==============================================================================

set -euo pipefail

VERSION="${1:-}"

if [[ -z "$VERSION" ]]; then
    echo "Error: Version argument is required."
    echo "Usage: RELEASE_NOTES=\"...\" ./publish-release.sh <version>"
    exit 1
fi

VERSION_NUM="${VERSION#v}"

echo "Publishing release for version: $VERSION"

GIT_USER_NAME="${GIT_USER_NAME:-github-actions[bot]}"
GIT_USER_EMAIL="${GIT_USER_EMAIL:-41898282+github-actions[bot]@users.noreply.github.com}"
RELEASE_BRANCH="${RELEASE_BRANCH:-main}"

echo "Configuring git with user: $GIT_USER_NAME <$GIT_USER_EMAIL>"
echo "Release branch: $RELEASE_BRANCH"
git config user.name "$GIT_USER_NAME"
git config user.email "$GIT_USER_EMAIL"

# Snapshot tag for semantic-release to track changes
SNAPSHOT_TAG="snapshot-v${VERSION_NUM}"
echo "Creating snapshot tag ${SNAPSHOT_TAG} on current HEAD..."
if git rev-parse "$SNAPSHOT_TAG" >/dev/null 2>&1; then
    echo "Snapshot tag already exists, deleting..."
    git tag -d "$SNAPSHOT_TAG"
    git push origin ":refs/tags/$SNAPSHOT_TAG" 2>/dev/null || true
fi
git tag "$SNAPSHOT_TAG"
git push origin "$SNAPSHOT_TAG"

# Stage CHANGELOG.md changes for the detached commit (will be committed in release.sh)
if [[ -f "CHANGELOG.md" ]]; then
    echo "Staging CHANGELOG.md for detached commit..."
    git add CHANGELOG.md
fi

echo "Creating detached release..."
./.github/scripts/release.sh "$VERSION_NUM" "$RELEASE_NOTES"

IFS='.' read -r MAJOR MINOR PATCH <<< "$VERSION_NUM"

echo "Checking out $RELEASE_BRANCH branch..."
git fetch origin "$RELEASE_BRANCH"
git checkout "$RELEASE_BRANCH"
git pull origin "$RELEASE_BRANCH"

# Calculate next version
NEXT_PATCH=$((PATCH + 1))
NEXT_VERSION="${MAJOR}.${MINOR}.${NEXT_PATCH}-SNAPSHOT"

echo "Bumping to next snapshot version..."
./.github/scripts/bump-version.sh "$MAJOR" "$MINOR" "$PATCH"

if git diff --quiet version.go; then
    echo "Version is already set to ${NEXT_VERSION}, no changes needed."
    echo "Release published successfully (no version bump needed)!"
    if [[ -n "${GITHUB_OUTPUT:-}" ]]; then
        echo "HAS_CHANGES=false" >> "$GITHUB_OUTPUT"
    fi
else
    echo "Version bumped to ${NEXT_VERSION}"
    echo "Changes will be committed by create-pull-request action"
    if [[ -n "${GITHUB_OUTPUT:-}" ]]; then
        echo "HAS_CHANGES=true" >> "$GITHUB_OUTPUT"
        echo "NEXT_VERSION=${NEXT_VERSION}" >> "$GITHUB_OUTPUT"
    fi
fi

echo "Release published successfully!"

