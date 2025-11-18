#!/usr/bin/env bash

# Licensed to Elasticsearch B.V. under one or more agreements.
# Elasticsearch B.V. licenses this file to you under the Apache 2.0 License.
# See the LICENSE file in the project root for more information.

# ==============================================================================
# Script: release.sh
# Description: Creates a detached commit with a specific version, tags it, 
#              and publishes a GitHub Release.
#
# Usage: ./release.sh <version> <release_notes> [target_file]
# Env Vars: 
#   - GITHUB_TOKEN (Required for gh cli)
#   - GIT_USER_NAME (Optional: Git committer name, default: "Elastic Machine")
#   - GIT_USER_EMAIL (Optional: Git committer email, default: "elasticmachine@users.noreply.github.com")
# ==============================================================================

set -euo pipefail

# --- Configuration ---
VERSION="${1:-}"
RELEASE_NOTES="${2:-}"
TARGET_FILE="${3:-version.go}" # Default to elastictransport/version/version.go if not provided

# --- Helper Functions ---
log() {
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')] [INFO] $1"
}

error() {
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')] [ERROR] $1" >&2
    exit 1
}

check_dependencies() {
    command -v git >/dev/null 2>&1 || error "git is required but not installed."
    command -v gh >/dev/null 2>&1 || error "gh (GitHub CLI) is required but not installed."
    
    if [[ -z "${GITHUB_TOKEN:-}" ]]; then
        error "GITHUB_TOKEN environment variable is required."
    fi
}

# --- Main Execution ---

# 1. Validation
check_dependencies

if [[ -z "$VERSION" ]]; then
    error "Version argument is missing. Usage: ./script.sh <version> <notes>"
fi

if [[ ! -f "$TARGET_FILE" ]]; then
    error "Target file '$TARGET_FILE' does not exist."
fi

log "Starting floating release for version: $VERSION"

# 2. Configure Git Identity (Required in CI)
# Check if user.name is set, if not, set from environment or use defaults
if ! git config user.name >/dev/null; then
    GIT_USER_NAME="${GIT_USER_NAME:-Elastic Machine}"
    GIT_USER_EMAIL="${GIT_USER_EMAIL:-elasticmachine@users.noreply.github.com}"
    log "Configuring git user identity: $GIT_USER_NAME <$GIT_USER_EMAIL>"
    git config user.name "$GIT_USER_NAME"
    git config user.email "$GIT_USER_EMAIL"
fi

# 3. Create Detached State
log "Creating detached HEAD state..."
git checkout --detach HEAD

# 4. Modify Version File
# Strip leading 'v' from version if present
FILE_VERSION="${VERSION#v}"
log "Updating $TARGET_FILE to version $FILE_VERSION..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS (BSD sed) requires an empty string for in-place edit without backup
    sed -i '' "s/Transport *= *\".*\"/Transport = \"$FILE_VERSION\"/" "$TARGET_FILE"
else
    sed -i "s/Transport *= *\".*\"/Transport = \"$FILE_VERSION\"/" "$TARGET_FILE"
fi

if [[ $? -eq 0 ]]; then
    log "File updated successfully."
else
    error "Failed to update $TARGET_FILE using sed."
fi

RELEASE_VERSION="v${VERSION#v}"

# 5. Commit and Tag
log "Committing and tagging..."
# Add any staged files (like CHANGELOG.md) plus the modified version file
git add -A
git commit -m "chore: release $RELEASE_VERSION"

# Check if tag already exists
if git rev-parse "$RELEASE_VERSION" >/dev/null 2>&1; then
    log "Tag $RELEASE_VERSION already exists, deleting and recreating..."
    git tag -d "$RELEASE_VERSION"
    # Also delete from remote if it exists
    git push origin ":refs/tags/$RELEASE_VERSION" 2>/dev/null || true
fi

git tag "$RELEASE_VERSION"

# Get the commit SHA that we just tagged
COMMIT_SHA=$(git rev-parse HEAD)
log "Tagged commit SHA: $COMMIT_SHA"

# 6. Push Tag
log "Pushing tag $RELEASE_VERSION..."
# We purposefully only push the tag, not the detached commit, to the branch
git push origin "$RELEASE_VERSION"

# 7. Create GitHub Release
log "Creating GitHub Release..."
gh release create "$RELEASE_VERSION" \
    --title "$RELEASE_VERSION" \
    --notes "$RELEASE_NOTES" \
    --target "$COMMIT_SHA"

log "Release $RELEASE_VERSION created successfully!"
