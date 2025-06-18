from fastapi import FastAPI
from app.api.endpoints.auth import router as auth_router
from shared.core.config import settings

print(settings.DB_URL)

app = FastAPI()

app.include_router(auth_router)
