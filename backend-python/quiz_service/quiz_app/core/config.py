from enum import Enum

from pydantic import ValidationError
from pydantic_settings import BaseSettings


class Environment(str, Enum):
    PROD = "production"
    DEV = "development"
    TEST = "test"
    STAGING = "staging"


class Settings(BaseSettings):
    S3_REGION: str
    S3_ENDPOINT_URL: str
    S3_BUCKET: str
    AWS_ACCESS_KEY_ID: str
    AWS_SECRET_ACCESS_KEY: str
    MAX_IMAGE_SIZE: int

    APP_NAME: str = "Quiz Service API"
    APP_VERSION: str = "1.0.0"

    ENVIRONMENT: Environment = Environment.PROD
    DEBUG: bool = False


try:
    settings = Settings()
except ValidationError:
    from dotenv import load_dotenv
    load_dotenv()
    settings = Settings()
