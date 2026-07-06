---
name: "Praxis Copilot Instructions"
description: "Active Copilot instruction entry point. Coordinates 11 phases for RFC-first, Extract-Don't-Invent implementation."
applyTo: ["services/**", "packages/**", "apps/**", "scripts/**", "infra/**"]
---

# Praxis Copilot Instructions

**ONLY active Copilot instruction file.** Single entry point for all implementation tasks.

**Default assumption:** User expects complete implementation via 11-phase autonomous workflow.

---

# 🎯 11-PHASE WORKFLOW

**Strict execution order. Do not skip. Do not proceed past stop conditions.**

## Phase 1: Discover

**Do this first.** Never start without understanding what already exists.

- List RFCs that might apply
- Read each completely, quote sections governing this work
- **STOP if RFCs conflict or ambiguous**
- Read relevant ADRs and Policies
- Use Graphify first (graphify query, path, explain)
- Search for reference implementations (mapper, repository, adapter, etc.)
- **Prefer evolving existing code**

---

## Phase 2: Architectural Feasibility Gate

**IF creating** new package/service/engine/adapter/projection/interface/DTO/manager: Apply **Extract, Don't Invent** test.

**Abstraction justified ONLY IF:**

1. Two+ independent implementations already exist, OR
2. Duplicated behavior exists and extraction reduces complexity, OR
3. Approved RFC explicitly defines it

**Prohibited:** "We may need it later", speculative extensibility, future integrations.

**If gate fails** → do not create abstraction. Use concrete implementation.

**Skip to Phase 4** if only modifying existing components.

---

## Phase 3: Architecture Review (If Gate = YES)

Complete 11-phase review. See [docs/copilot/architecture-review.md](docs/copilot/architecture-review.md).

Phases: RFC Review → Existing Arch → Responsibility → Ownership → Layer Validation → Component Categorization → Storage Rules → Dependency Graph → Runtime Data Flow → Public API → Abstraction Review (Extract/Invent).

**Stop if any phase fails.** Abstraction MUST be removable or justified.

---

## Phase 4: Implementation Planning

**Do not write code yet.**

Plan must include:

- Relevant RFCs
- Impacted invariants (see [docs/copilot/engineering-laws.md](docs/copilot/engineering-laws.md))
- Existing components to reuse
- New components (each justified)
- Minimal vertical slice
- Verification strategy
- Risks

**Cannot map to RFC?** Stop. Do not invent behavior.

---

## Phase 5: Implementation — Minimal Vertical Slice

Reuse-before-build. Vertical slices: `Command → Event → Projection → Query → Verification`

### 5.1 Pure Mapper (If applicable)

Every transport mapper MUST be pure. See [docs/copilot/mapper.md](docs/copilot/mapper.md).

7-point verification checklist. Referential transparency. Structural triviality test. Reference: `docs/architecture/GOLDEN_MAPPER.md`

---

## Phase 6: Continuous Validation

Before committing code:

- Does this violate any RFC? → **STOP if YES**
- Does this duplicate existing concept? → **STOP if YES**
- Does this introduce second source of truth? → **STOP if YES**
- Does this bypass Review → Decision → Action? → **STOP if YES**
- Does this introduce mutable canonical state? → **STOP if YES**
- Does this weaken auditability? → **STOP if YES**

---

## Phase 7: Self-Review

After implementation, review for:

- RFC violations
- ADR violations
- Policy violations
- Architecture drift
- Duplicated ownership
- Dependency violations
- Speculative abstractions
- Unnecessary interfaces/DTOs/services
- Hidden business logic

**Immediately repair detected issues.**

---

## Phase 8: Verification

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

See [docs/copilot/validation.md](docs/copilot/validation.md) for full verification rules.

---

## Phase 9: Repair Loop

If verification fails:

1. Analyze failure
2. Repair
3. Verify again
4. Repeat

**Maximum iterations: 5.** If unfixable, surface blocker.

---

## Phase 10: Final Review

Verify:

- Architecture simpler (fewer files, abstractions, layers)
- Ownership correct (single owner per responsibility)
- Boundaries intact (no cross-layer violations)
- No duplication
- No unnecessary abstractions
- Code follows conventions

---

## Phase 11: Final Report

Document:

- **Summary** — what, why
- **Files Changed**
- **Architecture Decisions** — key choices and why
- **RFC Compliance** — which RFCs, satisfied?
- **ADR Compliance**
- **Policy Compliance**
- **Reference Implementations Used**
- **Commands Executed** — all build/verification commands run
- **Test Results**
- **Validation Results**
- **Remaining Technical Debt**
- **Risks** — what could break, mitigations
- **Next Recommended Slice**

**Never finish with only "Done".**

---

# 🚫 STOP CONDITIONS

These always stop the workflow:

- **RFCs conflict or ambiguous** → Stop, wait for clarification
- **Multiple equally valid implementations** → Stop, ask
- **Violates established architecture** → Stop immediately
- **Blocked by missing information** → Stop, explain
- **Abstraction fails Extract, Don't Invent** → Stop, do not create
- **Architecture review phase fails** → Stop, do not proceed
- **RFC compliance fails** → Stop, fix
- **Repair loop exceeds 5 iterations** → Stop, document

