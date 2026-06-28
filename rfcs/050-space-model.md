# RFC-050 Space Model

**Status:** Draft
**Authors:** Tiroq + ChatGPT
**Last Updated:** 2026-06-28

## 1. Summary

This RFC defines the **Space** as the primary bounded context of Praxis.

Every Canonical Object, Event, Artifact, Decision, Review, Agent, Prompt, Memory and Workflow exists inside a Space.

## 2. Goals

- Unify all domains under one architectural model.
- Isolate runtime concerns.
- Scope memory, prompts and agents.
- Support collaboration through explicit cross-space relationships.

## 3. Space Definition

A Space is an isolated operational environment with its own identity, policies, memory, agents and knowledge.

Examples:

- Personal
- Work
- Product
- Freelance
- Education
- Finance

## 4. Architecture

```
Space
├── Canonical Objects
├── Events
├── Artifacts
├── Reviews
├── Decisions
├── Agents
├── Prompts
├── Memory
├── Knowledge Graph
├── Policies
└── Workflows
```

## 5. Identity

| Field | Description |
|------|-------------|
| Space ID | Stable identifier |
| Name | Display name |
| Type | Personal, Work, Product... |
| Owner | User or Organization |
| Visibility | Private, Shared, Public |
| Status | Draft, Active, Archived |

## 6. Responsibilities

A Space owns:
- Policies
- Agents
- Prompt Registry
- Memory
- Knowledge Graph
- Reviews
- Decisions
- Workflows

## 7. Isolation

Spaces isolate:
- Memory
- Prompts
- Agents
- Reviews
- Decisions
- Permissions

Cross-space access must be explicit.

## 8. Lifecycle

Draft → Active → Suspended → Archived

## 9. Default Spaces

- Personal Space
- Work Space
- Product Space
- Freelance Space
- Education Space
- Finance Space

## 10. Invariants

- Every Canonical Object belongs to one primary Space.
- Every Agent belongs to a Space.
- Every Prompt belongs to a Space.
- Memory is space-scoped.
- Policies are space-scoped.
- Cross-space communication is explicit.
- Space boundaries are auditable.

## 11. Dependencies

Depends on RFC-000 through RFC-043.

Required before RFC-051 through RFC-056.

## 12. Acceptance Criteria

- Space identity defined.
- Isolation defined.
- Responsibilities defined.
- Cross-space communication defined.

## 13. Decision Log

| Date | Decision |
|------|----------|
| 2026-06-28 | Introduced Space as the primary bounded context. |

> Everything in Praxis exists within a Space.
