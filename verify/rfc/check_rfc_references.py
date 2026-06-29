#!/usr/bin/env python3
"""Phase 0 check: RFC cross-references resolve, and parenthetical title claims match.

Read-only. Scans ``rfcs/*.md`` for ``RFC-NNN`` references and reports:

  - broken_reference        (ERROR): RFC-NNN points to a number with no file.
  - reference_title_mismatch (ERROR): "RFC-NNN (Some Title)" names a title that does not
                              match the referenced file's heading title.

Both are structural errors and cause a non-zero exit. Ambiguities (e.g. Canonical Object vs
Artifact) are intentionally out of scope here and are tracked in
docs/ARCHITECTURE_DECISION_QUEUE.md.

Usage:
    python verify/rfc/check_rfc_references.py [--rfcs-dir DIR] [--json]
"""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path

try:
    from . import _common as c
except ImportError:  # standalone execution
    sys.path.insert(0, str(Path(__file__).resolve().parent))
    import _common as c  # type: ignore


NAME = "references"


def run_check(
    files_map: dict[str, Path],
    title_map: dict[str, str],
) -> dict:
    errors: list[dict] = []
    warnings: list[dict] = []

    for number, path in sorted(files_map.items()):
        rel = path.name
        for lineno, line in enumerate(c.read_text(path).splitlines(), start=1):
            for match in c.RFC_REF_RE.finditer(line):
                ref_num = match.group(1)
                paren = match.group(2)

                if ref_num not in files_map:
                    errors.append(
                        {
                            "code": "broken_reference",
                            "rfc": f"RFC-{number}",
                            "file": rel,
                            "line": lineno,
                            "reference": f"RFC-{ref_num}",
                            "message": f"RFC-{number} references RFC-{ref_num}, which has no file in rfcs/.",
                        }
                    )
                    continue

                if paren and c.is_probable_title(paren):
                    actual = title_map.get(ref_num, "")
                    if actual and c.title_key(paren) != c.title_key(actual):
                        errors.append(
                            {
                                "code": "reference_title_mismatch",
                                "rfc": f"RFC-{number}",
                                "file": rel,
                                "line": lineno,
                                "reference": f"RFC-{ref_num}",
                                "claimed_title": c.normalize_title(paren),
                                "actual_title": c.normalize_title(actual),
                                "message": (
                                    f"RFC-{number} refers to RFC-{ref_num} as "
                                    f"'{c.normalize_title(paren)}', but its title is "
                                    f"'{c.normalize_title(actual)}'."
                                ),
                            }
                        )

    return {
        "name": NAME,
        "ok": not errors,
        "errors": errors,
        "warnings": warnings,
        "details": {"files_scanned": len(files_map)},
    }


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--rfcs-dir", type=Path, default=None)
    parser.add_argument("--json", action="store_true", help="Print the raw JSON result.")
    args = parser.parse_args(argv)

    root = c.find_repo_root()
    rfcs_path = args.rfcs_dir or c.rfcs_dir(root)
    files_map = c.list_rfc_files(rfcs_path)
    title_map = c.first_h1_title_map(files_map)

    result = run_check(files_map, title_map)

    if args.json:
        print(json.dumps(result, indent=2))
    else:
        print(f"[{NAME}] scanned {len(files_map)} files: "
              f"{len(result['errors'])} error(s), {len(result['warnings'])} warning(s)")
        for err in result["errors"]:
            print(f"  ERROR {err['file']}:{err['line']} {err['message']}")

    return 0 if result["ok"] else 1


if __name__ == "__main__":
    raise SystemExit(main())
