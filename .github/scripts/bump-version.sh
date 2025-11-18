#!/bin/bash

# ==============================================================================
# Script: bump-version.sh
# Description: Calculates the next snapshot version (Patch + 1), updates the 
#              target file, and outputs the new version for GitHub Actions.
#
# Usage: ./bump-version.sh <major> <minor> <patch> [target_file]
# ==============================================================================

set -euo pipefail

# --- Inputs ---
MAJOR="${1:-}"
MINOR="${2:-}"
PATCH="${3:-}"
TARGET_FILE="${4:-version.go}"

# --- Validation ---
if [[ -z "$MAJOR" || -z "$MINOR" || -z "$PATCH" ]]; then
    echo "Error: Major, Minor, and Patch versions are required."
    echo "Usage: ./bump-version.sh <major> <minor> <patch> [target_file]"
    exit 1
fi

if [[ ! -f "$TARGET_FILE" ]]; then
    echo "Error: Target file '$TARGET_FILE' does not exist."
    exit 1
fi

# --- Logic ---
# Calculate next patch version
NEXT_PATCH=$((PATCH + 1))
NEXT_VERSION="${MAJOR}.${MINOR}.${NEXT_PATCH}-SNAPSHOT"

echo "Current Release: $MAJOR.$MINOR.$PATCH"
echo "Next Snapshot:   $NEXT_VERSION"

# Update the file
# Regex: Looks for 'Transport = "..."' and replaces the value
# This sed command is compatible with both GNU (Linux) and BSD (macOS) sed
echo "Updating $TARGET_FILE..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS (BSD sed) requires an empty string for in-place edit without backup
    sed -i '' "s/Transport *= *\".*\"/Transport = \"$NEXT_VERSION\"/" "$TARGET_FILE"
else
    # Linux (GNU sed)
    sed -i "s/Transport *= *\".*\"/Transport = \"$NEXT_VERSION\"/" "$TARGET_FILE"
fi

# --- Output for GitHub Actions ---
# If running in a GitHub Action, export the variable
if [[ -n "${GITHUB_OUTPUT:-}" ]]; then
    echo "next_version=$NEXT_VERSION" >> "$GITHUB_OUTPUT"
fi

echo "Bump complete."