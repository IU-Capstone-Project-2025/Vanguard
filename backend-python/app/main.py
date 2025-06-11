from fastapi import FastAPI
from app.api.endpoints.quiz import router as quiz_router

app = FastAPI()

app.include_router(quiz_router)
