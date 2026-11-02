"""grpc.aio server: AiService + grpc.health.v1 (+ reflection in dev) + metrics.

Runs on the same event loop as FastAPI (started in its lifespan). Keep handlers async;
push blocking or CPU-heavy work to a thread (asyncio.to_thread) so neither server stalls.
"""

import time
from collections.abc import Awaitable, Callable
from dataclasses import asdict
from typing import Any, cast

import grpc
from grpc_health.v1 import health, health_pb2, health_pb2_grpc
from grpc_reflection.v1alpha import reflection
from prometheus_client import Counter, Histogram

from coolinary.ai.v1 import ai_pb2, ai_pb2_grpc
from coolinary.common.v1 import common_pb2

from . import cache, jobs, tracing
from .config import settings
from .llm import Provider, get_provider
from .log import log

SERVICE_NAME = ai_pb2.DESCRIPTOR.services_by_name["AiService"].full_name
MAX_TEXT = 2000

# Metric names are shared by all services (docs/service-conventions.md).
RPC_REQUESTS = Counter(
    "rpc_server_requests_total", "RPCs handled.", ["rpc_service", "rpc_method", "rpc_code"]
)
RPC_DURATION = Histogram(
    "rpc_server_duration_seconds",
    "RPC handling time.",
    ["rpc_service", "rpc_method"],
    buckets=(0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 30),
)


class MetricsInterceptor(grpc.aio.ServerInterceptor):
    async def intercept_service(
        self,
        continuation: Callable[[grpc.HandlerCallDetails], Awaitable[grpc.RpcMethodHandler[Any, Any] | None]],
        handler_call_details: grpc.HandlerCallDetails,
    ) -> grpc.RpcMethodHandler[Any, Any] | None:
        handler = await continuation(handler_call_details)
        if handler is None or handler.unary_unary is None:
            return handler
        service, _, method = handler_call_details.method.lstrip("/").rpartition("/")
        # Under grpc.aio the handler is a coroutine function taking an aio context;
        # grpcio's stubs only describe the sync API.
        inner = cast(Callable[[Any, grpc.aio.ServicerContext[Any, Any]], Awaitable[Any]], handler.unary_unary)

        async def wrapped(request: Any, context: grpc.aio.ServicerContext[Any, Any]) -> Any:
            start = time.perf_counter()
            failed = False
            try:
                return await inner(request, context)
            except Exception:
                failed = True
                raise
            finally:
                code = context.code()
                name = code.name.lower() if code else ("unknown" if failed else "ok")
                RPC_REQUESTS.labels(service, method, name).inc()
                RPC_DURATION.labels(service, method).observe(time.perf_counter() - start)
                log.debug("rpc", extra={"method": handler_call_details.method, "code": name})

        return grpc.unary_unary_rpc_method_handler(
            wrapped,
            request_deserializer=handler.request_deserializer,
            response_serializer=handler.response_serializer,
        )


def _to_proto(items: list[dict[str, Any]]) -> list[ai_pb2.ExtractedItem]:
    out = []
    for i in items:
        q = None
        if i["quantity"] is not None:
            q = common_pb2.Quantity(value=i["quantity"], unit=cast(common_pb2.Unit, int(i["unit"])))
        out.append(ai_pb2.ExtractedItem(name=i["name"], quantity=q, confidence=i["confidence"]))
    return out


async def _abort_disabled(context: grpc.aio.ServicerContext[Any, Any], job: jobs.Job) -> None:
    """A job switched off in the registry: UNAVAILABLE, so the caller runs its non-AI fallback."""
    await context.abort(
        grpc.StatusCode.UNAVAILABLE,
        f"ai_disabled: {job.name} is switched off",
        trailing_metadata=(("x-error-code", "ai_disabled"),),
    )


class AiServicer(ai_pb2_grpc.AiServiceServicer):
    def __init__(self, registry: jobs.Registry, provider: Provider) -> None:
        self.registry = registry
        self.provider = provider

    async def ExtractItemsFromText(
        self,
        request: ai_pb2.ExtractItemsFromTextRequest,
        context: grpc.aio.ServicerContext[
            ai_pb2.ExtractItemsFromTextRequest, ai_pb2.ExtractItemsFromTextResponse
        ],
    ) -> ai_pb2.ExtractItemsFromTextResponse:
        text = request.text.strip()
        if not text:
            await context.abort(grpc.StatusCode.INVALID_ARGUMENT, "text is required")
        if len(text) > MAX_TEXT:
            await context.abort(
                grpc.StatusCode.INVALID_ARGUMENT, f"text is longer than {MAX_TEXT} characters"
            )
        # Text metadata only; binary (-bin) keys carry bytes.
        md = {k: v for k, v in context.invocation_metadata() or () if isinstance(v, str)}
        job = self.registry.job("text_to_items", md.get("x-plan", ""))
        if not job.enabled:
            await _abort_disabled(context, job)
        language = request.language or "en"

        k = cache.key("text_to_items", self.provider.name, language, text)
        if hit := await cache.get_json(k):
            return ai_pb2.ExtractItemsFromTextResponse(
                items=_to_proto(hit["items"]), model=hit["model"], cached=True
            )

        with tracing.generation(
            feature="text_to_items",
            plan=md.get("x-plan", ""),
            household_id=request.household_id or md.get("x-household-id", ""),
            model=self.provider.name,
            input={"feature": "text_to_items", "language": language, "chars": len(text)},
        ) as gen:
            result = await self.provider.extract_items(text, language)
            if gen is not None:
                gen.update(
                    model=result.model, output={"items": len(result.items)}, usage_details=result.usage
                )

        items = [asdict(i) | {"unit": int(i.unit)} for i in result.items]
        await cache.set_json(k, {"items": items, "model": result.model})
        return ai_pb2.ExtractItemsFromTextResponse(items=_to_proto(items), model=result.model, cached=False)


async def start(registry: jobs.Registry) -> grpc.aio.Server:
    server = grpc.aio.server(
        interceptors=[MetricsInterceptor()],
        options=[
            ("grpc.so_reuseport", 0),  # one process owns the port
            ("grpc.keepalive_permit_without_calls", 1),
            ("grpc.http2.min_ping_interval_without_data_ms", 10_000),
        ],
    )
    ai_pb2_grpc.add_AiServiceServicer_to_server(
        AiServicer(registry, get_provider(settings.llm_provider)), server
    )
    health_servicer = health.aio.HealthServicer()
    health_pb2_grpc.add_HealthServicer_to_server(health_servicer, server)
    for name in ("", SERVICE_NAME):
        await health_servicer.set(name, health_pb2.HealthCheckResponse.SERVING)
    if settings.grpc_reflection:
        reflection.enable_server_reflection(
            (SERVICE_NAME, health.SERVICE_NAME, reflection.SERVICE_NAME), server
        )
    server.add_insecure_port(f"0.0.0.0:{settings.grpc_port}")
    await server.start()
    log.info(
        "gRPC server listening",
        extra={
            "port": settings.grpc_port,
            "services": [SERVICE_NAME],
            "reflection": settings.grpc_reflection,
            "llm_provider": settings.llm_provider,
        },
    )
    return server
