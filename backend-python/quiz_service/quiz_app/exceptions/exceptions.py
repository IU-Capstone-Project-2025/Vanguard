from uuid import UUID


class NotFoundError(Exception):
    """Base exception for when a resource is not found"""
    def __init__(self, resource: str, identifier: str | int | UUID):
        self.resource = resource
        self.identifier = identifier
        super().__init__(f"{resource} with id '{identifier}' not found")


class ForbiddenError(Exception):
    """Base exception for when access is forbidden"""
    def __init__(self, message: str = "Action is forbidden"):
        super().__init__(message)


class BadRequestError(Exception):
    """Base exception for when a request is bad request"""
    def __init__(self, message: str = "Bad Request"):
        super().__init__(message)


class ImageS3Error(Exception):
    """Base exception for image upload related errors"""
    pass


class QuizNotFoundError(NotFoundError):
    """Raised when a quiz is not found"""
    def __init__(self, identifier: str | int | UUID):
        super().__init__("Quiz", identifier)


class InvalidQuizQueryParametersError(BadRequestError):
    """Raised when a quiz parameters are invalid"""
    def __init__(self, message: str = "Invalid quiz query parameters"):
        super().__init__(message)


class InvalidImageError(ImageS3Error):
    """Raised when invalid image file is provided"""
    pass


class FileTooLargeError(ImageS3Error):
    """Raised when file exceeds size limit"""
    pass


class ImageNotFoundError(ImageS3Error):
    """Raised when image to delete is not found"""
    pass
