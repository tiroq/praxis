# Architecture Decision Queue (ADQ)

> Source of truth: the **Corrected Praxis RFC Audit**.
> The ADQ is the **authoritative list of unresolved architecture questions**. As long as an
> item lives here in any non-`CLOSED` state, the question it names is considered **open** — no
> ADR, preference, or implementation detail overrides that until the item is closed by the
> lifecycle below.
> This file **documents** open architecture decisions. It does **not** resolve them.
> No item here is decided. "Recommended option" is a *preference with reasoning*, not a ruling.
> An item is only "decided" when an RFC is updated to state the decision explicitly.

## Lifecycle

Every ADQ moves through the following states:

```
OPEN
  → PROPOSED ADR
    → ACCEPTED ADR
      → RFC UPDATE REQUIRED
        → CLOSED
```

- **OPEN** — the question is recorded; no ADR proposes an answer yet.
- **PROPOSED ADR** — one or more ADRs propose an answer but are not accepted.
- **ACCEPTED ADR** — an ADR resolving this ADQ has been accepted.
- **RFC UPDATE REQUIRED** — the ADR is accepted; the listed RFC updates are pending.
- **CLOSED** — the required RFC updates have landed and the question is settled.

### Rules

- A **PROPOSED ADR does not close an ADQ.** The ADQ stays open (state `PROPOSED ADR`).
- An **ACCEPTED ADR may close an ADQ only after the required RFC updates are listed** on this
  ADQ. Acceptance advances the item to `ACCEPTED ADR` → `RFC UPDATE REQUIRED`; `CLOSED` is
  reached only once those RFC updates land.
- **Every ADQ must link to its related ADR(s).** An ADQ with no ADR yet stays `OPEN`.
- **Every ADR must link back to the ADQ(s) it addresses** (via its `Related:` line).
- ADQ items are **not deleted**. Resolved items remain in the table marked `CLOSED`.

Status / lifecycle values: `OPEN` · `PROPOSED ADR` · `ACCEPTED ADR` · `RFC UPDATE REQUIRED` ·
`CLOSED`. An item is only truly "decided" once it is `CLOSED` and an RFC explicitly records it.

| ID | Decision | Recommended (preference, not decided) | Related ADRs | Lifecycle state | Blocks | Owner |
|----|----------|----------------------------------------|--------------|-----------------|--------|-------|
| ADQ-001 | Canonical Object vs Artifact model | B (representation w/ own id) + D (interim ReviewTarget); A is the simplest alt | ADR-001, ADR-007 | PROPOSED ADR | Domain model, storage ownership, review/decision typing, space persistence | Architecture |
| ADQ-002 | Memory / Knowledge ownership | A or D (standalone / platform services via RFC-031) | ADR-004, ADR-005, ADR-007 | PROPOSED ADR | Memory subsystem, service roster | Architecture |
| ADQ-003 | Workflow model | D for Phases 1–2, then B (orchestration pattern) via a new RFC | ADR-002, ADR-003 | PROPOSED ADR | Orchestration layer (Phase 3+) | Architecture |
| ADQ-004 | Action state machine reconciliation | A (RFC-022 authoritative) | — (no ADR yet) | OPEN | Action Service implementation | Architecture |
| ADQ-005 | Invariant ID registry | A or B (registry in RFC-060, or per-RFC IDs aggregated) | — (no ADR yet) | OPEN | Verification script manifest mapping (RFC-061) | Documentation / Verification |
| ADQ-006 | Space storage mapping consistency | A (rewrite 051/052/053/054 to RFC-033 categories) | ADR-007 | PROPOSED ADR | Space-specific persistence design | Documentation |
| ADQ-007 | External Workflow Orchestrator (operational pipelines) | Integrate Kestra as infrastructure (not a Praxis component) | ADR-002 | PROPOSED ADR | Operational scheduling/retry/observability for maintenance pipelines | Architecture |

---

## ADQ-001 — Canonical Object vs Artifact model

**Problem**
The RFCs do not unambiguously state the relationship between a Canonical Object and an
Artifact. They are sometimes treated as the same identity-bearing, reviewable unit and
sometimes as distinct layers.

**Evidence**
- RFC-002 (Provisional Concepts): "**Artifact** — A durable representation or document that
  may represent, derive from, or reference a Canonical Object. **Artifacts do not own
  canonical identity.**"
- RFC-012 §8/§18: "Each Artifact must have a stable identity." Artifact is also the unit that
  is reviewed, decided upon, revised, and lifecycle-managed.
- RFC-011 §10: "Reviews and Decisions operate on **Canonical Objects**."
- RFC-014 §9/§15: the Canonical Object is the sole owner of business identity; representations
  (Projection / Read Model / View) never own identity. Artifact is not placed in this layer list.
