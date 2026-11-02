import json
from collections.abc import Mapping
from pathlib import Path

import pytest

import coolinary_ai
from coolinary_ai import jobs

HAIKU = "claude-haiku-5-5"
SONNET = "claude-sonnet-5-5"
MODELS = frozenset({HAIKU, SONNET})
SHIPPED = Path(coolinary_ai.__file__).with_name("jobs.toml")
HOUR = 3600
DAY = 24 * HOUR
EFFORTS = "low, medium, high, xhigh, max"


def job(name: str = "text_to_items", /, **fields: object) -> str:
    """A valid job table as TOML; `fields` replace the defaults, and a None drops the field."""
    values: dict[str, object] = {
        "model": HAIKU,
        "effort": "low",
        "max_tokens": 4000,
        "timeout_seconds": 15,
        "max_retries": 2,
        "cache_ttl_seconds": HOUR,
        "enabled": True,
        "fallbacks": False,
    } | fields
    lines = [f"[{name}]", *(f"{k} = {json.dumps(v)}" for k, v in values.items() if v is not None)]
    return "\n".join(lines) + "\n"


def plan_override(plan: str, name: str = "text_to_items", /, **fields: object) -> str:
    lines = [f"[{name}.by_plan.{plan}]", *(f"{k} = {json.dumps(v)}" for k, v in fields.items())]
    return "\n".join(lines) + "\n"


def load(tmp_path: Path, *tables: str, env: Mapping[str, str] | None = None) -> jobs.Registry:
    path = tmp_path / "jobs.toml"
    path.write_text("\n".join(tables))
    return jobs.load(path, env or {}, MODELS)


def load_error(tmp_path: Path, *tables: str, env: Mapping[str, str] | None = None) -> str:
    with pytest.raises(jobs.RegistryError) as err:
        load(tmp_path, *tables, env=env)
    return str(err.value)


# Shipped registry and price table


def test_price_table_lists_haiku_and_sonnet_5_5() -> None:
    assert jobs.priced_models() == {HAIKU, SONNET}


def test_shipped_registry_holds_the_starting_values() -> None:
    registry = jobs.load(SHIPPED, {}, jobs.priced_models())
    got = {
        name: (j.model, j.effort, j.max_tokens, j.timeout_seconds, j.cache_ttl_seconds, j.fallbacks)
        for name in registry.names()
        for j in [registry.job(name, "plus")]
    }
    assert got == {
        "text_to_items": (HAIKU, "low", 4000, 15, 7 * DAY, False),
        "photo_to_items": (HAIKU, "medium", 8000, 25, DAY, False),
        "read_label": (HAIKU, "low", 3000, 15, 30 * DAY, False),
        "pick_matches": (HAIKU, "low", 6000, 20, 30 * DAY, False),
        "compose_day_plan": (HAIKU, "medium", 16000, 55, HOUR, False),
        "modify_meal": (SONNET, "low", 8000, 40, HOUR, True),
        "rewrite_dish": (SONNET, "low", 8000, 40, HOUR, True),
        "compose_family_day": (SONNET, "medium", 16000, 55, HOUR, True),
        "compose_from_deals": (SONNET, "medium", 16000, 55, HOUR, True),
        "compose_week_plan": (SONNET, "high", 32000, 110, HOUR, True),
        "cheapest_top_up": (SONNET, "medium", 12000, 55, HOUR, True),
        "extract_flyer_offers": (HAIKU, "medium", 12000, 40, 30 * DAY, False),
        "extract_receipt": (HAIKU, "medium", 6000, 25, DAY, False),
    }


def test_shipped_jobs_are_enabled_with_two_retries() -> None:
    registry = jobs.load(SHIPPED, {}, jobs.priced_models())
    got = {
        (registry.job(n, p).enabled, registry.job(n, p).max_retries)
        for n in registry.names()
        for p in ("free", "plus")
    }
    assert got == {(True, 2)}


def test_env_model_moves_plus_day_plans_while_free_stays_on_haiku() -> None:
    registry = jobs.load(SHIPPED, {"AI_JOB_COMPOSE_DAY_PLAN_MODEL": SONNET}, jobs.priced_models())
    models = [registry.job("compose_day_plan", plan).model for plan in ("plus", "free", "", "gold")]
    assert models == [SONNET, HAIKU, HAIKU, HAIKU]


def test_plus_day_plans_moved_to_sonnet_get_fallbacks() -> None:
    registry = jobs.load(SHIPPED, {"AI_JOB_COMPOSE_DAY_PLAN_MODEL": SONNET}, jobs.priced_models())
    fallbacks = [registry.job("compose_day_plan", plan).fallbacks for plan in ("plus", "free")]
    assert fallbacks == [True, False]


