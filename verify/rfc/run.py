#!/usr/bin/env python3
"""Phase 0 RFC hygiene runner.

Read-only. Aggregates the structural checks and the extractors, then writes:

    verify/rfc/out/report.json
    verify/rfc/out/report.md
    verify/rfc/out/rfc_graph.json
    verify/rfc/out/invariants.json

Exit code:
    0  no structural errors
    1  one or more structural errors (duplicate bodies, missing files, broken references,
       reference title mismatch)

Warnings (e.g. README title drift) and architecture-decision advisories (ADQ-001..006) NEVER
affect the exit code. This runner does not resolve any architecture decision and does not modify
RFC files.

Usage:
    python verify/rfc/run.py [--rfcs-dir DIR] [--out-dir DIR]
"""

from __future__ import annotations

import argparse
import datetime as _dt
import json
import sys
from pathlib import Path

try:
    from . import _common as c
    from . import check_duplicate_bodies, check_rfc_references, check_rfc_titles
    from . import extract_invariants, extract_rfc_graph
except ImportError:
    sys.path.insert(0, str(Path(__file__).resolve().parent))
    import _common as c  # type: ignore
    import check_duplicate_bodies  # type: ignore
    import check_rfc_references  # type: ignore
    import check_rfc_titles  # type: ignore
    import extract_invariants  # type: ignore
    import extract_rfc_graph  # type: ignore


# Open architecture decisions surfaced as advisories (never failures). These mirror
# docs/ARCHITECTURE_DECISION_QUEUE.md and are intentionally NOT resolved here.
ADQ_ADVISORIES = [
    ("ADQ-001", "Canonical Object vs Artifact model is unresolved."),
    ("ADQ-002", "Memory / Knowledge service ownership is undefined in RFC-030/031."),
    ("ADQ-003", "Workflow model has no RFC; deferred for Phases 1-2."),
    ("ADQ-004", "Action state machine differs between RFC-022 and RFC-023."),
    ("ADQ-005", "Invariant ID registry referenced by RFC-061 is undefined in RFC-060."),
    ("ADQ-006", "Some Space RFC storage mappings contradict RFC-033."),
]
ADQ_DOC = "docs/ARCHITECTURE_DECISION_QUEUE.md"


