import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from sqlalchemy import text

from shared.core.config import settings as shared_settings
from shared.db.database import async_session_maker

from auth_app.api.endpoints.auth import router as auth_router
from auth_app.core.config import settings
from auth_app.core.logging import setup_logging
from auth_app.exceptions.handlers import register_exception_handlers

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
    lifespan=lifespan,
    root_path="/api/auth"
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=shared_settings.cors_origins_list,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"]
)

app.include_router(auth_router)

register_exception_handlers(app)
