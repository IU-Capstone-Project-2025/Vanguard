import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from prometheus_client import make_asgi_app
from sqlalchemy import text

from shared.core.config import settings as shared_settings
from shared.db.database import async_session_maker

from auth_app.api.endpoints import auth_router, health_router
from auth_app.core.config import settings, Environment
from auth_app.core.loggers import setup_loggers
from auth_app.core.middleware import MetricsMiddleware, LoggingMiddleware
from auth_app.exceptions.handlers import register_exception_handlers

setup_loggers(debug=settings.DEBUG)
logger = logging.getLogger("app")

@asynccontextmanager
async def lifespan(_: FastAPI):
    startup_msg = f"Starting {settings.APP_NAME} in {settings.ENVIRONMENT.value} mode"
    logger.info(startup_msg)

    try:
        async with async_session_maker() as db:
            await db.execute(text("SELECT 1"))
        logger.info("Database connection verified")
    except Exception as e:
        logger.exception("Database connection failed")
        raise e

    yield

    shutdown_msg = f"Shutting down {settings.APP_NAME}"
    logger.info(shutdown_msg)

app = FastAPI(
    title=settings.APP_NAME,
    version=settings.APP_VERSION,
    debug=settings.DEBUG,
    lifespan=lifespan,
    docs_url="/docs" if settings.ENVIRONMENT == Environment.DEV else None,
    redoc_url="/redoc" if settings.ENVIRONMENT == Environment.DEV else None,
    openapi_url="/openapi.json" if settings.ENVIRONMENT == Environment.DEV else None
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=shared_settings.cors_origins_list,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
    max_age=600 if settings.ENVIRONMENT == Environment.PROD else 0
)
app.add_middleware(MetricsMiddleware)
app.add_middleware(LoggingMiddleware)

metrics_app = make_asgi_app()
app.mount("/metrics", metrics_app)

app.include_router(health_router)
app.include_router(auth_router)

register_exception_handlers(app)
