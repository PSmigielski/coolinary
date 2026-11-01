# syntax=docker/dockerfile:1.7
# Stages: base (python + uv) -> deps -> dev (debugpy + watchfiles) ; deps -> build -> prod (slim, non-root)
ARG PYTHON_VERSION=3.14
ARG UV_VERSION=0.13

FROM ghcr.io/astral-sh/uv:${UV_VERSION} AS uv

FROM python:${PYTHON_VERSION}-slim AS base
COPY --from=uv /uv /uvx /usr/local/bin/
ENV UV_COMPILE_BYTECODE=1 \
    UV_LINK_MODE=copy \
    UV_PYTHON_DOWNLOADS=never \
    UV_PROJECT_ENVIRONMENT=/opt/venv \
    PYTHONUNBUFFERED=1
WORKDIR /app

FROM base AS deps
COPY pyproject.toml uv.lock ./
RUN --mount=type=cache,target=/root/.cache/uv uv sync --frozen --no-install-project --no-dev

FROM deps AS dev
RUN --mount=type=cache,target=/root/.cache/uv uv sync --frozen --no-install-project
COPY . .
RUN --mount=type=cache,target=/root/.cache/uv uv sync --frozen
# Source is bind-mounted over /app; the venv lives in /opt/venv so it stays visible.
ENV PATH=/opt/venv/bin:$PATH PYTHONDONTWRITEBYTECODE=1
EXPOSE 50051 8081 5678
HEALTHCHECK --interval=10s --timeout=3s --start-period=60s --retries=5 \
  CMD ["python", "-c", "import sys,urllib.request; sys.exit(0 if urllib.request.urlopen('http://127.0.0.1:8081/healthz', timeout=3).status == 200 else 1)"]
# watchfiles restarts the whole process (with debugpy) on change; attach from Zed on :5678.
CMD ["watchfiles", "--filter", "python", "python -m debugpy --listen 0.0.0.0:5678 -m coolinary_ai", "src"]

FROM deps AS build
COPY . .
RUN --mount=type=cache,target=/root/.cache/uv uv sync --frozen --no-dev --no-editable

FROM python:${PYTHON_VERSION}-slim AS prod
RUN groupadd --system --gid 65532 app && useradd --system --uid 65532 --gid app --no-create-home app
COPY --from=build /opt/venv /opt/venv
ENV PATH=/opt/venv/bin:$PATH PYTHONUNBUFFERED=1 PYTHONDONTWRITEBYTECODE=1
USER 65532:65532
EXPOSE 50051 8081
HEALTHCHECK --interval=15s --timeout=3s --start-period=15s --retries=3 \
  CMD ["python", "-c", "import sys,urllib.request; sys.exit(0 if urllib.request.urlopen('http://127.0.0.1:8081/healthz', timeout=3).status == 200 else 1)"]
CMD ["python", "-m", "coolinary_ai"]
