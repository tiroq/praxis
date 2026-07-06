# Architecture Checklist Catalog

Status: ACTIVE

---

# Purpose

Checklists transform architectural knowledge into repeatable engineering decisions.

RFCs define truth.

ADRs define decisions.

Patterns define solutions.

Reference Implementations define examples.

Checklists verify correctness.

---

# Why

Architectural reviews should not depend on memory.

They should depend on explicit verification.

Every checklist represents a gate.

Passing all gates is required before implementation.

---

# Checklist Categories

Current checklists:

Architecture Review

Implementation Readiness

Transport Adapter

Repository

Worker

Storage

Configuration

Service Bootstrap

Observability

Testing

Migration

Documentation

Release

RFC Compliance

ADR Compliance

Pattern Compliance

Reference Compliance

Golden Mapper

Performance

Security

---

# Checklist Template

Every checklist contains:

Purpose

Trigger

Questions

Stop Conditions

Required Evidence

Approval Criteria

---

# Architecture Review Checklist

Trigger

Any architectural change.

Verify

Relevant RFCs reviewed.

Relevant ADRs reviewed.

Reference Implementations searched.

Patterns identified.

Ownership verified.

Responsibilities verified.

Layer boundaries verified.

No duplicated abstractions.

Extract-Don't-Invent satisfied.

Composition Root identified.

Decision documented.

Stop if

Unknown ownership.

Unknown responsibility.

Layer violation.

Duplicate abstraction.

Speculative interface.

Missing RFC.

---

# Transport Adapter Checklist

Trigger

New transport.

Verify

Golden Mapper reviewed.

Pure mapper.

Exactly one translation.

No business logic.

No infrastructure.

Reference Implementation reused.

Output traceable.

Wire contract preserved.

No framework.

Stop if

Mapper grows.

Second translation.

Repository access.

LLM call.

Publish inside mapper.

Retry inside mapper.

---

# Repository Checklist

Trigger

New repository.

Verify

Repository owns persistence only.

No business rules.

No orchestration.

No transport.

Reference searched.

Pattern followed.

---

# Worker Checklist

Trigger

New worker.

Verify

Composition Root only.

Dependencies injected.

No business rules.

Pipeline ownership clear.

---

# Configuration Checklist

Trigger

New environment variable.

Verify

PRAXIS_* naming.

Subsystem owner.

ADR-012 compliant.

Default documented.

Migration considered.

---

# Observability Checklist

Trigger

New service.

Verify

Trace propagation.

Correlation ID.

Structured logging.

Metrics.

Audit events.

---

# Testing Checklist

Trigger

New feature.

Verify

Unit tests.

Integration tests.

Smoke tests.

RFC verification.

Golden path.

Failure path.

---

# Documentation Checklist

Trigger

Architecture change.

Verify

RFC updated.

ADR updated.

Policy updated.

Pattern updated.

Reference updated.

Guide updated.

---

# Reference Compliance Checklist

Trigger

New implementation.

Verify

Reference searched.

Reference reused.

Deviation justified.

Pattern reused.

No new style introduced.

---

# Approval Rule

Every checklist produces one of:

PASS

FAIL

BLOCKED

REQUIRES REVIEW

Nothing else.

---

# Metrics

Track:

Checklist failures

Most violated rules

Missing references

Architecture drift

Average review time

---

# Final Principle

Architecture is enforced through checklists.

Not memory.

Not intuition.

Not experience.

Consistency is measurable.