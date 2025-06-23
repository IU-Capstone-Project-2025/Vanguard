import logging

from fastapi import Request, FastAPI
from fastapi.responses import JSONResponse
from starlette import status

from quiz_app.exceptions import NotFoundError, ForbiddenError

logger = logging.getLogger("app")


def register_exception_handlers(app: FastAPI):
    @app.exception_handler(NotFoundError)
    async def not_found_exception_handler(_: Request, exc: NotFoundError):
        logger.warning(f"NotFoundError: {str(exc)}")
        return JSONResponse(
            status_code=status.HTTP_404_NOT_FOUND,
            content={"detail": str(exc)}
        )

    @app.exception_handler(ForbiddenError)
    async def forbidden_exception_handler(_: Request, exc: ForbiddenError):
        logger.warning(f"ForbiddenError: {str(exc)}")
        return JSONResponse(
            status_code=status.HTTP_403_FORBIDDEN,
            content={"detail": str(exc) or "Access forbidden"},
        )

    @app.exception_handler(Exception)
    async def generic_exception_handler(_: Request, _exc: Exception):
        logger.exception("Unhandled exception occurred")
        return JSONResponse(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            content={"detail": "An unexpected error occurred"},
        )
