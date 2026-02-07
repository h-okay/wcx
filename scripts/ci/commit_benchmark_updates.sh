#!/usr/bin/env bash
set -euo pipefail

if git diff --quiet README.md .github/badges/benchmark.json; then
  exit 0
fi

git config user.name "github-actions[bot]"
git config user.email "41898282+github-actions[bot]@users.noreply.github.com"
git add README.md .github/badges/benchmark.json
git commit -m "chore: update benchmark badge and README [skip ci]"
git push
