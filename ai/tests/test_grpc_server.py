from collections.abc import Sequence
from pathlib import Path
from typing import Any, cast

import grpc
import pytest

import coolinary_ai
from coolinary.ai.v1 import ai_pb2
from coolinary_ai import jobs
from coolinary_ai.grpc_server import AiServicer
from coolinary_ai.llm import ExtractResult

SHIPPED = Path(coolinary_ai.__file__).with_name("jobs.toml")


class Aborted(Exception):
    pass


class FakeContext:
    """The parts of grpc.aio.ServicerContext the handlers use; abort records the status and raises."""

    def __init__(self, metadata: Sequence[tuple[str, str]] = ()) -> None:
        self.metadata = metadata
        self.code: grpc.StatusCode | None = None
        self.details = ""
        self.trailing: tuple[tuple[str, str], ...] = ()

    def invocation_metadata(self) -> Sequence[tuple[str, str]]:
        return self.metadata

    async def abort(
        self, code: grpc.StatusCode, details: str = "", trailing_metadata: Sequence[tuple[str, str]] = ()
    ) -> None:
        self.code, self.details, self.trailing = code, details, tuple(trailing_metadata)
        raise Aborted


class RecordingProvider:
    name = "recording"

    def __init__(self) -> None:
        self.calls: list[str] = []

    async def extract_items(self, text: str, language: str) -> ExtractResult:
        self.calls.append(text)
        return ExtractResult(items=[], model="recording")


async def test_disabled_job_returns_unavailable_ai_disabled_without_calling_the_provider() -> None:
    registry = jobs.load(SHIPPED, {"AI_JOB_TEXT_TO_ITEMS_ENABLED": "false"}, jobs.priced_models())
    provider = RecordingProvider()
    context = FakeContext(metadata=[("x-plan", "plus")])

    with pytest.raises(Aborted):
        await AiServicer(registry, provider).ExtractItemsFromText(
            ai_pb2.ExtractItemsFromTextRequest(text="two eggs", language="en"), cast(Any, context)
        )

    assert (context.code, context.details, context.trailing, provider.calls) == (
        grpc.StatusCode.UNAVAILABLE,
        "ai_disabled: text_to_items is switched off",
        (("x-error-code", "ai_disabled"),),
        [],
    )
