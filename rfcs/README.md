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
