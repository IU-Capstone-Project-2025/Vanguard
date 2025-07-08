import os

from pydantic_settings import BaseSettings

ENVIRONMENT = os.getenv("ENVIRONMENT", "dev").lower()

if ENVIRONMENT != "prod" and ENVIRONMENT != "test":
    from dotenv import load_dotenv
    load_dotenv()


class Settings(BaseSettings):
    APP_NAME: str = "Auth Service API"
    APP_VERSION: str = "1.0.0"
    DEBUG: bool = False
    ACCESS_TOKEN_EXPIRE_MINUTES: int = 15
    REFRESH_TOKEN_EXPIRE_DAYS: int = 7


settings = Settings()
