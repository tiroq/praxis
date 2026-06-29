#!/usr/bin/env python3
"""Phase 0 extractor: draft invariant catalog from the RFC corpus.

Read-only. Scans each ``rfcs/NNN-*.md`` for sections whose heading mentions "Invariants" or
"Acceptance Criteria" and extracts their list items into a draft, machine-readable catalog.

The output is DELIBERATELY non-authoritative. It exists only to inform ADQ-005 (Invariant ID
registry). The canonical invariant ID scheme is an open architecture decision and is NOT
established here. Draft IDs use the form ``RFC-NNN-INV-###`` / ``RFC-NNN-AC-###`` and may change
once ADQ-005 is decided.

Usage:
    python verify/rfc/extract_invariants.py [--rfcs-dir DIR] [--json] [--out FILE]
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


NAME = "invariants"

_SECTION_RE = re.compile(r"(invariants|acceptance criteria)", re.IGNORECASE)
_LIST_ITEM_RE = re.compile(r"^\s*(?:[-*+]|\d+[.)])\s+(.*\S)\s*$")


def _clean(text: str) -> str:
    text = text.replace("**", "").replace("`", "")
    return c.normalize_title(text)


def _section_kind(heading_text: str) -> str:
    return "acceptance_criteria" if "acceptance" in heading_text.lower() else "invariant"


def extract_from_text(number: str, text: str) -> list[dict]:
    items: list[dict] = []
    counters: dict[str, int] = {"invariant": 0, "acceptance_criteria": 0}

    lines = text.splitlines()
    capturing = False
    capture_level = 0
    capture_kind = ""
    capture_heading = ""

    for line in lines:
        heading = c.HEADING_RE.match(line)
        if heading:
            level = len(heading.group(1))
            heading_text = heading.group(2)
            if capturing and level <= capture_level:
                capturing = False  # left the captured section
            if _SECTION_RE.search(heading_text):
                capturing = True
                capture_level = level
                capture_kind = _section_kind(heading_text)
                capture_heading = c.normalize_title(heading_text)
            continue

        if not capturing:
            continue

        item = _LIST_ITEM_RE.match(line)
        if item:
            body = _clean(item.group(1))
            if not body:
                continue
            counters[capture_kind] += 1
            tag = "INV" if capture_kind == "invariant" else "AC"
            items.append(
                {
                    "id": f"RFC-{number}-{tag}-{counters[capture_kind]:03d}",
                    "rfc": f"RFC-{number}",
                    "kind": capture_kind,
                    "source_heading": capture_heading,
                    "text": body,
                }
            )

    return items


def run_extract(files_map: dict[str, Path]) -> dict:
    all_items: list[dict] = []
    per_rfc: dict[str, int] = {}
    for number, path in sorted(files_map.items()):
        items = extract_from_text(number, c.read_text(path))
        per_rfc[f"RFC-{number}"] = len(items)
        all_items.extend(items)

    invariant_count = sum(1 for i in all_items if i["kind"] == "invariant")
    ac_count = sum(1 for i in all_items if i["kind"] == "acceptance_criteria")

    return {
        "name": NAME,
        "authoritative": False,
        "note": (
            "DRAFT catalog only. Draft IDs are not the canonical invariant ID registry. "
            "The registry is an open decision tracked as ADQ-005 in "
            "docs/ARCHITECTURE_DECISION_QUEUE.md."
        ),
        "summary": {
            "rfcs": len(files_map),
            "invariants": invariant_count,
            "acceptance_criteria": ac_count,
            "total_items": len(all_items),
        },
        "per_rfc": per_rfc,
        "items": all_items,
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
        print(f"[{NAME}] wrote {result['summary']['total_items']} draft items to {args.out}")
    if args.json or not args.out:
        print(json.dumps(result, indent=2))

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
