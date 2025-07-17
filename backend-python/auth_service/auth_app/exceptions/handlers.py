import logging

from fastapi import Request, FastAPI
from fastapi.responses import JSONResponse
from starlette import status

from auth_app.exceptions import AuthServiceError

logger = logging.getLogger("app")


def register_exception_handlers(app: FastAPI):
    @app.exception_handler(AuthServiceError)
    async def auth_service_exception_handler(_: Request, exc: AuthServiceError):
        detail = getattr(exc, "detail", str(exc))
        logger.warning(f"AuthServiceError: {detail}")
        return JSONResponse(
            status_code=getattr(exc, "status_code", status.HTTP_500_INTERNAL_SERVER_ERROR),
            content={"detail": detail}
        )

    @app.exception_handler(Exception)
    async def generic_exception_handler(_: Request, exc: Exception):
        logger.exception(f"Unhandled exception occurred: {str(exc)}")
        return JSONResponse(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            content={"detail": "An unexpected error occurred"},
        )
