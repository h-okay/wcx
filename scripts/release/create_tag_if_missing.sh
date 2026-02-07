#!/usr/bin/env bash
set -euo pipefail

version="$1"

if [ -z "$(git ls-remote --tags origin "refs/tags/${version}")" ]; then
  git tag "$version"
  git push origin "refs/tags/${version}"
fi
