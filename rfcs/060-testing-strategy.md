

# RFC-060 Testing Strategy

**Status:** Draft  
**Authors:** Tiroq + ChatGPT  
**Last Updated:** 2026-06-28

---

## 1. Summary

This RFC defines the testing strategy for Praxis.

Praxis is not a conventional CRUD application. It combines Canonical Objects, Events, Reviews, Decisions, Agents, LLM Routing, Prompt Versioning, Memory, Knowledge Graphs, Spaces, external integrations, and human-in-the-loop workflows.

Testing must therefore verify more than functions and endpoints. It must verify architecture, invariants, policies, replayability, agent behavior, prompt behavior, memory behavior, and cross-space boundaries.

This RFC defines the quality model and testing layers required to make Praxis trustworthy, evolvable, and safe to automate.

---

## 2. Relationship to Previous RFCs

This RFC depends on:

- RFC-000 Vision
- RFC-001 Principles
- RFC-002 Terminology
- RFC-003 Concept Model
- RFC-010 Capability Map
- RFC-011 Domain Model
- RFC-012 Artifact Model
- RFC-013 Event Model
- RFC-014 Identity & Representation Model
- RFC-015 Object Lifecycle Model
- RFC-020 Review System
- RFC-021 Decision Model
- RFC-022 State Machine
- RFC-030 System Architecture
- RFC-031 Service Contracts
- RFC-032 Data Flow
- RFC-033 Storage Model
- RFC-040 Agent Architecture
- RFC-041 LLM Routing
- RFC-042 Prompt Versioning
- RFC-043 Memory & Knowledge Graph
- RFC-050 Space Model
- RFC-051 Personal Space
- RFC-052 Work Space
- RFC-053 Product Space
- RFC-054 Freelance Space
- RFC-055 Education Space
- RFC-056 Finance Space

This RFC is required before:

- RFC-061 Verification Scripts
- RFC-062 Benchmarking

RFC-061 will define executable verification scripts.

RFC-062 will define benchmarks and evaluation datasets.

This RFC defines the strategy that those later RFCs implement.

---

## 3. Goals

The goals of this RFC are to:

- Define the testing philosophy for Praxis.
- Define quality dimensions and test layers.
- Define what must be verified for core architecture, Spaces, Agents, Prompts, LLM Routing, Memory, and Decisions.
- Define how tests relate to RFC invariants.
- Define deterministic and non-deterministic test approaches.
- Define regression, acceptance, integration, contract, security, and resilience testing expectations.
- Establish the foundation for verification scripts and benchmarking.

---

## 4. Non-Goals

This RFC does not:

- Define exact test framework choices.
- Define all test cases.
- Define CI/CD implementation.
- Define benchmark datasets.
- Define production monitoring dashboards.
- Define compliance certification.
- Replace manual review.

Those concerns are handled by implementation documents, RFC-061, and RFC-062.

---

## 5. Testing Philosophy

Praxis must be tested as a system of decisions, boundaries, and evidence.

The most important failures are not only crashes.

More dangerous failures include:

- agent overreach;
- prompt drift;
- memory leakage;
- lost decisions;
- broken replayability;
- invalid lifecycle transitions;
- cross-space boundary leaks;
- untraceable LLM outputs;
- silent corruption of derived stores;
- external action without approval;
- false confidence in AI-generated output.

Testing must therefore verify both behavior and governance.

---

## 6. Quality Dimensions

Praxis quality is evaluated across multiple dimensions.

| Dimension | Meaning |
|----------|---------|
| Correctness | System produces expected valid behavior. |
| Traceability | Outputs can be traced to inputs, Events, Reviews, Decisions, Prompts, and Policies. |
| Replayability | Derived state can be rebuilt from immutable history. |
| Safety | Agents and actions do not bypass review, decision, or policy boundaries. |
| Privacy | Space, memory, and sensitive data boundaries are preserved. |
| Determinism | Deterministic parts behave predictably. |
| Robustness | System handles bad inputs, failures, and partial outages. |
| Observability | Failures and decisions are visible. |
| Evolvability | RFC changes, prompt changes, and model changes can be tested. |
| Human Control | Human override and approval remain available where required. |

---

## 7. Testing Pyramid

Praxis uses a layered testing model.

```text
Manual Review / Human Evaluation
↑
Benchmark and Evaluation Suites
↑
End-to-End Workflow Tests
↑
Integration Tests
↑
Contract Tests
↑
Policy and Invariant Tests
↑
Unit Tests
```

