from fastapi import APIRouter

router = APIRouter(prefix="/api", tags=["auth-service"])

@router.get("/auth")
async def auth():
    return {"status": "Working auth service"}
