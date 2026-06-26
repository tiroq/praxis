from __future__ import annotations

from dataclasses import asdict, is_dataclass
from typing import Any


class InMemoryRepository:
    def __init__(self) -> None:
        self._items: dict[str, dict[str, object]] = {}

    def put(self, item_id: str, item: Any) -> None:
        if not is_dataclass(item):
            raise TypeError("InMemoryRepository.put expects a dataclass instance")
        self._items[item_id] = asdict(item)

    def list(self) -> list[dict[str, object]]:
        return list(self._items.values())
