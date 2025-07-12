import logging
from uuid import uuid4

import boto3
from botocore.exceptions import ClientError, BotoCoreError
from fastapi import UploadFile

from shared.schemas.image import ImageUploadResponse

from quiz_app.core.config import settings
from quiz_app.core.metrics import QUIZ_IMAGE_UPLOADS_TOTAL, QUIZ_IMAGE_UPLOAD_SIZE_BYTES, SERVICE
from quiz_app.exceptions import InvalidImageError, ImageNotFoundError, FileTooLargeError, ImageS3Error

logger = logging.getLogger("app")


class S3ImageService:
    def __init__(self):
        self.region = settings.S3_REGION
        self.endpoint_url = settings.S3_ENDPOINT_URL
        self.bucket = settings.S3_BUCKET
        self.aws_access_key_id = settings.AWS_ACCESS_KEY_ID
        self.aws_secret_access_key = settings.AWS_SECRET_ACCESS_KEY
        self.max_size = settings.MAX_IMAGE_SIZE

        self.client = boto3.client(
            "s3",
            region_name=self.region,
            aws_access_key_id=self.aws_access_key_id,
            aws_secret_access_key=self.aws_secret_access_key,
            endpoint_url=self.endpoint_url
        )
        self.bucket_url = f"{self.client.meta.endpoint_url}/{self.bucket}"

    def upload_file(self, file: UploadFile, folder: str = "uploads") -> ImageUploadResponse:
        """
        Uploads a file to S3 storage.

        Args:
            file: FastAPI UploadFile object
            folder: Target folder in bucket

        Returns:
            Public URL of the uploaded file in ImageUploadResponse object.

        Raises:
            InvalidImageError: 400 - For invalid file types
            FileTooLargeError: 413 - For files > 5MB
            ImageS3Error: 500 - For other upload failures
        """
        logger.debug(f"Received image upload request: filename={file.filename}, content_type={file.content_type}")

        if not file.content_type.startswith("image/"):
            QUIZ_IMAGE_UPLOADS_TOTAL.labels(service=SERVICE, status="invalid_type").inc()
            raise InvalidImageError("Only image files are allowed")

        if file.size > self.max_size:
            QUIZ_IMAGE_UPLOADS_TOTAL.labels(service=SERVICE, status="too_large").inc()
            raise FileTooLargeError(f"File size exceeds {self.max_size} bytes limit")

        file_ext = file.filename.split(".")[-1]
        key = f"{folder}/{uuid4()}.{file_ext}"

        try:
            self.client.upload_fileobj(
                file.file,
                self.bucket,
                key,
                ExtraArgs={"ACL": "public-read", "ContentType": file.content_type}
            )
        except (ClientError, BotoCoreError) as e:
            QUIZ_IMAGE_UPLOADS_TOTAL.labels(service=SERVICE, status="s3_error").inc()
            logger.error(f"S3 upload error: {str(e)}")
            raise ImageS3Error(f"S3 upload failed: {str(e)}")
        except Exception as e:
            QUIZ_IMAGE_UPLOADS_TOTAL.labels(service=SERVICE, status="unexpected_error").inc()
            logger.error(f"Unexpected upload error: {str(e)}")
            raise ImageS3Error(f"Upload failed: {str(e)}")

        QUIZ_IMAGE_UPLOADS_TOTAL.labels(service=SERVICE, status="success").inc()
        QUIZ_IMAGE_UPLOAD_SIZE_BYTES.labels(service=SERVICE, status="success").observe(file.size)
        logger.info(f"Uploaded image {file.filename} as {key} ({file.size} bytes)")

        url = self.get_file_url(key)
        return ImageUploadResponse(url=url)

    def delete_file(self, img_url: str) -> None:
        """
        Deletes a file from S3 storage.

        Args:
            img_url: Public URL of the file to delete

        Raises:
            ImageNotFoundError: 404 - When file doesn't exist
            ImageS3Error: 500 - For other deletion failures
        """
        logger.debug(f"Attempting to delete image: {img_url}")
        try:
            key = self.get_key_from_url(img_url)
            self.client.delete_object(Bucket=self.bucket, Key=key)
            logger.info(f"Deleted image {key} from bucket {self.bucket}")
        except self.client.exceptions.NoSuchKey:
            raise ImageNotFoundError("The specified image does not exist")
        except (ClientError, BotoCoreError) as e:
            logger.error(f"S3 deletion error: {str(e)}")
            raise ImageS3Error(f"S3 deletion failed: {str(e)}")
        except Exception as e:
            logger.error(f"Unexpected deletion error: {str(e)}")
            raise ImageS3Error(f"Deletion failed: {str(e)}")

    def get_file_url(self, key: str) -> str:
        """Generates public URL for a given S3 key"""
        return f"{self.bucket_url}/{key}"

    def get_key_from_url(self, url: str) -> str:
        """Extracts S3 key from public URL"""
        if not url.startswith(self.bucket_url):
            logger.error(f"Invalid image URL: {url}")
            raise ValueError("Invalid URL - does not belong to this bucket")
        return url.replace(f"{self.bucket_url}/", "")
