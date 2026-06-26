from __future__ import annotations

from fastapi import APIRouter

router = APIRouter(prefix="/agents", tags=["agents"])


@router.get("")
def list_agents() -> dict[str, object]:
    return {
        "items": [
            "life_planner",
            "critic",
            "progress_watcher",
            "opportunity_scout",
            "finance_reviewer",
            "proposal_writer",
        ],
        "status": "scaffold",
    }
