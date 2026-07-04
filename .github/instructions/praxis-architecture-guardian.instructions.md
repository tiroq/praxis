---
description: "Mandatory architecture review gate before ANY implementation of new abstractions in Praxis. Enforces RFC fidelity, clean architecture boundaries, component categorization, abstraction justification (Extract Don't Invent), and justification requirements. Architecture is always more important than implementation."
name: "Praxis Architecture Guardian"
applyTo: ["services/**", "packages/**", "apps/**", "scripts/**", "infra/**"]
---

# Praxis Architecture Guardian

**Architecture is always more important than implementation.**

Mandatory architecture review gate. **MUST** complete before implementing any new abstraction. Never skip. Never assume approval.

---

## When This Gate Applies

Complete full architecture review before creating any:

- new package
- new repository
- new engine
- new runtime
- new storage mechanism
- new adapter
- new public interface
- new exported type
- new interface
- new DTO or intermediate model
- new service
- new manager

**Only modifying existing components? Gate does not apply.**

---

## Core Principle: RFCs Are Architectural Contracts

RFCs in `./rfcs` are **immutable architectural contracts**.

**Never:**

- Reinterpret them
- Extend with assumptions
- Fill gaps with "reasonable" guesses
- Implement behavior not explicitly specified

**If RFC is ambiguous, incomplete, or contradictory:**

- **STOP immediately**
- Document specific ambiguity
- Explain what cannot be determined
- Wait for clarification

Do not proceed when architecture is unclear.

---

## Core Principle: Extract, Don't Invent

Abstractions **MUST** be extracted from existing code or approved RFCs. Abstractions **MUST NOT** be invented from anticipated future needs.

This rule applies to every abstraction in Praxis without exception:

- interfaces
- DTOs and intermediate models
- repositories
- services
- adapters
- engines
- managers
- public APIs

### Validity Criteria

A new abstraction **MUST** satisfy at least one of the following conditions:

1. **Two or more independent implementations already exist** in the current codebase.
2. **Duplicated behavior already exists** and extracting it measurably reduces complexity.
3. **An approved RFC explicitly defines the abstraction** by name and contract.

### Prohibited Justifications

The following are **NEVER** valid justifications for introducing an abstraction:

- Future integrations (Kafka, Slack, Email, Postgres, Redis, etc.) that do not yet exist
- Speculative extensibility ("we may need this later")
- Anticipated requirements without an approved RFC
- Theoretical flexibility without a current concrete use

**"We may need it later" is never sufficient justification.**

### Presumption of Concreteness

- Favor concrete implementations until duplication actually appears in the codebase.
- The burden of proof belongs to the author introducing the abstraction.
- If removing the abstraction simplifies the architecture, remove it.

### Review Obligation

During architecture review, the agent **MUST** explicitly justify every newly introduced abstraction using the Abstraction Review checklist (Phase 11).

If the justification cannot be proven using the current codebase or approved RFCs, the abstraction **MUST** be removed before implementation proceeds.

---

## 11-Phase Architecture Review

Complete all phases in order. Document findings. If any phase reveals problem, **STOP** and explain issue.

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

**Reject designs mixing categories.**

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

### Phase 8: Dependency Graph

**Objective:** Visualize and validate dependencies.

**Deliverables:**

Produce dependency graph showing:

- New component
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

**Reject:**

- Speculative abstractions (no caller exists yet)
- Generic frameworks ("we might need this later")
- "Future-proof" interfaces (implement YAGNI)
- Exported internals (make them package-private)

**Stop condition:** If public API includes speculative elements, stop and reduce surface area.

### Phase 11: Abstraction Review

**Objective:** Enforce the Extract, Don't Invent principle for every new abstraction.

**Deliverables:**

For every new interface, DTO, repository, service, adapter, engine, manager, intermediate model, or public API introduced, answer all of the following:

| Question | Answer |
|---|---|
| Why does this abstraction exist today? | _(required)_ |
| Which existing duplication does it remove? | _(required)_ |
| How many concrete implementations currently exist? | _(required)_ |
| Which approved RFC requires it? | _(required)_ |
| Can a concrete implementation be used instead? | _(required)_ |
| Would removing this abstraction simplify the architecture? | _(required)_ |

**Validity check:** The abstraction is justified only if at least one of the following is true:

