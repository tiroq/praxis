

# RFC-051 Personal Space

**Status:** Draft  
**Authors:** Praxis Core Team  
**Last Updated:** 2026-06-28

---

## Summary

Personal Space is the reference implementation of [RFC-050 Space Model](./050-space-model.md), designed as the canonical environment for an individual to manage all aspects of their personal life. It provides structure, objects, agents, and policies for life management, distinct from work, finance, or product domains.

---

## Relationship to Previous RFCs

- **Depends on:** RFC-000 through RFC-050 (core abstractions, canonical objects, space model)
- **Required before:** Higher-level personal features (journaling, life review, personal agents, cross-space sync)

---

## Goals

- Provide a private, unified digital environment for an individual's life management.
- Offer a reference implementation of the Space Model tailored to personal needs.
- Enable canonical objects for personal tasks, goals, people, places, and assets.
- Support personal agents for planning, reminders, and life optimization.
- Ensure privacy, explicit sharing, and auditability of personal data.

---

## Non-Goals

- Managing organizational, financial, educational, or product-specific spaces.
- Implicit data sharing with external parties or spaces.
- Defining non-canonical or experimental object types.
- Implementing third-party integrations beyond the interface level.

---

## Personal Space Philosophy

Personal Space is the individual's sovereign digital environment. It is private by default, user-owned, and user-controlled. All personal information, decisions, and memories reside here, with explicit boundaries and clear separation from other spaces (work, finance, education, product). Agents and automations serve the individual's interests, never external ones.

---

## Scope

Personal Space covers the following life areas:

- **Health:** Wellbeing, fitness, medical, sleep, nutrition.
- **Family:** Relationships, dependents, caregiving, milestones.
- **Home:** Residence, household management, possessions.
- **Travel:** Trips, itineraries, mobility, locations.
- **Learning:** Education, reading, skills, certifications.
- **Hobbies:** Recreation, interests, clubs, personal projects.
- **Documents:** Personal files, IDs, records, contracts.
- **Goals:** Life goals, aspirations, personal projects.

---

## Identity Model

| Field        | Description                                      |
|--------------|--------------------------------------------------|
| Owner        | The primary individual (user)                    |
| Household    | Group of people sharing a home (optional)        |
| Timezone     | Default timezone for scheduling, reminders       |
| Locale       | Language, region, formatting preferences         |
| Preferences  | User settings (privacy, notifications, etc.)     |

---

## Core Canonical Objects

| Object    | Description                                 |
|-----------|---------------------------------------------|
| Task      | Actionable item with due date or context    |
| Goal      | Desired outcome, milestone, or aspiration   |
| Habit     | Recurring behavior to track or reinforce    |
| Person    | Individual known to the owner               |
| Place     | Location relevant to the owner              |
| Asset     | Owned item (device, vehicle, home, etc.)    |
| Document  | File, contract, ID, or record               |
| Event     | Scheduled occurrence or appointment         |
| Decision  | Explicit choice or fork in personal context |
| Review    | Periodic reflection or summary              |

---

## Personal Agents

| Agent             | Role / Function                                    |
|-------------------|---------------------------------------------------|
| Planner           | Suggests daily/weekly plans, prioritizes tasks    |
| Reminder          | Notifies of upcoming events, deadlines, routines  |
| Health Coach      | Tracks habits, health metrics, offers suggestions |
| Travel Planner    | Organizes trips, bookings, and itineraries        |
| Finance Assistant | Tracks expenses, budgets, and assets              |
| Family Coordinator| Manages family events, reminders, shared tasks    |

---

## Memory Model

- **Episodic Memory:** Chronological record of events, actions, and experiences within Personal Space.
- **Semantic Memory:** Structured facts about people, places, assets, and routines.
- **Preference Memory:** User choices, habits, and configuration preferences scoped to the Personal Space.

All memory is private, scoped to the individual, and never shared across spaces without explicit approval.

---

