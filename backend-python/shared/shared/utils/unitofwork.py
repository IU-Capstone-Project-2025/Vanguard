import logging
from abc import ABC, abstractmethod

from shared.db.database import async_session_maker
from shared.repositories import UserRepository, QuizRepository

logger = logging.getLogger("app")


class IUnitOfWork(ABC):
    @abstractmethod
    async def __aenter__(self):
        ...

    @abstractmethod
    async def __aexit__(self, *args):
        ...

    @abstractmethod
    async def commit(self):
        ...

    @abstractmethod
    async def rollback(self):
        ...


class UnitOfWork(IUnitOfWork):
    def __init__(self):
        self.session_factory = async_session_maker

    async def __aenter__(self):
        self.session = self.session_factory()
        self.user = UserRepository(self.session)
        self.quiz = QuizRepository(self.session)
        logger.debug("UoW session started")
        return self

    async def __aexit__(self, *args):
        await self.rollback()
        await self.session.close()
        self.session = None
        logger.debug("UoW session closed")

    async def commit(self):
        await self.session.commit()
        logger.debug("UoW committed")

    async def rollback(self):
        await self.session.rollback()
        logger.debug("UoW rolled back")
