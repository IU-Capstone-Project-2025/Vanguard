import os

import asyncio
import pytest
from httpx import AsyncClient, ASGITransport
from sqlalchemy import text
from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker

from shared.db.models import Base, Quiz, User
from shared.repositories import UserRepository, QuizRepository
from shared.utils.unitofwork import UnitOfWork

from quiz_app.api.dependencies.dependencies import get_uow
from quiz_app.main import app

ADMIN_DB_URL = os.environ.get("ADMIN_DB_URL")
TEST_DB_URL = os.getenv("TEST_DB_URL")

# Session-wide: create and drop test DB
@pytest.fixture(scope="session", autouse=True)
async def create_test_database():
    admin_engine = create_async_engine(ADMIN_DB_URL, isolation_level="AUTOCOMMIT")
    async with admin_engine.connect() as conn:
        await conn.execute(text("DROP DATABASE IF EXISTS test_kahoot_clone"))
        await conn.execute(text("CREATE DATABASE test_kahoot_clone"))
    await admin_engine.dispose()

    yield

    # Final teardown: drop test DB
    admin_engine = create_async_engine(ADMIN_DB_URL, isolation_level="AUTOCOMMIT")
    async with admin_engine.connect() as conn:
        await conn.execute(text("DROP DATABASE IF EXISTS test_kahoot_clone"))
    await admin_engine.dispose()

# Function-scoped: clean tables between each test
@pytest.fixture(scope="function", autouse=True)
async def clean_test_database():
    engine = create_async_engine(TEST_DB_URL)
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.drop_all)
        await conn.run_sync(Base.metadata.create_all)
    await engine.dispose()

# Create event loop for entire test session
@pytest.fixture(scope="session")
def event_loop():
    loop = asyncio.get_event_loop_policy().new_event_loop()
    yield loop
    loop.close()

# Session-scoped engine for all tests
@pytest.fixture(scope="session")
async def test_engine():
    engine = create_async_engine(TEST_DB_URL, echo=False)
    yield engine
    await engine.dispose()

@pytest.fixture
async def session_maker(test_engine):
    session_maker = async_sessionmaker(bind=test_engine, expire_on_commit=False)
    yield session_maker
    await test_engine.dispose()

@pytest.fixture
async def uow_test(session_maker):
    return UnitOfWork(session_maker)

@pytest.fixture
async def test_client(uow_test):
    app.dependency_overrides = {get_uow: lambda: uow_test}
    async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as client:
        yield client

@pytest.fixture
async def test_user(uow_test):
    user_data = {
        "username": "test_user",
        "email": "test@mail.com",
        "password_hash": "password_hash"
    }

    async with uow_test.transaction() as session:
        repo = UserRepository(session)
        user_db = await repo.create(User(**user_data))
        return user_db

@pytest.fixture
async def test_quiz(uow_test, test_user):
    quiz_data = {
        "title": "Basic Python Knowledge",
        "description": "A quiz to test your basic Python knowledge.",
        "is_public": True,
        "tags": ["python", "beginner"],
        "questions": [
            {
                "type": "single_choice",
                "text": "What is the output of print(2 ** 3)?",
                "options": [
                    {"text": "6", "is_correct": False},
                    {"text": "8", "is_correct": True},
                    {"text": "9", "is_correct": False},
                    {"text": "5", "is_correct": False}
                ]
            },
            {
                "type": "single_choice",
                "text": "Which keyword is used to create a function in Python?",
                "options": [
                    {"text": "func", "is_correct": False},
                    {"text": "function", "is_correct": False},
                    {"text": "def", "is_correct": True},
                    {"text": "define", "is_correct": False}
                ]
            },
            {
                "type": "single_choice",
                "text": "What data type is the result of: 3 / 2 in Python 3?",
                "options": [
                    {"text": "int", "is_correct": False},
                    {"text": "float", "is_correct": True},
                    {"text": "str", "is_correct": False},
                    {"text": "decimal", "is_correct": False}
                ]
            }
        ]
    }

    async with uow_test.transaction() as session:
        repo = QuizRepository(session)
        quiz = Quiz(
            title=quiz_data["title"],
            description=quiz_data["description"],
            is_public=quiz_data["is_public"],
            tags=[],
            questions=quiz_data["questions"],
            owner_id=test_user.id
        )
        quiz_db = await repo.create(quiz)
        return quiz_db
