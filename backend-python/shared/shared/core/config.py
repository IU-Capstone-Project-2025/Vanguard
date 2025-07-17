from enum import Enum

from pydantic import ValidationError
from pydantic_settings import BaseSettings


class Environment(str, Enum):
    PROD = "production"
    DEV = "development"
    TEST = "test"
    STAGING = "staging"


class Settings(BaseSettings):
    DB_URL: str
    TEST_DB_URL: str
    CORS_ORIGINS: str
    JWT_SECRET_KEY: str
    JWT_ALGORITHM: str = "HS256"
    ENVIRONMENT: Environment = Environment.PROD

    @property
    def cors_origins_list(self) -> list[str]:
        return [origin.strip() for origin in self.CORS_ORIGINS.split(",") if origin.strip()]


try:
    settings = Settings()
except ValidationError:
    from dotenv import load_dotenv
    load_dotenv()
    settings = Settings()
