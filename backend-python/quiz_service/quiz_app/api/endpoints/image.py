from typing import Annotated

from fastapi import Depends, File, UploadFile, APIRouter, Query, status

from shared.schemas.image import ImageUploadResponse

from quiz_app.services.image_service import S3ImageService
from quiz_app.core.dependencies import get_image_service

router = APIRouter(prefix="/api/images", tags=["images"])


@router.post(
    "/upload",
    summary="Upload an image to S3",
    description="Uploads an image file to the S3-compatible storage bucket. Supported formats: JPEG, PNG, GIF.",
    response_model=ImageUploadResponse,
    status_code=status.HTTP_201_CREATED,
    response_description="Public URL of the uploaded image",
    responses={
        201: {"description": "Image successfully uploaded", "content": {
            "application/json": {"example": {"url": "https://storage.yandexcloud.net/tryit/uploads/uuid.png"}}
        }},
        400: {"description": "Invalid file type or corrupted file"},
        413: {"description": "File size exceeds maximum allowed"},
        500: {"description": "Internal server error during upload"}
    }
)
async def upload_image(
        file: Annotated[UploadFile, File(description="Image file to upload (JPEG, PNG, GIF supported)")],
        image_service: S3ImageService = Depends(get_image_service)
):
    return image_service.upload_file(file)


@router.delete(
    "/",
    summary="Delete Image from S3",
    description="Deletes an uploaded image from S3 storage using its full public URL.",
    status_code=status.HTTP_204_NO_CONTENT,
    response_description="Image deleted successfully",
    responses={
        204: {"description": "Image successfully deleted"},
        400: {"description": "Invalid image URL"},
        404: {"description": "Image not found"},
        500: {"description": "Internal server error during deletion"}
    }
)
async def delete_image(
        img_url: Annotated[str, Query(description="Full URL of the image to delete")],
        image_service: S3ImageService = Depends(get_image_service)
):
    image_service.delete_file(img_url)
