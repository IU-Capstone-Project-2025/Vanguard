from fastapi import Depends

from shared.db.database import async_session_maker
from shared.utils.unitofwork import UnitOfWork

from auth_app.services.auth_service import AuthService

uow = UnitOfWork(async_session_maker)


async def get_uow() -> UnitOfWork:
    """Dependency that provides a UnitOfWork instance."""
    return uow


async def get_auth_service(_uow: UnitOfWork = Depends(get_uow)) -> AuthService:
    """Dependency that provides a QuizService instance."""
    return AuthService(_uow)
