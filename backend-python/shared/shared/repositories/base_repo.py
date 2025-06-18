from abc import ABC, abstractmethod
from typing import Any, Dict, List, Generic, TypeVar, Optional, cast

from sqlalchemy import insert, select, delete
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import DeclarativeBase

ModelType = TypeVar("ModelType", bound=DeclarativeBase)


class AbstractRepository(ABC, Generic[ModelType]):
    """Abstract base class defining the repository interface."""

    @abstractmethod
    async def find_all(self) -> List[ModelType]:
        """Retrieve all records."""
        ...

    @abstractmethod
    async def find_one(self) -> Optional[ModelType]:
        """Retrieve a single record."""
        ...

    @abstractmethod
    async def add(self, data: Dict[str, Any]) -> ModelType:
        """Add a new record."""
        ...

    @abstractmethod
    async def update(self, data: Dict[str, Any]) -> Optional[ModelType]:
        """Update an existing record."""
        ...

    @abstractmethod
    async def delete(self) -> bool:
        """Delete a record."""
        ...


class SQLAlchemyRepository(AbstractRepository[ModelType], Generic[ModelType]):
    """Generic SQLAlchemy repository implementation."""
    model: ModelType

    def __init__(self, session: AsyncSession):
        self.session = session
        if self.model is None:
            raise NotImplementedError("Repository must have a 'model' class attribute defined.")

    async def find_all(
            self,
            skip: int = 0,
            limit: int | None = None,
            order_by: Dict[str, str] | None = None,
            **filter_by: Any
    ) -> List[ModelType]:
        stmt = select(self.model).filter_by(**filter_by)

        if order_by:
            for column, direction in order_by.items():
                sort_column = getattr(self.model, column)
                if direction.lower() == "desc":
                    stmt = stmt.order_by(sort_column.desc())
                else:
                    stmt = stmt.order_by(sort_column.asc())

        if limit is not None and limit > 0:
            stmt = stmt.limit(limit)
            if skip > 0:
                stmt = stmt.offset(skip)

        result = await self.session.execute(stmt)
        return list(result.scalars().all())

    async def find_one(self, **filter_by: Any) -> Optional[ModelType]:
        stmt = select(self.model).filter_by(**filter_by)
        result = await self.session.execute(stmt)
        return result.scalars().first()

    async def add(self, data: Dict[str, Any]) -> ModelType:
        stmt = insert(self.model).values(**data).returning(self.model)
        result = await self.session.execute(stmt)
        return result.scalars().first()

    async def update(self, data: Dict[str, Any], **filter_by: Any) -> Optional[ModelType]:
        obj = await self.find_one(**filter_by)
        if not obj:
            return None

        for field, value in data.items():
            if hasattr(obj, field):
                setattr(obj, field, value)
        return obj

    async def delete(self, **filter_by: Any) -> bool:
        stmt = delete(self.model).filter_by(**filter_by)
        result = await self.session.execute(stmt)
        rowcount = cast(int, result.rowcount)  # IDE thinks `rowcount` is method for some reason
        return rowcount > 0
