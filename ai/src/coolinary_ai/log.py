"""JSON lines on stdout (docs/service-conventions.md)."""

import json
import logging
import sys
from datetime import UTC, datetime

from .config import settings

_STANDARD = set(vars(logging.makeLogRecord({})).keys()) | {"message", "asctime", "color_message"}


class JsonFormatter(logging.Formatter):
    def format(self, record: logging.LogRecord) -> str:
        line = {
            "time": datetime.fromtimestamp(record.created, UTC).isoformat(),
            "level": record.levelname.lower(),
            "service": settings.service_name,
            "msg": record.getMessage(),
        }
        line.update({k: v for k, v in vars(record).items() if k not in _STANDARD})
        if record.exc_info:
            line["error"] = self.formatException(record.exc_info)
        return json.dumps(line, default=str)


def setup() -> None:
    handler = logging.StreamHandler(sys.stdout)
    handler.setFormatter(JsonFormatter())
    root = logging.getLogger()
    root.handlers[:] = [handler]
    root.setLevel(settings.log_level.upper())
    for noisy in ("uvicorn.access", "httpx", "opentelemetry", "grpc", "grpc._cython.cygrpc", "redis"):
        logging.getLogger(noisy).setLevel(logging.WARNING)


log = logging.getLogger("coolinary_ai")
