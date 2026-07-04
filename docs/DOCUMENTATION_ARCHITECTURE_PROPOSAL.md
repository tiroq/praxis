# Documentation Architecture Proposal: ADR-Policy Separation

**Date:** 2026-07-04  
**Status:** PROPOSAL (awaiting approval)  
**Prepared by:** Architecture  
**Related:** ADR-001 through ADR-012, RFC-030, RFC-060

---

## Executive Summary

Current Praxis documentation mixes architectural decisions with operational guidance, creating three problems:

1. **Ambiguous authority:** ADRs contain both immutable decisions and evolving procedures, blurring governance
2. **Difficult updates:** Operational improvements require ADR amendments or create duplicate guidance
3. **Poor navigation:** Developers searching for "how do I?" get decision rationale; architects searching for "what was decided?" get procedure checklists

**Proposal:** Introduce a **Policies layer** (POL-XXX) that organizes operational guidance thematically across ADRs, creating clean separation between immutable decisions and living operational procedures.

**Expected outcome:** Documentation that is easier to navigate, maintain, and evolve without diluting architectural decisions.

---

## 1. CURRENT STATE ASSESSMENT

### 1.1 Four Existing Documentation Layers

| Layer | Purpose | Location | Authority | Update Frequency | Example |
|-------|---------|----------|-----------|------------------|---------|
| **RFC** | Specify domain concepts, contracts, subsystem boundaries | `docs/rfcs/` | Architecture | Rare (design amendments) | RFC-030 (Architecture), RFC-033 (Storage) |
| **ADR** | Record architectural decisions and rationale | `docs/adr/` | Architecture | Rare (new decisions or supersession) | ADR-006 (Transport boundaries) |
| **Architecture Instruction** | Mandatory practices for development | `.github/instructions/` | Architecture | Occasional (governance updates) | praxis-architecture-guardian.instructions.md |
| **Reference Implementation** | Working code examples and patterns | `docs/architecture/` | Architecture + Engineering | Frequent (as practices evolve) | GOLDEN_MAPPER.md |

### 1.2 ADR Content Distribution Analysis

Reviewed all 12 ADRs (ADR-001 through ADR-012):

**Architectural Decision Content (immutable, RFC-backed):**
- What decision was made?
- Why? (options considered, rationale)
- What are the consequences?
- Which RFC(s) justify this?
- What invariants does it establish?

**Operational Content (evolving, practice-based):**
- How to implement the decision?
- Responsibilities and ownership?
- Step-by-step procedures?
- Review checklists and verification?
- Examples and anti-patterns?
- Configuration naming conventions?
- Directory placement rules?
- Build target model and profiles?
- Testing procedures?
- Best practices for compliance?

### 1.3 Operational Content by ADR

| ADR | Title | Op Content % | Primary Op Content |
|-----|-------|--------------|-------------------|
| ADR-001 | Canonical Object vs Artifact | ~5% | Minimal (pure decision) |
| ADR-002 | External Workflow Orchestrator | ~25% | Integration boundaries, responsibility matrix |
| ADR-003 | Internal Event Bus | ~10% | Technology choice (minimal operational content) |
| ADR-004 | Agent Runtime | ~30% | Agent development responsibilities, separation rules |
| ADR-005 | LLM Provider Abstraction | ~35% | Provider adapter pattern, implementation checklist |
| ADR-006 | Transport Boundaries | ~40% | Adapter implementation pattern, responsibility matrix |
| ADR-007 | Storage Architecture | ~45% | Storage ownership matrix, migration governance |
| ADR-008 | Configuration Strategy | ~50% | Precedence rules, layer semantics, examples |
| ADR-009 | Observability | ~50% | Signal model, instrumentation guidelines, ID propagation |
| ADR-010 | Repository Structure | ~55% | Directory responsibilities, placement rules |
| ADR-011 | Build & Development Workflow | ~60% | Target model, profile definitions, developer procedures |
| ADR-012 | Runtime Configuration Naming | ~70% | Subsystems table, setting names rules, compliance checks |

**Average operational content:** ~38% across all ADRs

**Total operational content:** Equivalent to 4-5 full comprehensive policy documents

---

## 2. PROBLEMS WITH CURRENT STATE

### 2.1 Separation of Concerns Violation

**Problem:** ADRs serve dual roles (decision record + operational manual), violating the principle that documents should have a single, clear purpose.

**Impact:** 
- Architects must read 70% operational content in ADR-012 to find the 30% decision
- Developers implementing a feature must navigate across multiple ADRs searching for "how do I?"
- Updates to procedures require ADR changes, diluting the immutability principle

