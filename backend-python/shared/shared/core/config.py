from pydantic_settings import BaseSettings
from dotenv import load_dotenv

load_dotenv()


class Settings(BaseSettings):
    DB_URL: str
    TEST_DB_URL: str
    CORS_ORIGINS: str
    JWT_SECRET_KEY: str
    JWT_ALGORITHM: str = "HS256"

    @property
    def cors_origins_list(self) -> list[str]:
        return [origin.strip() for origin in self.CORS_ORIGINS.split(",") if origin.strip()]


settings = Settings()
