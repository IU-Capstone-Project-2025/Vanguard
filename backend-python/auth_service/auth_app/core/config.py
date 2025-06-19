from pydantic_settings import BaseSettings

# from dotenv import load_dotenv
# load_dotenv()


class Settings(BaseSettings):
    APP_NAME: str = "Auth Service API"
    APP_VERSION: str = "1.0.0"
    DEBUG: bool = False


settings = Settings()
