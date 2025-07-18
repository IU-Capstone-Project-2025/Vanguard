from typing import cast
from uuid import UUID

from sqlalchemy import delete

from shared.db.models import RefreshToken
from shared.repositories.base_repo import BaseRepository


class RefreshTokenRepository(BaseRepository[RefreshToken]):
    @property
    def model(self) -> type[RefreshToken]:
        return RefreshToken

    async def revoke(self, token: RefreshToken) -> None:
        await self.delete(token)
        await self._session.flush()

    async def revoke_all_for_user(self, user_id: UUID) -> int:
        stmt = delete(RefreshToken).where(RefreshToken.user_id == user_id)
        result = await self._session.execute(stmt)
        rowcount = cast(int, result.rowcount)  # Number of rows deleted
        return rowcount