## Knowledge Graph

Personal Space maintains a knowledge graph linking:

- **People:** Relationships, contacts, family, friends.
- **Places:** Home, frequent locations, visited places.
- **Assets:** Owned items, devices, vehicles, documents.
- **Goals:** Aspirations, progress, dependencies.
- **Routines:** Habits, recurring events, daily structure.

Edges encode relationships (ownership, membership, proximity, dependency, history).

---

## Policies

- **Privacy-First:** All data is private by default.
- **Approval:** Sharing, export, or cross-space references require explicit user approval.
- **Notifications:** User controls notification preferences and channels.
- **Retention:** User can archive, export, or delete any personal data.

---

## Integrations

Personal Space supports the following integrations (optional, user-controlled):

- **Calendar:** Import/export events, sync with external calendars.
- **Email:** Reference or link emails to tasks, contacts, events.
- **Contacts:** Sync with device or cloud address books.
- **Tasks:** Import/export with external task managers.
- **Notes:** Link or embed notes from external sources.
- **Health:** Sync health metrics from devices or apps.
- **Maps:** Geolocation, travel routes, place lookups.

---

## Projections

Personal Space supports multiple views ("projections") for life management:

- **Today:** Tasks, events, and reminders for the current day.
- **This Week:** Aggregated view of upcoming tasks, events, goals.
- **Goals:** Progress and next actions for personal goals.
- **Family:** Shared events, milestones, and tasks involving family.
- **Home:** Household management, maintenance, assets.
- **Health:** Habits, metrics, health-related tasks.
- **Travel:** Upcoming trips, itineraries, and logistics.

---

## Cross-Space Communication

Personal Space interacts with Work, Finance, Education, and Product Spaces **only** via explicit references (e.g., linking a work event to a personal calendar, or referencing a financial asset). No implicit sharing or data leakage occurs; all cross-space communication is user-initiated and auditable.

---

## Security and Privacy

Personal Space is **private by default**. All data, agents, and automations operate solely for the owner's benefit. Explicit user action is required for any sharing, export, or cross-space reference. All access is logged and auditable by the owner.

---

## Lifecycle

- **Draft:** Space is being set up or imported.
- **Active:** Space is in regular use.
- **Archived:** Space is frozen for reference or export.

Transitions are explicit and reversible (except for deletion).

---

## Storage Mapping

| Store            | Purpose                            |
|------------------|------------------------------------|
| Canonical Store  | Core objects (tasks, people, etc.) |
| Event Store      | Episodic memory, actions, events   |
| Review Store     | Periodic reviews, reflections      |
| Decision Store   | Choices, forks, personal decisions |
| Vector Store     | Embeddings for search, agents      |
| Knowledge Graph  | Relationships, links, facts        |

---

## Invariants

- All objects belong to the Personal Space or are explicitly referenced.
- Memory and data **never leak** to other spaces by default.
- Cross-space sharing is always explicit and logged.
- All personal data is auditable by the owner.

---

## Architectural Consequences

- Personal Space is the foundational context for all individual-centric features.
- Clear separation from organizational, financial, and product spaces.
- All agents and automations are scoped to the individual's interests.
- Enables safe experimentation and extension without risk to privacy.

---

## Dependencies

- RFC-000 through RFC-050 (core abstractions, canonical objects, space model)
- Underlying storage and identity infrastructure
- Agent framework (for personal agents)

---

## Acceptance Criteria

- Canonical objects for all defined types are implemented and scoped to Personal Space.
- Personal agents operate only within the Personal Space context.
- All integrations require explicit user approval.
- Privacy, auditability, and retention policies are enforced by default.
- Cross-space references are explicit, logged, and user-controlled.
- Projections (Today, This Week, etc.) are available for life management.

---

## Decision Log

| Date       | Decision                                                   |
|------------|------------------------------------------------------------|
| 2026-06-28 | Initial draft and definition of Personal Space specialization. |
