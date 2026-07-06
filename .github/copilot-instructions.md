---
name: "Praxis Copilot Instructions"
description: "Complete autonomous implementation workflow for Praxis. Coordinates 11 phases from discovery through final report with architecture gates, RFC compliance, and verification rules."
applyTo: ["services/**", "packages/**", "apps/**", "scripts/**", "infra/**"]
---

# Praxis Copilot Instructions

**Single active entry point for all Copilot-driven implementation in Praxis.**

**Default Assumption:** The user expects complete implementation, not just analysis or planning.

---

# 🎯 THE 11-PHASE AUTONOMOUS WORKFLOW

**Phase execution order is strict. Do not skip phases. Do not proceed if stop conditions are reached.**

## Phase 1: Discover — RFCs, Architecture, and Existing Code

**Do this first.** Never start coding without understanding what already exists.

### 1.1 RFC Discovery

- List all RFCs that might apply to this task
- Read each RFC completely
- Quote specific sections governing this work
- Identify RFC ambiguities, gaps, or contradictions
- **STOP if RFCs conflict.** Surface the conflict and wait for clarification

### 1.2 Architecture Context

- Read relevant ADRs (in `./docs/adr/`)
- Read relevant Policies
- Identify architectural invariants (see Engineering Laws section)
- Identify project instructions

### 1.3 Graphify for Architecture Questions

Use Graphify before manual search: `graphify query`, `graphify path`, `graphify explain`

For architecture, ownership, dependencies: Graphify first.

### 1.4 Reference Implementation Search

Search for existing implementations:

- mapper, repository, worker, adapter, storage backend, transport, bootstrap, configuration loader, service, reducer, projection, event, command, query, aggregate

If reference implementation exists → **prefer evolving it**.

---

## Phase 2: Architectural Feasibility Gate — Is New Abstraction Needed?

**MANDATORY gate if creating:** new package, repository, service, engine, adapter, projection, public interface, composition root, interface, DTO, exported type, or manager.

**Skip to Phase 4 if only modifying existing components.**

### 2.1 Three-Point Validity Check

New abstraction justified ONLY IF at least one applies:

1. **Two or more independent implementations already exist** in the current codebase
2. **Duplicated behavior exists** and extraction reduces complexity
3. **Approved RFC explicitly defines** the abstraction by name and contract

### 2.2 Prohibited Justifications

- "We may need it later"
- Speculative extensibility
- Future integrations (Kafka, Slack, Email, Postgres, Redis, etc.)
- Theoretical flexibility without current concrete use

**If none apply** → do not create abstraction. Use concrete implementation.

---

## Phase 3: Complete Architecture Review (If Gate Applied)

**Only if Phase 2 = YES. Complete all 11 phases in order. Do not skip any.**

### 3.1 RFC Review

- List all relevant RFCs
- Quote specific sections governing this work
- Confirm no RFC prohibits proposed approach
- **STOP** if RFCs conflict or ambiguous

### 3.2 Existing Architecture Review

- List existing services, packages, repositories, engines, adapters
- List existing events, commands, queries, aggregates, projections
- List existing public interfaces and exported types
- **STOP** if existing component satisfies need (reuse instead)

### 3.3 Responsibility Analysis

For proposed component, explicitly state:

- What it owns (data, logic, behavior)
- What it does NOT own
- What it depends on
- What depends on it

**STOP** if responsibilities overlap with existing components.

### 3.4 Ownership Analysis

- Identify single owner of each responsibility
- Confirm no shared ownership
- Confirm no orphaned responsibilities

**STOP** if ownership is ambiguous or duplicated.

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

Validate:

- Infrastructure never depends on Domain
- Domain has zero infrastructure imports
- Application depends only on Domain interfaces
- Only Composition Root wires unrelated layers

**STOP** if violations exist.

### 3.6 Component Categorization

Assign component to exactly one category:

- **Repository** — storage interface, no business logic
- **Engine** — stateless logic processor, no storage
- **Reducer** — event → state transformation, pure function
- **Coordinator** — orchestrates across services, no domain logic
- **Adapter** — translates between layers, no domain logic
- **Composition Root** — wiring only, no behavior

Reject designs mixing categories.

Invalid: Repository containing business logic, Engine storing state, Reducer calling services, Coordinator implementing domain logic.

**STOP** if mixed.

### 3.7 Storage Rules Validation

If storage-related, confirm it does NOT own:

