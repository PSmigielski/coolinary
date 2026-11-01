"""Pick the LLM provider from LLM_PROVIDER. Only the stub exists until a provider is chosen."""

from .base import ExtractResult, Item, Provider, Unit
from .stub import StubProvider

__all__ = ["ExtractResult", "Item", "Provider", "Unit", "get_provider"]


def get_provider(name: str) -> Provider:
    if name == "stub":
        return StubProvider()
    raise ValueError(f"unknown LLM_PROVIDER {name!r}; add it in coolinary_ai/llm/")
