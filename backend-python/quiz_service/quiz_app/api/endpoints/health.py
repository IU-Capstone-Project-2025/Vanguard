from datetime import datetime, UTC

from fastapi import APIRouter, status

from quiz_app.core.config import settings

router = APIRouter(tags=["quiz-service"])


@router.get(
    "/health",
    status_code=status.HTTP_200_OK,
    summary="Quiz Service Health Check",
    response_description="Service status and dependencies",
    include_in_schema=False
)
async def health_check():
    return {
        "status": "OK",
        "version": settings.APP_VERSION,
        "timestamp": datetime.now(UTC).isoformat()
    }
