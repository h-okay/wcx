#!/usr/bin/env python3

import json
import pathlib
import re
import statistics
import subprocess
import time


def parse_benchmark_output() -> tuple[str, str, str]:
    benchmark_text = pathlib.Path("benchmark.txt").read_text()
    match = re.search(
        r"^BenchmarkCountReader[^\n]*\s+(\d+)\s+([0-9.]+\s+\w+/op)\s+([0-9.]+\s+B/op)\s+(\d+\s+allocs/op)",
        benchmark_text,
        re.M,
    )
    if not match:
        raise SystemExit("could not parse benchmark output")

    ns_per_op = match.group(2)
    bytes_per_op = match.group(3)
    allocs_per_op = match.group(4)
    return ns_per_op, bytes_per_op, allocs_per_op


def parse_counts(text: str) -> list[int]:
    ints: list[int] = []
    for token in text.strip().split():
        if token.isdigit():
            ints.append(int(token))
        if len(ints) == 5:
            break
    if len(ints) != 5:
        raise ValueError(f"unable to parse counts from: {text!r}")
    return ints


def median_ms(command: list[str], runs: int = 20) -> str:
    samples: list[float] = []
    for _ in range(runs):
        start = time.perf_counter()
        subprocess.run(
            command, check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL
        )
        samples.append((time.perf_counter() - start) * 1000.0)
    return f"{statistics.median(samples):.2f}"


def verify_cli_output() -> tuple[str, str]:
    wcx_cmd = ["./wcx", "-l", "-w", "-m", "-c", "-L", "benchmark-input.txt"]
    gnu_cmd = ["wc", "-l", "-w", "-m", "-c", "-L", "benchmark-input.txt"]

    wcx_out = subprocess.check_output(wcx_cmd, text=True)
    gnu_out = subprocess.check_output(gnu_cmd, text=True)
    if parse_counts(wcx_out) != parse_counts(gnu_out):
        raise SystemExit("wcx output does not match GNU wc for benchmark input")

    return median_ms(wcx_cmd), median_ms(gnu_cmd)


def write_badge(ns_per_op: str) -> None:
    badge = {
        "schemaVersion": 1,
        "label": "benchmark",
        "message": ns_per_op,
        "color": "blue",
    }
    pathlib.Path(".github/badges").mkdir(parents=True, exist_ok=True)
    pathlib.Path(".github/badges/benchmark.json").write_text(
        json.dumps(badge, indent=2) + "\n"
    )


def write_readme(
    ns_per_op: str, bytes_per_op: str, allocs_per_op: str, wcx_ms: str, gnu_ms: str
) -> None:
    readme_path = pathlib.Path("README.md")
    readme = readme_path.read_text()
    start = "<!-- BENCHMARKS:START -->"
    end = "<!-- BENCHMARKS:END -->"

    body = (
        "### Go Micro-Benchmark\n"
        "| Benchmark | ns/op | B/op | allocs/op |\n"
        "| --- | ---: | ---: | ---: |\n"
        f"| `BenchmarkCountReader` | {ns_per_op} | {bytes_per_op} | {allocs_per_op} |\n\n"
        "### CLI Comparison (median of 20 runs)\n"
        "| Tool | ms/op |\n"
        "| --- | ---: | --- |\n"
        f"| `wcx -l -w -m -c -L benchmark-input.txt` | {wcx_ms} |\n"
        f"| `wc -l -w -m -c -L benchmark-input.txt` | {gnu_ms} |"
    )

    replacement = f"{start}\n{body}\n{end}"
    readme = re.sub(
        f"{re.escape(start)}.*?{re.escape(end)}", replacement, readme, flags=re.S
    )
    readme_path.write_text(readme)


def main() -> None:
    ns_per_op, bytes_per_op, allocs_per_op = parse_benchmark_output()
    wcx_ms, gnu_ms = verify_cli_output()
    write_badge(ns_per_op)
    write_readme(ns_per_op, bytes_per_op, allocs_per_op, wcx_ms, gnu_ms)


if __name__ == "__main__":
    main()
