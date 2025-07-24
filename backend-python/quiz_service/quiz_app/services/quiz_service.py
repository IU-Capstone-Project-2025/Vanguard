import logging
from uuid import UUID

from shared.db.models import Quiz as DBQuiz
from shared.repositories import QuizRepository, TagRepository
from shared.schemas.quiz import QuizCreate, QuizUpdate, QuizResponse, QuizFilterMode
from shared.schemas.tag import TagBase
from shared.utils.unitofwork import UnitOfWork

from quiz_app.exceptions import QuizNotFoundError, QuizForbiddenError, InvalidQuizQueryParametersError
from quiz_app.core.metrics import (
    QUIZ_CREATIONS_TOTAL, QUIZ_FETCHES_TOTAL, QUIZ_UPDATES_TOTAL,
    QUIZ_DELETES_TOTAL, QUIZ_LISTING_REQUESTS_TOTAL, SERVICE
)

logger = logging.getLogger("app")


class QuizService:
    def __init__(self, uow: UnitOfWork):
        self.uow = uow

    async def create_quiz(self, user_id: UUID, quiz_in: QuizCreate) -> QuizResponse:
        logger.debug(f"Creating quiz for user {user_id} with visibility={quiz_in.is_public}")
        async with self.uow.transaction() as session:
            quiz_repo = QuizRepository(session)
            tag_repo = TagRepository(session)

            # Create the Quiz object (without tags yet)
            quiz_data = quiz_in.model_dump(exclude={"tags"})
            quiz_data["owner_id"] = user_id
            quiz_db = await quiz_repo.create(DBQuiz(**quiz_data))

            # Normalize & validate tag names via Pydantic
            normalized_names = {TagBase.model_validate({"name": raw}).name for raw in quiz_in.tags if raw}

            # Bulk get-or-create Tag records
            tags = await tag_repo.get_or_create_bulk(list(normalized_names))

            # Associate tags onto the quiz
            quiz_db.tags = tags
            await session.flush()
            await session.refresh(quiz_db)

            visibility = "public" if quiz_db.is_public else "private"
            QUIZ_CREATIONS_TOTAL.labels(service=SERVICE, status="success", visibility=visibility).inc()
            logger.info(
                "Created quiz successfully",
                extra={
                    "metadata": {
                        "user_id": user_id,
                        "quiz_id": quiz_db.id,
                        "visibility": visibility
                    }
                }
            )
            return QuizResponse.model_validate(quiz_db)

    async def get_quiz_by_id(self, quiz_id: UUID, user_id: UUID | None) -> QuizResponse:
        logger.debug(f"Fetching quiz {quiz_id} (requested by {user_id})")
        async with self.uow.readonly() as session:
            repo = QuizRepository(session)
            quiz = await repo.get(_id=quiz_id)

            public_only = "true" if user_id is None else "false"

            if not quiz:
                QUIZ_FETCHES_TOTAL.labels(service=SERVICE, status="not_found", public_only=public_only).inc()
                raise QuizNotFoundError(quiz_id)

            if quiz.owner_id != user_id and not quiz.is_public:
                QUIZ_FETCHES_TOTAL.labels(service=SERVICE, status="forbidden", public_only=public_only).inc()
                raise QuizForbiddenError("You do not own this quiz.")

            QUIZ_FETCHES_TOTAL.labels(service=SERVICE, status="success", public_only=public_only).inc()
            logger.info(
                "Quiz retrieved successfully",
                extra={
                    "metadata": {
                        "user_id": user_id,
                        "quiz_id": quiz_id
                    }
                }
            )
            return QuizResponse.model_validate(quiz)

    async def update_quiz(self, quiz_id: UUID, user_id: UUID, data: QuizUpdate) -> QuizResponse:
        logger.debug(f"Updating quiz {quiz_id} by user {user_id}")
        async with self.uow.transaction() as session:
            quiz_repo = QuizRepository(session)
            tag_repo = TagRepository(session)

            quiz = await quiz_repo.get(_id=quiz_id)
            if not quiz:
                QUIZ_UPDATES_TOTAL.labels(service=SERVICE, status="not_found").inc()
                raise QuizNotFoundError(quiz_id)
            if quiz.owner_id != user_id:
                QUIZ_UPDATES_TOTAL.labels(service=SERVICE, status="forbidden").inc()
                raise QuizForbiddenError("You do not own this quiz.")

            # Update scalar fields
            update_data = data.model_dump(exclude_none=True, exclude={'tags'})
            for field, value in update_data.items():
                setattr(quiz, field, value)

            # If tags provided, re-normalize and re-associate
            if data.tags is not None:
                normalized = {TagBase.model_validate({"name": raw}).name for raw in data.tags if raw}
                tags = await tag_repo.get_or_create_bulk(list(normalized))
                quiz.tags = tags

            await session.flush()
            await session.refresh(quiz)

            QUIZ_UPDATES_TOTAL.labels(service=SERVICE, status="success").inc()
            logger.info(
                "Quiz updated successfully",
                extra={
                    "metadata": {
                        "user_id": user_id,
                        "quiz_id": quiz_id
                    }
                }
            )
            return QuizResponse.model_validate(quiz)

    async def delete_quiz(self, quiz_id: UUID, user_id: UUID) -> None:
        logger.debug(f"Deleting quiz {quiz_id} by user {user_id}")
        async with self.uow.transaction() as session:
            repo = QuizRepository(session)
            quiz = await repo.get(_id=quiz_id)
            if not quiz:
                QUIZ_DELETES_TOTAL.labels(service=SERVICE, status="not_found").inc()
                raise QuizNotFoundError(quiz_id)
            if quiz.owner_id != user_id:
                QUIZ_DELETES_TOTAL.labels(service=SERVICE, status="forbidden").inc()
                raise QuizForbiddenError("You do not own this quiz.")

            await repo.delete(quiz)
            QUIZ_DELETES_TOTAL.labels(service=SERVICE, status="success").inc()
            logger.info(
                "Quiz deleted successfully",
                extra={
                    "metadata": {
                        "user_id": user_id,
                        "quiz_id": quiz_id
                    }
                }
            )

    async def list_quizzes(
            self,
            *,
            requester_id: UUID | None,
            public: bool | None = None,
            mine: bool | None = None,
            user_id: UUID | None = None,
            search: str | None = None,
            tags: list[str] | None = None,
            page: int = 1,
            size: int = 20
    ) -> list[QuizResponse]:
        logger.debug(f"Listing quizzes for requester={requester_id}, public={public}, mine={mine}, user_id={user_id}, search={search}, tags={tags}")

        filter_type = "unknown"
        try:
            filter_mode = self.resolve_filter_mode(requester_id, public, mine, user_id)
            filter_type = filter_mode.value
        except InvalidQuizQueryParametersError:
            QUIZ_LISTING_REQUESTS_TOTAL.labels(service=SERVICE, status="invalid", filter_type=filter_type).inc()
            raise

        if tags:
            tags = list({TagBase.model_validate({"name": raw}).name for raw in tags if raw})

        async with self.uow.readonly() as session:
            repo = QuizRepository(session)
            quizzes = await repo.list_quizzes(
                mode=filter_mode,
                requester_id=requester_id,
                user_id=user_id,
                search=search,
                tags=tags,
                page=page,
                size=size
            )

        QUIZ_LISTING_REQUESTS_TOTAL.labels(service=SERVICE, status="success", filter_type=filter_type).inc()
        logger.info(
            "Returned quizzes for requester",
            extra={
                "metadata": {
                    "user_id": requester_id,
                    "count": len(quizzes)
                }
            }
        )
        return [QuizResponse.model_validate(q) for q in quizzes]

    @staticmethod
    def resolve_filter_mode(
            requester_id: UUID | None,
            public: bool | None = None,
            mine: bool | None = None,
            user_id: UUID | None = None
    ) -> QuizFilterMode:
        if requester_id is None:
            if public is False:
                raise InvalidQuizQueryParametersError("Cannot filter `public=false` when unauthenticated")
            if mine:
                raise InvalidQuizQueryParametersError("Cannot filter `mine=true` when unauthenticated")
            if user_id is not None:
                raise InvalidQuizQueryParametersError("Cannot filter by `user_id` when unauthenticated")
            return QuizFilterMode.ALL_PUBLIC

        # mine=True
        if mine:
            if user_id is not None and user_id != requester_id:
                raise InvalidQuizQueryParametersError("Cannot filter by other users with `mine=true`")

            if public is None:
                return QuizFilterMode.ALL_MINE
            elif public is False:
                return QuizFilterMode.MINE_PRIVATE
            else:
                return QuizFilterMode.MINE_PUBLIC

        # mine=False
        if mine is False:
            if public is False:
                raise InvalidQuizQueryParametersError("Cannot request private quizzes of other users")
            if user_id == requester_id:
                raise InvalidQuizQueryParametersError("`mine=false` is incompatible with own `user_id`")
            return QuizFilterMode.OTHER_PUBLIC

        # Default visibility: mine + public
        if public is None:
            return QuizFilterMode.VISIBLE_TO_ME
        elif public is False:
            if user_id is not None and user_id != requester_id:
                raise InvalidQuizQueryParametersError("Cannot request private quizzes of other users")
            return QuizFilterMode.MINE_PRIVATE
        else:
            return QuizFilterMode.ALL_PUBLIC
