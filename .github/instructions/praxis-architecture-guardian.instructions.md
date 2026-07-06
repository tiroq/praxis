---
description: "CONSOLIDATED — Refer to COPILOT_AUTONOMOUS_WORKFLOW.instructions.md (Phase 3 and Phase 2.2)"
name: "Praxis Architecture Guardian [DEPRECATED]"
applyTo: ["services/**", "packages/**", "apps/**", "scripts/**", "infra/**"]
---

# ⚠️ CONSOLIDATED — Use New Unified System

This file has been consolidated into the unified instruction system:

→ **Read here:** `.github/instructions/COPILOT_AUTONOMOUS_WORKFLOW.instructions.md`

Specifically:

- **Phase 2.2 - 2.3:** Extract, Don't Invent principle and abstraction validation
- **Phase 3:** Complete 11-phase Architecture Review (all phases included)
- **Section: Pure Mapping Verification:** Mapper purity rules and verification checklist
- **Section: Architectural Invariants:** All hard rules preserved

**Why consolidate?** Eliminated duplication across three files, established single entry point, created clear execution order for the complete autonomous workflow.

All content from this file is preserved in the new unified system.



### Mapper Size Guideline

Mapper size is not enforced by line count.

Instead evaluate responsibility.

A mapper should normally:

- perform one transformation
- have one input
- produce one output
- contain no nested decision trees
- remain understandable in a single screen

If a mapper requires sections, regions, helper classes, or extensive comments, architecture review is required.

### Behavior Accumulation Test

Every mapper review must answer:

1. Has this mapper become more complex than the transport object it maps?
2. Are new business cases being added to the mapper over time?
3. Is this becoming the easiest place to add unrelated logic?
4. Would another subsystem reasonably want to reuse this logic?
5. Does this logic have a better owner?

If any answer is "Yes", stop implementation and perform an architecture review.

### One Translation Rule

A mapper performs exactly one translation.

Examples:

```
Telegram Update → Praxis InputMessage

HTTP Request → Command DTO

Database Row → Projection DTO
```

A mapper must never chain multiple conceptual transformations.

If two translations exist, they should be two mappers owned by their respective boundaries.

### No Helper Gravity

Do not slowly convert mappers into utility modules.

Avoid adding:

- helper classes
- shared mapping utilities
- generic transformation libraries
- base mapper classes
- mapper inheritance
- mapper registries

Extract only after demonstrated duplication.

### Golden Mapper Comparison

**For every new transport adapter with a mapping function, mandatory comparison against the Golden Mapper reference.**

Reference implementation: `apps/telegram/main.py::telegram_update_to_payload()`

Documentation: [docs/architecture/GOLDEN_MAPPER.md](../architecture/GOLDEN_MAPPER.md)

**Comparison Rule:**

Before approving any new adapter mapper, compare it directly against the Golden Mapper.

- If the new mapper is **simpler or equal in complexity**, approve.
- If the new mapper is **more complex**, require explicit justification:
  - What additional responsibility necessitates the complexity?
  - Is this responsibility documented in the relevant RFC?
  - Could this complexity be extracted to a different owner?

**Complexity Justification Rule:**

If the new mapper is more complex than the reference and the complexity is **not justified by current code patterns or approved RFCs**, reject the mapper. Require refactoring to match the Golden Mapper pattern.

**Reference Mapper Properties:**

- One translation (external → wire contract)
- One input, one output
- 13 lines of code (excluding docstring)
- Comprehensible in under one minute
- All output fields traceable to input
- No side effects, no I/O, no state
- No business logic, no infrastructure calls
- No accumulated behavior

New mappers should match or exceed this standard.

---

## Reference Implementations

**Architectural consistency through proven examples.**

When Praxis already contains an implementation that is considered architecturally correct, new implementations of the same category should be compared against it before introducing new patterns.

A Reference Implementation is **not reusable code**. It is a reviewed example of correct ownership, responsibility, and boundaries. It exists to reduce architectural drift and prevent inconsistent design evolution.

### Categories with Reference Implementations

Reference Implementations may be defined for:

