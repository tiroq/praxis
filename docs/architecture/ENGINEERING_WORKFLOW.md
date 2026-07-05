# Engineering Workflow

Status: ACTIVE

---

# Purpose

This document defines the engineering workflow used to design, implement, review, verify, and continuously improve Praxis.

It is independent of:

- LLM vendor
- AI framework
- Agent framework
- IDE
- Runtime

The workflow is the engineering process.

Agents are only implementations of this workflow.

---

# Goals

The workflow exists to maximize:

- correctness
- consistency
- architectural integrity
- engineering velocity
- accumulated knowledge

The workflow minimizes:

- architectural drift
- duplicated solutions
- speculative abstractions
- review effort
- regressions

---

# Workflow Overview

Every task follows the same lifecycle.

```
Task
    │
    ▼
Intake
    │
    ▼
Knowledge Search
    │
    ▼
Architecture Review
    │
    ▼
Architecture Diff
    │
    ▼
RFC Impact
    │
    ▼
Reference Comparison
    │
    ▼
Complexity Budget
    │
    ▼
Implementation
    │
    ▼
Verification
    │
    ▼
Architecture Guardian
    │
    ▼
Self Review
    │
    ▼
Documentation
    │
    ▼
Engineering Learning
```

---

# Workflow Roles

The workflow defines responsibilities, not people.

One AI may execute multiple roles.

Multiple AIs may execute one role.

---

## Intake

Purpose

Normalize requirements.

Produces

- normalized task
- assumptions
- risks
- constraints

Never writes code.

---

## Knowledge Search

Purpose

Discover existing knowledge.

Searches

- RFC
- ADR
- Policy
- Reference Implementation
- Golden Example
- existing packages

Never proposes new architecture before search.

---

## Architect

Purpose

Design the implementation.

Produces

- architecture review
- ownership
- dependency graph
- architecture diff
- RFC impact
- complexity budget

Never writes production code.

---

## Implementer

Purpose

Implement approved architecture.

Responsibilities

- smallest vertical slice
- follow reference implementations
- preserve ownership

Never expands scope.

---

## Verifier

Purpose

Verify implementation.

Runs

- build
- tests
- lint
- smoke
- reports

Never fixes code.

Produces only facts.

---

## Architecture Guardian

Purpose

Protect architecture.

Checks

- RFCs
- ADRs
- Policies
- Reference Implementations
- ownership
- boundaries
- dependency direction

Rejects violations.

---

## Reviewer

Purpose

Independent engineering review.

Looks for

- duplication
- complexity
- coupling
- hidden behavior
- technical debt

May request refactoring.

---

## Documentation

Purpose

Update engineering knowledge.

Evaluates

- README
- ADR
- Policy
- Reference
- Golden Example

Produces documentation only.

---

## Chief Engineer

Purpose

Improve the engineering system.

Evaluates

Should this task update:

- instructions
- prompts
- workflow
- quality gates
- policies
- references
- ADR
- RFC

Engineering improvement is mandatory.

---

# Workflow Artifacts

Each stage produces artifacts.

| Stage | Artifact |
|--------|----------|
| Intake | Normalized Task |
| Knowledge | Search Report |
| Architect | Architecture Proposal |
| Architect | Architecture Diff |
| Architect | RFC Impact Matrix |
| Architect | Complexity Budget |
| Implementer | Code Changes |
| Verifier | Verification Report |
| Guardian | Architecture Report |
| Reviewer | Review Report |
| Documentation | Documentation Updates |
| Chief Engineer | Engineering Learning Report |

Artifacts are immutable once accepted.

---

# Stage Inputs / Outputs

## Intake

Input

Task

Output

Normalized Task

---

## Knowledge

Input

Normalized Task

Output

Knowledge Report

---

## Architecture

Input

Knowledge Report

Output

Architecture Proposal

---

## Implementation

Input

Approved Proposal

Output

Implementation

---

## Verification

Input

Implementation

Output

Verification Report

---

## Guardian

Input

Implementation

Output

Architecture Report

---

## Review

Input

Implementation

Output

Review Report

---

## Learning

Input

Everything

Output

Engineering Improvements

---

# Feedback Loops

The workflow is not linear.

Verification failure

↓

Implementation

Guardian failure

↓

Architecture

Review failure

↓

Implementation

Documentation failure

↓

Documentation

Learning never loops.

It always executes once.

---

# Stop Conditions

Immediately stop when

- RFC conflict
- ownership conflict
- architecture violation
- circular dependency
- failed tests
- failed verification
- unjustified abstraction
- complexity exceeds budget

---

# Completion Conditions

A task completes only when

- every Quality Gate passed
- documentation updated
- engineering learning completed

---

# Workflow Independence

This workflow does not depend on

- Copilot
- Cursor
- Claude
- ChatGPT
- Gemini
- OpenAI
- Anthropic

Any implementation capable of executing the workflow may be used.

The workflow is the asset.

The implementation is replaceable.

---

# Evolution

The workflow evolves through

Engineering Learning.

Changes require

- architecture review
- documented rationale
- updated reference documents

The workflow must become simpler over time.

Never more complicated.

---

# Final Principle

The workflow is not designed to produce code.

The workflow is designed to produce correct engineering decisions that naturally lead to correct code.