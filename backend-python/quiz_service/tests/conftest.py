import os

import pytest
from httpx import AsyncClient, ASGITransport
from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker

from shared.db.models import Base, Quiz, User
from shared.repositories import UserRepository, QuizRepository
from shared.utils.unitofwork import UnitOfWork

from quiz_app.api.dependencies.dependencies import get_uow
from quiz_app.main import app

TEST_DB_URL = os.getenv("TEST_DB_URL")

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
            questions=quiz_data["questions"],
            owner_id=test_user.id
        )
        quiz_db = await repo.create(quiz)
        return quiz_db
