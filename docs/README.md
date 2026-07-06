# Praxis Documentation

This is the canonical entry point for all long-lived repository documentation.

## Authority Hierarchy

When documents disagree, authority is resolved in this exact order:

1. Engineering Constitution
2. RFC
3. ADR
4. Policies
5. Principles
6. Reference
7. Development
8. Operations
9. Examples
10. Reports

## Canonical Sources by Authority

- Engineering Constitution: [MANIFESTO.md](../MANIFESTO.md)
- RFC: [rfcs/README.md](../rfcs/README.md)
- ADR: [docs/adr](adr)
- Policies: [docs/architecture/policies/README.md](architecture/policies/README.md)
- Principles: [docs/architecture/principles/engineering-laws.md](architecture/principles/engineering-laws.md), [docs/architecture/principles/GOLDEN_MAPPER.md](architecture/principles/GOLDEN_MAPPER.md)
- Reference: [docs/architecture/reference/README.md](architecture/reference/README.md)
- Development: [docs/development/architecture-review.md](development/architecture-review.md), [docs/development/implementation-guide.md](development/implementation-guide.md), [docs/development/validation.md](development/validation.md)
- Operations: [docs/operations](operations)
- Examples: service and app READMEs (for example [services/worker/README.md](../services/worker/README.md), [apps/telegram/README.md](../apps/telegram/README.md))
- Reports: temporary outputs only; not authoritative

## Taxonomy

Every document must belong to exactly one category from the hierarchy above.

- Constitution and RFC/ADR documents define what must be true.
- Policies, Principles, and Reference define how engineering work must remain consistent.
- Development and Operations define execution and runtime operation.
- Examples are explanatory and non-authoritative.
- Reports are temporary and must be removed after use.

## Single-Owner Concepts

- Engineering Laws: [docs/architecture/principles/engineering-laws.md](architecture/principles/engineering-laws.md)
- Mapper Rules: [docs/architecture/principles/GOLDEN_MAPPER.md](architecture/principles/GOLDEN_MAPPER.md)
- Architecture Review: [docs/development/architecture-review.md](development/architecture-review.md)
- Reference Registry: [docs/architecture/reference/REFERENCE_REGISTRY.md](architecture/reference/REFERENCE_REGISTRY.md)
- Quality Gates: [docs/development/QUALITY_GATES.md](development/QUALITY_GATES.md)
- Implementation Workflow: [docs/development/implementation-guide.md](development/implementation-guide.md)

## Lifecycle Rules

- Add new long-lived architecture requirements to RFC or ADR first.
- Add implementation constraints to Policies, Principles, or Development docs.
- Add style guidance and examples to Reference.
- Put runtime/runbook content in Operations.
- Keep report and migration artifacts out of long-lived documentation.

## What Must Never Be Added

- Proposal-only documents after decisions are finalized.
- Historical implementation status reports as permanent docs.
- Migration playbooks after migration completion.
- Duplicate concept definitions across multiple files.
- Tool-specific behavior described as product architecture.
