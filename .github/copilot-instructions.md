# RFC-first

RFCs in `./rfcs` are the architectural source of truth.

Never implement behavior that contradicts an accepted RFC.

Before changing code:

- identify all relevant RFCs
- identify affected ADRs and Policies
- identify affected architectural invariants
- identify existing reference implementations
- propose the smallest executable vertical slice
- verify implementation against RFCs before completion

When editing code under:

- services/
- packages/
- apps/
- scripts/
- infra/

follow:

instructions/praxis-implementation.instructions.md

and

instructions/praxis-architecture-guardian.instructions.md

---

# Graphify

For any question about repository architecture, ownership, dependencies or implementation locations:

FIRST use Graphify.

Prefer:

graphify query

before grep.

Use:

graphify path

for dependency questions.

Use:

graphify explain

for architectural concepts.

Only inspect source files when:

- implementation is required
- debugging
- Graphify lacks enough detail

Never manually explore the repository before Graphify unless Graphify is unavailable.

---

# Autonomous Engineering Workflow

Default assumption:

The user expects implementation.

Do not stop after analysis.

Do not stop after planning.

Execute the complete engineering workflow automatically.

## Phase 1 — Discover

Read:

- relevant RFCs
- relevant ADRs
- relevant Policies
- project instructions
- Architecture Guardian
- Reference Implementations
- Golden Mapper (if applicable)

Search before creating.

Never assume.

---

## Phase 2 — Architecture Review

Identify:

- ownership
- responsibilities
- dependency direction
- boundaries
- composition root
- architectural invariants

Reject speculative designs.

---

## Phase 3 — Reference Search

Search for existing implementations.

Categories include:

- mapper
- repository
- worker
- adapter
- storage backend
- transport
- bootstrap
- configuration loader
- service

If one exists:

prefer evolving it.

Do not introduce a second implementation style without justification.

---

## Phase 4 — Implementation

Implement the smallest executable vertical slice.

Rules:

- reuse existing code
- reuse existing packages
- reuse existing architecture
- prefer deletion over addition
- prefer concrete implementations
- avoid speculative abstractions

---

## Phase 5 — Self Review

Review implementation independently.

Search for:

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

Immediately repair detected issues.

---

## Phase 6 — Verification

Run every applicable verification.

Examples:

go test ./...

golangci-lint run

task dev

task smoke:*

task report

Only skip commands that genuinely do not exist.

Never claim commands passed unless executed.

---

## Phase 7 — Repair Loop

If verification fails:

analyse

repair

verify again

Repeat until:

- implementation succeeds
- or further repair is impossible

Maximum repair iterations: 5.

---

## Phase 8 — Final Review

Verify:

- architecture remains simpler
- ownership remains correct
- boundaries remain intact
- no duplication introduced
- no unnecessary abstractions

---

## Phase 9 — Final Report

Always finish with:

- Summary
- Files Changed
- Architecture Decisions
- RFC Compliance
- ADR Compliance
- Policy Compliance
- Reference Implementations Used
- Commands Executed
- Test Results
- Remaining Technical Debt
- Risks
- Next Recommended Slice

Never finish with only "Done".

---

# Engineering Laws

## Extract, Don't Invent

Abstractions are extracted.

Never predicted.

Interfaces require:

- multiple implementations
- existing duplication
- explicit RFC

Otherwise use concrete implementations.

---

## Reference First

If a Reference Implementation exists:

follow it.

Consistency is preferred over novelty.

---

## Pure Mapper Rule

Transport mappers must:

- perform exactly one translation
- remain pure
- contain no business logic
- contain no infrastructure
- contain no retries
- contain no orchestration

Golden Mapper is the architectural reference.

---

## Single Ownership Rule

Every responsibility has exactly one owner.

Never duplicate ownership.

---

## Boundary Rule

Kernel never knows transport.

Kernel never knows storage implementation.

Storage never owns business logic.

Transport never owns business logic.

Composition roots own wiring only.

---

## Simplicity Rule

Prefer:

- fewer files
- fewer abstractions
- fewer interfaces
- fewer services
- fewer layers
- less indirection

Architecture should become simpler after every implementation.

---

## Documentation Rule

Update documentation only when implementation changes architecture or behaviour.

Do not generate speculative documentation.

---

# Default Behaviour

When the user asks to implement a feature:

assume implementation is requested.

Do not ask unnecessary clarification questions.

Make reasonable architectural decisions using:

- RFCs
- ADRs
- Policies
- existing repository patterns
- reference implementations

Only stop when:

- RFCs conflict
- architectural ambiguity makes multiple incompatible implementations equally valid
- implementation would violate established architecture