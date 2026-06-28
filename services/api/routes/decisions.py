from __future__ import annotations

from fastapi import APIRouter

router = APIRouter(prefix="/decisions", tags=["decisions"])


@router.get("")
def list_decisions() -> dict[str, object]:
    return {
        "items": [],
        "status": "scaffold",
        "todo": ["Add decision persistence", "Add review-cycle linkage"],
    }
