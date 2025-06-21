import logging

from fastapi import Request, FastAPI
from fastapi.responses import JSONResponse
from starlette import status

from quiz_app.exceptions import (NotFoundError, ForbiddenError, BadRequestError, InvalidImageError, FileTooLargeError, ImageNotFoundError, ImageS3Error)

logger = logging.getLogger("app")


def register_exception_handlers(app: FastAPI):
    @app.exception_handler(InvalidImageError)
    async def invalid_image_exception_handler(_: Request, exc: InvalidImageError):
        logger.warning(f"InvalidImageError: {str(exc)}")
        return JSONResponse(
            status_code=status.HTTP_400_BAD_REQUEST,
            content={"detail": str(exc)}
        )

    @app.exception_handler(FileTooLargeError)
    async def file_too_large_exception_handler(_: Request, exc: FileTooLargeError):
        logger.warning(f"FileTooLargeError: {str(exc)}")
        return JSONResponse(
            status_code=status.HTTP_413_REQUEST_ENTITY_TOO_LARGE,
            content={"detail": str(exc)}
        )

    @app.exception_handler(ImageNotFoundError)
    async def file_not_found_exception_handler(_: Request, exc: ImageNotFoundError):
        logger.warning(f"ImageNotFoundError: {str(exc)}")
        return JSONResponse(
            status_code=status.HTTP_404_NOT_FOUND,
            content={"detail": str(exc)}
        )

    @app.exception_handler(ImageS3Error)
    async def image_exception_handler(_: Request, exc: ImageS3Error):
        logger.warning(f"ImageS3Error: {str(exc)}")
        return JSONResponse(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            content={"detail": str(exc)}
        )

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

    @app.exception_handler(BadRequestError)
    async def bad_request_exception_handler(_: Request, exc: BadRequestError):
        logger.warning(f"BadRequestError: {str(exc)}")
        return JSONResponse(
            status_code=status.HTTP_400_BAD_REQUEST,
            content={"detail": str(exc)},
        )

    @app.exception_handler(Exception)
    async def generic_exception_handler(_: Request, _exc: Exception):
        logger.exception("Unhandled exception occurred")
        return JSONResponse(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            content={"detail": "An unexpected error occurred"},
        )
