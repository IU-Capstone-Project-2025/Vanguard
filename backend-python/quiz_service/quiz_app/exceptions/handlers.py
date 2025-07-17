import logging

from fastapi import Request, FastAPI
from fastapi.responses import JSONResponse
from starlette import status

from quiz_app.exceptions import QuizServiceError, ImageServiceError

logger = logging.getLogger("app")


def register_exception_handlers(app: FastAPI):
    @app.exception_handler(ImageServiceError)
    async def image_service_exception_handler(_: Request, exc: ImageServiceError):
        detail = getattr(exc, "detail", str(exc))
        logger.warning(f"ImageServiceError: {detail}")
        return JSONResponse(
            status_code=getattr(exc, "status_code", status.HTTP_500_INTERNAL_SERVER_ERROR),
            content={"detail": detail}
        )

    @app.exception_handler(QuizServiceError)
    async def quiz_service_exception_handler(_: Request, exc: QuizServiceError):
        detail = getattr(exc, "detail", str(exc))
        logger.warning(f"QuizServiceError: {detail}")
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