- ❌ Business logic
- ❌ Replay logic
- ❌ Orchestration
- ❌ Reducers
- ❌ Workflows
- ❌ Event sourcing logic

Storage owns: Persist data, Retrieve data, Transactions, Referential integrity.

**Projection-specific:** Projection Repository stores snapshots only. Replay → Replay Engine. Projection construction → Reducers. Checkpoints → Checkpoint Store.

**STOP** if violations.

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

**STOP** if violations.

### 3.9 Runtime Data Flow

Trace how data moves at runtime:

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
- No bypasses of boundaries
- State changes are auditable
- Events drive projections (not direct writes)

**STOP** if violations.

### 3.10 Public API Review

For every exported function, type, interface, method:

- Justify why public
- Confirm required by immediate production caller
- Confirm does not expose internals
- Confirm follows RFC naming

Reject: Speculative abstractions, generic frameworks, "future-proof" interfaces, exported internals.

**STOP** if speculative elements.

### 3.11 Abstraction Review — Extract, Don't Invent

For every new interface, DTO, repository, service, adapter, engine, manager, or public API:

| Question | Answer |
|---|---|
| Why does this abstraction exist today? | _(required)_ |
| Which existing duplication does it remove? | _(required)_ |
| How many concrete implementations currently exist? | _(required)_ |
| Which approved RFC requires it? | _(required)_ |
| Can a concrete implementation be used instead? | _(required)_ |
| Would removing this abstraction simplify the architecture? | _(required)_ |

**Validity:** At least one must be true:

1. Two or more independent implementations already exist
2. Duplicated behavior exists and extraction reduces complexity
3. Approved RFC explicitly defines this abstraction

**STOP** if any abstraction cannot be justified. The abstraction MUST be removed.

---

## Phase 4: Implementation Planning

**Do not write code yet.**

Plan must include:

- Relevant RFCs
- Impacted invariants (see Engineering Laws)
- Existing components to reuse
- New components to create (each justified)
- Minimal implementation slice
- Verification strategy
- Risks

**Cannot map to RFC?** Stop. Do not invent behavior.

---

## Phase 5: Implementation — Build Minimal Vertical Slice

Key rules:

- Reuse existing code, packages, architecture
- Prefer deletion over addition
- Prefer concrete implementations
- Avoid speculative abstractions
- Smallest end-to-end slice: `Command → Event → Projection → Query → Verification`

### 5.1 Pure Mapper Rule (If Applicable)

**Every transport mapping function MUST be pure with no side effects.**

#### Mapper Shape (Required)

```
Input:  external transport object (e.g., Telegram Update, HTTP request)
Output: Praxis wire contract dict (internal/transport/nats.InputMessage)
```

#### Mapper MUST NOT

- ❌ Publish messages
- ❌ Call HTTP endpoints
- ❌ Call NATS or storage
- ❌ Access environment variables
- ❌ Perform retries
- ❌ Mutate global state
- ❌ Generate business identifiers
- ❌ Execute business logic
- ❌ Call repositories, databases, HTTP services, LLMs

#### Allowed Complexity Only

- ✅ Rename fields
- ✅ Copy fields
- ✅ Drop unused fields
- ✅ Normalize formatting
- ✅ Convert primitive types
- ✅ Construct transport DTOs
- ✅ Perform deterministic serialization
- ✅ Validate transport-level preconditions
- ✅ Create deterministic identifiers derived from input only

#### Referential Transparency Invariant

Given the same input, MUST always produce the same output. No side effects permitted.

#### Verification Checklist

For every mapper:

| Question | Answer |
|---|---|
| Does mapper accept only the external transport object? | _(must be YES)_ |
| Does mapper return only a dict (wire contract)? | _(must be YES)_ |
| Does mapper call any I/O function? | _(must be NO)_ |
| Does mapper access any global state? | _(must be NO)_ |
| Does mapper perform any retry logic? | _(must be NO)_ |
| Can mapper be called 1,000,000 times with same input and produce identical output? | _(must be YES)_ |
| Is mapper completely free of business logic? | _(must be YES)_ |

**STOP** if any fails purity. Refactor before proceeding.

#### Structural Triviality Test

- Can every output field be traced directly to input fields?
- Is every transformation deterministic?
- Would another engineer understand it in under one minute?
- Could it be rewritten as simple field mappings?
- If removed, would business behavior remain unchanged?

If any "No", stop and perform architecture review.

Reference: `docs/architecture/GOLDEN_MAPPER.md`

---

## Phase 6: Continuous Validation — RFC and Architecture Compliance

