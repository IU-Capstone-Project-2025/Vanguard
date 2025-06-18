from fastapi import FastAPI
from app.api.endpoints.quiz import router as quiz_router
from shared.core.config import settings

print(settings.DB_URL)

app = FastAPI()

app.include_router(quiz_router)
