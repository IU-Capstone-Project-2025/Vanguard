import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from sqlalchemy import text

from shared.core.config import settings as shared_settings
from shared.db.database import async_session_maker

from quiz_app.init_sample_data import init_sample_data
from quiz_app.api.endpoints import health_router, quiz_router, image_router
from quiz_app.core.config import settings
from quiz_app.core.logging import setup_logging
from quiz_app.exceptions.handlers import register_exception_handlers

setup_logging(debug=settings.DEBUG)
logger = logging.getLogger("app")

@asynccontextmanager
async def lifespan(_: FastAPI):
    logger.info("Application startup.")

    try:
        db = async_session_maker()
        await db.execute(text("SELECT 1"))
        await db.close()
        logger.info("Database connected successfully.")
    except Exception as e:
        logger.exception("Failed to connect to the database.")
        raise e

    await init_sample_data()
    logger.info("Sample user and quiz was successfully initialized.")

    yield

    logger.info("Application shutdown.")

app = FastAPI(
    title=settings.APP_NAME,
    version=settings.APP_VERSION,
    debug=settings.DEBUG,
    lifespan=lifespan,
    root_path="/api/quiz"
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=shared_settings.cors_origins_list,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"]
)

app.include_router(health_router)
app.include_router(quiz_router)
app.include_router(image_router)

register_exception_handlers(app)
