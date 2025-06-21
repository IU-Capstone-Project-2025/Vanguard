from sqlalchemy.exc import IntegrityError

from shared.db.models import Quiz, User
from shared.repositories import UserRepository, QuizRepository

from quiz_app.api.dependencies.dependencies import get_uow


async def init_sample_data():
    uow = await get_uow()
    try:
        async with uow.transaction() as session:
            user_repo = UserRepository(session)
            quiz_repo = QuizRepository(session)

            user_data = {
                "username": "test_user",
                "email": "test@mail.com",
                "password_hash": "password_hash"
            }
            test_user = await user_repo.create(User(**user_data))

            quiz_data = {
                "title": "Basic Python Knowledge",
                "description": "A quiz to test your basic Python knowledge.",
                "is_public": True,
                "owner_id": test_user.id,
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
            await quiz_repo.create(Quiz(**quiz_data))
    except IntegrityError:
        pass  # sample data already exists in db, so no need to add it
