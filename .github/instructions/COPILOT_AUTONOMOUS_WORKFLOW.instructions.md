---
name: "Copilot Autonomous Workflow"
description: "Orchestration layer for autonomous implementation in Praxis. Coordinates 11 phases from discovery through validation and final report. Delegates rule enforcement to specialized instruction files."
applyTo: ["services/**", "packages/**", "apps/**", "scripts/**", "infra/**"]
---

# Copilot Autonomous Workflow — Orchestration Layer

**Default Assumption:** The user expects complete implementation, not just analysis or planning.

This workflow orchestrates autonomous implementation from discovery through final report. Execute all phases automatically without stopping for intermediate confirmations unless genuinely blocked by missing information.

---

# 🎯 THE 11-PHASE AUTONOMOUS WORKFLOW

**Phase execution order is strict. Do not skip phases. Do not proceed if stop conditions are reached.**

## Phase 1: Discover — RFCs, Architecture, and Existing Code

**Do this first.** Never start coding without understanding what already exists.

### 1.1 RFC Discovery

- List all RFCs that might apply to this task
- Read each RFC completely
- Quote specific sections governing this work
- Identify RFC ambiguities, gaps, or contradictions
- **STOP if RFCs conflict.** Surface the conflict and wait for clarification

### 1.2 Architecture Context

- Read relevant ADRs (in `./docs/adr/`)
- Read relevant Policies
- Identify architectural invariants (see engineering-laws.instructions.md)
- Identify project instructions

### 1.3 Reference Implementation Search

Use Graphify first, then code search for:

- mapper, repository, worker, adapter, storage backend, transport, bootstrap, configuration loader, service, reducer, projection, event, command, query, aggregate

If reference implementation exists → **prefer evolving it**.

### 1.4 Graphify for Architecture Questions

Use Graphify before manual search: `graphify query`, `graphify path`, `graphify explain`

---

## Phase 2: Architectural Feasibility Gate — Is New Abstraction Needed?

**MANDATORY gate if creating:** new package, repository, service, engine, adapter, projection, public interface, composition root, interface, DTO, exported type, or manager.

**Skip to Phase 4 if only modifying existing components.**

### 2.1 Three-Point Validity Check

New abstraction justified ONLY IF at least one:

1. Two or more independent implementations already exist
2. Duplicated behavior exists and extraction reduces complexity
3. Approved RFC explicitly defines the abstraction

**Prohibited justifications:** "We may need it later", speculative extensibility, future integrations, theoretical flexibility.

**If none apply** → do not create abstraction. Use concrete implementation.

### 2.2 Architecture Review Gate Decision

**Does abstraction satisfy 2.1?**

- **YES** → Proceed to Phase 3 (11-phase architecture review)
- **NO** → Stop. Do not create abstraction.

---

## Phase 3: Complete Architecture Review (If Gate Applied)

**Only if Phase 2 = YES. Complete all phases. Do not skip.**

**Reference:** architecture-review.instructions.md

All 11 phases:

- 3.1: RFC Review
- 3.2: Existing Architecture Review
- 3.3: Responsibility Analysis
- 3.4: Ownership Analysis
- 3.5: Layer Validation
- 3.6: Component Categorization
- 3.7: Storage Rules Validation
- 3.8: Dependency Graph Validation
- 3.9: Runtime Data Flow
- 3.10: Public API Review
- 3.11: Abstraction Review — Extract, Don't Invent checklist

**After all 11 phases complete:** Explicit architectural approval to proceed to implementation.

---

## Phase 4: Implementation Planning

**Do not write code yet.**

**Reference:** implementation.instructions.md — Plan Before Coding

Plan must include:

- Relevant RFCs
- Impacted invariants
- Existing components to reuse
- New components to create (justified)
- Minimal implementation slice
- Verification strategy
- Risks

**Cannot map to RFC?** Stop. Do not invent behavior.

---

## Phase 5: Implementation — Build Minimal Vertical Slice

**Reference:** implementation.instructions.md — Reuse Before Build & Vertical Slices

Key rules:

- Reuse existing code, packages, architecture
- Prefer deletion over addition
- Prefer concrete implementations
- Avoid speculative abstractions
- Smallest end-to-end slice: `Command → Event → Projection → Query → Verification`

### 5.1 Pure Mapper Rule (If Applicable)

**Reference:** mapper.instructions.md

Mapper must be pure, single-responsibility, business-logic-free, understandable in <1 minute.

Mapper MUST NOT: call repos/databases/HTTP/LLMs, publish messages, access env vars, perform retries, generate business IDs, execute business logic.

Reference: `docs/architecture/GOLDEN_MAPPER.md`

---

## Phase 6: Continuous Validation — RFC and Architecture Compliance

**Reference:** validation.instructions.md — Continuous Validation

Validate before committing:

- Does this violate any RFC? ❌ STOP if YES
- Does this duplicate existing concept? ❌ STOP if YES
- Does this introduce second source of truth? ❌ STOP if YES
- Does this bypass Review → Decision → Action? ❌ STOP if YES
- Does this introduce mutable canonical state? ❌ STOP if YES
- Does this weaken auditability? ❌ STOP if YES

---

## Phase 7: Self-Review — Architecture Verification

**After implementation,** review independently for:

- RFC violations
- ADR violations
- Policy violations
- Architecture drift
- Duplicated ownership
- Dependency violations
- Speculative abstractions
- Unnecessary interfaces, DTOs, services
- Hidden business logic