Before committing code, validate:

- Does this violate any RFC? → **STOP if YES**
- Does this duplicate existing concept? → **STOP if YES**
- Does this introduce second source of truth? → **STOP if YES**
- Does this bypass Review → Decision → Action? → **STOP if YES**
- Does this introduce mutable canonical state? → **STOP if YES**
- Does this weaken auditability? → **STOP if YES**

---

## Phase 7: Self-Review — Architecture Verification

After implementation, review independently for:

- RFC violations
- ADR violations
- Policy violations
- Architecture drift
- Duplicated ownership
- Dependency violations
- Speculative abstractions
- Unnecessary interfaces, DTOs, services
- Hidden business logic

**Immediately repair detected issues.**

---

## Phase 8: Verification — Run All Tests and Validation

Execute every applicable command:

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

Rules:

- Only skip commands that genuinely do not exist
- Never claim commands passed unless executed
- Build/graphify-out should never be committed

---

## Phase 9: Repair Loop — Fix Issues Automatically

If verification fails:

1. Analyze failure
2. Repair issue
3. Verify again
4. Repeat

**Maximum iterations: 5.**

If cannot fix in 5 iterations, surface blocker and document.

---

## Phase 10: Final Review — Verify Architecture Health

Verify:

- Architecture remains simpler (fewer files, fewer abstractions, fewer layers)
- Ownership remains correct (single owner per responsibility)
- Boundaries remain intact (no cross-layer violations)
- No duplication introduced
- No unnecessary abstractions added
- Code follows project conventions

---

## Phase 11: Final Report — Document Complete Implementation

Always finish with:

- **Summary** — what was implemented, what problem it solves
- **Files Changed** — all modified/added files
- **Architecture Decisions** — key choices and why
- **RFC Compliance** — which RFCs implemented, were they satisfied
- **ADR Compliance** — which ADRs followed
- **Policy Compliance** — which policies applied
- **Reference Implementations Used** — existing code extended/reused
- **Commands Executed** — all build and verification commands run
- **Test Results** — pass/fail status
- **Validation Results** — lint, style, coverage results
- **Remaining Technical Debt** — incomplete work, workarounds, future improvements
- **Risks** — what could break, mitigation strategies
- **Next Recommended Slice** — what to build next

**Never finish with only "Done".**

---

# 🚫 STOP CONDITIONS — These Always Stop Workflow

- **RFCs conflict or ambiguous** → Stop and wait for clarification
- **Multiple equally valid implementations possible** → Stop and ask
- **Would violate established architecture** → Stop immediately
- **Genuinely blocked by missing information** → Stop and explain
- **New abstraction fails Extract, Don't Invent test** → Stop (do not create)
- **Any phase of architecture review fails** → Stop (do not proceed)
- **RFC compliance fails** → Stop and fix
- **Repair loop exceeds 5 iterations** → Stop and document

---

# 🏛️ ARCHITECTURAL INVARIANTS

These are non-negotiable. Any code violating one is wrong by definition. Repair immediately.

- **Events are immutable.** Never mutate or delete; correct via new events.
- **Decisions are explicit and auditable.** Every Decision records who/what, why, inputs.
- **Reviews never commit Decisions.** Reviews produce findings; cannot enact Decision.
- **Agents never mutate canonical state directly.** They propose; system commits.
- **Agents never call LLM providers directly.** All model access through LLM router.
- **Prompt versions are immutable after release.** Released prompts frozen; ship new version.
- **Memory is policy-bound.** Reads/writes honor governing policy; no ad-hoc access.
- **Spaces are bounded contexts.** Keep models, data, logic within their space.
- **Cross-space communication is explicit.** Use defined contracts/events; no hidden coupling.
- **Derived stores are rebuildable.** Never treat projection/cache as source of truth.

---

# 💡 CORE PRINCIPLES

### RFC Is Source of Truth

RFCs in `./rfcs/` are **immutable architectural contracts**.

- Never reinterpret them
- Never extend with assumptions
- Never fill gaps with "reasonable" guesses
- Never implement behavior not explicitly specified

If RFC is ambiguous, incomplete, or contradictory:

- **STOP immediately**
- Document specific ambiguity
- Explain what cannot be determined
- Wait for clarification

Do not proceed when architecture is unclear.

### Extract, Don't Invent

Abstractions **MUST** be extracted from existing code or approved RFCs. Never invented from anticipated future needs.

**Validity Criteria:** At least one must be true:

