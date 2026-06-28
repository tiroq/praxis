

# RFC-054 Freelance Space

---

## Status

**Status:** Draft  
**Authors:** Praxis Core Team  
**Last Updated:** 2026-06-28

---

## Summary

Freelance Space is a specialization of RFC-050 (Business Space) focused on supporting independent consulting, contracting, and client delivery workflows. It provides a canonical, repeatable business architecture for freelancers to manage the entire client lifecycle—from lead generation through delivery, invoicing, and reputation management—enabling freelancers to operate as sustainable, knowledge-driven businesses rather than as ad hoc gig workers.

---

## Relationship to Previous RFCs

- **Depends on:** RFC-000 through RFC-053 (Core, Data, Product, Work, Business, etc.)
- **Specializes:** RFC-050 (Business Space)
- **Required before:** Finance RFCs, Benchmark RFCs, and any domain-specific freelance automation.

---

## Goals

- Provide a reusable, end-to-end business workflow for independent consulting and contracting.
- Standardize the freelance lifecycle for repeatable, scalable client delivery.
- Accumulate business knowledge and reputation across engagements.
- Enable integration with external platforms (Upwork, LinkedIn, Stripe, etc.).
- Support robust opportunity qualification, proposal generation, and contract management.
- Ensure client trust, transparency, and business continuity.

---

## Non-Goals

- Does not prescribe implementation details for each integration.
- Does not cover employee/agency models (see RFC-050).
- Does not address personal finances or tax automation (see Finance RFCs).
- Does not automate domain-specific delivery (see Product/Work Spaces).

---

## Freelance Philosophy

Freelance work is treated not as a series of isolated gigs, but as a repeatable, knowledge-driven business system. The freelancer is positioned as a business operator, building a pipeline of opportunities, accumulating reputation, and delivering consistent value. Each engagement is a building block in a compounding business, not a one-off transaction.

---

## Scope

- Lead generation and intake
- Opportunity management and qualification
- Proposal and estimate creation
- Contract negotiation and agreement
- Project and milestone tracking
- Delivery management
- Invoicing and payment processing
- Reputation and review management
- Referral generation and tracking

---

## Freelance Lifecycle

`Lead → Qualification → Proposal → Negotiation → Contract → Delivery → Review → Payment → Referral`

Each stage is tracked, with transitions recorded and knowledge accumulated for future opportunities.

---

## Business Identity Model

| Entity      | Description                                         |
|------------|-----------------------------------------------------|
| Business   | The freelance business entity (sole proprietorship, LLC, etc.) |
| Brand      | Public-facing identity (name, logo, positioning)    |
| Consultant | The freelancer as service provider                  |
| Client     | The organization or person purchasing services      |
| Partner    | Collaborators, subcontractors, or referral sources |

---

## Sales Pipeline Model

The pipeline is a staged funnel:

1. **Leads:** Raw inbound/outbound opportunities.
2. **Qualified Opportunities:** Screened for fit and viability.
3. **Proposals:** Formal offers sent to clients.
4. **Negotiation:** Terms and scope are refined.
5. **Contracts:** Signed agreements.
6. **Active Engagements:** Work in progress.
7. **Closed (Won/Lost):** Engagements completed or lost.

---

## Core Canonical Objects

| Object        | Description                                                         |
|--------------|---------------------------------------------------------------------|
| Lead         | Initial contact or inquiry                                          |
| Opportunity  | Qualified business opportunity                                      |
| Proposal     | Formal scope, pricing, and terms offer                              |
| Contract     | Signed, legally binding agreement                                   |
| Client       | Organization or individual purchasing services                      |
| Engagement   | Specific project or retainer agreement                              |
| Milestone    | Key delivery checkpoints within an engagement                       |
| Deliverable  | Output or artifact delivered to client                              |
| Invoice      | Bill for services rendered                                          |
| Payment      | Record of payment received                                          |
| Review       | Client feedback on engagement                                       |
| Referral     | New lead or client sourced from an existing relationship            |

