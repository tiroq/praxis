# Praxis Engineering System

Status: ACTIVE

---

# Purpose

This document defines how Praxis itself is engineered.

It is not part of the runtime architecture.

It governs the engineering process used to design, implement, review, verify, and evolve Praxis.

The engineering system exists to continuously improve:

- product quality
- architectural consistency
- development speed
- engineering discipline

Every completed task must improve both:

1. Praxis
2. the engineering system that builds Praxis.

---

# Engineering Principles

The engineering system follows these principles.

## Search Before Create

Always search first.

Look for:

- RFC
- ADR
- Policy
- Reference Implementation
- Existing package
- Existing API

Do not invent when something already exists.

---

## Extract Don't Invent

Abstractions must emerge.

Never create:

- interface
- DTO
- repository
- service
- framework

unless justified by:

- existing duplication
- multiple implementations
- explicit RFC

---

## Consistency Over Novelty

Prefer existing patterns.

Deviation requires explicit justification.

Reference implementations are architectural assets.

---

## Small Vertical Slices

Every implementation should produce the smallest possible working slice.

Avoid horizontal implementation.

---

## Architecture First

No implementation begins before architecture review.

---

## Continuous Learning

Every completed task must improve the engineering system.

---

# Engineering Artifacts

The engineering system is composed of:

RFC
↓

ADR
↓

Policy
↓

Reference Implementation
↓

Golden Example
↓

Instructions
↓

Workflow
↓

Prompt

Each layer has a different purpose.

---

# Responsibilities

## RFC

Defines domain rules.

Immutable.

---

## ADR

Defines architectural decisions.

Immutable.

---

## Policy

Defines operational rules.

Living document.

---

## Reference Implementation

Defines proven implementation pattern.

May evolve.

---

## Golden Example

Minimal exemplary implementation.

Used for comparison.

---

## Instructions

Guide AI engineering behaviour.

---

## Workflow

Defines engineering stages.

---

## Prompt

Executes a workflow.

---

# Engineering Workflow

Every task follows the same lifecycle.

Task
↓

Architecture Review
↓

Reference Search
↓

Architecture Diff
↓

RFC Impact
↓

Complexity Budget
↓

Implementation
↓

Verification
↓

Architecture Guardian
↓

Self Review
↓

Auto Fix Loop
↓

Final Report
↓

Engineering Learning

No stages may be skipped.

---

# Quality Gates

Every task must pass:

- Requirements
- Architecture
- References
- Complexity
- Implementation
- Tests
- Verification
- Guardian
- Self Review
- Engineering Learning

Failure blocks completion.

---

# Reference Implementations

Reference implementations define architectural consistency.

Examples:

- Transport Mapper
- Repository
- Worker
- Adapter
- Bootstrap
- Configuration Loader
- Storage Backend

Reference implementations are preferred over creating new patterns.

---

# Golden Examples

Golden Examples are intentionally minimal.

They demonstrate:

- ownership
- responsibility
- simplicity

Example:

telegram_update_to_payload()

Golden examples should remain boring.

Interesting code usually belongs elsewhere.

---

# Complexity Budget

Every task must estimate:

- files
- packages
- interfaces
- services
- repositories
- adapters
- public APIs

Complexity growth must always be justified.

---

# Architecture Review

Architecture review verifies:

- RFC compliance
- ADR compliance
- Policies
- Layer boundaries
- Ownership
- Dependency direction
- Composition roots
- Reference implementations

---

# Verification

Verification includes:

- build
- tests
- lint
- smoke
- reports

Failures are never ignored.

---

# Engineering Learning

Every completed task evaluates:

Should Instructions change?

Should Policies change?

Should ADRs change?

Should RFCs change?

Should a Reference Implementation be created?

Should a Golden Example be created?

Should the Workflow improve?

Should Prompts improve?

Engineering learning is mandatory.

---

# Engineering Debt

Technical debt is classified:

NOW

Must fix immediately.

LATER

Planned improvement.

NEVER

Rejected ideas.

Engineering debt must be visible.

---

# Success Criteria

The engineering system succeeds when:

- architecture remains consistent
- implementation speed increases
- review effort decreases
- duplication decreases
- onboarding becomes easier
- AI agents produce increasingly consistent code
- engineering knowledge accumulates instead of disappearing

---

# Final Principle

Praxis is not only an AI system.

Praxis is also the engineering system that continuously improves the way Praxis itself is built.