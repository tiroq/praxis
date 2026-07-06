---
name: "Architecture Review Discipline — Reference"
description: "Reference documentation. For active Copilot instructions, see .github/copilot-instructions.md"
---

# Architecture Review Discipline

**Architecture is always more important than implementation.**

This gate applies when creating any new abstraction. Complete all phases in order. If any phase reveals a problem, **STOP** and explain the issue.

---

## When This Gate Applies

Complete full architecture review before creating:

- new package
- new repository
- new service
- new engine
- new adapter
- new projection
- new public interface
- new composition root
- new interface
- new DTO or intermediate model
- new exported type
- new manager

**Only modifying existing components? Gate does not apply.**

---

## Core Principle: Extract, Don't Invent

Abstractions **MUST** be extracted from existing code or approved RFCs. Never invented from anticipated future needs.

### Validity Criteria

A new abstraction **MUST** satisfy at least one of these conditions:

1. **Two or more independent implementations already exist** in the current codebase
2. **Duplicated behavior already exists** and extracting it measurably reduces complexity
3. **An approved RFC explicitly defines the abstraction** by name and contract

### Prohibited Justifications

- Future integrations (Kafka, Slack, Email, Postgres, Redis, etc.) that do not yet exist
- Speculative extensibility ("we may need this later")
- Anticipated requirements without an approved RFC
- Theoretical flexibility without a current concrete use

**"We may need it later" is never sufficient justification.**

---

## 11-Phase Architecture Review

### Phase 1: RFC Review

**Objective:** Establish architectural foundation and constraints.

**Deliverables:**

- List all relevant RFCs (by number and title)
- Quote specific sections that govern this work
- Identify any RFC ambiguities, gaps, or contradictions
- Confirm no RFC prohibits proposed approach

**Stop condition:** If RFCs conflict or are ambiguous, stop here.

### Phase 2: Existing Architecture Review

**Objective:** Understand what already exists to avoid duplication.

**Deliverables:**

- List existing services, packages, repositories, engines, adapters
- List existing events, commands, queries, aggregates, projections
- List existing public interfaces and exported types
- Identify reusable components

**Stop condition:** If existing component satisfies need, stop and reuse it.

### Phase 3: Responsibility Analysis

**Objective:** Define what new component owns and does not own.

**Deliverables:**

For proposed component, explicitly state:

- What it owns (data, logic, behavior)
- What it does NOT own (delegated responsibilities)
- What it depends on
- What depends on it

**Stop condition:** If responsibilities overlap with existing components, stop and clarify.

### Phase 4: Ownership Analysis

**Objective:** Confirm single ownership for every capability.

**Deliverables:**

- Identify single owner of each responsibility
- Confirm no shared ownership
- Confirm no orphaned responsibilities

**Stop condition:** If ownership is ambiguous or duplicated, stop and resolve.

### Phase 5: Layer Validation

**Objective:** Enforce clean architecture boundaries.

**Deliverables:**

Classify component into exactly one layer:

```
Domain          ← business logic, entities, value objects
  ↓
Application     ← use cases, orchestration, commands, queries
  ↓
Infrastructure  ← storage, transport, adapters, external systems
  ↓
Composition Root ← wiring, dependency injection, main()
```

Confirm:

- Infrastructure never depends on Domain
- Domain has zero infrastructure imports
- Application depends only on Domain interfaces
- Only Composition Root wires unrelated layers

**Stop condition:** If component violates layer boundaries, stop and redesign.

### Phase 6: Component Categorization

**Objective:** Classify component into exactly one architectural category.

**Deliverables:**

Assign component to exactly one category:

- **Repository** — storage interface, no business logic
- **Engine** — stateless logic processor, no storage
- **Reducer** — event → state transformation, pure function
- **Coordinator** — orchestrates across services, no domain logic
- **Adapter** — translates between layers, no domain logic
- **Composition Root** — wiring only, no behavior

Reject designs mixing categories.

Invalid designs:

