"""Langfuse tracing for every LLM call.

Privacy (docs/infrastructure.md): traces carry the feature, plan tier, a hashed household
id, model, token usage and sizes. Never user text, item names or image bytes; the mask
function below redacts any string that slips through.

Tracing is on only when Langfuse keys are set and the Langfuse host resolves, so
`make dev` (no observability profile) runs without it.
"""

import hashlib
import socket
import time
from collections.abc import Iterator
from contextlib import contextmanager
from typing import Any
from urllib.parse import urlparse

from langfuse import Langfuse, propagate_attributes

from .config import settings
from .log import log

SAFE_KEYS = {"feature", "language", "chars", "items", "cached", "plan", "model"}
_client: Langfuse | None = None
_checked_at = 0.0
_reachable = False


def _mask(*, data: Any, **_: Any) -> Any:
    if isinstance(data, dict):
        return {k: (v if k in SAFE_KEYS else _mask(data=v)) for k, v in data.items()}
    if isinstance(data, list):
        return [_mask(data=v) for v in data]
    if isinstance(data, str):
        return "[redacted]"
    return data


def _available() -> bool:
    global _checked_at, _reachable
    if not (settings.langfuse_public_key and settings.langfuse_secret_key):
        return False
    if time.monotonic() - _checked_at > 60:
        _checked_at = time.monotonic()
        host = urlparse(settings.langfuse_base_url).hostname or ""
        try:
            socket.getaddrinfo(host, None)
            _reachable = True
        except OSError:
            if _reachable or _checked_at == 0:
                log.info("langfuse not reachable; tracing off", extra={"host": host})
            _reachable = False
    return _reachable


def client() -> Langfuse | None:
    global _client
    if not _available():
        return None
    if _client is None:
        _client = Langfuse(
            public_key=settings.langfuse_public_key,
            secret_key=settings.langfuse_secret_key,
            base_url=settings.langfuse_base_url,
            environment=settings.environment,
            mask=_mask,
        )
        log.info("langfuse tracing on", extra={"base_url": settings.langfuse_base_url})
    return _client


def household_hash(household_id: str) -> str | None:
    if not household_id:
        return None
    return hashlib.sha256(f"{settings.trace_hash_salt}:{household_id}".encode()).hexdigest()[:16]


@contextmanager
def generation(
    *, feature: str, plan: str, household_id: str, model: str, input: dict[str, Any]
) -> Iterator[Any]:
    """Yields a Langfuse generation (or None when tracing is off)."""
    lf = client()
    if lf is None:
        yield None
        return
    attrs = propagate_attributes(
        user_id=household_hash(household_id),
        tags=[feature, f"plan:{plan or 'unknown'}"],
        trace_name=feature,
        metadata={"feature": feature, "plan": plan or "unknown"},
    )
    with (
        attrs,
        lf.start_as_current_observation(name=feature, as_type="generation", model=model, input=input) as gen,
    ):
        yield gen


def shutdown() -> None:
    if _client is not None:
        _client.shutdown()
