#!/usr/bin/env python3
"""Phase 0 check: each RFC file contains exactly one RFC body.

Read-only. Counts top-level "# RFC-NNN ..." headings per file and reports:

  - duplicate_body (ERROR): a file contains more than one top-level RFC heading.

A duplicate body usually means two RFC drafts were concatenated into one file (e.g. RFC-052).
This is a structural error and causes a non-zero exit.

Usage:
    python verify/rfc/check_duplicate_bodies.py [--rfcs-dir DIR] [--json]
"""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path

try:
    from . import _common as c
except ImportError:
    sys.path.insert(0, str(Path(__file__).resolve().parent))
    import _common as c  # type: ignore


NAME = "duplicate_bodies"


def run_check(files_map: dict[str, Path]) -> dict:
    errors: list[dict] = []

    for number, path in sorted(files_map.items()):
        headings = c.find_h1_rfc_headings(c.read_text(path))
        if len(headings) > 1:
            lines = [ln for ln, _, _ in headings]
            titles = [f"RFC-{num} {title}".strip() for _, num, title in headings]
            errors.append(
                {
                    "code": "duplicate_body",
                    "rfc": f"RFC-{number}",
                    "file": path.name,
                    "heading_lines": lines,
                    "headings": titles,
                    "message": (
                        f"{path.name} contains {len(headings)} top-level RFC headings "
                        f"(lines {lines}); expected exactly one."
                    ),
                }
            )

    return {
        "name": NAME,
        "ok": not errors,
        "errors": errors,
        "warnings": [],
        "details": {"files": len(files_map)},
    }


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--rfcs-dir", type=Path, default=None)
    parser.add_argument("--json", action="store_true")
    args = parser.parse_args(argv)

    root = c.find_repo_root()
    rfcs_path = args.rfcs_dir or c.rfcs_dir(root)
    files_map = c.list_rfc_files(rfcs_path)

    result = run_check(files_map)

    if args.json:
        print(json.dumps(result, indent=2))
    else:
        print(f"[{NAME}] {len(result['errors'])} error(s)")
        for err in result["errors"]:
            print(f"  ERROR {err['file']} {err['message']}")

    return 0 if result["ok"] else 1


if __name__ == "__main__":
    raise SystemExit(main())
