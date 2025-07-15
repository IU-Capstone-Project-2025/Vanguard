from uuid import UUID

from fastapi import APIRouter, Depends, Query, status

from shared.core.dependencies import get_current_user_id
from shared.schemas.quiz import QuizCreate, QuizResponse, QuizUpdate

from quiz_app.core.dependencies import get_quiz_service, get_potential_user_id
from quiz_app.services.quiz_service import QuizService

router = APIRouter(tags=["quiz-service"])


@router.post(
    "/",
    status_code=status.HTTP_201_CREATED,
    response_model=QuizResponse,
    summary="Create a new quiz",
    description="Creates a new quiz with the provided details."
)
async def create_quiz(
        quiz: QuizCreate,
        user_id: UUID = Depends(get_current_user_id),
        quiz_service: QuizService = Depends(get_quiz_service)
):
    print(user_id)
    return await quiz_service.create_quiz(user_id=user_id, quiz_in=quiz)


@router.get(
    "/{quiz_id}",
    response_model=QuizResponse,
    summary="Get quiz by ID",
    description="Retrieves a quiz by its unique ID."
)
async def get_quiz_by_id(
        quiz_id: UUID,
        user_id: UUID | None = Depends(get_potential_user_id),
        quiz_service: QuizService = Depends(get_quiz_service)
):
    return await quiz_service.get_quiz_by_id(quiz_id=quiz_id, user_id=user_id)


@router.put(
    "/{quiz_id}",
    response_model=QuizResponse,
    summary="Update quiz",
    description="Updates a quiz by its unique ID."
)
async def update_quiz(
        quiz_id: UUID,
        quiz: QuizUpdate,
        user_id: UUID = Depends(get_current_user_id),
        quiz_service: QuizService = Depends(get_quiz_service)
):
    return await quiz_service.update_quiz(quiz_id=quiz_id, user_id=user_id, data=quiz)


@router.delete(
    "/{quiz_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    summary="Delete quiz",
    description="Deletes a quiz by its unique ID."
)
async def delete_quiz(
        quiz_id: UUID,
        user_id: UUID = Depends(get_current_user_id),
        quiz_service: QuizService = Depends(get_quiz_service)
):
    await quiz_service.delete_quiz(quiz_id=quiz_id, user_id=user_id)


@router.get(
    "/",
    response_model=list[QuizResponse],
    summary="Get list of quizzes",
    description="Get public quizzes and ones owned by the user"
)
async def list_quizzes(
        public: bool | None = Query(None),
        mine: bool | None = Query(None),
        user_id: UUID | None = Query(None),
        search: str | None = Query(None, min_length=1),
        tag: list[str] = Query([], alias="tag"),
        page: int = Query(1, ge=1),
        size: int = Query(20, ge=1, le=100),
        user_id_req: UUID | None = Depends(get_potential_user_id),
        quiz_service: QuizService = Depends(get_quiz_service)
):
    """
    List/filter quizzes.
    - Unauthenticated: only public.
    - Authenticated default: public + own.
    - public=true → only public.
    - mine=true → only own.
    - user_id=… → public by that user.
    - search → title/description ilike.
    - tag → AND filter (must have all).
    """
    return await quiz_service.list_quizzes(
        requester_id=user_id_req,
        public=public,
        mine=mine,
        user_id=user_id,
        search=search,
        tags=tag or None,
        page=page,
        size=size
    )
