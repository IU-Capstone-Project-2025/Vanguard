import logging
from datetime import datetime, timedelta, UTC
from uuid import UUID, uuid4

from fastapi import Request

import bcrypt, jwt

from shared.core.config import settings as shared_settings
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
from auth_app.core.metrics import (
    AUTH_USER_REGISTRATIONS_TOTAL, AUTH_USER_LOGINS_TOTAL, AUTH_TOKEN_REFRESHES_TOTAL,
    AUTH_USER_LOGOUTS_TOTAL, AUTH_ACTIVE_SESSIONS, SERVICE
)

logger = logging.getLogger("app")


class AuthService:
    def __init__(self, uow: UnitOfWork):
        self.uow = uow

    @staticmethod
    def _hash(pw: str) -> str:
        logger.debug("Hashing password")
        return bcrypt.hashpw(pw.encode(), bcrypt.gensalt()).decode()

    @staticmethod
    def _verify(pw: str, hashed: str) -> bool:
        try:
            result = bcrypt.checkpw(pw.encode(), hashed.encode())
            logger.debug(f"Password verification result: {result}")
            return result
        except (ValueError, TypeError) as e:
            logger.warning(f"Password verification failed due to error: {e}")
            return False

    @staticmethod
    def _create_jwt(user_id: UUID, expires_delta: timedelta) -> str:
        exp = datetime.now(UTC) + expires_delta
        logger.debug(f"Creating JWT for user: {user_id}, expires at: {exp.isoformat()}")

        payload = {
            "sub": str(user_id),
            "iat": datetime.now(UTC),
            "exp": exp,
            "jti": str(uuid4())
        }
        return jwt.encode(payload, shared_settings.JWT_SECRET_KEY, algorithm=shared_settings.JWT_ALGORITHM)

    def _create_tokens(self, user_id: UUID, request: Request):
        logger.debug(f"Generating access/refresh tokens for user: {user_id}")
        access = self._create_jwt(user_id, timedelta(minutes=settings.ACCESS_TOKEN_EXPIRE_MINUTES))
        refresh_expires = datetime.now(UTC) + timedelta(days=settings.REFRESH_TOKEN_EXPIRE_DAYS)

        ua = request.headers.get("user-agent", "")
        ip = request.client.host or "0.0.0.0"
        logger.debug(f"Token context — IP: {ip}, User-Agent: {ua}")

        return access, refresh_expires, ua, ip

    async def register(self, data: UserCreate) -> UserResponse:
        logger.debug(f"Starting user registration for email: {data.email}")
        async with self.uow.transaction() as session:
            user_repo = UserRepository(session)

            if await user_repo.get_by_email(str(data.email).lower()):
                AUTH_USER_REGISTRATIONS_TOTAL.labels(service=SERVICE, status="fail").inc()
                raise EmailAlreadyExists
            if await user_repo.get_by_username(data.username.lower()):
                AUTH_USER_REGISTRATIONS_TOTAL.labels(service=SERVICE, status="fail").inc()
                raise UsernameAlreadyExists

            hashed_pw = self._hash(data.password.get_secret_value())
            user = await user_repo.create(User(
                username=data.username.lower(),
                email=str(data.email).lower(),
                password_hash=hashed_pw
            ))

            AUTH_USER_REGISTRATIONS_TOTAL.labels(service=SERVICE, status="success").inc()
            logger.info(f"User registered successfully: {user.id}")
            return UserResponse.model_validate(user)

    async def login(self, data: UserLogin, request: Request) -> TokenResponse:
        logger.debug(f"Login attempt for email: {data.email}")
        async with self.uow.transaction() as session:
            user_repo = UserRepository(session)
            user = await user_repo.get_by_email(str(data.email).lower())
            if not user or not self._verify(data.password.get_secret_value(), user.password_hash):
                AUTH_USER_LOGINS_TOTAL.labels(service=SERVICE, status="fail").inc()
                raise InvalidCredentials()

            access, refresh_expires, ua, ip = self._create_tokens(user.id, request)
            rt_repo = RefreshTokenRepository(session)
            rt = await rt_repo.create(RefreshToken(
                user_id=user.id,
                user_agent=ua,
                ip_address=ip,
                expires_at=refresh_expires
            ))

            AUTH_USER_LOGINS_TOTAL.labels(service=SERVICE, status="success").inc()
            AUTH_ACTIVE_SESSIONS.labels(service=SERVICE).inc()
            logger.info(f"User logged in successfully: {user.id}")

            return TokenResponse(
                access_token=access,
                refresh_token=str(rt.token)
            )

    async def refresh(self, req: RefreshTokenRequest, request: Request) -> TokenResponse:
        logger.debug(f"Refreshing token: {req.refresh_token}")
        async with self.uow.transaction() as session:
            rt_repo = RefreshTokenRepository(session)
            token = await rt_repo.get(req.refresh_token)

            if not token:
                AUTH_TOKEN_REFRESHES_TOTAL.labels(service=SERVICE, status="fail").inc()
                raise InvalidRefreshToken()
            if token.expires_at < datetime.now(UTC):
                AUTH_TOKEN_REFRESHES_TOTAL.labels(service=SERVICE, status="expired").inc()
                raise ExpiredRefreshToken()

            await rt_repo.revoke(token)

            access, refresh_expires, ua, ip = self._create_tokens(token.user_id, request)
            new_rt = await rt_repo.create(RefreshToken(
                user_id=token.user_id,
                user_agent=ua,
                ip_address=ip,
                expires_at=refresh_expires
            ))

            AUTH_TOKEN_REFRESHES_TOTAL.labels(service=SERVICE, status="success").inc()
            logger.info(f"Token refreshed successfully for user: {token.user_id}")

            return TokenResponse(
                access_token=access,
                refresh_token=str(new_rt.token)
            )

    async def me(self, user_id: UUID) -> UserResponse:
        logger.debug(f"Fetching profile for user: {user_id}")
        async with self.uow.readonly() as session:
            repo = UserRepository(session)
            user = await repo.get(_id=user_id)
            if not user:
                raise UserNotFound()
            logger.debug(f"User profile fetched successfully: {user_id}")
            return UserResponse.model_validate(user)

    async def update_me(self, user_id: UUID, data: UpdateProfile) -> UserResponse:
        logger.debug(f"Updating profile for user: {user_id}")
        async with self.uow.transaction() as session:
            repo = UserRepository(session)
            user = await repo.get(_id=user_id)
            if not user:
                raise UserNotFound()

            updated = False

            if data.email is not None:
                normalized_email = str(data.email).lower()
                if normalized_email != user.email:
                    logger.debug(f"Attempting email update: {user.email} -> {normalized_email}")
                    if await repo.get_by_email(normalized_email):
                        raise EmailInUse()
                    user.email = normalized_email
                    updated = True

            if data.username is not None:
                normalized_username = data.username.lower()
                if normalized_username != user.username:
                    logger.debug(f"Attempting username update: {user.username} -> {normalized_username}")
                    if await repo.get_by_username(normalized_username):
                        raise UsernameInUse()
                    user.username = normalized_username
                    updated = True

            if data.password is not None:
                logger.debug(f"Updating password for user: {user_id}")
                user.password_hash = self._hash(data.password.get_secret_value())
                updated = True

            if updated:
                logger.info(f"Profile updated for user: {user_id}")
                user = await repo.update(user)
            else:
                logger.debug(f"No changes detected for user: {user_id}")

            return UserResponse.model_validate(user)

    async def logout(self, token: UUID) -> None:
        logger.debug(f"Logging out refresh token: {token}")
        async with self.uow.transaction() as session:
            rt_repo = RefreshTokenRepository(session)
            token_obj = await rt_repo.get(token)
            if token_obj:
                await rt_repo.revoke(token_obj)
                AUTH_USER_LOGOUTS_TOTAL.labels(service=SERVICE, scope="single").inc()
                AUTH_ACTIVE_SESSIONS.labels(service=SERVICE).dec()
                logger.info(f"Logged out single session for user: {token_obj.user_id}")

    async def logout_all(self, user_id: UUID) -> None:
        logger.debug(f"Logging out all sessions for user: {user_id}")
        async with self.uow.transaction() as session:
            rt_repo = RefreshTokenRepository(session)
            count = await rt_repo.revoke_all_for_user(user_id)
            AUTH_USER_LOGOUTS_TOTAL.labels(service=SERVICE, scope="all").inc()
            AUTH_ACTIVE_SESSIONS.labels(service=SERVICE).dec(count)
            logger.info(f"Logged out all sessions for user: {user_id}, revoked {count} tokens")
