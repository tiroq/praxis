# Praxis

Praxis is a self-hosted personal AI operating platform for turning intent into action.

Chiefly is the primary assistant and chief-of-staff experience inside Praxis. Praxis captures work items, events, decisions, reviews, artifacts, and relationship context, then routes that information through agents, review flows, and external integrations without giving up ownership of the source of truth.

## Why Praxis exists

Personal operating systems usually fragment across notes, inboxes, task apps, CRMs, chat threads, and freelance tools. Praxis exists to unify that context in one local-first platform that can coordinate action across life domains while remaining self-hosted and extensible.

## High-level architecture

- **Praxis** is the platform and system of record.
- **Chiefly** is the main assistant interface inside Praxis.
- **LLM routing** is centralized through an OpenAI-compatible gateway such as OmniRoute.
- **External systems** like Google Tasks, Telegram, Upwork, GitHub, CRM, and calendars are projections or integrations, not the source of truth.
- **Core entities** include Work Items, Events, Decisions, Agent Reviews, Artifacts, Relations, Review Cycles, and Sync Links.
- **Runtime target** is a home PC or mini PC using Docker Compose.

## Repository layout

- `apps/` user-facing and channel-specific entrypoints
- `services/` runtime services such as API, worker, scheduler, and LLM router
- `packages/` shared domain models, agents, connectors, and workflows
- `domains/` domain-specific prompts, schemas, and workflow notes
- `agents/` prompt and policy scaffolds for future agents
- `configs/` life areas, projects, goals, review cycles, integrations, and routing config
- `infra/` local infrastructure bootstrapping files
- `scripts/verify/` executable scaffold verification placeholders
- `docs/` architectural and planning documents

## Local development

```bash
cp .env.example .env
make verify
docker compose up --build
```

Useful commands:

```bash
make run-api
make run-worker
make run-chiefly
make run-telegram
```

## Current status

This repository is **scaffold only**. It includes a monorepo structure, placeholder services, basic models, configuration files, and verification scripts, but it does not yet implement full business logic or production workflows.
