# Engineering Learning

Status: ACTIVE

---

# Purpose

Engineering Learning defines how Praxis continuously improves its own engineering system.

Every completed task is an opportunity to improve not only the product, but also the process that builds the product.

Engineering Learning is mandatory.

No task is considered fully complete until Engineering Learning has been evaluated.

---

# Principle

Every implementation produces two outputs:

1. Product Improvement

2. Engineering Improvement

Both are first-class engineering artifacts.

---

# Learning Cycle

Every completed task enters the Engineering Learning stage.

```
Completed Task

↓

Engineering Reflection

↓

Knowledge Extraction

↓

Improvement Candidates

↓

Review

↓

Accepted Improvements

↓

Engineering System Updated
```

---

# Goals

Engineering Learning exists to:

- reduce future complexity
- eliminate repeated mistakes
- capture reusable knowledge
- improve engineering consistency
- reduce review effort
- improve AI performance
- improve onboarding

---

# Engineering Assets

Engineering Learning may update:

- Instructions
- Prompts
- Workflow
- Quality Gates
- Policies
- Reference Implementations
- Golden Examples
- ADRs
- RFCs
- Architecture Guides

No other artifacts.

---

# Mandatory Reflection

Every completed task must answer:

## Instructions

Should engineering instructions improve?

Why?

---

## Prompt

Should prompts improve?

Which prompt?

Why?

---

## Workflow

Should workflow change?

Which stage?

Why?

---

## Quality Gates

Was a missing Quality Gate discovered?

Should a new gate exist?

---

## Reference Implementation

Does this implementation become:

- candidate
- replacement
- no change

Explain.

---

## Golden Example

Does a new Golden Example exist?

If yes:

why?

---

## Policy

Did operational guidance emerge?

Should a Policy be updated?

---

## ADR

Was an architectural decision missing?

Should an ADR change?

---

## RFC

Was domain knowledge missing?

Should an RFC change?

---

# Knowledge Extraction

Extract reusable knowledge.

Examples

- repeated review comments

- recurring mistakes

- successful implementation patterns

- successful architecture patterns

- anti-patterns

- common pitfalls

Knowledge must be concrete.

---

# Improvement Candidates

Every suggestion belongs to exactly one category.

Allowed categories:

Instruction

Prompt

Workflow

Quality Gate

Policy

Reference Implementation

Golden Example

ADR

RFC

Anything else is rejected.

---

# Promotion Rules

## Reference Implementation

Candidate only if:

- production quality
- architecturally reviewed
- simpler than existing reference
- follows all RFCs

---

## Golden Example

Candidate only if:

- minimal
- obvious
- single responsibility
- educational

---

## Policy

Candidate only if:

- operational guidance repeated
- architecture already decided
- future developers benefit

---

## ADR

Candidate only if:

- architectural decision required
- alternatives considered
- decision affects multiple components

---

## RFC

Candidate only if:

- domain model changes
- business semantics change
- architectural specification insufficient

---

# Anti-Learning

Engineering Learning must reject:

- personal preferences

- coding style

- temporary hacks

- speculative ideas

- "might be useful"

- future-proofing

Only demonstrated improvements are accepted.

---

# Learning Report

Every task produces:

```
Engineering Learning

Instruction Updates

Prompt Updates

Workflow Updates

Quality Gates

Reference Candidates

Golden Examples

Policy Candidates

ADR Candidates

RFC Candidates

Rejected Ideas
```

---

# Acceptance Criteria

Engineering Learning succeeds when:

Every suggestion is:

- justified
- categorized
- reviewable
- actionable

---

# Metrics

Track over time:

New Reference Implementations

New Golden Examples

Instruction Updates

Workflow Improvements

Review Time Reduction

Architecture Violations Prevented

Duplicate Implementations Prevented

Engineering metrics are product metrics.

---

# Evolution

Engineering Learning must make the engineering system:

simpler

never

more complicated.

Complexity requires explicit architectural justification.

---

# Final Principle

Engineering knowledge is a product.

Every completed task must leave the engineering system better than it was before.

Praxis improves itself by improving the way Praxis is engineered.