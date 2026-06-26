from __future__ import annotations

from enum import Enum


class WorkItemStatus(str, Enum):
    INBOX = "inbox"
    READY = "ready"
    IN_PROGRESS = "in_progress"
    BLOCKED = "blocked"
    DONE = "done"


class ReviewStatus(str, Enum):
    PENDING = "pending"
    APPROVED = "approved"
    REVISE = "revise"
    REJECTED = "rejected"
