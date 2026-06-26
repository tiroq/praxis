from __future__ import annotations

from dataclasses import asdict


class InMemoryRepository:
    def __init__(self) -> None:
        self._items: dict[str, dict[str, object]] = {}

    def put(self, item_id: str, item: object) -> None:
        self._items[item_id] = asdict(item)

    def list(self) -> list[dict[str, object]]:
        return list(self._items.values())
