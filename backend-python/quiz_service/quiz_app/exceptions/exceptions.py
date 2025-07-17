from fastapi import status
from uuid import UUID


class QuizServiceError(Exception):
    """Base exception for quiz service."""
    status_code = status.HTTP_500_INTERNAL_SERVER_ERROR
    detail = "An error occurred with the quiz service"


class QuizNotFoundError(QuizServiceError):
    """Raised when a quiz is not found"""
    status_code = status.HTTP_404_NOT_FOUND

    def __init__(self, identifier: str | int | UUID):
        self.detail = f"Quiz with id '{identifier}' not found"


class QuizForbiddenError(QuizServiceError):
    """Base exception for when access is forbidden"""
    status_code = status.HTTP_403_FORBIDDEN

    def __init__(self, message: str = "Action with quiz is forbidden"):
        self.detail = message


class InvalidQuizQueryParametersError(QuizServiceError):
    """Raised when a quiz parameters are invalid"""
    status_code = status.HTTP_400_BAD_REQUEST

    def __init__(self, message: str = "Invalid quiz query parameters"):
        self.detail = message


class ImageServiceError(Exception):
    """Base exception for image service."""
    status_code = status.HTTP_500_INTERNAL_SERVER_ERROR

    def __init__(self, message: str = "An error occurred with the image service"):
        self.detail = message


class InvalidImageError(ImageServiceError):
    """Raised when invalid image file is provided"""
    status_code = status.HTTP_400_BAD_REQUEST

    def __init__(self, message: str = "Invalid image file"):
        self.detail = message


class FileTooLargeError(ImageServiceError):
    """Raised when file exceeds size limit"""
    status_code = status.HTTP_413_REQUEST_ENTITY_TOO_LARGE

    def __init__(self, message: str = "File size exceeds maximum allowed limit"):
        self.detail = message


class ImageNotFoundError(ImageServiceError):
    """Raised when image to delete is not found"""
    status_code = status.HTTP_404_NOT_FOUND

    def __init__(self, message: str = "Image not found"):
        self.detail = message


class InvalidImageURL(ImageServiceError):
    """Raised when an invalid S3 image URL is provided"""
    status_code = status.HTTP_400_BAD_REQUEST

    def __init__(self, message: str = "Invalid image URL"):
        self.detail = message
