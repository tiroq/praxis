# Reference Registry

Status: ACTIVE

This file is the single source of truth for approved architectural patterns and reference implementations.

## Purpose

The registry prevents architectural drift by making the preferred implementation style explicit.

A reference implementation is not a framework or reusable base class. It is an approved example of ownership, boundaries, and implementation style.

## Principles

Reference implementations must be:

- concrete
- reviewed
- minimal
- production quality
- RFC compliant

Reference implementations are never:

- speculative templates
- inheritance hierarchies
- mandatory copy-paste

## Pattern Catalog

Approved architectural pattern categories:

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

## Registry

### Transport Mapper

- Status: ACTIVE
- Reference: `apps/telegram/main.py`
- Symbol: `telegram_update_to_payload()`
- Documentation: `docs/architecture/principles/GOLDEN_MAPPER.md`
- Purpose: Translate Telegram Update into Praxis InputMessage
- Properties: pure, deterministic, single-responsibility, no business logic, no infrastructure side effects

### Repository

- Status: TBD
- Reference: None

### Worker

- Status: TBD
- Reference: None

### Adapter

- Status: TBD
- Reference: None

### Configuration Loader

- Status: TBD
- Reference: None

### Event Publisher

- Status: TBD
- Reference: None

### Event Subscriber

- Status: TBD
- Reference: None

### Storage Backend

- Status: TBD
- Reference: None

### Projection Builder

- Status: TBD
- Reference: None

## Usage Rules

Before implementing any new component:

1. Search this registry for a matching category.
2. Reuse the documented style where possible.
3. Justify every deviation in architecture review.

## Promotion Rules

A candidate becomes a reference implementation only if it is:

- architecturally reviewed
- RFC and ADR compliant
- simpler than alternatives
- clear about ownership and boundaries

## Replacement Rules

A new reference replaces an existing one only if it is clearly simpler, clearer, and at least equally compliant.

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