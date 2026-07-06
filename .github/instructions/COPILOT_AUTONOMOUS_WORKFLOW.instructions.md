---
name: "Copilot Autonomous Workflow"
description: "Orchestration layer for autonomous implementation in Praxis. Coordinates 11 phases from discovery through validation and final report. Delegates rule enforcement to specialized instruction files."
applyTo: ["services/**", "packages/**", "apps/**", "scripts/**", "infra/**"]
---

# Copilot Autonomous Workflow — Orchestration Layer

**Default Assumption:** The user expects complete implementation, not just analysis or planning.

This workflow orchestrates autonomous implementation from discovery through final report. Execute all phases automatically without stopping for intermediate confirmations unless genuinely blocked by missing information.

---

# 🎯 THE NINE-PHASE AUTONOMOUS WORKFLOW

## Phase 1: Discover — RFCs, Architecture, and Existing Code

**Do this first.** Never start coding without understanding what already exists.

### 1.1 RFC Discovery

- List all RFCs that might apply to this task
- Read each RFC completely
- Quote specific sections governing this work
- Identify RFC ambiguities, gaps, or contradictions
- **STOP if RFCs conflict.** Surface the conflict and wait for clarification

**Tools:** Use `grep` or file search for keywords. Most RFCs are in `./rfcs/`.

### 1.2 Architecture Context

- Read relevant ADRs (in `./docs/adr/`)
- Read relevant Policies
- Identify all architectural invariants that must not be violated (see Invariants section below)
- Identify project instructions (in `./.github/instructions/`)

### 1.3 Reference Implementation Search

Use Graphify first, then code search:

- **Graphify query** — "show me existing [service/adapter/repository/mapper]"
- **Code search** — grep for existing patterns, services, engines, adapters

Categories to search for reusable implementations:

- mapper (especially transport adapters)
- repository
- worker
- adapter
- storage backend
- transport
- bootstrap
- configuration loader
- service
- reducer
- projection
- event
- command
- query
- aggregate

If a reference implementation exists → **prefer evolving it** over building new.

### 1.4 Graphify for Architecture Questions

For any question about repository architecture, ownership, dependencies, or implementation locations:

**Use Graphify first** before manual code search.

- `graphify query` — semantic search
- `graphify path` — dependency relationships
- `graphify explain` — architectural concepts

Only inspect source files when implementation is required, debugging, or Graphify lacks detail.

---

## Phase 2: Architectural Feasibility Gate — Is New Abstraction Needed?

**MANDATORY gate before creating anything new.**

Stop here if you're only modifying existing components. Gate only applies when creating:

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

### 2.1 Check: Does an RFC Explicitly Require This?

Read the relevant RFC. Does it name and define the abstraction you plan to create?

- **YES** → continue to 2.2
- **NO** → continue to 2.2 (may still be justified)

### 2.2 Extract, Don't Invent: Apply the Three-Point Rule

A new abstraction is justified **ONLY IF** at least one of these is true:

1. **Two or more independent implementations already exist** in the codebase
2. **Duplicated behavior already exists** and extracting it measurably reduces complexity
3. **An approved RFC explicitly defines the abstraction** by name and contract

**Prohibited justifications:**

- "We may need this later"
- "Future integrations (Kafka, Slack, Email, etc.) we don't have yet"
- "Speculative extensibility"
- "Theoretical flexibility without current concrete use"

**If none of these three apply** → do not create the abstraction. Use concrete implementation instead.

### 2.3 Abstraction Validation Checklist

For every new abstraction you plan to introduce, answer all of these:

| Question | Answer |
|---|---|
| Why does this abstraction exist today (not "why might we need it")? | _(required)_ |
| Which existing duplication does it remove? | _(required)_ |
| How many concrete implementations currently exist? | _(required)_ |
| Which approved RFC requires it by name? | _(required)_ |
| Can a concrete implementation satisfy the need instead? | _(required)_ |
| Would removing this abstraction make the code simpler? | _(required)_ |