### 2.2 Authority Confusion

**Problem:** ADRs are simultaneously immutable (RFC-backed architectural decisions) and living (evolving operational guidance).

**Impact:**
- Code reviewers unsure: Is this check from an immutable architectural law or an evolving best practice?
- New practices can't be introduced without ADR amendments
- Documentation feels static when procedures need iteration

### 2.3 Navigation Difficulty

**Problem:** Related operational content scattered across 12 separate ADRs; no thematic organization.

**Impact:**
- "How do I add a new transport?" → Read ADR-006 + ADR-010 + check Reference Implementation
- "How do I name a configuration variable?" → Read ADR-008 + ADR-012 + check examples
- Developers waste time cross-referencing instead of finding consolidated guidance

### 2.4 Update Friction

**Problem:** Operational improvements and best practices require ADR amendments or stay undocumented.

**Impact:**
- Better patterns discovered in practice can't be captured as quickly as code can evolve
- Consistency drift as teams create local wikis instead of updating shared documentation
- Learning from real deployments doesn't flow back to documentation

### 2.5 Governance Granularity

**Problem:** Current approach requires either:
- Keep all operational content in ADRs (dilutes decision importance)
- Create one Policy per ADR (mirrors ADR structure, defeats purpose of thematic organization)

**Impact:**
- False choice between two extremes
- No middle ground for thematic organization across multiple ADRs

---

## 3. TARGET DOCUMENTATION ARCHITECTURE

### 3.1 Proposed Hierarchy

```
docs/
├── rfcs/                          ← Domain/subsystem specifications (immutable, RFC-reviewed)
│   ├── 000-vision.md
│   ├── 001-principles.md
│   ├── ...
│   └── 060-verification.md
│
├── adr/                           ← Architectural decisions (immutable, decision-driven)
│   ├── ADR-001-canonical-object-vs-artifact.md
│   ├── ADR-002-external-workflow-orchestrator.md
│   ├── ...
│   └── ADR-012-runtime-configuration-naming.md
│
├── architecture/
│   │
│   ├── policies/                  ← Operational guidance (living, thematic)
│   │   ├── POL-Transport-Adapters.md
│   │   ├── POL-Storage.md
│   │   ├── POL-Runtime-Configuration.md
│   │   ├── POL-Developer-Workflow.md
│   │   ├── POL-Observability.md
│   │   ├── POL-Repository-Layout.md
│   │   ├── POL-Dependency-Management.md
│   │   ├── POL-Testing.md
│   │   └── POL-Provider-Adapters.md
│   │
│   ├── reference/                 ← Reference implementations (code examples)
│   │   ├── GOLDEN_MAPPER.md
│   │   ├── Transport-Adapter-Example.go
│   │   ├── Provider-Adapter-Example.go
│   │   └── ...
│   │
│   └── ARCHITECTURE_README.md     ← Navigation and cross-reference hub
│
├── ARCHITECTURE.md                ← Existing guide (updated to link to new structure)
└── ...
```

### 3.2 Governance Model

| Artifact | Authority | Update Model | Update Frequency | Example |
|----------|-----------|--------------|------------------|---------|
| **RFC** | Domain specifications | RFC amendment (formal review) | Rare | RFC-030 subsystem changes |
| **ADR** | Architectural decision | ADR supersession (formal review) | Rare | New decision → new ADR |
| **Policy** | Operational procedure | Direct update (architecture review) | Frequent | Better pattern discovered |
| **Reference** | Example implementation | Direct update (pull request) | Frequent | New adapter pattern |

**Key distinction:**
- **ADRs are immutable** (can only be superseded, never amended)
- **Policies are living** (amended when better practices emerge)
- **References are illustrative** (updated as patterns evolve)

---

## 4. PROPOSED POLICY CATALOG

### 4.1 POL-Transport-Adapters

**Purpose:** Operational guidance for implementing Praxis transport adapters (HTTP, gRPC, MCP, Telegram, CLI, etc.)

**Related ADRs:** ADR-006 (Transport Boundaries), ADR-010 (Repository Structure)  
**Related RFCs:** RFC-030 (Architecture), RFC-031 (Service Contracts)  
**Related Reference:** Transport-Adapter-Example.go (to be created)

**Content to extract:**
- Hexagonal/ports-and-adapters pattern as operational guidance
- Responsibility matrix: What adapter owns, what Kernel owns, what transport owns
- Transport adapter implementation checklist
- "No business logic in adapters" enforcement rules with examples
- Contract verification at adapter boundary
- Testing strategy for adapters
- Step-by-step: "How to add a new transport"

