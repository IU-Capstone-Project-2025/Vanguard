import logging
from uuid import UUID

from shared.db.models import User as DBUser
from shared.repositories import UserRepository
from shared.schemas.user import UserCreate, UserResponse
from shared.utils.unitofwork import UnitOfWork

logger = logging.getLogger("app")


class UserService:
    def __init__(self, uow: UnitOfWork):
        self.uow = uow

    async def create_user(self, user_in: UserCreate, hash_password_func: callable) -> UserResponse:
        async with self.uow.transaction() as session:
            repo = UserRepository(session)
            user_data = user_in.model_dump(exclude={"password"})
            user_data["password_hash"] = hash_password_func(user_in.password)
            user_db = await repo.create(DBUser(**user_data))
            return UserResponse.model_validate(user_db)

    async def get_user_by_id(self, user_id: UUID) -> UserResponse:
        async with self.uow.readonly() as session:
            repo = UserRepository(session)
            user = await repo.get(_id=user_id)
            if not user:
                raise ValueError()
            return UserResponse.model_validate(user)

    async def get_user_by_email(self, email: str) -> UserResponse:
        async with self.uow.readonly() as session:
            repo = UserRepository(session)
            user = await repo.get_by_email(email=email)
            if not user:
                raise ValueError()
            return UserResponse.model_validate(user)

    async def get_user_by_username(self, username: str) -> UserResponse:
        async with self.uow.readonly() as session:
            repo = UserRepository(session)
            user = await repo.get_by_username(username=username)
            if not user:
                raise ValueError()
            return UserResponse.model_validate(user)
