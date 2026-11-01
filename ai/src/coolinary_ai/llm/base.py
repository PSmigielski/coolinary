"""Provider-agnostic LLM interface. Each provider turns a request into structured results;
callers never see provider-specific types."""

from dataclasses import dataclass, field
from enum import IntEnum
from typing import Protocol


class Unit(IntEnum):
    """Mirrors coolinary.common.v1.Unit."""

    UNSPECIFIED = 0
    GRAM = 1
    MILLILITRE = 2
    PIECE = 3


@dataclass(frozen=True)
class Item:
    name: str
    quantity: float | None
    unit: Unit
    confidence: float


@dataclass(frozen=True)
class ExtractResult:
    items: list[Item]
    model: str
    # Token counts for cost tracking in Langfuse, e.g. {"input": 120, "output": 40}.
    usage: dict[str, int] = field(default_factory=dict)


class Provider(Protocol):
    name: str

    async def extract_items(self, text: str, language: str) -> ExtractResult: ...
