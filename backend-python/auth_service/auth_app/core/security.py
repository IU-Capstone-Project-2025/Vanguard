from datetime import datetime, timedelta, UTC

import jwt
from passlib.context import CryptContext

from auth_app.core.config import settings
from auth_app.exceptions import TokenExpiredError, InvalidTokenError, MissingTokenError, InvalidTokenPayloadError

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")


def get_password_hash(password: str) -> str:
    return pwd_context.hash(password)


def verify_password(plain_password: str, hashed_password: str) -> bool:
    return pwd_context.verify(plain_password, hashed_password)


def create_jwt_token(data: dict) -> str:
    to_encode = data.copy()
    expire = datetime.now(UTC) + timedelta(minutes=settings.ACCESS_TOKEN_EXPIRE_MINUTES)
    to_encode.update({"exp": expire})
    encoded_jwt = jwt.encode(to_encode, settings.SECRET_KEY, algorithm=settings.ALGORITHM)
    return encoded_jwt


def decode_jwt_token(token: str):
    return jwt.decode(token, settings.SECRET_KEY, algorithms=[settings.ALGORITHM])


def verify_jwt_token(token: str | None) -> dict[str, str]:
    if not token:
        raise MissingTokenError()

    try:
        payload: dict[str, str] = decode_jwt_token(token)
        username = payload.get("sub")
        if not username:
            raise InvalidTokenPayloadError("sub")
        return payload

    except jwt.ExpiredSignatureError:
        raise TokenExpiredError()

    except (jwt.PyJWTError, ValueError):
        raise InvalidTokenError()
