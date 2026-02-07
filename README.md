# wcx

[![Tests](https://github.com/h-okay/wc-tool/actions/workflows/ci.yml/badge.svg)](https://github.com/h-okay/wc-tool/actions/workflows/ci.yml)
[![Release](https://github.com/h-okay/wc-tool/actions/workflows/release.yml/badge.svg)](https://github.com/h-okay/wc-tool/actions/workflows/release.yml)
[![Benchmark](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/h-okay/wc-tool/main/.github/badges/benchmark.json)](https://github.com/h-okay/wc-tool/actions/workflows/ci.yml)

`wcx` is a drop-in replacement for GNU `wc` with one extra extension: `--json`.

## What it supports

| Capability | GNU `wc` | `wcx` |
| --- | --- | --- |
| Core options `-c -m -l -L -w` | yes | yes |
| Multi-file output + totals | yes | yes |
| --files0-from=F | yes | yes |
| --total=`auto\|always\|only\|never` | yes | yes |
| --version | yes | yes |
| Stdin with no file args | yes | yes |
| Stdin via `-` file operand | yes | yes |
| JSON output (`--json`) | no | yes |

`--json` outputs machine-readable counts while preserving normal GNU behavior unless explicitly enabled.

Example:

```bash
./wcx --json internal/wc/testdata/test.txt
```

## Usage

```txt
wcx [OPTION]... [FILE]...
```

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

## Demo

![wcx demo](wcx-demo.gif)

## Benchmarks

The benchmark badge and table are updated automatically by the CI pipeline on `main`.

<!-- BENCHMARKS:START -->
### Go Micro-Benchmark
| Benchmark | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| `BenchmarkCountReader` | 79739857 ns/op | 65584 B/op | 2 allocs/op |

### CLI Comparison (median of 20 runs)
| Tool | ms/op |
| --- | ---: |
| `wcx -l -w -m -c -L benchmark-input.txt` | 83.43 |
| `wc -l -w -m -c -L benchmark-input.txt` | 10.68 |
<!-- BENCHMARKS:END -->
