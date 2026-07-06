# Praxis Documentation

This directory contains architecture and engineering documentation used to build and operate Praxis.

## Authority Hierarchy

Use this order when documents disagree:

1. `rfcs/` - domain and subsystem contracts (source of truth)
2. `docs/adr/` - architectural decisions and rationale
3. `docs/architecture/principles/` - non-negotiable engineering invariants and rules
4. `docs/architecture/reference/REFERENCE_REGISTRY.md` - approved patterns and reference implementations
5. `docs/development/` - implementation and verification guidance

## Documentation Ownership

- Architecture team owns RFC and ADR correctness.
- Engineering team owns development workflow and validation guidance.
- Service owners own runtime configuration, deployment, and security operating notes.

## Active Documentation Map

- `docs/architecture.md` - platform architecture overview
- `docs/adr/` - architecture decision records
- `docs/architecture/principles/engineering-laws.md` - core invariants and laws
- `docs/architecture/principles/GOLDEN_MAPPER.md` - transport mapper purity contract
- `docs/architecture/reference/REFERENCE_REGISTRY.md` - canonical pattern/reference catalog
- `docs/development/architecture-review.md` - abstraction gate and architecture review
- `docs/development/implementation-guide.md` - implementation planning and slice discipline
- `docs/development/validation.md` - verification commands and repair loop
- `docs/development/QUALITY_GATES.md` - mandatory quality gates
- `docs/configuration-reference.md` - runtime configuration catalog
- `docs/deployment.md` - deployment notes
- `docs/security.md` - security baseline
- `docs/review-cycles.md` - review-cycle concept summary

## Deliberate Scope

This directory does not contain temporary audits, migration-only playbooks, or one-off implementation reports. Those belong in pull requests, issue trackers, or generated build artifacts.
