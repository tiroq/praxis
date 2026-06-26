from __future__ import annotations

from dataclasses import asdict
from threading import Lock

from fastapi import APIRouter
from pydantic import BaseModel, Field

from packages.praxis_core.enums import WorkItemStatus
from packages.praxis_core.models import WorkItem

router = APIRouter(tags=["work-items"])
# Scaffold-only local store; not safe for multi-worker or multi-process deployment.
_WORK_ITEMS: list[WorkItem] = []
_WORK_ITEMS_LOCK = Lock()


class WorkItemCreateRequest(BaseModel):
    title: str = "Untitled work item"
    description: str = ""
    domain: str = "work"
    status: WorkItemStatus = WorkItemStatus.INBOX
    source: str = "api"
    tags: list[str] = Field(default_factory=list)


@router.get("/work-items")
def list_work_items() -> list[dict[str, object]]:
    with _WORK_ITEMS_LOCK:
        items = list(_WORK_ITEMS)
    return [asdict(item) for item in items]


@router.post("/work-items")
def create_work_item(payload: WorkItemCreateRequest) -> dict[str, object]:
    item = WorkItem(
        title=payload.title,
        description=payload.description,
        domain=payload.domain,
        status=payload.status,
        source=payload.source,
        tags=payload.tags,
    )
    with _WORK_ITEMS_LOCK:
        _WORK_ITEMS.append(item)
    return asdict(item)
