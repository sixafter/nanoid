#!/bin/bash
# Copyright (c) 2024-2025 Six After, Inc.
#
# This source code is licensed under the Apache 2.0 License found in the
# LICENSE file in the root directory of this source tree.

set -e

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${__dir}"/os-type.sh

# ---------------------------------------------------------------------------
# Platform Guard
# ---------------------------------------------------------------------------
if is_windows; then
    echo "[ERROR] Windows is not currently supported." >&2
    exit 1
fi

# ---------------------------------------------------------------------------
# Setup
# ---------------------------------------------------------------------------
mkdir -p tmp
rm -f tmp/*.zip 2>/dev/null || true

REPO_OWNER="sixafter"
REPO_NAME="nanoid"
MODULE="github.com/${REPO_OWNER}/${REPO_NAME}"

# ---------------------------------------------------------------------------
# Select TAG
#   If TAG env var is provided, use it.
#   Otherwise detect latest release from GitHub.
# ---------------------------------------------------------------------------
if [ -n "${TAG:-}" ]; then
    echo "Using provided TAG: ${TAG}"
else
    echo "No TAG provided — detecting latest GitHub release..."
    TAG=$(curl -s "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" | jq -r .tag_name)

    if [ -z "$TAG" ] || [ "$TAG" = "null" ]; then
        echo "[ERROR] Could not detect latest release tag from GitHub." >&2
        exit 1
    fi
fi

VERSION=${TAG#v}
echo "Verifying module version: ${TAG} (parsed version: ${VERSION})"
echo

# ---------------------------------------------------------------------------
# Portable SHA-256 (Linux/macOS)
# ---------------------------------------------------------------------------
if command -v sha256sum >/dev/null 2>&1; then
  SHA256="sha256sum"
else
  SHA256="shasum -a 256"
fi

# ---------------------------------------------------------------------------
# 1. GitHub UI Tag Archive
#    This is the human-facing ZIP ("Source code (zip)")
#    Note: This archive is NOT used by the Go module system.
# ---------------------------------------------------------------------------
echo "STEP 1: Downloading GitHub UI ZIP (human-facing tag archive)..."
echo "        URL: https://github.com/${REPO_OWNER}/${REPO_NAME}/archive/refs/tags/${TAG}.zip"
curl -sSfL -o tmp/github-ui.zip \
  "https://github.com/${REPO_OWNER}/${REPO_NAME}/archive/refs/tags/${TAG}.zip"

GITHUB_UI_SHA=$($SHA256 tmp/github-ui.zip | awk '{print $1}')
echo "GitHub UI ZIP SHA256:       ${GITHUB_UI_SHA}"
echo

# ---------------------------------------------------------------------------
# 2. Direct Module ZIP
#    The EXACT bytes Go would fetch using: GOPROXY=direct
#    This downloads: https://github.com/<repo>/@v/<tag>.zip
# ---------------------------------------------------------------------------
echo "STEP 2: Downloading DIRECT module ZIP via go mod (GOPROXY=direct)..."
MOD_JSON=$(GOPROXY=direct go mod download -json "${MODULE}@${TAG}")
MOD_ZIP_PATH=$(echo "$MOD_JSON" | jq -r '.Zip')

if [ ! -f "$MOD_ZIP_PATH" ]; then
  echo "[ERROR] go mod direct ZIP not found at:"
  echo "$MOD_ZIP_PATH"
  exit 1
fi

cp "$MOD_ZIP_PATH" tmp/direct.zip
DIRECT_SHA=$($SHA256 tmp/direct.zip | awk '{print $1}')
echo "Direct module ZIP SHA256:   ${DIRECT_SHA}"
echo

# ---------------------------------------------------------------------------
# 3. Go Proxy Module ZIP
#    This is what proxy.golang.org serves and what Go’s checksum database uses.
# ---------------------------------------------------------------------------
echo "STEP 3: Downloading Go PROXY module ZIP..."
echo "        URL: https://proxy.golang.org/${MODULE}/@v/${TAG}.zip"
curl -sSfL -o tmp/proxy.zip \
  "https://proxy.golang.org/${MODULE}/@v/${TAG}.zip"

PROXY_SHA=$($SHA256 tmp/proxy.zip | awk '{print $1}')
echo "Go PROXY module ZIP SHA256: ${PROXY_SHA}"
echo

# ---------------------------------------------------------------------------
# Comparison
# ---------------------------------------------------------------------------
echo "=============================================================="
echo "CHECKSUM COMPARISON"
echo "=============================================================="
printf "GitHub UI ZIP (human-facing)......: %s\n" "$GITHUB_UI_SHA"
printf "Direct module ZIP (GOPROXY=direct): %s\n" "$DIRECT_SHA"
printf "Go PROXY ZIP (checksum DB)........: %s\n" "$PROXY_SHA"
echo

# ---------------------------------------------------------------------------
# Authoritative Integrity Check
# DIRECT vs PROXY must match for module reproducibility.
# ---------------------------------------------------------------------------
if [ "$DIRECT_SHA" != "$PROXY_SHA" ]; then
  echo "[ERROR] ❌ DIRECT and PROXY module ZIPs DO NOT MATCH!"
  echo "        This version is NOT reproducible across environments."
  exit 1
fi

echo "✔ DIRECT and PROXY module ZIPs MATCH (authoritative)."
echo

# ---------------------------------------------------------------------------
# Informational Check: GitHub UI ZIP (not used by Go)
# ---------------------------------------------------------------------------
if [ "$GITHUB_UI_SHA" != "$DIRECT_SHA" ]; then
  echo "⚠ INFO: GitHub UI ZIP does NOT match the module ZIP."
  echo "        This is normal — GitHub generates UI tag archives separately."
else
  echo "✔ GitHub UI ZIP matches module ZIP (rare but valid)."
fi

echo
echo "✔ Module ${TAG} is fully reproducible across GOPROXY=direct and proxy.golang.org."
echo "   (Human-facing GitHub UI ZIP differences are expected and non-breaking.)"