- RFC-020 §15 / RFC-021 §15: each Review/Decision targets "exactly one **Canonical Object or
  Artifact**."
- RFC-033 has a Canonical Store and a Revision Store but **no Artifact Store**.

This is an AMBIGUITY with an embedded CONTRADICTION (RFC-002 "no canonical identity" vs
RFC-012 "stable identity"), depending on whether "stable identity" means "canonical identity".

**Options**
- **A. Artifact is a typed Canonical Object** (a CO subtype).
- **B. Artifact is a durable, separately-identified representation** that references/derives
  from a Canonical Object but does not own the *canonical* business identity.
- **C. Artifact is an input/output document** that may be *promoted* into a Canonical Object.
- **D. Artifact and Canonical Object are distinct**, both implementing a common
  `ReviewTarget` / `DecisionTarget` interface.

**Evaluation (summary)**
- A: simplest to build; best fit with RFC-012; weak fit with RFC-002 (which says Artifact owns
  no canonical identity); Artifacts live cleanly in the Canonical Store.
- B: best fit with RFC-002; reconciles both identity statements (own document-identity ≠
  canonical-identity); needs an Artifact store or careful non-derived modeling.
- C: fits RFC-002 wording; adds promotion machinery; storage similar to B.
- D: directly explains the "CO or Artifact" review/decision phrasing; most extensible; orthogonal
  to the conceptual A/B/C choice (it is an interface mechanism).

**Recommended option (preference, not decided)**
**B** for the conceptual model (strongest foundation/RFC-002 consistency and reconciles both
literal statements), with **D** as the interim review-typing mechanism. **A** is the simplest
alternative if the team chooses instead to revise RFC-002's wording. This is a genuine fork.

**Decision status:** OPEN
**Blocks:** domain model, storage ownership, review/decision target typing, space-scoped
persistence, plus the doc rewrites in RFC-003 (stale pipeline) and Space object lists.
**Owner:** Architecture
**Required RFC updates once decided:** RFC-002, RFC-012, RFC-014 (and downstream RFC-003,
RFC-011, RFC-020, RFC-021, RFC-033, RFC-050) to state the chosen model and align identity,
storage, and review-target wording.

---

## ADQ-002 — Memory / Knowledge ownership

**Problem**
RFC-033/043/050 name a "Memory Service" and a "Knowledge Service" as owners, but RFC-030/031
(the authoritative service roster and contract list) do not define them.

**Evidence**
- RFC-030 §7 and RFC-031 §11 list the runtime services; neither includes Memory or Knowledge.
- RFC-033 §6/§30 assign stores to "Knowledge Service" and "Memory/Search services".
- RFC-050 §12 (Space Runtime Mapping) lists "Memory Service" and "Knowledge Service".
- RFC-040 §23: agents "do not own memory"; "Knowledge services own persistence and validation".
- RFC-043 §28: "Memory never owns canonical truth"; "Knowledge Graph references Canonical Objects".

**Options**
- A. Standalone runtime services (add to RFC-030/031).
- B. Capabilities inside the Agent Runtime.
- C. Storage-backed subsystems owned by the Canonical Domain Service.
- D. Platform services exposed through RFC-031 contracts.

**Evaluation (summary)**
- B contradicts RFC-040 §23 / RFC-043 §28 (agents do not own memory).
- C contradicts RFC-043 §28 (memory/knowledge are separate from canonical truth).
- A and D match RFC-033 §30 and RFC-050 §12, which already treat them as named owners.

**Recommended option (preference, not decided)**
**A or D** — add both as first-class (or platform-tier) services with RFC-031 contracts. The
A-vs-D distinction (tier/label) is minor.

**Decision status:** OPEN
**Blocks:** memory subsystem implementation, service-roster completeness.
**Owner:** Architecture (with Documentation)
**Required RFC updates once decided:** RFC-030 §7 and RFC-031 §11 to add the services; confirm
ownership wording in RFC-033 §30 and RFC-050 §12.

---

## ADQ-003 — Workflow model

**Problem**
"Workflow" is a canonical term and is referenced across RFCs, but no RFC defines its identity,
lifecycle, state machine, or contract.

**Evidence**
- RFC-002 §4 and RFC-003 taxonomy define "Workflow" (Process category).
- RFC-010, RFC-030 (Application layer), RFC-040, RFC-050 reference workflows.
- RFC-022 enumerates ten runtime state machines; Workflow is not among them.
- RFC-013 §21: "Workflow execution begins only after Understanding has produced one or more
  Artifacts."

**Options**
- A. Workflow as a first-class Canonical Object.
- B. Workflow as an orchestration pattern over Commands/Events/Decisions.
- C. Workflow as a projection only.
- D. Workflow as a deferred implementation concern (no RFC yet).

