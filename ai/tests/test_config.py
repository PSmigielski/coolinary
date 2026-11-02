import logging
from pathlib import Path

import pytest

import coolinary_ai
from coolinary_ai import config
from coolinary_ai.log import JsonFormatter

FAKE_KEY = "sk-ant-test-0000-not-a-real-key"


@pytest.fixture(autouse=True)
def clean_env(monkeypatch: pytest.MonkeyPatch) -> None:
    for key in ("LLM_PROVIDER", "ANTHROPIC_API_KEY", "AI_JOBS_FILE", "AI_USD_TO_EUR"):
        monkeypatch.delenv(key, raising=False)
        monkeypatch.delenv(f"{key}_FILE", raising=False)


def test_stub_is_the_default_provider_and_needs_no_api_key() -> None:
    settings = config.load()
    assert (settings.llm_provider, settings.anthropic_api_key) == ("stub", "")


def test_unknown_llm_provider_refuses_to_start(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv("LLM_PROVIDER", "openai")
    with pytest.raises(RuntimeError, match=r"^LLM_PROVIDER must be anthropic or stub, got 'openai'$"):
        config.load()


def test_anthropic_provider_without_api_key_refuses_to_start(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv("LLM_PROVIDER", "anthropic")
    with pytest.raises(RuntimeError, match=r"^ANTHROPIC_API_KEY is required$"):
        config.load()


def test_anthropic_provider_with_empty_api_key_refuses_to_start(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv("LLM_PROVIDER", "anthropic")
    monkeypatch.setenv("ANTHROPIC_API_KEY", "")
    with pytest.raises(RuntimeError, match=r"^ANTHROPIC_API_KEY is required$"):
        config.load()


def test_api_key_is_read_from_a_secret_file(monkeypatch: pytest.MonkeyPatch, tmp_path: Path) -> None:
    secret = tmp_path / "anthropic_api_key"
    secret.write_text(FAKE_KEY + "\n")
    monkeypatch.setenv("LLM_PROVIDER", "anthropic")
    monkeypatch.setenv("ANTHROPIC_API_KEY_FILE", str(secret))
    assert config.load().anthropic_api_key == FAKE_KEY


def test_api_key_never_shows_in_the_settings_text(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv("LLM_PROVIDER", "anthropic")
    monkeypatch.setenv("ANTHROPIC_API_KEY", FAKE_KEY)
    settings = config.load()
    assert FAKE_KEY not in repr(settings)
    assert FAKE_KEY not in str(settings)


def test_api_key_never_shows_in_a_log_line(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv("LLM_PROVIDER", "anthropic")
    monkeypatch.setenv("ANTHROPIC_API_KEY", FAKE_KEY)
    record = logging.makeLogRecord({"msg": "settings", "settings": config.load()})
    assert FAKE_KEY not in JsonFormatter().format(record)


def test_jobs_file_defaults_to_the_registry_shipped_with_the_code() -> None:
    assert Path(config.load().ai_jobs_file) == Path(coolinary_ai.__file__).parent / "jobs.toml"


def test_usd_to_eur_defaults_to_0_92() -> None:
    assert config.load().ai_usd_to_eur == 0.92


def test_usd_to_eur_can_be_set(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv("AI_USD_TO_EUR", "0.9")
    assert config.load().ai_usd_to_eur == 0.9
