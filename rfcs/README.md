# RFC Index

RFCs are the architectural source of truth for Praxis. They define the standards, designs, and decisions that guide the implementation of the system. Implementation follows accepted RFCs to ensure consistency and alignment across the project.

## Purpose

RFCs exist to document architectural decisions, provide a clear rationale for design choices, and facilitate communication among stakeholders. They help ensure that the system evolves in a coherent and maintainable way.

## RFC Lifecycle

- **Draft**: Initial proposal, open for discussion and refinement.
- **Review**: Under active consideration by the community or architecture team.
- **Accepted**: Approved as the official architectural guidance.
- **Implemented**: The design has been realized in code or infrastructure.
- **Deprecated**: Marked as outdated and should no longer be used for new work.
- **Superseded**: Replaced by a newer RFC that updates or corrects the original.

## RFC Template

Every RFC must contain the following sections:

- **Summary**: A brief overview of the proposal.
- **Motivation**: The reasons driving the need for this RFC.
- **Goals**: What the RFC aims to achieve.
- **Non-Goals**: What is explicitly out of scope.
- **Definitions**: Key terms and concepts used.
- **Design**: The proposed architecture or solution.
- **Diagrams**: Visual aids illustrating the design.
- **Examples**: Use cases or scenarios demonstrating the design.
- **Alternatives Considered**: Other options explored and reasons for rejection.
- **Open Questions**: Unresolved issues or areas needing further discussion.
- **Future Work**: Potential extensions or follow-ups.
- **Decision Log**: Record of decisions and revisions.

## Governance Rules

- Architecture before implementation.
- No major feature without an RFC.
- One source of truth.
- Human-reviewed decisions.
- Backward-compatible evolution where practical.
- Update RFC before changing behavior.

## RFC Numbering

**Foundation:**  
- RFC-000 Vision  
- RFC-001 Principles  
- RFC-002 Terminology  

**Core Architecture:**  
- RFC-010 Capability Map  
- RFC-011 Domain Model  
- RFC-012 Artifact Model  
- RFC-013 Event Model  

**Decision Engine:**  
- RFC-020 Review System  
- RFC-021 Decision Model  
- RFC-022 State Machine  

**Platform:**  
- RFC-030 System Architecture  
- RFC-031 Service Catalog  
- RFC-032 Data Flow  
- RFC-033 Storage Model  

**Intelligence:**  
- RFC-040 Agent Architecture  
- RFC-041 LLM Routing  
- RFC-042 Prompt Versioning  
- RFC-043 Memory & Knowledge Graph  

**Domains:**  
- RFC-050 Freelance Domain  
- RFC-051 Freelance CRM  
- RFC-052 Proposal Workflow  
- RFC-053 Personal Domain  
- RFC-054 Work Domain  
- RFC-055 Product Domain  

**Quality:**  
- RFC-060 Testing Strategy  
- RFC-061 Verification Scripts  
- RFC-062 Benchmarking  

## Recommended Reading Order

Newcomers should start with the Foundation RFCs, then proceed through Core Architecture, Decision Engine, Platform, Intelligence, Domains, and finally Quality to gain a comprehensive understanding of Praxis.

## Contribution Workflow

The process for contributing RFCs is: discuss → draft RFC → review → accepted → implementation → verification → documentation update.

## Relationship Between RFCs and Code

Every major implementation references one or more accepted RFCs. Code should not diverge from the architecture defined in these RFCs to maintain consistency and integrity across the system.

## Philosophy

RFCs document reasoning, not just decisions. They preserve long-term architectural consistency for Praxis by capturing the context, motivation, and trade-offs behind each choice.

# RFC Index

RFCs are the architectural constitution of Praxis. They are the single source of truth for all architectural decisions, principles, and standards. Implementation always follows architecture—not the other way around. This repository of RFCs ensures that the evolution of Praxis is deliberate, consistent, and well-reasoned, not accidental or ad hoc.

## Documentation Hierarchy

Praxis documentation is organized in a strict hierarchy, with each layer building on the one above it:

```
     README.md         ← Project intro, navigation
        │
   MANIFESTO.md        ← Vision, purpose, values
        │
   ROADMAP.md          ← High-level milestones, priorities
        │
      RFCs/            ← Architectural constitution (this folder)
        │
  Implementation       ← Code, tests, infrastructure
```