Unlike a classic testing pyramid, Praxis requires explicit invariant and policy testing because many failures occur at architectural boundaries.

---

## 8. Test Categories

| Category | Purpose |
|---------|---------|
| Unit Tests | Verify isolated functions and domain logic. |
| Contract Tests | Verify service Commands, Events, Queries, and schemas. |
| Invariant Tests | Verify RFC-defined rules always hold. |
| Policy Tests | Verify access, approval, routing, and cross-space rules. |
| Integration Tests | Verify service interactions and external adapters. |
| Workflow Tests | Verify complete user/system flows. |
| Replay Tests | Verify rebuild of derived state. |
| Agent Tests | Verify agent behavior and boundaries. |
| Prompt Tests | Verify prompt versions and output schemas. |
| Memory Tests | Verify retrieval, privacy, provenance, and consolidation. |
| Security Tests | Verify access control and sensitive data handling. |
| Resilience Tests | Verify behavior under failures and degraded dependencies. |
| Benchmark Tests | Measure quality, latency, cost, and consistency over datasets. |

---

## 9. RFC Invariant Testing

Every RFC that defines invariants must produce corresponding tests.

Examples:

| RFC | Example Invariant | Test Type |
|-----|-------------------|-----------|
| RFC-013 | Events are immutable. | Event Store invariant test. |
| RFC-014 | Canonical identity is stable. | Identity invariant test. |
| RFC-020 | Reviews do not commit Decisions. | Review/Decision boundary test. |
| RFC-021 | Decisions are immutable after commit. | Decision lifecycle test. |
| RFC-033 | Derived stores are rebuildable. | Replay test. |
| RFC-040 | Agents do not execute irreversible actions directly. | Agent policy test. |
| RFC-041 | Agents never call providers directly. | Routing boundary test. |
| RFC-042 | Prompt versions are immutable after release. | Prompt registry test. |
| RFC-043 | Memory retrieval is policy-bound. | Memory privacy test. |
| RFC-050 | Cross-Space communication is explicit. | Space boundary test. |

RFC invariants are not documentation only. They are test obligations.

---

## 10. Domain Model Tests

Domain Model tests verify Canonical Objects, lifecycle states, relationships, and transitions.

They must cover:

- object creation;
- object identity;
- lifecycle transitions;
- invalid transitions;
- revision creation;
- external identity mapping;
- canonical vs projection separation;
- relationship constraints;
- object ownership by Space.

---

## 11. Event Tests

Event tests verify that Events remain immutable, traceable, and replayable.

They must cover:

- event append;
- event schema validation;
- event versioning;
- correlation ID;
- causation ID;
- actor identity;
- timestamp presence;
- Space ID where applicable;
- event replay order;
- correction via new Event rather than mutation.

---

## 12. State Machine Tests

State Machine tests verify allowed object and workflow transitions.

They must cover:

- allowed transitions;
- forbidden transitions;
- idempotent transition handling;
- terminal states;
- rollback or compensation where allowed;
- audit events for transitions;
- policy-gated transitions.

---

## 13. Review and Decision Tests

Review and Decision tests verify governance flow.

They must cover:

- Review Request creation;
- Review Package generation;
- Review immutability after completion;
- Decision Request creation;
- Decision commit;
- Decision immutability;
- superseding Decisions;
- evidence references;
- approval policy;
- human override;
- no direct Review-to-Action bypass.

---

## 14. Agent Tests

Agent tests verify behavior of bounded agents.

They must cover:

- agent identity;
- agent configuration;
- allowed capabilities;
- tool permissions;
- context policy;
- output schema;
- failure handling;
- retry behavior;
- human approval requirement;
- no direct provider calls;
- no direct canonical store mutation;
- no irreversible action execution without Decision.

Agent tests must include both positive and negative scenarios.

---

## 15. LLM Routing Tests

LLM Routing tests verify provider-independent inference.

They must cover:

- capability-based routing;
- model registry lookup;
- provider adapter selection;
- local-first policy;
- privacy constraints;
- fallback behavior;
- retry behavior;
- structured output validation;
- routing decision records;
- provider health handling;
- cost and latency limits;
- no agent direct provider invocation.

---

## 16. Prompt Tests

Prompt tests verify Prompt Versioning behavior.

They must cover:

