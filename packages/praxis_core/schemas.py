from __future__ import annotations

from packages.praxis_core.models import AgentReview, Artifact, Decision, Event, Relation, SyncLink, WorkItem

CORE_SCHEMAS = {
    "WorkItem": WorkItem,
    "Event": Event,
    "Decision": Decision,
    "AgentReview": AgentReview,
    "Artifact": Artifact,
    "Relation": Relation,
    "SyncLink": SyncLink,
}
