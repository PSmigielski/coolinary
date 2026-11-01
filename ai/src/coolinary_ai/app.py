"""Admin HTTP app on :8081 (/healthz, /metrics). Its lifespan starts and stops the gRPC server,
so both share one event loop and one process (run with a single uvicorn worker)."""

from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.responses import JSONResponse
from prometheus_client import make_asgi_app

from . import cache, grpc_server, log, tracing


@asynccontextmanager
async def lifespan(app: FastAPI):
    log.setup()
    app.state.grpc = await grpc_server.start()
    yield
    await app.state.grpc.stop(grace=5)
    await cache.client.aclose()
    tracing.shutdown()


app = FastAPI(title="coolinary-ai admin", lifespan=lifespan, docs_url=None, redoc_url=None, openapi_url=None)
app.mount("/metrics", make_asgi_app())


@app.get("/healthz")
async def healthz() -> JSONResponse:
    checks = {"grpc": "ok" if getattr(app.state, "grpc", None) else "starting"}
    try:
        await cache.client.ping()
        checks["redis_cache"] = "ok"
    except Exception as err:  # noqa: BLE001
        checks["redis_cache"] = f"down: {err}"
    ok = all(v == "ok" for v in checks.values())
    return JSONResponse({"ok": ok, "checks": checks}, status_code=200 if ok else 503)
