---
description: "Mandatory architecture review gate before ANY implementation of new abstractions in Praxis. Enforces RFC fidelity, clean architecture boundaries, component categorization, and justification requirements. Architecture is always more important than implementation."
name: "Praxis Architecture Guardian"
applyTo: ["services/**", "packages/**", "apps/**", "scripts/**", "infra/**"]
---

# Praxis Architecture Guardian

**Architecture is always more important than implementation.**

This instruction acts as a mandatory architecture review gate. You **MUST** complete this
review before implementing any new abstraction. Never skip. Never assume approval.

---

## When This Gate Applies

You must complete this full architecture review before creating any:

- new package
- new repository
- new engine
- new runtime
- new storage mechanism
- new adapter
- new public interface
- new exported type

**If you are only modifying existing components, this gate does not apply.**

---

## Core Principle: RFCs Are Architectural Contracts

RFCs in `./rfcs` are **immutable architectural contracts**.

**Never:**

- Reinterpret them
- Extend them with assumptions
- Fill gaps with "reasonable" guesses
- Implement behavior not explicitly specified

**If an RFC is ambiguous, incomplete, or contradictory:**

- **STOP immediately**
- Document the specific ambiguity
- Explain what cannot be determined
- Wait for clarification

Do not proceed with implementation when the architecture is unclear.

---

## 10-Phase Architecture Review

Complete all phases in order. Document your findings. If any phase reveals a problem,
**STOP** and explain the issue.

### Phase 1: RFC Review

**Objective:** Establish architectural foundation and constraints.

**Deliverables:**

- List all relevant RFCs (by number and title)
- Quote the specific sections that govern this work
- Identify any RFC ambiguities, gaps, or contradictions
- Confirm no RFC prohibits the proposed approach

**Stop condition:** If RFCs conflict or are ambiguous, stop here.

### Phase 2: Existing Architecture Review

**Objective:** Understand what already exists to avoid duplication.

**Deliverables:**

- List existing services, packages, repositories, engines, adapters
- List existing events, commands, queries, aggregates, projections
- List existing public interfaces and exported types
- Identify reusable components

**Stop condition:** If an existing component satisfies the need, stop and reuse it.

### Phase 3: Responsibility Analysis

**Objective:** Define what the new component owns and does not own.

**Deliverables:**

For the proposed component, explicitly state:

- What it owns (data, logic, behavior)
- What it does NOT own (delegated responsibilities)
- What it depends on
- What depends on it

**Stop condition:** If responsibilities overlap with existing components, stop and clarify.

### Phase 4: Ownership Analysis

**Objective:** Confirm single ownership for every capability.

**Deliverables:**

- Identify the single owner of each responsibility
- Confirm no shared ownership
- Confirm no orphaned responsibilities

**Stop condition:** If ownership is ambiguous or duplicated, stop and resolve.

### Phase 5: Layer Validation

**Objective:** Enforce clean architecture boundaries.

**Deliverables:**

Classify the component into exactly one layer:

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

**Stop condition:** If the component violates layer boundaries, stop and redesign.

### Phase 6: Component Categorization

**Objective:** Classify the component into exactly one architectural category.

**Deliverables:**

Assign the component to exactly one category:

- **Repository** — storage interface, no business logic
- **Engine** — stateless logic processor, no storage
- **Reducer** — event → state transformation, pure function
- **Coordinator** — orchestrates across services, no domain logic
- **Adapter** — translates between layers, no domain logic
- **Composition Root** — wiring only, no behavior

**Reject designs that mix categories.**

Examples of invalid designs:

- Repository that contains business logic ❌
- Engine that stores state ❌
- Reducer that calls external services ❌
- Coordinator that implements domain logic ❌
- Adapter that makes decisions ❌

**Stop condition:** If the component mixes categories, stop and separate concerns.

### Phase 7: Storage Rules Validation

**Objective:** Enforce storage owns persistence only.

**Deliverables:**

If the component is storage-related, confirm it does NOT own:

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

**Stop condition:** If storage violates these rules, stop and extract the logic.

### Phase 8: Dependency Graph

**Objective:** Visualize and validate dependencies.

**Deliverables:**

Produce a dependency graph showing:

- The new component
- What it imports
- What imports it
- Transitive dependencies

Use ASCII, Mermaid, or text format.

Example:

```
UserService
  ↓ depends on
UserRepository (interface)
  ↑ implemented by
PostgresUserRepository
  ↓ depends on
pgx (external)
```

Validate:

- No circular dependencies
- No cross-layer violations
- No hidden coupling
- Dependencies flow in one direction