- prompt identity;
- prompt version immutability;
- variable declaration;
- variable resolution;
- context injection;
- output schema binding;
- prompt rendering determinism where possible;
- prompt lifecycle transitions;
- prompt rollback;
- golden examples;
- prompt regression tests;
- prompt injection defenses.

---

## 17. Memory and Knowledge Tests

Memory and Knowledge tests verify retrieval, provenance, consolidation, and privacy.

They must cover:

- memory scope;
- retrieval policy;
- privacy zones;
- provenance presence;
- evidence references;
- confidence metadata;
- temporal graph relationships;
- contradiction preservation;
- consolidation pipeline;
- embedding regeneration;
- graph traversal bounds;
- no memory leakage across Spaces.

---

## 18. Space Boundary Tests

Space tests verify RFC-050 and Space specialization boundaries.

They must cover:

- object primary Space ownership;
- Space-scoped Events;
- Space-scoped Agents;
- Space-scoped Prompts;
- Space-scoped Memory;
- Space-scoped Reviews;
- Space-scoped Decisions;
- Cross-Space References;
- Cross-Space Events;
- forbidden implicit sharing;
- Space hierarchy access rules;
- archived Space behavior.

---

## 19. Space Specialization Tests

Each Space specialization must have tests for its core invariants.

| Space | Test Focus |
|------|------------|
| Personal | Privacy, family context, personal projections. |
| Work | decisions, incidents, service ownership, QA, delivery. |
| Product | roadmap, experiments, feedback, outcomes, governance. |
| Freelance | proposals, contracts, invoices, client delivery. |
| Education | skill evidence, assessments, child privacy, feedback. |
| Finance | transactions, budgets, risk, tax evidence, privacy. |

---

## 20. Integration Tests

Integration tests verify adapters and external system boundaries.

They must cover:

- authentication failure;
- rate limiting;
- malformed external data;
- duplicate external records;
- external ID mapping;
- sync idempotency;
- partial sync failure;
- retry behavior;
- archive/delete mismatch;
- external side-effect audit.

External systems are never the source of Praxis identity unless a specific RFC permits it.

---

## 21. API and Service Contract Tests

Service Contract tests verify RFC-031.

They must cover:

- Command schema validation;
- Event schema validation;
- Query schema validation;
- error contract;
- idempotency key behavior;
- authentication and authorization;
- version compatibility;
- backward compatibility where required.

---

## 22. Replay and Rebuild Tests

Replay tests verify that derived state can be rebuilt.

They must cover rebuild of:

- projections;
- read models;
- search indexes;
- vector indexes;
- derived knowledge graph relationships;
- analytics views;
- caches.

Replay must never rewrite immutable history.

---

## 23. Data Migration Tests

Migration tests verify that storage and schema changes are safe.

They must cover:

- forward migration;
- rollback where supported;
- data preservation;
- event compatibility;
- projection rebuild after migration;
- prompt version compatibility;
- old object lifecycle compatibility;
- migration observability.

---

## 24. Security and Privacy Tests

Security and Privacy tests verify trust boundaries.

They must cover:

- authentication;
- authorization;
- Space-level access;
- object-level access;
- memory-level access;
- prompt access;
- tool access;
- integration token handling;
- sensitive field redaction;
- logging without leaking secrets;
- privacy zone enforcement;
- cross-space data leakage.

---

## 25. Resilience Tests

Resilience tests verify behavior under failures.

Failure scenarios:

- LLM provider outage;
- external integration outage;
- database partial outage;
- message bus outage;
- vector store unavailable;
- search index unavailable;
- slow provider response;
- invalid model output;
- duplicate events;
- retry storms;
- stale projection;
- corrupted cache.

System behavior must degrade safely.

---

## 26. Human-in-the-Loop Tests

Human-in-the-loop tests verify human control remains effective.

They must cover:

- approval request creation;
- approval denial;
- manual override;
- human review requirement;
- escalation on low confidence;
- audit of human decisions;
- stale approval expiration;
- re-review after material change.

---

## 27. Non-Deterministic Output Testing

LLM and agent outputs may be non-deterministic.

Testing must therefore use:

- schema validation;
- property-based expectations;
- rubric-based evaluation;
- golden examples with allowed variance;
- regression thresholds;
- multiple-run consistency checks;
- human review for critical workflows.

Exact string comparison is acceptable only for deterministic components.

---

## 28. Test Data Strategy

Test data must include:

- minimal valid fixtures;
- invalid fixtures;
- edge cases;
- privacy-sensitive fixtures;
- cross-space fixtures;
- replay fixtures;
- prompt fixtures;
- agent fixtures;
- integration fixtures;
- benchmark datasets.