- **Transport Mapper** — translates external object to wire contract
- **Repository** — storage interface for a specific domain entity
- **Adapter** — translates between external system and internal services
- **Worker** — background job processor, event handler
- **Transport** — connection, polling, acknowledgment mechanism
- **Storage Backend** — persistence engine (event store, projection store, etc.)
- **Service Bootstrap** — service initialization, dependency injection, lifecycle
- **Configuration Loader** — environment variable parsing, validation, defaults

Other categories may be added when multiple implementations exist.

### Review Rule: Search Before Creating

**Before creating a new implementation in any category listed above:**

1. **Search** for an existing Reference Implementation in that category.
2. **If one exists:**
   - Compare the new implementation directly against it.
   - Prefer following the same structure.
   - Any deviation requires explicit justification.
   - If the new implementation is significantly more complex, identify which current requirement makes that complexity necessary.
   - If complexity exists only because of future predictions, reject the design.

3. **If none exists:**
   - Implement the new category following the general principles (purity, single responsibility, clear ownership).
   - Document whether this implementation should become the Reference Implementation for future use.

### Currently Defined Reference Implementations

#### Transport Mapper

| Aspect | Value |
|---|---|
| Reference | `apps/telegram/main.py` |
| Function | `telegram_update_to_payload()` |
| Documentation | [docs/architecture/GOLDEN_MAPPER.md](../architecture/GOLDEN_MAPPER.md) |
| Properties | One translation; one input, one output; 13 lines; pure; no side effects; no business logic |

### Mandatory Review Questions

For every new implementation in a category with a Reference Implementation, answer:

1. Does a Reference Implementation already exist for this category?
2. If yes, why is this implementation different?
3. Can the existing reference simply be copied and adapted?
4. What new requirement forces deviation from the reference?
5. Would copying the reference produce a simpler implementation?

**Stop-and-Review Rule:** If answers reveal that the new implementation is more complex without justification, stop and refactor to match the reference pattern.

### Architecture Law: Consistency Through Reference

> Prefer evolving a proven implementation over inventing a new implementation style.
>
> Consistency is an architectural asset.
>
> When a Reference Implementation exists, follow it unless you can justify deviation by current requirements or approved RFCs.

---

## Justification Requirements

For every new component, provide explicit justification:

### New Package

- Why existing packages cannot contain this?
- What boundary does this enforce?
- Which RFC mandates this separation?

### New Repository

- What aggregate or entity does it persist?
- Why existing repositories cannot handle this?
- Which RFC defines this storage boundary?

### New Engine

- What stateless transformation does it perform?
- Why existing engines cannot handle this?
- Which RFC defines this processing?

### New Runtime

- What execution model does it provide?
- Why existing runtimes cannot support this?
- Which RFC requires this runtime?

### New Storage

- What data does it persist?
- Why existing storage cannot handle this?
- Which RFC mandates this storage?

### New Adapter

- What external system does it integrate?
- Why existing adapters cannot support this?
- Which RFC requires this integration?

### New Public Interface

- What contract does it define?
- Which production caller requires it?
- Why must it be public?

### New Exported Type

- Why must this be exported?
- Which external package imports it?
- Can it be package-private instead?

**No component may exist without justification.**

---

## Architecture Self-Review

Before finalizing review, ask:

### Simplicity

- Is this simplest design satisfying RFCs?
- Can I remove any abstraction?
- Can I reuse instead of create?

### Boundaries

- Does this respect clean architecture layers?
- Does this belong to exactly one category?
- Does this violate any ownership rules?

### Future Impact

- Does this introduce second sources of truth?
- Does this create hidden coupling?
- Does this weaken auditability?

### RFC Fidelity

- Does this implement exactly what RFC specifies?
- Does this add behavior not in RFC?
- Does this contradict any RFC?

### Testability

- Can this be tested in isolation?
- Can this be verified automatically?
- Does this require manual validation?

---

## Required Deliverables

Before proceeding to implementation, produce:

### 1. RFC Summary

- List of relevant RFCs
- Quoted sections governing this work
- Any ambiguities or gaps identified

### 2. Ownership Analysis

- Single owner for each responsibility
- Clear delegation boundaries
- No overlapping ownership

### 3. Dependency Graph

- ASCII or Mermaid diagram
- All dependencies shown
- Direction validated

### 4. Runtime Data Flow