**Sections to move from ADR-006:**
- (Most of) Responsibilities section
- Responsibility matrix from Options Considered
- Implementation impact details
- Verification impact procedures

**Expected lifetime:** 3-5 years (evolves as transport patterns mature)

**Why separate:** Different from ADR-006's boundary decision; guidance needed for implementing multiple transports. Thematic: teams implementing HTTP, gRPC, or MCP adapters all follow the same pattern.

---

### 4.2 POL-Storage

**Purpose:** Operational guidance for storage ownership, repository patterns, migration procedures, and engine selection decisions

**Related ADRs:** ADR-007 (Storage Architecture), ADR-003 (Internal Event Bus)  
**Related RFCs:** RFC-033 (Storage Model), RFC-031 (Service Contracts)  
**Related Reference:** Repository implementation example (to be created)

**Content to extract:**
- Service-store ownership matrix (which service owns which store category)
- Repository interface pattern and implementation
- Per-store migration governance: versioning, reversibility, replay-safety
- Identity ownership rules per service
- Storage selection decision tree (when to use Postgres, S3, vector store, graph store, Redis, etc.)
- Data corruption recovery procedures
- Backup and restore strategies per store
- Query optimization patterns

**Sections to move from ADR-007:**
- Service-store ownership matrix and responsibility section
- Migration cost section (as migration procedures)
- Database-specific operational impact
- Reference engine mapping table (as deployment reference)

**Sections to reference from ADR-003:**
- Event Store durability procedures

**Expected lifetime:** 2-3 years (evolves as storage engines mature and migration procedures improve)

**Why separate:** Storage is a deep operational domain across multiple services. Teams adding new services, migrating stores, or tuning engines all need consolidated guidance. Not just a rehash of ADR-007; brings in RFC-033 §30 ownership matrix and operational procedures.

---

### 4.3 POL-Runtime-Configuration

**Purpose:** Operational guidance for configuration naming, layer precedence, compliance verification, and deployment patterns

**Related ADRs:** ADR-012 (Runtime Configuration Naming), ADR-008 (Configuration Strategy)  
**Related RFCs:** RFC-030 (Cross-Cutting Concerns)  
**Related Reference:** Configuration examples (docker-compose, Kubernetes, CLI)

**Content to extract from ADR-012:**
- Subsystems reference table (prefixes, examples, owners)
- Setting Names rules (SNAKE_CASE, specificity, `_ENABLED` suffix, `_URL` conventions)
- Namespace Growth Rule (justified subsystems only; forbidden patterns)
- ADR Compliance Rules (7-point checklist for every new variable)
- Review Checklist for Configuration Changes
- Compliance Rules table
- Examples (before/after)
- Migration strategy (Phase 1 parallel support, Phase 2 breaking change)

**Content to extract from ADR-008:**
- Precedence rules table and per-layer semantics
- Layer implementation rules (override semantics, secrets boundary, validation)
- Deployment examples (docker-compose, Kubernetes, CLI)
- Troubleshooting common precedence mistakes
- Secrets handling procedures

**Expected lifetime:** 3-5 years (stabilizes as platform configuration matures)

**Why separate:** Already heavily policy-focused in ADR-012. New services, new subsystems, and deployment questions all converge here. POL becomes the operational handbook; ADR stays as the naming decision.

---

### 4.4 POL-Developer-Workflow

**Purpose:** Operational guidance for build tasks, testing profiles, local development, and CI/local parity

**Related ADRs:** ADR-011 (Build & Development Workflow), ADR-009 (Observability)  
**Related RFCs:** RFC-061 (Verification), RFC-062 (Benchmarking)  
**Related Reference:** Taskfile.yml (reference implementation)

**Content to extract from ADR-011:**
- Target Model table: what each task does, which RFC obligation it fulfills
- Profile definitions: fast (static + schema + contract), standard (+ invariant + security), full (all except benchmark), release (all)
- Developer workflow steps: when to run `task dev` vs `task verify:full` vs `task verify:release`
- CI/local parity rule and enforcement
- Profile-failure troubleshooting
- Pre-commit checklist for developers
- Release readiness verification

**Content to extract from ADR-009:**
- Instrumentation checklist: what every new service must emit
- How-to: "Add observability to a new service"

**Expected lifetime:** 2-3 years (evolves as profile definitions improve and new tasks emerge)

**Why separate:** Developers need this daily; architects need ADR-011 rarely. Thematic: all Praxis development follows one workflow. POL becomes the mandatory reference; ADR stays as the Taskfile decision.

---

### 4.5 POL-Observability

**Purpose:** Operational guidance for instrumentation, trace propagation, log structures, and signal correlation