- Two or more independent implementations already exist in the codebase
- Duplicated behavior already exists and extraction reduces complexity
- An approved RFC explicitly defines this abstraction

**Stop condition:** If any abstraction cannot be justified by answering these questions satisfactorily using the current codebase or approved RFCs, **STOP**. The abstraction **MUST** be removed before implementation proceeds. Do not proceed until explicit approval is received.

---

## Pure Mapping Verification

**Mandatory rule for all transport mapping functions (adapters).**

Every mapping function that translates an external transport object to a Praxis wire contract **MUST** be a pure function with no side effects, no external dependencies, and no business logic.

### Mapper Shape (Required)

```
Input:  external transport object (e.g., Telegram Update, HTTP request, AMQP message)
Output: Praxis wire contract dict (matching internal/transport/nats.InputMessage)
```

### Mapper MUST NOT

- ❌ Publish messages
- ❌ Acknowledge messages
- ❌ Call HTTP endpoints
- ❌ Call NATS
- ❌ Call storage (databases, caches, files)
- ❌ Access environment variables
- ❌ Perform retries
- ❌ Mutate global state
- ❌ Generate business identifiers
- ❌ Execute business logic
- ❌ Classify intent
- ❌ Infer semantics
- ❌ Enrich domain objects
- ❌ Create workflows

### Forbidden Complexity in Mappers

A mapper MUST NOT:

- Call repositories
- Call databases
- Call HTTP services
- Call LLMs
- Call event stores
- Call action planners
- Call reviewers
- Call decision makers
- Perform retries
- Perform caching
- Maintain state
- Generate business identifiers
- Compute business outcomes
- Infer intent
- Classify content
- Enrich with external knowledge
- Coordinate workflows

If any of these become necessary, stop and move the logic into the correct owner.

### Allowed Complexity

A mapper MAY ONLY:

- ✅ Rename fields
- ✅ Copy fields
- ✅ Drop unused fields
- ✅ Normalize formatting
- ✅ Convert primitive types
- ✅ Construct transport DTOs
- ✅ Perform deterministic serialization
- ✅ Validate transport-level preconditions
- ✅ Create deterministic identifiers derived solely from the input

Nothing else.

### Referential Transparency Invariant

A mapper must be **referentially transparent**:

> Given the same input object, it MUST always produce the same output object.
>
> No observable side effects are permitted.

### Verification Checklist (Phase 11 Extension)

For every adapter with a mapping function, answer:

| Question | Answer |
|---|---|
| Does the mapper accept only the external transport object? | _(required)_ |
| Does the mapper return only a dict (wire contract)? | _(required)_ |
| Does the mapper call any I/O function (publish, HTTP, storage)? | _(required: must be NO)_ |
| Does the mapper access any global state? | _(required: must be NO)_ |
| Does the mapper perform any retry logic? | _(required: must be NO)_ |
| Can this mapper be called 1,000,000 times with the same input and produce identical output? | _(required: must be YES)_ |
| Is the mapper completely free of business logic? | _(required: must be YES)_ |

**Stop condition:** If any answer fails purity verification, the mapper **MUST** be refactored before implementation proceeds.

### Structural Triviality Test

Review every mapper using these questions:

1. Can every output field be traced directly to one or more input fields?
2. Is every transformation deterministic?
3. Would another engineer understand the mapper in under one minute?
4. Could the mapper be rewritten as a simple table of field mappings?
5. If removed, would business behavior remain unchanged?

If any answer is "No", the mapper likely owns behavior it should not own.

Stop implementation and perform an architecture review.

### Extraction Rule

A mapper is not a reusable abstraction.

Do not extract mapper frameworks, mapper interfaces, generic converters, transformation pipelines, or reusable mapping engines unless duplication already exists.

Prefer explicit mapping functions.

Extract only after repeated duplication.

### Wire Contract Rule

Transport mappers own only the translation into the published wire contract.

They do not own the contract itself.

Changing a mapper must never silently change a transport contract.

Transport contracts remain owned by RFCs.

### Mapper Size Guideline

Mapper size is not enforced by line count.

Instead evaluate responsibility.

A mapper should normally:

- perform one transformation
- have one input
- produce one output
- contain no nested decision trees
- remain understandable in a single screen

If a mapper requires sections, regions, helper classes, or extensive comments, architecture review is required.

