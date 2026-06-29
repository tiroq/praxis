# Phase 0 — RFC Cleanup Plan

> Source of truth: the **Corrected Praxis RFC Audit** and
> [ARCHITECTURE_DECISION_QUEUE.md](./ARCHITECTURE_DECISION_QUEUE.md).
>
> Phase 0 is **hygiene only**. It does **not** resolve any ADQ item, does **not** change
> Canonical Object / Artifact semantics, and does **not** implement any domain model,
> store, service, workflow engine, or Action state machine.
>
> The tooling in [`verify/rfc/`](../verify/rfc/README.md) is **read-only** and never edits RFC files.

## Goal

Make the RFC corpus structurally trustworthy before any implementation:

- references resolve,
- titles are consistent,
- no file contains more than one RFC,
- invariants are extractable into a draft catalog (input to ADQ-005),
- the RFC dependency graph is emittable as simple JSON.

## Safe vs Unsafe in Phase 0

**Safe (in scope):**
- read-only RFC checkers (references, titles, duplicate bodies),
- draft invariant extraction (non-authoritative),
- simple-JSON RFC dependency graph,
- documentation of open architecture decisions.

**Unsafe (NOT in scope — blocked until ADQ decisions):**
- Canonical Object store shape, Artifact store, review-target typing (ADQ-001),
- Memory / Knowledge service implementation (ADQ-002),
- Workflow engine (ADQ-003),
- Action persisted state machine (ADQ-004),
- the canonical invariant ID registry itself (ADQ-005 — only a *draft* is produced),
- editing RFC files to fix the findings below.

## Structural checks (cause non-zero exit)

Run via `python verify/rfc/run.py`. The following are **errors**:

1. **Duplicate RFC bodies** — a file contains more than one top-level `# RFC-NNN` heading.
2. **Missing RFC files** — a number listed in the README numbering table has no `rfcs/NNN-*.md`.
3. **Broken references** — an `RFC-NNN` reference points to a number with no file.
4. **Reference title mismatch** — a reference of the form `RFC-NNN (Some Title)` names a title
   that does not match the referenced file's heading title.

## Warnings (non-blocking)

- **Title-vs-index drift** — a file's `# RFC-NNN <Title>` differs from the README numbering
  table title.
- **Architecture Decision advisories** — the ambiguities tracked as ADQ-001…ADQ-006 are emitted
  as advisories, never as failures (the tooling does not resolve them).

## Findings this plan covers (from the audit)

These are documented here as cleanup tasks. **They require editing RFC files and are therefore
out of Phase 0 tooling scope** — they are listed so the structural checkers can confirm when each
is fixed. Editing the RFCs is a separate, explicitly-approved step.

| Finding | Type | Detected by | Cleanup task (separate step) |
|---------|------|-------------|------------------------------|
| F7 RFC-052 contains two RFC bodies | FACT | `check_duplicate_bodies.py` | Merge unique content, remove the first body. |
| F8 RFC-053 forward-refs wrong titles | FACT | `check_rfc_references.py` | Correct titles for RFC-054/060/061/062; drop false prerequisite. |
| F9 RFC-054 names RFC-050 "Business Space" | FACT | `check_rfc_references.py` | Rename references to "RFC-050 Space Model". |
| F10 RFC-003 stale pipeline/laws | CONTRADICTION | manual review | Update after ADQ-001 (it encodes the contested Artifact model). |
| F11 Author/date metadata drift | FACT | manual review | Normalize headers (RFC-023, 051, 053, 054). |
| F12 "Workspace" reused as memory scope | AMBIGUITY | manual review | Rename RFC-043 memory scope to avoid the RFC-002 term clash. |
| F13 Pipeline objects listed as Canonical Objects | FACT | manual review | Remove Event/Decision/Review from RFC-051/053 object lists (after ADQ-001). |
| Title-vs-index drift (RFC-043) | LOW | `check_rfc_titles.py` (warning) | Align README numbering title with the file H1. |

## Sequencing

1. Run the checkers; triage the structural errors.
2. Fix the **pure-documentation** facts that carry no architectural choice: F7, F8, F9, F11
   (and the RFC-043 title drift). These do not depend on any ADQ.
3. **Hold** F10 and F13 until **ADQ-001** is decided (they encode the contested Artifact model).
4. Decide ADQ-002, ADQ-004, ADQ-005, ADQ-006 before the corresponding subsystems are built.
5. Re-run the checkers until structural errors are zero.

## What stays open after Phase 0

Phase 0 produces a clean, machine-checkable corpus and a documented decision queue. It does
**not** unblock the domain model: that remains gated on **ADQ-001**. See the audit's readiness
matrix and the ADQ file for the full list.