**Related ADRs:** ADR-009 (Observability), ADR-008 (Configuration Strategy)  
**Related RFCs:** RFC-030 (Cross-Cutting Concerns), RFC-062 (Benchmarking)  
**Related Reference:** Instrumentation example code (to be created)

**Content to extract:**
- Signal Model table: what each signal carries, RFC basis, when to use each
- ID propagation rules: Trace ID, Correlation ID, Causation ID through every layer
- Instrumentation checklist: what every new service must emit
- Structured logging guidelines: slog field conventions, context attachment, log levels
- Query examples: How to find traces by Correlation ID, construct business event queries
- Metrics naming conventions and cardinality guards
- Error tagging and sampling procedures

**Sections to move from ADR-009:**
- Signal Model and responsibilities
- ID propagation section
- Instrumentation procedures
- Structured logging section

**Expected lifetime:** 3-5 years (evolves as observability stack matures)

**Why separate:** Observability spans all services and many RFCs. Teams implementing new services, debugging production issues, and adding observability to new features all converge here. POL becomes the operational manual.

---

### 4.6 POL-Repository-Layout

**Purpose:** Operational guidance for code organization, package placement, import boundaries, and structural consistency

**Related ADRs:** ADR-010 (Repository Structure), ADR-006 (Transport Boundaries)  
**Related RFCs:** RFC-030 (Architecture), RFC-061 (Verification)  
**Related Reference:** Repository structure walkthrough

