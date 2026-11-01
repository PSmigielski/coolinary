"""Entry point: `python -m coolinary_ai`."""

import uvicorn

from .config import settings


def main() -> None:
    uvicorn.run(
        "coolinary_ai.app:app",
        # Every interface of the container; only the Docker network reaches it.
        host="0.0.0.0",  # noqa: S104
        port=settings.admin_port,
        workers=1,  # one process: it also owns the gRPC port
        loop="asyncio",
        log_config=None,
        access_log=False,
    )


if __name__ == "__main__":
    main()
