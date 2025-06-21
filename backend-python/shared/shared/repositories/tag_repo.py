from typing import List, Sequence
from uuid import UUID

from sqlalchemy import select

from shared.db.models import Tag, QuizTag
from shared.repositories.base_repo import BaseRepository


class TagRepository(BaseRepository[Tag]):
    @property
    def model(self) -> type[Tag]:
        return Tag

    async def get_by_name(self, name: str) -> Tag | None:
        stmt = select(Tag).where(Tag.name == name)
        result = await self._session.execute(stmt)
        return result.scalar_one_or_none()

    async def get_or_create_bulk(self, tag_names: List[str]) -> List[Tag]:
        """
        Given a list of normalized tag names, returns Tag objects,
        creating any that don’t already exist.
        """
        stmt = select(Tag).where(Tag.name.in_(tag_names))
        result = await self._session.execute(stmt)
        existing_tags = {t.name: t for t in result.scalars().all()}

        tags_result: List[Tag] = []
        for name in tag_names:
            if name in existing_tags:
                tags_result.append(existing_tags[name])
            else:
                tag = Tag(name=name)
                await self.create(tag)
                tags_result.append(tag)
        return tags_result

    async def get_tags_for_quiz(self, quiz_id: UUID) -> Sequence[Tag]:
        stmt = select(Tag).join(QuizTag).where(QuizTag.quiz_id == quiz_id)
        result = await self._session.execute(stmt)
        return result.scalars().all()
