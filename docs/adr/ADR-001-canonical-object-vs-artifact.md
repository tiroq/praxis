# ADR-001 — Canonical Object vs Artifact

**Status:** PROPOSED (not accepted — see ADQ-001)
**Date:** 2026-06-29
**Authors:** Architecture
**Supersedes:** —
**Superseded by:** —
**Related:** ADQ-001, RFC-002, RFC-012, RFC-014, RFC-020, RFC-021, RFC-033

---

## Context

The Praxis RFCs do not state a single, consistent relationship between a **Canonical Object**
and an **Artifact**. Four different readings are simultaneously supported by the current text,
and two of them contradict each other directly. This creates an unresolvable ambiguity in every
downstream RFC that touches identity, storage, review targets, or space persistence.

### Evidence of ambiguity and contradiction

| RFC | Quote / Claim | Implied model |
|-----|---------------|---------------|
| RFC-002 §4 Provisional | "Artifact — A durable representation … **Artifacts do not own canonical identity.**" | B or C: Artifact ≠ Canonical Object |
| RFC-012 §8 | "Each Artifact must have a stable identity." | A or B: Artifact has its own persistent identity |
| RFC-012 §7 | "Artifact **may become or reference** a Canonical Object." | C (promotable) or B (reference) |
| RFC-014 §5–§9 | Identity layer model: Concept → Identity → **Canonical Object** → Projection → Read Model → View. Artifact is **not placed** in this hierarchy. | Artifact is outside the layer stack entirely |
| RFC-011 §10 *(cited in ADQ)* | "Reviews and Decisions operate on **Canonical Objects**." | A: Artifact is a kind of CO |
| RFC-020 §1, §8 | Review evaluates "Canonical Objects and Artifacts"; Object ID field references "the reviewed Canonical Object or Artifact." | D: both are review-targetable peers |
| RFC-021 §15 | "Every Decision targets exactly one **Canonical Object or Artifact**." | D: both are decision-targetable peers |
| RFC-033 §6, §9–§10 | Defines Canonical Store and Revision Store; Revision Store covers "Canonical Objects and Artifacts." **No Artifact Store is defined.** | Ambiguous: Artifact may live in Canonical Store, or has no assigned home |

The core contradiction is between RFC-002's "Artifacts do not own canonical identity" and
RFC-012's "Each Artifact must have a stable identity." Whether these statements conflict
depends entirely on what "canonical identity" means — a question no RFC currently answers.

---

## Decision

> **This ADR proposes Option B, augmented by Option D as the review/decision-typing
> mechanism. It is PROPOSED, not accepted. No RFC shall be changed until this ADR is
> promoted to ACCEPTED.**

---

## Options Considered

### Option A — Artifact is a typed Canonical Object

Artifact is a subtype (or role) of Canonical Object. RFC-002's "no canonical identity"
statement is rewritten to remove the exclusion; RFC-012's "stable identity" is canonical
identity. Artifacts live in the Canonical Store alongside all other Canonical Objects.

**Fits:** RFC-012 (full lifecycle, identity, revision history, review/decision participation).
RFC-011 ("reviews operate on Canonical Objects" — because Artifact is one).
RFC-033 (Canonical Store already defined; no new store needed).

**Conflicts:** RFC-002 provisional definition, which explicitly excludes canonical identity
from Artifact. Accepting A requires rewriting the RFC-002 definition — the clearest statement
of intent the system has.

**Implementation cost:** Lowest. No new storage category. Artifact becomes a tagged subtype;
all existing Canonical Object machinery applies directly.

**Verdict:** Simplest. Best fit with implementation. Weakest conceptual foundation because it
discards the most explicit definitional statement in the terminology RFC.

---

### Option B — Artifact is a separately-identified durable representation *(PROPOSED)*

