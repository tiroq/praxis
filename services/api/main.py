"""Minimal FastAPI scaffold for the Praxis API.

This app is scaffold-only and intended for local/private development.
Authentication, authorization, and production-grade deployment hardening
remain TODO items before any wider network exposure.
"""

from __future__ import annotations

from fastapi import FastAPI

from services.api.routes.agents import router as agents_router
from services.api.routes.decisions import router as decisions_router
from services.api.routes.events import router as events_router
from services.api.routes.health import router as health_router
from services.api.routes.work_items import router as work_items_router

app = FastAPI(title="Praxis API", version="0.1.0", description="Scaffold-only Praxis API")
app.include_router(health_router)
app.include_router(work_items_router)
app.include_router(events_router)
app.include_router(decisions_router)
app.include_router(agents_router)


if __name__ == "__main__":
    import uvicorn

    uvicorn.run("services.api.main:app", host="0.0.0.0", port=8080, reload=False)