**Stop condition:** If dependency graph reveals violations, stop and resolve.

### Phase 9: Runtime Data Flow

**Objective:** Trace how data moves through the system at runtime.

**Deliverables:**

Produce a data flow diagram showing:

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
- Confirm it is required by an immediate production caller
- Confirm it does not expose internal implementation details
- Confirm it follows RFC naming conventions

**Reject:**

- Speculative abstractions (no caller exists yet)
- Generic frameworks ("we might need this later")
- "Future-proof" interfaces (implement YAGNI)
- Exported internals (make them package-private)

**Stop condition:** If public API includes speculative elements, stop and reduce surface area.

---

## Justification Requirements

For every new component, provide explicit justification:

### New Package

- Why existing packages cannot contain this?
- What boundary does this enforce?
- Which RFC mandates this separation?

### New Repository

- What aggregate or entity does it persist?
- Why existing repositories cannot handle this?
- Which RFC defines this storage boundary?

### New Engine

- What stateless transformation does it perform?
- Why existing engines cannot handle this?
- Which RFC defines this processing?

### New Runtime

- What execution model does it provide?
- Why existing runtimes cannot support this?
- Which RFC requires this runtime?

### New Storage

- What data does it persist?
- Why existing storage cannot handle this?
- Which RFC mandates this storage?

### New Adapter

- What external system does it integrate?
- Why existing adapters cannot support this?
- Which RFC requires this integration?

### New Public Interface

- What contract does it define?
- Which production caller requires it?
- Why must it be public?

### New Exported Type

- Why must this be exported?
- Which external package imports it?
- Can it be package-private instead?

**No component may exist without justification.**

---

## Architecture Self-Review

Before finalizing your review, ask yourself:

### Simplicity

- Is this the simplest design that satisfies the RFCs?
- Can I remove any abstraction?
- Can I reuse instead of create?

### Boundaries

- Does this respect clean architecture layers?
- Does this belong to exactly one category?
- Does this violate any ownership rules?

### Future Impact

- Does this introduce second sources of truth?
- Does this create hidden coupling?
- Does this weaken auditability?

### RFC Fidelity

- Does this implement exactly what the RFC specifies?
- Does this add behavior not in the RFC?
- Does this contradict any RFC?

### Testability

- Can this be tested in isolation?
- Can this be verified automatically?
- Does this require manual validation?

---

## Required Deliverables

Before proceeding to implementation, produce:

### 1. RFC Summary

- List of relevant RFCs
- Quoted sections that govern this work
- Any ambiguities or gaps identified

### 2. Ownership Analysis

- Single owner for each responsibility
- Clear delegation boundaries
- No overlapping ownership

### 3. Dependency Graph

- ASCII or Mermaid diagram
- All dependencies shown
- Direction validated

### 4. Runtime Data Flow

- Request → Response path
- Data transformations
- Storage operations
- Event flows

### 5. Architecture Debt Section

Explicitly document:

- Shortcuts taken (if any)
- Technical debt introduced (if any)
- Future refactoring required (if any)
- RFC gaps this work exposed (if any)

Even if none, state: "No architecture debt introduced."

### 6. Smallest Vertical Slice

Define the minimal end-to-end implementation:

- Entry point
- Core logic
- Storage
- Exit point
- Verification

Ship one working path, not multiple partial paths.

---

## Stop-Before-Code Rule

**After completing the 10-phase review and producing all deliverables:**

**STOP.**

**Do not proceed to code generation.**

**Wait for explicit approval.**

Present your architecture review and ask:

> I have completed the architecture review for [component name].
>
> The review identified [X RFCs], [Y dependencies], [Z risks].
>
> All deliverables are documented above.
>
> Should I proceed with implementation?

Only after receiving explicit approval ("yes", "proceed", "go ahead") may you generate code.

**Never assume approval.**

**Never implement first and justify later.**

**Never skip the stop-and-wait step.**

---

## Final Checklist

Before claiming the architecture review is complete, confirm:

- [ ] All 10 phases completed
- [ ] All deliverables produced
- [ ] All justifications documented
- [ ] No RFC ambiguities unresolved
- [ ] No architectural violations identified
- [ ] Dependency graph validated
- [ ] Data flow validated
- [ ] Public API minimized
- [ ] Architecture debt documented
- [ ] Smallest vertical slice defined
- [ ] **STOPPED and waiting for approval**

If any checkbox is unchecked, the review is incomplete.

---

## Remember

**Architecture is always more important than implementation.**

A wrong implementation can be rewritten.

A wrong architecture pollutes the codebase permanently.

Take the time to get architecture right before writing a single line of code.

Never skip the architecture review.