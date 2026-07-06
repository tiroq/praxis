# Engineering Laws & Invariants

These are non-negotiable. Any code violating one is wrong by definition. Repair immediately.

---

## Architectural Invariants

- **Events are immutable.** Never mutate or delete emitted event; correct via new events.
- **Decisions are explicit and auditable.** Every Decision records who/what, why, inputs.
- **Reviews never commit Decisions.** Review produces findings; cannot enact Decision.
- **Agents never mutate canonical state directly.** They propose; system commits.
- **Agents never call LLM providers directly.** All model access goes through LLM router.
- **Prompt versions are immutable after release.** Released prompts frozen; ship new version.
- **Memory is policy-bound.** Reads/writes to memory honor governing policy; no ad-hoc access.
- **Spaces are bounded contexts.** Keep models, data, logic within their space.
- **Cross-space communication is explicit.** Use defined contracts/events; no hidden coupling.
- **Derived stores are rebuildable.** Never treat projection/cache as source of truth.

---

## Core Principles

### Extract, Don't Invent

Abstractions are extracted from existing code or approved RFCs. Never predicted or invented from future needs.

See: [architecture-review.md](architecture-review.md) (Phase 11 — Abstraction Review)

### RFC Is Source of Truth

RFCs in `./rfcs/` are **immutable architectural contracts**.

- Never reinterpret them
- Never extend with assumptions
- Never fill gaps with "reasonable" guesses
- Never implement behavior not explicitly specified

If RFC is ambiguous, incomplete, or contradictory:

- **STOP immediately**
- Document specific ambiguity
- Explain what cannot be determined
- Wait for clarification

Do not proceed when architecture is unclear.

### Reference First

If a Reference Implementation exists, follow it. Consistency is preferred over novelty.

### Single Ownership Rule

Every responsibility has exactly one owner. Never duplicate ownership.

### Boundary Rule

- Kernel never knows transport
- Kernel never knows storage implementation
- Storage never owns business logic
- Transport never owns business logic
- Composition roots own wiring only

---

## Engineering Laws

### Simplicity Rule

Prefer:

- Fewer files
- Fewer abstractions
- Fewer interfaces
- Fewer services
- Fewer layers
- Less indirection

**Architecture should become simpler after every implementation.**

### Minimize Architectural Debt

- Prefer deleting code over adding abstractions
- Prefer composition over inheritance
- Prefer explicit behavior over magic
- Prefer deterministic behavior over convenience

### Challenge Assumptions

Do not blindly implement requested design. If simpler solution satisfies RFCs, explain and recommend it before implementing. If request contradicts long-term architecture, explain why instead of complying silently.

### Leave the Repository Healthier

Every change should improve at least one of: documentation, verification, naming, comments, tests, architecture consistency, or dead-code removal.

---

## Avoid

- Speculative implementation of future features
- Generic frameworks or unused abstractions
- Placeholder services
- Speculative config
- Ad-hoc build scripts
- Exporting internal implementation details
- "Future-proof" interfaces

Everything must be justified by existing RFC or accepted implementation slice.

---

## Default Assumption

When user asks to implement a feature: **assume implementation is requested.**

Do not ask unnecessary clarification questions.

Make reasonable architectural decisions using:

- RFCs
- ADRs
- Policies
- existing repository patterns
- reference implementations

**Only stop when:**

- RFCs conflict
- architectural ambiguity makes multiple incompatible implementations equally valid
- implementation would violate established architecture