### Behavior Accumulation Test

Every mapper review must answer:

1. Has this mapper become more complex than the transport object it maps?
2. Are new business cases being added to the mapper over time?
3. Is this becoming the easiest place to add unrelated logic?
4. Would another subsystem reasonably want to reuse this logic?
5. Does this logic have a better owner?

If any answer is "Yes", stop implementation and perform an architecture review.

### One Translation Rule

A mapper performs exactly one translation.

Examples:

```
Telegram Update → Praxis InputMessage

HTTP Request → Command DTO

Database Row → Projection DTO
```

A mapper must never chain multiple conceptual transformations.

If two translations exist, they should be two mappers owned by their respective boundaries.

### No Helper Gravity

Do not slowly convert mappers into utility modules.

Avoid adding:

- helper classes
- shared mapping utilities
- generic transformation libraries
- base mapper classes
- mapper inheritance
- mapper registries

Extract only after demonstrated duplication.

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

Before finalizing review, ask:

### Simplicity

- Is this simplest design satisfying RFCs?
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

- Does this implement exactly what RFC specifies?
- Does this add behavior not in RFC?
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
- Quoted sections governing this work
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

Define minimal end-to-end implementation:

- Entry point
- Core logic
- Storage
- Exit point
- Verification

Ship one working path, not multiple partial paths.

---

## Stop-Before-Code Rule

**After completing 10-phase review and producing all deliverables:**

**STOP.**

**Do not proceed to code generation.**

**Wait for explicit approval.**

Present architecture review and ask:

> I have completed architecture review for [component name].
>
> Review identified [X RFCs], [Y dependencies], [Z risks].
>
> All deliverables documented above.
>
> Should I proceed with implementation?

Only after receiving explicit approval ("yes", "proceed", "go ahead") may you generate code.

**Never assume approval.**

**Never implement first and justify later.**

**Never skip stop-and-wait step.**

---

## Final Checklist

Before claiming architecture review complete, confirm:

- [ ] All 11 phases completed
- [ ] All deliverables produced
- [ ] All justifications documented
- [ ] No RFC ambiguities unresolved
- [ ] No architectural violations identified
- [ ] Dependency graph validated
- [ ] Data flow validated
- [ ] Public API minimized
- [ ] Architecture debt documented
- [ ] Smallest vertical slice defined
- [ ] **Every new abstraction justified via Phase 11 Abstraction Review**
- [ ] **No abstraction introduced on "we may need it later" grounds**
- [ ] **If adapter: mapper purity verification completed (or no mapping function required)**
- [ ] **Mapper is referentially transparent**
- [ ] **Mapper is structurally trivial**
- [ ] **Every output field is traceable to input**
- [ ] **No business decisions exist**
- [ ] **No infrastructure calls exist**
- [ ] **No state is maintained**
- [ ] **Mapper performs exactly one translation**
- [ ] **Mapper has not accumulated unrelated behavior**
- [ ] **Mapper remains the obvious place only for representation translation**
- [ ] **Architecture review completed if mapper scope expanded**
- [ ] **No speculative abstractions were introduced**
- [ ] **STOPPED and waiting for approval**

If any checkbox unchecked, review is incomplete.

---

## Architecture Laws

These laws are non-negotiable and permanent:

**Mappers are pure.**

Transport mapping functions contain no side effects, no business logic, no external dependencies.

**Adapters translate only.**

Adapters bridge external systems to internal wire contracts. They do not decide, enrich, or reason about data.

**Business logic belongs outside transport mappers.**

Validation, enrichment, intent classification, and workflow creation happen in engines, reducers, or coordinators—never in mappers.

**Mappers must remain structurally trivial.**

A mapper exists only to translate one representation into another.

A mapper must never become a place where policy accumulates.

If a mapper begins making business decisions, branching on domain state, or coordinating behavior, the logic belongs elsewhere.

**Mapper growth is an architectural smell.**

A mapper that continuously grows is evidence that responsibility is misplaced.

Large mappers are architectural smells, not implementation achievements.

When a mapper keeps accumulating rules, stop extending it and identify the missing owner.

---

## Remember

**Architecture is always more important than implementation.**

Wrong implementation can be rewritten.

Wrong architecture pollutes codebase permanently.

Take time to get architecture right before writing single line of code.

Never skip architecture review.