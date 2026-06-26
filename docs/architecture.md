# Architecture

## Platform overview

Praxis is the platform and source of truth. Chiefly is the primary assistant interface. All other tools act as downstream projections, capture channels, or synchronization targets.

## Runtime components

- `services/api`: API boundary for reading and writing Praxis records
- `services/worker`: background processing for workflows and reviews
- `services/scheduler`: cadence-based triggers for recurring reviews
- `services/llm-router`: logical model-role routing through OmniRoute or compatible providers
- `apps/chiefly`: primary assistant entrypoint
- `apps/telegram`: message-based capture and notification channel
- `apps/admin`: private web UI scaffold
- Postgres, NATS, Ollama, and Caddy as initial supporting runtime services

## Data flow

1. Inputs arrive through Chiefly, Telegram, API calls, or future integrations.
2. Praxis normalizes inputs into internal records such as Work Items, Events, and Decisions.
3. Workers and agents enrich, review, score, or draft outputs.
4. Review cycles gate important actions before external projection.
5. Approved outputs sync to external systems while Praxis remains the source of truth.

## Event-driven model

Praxis is structured around events and workflow transitions instead of a task-only abstraction. Work items, reviews, decisions, and artifacts can emit events that drive classification, review, follow-up, and synchronization.

## LLM gateway model

An OpenAI-compatible gateway such as OmniRoute sits between Praxis and local or hosted models. Praxis assigns logical roles like `json_controller`, `fast_classifier`, `long_reasoner`, and `critic`, then the router resolves those roles to concrete providers and models.

## External integrations as projections

Google Tasks, CRM, Telegram, Upwork, GitHub, and Calendar are treated as projections or integrations. They can feed data into Praxis and receive approved updates back out, but they do not own the canonical state.
