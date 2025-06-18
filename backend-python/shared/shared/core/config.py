from pydantic_settings import BaseSettings

# from dotenv import load_dotenv
# load_dotenv()


class Settings(BaseSettings):
    DEBUG: bool = True
    DB_URL: str


settings = Settings()