**Evaluation (summary)**
- Phases 1–2 (kernel, object lifecycle, review/decision/action) do not require a Workflow
  object; the pipeline is expressed as discrete runtime objects with their own state machines.
- B aligns with RFC-030 ("Application coordinates but does not own domain invariants") and
  RFC-031's command/event model.

**Recommended option (preference, not decided)**
**D for Phases 1–2**, then author a Workflow RFC (leaning **B**) before orchestration-heavy work.
A new RFC is **not** a precondition for early implementation.

**Decision status:** OPEN (deferrable)
**Blocks:** orchestration layer (Phase 3+).
**Owner:** Architecture
**Required RFC updates once decided:** new Workflow RFC (e.g., RFC-024) or an orchestration
section in RFC-030.

---

## ADQ-004 — Action state machine reconciliation

**Problem**
RFC-022 and RFC-023 define different state machines for the same Action and Action Request.

**Evidence**
- RFC-022 §9.6 Action Request: `Issued → Dispatched → Completed/Cancelled → Archived`.
  §9.7 Action: `Pending → Running → Succeeded/Failed → Archived`.
- RFC-023 Action Request: `Created → Planned → Ready → Queued → Executed → Completed`.
  Action: `Planned → Assigned → In Progress → Succeeded/Failed → Compensated`.
- RFC-022 §9: specialized RFCs "may extend ... but must not contradict" the baseline.

**Options**
- A. RFC-022 authoritative; express RFC-023 states as documented sub-states/extensions.
- B. RFC-023 authoritative; update RFC-022 baseline to match.
- C. Define an explicit mapping table, keep both.

**Recommended option (preference, not decided)**
**A**, because RFC-022 is the dedicated state-machine RFC and explicitly claims precedence.
B is acceptable if the team prefers RFC-023's granularity as canonical.

**Decision status:** OPEN
**Blocks:** Action Service implementation (HIGH rework risk if built before reconciliation).
**Owner:** Architecture
**Required RFC updates once decided:** RFC-022 and/or RFC-023 to align the Action and Action
Request transitions.

---

## ADQ-005 — Invariant ID registry

**Problem**
RFC-061 maps verification scripts to invariant IDs (`RFC-060-INV-001`, `-010`) that RFC-060
never defines. There is no stable invariant ID scheme.

**Evidence**
- RFC-061 §10 manifest example references `RFC-060-INV-001` / `RFC-060-INV-010`.
- RFC-061 §11 requires every RFC-060 invariant to be covered by a script.
- RFC-060 lists invariants as prose bullets with no IDs.

**Options**
- A. Add an invariant ID registry to RFC-060 (`RFC-060-INV-###`).
- B. Assign IDs in each source RFC; RFC-060 aggregates references.

**Recommended option (preference, not decided)**
**A or B** — both yield a stable, testable catalog. The Phase 0 tool `extract_invariants.py`
produces a **draft, non-authoritative** catalog to inform this decision; it does not define the
canonical registry.

**Decision status:** OPEN (mechanical, low risk)
**Blocks:** RFC-061 manifest-to-invariant mapping.
**Owner:** Documentation / Verification
**Required RFC updates once decided:** RFC-060 (and possibly each invariant-bearing RFC) to
assign and freeze invariant IDs.

---

## ADQ-006 — Space storage mapping consistency

**Problem**
Some Space RFCs specify concrete database engines and per-object retention, contradicting the
engine-agnostic, service-owned storage model of RFC-033.

**Evidence**
- RFC-033 §4 Non-Goals: "does not Select a concrete database engine."
- RFC-053 storage mapping: PostgreSQL, Neo4j/JanusGraph, S3/Azure Blob, Snowflake/BigQuery.
- RFC-054 storage mapping: "CRM database", "secure contract vault", "Finance system" + retention years.
- RFC-052 first body storage mapping: "Relational DB", "Document Store" + retention years.
- RFC-055 §31 and RFC-056 §35 already use RFC-033 store categories — the intended pattern.

**Options**
- A. Rewrite 051/052(first)/053/054 storage mappings to RFC-033 categories (match 055/056).
- B. Amend RFC-033 to permit engine hints (not recommended; weakens the boundary).

**Recommended option (preference, not decided)**
**A** — the correct pattern already exists in RFC-055/056; this is cleanup, not a fork.

**Decision status:** OPEN (cleanup; tracked to ensure it happens)
**Blocks:** space-specific persistence design; documentation trust.
**Owner:** Documentation
**Required RFC updates once decided:** RFC-051, RFC-052, RFC-053, RFC-054 storage-mapping
sections.
