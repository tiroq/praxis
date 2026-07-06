# Development Process & Validation

This directory contains the engineering process, workflow, and verification discipline used to design, implement, and validate all Praxis changes.

## Core Process Documents

### [architecture-review.md](architecture-review.md)

**11-phase architecture review gate for new abstractions.**

**Use when:**
- Creating new packages, services, adapters, interfaces, DTOs
- Proposing new abstractions
- Need to validate architectural feasibility

**Phases:**
1. RFC Review
2. Existing Architecture Analysis
3. Responsibility Analysis
4. Ownership Analysis
5. Layer Validation
6. Component Categorization
7. Storage Rules
8. Dependency Graph
9. Runtime Data Flow
10. Public API Analysis
11. Abstraction Review (Extract/Don't Invent)

**Stop condition:** If any phase fails, abstraction must be redesigned or justified.

### [implementation-guide.md](implementation-guide.md)

**Planning and execution discipline for implementation tasks.**

- **Plan before coding** — document all decisions
- **Reuse before build** — search for existing solutions first
- **Vertical slices** — implement minimum complete feature
- **7-layer implementation order** — canonical sequence

**Use when:**
- Starting a new task
- Planning a feature
- Deciding what to build vs. reuse

### [validation.md](validation.md)

**Verification operations and repair discipline.**

- **Build verification** — tests, lint, build commands
- **Repair loop** — maximum 5 iterations before blocking
- **Final review checklist** — 6-point verification before completion

**Use when:**
- Verifying implementation
- Running tests
- Fixing failures
- Completing a task

### [QUALITY_GATES.md](QUALITY_GATES.md)

**Mandatory quality gates for every engineering task.**

Each gate provides objective, deterministic, repeatable completion criteria.

**Gates:**
- Gate 0 — Requirements (goal, constraints, assumptions, unknowns, RFCs)
- Gate 1 — Existing Knowledge (RFC, ADR, Policy, Reference Impl searches)
- Gate 2 — Architecture (ownership, dependencies, boundaries, invariants)
- Gate 3 — Reference Comparison (use ref impl or justify deviation)
- Gate 4 — Complexity Budget (scope and risk assessment)
- ...plus verification gates

**Use when:**
- Starting a task
- Verifying completion criteria
- Reviewing architecture

## Reference & Learning

### [workflow-roles.md](workflow-roles.md)

**Role-based perspective on the engineering workflow.**

Describes workflow from role perspective (Intake, Knowledge Search, Architect, Implementer, Verifier, Guardian, etc.).

**Use when:**
- Understanding role responsibilities
- Designing workflows
- Delegating tasks

### [learning.md](learning.md)

**Continuous learning and knowledge evolution practices.**

### [knowledge-evolution.md](knowledge-evolution.md)

**How engineering knowledge evolves in Praxis.**

### [operating-system.md](operating-system.md)

**Engineering operating system layers.**

The hierarchy of engineering authority:

```
RFC
  ↓
ADR
  ↓
Policy
  ↓
Pattern
  ↓
Reference Implementation
  ↓
Checklist
  ↓
Playbook
```

### [verification-checklists.md](verification-checklists.md)

**Checklist catalog for repeatable quality gates.**

Checklists transform architectural knowledge into objective verification steps.

## Workflow Lifecycle

Every task follows this sequence:

```
Requirements → Architecture → Implementation → Verification → Review → Learning
```

See: [workflow-roles.md](workflow-roles.md) for role-based details.

See: [QUALITY_GATES.md](QUALITY_GATES.md) for verification gates.

## Integration with Copilot

The 11-phase workflow in [../../.github/copilot-instructions.md](../../.github/copilot-instructions.md) implements this engineering process for autonomous AI execution.

The documents in this directory are the **engineering reference**; the Copilot file is the **orchestration layer**.
