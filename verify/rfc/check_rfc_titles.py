#!/usr/bin/env python3
"""Phase 0 check: RFC file titles are self-consistent and match the README index.

Read-only. For each ``rfcs/NNN-*.md`` reports:

  - missing_heading       (ERROR): the file has no top-level "# RFC-NNN ..." heading.
  - header_number_mismatch (ERROR): the heading's RFC number != the filename number.
  - missing_file          (ERROR): a number listed in the README numbering table has no file.
  - title_index_drift     (WARNING): the file heading title differs from the README title.
  - not_in_index          (WARNING): an RFC file is absent from the README numbering table.

Only the ERROR classes affect the exit code. README drift is a non-blocking warning because the
README index is a secondary listing, not the authoritative title source.

Usage:
    python verify/rfc/check_rfc_titles.py [--rfcs-dir DIR] [--readme FILE] [--json]
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


NAME = "titles"


def run_check(
    files_map: dict[str, Path],
    readme_titles: dict[str, str],
) -> dict:
    errors: list[dict] = []
    warnings: list[dict] = []

    for number, path in sorted(files_map.items()):
        rel = path.name
        headings = c.find_h1_rfc_headings(c.read_text(path))

        if not headings:
            errors.append(
                {
                    "code": "missing_heading",
                    "rfc": f"RFC-{number}",
                    "file": rel,
                    "message": f"{rel} has no top-level '# RFC-{number} ...' heading.",
                }
            )
            continue

        first_lineno, first_num, first_title = headings[0]

        if first_num != number:
            errors.append(
                {
                    "code": "header_number_mismatch",
                    "rfc": f"RFC-{number}",
                    "file": rel,
                    "line": first_lineno,
                    "message": (
                        f"{rel} filename is RFC-{number} but its first heading is RFC-{first_num}."
                    ),
                }
            )

        index_title = readme_titles.get(number)
        if index_title and c.title_key(index_title) != c.title_key(first_title):
            warnings.append(
                {
                    "code": "title_index_drift",
                    "rfc": f"RFC-{number}",
                    "file": rel,
                    "line": first_lineno,
                    "heading_title": c.normalize_title(first_title),
                    "index_title": c.normalize_title(index_title),
                    "message": (
                        f"RFC-{number} heading title '{c.normalize_title(first_title)}' differs "
                        f"from README index title '{c.normalize_title(index_title)}'."
                    ),
                }
            )

    # README numbering completeness.
    for number in sorted(readme_titles):
        if number not in files_map:
            errors.append(
                {
                    "code": "missing_file",
                    "rfc": f"RFC-{number}",
                    "file": None,
                    "message": f"README lists RFC-{number} but no rfcs/{number}-*.md exists.",
                }
            )
    for number in sorted(files_map):
        if readme_titles and number not in readme_titles:
            warnings.append(
                {
                    "code": "not_in_index",
                    "rfc": f"RFC-{number}",
                    "file": files_map[number].name,
                    "message": f"RFC-{number} is not listed in the README numbering table.",
                }
            )

    return {
        "name": NAME,
        "ok": not errors,
        "errors": errors,
        "warnings": warnings,
        "details": {"files": len(files_map), "index_entries": len(readme_titles)},
    }


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--rfcs-dir", type=Path, default=None)
    parser.add_argument("--readme", type=Path, default=None)
    parser.add_argument("--json", action="store_true")
    args = parser.parse_args(argv)

    root = c.find_repo_root()
    rfcs_path = args.rfcs_dir or c.rfcs_dir(root)
    readme_path = args.readme or (rfcs_path / "README.md")

    files_map = c.list_rfc_files(rfcs_path)
    readme_titles = c.parse_readme_titles(readme_path)

    result = run_check(files_map, readme_titles)

    if args.json:
        print(json.dumps(result, indent=2))
    else:
        print(f"[{NAME}] {len(result['errors'])} error(s), {len(result['warnings'])} warning(s)")
        for err in result["errors"]:
            loc = err.get("file") or "-"
            print(f"  ERROR {loc} {err['message']}")
        for warn in result["warnings"]:
            print(f"  WARN  {warn.get('file') or '-'} {warn['message']}")

    return 0 if result["ok"] else 1


if __name__ == "__main__":
    raise SystemExit(main())