**Content to extract:**
- Directory Responsibilities table: where to put each type of code
- Placement rules: no business logic in services/*, no transport in internal/core, no frameworks in adapters
- Co-location rules: tests live with code
- Import boundaries: internal/* may not import services/*, services/* may not import each other
- Package naming conventions
- File organization guidelines
- Checklist: "Is my code in the right place?"
- Cross-cutting concern placement (config, observability, shared utilities)

**Sections to move from ADR-010:**
- Directory Responsibilities section
- Invariants section (as enforcement rules)
- Implementation impact details

**Expected lifetime:** 2-3 years (stabilizes as repository patterns settle)

**Why separate:** Developers use this daily for code organization. Thematic across multiple ADRs (transport adapters, storage, configuration all reference placement rules). POL becomes the coding guide.

---

### 4.7 POL-Dependency-Management

**Purpose:** Operational guidance for managing shared libraries, versioning, dependency cycles, and the pkg/ directory

**Related ADRs:** ADR-010 (Repository Structure)  
**Related RFCs:** RFC-030 (Architecture), RFC-061 (Verification)  
**Related Reference:** shared library examples

**Content to extract:**
- When to extract code to pkg/ (proven duplication, no cycles)
- Versioning strategy for shared code
- Cycle detection and breaking strategies
- Export rules: What's public API, what's internal
- Documentation requirements for shared packages
- Breaking change procedures for pkg/
- Test requirements for libraries

**Sections to move from ADR-010:**
- Parts of Responsibilities related to pkg/
- `pkg/` vs `internal/` placement guidance

**Expected lifetime:** 3-5 years (evolves as library patterns mature)

**Why separate:** Shared library management is specialized operational guidance, distinct from repository structure decision. Thematic: all developers who extract or use shared code follow the same rules.

---

### 4.8 POL-Testing

**Purpose:** Operational guidance for testing strategy, profiles, coverage requirements, and test organization

**Related ADRs:** ADR-011 (Build & Development Workflow), ADR-006 (Transport Boundaries)  
**Related RFCs:** RFC-061 (Verification)  
**Related Reference:** Test examples and fixtures

**Content to extract:**
- Testing strategy by component: Kernel, services, adapters, agents
- Profile testing requirements: fast, standard, full
- Coverage expectations and cardinality
- Fixture and mock strategy
- Integration test organization
- Contract testing procedures for services
- Performance testing procedures

**Sections to move from ADR-011:**
- Test-related parts of profile definitions
- Verification impact details

**Expected lifetime:** 3-5 years (evolves as testing patterns improve)

**Why separate:** Testing spans multiple ADRs and is a specialized operational domain. Thematic: all developers follow the same testing strategy. POL becomes the testing handbook.

---

### 4.9 POL-Provider-Adapters

**Purpose:** Operational guidance for implementing LLM provider adapters and extending the routing abstraction

**Related ADRs:** ADR-005 (LLM Provider Abstraction), ADR-004 (Agent Runtime)  
**Related RFCs:** RFC-041 (LLM Routing)  
**Related Reference:** Provider-Adapter-Example.go (to be created)

**Content to extract:**
- Provider adapter implementation pattern (SDKs, retry, fallback, cost tracking)
- "How to add a new LLM provider" step-by-step
- Forbidden patterns: no direct SDK import in business code, no provider branching
- Adapter responsibilities and ownership boundaries
- Fallback strategy and retry procedures
- Cost tracking and quota management per provider
- Testing strategy for new providers
- Performance profiling procedures

**Sections to move from ADR-005:**
- Responsibilities section
- Implementation impact details
- Provider adapter pattern from Options Considered

**Expected lifetime:** 3-5 years (evolves as provider integrations mature)

**Why separate:** Provider adaptation is a specialized operational domain needed by teams adding new LLM models. Thematic: all provider adapter implementations follow the same pattern. POL becomes the extension handbook.

---

### 4.10 POL-Agent-Development

**Purpose:** Operational guidance for implementing agents, managing agent lifecycle, and coordinating with Kernel

**Related ADRs:** ADR-004 (Agent Runtime), ADR-005 (LLM Provider Abstraction)  
**Related RFCs:** RFC-040 (Agent Architecture), RFC-041 (LLM Routing)  
**Related Reference:** Agent implementation example

**Content to extract:**
- Business execution vs agent execution separation and responsibilities
- Context assembly and preparation procedures
- Tool calling patterns and error handling
- Agent testing and debugging procedures
- Cost tracking and token budgeting
- Timeout and cancellation procedures
- Agent state management and recovery

**Sections to move from ADR-004:**
- Responsibilities section
- Separation rules and enforcement
- Implementation procedures

**Expected lifetime:** 2-3 years (evolves as agent patterns mature)

**Why separate:** Agent development is a specialized operational domain. Teams implementing new agents or agent types follow the same procedures. POL becomes the agent development handbook.

---

## 5. ADR EXTRACTION PLAN

### 5.1 ADRs Requiring No Changes (Pure Decisions)

**ADR-001: Canonical Object vs Artifact**
- **Status:** Keep unchanged
- **Rationale:** Pure architectural decision with options analysis; minimal operational content
- **Policy extraction:** None

**ADR-003: Internal Event Bus**
- **Status:** Keep unchanged  
- **Rationale:** Technology selection decision; operational setup (clustering, persistence) belongs in deployment docs, not extracted from this ADR
- **Policy extraction:** None

---

### 5.2 ADRs Requiring Shortening + Policy Extraction

| ADR | Current Focus | After Extraction | Policy Destination | Extraction % |
|-----|---------------|------------------|-------------------|--------------|
| ADR-002 | Orchestrator integration rules | Decision + invariants only | POL-Transport-Adapters (reference) | ~20% |
| ADR-004 | Agent runtime default | Separation rules + pluggability guarantee | POL-Agent-Development | ~25% |
| ADR-005 | LLM provider abstraction | Abstraction boundary | POL-Provider-Adapters | ~30% |
| ADR-006 | Transport boundaries | Hexagonal principle | POL-Transport-Adapters | ~35% |
| ADR-007 | Storage architecture | Repository pattern + engine mapping | POL-Storage | ~40% |
| ADR-008 | Configuration precedence | Layered model principle | POL-Runtime-Configuration + POL-Developer-Workflow | ~45% |
| ADR-009 | Observability architecture | Signal plane decision | POL-Observability + POL-Developer-Workflow | ~50% |
| ADR-010 | Repository structure | Monorepo decision | POL-Repository-Layout + POL-Dependency-Management | ~50% |
| ADR-011 | Build workflow decision | Taskfile decision | POL-Developer-Workflow | ~55% |
| ADR-012 | Configuration naming | Naming convention language | POL-Runtime-Configuration | ~65% |

---

### 5.3 Extraction Scope by ADR

#### ADR-002: External Workflow Orchestrator
**Remove:** Integration boundaries table, responsibility matrix details, orchestrator evaluation checklist  
**Keep:** Kestra decision, six invariants, reasoning for orchestrator-agnosticism  
**Destination:** Principles go to POL-Transport-Adapters; integration checklists reference ADR-002

#### ADR-004: Agent Runtime
**Remove:** Agent development procedures, context assembly guidance, tool calling patterns  
**Keep:** OpenAI ADK decision, pluggability guarantee from RFC-040  
**Destination:** POL-Agent-Development

#### ADR-005: LLM Provider Abstraction
**Remove:** Provider adapter implementation pattern, "how to add a provider" steps, forbidden patterns enforcement  
**Keep:** Custom adapters decision, why not OpenRouter/other routers  
**Destination:** POL-Provider-Adapters

#### ADR-006: Transport Boundaries
**Remove:** Hexagonal pattern explanation as operational guide, adapter implementation checklist, "no business logic" enforcement  
**Keep:** Zero transport knowledge principle, ports & adapters decision  
**Destination:** POL-Transport-Adapters

#### ADR-007: Storage Architecture
**Remove:** Service-store ownership matrix, migration procedures, engine-specific guidance, repository implementation details  
**Keep:** Repository pattern decision, reference engine mapping rationale  
**Destination:** POL-Storage

#### ADR-008: Configuration Strategy
**Remove:** Precedence rules table, per-layer semantics, deployment examples, layer implementation details  
**Keep:** Layered precedence model principle, secrets boundary principle  
**Destination:** POL-Runtime-Configuration (precedence) + POL-Developer-Workflow (layer examples)

#### ADR-009: Observability
**Remove:** Signal Model table, ID propagation procedures, instrumentation checklist, structured logging conventions  
**Keep:** OpenTelemetry + slog + business events decision  
**Destination:** POL-Observability

#### ADR-010: Repository Structure
**Remove:** Directory responsibilities table, placement rules, import boundary enforcement, file organization guidelines  
**Keep:** Standard Go monorepo decision, internal/ boundary principle  
**Destination:** POL-Repository-Layout + POL-Dependency-Management

#### ADR-011: Build & Development Workflow
**Remove:** Target Model table, profile definitions, developer workflow steps, CI/local parity procedures  
**Keep:** Taskfile as single entry point decision  
**Destination:** POL-Developer-Workflow

#### ADR-012: Runtime Configuration Naming
**Remove:** Subsystems table, setting names rules, namespace growth rule, compliance rules, review checklists  
**Keep:** PRAXIS_<SUBSYSTEM>_<SETTING_NAME> naming decision, naming language authority  
**Destination:** POL-Runtime-Configuration

---

## 6. MIGRATION STRATEGY

### 6.1 Phase 1: Create Policies (Non-Breaking)

**Timeline:** 1-2 weeks

**Steps:**
1. Create `docs/architecture/policies/` directory
2. Create 10 POL-XXX files based on content extracted from ADRs
3. Extract content from ADRs into Policies (copy, do not delete from ADRs yet)
4. Update each POL to include "Related ADRs" cross-references
5. Update each ADR to include "See also: POL-XXX" at the top

**Deliverables:**
- 10 new Policy documents with extracted content
- All ADRs updated with Policy cross-references
- `docs/architecture/ARCHITECTURE_README.md` created with navigation guide

**No breaking changes:** All original ADR content remains; Policies are additive.

### 6.2 Phase 2: Rewrite ADRs (Design Only, No Changes Yet)

**Timeline:** 1-2 weeks (design phase)

**Design each ADR rewrite:**
1. Decision only (1-2 pages): What was decided and why?
2. Invariants: Non-negotiable constraints from this decision
3. Related Policies: Links to extracted operational content
4. Open Questions: Unresolved tensions
5. References: Related RFCs and ADRs

**For each ADR, create a design document showing:**
- Which sections will be removed (with destination POL)
- Which sections will remain
- How the ADR will be reshaped for clarity

**Stakeholder review:** Architecture team reviews rewritten ADR designs before execution.

### 6.3 Phase 3: Rewrite ADRs (Implementation)

**Timeline:** 2-3 weeks (after design approval)

**Execute rewrite:**
1. Remove operational content that has been moved to Policies
2. Sharpen decision statements
3. Update cross-references to Policies
4. Verify no content loss (all moved content is in a Policy)

**Validation:** Ensure each rewritten ADR answers:
- What decision was made?
- Why? (options, rationale, consequences)
- What invariants does it establish?
- Which RFCs justify it?

### 6.4 Phase 4: Create Reference Implementations

**Timeline:** Ongoing (3-4 weeks)

**For each Policy, create reference implementation:**
- POL-Transport-Adapters → Transport-Adapter-Example.go (HTTP adapter)
- POL-Storage → Repository-Example.go (repository interface + implementation)
- POL-Provider-Adapters → Provider-Adapter-Example.go (OpenAI adapter example)
- etc.

**Validation:** Each reference should be executable/runnable code; not pseudocode.

### 6.5 Phase 5: Update Documentation Hub

**Timeline:** 1 week (final phase)

**Create `docs/architecture/ARCHITECTURE_README.md`:**
- Navigation by role: Developer, Architect, Reviewer, DevOps
- Quick links to RFCs, ADRs, Policies, References
- Decision flowchart: "How do I find…?"
- Governance model diagram

**Update `docs/ARCHITECTURE.md`:**
- Link to new structure
- Update navigation examples

**Update `.github/instructions/praxis-architecture-guardian.instructions.md`:**
- Add Policy references to relevant sections
- Update architecture review checklist to reference Policies

---

## 7. RISK ANALYSIS

### 7.1 Documentation Fragmentation

**Risk:** Creating 10 Policies across multiple files increases complexity instead of reducing it.

**Severity:** Medium  
**Mitigation:**
- Create comprehensive `ARCHITECTURE_README.md` as navigation hub
- Organize Policies by theme, not by ADR
- Include "See Also" cross-references in every Policy
- Use consistent POL template across all documents
- Validate navigation with user testing (developers finding "how do I?" questions)

**Acceptance Criteria:** New developer can find guidance for any task in ≤2 clicks from ARCHITECTURE_README.

---

### 7.2 Duplicated Authority

**Risk:** Same guidance appearing in both ADR and Policy creates confusion about which is authoritative.

**Severity:** Medium  
**Mitigation:**
- Rewrite ADRs to remove all extracted content (not just reference it)
- Update extraction plan to ensure no content remains in both places
- Document authority hierarchy: "This ADR is authoritative; see POL-XXX for operational guidance"
- Add validation checks in code review: "Is this content duplicated across ADR and Policy?"

**Acceptance Criteria:** Content appears in exactly one place (ADR for decisions, Policy for procedures).

---

### 7.3 Circular References

**Risk:** Policies reference ADRs, ADRs reference Policies, References reference both → difficult to follow.

**Severity:** Low  
**Mitigation:**
- Establish clear reference hierarchy: RFC → ADR → Policy → Reference (one direction only)
- No Policy should reference another Policy except for "See Also" suggestions
- Each document should stand alone but provide context pointers
- Create reference flowchart showing navigation paths

**Acceptance Criteria:** No cycles detected; "See Also" only goes in one direction.

---

### 7.4 Unstable Ownership

**Risk:** Policies become orphaned if ADRs are superseded or new ADRs create conflicting guidance.

**Severity:** Medium  
**Mitigation:**
- Link each Policy to its source ADR(s) at the top
- Establish ownership rules: "When ADR-X is superseded by ADR-Y, update POL-* to reference ADR-Y"
- Create a Policy update checklist: "Are all related ADRs still valid?"
- Archive superseded Policies if their related ADR is superseded
- Establish governance: Architecture team reviews Policy updates

**Acceptance Criteria:** Every Policy has exactly one authoritative source ADR (and vice versa).

---

### 7.5 Excessive Policy Granularity

**Risk:** Too many Policies (10) defeats the purpose; too few misses thematic organization.

**Severity:** Low  
**Mitigation:**
- Propose 10 Policies organized by theme, not by ADR count
- Each Policy should address multiple ADRs when related thematically
- Example: POL-Transport-Adapters covers HTTP, gRPC, Telegram, CLI (not separate policies)
- Later consolidation possible: If two Policies become inseparable, merge them
- Monitor user feedback: Are developers finding what they need?

**Acceptance Criteria:** User testing shows developers find guidance easily; no Policies feel orphaned.

---

### 7.6 Missing Governance

**Risk:** Policies evolve without architecture oversight; diverge from ADR principles.

**Severity:** Medium  
**Mitigation:**
- Establish governance rule: Policy updates require architecture review
- Create Policy update template: "Which ADR principle does this enforce? Is this a new pattern?"
- Add Policy review checklist: "Does this contradict any ADR or RFC?"
- Document when to propose new Policies vs. amending existing ones
- Establish Policy lifetime expectations (3-5 years) and review triggers

**Acceptance Criteria:** Architecture team can audit any Policy and confirm alignment with ADRs/RFCs in ≤30 minutes.

---

### 7.7 Outdated Cross-References

**Risk:** ADRs updated, Policies not updated; cross-references break.

**Severity:** Low  
**Mitigation:**
- Create cross-reference validation script (lint check)
- Add to CI pipeline: Verify all ADR → Policy references are valid
- Document update procedure: "When you update an ADR, audit all related Policies"
- Include in code review checklist: "Are cross-references still valid?"

**Acceptance Criteria:** No broken links; cross-reference consistency verified in CI.

---

### 7.8 Reduced Discoverability

**Risk:** Breaking up operational content across multiple Policies makes it harder to discover related guidance.

**Severity:** Low  
**Mitigation:**
- Create comprehensive ARCHITECTURE_README.md with navigation by role and by task
- Use consistent cross-referencing ("See Also" sections in every Policy)
- Create thematic maps: "I'm working on storage" → shows all related Policies/ADRs/References
- Implement search-friendly structure: Consistent naming, clear headers

**Acceptance Criteria:** Developer can navigate from any entry point to any related document in ≤2 clicks.

---

## 8. RECOMMENDATION

### 8.1 Proceed with ADR-Policy Separation

**Recommendation:** ACCEPT the proposal with the following conditions:

1. **Proceed with Phase 1 (Create Policies):** Non-breaking changes; enables parallel review
2. **Design reviews before Phase 2 (Rewrite ADRs):** Architecture team signs off on each ADR rewrite
3. **Pilot with 3 ADRs:** Start with ADR-012, ADR-010, ADR-008 (highest operational content)
4. **Gather feedback after pilot:** Before completing remaining ADR rewrites, validate user experience
5. **Establish Policy governance:** Document who reviews/updates Policies, how often, and under what conditions

### 8.2 Expected Outcomes

**Immediate benefits (Phase 1-2):**
- Clearer navigation: "How do I?" questions answered in consolidated Policies
- Reduced ADR length: Developers can find architectural decision without reading procedures
- Easier updates: Operational improvements can be captured as Policy amendments, not ADR changes

**Long-term benefits (Phase 3+):**
- Scalable documentation: New services/features add Policy sections without diluting ADRs
- Sustainable governance: Immutable architectural decisions separated from evolving practices
- Better maintainability: Teams can update Policies independently; ADRs remain stable references

### 8.3 Success Criteria

**Documentation quality:**
- [ ] New developer can find any guidance in ≤2 clicks from ARCHITECTURE_README
- [ ] No content duplicated across ADR and Policy
- [ ] All ADR cross-references to Policies are valid
- [ ] All Policy cross-references to ADRs are valid

**Governance clarity:**
- [ ] Architecture team agrees on Policy review/update process
- [ ] No conflicting guidance between ADRs and Policies
- [ ] Every Policy has a clear scope and authority

**Adoption:**
- [ ] Developers actively use Policies for implementation guidance
- [ ] Code reviews reference specific Policies (not just ADRs)
- [ ] New Policies proposed within 6 months for emerging patterns

---

## 9. FINAL DOCUMENTATION HIERARCHY

```
docs/

├── rfcs/                          ← Domain specifications (immutable, RFC-reviewed)
│   ├── 000-vision.md
│   ├── 030-architecture.md
│   ├── 033-storage.md
│   └── ...

├── adr/                           ← Architectural decisions (immutable, decision-driven)
│   ├── ADR-001-canonical-object-vs-artifact.md
│   ├── ADR-006-transport-boundaries.md
│   ├── ADR-007-storage-architecture.md
│   ├── ADR-010-repository-structure.md
│   ├── ADR-012-runtime-configuration-naming.md
│   └── ... (12 ADRs, trimmed to decision core only)

├── architecture/
│   ├── ARCHITECTURE_README.md      ← Navigation hub (discovery starting point)
│   │
│   ├── policies/                  ← Operational guidance (living, thematic)
│   │   ├── POL-Transport-Adapters.md
│   │   ├── POL-Storage.md
│   │   ├── POL-Runtime-Configuration.md
│   │   ├── POL-Developer-Workflow.md
│   │   ├── POL-Observability.md
│   │   ├── POL-Repository-Layout.md
│   │   ├── POL-Dependency-Management.md
│   │   ├── POL-Testing.md
│   │   ├── POL-Provider-Adapters.md
│   │   └── POL-Agent-Development.md
│   │
│   └── reference/                 ← Reference implementations (code examples)
│       ├── GOLDEN_MAPPER.md
│       ├── Transport-Adapter-Example.go
│       ├── Provider-Adapter-Example.go
│       ├── Repository-Example.go
│       └── ...

├── ARCHITECTURE.md                ← Updated to link to new hierarchy

└── ...

.github/

└── instructions/
    └── praxis-architecture-guardian.instructions.md  ← Updated to reference Policies
```

---

## 10. APPROVAL CHECKPOINT

**This proposal is ready for architecture team review.**

**Next steps:**
1. Review and approve (or request modifications) this proposal
2. If approved: Begin Phase 1 (Create Policies)
3. If modifications requested: Return to proposal refinement

**Decision deadline:** [To be set by architecture team]

**Questions for team:**
1. Does the proposed Policy set capture all operational guidance from ADRs?
2. Are the Policy themes correct, or should we reorganize?
3. Should any Policies be combined or split further?
4. What is the timeline for implementation?
5. Who will lead the Phase 2 ADR rewrite review process?

---

**End of Proposal**
