# verify/rfc — Phase 0 RFC Hygiene Tools

Read-only tooling that checks the structural health of the RFC corpus in [`rfcs/`](../../rfcs)
and produces machine-readable artifacts. This is **Phase 0 only**.

## Guarantees

- **Read-only.** These tools never write to or mutate any `rfcs/*.md` file.
- **Standard library only.** No external dependencies. No Neo4j, Graphify, MCP, or database.
- **No architecture decisions.** Ambiguities (e.g. Canonical Object vs Artifact) are reported as
  advisories, never resolved. Open decisions live in
  [`docs/ARCHITECTURE_DECISION_QUEUE.md`](../../docs/ARCHITECTURE_DECISION_QUEUE.md).

This implements the hygiene layer of RFC-060 (Testing Strategy) and follows the directory/manifest
spirit of RFC-061 (Verification Scripts). It is intentionally narrow.

## Run

```sh
python verify/rfc/run.py
```

Outputs (written to `verify/rfc/out/`):

- `report.json` — full machine-readable result.
- `report.md` — human-readable summary.
- `rfc_graph.json` — simple-JSON RFC dependency graph (`{nodes, edges}`).
- `invariants.json` — **draft, non-authoritative** invariant catalog (input to ADQ-005).

Each tool can also run standalone:

```sh
python verify/rfc/check_rfc_references.py --json
python verify/rfc/check_rfc_titles.py --json
python verify/rfc/check_duplicate_bodies.py --json
python verify/rfc/extract_invariants.py --out verify/rfc/out/invariants.json
python verify/rfc/extract_rfc_graph.py --out verify/rfc/out/rfc_graph.json
```

## What is checked

### Structural errors (exit code 1)

| Check | Code | Meaning |
|-------|------|---------|
| `check_duplicate_bodies.py` | `duplicate_body` | A file contains more than one `# RFC-NNN` heading. |
| `check_rfc_titles.py` | `missing_heading` | A file has no top-level RFC heading. |
| `check_rfc_titles.py` | `header_number_mismatch` | Heading RFC number ≠ filename number. |
| `check_rfc_titles.py` | `missing_file` | README lists an RFC number with no file. |
| `check_rfc_references.py` | `broken_reference` | `RFC-NNN` reference points to a nonexistent file. |
| `check_rfc_references.py` | `reference_title_mismatch` | `RFC-NNN (Title)` names the wrong title. |

### Warnings (exit code unaffected)

| Check | Code | Meaning |
|-------|------|---------|
| `check_rfc_titles.py` | `title_index_drift` | File heading title differs from the README index title. |
| `check_rfc_titles.py` | `not_in_index` | An RFC file is missing from the README numbering table. |

### Advisories (never fail)

The six open architecture decisions (ADQ-001…ADQ-006) are emitted as advisories so they stay
visible without blocking. They are documented — not resolved — in
[`docs/ARCHITECTURE_DECISION_QUEUE.md`](../../docs/ARCHITECTURE_DECISION_QUEUE.md).

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | No structural errors. |
| 1 | One or more structural errors. |

## Scope boundaries

**Safe (covered here):** reference resolution, title consistency, duplicate-body detection,
draft invariant extraction, simple RFC dependency graph.

**Unsafe (NOT here — blocked on ADQ decisions):** Canonical Object store, Artifact store,
review-target typing, Memory/Knowledge services, Workflow engine, Action persisted state machine,
the canonical invariant ID registry, and any edits to RFC files. These are tracked in
[`docs/ARCHITECTURE_DECISION_QUEUE.md`](../../docs/ARCHITECTURE_DECISION_QUEUE.md).

## Note on `_common.py`

`_common.py` is internal Phase 0 scaffolding shared by the checkers (RFC discovery, heading/title
parsing). It encodes no architectural decision.
