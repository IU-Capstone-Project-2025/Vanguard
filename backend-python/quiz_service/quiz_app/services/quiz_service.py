import logging
from uuid import UUID

from shared.db.models import Quiz as DBQuiz
from shared.repositories import QuizRepository
from shared.schemas.quiz import QuizCreate, QuizResponse
from shared.utils.unitofwork import UnitOfWork

from quiz_app.exceptions import QuizNotFoundError, ForbiddenError

logger = logging.getLogger("app")


class QuizService:
    def __init__(self, uow: UnitOfWork):
        self.uow = uow

    async def create_quiz(self, user_id: UUID, quiz_in: QuizCreate) -> QuizResponse:
        async with self.uow.transaction() as session:
            repo = QuizRepository(session)
            quiz_data = quiz_in.model_dump()
            quiz_data["owner_id"] = user_id
            quiz_db = await repo.create(DBQuiz(**quiz_data))

            logger.info(f"Created quiz {quiz_db.id} by user {user_id}")
            return QuizResponse.model_validate(quiz_db)

    async def get_public_quizzes(self) -> list[QuizResponse]:
        async with self.uow.readonly() as session:
            repo = QuizRepository(session)
            public_quizzes = await repo.get_public_quizzes()
            return [QuizResponse.model_validate(q) for q in public_quizzes]

    async def get_user_quizzes(self, user_id: UUID) -> list[QuizResponse]:
        async with self.uow.readonly() as session:
            repo = QuizRepository(session)
            user_quizzes = await repo.get_user_quizzes(user_id)
            return [QuizResponse.model_validate(q) for q in user_quizzes]

    async def get_visible_quizzes(self, user_id: UUID) -> list[QuizResponse]:
        async with self.uow.readonly() as session:
            repo = QuizRepository(session)
            visible_quizzes = await repo.get_visible_quizzes(user_id)
            return [QuizResponse.model_validate(q) for q in visible_quizzes]

    async def get_quiz_by_id(self, quiz_id: UUID, user_id: UUID | None) -> QuizResponse:
        async with self.uow.readonly() as session:
            repo = QuizRepository(session)
            quiz = await repo.get(_id=quiz_id)
            if not quiz:
                raise QuizNotFoundError(quiz_id)
            if quiz.owner_id != user_id and not quiz.is_public:
                raise ForbiddenError("You do not own this quiz.")
            return QuizResponse.model_validate(quiz)
