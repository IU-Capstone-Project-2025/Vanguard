import logging
from uuid import UUID

import jwt
from fastapi import HTTPException, Request, status

from shared.core.config import settings

logger = logging.getLogger("app")


async def get_current_user_id(request: Request) -> UUID:
    """Dependency that provides a UUID."""
    auth: str = request.headers.get("Authorization", "")

    if not auth.startswith("Bearer "):
        logger.warning("Missing or invalid Authorization header format")
        raise HTTPException(status.HTTP_401_UNAUTHORIZED, "Missing or invalid Authorization header")

    token = auth.removeprefix("Bearer ").strip()
    logger.debug(f"Processing JWT token: {token[:10]}...")

    try:
        payload = jwt.decode(token, settings.JWT_SECRET_KEY, algorithms=[settings.JWT_ALGORITHM])
        user_id = payload.get("sub")

        if not user_id:
            logger.warning("JWT token missing 'sub' claim")
            raise HTTPException(status.HTTP_401_UNAUTHORIZED, "Invalid token")

        logger.debug(f"Successfully authenticated user: {user_id}")
        return UUID(user_id)

    except jwt.ExpiredSignatureError:
        logger.warning("Expired JWT token")
        raise HTTPException(status.HTTP_401_UNAUTHORIZED, "Expired token")

    except (jwt.PyJWTError, ValueError) as exc:
        logger.warning(f"Invalid JWT token: {str(exc)}")
        raise HTTPException(status.HTTP_401_UNAUTHORIZED, "Invalid token")
