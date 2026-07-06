# Reference Implementations & Patterns

This directory contains approved reference implementations and architectural patterns that demonstrate correct engineering decisions.

**Use these as templates for consistency.** Prefer evolving a reference implementation over inventing new patterns.

## Contents

### [REFERENCE_IMPLEMENTATIONS.md](REFERENCE_IMPLEMENTATIONS.md)

**Defines what reference implementations are and the process for creating them.**

- **Principles** — concrete, reviewed, minimal, production quality, architecturally correct
- **Selection criteria** — when a component becomes a reference implementation
- **Review process** — how candidates are approved
- **Architecture law** — "Prefer evolving reference implementation over inventing new style"

### [REFERENCE_REGISTRY.md](REFERENCE_REGISTRY.md)

**Canonical catalog of all approved reference implementations by architectural role.**

Includes official reference implementations for:

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

### [PATTERN_CATALOG.md](PATTERN_CATALOG.md)

**Approved architectural patterns used throughout Praxis.**

Patterns exist to maximize consistency. Do not invent new patterns without RFC approval.

## Using Reference Implementations

1. **Before designing** — Check if a reference implementation exists
2. **During implementation** — Follow the reference exactly
3. **On deviation** — Document explicit justification
4. **After completion** — Consider if your implementation should become a reference

See: [../../development/architecture-review.md](../../development/architecture-review.md#phase-3-reference-comparison) for Reference Comparison gate.

## Adding New Reference Implementations

A component may become a Reference Implementation only if:

- Used successfully in production or accepted architecture
- Reviewed by Architecture Guardian
- Follows all applicable RFCs
- Follows all ADRs
- Has no known architectural violations
- Remains understandable without additional explanation
- Approved by maintainers

Process: RFC → Architecture Review → Implementation → Approval
