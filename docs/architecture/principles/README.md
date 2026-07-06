# Architectural Principles

This directory contains the core architectural principles governing all Praxis design and implementation.

These are **non-negotiable rules**. Violation is wrong by definition and must be repaired immediately.

## Contents

### [engineering-laws.md](engineering-laws.md)

**Core reference for all engineering decisions.**

- **10 Architectural Invariants** — immutable rules that code cannot violate
- **5 Core Principles** — Extract/Don't Invent, RFC Is Source of Truth, Reference First, Single Ownership, Boundary Rule
- **4 Engineering Laws** — Simplicity, Minimize Debt, Challenge Assumptions, Leave Healthier
- **Default assumptions** for implementation

**Use when:**
- Designing new abstractions
- Reviewing architecture
- Checking implementation against rules
- Understanding why something is forbidden

### [GOLDEN_MAPPER.md](GOLDEN_MAPPER.md)

**Pure function discipline for all transport mapping.**

- **Mapper shape** — canonical structure for transport adapters
- **23 Forbidden behaviors** — what mappers cannot do
- **9 Allowed behaviors** — what mappers can do
- **7-point verification checklist** — validate mapper purity
- **Referential transparency invariant** — immutable, deterministic, no side effects
- **Structural triviality test** — validate simplicity

**Use when:**
- Implementing transport adapters
- Reviewing mapper code
- Defining new transport contracts
- Understanding pure function requirements

## Architecture Domains

These principles apply across all architectural domains:

- **Kernel** — Domain logic, commands, events, projections, queries
- **Transport** — HTTP, gRPC, CLI adapters
- **Storage** — Event stores, projections, repositories
- **Integration** — External systems, agents, connectors

## Enforcement

Every architectural review must verify compliance with these principles.

See: [../development/architecture-review.md](../development/architecture-review.md) for the 11-phase review gate.
