#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/lib.sh"

VERSION="${1:-}"
RELEASE_NOTES="${2:-}"
TARGET_FILE="${3:-version.go}"

log() {
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')] [INFO] $1"
}

command -v git >/dev/null 2>&1 || error "git is required but not installed"
command -v gh >/dev/null 2>&1 || error "gh (GitHub CLI) is required but not installed"
[[ -n "${GITHUB_TOKEN:-}" ]] || error "GITHUB_TOKEN environment variable is required"
[[ -n "$VERSION" ]] || error "Version argument is missing. Usage: ./release.sh <version> <notes>"
[[ -f "$TARGET_FILE" ]] || error "Target file '$TARGET_FILE' does not exist"

log "Starting floating release for version: $VERSION"

setup_git_user "Elastic Machine" "elasticmachine@users.noreply.github.com"

log "Creating detached HEAD state..."
git checkout --detach HEAD

FILE_VERSION="${VERSION#v}"
update_version_file "$TARGET_FILE" "$FILE_VERSION"

RELEASE_VERSION="v${VERSION#v}"

log "Committing and tagging..."
git add -A
git commit -m "chore: release $RELEASE_VERSION"

delete_tag_if_exists "$RELEASE_VERSION"
git tag "$RELEASE_VERSION"

COMMIT_SHA=$(git rev-parse HEAD)
log "Tagged commit SHA: $COMMIT_SHA"

log "Pushing tag $RELEASE_VERSION..."
git push origin "$RELEASE_VERSION"

log "Creating GitHub Release..."
gh release create "$RELEASE_VERSION" \
    --title "$RELEASE_VERSION" \
    --notes "$RELEASE_NOTES" \
    --target "$COMMIT_SHA"

log "Release $RELEASE_VERSION created successfully!"
