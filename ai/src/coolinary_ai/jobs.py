"""Job registry: every AI job's model, effort, limits and on/off switch, from jobs.toml plus env overrides.

Loaded once at start-up. An invalid value raises RegistryError naming the job and field, which stops the
service. Each job is resolved for every plan here, so a call only looks up its job and the caller's plan.
"""

import tomllib
from collections.abc import Collection, Mapping
from dataclasses import dataclass
from pathlib import Path
from typing import Any

PRICES_FILE = Path(__file__).with_name("prices.toml")
EFFORTS = ("low", "medium", "high", "xhigh", "max")
PLANS = ("free", "plus")
# Decision 15 (Q240): Free day plans compose on Haiku 5.5; the file must say so and env can't change it.
PINNED_MODELS = {("compose_day_plan", "free"): "claude-haiku-5-5"}

FIELD_TYPES: dict[str, type] = {
    "model": str,
    "effort": str,
    "max_tokens": int,
    "timeout_seconds": int,
    "max_retries": int,
    "cache_ttl_seconds": int,
    "enabled": bool,
    "fallbacks": bool,
}
PLAN_FIELDS = ("model", "effort")
TYPE_NAMES = {str: "text", int: "a whole number", bool: "true or false"}
LIMITS: dict[str, tuple[int, int | None]] = {
    "max_tokens": (256, 64000),
    "timeout_seconds": (5, 120),
    "max_retries": (0, 3),
    "cache_ttl_seconds": (0, None),
}


class RegistryError(ValueError):
    """An invalid registry value; the message starts with `<job>.<field>:`."""

    def __init__(self, job: str, field: str, problem: str) -> None:
        super().__init__(f"{job}.{field}: {problem}")
        self.job = job
        self.field = field


@dataclass(frozen=True)
class Job:
    """One job's settings for one plan."""

    name: str
    model: str
    effort: str
    max_tokens: int
    timeout_seconds: int
    max_retries: int
    cache_ttl_seconds: int
    enabled: bool
    # Server-side refusal fallbacks; only ever true while the job runs on a Sonnet model (Q243).
    fallbacks: bool


class Registry:
    def __init__(self, jobs: Mapping[str, Mapping[str, Job]]) -> None:
        self._jobs = jobs

    def names(self) -> list[str]:
        return list(self._jobs)

    def job(self, name: str, plan: str) -> Job:
        """The job's settings for the caller's plan; an empty or unknown plan gets Free's."""
        plans = self._jobs[name]
        return plans.get(plan, plans["free"])


def priced_models(path: Path = PRICES_FILE) -> frozenset[str]:
    """Model ids in the price table; the registry accepts no others."""
    with path.open("rb") as f:
        return frozenset(tomllib.load(f))


def load(path: Path, environ: Mapping[str, str], models: Collection[str]) -> Registry:
    with path.open("rb") as f:
        data = tomllib.load(f)
    return Registry({name: _resolve(name, table, environ, models) for name, table in data.items()})


def _resolve(name: str, table: Any, environ: Mapping[str, str], models: Collection[str]) -> dict[str, Job]:
    if not isinstance(table, dict):
        raise RegistryError(name, "*", "must be a table")
    base = _fields(name, table, models)
    by_plan = _plan_overrides(name, table.get("by_plan", {}), models)
    env = _env_overrides(name, environ, models)
    for (job, plan), model in PINNED_MODELS.items():
        if job == name and by_plan.get(plan, {}).get("model") != model:
            raise RegistryError(name, f"by_plan.{plan}.model", f"must be {model} (decision 15)")

    resolved = {}
    for plan in PLANS:
        values = base | by_plan.get(plan, {}) | env
        if (name, plan) in PINNED_MODELS:
            values["model"] = PINNED_MODELS[(name, plan)]
        values["fallbacks"] = values["fallbacks"] and values["model"].startswith("claude-sonnet-")
        resolved[plan] = Job(name=name, **values)
    return resolved


def _fields(name: str, table: dict[str, Any], models: Collection[str]) -> dict[str, Any]:
    for key in table:
        if key not in FIELD_TYPES and key != "by_plan":
            raise RegistryError(name, key, "unknown field")
    values = {}
    for key in FIELD_TYPES:
        if key not in table:
            raise RegistryError(name, key, "is required")
        values[key] = _check(name, key, table[key], models)
    return values


def _plan_overrides(name: str, by_plan: Any, models: Collection[str]) -> dict[str, dict[str, Any]]:
    if not isinstance(by_plan, dict):
        raise RegistryError(name, "by_plan", "must be a table")
    overrides = {}
    for plan, table in by_plan.items():
        if plan not in PLANS:
            raise RegistryError(name, f"by_plan.{plan}", f"unknown plan, use {' or '.join(PLANS)}")
        if not isinstance(table, dict):
            raise RegistryError(name, f"by_plan.{plan}", "must be a table")
        for key, value in table.items():
            if key not in PLAN_FIELDS:
                raise RegistryError(name, f"by_plan.{plan}.{key}", "unknown field")
            _check(name, key, value, models, field=f"by_plan.{plan}.{key}")
        overrides[plan] = dict(table)
    return overrides


def _env_overrides(name: str, environ: Mapping[str, str], models: Collection[str]) -> dict[str, Any]:
    """AI_JOB_<JOB>_MODEL, _EFFORT and _ENABLED; they win over the file for every plan (Q240)."""
    overrides: dict[str, Any] = {}
    for key in ("model", "effort", "enabled"):
        var = f"AI_JOB_{name.upper()}_{key.upper()}"
        raw = environ.get(var, "")
        if not raw:
            continue
        value: Any = raw
        if key == "enabled":
            if raw not in ("true", "false"):
                raise RegistryError(name, key, f"must be true or false, got {raw!r} (from {var})")
            value = raw == "true"
        overrides[key] = _check(name, key, value, models, source=var)
    return overrides


def _check(
    name: str, key: str, value: Any, models: Collection[str], *, field: str = "", source: str = ""
) -> Any:
    field = field or key
    suffix = f" (from {source})" if source else ""
    expected = FIELD_TYPES[key]
    if not _is_type(value, expected):
        raise RegistryError(name, field, f"must be {TYPE_NAMES[expected]}, got {value!r}{suffix}")
    if key == "model" and value not in models:
        raise RegistryError(name, field, f"{value!r} is not in the price table{suffix}")
    if key == "effort" and value not in EFFORTS:
        raise RegistryError(name, field, f"{value!r} is not one of {', '.join(EFFORTS)}{suffix}")
    if key in LIMITS:
        low, high = LIMITS[key]
        if high is None and value < low:
            raise RegistryError(name, field, f"{value} is under {low}")
        if high is not None and not low <= value <= high:
            raise RegistryError(name, field, f"{value} is not between {low} and {high}")
    return value


def _is_type(value: Any, expected: type) -> bool:
    # TOML true is a Python int too; it must not pass as a number.
    return isinstance(value, expected) and not (expected is int and isinstance(value, bool))
