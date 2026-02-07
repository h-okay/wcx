# wcx

[![Tests](https://github.com/h-okay/wc-tool/actions/workflows/ci.yml/badge.svg)](https://github.com/h-okay/wc-tool/actions/workflows/ci.yml)
[![Release](https://github.com/h-okay/wc-tool/actions/workflows/release.yml/badge.svg)](https://github.com/h-okay/wc-tool/actions/workflows/release.yml)
[![Benchmark](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/h-okay/wc-tool/main/.github/badges/benchmark.json)](https://github.com/h-okay/wc-tool/actions/workflows/ci.yml)

`wcx` is a drop-in replacement for GNU `wc` with one extra extension: `--json`.

## What it supports

- GNU-compatible core options: `-c`, `-m`, `-l`, `-L`, `-w`
- GNU-compatible multi-file behavior and totals
- `--files0-from=F`
- `--total=auto|always|only|never`
- `--version`
- zero third-party dependencies (stdlib only)
- CI runs unit tests and short fuzz smoke tests on PRs and `main`
- Stdin behavior:
  - no file args -> read stdin
  - `-` file operand -> read stdin at that position

## wcx extension

- `--json` outputs machine-readable counts while preserving normal GNU behavior unless explicitly enabled.

Example:

```bash
./wcx --json internal/wc/testdata/test.txt
```

## Usage

```txt
wcx [OPTION]... [FILE]...
```

Default output matches GNU `wc`: `lines words bytes`.

If one or more count options are provided, output order is always:

`lines words chars bytes max-line-length`

## Examples

```bash
# default counts (lines, words, bytes)
./wcx internal/wc/testdata/test.txt

# chars only
./wcx -m internal/wc/testdata/test.txt

# max line length only
./wcx -L internal/wc/testdata/test.txt

# multiple files + automatic total
./wcx internal/wc/testdata/test.txt internal/wc/testdata/test.txt

# files listed in a NUL-delimited file list
./wcx --files0-from=filelist.txt

# only print total counts
./wcx --total=only internal/wc/testdata/test.txt internal/wc/testdata/test.txt
```

## Build

```bash
go build -o wcx ./cmd/wcx
```

## Benchmarks

The benchmark badge and table are updated automatically by the CI pipeline on `main`.

<!-- BENCHMARKS:START -->
### Go Micro-Benchmark
| Benchmark | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| `BenchmarkCountReader` | _pending CI run_ | _pending_ | _pending_ |

### CLI Comparison (median of 20 runs)
| Tool | ms/op | Notes |
| --- | ---: | --- |
| `wcx -l -w -m -c -L benchmark-input.txt` | _pending_ | _pending_ |
| `wc -l -w -m -c -L benchmark-input.txt` | _pending_ | GNU reference |
<!-- BENCHMARKS:END -->
