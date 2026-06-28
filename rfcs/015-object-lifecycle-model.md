# RFC-015 Object Lifecycle Model

**Status:** Draft  
**Authors:** Tiroq + ChatGPT  
**Last Updated:** 2026-06-28

---

## 1. Summary

This RFC defines the lifecycle of Canonical Objects from creation to archival and connects Events, Artifacts, Revisions, Projections, and Views. It establishes a comprehensive model describing how objects evolve through various stages, how their revisions are managed immutably, and how projections and views are synchronized and rendered. This lifecycle model forms a foundational component of the overall system architecture, ensuring consistency and clarity in object state transitions and representations.

## 2. Relationship to Previous RFCs

This RFC builds upon the foundations laid out in RFC-000 through RFC-014. While those RFCs established core concepts such as Events, Artifacts, and basic object models, RFC-015 completes the Foundation architecture by formalizing the lifecycle of Canonical Objects and their associated representations. It integrates and extends prior work to provide a unified lifecycle framework essential for subsequent RFCs.

## 3. Goals

- Define a clear and comprehensive lifecycle for Canonical Objects.
- Separate the processing lifecycle (how data is handled) from the business lifecycle (how objects evolve in business terms).
- Define an immutable revision model for object versions.
- Define the projection lifecycle for derived data representations.
- Define archival procedures and states.

## 4. Non-Goals

- This RFC does not specify persistence implementations.
- It does not define workflow engines or storage details.
- It does not prescribe UI or UX specifics.

## 5. Lifecycle Philosophy

The lifecycle model distinguishes three interrelated but distinct lifecycles:

- **Events** have *processing states* that track how events are handled within the system.
- **Canonical Objects** have *business lifecycles* representing their real-world or domain-related state transitions.
- **Projections** have *synchronization lifecycles* describing how derived representations stay current or become stale.

This separation ensures clarity in responsibilities and enables flexible system design.

## 6. Canonical Lifecycle

```mermaid
flowchart LR
    UND[Understanding]
    --> CO[Canonical Object]
    --> REV[Revision]
    --> PROJ[Projection]
    --> RM[Read Model]
    --> VIEW[View]

    CO --> ARCH[Archived]
```

This RFC intentionally begins after Understanding has produced or updated a Canonical Object. Events and Event processing are defined by RFC-013. This RFC defines how Canonical Objects evolve throughout their business lifetime and how their derived representations are managed.

## 7. Lifecycle Stages

| Lifecycle Stage | Description |
|-----------------|-------------|
| Created | Object has been materialized. |
| Active | Object participates in business processes. |
| Suspended | Object is temporarily inactive without losing identity. |
| Archived | Object is no longer active but remains historically accessible. |

Individual Artifact types may define additional business-specific states such as Approved, Scheduled, Executing or Completed without redefining the canonical lifecycle.

## 8. Processing vs Business Lifecycle

| Event Processing | Object Lifecycle | Business State |
|------------------|------------------|----------------|
| Captured | Created | Draft |
| Validated | Active | Planned |
| Accepted | Active | Approved |
| Understood | Active | Executing |
| Archived | Archived | Completed |

This table distinguishes the transient event processing states, the persistent object lifecycle states, and the business-specific states which describe domain semantics. Processing states reflect how events move through the system, lifecycle states define object existence and activity, and business states capture domain-specific meaning.

## 9. Revision Model

Revisions are immutable snapshots of a Canonical Object at a point in time. Each revision records:

- A unique identifier.
- The lineage (parent revision(s)).
- The author of the revision.
- Timestamp of creation.
- Reason for the revision.
- Source decision or event that triggered the revision.

Revisions form a chain representing the full history of an object.

```mermaid
graph TD
    Rev1["Revision 1"]
    Rev2["Revision 2"]
    Rev3["Revision 3"]
    Rev1 --> Rev2 --> Rev3
```

## 10. Projection Lifecycle

Projections are derived representations of revisions optimized for specific queries or use cases. Their lifecycle stages include:

- **Created:** Projection instance is created.
- **Synchronized:** Projection is up-to-date with the latest revisions.
- **Stale:** Projection is out-of-date due to new revisions.
- **Rebuilt:** Projection is refreshed or rebuilt from source revisions.
- **Removed:** Projection is deleted or discarded.

Projections are considered disposable and replaceable.

## 11. Read Model Lifecycle

Read Models are generated from projections to support querying and UI rendering. Their lifecycle stages are:

- **Generated:** Initial creation.
- **Refreshed:** Updated to reflect changes.
- **Invalidated:** Marked outdated due to data changes.
- **Rebuilt:** Fully reconstructed.

## 12. View Lifecycle

Views represent rendered presentations of Read Models, typically consumed by clients or UIs. Their lifecycle stages are:

- **Requested:** A view is requested.
- **Rendered:** The view is generated.
- **Expired:** The view is no longer valid or cached.

Views do not own state and are ephemeral by nature.

## 13. Lifecycle Invariants

- Canonical Objects own their lifecycle and transitions.
- Revisions never overwrite history; they append immutable state.
- Projections are replaceable and disposable.
- Read Models are rebuildable to ensure consistency.
- Views are ephemeral and stateless.
- Object identity persists across lifecycle transitions.

## 13.1 Lifecycle Ownership

- Only Canonical Objects own business lifecycle.
- Revisions inherit lifecycle context but never define lifecycle.
- Projections reflect lifecycle but never control it.
- Read Models expose lifecycle information but may be rebuilt.
- Views render lifecycle state but never persist it.

## 14. Lifecycle Relationships

```mermaid
graph LR
    Canonical_Object --> Revision
    Revision --> Projection
    Projection --> Read_Model
    Read_Model --> View
```

This diagram shows one Canonical Object producing multiple revisions, which feed projections and read models.

## 15. Replay & Recovery

Replay reconstructs derived representations from immutable Event Records and Canonical Objects without modifying identity.

Rebuild targets include:

- Projections
- Read Models
- Search Indexes
- Knowledge Graph relationships
- Analytics Views
- Cached representations

Replay must never mutate:

- Canonical Identity
- Revision History
- Historical Decisions

## 16. Dependencies

Depends on RFC-000 through RFC-014.  
Required before:  
- RFC-020 Review System  
- RFC-021 Decision Model  
- RFC-030 System Architecture  
- RFC-032 Data Flow  
- RFC-043 Memory & Knowledge Graph

## Architectural Summary

- Event processing and object lifecycle are separate concerns.  
- Canonical Objects own lifecycle.  
- Identity survives every lifecycle transition.  
- Revisions preserve history.  
- Projections are replaceable.  
- Read Models are rebuildable.  
- Views are ephemeral.

## 17. Acceptance Criteria

- Lifecycle stages and transitions are clearly defined and implemented.
- Revision immutability is enforced.
- Projection and read model synchronization mechanisms are reliable.
- Archival processes are defined and tested.
- Documentation and diagrams are complete and accurate.

## 18. Decision Log

| Date       | Decision                          | Notes                       |
|------------|---------------------------------|-----------------------------|
| 2026-06-28 | Initial draft completed          | RFC-015 submitted for review|

---

"Objects evolve. Identity endures. History is never lost."
