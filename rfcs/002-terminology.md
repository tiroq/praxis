

# RFC-002 Terminology

**Status:** Draft  
**Authors:** Ivan + ChatGPT  
**Last Updated:** 2026-06-28

---

## 1. Summary

This RFC establishes the canonical language of Praxis. All Praxis documents, APIs, user interfaces, database schemas, and source code **must** use these definitions consistently. A shared, unambiguous vocabulary is foundational to collaboration, clarity, and system integrity. This RFC is the single source of truth for all terminology in Praxis.

---

## 2. Relationship to Previous RFCs

- **RFC-000** defines *why* Praxis exists.
- **RFC-001** defines the *principles* guiding Praxis.
- **RFC-002 (this RFC)** defines the *common vocabulary* that underpins all Praxis artifacts.

---

## 3. Terminology Rules

- **One concept → one canonical term.** Each important concept has a single, unambiguous name.
- **Avoid synonyms.** Synonyms, abbreviations, and aliases are forbidden unless explicitly documented.
- **API, UI, and documentation use identical names.** The same term must be used everywhere, in the same form.
- **Internal names should be stable.** Once a canonical term is established, it should not change without strong justification.
- **Business terminology has priority over implementation terminology.** Use domain language wherever possible.

---

## 4. Canonical Terms

| Term        | Definition | Examples | Notes |
|------------ |----------- |----------|-------|
| **Event**   | A discrete occurrence representing a change, signal, or observation in the system. | "ArtifactCreated", "ReviewSubmitted" | Immutable; always in past tense. |
| **Artifact** | A durable, versioned object produced, modified, or reviewed within Praxis. | Proposal, Document, Model, Report | The primary unit of work and review. |
| **Review**  | A formal evaluation or feedback process applied to an Artifact. | Peer review, Automated review | Can be human or automated. |
| **Decision** | A recorded choice or resolution made after reviewing Artifacts. | Approve, Reject, Request changes | Must be explicit and recorded. |
| **Action**  | An operation carried out in response to a Decision. | Merge, Deploy, Notify | Actions are idempotent and auditable. |
| **Workflow** | A defined sequence of Events, Actions, and Decisions to achieve an outcome. | Artifact approval workflow | Automates multi-step processes. |
| **Agent**   | An autonomous or semi-autonomous actor (human, AI, or service) that performs Actions. | User, Bot, Service | Not called "Bot" or "Worker". |
| **Domain**  | A bounded area of knowledge, responsibility, or expertise. | HR, Finance, ML | Used for context boundaries. |
| **Capability** | A discrete skill or function that an Agent or system can perform. | Summarization, Translation | Used for access and routing. |
| **Project** | A collection of related Artifacts, Workflows, and Agents pursuing a shared goal. | Research project, Client implementation | |
| **Lead**    | The primary Agent responsible for a Project or Artifact. | Project Lead, Review Lead | |
| **Proposal** | An initial Artifact suggesting a change, addition, or new work. | RFC draft, Feature request | Must be reviewed. |
| **Client**  | An external system or user consuming Praxis services or APIs. | Web app, External integration | |
| **Integration** | A connection to an external system or service. | Slack integration, GitHub integration | Not called "Plugin". |
| **Provider** | An external or internal system supplying data or services to Praxis. | LLM provider, Identity provider | |
| **Policy**  | A formal rule or constraint governing behavior in Praxis. | Access policy, Review policy | |
| **Context** | The set of relevant information or state for processing an Event or Action. | User context, Project context | |
| **Knowledge** | Curated, persistent information used for reasoning or decision-making. | Knowledge base, Ontology | |
| **Memory**  | A structured, system-managed store of interactions, states, or facts. | Conversation memory, Event memory | Not called "Cache". |
| **Prompt**  | A structured input to a language model, often including system and user instructions. | System prompt, User prompt | |
| **Model**   | A machine learning or AI model used for reasoning or generation. | LLM, Classifier | |
| **Router**  | A component that directs requests or Events to the correct Capability or Agent. | Workflow router | |
| **Adapter** | A component that transforms data or interfaces between systems. | API adapter, Data adapter | |
| **Projection** | A read-optimized, derived view of state, built from Events. | Artifact list view, Review summary | Not a source of truth. |
| **State**   | The current, authoritative representation of an Artifact, Agent, or Domain. | Artifact state, Agent state | |
| **Revision** | A specific, versioned snapshot of an Artifact or state. | Document revision, Model revision | |

