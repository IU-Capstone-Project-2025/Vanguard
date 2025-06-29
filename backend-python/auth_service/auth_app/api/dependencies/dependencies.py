from uuid import UUID

import jwt
from fastapi import Depends, HTTPException, Request, status

from shared.db.database import async_session_maker
from shared.utils.unitofwork import UnitOfWork

from auth_app.core.config import settings
from auth_app.services.auth_service import AuthService

uow = UnitOfWork(async_session_maker)


async def get_uow() -> UnitOfWork:
    """Dependency that provides a UnitOfWork instance."""
    return uow


async def get_auth_service(_uow: UnitOfWork = Depends(get_uow)) -> AuthService:
    """Dependency that provides a QuizService instance."""
    return AuthService(_uow)


async def get_current_user_id(request: Request) -> UUID:
    """Dependency that provides a UUID."""
    auth: str = request.headers.get("Authorization", "")
    if not auth.startswith("Bearer "):
        raise HTTPException(status.HTTP_401_UNAUTHORIZED, "Missing or invalid Authorization header")
    token = auth.removeprefix("Bearer ").strip()
    try:
        payload = jwt.decode(token, settings.SECRET_KEY, algorithms=[settings.ALGORITHM])
        user_id = payload.get("sub")
        if not user_id:
            raise HTTPException(status.HTTP_401_UNAUTHORIZED, "Invalid token")
        return UUID(user_id)
    except jwt.ExpiredSignatureError:
        raise HTTPException(status.HTTP_401_UNAUTHORIZED, "Expired token")
    except (jwt.PyJWTError, ValueError):
        raise HTTPException(status.HTTP_401_UNAUTHORIZED, "Invalid token")
