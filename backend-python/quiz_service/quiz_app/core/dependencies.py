import logging
from uuid import UUID

from fastapi import Depends, Request, HTTPException

from shared.core.dependencies import get_current_user_id
from shared.db.database import async_session_maker
from shared.utils.unitofwork import UnitOfWork

from quiz_app.services.image_service import S3ImageService
from quiz_app.services.quiz_service import QuizService

logger = logging.getLogger("app")

uow = UnitOfWork(async_session_maker)


async def get_uow() -> UnitOfWork:
    """Dependency that provides a UnitOfWork instance."""
    return uow


async def get_quiz_service(_uow: UnitOfWork = Depends(get_uow)) -> QuizService:
    """Dependency that provides a QuizService instance."""
    return QuizService(_uow)


async def get_image_service() -> S3ImageService:
    """Dependency that provides a S3ImageService instance."""
    return S3ImageService()


async def get_potential_user_id(request: Request) -> UUID | None:
    try:
        user_id = await get_current_user_id(request)
        logger.debug(f"Authenticated user: {user_id}")
        return user_id
    except HTTPException as exc:
        logger.debug(f"Unauthenticated request: {exc.detail}")
        return None
