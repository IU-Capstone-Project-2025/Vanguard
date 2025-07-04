import os

from pydantic_settings import BaseSettings

ENVIRONMENT = os.getenv("ENVIRONMENT", "dev").lower()

if ENVIRONMENT != "prod" and ENVIRONMENT != "test":
    from dotenv import load_dotenv
    load_dotenv()


class Settings(BaseSettings):
    S3_REGION: str
    S3_ENDPOINT_URL: str
    S3_BUCKET: str
    AWS_ACCESS_KEY_ID: str
    AWS_SECRET_ACCESS_KEY: str
    MAX_IMAGE_SIZE: int

    APP_NAME: str = "Quiz Service API"
    APP_VERSION: str = "1.0.0"
    DEBUG: bool = False


settings = Settings()