def _now_iso() -> str:
    return _dt.datetime.now(_dt.timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def _md_table(headers: list[str], rows: list[list[str]]) -> str:
    out = ["| " + " | ".join(headers) + " |",
           "| " + " | ".join("---" for _ in headers) + " |"]
    for row in rows:
        out.append("| " + " | ".join(row) + " |")
    return "\n".join(out)


def build_markdown(report: dict) -> str:
    s = report["summary"]
    status = "PASS" if s["errors"] == 0 else "FAIL"
    lines: list[str] = []
    lines.append("# Praxis RFC Hygiene Report")
    lines.append("")
    lines.append(f"- Generated: {report['generated_at']}")
    lines.append(f"- Status: **{status}** ({s['errors']} error(s), {s['warnings']} warning(s))")
    lines.append(f"- RFC files scanned: {s['rfc_files']}")
    lines.append("")
    lines.append("> Read-only Phase 0 tooling. It does not modify RFC files and does not resolve")
    lines.append(f"> any architecture decision. See [{ADQ_DOC}](../../../{ADQ_DOC}).")
    lines.append("")

    lines.append("## Check Summary")
    lines.append("")
    rows = []
    for name, chk in report["checks"].items():
        rows.append([
            name,
            str(len(chk["errors"])),
            str(len(chk["warnings"])),
            "PASS" if chk["ok"] else "FAIL",
        ])
    lines.append(_md_table(["Check", "Errors", "Warnings", "Status"], rows))
    lines.append("")

    # Errors
    lines.append("## Structural Errors (cause non-zero exit)")
    lines.append("")
    errors = [e for chk in report["checks"].values() for e in chk["errors"]]
    if not errors:
        lines.append("_None._")
    else:
        for e in errors:
            loc = e.get("file") or "-"
            line = f":{e['line']}" if e.get("line") else ""
            lines.append(f"- **{e['code']}** `{loc}{line}` — {e['message']}")
    lines.append("")

    # Warnings
    lines.append("## Warnings (non-blocking)")
    lines.append("")
    warnings = [w for chk in report["checks"].values() for w in chk["warnings"]]
    if not warnings:
        lines.append("_None._")
    else:
        for w in warnings:
            loc = w.get("file") or "-"
            line = f":{w['line']}" if w.get("line") else ""
            lines.append(f"- **{w['code']}** `{loc}{line}` — {w['message']}")
    lines.append("")

    # Advisories
    lines.append("## Architecture Decision Advisories (open — not resolved)")
    lines.append("")
    lines.append(f"These ambiguities are tracked in [{ADQ_DOC}](../../../{ADQ_DOC}) and are "
                 "reported as advisories only. They never fail the build.")
    lines.append("")
    adv_rows = [[a["id"], a["status"], a["summary"]] for a in report["advisories"]]
    lines.append(_md_table(["ID", "Status", "Summary"], adv_rows))
    lines.append("")

    # Artifacts
    lines.append("## Artifacts")
    lines.append("")
    art = report["artifacts"]
    lines.append(f"- `{art['rfc_graph']['path']}` — "
                 f"{art['rfc_graph']['nodes']} nodes, {art['rfc_graph']['edges']} edges")
    lines.append(f"- `{art['invariants']['path']}` — "
                 f"{art['invariants']['total_items']} draft items "
                 "(non-authoritative; see ADQ-005)")
    lines.append("")
    return "\n".join(lines) + "\n"


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--rfcs-dir", type=Path, default=None)
    parser.add_argument("--out-dir", type=Path, default=None)
    args = parser.parse_args(argv)

    root = c.find_repo_root()
    rfcs_path = args.rfcs_dir or c.rfcs_dir(root)
    out_dir = args.out_dir or (Path(__file__).resolve().parent / "out")
    out_dir.mkdir(parents=True, exist_ok=True)

    files_map = c.list_rfc_files(rfcs_path)
    title_map = c.first_h1_title_map(files_map)
    readme_titles = c.parse_readme_titles(rfcs_path / "README.md")

    # Structural checks.
    checks = {
        "duplicate_bodies": check_duplicate_bodies.run_check(files_map),
        "titles": check_rfc_titles.run_check(files_map, readme_titles),
        "references": check_rfc_references.run_check(files_map, title_map),
    }

    # Extractors (artifacts).
    graph = extract_rfc_graph.run_extract(files_map)
    invariants = extract_invariants.run_extract(files_map)

    graph_path = out_dir / "rfc_graph.json"
    inv_path = out_dir / "invariants.json"
    graph_path.write_text(json.dumps(graph, indent=2) + "\n", encoding="utf-8")
    inv_path.write_text(json.dumps(invariants, indent=2) + "\n", encoding="utf-8")

    total_errors = sum(len(chk["errors"]) for chk in checks.values())
    total_warnings = sum(len(chk["warnings"]) for chk in checks.values())

    def rel(p: Path) -> str:
        try:
            return str(p.relative_to(root))
        except ValueError:
            return str(p)

    report = {
        "tool": "verify/rfc",
        "version": "0.1.0",
        "phase": 0,
        "generated_at": _now_iso(),
        "repo_root": str(root),
        "read_only": True,
        "resolves_architecture_decisions": False,
        "summary": {
            "rfc_files": len(files_map),
            "errors": total_errors,
            "warnings": total_warnings,
            "status": "PASS" if total_errors == 0 else "FAIL",
        },
        "checks": checks,
        "advisories": [
            {"id": adq_id, "status": "open", "severity": "warning",
             "summary": summary, "see": ADQ_DOC}
            for adq_id, summary in ADQ_ADVISORIES
        ],
        "artifacts": {
            "rfc_graph": {
                "path": rel(graph_path),
                "nodes": graph["summary"]["nodes"],
                "edges": graph["summary"]["edges"],
            },
            "invariants": {
                "path": rel(inv_path),
                "total_items": invariants["summary"]["total_items"],
                "authoritative": False,
            },
        },
    }

    report_json = out_dir / "report.json"
    report_md = out_dir / "report.md"
    report_json.write_text(json.dumps(report, indent=2) + "\n", encoding="utf-8")
    report_md.write_text(build_markdown(report), encoding="utf-8")

    status = report["summary"]["status"]
    print(f"RFC hygiene: {status} — {total_errors} error(s), {total_warnings} warning(s)")
    print(f"  report.json   -> {rel(report_json)}")
    print(f"  report.md     -> {rel(report_md)}")
    print(f"  rfc_graph.json-> {rel(graph_path)} "
          f"({graph['summary']['nodes']} nodes, {graph['summary']['edges']} edges)")
    print(f"  invariants.json-> {rel(inv_path)} "
          f"({invariants['summary']['total_items']} draft items)")
    if total_errors:
        print("Structural errors (see report for detail):")
        for chk in checks.values():
            for e in chk["errors"]:
                loc = e.get("file") or "-"
                line = f":{e['line']}" if e.get("line") else ""
                print(f"  ERROR [{e['code']}] {loc}{line} {e['message']}")

    return 0 if total_errors == 0 else 1


if __name__ == "__main__":
    raise SystemExit(main())