---

# 🏛️ ARCHITECTURAL INVARIANTS

Non-negotiable. Repair immediately if violated:

- Events are immutable (never mutate; correct via new events)
- Decisions are explicit and auditable
- Reviews never commit Decisions (only findings)
- Agents never mutate canonical state directly (propose only)
- Agents never call LLM providers directly (use LLM router)
- Prompt versions immutable after release
- Memory is policy-bound (no ad-hoc access)
- Spaces are bounded contexts
- Cross-space communication explicit
- Derived stores rebuildable (never source of truth)

---

# 💡 CORE PRINCIPLES

### RFC Is Source of Truth

RFCs in `./rfcs/` are immutable contracts.

- Never reinterpret
- Never extend with assumptions
- Never fill gaps with "reasonable" guesses
- If ambiguous → **STOP, wait for clarification**

### Extract, Don't Invent

Abstractions extracted from code or RFCs. Never invented.

**Validity:** At least one must be true:

1. Two+ implementations exist
2. Duplicated behavior reduced by extraction
3. RFC explicitly defines abstraction

### Reference First

If reference implementation exists, follow it. Consistency > novelty.

### Single Ownership Rule

Every responsibility has exactly one owner.

### Boundary Rule

- Kernel never knows transport
- Kernel never knows storage implementation
- Storage never owns business logic
- Transport never owns business logic
- Composition roots own wiring only

### Simplicity Rule

Prefer: fewer files, abstractions, interfaces, services, layers, indirection.

**Architecture should become simpler after every implementation.**

### Minimize Architectural Debt

- Prefer deleting over adding abstractions
- Prefer composition over inheritance
- Prefer explicit over magic
- Prefer deterministic over convenient

### Challenge Assumptions

Don't blindly implement. If simpler solution satisfies RFCs, explain and recommend. If request contradicts architecture, explain why.

### Leave Repository Healthier

Every change improves: documentation, verification, naming, comments, tests, architecture, or dead-code removal.

---

---

# 🔍 CONTEXT GATHERING

### Use Graphify First

For repository architecture, ownership, dependencies, implementation locations:

- `graphify query` before grep
- `graphify path` for dependency questions
- `graphify explain` for architectural concepts

### When Unsure

- Prefer reading RFC over guessing
- If RFC ambiguous/missing, surface gap (don't fill with assumptions)
- Do not implement when blocked by ambiguity
- Ask for clarification instead of inventing

---

# 📋 BUILD & VERIFICATION RULES

- **Prefer Taskfile commands** over raw go/python commands
- **Use `task test`** for normal validation
- **Use `task build`** before claiming binaries compile
- **Use `task verify:rfc`** after RFC changes
- **Use `task graph:rebuild`** after graph-relevant changes
- **Never commit** `build/` or `graphify-out/`
- **No ad-hoc build scripts** unless Taskfile cannot express operation

---

# ✅ DEFAULT CHECKLIST

Before finishing, verify:

- [ ] All relevant RFCs read
- [ ] Phase 1: Discovery complete
- [ ] If new abstraction: Phase 2 gate passed
- [ ] If new abstraction: Phase 3 review complete
- [ ] Phase 4: Plan produced and followed
- [ ] Phase 5: Vertical slice implemented
- [ ] Phase 6: RFC compliance verified
- [ ] Phase 7: Self-review completed
- [ ] Phase 8: Verification commands executed
- [ ] Phase 9: Repair loop complete
- [ ] Phase 10: Final review passed
- [ ] Phase 11: Final report produced
- [ ] No RFC violations
- [ ] No ADR violations
- [ ] No invariant violations
- [ ] Architecture simpler
- [ ] Ownership correct
- [ ] Boundaries intact

---

# 🚀 DEVELOPER WORKFLOW

1. **Provide Sprint Task** — developer gives feature
2. **Execute Phases 1-11** — copilot autonomous execution
3. **Stop only if blocked** — RFC ambiguity, architectural ambiguity, or architecture violation
4. **Produce Final Report** — always deliver complete report

**Expectation:** Provide task. Copilot executes complete workflow autonomously. No intermediate confirmations unless blocked.

---

# 📚 REFERENCE DOCUMENTS

For detailed guidance, see:

- [docs/copilot/architecture-review.md](docs/copilot/architecture-review.md) — 11-phase architecture review gate
- [docs/copilot/engineering-laws.md](docs/copilot/engineering-laws.md) — Core principles, invariants, laws
- [docs/copilot/implementation.md](docs/copilot/implementation.md) — Planning, reuse, vertical slices
- [docs/copilot/mapper.md](docs/copilot/mapper.md) — Pure mapper verification
- [docs/copilot/validation.md](docs/copilot/validation.md) — Tests, build, verification
- [rfcs/](rfcs/) — Architectural source of truth
- [docs/adr/](docs/adr/) — Architecture Decision Records
- [docs/architecture/GOLDEN_MAPPER.md](docs/architecture/GOLDEN_MAPPER.md) — Reference mapper
- [docs/architecture.md](docs/architecture.md) — Overall architecture
