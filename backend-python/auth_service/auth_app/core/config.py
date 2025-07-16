from enum import Enum

from pydantic import ValidationError
from pydantic_settings import BaseSettings


class Environment(str, Enum):
    PROD = "production"
    DEV = "development"
    TEST = "test"
    STAGING = "staging"


class Settings(BaseSettings):
    APP_NAME: str = "Auth Service API"
    APP_VERSION: str = "1.0.0"

    ENVIRONMENT: Environment = Environment.PROD
    DEBUG: bool = False

    ACCESS_TOKEN_EXPIRE_MINUTES: int = 15
    REFRESH_TOKEN_EXPIRE_DAYS: int = 7


try:
    settings = Settings()
except ValidationError:
    from dotenv import load_dotenv
    load_dotenv()
    settings = Settings()
