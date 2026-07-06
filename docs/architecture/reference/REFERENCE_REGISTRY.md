# Reference Registry

Status: ACTIVE

---

# Purpose

Reference Registry is the canonical catalog of all approved Reference Implementations.

It exists to prevent architectural drift.

A Reference Implementation is not reusable code.

It is an approved example of ownership, responsibility,
boundaries and implementation style.

Every engineer must search the registry before inventing a new implementation.

---

# Principles

Reference Implementations are:

- concrete
- reviewed
- production quality
- minimal
- RFC compliant

Reference Implementations are NOT:

- frameworks
- reusable libraries
- templates
- inheritance hierarchies
- mandatory copy-paste

They exist to demonstrate correct engineering decisions.

---

# Categories

Reference Implementations are organized by architectural role.

Current categories:

- Transport Mapper
- Repository
- Storage Backend
- Worker
- Adapter
- Service Bootstrap
- Configuration Loader
- Event Publisher
- Event Subscriber
- Projection Builder
- Reducer
- Event Store
- Projection Store
- CLI Command
- HTTP Handler
- Background Job
- Scheduler
- Migration
- Tests
- Benchmark

---

# Registry

---

## Transport Mapper

Status

ACTIVE

Reference

apps/telegram/main.py

Function

telegram_update_to_payload()

Documentation

architecture/GOLDEN_MAPPER.md

Purpose

Translate Telegram Update into Praxis InputMessage.

Properties

- pure
- deterministic
- one translation
- one input
- one output
- no business logic
- no infrastructure

---

## Repository

Status

TBD

Reference

None

---

## Worker

Status

TBD

Reference

None

---

## Adapter

Status

TBD

Reference

None

---

## Configuration Loader

Status

TBD

Reference

None

---

## Event Publisher

Status

TBD

Reference

None

---

## Event Subscriber

Status

TBD

Reference

None

---

## Storage Backend

Status

TBD

Reference

None

---

## Projection Builder

Status

TBD

Reference

None

---

# Promotion Rules

A candidate becomes a Reference Implementation only if:

- production quality
- architecturally reviewed
- RFC compliant
- simpler than alternatives
- reviewed multiple times
- demonstrates correct ownership
- demonstrates correct boundaries

Anything else is rejected.

---

# Replacement Rules

A Reference Implementation may replace an older reference only if it is:

- simpler
- clearer
- less coupled
- easier to review
- equally RFC compliant

Novelty alone never justifies replacement.

---

# Search Rule

Before implementing:

Search the Reference Registry.

If a matching category exists:

1. Study it.
2. Compare with current task.
3. Reuse implementation style.
4. Justify every deviation.

Failure to search is an architecture violation.

---

# Review Questions

Every implementation review asks:

Does this category already have a Reference Implementation?

Can it be adapted?

Why is this implementation different?

Is the difference required?

Would copying the reference be simpler?

---

# Candidate Promotion

Engineering Learning may nominate:

Reference Candidate

Reason

Advantages

Differences

Review Outcome

Accepted

Rejected

Needs More Evidence

---

# Anti-Patterns

Reference Implementations must never become:

utility libraries

frameworks

base classes

shared abstractions

generic engines

registries

factories

They remain examples.

Nothing more.

---

# Metrics

Track:

Number of categories

Number of approved references

Categories without references

Rejected candidates

Reference replacements

Average implementation consistency

Architecture review findings prevented

---

# Final Principle

A Reference Implementation captures engineering knowledge.

It exists so the next implementation starts from proven architecture instead of invention.