- Request → Response path
- Data transformations
- Storage operations
- Event flows

### 5. Architecture Debt Section

Explicitly document:

- Shortcuts taken (if any)
- Technical debt introduced (if any)
- Future refactoring required (if any)
- RFC gaps this work exposed (if any)

Even if none, state: "No architecture debt introduced."

### 6. Smallest Vertical Slice

Define minimal end-to-end implementation:

- Entry point
- Core logic
- Storage
- Exit point
- Verification

Ship one working path, not multiple partial paths.

---

## Stop-Before-Code Rule

**After completing 10-phase review and producing all deliverables:**

**STOP.**

**Do not proceed to code generation.**

**Wait for explicit approval.**

Present architecture review and ask:

> I have completed architecture review for [component name].
>
> Review identified [X RFCs], [Y dependencies], [Z risks].
>
> All deliverables documented above.
>
> Should I proceed with implementation?

Only after receiving explicit approval ("yes", "proceed", "go ahead") may you generate code.

**Never assume approval.**

**Never implement first and justify later.**

**Never skip stop-and-wait step.**

---

## Final Checklist

Before claiming architecture review complete, confirm:

- [ ] All 11 phases completed
- [ ] All deliverables produced
- [ ] All justifications documented
- [ ] No RFC ambiguities unresolved
- [ ] No architectural violations identified
- [ ] Dependency graph validated
- [ ] Data flow validated
- [ ] Public API minimized
- [ ] Architecture debt documented
- [ ] Smallest vertical slice defined
- [ ] **Every new abstraction justified via Phase 11 Abstraction Review**
- [ ] **No abstraction introduced on "we may need it later" grounds**
- [ ] **If adapter: mapper purity verification completed (or no mapping function required)**
- [ ] **Mapper is referentially transparent**
- [ ] **Mapper is structurally trivial**
- [ ] **Every output field is traceable to input**
- [ ] **No business decisions exist**
- [ ] **No infrastructure calls exist**
- [ ] **No state is maintained**
- [ ] **Mapper performs exactly one translation**
- [ ] **Mapper has not accumulated unrelated behavior**
- [ ] **Mapper remains the obvious place only for representation translation**
- [ ] **Architecture review completed if mapper scope expanded**
- [ ] **Golden Mapper comparison completed**
- [ ] **New mapper is not more complex than the reference, or complexity is explicitly justified by RFC or existing code patterns**
- [ ] **No mapper framework, interfaces, base classes, or shared utilities introduced without demonstrated duplication**
- [ ] **Existing Reference Implementation searched (if applicable to implementation category)**
- [ ] **Deviations from Reference Implementation explicitly justified**
- [ ] **No new implementation pattern introduced when a Reference Implementation already exists**
- [ ] **No speculative abstractions were introduced**
- [ ] **STOPPED and waiting for approval**

If any checkbox unchecked, review is incomplete.

---

## Architecture Laws

These laws are non-negotiable and permanent:

**Mappers are pure.**

Transport mapping functions contain no side effects, no business logic, no external dependencies.

**Adapters translate only.**

Adapters bridge external systems to internal wire contracts. They do not decide, enrich, or reason about data.

**Business logic belongs outside transport mappers.**

Validation, enrichment, intent classification, and workflow creation happen in engines, reducers, or coordinators—never in mappers.

**Mappers must remain structurally trivial.**

A mapper exists only to translate one representation into another.

A mapper must never become a place where policy accumulates.

If a mapper begins making business decisions, branching on domain state, or coordinating behavior, the logic belongs elsewhere.

**Mapper growth is an architectural smell.**

A mapper that continuously grows is evidence that responsibility is misplaced.

Large mappers are architectural smells, not implementation achievements.

When a mapper keeps accumulating rules, stop extending it and identify the missing owner.

**Consistency over novelty.**

Prefer evolving a proven implementation over inventing a new implementation style.

Consistency is an architectural asset. When a Reference Implementation exists, follow it unless you can justify deviation by current requirements or approved RFCs.

---

## Remember

**Architecture is always more important than implementation.**

Wrong implementation can be rewritten.

Wrong architecture pollutes codebase permanently.

Take time to get architecture right before writing single line of code.

Never skip architecture review.