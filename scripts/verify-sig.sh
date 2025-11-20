#!/usr/bin/env bash
# Copyright (c) 2024-2025 Six After, Inc.
#
# This source code is licensed under the Apache 2.0 License found in the
# LICENSE file in the root directory of this source tree.

set -euo pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${__dir}/os-type.sh"

# Windows
if is_windows; then
    echo "[ERROR] Windows is not currently supported." >&2
    exit 1
fi

# ------------------------------------------------------------
# Project / repository name (portable)
# ------------------------------------------------------------
PROJECT="nanoid"
REPO="sixafter/${PROJECT}"
MODULE="github.com/${REPO}"

# tmp directory for artifacts
TMP="${__dir}/tmp"
mkdir -p "${TMP}"

echo "Project: ${PROJECT}"
echo "Repository: ${REPO}"
echo "Module path: ${MODULE}"
echo "Artifact directory: ${TMP}"
echo

# ------------------------------------------------------------
# Detect latest release
# ------------------------------------------------------------
TAG=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | jq -r .tag_name)
VERSION=${TAG#v}

echo "Latest release: ${TAG} (version: ${VERSION})"

# ------------------------------------------------------------
# Determine SHA-256 tool
# ------------------------------------------------------------
if command -v sha256sum >/dev/null 2>&1; then
  SHA256="sha256sum"
else
  SHA256="shasum -a 256"
fi

# ------------------------------------------------------------
# Download release artifacts → tmp/
# ------------------------------------------------------------
echo
echo "Downloading release artifacts into ${TMP}..."

# Core tarball
curl -sSfL -o "${TMP}/${PROJECT}-${VERSION}.tar.gz" \
  "https://github.com/${REPO}/releases/download/${TAG}/${PROJECT}-${VERSION}.tar.gz"

# Tarball signature
curl -sSfL -o "${TMP}/${PROJECT}-${VERSION}.tar.gz.sigstore.json" \
  "https://github.com/${REPO}/releases/download/${TAG}/${PROJECT}-${VERSION}.tar.gz.sigstore.json"

# SBOM (required for checksum verification)
curl -sSfL -o "${TMP}/${PROJECT}-${VERSION}.tar.gz.sbom.json" \
  "https://github.com/${REPO}/releases/download/${TAG}/${PROJECT}-${VERSION}.tar.gz.sbom.json"

# checksums.txt + signature
curl -sSfL -o "${TMP}/checksums.txt" \
  "https://github.com/${REPO}/releases/download/${TAG}/checksums.txt"

curl -sSfL -o "${TMP}/checksums.txt.sigstore.json" \
  "https://github.com/${REPO}/releases/download/${TAG}/checksums.txt.sigstore.json"

# ------------------------------------------------------------
# Verify tarball with Cosign
# ------------------------------------------------------------
echo
echo "Verifying tarball signature..."

cosign verify-blob \
  --key "${__dir}/../cosign.pub" \
  --bundle "${TMP}/${PROJECT}-${VERSION}.tar.gz.sigstore.json" \
  "${TMP}/${PROJECT}-${VERSION}.tar.gz"

echo "Tarball signature OK."

# ------------------------------------------------------------
# Verify checksums manifest signature
# ------------------------------------------------------------
echo
echo "Verifying checksums.txt signature..."

cosign verify-blob \
  --key "${__dir}/../cosign.pub" \
  --bundle "${TMP}/checksums.txt.sigstore.json" \
  "${TMP}/checksums.txt"

echo "Checksums signature OK."

# ------------------------------------------------------------
# Validate local artifact integrity
# ------------------------------------------------------------
echo
echo "Verifying file checksums locally..."
(
  cd "${TMP}"
  $SHA256 -c checksums.txt
) || {
    echo
    echo "❌ Release verification FAILED."
    exit 1
}

echo
echo "✔ Release verification succeeded."
