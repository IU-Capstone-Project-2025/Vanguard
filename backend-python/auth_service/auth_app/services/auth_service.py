from datetime import datetime, timedelta, UTC
from uuid import UUID, uuid4

from fastapi import Request

import bcrypt, jwt

from shared.db.models import User, RefreshToken
from shared.repositories import UserRepository, RefreshTokenRepository
from shared.schemas.auth import (
    UserCreate, UserLogin, TokenResponse, UserResponse, UpdateProfile, RefreshTokenRequest
)
from shared.utils.unitofwork import UnitOfWork

from auth_app.core.config import settings
from auth_app.exceptions import (
    EmailAlreadyExists, UsernameAlreadyExists, InvalidCredentials,
    UserNotFound, EmailInUse, UsernameInUse,
    InvalidRefreshToken, ExpiredRefreshToken
)


class AuthService:
    def __init__(self, uow: UnitOfWork):
        self.uow = uow

    @staticmethod
    def _hash(pw: str) -> str:
        return bcrypt.hashpw(pw.encode(), bcrypt.gensalt()).decode()

    @staticmethod
    def _verify(pw: str, hashed: str) -> bool:
        try:
            return bcrypt.checkpw(pw.encode(), hashed.encode())
        except (ValueError, TypeError):
            return False

    @staticmethod
    def _create_jwt(user_id: UUID, expires_delta: timedelta) -> str:
        payload = {
            "sub": str(user_id),
            "iat": datetime.now(UTC),
            "exp": datetime.now(UTC) + expires_delta,
            "jti": str(uuid4())
        }
        return jwt.encode(payload, settings.SECRET_KEY, algorithm=settings.ALGORITHM)

    def _create_tokens(self, user_id: UUID, request: Request):
        access = self._create_jwt(user_id, timedelta(minutes=settings.ACCESS_TOKEN_EXPIRE_MINUTES))
        refresh_expires = datetime.now(UTC) + timedelta(days=settings.REFRESH_TOKEN_EXPIRE_DAYS)

        ua = request.headers.get("user-agent", "")
        ip = request.client.host or "0.0.0.0"

        return access, refresh_expires, ua, ip

    async def register(self, data: UserCreate) -> UserResponse:
        async with self.uow.transaction() as session:
            user_repo = UserRepository(session)

            if await user_repo.get_by_email(str(data.email).lower()):
                raise EmailAlreadyExists
            if await user_repo.get_by_username(data.username.lower()):
                raise UsernameAlreadyExists

            hashed_pw = self._hash(data.password.get_secret_value())
            user = await user_repo.create(User(
                username=data.username.lower(),
                email=str(data.email).lower(),
                password_hash=hashed_pw
            ))

            return UserResponse.model_validate(user)

    async def login(self, data: UserLogin, request: Request) -> TokenResponse:
        async with self.uow.transaction() as session:
            user_repo = UserRepository(session)
            user = await user_repo.get_by_email(str(data.email).lower())
            if not user or not self._verify(data.password.get_secret_value(), user.password_hash):
                raise InvalidCredentials()

            access, refresh_expires, ua, ip = self._create_tokens(user.id, request)
            rt_repo = RefreshTokenRepository(session)
            rt = await rt_repo.create(RefreshToken(
                user_id=user.id,
                user_agent=ua,
                ip_address=ip,
                expires_at=refresh_expires
            ))

            return TokenResponse(
                access_token=access,
                refresh_token=str(rt.token)
            )

    async def refresh(self, req: RefreshTokenRequest, request: Request) -> TokenResponse:
        async with self.uow.transaction() as session:
            rt_repo = RefreshTokenRepository(session)
            token = await rt_repo.get(req.refresh_token)

            if not token:
                raise InvalidRefreshToken()
            if token.expires_at < datetime.now(UTC):
                raise ExpiredRefreshToken()

            await rt_repo.revoke(token)

            access, refresh_expires, ua, ip = self._create_tokens(token.user_id, request)
            new_rt = await rt_repo.create(RefreshToken(
                user_id=token.user_id,
                user_agent=ua,
                ip_address=ip,
                expires_at=refresh_expires
            ))

            return TokenResponse(
                access_token=access,
                refresh_token=str(new_rt.token)
            )

    async def me(self, user_id: UUID) -> UserResponse:
        async with self.uow.readonly() as session:
            repo = UserRepository(session)
            user = await repo.get(_id=user_id)
            if not user:
                raise UserNotFound()
            return UserResponse.model_validate(user)

    async def update_me(self, user_id: UUID, data: UpdateProfile) -> UserResponse:
        async with self.uow.transaction() as session:
            repo = UserRepository(session)
            user = await repo.get(_id=user_id)
            if not user:
                raise UserNotFound()

            updated = False

            if data.email is not None:
                normalized_email = str(data.email).lower()
                if normalized_email != user.email:
                    if await repo.get_by_email(normalized_email):
                        raise EmailInUse()
                    user.email = normalized_email
                    updated = True

            if data.username is not None:
                normalized_username = data.username.lower()
                if normalized_username != user.username:
                    if await repo.get_by_username(normalized_username):
                        raise UsernameInUse()
                    user.username = normalized_username
                    updated = True

            # TODO: what about verifying previous password?
            if data.password is not None:
                user.password_hash = self._hash(data.password.get_secret_value())
                updated = True

            if updated:
                user = await repo.update(user)

            return UserResponse.model_validate(user)

    async def logout(self, token: UUID) -> None:
        async with self.uow.transaction() as session:
            rt_repo = RefreshTokenRepository(session)
            token_obj = await rt_repo.get(token)
            if token_obj:
                await rt_repo.revoke(token_obj)

    async def logout_all(self, user_id: UUID) -> None:
        async with self.uow.transaction() as session:
            rt_repo = RefreshTokenRepository(session)
            await rt_repo.revoke_all_for_user(user_id)
