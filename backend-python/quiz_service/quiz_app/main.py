import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI
from sqlalchemy import text

from shared.db.database import async_session_maker

from quiz_app.init_sample_data import init_sample_data
from quiz_app.api.endpoints.quiz import router as quiz_router
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
    lifespan=lifespan
)

app.include_router(quiz_router)

register_exception_handlers(app)