1. Two or more independent implementations already exist in codebase
2. Duplicated behavior exists and extraction reduces complexity
3. Approved RFC explicitly defines the abstraction

### Reference First

If a Reference Implementation exists, follow it. Consistency is preferred over novelty.

### Single Ownership Rule

Every responsibility has exactly one owner. Never duplicate ownership.

### Boundary Rule

- Kernel never knows transport
- Kernel never knows storage implementation
- Storage never owns business logic
- Transport never owns business logic
- Composition roots own wiring only

### Simplicity Rule

Prefer:

- Fewer files
- Fewer abstractions
- Fewer interfaces
- Fewer services
- Fewer layers
- Less indirection

**Architecture should become simpler after every implementation.**

### Minimize Architectural Debt

- Prefer deleting code over adding abstractions
- Prefer composition over inheritance
- Prefer explicit behavior over magic
- Prefer deterministic behavior over convenience

### Challenge Assumptions

Do not blindly implement requested design. If simpler solution satisfies RFCs, explain and recommend it. If request contradicts long-term architecture, explain why.

### Leave the Repository Healthier

Every change should improve: documentation, verification, naming, comments, tests, architecture consistency, or dead-code removal.

---

# ❌ AVOID

- Speculative implementation of future features
- Generic frameworks or unused abstractions
- Placeholder services
- Speculative config
- Ad-hoc build scripts
- Exporting internal implementation details
- "Future-proof" interfaces

Everything must be justified by existing RFC or accepted implementation slice.

---

# 🔍 CONTEXT GATHERING DEFAULTS

### Use Graphify First

For repository architecture, ownership, dependencies, implementation locations:

- `graphify query` before grep
- `graphify path` for dependency questions
- `graphify explain` for architectural concepts

Only inspect source files when implementation required, debugging, or Graphify lacks detail.

### When Unsure

- Prefer reading RFC over guessing
- If RFC ambiguous or missing, surface gap rather than filling with assumed behavior
- Do not implement when blocked by ambiguity
- Ask for clarification instead of inventing

---

# 📋 BUILD AND VERIFICATION RULES

- **Prefer Taskfile commands** over raw go/python commands
- **Use `task test`** for normal validation
- **Use `task build`** before claiming binaries compile
- **Use `task verify:rfc`** after changing RFCs or docs under `rfcs/`
- **Use `task graph:rebuild`** after changing graph-relevant docs or code
- **Never commit** `build/` or `graphify-out/`
- **Do not create ad-hoc build scripts** unless Taskfile cannot express operation

---

# ✅ DEFAULT CHECKLIST FOR EVERY IMPLEMENTATION

Before finishing, verify:

- [ ] All relevant RFCs read and understood
- [ ] Phase 1: Discovery complete
- [ ] If new abstraction: Phase 2 feasibility gate passed
- [ ] If new abstraction: Phase 3 (11-phase architecture review) complete
- [ ] Phase 4: Implementation plan produced and followed
- [ ] Phase 5: Minimal vertical slice implemented
- [ ] Phase 6: RFC compliance verified
- [ ] Phase 7: Self-review completed
- [ ] Phase 8: All tests and verification commands executed
- [ ] Phase 9: Issues fixed in repair loop
- [ ] Phase 10: Final review passed
- [ ] Phase 11: Final report produced
- [ ] No RFC violations remain
- [ ] No ADR violations remain
- [ ] No invariant violations remain
- [ ] Architecture remains simpler
- [ ] Ownership remains correct
- [ ] Boundaries remain intact

---

# 📖 HOW TO USE THIS WORKFLOW

1. **Get a Sprint Task** — Developer provides feature to implement
2. **Run the 11 Phases** — Execute Phase 1 through Phase 11 automatically
3. **Stop only if genuinely blocked** — missing RFC clarity, architectural ambiguity, or implementation would violate architecture
4. **Produce Final Report** — Never stop at implementation; always produce complete report

**Developer expectation:** Provide task. Copilot executes complete workflow autonomously from discovery through final report.

---

# 🔗 REFERENCE DOCUMENTS

For detailed context on specific topics, see:

- `rfcs/` — Architectural source of truth
- `docs/adr/` — Architecture Decision Records
- `docs/architecture/GOLDEN_MAPPER.md` — Reference mapper implementation
- `docs/architecture.md` — Overall architecture
- `docs/domain-model.md` — Domain model reference
- `.github/instructions/` — Detailed reference documentation (not active instructions)