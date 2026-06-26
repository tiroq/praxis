from __future__ import annotations

from dataclasses import asdict
from threading import Lock

from fastapi import APIRouter
from pydantic import BaseModel, Field

from packages.praxis_core.models import Event

router = APIRouter(tags=["events"])
# Scaffold-only single-process placeholder until repository-backed persistence exists.
_EVENTS: list[Event] = []
_EVENTS_LOCK = Lock()


class EventCreateRequest(BaseModel):
    kind: str = "note"
    title: str = "Untitled event"
    summary: str = ""
    domain: str = "work"
    source: str = "api"
    payload: dict[str, object] = Field(default_factory=dict)


@router.get("/events")
def list_events() -> list[dict[str, object]]:
    with _EVENTS_LOCK:
        return [asdict(event) for event in _EVENTS]


@router.post("/events")
def create_event(payload: EventCreateRequest) -> dict[str, object]:
    event = Event(
        kind=payload.kind,
        title=payload.title,
        summary=payload.summary,
        domain=payload.domain,
        source=payload.source,
        payload=payload.payload,
    )
    with _EVENTS_LOCK:
        _EVENTS.append(event)
    return asdict(event)
