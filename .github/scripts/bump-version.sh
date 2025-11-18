#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/lib.sh"

MAJOR="${1:-}"
MINOR="${2:-}"
PATCH="${3:-}"
TARGET_FILE="${4:-version.go}"

[[ -n "$MAJOR" && -n "$MINOR" && -n "$PATCH" ]] || error "Major, Minor, and Patch versions are required. Usage: ./bump-version.sh <major> <minor> <patch> [target_file]"

NEXT_PATCH=$((PATCH + 1))
NEXT_VERSION="${MAJOR}.${MINOR}.${NEXT_PATCH}-SNAPSHOT"

echo "Current Release: $MAJOR.$MINOR.$PATCH"
echo "Next Snapshot:   $NEXT_VERSION"

update_version_file "$TARGET_FILE" "$NEXT_VERSION"
write_github_output "next_version" "$NEXT_VERSION"

echo "Bump complete."