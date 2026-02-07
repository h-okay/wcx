#!/usr/bin/env bash
set -euo pipefail

python3 - <<'PY'
from pathlib import Path

line = b"Lorem ipsum dolor sit amet, consectetur adipiscing elit.\n"
Path("benchmark-input.txt").write_bytes(line * 50_000)
PY

go test -run=^$ -bench=BenchmarkCountReader -benchmem ./internal/wc | tee benchmark.txt
go build -o wcx ./cmd/wcx
