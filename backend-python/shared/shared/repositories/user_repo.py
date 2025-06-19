from sqlalchemy import select

from shared.db.models import User
from shared.repositories.base_repo import BaseRepository


class UserRepository(BaseRepository[User]):
    @property
    def model(self) -> type[User]:
        return User

    async def get_by_email(self, email: str) -> User | None:
        stmt = select(User).where(User.email == email)
        result = await self._session.execute(stmt)
        return result.scalar_one_or_none()

    async def get_by_username(self, username: str) -> User | None:
        stmt = select(User).where(User.username == username)
        result = await self._session.execute(stmt)
        return result.scalar_one_or_none()
