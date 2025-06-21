from pydantic import BaseModel


class ImageUploadResponse(BaseModel):
    url: str
