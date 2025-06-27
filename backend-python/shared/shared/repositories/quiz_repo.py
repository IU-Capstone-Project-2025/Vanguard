import re
from typing import Sequence
from uuid import UUID

from sqlalchemy import select, or_, and_, func
from sqlalchemy.orm import selectinload

from shared.db.models import Quiz, QuizTag, Tag
from shared.repositories.base_repo import BaseRepository
from shared.schemas.quiz import QuizFilterMode


class QuizRepository(BaseRepository[Quiz]):
    @property
    def model(self) -> type[Quiz]:
        return Quiz

    async def list_quizzes(
            self,
            *,
            mode: QuizFilterMode,
            requester_id: UUID | None = None,
            user_id: UUID | None = None,
            search: str | None = None,
            tags: list[str] | None = None,
            page: int = 1,
            size: int = 20
    ) -> Sequence[Quiz]:
        # Eager‐load tags
        stmt = select(Quiz).options(selectinload(Quiz.tags))

        # Handle filtering modes
        if mode == QuizFilterMode.ALL_PUBLIC:
            stmt = stmt.where(Quiz.is_public.is_(True))
        elif mode == QuizFilterMode.ALL_MINE:
            stmt = stmt.where(Quiz.owner_id == requester_id)
        elif mode == QuizFilterMode.MINE_PUBLIC:
            stmt = stmt.where(and_(Quiz.owner_id == requester_id, Quiz.is_public.is_(True)))
        elif mode == QuizFilterMode.MINE_PRIVATE:
            stmt = stmt.where(and_(Quiz.owner_id == requester_id, Quiz.is_public.is_(False)))
        elif mode == QuizFilterMode.OTHER_PUBLIC:
            stmt = stmt.where(and_(Quiz.owner_id == user_id, Quiz.is_public.is_(True)))
        elif mode == QuizFilterMode.VISIBLE_TO_ME:
            stmt = stmt.where(or_(Quiz.owner_id == requester_id, Quiz.is_public.is_(True)))

        # Apply user_id filtering
        if user_id is not None:
            stmt = stmt.where(Quiz.owner_id == user_id)

        # Apply search
        if search:
            cleaned_search = re.sub(r"[^\w\s]", "", search).strip()
            if cleaned_search:
                tsvector = func.to_tsvector("english", Quiz.title + " " + Quiz.description)

                query_terms = " & ".join(cleaned_search.split())
                tsquery = func.to_tsquery("english", query_terms)

                stmt = stmt.where(tsvector.op("@@")(tsquery))
                stmt = stmt.order_by(
                    func.ts_rank(tsvector, tsquery).desc(),
                    Quiz.created_at.desc()
                )

        # Apply tags filtering
        if tags:
            stmt = (
                stmt
                .join(QuizTag, QuizTag.quiz_id == Quiz.id)
                .join(Tag, Tag.id == QuizTag.tag_id)
                .where(Tag.name.in_(tags))
                .group_by(Quiz.id)
                .having(func.count(Tag.id) == len(set(tags)))
            )

        # Apply paging
        offset = (page - 1) * size
        stmt = stmt.limit(size).offset(offset)

        res = await self._session.execute(stmt)
        return res.scalars().all()
