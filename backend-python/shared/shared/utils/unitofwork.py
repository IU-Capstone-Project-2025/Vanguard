from contextlib import asynccontextmanager
from typing import AsyncGenerator

from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker


class UnitOfWork:
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]):
        self.session_factory = session_factory

    @asynccontextmanager
    async def transaction(self) -> AsyncGenerator[AsyncSession, None]:
        """Provide a transactional scope around a series of operations"""
        session: AsyncSession = self.session_factory()
        try:
            yield session
            await session.commit()
        except Exception:
            await session.rollback()
            raise
        finally:
            await session.close()

    @asynccontextmanager
    async def readonly(self) -> AsyncGenerator[AsyncSession, None]:
        """Provide a read-only session (no transaction)"""
        session: AsyncSession = self.session_factory()
        try:
            yield session
        finally:
            await session.close()