Test data must not contain real secrets or uncontrolled personal data.

---

## 29. Golden Datasets

Golden datasets preserve expected behavior over time.

They should exist for:

- prompt outputs;
- agent reviews;
- classification;
- opportunity scoring;
- memory retrieval;
- product prioritization;
- finance categorization;
- education assessment;
- work incident analysis.

Golden datasets must be versioned.

---

## 30. Regression Strategy

Regression tests protect against unwanted behavior change.

Regression triggers:

- RFC changes;
- prompt changes;
- model routing policy changes;
- model provider changes;
- memory retrieval changes;
- schema migrations;
- agent configuration changes;
- integration adapter changes.

Regression results must be comparable across versions.

---

## 31. Acceptance Testing

Acceptance tests verify user-significant flows.

Examples:

- capture a task and route it to correct Space;
- review an artifact and create a Decision;
- generate proposal draft and require approval;
- classify finance transaction and update budget projection;
- create learning plan and track assessment evidence;
- reconstruct projection from Events;
- prevent unauthorized cross-space memory access.

---

## 32. Observability Testing

Observability itself must be tested.

Tests must verify:

- logs are emitted;
- metrics are emitted;
- traces preserve correlation IDs;
- routing decision records exist;
- prompt decision records exist;
- memory decision records exist;
- failures are visible;
- sensitive payloads are not leaked.

---

## 33. CI Strategy

A recommended CI structure:

```text
Fast Checks
↓
Unit Tests
↓
Contract Tests
↓
Invariant Tests
↓
Integration Tests
↓
Workflow Tests
↓
Evaluation / Benchmark Suite
```

Fast checks should block every change.

Longer evaluation suites may run on schedule or before release.

---

## 34. Test Ownership

Test ownership follows architectural ownership.

| Test Type | Owner |
|----------|-------|
| Unit Tests | Owning service or package. |
| Contract Tests | Service contract owner. |
| Invariant Tests | RFC owner or architecture owner. |
| Agent Tests | Agent owner. |
| Prompt Tests | Prompt owner. |
| Memory Tests | Memory/Knowledge owner. |
| Space Tests | Space specialization owner. |
| Integration Tests | Adapter owner. |
| Benchmarks | Evaluation owner. |

---

## 35. Testing Invariants

The following testing invariants must hold:

- Every RFC invariant should have at least one test or verification script.
- Tests must distinguish canonical truth from derived projections.
- Tests must verify cross-space boundaries.
- Tests must verify agent and tool permissions.
- Tests must verify prompt version traceability.
- Tests must verify memory privacy and provenance.
- Tests must verify Decisions are explicit and auditable.
- Tests must not rely on uncontrolled real personal data.
- Tests must make non-determinism explicit.
- Critical workflows must include negative tests.

---

## 36. Architectural Consequences

This testing strategy enables:

- architecture-level regression protection;
- safer AI-assisted automation;
- reliable prompt and agent evolution;
- stronger privacy guarantees;
- confidence in replay and rebuild;
- explicit governance testing;
- future benchmarking;
- safer local-first and cloud-assisted deployments.

The cost is that Praxis requires more testing discipline than a simple CRUD app. This is necessary because the system is designed to make decisions, preserve memory, and coordinate actions across life and work domains.

---

## 37. Dependencies

Depends on:

- RFC-000 through RFC-056

Required before:

- RFC-061 Verification Scripts
- RFC-062 Benchmarking

---

## 38. Acceptance Criteria

This RFC can be accepted when:

- Testing philosophy is defined.
- Quality dimensions are defined.
- Test categories are defined.
- RFC invariant testing is required.
- Agent, prompt, memory, routing, and Space testing are covered.
- Replay and rebuild testing are covered.
- Security and privacy testing are covered.
- Non-deterministic output testing is covered.
- Regression strategy is defined.
- CI strategy is defined.
- Testing invariants are agreed upon.

---

## 39. Decision Log

| Date | Decision | Author |
|------|----------|--------|
| 2026-06-28 | Initial draft of Testing Strategy. | Tiroq + ChatGPT |
| 2026-06-28 | Defined testing as architecture, policy, agent, memory, prompt, replay, and Space verification rather than unit tests only. | Tiroq + ChatGPT |

---

> **Praxis cannot be trusted because it runs. Praxis can be trusted when its invariants are continuously verified.**