- Repository containing business logic ❌
- Engine storing state ❌
- Reducer calling external services ❌
- Coordinator implementing domain logic ❌
- Adapter making decisions ❌

**Stop condition:** If component mixes categories, stop and separate concerns.

### Phase 7: Storage Rules Validation

**Objective:** Enforce storage owns persistence only.

**Deliverables:**

If component is storage-related, confirm it does NOT own:

- ❌ Business logic
- ❌ Replay logic
- ❌ Orchestration
- ❌ Reducers
- ❌ Workflows
- ❌ Event sourcing logic

Storage responsibilities:

- ✅ Persist data
- ✅ Retrieve data
- ✅ Provide transactions
- ✅ Enforce referential integrity

**Projection-specific rules:**

- Projection Repository stores **snapshots only**
- Replay belongs to **Replay Engine**
- Projection construction belongs to **Reducers**
- Checkpoint management belongs to **Checkpoint Store**

**Stop condition:** If storage violates these rules, stop and extract logic.

### Phase 8: Dependency Graph Validation

**Objective:** Visualize and validate dependencies.

**Deliverables:**

Produce dependency graph showing:

- New component
- What it imports
- What imports it
- Transitive dependencies

Validate:

- No circular dependencies
- No cross-layer violations
- No hidden coupling
- Dependencies flow in one direction

**Stop condition:** If dependency graph reveals violations, stop and resolve.

### Phase 9: Runtime Data Flow

**Objective:** Trace how data moves through system at runtime.

**Deliverables:**

Produce data flow diagram showing:

- Request entry point
- Data transformations
- Service calls
- Storage operations
- Response path

Example:

```
HTTP Request
  → Command Handler
    → Aggregate (load from EventStore)
    → Aggregate.Apply(command)
    → Emit Event
    → EventStore.Append(event)
    → EventBus.Publish(event)
  → HTTP Response

Async:
  EventBus
    → Projection Reducer
      → Projection Repository.Save(snapshot)
```

Validate:

- Data flow matches RFC specifications
- No bypasses of architectural boundaries
- State changes are auditable
- Events drive projections (not direct writes)

**Stop condition:** If data flow violates invariants, stop and redesign.

### Phase 10: Public API Review

**Objective:** Minimize surface area and prevent leaky abstractions.

**Deliverables:**

For every exported function, type, interface, or method:

- Justify why it must be public
- Confirm required by immediate production caller
- Confirm does not expose internal implementation details
- Confirm follows RFC naming conventions

Reject:

- Speculative abstractions (no caller exists yet)
- Generic frameworks ("we might need this later")
- "Future-proof" interfaces (implement YAGNI)
- Exported internals (make them package-private)

**Stop condition:** If public API includes speculative elements, stop and reduce surface area.

### Phase 11: Abstraction Review — Final Validation

**Objective:** Enforce Extract, Don't Invent for every new abstraction.

**Deliverables:**

For every new interface, DTO, repository, service, adapter, engine, manager, intermediate model, or public API, answer all:

| Question | Answer |
|---|---|
| Why does this abstraction exist today? | _(required)_ |
| Which existing duplication does it remove? | _(required)_ |
| How many concrete implementations currently exist? | _(required)_ |
| Which approved RFC requires it? | _(required)_ |
| Can a concrete implementation be used instead? | _(required)_ |
| Would removing this abstraction simplify the architecture? | _(required)_ |

**Validity check:** At least one must be true:

- Two or more independent implementations already exist in codebase
- Duplicated behavior exists and extraction reduces complexity
- An approved RFC explicitly defines this abstraction

**Stop condition:** If any abstraction cannot be justified, **STOP.** The abstraction **MUST** be removed before implementation proceeds.

---

## After Architecture Review Complete

All 11 phases passed. You have explicit architectural approval to proceed to implementation.

Document findings for the final report:

- RFC compliance
- Architectural decisions made
- Alternative designs considered and rejected
- Risks and mitigations

