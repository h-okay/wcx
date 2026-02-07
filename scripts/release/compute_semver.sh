#!/usr/bin/env bash
set -euo pipefail

latest_tag="$(git tag -l 'v[0-9]*.[0-9]*.[0-9]*' --sort=-v:refname | head -n 1)"
if [ -z "$latest_tag" ]; then
  latest_tag="v0.0.0"
  range=""
else
  range="$latest_tag..HEAD"
fi

if [ -n "$range" ]; then
  subjects="$(git log --format=%s "$range")"
  bodies="$(git log --format=%b "$range")"
else
  subjects="$(git log --format=%s)"
  bodies="$(git log --format=%b)"
fi

if [ -z "$subjects" ]; then
  echo "should_release=false" >> "$GITHUB_OUTPUT"
  echo "version=$latest_tag" >> "$GITHUB_OUTPUT"
  exit 0
fi

bump="patch"
if printf '%s\n%s\n' "$subjects" "$bodies" | grep -Eq '(BREAKING CHANGE|^[a-zA-Z]+(\(.+\))?!:)'; then
  bump="major"
elif printf '%s\n' "$subjects" | grep -Eq '^feat(\(.+\))?:'; then
  bump="minor"
fi

version="${latest_tag#v}"
IFS='.' read -r major minor patch <<< "$version"

case "$bump" in
  major)
    major=$((major + 1))
    minor=0
    patch=0
    ;;
  minor)
    minor=$((minor + 1))
    patch=0
    ;;
  patch)
    patch=$((patch + 1))
    ;;
esac

next="v${major}.${minor}.${patch}"
echo "should_release=true" >> "$GITHUB_OUTPUT"
echo "version=$next" >> "$GITHUB_OUTPUT"
