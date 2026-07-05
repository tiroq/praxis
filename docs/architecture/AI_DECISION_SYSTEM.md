# AI Decision System

Status: ACTIVE

Authority

Mandatory decision model for every AI agent operating inside Praxis.

This document defines HOW AI makes engineering decisions.

It does not define architecture.

It defines reasoning.

---

# Purpose

AI must not improvise.

AI must follow a deterministic decision process.

The same problem should produce the same reasoning.

---

# Decision Lifecycle

Every engineering task follows:

Understand

↓

Collect Context

↓

Identify Constraints

↓

Search Existing Knowledge

↓

Generate Alternatives

↓

Evaluate Alternatives

↓

Select Simplest Valid Solution

↓

Verify

↓

Implement

↓

Review

↓

Learn

---

# Step 1

Understand

Questions

What is the actual problem?

Who owns it?

What is the expected outcome?

What is NOT requested?

Stop if unclear.

---

# Step 2

Collect Context

Read

RFCs

ADRs

Policies

Patterns

References

Related code

Related tests

Never reason without context.

---

# Step 3

Identify Constraints

Architecture

Layering

Ownership

Performance

Security

Compatibility

Backward compatibility

Operational impact

Deployment impact

---

# Step 4

Search Existing Knowledge

Search

Reference Registry

Pattern Catalog

Policies

Previous implementations

Existing code

Golden examples

Never solve before searching.

---

# Step 5

Generate Alternatives

Produce at least

Alternative A

Alternative B

Alternative C

if complexity justifies.

Document trade-offs.

---

# Step 6

Evaluate

Evaluate each option by

Correctness

Simplicity

Coupling

Maintainability

RFC compliance

ADR compliance

Future complexity

Reuse

---

# Step 7

Selection Rule

Choose the solution with

lowest complexity

that satisfies

every requirement.

Never choose

the clever solution

over

the simple solution.

---

# Step 8

Verification

Before implementation verify

Architecture

Ownership

Responsibilities

Dependencies

Patterns

References

Golden Examples

Quality Gates

---

# Step 9

Implementation

Implement

minimum vertical slice

reuse existing code

avoid abstractions

keep diff minimal

---

# Step 10

Review

Review own work.

Review architecture.

Review RFC compliance.

Review simplicity.

Review maintainability.

---

# Step 11

Learning

Ask

Should Reference improve?

Should Pattern improve?

Should Policy improve?

Should Prompt improve?

Should Workflow improve?

Should Documentation improve?

If yes

create proposal.

---

# Decision Rules

When uncertain

Search.

Never guess.

When conflict exists

Stop.

Escalate.

When ownership unclear

Stop.

When RFC missing

Stop.

When architecture violated

Stop.

---

# Confidence

Confidence must come from

evidence

never intuition.

Evidence includes

RFC

ADR

Reference

Pattern

Code

Tests

Review

---

# Forbidden Decisions

Never choose

largest abstraction

because it looks future-proof.

Never choose

new technology

without requirement.

Never optimize

without measurement.

Never generalize

before duplication.

Never create

frameworks

for one implementation.

---

# Final Principle

Every engineering decision must be explainable.

If AI cannot explain WHY it selected a solution,

the decision is incomplete.