---

## Freelance Agents

| Agent              | Role                                                        |
|--------------------|------------------------------------------------------------|
| Lead Analyzer      | Screens and qualifies inbound leads                         |
| Proposal Writer    | Drafts proposals and estimates                             |
| Estimator          | Calculates pricing and delivery timelines                   |
| Contract Reviewer  | Ensures contract terms are clear and enforceable            |
| Delivery Planner   | Breaks engagements into milestones and deliverables         |
| Client Success Agent | Manages client relationship and satisfaction              |
| Invoice Assistant  | Prepares and tracks invoices and payments                   |

---

## Opportunity Qualification Model

Opportunities are qualified using criteria such as:

- Client fit (industry, size, values)
- Budget availability
- Timeline alignment
- Project scope clarity
- Strategic value (reputation, portfolio, learning)
- Likelihood of closing

Qualification is scored and recorded for each opportunity.

---

## Proposal Model

Proposals consist of:
- Executive summary
- Scope of work
- Timeline and milestones
- Pricing and payment terms
- Deliverables
- Assumptions and exclusions
- Acceptance and next steps

Proposals are versioned and tracked.

---

## Pricing Model

Supported pricing structures:
- **Hourly:** Time-based billing at a set rate.
- **Fixed Price:** Lump sum for defined scope.
- **Retainer:** Recurring monthly/periodic fee for ongoing services.
- **Value-Based:** Pricing aligned to client outcomes or value delivered.

Each engagement records its pricing model and rationale.

---

## Contract Model

Contracts include:
- Parties and legal entities
- Scope and deliverables
- Payment terms and schedule
- Intellectual property and confidentiality clauses
- Termination and dispute resolution
- Signatures and effective dates

Contracts are versioned and securely stored.

---

## Delivery Model

Delivery is structured by:
- Engagement (project or retainer)
- Milestones (major checkpoints)
- Deliverables (outputs per milestone)
- Communication cadence (updates, demos, reviews)
- Change management (scope changes, re-estimates)

Progress is tracked and reported.

---

## Client Relationship Model

Client relationships are managed across:
- Communication history
- Feedback and reviews
- Account status (active, dormant, former)
- Referral and upsell opportunities
- Knowledge of client preferences and context

All interactions are logged for future reference.

---

## Reputation Model

Reputation is accumulated via:
- Client reviews and ratings
- Engagement outcomes (on-time, on-budget, NPS)
- Public testimonials and case studies
- Referral activity
- Portfolio of completed work

Reputation is surfaced in proposals and lead generation.

---

## Referral Model

Referrals are tracked as:
- Source (client, partner, network)
- Status (pending, accepted, converted)
- Value (revenue, strategic fit)
- Attribution (who referred whom)

Referral performance is measured and rewarded.

---

## Freelance Knowledge Graph

The knowledge graph links:

- **Clients** to **Engagements**, **Reviews**, **Referrals**
- **Opportunities** to **Proposals**, **Contracts**, **Invoices**
- **Deliverables** to **Milestones** and **Engagements**
- **Referrals** to **Sources** and **Converted Clients**
- **Reputation** to **Public Profiles** and **Proposals**

**Example:**  
`Client A → Engagement X → Milestone 1 → Deliverable Y → Review Z`  
`Client B → Referral from Client A → Opportunity → Proposal → Contract`

---

## Freelance Memory Model

All business interactions, decisions, documents, and outcomes are stored in a structured, queryable memory. This enables:
- Reuse of successful proposals and contracts
- Analysis of win/loss patterns
- Data-driven improvements to qualification and delivery
- Preservation of client context across engagements

---

## Policies

- All client data is confidential by default.
- Proposals and contracts require explicit versioning and approval.
- Invoices must match delivered milestones.
- Reviews and referrals are solicited after successful delivery.
- AI agents must not take irreversible actions without human approval.

---

## AI Governance

