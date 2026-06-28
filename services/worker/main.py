"""Placeholder worker loop for Praxis."""

from __future__ import annotations

import time


def main() -> None:
    print("Praxis worker scaffold starting")
    print("Status: scaffold only")
    print("TODO: connect to NATS and process queued jobs")
    print("TODO: dispatch work item, review, and artifact handlers")
    for _ in range(1):
        print("Worker heartbeat: no-op loop for scaffold verification")
        time.sleep(0.1)


if __name__ == "__main__":
    main()
