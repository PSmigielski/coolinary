"""AI result cache in redis-cache: the same request never costs twice."""

import contextlib
import hashlib
import json
from typing import Any

import redis.asyncio as redis

from .config import settings

client = redis.from_url(settings.redis_cache_url, socket_timeout=0.5, socket_connect_timeout=0.5)


def key(feature: str, *parts: str) -> str:
    digest = hashlib.sha256("\x1f".join(parts).encode()).hexdigest()[:32]
    return f"ai:{feature}:{digest}"


async def get_json(k: str) -> Any | None:
    try:
        raw = await client.get(k)
    except redis.RedisError:
        return None  # cache is best-effort
    return json.loads(raw) if raw else None


async def set_json(k: str, value: Any, ttl: int = settings.cache_ttl_seconds) -> None:
    with contextlib.suppress(redis.RedisError):  # cache is best-effort
        await client.set(k, json.dumps(value), ex=ttl)
