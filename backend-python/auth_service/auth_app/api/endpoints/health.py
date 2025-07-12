from datetime import datetime, UTC

from fastapi import APIRouter, status

from auth_app.core.config import settings

router = APIRouter(tags=["auth-service"])

@router.get(
    "/health",
    summary="Health Check",
    response_description="Service status, version, and timestamp",
    status_code=status.HTTP_200_OK
)
async def health_check():
    return {
        "status": "OK",
        "version": settings.APP_VERSION,
        "timestamp": datetime.now(UTC).isoformat()
    }
