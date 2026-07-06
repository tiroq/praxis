# AI Tool Governance

This directory contains governance documents specific to AI tools (like GitHub Copilot) used in Praxis development.

These are **tool-specific** governance documents. They describe how AI tools should behave when working on Praxis, not the Praxis architecture itself.

## Contents

### [AI_ENGINEERING_CONSTITUTION.md](AI_ENGINEERING_CONSTITUTION.md)

**Highest authority governing AI behavior inside Praxis.**

- **Primary mission** — AI exists to improve engineering quality, not to produce code
- **Core principles** — understand, reason, verify, implement, review, improve, document, learn
- **Authority hierarchy** — RFCs govern the product; this Constitution governs AI
- **Rights and obligations** — what AI can and cannot do
- **Decision rules** — how AI should make engineering decisions

**Reference:** This is the constitution for AI tool behavior. RFCs govern the Praxis product itself.

### [AI_DECISION_SYSTEM.md](AI_DECISION_SYSTEM.md)

**Decision lifecycle and rules for AI decision-making.**

- **11-step decision process** — how AI tools make decisions
- **Decision rules** — constraints on decisions
- **Confidence framework** — how to express decision confidence
- **Forbidden decisions** — what AI tools cannot decide
- **Final principle** — when to stop and ask for human input

**Reference:** This governs how AI tools like Copilot should make architectural and implementation decisions.

## Relationship to Praxis Architecture

The **Praxis Architecture** (in `docs/architecture/`) is independent of any tool.

**AI Tool Governance** (this directory) describes how AI tools should implement that architecture.

Example:
- Architecture says: "Extract, Don't Invent" (principle)
- AI Governance says: "AI tools must validate Extract/Don't Invent before creating abstractions" (AI responsibility)

## Authority Hierarchy

1. **RFCs** (`docs/rfcs/`) — Govern Praxis product
2. **ADRs** (`docs/adr/`) — Govern Praxis architecture decisions
3. **Architecture** (`docs/architecture/`) — Govern Praxis design
4. **AI Constitution** (this directory) — Govern AI tool behavior
5. **Copilot Instructions** (`.github/copilot-instructions.md`) — Govern Copilot orchestration

An AI tool violates the Constitution if it ignores Architecture principles.

An AI tool violates the Copilot Instructions if it doesn't follow the 11-phase workflow.

## For Tool Maintainers

If adding or updating AI tool governance:

1. Ensure it does not contradict Praxis Architecture
2. Ensure it does not override RFCs or ADRs
3. Document any constraints on AI decision-making
4. Keep it independent of specific AI vendors or models
