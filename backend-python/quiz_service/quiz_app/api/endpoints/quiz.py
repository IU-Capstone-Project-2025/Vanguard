from fastapi import APIRouter

router = APIRouter(prefix="/api", tags=["quiz-service"])

@router.get("/quizzes")
async def get_quizzes():
    return {"status": "Working quizzes service"}
