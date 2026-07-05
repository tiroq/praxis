# Quality Gates

Status: ACTIVE

---

# Purpose

This document defines the mandatory Quality Gates for every engineering task in Praxis.

A task is not complete because code compiles.

A task is complete only after every required Quality Gate passes.

Quality Gates provide objective completion criteria.

---

# Principles

Quality Gates must be:

- objective
- deterministic
- repeatable
- automatable
- architecture driven

Quality Gates never rely on personal judgement alone.

---

# Engineering Lifecycle

Every task follows:

Task
↓

Architecture

↓

Implementation

↓

Verification

↓

Review

↓

Learning

Every phase has one or more gates.

---

# Gate 0 — Requirements

Purpose

Ensure the task is understood before implementation.

Must Verify

- goal understood
- constraints identified
- assumptions documented
- unknowns listed
- affected RFCs identified

Failure

STOP

Implementation must not begin.

---

# Gate 1 — Existing Knowledge

Purpose

Prevent reinventing existing solutions.

Must Verify

- RFC search completed
- ADR search completed
- Policy search completed
- Reference Implementation search completed
- existing packages reviewed

Failure

STOP

Search again.

---

# Gate 2 — Architecture

Purpose

Validate architecture before coding.

Must Verify

- ownership defined
- dependencies valid
- layer boundaries preserved
- composition root identified
- invariants preserved

Failure

STOP

Architecture review required.

---

# Gate 3 — Reference Comparison

Purpose

Ensure architectural consistency.

Must Verify

- reference implementation used
or
- deviation justified

Failure

STOP

Explain deviation or simplify.

---

# Gate 4 — Complexity Budget

Purpose

Prevent uncontrolled growth.

Must Estimate

- files
- packages
- interfaces
- DTOs
- repositories
- adapters
- services
- public APIs

Must Verify

Every abstraction justified.

Failure

STOP

Reduce complexity.

---

# Gate 5 — Implementation

Purpose

Ensure implementation matches architecture.

Must Verify

- approved plan followed
- no unrelated refactoring
- minimal slice delivered
- responsibilities preserved

Failure

STOP

Implementation rejected.

---

# Gate 6 — Build

Purpose

Ensure project compiles.

Examples

```
task build
```

Must Verify

- build passes

Failure

STOP

---

# Gate 7 — Tests

Purpose

Ensure functionality.

Examples

```
go test ./...
```

Must Verify

- unit tests
- integration tests
- smoke tests when applicable

Failure

STOP

---

# Gate 8 — Static Analysis

Purpose

Ensure code quality.

Examples

```
golangci-lint run

go vet

gofmt
```

Must Verify

- lint
- formatting
- vet

Failure

STOP

---

# Gate 9 — RFC Compliance

Purpose

Ensure architecture follows specifications.

Must Verify

- RFC compliance
- ADR compliance
- Policy compliance

Failure

STOP

Architecture review required.

---

# Gate 10 — Architecture Guardian

Purpose

Protect architectural integrity.

Must Verify

- ownership
- dependency direction
- composition roots
- Extract Don't Invent
- Reference Implementations
- Golden Examples

Failure

STOP

Implementation rejected.

---

# Gate 11 — Self Review

Purpose

Independent engineering review.

Must Verify

- duplication
- coupling
- complexity
- speculative abstractions
- dead code
- missing tests
- missing documentation

Failure

Return to implementation.

---

# Gate 12 — Documentation

Purpose

Engineering knowledge must remain current.

Must Verify

Documentation updated if needed.

Examples

- RFC
- ADR
- Policy
- README
- Reference Implementation
- Golden Example

Failure

Return to implementation.

---

# Gate 13 — Engineering Learning

Purpose

Improve the engineering system.

Must Evaluate

Should we update:

- Instructions
- Workflow
- Prompt
- Reference Implementation
- Golden Example
- Policy
- ADR
- RFC

Engineering improvement is mandatory.

Failure

Task cannot be considered fully complete.

---

# Auto Fix Loop

Whenever any gate fails:

Implementation

↓

Verification

↓

Guardian

↓

Self Review

↓

Repeat

until all gates pass.

---

# Gate Dependencies

```
Requirements
    ↓
Knowledge
    ↓
Architecture
    ↓
Reference
    ↓
Complexity
    ↓
Implementation
    ↓
Build
    ↓
Tests
    ↓
Static Analysis
    ↓
RFC Compliance
    ↓
Architecture Guardian
    ↓
Self Review
    ↓
Documentation
    ↓
Engineering Learning
```

A gate cannot execute before its dependencies pass.

---

# Mandatory Stop Rules

Immediately stop when:

- RFC conflict detected
- ownership unclear
- abstraction unjustified
- architecture violation
- circular dependency
- failed tests
- failed verification
- complexity exceeds budget

Never continue through failed gates.

---

# Gate Ownership

| Gate | Owner |
|-------|-------|
| Requirements | Architect |
| Knowledge | Architect |
| Architecture | Architecture Guardian |
| Reference | Architecture Guardian |
| Complexity | Architect |
| Implementation | Implementer |
| Build | Verification |
| Tests | Verification |
| Static Analysis | Verification |
| RFC Compliance | Architecture Guardian |
| Self Review | Reviewer |
| Documentation | Documentation Agent |
| Engineering Learning | Chief Engineer |

---

# Success Criteria

A task is complete only when:

✅ Every Quality Gate passes

AND

✅ The engineering system itself has improved.

---

# Final Principle

Quality is not inspected into Praxis.

Quality is enforced through mandatory engineering gates.