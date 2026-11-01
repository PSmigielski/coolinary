"""Settings from environment variables. Any X can be given as X_FILE (Docker secrets)."""

import os
from dataclasses import dataclass
from pathlib import Path


def env(key: str, default: str | None = None) -> str:
    path = os.environ.get(f"{key}_FILE")
    if path:
        return Path(path).read_text().strip()
    value = os.environ.get(key, default)
    if value is None:
        raise RuntimeError(f"{key} is required")
    return value


@dataclass(frozen=True)
class Settings:
    service_name: str
    grpc_port: int
    admin_port: int
    log_level: str
    grpc_reflection: bool
    environment: str
    redis_cache_url: str
    cache_ttl_seconds: int
    llm_provider: str
    langfuse_base_url: str
    langfuse_public_key: str
    langfuse_secret_key: str
    trace_hash_salt: str


def load() -> Settings:
    return Settings(
        service_name=env("SERVICE_NAME", "ai"),
        grpc_port=int(env("GRPC_PORT", "50051")),
        admin_port=int(env("ADMIN_PORT", "8081")),
        log_level=env("LOG_LEVEL", "info"),
        grpc_reflection=env("GRPC_REFLECTION", "false") == "true",
        environment=env("ENVIRONMENT", "dev"),
        redis_cache_url=env("REDIS_CACHE_URL", "redis://redis-cache:6379/0"),
        cache_ttl_seconds=int(env("AI_CACHE_TTL_SECONDS", str(7 * 24 * 3600))),
        llm_provider=env("LLM_PROVIDER", "stub"),
        langfuse_base_url=env("LANGFUSE_BASE_URL", "http://langfuse-web:3000"),
        langfuse_public_key=env("LANGFUSE_PUBLIC_KEY", ""),
        langfuse_secret_key=env("LANGFUSE_SECRET_KEY", ""),
        trace_hash_salt=env("TRACE_HASH_SALT", "coolinary-dev"),
    )


settings = load()