# Plans


def test_plan_override_sets_model_and_effort_for_that_plan_only(tmp_path: Path) -> None:
    registry = load(tmp_path, job(), plan_override("plus", model=SONNET, effort="medium"))
    got = {
        plan: (registry.job("text_to_items", plan).model, registry.job("text_to_items", plan).effort)
        for plan in ("free", "plus")
    }
    assert got == {"free": (HAIKU, "low"), "plus": (SONNET, "medium")}


def test_empty_or_unknown_plan_uses_the_free_override(tmp_path: Path) -> None:
    registry = load(
        tmp_path, job(), plan_override("free", effort="medium"), plan_override("plus", effort="high")
    )
    assert [registry.job("text_to_items", plan).effort for plan in ("", "gold")] == ["medium", "medium"]


def test_unknown_plan_in_by_plan_is_rejected(tmp_path: Path) -> None:
    message = load_error(tmp_path, job(), plan_override("gold", effort="high"))
    assert message == "text_to_items.by_plan.gold: unknown plan, use free or plus"


def test_plan_override_holds_only_model_and_effort(tmp_path: Path) -> None:
    message = load_error(tmp_path, job(), plan_override("plus", max_tokens=8000))
    assert message == "text_to_items.by_plan.plus.max_tokens: unknown field"


def test_plan_override_values_are_checked(tmp_path: Path) -> None:
    message = load_error(tmp_path, job(), plan_override("plus", effort="extreme"))
    assert message == f"text_to_items.by_plan.plus.effort: 'extreme' is not one of {EFFORTS}"


# Env overrides


def test_env_overrides_win_over_the_file_for_every_plan(tmp_path: Path) -> None:
    env = {"AI_JOB_TEXT_TO_ITEMS_MODEL": SONNET, "AI_JOB_TEXT_TO_ITEMS_EFFORT": "high"}
    registry = load(tmp_path, job(), plan_override("plus", model=HAIKU, effort="medium"), env=env)
    got = {
        (registry.job("text_to_items", plan).model, registry.job("text_to_items", plan).effort)
        for plan in ("free", "plus")
    }
    assert got == {(SONNET, "high")}


def test_env_can_switch_a_job_off(tmp_path: Path) -> None:
    registry = load(tmp_path, job(enabled=True), env={"AI_JOB_TEXT_TO_ITEMS_ENABLED": "false"})
    assert [registry.job("text_to_items", plan).enabled for plan in ("free", "plus")] == [False, False]


def test_env_can_switch_a_job_on(tmp_path: Path) -> None:
    registry = load(tmp_path, job(enabled=False), env={"AI_JOB_TEXT_TO_ITEMS_ENABLED": "true"})
    assert registry.job("text_to_items", "free").enabled is True


def test_env_enabled_must_be_true_or_false(tmp_path: Path) -> None:
    message = load_error(tmp_path, job(), env={"AI_JOB_TEXT_TO_ITEMS_ENABLED": "maybe"})
    assert (
        message
        == "text_to_items.enabled: must be true or false, got 'maybe' (from AI_JOB_TEXT_TO_ITEMS_ENABLED)"
    )


def test_env_model_must_be_in_the_price_table(tmp_path: Path) -> None:
    message = load_error(tmp_path, job(), env={"AI_JOB_TEXT_TO_ITEMS_MODEL": "claude-haiku-5-5-20261001"})
    assert message == (
        "text_to_items.model: 'claude-haiku-5-5-20261001' is not in the price table"
        " (from AI_JOB_TEXT_TO_ITEMS_MODEL)"
    )


def test_env_effort_must_be_in_the_allowed_set(tmp_path: Path) -> None:
    message = load_error(tmp_path, job(), env={"AI_JOB_TEXT_TO_ITEMS_EFFORT": "extreme"})
    assert (
        message
        == f"text_to_items.effort: 'extreme' is not one of {EFFORTS} (from AI_JOB_TEXT_TO_ITEMS_EFFORT)"
    )


# Decision 15: Free day plans compose on Haiku 5.5


def test_day_plans_without_the_free_haiku_override_stop_start_up(tmp_path: Path) -> None:
    message = load_error(tmp_path, job("compose_day_plan"))
    assert message == f"compose_day_plan.by_plan.free.model: must be {HAIKU} (decision 15)"


