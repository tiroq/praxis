from __future__ import annotations

from dataclasses import asdict

from fastapi import APIRouter, Body

from packages.praxis_core.models import Event

router = APIRouter(tags=["events"])
_EVENTS: list[Event] = []


@router.get("/events")
def list_events() -> list[dict[str, object]]:
    return [asdict(event) for event in _EVENTS]


@router.post("/events")
def create_event(payload: dict[str, object] = Body(...)) -> dict[str, object]:
    event = Event(
        kind=str(payload.get("kind", "note")),
        title=str(payload.get("title", "Untitled event")),
        summary=str(payload.get("summary", "")),
        domain=str(payload.get("domain", "work")),
        source=str(payload.get("source", "api")),
        payload=dict(payload.get("payload") or {}),
    )
    _EVENTS.append(event)
    return asdict(event)