**Stop condition:** If you cannot answer all these questions using the current codebase or approved RFCs, **STOP.** Do not proceed. The abstraction must be removed before implementation continues.

---

## Phase 3: Complete Architecture Review (If Gate Applied)

**Only if Phase 2 determined a new abstraction is needed.** Complete all 11 phases below. Do not skip. Do not assume approval.

### 3.1 RFC Review

- List all relevant RFCs (by number and title)
- Quote specific sections that govern this work
- Identify any RFC ambiguities, gaps, or contradictions
- Confirm no RFC prohibits proposed approach

**Stop if RFCs conflict or are ambiguous.**

### 3.2 Existing Architecture Review

- List existing services, packages, repositories, engines, adapters
- List existing events, commands, queries, aggregates, projections
- List existing public interfaces and exported types
- Identify what could be reused or extended

**Stop if existing component satisfies need.** Reuse it instead.

### 3.3 Responsibility Analysis

For the proposed component, explicitly state:

- What it owns (data, logic, behavior)
- What it does NOT own (delegated responsibilities)
- What it depends on
- What depends on it

**Stop if responsibilities overlap with existing components.**

### 3.4 Ownership Analysis

- Identify single owner of each responsibility
- Confirm no shared ownership
- Confirm no orphaned responsibilities

**Stop if ownership is ambiguous or duplicated.**

### 3.5 Layer Validation

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

**Stop if component violates layer boundaries.**

### 3.6 Component Categorization

Assign component to exactly one category:

- **Repository** — storage interface, no business logic
- **Engine** — stateless logic processor, no storage
- **Reducer** — event → state transformation, pure function
- **Coordinator** — orchestrates across services, no domain logic
- **Adapter** — translates between layers, no domain logic
- **Composition Root** — wiring only, no behavior

Reject designs mixing categories.

**Stop if component mixes categories.**

### 3.7 Storage Rules Validation

If component is storage-related, confirm it does NOT own:

- ❌ Business logic
- ❌ Replay logic
- ❌ Orchestration
- ❌ Reducers
- ❌ Workflows
- ❌ Event sourcing logic

Storage must only:

- ✅ Persist data
- ✅ Retrieve data
- ✅ Provide transactions
- ✅ Enforce referential integrity

Projection-specific rules:

- Projection Repository stores **snapshots only**
- Replay belongs to **Replay Engine**
- Projection construction belongs to **Reducers**
- Checkpoint management belongs to **Checkpoint Store**

**Stop if storage violates these rules.**

### 3.8 Dependency Graph Validation

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

**Stop if dependency graph reveals violations.**

### 3.9 Runtime Data Flow

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

**Stop if data flow violates invariants.**

### 3.10 Public API Review

For every exported function, type, interface, or method:

- Justify why it must be public
- Confirm required by immediate production caller
- Confirm does not expose internal implementation details
- Confirm follows RFC naming conventions

Reject:

- Speculative abstractions (no caller exists yet)
- Generic frameworks ("we might need this later")
- "Future-proof" interfaces (implement YAGNI)
- Exported internals (make package-private)

**Stop if public API includes speculative elements.**

### 3.11 Abstraction Review — Final Validation

For every new interface, DTO, repository, service, adapter, engine, manager, or public API, answer all:

| Question | Answer |
|---|---|
| Why does this abstraction exist today? | _(required)_ |
| Which existing duplication does it remove? | _(required)_ |
| How many concrete implementations currently exist? | _(required)_ |
| Which approved RFC requires it? | _(required)_ |
| Can a concrete implementation be used instead? | _(required)_ |
| Would removing this abstraction simplify the architecture? | _(required)_ |

**Validity:** At least one of these MUST be true:

- Two or more independent implementations already exist in codebase
- Duplicated behavior already exists and extraction reduces complexity
- An approved RFC explicitly defines this abstraction