- **README.md**: Entry point and navigation for the project.
- **MANIFESTO.md**: The vision, values, and high-level purpose of Praxis.
- **ROADMAP.md**: Major milestones, priorities, and planned work.
- **RFCs**: The architectural source of truth. All design and implementation must align with accepted RFCs.
- **Implementation**: The actual code, which must trace back to one or more RFCs.

## Why RFCs Exist

RFCs are essential for:

- **Architectural consistency**: Preventing accidental divergence and "architecture drift."
- **Preserving reasoning**: Capturing not only decisions, but the context, motivation, and trade-offs behind them.
- **Onboarding**: Providing new contributors with a clear, structured path to understanding the system.
- **Historical context**: Explaining why things are the way they are—even years later.
- **Avoiding architectural drift**: Ensuring that the implementation always follows the accepted architecture, not vice versa.

## RFC Lifecycle

| Status        | Description                                                                                   |
|-------------- |---------------------------------------------------------------------------------------------|
| Draft         | Initial proposal, open for discussion and feedback.                                          |
| Review        | Under active review by the architecture team and stakeholders.                               |
| Accepted      | Approved as the official architectural guidance; implementation may begin.                   |
| Implemented   | Fully realized in code, infrastructure, and documentation.                                   |
| Deprecated    | Outdated; should not be used for new work but may still affect legacy systems.               |
| Superseded    | Replaced by a newer RFC that updates, corrects, or replaces the original.                   |
| Rejected      | Not accepted; either withdrawn or not aligned with project goals.                            |

## RFC Template: Required Sections

Every RFC must contain the following sections, each with a clear heading and sufficient detail:

### Summary
A concise overview of the proposal and its intent.

### Motivation
Why is this RFC needed? What problems does it solve? What is the context?

### Goals
The specific objectives this RFC aims to achieve.

### Non-Goals
Explicitly state what is out of scope or not addressed by this RFC.

### Definitions
Key terms, concepts, and abbreviations used in the RFC.

### Design
The proposed architecture, solution, or approach. Include rationale for major choices.

### Architecture Decisions
Enumerate all significant decisions made, including trade-offs and justifications.

### Diagrams
Visual representations (UML, sequence, data flow, etc.) that clarify the design.

### Data Model
Description of data structures, schemas, or models involved.

### Examples
Concrete scenarios, sample flows, or use cases demonstrating the design in practice.

### Alternatives Considered
Other approaches that were evaluated, and reasons for their rejection.

### Migration Strategy
How to transition from the current state to the new architecture, including backward compatibility.

### Open Questions
Unresolved issues or areas needing further investigation or feedback.

### Future Work
Potential extensions, follow-ups, or related work that is out of scope for this RFC.

### Decision Log
Chronological record of major changes, feedback, and final acceptance.

## Architecture Dependency Graph

The architectural reasoning flows from vision all the way to implementation. Each layer depends on the integrity of those above it:

```
Vision
  └─> Principles
        └─> Terminology
              └─> Capability Map
                    └─> Domain Model
                          └─> Artifact/Event Model
                                └─> Review/Decision Model
                                      └─> System Architecture
                                            └─> Agent Architecture
                                                  └─> Domains
                                                        └─> Implementation
```

## RFC Numbering

