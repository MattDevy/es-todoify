#!/usr/bin/env bash

set -euo pipefail

VERSION_FILE_PATTERN="${VERSION_FILE_PATTERN:-Transport}"

error() {
    echo "[ERROR] $1" >&2
    exit 1
}

setup_git_user() {
    local default_name="${1:-github-actions[bot]}"
    local default_email="${2:-41898282+github-actions[bot]@users.noreply.github.com}"
    
    if ! git config user.name >/dev/null 2>&1; then
        local user_name="${GIT_USER_NAME:-$default_name}"
        local user_email="${GIT_USER_EMAIL:-$default_email}"
        echo "Configuring git user: $user_name <$user_email>"
        git config user.name "$user_name"
        git config user.email "$user_email"
    fi
}

update_version_file() {
    local target_file="$1"
    local version="$2"
    
    if [[ ! -f "$target_file" ]]; then
        error "Target file '$target_file' does not exist."
    fi
    
    echo "Updating $target_file to version $version..."
    
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "s/${VERSION_FILE_PATTERN} *= *\".*\"/${VERSION_FILE_PATTERN} = \"$version\"/" "$target_file"
    else
        sed -i "s/${VERSION_FILE_PATTERN} *= *\".*\"/${VERSION_FILE_PATTERN} = \"$version\"/" "$target_file"
    fi
    
    [[ $? -eq 0 ]] || error "Failed to update $target_file"
}

delete_tag_if_exists() {
    local tag="$1"
    
    if git rev-parse "$tag" >/dev/null 2>&1; then
        echo "Tag $tag already exists, deleting..."
        git tag -d "$tag"
        git push origin ":refs/tags/$tag" 2>/dev/null || true
    fi
}

write_github_output() {
    local key="$1"
    local value="$2"
    
    if [[ -n "${GITHUB_OUTPUT:-}" ]]; then
        echo "${key}=${value}" >> "$GITHUB_OUTPUT"
    fi
}