def test_day_plans_with_another_free_model_stop_start_up(tmp_path: Path) -> None:
    message = load_error(
        tmp_path, job("compose_day_plan"), plan_override("free", "compose_day_plan", model=SONNET)
    )
    assert message == f"compose_day_plan.by_plan.free.model: must be {HAIKU} (decision 15)"


# Fallbacks (Q243)


def test_fallbacks_apply_only_when_the_job_runs_on_sonnet(tmp_path: Path) -> None:
    registry = load(tmp_path, job(model=HAIKU, fallbacks=True), plan_override("plus", model=SONNET))
    assert [registry.job("text_to_items", plan).fallbacks for plan in ("free", "plus")] == [False, True]


def test_fallbacks_stay_off_on_sonnet_when_the_file_turns_them_off(tmp_path: Path) -> None:
    registry = load(tmp_path, job(model=SONNET, fallbacks=False))
    assert registry.job("text_to_items", "plus").fallbacks is False


# Validation at start-up


def test_effort_outside_the_allowed_set_stops_start_up(tmp_path: Path) -> None:
    with pytest.raises(jobs.RegistryError) as err:
        load(tmp_path, job(effort="extreme"))
    assert (err.value.job, err.value.field, str(err.value)) == (
        "text_to_items",
        "effort",
        f"text_to_items.effort: 'extreme' is not one of {EFFORTS}",
    )


def test_model_not_in_the_price_table_stops_start_up(tmp_path: Path) -> None:
    message = load_error(tmp_path, job(model="claude-haiku-5-5-20261001"))
    assert message == "text_to_items.model: 'claude-haiku-5-5-20261001' is not in the price table"


@pytest.mark.parametrize(
    ("field", "value"),
    [
        pytest.param("max_tokens", 256, id="max_tokens 256"),
        pytest.param("max_tokens", 64000, id="max_tokens 64000"),
        pytest.param("timeout_seconds", 5, id="timeout 5"),
        pytest.param("timeout_seconds", 120, id="timeout 120"),
        pytest.param("max_retries", 0, id="retries 0"),
        pytest.param("max_retries", 3, id="retries 3"),
        pytest.param("cache_ttl_seconds", 0, id="cache ttl 0"),
    ],
)
def test_values_at_the_limits_are_accepted(tmp_path: Path, field: str, value: int) -> None:
    registry = load(tmp_path, job(**{field: value}))
    assert getattr(registry.job("text_to_items", "free"), field) == value


@pytest.mark.parametrize(
    ("field", "value", "problem"),
    [
        pytest.param("max_tokens", 255, "255 is not between 256 and 64000", id="max_tokens 255"),
        pytest.param("max_tokens", 64001, "64001 is not between 256 and 64000", id="max_tokens 64001"),
        pytest.param("timeout_seconds", 4, "4 is not between 5 and 120", id="timeout 4"),
        pytest.param("timeout_seconds", 121, "121 is not between 5 and 120", id="timeout 121"),
        pytest.param("max_retries", -1, "-1 is not between 0 and 3", id="retries -1"),
        pytest.param("max_retries", 4, "4 is not between 0 and 3", id="retries 4"),
        pytest.param("cache_ttl_seconds", -1, "-1 is under 0", id="cache ttl -1"),
    ],
)
def test_values_past_the_limits_stop_start_up(tmp_path: Path, field: str, value: int, problem: str) -> None:
    assert load_error(tmp_path, job(**{field: value})) == f"text_to_items.{field}: {problem}"


@pytest.mark.parametrize(
    ("field", "value", "problem"),
    [
        pytest.param("max_tokens", "4000", "must be a whole number, got '4000'", id="number as text"),
        pytest.param("max_retries", True, "must be a whole number, got True", id="bool as number"),
        pytest.param("timeout_seconds", 15.5, "must be a whole number, got 15.5", id="fraction"),
        pytest.param("enabled", "yes", "must be true or false, got 'yes'", id="text as bool"),
        pytest.param("model", 5, "must be text, got 5", id="number as model"),
    ],
)
def test_values_of_the_wrong_type_stop_start_up(
    tmp_path: Path, field: str, value: object, problem: str
) -> None:
    assert load_error(tmp_path, job(**{field: value})) == f"text_to_items.{field}: {problem}"


def test_missing_field_stops_start_up(tmp_path: Path) -> None:
    assert load_error(tmp_path, job(timeout_seconds=None)) == "text_to_items.timeout_seconds: is required"


def test_unknown_field_stops_start_up(tmp_path: Path) -> None:
    assert load_error(tmp_path, job(max_token=4000)) == "text_to_items.max_token: unknown field"
