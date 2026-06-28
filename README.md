
# Praxis

> **Turn intent into action.**

> **Status:** Foundation (Architecture Phase)

## Overview

Praxis is an AI Work Operating System.

Unlike traditional productivity tools that focus on storing tasks, Praxis focuses on transforming intent into action through structured reasoning, AI collaboration, human review and execution.

The platform acts as an AI Chief of Staff.

---

# The Problem

Knowledge is fragmented across dozens of disconnected systems:

- Telegram
- Email
- Calendar
- GitHub
- Google Tasks
- Upwork
- Documents
- Notes
- CRM

Existing software organizes information.

Praxis orchestrates decisions.

---

# Vision

Every incoming event should become:

```text
Event
 ↓
Understanding
 ↓
Knowledge
 ↓
AI Reviews
 ↓
Human Decision
 ↓
Execution
 ↓
Learning
```

---

# Mission

Help people spend less time organizing work and more time doing meaningful work.

---

# Design Principles

- Human remains in control.
- AI proposes. Humans decide.
- Everything starts as an Event.
- Everything becomes an Artifact.
- Every important action is explainable.
- Local-first whenever practical.
- Provider agnostic.
- Event-driven.
- Composable.
- Observable.
- Versioned.

---

# Core Concepts

## Event

Any incoming information.

Examples:

- Telegram message
- Voice note
- Upwork job
- GitHub issue
- Email
- Calendar event

## Artifact

Normalized object inside Praxis.

Artifact types:

- Task
- Lead
- Proposal
- Client
- Project
- Meeting
- Research
- Decision
- Document
- Contract
- Reminder

## Review

Opinion generated independently by:

- Planner
- Critic
- Architect
- Risk
- Portfolio
- Human

## Decision

Possible outcomes:

- Approve
- Reject
- Edit
- Split
- Merge
- Archive
- Escalate

---

# Capability Map

```text
Capture
  ↓
Understand
  ↓
Enrich
  ↓
Think
  ↓
Review
  ↓
Decide
  ↓
Execute
  ↓
Learn
```

---

# User Journey

Telegram

↓

Raw Event

↓

Normalization

↓

Classification

↓

Project Routing

↓

Priority Estimation

↓

Planner Review

↓

Critic Review

↓

Human Approval

↓

Execution

↓

History & Learning

---

# Domains

- Personal
- Work
- Freelance
- Products
- Content

Future:

- Finance
- Health
- Knowledge
- Travel

---

# Freelance Domain

Pipeline:

```text
New Lead
 ↓
Scoring
 ↓
Risk Review
 ↓
Proposal Draft
 ↓
Human Review
 ↓
Applied
 ↓
Interview
 ↓
Won
 ↓
Project
```

---

# Agent Model

Planner
: Suggests execution strategy.

Critic
: Challenges assumptions.

Architect
: Reviews technical approach.

Risk
: Evaluates uncertainty.

Portfolio
: Suggests relevant experience.

Consensus
: Aggregates reviews.

Human
: Makes final decision.

---

# High-Level Architecture

```text
Sources
   ↓
Praxis Core
   ↓
Workflow Engine
   ↓
Agent Runtime
   ↓
Reviews
   ↓
Decision
   ↓
Execution
   ↓
External Systems
```

---

# Planned Integrations

- Telegram
- Google Tasks
- Google Calendar
- Gmail
- GitHub
- Upwork
- OmniRoute
- Ollama
- OpenRouter

---

# Technology

- Go
- Python
- React
- PostgreSQL
- pgvector
- NATS JetStream
- Docker
- OmniRoute

---

# Repository

```text
praxis/
├── README.md
├── MANIFESTO.md
├── ROADMAP.md
├── rfcs/
├── apps/
├── services/
├── packages/
├── domains/
├── agents/
├── workflows/
├── integrations/
├── infrastructure/
├── configs/
└── scripts/
```

---

# Why Praxis?

| Existing Tool | Primary Goal | Praxis Adds |
|---------------|--------------|-------------|
| Notion | Documentation | AI reasoning |
| Linear | Project tracking | Multi-domain orchestration |
| GitHub Projects | Engineering | Personal + Freelance + Work |
| Google Tasks | Todos | Decision lifecycle |
| ChatGPT | Conversation | Persistent workflows |
| n8n | Automation | AI-native orchestration |

---

# Roadmap

1. Foundation
2. Core Runtime
3. Agent Workflows
4. Knowledge Graph
5. Automation
6. Multi-user

---

# Current Status

Architecture and RFC phase.

Implementation intentionally starts only after architectural foundation is accepted.

---

# Contributing

Architecture-first.

Every major feature begins with an RFC before implementation.

---

# License

TBD

---

> **Praxis doesn't manage tasks. Praxis manages decisions.**

> **Turn intent into action.**
