#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/lib.sh"

VERSION="${1:-}"
[[ -n "$VERSION" ]] || error "Version argument is required. Usage: RELEASE_NOTES=\"...\" ./publish-release.sh <version>"

VERSION_NUM="${VERSION#v}"
RELEASE_BRANCH="${RELEASE_BRANCH:-main}"

echo "Publishing release for version: $VERSION"
echo "Release branch: $RELEASE_BRANCH"

setup_git_user

SNAPSHOT_TAG="snapshot-v${VERSION_NUM}"
echo "Creating snapshot tag ${SNAPSHOT_TAG}..."
delete_tag_if_exists "$SNAPSHOT_TAG"
git tag "$SNAPSHOT_TAG"
git push origin "$SNAPSHOT_TAG"

[[ -f "CHANGELOG.md" ]] && git add CHANGELOG.md

echo "Creating detached release..."
"${SCRIPT_DIR}/release.sh" "$VERSION_NUM" "$RELEASE_NOTES"

IFS='.' read -r MAJOR MINOR PATCH <<< "$VERSION_NUM"

echo "Checking out $RELEASE_BRANCH..."
git fetch origin "$RELEASE_BRANCH"
git checkout "$RELEASE_BRANCH"
git pull origin "$RELEASE_BRANCH"

NEXT_PATCH=$((PATCH + 1))
NEXT_VERSION="${MAJOR}.${MINOR}.${NEXT_PATCH}-SNAPSHOT"

echo "Bumping to next snapshot version..."
"${SCRIPT_DIR}/bump-version.sh" "$MAJOR" "$MINOR" "$PATCH"

if git diff --quiet version.go; then
    echo "Version already set to ${NEXT_VERSION}, no changes needed"
    write_github_output "HAS_CHANGES" "false"
else
    echo "Version bumped to ${NEXT_VERSION}"
    write_github_output "HAS_CHANGES" "true"
    write_github_output "NEXT_VERSION" "$NEXT_VERSION"
fi

echo "Release published successfully!"

