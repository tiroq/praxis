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

### Phase 1: Foundation

**Objectives**
- Establish core architectural principles and development workflows.
- Set up repository, documentation, and RFC process.
- Implement basic tooling and environment for development.

**Deliverables**
- Project repository with initial codebase.
- RFC repository and contribution guidelines.
- Continuous integration and testing infrastructure.

**Exit Criteria**
- Stable development environment.
- Documented architecture and processes.
- Team alignment on vision and principles.

---

### Phase 2: Core Platform

**Objectives**
- Develop core runtime and event-driven pipeline.
- Implement foundational data models and storage.
- Enable basic user interaction and data ingestion.

**Deliverables**
- Core runtime engine.
- Event processing infrastructure.
- User interface prototypes.
- Initial integrations with calendars and Git repositories.

**Exit Criteria**
- Reliable event pipeline processing.
- Basic user workflows functional.
- Data persistence and retrieval operational.

---

### Phase 3: Intelligence Layer

**Objectives**
- Integrate AI capabilities for analysis, summarization, and recommendations.
- Build explainability framework for AI outputs.
- Develop decision modeling and review systems.

**Deliverables**
- AI-driven recommendation engine.
- Explanation and reasoning modules.
- Review and decision tracking components.

**Exit Criteria**
- AI outputs accompanied by transparent reasoning.
- Decision lifecycle fully supported.
- Positive user feedback on AI assistance.

---

### Phase 4: Ecosystem

**Objectives**
- Expand integrations with external platforms (CRM, Upwork, automation tools).
- Enable multi-user collaboration and permissions.
- Launch plugin ecosystem and marketplace.

**Deliverables**
- CRM and freelance domain integrations.
- Multi-user support with roles and access control.
- Plugin API and developer documentation.
- Marketplace platform for extensions.

**Exit Criteria**
- Seamless integration with key external systems.
- Stable multi-user collaboration.
- Active plugin ecosystem with community contributions.

---

## Milestone Timeline

| Stage   | Name                 | Description                             |
|---------|----------------------|-------------------------------------|
| Stage 0 | Repository Foundation | Establish codebase and RFC process.  |
| Stage 1 | Architecture RFCs    | Define core architecture and principles. |
| Stage 2 | Core Runtime         | Build event-driven runtime engine.   |
| Stage 3 | Event Pipeline       | Implement event ingestion and processing. |
| Stage 4 | Review System        | Develop decision review and tracking. |
| Stage 5 | Upwork Domain        | Integrate freelance work domain.     |
| Stage 6 | CRM                  | Connect customer relationship management. |
| Stage 7 | Knowledge Graph      | Build interconnected knowledge models. |
| Stage 8 | Automation           | Enable automation and workflow orchestration. |
| Stage 9 | Multi-user           | Support collaboration and permissions. |
| Stage 10| Marketplace          | Launch plugin marketplace and ecosystem. |

---

## Definition of Done

For each major stage:

- **Complete Documentation**: All design and usage documents finalized.
- **Test Coverage**: Automated tests covering critical functionality.
- **Code Review**: Peer-reviewed and merged code.
- **User Validation**: Feedback collected and incorporated from initial users.
- **Performance Benchmarks**: Meeting defined performance criteria.
- **Deployment Ready**: Stable builds ready for production use.

---

## Future Vision

- **Team Collaboration**: Enable seamless coordination across distributed teams with shared decision contexts and transparent workflows.
- **Plugin Ecosystem**: Foster a vibrant ecosystem where developers can extend and customize Praxis with plugins.
- **Marketplace**: Provide a curated marketplace for plugins, templates, and integrations to accelerate adoption and innovation.
- **Enterprise Capabilities**: Support advanced security, compliance, and scalability features to serve enterprise customers.

---

## Roadmap Philosophy

- **Architecture Before Implementation**: Prioritize thoughtful design and planning to ensure a robust foundation.
- **RFC Before Code**: Formalize all significant changes through Request for Comments to foster consensus and clarity.
- **Tests Before Optimization**: Build reliable tests early to enable confident refactoring and performance tuning.

# RFC Index

## Purpose of RFCs

Request for Comments (RFCs) serve as the primary mechanism for proposing, discussing, and documenting significant architectural and design decisions within the Praxis project. They ensure transparency, community involvement, and maintain a clear historical record of the project's evolution.

## RFC Lifecycle

- **Draft**: Initial proposal stage, open for feedback and revisions.
- **Review**: Community and team review for technical and conceptual soundness.
- **Accepted**: Agreed upon and ready for implementation.
- **Implemented**: Fully realized in code and deployed.
- **Deprecated**: Marked obsolete but retained for historical context.
- **Superseded**: Replaced by a newer RFC.

## RFC Template Sections

- **Summary**: Brief overview of the proposal.
- **Motivation**: Why this change is necessary.
- **Goals**: What the RFC intends to achieve.
- **Non-Goals**: Clarifications on what is out of scope.
- **Design**: Detailed description of the proposed solution.
- **Diagrams**: Visual aids illustrating architecture or workflows.
- **Alternatives**: Other approaches considered.
- **Open Questions**: Unresolved issues or concerns.
- **Future Work**: Potential extensions or follow-ups.

## Contribution Workflow

1. Fork the RFC repository.
2. Create a new RFC file following the naming conventions.
3. Submit a pull request for review.
4. Address feedback and iterate until consensus.
5. Merge upon approval.

## Numbering Convention

- **000** Vision
- **001** Principles
- **002** Terminology
- **010** Capability Map
- **011** Domain Model
- **012** Artifact Model
- **013** Event Model
- **020** Review System
- **021** Decision Model
- **022** State Machine
- **030** System Architecture
- **031** Service Catalog
- **032** Data Flow
- **033** Storage Model
- **040** Agent Architecture
- **041** LLM Routing
- **042** Prompt Versioning
- **050** Freelance Domain
- **051** Freelance CRM
- **052** Proposal Workflow
- **060** Testing Strategy
- **061** Verification Scripts
- **062** Benchmarking

## Recommended Reading Order

Start with foundational RFCs (000-002), then proceed to capability and domain models (010-013), followed by system architecture and components (020-042), and finally domain-specific and testing strategies (050-062).

## Architectural Governance Rules

- **No Major Feature Without an RFC**: Every significant change must be documented and reviewed.
- **Architecture First**: Design decisions precede implementation efforts.
- **One Source of Truth**: RFCs serve as the definitive reference for design intent.
- **Backward-Compatible Evolution Where Practical**: Strive to maintain compatibility to minimize disruption.

## Mapping Accepted RFCs to Implementation Phases

Accepted RFCs guide the development roadmap by defining the scope and design for each phase. Implementation teams reference relevant RFCs to ensure alignment with project goals and architectural standards throughout the lifecycle.