# Pattern Catalog

Status: ACTIVE

---

# Purpose

Pattern Catalog defines the approved implementation patterns used throughout Praxis.

A Pattern is not code.

A Pattern is a proven architectural solution to a recurring problem.

Reference Implementations demonstrate patterns.

Patterns explain why those implementations exist.

---

# Relationship

RFC
    ↓

ADR
    ↓

Pattern
    ↓

Reference Implementation
    ↓

Concrete Code

RFC defines rules.

ADR defines decisions.

Patterns define recurring solutions.

Reference Implementations demonstrate them.

Code applies them.

---

# Pattern Template

Every pattern contains:

Purpose

Problem

Forces

Solution

Responsibilities

Ownership

Boundaries

Reference Implementation

Common Mistakes

Anti-patterns

When NOT to use

Related RFCs

Related ADRs

---

# Pattern Index

Current approved patterns:

Transport Mapper

Repository

Adapter

Composition Root

Worker

Storage Backend

Configuration Loader

Event Publisher

Event Subscriber

Projection Builder

Reducer

Service Bootstrap

CLI Command

HTTP Handler

Retry Loop

Background Worker

Migration

---

# Pattern

## Transport Mapper

Purpose

Translate one external representation into one internal representation.

Problem

External systems speak different protocols.

Kernel must never understand them.

Forces

Multiple transports.

Transport-specific objects.

Need deterministic conversion.

Need transport isolation.

Solution

Create a single pure mapper.

Exactly one translation.

No business logic.

No infrastructure.

Responsibilities

Own:

field mapping

renaming

primitive conversion

identifier formatting

default transport values

Do NOT own:

classification

validation

planning

decision making

storage

publishing

logging

retry

metrics

LLM

Reference

apps/telegram/main.py

telegram_update_to_payload()

Related

architecture/GOLDEN_MAPPER.md

---

## Repository Pattern

Purpose

Hide persistence.

Problem

Kernel cannot know databases.

Solution

Repository interface.

Storage implementation.

Composition root wiring.

Reference

(TBD)

---

## Adapter Pattern

Purpose

Translate external systems into Praxis contracts.

Problem

External APIs evolve independently.

Solution

Adapters terminate protocols.

Kernel receives canonical contracts.

Reference

Telegram Adapter

---

## Composition Root

Purpose

Wire implementations.

Problem

Kernel must not construct dependencies.

Solution

Construct everything once.

Inject downward.

Reference

Worker bootstrap

(TBD)

---

## Worker Pattern

Purpose

Consume transport.

Execute business pipeline.

Publish outcome.

Reference

Worker service.

(TBD)

---

# Pattern Selection

Before implementation ask:

Which Pattern solves this problem?

Does one already exist?

Can Reference Implementation be reused?

Would another Pattern fit better?

If no Pattern exists:

Document one.

Review it.

Approve it.

Only then implement repeatedly.

---

# Pattern Evolution

Patterns evolve.

Reference Implementations evolve.

Code evolves.

RFCs rarely evolve.

Patterns may become obsolete.

Deprecated patterns remain documented.

---

# Pattern Promotion

A solution becomes a Pattern only if:

appears multiple times

architecturally reviewed

proven in production

RFC compliant

simpler than alternatives

---

# Anti-pattern

Never create:

Pattern for one implementation.

Pattern from prediction.

Pattern from hypothetical future.

Extract after repetition.

Never before.

---

# Review Questions

Does this implementation follow a Pattern?

Which Pattern?

Why?

Does Reference Implementation match?

Does deviation improve architecture?

Would another Pattern be simpler?

---

# Metrics

Track:

Approved Patterns

Deprecated Patterns

Reference Implementations

Implementations using each Pattern

Violations

Review findings

---

# Final Principle

Patterns capture architectural thinking.

Reference Implementations capture engineering practice.

Both exist so engineers solve recurring problems consistently instead of inventing new solutions.