Artifact has its own **document-level identity** (stable `artifact_id`, immutable after
creation), but does not own the **canonical business identity** of the domain entity it
represents. It references or derives from a Canonical Object — or, in its initial unreviewed
state, exists prior to one being created. The RFC-002 statement ("Artifacts do not own
canonical identity") and the RFC-012 statement ("stable identity") are both true and
non-contradictory under this reading: "stable identity" refers to the artifact's own
persistent reference, not to the canonical business identity defined in RFC-014.

**Fits:** RFC-002 (Artifact does not own *canonical* identity, but has its own durable id).
RFC-012 (Artifact identity, lifecycle, revisions are fully defined). RFC-014 (Canonical Object
retains sole ownership of business identity; Artifact sits beside the layer hierarchy as a
durable document-class object rather than inside it). RFC-020 and RFC-021 (Artifact is a
first-class review/decision target independently of a Canonical Object).

**Conflicts:** RFC-033 has no Artifact Store. Under Option B, Artifacts require explicit
storage ownership — either their own store category or a documented sub-category of the
Canonical Store with distinct identity rules.

**Requires:** RFC-014 to add Artifact to the representation model with a note distinguishing
document-level identity from canonical business identity. RFC-033 to add an Artifact Store
or formally place Artifacts as a sub-category of the Canonical Store with their own identity
and ownership rules.

**Implementation cost:** Medium. Requires a new or extended storage category. Artifact Service
(or Canonical Domain Service) must distinguish artifact identity from canonical object
identity in persistence, APIs, and schemas. Review and Decision services remain unchanged
because both already reference "Canonical Object or Artifact."

**Verdict:** Strongest conceptual foundation. Reconciles all RFC statements without discarding
any. Adds one storage concern that must be resolved.

---

### Option C — Artifact is a promotable input/output document

Artifact starts as an unstructured input or output document (a "candidate"). It can be
*promoted* into a Canonical Object when the system commits to owning it as a first-class
business entity. Pre-promotion, the Artifact has only document identity; post-promotion, the
resulting Canonical Object owns canonical business identity. Pre-promotion artifacts may be
reviewed (hence RFC-020/021 "or Artifact" phrasing).

**Fits:** RFC-002 literally (pre-promotion Artifacts do not own canonical identity). The
RFC-012 §7 phrase "Artifact may become … a Canonical Object" reads most naturally under this
option.

**Conflicts:** RFC-012 §8–§13 describe Artifacts as having a full lifecycle (Draft →
Archived), revision history, review and decision participation, and knowledge graph
membership — behaviors that make sense only for a first-class, post-promotion object. The
promotion boundary is not defined by any RFC. RFC-012's "may become" language is ambiguous
between promotion (C) and reference (B).

**Implementation cost:** High. Requires a promotion state machine and a mechanism that either
transforms or replaces an Artifact with a Canonical Object on promotion, while preserving
all pre-promotion review and decision history.

**Verdict:** Adds machinery that no RFC requires. The RFC-012 lifecycle implies that the
promoted state is the normal operating state, making Option C an awkward fit.

---

### Option D — Artifact and Canonical Object are distinct; both implement ReviewTarget / DecisionTarget

This is an orthogonal design mechanism rather than a mutually exclusive conceptual option.
RFC-020 and RFC-021 already describe both as valid targets. Option D formalizes that phrasing
into a `ReviewTarget` and `DecisionTarget` interface (or type union) that both Canonical
Objects and Artifacts satisfy.

**Fits:** RFC-020 §1, §8 and RFC-021 §15 directly. Keeps the review and decision systems
generic without requiring resolution of the A/B/C conceptual question.

**Conflicts:** D alone does not answer what an Artifact *is* — only that it can be reviewed
and decided upon. The underlying conceptual question (A, B, or C) must still be answered.

**Implementation cost:** Low incremental cost if used alongside B. Slightly higher if used as
the sole resolution (because A/B/C ambiguity remains in storage and identity layers).

**Verdict:** D is not an alternative to A/B/C; it is a complementary interface mechanism. The
proposed recommendation adopts D alongside B.

---

## Proposed Decision: B + D

**Conceptual model:** Option B.
An Artifact is a durable document with its own stable `artifact_id` that references or
derives from a Canonical Object (or exists independently of one in its early lifecycle), but
does not own canonical business identity as defined by RFC-014.

**Review/Decision typing:** Option D.
Both Canonical Objects and Artifacts satisfy the `ReviewTarget` and `DecisionTarget`
contracts. RFC-020 and RFC-021 language ("Canonical Object or Artifact") is formalized as a
named type union or interface rather than an ad-hoc disjunction.

**Fallback:** If the architecture team chooses to minimise RFC changes, Option A is the
recommended alternative. It requires only a one-sentence change to RFC-002's provisional
Artifact definition and avoids introducing a new storage category.

---

## Consequences

### If B + D is accepted

**Positive:**
- RFC-002 and RFC-012 identity statements are reconciled without discarding either.
- RFC-014's representation model gains a clear place for Artifact as a parallel durable class
  that does not substitute into the Concept → CO → Projection hierarchy.
- RFC-020 and RFC-021 review/decision phrasing is formalized by the D interface mechanism.
- The model is extensible: future document classes can satisfy `ReviewTarget` without
  becoming Canonical Objects.

**Negative / costs:**
- RFC-033 must define an Artifact Store or formally document Artifact as a sub-category of
  the Canonical Store with distinct identity rules.
- RFC-014 must add Artifact to the representation model with a note on document-level vs
  canonical identity.
- Implementations must maintain two identity namespaces: `canonical_id` (RFC-014) and
  `artifact_id` (RFC-012). APIs and schemas must not conflate them.
- The `ReviewTarget` / `DecisionTarget` interface must be added to RFC-020/021 and
  implementation contracts.

### If A is accepted (fallback)

**Positive:**
- Lowest implementation cost. Single storage category. No new interface abstraction needed.

**Negative / costs:**
- RFC-002's provisional definition must be changed to remove "Artifacts do not own canonical
  identity." This is the clearest explicit statement in the terminology RFC and discarding it
  weakens the document's authority.
- RFC-014's layer model must be extended to position Artifact as a subtype of Canonical
  Object, including clarification of which CO invariants Artifacts inherit.
- The "or Artifact" language in RFC-020/021 becomes redundant (Artifact is a CO) and must be
  either kept as a domain term or removed.

---

## Required RFC Edits (if B + D is accepted)

| RFC | Required change | Scope |
|-----|-----------------|-------|
| RFC-002 | Revise Provisional Concepts definition of Artifact to distinguish document-level identity from canonical business identity | Terminology |
| RFC-012 | Add explicit statement: "Artifact identity (`artifact_id`) is document-level identity, distinct from canonical business identity as defined by RFC-014." Clarify §7 "may become or reference" to mean reference (not promotion). | Artifact model |
| RFC-014 | Add Artifact to the conceptual model as a durable document class that owns document identity but not canonical business identity. Distinguish from Projection (which is derived and disposable). | Identity model |
| RFC-020 | Add formal definition of `ReviewTarget` as the type satisfied by both Canonical Objects and Artifacts. Update §8 Object ID field description accordingly. | Review system |
| RFC-021 | Add formal definition of `DecisionTarget`. Update §15 invariant to reference `DecisionTarget` by name. | Decision model |
| RFC-033 | Define an Artifact Store (or formally specify Artifact as a Canonical Store sub-category with its own identity and ownership rules). Update storage category table. | Storage model |
| RFC-003 | Minor: update pipeline text to use the reconciled Artifact definition | Concept model |
| RFC-011 | Minor: align "Reviews and Decisions operate on Canonical Objects" to "… on ReviewTargets / DecisionTargets" | Domain model |
| RFC-050 | Minor: verify Space object lists use the RFC-012/RFC-014-aligned Artifact definition | Space model |

---

## Implementation Impact

### Not blocked (safe to proceed regardless of outcome)

- Event ingestion, event store, and event schemas (RFC-013 is unaffected)
- Review and Decision service internal logic (RFC-020/021 behavior is unchanged; only type
  contract formalization is needed)
- Action model (RFC-023)
- State machine (RFC-022)
- LLM routing (RFC-041)

### Blocked until accepted

- Canonical Domain Service schema: cannot define the table/document structure for Artifacts
  without knowing whether they share the Canonical Object table or have a separate store
- Storage model implementation (RFC-033): Artifact Store or sub-category must be specified
- Space persistence (RFC-051–056): Space object lists reference Artifacts; storage mapping
  depends on where Artifacts live
- API contracts that expose Artifact identity: `artifact_id` vs `canonical_id` field
  naming requires B vs A to be settled

---

## Verification Impact

### Existing verification scripts (RFC-061) affected

- Any invariant that asserts "every Canonical Object has exactly one entry in the Canonical
  Store" must be updated to handle the A vs B split (under B, Artifacts may not be in the
  Canonical Store)