- Agent actions are auditable and reversible where possible.
- Sensitive actions (contract signing, pricing changes) require explicit confirmation.
- AI must adhere to business policies and legal requirements.
- Human override is always available.

---

## Integrations

- **Upwork:** Lead and contract intake, reputation sync
- **LinkedIn:** Lead generation, client research
- **Email:** Communication tracking, proposal delivery
- **Calendar:** Scheduling milestones and meetings
- **Stripe:** Invoice and payment processing
- **GitHub/GitLab:** Delivery artifact tracking
- **CRM:** Client and opportunity management
- **Google Drive:** Document storage and sharing

---

## Projections

- **Pipeline:** Forecast of potential revenue and client workload
- **Revenue Forecast:** Expected earnings by period
- **Active Clients:** Current engagements and statuses
- **Deliverables:** Upcoming and overdue outputs
- **Payments:** Pending, received, overdue
- **Reputation:** Trend of reviews, ratings, and referrals

---

## Cross-Space Communication

- **Product Space:** Links freelance delivery to product development
- **Work Space:** Shares tasks, milestones, and delivery status
- **Finance Space:** Syncs invoices, payments, and financial reports
- **Personal Space:** Surfaces business events impacting personal schedule

---

## Security Model

- All business data encrypted at rest and in transit
- Role-based access control (freelancer, client, partner)
- Audit logs for sensitive actions (contract signing, payment)
- Secure integration tokens for third-party services
- Client data never shared without explicit consent

---

## Lifecycle

1. Intake lead
2. Qualify opportunity
3. Draft and send proposal
4. Negotiate and finalize contract
5. Plan and deliver engagement
6. Invoice and track payment
7. Capture review and referral
8. Archive and analyze engagement for future improvement

---

## Storage Mapping

| Object        | Storage Location              | Retention Policy         |
|--------------|------------------------------|-------------------------|
| Lead         | CRM database                  | 2 years or conversion   |
| Proposal     | Document store (Drive, CRM)   | 5 years                 |
| Contract     | Secure contract vault         | 7 years (legal)         |
| Engagement   | Project management system     | Active + 3 years        |
| Deliverable  | GitHub/Drive                  | Per client agreement    |
| Invoice      | Finance system                | 7 years (legal)         |
| Payment      | Finance system                | 7 years (legal)         |
| Review       | CRM + public profile          | Indefinite              |
| Referral     | CRM database                  | 5 years                 |

---

## Failure Modes

| Failure                   | Mitigation                                    |
|--------------------------|-----------------------------------------------|
| Lost lead or opportunity | Automated reminders, logging, and alerts      |
| Missed invoice/payment   | Payment tracking, overdue alerts              |
| Contract dispute         | Versioned contracts, audit trail              |
| Data loss                | Regular backups, encrypted storage            |
| AI agent error           | Human override, action logging                |

---

## Invariants

- Every engagement is linked to a client, contract, and at least one deliverable.
- No invoice is issued without a corresponding milestone or deliverable.
- All client-facing documents are versioned and auditable.
- Reputation is only updated after verified delivery.
- Referrals are attributed to their source.

---

## Architectural Consequences

- Enables repeatable, scalable freelance business operations.
- Reduces risk of lost knowledge between engagements.
- Increases client trust via transparency and professionalism.
- Facilitates automation of routine business tasks.
- Supports future benchmarking and financial analysis.

---

## Dependencies

- RFC-000 through RFC-053 (Core, Data, Product, Work, Business, etc.)
- Secure storage and integration infrastructure.
- Identity and access management.

---

## Acceptance Criteria

- All freelance lifecycle stages are modeled as canonical objects.
- Opportunity, proposal, contract, and delivery workflows are reusable.
- Knowledge and reputation accumulate across engagements.
- Integrations support at least two external platforms.
- All major failure modes are mitigated.
- Explicit decision log is maintained.

---

## Decision Log

- **2026-06-28:** Adopted canonical freelance lifecycle as business system, not gig sequence.
- **2026-06-28:** Required explicit versioning and audit trails for all proposals and contracts.