**Stop condition:** If any abstraction cannot be justified, **STOP.** The abstraction **MUST** be removed before implementation proceeds. Do not proceed until explicit approval is received.

**After completion of all 11 phases:** You have explicit architectural approval to proceed to implementation.

---

## Phase 4: Implementation Planning

**Do not write code yet.** Produce implementation plan first.

Plan must include:

- **Relevant RFCs** — specific RFCs this traces back to
- **Impacted invariants** — which hard rules (see Invariants section) this touches
- **Existing components to reuse** — services, events, commands, queries, aggregates, projections, adapters, storage already defined and applicable
- **New components to create** — only what reuse cannot cover; each justified by RFC or architecture review gate
- **Minimal implementation slice** — smallest end-to-end change delivering value
- **Verification strategy** — tests/commands that will validate behavior
- **Risks** — what could break, ambiguous RFC areas, second-source-of-truth dangers

**Cannot map work to RFC?** Stop and ask. Do not invent behavior.

---

## Phase 5: Implementation — Build Minimal Vertical Slice

**Do not build entire subsystems at once.** Implement smallest end-to-end slice providing value, typically:

```
Command → Event → Projection → Query → Verification
```

### 5.1 Reuse-Before-Build Rule

Before creating anything new, search for existing:

- service
- event
- command
- query
- aggregate
- projection
- agent
- storage
- adapter

Extend or compose what exists. Only create new when reuse genuinely impossible. Explain why if adding new.

### 5.2 Implementation Rules

- Reuse existing code
- Reuse existing packages
- Reuse existing architecture
- Prefer deletion over addition
- Prefer concrete implementations
- Avoid speculative abstractions

### 5.3 Vertical Slice Structure

Build through layers in order:

1. RFC-031 — service envelopes (`rfcs/031-service-contracts.md`)
2. RFC-013 — event model (`rfcs/013-event-model.md`)
3. RFC-014 — identity/representation model (`rfcs/014-identity-representation-model.md`)
4. RFC-015 — object lifecycle model (`rfcs/015-object-lifecycle-model.md`)
5. RFC-020 — review system (`rfcs/020-review-system.md`)
6. RFC-021 — decision model (`rfcs/021-decision-model.md`)
7. RFC-023 — action model (`rfcs/023-action-model.md`)

Do not implement later layer before prerequisites exist and are tested.

### 5.4 Transport Mappers — Pure Function Rule

**If implementing a transport adapter with a mapping function:**

Mapping function must be:

- ✅ Pure (no side effects, no external dependencies)
- ✅ Single responsibility (exactly one translation)
- ✅ Business-logic free
- ✅ Infrastructure-logic free
- ✅ Simple and understandable in under one minute

Mapper MUST NOT:

- ❌ Call repositories or databases
- ❌ Call HTTP services
- ❌ Call LLMs or external APIs
- ❌ Publish or acknowledge messages
- ❌ Access environment variables
- ❌ Perform retries
- ❌ Mutate global state
- ❌ Generate business identifiers
- ❌ Execute business logic
- ❌ Classify intent or infer semantics

Mapper MAY:

- ✅ Rename and copy fields
- ✅ Drop unused fields
- ✅ Normalize formatting
- ✅ Convert primitive types
- ✅ Create deterministic identifiers derived solely from input
- ✅ Validate transport-level preconditions

Reference: Review `docs/architecture/GOLDEN_MAPPER.md` before writing any mapper.

**Referential Transparency:** Given the same input, mapper must always produce identical output.

---

## Phase 6: Continuous Validation — Check RFC and Architecture Compliance

**During implementation,** before committing code, answer each. If any is **YES**, stop and fix:

- Does this violate any RFC?
- Does this duplicate existing concept?
- Does this introduce second source of truth?
- Does this bypass Review → Decision → Action?
- Does this introduce mutable canonical state?
- Does this weaken auditability?

