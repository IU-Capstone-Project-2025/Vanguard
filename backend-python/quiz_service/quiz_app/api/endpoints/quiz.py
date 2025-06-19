from datetime import datetime, UTC
from uuid import UUID, uuid4

from fastapi import APIRouter, Depends, status

from shared.schemas.quiz import QuizCreate, QuizResponse

from quiz_app.api.dependencies.dependencies import get_quiz_service
from quiz_app.core.config import settings
from quiz_app.services.quiz_service import QuizService

router = APIRouter(tags=["quiz-service"])


# TODO: move to separate router
@router.get(
    "/health",
    status_code=status.HTTP_200_OK,
    summary="Quiz Service Health Check",
    response_description="Service status and dependencies"
)
async def health_check():
    return {
        "status": "OK",
        "version": settings.APP_VERSION,
        "timestamp": datetime.now(UTC).isoformat()
    }


@router.post(
    "/",
    status_code=status.HTTP_201_CREATED,
    response_model=QuizResponse,
    summary="Create a new quiz",
    description="Creates a new quiz with the provided details."
)
async def create_quiz(
        quiz: QuizCreate,
        quiz_service: QuizService = Depends(get_quiz_service)
):
    # TODO: handle auth normally
    fake_auth_user = uuid4()
    quiz_response = await quiz_service.create_quiz(fake_auth_user, quiz)
    return quiz_response


@router.get(
    "/",
    response_model=list[QuizResponse],
    summary="Get list of quizzes",
    description="Get public quizzes and ones owned by the user"
)
async def get_quizzes(
        # TODO: handle auth normally
        user_id: UUID | None = None,
        only_my: bool = False,
        quiz_service: QuizService = Depends(get_quiz_service)
):
    if user_id:
        if only_my:
            return await quiz_service.get_user_quizzes(user_id=user_id)
        else:
            return await quiz_service.get_visible_quizzes(user_id=user_id)
    else:
        return await quiz_service.get_public_quizzes()


@router.get(
    "/{quiz_id}",
    response_model=QuizResponse,
    summary="Get quiz by ID",
    description="Retrieves a quiz by its unique ID."
)
async def get_quiz_by_id(
        quiz_id: UUID,
        # TODO: handle auth normally
        user_id: UUID | None = None,
        quiz_service: QuizService = Depends(get_quiz_service)
):
    quiz = await quiz_service.get_quiz_by_id(quiz_id=quiz_id, user_id=user_id)
    return quiz
