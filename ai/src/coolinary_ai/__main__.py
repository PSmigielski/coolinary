"""Entry point: `python -m coolinary_ai`."""

import uvicorn

from .config import settings


def main() -> None:
    uvicorn.run(
        "coolinary_ai.app:app",
        host="0.0.0.0",
        port=settings.admin_port,
        workers=1,  # one process: it also owns the gRPC port
        loop="asyncio",
        log_config=None,
        access_log=False,
    )


if __name__ == "__main__":
    main()
