
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

## Local NATS Smoke Test

Praxis includes a real JetStream smoke path that verifies:

input message
-> NATS JetStream
-> worker transport adapter
-> core kernel
-> output message

Run it only when you want live infrastructure verification:

```sh
task smoke:nats
```

This does not change kernel behavior. The kernel remains transport-agnostic and
the NATS-specific flow stays in `services/worker`, `internal/transport/nats`,
and `cmd/nats-smoke`.

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

Architecture and RFC phase complete. Core Kernel implemented.

Implementation intentionally starts only after architectural foundation is accepted.

---

# Core Kernel Demo

The `kernel-demo` CLI exercises the full `Event → Review → Decision → Action` pipeline
with no external dependencies.

**Requirements:** Go 1.24+, no other setup needed.

```sh
# From repo root
go run ./cmd/kernel-demo "нужно купить билеты в Шанхай"
go run ./cmd/kernel-demo "review this proposal urgently"
go run ./cmd/kernel-demo "nothing relevant here"
```

**Example output:**

```json
{
  "EventID": "demo-1782919033780095",
  "Review": {
    "Reviewer": "keyword-reviewer-v1",
    "Recommendation": "approve",
    "Explanation": "3 keyword(s) matched; review recommends approval with noted risks"
  },
  "Decision": {
    "Outcome": "approve",
    "Reasoning": "confidence 0.75 meets approval threshold 0.60",
    "Policy": "rule-based-policy-v1"
  },
  "Actions": [
    { "Type": "notify", "Priority": "medium", "Description": "notify actor of approval" }
  ]
}
```

Exit code is non-zero on validation errors (empty text, bad confidence) or runtime failures.

**Module structure:**

```
go.work                        ← workspace linking root + kernel modules
go.mod                         ← root module: github.com/tiroq/praxis
cmd/kernel-demo/main.go        ← demo CLI
internal/core/kernel/          ← standalone kernel module (go.mod)
  kernel.go                    ← Kernel.Run: Event → Review → Decision → Action
  default_reviewer.go          ← KeywordReviewer (deterministic, no LLM)
  default_decision_maker.go    ← RuleBasedDecisionMaker
  default_action_planner.go    ← SimpleActionPlanner
  kernel_test.go               ← 29 tests, 97.9% coverage
```

---

# Contributing

Architecture-first.

Every major feature begins with an RFC before implementation.

---

# Developer Commands

```sh
task test             # run all Go tests
task test:coverage    # run tests and write coverage to build/coverage.out
task test:storage     # run storage layer tests (eventstore and sqlite)
task build            # build all binaries into build/
task build:worker     # build services/worker into build/worker
task run:kernel-demo  # build and run kernel demo
task run:api-kernel   # run the HTTP kernel API service
task run:worker       # run the NATS worker (requires NATS server)
task dev              # full verification: test + coverage + build + RFC hygiene
task report           # run tests, write coverage + human-readable summary to build/report.md
task verify:rfc       # verify RFC consistency
task graph:rebuild    # rebuild knowledge graph
task clean            # remove build/
```

---

# License

TBD

---

> **Praxis doesn't manage tasks. Praxis manages decisions.**

> **Turn intent into action.**
