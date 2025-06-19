from typing import Sequence
from uuid import UUID

from sqlalchemy import select, or_

from shared.db.models import Quiz
from shared.repositories.base_repo import BaseRepository


class QuizRepository(BaseRepository[Quiz]):
    @property
    def model(self) -> type[Quiz]:
        return Quiz

    async def get_public_quizzes(self) -> Sequence[Quiz]:
        stmt = select(Quiz).where(Quiz.is_public.is_(True))
        result = await self._session.execute(stmt)
        return result.scalars().all()

    async def get_user_quizzes(self, user_id: UUID) -> Sequence[Quiz]:
        stmt = select(Quiz).where(Quiz.owner_id == user_id)
        result = await self._session.execute(stmt)
        return result.scalars().all()

    async def get_visible_quizzes(self, user_id: UUID) -> Sequence[Quiz]:
        stmt = select(Quiz).where(
            or_(
                Quiz.is_public.is_(True),
                Quiz.owner_id == user_id
            )
        )
        result = await self._session.execute(stmt)
        return result.scalars().all()
