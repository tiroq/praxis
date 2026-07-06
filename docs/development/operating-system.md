# Praxis Engineering Operating System

Status: ACTIVE

---

# Purpose

Praxis Engineering Operating System defines how engineering work is performed.

It is the operating model for humans and AI.

Everything else exists inside this system.

---

# Layers

Engineering operates through layers.

Each layer has one responsibility.

RFC

↓

ADR

↓

Policy

↓

Pattern

↓

Reference Implementation

↓

Checklist

↓

Playbook

↓

Prompt

↓

Workflow

↓

Quality Gate

↓

Report

---

# Responsibilities

RFC

Defines domain truth.

Never implementation.

---

ADR

Defines architecture.

Never operational guidance.

---

Policy

Defines engineering rules.

May evolve.

---

Pattern

Defines reusable architectural solutions.

---

Reference

Defines approved implementations.

---

Checklist

Defines verification.

---

Playbook

Defines execution.

---

Prompt

Defines AI behavior.

---

Workflow

Defines orchestration.

---

Quality Gate

Defines release criteria.

---

Report

Defines evidence.

---

# Engineering Lifecycle

Idea

↓

Architecture

↓

Review

↓

Implementation

↓

Verification

↓

Self Review

↓

Quality Gates

↓

Merge

↓

Release

↓

Observe

↓

Learn

↓

Improve References

↓

Improve Patterns

↓

Improve Policies

↓

Repeat

---

# Human Responsibilities

Humans own

architecture

tradeoffs

priorities

RFC approval

ADR approval

risk

product

vision

AI never owns these.

---

# AI Responsibilities

AI owns

analysis

implementation

verification

documentation

testing

review

comparison

consistency

AI proposes.

Humans approve.

---

# Decision Hierarchy

When documents disagree

RFC wins.

If RFC absent

ADR wins.

If ADR absent

Policy wins.

If Policy absent

Pattern wins.

If Pattern absent

Reference wins.

If nothing exists

Architecture Review is required.

---

# Learning Loop

Every implementation should improve:

Reference

Pattern

Checklist

Playbook

Prompt

Workflow

Documentation

The system becomes better after every change.

---

# Metrics

Track

Architecture violations

Review findings

RFC compliance

Pattern reuse

Reference reuse

Average implementation time

Average review time

Defect escape rate

Golden Mapper reuse

Policy violations

Prompt improvements

Workflow improvements

---

# Anti-Patterns

Never optimize implementation before optimizing engineering.

Never automate broken processes.

Never create documents without ownership.

Never create abstractions without evidence.

Never replace architecture with prompts.

Never replace engineering with AI.

---

# Evolution

Engineering itself evolves.

When recurring work appears:

implementation

↓

reference

↓

pattern

↓

policy

↓

automation

↓

workflow

↓

platform capability

Engineering knowledge becomes engineering infrastructure.

---

# Final Principle

Praxis is not only a software platform.

Praxis is an engineering platform.

The product improves.

The engineering system improves.

Both evolve continuously.