---

## Phase 7: Self-Review — Architecture Verification

**After implementation,** review independently. Search for:

- RFC violations
- ADR violations
- Policy violations
- architecture drift
- duplicated ownership
- dependency violations
- speculative abstractions
- unnecessary interfaces
- unnecessary DTOs
- unnecessary services
- hidden business logic

**Immediately repair detected issues.**

---

## Phase 8: Verification — Run All Tests and Validation

**Execute every applicable verification command.** Examples:

```bash
go test ./...
golangci-lint run
task test
task build
task verify:rfc
task graph:rebuild
task smoke:*
task report
```

**Rules:**

- Only skip commands that genuinely do not exist
- Never claim commands passed unless executed
- Build/graphify-out should never be committed

---

## Phase 9: Repair Loop — Fix Issues Automatically

If verification fails:

1. Analyze the failure
2. Repair the issue
3. Verify again
4. Repeat

**Maximum repair iterations: 5.** If issue cannot be fixed in 5 iterations, surface the blocker and document.

---

## Phase 10: Final Review — Verify Architecture Health

Verify that after your changes:

- Architecture remains simpler (fewer files, fewer abstractions, fewer layers)
- Ownership remains correct (single owner per responsibility)
- Boundaries remain intact (no cross-layer violations)
- No duplication introduced
- No unnecessary abstractions added
- Code follows project conventions

---

## Phase 11: Final Report — Document Complete Implementation

**Always finish with a report containing:**

- **Summary** — what was implemented, what problem it solves
- **Files Changed** — list of all modified/added files
- **Architecture Decisions** — key choices and why
- **RFC Compliance** — which RFCs were implemented, were they satisfied
- **ADR Compliance** — which ADRs were followed
- **Policy Compliance** — which policies were applied
- **Reference Implementations Used** — what existing code was extended/reused
- **Commands Executed** — all build and verification commands run
- **Test Results** — pass/fail status of all tests
- **Validation Results** — lint, style, coverage results
- **Remaining Technical Debt** — incomplete work, workarounds, future improvements needed
- **Risks** — what could break, mitigation strategies
- **Next Recommended Slice** — what to build next given this foundation

**Never finish with only "Done".**

---

# ⚙️ ARCHITECTURAL INVARIANTS — NON-NEGOTIABLE RULES

These rules are non-negotiable. Any code violating one is wrong by definition. Repair immediately.

- **Events are immutable.** Never mutate or delete emitted event; correct via new events.
- **Decisions are explicit and auditable.** Every Decision records who/what, why, inputs.
- **Reviews never commit Decisions.** Review produces findings; cannot enact Decision.
- **Agents never mutate canonical state directly.** They propose; system commits.
- **Agents never call LLM providers directly.** All model access goes through LLM router.
- **Prompt versions are immutable after release.** Released prompts frozen; ship new version.
- **Memory is policy-bound.** Reads/writes to memory honor governing policy; no ad-hoc access.
- **Spaces are bounded contexts.** Keep models, data, logic within their space.
- **Cross-space communication is explicit.** Use defined contracts/events; no hidden coupling.
- **Derived stores are rebuildable.** Never treat projection/cache as source of truth.

---

# 🏛️ ENGINEERING LAWS

## Reference First

If a Reference Implementation exists, follow it. Consistency is preferred over novelty.

## Single Ownership Rule

Every responsibility has exactly one owner. Never duplicate ownership.

## Boundary Rule

- Kernel never knows transport
- Kernel never knows storage implementation
- Storage never owns business logic
- Transport never owns business logic
- Composition roots own wiring only

## Simplicity Rule

Prefer:

- Fewer files
- Fewer abstractions
- Fewer interfaces
- Fewer services
- Fewer layers
- Less indirection

**Architecture should become simpler after every implementation.**

## Minimize Architectural Debt

