#!/bin/bash
# Copyright (c) 2024-2025 Six After, Inc.
#
# This source code is licensed under the Apache 2.0 License found in the
# LICENSE file in the root directory of this source tree.

# Verify the integrity of the latest release using Cosign + checksums
# Works on macOS and Linux

set -euo pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${__dir}"/os-type.sh

# Windows
if is_windows; then
    echo "[ERROR] Windows is not currently supported." >&2
    exit 1
fi

rm -fr dist
goreleaser --config .goreleaser.yaml release --snapshot