| RFC  | Title               | Status    | Purpose                                                    |
|------|---------------------|-----------|------------------------------------------------------------|
| 000  | Vision              | Planned   | High-level vision, purpose, and long-term direction        |
| 001  | Principles          | Planned   | Foundational architectural principles and values           |
| 002  | Terminology         | Planned   | Definitions of key terms and concepts                      |
| 010  | Capability Map      | Planned   | Top-level system capabilities and boundaries               |
| 011  | Domain Model        | Planned   | Core business/domain entities and relationships            |
| 012  | Artifact Model      | Planned   | Structure and lifecycle of system artifacts                |
| 013  | Event Model         | Planned   | Events and their role in system state transitions          |
| 020  | Review System       | Planned   | Architecture of review/decision workflows                  |
| 021  | Decision Model      | Planned   | How decisions are represented, tracked, and enforced       |
| 022  | State Machine       | Planned   | State transitions and logic for major entities             |
| 030  | System Architecture | Planned   | High-level technical architecture and components           |
| 031  | Service Catalog     | Planned   | List and responsibilities of all major services            |
| 032  | Data Flow           | Planned   | Data movement, pipelines, and integration points           |
| 033  | Storage Model       | Planned   | Storage technologies, schemas, and access patterns         |
| 040  | Agent Architecture  | Planned   | Design of intelligent agents and their interactions        |
| 041  | LLM Routing         | Planned   | How large language models are routed and orchestrated      |
| 042  | Prompt Versioning   | Planned   | Versioning and management of prompts for LLMs              |
| 043  | Memory & Knowledge  | Planned   | Memory, context, and knowledge graph architecture          |
| 050  | Freelance Domain    | Planned   | Domain-specific logic for freelancing                      |
| 051  | Freelance CRM       | Planned   | Customer relationship management for freelance workflows   |
| 052  | Proposal Workflow   | Planned   | Handling proposals and their lifecycle                     |
| 053  | Personal Domain     | Planned   | Personal information management                            |
| 054  | Work Domain         | Planned   | Work item, project, and task management                    |
| 055  | Product Domain      | Planned   | Product and offering management                            |
| 060  | Testing Strategy    | Planned   | Testing and quality assurance approaches                   |
| 061  | Verification Scripts| Planned   | Scripts and tools for automated verification               |
| 062  | Benchmarking        | Planned   | Benchmarking, performance, and reliability                 |

## Recommended Reading Order

For new contributors, the recommended path is:

1. **RFC-000 Vision**: Understand the high-level purpose and direction.
2. **RFC-001 Principles**: Learn the core architectural values and constraints.
3. **RFC-002 Terminology**: Get familiar with the project's vocabulary.
4. **RFC-010 Capability Map**: See what the system is fundamentally designed to do.
5. **RFC-011 Domain Model**: Understand the main entities and their relationships.
6. **RFC-012 Artifact Model** and **RFC-013 Event Model**: Learn about artifacts and events.
7. **RFC-020+**: Explore review, decision, and state models.
8. **RFC-030+**: Study system, service, and data architecture.
9. **RFC-040+**: Delve into intelligence, agent, and LLM architecture.
10. **RFC-050+**: Review domain-specific RFCs.
11. **RFC-060+**: Finish with quality and verification strategies.

## Contribution Workflow

1. **Discuss**: Raise an issue or start a discussion to propose a change or new feature.
2. **Draft**: Write an RFC draft following the template above.
3. **Architecture Review**: Submit the RFC for review by the architecture team and stakeholders.
4. **Revision**: Incorporate feedback; revise until consensus is reached.
5. **Acceptance**: The RFC is formally accepted and assigned the "Accepted" status.
6. **Implementation**: Implementation begins, referencing the accepted RFC.
7. **Verification**: Ensure the implementation matches the RFC; update status to "Implemented."
8. **Documentation Update**: Update documentation and cross-reference the RFC in code and docs.

## Governance Rules

The following rules are mandatory for all contributors and maintainers:

1. **Architecture before implementation**: No code or feature may be implemented before the corresponding RFC is accepted.
2. **No feature without RFC**: Every major feature or change must have an accepted RFC.
3. **One source of truth**: RFCs are the only authoritative source for architecture.
4. **Update RFC before changing behavior**: If implementation needs to diverge, the RFC must be updated and re-accepted first.
5. **Code references accepted RFCs**: Major modules and pull requests must reference relevant RFCs.
6. **Breaking changes require a new RFC**: Any backward-incompatible change must go through the RFC process.
7. **Architectural discussions belong in RFCs**: Major decisions and debates must be recorded in RFCs, not in chat or code comments.

## Relationship Between RFCs and Code

Traceability from RFCs to code is required. Every significant module, component, or pull request must reference the relevant RFC(s). This ensures that all implementation aligns with accepted architecture, and anyone can trace the "why" behind any major feature or decision.

## Long-Term Philosophy

RFCs are not just about recording what was decided—they are about preserving the reasoning and context behind every architectural choice. An RFC should remain understandable and relevant years later, even for someone new to the project. This is how Praxis maintains architectural integrity across time, contributors, and changing requirements.

## Next RFCs

The immediate next RFCs to be drafted and accepted are:

- [ ] RFC-000 Vision
- [ ] RFC-001 Principles
- [ ] RFC-002 Terminology
- [ ] RFC-010 Capability Map
- [ ] RFC-011 Domain Model