- Review and Decision invariant checks that assert "target is a Canonical Object" must be
  updated to accept "target is a ReviewTarget / DecisionTarget" (Option D)
- Identity uniqueness checks must be extended to cover both `canonical_id` and `artifact_id`
  namespaces under Option B

### New verification required

Under B + D:
- Assert that no Artifact carries a `canonical_id` (document identity must not be confused
  with canonical business identity)
- Assert that every ReviewTarget and DecisionTarget is either a Canonical Object or an
  Artifact (not any other object type)
- Assert that Artifact storage owner is consistent with the store category defined in RFC-033

---

## Rejected Alternatives

### Why not C alone (promotable document)

The RFC-012 lifecycle (Draft → Archived with full revision history, reviews, decisions) is
defined for Artifacts as they exist *now*, not for a post-promotion entity. Requiring
promotion machinery before an Artifact can participate in the lifecycle is inconsistent with
the current RFC-012 design and adds complexity that no RFC currently requires.

### Why not D alone (interface-only resolution)

D does not answer the identity or storage questions. It cannot resolve the RFC-002 vs RFC-012
contradiction because neither document is about review/decision typing. D is necessary but
not sufficient.

### Why not A + D

A discards the most explicit definitional statement in RFC-002 without a strong reason. The
B reconciliation is available at comparable implementation cost (one additional storage
category) and a stronger long-term foundation. A is listed as a fallback, not the primary
recommendation.

