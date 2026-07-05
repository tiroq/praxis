# Praxis Reference Implementations

Status: ACTIVE

---

# Purpose

This document defines the official Reference Implementations of Praxis.

A Reference Implementation is not reusable code.

It is an architecturally reviewed example demonstrating the preferred way to solve a specific engineering problem.

Reference Implementations exist to maximize architectural consistency.

---

# Principles

Reference Implementations are:

- concrete
- reviewed
- minimal
- production quality
- architecturally correct

They are not:

- frameworks
- libraries
- templates
- base classes

They are examples.

---

# Architecture Law

Prefer evolving an existing Reference Implementation over inventing a new implementation style.

Consistency is an architectural asset.

---

# Selection Criteria

A component may become a Reference Implementation only if:

- used successfully in production or accepted architecture
- reviewed by Architecture Guardian
- follows all applicable RFCs
- follows all ADRs
- has no known architectural violations
- remains understandable without additional explanation

---

# Review Process

Candidate

↓

Architecture Review

↓

RFC Compliance

↓

Reference Comparison

↓

Approval

↓

Registered

Only approved implementations belong in this document.

---

# Categories

Reference Implementations are organized by engineering category.

---

# Transport Mapper

Purpose

Translate one transport representation into one Praxis contract.

Golden Example

apps/telegram/main.py

Function

telegram_update_to_payload()

Properties

- pure
- deterministic
- one translation
- one responsibility
- no infrastructure
- no business logic

Reference Documentation

docs/architecture/GOLDEN_MAPPER.md

---

# Repository

Status

No official reference yet.

Selection Criteria

- repository only
- no business logic
- storage abstraction only

---

# Worker

Status

Candidate

Expected Responsibilities

- compose components
- dependency injection
- lifecycle
- startup
- shutdown

Must Not

- business logic
- storage logic
- transport logic

---

# Adapter

Status

Candidate

Expected Responsibilities

- boundary translation
- protocol handling
- representation mapping

Must Not

- orchestration
- planning
- business rules

---

# Bootstrap

Status

Candidate

Expected Responsibilities

- configuration loading
- dependency graph
- startup

---

# Configuration Loader

Status

Candidate

Expected Responsibilities

- environment variables
- defaults
- validation

Must Not

- business decisions

---

# Storage Backend

Status

Candidate

Expected Responsibilities

- persistence
- transactions
- migrations

Must Not

- replay
- business logic
- orchestration

---

# Service Bootstrap

Status

Candidate

Expected Responsibilities

- initialize runtime
- connect dependencies
- health
- shutdown

---

# Review Checklist

When implementing a component:

1.

Does a Reference Implementation exist?

↓

If YES

Compare against it.

2.

What differs?

3.

Why?

4.

Is deviation required?

5.

Can the existing implementation simply be copied and adapted?

If yes,

prefer adaptation.

---

# Complexity Rule

A new implementation should be:

- simpler than
or

- equivalent to

the Reference Implementation.

More complex implementations require explicit architectural justification.

---

# Deviation Rules

Allowed:

- RFC requires it
- Existing implementation cannot satisfy requirements
- Existing implementation contains accepted technical debt

Not Allowed:

- personal preference
- different coding style
- speculative flexibility
- future-proofing
- "cleaner architecture"

---

# Becoming a Reference

After completing a task evaluate:

Does this implementation improve the current reference?

Could it replace the existing reference?

Should it become the new standard?

If yes,

open Architecture Review.

Never self-promote an implementation.

---

# Lifecycle

Candidate

↓

Reference

↓

Deprecated

↓

Archived

Reference Implementations may evolve.

Deprecated references remain documented until replaced.

---

# Success Criteria

Reference Implementations should reduce:

- architectural drift
- unnecessary abstractions
- implementation diversity
- review effort

while increasing:

- consistency
- readability
- onboarding speed
- AI implementation quality

---

# Final Principle

A Reference Implementation is not the only possible implementation.

It is the implementation the architecture currently recommends.