#!/usr/bin/env bash
set -euo pipefail

version="$1"

if gh release view "$version" >/dev/null 2>&1; then
  gh release upload "$version" dist/wcx_* dist/checksums.txt --clobber
else
  gh release create "$version" dist/wcx_* dist/checksums.txt --title "$version" --generate-notes
fi
