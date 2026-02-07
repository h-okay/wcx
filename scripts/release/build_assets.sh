#!/usr/bin/env bash
set -euo pipefail

version="$1"

mkdir -p dist
for target in "linux amd64" "linux arm64" "darwin amd64" "darwin arm64"; do
  set -- $target
  GOOS="$1" GOARCH="$2" go build -trimpath -ldflags "-s -w -X main.version=${version}" -o "dist/wcx_${1}_${2}" ./cmd/wcx
done

shasum -a 256 dist/wcx_* > dist/checksums.txt