**Immediately repair detected issues.**

---

## Phase 8: Verification — Run All Tests and Validation

**Reference:** validation.instructions.md — Verification Operations

Execute every applicable command:

```bash
go test ./...
golangci-lint run
task test
task build
task verify:rfc
task graph:rebuild
task smoke:*
task report
```

Rules:

- Only skip commands that genuinely do not exist
- Never claim commands passed unless executed
- Build/graphify-out should never be committed

---

## Phase 9: Repair Loop — Fix Issues Automatically

If verification fails:

1. Analyze failure
2. Repair issue
3. Verify again
4. Repeat

**Maximum iterations: 5.**

If cannot fix in 5 iterations, surface blocker and document.

---

## Phase 10: Final Review — Verify Architecture Health

Verify:

- Architecture remains simpler (fewer files, fewer abstractions, fewer layers)
- Ownership remains correct (single owner per responsibility)
- Boundaries remain intact (no cross-layer violations)
- No duplication introduced
- No unnecessary abstractions added
- Code follows project conventions

---

## Phase 11: Final Report — Document Complete Implementation

**Always finish with:**

- **Summary** — what was implemented, what problem it solves
- **Files Changed** — list of all modified/added files
- **Architecture Decisions** — key choices and why
- **RFC Compliance** — which RFCs implemented, were they satisfied
- **ADR Compliance** — which ADRs followed
- **Policy Compliance** — which policies applied
- **Reference Implementations Used** — what existing code was extended/reused
- **Commands Executed** — all build and verification commands run
- **Test Results** — pass/fail status of all tests
- **Validation Results** — lint, style, coverage results
- **Remaining Technical Debt** — incomplete work, workarounds, future improvements
- **Risks** — what could break, mitigation strategies
- **Next Recommended Slice** — what to build next

**Never finish with only "Done".**

---

# 📋 INSTRUCTION FILE REFERENCE

Each instruction file has a single responsibility:

| File | Responsibility |
|------|---|
| architecture-review.instructions.md | 11-phase architecture review gate (Phase 3) |
| implementation.instructions.md | Planning, reuse-before-build, vertical slices (Phases 4-5) |
| mapper.instructions.md | Pure mapper verification (Phase 5.1) |
| validation.instructions.md | Tests, build, verification, repair loops (Phases 6-9) |
| engineering-laws.instructions.md | Invariants, core principles, engineering laws (used throughout) |

---

# 🚫 STOP CONDITIONS — These Always Stop Workflow

- **RFCs conflict or ambiguous** (Phase 1) → Stop and wait for clarification
- **Multiple equally valid implementations possible** (Phase 2-3) → Stop and ask
- **Would violate established architecture** (any phase) → Stop immediately
- **Genuinely blocked by missing information** (any phase) → Stop and explain
- **New abstraction fails Extract, Don't Invent test** (Phase 2) → Stop (do not create)
- **Any phase of architecture review fails** (Phase 3) → Stop (do not proceed)
- **RFC compliance fails** (Phase 6) → Stop and fix
- **Repair loop exceeds 5 iterations** (Phase 9) → Stop and document

---

# ✅ DEFAULT CHECKLIST FOR EVERY IMPLEMENTATION

Before finishing, verify:

- [ ] All relevant RFCs read and understood
- [ ] Phase 1: Discovery complete
- [ ] If new abstraction: Phase 2 feasibility gate passed
- [ ] If new abstraction: Phase 3 (11-phase architecture review) complete
- [ ] Phase 4: Implementation plan produced and followed
- [ ] Phase 5: Minimal vertical slice implemented
- [ ] Phase 6: RFC compliance verified
- [ ] Phase 7: Self-review completed
- [ ] Phase 8: All tests and verification commands executed
- [ ] Phase 9: Issues fixed in repair loop
- [ ] Phase 10: Final review passed
- [ ] Phase 11: Final report produced
- [ ] No RFC violations remain
- [ ] No ADR violations remain
- [ ] No invariant violations remain (see engineering-laws.instructions.md)
- [ ] Architecture remains simpler
- [ ] Ownership remains correct
- [ ] Boundaries remain intact

---

# 📖 HOW TO USE THIS WORKFLOW

1. **Get a Sprint Task** — Developer provides feature to implement
2. **Run the 11 Phases** — Execute Phase 1 through Phase 11 automatically
3. **Stop only if genuinely blocked** — missing RFC clarity, architectural ambiguity, or implementation would violate architecture
4. **Produce Final Report** — Never stop at implementation; always produce complete report

**Developer expectation:** Provide task. Copilot executes complete workflow autonomously from discovery through final report.

---

# 🔍 CONTEXT GATHERING DEFAULTS

### Use Graphify First

For repository architecture, ownership, dependencies, or implementation locations:

- `graphify query` before grep
- `graphify path` for dependency questions
- `graphify explain` for architectural concepts

Only inspect source files when implementation is required, debugging, or Graphify lacks detail.

### RFC Is Truth

RFCs in `./rfcs` are **architectural source of truth** (see engineering-laws.instructions.md).

- Never implement behavior contradicting an accepted RFC
- If code and RFC disagree, RFC wins
- Stop and flag conflicts immediately

### When Unsure

- Prefer reading RFC over guessing
- If RFC ambiguous or missing, surface gap rather than filling with assumed behavior
- Do not implement when blocked by ambiguity
- Ask for clarification instead of inventing