---

## Open Questions

1. **Artifact Store vs Canonical Store sub-category:** Under Option B, should Artifacts live
   in a dedicated Artifact Store (cleaner separation) or as a documented sub-type in the
   Canonical Store (lower operational cost)? RFC-033 must decide this once this ADR is
   accepted.

2. **RFC-003 Concept Model pipeline:** RFC-003 describes a stale pipeline that needs
   updating independent of this decision. Should RFC-003 be updated concurrently with this
   ADR, or in a separate follow-up?

3. **Promotion path (Option C residual):** RFC-012 §7 says "Artifact may become … a
   Canonical Object." Under Option B (reference model), does this path ever exist, or should
   that clause be removed? If it remains, a promotion invariant must be defined.

---

## References

- [docs/ARCHITECTURE_DECISION_QUEUE.md](../ARCHITECTURE_DECISION_QUEUE.md) — ADQ-001 entry
- [rfcs/002-terminology.md](../../rfcs/002-terminology.md) — Artifact provisional definition
- [rfcs/012-artifact-model.md](../../rfcs/012-artifact-model.md) — Artifact identity, lifecycle, review participation
- [rfcs/014-identity-representation-model.md](../../rfcs/014-identity-representation-model.md) — Canonical Object and layer model
- [rfcs/020-review-system.md](../../rfcs/020-review-system.md) — ReviewTarget phrasing
- [rfcs/021-decision-model.md](../../rfcs/021-decision-model.md) — DecisionTarget invariant
- [rfcs/033-storage-model.md](../../rfcs/033-storage-model.md) — Storage categories (no Artifact Store defined)