---

## 5. Forbidden Synonyms

| Forbidden Term | Canonical Replacement | Notes |
|---------------|----------------------|-------|
| Task          | Artifact             | |
| Job           | Workflow             | |
| Ticket        | Artifact             | |
| Bot           | Agent                | |
| Worker        | Agent                | |
| Memory Cache  | Memory               | |
| Plugin        | Integration          | |
| Extension     | Integration          | |
| Hook          | Integration          | |
| User Story    | Proposal             | |
| Issue         | Artifact             | |
| Snapshot      | Revision             | |
| View Model    | Projection           | |
| Fact Table    | Projection           | |
| Human-in-the-loop | Human Review      | |
| Workspace     | Project              | |
| Organization  | Domain               | |

---

## 6. Naming Conventions

- **RFC naming:** `RFC-XXX-title.md` (e.g. `RFC-002-terminology.md`)
- **Service naming:** Lowercase, hyphen-separated, e.g. `artifact-service`, `review-router`
- **Package naming:** Python: `praxis.artifact`, JS: `@praxis/artifact`
- **Event naming:** PascalCase, past tense, e.g. `ArtifactCreated`, `ReviewSubmitted`
- **Database tables:** Plural, snake_case, e.g. `artifacts`, `reviews`
- **API endpoints:** Plural, kebab-case, e.g. `/api/artifacts`, `/api/reviews`
- **Environment variables:** Uppercase, underscore, e.g. `PRAXIS_API_URL`
- **Prompt identifiers:** Lowercase, dot-separated, e.g. `artifact.summarize`, `review.request`

---

## 7. Lifecycle Vocabulary

The standard Praxis pipeline uses these terms, in order:

**Event → Understanding → Artifact → Reviews → Decision → Action → Learning**

Each term is defined above. "Learning" refers to system or Agent improvement based on outcomes.

---

## 8. Architectural Vocabulary

- **Event-Driven:** System design where Events trigger processing and state changes.
- **Bounded Context:** A clearly defined boundary within which a particular Domain model applies.
- **Projection:** A derived, read-optimized view of state, built from Events.
- **Source of Truth:** The canonical, authoritative store of a particular type of data (e.g., Event log or Artifact store).
- **Idempotency:** Property that an Action or Event can be applied multiple times without changing the result beyond the initial application.
- **Consensus:** Process by which multiple Agents or systems agree on a Decision or state.
- **Human Review:** Explicit, human-in-the-loop evaluation step required for certain Decisions or Actions.

---

## 9. Future Terminology Evolution

- **New terms** may only be introduced via new RFCs.
- **Existing canonical terms** must not change without strong justification and an approved RFC.
- **Proposed synonyms** must be justified and reviewed before inclusion.
- **Deprecation** of terms must be documented in this RFC with migration guidance.

---

## 10. Acceptance Criteria

- All Praxis documentation, APIs, UIs, schemas, and codebases use only canonical terms.
- No forbidden synonyms appear in any official artifact.
- Naming conventions are followed in all new code and documentation.
- Any new term or change to a canonical term is introduced via RFC.

---

## 11. Dependencies

- [RFC-000 Why Praxis Exists](./000-why-praxis-exists.md)
- [RFC-001 Principles](./001-principles.md)

---

## 12. Decision Log

- **2026-06-28:** Initial draft (Ivan + ChatGPT)

---

> "Shared language is the foundation of shared architecture."