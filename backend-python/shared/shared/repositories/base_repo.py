from abc import ABC, abstractmethod
from collections.abc import Sequence
from typing import Generic, TypeVar, cast

from sqlalchemy import Select, delete, select
from sqlalchemy.ext.asyncio import AsyncSession

T = TypeVar("T")
IdType = TypeVar("IdType")


class BaseRepository(ABC, Generic[T]):
    def __init__(self, session: AsyncSession):
        self._session = session

    @property
    @abstractmethod
    def model(self) -> type[T]:
        """Return the model class for this repository"""
        raise NotImplementedError

    async def create(self, entity: T) -> T:
        """Persist a new entity"""
        self._session.add(entity)
        await self._session.flush()
        await self._session.refresh(entity)
        return entity

    async def get(self, _id: IdType) -> T | None:
        """Get entity by primary key"""
        return await self._session.get(self.model, _id)

    async def get_many(self, *ids: IdType) -> Sequence[T]:
        """Get multiple entities by primary keys"""
        stmt = select(self.model).where(self.model.id.in_(ids))
        result = await self._session.execute(stmt)
        return result.scalars().all()

    async def find(self, stmt: Select) -> Sequence[T]:
        """Execute a custom SELECT query"""
        result = await self._session.execute(stmt)
        return result.scalars().all()

    async def update(self, entity: T) -> T:
        """Update an existing entity"""
        await self._session.flush()
        await self._session.refresh(entity)
        return entity

    async def delete(self, entity: T) -> None:
        """Delete an entity"""
        await self._session.delete(entity)

    async def delete_many(self, *ids: IdType) -> int:
        """Delete multiple entities by primary keys"""
        stmt = delete(self.model).where(self.model.id.in_(ids))
        result = await self._session.execute(stmt)
        rowcount = cast(int, result.rowcount)  # Number of rows deleted
        return rowcount
