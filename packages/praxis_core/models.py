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


@dataclass(slots=True)
class Conversation:
    """Canonical Conversation object (RFC-014 Canonical Object, RFC-033 Canonical Store).
    
    A Conversation represents a persistent dialogue between participants,
    identified by a correlation_id that groups related messages from an external
    system (e.g. Telegram chat_id). Conversation metadata is mutable; messages
    are append-only.
    
    Fields:
      id: Globally unique, immutable conversation identifier (conv_*).
      correlation_id: External grouping identifier (e.g. telegram-chat-{chat_id}).
      lifecycle: Conversation state (created, active, archived).
      created_at: ISO8601 timestamp of conversation creation.
      updated_at: ISO8601 timestamp of last metadata update.
      last_message_at: ISO8601 timestamp of last message in conversation (or None).
    """
    
    correlation_id: str
    lifecycle: str = "created"  # created, active, archived
    id: str = field(default_factory=lambda: f"conv_{uuid4().hex[:12]}")
    created_at: str = field(default_factory=utc_now)
    updated_at: str = field(default_factory=utc_now)
    last_message_at: str | None = None


@dataclass(slots=True)
class Message:
    """Message in a Conversation (RFC-013 Event reference, RFC-033 Projection Store).
    
    A Message represents a single message in a conversation. It references the
    immutable Event ID (source of truth) and contains extracted fields for efficient
    querying and display. Messages are append-only.
    
    Fields:
      id: Globally unique, immutable message identifier (msg_*).
      conversation_id: Foreign key to Conversation.
      event_id: Reference to immutable event in Event Store (never changes).
      role: "user" or "assistant" (role of message author).
      content: Message text content.
      timestamp: ISO8601 timestamp of message.
      metadata: Optional enrichment fields (e.g. username, message_id, chat_id).
    """
    
    conversation_id: str
    event_id: str
    role: str  # "user" or "assistant"
    content: str
    timestamp: str
    metadata: dict[str, str] = field(default_factory=dict)
    id: str = field(default_factory=lambda: f"msg_{uuid4().hex[:12]}")
    created_at: str = field(default_factory=utc_now)
