import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI
from sqlalchemy import text

from shared.db.database import async_session_maker

from auth_app.api.endpoints.auth import router as auth_router
from auth_app.core.config import settings
from auth_app.core.logging import setup_logging

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

    yield

    logger.info("Application shutdown.")

app = FastAPI(
    title=settings.APP_NAME,
    version=settings.APP_VERSION,
    debug=settings.DEBUG,
    lifespan=lifespan
)

app.include_router(auth_router)
