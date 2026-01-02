#!/usr/bin/env bash
set -euo pipefail

# Spectr Validate with Binary
# This script wraps the spectr binary for validation

# Check if spectr binary is available
if ! command -v spectr >/dev/null 2>&1; then
    echo "Error: spectr binary not found in PATH" >&2
    echo "Please ensure spectr is installed and available in PATH" >&2
    exit 1
fi

# Pass all arguments to spectr validate
spectr validate "$@"