from uuid import UUID

from fastapi import APIRouter, Depends, Header, Query, status

from shared.core.dependencies import get_current_user_id
from shared.schemas.quiz import QuizCreate, QuizResponse, QuizUpdate

from quiz_app.core.dependencies import get_quiz_service, get_potential_user_id
from quiz_app.services.quiz_service import QuizService

router = APIRouter(
    prefix="/api",
    tags=["quiz-service"],
    responses={
        401: {"description": "Unauthorized - Invalid or missing credentials"},
        500: {"description": "Internal Server Error"}
    }
)


@router.post(
    "/",
    summary="Create a new quiz",
    description="Creates a new quiz with the provided details.",
    response_model=QuizResponse,
    status_code=status.HTTP_201_CREATED,
    response_description="The created quiz object",
    responses={
        422: {"description": "Validation error in request body"}
    }
)
async def create_quiz(
        quiz: QuizCreate,
        _authorization: str | None = Header(None, alias="Authorization", description="Bearer <access_token>"),
        user_id: UUID = Depends(get_current_user_id),
        quiz_service: QuizService = Depends(get_quiz_service)
):
    return await quiz_service.create_quiz(user_id=user_id, quiz_in=quiz)


@router.get(
    "/{quiz_id}",
    summary="Get quiz by ID",
    description="Retrieves a quiz by its unique ID. Public quizzes are accessible to everyone; private ones only to the owner.",
    response_model=QuizResponse,
    status_code=status.HTTP_200_OK,
    response_description="The quiz data",
    responses={
        403: {"description": "Forbidden - Not allowed to access this quiz"},
        404: {"description": "Quiz not found"}
    }
)
async def get_quiz_by_id(
        quiz_id: UUID,
        _authorization: str | None = Header(None, alias="Authorization", description="Bearer <access_token>"),
        user_id: UUID | None = Depends(get_potential_user_id),
        quiz_service: QuizService = Depends(get_quiz_service)
):
    return await quiz_service.get_quiz_by_id(quiz_id=quiz_id, user_id=user_id)


@router.put(
    "/{quiz_id}",
    summary="Update quiz",
    description="Updates the quiz with the given ID. Only the owner is allowed to modify their quiz.",
    response_model=QuizResponse,
    status_code=status.HTTP_200_OK,
    response_description="The updated quiz object",
    responses={
        403: {"description": "Forbidden - Not allowed to access this quiz"},
        404: {"description": "Quiz not found"},
        422: {"description": "Validation error in request body"}
    }
)
async def update_quiz(
        quiz_id: UUID,
        quiz: QuizUpdate,
        _authorization: str | None = Header(None, alias="Authorization", description="Bearer <access_token>"),
        user_id: UUID = Depends(get_current_user_id),
        quiz_service: QuizService = Depends(get_quiz_service)
):
    return await quiz_service.update_quiz(quiz_id=quiz_id, user_id=user_id, data=quiz)


@router.delete(
    "/{quiz_id}",
    summary="Delete quiz",
    description="Deletes a quiz by ID. Only the quiz owner can delete it.",
    status_code=status.HTTP_204_NO_CONTENT,
    responses={
        204: {"description": "Quiz successfully deleted"},
        403: {"description": "Forbidden - Not allowed to access this quiz"},
        404: {"description": "Quiz not found"}
    }
)
async def delete_quiz(
        quiz_id: UUID,
        _authorization: str | None = Header(None, alias="Authorization", description="Bearer <access_token>"),
        user_id: UUID = Depends(get_current_user_id),
        quiz_service: QuizService = Depends(get_quiz_service)
):
    await quiz_service.delete_quiz(quiz_id=quiz_id, user_id=user_id)


@router.get(
    "/",
    summary="List & Filter Quizzes",
    description="""
Returns a paginated list of quizzes.

- Unauthenticated users see only public quizzes.
- Authenticated users see public and their own by default.
- Filters:
  - `public=true`: only public quizzes
  - `mine=true`: only your quizzes (requires auth)
  - `user_id=<uuid>`: quizzes from a specific user
  - `search=<term>`: fuzzy match on title/description
  - `tag=tag1&tag=tag2`: AND filter on tags
""",
    response_model=list[QuizResponse],
    status_code=status.HTTP_200_OK,
    response_description="List of quizzes matching the filter",
    responses={
        400: {"description": "Invalid query parameters"},
    }
)
async def list_quizzes(
        public: bool | None = Query(None, description="Return only public quizzes"),
        mine: bool | None = Query(None, description="Return only your quizzes (requires authentication)"),
        user_id: UUID | None = Query(None, description="Return public quizzes by this user ID"),
        search: str | None = Query(None, description="Search in quiz title or description", min_length=1),
        tag: list[str] = Query([], alias="tag", description="Tag filters (must match all given tags)"),
        page: int = Query(1, ge=1, description="Page number (pagination)"),
        size: int = Query(20, ge=1, le=100, description="Items per page (pagination)"),
        _authorization: str | None = Header(None, alias="Authorization", description="Bearer <access_token>"),
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
