from __future__ import annotations

from dataclasses import asdict

from fastapi import APIRouter, Body

from packages.praxis_core.enums import WorkItemStatus
from packages.praxis_core.models import WorkItem

router = APIRouter(tags=["work-items"])
_WORK_ITEMS: list[WorkItem] = []


@router.get("/work-items")
def list_work_items() -> list[dict[str, object]]:
    return [asdict(item) for item in _WORK_ITEMS]


@router.post("/work-items")
def create_work_item(payload: dict[str, object] = Body(...)) -> dict[str, object]:
    item = WorkItem(
        title=str(payload.get("title", "Untitled work item")),
        description=str(payload.get("description", "")),
        domain=str(payload.get("domain", "work")),
        status=WorkItemStatus(str(payload.get("status", WorkItemStatus.INBOX.value))),
        source=str(payload.get("source", "api")),
        tags=list(payload.get("tags") or []),
    )
    _WORK_ITEMS.append(item)
    return asdict(item)