- Prefer deleting code over adding abstractions
- Prefer composition over inheritance
- Prefer explicit behavior over magic
- Prefer deterministic behavior over convenience

## Challenge Assumptions

Do not blindly implement requested design. If simpler solution satisfies RFCs, explain and recommend it before implementing. If request contradicts long-term architecture, explain why instead of complying silently.

## Leave the Repository Healthier

Every change should improve at least one of: documentation, verification, naming, comments, tests, architecture consistency, or dead-code removal.

---

# 🔍 CONTEXT GATHERING DEFAULTS

### Use Graphify First

For any question about repository architecture, ownership, dependencies or implementation locations:

- `graphify query` before grep
- `graphify path` for dependency questions
- `graphify explain` for architectural concepts

Only inspect source files when implementation is required, debugging, or Graphify lacks detail.

### Search Before Creating

Never assume. Search for existing implementations in these categories:

- mapper
- repository
- worker
- adapter
- storage backend
- transport
- bootstrap
- configuration loader
- service
- reducer
- projection
- event
- command
- query
- aggregate

### RFC Is Truth

RFCs in `./rfcs` are **architectural source of truth**.

- Never implement behavior that contradicts an accepted RFC
- If code and RFC disagree, RFC wins
- Stop and flag conflicts immediately

---

# 📋 WHEN UNSURE

- Prefer reading RFC over guessing
- If RFC ambiguous or missing, surface gap rather than filling with assumed behavior
- Do not implement when blocked by ambiguity
- Ask for clarification instead of inventing

---

# 🛠️ BUILD AND VERIFICATION OPERATIONS

**Rules:**

- Prefer Taskfile commands over raw go/python commands
- Use `task test` for normal validation
- Use `task build` before claiming binaries compile
- Use `task verify:rfc` after changing RFCs or docs
- Use `task graph:rebuild` after changing graph-relevant docs or code
- Never commit `build/` or `graphify-out/`
- Do not create ad-hoc build scripts unless Taskfile cannot express operation
- No command is canonical unless in `Taskfile.yml`

---

# 🚫 AVOID

- Speculative implementation of future features
- Generic frameworks or unused abstractions
- Placeholder services
- Speculative config
- Ad-hoc build scripts
- Exporting internal implementation details
- "Future-proof" interfaces

Everything must be justified by existing RFC or accepted implementation slice.

---

# ✅ DEFAULT CHECKLIST FOR EVERY IMPLEMENTATION

Before finishing, verify:

- [ ] All relevant RFCs read and understood
- [ ] No new abstractions without 3-point justification (Phase 2.2)
- [ ] If new abstraction: 11-phase architecture review completed (Phase 3)
- [ ] Implementation plan produced and followed
- [ ] Minimal vertical slice implemented
- [ ] RFC compliance verified (Phase 6)
- [ ] Self-review completed (Phase 7)
- [ ] All tests and verification commands executed (Phase 8)
- [ ] Issues fixed in repair loop (Phase 9)
- [ ] Final review passed (Phase 10)
- [ ] Final report produced (Phase 11)
- [ ] No RFC violations remain
- [ ] No ADR violations remain
- [ ] Architecture remains simpler
- [ ] Ownership remains correct
- [ ] Boundaries remain intact

---

# 📖 HOW TO USE THIS WORKFLOW

1. **Get a Sprint Task** — Developer provides a feature to implement
2. **Run the Nine Phases** — Execute Phase 1 through Phase 11 automatically
3. **Stop only if genuinely blocked** — missing RFC clarity, architectural ambiguity that makes multiple implementations equally valid, or implementation would violate established architecture
4. **Produce Final Report** — Never stop at implementation; always produce complete report
5. **Make it the default** — This workflow IS the default Copilot behavior for implementation tasks in Praxis code

**Developer expectation:** Provide task. Copilot executes complete workflow autonomously from discovery through final report.

