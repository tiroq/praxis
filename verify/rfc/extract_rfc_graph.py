#!/usr/bin/env python3
"""Phase 0 extractor: simple-JSON RFC dependency graph.

Read-only. Builds a plain JSON graph of RFC dependency relationships parsed from the
"Depends on" / "Required before" cues in each ``rfcs/NNN-*.md``.

Deliberately simple: no Neo4j, no Graphify, no external libraries, no database. The output is a
flat ``{nodes, edges}`` document intended as a first, inspectable artifact. Richer graph backends
are explicitly out of Phase 0 scope.

Edge relations:
  - depends_on      : "Depends on ... RFC-NNN" (and "RFC-AAA through RFC-BBB" ranges)
  - required_before : "Required before ... RFC-NNN"

Usage:
    python verify/rfc/extract_rfc_graph.py [--rfcs-dir DIR] [--json] [--out FILE]
"""

from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path

try:
    from . import _common as c
except ImportError:
    sys.path.insert(0, str(Path(__file__).resolve().parent))
    import _common as c  # type: ignore


NAME = "rfc_graph"

_DEPENDS_RE = re.compile(r"depends on", re.IGNORECASE)
_REQUIRED_BEFORE_RE = re.compile(r"required before", re.IGNORECASE)
_RANGE_RE = re.compile(r"RFC-(\d{3})\s+through\s+RFC-(\d{3})", re.IGNORECASE)


def _refs_in_line(line: str) -> list[str]:
    return [m.group(1) for m in c.RFC_REF_RE.finditer(line)]


def _expand_range(a: str, b: str, existing: set[str]) -> list[str]:
    lo, hi = sorted((int(a), int(b)))
    return [f"{n:03d}" for n in range(lo, hi + 1) if f"{n:03d}" in existing]


def run_extract(files_map: dict[str, Path]) -> dict:
    existing = set(files_map)
    title_map = c.first_h1_title_map(files_map)

    nodes = [
        {
            "id": f"RFC-{number}",
            "number": number,
            "title": c.normalize_title(title_map.get(number, "")),
            "file": files_map[number].name,
        }
        for number in sorted(files_map)
    ]

    edge_set: set[tuple[str, str, str]] = set()

    for number, path in sorted(files_map.items()):
        source = f"RFC-{number}"
        relation: str | None = None
        for raw in c.read_text(path).splitlines():
            line = raw.strip()

            # Reset scope on blank lines and headings to avoid bleeding across sections.
            if not line or c.HEADING_RE.match(raw):
                relation = None

            if _DEPENDS_RE.search(line):
                relation = "depends_on"
            elif _REQUIRED_BEFORE_RE.search(line):
                relation = "required_before"

            if relation is None:
                continue

            targets: list[str] = []
            for ra, rb in _RANGE_RE.findall(line):
                targets.extend(_expand_range(ra, rb, existing))
            targets.extend(_refs_in_line(line))

            for tnum in targets:
                if tnum not in existing:
                    continue
                target = f"RFC-{tnum}"
                if target == source:
                    continue
                edge_set.add((source, target, relation))

    edges = [
        {"source": s, "target": t, "relation": r}
        for (s, t, r) in sorted(edge_set)
    ]

    return {
        "name": NAME,
        "format": "simple-json-v1",
        "note": (
            "Phase 0 dependency graph parsed from 'Depends on' / 'Required before' cues. "
            "No external graph backend; this is intentionally a flat nodes/edges document."
        ),
        "summary": {"nodes": len(nodes), "edges": len(edges)},
        "nodes": nodes,
        "edges": edges,
    }


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--rfcs-dir", type=Path, default=None)
    parser.add_argument("--out", type=Path, default=None, help="Write JSON to this path.")
    parser.add_argument("--json", action="store_true", help="Print JSON to stdout.")
    args = parser.parse_args(argv)

    root = c.find_repo_root()
    rfcs_path = args.rfcs_dir or c.rfcs_dir(root)
    files_map = c.list_rfc_files(rfcs_path)

    result = run_extract(files_map)

    if args.out:
        args.out.parent.mkdir(parents=True, exist_ok=True)
        args.out.write_text(json.dumps(result, indent=2) + "\n", encoding="utf-8")
        print(f"[{NAME}] wrote {result['summary']['nodes']} nodes / "
              f"{result['summary']['edges']} edges to {args.out}")
    if args.json or not args.out:
        print(json.dumps(result, indent=2))

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
