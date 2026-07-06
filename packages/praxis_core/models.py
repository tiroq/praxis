from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime, timezone
from uuid import uuid4

from packages.praxis_core.enums import ReviewStatus, WorkItemStatus


def utc_now() -> str:
    return datetime.now(timezone.utc).isoformat()


@dataclass(slots=True)
class WorkItem:
    title: str
    description: str = ""
    domain: str = "work"
    status: WorkItemStatus = WorkItemStatus.INBOX
    source: str = "manual"
    tags: list[str] = field(default_factory=list)
    id: str = field(default_factory=lambda: f"wi_{uuid4().hex[:12]}")
    created_at: str = field(default_factory=utc_now)


@dataclass(slots=True)
class Event:
    kind: str
    title: str
    summary: str = ""
    domain: str = "work"
    source: str = "manual"
    payload: dict[str, object] = field(default_factory=dict)
    id: str = field(default_factory=lambda: f"ev_{uuid4().hex[:12]}")
    created_at: str = field(default_factory=utc_now)


@dataclass(slots=True)
class Decision:
    title: str
    rationale: str = ""
    outcome: str = "proposed"
    domain: str = "work"
    id: str = field(default_factory=lambda: f"de_{uuid4().hex[:12]}")
    created_at: str = field(default_factory=utc_now)


@dataclass(slots=True)
class AgentReview:
    agent_name: str
    target_id: str
    summary: str = ""
    status: ReviewStatus = ReviewStatus.PENDING
    id: str = field(default_factory=lambda: f"ar_{uuid4().hex[:12]}")
    created_at: str = field(default_factory=utc_now)


@dataclass(slots=True)
class Artifact:
    kind: str
    uri: str
    title: str = ""
    content_type: str = "text/plain"
    id: str = field(default_factory=lambda: f"af_{uuid4().hex[:12]}")
    created_at: str = field(default_factory=utc_now)


@dataclass(slots=True)
class Relation:
    source_id: str
    target_id: str
    relation_type: str
    id: str = field(default_factory=lambda: f"re_{uuid4().hex[:12]}")
    created_at: str = field(default_factory=utc_now)


@dataclass(slots=True)
class SyncLink:
    internal_id: str
    external_system: str
    external_id: str
    last_synced_at: str | None = None
    id: str = field(default_factory=lambda: f"sl_{uuid4().hex[:12]}")
    created_at: str = field(default_factory=utc_now)
