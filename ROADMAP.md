# ROADMAP

## Vision

Praxis aims to be the definitive operating system for decisions, seamlessly connecting fragmented knowledge and workflows into a coherent, intelligent platform that empowers individuals and teams to turn intent into action with clarity and confidence.

## Guiding Principles

- **User-Centric**: Prioritize usability and meaningful impact on users' decision-making processes.
- **Modular Architecture**: Build components that are composable, replaceable, and extensible.
- **Transparency**: Ensure all recommendations and automations are explainable.
- **Local First**: Respect user data ownership and enable offline capabilities.
- **Vendor Independence**: Avoid lock-in by supporting multiple providers and integrations.
- **Open Collaboration**: Encourage community involvement through open RFCs and transparent design.

## Roadmap Phases

### Phase 0 — Architecture Foundation

**Focus**: Vision, Principles, Terminology, Concept Model, Capability Map, Domain Model.

**Exit Criteria**: All foundation RFCs accepted.

---

### Phase 1 — Core Platform

**Focus**: Event pipeline, storage, workflow engine, service architecture, local development.

**Exit Criteria**: Telegram message can travel through the complete pipeline into a stored decision.

---

### Phase 2 — Intelligence Layer

**Focus**: Multi-agent reasoning, review system, OmniRoute, prompt versioning, explainability.

**Exit Criteria**: Multiple AI agents collaborate and produce explainable reviews.

---

### Phase 3 — Domains

**Focus**: Freelance, Personal, Work, Products, Content.

**Exit Criteria**: Each domain has its own bounded context and workflows.

---

### Phase 4 — Automation

**Focus**: CRM, automation engine, scheduling, learning loop, knowledge graph.

**Exit Criteria**: Closed-loop automation with human approval.

---

### Phase 5 — Ecosystem

**Focus**: Plugins, API, SDK, Marketplace, Team collaboration, Enterprise.

**Exit Criteria**: External developers can extend Praxis.

---

## Milestone Timeline

### Foundation Milestones

- Repository
- Documentation
- RFCs

### Platform Milestones

- Core Runtime
- Event Pipeline
- Storage
- API
- Workers

### Product Milestones

- Upwork
- CRM
- Knowledge Graph
- Automation
- Multi-user

---

## Definition of Done

### Architecture DoD

- All foundation RFCs accepted and documented.
- Architecture diagrams and models finalized.
- Stakeholder alignment on design decisions.

### Implementation DoD

- Code meets design specifications.
- Automated tests cover critical paths.
- Peer-reviewed and merged changes.
- Development environment stable and reproducible.

### Production DoD

- Deployment pipelines fully automated.
- Performance benchmarks met or exceeded.
- User feedback incorporated.
- Monitoring and alerting in place.

---

## Current Focus

Implementation has intentionally not started. The immediate priority is completing Foundation RFCs in the following order:

README → MANIFESTO → ROADMAP → RFC-000 → RFC-001 → RFC-002 → RFC-003 → RFC-010 → RFC-011.

---

## Success Metrics

- New provider integrated without Core changes.
- New domain added without modifying existing domains.
- New workflow introduced without changing runtime.
- New model supported through OmniRoute only.
- All important decisions are traceable.
- 100% RFC coverage for major features.

---

## Roadmap Philosophy

- **Architecture Before Implementation**: Prioritize thoughtful design and planning to ensure a robust foundation.
- **RFC Before Code**: Formalize all significant changes through Request for Comments to foster consensus and clarity.
- **Tests Before Optimization**: Build reliable tests early to enable confident refactoring and performance tuning.
- **Documentation is Part of the Product**: Maintain comprehensive and up-to-date documentation as a core deliverable.
- **Architecture is a Long-Term Asset**: Invest in sustainable design to support future growth and adaptability.
- **Small Iterations Over Large Rewrites**: Favor incremental improvements to reduce risk and increase learning.
- **Build for Decades, Not Demos**: Design systems for longevity and real-world impact, not just prototypes.