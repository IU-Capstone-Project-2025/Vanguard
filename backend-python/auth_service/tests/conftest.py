import bcrypt
import pytest
from httpx import AsyncClient, ASGITransport
from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker

from shared.core.config import settings
from shared.db.models import Base, User
from shared.repositories import UserRepository
from shared.utils.unitofwork import UnitOfWork

from auth_app.core.dependencies import get_uow
from auth_app.main import app

TEST_DB_URL = settings.TEST_DB_URL

@pytest.fixture(scope="function")
async def uow_test():
    engine = create_async_engine(TEST_DB_URL, echo=False)
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.drop_all)
        await conn.run_sync(Base.metadata.create_all)
    session_maker = async_sessionmaker(bind=engine, expire_on_commit=False, autoflush=False)

    yield UnitOfWork(session_maker)

    await engine.dispose()

@pytest.fixture
async def test_client(uow_test):
    app.dependency_overrides = {get_uow: lambda: uow_test}
    async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test/api") as client:
        yield client

@pytest.fixture
async def registered_user(uow_test):
    async with uow_test.transaction() as sess:
        repo = UserRepository(sess)
        user = await repo.create(User(
            email="alice@example.com",
            username="alice",
            password_hash=bcrypt.hashpw(b"secret", bcrypt.gensalt()).decode()
        ))
        return user
