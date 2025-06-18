from shared.db.models import Quiz as DBQuiz
from shared.repositories.base_repo import SQLAlchemyRepository


class QuizRepository(SQLAlchemyRepository[DBQuiz]):
    model = DBQuiz